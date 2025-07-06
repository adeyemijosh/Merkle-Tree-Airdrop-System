// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract MerkleAirdrop is Ownable, ReentrancyGuard {
    IERC20 public immutable token;
    bytes32 public immutable merkleRoot;

    mapping(address => bool) public claimed;
    uint256 public totalClaimed;
    uint256 public claimDeadline;

    event Claimed(address indexed account, uint256 amount, uint256 index);
    event DeadlineUpdated(uint256 newDeadline);

    constructor(
        IERC20 _token,
        bytes32 _merkleRoot,
        uint256 _claimDeadline,
        address initialOwner
    ) Ownable(initialOwner) {
        token = _token;
        merkleRoot = _merkleRoot;
        claimDeadline = _claimDeadline;
    }

    function claim(
        uint256 index,
        address account,
        uint256 amount,
        bytes32[] calldata merkleProof
    ) external nonReentrant {
        require(block.timestamp <= claimDeadline, "Claim period ended");
        require(!claimed[account], "Already claimed");
        require(account == msg.sender, "Can only claim for yourself");

        // Verify the merkle proof
        bytes32 leaf = keccak256(
            bytes.concat(
                keccak256(
                    abi.encode(
                        keccak256(bytes.concat(bytes20(account))),
                        keccak256(abi.encode(amount)),
                        keccak256(abi.encode(index))
                    )
                )
            )
        );

        require(
            MerkleProof.verify(merkleProof, merkleRoot, leaf),
            "Invalid merkle proof"
        );

        // Mark as claimed and transfer tokens
        claimed[account] = true;
        totalClaimed += amount;

        require(token.transfer(account, amount), "Transfer failed");

        emit Claimed(account, amount, index);
    }

    function updateDeadline(uint256 _newDeadline) external onlyOwner {
        require(_newDeadline > block.timestamp, "Deadline must be in future");
        claimDeadline = _newDeadline;
        emit DeadlineUpdated(_newDeadline);
    }

    function withdrawUnclaimed() external onlyOwner {
        require(block.timestamp > claimDeadline, "Claim period not ended");
        uint256 balance = token.balanceOf(address(this));
        require(token.transfer(owner(), balance), "Transfer failed");
    }

    function getClaimStatus(address account) external view returns (bool) {
        return claimed[account];
    }
}
