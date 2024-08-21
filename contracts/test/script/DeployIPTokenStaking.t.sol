// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { Test } from "forge-std/Test.sol";

import { DeployIPTokenStaking } from "../../script/DeployIPTokenStaking.sol";

contract DeployIPTokenStakingTest is Test {
    DeployIPTokenStaking private deployIPTokenStaking;

    function setUp() public {
        deployIPTokenStaking = new DeployIPTokenStaking();
    }

    function testDeployIPTokenStaking_run() public {
        // Network shall not deploy the IPTokenStaking contract if not mainnet.
        vm.expectRevert("only mainnet deployment");
        deployIPTokenStaking.run();

        // Network shall not deploy the IPTokenStaking contract if IPTOKENSTAKING_DEPLOYER_KEY not set.
        vm.chainId(1);
        // solhint-disable
        vm.expectRevert("vm.envUint: environment variable \"IPTOKENSTAKING_DEPLOYER_KEY\" not found");
        deployIPTokenStaking.run();

        // Network shall deploy the IPTokenStaking contract.
        vm.setEnv("IPTOKENSTAKING_DEPLOYER_KEY", "0x123456789abcdef");
        deployIPTokenStaking.run();
    }
}
