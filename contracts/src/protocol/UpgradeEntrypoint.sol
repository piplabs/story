// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import { IUpgradeEntrypoint } from "../interfaces/IUpgradeEntrypoint.sol";

/**
 * @title UpgradeEntrypoint
 * @notice Entrypoint contract for submitting x/upgrade module actions.
 */
contract UpgradeEntrypoint is IUpgradeEntrypoint, Ownable2StepUpgradeable, UUPSUpgradeable {
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(address accessManager) public initializer {
        require(accessManager != address(0), "UpgradeEntrypoint: accessManager cannot be zero address");
        __UUPSUpgradeable_init();
        __Ownable_init(accessManager);
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

    /// @dev Hook to authorize the upgrade according to UUPSUpgradeable
    /// @param newImplementation The address of the new implementation
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
