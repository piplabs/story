// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { ChainIds } from "../utils/ChainIds.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { Script } from "forge-std/Script.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";
import { console2 } from "forge-std/console2.sol";

contract TransferTimelock is Script {
    TimelockController newTimelock;
    TimelockController currentTimelock;

    constructor() {}

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("ADMIN_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        currentTimelock = TimelockController(payable(Ownable2StepUpgradeable(Predeploys.Staking).owner()));

        address protocolAdmin = vm.envAddress("ADMIN_ADDRESS");

        require(protocolAdmin != address(0), "protocolAdmin not set");
        address timelockExecutor = vm.envAddress("TIMELOCK_EXECUTOR_ADDRESS");
        require(timelockExecutor != address(0), "executor not set");
        address timelockGuardian = vm.envAddress("TIMELOCK_GUARDIAN_ADDRESS");
        require(timelockGuardian != address(0), "canceller not set");
        address[] memory proposers = new address[](1);
        proposers[0] = protocolAdmin;
        address[] memory executors = new address[](1);
        executors[0] = timelockExecutor;


        bytes memory creationCode = abi.encodePacked(
            type(TimelockController).creationCode,
            abi.encode(currentTimelock.getMinDelay(), proposers, executors, deployer)
        );
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER");
        address newTimelockAddress = Create3(Predeploys.Create3).deploy(salt, creationCode);
        newTimelock = TimelockController(payable(newTimelockAddress));

        newTimelock.grantRole(
            newTimelock.CANCELLER_ROLE(),
            timelockGuardian
        );
        if (block.chainid == ChainIds.MAINNET) {
            newTimelock.renounceRole(
                newTimelock.DEFAULT_ADMIN_ROLE(),
                deployer
            );
        }

        console2.log("TimelockController deployed at:", address(newTimelock));


        vm.stopBroadcast();
    }
}
