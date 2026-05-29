// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
// solhint-disable-next-line max-line-length
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";

import { SGXValidationHook } from "../../src/protocol/SGXValidationHook.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";
import { EIP1967Helper } from "../utils/EIP1967Helper.sol";

/**
 * @title UpgradeSGXValidationHook_V1_0_1
 * @notice Deploys the new SGXValidationHook implementation via Create3 and prepares a Safe
 *         transaction that drives a TimelockController operation to upgrade the already-deployed
 *         proxy on Aeneid to the new implementation. The new implementation calls the
 *         parameterless `verifyAndAttestOnChain` overload, so no additional state needs to be set
 *         on the proxy after the upgrade.
 * @dev Run with `--rpc-url $AENEID --broadcast` for the first run (deploys the implementation and
 *      logs the timelock calldata). Subsequent runs without `--broadcast` skip the deploy (Create3
 *      address is deterministic) and just regenerate the calldata — useful for switching between
 *      SCHEDULE and EXECUTE modes after submitting the schedule.
 *
 *      Required env vars:
 *        DEPLOYER_PRIVATE_KEY — deployer's private key (must be the same address used for the
 *                               Create3 prediction across runs)
 *        TIMELOCK_ADDRESS     — TimelockController address
 *        TIMELOCK_MIN_DELAY   — Timelock min delay (seconds), used for schedule only
 *        SGX_HOOK_PROXY       — already-deployed SGXValidationHook proxy address
 */
contract UpgradeSGXValidationHook_V1_0_1 is Script {
    enum MODE {
        SCHEDULE,
        EXECUTE,
        CANCEL
    }

    MODE mode = MODE.SCHEDULE;

    bytes32 constant IMPL_SALT = keccak256(abi.encodePacked("SGXValidationHook_Implementation_v1_0_1"));

    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        address timelockAddr = vm.envAddress("TIMELOCK_ADDRESS");
        uint256 minDelay = vm.envUint("TIMELOCK_MIN_DELAY");
        address sgxHookProxy = vm.envAddress("SGX_HOOK_PROXY");

        Create3 create3 = Create3(Predeploys.Create3);
        address newImpl = create3.getDeployed(deployer, IMPL_SALT);

        // Deploy the new implementation if it has not been deployed yet at the predicted address.
        // Idempotent across re-runs: subsequent invocations (e.g. when switching from SCHEDULE to
        // EXECUTE mode) skip this branch.
        if (newImpl.code.length == 0) {
            vm.startBroadcast(deployerPrivateKey);
            bytes memory creationCode = abi.encodePacked(
                type(SGXValidationHook).creationCode,
                abi.encode(Predeploys.DKG)
            );
            address deployed = create3.deploy(IMPL_SALT, creationCode);
            require(deployed == newImpl, "Create3 address mismatch");
            vm.stopBroadcast();
            console2.log("New SGXValidationHook implementation deployed at:", newImpl);
        } else {
            console2.log("SGXValidationHook implementation already deployed at:", newImpl);
        }

        // Read ProxyAdmin from the proxy's EIP-1967 admin slot
        address proxyAdmin = EIP1967Helper.getAdmin(sgxHookProxy);

        console2.log("SGXValidationHook proxy:", sgxHookProxy);
        console2.log("SGXValidationHook ProxyAdmin:", proxyAdmin);

        // Single timelock operation: pure proxy upgrade. initialize() must not be re-run (the v1.0.1
        // impl does not introduce a reinitializer), and no extra owner state needs to be set after
        // the upgrade — the new implementation calls verifyAndAttestOnChain(bytes) which lets the
        // Automata verifier resolve the standard TCB Evaluation Data Number on-chain.
        bytes memory upgradePayload = abi.encodeCall(
            ProxyAdmin.upgradeAndCall,
            (ITransparentUpgradeableProxy(sgxHookProxy), newImpl, bytes(""))
        );

        string memory modeString;
        bytes memory data;

        if (mode == MODE.CANCEL) {
            revert("TODO");
        } else if (mode == MODE.SCHEDULE) {
            modeString = "Schedule";
            data = abi.encodeCall(
                TimelockController.schedule,
                (proxyAdmin, 0, upgradePayload, bytes32(0), bytes32(0), minDelay)
            );
        } else {
            modeString = "Execute";
            data = abi.encodeCall(TimelockController.execute, (proxyAdmin, 0, upgradePayload, bytes32(0), bytes32(0)));
        }

        // Output Safe transaction parameters
        console2.log("=== Safe Transaction Parameters ===");
        console2.log("To (TimelockController):", timelockAddr);
        console2.log("Mode:", modeString);
        console2.log("Value: 0");
        console2.log("Data (hex):");
        console2.logBytes(data);
    }
}
