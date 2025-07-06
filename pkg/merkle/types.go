// pkg/merkle/types.go
package merkle

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// AirdropClaim represents a single airdrop entry
type AirdropClaim struct {
	Address common.Address `json:"address"`
	Amount  *big.Int       `json:"amount"`
	Index   uint32         `json:"index"`
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Hash  []byte
	Left  *MerkleNode
	Right *MerkleNode
	Data  *AirdropClaim // Only for leaf nodes
}

// MerkleTree represents the complete Merkle tree
type MerkleTree struct {
	Root   *MerkleNode
	Leaves []*MerkleNode
	Claims []AirdropClaim
}

// MerkleProof represents the proof needed to verify a claim
type MerkleProof struct {
	Proof  []string `json:"proof"`
	Index  uint32   `json:"index"`
	Amount string   `json:"amount"`
}
