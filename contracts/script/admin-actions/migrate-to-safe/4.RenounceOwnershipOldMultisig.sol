/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { console2 } from "forge-std/console2.sol";
import { Script } from "forge-std/Script.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { IAccessControl } from "@openzeppelin/contracts/access/IAccessControl.sol";
import { Create3 } from "src/deploy/Create3.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { JSONTxWriter } from "../../utils/JSONTxWriter.s.sol";

/// @title RenounceOldMultisigRoles
/// @notice A script to generate a batch of transactions for the old multisig to renounce its roles in the new timelock controller
contract RenounceOldMultisigRoles is JSONTxWriter {
    constructor() JSONTxWriter("renounce-old-multisig-roles") {}

    // Role constants from TimelockController
    bytes32 public constant DEFAULT_ADMIN_ROLE = 0x00;

    address public newTimelockAddress;
    TimelockController public newTimelock;

    function run() public {
        // Get the new timelock address
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        newTimelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        newTimelock = TimelockController(payable(newTimelockAddress));
        console2.log("New Timelock address:", newTimelockAddress);

        // Check that the timelock exists
        require(newTimelockAddress.code.length > 0, "Timelock not deployed");

        // Get addresses of old multisigs
        address oldMultisigAddress = vm.envAddress("OLD_MULTISIG_ADDRESS");
        address oldSecurityCouncilAddress = vm.envAddress("OLD_TIMELOCK_GUARDIAN_ADDRESS");
        console2.log("Old Multisig address:", oldMultisigAddress);
        console2.log("Old Security Council address:", oldSecurityCouncilAddress);

        // Add transaction for old multisig to renounce PROPOSER_ROLE
        _saveTx(
            TimelockOp.EXECUTE,
            oldMultisigAddress,
            newTimelockAddress,
            0,
            abi.encodeWithSelector(
                IAccessControl.renounceRole.selector,
                newTimelock.PROPOSER_ROLE(),
                oldMultisigAddress
            ),
            "Old multisig renounces PROPOSER_ROLE"
        );

        // Add transaction for old multisig to renounce EXECUTOR_ROLE
        _saveTx(
            TimelockOp.EXECUTE,
            oldMultisigAddress,
            newTimelockAddress,
            0,
            abi.encodeWithSelector(
                IAccessControl.renounceRole.selector,
                newTimelock.EXECUTOR_ROLE(),
                oldMultisigAddress
            ),
            "Old multisig renounces EXECUTOR_ROLE"
        );

        // Add transaction for old multisig to renounce DEFAULT_ADMIN_ROLE
        _saveTx(
            TimelockOp.EXECUTE,
            oldMultisigAddress,
            newTimelockAddress,
            0,
            abi.encodeWithSelector(
                IAccessControl.renounceRole.selector,
                DEFAULT_ADMIN_ROLE,
                oldMultisigAddress
            ),
            "Old multisig renounces DEFAULT_ADMIN_ROLE"
        );

        // Add transaction for old security council to renounce CANCELLER_ROLE
        _saveTx(
            TimelockOp.EXECUTE,
            oldSecurityCouncilAddress,
            newTimelockAddress,
            0,
            abi.encodeWithSelector(
                IAccessControl.renounceRole.selector,
                newTimelock.CANCELLER_ROLE(),
                oldSecurityCouncilAddress
            ),
            "Old security council renounces CANCELLER_ROLE"
        );

        // Write all JSON files
        _writeFiles();
        
        console2.log("Generated batch transaction for old multisig and security council to renounce roles in new timelock");
    }
}