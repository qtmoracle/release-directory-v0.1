package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type MintResponse struct {
	AssetID    string `json:"asset_id"`
	CreatedAt  string `json:"created_at"`
}

func mintBaseAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp := MintResponse{
		AssetID:   "asset_demo_001",
		CreatedAt: "2026-01-25T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/mint_base_asset", mintBaseAsset)
	log.Println("Directory MVP listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
