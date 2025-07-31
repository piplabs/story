// // SPDX-License-Identifier: GPL-3.0-only
// pragma solidity 0.8.23;
// /* solhint-disable no-console */
// /* solhint-disable max-line-length */

// import { Script } from "forge-std/Script.sol";
// import { console2 } from "forge-std/console2.sol";
// import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
// import { IIPTokenStaking } from "src/interfaces/IIPTokenStaking.sol";
// import { IPTokenStaking } from "src/protocol/IPTokenStaking.sol";
// import { Create3 } from "src/deploy/Create3.sol";
// import { Predeploys } from "src/libraries/Predeploys.sol";
// import { EIP1967Helper } from "script/utils/EIP1967Helper.sol";
// import { DeployNewTimelock } from "script/admin-actions/migrate-to-safe/1.DeployNewTimelock.s.sol";
// import { BaseTransferOwnershipProxyAdmin } from "script/admin-actions/migrate-to-safe/BaseTransferOwnershipProxyAdmin.s.sol";

// /// @dev Minimal Initializable mock for proxy implementation
// contract MockInitializable {
//     bool public initialized;
//     function initialize() public {
//         initialized = true;
//     }
// }

// contract TestTransferOwnershipProxyAdmin is BaseTransferOwnershipProxyAdmin {

//     address public predeploy;

//     constructor(address _predeploy) BaseTransferOwnershipProxyAdmin("test proxy ownership", uint160(0), uint160(1)) {
//         predeploy = _predeploy;
//     }

//     function _generate() internal virtual override {
//         require(address(newTimelock) != address(0), "Timelock not deployed");
//         require(address(newTimelock) != address(currentTimelock()), "Timelock already set");
//         console2.log("targetsLength", targetsLength);

//         address[] memory targets = new address[](targetsLength);

//         address proxyAdmin = EIP1967Helper.getAdmin(predeploy);
//         targets[0] = proxyAdmin;

//         bytes4 selector = Ownable.transferOwnership.selector;
//         bytes[] memory data = new bytes[](targetsLength);
//         for (uint160 i = 0; i < targetsLength; i++) {
//             data[i] = abi.encodeWithSelector(selector, address(newTimelock));
//         }
//         uint256[] memory values = new uint256[](targetsLength);

//         _generateBatchAction(from, targets, values, data, bytes32(0), bytes32(0), minDelay);
//     }
// }

// /**
//  * @title DeployAlloc
//  * @dev A script to deploy IPTokenStaking implementation and proxy
//  */
// contract DeployAlloc is Script {
//     /// @dev Struct to store all deployed contract addresses
//     struct DeployedAddresses {
//         address ipTokenStakingImpl;
//         address ipTokenStakingProxy;
//         address ipTokenStakingProxyAdmin;
//     }

//     DeployedAddresses public deployed;
//     address sourceTimelockController;
//     address timelockControllerDestination;

//     /// @notice main script method
//     function run() public {
//         console2.log("--- running");
//         // Load deployer private key from .env
//         uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
//         console2.log("---- loaded private key");
//         vm.startBroadcast(deployerPrivateKey);

//         // Deploy Test Timelock to simulate old one controlled by old multisig
//         DeployNewTimelock deployNewTimelock = new DeployNewTimelock();
//         deployNewTimelock.run("STORY_TIMELOCK_TEST_AENEID_SOURCE");
//         sourceTimelockController = deployNewTimelock.newTimelockAddress();
//         timelockControllerDestination = deployNewTimelock.estimateTimelockAddress("STORY_TIMELOCK_CONTROLLER_SAFE");

//         console2.log("---- Test source timelock deployed at", sourceTimelockController);

//         // Deploy implementation using new
//         deployed.ipTokenStakingImpl = address(new IPTokenStaking(1 ether, 256));
//         console2.log("deployed implementations");

//         // Deploy a Create3 factory
//         Create3 create3 = Create3(Predeploys.Create3);
//         console2.log("loaded create3");

//         // Deploy proxy for IPTokenStaking using Create3
//         bytes32 saltProxy = keccak256(abi.encodePacked(bytes("PredeployProxy"), uint256(0)));
//         address impl = deployed.ipTokenStakingImpl;
//         IIPTokenStaking.InitializerArgs memory args = IIPTokenStaking.InitializerArgs({
//             owner: sourceTimelockController,
//             minStakeAmount: 1024 ether,
//             minUnstakeAmount: 1024 ether,
//             minCommissionRate: 5_00,
//             fee: 1 ether
//         });
//         bytes memory data = abi.encodeWithSelector(IPTokenStaking.initialize.selector, args);
//         address proxy = create3.deployDeterministic(
//             abi.encodePacked(
//                 type(TransparentUpgradeableProxy).creationCode,
//                 abi.encode(impl, timelockController, data)
//             ),
//             saltProxy
//         );
//         deployed.ipTokenStakingProxy = proxy;
//         address proxyAdmin = EIP1967Helper.getAdmin(proxy);
//         deployed.ipTokenStakingProxyAdmin = proxyAdmin;
//         vm.stopBroadcast();

//         // Log all named addresses
//         console2.log("IPTokenStaking Impl:", deployed.ipTokenStakingImpl);
//         console2.log("IPTokenStaking Proxy:", deployed.ipTokenStakingProxy);
//         console2.log("IPTokenStaking ProxyAdmin:", deployed.ipTokenStakingProxyAdmin);

//         // Log owners
//         console2.log("IPTokenStaking Proxy Owner:", IPTokenStaking(deployed.ipTokenStakingProxy).owner());
//         console2.log("IPTokenStaking ProxyAdmin Owner:", IPTokenStaking(deployed.ipTokenStakingProxyAdmin).owner());

//         // Generate json file with the timelocked operations to transfer the ownership of proxy admin to new timelock
//         TestTransferOwnershipProxyAdmin transferOwnershipProxyAdmin = new TestTransferOwnershipProxyAdmin(deployed.ipTokenStakingProxy);
//         transferOwnershipProxyAdmin.run();

                

//     }
// } 