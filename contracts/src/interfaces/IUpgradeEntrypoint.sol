// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

interface IUpgradeEntrypoint {
    /// PROTO: https://github.com/cosmos/cosmos-sdk/blob/v0.50.9/proto/cosmos/upgrade/v1beta1/upgrade.proto
    /// @notice Emitted when an upgrade is submitted.
    /// @param appVersion The app version of the upgrade.
    /// @param upgradeHeight The height at which the upgrade must be performed.
    /// @param info Any application specific upgrade info to be included on-chain such as a git commit that validators
    /// could automatically upgrade to.
    event SoftwareUpgrade(uint64 appVersion, uint64 upgradeHeight, string info);

    /// @notice Submits an upgrade plan.
    /// @param appVersion Sets the app version for the upgrade.
    /// @param upgradeHeight The height at which the upgrade must be performed.
    /// @param info Any application specific upgrade info to be included on-chain such as a git commit that validators
    /// could automatically upgrade to.
    function planUpgrade(uint64 appVersion, uint64 upgradeHeight, string calldata info) external;
}
