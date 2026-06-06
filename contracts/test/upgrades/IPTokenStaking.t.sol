// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Test } from "../utils/Test.sol";
import { Initializable } from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { PausableUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import { AccessControlUpgradeable } from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import { IAccessControl } from "@openzeppelin/contracts/access/IAccessControl.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { IPTokenStaking } from "../../src/upgrades/IPTokenStaking.sol";
import { IIPTokenStaking } from "../../src/interfaces/IIPTokenStaking.sol";
import { EIP1967Helper } from "../../script/utils/EIP1967Helper.sol";

import { console2 } from "forge-std/console2.sol";

/**
 * @title IPTokenStakingTest
 * @dev A test for the IPTokenStaking contract
 */
contract IPTokenStakingTest is Test {
    IPTokenStaking ipTokenStakingProxy;
    address safeGovernanceMultisig;
    address securityCouncilMultisig;
    address oldOwner;

    uint256 public monikerLengthBefore;
    uint256 public stakeRoundingBefore;
    bytes32 public pauserRoleBefore;
    uint256 public defaultMinFeeBefore;
    uint256 public maxDataLengthBefore;
    uint256 public minCommissionRateBefore;
    uint256 public minStakeAmountBefore;
    uint256 public minUnstakeAmountBefore;
    uint256 public feeBefore;

    function setUp() public override {
        // Fork the desired network where UMA contracts are deployed
        uint256 forkId = vm.createFork("https://mainnet.storyrpc.io/");
        vm.selectFork(forkId);

        // Mainnet related addresses
        ipTokenStakingProxy = IPTokenStaking(0xCCcCcC0000000000000000000000000000000001);
        safeGovernanceMultisig = 0xF07cA4b61022F0399C1511E7E668A57567f2138B;
        securityCouncilMultisig = 0x25D2605b2C768082A14E79713114389d0eC297D8;
        timelock = TimelockController(payable(0x6c7FA8DF1B8Dc29a7481Bb65ad590D2D16787a82));

        oldOwner = OwnableUpgradeable(address(ipTokenStakingProxy)).owner();
        monikerLengthBefore = ipTokenStakingProxy.MAX_MONIKER_LENGTH();
        stakeRoundingBefore = ipTokenStakingProxy.STAKE_ROUNDING();
        defaultMinFeeBefore = ipTokenStakingProxy.DEFAULT_MIN_FEE();
        maxDataLengthBefore = ipTokenStakingProxy.MAX_DATA_LENGTH();
        minCommissionRateBefore = ipTokenStakingProxy.minCommissionRate();
        minStakeAmountBefore = ipTokenStakingProxy.minStakeAmount();
        minUnstakeAmountBefore = ipTokenStakingProxy.minUnstakeAmount();
        feeBefore = ipTokenStakingProxy.fee();

        // deploy new implementation
        IPTokenStaking newImpl = new IPTokenStaking(1000000000000000000, 256);

        // upgrade the proxy to the new implementation
        ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(address(ipTokenStakingProxy)));
        console2.log("proxyAdmin", address(proxyAdmin));
        vm.startPrank(safeGovernanceMultisig);
        timelock.schedule(
            address(proxyAdmin),
            0,
            abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(address(ipTokenStakingProxy)),
                newImpl,
                abi.encodeWithSelector(
                    IPTokenStaking.initializeV2.selector,
                    safeGovernanceMultisig,
                    securityCouncilMultisig
                )
            ),
            bytes32(0),
            bytes32(0),
            timelock.getMinDelay()
        );

        vm.warp(block.timestamp + timelock.getMinDelay() + 1);

        timelock.execute(
            address(proxyAdmin),
            0,
            abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(address(ipTokenStakingProxy)),
                newImpl,
                abi.encodeWithSelector(
                    IPTokenStaking.initializeV2.selector,
                    safeGovernanceMultisig,
                    securityCouncilMultisig
                )
            ),
            bytes32(0),
            bytes32(0)
        );
        vm.stopPrank();
    }

    function testNewRoles() public {
        assertEq(
            AccessControlUpgradeable(address(ipTokenStakingProxy)).hasRole(
                AccessControlUpgradeable(address(ipTokenStakingProxy)).DEFAULT_ADMIN_ROLE(),
                address(timelock)
            ),
            true
        );
        assertEq(
            AccessControlUpgradeable(address(ipTokenStakingProxy)).hasRole(
                ipTokenStakingProxy.PAUSER_ROLE(),
                address(safeGovernanceMultisig)
            ),
            true
        );
        assertEq(
            AccessControlUpgradeable(address(ipTokenStakingProxy)).hasRole(
                ipTokenStakingProxy.PAUSER_ROLE(),
                address(securityCouncilMultisig)
            ),
            true
        );
    }

    function testInitializeRevertWhenCalledTwice() public {
        vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));
        ipTokenStakingProxy.initialize(
            IIPTokenStaking.InitializerArgs({
                minStakeAmount: 1,
                minUnstakeAmount: 1,
                minCommissionRate: 1,
                fee: 1,
                owner: address(1)
            })
        );
    }

    function testInitializeV2RevertWhenCalledTwice() public {
        vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));
        ipTokenStakingProxy.initializeV2(address(1), address(1));
    }

    function testPauseUnpauseSafeGovernanceMultisig() public {
        assertEq(ipTokenStakingProxy.paused(), false);
        vm.startPrank(safeGovernanceMultisig);
        ipTokenStakingProxy.pause();
        assertEq(ipTokenStakingProxy.paused(), true);
        ipTokenStakingProxy.unpause();
        assertEq(ipTokenStakingProxy.paused(), false);
    }

    function testPauseUnpauseSecurityCouncilMultisig() public {
        assertEq(ipTokenStakingProxy.paused(), false);
        vm.startPrank(securityCouncilMultisig);
        ipTokenStakingProxy.pause();
        assertEq(ipTokenStakingProxy.paused(), true);
        ipTokenStakingProxy.unpause();
        assertEq(ipTokenStakingProxy.paused(), false);
    }

    function testSetMinStakeAmountRevertWhenNotAdmin() public {
        vm.startPrank(address(1));
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                address(1),
                ipTokenStakingProxy.DEFAULT_ADMIN_ROLE()
            )
        );
        ipTokenStakingProxy.setMinStakeAmount(1 ether);
    }

    function testSetMinStakeAmount() public {
        vm.startPrank(address(timelock));
        ipTokenStakingProxy.setMinStakeAmount(1 ether);
        assertEq(ipTokenStakingProxy.minStakeAmount(), 1 ether);
    }

    function testSetMinUnstakeAmountRevertWhenNotAdmin() public {
        vm.startPrank(address(1));
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                address(1),
                ipTokenStakingProxy.DEFAULT_ADMIN_ROLE()
            )
        );
        ipTokenStakingProxy.setMinUnstakeAmount(1 ether);
    }

    function testSetMinUnstakeAmount() public {
        vm.startPrank(address(timelock));
        ipTokenStakingProxy.setMinUnstakeAmount(1 ether);
        assertEq(ipTokenStakingProxy.minUnstakeAmount(), 1 ether);
    }

    function testSetFeeRevertWhenNotAdmin() public {
        vm.startPrank(address(1));
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                address(1),
                ipTokenStakingProxy.DEFAULT_ADMIN_ROLE()
            )
        );
        ipTokenStakingProxy.setFee(1);
    }

    function testSetFee() public {
        vm.startPrank(address(timelock));
        ipTokenStakingProxy.setFee(1000000000000000001);
        assertEq(ipTokenStakingProxy.fee(), 1000000000000000001);
    }

    function testSetMinCommissionRateRevertWhenNotAdmin() public {
        vm.startPrank(address(1));
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                address(1),
                ipTokenStakingProxy.DEFAULT_ADMIN_ROLE()
            )
        );
        ipTokenStakingProxy.setMinCommissionRate(1);
    }

    function testSetMinCommissionRate() public {
        vm.startPrank(address(timelock));
        ipTokenStakingProxy.setMinCommissionRate(1);
        assertEq(ipTokenStakingProxy.minCommissionRate(), 1);
    }

    function testStakeRevertWhenPaused() public {
        vm.deal(safeGovernanceMultisig, 1 ether);
        vm.startPrank(safeGovernanceMultisig);
        ipTokenStakingProxy.pause();
        assertEq(ipTokenStakingProxy.paused(), true);
        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
        ipTokenStakingProxy.stake{ value: 1 ether }(bytes(""), IIPTokenStaking.StakingPeriod.FLEXIBLE, bytes(""));
    }

    function testStakeOnBehalfRevertWhenPaused() public {
        vm.deal(safeGovernanceMultisig, 1 ether);
        vm.startPrank(safeGovernanceMultisig);
        ipTokenStakingProxy.pause();
        assertEq(ipTokenStakingProxy.paused(), true);
        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
        ipTokenStakingProxy.stakeOnBehalf{ value: 1 ether }(
            address(1),
            bytes(""),
            IIPTokenStaking.StakingPeriod.FLEXIBLE,
            bytes("")
        );
    }

    // -------------------------------------------------------------------------
    // whenNotPaused coverage — every gated entry point must revert while paused.
    // Before this, only stake / stakeOnBehalf were covered. whenNotPaused is the
    // first modifier on each, so it reverts before any param validation; dummy
    // args + zero value are sufficient.
    // -------------------------------------------------------------------------

    function _pauseAsGov() internal {
        vm.prank(safeGovernanceMultisig);
        ipTokenStakingProxy.pause();
        assertEq(ipTokenStakingProxy.paused(), true);
    }

    function _expectEnforcedPause() internal {
        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    }

    function testSetOperatorRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.setOperator(address(1));
    }

    function testUnsetOperatorRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.unsetOperator();
    }

    function testSetWithdrawalAddressRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.setWithdrawalAddress(address(1));
    }

    function testSetRewardsAddressRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.setRewardsAddress(address(1));
    }

    function testCreateValidatorRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.createValidator(bytes(""), "", 0, 0, 0, false, bytes(""));
    }

    function testUpdateValidatorCommissionRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.updateValidatorCommission(bytes(""), 0);
    }

    function testRedelegateRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.redelegate(bytes(""), bytes(""), 0, 0);
    }

    function testRedelegateOnBehalfRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.redelegateOnBehalf(address(1), bytes(""), bytes(""), 0, 0);
    }

    function testUnstakeRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.unstake(bytes(""), 0, 0, bytes(""));
    }

    function testUnstakeOnBehalfRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.unstakeOnBehalf(address(1), bytes(""), 0, 0, bytes(""));
    }

    function testUnjailRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.unjail(bytes(""), bytes(""));
    }

    function testUnjailOnBehalfRevertWhenPaused() public {
        _pauseAsGov();
        _expectEnforcedPause();
        ipTokenStakingProxy.unjailOnBehalf(bytes(""), bytes(""));
    }

    // -------------------------------------------------------------------------
    // pause / unpause access control — only PAUSER_ROLE holders.
    // -------------------------------------------------------------------------

    function testPauseRevertWhenNotPauser() public {
        vm.startPrank(address(1));
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                address(1),
                ipTokenStakingProxy.PAUSER_ROLE()
            )
        );
        ipTokenStakingProxy.pause();
    }

    function testUnpauseRevertWhenNotPauser() public {
        _pauseAsGov();
        vm.startPrank(address(1));
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                address(1),
                ipTokenStakingProxy.PAUSER_ROLE()
            )
        );
        ipTokenStakingProxy.unpause();
    }

    // -------------------------------------------------------------------------
    // pause state machine edges.
    // -------------------------------------------------------------------------

    function testPauseRevertWhenAlreadyPaused() public {
        _pauseAsGov();
        vm.prank(securityCouncilMultisig);
        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
        ipTokenStakingProxy.pause();
    }

    function testUnpauseRevertWhenNotPaused() public {
        vm.prank(safeGovernanceMultisig);
        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.ExpectedPause.selector));
        ipTokenStakingProxy.unpause();
    }

    // -------------------------------------------------------------------------
    // Admin config setters are intentionally NOT gated by pause: they stay
    // callable by DEFAULT_ADMIN_ROLE (timelock) while the contract is paused.
    // -------------------------------------------------------------------------

    function testAdminSettersWorkWhilePaused() public {
        _pauseAsGov();
        vm.startPrank(address(timelock));
        ipTokenStakingProxy.setMinStakeAmount(2 ether);
        ipTokenStakingProxy.setFee(1000000000000000002);
        vm.stopPrank();
        assertEq(ipTokenStakingProxy.minStakeAmount(), 2 ether);
        assertEq(ipTokenStakingProxy.fee(), 1000000000000000002);
    }

    // -------------------------------------------------------------------------
    // Gated entry points work again after unpause (resume). setWithdrawalAddress
    // is fee-charged but has no pubkey validation, so a correct call succeeds.
    // -------------------------------------------------------------------------

    function testGatedFunctionWorksAfterUnpause() public {
        vm.startPrank(safeGovernanceMultisig);
        ipTokenStakingProxy.pause();
        ipTokenStakingProxy.unpause();
        vm.stopPrank();
        assertEq(ipTokenStakingProxy.paused(), false);

        uint256 currentFee = ipTokenStakingProxy.fee();
        vm.deal(address(2), currentFee);
        vm.prank(address(2));
        ipTokenStakingProxy.setWithdrawalAddress{ value: currentFee }(address(3));
    }

    function testStorage() public {
        assertEq(ipTokenStakingProxy.MAX_MONIKER_LENGTH(), monikerLengthBefore);
        assertEq(ipTokenStakingProxy.STAKE_ROUNDING(), stakeRoundingBefore);
        assertEq(ipTokenStakingProxy.DEFAULT_MIN_FEE(), defaultMinFeeBefore);
        assertEq(ipTokenStakingProxy.MAX_DATA_LENGTH(), maxDataLengthBefore);
        assertEq(ipTokenStakingProxy.minCommissionRate(), minCommissionRateBefore);
        assertEq(ipTokenStakingProxy.minStakeAmount(), minStakeAmountBefore);
        assertEq(ipTokenStakingProxy.minUnstakeAmount(), minUnstakeAmountBefore);
        assertEq(ipTokenStakingProxy.fee(), feeBefore);
    }
}
