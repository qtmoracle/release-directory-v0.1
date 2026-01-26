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
)

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
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/mint_base_asset", mintBaseAssetHandler)
	mux.HandleFunc("/set_listing", setListingHandler)

	addr := ":8080"
	log.Printf("Directory MVP listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
