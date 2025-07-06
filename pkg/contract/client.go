// pkg/contract/client.go
package contract

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ContractClient handles Ethereum contract interactions
type ContractClient struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
}

// NewContractClient creates a new contract client
func NewContractClient(rpcURL, privateKeyHex string) (*ContractClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	return &ContractClient{
		client:     client,
		privateKey: privateKey,
		chainID:    chainID,
	}, nil
}

// DeployAirdrop deploys the airdrop contract
func (cc *ContractClient) DeployAirdrop(tokenAddress common.Address, merkleRoot [32]byte) (common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(cc.privateKey, cc.chainID)
	if err != nil {
		return common.Address{}, err
	}

	// Set gas limit and price
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(20000000000) // 20 gwei

	// Deploy contract (you'll need to generate Go bindings from the Solidity contract)
	// This is a simplified example
	address := common.HexToAddress("0x...") // Contract address after deployment

	return address, nil
}
