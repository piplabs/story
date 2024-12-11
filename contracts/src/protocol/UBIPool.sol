// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { MulticallUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/MulticallUpgradeable.sol";
import { PubKeyVerifier } from "./PubKeyVerifier.sol";
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
    PubKeyVerifier,
    MulticallUpgradeable
{
    /// @notice The maximum UBI percentage
    uint32 public immutable MAX_UBI_PERCENTAGE;

    /// @notice The current distribution id, incremented for each new distribution
    uint256 public currentDistributionId;

    /// @notice The amount of UBI for each validator for a given distribution
    mapping(uint256 distributionId => mapping(bytes validatorUncmpPubkey => uint256 amount)) public validatorUBIAmounts;

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
    /// @param validatorUncmpPubKeys The validator uncompressed public keys
    /// @param amounts The amounts of UBI for each validator
    /// @return distributionId The distribution id
    function setUBIDistribution(
        uint256 totalUBI,
        bytes[] calldata validatorUncmpPubKeys,
        uint256[] calldata amounts
    ) external onlyOwner returns (uint256) {
        require(validatorUncmpPubKeys.length > 0, "UBIPool: validatorUncmpPubKeys cannot be empty");
        require(validatorUncmpPubKeys.length == amounts.length, "UBIPool: length mismatch");
        require(totalUBI <= address(this).balance, "UBIPool: not enough balance");
        uint256 accAmount;
        currentDistributionId++;
        for (uint256 i = 0; i < amounts.length; i++) {
            require(amounts[i] > 0, "UBIPool: amounts cannot be zero");
            _verifyUncmpPubkey(validatorUncmpPubKeys[i]);
            validatorUBIAmounts[currentDistributionId][validatorUncmpPubKeys[i]] = amounts[i];
            accAmount += amounts[i];
        }
        require(accAmount == totalUBI, "UBIPool: total amount mismatch");
        emit UBIDistributionSet(currentDistributionId, totalUBI, validatorUncmpPubKeys, amounts);
        return currentDistributionId;
    }

    /// @notice Claims the UBI for a given month for a validator
    /// @dev The validator address must be the one who is set to receive the UBI
    /// @param distributionId The distribution id
    /// @param validatorUncmpPubkey The validator uncompressed public key
    function claimUBI(
        uint256 distributionId,
        bytes calldata validatorUncmpPubkey
    ) external nonReentrant verifyUncmpPubkeyWithExpectedAddress(validatorUncmpPubkey, msg.sender) {
        uint256 amount = validatorUBIAmounts[distributionId][validatorUncmpPubkey];
        require(amount > 0, "UBIPool: no UBI to claim");
        validatorUBIAmounts[distributionId][validatorUncmpPubkey] = 0;
        (bool success, ) = msg.sender.call{ value: amount }("");
        require(success, "UBIPool: failed to send UBI");
    }
}
