// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import { IIPTokenStaking } from "../src/interfaces/IIPTokenStaking.sol";
import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";
import { UpgradeEntrypoint } from "../src/protocol/UpgradeEntrypoint.sol";

import { EIP1967Helper } from "./utils/EIP1967Helper.sol";
import { InitializableHelper } from "./utils/InitializableHelper.sol";
import { Predeploys } from "../src/libraries/Predeploys.sol";

/**
 * @title GenerateAlloc
 * @dev A script to generate the alloc section of EL genesis
 * - Predeploys (See src/libraries/Predeploys.sol)
 * - Genesis $IP allocations (chain id dependent)
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

    // Upgrade admin controls upgradeability (by being Owner of each ProxyAdmin),
    // protocol admin is Owner of precompiles (admin/governance methods).
    // To disable upgradeability, we transfer ProxyAdmin ownership to a dead address
    address internal upgradeAdmin;
    address internal protocolAdmin;
    string internal dumpPath = getDumpPath();
    bool public saveState = true;
    uint256 public constant MAINNET_CHAIN_ID = 0; // TBD

    /// @notice call from Test.sol to run test fast (no json saving)
    function disableStateDump() external {
        require(block.chainid == 31337, "Only for local tests");
        saveState = false;
    }

    /// @notice call from Test.sol only
    function setAdminAddresses(address upgrade, address protocol) external {
        require(block.chainid == 31337, "Only for local tests");
        upgradeAdmin = upgrade;
        protocolAdmin = protocol;
    }

    /// @notice path where alloc file will be stored
    function getDumpPath() internal view returns (string memory) {
        if (block.chainid == 1513) {
            return "./iliad-alloc.json";
        } else if (block.chainid == 1512) {
            return "./mininet-alloc.json";
        } else if (block.chainid == 1315) {
            return "./odyssey-devnet-alloc.json";
        } else if (block.chainid == 31337) {
            return "./local-alloc.json";
        } else {
            revert("Unsupported chain id");
        }
    }

    /// @notice main script method
    function run() public {
        if (upgradeAdmin == address(0)) {
            upgradeAdmin = vm.envAddress("UPGRADE_ADMIN_ADDRESS");
        }
        require(upgradeAdmin != address(0), "upgradeAdmin not set");

        if (protocolAdmin == address(0)) {
            protocolAdmin = vm.envAddress("ADMIN_ADDRESS");
        }
        require(protocolAdmin != address(0), "protocolAdmin not set");

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

    function setPredeploys() internal {
        setProxy(Predeploys.Staking);
        setProxy(Predeploys.Upgrades);

        setStaking();
        setUpgrade();
    }

    function setProxy(address proxyAddr) internal {
        address impl = Predeploys.getImplAddress(proxyAddr);

        // set impl code to non-zero length, so it passes TransparentUpgradeableProxy constructor check
        // assert it is not already set
        require(impl.code.length == 0, "impl already set");
        vm.etch(impl, "00");

        // use new, so that the immutable variable the holds the ProxyAdmin proxyAddr is set in properly in bytecode
        address tmp = address(new TransparentUpgradeableProxy(impl, upgradeAdmin, ""));
        vm.etch(proxyAddr, tmp.code);

        // set implempentation storage manually
        EIP1967Helper.setImplementation(proxyAddr, impl);

        // set admin storage, to follow EIP1967 standard
        EIP1967Helper.setAdmin(proxyAddr, EIP1967Helper.getAdmin(tmp));

        // reset impl & tmp
        vm.etch(impl, "");
        vm.etch(tmp, "");

        // can we reset nonce here? we are using "deployer" proxyAddr
        vm.resetNonce(tmp);
        vm.deal(impl, 1);
        vm.deal(proxyAddr, 1);
    }

    /**
     * @notice Setup Staking predeploy
     */
    function setStaking() internal {
        address impl = Predeploys.getImplAddress(Predeploys.Staking);

        address tmp = address(new IPTokenStaking(
            1 gwei, // stakingRounding
            1 ether // defaultMinUnjailFee, 1 IP
        ));
        console2.log("tpm", tmp);
        vm.etch(impl, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        InitializableHelper.disableInitializers(impl);
        IIPTokenStaking.InitializerArgs memory args = IIPTokenStaking.InitializerArgs({
            owner: protocolAdmin,
            minStakeAmount: 1 ether,
            minUnstakeAmount: 1 ether,
            minCommissionRate: 5_00, // 5% in basis points
            shortStakingPeriod: 1 days, // TBD
            mediumStakingPeriod: 2 days, // TBD
            longStakingPeriod: 3 days, // TBD
            unjailFee: 1 ether
        });

        // Testnet timing values
        if (block.chainid != MAINNET_CHAIN_ID) {
            args.minCommissionRate = 5_00 seconds;
            args.shortStakingPeriod = 10 seconds;
            args.mediumStakingPeriod = 15 seconds;
            args.longStakingPeriod = 20 seconds;
        }

        IPTokenStaking(Predeploys.Staking).initialize(args);

        console2.log("IPTokenStaking proxy deployed at:", Predeploys.Staking);
        console2.log("IPTokenStaking ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.Staking));
        console2.log("IPTokenStaking impl at:", EIP1967Helper.getImplementation(Predeploys.Staking));
        console2.log("IPTokenStaking owner:", IPTokenStaking(Predeploys.Staking).owner());
    }

    /**
     * @notice Setup Upgrade predeploy
     */
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
        UpgradeEntrypoint(Predeploys.Upgrades).initialize(protocolAdmin);

        console2.log("UpgradeEntrypoint proxy deployed at:", Predeploys.Upgrades);
        console2.log("UpgradeEntrypoint ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.Upgrades));
        console2.log("UpgradeEntrypoint impl at:", EIP1967Helper.getImplementation(Predeploys.Upgrades));
    }

    function setAllocations() internal {
        // EL Predeploys
        vm.deal(0x0000000000000000000000000000000000000001, 1);
        vm.deal(0x0000000000000000000000000000000000000001, 1);
        vm.deal(0x0000000000000000000000000000000000000002, 1);
        vm.deal(0x0000000000000000000000000000000000000003, 1);
        vm.deal(0x0000000000000000000000000000000000000004, 1);
        vm.deal(0x0000000000000000000000000000000000000005, 1);
        vm.deal(0x0000000000000000000000000000000000000006, 1);
        vm.deal(0x0000000000000000000000000000000000000007, 1);
        vm.deal(0x0000000000000000000000000000000000000008, 1);
        vm.deal(0x0000000000000000000000000000000000000009, 1);
        vm.deal(0x000000000000000000000000000000000000001a, 1);
        // Allocation
        if (block.chainid == MAINNET_CHAIN_ID) {
            // TBD
        } else {
            // Testnet alloc
            vm.deal(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, 100000000 ether);
            vm.deal(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab, 100000000 ether);
            vm.deal(0xEcB1D051475A7e330b1DD6683cdC7823Bbcf8Dcf, 100000000 ether);
            vm.deal(0x5518D1BD054782792D2783509FbE30fa9D888888, 100000000 ether);
            vm.deal(0xbd39FAe873F301b53e14d365383118cD4a222222, 100000000 ether);
            vm.deal(0x00FCeC044cD73e8eC6Ad771556859b00C9011111, 100000000 ether);
            vm.deal(0xb5350B7CaE94C2bF6B2b56Ef6A06cC1153900000, 100000000 ether);
        }
    }
}
