// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { IIPTokenStaking } from "src/interfaces/IIPTokenStaking.sol";
import { IPTokenStaking } from "src/protocol/IPTokenStaking.sol";
import { Create3 } from "src/deploy/Create3.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { EIP1967Helper } from "script/utils/EIP1967Helper.sol";
import { DeployNewTimelock } from "script/admin-actions/migrate-to-safe/1.DeployNewTimelock.s.sol";
import { BaseTransferOwnershipProxyAdmin } from "script/admin-actions/migrate-to-safe/BaseTransferOwnershipProxyAdmin.s.sol";
import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";
import { TransferOwnershipIPTokenStaking } from "script/admin-actions/migrate-to-safe/3.5.TransferOwnershipIPTokenStaking.s.sol";
import { ReceiveOwnershipIPTokenStaking } from "script/admin-actions/migrate-to-safe/3.6.ReceiveOwnershipIPTokenStaking.s.sol";

/// @dev Minimal Initializable mock for proxy implementation
contract MockInitializable {
    bool public initialized;
    function initialize() public {
        initialized = true;
    }
}

contract TestTransferOwnershipProxyAdmin is BaseTransferOwnershipProxyAdmin {

    address public predeploy;

    constructor(address _predeploy, string memory _action) BaseTransferOwnershipProxyAdmin(_action, uint160(0), uint160(1)) {
        predeploy = _predeploy;
    }

    function _generate() internal virtual override {
        require(address(newTimelock) != address(0), "Timelock not deployed");
        require(address(newTimelock) != address(currentTimelock()), "Timelock already set");

        address proxyAdmin = EIP1967Helper.getAdmin(predeploy);
        bytes4 selector = Ownable.transferOwnership.selector;
        _generateAction(from, proxyAdmin, 0, abi.encodeWithSelector(selector, address(newTimelock)), bytes32(0), bytes32(0), minDelay);
    }
}

contract TestTransferIPTokenStakingOwnership is TransferOwnershipIPTokenStaking {
    constructor(address _predeploy, string memory _action) TransferOwnershipIPTokenStaking() {
        action = _action;
        ipTokenStakingProxy = _predeploy;
    }
}

contract TestReceiveOwnershipIPTokenStaking is ReceiveOwnershipIPTokenStaking {
    constructor(address _predeploy, string memory _action) ReceiveOwnershipIPTokenStaking() {
        action = _action;
        ipTokenStakingProxy = _predeploy;
    }
}


/**
 * @title DeployTestAndGenJSONs
 * @dev A script to deploy IPTokenStaking implementation and proxy
 */
contract DeployTestAndGenJSONs is Script {
    /// @dev Struct to store all deployed contract addresses
    struct DeployedAddresses {
        address ipTokenStakingImpl;
        address ipTokenStakingProxy;
        address ipTokenStakingProxyAdmin;
    }

    DeployedAddresses public deployed;
    address sourceTimelockController;
    address timelockControllerDestination;

    function estimateTimelockAddress(string memory _unhashedSalt) public view returns (address) {
        bytes32 salt = keccak256(bytes(_unhashedSalt));
        return Create3(Predeploys.Create3).predictDeterministicAddress(salt);
    }

    /// @notice main script method
    function run() public {
        require(block.chainid == 1315, "Only on Aeneid");
        console2.log("--- running");
        // Load deployer private key from .env
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        console2.log("---- loaded private key");
        vm.startBroadcast(deployerPrivateKey);
        address deployer = vm.addr(deployerPrivateKey);
        console2.log("---- deployer", deployer);

        // Both previously deployed timelocks are used for testing
        sourceTimelockController = estimateTimelockAddress("STORY_TIMELOCK_TEST_AENEID_SOURCE");
        timelockControllerDestination = estimateTimelockAddress("STORY_TIMELOCK_CONTROLLER_SAFE");

        console2.log("---- Test source timelock deployed at", sourceTimelockController);

        // Deploy implementation using new
        deployed.ipTokenStakingImpl = address(new IPTokenStaking(1 ether, 256));
        console2.log("deployed implementations");

        // Deploy a Create3 factory
        Create3 create3 = Create3(Predeploys.Create3);
        console2.log("loaded create3");

        // Deploy proxy for IPTokenStaking using Create3
        bytes32 saltProxy = keccak256("TestPredeployProxy");
        
        address impl = deployed.ipTokenStakingImpl;
        IIPTokenStaking.InitializerArgs memory args = IIPTokenStaking.InitializerArgs({
            owner: sourceTimelockController,
            minStakeAmount: 1024 ether,
            minUnstakeAmount: 1024 ether,
            minCommissionRate: 5_00,
            fee: 1 ether
        });
        bytes memory data = abi.encodeWithSelector(IPTokenStaking.initialize.selector, args);

        // Estimate address of proxy first
        address proxy = create3.predictDeterministicAddress(saltProxy);
        console2.log("estimated proxy address", proxy);

        proxy = create3.deployDeterministic(
            abi.encodePacked(
                type(TransparentUpgradeableProxy).creationCode,
                abi.encode(impl, address(sourceTimelockController), data)
            ),
            saltProxy
        );
        deployed.ipTokenStakingProxy = proxy;
        address proxyAdmin = EIP1967Helper.getAdmin(proxy);
        deployed.ipTokenStakingProxyAdmin = proxyAdmin;
        vm.stopBroadcast();

        // Log all named addresses
        console2.log("IPTokenStaking Impl:", deployed.ipTokenStakingImpl);
        console2.log("IPTokenStaking Proxy:", deployed.ipTokenStakingProxy);
        console2.log("IPTokenStaking ProxyAdmin:", deployed.ipTokenStakingProxyAdmin);

        // Log owners
        console2.log("IPTokenStaking Proxy Owner:", IPTokenStaking(deployed.ipTokenStakingProxy).owner());
        console2.log("IPTokenStaking ProxyAdmin Owner:", IPTokenStaking(deployed.ipTokenStakingProxyAdmin).owner());

        // Generate json files with the timelocked operations to transfer the ownership of proxy admin and IPTokenStaking to new timelock

        TestTransferOwnershipProxyAdmin transferOwnershipProxyAdmin = new TestTransferOwnershipProxyAdmin(deployed.ipTokenStakingProxy, "test_1.-proxy-ownership");
        transferOwnershipProxyAdmin.run();

        TestTransferIPTokenStakingOwnership transferIPTokenStakingOwnership = new TestTransferIPTokenStakingOwnership(deployed.ipTokenStakingProxy, "test_2.-ip-token-staking-ownership");
        transferIPTokenStakingOwnership.run();

        TestReceiveOwnershipIPTokenStaking receiveOwnershipIPTokenStaking = new TestReceiveOwnershipIPTokenStaking(deployed.ipTokenStakingProxy, "test_3.-receive-ip-token-staking-ownership");
        receiveOwnershipIPTokenStaking.run();
    }
} 