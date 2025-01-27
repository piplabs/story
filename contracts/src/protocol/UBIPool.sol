// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { MulticallUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/MulticallUpgradeable.sol";
import { Secp256k1Verifier } from "./Secp256k1Verifier.sol";
import { IUBIPool } from "../interfaces/IUBIPool.sol";

/// @title UBIPool
/// @notice Contract for distributing UBI to validators
/// UBI comes from a percentage of the protocol's emmission during a defined period
/// This % is minted to the UBIPool contract and can be claimed by validators
/// Each validator can claim their UBI for a given month
/// Distributions will be made public monthly by Story for community scrutiny,
/// and if correct, will be distributed to validators via this contract
contract UBIPool is
    IUBIPool,
    Ownable2StepUpgradeable,
    ReentrancyGuardUpgradeable,
    Secp256k1Verifier,
    MulticallUpgradeable
{
    /// @notice The maximum UBI percentage
    uint32 public immutable MAX_UBI_PERCENTAGE;

    /// @notice The current distribution id, incremented for each new distribution
    uint256 public currentDistributionId;

    /// @notice The amount of UBI for each validator for a given distribution
    mapping(uint256 distributionId => mapping(bytes validatorCmpPubkey => uint256 amount)) public validatorUBIAmounts;

    /// @notice The total amount of pending tokens to claim.
    /// Added when a distribution is set, and subtracted when a validator claims their UBI.
    uint256 public totalPendingClaims;

    constructor(uint32 maxUBIPercentage) {
        MAX_UBI_PERCENTAGE = maxUBIPercentage;
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    /// @param owner The owner of the contract
    function initialize(address owner) public initializer {
        require(owner != address(0), "UBIPool: owner cannot be zero address");
        __Ownable_init(owner);
    }

    /// @notice Sets the UBI percentage distribution in CL
    /// @param percentage The percentage of the UBI
    function setUBIPercentage(uint32 percentage) external onlyOwner {
        require(percentage <= MAX_UBI_PERCENTAGE, "UBIPool: percentage too high");
        emit UBIPercentageSet(percentage);
    }

    /// @notice Sets the UBI distribution for a given month
    /// @param totalUBI The total amount of UBI
    /// @param validatorCmpPubKeys The validator compressed public keys
    /// @param amounts The amounts of UBI for each validator
    /// @return distributionId The distribution id
    function setUBIDistribution(
        uint256 totalUBI,
        bytes[] calldata validatorCmpPubKeys,
        uint256[] calldata amounts
    ) external onlyOwner returns (uint256) {
        require(validatorCmpPubKeys.length > 0, "UBIPool: validatorCmpPubKeys cannot be empty");
        require(validatorCmpPubKeys.length == amounts.length, "UBIPool: length mismatch");
        require(totalUBI + totalPendingClaims <= address(this).balance, "UBIPool: not enough balance");
        totalPendingClaims += totalUBI;
        uint256 accAmount;
        currentDistributionId++;
        for (uint256 i = 0; i < amounts.length; i++) {
            require(amounts[i] > 0, "UBIPool: amounts cannot be zero");
            _verifyCmpPubkey(validatorCmpPubKeys[i]);
            validatorUBIAmounts[currentDistributionId][validatorCmpPubKeys[i]] = amounts[i];
            accAmount += amounts[i];
        }
        require(accAmount == totalUBI, "UBIPool: total amount mismatch");
        emit UBIDistributionSet(currentDistributionId, totalUBI, validatorCmpPubKeys, amounts);
        return currentDistributionId;
    }

    /// @notice Claims the UBI for a given month for a validator
    /// @dev The validator address must be the one who is set to receive the UBI
    /// @param distributionId The distribution id
    /// @param validatorCmpPubkey The validator compressed public key
    function claimUBI(
        uint256 distributionId,
        bytes calldata validatorCmpPubkey
    ) external nonReentrant verifyCmpPubkeyWithExpectedAddress(validatorCmpPubkey, msg.sender) {
        uint256 amount = validatorUBIAmounts[distributionId][validatorCmpPubkey];
        require(amount > 0, "UBIPool: no UBI to claim");
        validatorUBIAmounts[distributionId][validatorCmpPubkey] = 0;
        (bool success, ) = msg.sender.call{ value: amount }("");
        require(success, "UBIPool: failed to send UBI");
        totalPendingClaims -= amount;
    }
}
