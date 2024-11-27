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
        uint256 deployerPrivateKey = vm.envUint("NEW_TIMELOCK_DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        currentTimelock = TimelockController(payable(Ownable2StepUpgradeable(Predeploys.Staking).owner()));

        address protocolAdmin = vm.envAddress("ADMIN_ADDRESS");

        if (protocolAdmin == address(0)) {
            protocolAdmin = vm.envAddress("ADMIN_ADDRESS");
        }
        require(protocolAdmin != address(0), "protocolAdmin not set");

        address[] memory timelockExecutors = vm.envAddress("TIMELOCK_EXECUTOR_ADDRESSES", ",");
        if (timelockExecutors.length == 0) {
            console2.log("Using public timelock executions");
        }

        address[] memory timelockCancellers = vm.envAddress("TIMELOCK_CANCELLER_ADDRESSES", ",");
        require(timelockCancellers.length > 0, "No timelock cancellers set");
        for (uint256 i = 0; i < timelockCancellers.length; i++) {
            require(timelockCancellers[i] != address(0), "Zero address astimelock guardian");
        }
        address[] memory proposers = new address[](1);
        proposers[0] = protocolAdmin;

        bytes memory creationCode = abi.encodePacked(
            type(TimelockController).creationCode,
            abi.encode(5 minutes, proposers, timelockExecutors, deployer)
        );
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER");
        address newTimelockAddress = Create3(Predeploys.Create3).deploy(salt, creationCode);
        newTimelock = TimelockController(payable(newTimelockAddress));

        for (uint256 i = 0; i < timelockCancellers.length; i++) {
            TimelockController(payable(newTimelock)).grantRole(
                TimelockController(payable(newTimelock)).CANCELLER_ROLE(),
                timelockCancellers[i]
            );
        }

        TimelockController(payable(newTimelock)).grantRole(
            TimelockController(payable(newTimelock)).DEFAULT_ADMIN_ROLE(),
            protocolAdmin
        );

        TimelockController(payable(newTimelock)).renounceRole(
            TimelockController(payable(newTimelock)).DEFAULT_ADMIN_ROLE(),
            deployer
        );

        console2.log("TimelockController deployed at:", address(newTimelock));


        vm.stopBroadcast();
    }
}
