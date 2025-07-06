// pkg/merkle/optimized.go
package merkle

import (
	"encoding/binary"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// HashPool manages reusable byte slices
var HashPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 128)
		return &b
	},
}

func OptimizedHashLeaf(address common.Address, amount *big.Int, index uint32) []byte {
	dataPtr := HashPool.Get().(*[]byte)
	defer HashPool.Put(dataPtr)
	data := *dataPtr

	// Use the existing data slice
	data = data[:0] // Reset length but keep capacity
	data = data[:0] // Reset length but keep capacity

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

	// Return a copy since we're returning the slice to pool
	result := make([]byte, 32)
	hash := crypto.Keccak256(data)
	copy(result, hash)
	return result
}

// BatchProcessor handles batch processing of claims
type BatchProcessor struct {
	BatchSize int
	Workers   int
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize, workers int) *BatchProcessor {
	return &BatchProcessor{
		BatchSize: batchSize,
		Workers:   workers,
	}
}

// ProcessClaims processes claims in batches with worker pools
func (bp *BatchProcessor) ProcessClaims(claims []AirdropClaim, processFn func([]AirdropClaim) error) error {
	jobs := make(chan []AirdropClaim, bp.Workers)
	results := make(chan error, bp.Workers)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < bp.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range jobs {
				results <- processFn(batch)
			}
		}()
	}

	// Send batches
	go func() {
		defer close(jobs)
		for i := 0; i < len(claims); i += bp.BatchSize {
			end := i + bp.BatchSize
			if end > len(claims) {
				end = len(claims)
			}
			jobs <- claims[i:end]
		}
	}()

	// Wait for completion
	go func() {
		wg.Wait()
		close(results)
	}()

	// Check for errors
	for err := range results {
		if err != nil {
			return err
		}
	}

	return nil
}
