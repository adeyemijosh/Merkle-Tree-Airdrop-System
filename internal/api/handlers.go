// internal/api/handlers.go
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"merkle-airdrop/pkg/merkle"

	"github.com/ethereum/go-ethereum/common"
)

type APIServer struct {
	tree   *merkle.MerkleTree
	proofs map[string]*merkle.MerkleProof
}

func NewAPIServer(tree *merkle.MerkleTree, proofs map[string]*merkle.MerkleProof) *APIServer {
	return &APIServer{
		tree:   tree,
		proofs: proofs,
	}
}

// GetRootHash returns the Merkle root hash
func (s *APIServer) GetRootHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"merkleRoot": s.tree.GetRootHash(),
		"success":    true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProof returns the Merkle proof for a specific address
func (s *APIServer) GetProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := strings.TrimPrefix(r.URL.Path, "/api/proof/")
	if !common.IsHexAddress(address) {
		http.Error(w, "Invalid address format", http.StatusBadRequest)
		return
	}

	// Normalize address
	normalizedAddr := common.HexToAddress(address).Hex()

	proof, exists := s.proofs[normalizedAddr]
	if !exists {
		response := map[string]interface{}{
			"error":   "Address not found in airdrop",
			"success": false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"address":    normalizedAddr,
		"proof":      proof.Proof,
		"amount":     proof.Amount,
		"index":      proof.Index,
		"merkleRoot": s.tree.GetRootHash(),
		"success":    true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetStats returns airdrop statistics
func (s *APIServer) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"totalClaims": len(s.tree.Claims),
		"totalProofs": len(s.proofs),
		"merkleRoot":  s.tree.GetRootHash(),
		"proofDepth":  calculateTreeDepth(len(s.tree.Claims)),
		"success":     true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VerifyProof verifies a Merkle proof
func (s *APIServer) VerifyProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string   `json:"address"`
		Amount  string   `json:"amount"`
		Proof   []string `json:"proof"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !common.IsHexAddress(req.Address) {
		http.Error(w, "Invalid address format", http.StatusBadRequest)
		return
	}

	// TODO: Implement actual proof verification logic
	isValid := len(req.Proof) > 0 // Simplified verification

	response := map[string]interface{}{
		"valid":      isValid,
		"address":    req.Address,
		"amount":     req.Amount,
		"merkleRoot": s.tree.GetRootHash(),
		"success":    true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetupRoutes configures HTTP routes
func (s *APIServer) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/root", s.GetRootHash)
	mux.HandleFunc("/api/proof/", s.GetProof)
	mux.HandleFunc("/api/stats", s.GetStats)
	mux.HandleFunc("/api/verify", s.VerifyProof)

	// CORS middleware
	return addCORS(mux)
}

// addCORS adds CORS headers
func addCORS(handler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		handler.ServeHTTP(w, r)
	})

	return mux
}

func calculateTreeDepth(leaves int) int {
	if leaves <= 1 {
		return 0
	}
	depth := 0
	for leaves > 1 {
		leaves = (leaves + 1) / 2
		depth++
	}
	return depth
}
