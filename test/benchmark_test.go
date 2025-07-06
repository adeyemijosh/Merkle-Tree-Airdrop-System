// test/benchmark_test.go
package test

import (
	"merkle-airdrop/pkg/data"
	"merkle-airdrop/pkg/merkle"
	"testing"
)

func BenchmarkTreeConstruction(b *testing.B) {
	claims := data.GenerateTestData(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := merkle.NewMerkleTree(claims)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProofGeneration(b *testing.B) {
	claims := data.GenerateTestData(10000)
	tree, _ := merkle.NewMerkleTree(claims)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tree.GenerateAllProofs()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMerkleTreeCorrectness(t *testing.T) {
	claims := data.GenerateTestData(100)
	tree, err := merkle.NewMerkleTree(claims)
	if err != nil {
		t.Fatal(err)
	}

	// Test proof generation for all claims
	for _, claim := range claims {
		proof, err := tree.GenerateProof(claim.Address)
		if err != nil {
			t.Errorf("Failed to generate proof for %s: %v", claim.Address.Hex(), err)
		}

		if proof.Amount != claim.Amount.String() {
			t.Errorf("Amount mismatch for %s", claim.Address.Hex())
		}
	}
}
