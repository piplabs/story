// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { IPTokenStaking, IIPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { Errors } from "../../src/libraries/Errors.sol";
import { Test } from "../utils/Test.sol";

contract IPTokenStakingTest is Test {
    bytes private delegatorUncmpPubkey =
        hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b"; // pragma: allowlist-secret
    // Address matching delegatorCmpPubkey
    address private delegatorAddr = address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab);

    event Received(address, uint256);

    // For some tests, we need to receive the native token to this contract
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }

    function setUp() public virtual override {
        super.setUp();
    }

    function testIPTokenStaking_Constructor() public {
        vm.expectRevert(Errors.IPTokenStaking__InvalidDefaultMinFee.selector);
        new IPTokenStaking(
            1 gwei, // stakingRounding
            0 ether
        );
        vm.expectRevert(Errors.IPTokenStaking__ZeroStakingRounding.selector);
        address impl = address(
            new IPTokenStaking(
                0, // stakingRounding
                1 ether // Default min fee charged for adding to CL storage, 1 eth
            )
        );

        IIPTokenStaking.InitializerArgs memory args = IIPTokenStaking.InitializerArgs({
            owner: admin,
            minStakeAmount: 0,
            minUnstakeAmount: 1 ether,
            minCommissionRate: 5_00,
            shortStakingPeriod: 1,
            mediumStakingPeriod: 2,
            longStakingPeriod: 3,
            fee: 1 ether
        });
        impl = address(
            new IPTokenStaking(
                1 gwei, // stakingRounding
                1 ether // Default min fee charged for adding to CL storage, 1 eth
            )
        );
        // IPTokenStaking: minStakeAmount cannot be 0
        vm.expectRevert(Errors.IPTokenStaking__ZeroMinStakeAmount.selector);
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        // IPTokenStaking: minUnstakeAmount cannot be 0
        vm.expectRevert(Errors.IPTokenStaking__ZeroMinUnstakeAmount.selector);
        args.minStakeAmount = 1 ether;
        args.minUnstakeAmount = 0;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        // IPTokenStaking: newWithdrawalAddressChangeInterval cannot be 0
        vm.expectRevert(Errors.IPTokenStaking__ZeroMinCommissionRate.selector);
        args.minUnstakeAmount = 1 ether;
        args.minCommissionRate = 0;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        // TODO test short
        vm.expectRevert(Errors.IPTokenStaking__ZeroShortPeriodDuration.selector);
        args.minCommissionRate = 5_00;
        args.shortStakingPeriod = 0;
        args.mediumStakingPeriod = 10;
        args.longStakingPeriod = 100;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        vm.expectRevert(Errors.IPTokenStaking__ShortPeriodLongerThanMedium.selector);
        args.shortStakingPeriod = 1;
        args.mediumStakingPeriod = 1;
        args.longStakingPeriod = 100;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        vm.expectRevert(Errors.IPTokenStaking__ShortPeriodLongerThanMedium.selector);
        args.shortStakingPeriod = 2;
        args.mediumStakingPeriod = 1;
        args.longStakingPeriod = 100;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        vm.expectRevert(Errors.IPTokenStaking__MediumLongerThanLong.selector);
        args.shortStakingPeriod = 2;
        args.mediumStakingPeriod = 100;
        args.longStakingPeriod = 100;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        vm.expectRevert(); // todo
        args.shortStakingPeriod = 2;
        args.mediumStakingPeriod = 3;
        args.longStakingPeriod = 2;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        vm.expectRevert(); // todo
        args.shortStakingPeriod = 1;
        args.mediumStakingPeriod = 2;
        args.longStakingPeriod = 3;
        args.fee = 10;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));
    }

    function testIPTokenStaking_Parameters() public view {
        assertEq(ipTokenStaking.minStakeAmount(), 1024 ether);
        assertEq(ipTokenStaking.minUnstakeAmount(), 1024 ether);
        assertEq(ipTokenStaking.STAKE_ROUNDING(), 1 gwei);
        assertEq(ipTokenStaking.minCommissionRate(), 5_00);
        assertEq(ipTokenStaking.DEFAULT_MIN_FEE(), 1 ether);
    }

    function testIPTokenStaking_CreateValidator() public {
        uint256 stakeAmount = 0.5 ether;
        bytes memory validatorUncmpPubkey = delegatorUncmpPubkey;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__StakeAmountUnderMin.selector);
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: validatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: false,
            data: ""
        });

        // Network shall not allow anyone to create a new validator on behalf if the msg.value < min
        bytes
            memory validator1Pubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist-secret
        stakeAmount = 0.5 ether;
        vm.deal(delegatorAddr, 1 ether);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__StakeAmountUnderMin.selector);
        ipTokenStaking.createValidatorOnBehalf{ value: stakeAmount }({
            validatorUncmpPubkey: validator1Pubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: false,
            data: ""
        });

        // Network shall allow anyone to create a new validator by staking validator’s own tokens (self-delegation)
        stakeAmount = ipTokenStaking.minStakeAmount();
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(
            validatorUncmpPubkey,
            "delegator's validator",
            stakeAmount,
            1000,
            5000,
            100,
            1, // supportsUnlocked
            delegatorAddr,
            abi.encode("data")
        );
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: delegatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: true,
            data: abi.encode("data")
        });

        // Network shall allow anyone to create a new validator on behalf of a validator.
        // Note that the operation stakes sender’s tokens to the validator, and the delegator will still be the validator itself.
        bytes
            memory validator2UncmpPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe222"; // pragma: allowlist-secret
        stakeAmount = ipTokenStaking.minStakeAmount();
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(
            validator2UncmpPubkey,
            "delegator's validator",
            stakeAmount,
            1000,
            5000,
            100,
            0, // supportsUnlocked
            delegatorAddr,
            abi.encode("data")
        );
        ipTokenStaking.createValidatorOnBehalf{ value: stakeAmount }({
            validatorUncmpPubkey: validator2UncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: false,
            data: abi.encode("data")
        });

        // Network shall not allow anyone to create a new validator if the provided public key doesn’t match sender’s address.
        bytes
            memory delegatorUncmpPubkeyChanged = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe222"; // pragma: allowlist-secret
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: delegatorUncmpPubkeyChanged,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: false,
            data: ""
        });
    }

    function testIPTokenStaking_Stake_Periods() public {
        // Flexible should produce 0 delegationId
        bytes memory validatorPubkey = delegatorUncmpPubkey;
        IIPTokenStaking.StakingPeriod stkPeriod = IIPTokenStaking.StakingPeriod.FLEXIBLE;
        vm.deal(delegatorAddr, 10000 ether);
        vm.prank(delegatorAddr);
        uint256 delegationId = ipTokenStaking.stake{ value: 1024 ether }(
            delegatorUncmpPubkey,
            validatorPubkey,
            stkPeriod,
            ""
        );
        assertEq(delegationId, 0);
        // Staking for short period should produce incremented delegationId and correct duration
        // emitted event
        uint32 shortDuration = ipTokenStaking.stakingDurations(IIPTokenStaking.StakingPeriod.SHORT);
        uint256 stakeAmount = ipTokenStaking.minUnstakeAmount();
        uint256 expectedDelegationId = 1;
        vm.deal(delegatorAddr, 10000000000 ether);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Deposit(
            delegatorUncmpPubkey,
            validatorPubkey,
            stakeAmount,
            shortDuration,
            expectedDelegationId,
            delegatorAddr,
            ""
        );
        delegationId = ipTokenStaking.stake{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorPubkey,
            IIPTokenStaking.StakingPeriod.SHORT,
            ""
        );
        assertEq(delegationId, expectedDelegationId);
        expectedDelegationId++;
        // Staking for medium period should produce incremented delegationId and correct duration
        // emitted event
        uint32 mediumDuration = ipTokenStaking.stakingDurations(IIPTokenStaking.StakingPeriod.MEDIUM);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Deposit(
            delegatorUncmpPubkey,
            validatorPubkey,
            stakeAmount,
            mediumDuration,
            expectedDelegationId,
            delegatorAddr,
            ""
        );
        delegationId = ipTokenStaking.stake{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorPubkey,
            IIPTokenStaking.StakingPeriod.MEDIUM,
            ""
        );
        assertEq(delegationId, expectedDelegationId);
        expectedDelegationId++;
        // Staking for long period should produce incremented delegationId and correct duration
        // emitted event
        uint32 longDuration = ipTokenStaking.stakingDurations(IIPTokenStaking.StakingPeriod.LONG);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Deposit(
            delegatorUncmpPubkey,
            validatorPubkey,
            stakeAmount,
            longDuration,
            expectedDelegationId,
            delegatorAddr,
            ""
        );
        delegationId = ipTokenStaking.stake{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorPubkey,
            IIPTokenStaking.StakingPeriod.LONG,
            ""
        );
        assertEq(delegationId, expectedDelegationId);
    }

    function testIPTokenStaking_stake_remainder() public {
        // No remainder if the stake amount has no values under STAKE_ROUNDING
        bytes memory validatorPubkey = delegatorUncmpPubkey;
        uint256 stakeAmount = 1024 ether;
        uint256 predeployInitialBalance = 1; // 1 wei, needed to have predeploy at genesis

        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorPubkey,
            IIPTokenStaking.StakingPeriod.FLEXIBLE,
            ""
        );
        assertEq(
            address(ipTokenStaking).balance,
            predeployInitialBalance,
            "IPTokenStaking: Stake amount should be burned"
        );
        assertEq(address(delegatorAddr).balance, 0, "Delegator: No remainder should be sent back");

        // Remainder if the stake amount has values under STAKE_ROUNDING
        stakeAmount = 1024 ether + 1 wei;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorPubkey,
            IIPTokenStaking.StakingPeriod.FLEXIBLE,
            ""
        );
        assertEq(address(ipTokenStaking).balance, predeployInitialBalance);
        assertEq(address(delegatorAddr).balance, 1 wei);
    }

    function testIPTokenStaking_Unstake_Flexible() public {
        bytes memory validatorPubkey = delegatorUncmpPubkey;

        // Network shall only allow the stake owner to withdraw from their stake pubkey
        uint256 stakeAmount = ipTokenStaking.minUnstakeAmount();
        uint256 delegationId = 1337;
        // Use VM setStorage to set the counter to delegationId + 1
        vm.store(
            address(ipTokenStaking),
            bytes32(uint256(3)), // _delegationIdCounter
            bytes32(uint256(1338))
        );

        vm.startPrank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Withdraw(
            delegatorUncmpPubkey,
            validatorPubkey,
            stakeAmount,
            delegationId,
            delegatorAddr,
            ""
        );
        ipTokenStaking.unstake(delegatorUncmpPubkey, validatorPubkey, delegationId, stakeAmount, "");
        vm.stopPrank();

        vm.startPrank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__LowUnstakeAmount.selector);
        ipTokenStaking.unstake(delegatorUncmpPubkey, validatorPubkey, delegationId, stakeAmount - 1, "");
        vm.stopPrank();

        // Smart contract allows non-operators of a stake owner to withdraw from the stake owner’s public key,
        // but this operation will fail in CL. Testing the event here
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);

        vm.startPrank(operator);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Withdraw(delegatorUncmpPubkey, validatorPubkey, stakeAmount, delegationId, operator, "");
        ipTokenStaking.unstakeOnBehalf(delegatorUncmpPubkey, validatorPubkey, delegationId, stakeAmount, "");
        vm.stopPrank();

        // Revert if delegationId is invalid
        delegationId++;
        vm.startPrank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidDelegationId.selector);
        ipTokenStaking.unstake(delegatorUncmpPubkey, validatorPubkey, delegationId + 2, stakeAmount, "");
        vm.stopPrank();
    }

    function testIPTokenStaking_Redelegation() public {
        uint256 stakeAmount = ipTokenStaking.minStakeAmount();
        uint256 delegationId = 1;
        // Use VM setStorage to set the counter to delegationId == 1
        vm.store(
            address(ipTokenStaking),
            bytes32(uint256(3)), // _delegationIdCounter
            bytes32(uint256(1))
        );
        bytes
            memory validatorUncmpSrcPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe222"; // pragma: allowlist-secret
        bytes
            memory validatorUncmpDstPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist-secret

        vm.expectEmit(true, true, true, true);
        emit IIPTokenStaking.Redelegate(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            delegatorAddr,
            stakeAmount
        );
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        ipTokenStaking.redelegate{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );

        // Can only be called by delegator
        vm.deal(address(0x4545), stakeAmount);
        vm.prank(address(0x4545));
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.redelegate{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );

        // Redelegating to same validator
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__RedelegatingToSameValidator.selector);
        ipTokenStaking.redelegate{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpSrcPubkey,
            delegationId,
            stakeAmount
        );
        // Malformed Src
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyLength.selector);
        ipTokenStaking.redelegate{ value: stakeAmount }(
            delegatorUncmpPubkey,
            hex"04e38d15ae6cc5d41cce27a2307903cb", // pragma: allowlist secret
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );
        // Malformed Dst
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyLength.selector);
        ipTokenStaking.redelegate{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            hex"04e38d15ae6cc5d41cce27a2307903cb", // pragma: allowlist secret
            delegationId,
            stakeAmount
        );
        // Stake < Min
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__StakeAmountUnderMin.selector);
        ipTokenStaking.redelegate{ value: stakeAmount - 1 }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount + 100
        );

        // Revert if delegationId is invalid
        delegationId++;
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidDelegationId.selector);
        ipTokenStaking.redelegate{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );
    }

    function testIPTokenStaking_RedelegationOnBehalf() public {
        uint256 stakeAmount = ipTokenStaking.minStakeAmount();
        uint256 delegationId = 1;
        // Use VM setStorage to set the counter to delegationId == 1
        vm.store(
            address(ipTokenStaking),
            bytes32(uint256(3)), // _delegationIdCounter
            bytes32(uint256(1))
        );
        bytes
            memory validatorUncmpSrcPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe222"; // pragma: allowlist-secret
        bytes
            memory validatorUncmpDstPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist-secret

        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);

        vm.expectEmit(true, true, true, true);
        emit IIPTokenStaking.Redelegate(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            operator,
            stakeAmount
        );
        vm.deal(operator, stakeAmount);
        vm.prank(operator);
        ipTokenStaking.redelegateOnBehalf{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );

        // Redelegating to same validator
        vm.deal(operator, stakeAmount);
        vm.prank(operator);
        vm.expectRevert(Errors.IPTokenStaking__RedelegatingToSameValidator.selector);
        ipTokenStaking.redelegateOnBehalf{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpSrcPubkey,
            delegationId,
            stakeAmount
        );
        // Malformed Src
        vm.deal(operator, stakeAmount);
        vm.prank(operator);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyLength.selector);
        ipTokenStaking.redelegateOnBehalf{ value: stakeAmount }(
            delegatorUncmpPubkey,
            hex"04e38d15ae6cc5d41cce27a2307903cb", // pragma: allowlist secret
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );
        // Malformed Dst
        vm.deal(operator, stakeAmount);
        vm.prank(operator);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyLength.selector);
        ipTokenStaking.redelegateOnBehalf{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            hex"04e38d15ae6cc5d41cce27a2307903cb", // pragma: allowlist secret
            delegationId,
            stakeAmount
        );
        // Stake < Min
        vm.deal(operator, stakeAmount);
        vm.prank(operator);
        vm.expectRevert(Errors.IPTokenStaking__StakeAmountUnderMin.selector);
        ipTokenStaking.redelegateOnBehalf{ value: stakeAmount - 1 }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount + 100
        );

        // Revert if delegationId is invalid
        delegationId++;
        vm.prank(operator);
        vm.expectRevert(Errors.IPTokenStaking__InvalidDelegationId.selector);
        ipTokenStaking.redelegateOnBehalf{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );
    }

    function testIPTokenStaking_SetWithdrawalAddress() public {
        uint256 feeAmount = ipTokenStaking.fee();
        // Network shall allow the delegators to set their withdrawal address
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.SetWithdrawalAddress(
            delegatorUncmpPubkey,
            0x0000000000000000000000000000000000000000000000000000000000000b0b
        );
        vm.prank(delegatorAddr);
        ipTokenStaking.setWithdrawalAddress{ value: feeAmount }(delegatorUncmpPubkey, address(0xb0b));

        // Network shall not allow anyone to set withdrawal address for other delegators
        bytes
            memory delegatorUncmpPubkey1 = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist secret
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.setWithdrawalAddress{ value: feeAmount }(delegatorUncmpPubkey1, address(0xb0b));

        // Network shall not allow anyone to set withdrawal address with insufficient fee
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.setWithdrawalAddress{ value: feeAmount - 1 }(delegatorUncmpPubkey, address(0xb0b));
    }

    function testIPTokenStaking_SetRewardsAddress() public {
        uint256 feeAmount = ipTokenStaking.fee();
        // Network shall allow the delegators to set their withdrawal address
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.SetRewardAddress(
            delegatorUncmpPubkey,
            0x0000000000000000000000000000000000000000000000000000000000000b0b
        );
        vm.prank(delegatorAddr);
        ipTokenStaking.setRewardsAddress{ value: feeAmount }(delegatorUncmpPubkey, address(0xb0b));

        // Network shall not allow anyone to set withdrawal address for other delegators
        bytes
            memory delegatorUncmpPubkey1 = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist secret
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.setRewardsAddress{ value: feeAmount }(delegatorUncmpPubkey1, address(0xb0b));

        // Network shall not allow anyone to set withdrawal address with insufficient fee
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.setRewardsAddress{ value: feeAmount - 1 }(delegatorUncmpPubkey, address(0xb0b));
    }

    function testIPTokenStaking_updateValidatorCommission() public {
        uint32 commissionRate = 100000000;
        uint256 feeAmount = ipTokenStaking.fee();
        vm.deal(delegatorAddr, feeAmount * 10);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.UpdateValidatorCommssion(delegatorUncmpPubkey, commissionRate);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(delegatorUncmpPubkey, commissionRate);

        // Network shall not allow anyone to update the commission rate of a validator if it is not the validator itself.
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.deal(otherAddress, feeAmount * 10);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(delegatorUncmpPubkey, commissionRate);

        // Network shall not allow anyone to update the commission rate of a validator if it is less than minCommissionRate.
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__CommissionRateUnderMin.selector);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(delegatorUncmpPubkey, 0);

        // Network shall not allow anyone to update the commission rate of a validator if the fee is not paid.
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount - 1 }(delegatorUncmpPubkey, commissionRate);
    }

    function testIPTokenStaking_updateValidatorCommission() public {
        uint32 commissionRate = 100000000;
        uint256 feeAmount = ipTokenStaking.validatorUpdateFee();
        vm.deal(delegatorAddr, feeAmount * 10);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.UpdateValidatorCommssion(delegatorUncmpPubkey, commissionRate);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(delegatorUncmpPubkey, commissionRate);

        // Network shall not allow anyone to update the commission rate of a validator if it is not the validator itself.
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.deal(otherAddress, feeAmount * 10);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(delegatorUncmpPubkey, commissionRate);

        // Network shall not allow anyone to update the commission rate of a validator if it is less than minCommissionRate.
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__CommissionRateUnderMin.selector);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(delegatorUncmpPubkey, 0);

        // Network shall not allow anyone to update the commission rate of a validator if the fee is not paid.
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount - 1 }(delegatorUncmpPubkey, commissionRate);
    }

    function testIPTokenStaking_addOperator() public {
        // Network shall not allow others to add operators for a delegator
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        bytes
            memory otherDelegatorUncmpPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist secret
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.addOperator(otherDelegatorUncmpPubkey, operator);

        // Network shall not allow anyone to add operators for a delegator if the fee is not paid.
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.addOperator(delegatorUncmpPubkey, operator);

        // Network shall not allow anyone to add operators for a delegator if the fee is wrong
        uint256 feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount + 1);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.addOperator{ value: feeAmount - 1 }(delegatorUncmpPubkey, operator);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.addOperator{ value: feeAmount + 1 }(delegatorUncmpPubkey, operator);

        // Network should allow delegators to add operators for themselves
        feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.AddOperator(delegatorUncmpPubkey, operator);
        ipTokenStaking.addOperator{ value: feeAmount }(delegatorUncmpPubkey, operator);
    }

    function testIPTokenStaking_removeOperator() public {
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        // Network shall not allow others to remove operators for a delegator
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.removeOperator(delegatorUncmpPubkey, operator);

        // Network shall allow delegators to remove their operators
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.RemoveOperator(delegatorUncmpPubkey, operator);
        ipTokenStaking.removeOperator(delegatorUncmpPubkey, operator);
    }

    function testIPTokenStaking_setMinStakeAmount() public {
        // Set amount that will be rounded down to 0
        vm.prank(admin);
        ipTokenStaking.setMinStakeAmount(5 wei);
        assertEq(ipTokenStaking.minStakeAmount(), 0);

        // Set amount that will not be rounded
        vm.prank(admin);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.MinStakeAmountSet(1 ether);
        ipTokenStaking.setMinStakeAmount(1 ether);
        assertEq(ipTokenStaking.minStakeAmount(), 1 ether);

        // Set 0
        vm.prank(admin);
        vm.expectRevert(Errors.IPTokenStaking__ZeroMinStakeAmount.selector);
        ipTokenStaking.setMinStakeAmount(0 ether);

        // Set using a non-owner address
        vm.prank(delegatorAddr);
        vm.expectRevert();
        ipTokenStaking.setMinStakeAmount(1 ether);
    }

    function testIPTokenStaking_setMinUnstakeAmount() public {
        // Set amount that will be rounded down to 0
        vm.prank(admin);
        ipTokenStaking.setMinUnstakeAmount(5 wei);
        assertEq(ipTokenStaking.minUnstakeAmount(), 0);

        // Set amount that will not be rounded
        vm.prank(admin);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.MinUnstakeAmountSet(1 ether);
        ipTokenStaking.setMinUnstakeAmount(1 ether);
        assertEq(ipTokenStaking.minUnstakeAmount(), 1 ether);

        // Set 0
        vm.prank(admin);
        vm.expectRevert(Errors.IPTokenStaking__ZeroMinUnstakeAmount.selector);
        ipTokenStaking.setMinUnstakeAmount(0 ether);

        // Set using a non-owner address
        vm.prank(delegatorAddr);
        vm.expectRevert();
        ipTokenStaking.setMinUnstakeAmount(1 ether);
    }

    function testIPTokenStaking_Unjail() public {
        uint256 feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);

        // Network shall not allow anyone to unjail a validator if it is not the validator itself.
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyDerivedAddress.selector);
        ipTokenStaking.unjail(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.unjail(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.unjail{ value: feeAmount }(delegatorUncmpPubkey, "");

        // Network shall allow anyone to unjail a validator if the fee is paid.
        feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Unjail(delegatorAddr, delegatorUncmpPubkey, "");
        ipTokenStaking.unjail{ value: feeAmount }(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is over.
        feeAmount = 1.1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.unjail{ value: feeAmount }(delegatorUncmpPubkey, "");
    }

    function testIPTokenStaking_UnjailOnBehalf() public {
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);

        // Network shall not allow anyone to unjail an non-existing validator.
        uint256 feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);

        // Network shall not allow anyone to unjail with compressed pubkey of incorrect length.
        bytes memory delegatorCmpPubkeyShortLen = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ce"; // pragma: allowlist secret
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyLength.selector);
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkeyShortLen, "");

        // Network shall not allow anyone to unjail with compressed pubkey of incorrect prefix.
        bytes
            memory delegatorCmpPubkeyWrongPrefix = hex"05e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b"; // pragma: allowlist secret
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidPubkeyPrefix.selector);
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkeyWrongPrefix, "");

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.unjailOnBehalf(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorUncmpPubkey, "");

        // Network shall allow anyone to unjail a validator on behalf if the fee is paid.
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Unjail(otherAddress, delegatorUncmpPubkey, "");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is over.
        feeAmount = 1.1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert(Errors.IPTokenStaking__InvalidFeeAmount.selector);
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorUncmpPubkey, "");
    }

    function testIPTokenStaking_SetFee() public {
        // Network shall allow the owner to set the fee charged for adding to CL storage.
        uint256 newFee = 2 ether;
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.FeeSet(newFee);
        vm.prank(admin);
        ipTokenStaking.setFee(newFee);
        assertEq(ipTokenStaking.fee(), newFee);

        // Network shall not allow non-owner to set the fee charged for adding to CL storage.
        vm.prank(address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA));
        vm.expectRevert();
        ipTokenStaking.setFee(1 ether);
        assertEq(ipTokenStaking.fee(), newFee);

        // Network shall not allow fees < default
        vm.expectRevert();
        ipTokenStaking.setFee(1);
    }
}
