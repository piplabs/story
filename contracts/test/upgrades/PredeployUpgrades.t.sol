// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
// solhint-disable-next-line max-line-length
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";
import { UBIPool } from "../../src/protocol/UBIPool.sol";

import { EIP1967Helper } from "../../script/utils/EIP1967Helper.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Test } from "../utils/Test.sol";

import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";

abstract contract MockNewFeatures {
    function foo() external pure returns (string memory) {
        return "bar";
    }
}

contract IPTokenStakingV2 is IPTokenStaking, MockNewFeatures {
    constructor(uint256 stakingRounding, uint256 defaultMinFee) IPTokenStaking(stakingRounding, defaultMinFee) {}
}

contract UpgradeEntrypointV2 is UpgradeEntrypoint, MockNewFeatures {}

contract UBIPoolV2 is UBIPool, MockNewFeatures {
    constructor(uint32 maxUBIPercentage) UBIPool(maxUBIPercentage) {}
}

contract InitialImplementation {
    function foo() external pure returns (string memory) {
        return "bar";
    }
}

/**
 * @title PredeployUpgrades
 * @dev A script to test upgrading the precompile contracts
 */
contract PredeployUpgrades is Test {
    function testUpgradeStaking() public {
        // ---- Staking
        address newImpl = address(
            new IPTokenStakingV2(
                1 gwei, // stakingRounding
                1 ether
            )
        );
        ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(Predeploys.Staking));
        assertEq(proxyAdmin.owner(), address(timelock));

        performTimelocked(
            address(proxyAdmin),
            abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(Predeploys.Staking),
                newImpl,
                ""
            )
        );

        assertEq(EIP1967Helper.getImplementation(Predeploys.Staking), newImpl, "Staking not upgraded");
        assertEq(
            keccak256(abi.encode(IPTokenStakingV2(Predeploys.Staking).foo())),
            keccak256(abi.encode("bar")),
            "Upgraded to wrong iface"
        );
    }

    function testUpgradeUpgradeEntrypoint() public {
        // ---- Upgrades
        address newImpl = address(new UpgradeEntrypointV2());
        ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(Predeploys.Upgrades));
        assertEq(proxyAdmin.owner(), address(timelock));

        performTimelocked(
            address(proxyAdmin),
            abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(Predeploys.Upgrades),
                newImpl,
                ""
            )
        );
        assertEq(EIP1967Helper.getImplementation(Predeploys.Upgrades), newImpl, "Upgrades not upgraded");
        assertEq(
            keccak256(abi.encode(IPTokenStakingV2(Predeploys.Upgrades).foo())),
            keccak256(abi.encode("bar")),
            "Upgraded to wrong iface"
        );
    }

    function testUpgradeUBIPool() public {
        // ---- UBIPool
        address newImpl = address(new UBIPoolV2(10_00));
        ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(Predeploys.UBIPool));
        assertEq(proxyAdmin.owner(), address(timelock));

        performTimelocked(
            address(proxyAdmin),
            abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(Predeploys.UBIPool),
                newImpl,
                ""
            )
        );
        assertEq(EIP1967Helper.getImplementation(Predeploys.UBIPool), newImpl, "Upgrades not upgraded");
        assertEq(
            keccak256(abi.encode(IPTokenStakingV2(Predeploys.UBIPool).foo())),
            keccak256(abi.encode("bar")),
            "Upgraded to wrong iface"
        );
    }

    function testUpgradeUnusedProxies() public {
        address initialImpl = address(new InitialImplementation());

        for (
            uint160 i = uint160(Predeploys.Upgrades) + uint160(1);
            i <= uint160(Predeploys.Namespace) + Predeploys.NamespaceSize;
            i++
        ) {
            // Verify predeploy is proxied and not upgraded
            address predeploy = address(i);
            assertTrue(Predeploys.proxied(predeploy), "Predeploy not proxied");
            assertEq(
                EIP1967Helper.getImplementation(predeploy),
                Predeploys.getImplAddress(predeploy),
                "Predeploy upgraded already"
            );
            ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(predeploy));
            assertEq(proxyAdmin.owner(), address(timelock));

            // Revert if not owner
            vm.expectRevert();
            proxyAdmin.upgradeAndCall(ITransparentUpgradeableProxy(predeploy), initialImpl, "");

            // Upgrade predeploy
            performTimelocked(
                address(proxyAdmin),
                abi.encodeWithSelector(
                    ProxyAdmin.upgradeAndCall.selector,
                    ITransparentUpgradeableProxy(predeploy),
                    initialImpl,
                    ""
                )
            );

            // Verify predeploy is upgraded
            assertEq(EIP1967Helper.getImplementation(predeploy), initialImpl, "Predeploy not upgraded");
            assertEq(
                keccak256(abi.encode(InitialImplementation(predeploy).foo())),
                keccak256(abi.encode("bar")),
                "Upgraded to wrong iface"
            );
        }
    }

    function testRenounceUpgradeability() public {
        address newImpl = address(new InitialImplementation());
        // For each predeploy, renounce upgradeability
        for (
            uint160 i = uint160(Predeploys.Namespace) + uint160(1);
            i <= uint160(Predeploys.Namespace) + Predeploys.NamespaceSize;
            i++
        ) {
            address predeploy = address(i);
            ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(predeploy));
            assertEq(proxyAdmin.owner(), address(timelock));
            // Renounce ownership
            performTimelocked(address(proxyAdmin), abi.encodeWithSelector(Ownable.renounceOwnership.selector));

            // Verify ownership was renounced
            assertEq(proxyAdmin.owner(), address(0), "Ownership not renounced");

            // Try to upgrade
            bytes memory upgradeAction = abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(predeploy),
                newImpl,
                ""
            );
            bytes memory expectedReason = abi.encodeWithSelector(
                Ownable.OwnableUnauthorizedAccount.selector,
                address(timelock)
            );
            expectRevertTimelocked(address(proxyAdmin), upgradeAction, expectedReason);
        }
    }
}
