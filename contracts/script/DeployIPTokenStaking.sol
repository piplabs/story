// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";

/**
 * @title DeployIPTokenStaking
 * @dev A script + utilities to deploy the IPTokenStaking contract
 */
contract DeployIPTokenStaking is Script {
    function run() public {
        // TODO: read env
        address protocolAccessManagerAddr = address(0x438d47Bc1e184b976fB9EB06Fa4EA27e1247674E);

        require(block.chainid == 1, "only mainnet deployment");
        require(protocolAccessManagerAddr != address(0), "address not set");

        uint256 deployerKey = vm.envUint("IPTOKENSTAKING_DEPLOYER_KEY");

        vm.startBroadcast(deployerKey);
        address impl = address(
            new IPTokenStaking(
                1 gwei, // stakingRounding
                1000, // defaultCommissionRate, 10%
                5000, // defaultMaxCommissionRate, 50%
                500 // defaultMaxCommissionChangeRate, 5%
            )
        );
        //1 ether, // minStakeAmount
        //1 ether, // minUnstakeAmount
        //1 ether, // minRedelegateAmount
        //7 days, // withdrawalAddressChangeInterval

        IPTokenStaking ipTokenStaking = IPTokenStaking(address(new ERC1967Proxy(impl, "")));
        vm.stopBroadcast();

        console2.log("IPTokenStaking deployed at:", address(ipTokenStaking));
    }
}
