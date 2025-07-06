// pkg/data/helpers.go
package data

import (
	"encoding/csv"
	"fmt"
	"os"

	"merkle-airdrop/pkg/merkle"
)

// SaveClaimsToCSV saves airdrop claims to a CSV file
func SaveClaimsToCSV(claims []merkle.AirdropClaim, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"address", "amount"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write claims
	for _, claim := range claims {
		record := []string{
			claim.Address.Hex(),
			claim.Amount.String(),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// ValidateClaimsData validates airdrop claims data
func ValidateClaimsData(claims []merkle.AirdropClaim) error {
	if len(claims) == 0 {
		return fmt.Errorf("no claims provided")
	}

	addressMap := make(map[string]bool)
	for i, claim := range claims {
		// Check for duplicate addresses
		addrHex := claim.Address.Hex()
		if addressMap[addrHex] {
			return fmt.Errorf("duplicate address at index %d: %s", i, addrHex)
		}
		addressMap[addrHex] = true

		// Check for zero amounts
		if claim.Amount.Sign() <= 0 {
			return fmt.Errorf("invalid amount at index %d: %s", i, claim.Amount.String())
		}

		// Check for zero address
		if claim.Address.Hex() == "0x0000000000000000000000000000000000000000" {
			return fmt.Errorf("zero address at index %d", i)
		}
	}

	return nil
}

// FilterClaims filters claims based on various criteria
func FilterClaims(claims []merkle.AirdropClaim, filters ClaimFilters) []merkle.AirdropClaim {
	var filtered []merkle.AirdropClaim

	for _, claim := range claims {
		if shouldIncludeClaim(claim, filters) {
			filtered = append(filtered, claim)
		}
	}

	return filtered
}

// ClaimFilters defines filtering criteria for claims
type ClaimFilters struct {
	MinAmount      string   // Minimum amount to include
	MaxAmount      string   // Maximum amount to include
	ExcludeAddress []string // Addresses to exclude
	IncludeAddress []string // Only include these addresses (if specified)
}

func shouldIncludeClaim(claim merkle.AirdropClaim, filters ClaimFilters) bool {
	// Check exclude list
	addrHex := claim.Address.Hex()
	for _, excluded := range filters.ExcludeAddress {
		if addrHex == excluded {
			return false
		}
	}

	// Check include list (if specified)
	if len(filters.IncludeAddress) > 0 {
		found := false
		for _, included := range filters.IncludeAddress {
			if addrHex == included {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// TODO: Add amount filtering logic

	return true
}

// DeduplicateClaims removes duplicate claims (keeps first occurrence)
func DeduplicateClaims(claims []merkle.AirdropClaim) []merkle.AirdropClaim {
	seen := make(map[string]bool)
	var deduplicated []merkle.AirdropClaim

	for _, claim := range claims {
		addrHex := claim.Address.Hex()
		if !seen[addrHex] {
			seen[addrHex] = true
			deduplicated = append(deduplicated, claim)
		}
	}

	return deduplicated
}

// SplitClaims splits claims into batches for processing
func SplitClaims(claims []merkle.AirdropClaim, batchSize int) [][]merkle.AirdropClaim {
	if batchSize <= 0 {
		return [][]merkle.AirdropClaim{claims}
	}

	var batches [][]merkle.AirdropClaim
	for i := 0; i < len(claims); i += batchSize {
		end := i + batchSize
		if end > len(claims) {
			end = len(claims)
		}
		batches = append(batches, claims[i:end])
	}

	return batches
}
