// main.go
package main

import (
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"merkle-airdrop/pkg/data"
	"merkle-airdrop/pkg/merkle"
)

func main() {
	fmt.Println(" Merkle Tree Airdrop System")
	fmt.Println("============================")

	// Configuration
	const (
		dataFile   = "airdrop_data.csv"
		outputFile = "merkle_proofs.json"
		numClaims  = 10000 // For testing
	)

	// Step 1: Load or generate airdrop data
	fmt.Printf(" Loading airdrop data...\n")
	var claims []merkle.AirdropClaim
	var err error

	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		fmt.Printf(" Generating %d test claims...\n", numClaims)
		claims = data.GenerateTestData(numClaims)

		// Save test data to CSV
		if err := saveToCSV(claims, dataFile); err != nil {
			log.Fatal("Failed to save test data:", err)
		}
	} else {
		claims, err = data.LoadAirdropFromCSV(dataFile)
		if err != nil {
			log.Fatal("Failed to load data:", err)
		}
	}

	fmt.Printf(" Loaded %d claims\n", len(claims))

	// Step 2: Building Merkle tree
	fmt.Printf(" Building Merkle tree...\n")
	start := time.Now()

	tree, err := merkle.NewMerkleTree(claims)
	if err != nil {
		log.Fatal("Failed to build tree:", err)
	}

	buildTime := time.Since(start)
	fmt.Printf(" Tree built in %v\n", buildTime)
	fmt.Printf(" Root hash: %s\n", tree.GetRootHash())

	// Step 3: Generate all proofs in parallel
	fmt.Printf(" Generating proofs with goroutines...\n")
	start = time.Now()

	proofs, err := tree.GenerateAllProofs()
	if err != nil {
		log.Fatal("Failed to generate proofs:", err)
	}

	proofTime := time.Since(start)
	fmt.Printf(" Generated %d proofs in %v\n", len(proofs), proofTime)
	fmt.Printf(" Performance: %.2f proofs/second\n", float64(len(proofs))/proofTime.Seconds())

	// Step 4: Save results
	fmt.Printf(" Saving results...\n")

	result := map[string]interface{}{
		"merkleRoot":  tree.GetRootHash(),
		"proofs":      proofs,
		"totalClaims": len(claims),
		"generatedAt": time.Now().Unix(),
		"buildTime":   buildTime.String(),
		"proofTime":   proofTime.String(),
	}

	if err := saveToJSON(result, outputFile); err != nil {
		log.Fatal("Failed to save results:", err)
	}

	fmt.Printf(" Results saved to %s\n", outputFile)

	// Step 5: Verify multiple random proofs
	fmt.Printf(" Verifying proofs...\n")

	if len(claims) > 0 {
		// Test first claim
		testClaim := claims[0]
		proof, exists := proofs[testClaim.Address.Hex()]
		if !exists {
			log.Fatal("Proof not found for test address")
		}

		valid := verifyProof(proof, testClaim, tree.GetRootHash())
		if valid {
			fmt.Printf(" Proof verification successful for %s!\n", testClaim.Address.Hex())
		} else {
			fmt.Printf(" Proof verification failed for %s!\n", testClaim.Address.Hex())
		}

		// Test a few more random claims
		testIndices := []int{len(claims) / 4, len(claims) / 2, len(claims) * 3 / 4}
		for _, idx := range testIndices {
			if idx < len(claims) {
				claim := claims[idx]
				proof, exists := proofs[claim.Address.Hex()]
				if exists {
					valid := verifyProof(proof, claim, tree.GetRootHash())
					if valid {
						fmt.Printf(" Proof verification successful for claim %d\n", idx)
					} else {
						fmt.Printf(" Proof verification failed for claim %d\n", idx)
					}
				}
			}
		}
	}

	// Step 6: Display summary statistics
	fmt.Printf("\n Summary Statistics:\n")
	fmt.Printf("   - Total Claims: %d\n", len(claims))
	fmt.Printf("   - Merkle Root: %s\n", tree.GetRootHash())
	fmt.Printf("   - Tree Height: %d\n", calculateTreeHeight(len(claims)))
	fmt.Printf("   - Average Proof Length: %.1f hashes\n", calculateAverageProofLength(proofs))

	fmt.Printf("\n Merkle Tree Airdrop System completed successfully!\n")
	fmt.Printf(" Performance Summary:\n")
	fmt.Printf("   - Tree construction: %v\n", buildTime)
	fmt.Printf("   - Proof generation: %v\n", proofTime)
	fmt.Printf("   - Total time: %v\n", buildTime+proofTime)
	fmt.Printf("   - Proofs per second: %.2f\n", float64(len(proofs))/proofTime.Seconds())

	fmt.Printf("\n Next Steps:\n")
	fmt.Printf("   1. Deploy smart contract with root hash: %s\n", tree.GetRootHash())
	fmt.Printf("   2. Fund contract with tokens\n")
	fmt.Printf("   3. Set up web interface with contract address\n")
	fmt.Printf("   4. Test claim functionality\n")
}

// verifyProof verifies a Merkle proof against a claim and root hash
func verifyProof(proof *merkle.MerkleProof, claim merkle.AirdropClaim, rootHash string) bool {
	// Reconstruct the leaf hash
	leafHash := merkle.HashLeaf(claim.Address, claim.Amount, claim.Index)

	// Verify proof path
	currentHash := leafHash
	for _, proofHash := range proof.Proof {
		// Removes "0x" prefix and decode hex
		proofBytes, err := hex.DecodeString(proofHash[2:])
		if err != nil {
			fmt.Printf("Error decoding proof hash %s: %v\n", proofHash, err)
			return false
		}

		// Hash with sibling
		currentHash = merkle.HashInternal(currentHash, proofBytes)
	}

	// Compare with root hash
	reconstructedRoot := fmt.Sprintf("0x%x", currentHash)
	return reconstructedRoot == rootHash
}

// saveToCSV saves claims to CSV file
func saveToCSV(claims []merkle.AirdropClaim, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"address", "amount", "index"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write claims
	for _, claim := range claims {
		record := []string{
			claim.Address.Hex(),
			claim.Amount.String(),
			strconv.FormatUint(uint64(claim.Index), 10),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// saveToJSON saves data to JSON file
func saveToJSON(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// calculateTreeHeight calculates the height of a binary tree given number of leaves
func calculateTreeHeight(numLeaves int) int {
	if numLeaves <= 1 {
		return 0
	}

	height := 0
	nodes := numLeaves

	for nodes > 1 {
		nodes = (nodes + 1) / 2 // Round up division
		height++
	}

	return height
}

// calculateAverageProofLength calculates the average length of proofs
func calculateAverageProofLength(proofs map[string]*merkle.MerkleProof) float64 {
	if len(proofs) == 0 {
		return 0
	}

	totalLength := 0
	for _, proof := range proofs {
		totalLength += len(proof.Proof)
	}

	return float64(totalLength) / float64(len(proofs))
}
