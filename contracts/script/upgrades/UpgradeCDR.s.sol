// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
// solhint-disable-next-line max-line-length
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";

import { CDR } from "../../src/protocol/CDR.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";
import { EIP1967Helper } from "../utils/EIP1967Helper.sol";

/**
 * @title UpgradeCDR
 * @notice Prepares a Safe transaction that calls TimelockController.scheduleBatch
 *         to upgrade the CDR predeploy proxy to the v1.1.0 implementation and
 *         initialize the new maxBatchSize storage slot.
 *
 * @dev Why no reinitializer call on upgradeAndCall:
 *      CDR.initialize uses the `initializer` modifier which rejects re-execution
 *      once the proxy is already initialized (version counter > 0). The new
 *      maxBatchSize field is instead set via CDR.setMaxBatchSize in the same
 *      timelock batch, atomically alongside the implementation upgrade.
 *
 * Prerequisites:
 *   1. Run DeployNewCDR.s.sol to deploy the v1.1.0 implementation via Create3
 *   2. Set env vars (see below)
 *   3. Run this script to generate the Safe transaction calldata
 *   4. Submit calldata to the Safe multisig (proposer of the timelock)
 *
 * Required env vars:
 *   DEPLOYER_PRIVATE_KEY   — address used to derive the Create3 implementation address
 *   TIMELOCK_ADDRESS       — address of the TimelockController (= CDR owner)
 *   TIMELOCK_MIN_DELAY     — minimum delay (seconds); used for scheduleBatch
 *   CDR_MAX_BATCH_SIZE     — value for the new maxBatchSize field (e.g. 20)
 */
contract UpgradeCDR is Script {
    enum MODE {
        SCHEDULE,
        EXECUTE,
        CANCEL
    }

    MODE mode = MODE.EXECUTE;

    function run() external view {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        address timelockAddr = vm.envAddress("TIMELOCK_ADDRESS");
        uint256 minDelay = vm.envUint("TIMELOCK_MIN_DELAY");
        uint256 maxBatchSize = vm.envUint("CDR_MAX_BATCH_SIZE");

        Create3 create3 = Create3(Predeploys.Create3);

        address newCDRImpl = create3.getDeployed(deployer, keccak256(abi.encodePacked("CDR_Implementation_v1_1_0")));

        console2.log("New CDR implementation:", newCDRImpl);

        address cdrProxyAdmin = EIP1967Helper.getAdmin(Predeploys.CDR);
        console2.log("CDR ProxyAdmin:", cdrProxyAdmin);
        console2.log("CDR proxy:", Predeploys.CDR);
        console2.log("maxBatchSize to set:", maxBatchSize);

        // Batch:
        //   [0] upgradeAndCall — swap implementation, no reinit (see dev note above)
        //   [1] setMaxBatchSize — initialize the new storage slot atomically
        address[] memory targets = new address[](2);
        targets[0] = cdrProxyAdmin;
        targets[1] = Predeploys.CDR;

        uint256[] memory values = new uint256[](2);

        bytes[] memory payloads = new bytes[](2);
        payloads[0] = abi.encodeWithSelector(
            ProxyAdmin.upgradeAndCall.selector,
            ITransparentUpgradeableProxy(Predeploys.CDR),
            newCDRImpl,
            bytes("") // no reinitializer — initialize is already locked
        );
        payloads[1] = abi.encodeWithSelector(CDR.setMaxBatchSize.selector, maxBatchSize);

        string memory modeString;
        bytes memory data;

        if (mode == MODE.CANCEL) {
            revert("TODO");
        } else if (mode == MODE.SCHEDULE) {
            modeString = "Schedule";
            // scheduleBatch(address[], uint256[], bytes[], bytes32 predecessor, bytes32 salt, uint256 delay)
            data = abi.encodeCall(
                TimelockController.scheduleBatch,
                (targets, values, payloads, bytes32(0), bytes32(0), minDelay)
            );
        } else {
            modeString = "Execute";
            // executeBatch(address[], uint256[], bytes[], bytes32 predecessor, bytes32 salt) — no delay
            data = abi.encodeCall(TimelockController.executeBatch, (targets, values, payloads, bytes32(0), bytes32(0)));
        }

        console2.log("=== Safe Transaction Parameters ===");
        console2.log("To (TimelockController):", timelockAddr);
        console2.log("Mode:", modeString);
        console2.log("Value: 0");
        console2.log("Data (hex):");
        console2.logBytes(data);
    }
}
