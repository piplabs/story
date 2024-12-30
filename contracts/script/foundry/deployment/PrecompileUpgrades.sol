// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { IPTokenStaking } from "src/protocol/IPTokenStaking.sol";
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { EIP1967Helper } from "script/utils/EIP1967Helper.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
// example usage:
// export UPGRADE_ADMIN_KEY=0x<your_admin_private_key>
// export EXECUTOR_KEY=0x<your_executor_private_key>

// gcp
// export IS_EXECUTE=false
// forge script src/script/TestPrecompileUpgrades.sol --rpc-url http://r1-d.odyssey-devnet.storyrpc.io:8545 --broadcast

// testnet
// export IS_EXECUTE=true
// forge script src/script/TestPrecompileUpgrades.sol --rpc-url https://odyssey.storyrpc.io/  #--broadcast
contract PrecompileUpgrades is Script {
    TimelockController internal timelock;
    address public newImpl = address(0x7038d1f55D789c25581ff72DA4C1EA543A61673B); // replace
    bytes32 public salt = keccak256(abi.encodePacked("StakingUpgrade"));

    function run() public {
        bool isExecution = vm.envBool("IS_EXECUTE");
        if (isExecution) {
            executeUpgrade();
        } else {
            scheduleUpgrade();
        }
    }

    function scheduleUpgrade() public {
        timelock = TimelockController(payable(0x4827c76bD61A223Ddd36D013c78F825eb0bb3Be3));
        uint256 upgradeKey = vm.envUint("UPGRADE_ADMIN_KEY");
        address upgrader = vm.addr(upgradeKey);
        
        vm.startBroadcast(upgradeKey);

        console2.log("=== Scheduling Upgrade ===");
        console2.log("Upgrader Address:", upgrader);
        console2.log("New Implementation Address:", newImpl);

        ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(Predeploys.Staking));
        console2.log("ProxyAdmin Address:", address(proxyAdmin));

        bytes memory data = abi.encodeWithSelector(
            proxyAdmin.upgradeAndCall.selector,
            ITransparentUpgradeableProxy(Predeploys.Staking),
            newImpl,
            ""
        );

        bytes32 operationId = keccak256(abi.encode(address(proxyAdmin), 0, data, bytes32(0), salt));
        console2.log("Operation ID:", vm.toString(operationId));

        uint256 minDelay = timelock.getMinDelay();
        require(minDelay > 0, "Invalid Min Delay");

        timelock.schedule(address(proxyAdmin), 0, data, bytes32(0), salt, minDelay);
        console2.log("Scheduled Upgrade with Min Delay:", minDelay);

        vm.stopBroadcast();
    }

    function executeUpgrade() public {
        timelock = TimelockController(payable(0x4827c76bD61A223Ddd36D013c78F825eb0bb3Be3));
        uint256 executorKey = vm.envUint("EXECUTOR_KEY");
        address executor = vm.addr(executorKey);

        vm.startBroadcast(executorKey);

        console2.log("=== Executing Upgrade ===");
        console2.log("Executor Address:", executor);

        ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(Predeploys.Staking));
        bytes memory data = abi.encodeWithSelector(
            proxyAdmin.upgradeAndCall.selector,
            ITransparentUpgradeableProxy(Predeploys.Staking),
            newImpl,
            ""
        );

        bytes32 operationId = keccak256(abi.encode(address(proxyAdmin), 0, data, bytes32(0), salt));
        console2.log("Operation ID:", vm.toString(operationId));

        require(timelock.isOperationReady(operationId), "Operation not ready for execution.");

        timelock.execute(address(proxyAdmin), 0, data, bytes32(0), salt);
        console2.log("Upgrade Executed Successfully.");

        verifyUpgrade();
        vm.stopBroadcast();
    }

    function verifyUpgrade() internal view {
        console2.log("Verifying Upgrade...");
        uint256 minStake = IPTokenStaking(Predeploys.Staking).minStakeAmount();
        address implAddress = EIP1967Helper.getImplementation(Predeploys.Staking);
        console2.log("implAddress: ", implAddress);
        
        require(minStake == 1024 ether, "Min stake amount mismatch.");
        require(implAddress == newImpl, "Implementation address mismatch.");

        console2.log("Upgrade Verified Successfully!");
    }
}