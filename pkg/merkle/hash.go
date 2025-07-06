// pkg/merkle/hash.go
package merkle

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// HashLeaf creates a hash for a leaf node (address + amount)
func HashLeaf(address common.Address, amount *big.Int, index uint32) []byte {
	// Create a buffer to hold our data
	data := make([]byte, 0, 32+32+4) // address(20) + amount(32) + index(4)

	// Add address (pad to 32 bytes)
	addressBytes := make([]byte, 32)
	copy(addressBytes[12:], address.Bytes()) // Ethereum addresses are 20 bytes
	data = append(data, addressBytes...)

	// Add amount (pad to 32 bytes)
	amountBytes := make([]byte, 32)
	amount.FillBytes(amountBytes)
	data = append(data, amountBytes...)

	// Add index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	data = append(data, indexBytes...)

	// Return Keccak256 hash (Ethereum standard)
	return crypto.Keccak256(data)
}

// HashInternal creates a hash for internal nodes
func HashInternal(left, right []byte) []byte {
	// Sort hashes to ensure deterministic tree
	if len(left) != 32 || len(right) != 32 {
		panic("Invalid hash length")
	}

	var data []byte
	// Smaller hash goes first for deterministic ordering
	if string(left) < string(right) {
		data = append(left, right...)
	} else {
		data = append(right, left...)
	}

	return crypto.Keccak256(data)
}
