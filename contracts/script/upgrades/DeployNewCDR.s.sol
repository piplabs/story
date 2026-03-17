// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { CDR } from "../../src/protocol/CDR.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";

/**
 * @title DeployNewCDR
 * @notice Deploys a new implementation of CDR contract to be used for upgrading the predeploy at
 *         0xCcCcCC0000000000000000000000000000000005
 * @dev This script only deploys the implementation contract, it does not perform the upgrade.
 *      The actual proxy upgrade must be executed through governance (timelock).
 */
contract DeployNewCDR is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        Create3 create3 = Create3(Predeploys.Create3);

        // CDR has no constructor args (only _disableInitializers in constructor)
        bytes memory creationCode = type(CDR).creationCode;

        bytes32 salt = keccak256(abi.encodePacked("CDR_Implementation_v1_0_0"));

        // Deploy using Create3 for deterministic address
        address newImplementation = create3.deploy(salt, creationCode);
        if (create3.getDeployed(deployer, salt) != newImplementation) {
            revert("Deployment failed");
        }

        vm.stopBroadcast();

        console2.log("New CDR implementation deployed at:", newImplementation);
        console2.log("CDR proxy address:", Predeploys.CDR);
    }
}
