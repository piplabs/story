// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { TimelockedOperations } from "../utils/TimelockedOperations.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { ChainIds } from "../utils/ChainIds.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";

contract TransferTimelock is TimelockedOperations {
    constructor() TimelockedOperations(
        "transfer-timelock",
        UpgradeModes.SCHEDULE,
        Output.BATCH_TX_EXECUTION
    ) {}

    TimelockController newTimelock;
    TimelockController currentTimelock;

    function run() public override {
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

        newTimelock = new TimelockController(currentTimelock.getMinDelay(), proposers, executors, deployer);


        // Add your timelock transfer logic here
        // Example: timelock.transferRole(ROLE, newAddress);
        if (block.chainid != ChainIds.MAINNET) {
            vm.stopBroadcast();
        }
    }

    function _scheduleActions() internal virtual override {}

    function _executeActions() internal virtual override {}

    function _cancelActions() internal virtual override {}
}
