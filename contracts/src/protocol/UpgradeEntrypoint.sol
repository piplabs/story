// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

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
    /// @param name Sets the name for the upgrade. This name will be used by the upgraded version of the software to
    /// apply any special "on-upgrade" commands during the first BeginBlock method after the upgrade is applied. It is
    /// also used to detect whether a software version can handle a given upgrade. If no upgrade handler with this name
    /// has been set in the software, it will be assumed that the software is out-of-date when the upgrade Time or
    /// Height is reached and the software will exit.
    /// @param height The height at which the upgrade must be performed.
    /// @param info Any application specific upgrade info to be included on-chain such as a git commit that validators
    /// could automatically upgrade to.
    function planUpgrade(string calldata name, int64 height, string calldata info) external onlyOwner {
        emit SoftwareUpgrade({ name: name, height: height, info: info });
    }
}
