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
        if (!isTimelockDeployed()) {

            deployTimelock();
        } else {
            console2.log("TimelockController already deployed");
        }

        vm.stopBroadcast();
    }

    /// @notice Check if the TimelockController is deployed
    /// @return True if the TimelockController is deployed, false otherwise
    function isTimelockDeployed() internal view returns (bool) {
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        address timelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        // Return false if there is no code at the predicted timelock address
        if (timelockAddress.code.length > 0) {
            return false;
        }
        TimelockController timelock = TimelockController(payable(timelockAddress));
        // Check if timelock has a minDelay and assigned proposer role as proof of deployment
        return timelock.getMinDelay() > 0 && timelock.hasRole(
            timelock.PROPOSER_ROLE(),
            vm.envAddress("SAFE_ADMIN_ADDRESS")
        );
    }

    /// @notice Deploy a new TimelockController deterministically
    function deployTimelock() internal {
        uint256 deployerPrivateKey = vm.envUint("NEW_TIMELOCK_DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        address protocolAdmin = vm.envAddress("SAFE_ADMIN_ADDRESS");
        require(protocolAdmin != address(0), "safe admin address not set");
        console2.log("protocolAdmin", protocolAdmin);

        // Executor can be address(0), public execution
        address timelockExecutor = vm.envAddress("SAFE_TIMELOCK_EXECUTOR_ADDRESS");
        address[] memory timelockExecutors = new address[](1);
        timelockExecutors[0] = timelockExecutor;
        console2.log("timelockExecutor", timelockExecutor);

        address timelockGuardian = vm.envAddress("SAFE_TIMELOCK_GUARDIAN_ADDRESS");
        require(timelockGuardian != address(0), "Zero address as timelock guardian");
        address[] memory timelockCancellers = new address[](1);
        timelockCancellers[0] = timelockGuardian;
        console2.log("timelockGuardian", timelockGuardian);

        address[] memory proposers = new address[](1);
        proposers[0] = protocolAdmin;
        console2.log("proposers", proposers[0]);

        bytes memory creationCode = abi.encodePacked(
            type(TimelockController).creationCode,
            abi.encode(5 minutes, proposers, timelockExecutors, deployer)
        );

        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        address newTimelockAddress = Create3(Predeploys.Create3).deploy(salt, creationCode);
        TimelockController newTimelock = TimelockController(payable(newTimelockAddress));

        for (uint256 i = 0; i < timelockCancellers.length; i++) {
            newTimelock.grantRole(
                newTimelock.CANCELLER_ROLE(),
                timelockCancellers[i]
            );
        }

        // Uncomment to grant admin role to protocol admin
        // TimelockController(payable(newTimelock)).grantRole(
        //     TimelockController(payable(newTimelock)).DEFAULT_ADMIN_ROLE(),
        //     protocolAdmin
        // );

        newTimelock.renounceRole(
            newTimelock.DEFAULT_ADMIN_ROLE(),
            deployer
        );

        console2.log("TimelockController deployed at:", address(newTimelock));
    }

}