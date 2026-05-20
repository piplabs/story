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
 * @notice Deploys TDXValidationHook (implementation + TransparentUpgradeableProxy) and logs the
 *         data required to whitelist it on DKG.
 * @dev After running, call DKG.whitelistEnclaveType through governance using the logged
 *      arguments, then approvePlatform for each (MRTD, RTMR0, RTMR1) tuple.
 *
 *      TDX_CODE_COMMITMENT must equal keccak256(RTMR2) for the running TDX kernel
 *      (RTMR2 measures initrd + cmdline); compute off-chain and paste before deployment.
 *
 *      Required env vars:
 *        DEPLOYER_PRIVATE_KEY      — deployer's private key
 *        TIMELOCK_ADDRESS          — owner of TDXValidationHook and ProxyAdmin
 *        AUTOMATA_VALIDATION_ADDR  — address of Automata DCAP attestation contract
 */
contract DeployTDXValidationHook is Script {
    /*//////////////////////////////////////////////////////////////////////////
    //                    Edit before running the script                      //
    //////////////////////////////////////////////////////////////////////////*/

    // keccak256(RTMR2) for the target TDX kernel. Update before deployment.
    bytes32 internal constant TDX_CODE_COMMITMENT =
        hex"0000000000000000000000000000000000000000000000000000000000000002";

    // DKG enclave type identifier for TDX; SGX genesis uses bytes32(uint256(1)).
    bytes32 internal constant TDX_ENCLAVE_TYPE = bytes32(uint256(2));

    /*//////////////////////////////////////////////////////////////////////////
    //                          Internal constants                            //
    //////////////////////////////////////////////////////////////////////////*/

    bytes32 internal constant TDX_IMPL_SALT = keccak256(abi.encodePacked("TDXValidationHook_Implementation_v1_0_0"));
    bytes32 internal constant TDX_PROXY_SALT = keccak256(abi.encodePacked("TDXValidationHook_Proxy_v1_0_0"));

    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address owner = vm.envAddress("TIMELOCK_ADDRESS");
        address automataValidationAddr = vm.envAddress("AUTOMATA_VALIDATION_ADDR");

        vm.startBroadcast(deployerPrivateKey);

        Create3 create3 = Create3(Predeploys.Create3);

        // Implementation.
        bytes memory implCreationCode = abi.encodePacked(
            type(TDXValidationHook).creationCode,
            abi.encode(Predeploys.DKG)
        );
        address tdxHookImpl = create3.deploy(TDX_IMPL_SALT, implCreationCode);

        // TransparentUpgradeableProxy.
        bytes memory initData = abi.encodeCall(TDXValidationHook.initialize, (owner, automataValidationAddr));
        bytes memory proxyCreationCode = abi.encodePacked(
            type(TransparentUpgradeableProxy).creationCode,
            abi.encode(tdxHookImpl, owner, initData)
        );
        address tdxHookProxy = create3.deploy(TDX_PROXY_SALT, proxyCreationCode);

        vm.stopBroadcast();

        console2.log("TDXValidationHook impl deployed at:", tdxHookImpl);
        console2.log("TDXValidationHook proxy deployed at:", tdxHookProxy);

        console2.log("---");
        console2.log("To whitelist this hook on DKG, call DKG.whitelistEnclaveType through governance with:");
        console2.log("  target:", Predeploys.DKG);
        console2.log("  enclaveType:");
        console2.logBytes32(TDX_ENCLAVE_TYPE);
        console2.log("  enclaveTypeData.codeCommitment (= keccak256(RTMR2)):");
        console2.logBytes32(TDX_CODE_COMMITMENT);
        console2.log("  enclaveTypeData.validationHookAddr:", tdxHookProxy);
        console2.log("  isWhitelisted: true");
        console2.log("---");
        console2.log("After whitelisting, call TDXValidationHook.approvePlatform(platformCommitment, label)");
        console2.log("for each approved (MRTD, RTMR0, RTMR1) tuple, where:");
        console2.log("  platformCommitment = keccak256(MRTD || RTMR0 || RTMR1)");
    }
}
