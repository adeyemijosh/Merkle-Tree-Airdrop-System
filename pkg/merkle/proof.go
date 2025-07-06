// pkg/merkle/proof.go
package merkle

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// GenerateProof creates a Merkle proof for a specific address
func (mt *MerkleTree) GenerateProof(address common.Address) (*MerkleProof, error) {
	// Find the leaf for this address
	var targetLeaf *MerkleNode
	var targetIndex uint32

	for i, leaf := range mt.Leaves {
		if leaf.Data.Address == address {
			targetLeaf = leaf
			targetIndex = uint32(i)
			break
		}
	}

	if targetLeaf == nil {
		return nil, fmt.Errorf("address not found in tree")
	}

	// Generate proof path
	proof := mt.generateProofPath(targetLeaf, targetIndex)

	return &MerkleProof{
		Proof:  proof,
		Index:  targetIndex,
		Amount: targetLeaf.Data.Amount.String(),
	}, nil
}

// generateProofPath generates the proof path for a leaf
func (mt *MerkleTree) generateProofPath(_ *MerkleNode, index uint32) []string {
	var proof []string

	// Start from leaves and work up
	nodes := mt.Leaves
	currentIndex := index

	for len(nodes) > 1 {
		var nextLevel []*MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			var right *MerkleNode

			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				right = left // Duplicate for odd number
			}

			// If current index is at this level, add sibling to proof
			if uint32(i) == currentIndex {
				if currentIndex%2 == 0 {
					// We're left child, add right sibling
					proof = append(proof, fmt.Sprintf("0x%x", right.Hash))
				} else {
					// We're right child, add left sibling
					proof = append(proof, fmt.Sprintf("0x%x", left.Hash))
				}
			} else if uint32(i+1) == currentIndex {
				// We're right child, add left sibling
				proof = append(proof, fmt.Sprintf("0x%x", left.Hash))
			}

			// Create parent for next level
			parentHash := HashInternal(left.Hash, right.Hash)
			parent := &MerkleNode{Hash: parentHash}
			nextLevel = append(nextLevel, parent)
		}

		nodes = nextLevel
		currentIndex = currentIndex / 2
	}

	return proof
}

// GenerateAllProofs generates proofs for all addresses using goroutines
func (mt *MerkleTree) GenerateAllProofs() (map[string]*MerkleProof, error) {
	numWorkers := runtime.NumCPU()
	if numWorkers > len(mt.Claims) {
		numWorkers = len(mt.Claims)
	}

	// Result type for channel
	type proofResult struct {
		Address string
		Proof   *MerkleProof
		Error   error
	}

	// Create channels
	jobs := make(chan AirdropClaim, len(mt.Claims))
	results := make(chan proofResult, len(mt.Claims))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for claim := range jobs {
				proof, err := mt.GenerateProof(claim.Address)
				results <- proofResult{
					Address: claim.Address.Hex(),
					Proof:   proof,
					Error:   err,
				}
			}
		}()
	}

	// Send jobs
	go func() {
		for _, claim := range mt.Claims {
			jobs <- claim
		}
		close(jobs)
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	proofs := make(map[string]*MerkleProof)
	for result := range results {
		if result.Error != nil {
			return nil, fmt.Errorf("failed to generate proof for %s: %w", result.Address, result.Error)
		}
		proofs[result.Address] = result.Proof
	}

	return proofs, nil
}
