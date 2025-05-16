/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { console2 } from "forge-std/console2.sol";
import { Script } from "forge-std/Script.sol";

import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";

import { Create3 } from "src/deploy/Create3.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";

/// @title DeployNewTimelock
/// @notice Deploy a new TimelockController governed by the new multisigs
contract DeployNewTimelock is Script {

    function run() public virtual {
        require(!isTimelockDeployed(), "TimelockController already deployed");
        deployTimelock();

    }

    /// @notice Check if the TimelockController is deployed
    /// @return True if the TimelockController is deployed, false otherwise
    function isTimelockDeployed() internal view returns (bool) {
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        address timelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);

        if (timelockAddress.code.length == 0) {
            // Not deployed
            return false;
        }
        TimelockController timelock = TimelockController(payable(timelockAddress));

        // Check if timelock has a minDelay and assigned proposer role as proof of deployment
        if (timelock.getMinDelay() == 0 || !timelock.hasRole(
            timelock.PROPOSER_ROLE(),
            vm.envAddress("SAFE_TIMELOCK_PROPOSER")
        )) {
            revert("wrong timelock controller deployment");
        }
        return true;
    }

    /// @notice Deploy a new TimelockController deterministically
    function deployTimelock() internal {
        uint256 deployerPrivateKey = vm.envUint("NEW_TIMELOCK_DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        address timelockProposer = vm.envAddress("SAFE_TIMELOCK_PROPOSER");
        require(timelockProposer != address(0), "safe admin address not set");
        console2.log("timelockProposer", timelockProposer);

        // Old timelock proposer during migration
        address oldTimelockProposer = vm.envAddress("OLD_TIMELOCK_PROPOSER");
        require(oldTimelockProposer != address(0), "old protocol admin address not set");
        console2.log("oldTimelockProposer", oldTimelockProposer);

        address[] memory proposers = new address[](2);
        proposers[0] = timelockProposer;
        proposers[1] = oldTimelockProposer;
        console2.log("proposers", proposers[0], proposers[1]);


        // Executor can be address(0), public execution
        address timelockExecutor = vm.envAddress("SAFE_TIMELOCK_EXECUTOR");
        address[] memory timelockExecutors = new address[](2);
        timelockExecutors[0] = timelockExecutor;
        timelockExecutors[1] = oldTimelockProposer; // Old multisig during migration
        console2.log("timelockExecutor", timelockExecutor);

        address timelockGuardian = vm.envAddress("SAFE_TIMELOCK_GUARDIAN");
        require(timelockGuardian != address(0), "Zero address as timelock guardian");

        address oldTimelockGuardian = vm.envAddress("OLD_TIMELOCK_GUARDIAN");
        require(oldTimelockGuardian != address(0), "old timelock guardian address not set");
        console2.log("oldTimelockGuardian", oldTimelockGuardian);

        address[] memory timelockCancellers = new address[](2);
        timelockCancellers[0] = timelockGuardian;
        timelockCancellers[1] = oldTimelockGuardian;
        console2.log("timelockGuardian", timelockGuardian);

        bytes memory creationCode = abi.encodePacked(
            type(TimelockController).creationCode,
            abi.encode(
                vm.envUint("MIN_DELAY"),
                proposers,
                timelockExecutors,
                deployer // Warning: root admin. Must renounce by end of script
            )
        );

        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        
        address newTimelockAddress = Create3(Predeploys.Create3).deployDeterministic(creationCode, salt);
        console2.log("Deployed TimelockController at address:", newTimelockAddress);
        console2.log("Temporary root admin:", deployer);
        console2.log("Proposers", proposers[0], proposers[1]);
        console2.log("timelockExecutors", timelockExecutors[0], timelockExecutors[1]);
        TimelockController newTimelock = TimelockController(payable(newTimelockAddress));

        for (uint256 i = 0; i < timelockCancellers.length; i++) {
            newTimelock.grantRole(
                newTimelock.CANCELLER_ROLE(),
                timelockCancellers[i]
            );
        }
        console2.log("timelockCancellers", timelockCancellers[0], timelockCancellers[1]);

        // WARNING: Revoke this role after migration
        newTimelock.grantRole(
            newTimelock.DEFAULT_ADMIN_ROLE(),
            oldTimelockProposer
        );
        console2.log("Granted DEFAULT_ADMIN_ROLE to oldTimelockProposer", oldTimelockProposer);
        // WARNING: Revoke this role after migration
        newTimelock.grantRole(
            newTimelock.DEFAULT_ADMIN_ROLE(),
            timelockProposer
        );
        console2.log("Granted DEFAULT_ADMIN_ROLE to timelockProposer", timelockProposer);

        newTimelock.renounceRole(
            newTimelock.DEFAULT_ADMIN_ROLE(), // DEFAULT_ADMIN_ROLE
            deployer
        );
        console2.log("Renounced DEFAULT_ADMIN_ROLE from deployer", deployer);

        vm.stopBroadcast();
    }

}