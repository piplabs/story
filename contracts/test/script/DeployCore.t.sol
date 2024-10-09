// // SPDX-License-Identifier: BUSL-1.1
// pragma solidity ^0.8.23;
// /* solhint-disable no-console */
// /* solhint-disable max-line-length */
// /// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
// /// flag "Hex High Entropy String" in CI run detect-secrets

// import { Test } from "forge-std/Test.sol";

// import { DeployCore } from "../../script/DeployCore.s.sol";

// contract DeployCoreTest is Test {
//     DeployCore private deployCore;

//     function setUp() public {
//         deployCore = new DeployCore();
//     }

//     function testDeployDeployCore_run() public {
//         // Network shall not deploy the IPTokenStaking contract if IPTOKENSTAKING_DEPLOYER_KEY not set.
//         vm.chainId(1513);
//         // solhint-disable
//         vm.expectRevert('vm.envUint: environment variable "IPTOKENSTAKING_DEPLOYER_KEY" not found');
//         deployCore.run();

//         // Network shall deploy the IPTokenStaking contract.
//         vm.setEnv("IPTOKENSTAKING_DEPLOYER_KEY", "0x123456789abcdef");
//         deployCore.run();
//     }
// }
