// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

import { IUpgradeEntrypoint } from "../interfaces/IUpgradeEntrypoint.sol";

/**
 * @title UpgradeEntrypoint
 * @notice Entrypoint contract for submitting x/upgrade module actions.
 */
contract UpgradeEntrypoint is IUpgradeEntrypoint, Ownable2StepUpgradeable {
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(address owner) public initializer {
        require(owner != address(0), "UpgradeEntrypoint: owner cannot be zero address");
        __Ownable_init(owner);
    }

    /// @notice Submits an upgrade plan.
    /// @param appVersion Sets the app version for the upgrade.
    /// @param upgradeHeight The height at which the upgrade must be performed.
    /// @param info Any application specific upgrade info to be included on-chain such as a git commit that validators
    /// could automatically upgrade to.
    function planUpgrade(uint64 appVersion, uint64 upgradeHeight, string calldata info) external onlyOwner {
        emit SoftwareUpgrade({ appVersion: appVersion, upgradeHeight: upgradeHeight, info: info });
    }
}
