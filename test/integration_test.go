package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"merkle-airdrop/internal/api"
	"merkle-airdrop/pkg/data"
	"merkle-airdrop/pkg/merkle"
)

func TestFullWorkflow(t *testing.T) {
	// Step 1: Generate test data
	claims := data.GenerateTestData(100)

	// Step 2: Validate data
	if err := data.ValidateClaimsData(claims); err != nil {
		t.Fatalf("Data validation failed: %v", err)
	}

	// Step 3: Build Merkle tree
	tree, err := merkle.NewMerkleTree(claims)
	if err != nil {
		t.Fatalf("Failed to build tree: %v", err)
	}

	// Step 4: Generate proofs
	proofs, err := tree.GenerateAllProofs()
	if err != nil {
		t.Fatalf("Failed to generate proofs: %v", err)
	}

	// Step 5: Verify proofs count
	if len(proofs) != len(claims) {
		t.Errorf("Expected %d proofs, got %d", len(claims), len(proofs))
	}

	// Step 6: Test API endpoints
	testAPIEndpoints(t, tree, proofs)
}

func testAPIEndpoints(t *testing.T, tree *merkle.MerkleTree, proofs map[string]*merkle.MerkleProof) {
	server := api.NewAPIServer(tree, proofs)
	handler := server.SetupRoutes()

	// Test root endpoint
	t.Run("GetRootHash", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/root", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response["success"] != true {
			t.Error("Expected success to be true")
		}

		if response["merkleRoot"] == "" {
			t.Error("Expected merkleRoot to be non-empty")
		}
	})

	// Test proof endpoint
	t.Run("GetProof", func(t *testing.T) {
		// Get first claim address
		testAddr := tree.Claims[0].Address.Hex()

		req := httptest.NewRequest(http.MethodGet, "/api/proof/"+testAddr, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response["success"] != true {
			t.Error("Expected success to be true")
		}

		if response["address"] != testAddr {
			t.Errorf("Expected address %s, got %s", testAddr, response["address"])
		}
	})

	// Test stats endpoint
	t.Run("GetStats", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response["success"] != true {
			t.Error("Expected success to be true")
		}

		totalClaims := int(response["totalClaims"].(float64))
		if totalClaims != len(tree.Claims) {
			t.Errorf("Expected %d total claims, got %d", len(tree.Claims), totalClaims)
		}
	})

	// Test verify endpoint
	t.Run("VerifyProof", func(t *testing.T) {
		testAddr := tree.Claims[0].Address.Hex()
		testProof := proofs[testAddr]

		payload := map[string]interface{}{
			"address": testAddr,
			"amount":  testProof.Amount,
			"proof":   testProof.Proof,
		}

		payloadBytes, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/verify", strings.NewReader(string(payloadBytes)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response["success"] != true {
			t.Error("Expected success to be true")
		}
	})
}

func TestDataValidation(t *testing.T) {
	t.Run("ValidData", func(t *testing.T) {
		claims := data.GenerateTestData(10)
		if err := data.ValidateClaimsData(claims); err != nil {
			t.Errorf("Valid data failed validation: %v", err)
		}
	})

	t.Run("EmptyData", func(t *testing.T) {
		var claims []merkle.AirdropClaim
		if err := data.ValidateClaimsData(claims); err == nil {
			t.Error("Empty data should fail validation")
		}
	})

	t.Run("DuplicateAddresses", func(t *testing.T) {
		claims := data.GenerateTestData(2)
		claims[1].Address = claims[0].Address // Create duplicate

		if err := data.ValidateClaimsData(claims); err == nil {
			t.Error("Duplicate addresses should fail validation")
		}
	})
}

func TestTreeProperties(t *testing.T) {
	claims := data.GenerateTestData(100)
	tree, err := merkle.NewMerkleTree(claims)
	if err != nil {
		t.Fatalf("Failed to build tree: %v", err)
	}

	// Test root hash consistency
	rootHash1 := tree.GetRootHash()

	// Rebuild tree with same data

	tree2, err := merkle.NewMerkleTree(claims)
	if err != nil {
		t.Fatalf("Failed to rebuild tree: %v", err)
	}
	rootHash2 := tree2.GetRootHash()
	if rootHash1 != rootHash2 {
		t.Errorf("Expected root hash %s, got %s", rootHash1, rootHash2)
	}
}
