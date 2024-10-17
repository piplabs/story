// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

interface IUBIPool {
    event UBIPercentageSet(uint256 percentage);
    event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] percentages);

    function setUBIPercentage(uint256 percentage) external;

    function setUBIDistribution(
        uint256 month,
        uint256 totalUBI,
        bytes[] calldata validatorUncmpPubKeys,
        uint256[] calldata percentages
    ) external;

    function claimUBI(uint256 month, bytes calldata validatorUncmpPubkey) external;
}
