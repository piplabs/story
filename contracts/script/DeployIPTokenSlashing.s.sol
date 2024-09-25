// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { IPTokenSlashing } from "../src/protocol/IPTokenSlashing.sol";

/**
 * @title DeployIPTokenSlashing
 * @dev A script to deploy IPTokenSlashing for Illiad
 */
contract DeployIPTokenSlashing is Script {
    // To run the script:
    // - Dry run
    // forge script script/DeployIPTokenSlashing.s.sol --fork-url <fork-url>
    //
    // - Deploy (OK for devnet)
    // forge script script/DeployIPTokenSlashing.s.sol --fork-url <fork-url> --broadcast
    // 
    // - Deploy and Verify (for testnet)
    // forge script script/DeployIPTokenSlashing.s.sol --fork-url https://testnet.storyrpc.io --broadcast --verify --verifier blockscout --verifier-url https://testnet.storyscan.xyz/api\?
    function run() public {
        // Read env for admin address
        address protocolAccessManagerAddr = vm.envAddress("ADMIN_ADDRESS");
        require(protocolAccessManagerAddr != address(0), "address not set");
        // Read env for deployer private key
        uint256 deployerKey = vm.envUint("IPTOKENSTAKING_DEPLOYER_KEY");

        vm.startBroadcast(deployerKey);

        address ipTokenStaking = 0xCCcCcC0000000000000000000000000000000001;

        address impl = address(new IPTokenSlashing(ipTokenStaking));
        IPTokenSlashing ipTokenSlashing = IPTokenSlashing(
            address(
                new ERC1967Proxy(
                    impl,
                    abi.encodeCall(
                        IPTokenSlashing.initialize, (
                            protocolAccessManagerAddr,
                            1 ether // unjailFee
                        )
                    )
                )
            )
        );

        vm.stopBroadcast();

        console2.log("IP_TOKEN_STAKING", address(ipTokenSlashing.IP_TOKEN_STAKING()));
        console2.log("IPTokenSlashing deployed at:", address(ipTokenSlashing));

    }
}
