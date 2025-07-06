// pkg/data/loader.go
package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"os"

	"merkle-airdrop/pkg/merkle"

	"github.com/ethereum/go-ethereum/common"
)

// LoadAirdropFromCSV loads airdrop data from CSV file
// Expected format: address,amount
func LoadAirdropFromCSV(filename string) ([]merkle.AirdropClaim, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2 // address, amount

	var claims []merkle.AirdropClaim

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	index := uint32(0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		// Parse address
		if !common.IsHexAddress(record[0]) {
			return nil, fmt.Errorf("invalid address: %s", record[0])
		}
		address := common.HexToAddress(record[0])

		// Parse amount
		amount, ok := new(big.Int).SetString(record[1], 10)
		if !ok {
			return nil, fmt.Errorf("invalid amount: %s", record[1])
		}

		claims = append(claims, merkle.AirdropClaim{
			Address: address,
			Amount:  amount,
			Index:   index,
		})

		index++
	}

	return claims, nil
}

// GenerateTestData creates test airdrop data
func GenerateTestData(count int) []merkle.AirdropClaim {
	claims := make([]merkle.AirdropClaim, count)

	for i := 0; i < count; i++ {
		// Generate pseudo-random address
		address := common.HexToAddress(fmt.Sprintf("0x%040d", i+1))

		// Generate amount (1-1000 tokens with 18 decimals)
		amount := big.NewInt(int64(i%1000 + 1))
		amount.Mul(amount, big.NewInt(1e18)) // 18 decimals

		claims[i] = merkle.AirdropClaim{
			Address: address,
			Amount:  amount,
			Index:   uint32(i),
		}
	}

	return claims
}
