/* solhint-disable no-console */
// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { IPTokenStaking } from "src/protocol/IPTokenStaking.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Create3 } from "src/deploy/Create3.sol";

/**
 * @title DeployNewIPTokenStakingImpl
 * @notice Deploys a new implementation of IPTokenStaking contract to be used for upgrading
 * @dev This script only deploys the implementation contract, it does not perform the upgrade
 */
contract DeployNewIPTokenStaking_V1_0_1 is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        Create3 create3 = Create3(Predeploys.Create3);

        // Generate creation code for IPTokenStaking
        bytes memory creationCode = abi.encodePacked(
            type(IPTokenStaking).creationCode,
            abi.encode(1 ether, 256) // Constructor args: defaultMinFee (1 IP), maxDataLength
        );

        bytes32 salt = keccak256(abi.encodePacked("IPTokenStaking_Implementation_v1_0_1"));

        // Deploy using Create3
        address newImplementation = create3.deploy(salt, creationCode);
        if (create3.getDeployed(deployer, salt) != newImplementation) {
            revert("Deployment failed");
        }

        vm.stopBroadcast();

        console2.log("New IPTokenStaking implementation deployed at:", newImplementation);
    }
}
