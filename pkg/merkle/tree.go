// pkg/merkle/tree.go
package merkle

import (
	"fmt"
	"sort"
)

// NewMerkleTree creates a new Merkle tree from airdrop claims
func NewMerkleTree(claims []AirdropClaim) (*MerkleTree, error) {
	if len(claims) == 0 {
		return nil, fmt.Errorf("no claims provided")
	}

	// Sort claims by address for deterministic tree
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Address.Hex() < claims[j].Address.Hex()
	})

	// Update indices after sorting
	for i := range claims {
		claims[i].Index = uint32(i)
	}

	tree := &MerkleTree{
		Claims: claims,
	}

	// Create leaf nodes
	leaves := make([]*MerkleNode, len(claims))
	for i, claim := range claims {
		hash := HashLeaf(claim.Address, claim.Amount, claim.Index)
		leaves[i] = &MerkleNode{
			Hash: hash,
			Data: &claims[i],
		}
	}

	tree.Leaves = leaves

	// Build the tree bottom-up
	tree.Root = tree.buildTree(leaves)

	return tree, nil
}

// buildTree recursively builds the Merkle tree
func (mt *MerkleTree) buildTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var nextLevel []*MerkleNode

	// Process pairs of nodes
	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		var right *MerkleNode

		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			// Odd number of nodes, duplicate the last one
			right = left
		}

		// Create parent node
		parentHash := HashInternal(left.Hash, right.Hash)
		parent := &MerkleNode{
			Hash:  parentHash,
			Left:  left,
			Right: right,
		}

		nextLevel = append(nextLevel, parent)
	}

	return mt.buildTree(nextLevel)
}

// GetRootHash returns the root hash as hex string
func (mt *MerkleTree) GetRootHash() string {
	if mt.Root == nil {
		return ""
	}
	return fmt.Sprintf("0x%x", mt.Root.Hash)
}
