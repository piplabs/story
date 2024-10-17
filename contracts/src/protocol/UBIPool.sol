// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { Math } from "@openzeppelin/contracts/utils/math/Math.sol";
import { PubKeyVerification } from "./PubKeyVerification.sol";
import { IUBIPool } from "../interfaces/IUBIPool.sol";

contract UBIPool is IUBIPool, Ownable2StepUpgradeable, ReentrancyGuardUpgradeable, PubKeyVerification {

    uint256 public immutable MAX_UBI_PERCENTAGE;

    mapping(uint256 month => mapping(bytes validatorUncmpPubkey => uint256 amount)) public validatorUBIAmounts;

    constructor(uint256 maxUBIPercentage) {
        MAX_UBI_PERCENTAGE = maxUBIPercentage;
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(address owner) public initializer {
        require(owner != address(0), "UBIPool: owner cannot be zero address");
        __Ownable_init(owner);
    }

    function setUBIPercentage(uint256 percentage) external onlyOwner {
        require(percentage < MAX_UBI_PERCENTAGE, "UBIPool: percentage too high");
        emit UBIPercentageSet(percentage);
    }

    function setUBIDistribution(
        uint256 month,
        uint256 totalUBI,
        bytes[] calldata validatorUncmpPubKeys,
        uint256[] calldata percentages
    ) external onlyOwner {
        require(validatorUncmpPubKeys.length > 0, "UBIPool: validatorUncmpPubKeys cannot be empty");
        require(validatorUncmpPubKeys.length == percentages.length, "UBIPool: length mismatch");
        require(totalUBI <= address(this).balance, "UBIPool: not enough balance");
        uint256 totalPercent;
        uint256 accAmount;
        for (uint256 i = 0; i < percentages.length; i++) {
            require(percentages[i] > 0, "UBIPool: percentage cannot be zero");
            totalPercent += percentages[i];
            uint256 amount = Math.mulDiv(totalUBI, percentages[i], 100 ether);
            validatorUBIAmounts[month][validatorUncmpPubKeys[i]] = amount;
            accAmount += amount;
        }
        require(totalPercent == 100 ether, "UBIPool: Percentages must sum to 100%");
        require(accAmount == totalUBI, "UBIPool: total amount mismatch");
        emit UBIDistributionSet(month, totalUBI, validatorUncmpPubKeys, percentages);
    }

    function claimUBI(
        uint256 month,
        bytes calldata validatorUncmpPubkey
    ) external nonReentrant verifyUncmpPubkeyWithExpectedAddress(validatorUncmpPubkey, msg.sender) {
        uint256 amount = validatorUBIAmounts[month][validatorUncmpPubkey];
        require(amount > 0, "UBIPool: no UBI to claim");
        validatorUBIAmounts[month][validatorUncmpPubkey] = 0;
        (bool success, ) = msg.sender.call{ value: amount }("");
        require(success, "UBIPool: failed to send UBI");
    }
}
