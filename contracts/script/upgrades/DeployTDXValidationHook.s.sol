// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { TDXValidationHook } from "../../src/protocol/TDXValidationHook.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";

/**
 * @title DeployTDXValidationHook
 * @notice Deploys TDXValidationHook (implementation + TransparentUpgradeableProxy) and outputs
 *         the data required to whitelist it on DKG at 0xCcCcCC0000000000000000000000000000000004
 * @dev This script only deploys the TDXValidationHook contracts.
 *      After running, call DKG.whitelistEnclaveType through governance (timelock) using the
 *      logged arguments below.
 *
 *      TDX_CODE_COMMITMENT must be set to keccak256(MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3)
 *      computed against the running TDX kernel's measured identity. The kernel surfaces
 *      MRTD and RTMR0..3 in its startup logs and via Identifier.GetSelfIdentity; the operator
 *      computes the keccak256 off-chain and pastes the result here before deployment.
 *
 *      Required env vars:
 *        DEPLOYER_PRIVATE_KEY          — deployer's private key
 *        OWNER_ADDRESS                 — owner of TDXValidationHook and ProxyAdmin (typically timelock)
 *        AUTOMATA_VALIDATION_ADDR      — address of automata DCAP attestation contract
 *        TCB_EVALUATION_DATA_NUMBER    — TCB evaluation data number (uint)
 *
 *      Edit TDX_CODE_COMMITMENT and TDX_ENCLAVE_TYPE constants below before running.
 */
contract DeployTDXValidationHook is Script {
    /*//////////////////////////////////////////////////////////////////////////
    //                    Edit before running the script                      //
    //////////////////////////////////////////////////////////////////////////*/

    // TDX enclave code commitment — keccak256(MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3).
    // Update this hex value before deployment. Must match the digest the running
    // TDX kernel computes for its own identity.
    bytes32 constant TDX_CODE_COMMITMENT = hex"0000000000000000000000000000000000000000000000000000000000000002";

    // DKG enclave type identifier for TDX. Distinct from the SGX genesis enclave type
    // (bytes32(uint256(1))) used by GenerateAlloc.setSGXValidationHook(). Operators
    // must whitelist this enclaveType through governance after deployment.
    bytes32 constant TDX_ENCLAVE_TYPE = bytes32(uint256(2));

    /*//////////////////////////////////////////////////////////////////////////
    //                          Internal constants                            //
    //////////////////////////////////////////////////////////////////////////*/

    bytes32 constant TDX_IMPL_SALT = keccak256(abi.encodePacked("TDXValidationHook_Implementation_v1_0_0"));
    bytes32 constant TDX_PROXY_SALT = keccak256(abi.encodePacked("TDXValidationHook_Proxy_v1_0_0"));

    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address owner = vm.envAddress("OWNER_ADDRESS");
        address automataValidationAddr = vm.envAddress("AUTOMATA_VALIDATION_ADDR");
        uint32 tcbEvaluationDataNumber = uint32(vm.envUint("TCB_EVALUATION_DATA_NUMBER"));

        vm.startBroadcast(deployerPrivateKey);

        Create3 create3 = Create3(Predeploys.Create3);

        // Deploy TDXValidationHook implementation via Create3
        bytes memory implCreationCode = abi.encodePacked(
            type(TDXValidationHook).creationCode,
            abi.encode(Predeploys.DKG)
        );
        address tdxHookImpl = create3.deploy(TDX_IMPL_SALT, implCreationCode);

        // Deploy TransparentUpgradeableProxy wrapping the implementation via Create3
        bytes memory initData = abi.encodeCall(
            TDXValidationHook.initialize,
            (owner, automataValidationAddr, tcbEvaluationDataNumber)
        );
        bytes memory proxyCreationCode = abi.encodePacked(
            type(TransparentUpgradeableProxy).creationCode,
            abi.encode(tdxHookImpl, owner, initData)
        );
        address tdxHookProxy = create3.deploy(TDX_PROXY_SALT, proxyCreationCode);

        vm.stopBroadcast();

        // Log deployment results
        console2.log("TDXValidationHook impl deployed at:", tdxHookImpl);
        console2.log("TDXValidationHook proxy deployed at:", tdxHookProxy);

        // Log whitelisting data for DKG governance action
        console2.log("---");
        console2.log("To whitelist this hook on DKG, call DKG.whitelistEnclaveType through governance with:");
        console2.log("  target:", Predeploys.DKG);
        console2.log("  enclaveType:");
        console2.logBytes32(TDX_ENCLAVE_TYPE);
        console2.log("  enclaveTypeData.codeCommitment:");
        console2.logBytes32(TDX_CODE_COMMITMENT);
        console2.log("  enclaveTypeData.validationHookAddr:", tdxHookProxy);
        console2.log("  isWhitelisted: true");
    }
}
