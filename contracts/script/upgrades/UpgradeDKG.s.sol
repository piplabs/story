// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
// solhint-disable-next-line max-line-length
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";

import { IDKG } from "../../src/interfaces/IDKG.sol";
import { CDR } from "../../src/protocol/CDR.sol";
import { DKG } from "../../src/protocol/DKG.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";
import { EIP1967Helper } from "../utils/EIP1967Helper.sol";

/**
 * @title UpgradeDKG
 * @notice Prepares a Safe transaction that calls TimelockController.scheduleBatch
 *         to upgrade both CDR and DKG predeploy proxies to new implementations.
 * @dev Prerequisites:
 *      1. Run DeployNewCDR.s.sol to deploy the new CDR implementation via Create3
 *      2. Run DeployNewDKG.s.sol to deploy the new DKG implementation via Create3
 *      3. Run this script to generate the Safe transaction calldata
 *      4. Submit the calldata to the Safe multisig (proposer of the timelock)
 */
contract UpgradeDKG is Script {
    enum MODE {
        SCHEDULE,
        EXECUTE,
        CANCEL
    }

    MODE mode = MODE.EXECUTE;

    bytes32 constant ENCLAVE_TYPE = 0x0000000000000000000000000000000000000000000000000000000000000001;
    bytes32 constant COMMITMENT = 0x6b2fb25e0084ad6ecbf6cfcefe09e2fa0fca2b84092f72c01f8fe98e9d7db5cd;
    bool constant WHITELISTED_VALUE = true;

    function run() external view {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        address timelockAddr = vm.envAddress("TIMELOCK_ADDRESS");
        uint256 minDelay = vm.envUint("TIMELOCK_MIN_DELAY");

        Create3 create3 = Create3(Predeploys.Create3);

        // Derive new implementation addresses from Create3
        address newDKGImpl = create3.getDeployed(
            deployer,
            keccak256(abi.encodePacked("DKG_Implementation_v1_0_0"))
        );
        address newCDRImpl = create3.getDeployed(
            deployer,
            keccak256(abi.encodePacked("CDR_Implementation_v1_0_0"))
        );

        console2.log("New DKG implementation:", newDKGImpl);
        console2.log("New CDR implementation:", newCDRImpl);

        // Get ProxyAdmin addresses
        address dkgProxyAdmin = EIP1967Helper.getAdmin(Predeploys.DKG);
        address cdrProxyAdmin = EIP1967Helper.getAdmin(Predeploys.CDR);

        console2.log("DKG ProxyAdmin:", dkgProxyAdmin);
        console2.log("CDR ProxyAdmin:", cdrProxyAdmin);

        // Build scheduleBatch arrays
        address[] memory targets = new address[](3);
        targets[0] = dkgProxyAdmin;
        targets[1] = cdrProxyAdmin;
        targets[2] = Predeploys.DKG;

        uint256[] memory values = new uint256[](3);

        bytes[] memory payloads = new bytes[](3);
        payloads[0] = _buildDKGUpgradePayload(newDKGImpl);
        payloads[1] = _buildCDRUpgradePayload(newCDRImpl);
        payloads[2] = _buildWhitelistPayload();

        bytes4 selector;
        string memory modeString;
        bytes memory data;

        if (mode == MODE.CANCEL) {
            revert("TODO");
        } else {
            if (mode == MODE.SCHEDULE) {
                selector = TimelockController.scheduleBatch.selector;
                modeString = "Schedule";
            } else {
                selector = TimelockController.executeBatch.selector;
                modeString = "Execute";
            }
            // Encode the full scheduleBatch call
            data = abi.encodeWithSelector(
                selector,
                targets,
                values,
                payloads,
                bytes32(0), // predecessor
                bytes32(0), // salt
                minDelay
            );
        }

        // Output Safe transaction parameters
        console2.log("=== Safe Transaction Parameters ===");
        console2.log("To (TimelockController):", timelockAddr);
        console2.log("Mode:", modeString);
        console2.log("Value: 0");
        console2.log("Data (hex):");
        console2.logBytes(data);
    }

    function _buildDKGUpgradePayload(address newImpl) internal view returns (bytes memory) {
        bytes memory initData = abi.encodeCall(
            DKG.initialize,
            (
                vm.envAddress("TIMELOCK_ADDRESS"), // DKG_OWNER
                vm.envUint("DKG_MIN_REQ_REGISTERED"),
                vm.envUint("DKG_MIN_REQ_FINALIZED"),
                vm.envUint("DKG_OPERATIONAL_THRESHOLD"),
                vm.envUint("DKG_FEE")
            )
        );

        return abi.encodeWithSelector(
            ProxyAdmin.upgradeAndCall.selector,
            ITransparentUpgradeableProxy(Predeploys.DKG),
            newImpl,
            initData
        );
    }

    function _buildCDRUpgradePayload(address newImpl) internal view returns (bytes memory) {
        bytes memory initData = abi.encodeCall(
            CDR.initialize,
            (
                vm.envAddress("TIMELOCK_ADDRESS"), // CDR_OWNER
                vm.envUint("CDR_BASE_FEE"),
                vm.envUint("CDR_WRITE_FEE"),
                vm.envUint("CDR_READ_FEE"),
                vm.envUint("CDR_ALLOCATE_FEE"),
                vm.envUint("CDR_MAX_ENCRYPTED_DATA_SIZE"),
                vm.envUint("CDR_MAX_ENCRYPTED_PARTIAL_SIZE")
            )
        );

        return abi.encodeWithSelector(
            ProxyAdmin.upgradeAndCall.selector,
            ITransparentUpgradeableProxy(Predeploys.CDR),
            newImpl,
            initData
        );
    }

    function _buildWhitelistPayload() internal view returns (bytes memory) {
        return abi.encodeWithSelector(
            DKG.whitelistEnclaveType.selector,
            ENCLAVE_TYPE,
            IDKG.EnclaveTypeData({
                codeCommitment: COMMITMENT,
                validationHookAddr: vm.envAddress("SGX_HOOK_PROXY")
            }),
            WHITELISTED_VALUE //true
        );
    }
}
