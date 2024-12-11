// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IUpgradeEntrypoint {
    /// PROTO: https://github.com/cosmos/cosmos-sdk/blob/v0.50.9/proto/cosmos/upgrade/v1beta1/upgrade.proto
    /// @notice Emitted when an upgrade is submitted.
    /// @param name Sets the name for the upgrade. This name will be used by the upgraded version of the software to
    /// apply any special "on-upgrade" commands during the first BeginBlock method after the upgrade is applied. It is
    /// also used to detect whether a software version can handle a given upgrade. If no upgrade handler with this name
    /// has been set in the software, it will be assumed that the software is out-of-date when the upgrade Time or
    /// Height is reached and the software will exit.
    /// @param height The height at which the upgrade must be performed.
    /// @param info Any application specific upgrade info to be included on-chain such as a git commit that validators
    /// could automatically upgrade to.
    event SoftwareUpgrade(string name, int64 height, string info);

    /// @notice Emitted when a planned upgrade is to be cancelled.
    event CancelUpgrade();

    /// @notice Submits an upgrade plan.
    /// @param name Sets the name for the upgrade. This name will be used by the upgraded version of the software to
    /// apply any special "on-upgrade" commands during the first BeginBlock method after the upgrade is applied. It is
    /// also used to detect whether a software version can handle a given upgrade. If no upgrade handler with this name
    /// has been set in the software, it will be assumed that the software is out-of-date when the upgrade Time or
    /// Height is reached and the software will exit.
    /// @param height The height at which the upgrade must be performed.
    /// @param info Any application specific upgrade info to be included on-chain such as a git commit that validators
    /// could automatically upgrade to.
    function planUpgrade(string calldata name, int64 height, string calldata info) external;

    /// @notice Cancels an upgrade plan if there is one planned. Otherwise, it does nothing.
    function cancelUpgrade() external;
}
