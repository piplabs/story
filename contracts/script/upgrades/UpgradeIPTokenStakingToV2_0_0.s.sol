// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { IPTokenStaking } from "../../src/upgrades/IPTokenStaking.sol";
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { EIP1967Helper } from "../utils/EIP1967Helper.sol";

/**
 * @title UpgradeIPTokenStakingToV2_0_0
 * @notice Schedules or executes the IPTokenStaking proxy upgrade to V2.0.0 (with pausing).
 * @dev Run twice: first with UPGRADE_ACTION=schedule, then after timelock minDelay with UPGRADE_ACTION=execute.
 *
 * Env:
 *   TIMELOCK_CONTROLLER_ADDRESS   - TimelockController contract address.
 *   TIMELOCK_PROPOSER_PRIVATE_KEY - Private key of timelock proposer (for schedule).
 *   TIMELOCK_EXECUTOR_PRIVATE_KEY - Private key of timelock executor (for execute).
 *   NEW_IMPLEMENTATION_ADDRESS    - Address of the new IPTokenStaking implementation (e.g. from DeployNewIPTokenStaking_V2_0_0).
 *   UPGRADE_ACTION                - 1 = schedule (uses proposer key), 2 = execute (uses executor key).
 *
 * Query minDelay: cast call $TIMELOCK_CONTROLLER_ADDRESS "getMinDelay()(uint256)" --rpc-url $RPC_URL
 *
 * Example:
 *   # Step 1: Schedule upgrade (proposer)
 *   UPGRADE_ACTION=1 NEW_IMPLEMENTATION_ADDRESS=0x... forge script ... --rpc-url $RPC_URL --broadcast -vvv
 *
 *   # Step 2: After timelock delay, execute (executor)
 *   UPGRADE_ACTION=2 NEW_IMPLEMENTATION_ADDRESS=0x... forge script ... --rpc-url $RPC_URL --broadcast -vvv
 */
contract UpgradeIPTokenStakingToV2_0_0 is Script {
    address internal constant STAKING_PROXY = Predeploys.Staking;
    address internal constant DEFAULT_SAFE_GOVERNANCE = 0xDc114b7aF2E3b298E212c88BefF2fECd849F9026;
    address internal constant DEFAULT_SECURITY_COUNCIL = 0xe4E3C9D65eEC4175742Eaa0917423C11Aa86b1AE;

    function run() external {
        address timelockAddress = vm.envAddress("TIMELOCK_CONTROLLER_ADDRESS");
        address newImplementation = vm.envAddress("NEW_IMPLEMENTATION_ADDRESS");
        uint256 action = vm.envUint("UPGRADE_ACTION");

        require(timelockAddress != address(0), "TIMELOCK_CONTROLLER_ADDRESS required");
        require(newImplementation != address(0), "NEW_IMPLEMENTATION_ADDRESS required");
        require(action == 1 || action == 2, "UPGRADE_ACTION must be 1 (schedule) or 2 (execute)");

        uint256 signerKey = action == 1
            ? vm.envUint("TIMELOCK_PROPOSER_PRIVATE_KEY")
            : vm.envUint("TIMELOCK_EXECUTOR_PRIVATE_KEY");
        require(signerKey != 0, action == 1 ? "TIMELOCK_PROPOSER_PRIVATE_KEY required" : "TIMELOCK_EXECUTOR_PRIVATE_KEY required");

        address proxyAdminAddr = EIP1967Helper.getAdmin(STAKING_PROXY);
        ProxyAdmin proxyAdmin = ProxyAdmin(proxyAdminAddr);
        TimelockController timelock = TimelockController(payable(timelockAddress));

        bytes memory upgradeCalldata = abi.encodeWithSelector(
            ProxyAdmin.upgradeAndCall.selector,
            ITransparentUpgradeableProxy(STAKING_PROXY),
            IPTokenStaking(newImplementation),
            abi.encodeWithSelector(
                IPTokenStaking.initializeV2.selector,
                DEFAULT_SAFE_GOVERNANCE,
                DEFAULT_SECURITY_COUNCIL
            )
        );

        if (action == 1) {
            vm.startBroadcast(signerKey);
            timelock.schedule(
                address(proxyAdmin),
                0,
                upgradeCalldata,
                bytes32(0),
                bytes32(0),
                timelock.getMinDelay()
            );
            console2.log("Upgrade scheduled. Execute after minDelay (seconds):", timelock.getMinDelay());
            vm.stopBroadcast();
        } else {
            vm.startBroadcast(signerKey);
            timelock.execute(
                address(proxyAdmin),
                0,
                upgradeCalldata,
                bytes32(0),
                bytes32(0)
            );
            console2.log("Upgrade executed. Staking proxy now at implementation:", newImplementation);
            vm.stopBroadcast();
        }
    }
}
