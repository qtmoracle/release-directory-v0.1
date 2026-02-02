package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"os"
	"path/filepath"


)
// ---------- Demo asset lookup (fixture-backed, no DB) ----------
var assetFixtureMap = map[string]string{
	// Original academic asset
	"asset:bf90c309b76abbad8d0329a8fc1861404b9e134095a107f97d16884494ff5a25":
		"/home/qtm/qtm-workspaces/academic-demo/directory_mvd.json",

	// Changed academic asset (sensitivity test)
	"asset:897c5267976af3abcefed38cca4177f65a5ee540f77dd5cad80e849ff6074992":
		"/home/qtm/qtm-workspaces/academic-demo/directory_mvd_changed.json",
}
// ---------- Minted asset persistence (append-only JSONL + in-memory index) ----------

const mintedLogPath = "./data/minted_assets.jsonl"

type mintedRecord struct {
	AssetID           string         `json:"asset_id"`
	FingerprintSHA256 string         `json:"fingerprint_sha256"`
	MVD               map[string]any `json:"mvd"`
	CreatedAtUTC      string         `json:"created_at_utc"`
}

var mintedIndex = map[string]mintedRecord{}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o755)
}

func loadMintedIndex(path string) error {
	mintedIndex = map[string]mintedRecord{}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no log yet, that's fine
		}
		return err
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var rec mintedRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			return fmt.Errorf("parse minted log line: %w", err)
		}
		if rec.AssetID != "" {
			mintedIndex[rec.AssetID] = rec
		}
	}
	return nil
}

func appendMintedRecord(path string, rec mintedRecord) error {
	if err := ensureParentDir(path); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	lineBytes, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(lineBytes, '\n')); err != nil {
		return err
	}
	return nil
}
// ---------- Listing persistence (append-only JSONL + in-memory index) ----------

const listingLogPath = "./data/listings.jsonl"

type listingRecord struct {
	AssetID       string `json:"asset_id"`
	Visibility    string `json:"visibility"`      // "private" | "public"
	Paid          bool   `json:"paid"`            // client assertion for demo
	ListingStatus string `json:"listing_status"`  // "private" | "payment_required" | "public"
	FeeRequired   bool   `json:"fee_required"`
	CreatedAtUTC  string `json:"created_at_utc"`
}

var listingIndex = map[string]listingRecord{}

func loadListingIndex(path string) error {
	listingIndex = map[string]listingRecord{}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no log yet
		}
		return err
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var rec listingRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			return fmt.Errorf("parse listing log line: %w", err)
		}
		if rec.AssetID != "" {
			listingIndex[rec.AssetID] = rec
		}
	}
	return nil
}

func appendListingRecord(path string, rec listingRecord) error {
	if err := ensureParentDir(path); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	lineBytes, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(lineBytes, '\n')); err != nil {
		return err
	}
	return nil
}
// ---------- Canonical JSON (minimal, deterministic) ----------

// canonicalize recursively sorts map keys and preserves array order.
// It returns a structure that will serialize deterministically with encoding/json.
func canonicalize(v any) any {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		out := make(map[string]any, len(t))
		for _, k := range keys {
			out[k] = canonicalize(t[k])
		}
		return out

	case []any:
		out := make([]any, len(t))
		for i := range t {
			out[i] = canonicalize(t[i])
		}
		return out

	default:
		return v
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// ---------- Idempotency spec: MVD v0.1 fingerprint + asset_id ----------

// Select only the fields that define identity for idempotency.
// IMPORTANT: Do not include server-generated fields (like created_at).
func buildFingerprintInput(m map[string]any) (map[string]any, error) {
	out := map[string]any{}

	// mvd_version is required for stable semantics
	if v, ok := m["mvd_version"]; ok {
		out["mvd_version"] = v
	} else {
		return nil, fmt.Errorf("missing required field: mvd_version")
	}

	// labels.display_name (human label; part of deterministic fingerprint by your prior mint request)
	// If labels exists, include it as-is (canonicalized later).
	if v, ok := m["labels"]; ok {
		out["labels"] = v
	}

	// provenance.nonce (your prior mint request includes this; it makes each asset creation unique)
	if v, ok := m["provenance"]; ok {
		out["provenance"] = v
	} else {
		return nil, fmt.Errorf("missing required field: provenance")
	}

	// p5 block (asset type/class/utility/etc. as declared)
	if v, ok := m["p5"]; ok {
		out["p5"] = v
	}

	// declared block (status/custodian)
	if v, ok := m["declared"]; ok {
		out["declared"] = v
	}

	// validity window
	if v, ok := m["validity"]; ok {
		out["validity"] = v
	}

	// external_ids (integrity pointers: hashes, signatures, URLs, etc.)
	if v, ok := m["external_ids"]; ok {
		out["external_ids"] = v
	}

	// creator/curator DID (if present in your schema)
	if v, ok := m["actors"]; ok {
		out["actors"] = v
	}

	return out, nil
}

func computeFingerprintAndAssetID(m map[string]any) (fingerprintSHA string, assetID string, err error) {
	fin, err := buildFingerprintInput(m)
	if err != nil {
		return "", "", err
	}

	// Canonicalize then marshal
	canon := canonicalize(fin)
	canonBytes, err := json.Marshal(canon)
	if err != nil {
		return "", "", fmt.Errorf("marshal canonical fingerprint input: %w", err)
	}

	// Fingerprint = sha256(canonical_bytes)
	fingerprintSHA = sha256Hex(canonBytes)

	// asset_id = sha256("qtm.directory.mvd.0.1|" + fingerprintSHA)
	seed := "qtm.directory.mvd.0.1|" + fingerprintSHA
	assetHash := sha256Hex([]byte(seed))
	assetID = "asset:" + assetHash

	return fingerprintSHA, assetID, nil
}

func loadJSONFile(path string) (map[string]any, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ---------- HTTP handlers ----------

type mintResponse struct {
	AssetID            string `json:"asset_id"`
	FingerprintSHA256  string `json:"fingerprint_sha256"`
	CreatedAtUTC       string `json:"created_at_utc"`
	MVDVersion         string `json:"mvd_version"`
	DeterminismSummary string `json:"determinism"`
}

type listingRequest struct {
	AssetID    string `json:"asset_id"`
	Visibility string `json:"visibility"` // "private" | "public"
	Paid       bool   `json:"paid"`
}

type listingResponse struct {
	AssetID       string `json:"asset_id"`
	ListingStatus string `json:"listing_status"` // "private" | "payment_required" | "public"
	FeeRequired   bool   `json:"fee_required"`
	CreatedAtUTC  string `json:"created_at_utc"`
}

func normalizeVisibility(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

type assetViewResponse struct {
	View        string         `json:"view"`
	AssetID     string         `json:"asset_id"`
	MVD         map[string]any `json:"directory_mvd"`
	Listing     listingResponse `json:"listing"`
	RenderedAt  string         `json:"rendered_at_utc"`
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assetID := strings.TrimSpace(r.URL.Query().Get("asset_id"))
	if assetID == "" {
		http.Error(w, "missing required query param: asset_id", http.StatusBadRequest)
		return
	}

		// 1) If minted in this instance (or previously persisted), return minted MVD
	if rec, ok := mintedIndex[assetID]; ok {
		resp := assetViewResponse{
			View:       "asset_card",
			AssetID:    assetID,
			MVD:        rec.MVD,
			Listing: func() listingResponse {
		if lr, ok := listingIndex[assetID]; ok {
			return listingResponse{
				AssetID:       lr.AssetID,
				ListingStatus: lr.ListingStatus,
				FeeRequired:   lr.FeeRequired,
				CreatedAtUTC:  lr.CreatedAtUTC,
			}
		}
			return listingResponse{
				AssetID:       assetID,
				ListingStatus: "private",
				FeeRequired:   false,
				CreatedAtUTC:  time.Now().UTC().Format(time.RFC3339),
			}
		}(),


		RenderedAt: time.Now().UTC().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
		return
	}

	// 2) Otherwise fall back to fixtures
	path, ok := assetFixtureMap[assetID]
	if !ok {
		http.Error(w, "asset_id not found (minted store + fixture demo)", http.StatusNotFound)
		return
	}

	mvd, err := loadJSONFile(path)
	if err != nil {
		http.Error(w, "failed to load fixture json", http.StatusInternalServerError)
		return
	}

	// Listing: for this demo endpoint, we DO persist listing state.
	// We return a deterministic default: private + unpaid.
	// Listing: return persisted listing if present; otherwise default private.
	listing := func() listingResponse {
		if lr, ok := listingIndex[assetID]; ok {
			return listingResponse{
				AssetID:       lr.AssetID,
				ListingStatus: lr.ListingStatus,
				FeeRequired:   lr.FeeRequired,
				CreatedAtUTC:  lr.CreatedAtUTC,
			}
		}
		return listingResponse{
			AssetID:       assetID,
			ListingStatus: "private",
			FeeRequired:   false,
			CreatedAtUTC:  time.Now().UTC().Format(time.RFC3339),
		}
	}()

	resp := assetViewResponse{
		View:       "asset_card",
		AssetID:    assetID,
		MVD:        mvd,
		Listing:    listing,
		RenderedAt: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func mintBaseAssetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body failed", http.StatusBadRequest)
		return
	}

	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	fp, assetID, err := computeFingerprintAndAssetID(req)
	if err != nil {
		http.Error(w, "mint rejected: "+err.Error(), http.StatusBadRequest)
		return
	}
		// Persist minted asset (append-only) + update in-memory index.
	rec := mintedRecord{
		AssetID:           assetID,
		FingerprintSHA256: "sha256:" + fp,
		MVD:               req,
		CreatedAtUTC:      time.Now().UTC().Format(time.RFC3339),
	}
	if err := appendMintedRecord(mintedLogPath, rec); err != nil {
		http.Error(w, "failed to persist minted asset", http.StatusInternalServerError)
		return
	}
	mintedIndex[assetID] = rec


	mvdVersion, _ := req["mvd_version"].(string)

	resp := mintResponse{
		AssetID:           assetID,
		FingerprintSHA256: "sha256:" + fp,
		CreatedAtUTC:      time.Now().UTC().Format(time.RFC3339),
		MVDVersion:        mvdVersion,
		DeterminismSummary: strings.Join([]string{
			"fingerprint_sha256 = sha256(canonical_json(subset))",
			"asset_id = 'asset:' + sha256('qtm.directory.mvd.0.1|' + fingerprint_sha256)",
			"canonical_json = map keys sorted recursively; arrays preserved",
		}, "; "),
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func setListingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body failed", http.StatusBadRequest)
		return
	}

	var req listingRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	req.AssetID = strings.TrimSpace(req.AssetID)
	req.Visibility = normalizeVisibility(req.Visibility)

	if req.AssetID == "" {
		http.Error(w, "missing required field: asset_id", http.StatusBadRequest)
		return
	}

	if req.Visibility != "private" && req.Visibility != "public" {
		http.Error(w, "invalid visibility (must be 'private' or 'public')", http.StatusBadRequest)
		return
	}

	feeRequired := (req.Visibility == "public")

	status := ""
	switch req.Visibility {
	case "private":
		status = "private"
	case "public":
		if req.Paid {
			status = "public"
		} else {
			status = "payment_required"
		}
	}

	resp := listingResponse{
		AssetID:       req.AssetID,
		ListingStatus: status,
		FeeRequired:   feeRequired,
		CreatedAtUTC:  time.Now().UTC().Format(time.RFC3339),
	}
		// Persist listing state (append-only) + update in-memory index.
	rec := listingRecord{
		AssetID:       resp.AssetID,
		Visibility:    req.Visibility,
		Paid:          req.Paid,
		ListingStatus: resp.ListingStatus,
		FeeRequired:   resp.FeeRequired,
		CreatedAtUTC:  resp.CreatedAtUTC,
	}
	if err := appendListingRecord(listingLogPath, rec); err != nil {
		http.Error(w, "failed to persist listing", http.StatusInternalServerError)
		return
	}
	listingIndex[resp.AssetID] = rec


	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func main() {
	mux := http.NewServeMux()

	if err := loadMintedIndex(mintedLogPath); err != nil {
		log.Fatalf("failed to load minted index: %v", err)
	}
	if err := loadListingIndex(listingLogPath); err != nil {
		log.Fatalf("failed to load listing index: %v", err)
}
	// --- UI (static) ---
	// Serves files from ./ui at /ui/
	fs := http.FileServer(http.Dir("./ui"))
	mux.Handle("/ui/", http.StripPrefix("/ui/", fs))

	// Friendly root redirect to UI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/mint_base_asset", mintBaseAssetHandler)
	mux.HandleFunc("/set_listing", setListingHandler)
	mux.HandleFunc("/asset", assetHandler)

	addr := ":8080"
	log.Printf("Directory MVP listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
