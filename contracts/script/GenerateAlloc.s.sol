// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { IIPTokenStaking } from "../src/interfaces/IIPTokenStaking.sol";
import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";
import { UpgradeEntrypoint } from "../src/protocol/UpgradeEntrypoint.sol";
import { UBIPool } from "../src/protocol/UBIPool.sol";

import { ChainIds } from "./utils/ChainIds.sol";
import { EIP1967Helper } from "./utils/EIP1967Helper.sol";
import { InitializableHelper } from "./utils/InitializableHelper.sol";
import { Predeploys } from "../src/libraries/Predeploys.sol";
import { Create3 } from "../src/deploy/Create3.sol";
import { ERC6551Registry } from "erc6551/ERC6551Registry.sol";
import { WIP } from "../src/token/WIP.sol";

/**
 * @title GenerateAlloc
 * @dev A script to generate the alloc section of EL genesis file
 * - Predeploys (See src/libraries/Predeploys.sol)
 * - Genesis $IP allocations (chain id dependent)
 * - If you want to allocate 10k test accounts with funds,
 * set this contract's property ALLOCATE_10K_TEST_ACCOUNTS to true
 * Run it by
 *  forge script script/GenerateAlloc.s.sol -vvvv --chain-id <CHAIN_ID>
 * Then, replace the contents of alloc field in EL genesis.json for the contents
 * of the generated json before starting the network.
 * This contract is also used by forge tests, to unify the process.
 */
contract GenerateAlloc is Script {
    /**
     * @notice Predeploy deployer address, used for each `new` call in this script
     */
    address internal deployer = 0xDDdDddDdDdddDDddDDddDDDDdDdDDdDDdDDDDDDd;

    // TimelockController
    address internal timelock;
    // Governance multi-sig
    address internal protocolAdmin;
    // Executor of scheduled operations
    address internal timelockExecutor;
    // Guardian of timelock
    address internal timelockGuardian;

    string internal dumpPath = getDumpPath();
    bool public saveState = true;
    // Optionally allocate 10k test accounts for devnets/testnets
    bool private constant ALLOCATE_10K_TEST_ACCOUNTS = false;
    // Optionally keep the timelock admin role for testnets
    bool private constant KEEP_TIMELOCK_ADMIN_ROLE = true;

    /// @notice this call should only be available from Test.sol, for speed
    function disableStateDump() external {
        require(block.chainid == ChainIds.FOUNDRY, "Only for local tests");
        saveState = false;
    }

    /// @dev this call should only be available from Test.sol
    function setAdminAddresses(address protocol, address executor, address guardian) external {
        require(block.chainid == ChainIds.FOUNDRY, "Only for local tests");
        protocolAdmin = protocol;
        timelockExecutor = executor;
        timelockGuardian = guardian;
    }

    /// @notice path where alloc file will be stored
    function getDumpPath() internal view returns (string memory) {
        if (block.chainid == ChainIds.ILIAD) {
            return "./iliad-alloc.json";
        } else if (block.chainid == ChainIds.MININET) {
            return "./mininet-alloc.json";
        } else if (block.chainid == ChainIds.ODYSSEY_DEVNET) {
            return "./odyssey-devnet-alloc.json";
        } else if (block.chainid == ChainIds.ODYSSEY_TESTNET) {
            return "./odyssey-testnet-alloc.json";
        } else if (block.chainid == ChainIds.LOCAL) {
            return "./local-alloc.json";
        } else if (block.chainid == ChainIds.FOUNDRY) {
            return "./foundry-alloc.json";
        } else if (block.chainid == ChainIds.STORY_MAINNET) {
            return "./mainnet-alloc.json";
        } else {
            revert("Unsupported chain id");
        }
    }

    /// @notice Get the minimum delay for the timelock
    function getTimelockMinDelay() internal view returns (uint256) {
        if (block.chainid == ChainIds.ILIAD) {
            // Iliad
            return 1 days;
        } else if (block.chainid == ChainIds.MININET) {
            // Mininet
            return 10 seconds;
        } else if (block.chainid == ChainIds.ODYSSEY_DEVNET) {
            // Odyssey devnet
            return 10 seconds;
        } else if (block.chainid == ChainIds.ODYSSEY_TESTNET) {
            // Odyssey testnet
            return 1 days;
        } else if (block.chainid == ChainIds.LOCAL) {
            // Local
            return 10 seconds;
        } else if (block.chainid == ChainIds.FOUNDRY) {
            // Foundry
            return 10 seconds;
        } else if (block.chainid == ChainIds.STORY_MAINNET) {
            // Mainnet
            return 2 days;
        } else {
            revert("Unsupported chain id");
        }
    }

    /// @notice main script method
    function run() public {
        // Tests should set these addresses first
        if (protocolAdmin == address(0)) {
            protocolAdmin = vm.envAddress("ADMIN_ADDRESS");
        }
        require(protocolAdmin != address(0), "protocolAdmin not set");

        if (timelockExecutor == address(0)) {
            timelockExecutor = vm.envAddress("TIMELOCK_EXECUTOR_ADDRESS");
        }
        if (timelockExecutor == address(0)) {
            console2.log("TimelockExecutor not set, executing timelock operations is public");
        }

        if (timelockGuardian == address(0)) {
            timelockGuardian = vm.envAddress("TIMELOCK_GUARDIAN_ADDRESS");
        }
        require(timelockGuardian != address(0), "canceller not set");

        if (block.chainid == ChainIds.STORY_MAINNET) {
            require(!KEEP_TIMELOCK_ADMIN_ROLE, "Timelock admin role not allowed on mainnet");
        } else {
            console2.log("Will timelock admin role be assigned?", KEEP_TIMELOCK_ADMIN_ROLE);
        }

        vm.startPrank(deployer);

        setPredeploys();
        setAllocations();

        // Necessary to skip for tests
        if (saveState) {
            // Reset so its not included state dump
            vm.etch(msg.sender, "");
            vm.resetNonce(msg.sender);
            vm.deal(msg.sender, 0);

            vm.etch(deployer, "");
            // Not resetting nonce
            vm.deal(deployer, 0);
        }

        vm.stopPrank();
        if (saveState) {
            vm.dumpState(dumpPath);
            console2.log("Alloc saved to:", dumpPath);
        }
    }

    /// @notice Prepares the bytecode and storage for predeployed contracts in genesis file
    function setPredeploys() internal {
        // Set predeploys that are outside of the proxied Namespace and Timelock
        setCreate3();
        deployTimelock();
        setERC6551();
        setWIP();

        // Set proxies for all predeploys in the proxied Namespace
        setProxies();

        // Set implementations for predeploys that are used since genesis
        setStaking();
        setUpgrade();
        setUBIPool();
    }

    /// @dev Populates the upgradeable predeploys namespace with proxies, to reserve the addresses
    /// for future use. Implementations are deterministically determined, but won't have code
    /// unless explicitly set in setPredeploys(). Later on, they can be upgraded to new
    /// implementations by governance.
    function setProxies() internal {
        for (uint160 i = 1; i <= Predeploys.NamespaceSize; i++) {
            address addr = address(uint160(Predeploys.Namespace) + i);
            setProxy(addr);
        }
    }

    /// @notice Deploy TimelockController to manage upgrades and admin actions
    /// @dev this is a deterministic deployment, not a predeploy (won't show in genesis file).
    function deployTimelock() internal {
        // WARNING: Make sure protocolAdmin and timelockGuardian are multisigs on mainnet
        uint256 minDelay = getTimelockMinDelay();
        address[] memory proposers = new address[](1);
        proposers[0] = protocolAdmin;
        address[] memory executors = new address[](1);
        executors[0] = timelockExecutor;
        address canceller = timelockGuardian;

        bytes memory creationCode = abi.encodePacked(
            type(TimelockController).creationCode,
            abi.encode(minDelay, proposers, executors, protocolAdmin)
        );
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER");
        // We deploy this with Create3 because we can't set storage variables in constructor with vm.etch
        timelock = Create3(Predeploys.Create3).deploy(salt, creationCode);
        vm.stopPrank();
        bytes32 cancellerRole = TimelockController(payable(timelock)).CANCELLER_ROLE();
        vm.prank(protocolAdmin);
        TimelockController(payable(timelock)).grantRole(cancellerRole, canceller);
        if (!KEEP_TIMELOCK_ADMIN_ROLE) {
            bytes32 adminRole = TimelockController(payable(timelock)).DEFAULT_ADMIN_ROLE();
            TimelockController(payable(timelock)).renounceRole(adminRole, protocolAdmin);
        }
        vm.stopPrank();
        vm.startPrank(deployer);

        console2.log("TimelockController deployed at:", timelock);
    }

    /// @notice Set a TransparentUpgradeableProxy bytecode and storage for a predeploy address,
    /// within the proxied Namespace
    /// @dev We use a deterministic implementation address
    function setProxy(address proxyAddr) internal {
        address impl = Predeploys.getImplAddress(proxyAddr);

        // set impl code to non-zero length, so it passes TransparentUpgradeableProxy constructor check
        // assert it is not already set
        require(impl.code.length == 0, "impl already set");
        vm.etch(impl, "00");

        // use new, so that the immutable variable the holds the ProxyAdmin proxyAddr is set in properly in bytecode
        address tmp = address(new TransparentUpgradeableProxy(impl, timelock, ""));
        vm.etch(proxyAddr, tmp.code);

        // set implempentation storage manually
        EIP1967Helper.setImplementation(proxyAddr, impl);

        // set admin storage, to follow EIP1967 standard
        EIP1967Helper.setAdmin(proxyAddr, EIP1967Helper.getAdmin(tmp));

        // reset impl & tmp
        vm.etch(impl, "");
        vm.etch(tmp, "");

        vm.resetNonce(tmp);
        vm.deal(impl, 1);
        vm.deal(proxyAddr, 1);
    }

    /// @notice Sets the bytecode for the implementation of IPTokenStaking predeploy
    function setStaking() internal {
        address impl = Predeploys.getImplAddress(Predeploys.Staking);

        address tmp = address(
            new IPTokenStaking(
                1 ether // defaultMinFee, 1 IP
            )
        );
        console2.log("tpm", tmp);
        vm.etch(impl, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        InitializableHelper.disableInitializers(impl);
        IIPTokenStaking.InitializerArgs memory args = IIPTokenStaking.InitializerArgs({
            owner: timelock,
            minStakeAmount: 1024 ether,
            minUnstakeAmount: 1024 ether,
            minCommissionRate: 5_00, // 5% in basis points
            fee: 1 ether // 1 IP
        });

        IPTokenStaking(Predeploys.Staking).initialize(args);

        console2.log("IPTokenStaking proxy deployed at:", Predeploys.Staking);
        console2.log("IPTokenStaking ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.Staking));
        console2.log("IPTokenStaking impl at:", EIP1967Helper.getImplementation(Predeploys.Staking));
        console2.log("IPTokenStaking owner:", IPTokenStaking(Predeploys.Staking).owner());
    }

    /// @notice Sets the bytecode for the implementation of UpgradeEntrypoint predeploy
    function setUpgrade() internal {
        address impl = Predeploys.getImplAddress(Predeploys.Upgrades);
        address tmp = address(new UpgradeEntrypoint());

        console2.log("tpm", tmp);
        vm.etch(impl, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        InitializableHelper.disableInitializers(impl);
        UpgradeEntrypoint(Predeploys.Upgrades).initialize(timelock);

        console2.log("UpgradeEntrypoint proxy deployed at:", Predeploys.Upgrades);
        console2.log("UpgradeEntrypoint ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.Upgrades));
        console2.log("UpgradeEntrypoint impl at:", EIP1967Helper.getImplementation(Predeploys.Upgrades));
        console2.log("UpgradeEntrypoint owner:", UpgradeEntrypoint(Predeploys.Upgrades).owner());
    }

    /// @notice Sets the bytecode for the implementation of UBIPool predeploy
    function setUBIPool() internal {
        address impl = Predeploys.getImplAddress(Predeploys.UBIPool);
        address tmp = address(new UBIPool(20_00)); // 20% UBI
        vm.etch(impl, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        InitializableHelper.disableInitializers(impl);
        UBIPool(Predeploys.UBIPool).initialize(timelock);

        console2.log("UBIPool proxy deployed at:", Predeploys.UBIPool);
        console2.log("UBIPool ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.UBIPool));
        console2.log("UBIPool impl at:", EIP1967Helper.getImplementation(Predeploys.UBIPool));
        console2.log("UBIPool owner:", UBIPool(Predeploys.UBIPool).owner());
    }

    /// @notice Sets the bytecode for Create3 factory as a predeploy
    /// @dev Create3 factory address https://github.com/ZeframLou/create3-factory
    function setCreate3() internal {
        address tmp = address(new Create3());
        vm.etch(Predeploys.Create3, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        vm.deal(Predeploys.Create3, 1);
        console2.log("Create3 deployed at:", Predeploys.Create3);
    }

    /// @notice Sets the bytecode for ERC6551Registry as a predeploy
    /// @dev ERC6551Registry as defined by ERC-6551
    function setERC6551() internal {
        address tmp = address(new ERC6551Registry());
        vm.etch(Predeploys.ERC6551Registry, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        vm.deal(Predeploys.ERC6551Registry, 1);
        console2.log("ERC6551 deployed at:", Predeploys.ERC6551Registry);
    }

    /// @notice Sets the bytecode for WIP as a predeploy
    /// @dev WIP is the ERC20 wrapper for IP token
    function setWIP() internal {
        address tmp = address(new WIP());
        vm.etch(Predeploys.WIP, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        vm.deal(Predeploys.WIP, 1);
        console2.log("WIP deployed at:", Predeploys.WIP);
    }

    /// @notice Sets initial balances for predeploys and genesis allocations
    function setAllocations() internal {
        // EL Predeploys
        // Geth precompile 1 wei allocation (Accounts with 0 balance and no EVM code may be removed from
        // the state trie, 1 wei balance prevents this).
        vm.deal(0x0000000000000000000000000000000000000001, 1);
        vm.deal(0x0000000000000000000000000000000000000002, 1);
        vm.deal(0x0000000000000000000000000000000000000003, 1);
        vm.deal(0x0000000000000000000000000000000000000004, 1);
        vm.deal(0x0000000000000000000000000000000000000005, 1);
        vm.deal(0x0000000000000000000000000000000000000006, 1);
        vm.deal(0x0000000000000000000000000000000000000007, 1);
        vm.deal(0x0000000000000000000000000000000000000008, 1);
        vm.deal(0x0000000000000000000000000000000000000009, 1);
        // p256 verification precompile
        vm.deal(0x0000000000000000000000000000000000000100, 1);
        // Story's IPGraph precompile
        vm.deal(0x0000000000000000000000000000000000000101, 1);
        // Allocation
        if (block.chainid == ChainIds.STORY_MAINNET) {
            // TBD
        } else if (block.chainid == ChainIds.ODYSSEY_DEVNET) {
            // Odyssey devnet alloc
            vm.deal(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, 100000000 ether);
            vm.deal(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab, 100000000 ether);
            vm.deal(0xEcB1D051475A7e330b1DD6683cdC7823Bbcf8Dcf, 100000000 ether);
            vm.deal(0x5518D1BD054782792D2783509FbE30fa9D888888, 100000000 ether);
            vm.deal(0xbd39FAe873F301b53e14d365383118cD4a222222, 100000000 ether);
            vm.deal(0x00FCeC044cD73e8eC6Ad771556859b00C9011111, 100000000 ether);
            vm.deal(0xb5350B7CaE94C2bF6B2b56Ef6A06cC1153900000, 100000000 ether);
            vm.deal(0x13919a0d8603c35DAC923f92D7E4e1D55e993898, 100000000 ether);
        } else if (block.chainid == ChainIds.ODYSSEY_TESTNET) {
            // Odyssey testnet alloc
            vm.deal(0x5687400189B13551137e330F7ae081142EdfD866, 200000000 ether);
            vm.deal(0x56A26642ad963D3542DdAe4d8fdECC396153c2f6, 200000000 ether);
            vm.deal(0x12cBb8F6F2F7d48bB22B6A1b12452381A45bEb7c, 100000000 ether);
            vm.deal(0xD26078bA39afccec71E0D68a151a853d21950FF0, 200000000 ether);
            vm.deal(0xcA93A8f7a3971D208670876202D8353Ca3D6869a, 200000000 ether);
            vm.deal(0x8Ffc89da28DD2F5f7582B0459505E9a615623791, 10000000 ether);
            vm.deal(0xE8DA8e345Ab1556E5DeE19F9c369C827561Ff712, 10000000 ether);
        } else {
            // Default network alloc
            vm.deal(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, 100000000 ether);
            vm.deal(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab, 100000000 ether);
            vm.deal(0xEcB1D051475A7e330b1DD6683cdC7823Bbcf8Dcf, 100000000 ether);
            vm.deal(0x5518D1BD054782792D2783509FbE30fa9D888888, 100000000 ether);
            vm.deal(0xbd39FAe873F301b53e14d365383118cD4a222222, 100000000 ether);
            vm.deal(0x00FCeC044cD73e8eC6Ad771556859b00C9011111, 100000000 ether);
            vm.deal(0xb5350B7CaE94C2bF6B2b56Ef6A06cC1153900000, 100000000 ether);
            vm.deal(0x13919a0d8603c35DAC923f92D7E4e1D55e993898, 100000000 ether);
            vm.deal(0x64a2fdc6f7CD8AA42e0bb59bf80bC47bFFbe4a73, 100000000 ether);
        }
        if (ALLOCATE_10K_TEST_ACCOUNTS && block.chainid != ChainIds.STORY_MAINNET) {
            setTestAllocations();
        }
    }

    /// @notice Sets 10,000 test accounts with increasing balances
    function setTestAllocations() internal {
        address allocSpace = address(0xBBbbbB0000000000000000000000000000000000);
        for (uint160 i = 1; i <= 10_000; i++) {
            vm.deal(address(uint160(allocSpace) + i), i * 1 ether);
        }
    }
}
