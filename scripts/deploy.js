const { ethers, network, run } = require("hardhat");
const fs = require("fs");

async function main() {
    console.log("🚀 Starting Merkle Airdrop deployment...");

    const merkleData = JSON.parse(fs.readFileSync("merkle_proofs.json", "utf8"));
    const merkleRoot = merkleData.merkleRoot;

    console.log("🌱 Merkle Root:", merkleRoot);
    console.log("📦 Total Claims:", merkleData.totalClaims);

    const tokenAddress = "0x3eDF60dd017aCe33A0220F78741b5581C385A1BA"; // Replace with actual address

    const claimDeadline = Math.floor(Date.now() / 1000) + (30 * 24 * 60 * 60); // 30 days from now
    const [deployer] = await ethers.getSigners();

    const MerkleAirdrop = await ethers.getContractFactory("MerkleAirdrop");

    const airdrop = await MerkleAirdrop.deploy(
        tokenAddress,
        merkleRoot,
        claimDeadline,
        deployer.address
    );

    // ✅ Wait for deployment to finish
    await airdrop.waitForDeployment();

    const deployedAddress = await airdrop.getAddress(); // ✅ use getAddress() with ethers v6+

    console.log("✅ MerkleAirdrop deployed to:", deployedAddress);
    console.log("📍 Merkle Root:", await airdrop.merkleRoot());
    console.log("⏳ Claim Deadline:", new Date(claimDeadline * 1000).toISOString());

    // Save deployment info
    const deploymentInfo = {
        contractAddress: deployedAddress,
        merkleRoot,
        tokenAddress,
        claimDeadline,
        totalClaims: merkleData.totalClaims,
        deployedAt: new Date().toISOString(),
        network: network.name
    };

    fs.writeFileSync("deployment.json", JSON.stringify(deploymentInfo, null, 2));
    console.log("📝 Deployment info saved to deployment.json");

    // Optional: Verify
    if (network.name !== "hardhat" && network.name !== "localhost") {
        console.log("⛓️ Waiting for block confirmations...");
        await airdrop.deploymentTransaction().wait(6);

        console.log("🔍 Verifying contract on Etherscan...");
        await run("verify:verify", {
            address: deployedAddress,
            constructorArguments: [
                tokenAddress,
                merkleRoot,
                claimDeadline,
                deployer.address
            ]
        });
    }
}

main().catch((err) => {
    console.error("❌ Error during deployment:", err);
    process.exit(1);
});
