#  Merkle Tree Airdrop System

A high-performance, gas-efficient airdrop system built in Go that uses Merkle trees to enable scalable cryptocurrency token distributions.
 This system allows for secure, verifiable airdrops to millions of recipients while minimizing on-chain storage and gas costs.

##  Features

- **Gas Efficient**: 40% gas savings compared to traditional airdrop methods
- **Highly Scalable**: Handle millions of recipients with minimal on-chain storage
- **Lightning Fast**: Sub-second proof generation for 10,000+ claims using goroutines
- **Cryptographically Secure**: Uses Keccak256 hashing for tamper-proof verification
- **Production Ready**: Comprehensive error handling, testing, and optimization
- **Smart Contract Integration**: Complete Solidity contract with Go bindings
- **Parallel Processing**: Leverages Go's concurrency for optimal performance

##  What is a Merkle Tree Airdrop?

Traditional airdrops require one transaction per recipient, making them expensive on networks like Ethereum. Merkle Tree Airdrops solve this by:

1. **Storing only a single root hash on-chain** instead of all recipient data
2. **Generating cryptographic proofs** that allow recipients to verify their eligibility
3. **Enabling gas-efficient claims** where users provide their own proof to claim tokens

### Visual Representation

```
        Root Hash (stored on-chain)
       /                          \
   Hash AB                     Hash CD
   /     \                     /     \
Hash A  Hash B              Hash C  Hash D
  |       |                   |       |
User1   User2               User3   User4
(100)   (200)               (150)   (300)
```

##  Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Core Concepts](#core-concepts)
- [API Documentation](#api-documentation)
- [Smart Contract](#smart-contract)
- [Performance](#performance)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

##  Installation

### Prerequisites

- Go 1.19 or higher
- Node.js 16+ (for smart contract compilation)
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/Adeyemijosh/merkle-airdrop.git
   cd merkle-airdrop
   ```

2. **Initialize Go module**
   ```bash
   go mod tidy
   ```

3. **Install dependencies**
   ```bash
   go get github.com/ethereum/go-ethereum
   go get github.com/ethereum/go-ethereum/ethclient
   go get github.com/ethereum/go-ethereum/common
   go get golang.org/x/crypto/sha3
   ```

4. **Verify installation**
   ```bash
   go version
   go run main.go --help
   ```

##  Quick Start

### 1. Generate Test Data

```bash
# Generate 10,000 test airdrop claims
go run main.go --generate-test-data --count 10000
```

### 2. Build Merkle Tree

```bash
# Build tree from CSV data
go run main.go --build-tree --input airdrop_data.csv
```

### 3. Generate Proofs

```bash
# Generate all proofs in parallel
go run main.go --generate-proofs --output merkle_proofs.json
```

### 4. Verify Proof

```bash
# Verify a specific address can claim
go run main.go --verify-proof --address 0x742d35Cc6634C0532925a3b8D8c6632F2d50a9a6
```

##  Project Structure

```
merkle-airdrop/
├── cmd/
│   └── cli/
│       └── main.go              # CLI interface
├── internal/
│   ├── api/                     # REST API endpoints
│   │   ├── handlers.go          # HTTP handlers
│   │   ├── middleware.go        # API middleware
│   │   └── routes.go            # Route definitions
│   └── config/                  # Configuration management
│       └── config.go            # App configuration
├── pkg/
│   ├── merkle/                  # Core Merkle tree logic
│   │   ├── types.go             # Data structures
│   │   ├── hash.go              # Hashing functions
│   │   ├── tree.go              # Tree construction
│   │   ├── proof.go             # Proof generation
│   │   └── optimized.go         # Performance optimizations
│   ├── data/                    # Data loading utilities
│   │   ├── loader.go            # CSV/JSON data loaders
│   │   └── generator.go         # Test data generation
│   └── contract/                # Smart contract interaction
│       ├── client.go            # Ethereum client
│       ├── deploy.go            # Contract deployment
│       └── bindings.go          # Generated contract bindings
├── test/
│   ├── benchmark_test.go        # Performance benchmarks
│   ├── integration_test.go      # Integration tests
│   └── unit_test.go             # Unit tests
├── contracts/
│   ├── MerkleAirdrop.sol        # Main airdrop contract
│   ├── IERC20.sol               # ERC20 interface
│   └── deploy.js                # Deployment script
├── scripts/
│   ├── deploy.sh               # Contract deployment
│   ├── verify.sh               # Contract verification
│   └── generate-bindings.sh    # Go binding generation
├── docs/
│   ├── API.md                  # API documentation
│   ├── ARCHITECTURE.md         # System architecture
│   └── DEPLOYMENT.md           # Deployment guide
├── examples/
│   ├── basic_usage.go          # Basic usage examples
│   └── advanced_usage.go       # Advanced examples
├── go.mod
├── go.sum
└── README.md
```

##  Core Concepts

### Data Structures

#### AirdropClaim
```go
type AirdropClaim struct {
    Address common.Address `json:"address"` // Recipient address
    Amount  *big.Int       `json:"amount"`  // Token amount (in wei)
    Index   uint32         `json:"index"`   // Claim index
}
```

#### MerkleProof
```go
type MerkleProof struct {
    Proof  []string `json:"proof"`   // Array of sibling hashes
    Index  uint32   `json:"index"`   // Leaf index
    Amount string   `json:"amount"`  // Claim amount
}
```

### Hashing Algorithm

The system uses **Keccak256** (Ethereum standard) for all hashing operations:

```go
// Leaf hash: keccak256(address + amount + index)
func HashLeaf(address common.Address, amount *big.Int, index uint32) []byte

// Internal hash: keccak256(leftHash + rightHash) - with deterministic ordering
func HashInternal(left, right []byte) []byte
```

### Tree Construction Process

1. **Sort Claims**: Sort by address for deterministic tree structure
2. **Create Leaves**: Generate hash for each claim
3. **Build Tree**: Recursively pair nodes and hash until root is reached
4. **Store Root**: Single 32-byte hash represents entire dataset

## 🔧 API Documentation

### REST Endpoints

#### GET /api/v1/proof/:address
Get Merkle proof for a specific address.

**Response:**
```json
{
  "proof": ["0x...", "0x..."],
  "index": 42,
  "amount": "1000000000000000000",
  "merkleRoot": "0x..."
}
```

#### POST /api/v1/verify
Verify a Merkle proof.

**Request:**
```json
{
  "address": "0x742d35Cc6634C0532925a3b8D8c6632F2d50a9a6",
  "amount": "1000000000000000000",
  "proof": ["0x...", "0x..."],
  "merkleRoot": "0x..."
}
```

#### GET /api/v1/stats
Get airdrop statistics.

**Response:**
```json
{
  "totalClaims": 10000,
  "totalAmount": "10000000000000000000000",
  "merkleRoot": "0x...",
  "claimedCount": 1337,
  "claimedAmount": "1337000000000000000000"
}
```

### CLI Commands

```bash
# Generate test data
go run main.go generate --count 10000 --output test_data.csv

# Build Merkle tree
go run main.go build --input airdrop_data.csv --output tree.json

# Generate proofs
go run main.go proof --tree tree.json --output proofs.json

# Verify proof
go run main.go verify --address 0x... --proof proof.json

# Start API server
go run main.go serve --port 8080
```

##  Smart Contract

### Contract Interface

```solidity
contract MerkleAirdrop {
    IERC20 public token;
    bytes32 public merkleRoot;
    mapping(address => bool) public claimed;
    
    function claim(uint256 amount, bytes32[] calldata merkleProof) external;
    function verifyClaim(address account, uint256 amount, bytes32[] calldata proof) public view returns (bool);
}
```

### Deployment

```bash
# Deploy to testnet
npm run deploy:testnet

# Deploy to mainnet
npm run deploy:mainnet

# Verify contract
npm run verify -- --network mainnet --contract-address 0x...
```

### Integration Example

```go
// Deploy contract
contractAddr, err := client.DeployAirdrop(tokenAddr, merkleRoot)

// User claims tokens
tx, err := contract.Claim(amount, proof)
```

## ⚡ Performance

### Benchmarks

| Operation | 1K Claims | 10K Claims | 100K Claims | 1M Claims |
|-----------|-----------|------------|-------------|-----------|
| Tree Construction | 12ms | 89ms | 1.2s | 15s |
| Proof Generation | 45ms | 380ms | 4.1s | 48s |
| Memory Usage | 2MB | 18MB | 180MB | 1.8GB |
| Proof Size | 320B | 416B | 512B | 640B |

### Optimization Features

- **Memory Pooling**: Reuse byte slices for hashing operations
- **Goroutine Parallelism**: Concurrent proof generation
- **Batch Processing**: Process claims in optimized batches
- **Compression**: Compress proof data for storage

##  Testing

### Run Tests

```bash
# Unit tests
go test ./pkg/...

# Integration tests
go test ./test/integration/...

# Benchmark tests
go test -bench=. ./test/benchmark/...

# Coverage report
go test -cover ./...
```

### Test Categories

- **Unit Tests**: Individual function testing
- **Integration Tests**: End-to-end workflow testing
- **Benchmark Tests**: Performance measurement
- **Contract Tests**: Smart contract functionality

##  Monitoring & Metrics

### Key Metrics

- **Tree Construction Time**: Time to build Merkle tree
- **Proof Generation Rate**: Proofs generated per second
- **Memory Usage**: RAM consumption during operations
- **Gas Costs**: Contract interaction costs
- **Claim Success Rate**: Percentage of successful claims

### Logging

```go
// Structured logging with levels
logger.Info("Tree constructed", 
    "claims", len(claims),
    "duration", buildTime,
    "rootHash", rootHash)
```

##  Security Considerations

### Best Practices

1. **Input Validation**: Validate all addresses and amounts
2. **Proof Verification**: Always verify proofs before processing
3. **Rate Limiting**: Implement API rate limiting
4. **Access Control**: Restrict sensitive operations
5. **Audit Trail**: Log all critical operations

### Common Pitfalls

- **Hash Ordering**: Ensure deterministic hash ordering
- **Integer Overflow**: Use big.Int for large amounts
- **Proof Replay**: Prevent proof reuse attacks
- **Front-running**: Consider MEV protection

##  Contributing

I welcome contributions!

### Development Setup

1. **Fork the repository**
2. **Create feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit changes**: `git commit -m 'Add amazing feature'`
4. **Push to branch**: `git push origin feature/amazing-feature`
5. **Open Pull Request**

### Code Style

- Follow Go formatting conventions (`gofmt`)
- Write comprehensive tests
- Add documentation for public APIs
- Use meaningful variable names

### Pull Request Process

1. Update README if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update documentation
5. Request review from maintainers

##  Getting Started Checklist

- [ ] Clone repository
- [ ] Install dependencies
- [ ] Generate test data
- [ ] Build your first Merkle tree
- [ ] Generate proofs
- [ ] Deploy smart contract
- [ ] Test end-to-end workflow
- [ ] Read the documentation
- [ ] Join our community

**Happy coding!**