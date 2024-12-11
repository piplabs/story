// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

/// @title IUBIPool
/// @notice Interface for the UBI Pool contract
interface IUBIPool {
    /// @notice Emitted when the UBI percentage is set
    /// @param percentage The percentage of the UBI
    event UBIPercentageSet(uint32 percentage);

    /// @notice Emitted when the UBI distribution is set
    /// @param totalUBI The total amount of UBI
    /// @param validatorUncmpPubKeys The validator uncompressed public keys
    /// @param amounts The amounts of the UBI for each validator
    event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] amounts);

    /// @notice Sets the UBI percentage
    /// @param percentage The percentage of the UBI
    function setUBIPercentage(uint32 percentage) external;

    /// @notice Sets the UBI distribution
    /// @param totalUBI The total amount of UBI
    /// @param validatorUncmpPubKeys The validator uncompressed public keys
    /// @param amounts The amounts of the UBI for each validator
    /// @return distributionId The distribution id
    function setUBIDistribution(
        uint256 totalUBI,
        bytes[] calldata validatorUncmpPubKeys,
        uint256[] calldata amounts
    ) external returns (uint256);

    /// @notice Claims the UBI for a validator
    /// @dev The validator address must be the one who is set to receive the UBI
    /// @param distributionId The distribution id
    /// @param validatorUncmpPubkey The validator uncompressed public key
    function claimUBI(uint256 distributionId, bytes calldata validatorUncmpPubkey) external;
}
