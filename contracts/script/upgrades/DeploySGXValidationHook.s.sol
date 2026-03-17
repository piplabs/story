// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { SGXValidationHook } from "../../src/protocol/SGXValidationHook.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";

/**
 * @title DeploySGXValidationHook
 * @notice Deploys SGXValidationHook (implementation + TransparentUpgradeableProxy) and outputs
 *         the data required to whitelist it on DKG at 0xCcCcCC0000000000000000000000000000000004
 * @dev This script only deploys the SGXValidationHook contracts.
 *      After running, call DKG.whitelistEnclaveType through governance (timelock) using the
 *      logged arguments below.
 *
 *      Required env vars:
 *        DEPLOYER_PRIVATE_KEY          — deployer's private key
 *        OWNER_ADDRESS                 — owner of SGXValidationHook and ProxyAdmin (typically timelock)
 *        AUTOMATA_VALIDATION_ADDR      — address of automata DCAP attestation contract
 *        TCB_EVALUATION_DATA_NUMBER    — TCB evaluation data number (uint)
 *
 *      Edit SGX_CODE_COMMITMENT constant below before running.
 */
contract DeploySGXValidationHook is Script {
    /*//////////////////////////////////////////////////////////////////////////
    //                    Edit before running the script                      //
    //////////////////////////////////////////////////////////////////////////*/

    // SGX enclave code commitment — update this hex value before deployment
    bytes32 constant SGX_CODE_COMMITMENT = hex"0000000000000000000000000000000000000000000000000000000000000001";

    /*//////////////////////////////////////////////////////////////////////////
    //                          Internal constants                            //
    //////////////////////////////////////////////////////////////////////////*/

    bytes32 constant SGX_IMPL_SALT = keccak256(abi.encodePacked("SGXValidationHook_Implementation_v1_0_0"));
    bytes32 constant SGX_PROXY_SALT = keccak256(abi.encodePacked("SGXValidationHook_Proxy_v1_0_0"));

    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address owner = vm.envAddress("OWNER_ADDRESS");
        address automataValidationAddr = vm.envAddress("AUTOMATA_VALIDATION_ADDR");
        uint32 tcbEvaluationDataNumber = uint32(vm.envUint("TCB_EVALUATION_DATA_NUMBER"));

        vm.startBroadcast(deployerPrivateKey);

        Create3 create3 = Create3(Predeploys.Create3);

        // Deploy SGXValidationHook implementation via Create3
        bytes memory implCreationCode = abi.encodePacked(
            type(SGXValidationHook).creationCode,
            abi.encode(Predeploys.DKG)
        );
        address sgxHookImpl = create3.deploy(SGX_IMPL_SALT, implCreationCode);

        // Deploy TransparentUpgradeableProxy wrapping the implementation via Create3
        bytes memory initData = abi.encodeCall(
            SGXValidationHook.initialize,
            (owner, automataValidationAddr, tcbEvaluationDataNumber)
        );
        bytes memory proxyCreationCode = abi.encodePacked(
            type(TransparentUpgradeableProxy).creationCode,
            abi.encode(sgxHookImpl, owner, initData)
        );
        address sgxHookProxy = create3.deploy(SGX_PROXY_SALT, proxyCreationCode);

        vm.stopBroadcast();

        // Log deployment results
        console2.log("SGXValidationHook impl deployed at:", sgxHookImpl);
        console2.log("SGXValidationHook proxy deployed at:", sgxHookProxy);

        // Log whitelisting data for DKG governance action
        console2.log("---");
        console2.log("To whitelist this hook on DKG, call DKG.whitelistEnclaveType through governance with:");
        console2.log("  target:", Predeploys.DKG);
        console2.log("  enclaveType:", uint256(1));
        console2.log("  enclaveTypeData.codeCommitment:");
        console2.logBytes32(SGX_CODE_COMMITMENT);
        console2.log("  enclaveTypeData.validationHookAddr:", sgxHookProxy);
        console2.log("  isWhitelisted: true");
    }
}
