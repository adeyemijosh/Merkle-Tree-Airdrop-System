require("@nomicfoundation/hardhat-toolbox");

module.exports = {
  solidity: "0.8.20",
  networks: {
    // Local development network (uses built-in test accounts)
    hardhat: {
      chainId: 31337,
      allowUnlimitedContractSize: true,
      mining: {
        auto: true, // Blocks mine automatically
        interval: 1000 // Optional: mine a block every 1s
      }
    },

    // Sepolia testnet with hardcoded values (replace with your actual values)
    sepolia: {
      url: "https://sepolia.infura.io/v3/your-infura-key-here",
      accounts: [
        "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // Hardhat test private key
      ],
      chainId: 11155111,
    }
  }
};

/* require("@nomicfoundation/hardhat-toolbox");

module.exports = {
  solidity: "0.8.19",
  networks: {
    sepolia: {
      url: "https://sepolia.infura.io/v3/YOUR_INFURA_KEY",
      accounts: ["YOUR_PRIVATE_KEY"]
    }
  }
};
  */
