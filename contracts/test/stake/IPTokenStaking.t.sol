// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { IPTokenStaking, IIPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { Secp256k1 } from "../../src/libraries/Secp256k1.sol";

import { Test } from "../utils/Test.sol";

contract IPTokenStakingTest is Test {
    address admin = address(0x123);
    bytes private delegatorUncmpPubkey =
        hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b"; // pragma: allowlist-secret
    // Address matching delegatorUncmpPubkey
    bytes private delegatorCmpPubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced1"; // pragma: allowlist-secret
    // Address matching delegatorCmpPubkey
    address private delegatorAddr = address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab);

    event Received(address, uint256);

    // For some tests, we need to receive the native token to this contract
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }

    function setUp() public override {
        setStaking();

        vm.assertEq(delegatorCmpPubkey, Secp256k1.compressPublicKey(delegatorUncmpPubkey));
    }

    function testIPTokenStaking_Constructor() public {
        // vm.expectRevert("IPTokenStaking: minStakeAmount cannot be 0");
        // new IPTokenStaking(
        //     0, // minStakeAmount
        //     1 ether, // minUnstakeAmount
        //     1 ether, // minRedelegateAmount
        //     1 gwei, // stakingRounding
        //     7 days, // withdrawalAddressChangeInterval
        //     1000, // defaultCommissionRate, 10%
        //     5000, // defaultMaxCommissionRate, 50%
        //     500 // defaultMaxCommissionChangeRate, 5%
        // );
        // vm.expectRevert("IPTokenStaking: minUnstakeAmount cannot be 0");
        // new IPTokenStaking(
        //     1 ether, // minStakeAmount
        //     0, // minUnstakeAmount
        //     1 ether, // minRedelegateAmount
        //     1 gwei, // stakingRounding
        //     7 days, // withdrawalAddressChangeInterval
        //     1000, // defaultCommissionRate, 10%
        //     5000, // defaultMaxCommissionRate, 50%
        //     500 // defaultMaxCommissionChangeRate, 5%
        // );
        // vm.expectRevert("IPTokenStaking: minRedelegateAmount cannot be 0");
        // new IPTokenStaking(
        //     1 ether, // minStakeAmount
        //     1 ether, // minUnstakeAmount
        //     0, // minRedelegateAmount
        //     1 gwei, // stakingRounding
        //     7 days, // withdrawalAddressChangeInterval
        //     1000, // defaultCommissionRate, 10%
        //     5000, // defaultMaxCommissionRate, 50%
        //     500 // defaultMaxCommissionChangeRate, 5%
        // );
        vm.expectRevert();
        new IPTokenStaking(
            0, // stakingRounding
            1000, // defaultCommissionRate, 10%
            5000, // defaultMaxCommissionRate, 50%
            500 // defaultMaxCommissionChangeRate, 5%
        );
        // vm.expectRevert("IPTokenStaking: newWithdrawalAddressChangeInterval cannot be 0");
        // new IPTokenStaking(
        //     1 ether, // minStakeAmount
        //     1 ether, // minUnstakeAmount
        //     1 ether, // minRedelegateAmount
        //     1 gwei, // stakingRounding
        //     0, // withdrawalAddressChangeInterval
        //     1000, // defaultCommissionRate, 10%
        //     5000, // defaultMaxCommissionRate, 50%
        //     500 // defaultMaxCommissionChangeRate, 5%
        // );
        vm.expectRevert("IPTokenStaking: Invalid default commission rate");
        new IPTokenStaking(
            1 gwei, // stakingRounding
            10_001, // defaultCommissionRate, 10%
            5000, // defaultMaxCommissionRate, 50%
            500 // defaultMaxCommissionChangeRate, 5%
        );
        vm.expectRevert("IPTokenStaking: Invalid default max commission rate");
        new IPTokenStaking(
            1 gwei, // stakingRounding
            1000, // defaultCommissionRate, 10%
            10_001, // defaultMaxCommissionRate, 50%
            500 // defaultMaxCommissionChangeRate, 5%
        );
        vm.expectRevert("IPTokenStaking: Invalid default max commission rate");
        new IPTokenStaking(
            1 gwei, // stakingRounding
            1000, // defaultCommissionRate, 10%
            1, // defaultMaxCommissionRate, 50%
            500 // defaultMaxCommissionChangeRate, 5%
        );
        vm.expectRevert("IPTokenStaking: Invalid default max commission change rate");
        new IPTokenStaking(
            1 gwei, // stakingRounding
            1000, // defaultCommissionRate, 10%
            5000, // defaultMaxCommissionRate, 50%
            10_001 // defaultMaxCommissionChangeRate, 5%
        );
    }

    function testIPTokenStaking_Parameters() public view {
        assertEq(ipTokenStaking.minStakeAmount(), 1 ether);
        assertEq(ipTokenStaking.minUnstakeAmount(), 1 ether);
        assertEq(ipTokenStaking.minRedelegateAmount(), 1 ether);
        assertEq(ipTokenStaking.STAKE_ROUNDING(), 1 gwei);
        assertEq(ipTokenStaking.withdrawalAddressChangeInterval(), 7 days);
        assertEq(ipTokenStaking.DEFAULT_COMMISSION_RATE(), 1000);
        assertEq(ipTokenStaking.DEFAULT_MAX_COMMISSION_RATE(), 5000);
        assertEq(ipTokenStaking.DEFAULT_MAX_COMMISSION_CHANGE_RATE(), 500);
    }

    function testIPTokenStaking_CreateValidator() public {
        // Network shall not allow anyone to create a new validator if the validator account’s balance is 0.
        // Note that this restriction doesn’t apply to validator creation on behalf.
        uint256 stakeAmount = 0 ether;
        bytes memory validatorUncmpPubkey = delegatorUncmpPubkey;
        bytes memory validatorCmpPubkey = delegatorCmpPubkey;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount too low");
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: validatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100
        });
        // Check that no stakes are put on the validator
        assertEq(ipTokenStaking.delegatorTotalStakes(validatorCmpPubkey), 0 ether);
        assertEq(ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorCmpPubkey), 0 ether);

        // Network shall not allow anyone to create a new validator on behalf if the sender account’s balance is 0.
        bytes memory validator1Pubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9111111"; // pragma: allowlist-secret
        stakeAmount = 0 ether;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount too low");
        ipTokenStaking.createValidatorOnBehalf{ value: stakeAmount }({ validatorPubkey: validator1Pubkey });
        // Check that no stakes are put on the validator
        assertEq(ipTokenStaking.delegatorTotalStakes(validator1Pubkey), 0 ether);
        assertEq(ipTokenStaking.delegatorValidatorStakes(validator1Pubkey, validator1Pubkey), 0 ether);

        // Network shall allow anyone to create a new validator by staking validator’s own tokens (self-delegation)
        stakeAmount = 1 ether;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(delegatorCmpPubkey, "delegator's validator", stakeAmount, 1000, 5000, 100);
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: delegatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100
        });
        // Check that stakes are correctly put on the validator
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), 1 ether);
        assertEq(ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, delegatorCmpPubkey), 1 ether);

        // NOTE: We have removed the validator existence check in createValidator, thus this test is not valid anymore.
        // // Adding a validator twice should not be allowed
        // vm.deal(delegatorAddr, stakeAmount);
        // vm.prank(delegatorAddr);
        // vm.expectRevert("IPTokenStaking: Validator already exists");
        // ipTokenStaking.createValidator{ value: stakeAmount }({
        //     validatorUncmpPubkey: delegatorUncmpPubkey,
        //     moniker: "delegator's validator",
        //     commissionRate: 1000,
        //     maxCommissionRate: 5000,
        //     maxCommissionChangeRate: 100
        // });

        // Network shall allow anyone to create a new validator on behalf of a validator.
        // Note that the operation stakes sender’s tokens to the validator, and the delegator will still be the validator itself.
        bytes memory validator2Pubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9222222"; // pragma: allowlist-secret
        stakeAmount = 1000 ether;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(validator2Pubkey, "validator", stakeAmount, 1000, 5000, 500);
        ipTokenStaking.createValidatorOnBehalf{ value: stakeAmount }({ validatorPubkey: validator2Pubkey });
        // Check that stakes are correctly put on the validator
        assertEq(ipTokenStaking.delegatorTotalStakes(validator2Pubkey), 1000 ether);
        // Check that the delegator is the validator itself
        assertEq(ipTokenStaking.delegatorValidatorStakes(validator2Pubkey, validator2Pubkey), 1000 ether);
        assertEq(ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validator2Pubkey), 0 ether);

        // NOTE: We have removed the validator existence check in createValidator, thus this test is not valid anymore.
        // // Network shall not allow anyone to create a new validator with existing validators’ public keys.
        // bytes memory validator3Pubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9222222"; // pragma: allowlist-secret
        // stakeAmount = 1 ether;
        // vm.deal(delegatorAddr, stakeAmount);
        // vm.prank(delegatorAddr);
        // vm.expectRevert("IPTokenStaking: Validator already exists");
        // ipTokenStaking.createValidatorOnBehalf{ value: stakeAmount }({ validatorPubkey: validator3Pubkey });
        // // Check that stakes are changing for the existing validator
        // assertEq(ipTokenStaking.delegatorTotalStakes(validator3Pubkey), 1000 ether);
        // assertEq(ipTokenStaking.delegatorValidatorStakes(validator3Pubkey, validator3Pubkey), 1000 ether);

        // Network shall not allow anyone to create a new validator if the provided public key doesn’t match sender’s address.
        bytes
            memory delegatorUncmpPubkeyChanged = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe222"; // pragma: allowlist-secret
        stakeAmount = 1 ether;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid pubkey derived address");
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: delegatorUncmpPubkeyChanged,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100
        });
        // Check that no stakes are put on the validator
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorUncmpPubkeyChanged), 0 ether);
        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorUncmpPubkeyChanged, delegatorUncmpPubkeyChanged),
            0 ether
        );
    }

    function testIPTokenStaking_CreateValidator_MultipleTimes() public {
        // When creating an existing validator (second time), it should emit CreateValidator event with existing values but updated stake amount.
        uint256 stakeAmount = ipTokenStaking.minStakeAmount();
        bytes memory validatorUncmpPubkey = delegatorUncmpPubkey;
        bytes memory validatorCmpPubkey = delegatorCmpPubkey;
        vm.deal(delegatorAddr, stakeAmount * 2);

        uint256 beforeDelegatorTotalStakes = ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey);
        uint256 beforeDelegatorValidatorStakes = ipTokenStaking.delegatorValidatorStakes(
            delegatorCmpPubkey,
            validatorCmpPubkey
        );

        // Create initially
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(validatorCmpPubkey, "delegator's validator", stakeAmount, 1000, 5000, 100);
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: validatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100
        });

        // Check that more stakes are put on the validator
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), beforeDelegatorTotalStakes + stakeAmount);
        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorCmpPubkey),
            beforeDelegatorValidatorStakes + stakeAmount
        );

        // Create again
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(validatorCmpPubkey, "delegator's validator", stakeAmount, 1000, 5000, 100);
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: validatorUncmpPubkey,
            moniker: "bad name validator",
            commissionRate: 100,
            maxCommissionRate: 100,
            maxCommissionChangeRate: 100
        });

        // Check that more stakes are put on the validator
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), beforeDelegatorTotalStakes + stakeAmount * 2);
        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorCmpPubkey),
            beforeDelegatorValidatorStakes + stakeAmount * 2
        );
    }

    modifier withDefaultValidator() {
        vm.deal(delegatorAddr, 1 ether);
        vm.prank(delegatorAddr);
        ipTokenStaking.createValidator{ value: 1 ether }({
            validatorUncmpPubkey: delegatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100
        });
        _;
    }

    function testIPTokenStaking_Stake() public withDefaultValidator {
        // Network shall allow anyone to deposit stake ≥ minimum stake amount into an existing validator for a delegator pubkey.
        bytes memory validatorPubkey = delegatorCmpPubkey;
        uint256 stakeAmount = 1 ether;
        vm.deal(delegatorAddr, stakeAmount);

        uint256 delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey);
        uint256 delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey);
        (
            bool isActive,
            string memory moniker,
            uint256 totalStake,
            uint32 commissionRate,
            uint32 maxCommissionRate,
            uint32 maxCommissionChangeRate
        ) = ipTokenStaking.validatorMetadata(validatorPubkey);
        uint256 validatorTotalBefore = totalStake;

        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(delegatorUncmpPubkey, validatorPubkey);

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey),
            delegatorValidatorBefore + stakeAmount
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), delegatorTotalBefore + stakeAmount);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore + stakeAmount);

        // (TODO) Network shall refund money to the staker in EL if execution in EL succeeds but CL fails

        // Network shall allow anyone to stake on behalf of another delegator.
        validatorPubkey = delegatorCmpPubkey;
        bytes memory delegator1Pubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9111111"; // pragma: allowlist secret
        stakeAmount = 1000 ether;
        vm.deal(delegatorAddr, stakeAmount);

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegator1Pubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.prank(delegatorAddr);
        ipTokenStaking.stakeOnBehalf{ value: stakeAmount }(delegator1Pubkey, validatorPubkey);

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey),
            delegatorValidatorBefore + stakeAmount
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegator1Pubkey), delegatorTotalBefore + stakeAmount);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore + stakeAmount);

        // (TODO) Network shall prevent depositing stake into a validator pubkey that has not been created (stake).

        // Network shall prevent depositing stake into a validator pubkey that has not been created (stakeOnBehalf).
        validatorPubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9777777"; // pragma: allowlist secret
        delegator1Pubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9111111"; // pragma: allowlist secret
        stakeAmount = 1000 ether;
        vm.deal(delegatorAddr, stakeAmount);

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegator1Pubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Validator does not exist");
        ipTokenStaking.stakeOnBehalf{ value: stakeAmount }(delegator1Pubkey, validatorPubkey);

        assertEq(ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey), delegatorValidatorBefore);
        assertEq(ipTokenStaking.delegatorTotalStakes(delegator1Pubkey), delegatorTotalBefore);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore);

        // Network shall not allow anyone to deposit stake < minimum stake amount
        validatorPubkey = delegatorCmpPubkey;
        delegator1Pubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9111111"; // pragma: allowlist secret
        stakeAmount = 100 gwei;
        vm.deal(delegatorAddr, stakeAmount);

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegator1Pubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount too low");
        ipTokenStaking.stakeOnBehalf{ value: stakeAmount }(delegator1Pubkey, validatorPubkey);

        assertEq(ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey), delegatorValidatorBefore);
        assertEq(ipTokenStaking.delegatorTotalStakes(delegator1Pubkey), delegatorTotalBefore);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore);

        // Network shall round the input stake amount by 1 gwei and send the remainder back to the sender.
        validatorPubkey = delegatorCmpPubkey;
        delegator1Pubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9111111"; // pragma: allowlist secret
        stakeAmount = 1_000_000_000_000_000_001 wei;
        vm.deal(delegatorAddr, stakeAmount);

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegator1Pubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.prank(delegatorAddr);
        ipTokenStaking.stakeOnBehalf{ value: stakeAmount }(delegator1Pubkey, validatorPubkey);

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegator1Pubkey, validatorPubkey),
            delegatorValidatorBefore + stakeAmount - 1 wei
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegator1Pubkey), delegatorTotalBefore + stakeAmount - 1 wei);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore + stakeAmount - 1 wei);
        // (TODO) Check that sender got 1 wei back.
    }

    function testIPTokenStaking_Redelegate() public withDefaultValidator {
        // Network shall allow the delegators to move their staked token from source validator to destination validator.
        bytes memory validatorSrcPubkey = delegatorCmpPubkey;

        uint256 stakeAmount = 5 ether;

        vm.deal(delegatorAddr, stakeAmount + 1 gwei);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(delegatorUncmpPubkey, validatorSrcPubkey);

        // last character modified from e => f
        bytes
            memory validatorDstUncmpPubkey = hex"03eb8e065336169de70e591e397b76600a71b356c9c3c629a8d0987e2169588e5b64d5f0c60f03ec8f5b13ba133b0a8e0f03bbaa8e678c0d03bb9dab42626be04f"; // pragma: allowlist-secret
        bytes memory validatorDstCmpPubkey = Secp256k1.compressPublicKey(validatorDstUncmpPubkey);

        // Create the new validator
        ipTokenStaking.createValidatorOnBehalf{ value: 1 gwei }(validatorDstCmpPubkey);

        uint256 srcValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorSrcPubkey);
        uint256 dstValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorDstCmpPubkey);

        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Redelegate(delegatorCmpPubkey, validatorSrcPubkey, validatorDstCmpPubkey, stakeAmount);
        ipTokenStaking.redelegate(
            IIPTokenStaking.RedelegateParams({
                delegatorUncmpPubkey: delegatorUncmpPubkey,
                validatorSrcPubkey: validatorSrcPubkey,
                validatorDstPubkey: validatorDstCmpPubkey,
                amount: stakeAmount
            })
        );

        // Check the amount for the source and destination validator
        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorSrcPubkey),
            srcValidatorBefore - stakeAmount
        );
        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorDstCmpPubkey),
            dstValidatorBefore + stakeAmount
        );

        // Network shall not allow non-operators of a stake owner to redelegate from the stake owner’s public key
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        validatorSrcPubkey = delegatorCmpPubkey;
        stakeAmount = 5 ether;
        vm.deal(delegatorAddr, stakeAmount + 1 gwei);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(delegatorUncmpPubkey, validatorSrcPubkey);

        validatorDstUncmpPubkey = hex"03eb8e065336169de70e591e397b76600a71b356c9c3c629a8d0987e2169588e5b64d5f0c60f03ec8f5b13ba133b0a8e0f03bbaa8e678c0d03bb9dab42626be04f"; // pragma: allowlist-secret
        validatorDstCmpPubkey = Secp256k1.compressPublicKey(validatorDstUncmpPubkey);

        srcValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorSrcPubkey);
        dstValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorDstCmpPubkey);

        vm.prank(operator);
        vm.expectRevert("IPTokenStaking: Caller is not an operator");
        ipTokenStaking.redelegateOnBehalf(
            IIPTokenStaking.RedelegateParams({
                delegatorUncmpPubkey: delegatorUncmpPubkey,
                validatorSrcPubkey: validatorSrcPubkey,
                validatorDstPubkey: validatorDstCmpPubkey,
                amount: stakeAmount
            })
        );

        // Network shall allow operators of a stake owner to redelegate from the stake owner’s public key
        vm.prank(delegatorAddr);
        ipTokenStaking.addOperator(delegatorUncmpPubkey, operator);
        validatorSrcPubkey = delegatorCmpPubkey;
        stakeAmount = 5 ether;
        vm.deal(delegatorAddr, stakeAmount + 1 gwei);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(delegatorUncmpPubkey, validatorSrcPubkey);

        validatorDstUncmpPubkey = hex"03eb8e065336169de70e591e397b76600a71b356c9c3c629a8d0987e2169588e5b64d5f0c60f03ec8f5b13ba133b0a8e0f03bbaa8e678c0d03bb9dab42626be04f"; // pragma: allowlist-secret
        validatorDstCmpPubkey = Secp256k1.compressPublicKey(validatorDstUncmpPubkey);

        srcValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorSrcPubkey);
        dstValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorDstCmpPubkey);

        vm.prank(operator);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Redelegate(delegatorCmpPubkey, validatorSrcPubkey, validatorDstCmpPubkey, stakeAmount);
        ipTokenStaking.redelegateOnBehalf(
            IIPTokenStaking.RedelegateParams({
                delegatorUncmpPubkey: delegatorUncmpPubkey,
                validatorSrcPubkey: validatorSrcPubkey,
                validatorDstPubkey: validatorDstCmpPubkey,
                amount: stakeAmount
            })
        );

        // Check the amount for the source and destination validator
        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorSrcPubkey),
            srcValidatorBefore - stakeAmount
        );
        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorDstCmpPubkey),
            dstValidatorBefore + stakeAmount
        );

        // Network shall not allow anyone to redelegate from non-existing-validator
        validatorSrcPubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9888888"; // pragma: allowlist secret
        stakeAmount = 5 ether;

        validatorDstUncmpPubkey = hex"03eb8e065336169de70e591e397b76600a71b356c9c3c629a8d0987e2169588e5b64d5f0c60f03ec8f5b13ba133b0a8e0f03bbaa8e678c0d03bb9dab42626be04f"; // pragma: allowlist-secret
        validatorDstCmpPubkey = Secp256k1.compressPublicKey(validatorDstUncmpPubkey);

        srcValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorSrcPubkey);
        dstValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorDstCmpPubkey);

        vm.prank(operator);
        vm.expectRevert("IPTokenStaking: Validator does not exist");
        ipTokenStaking.redelegateOnBehalf(
            IIPTokenStaking.RedelegateParams({
                delegatorUncmpPubkey: delegatorUncmpPubkey,
                validatorSrcPubkey: validatorSrcPubkey,
                validatorDstPubkey: validatorDstCmpPubkey,
                amount: stakeAmount
            })
        );

        // Network shall not allow anyone to redelegate to non-existing-validator
        validatorSrcPubkey = delegatorCmpPubkey;
        stakeAmount = 5 ether;
        vm.deal(delegatorAddr, stakeAmount + 1 gwei);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(delegatorUncmpPubkey, validatorSrcPubkey);

        validatorDstUncmpPubkey = hex"03eb8e065336169de70e591e397b76600a71b356c9c3c629a8d0987e2169588e5b64d5f0c60f03ec8f5b13ba133b0a8e0f03bbaa8e678c0d03bb9dab4262000000"; // pragma: allowlist-secret
        validatorDstCmpPubkey = Secp256k1.compressPublicKey(validatorDstUncmpPubkey);

        vm.prank(operator);
        vm.expectRevert("IPTokenStaking: Validator does not exist");
        ipTokenStaking.redelegateOnBehalf(
            IIPTokenStaking.RedelegateParams({
                delegatorUncmpPubkey: delegatorUncmpPubkey,
                validatorSrcPubkey: validatorSrcPubkey,
                validatorDstPubkey: validatorDstCmpPubkey,
                amount: stakeAmount
            })
        );

        // Network shall not allow operators or stake owners to redelegate more than the delegator staked on the source validator
        validatorSrcPubkey = delegatorCmpPubkey;
        stakeAmount = 5 ether;
        vm.deal(delegatorAddr, stakeAmount + 1 gwei);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(delegatorUncmpPubkey, validatorSrcPubkey);

        validatorDstUncmpPubkey = hex"03eb8e065336169de70e591e397b76600a71b356c9c3c629a8d0987e2169588e5b64d5f0c60f03ec8f5b13ba133b0a8e0f03bbaa8e678c0d03bb9dab42626be04f"; // pragma: allowlist-secret
        validatorDstCmpPubkey = Secp256k1.compressPublicKey(validatorDstUncmpPubkey);

        srcValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorSrcPubkey);
        dstValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorDstCmpPubkey);

        vm.prank(operator);
        vm.expectRevert("IPTokenStaking: Insufficient staked amount");
        ipTokenStaking.redelegateOnBehalf(
            IIPTokenStaking.RedelegateParams({
                delegatorUncmpPubkey: delegatorUncmpPubkey,
                validatorSrcPubkey: validatorSrcPubkey,
                validatorDstPubkey: validatorDstCmpPubkey,
                amount: stakeAmount + 100 ether
            })
        );
    }

    function testIPTokenStaking_Unstake() public withDefaultValidator {
        bytes memory validatorPubkey = delegatorCmpPubkey;

        vm.deal(delegatorAddr, 100 ether);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: 50 ether }(delegatorUncmpPubkey, validatorPubkey);

        // Network shall only allow the stake owner to withdraw from their stake pubkey
        uint256 stakeAmount = 1 ether;

        uint256 delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey);
        uint256 delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey);
        (
            bool isActive,
            string memory moniker,
            uint256 totalStake,
            uint32 commissionRate,
            uint32 maxCommissionRate,
            uint32 maxCommissionChangeRate
        ) = ipTokenStaking.validatorMetadata(validatorPubkey);
        uint256 validatorTotalBefore = totalStake;

        vm.warp(vm.getBlockTimestamp() + ipTokenStaking.withdrawalAddressChangeInterval() + 1);

        vm.startPrank(delegatorAddr);
        ipTokenStaking.setWithdrawalAddress(delegatorUncmpPubkey, address(0xb0b));
        ipTokenStaking.unstake(delegatorUncmpPubkey, validatorPubkey, stakeAmount);
        vm.stopPrank();

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey),
            delegatorValidatorBefore - stakeAmount
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), delegatorTotalBefore - stakeAmount);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore - stakeAmount);

        // Network shall not allow non-operators of a stake owner to withdraw from the stake owner’s public key
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        stakeAmount = 1 ether;

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.startPrank(operator);
        vm.expectRevert("IPTokenStaking: Caller is not an operator");
        ipTokenStaking.unstakeOnBehalf(delegatorCmpPubkey, validatorPubkey, stakeAmount);
        vm.stopPrank();

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey),
            delegatorValidatorBefore
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), delegatorTotalBefore);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore);

        // Network shall allow operators of a stake owner to withdraw from the stake owner’s public key
        vm.prank(delegatorAddr);
        ipTokenStaking.addOperator(delegatorUncmpPubkey, operator);
        stakeAmount = 1 ether;

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.startPrank(operator);
        ipTokenStaking.unstakeOnBehalf(delegatorCmpPubkey, validatorPubkey, stakeAmount);
        vm.stopPrank();

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey),
            delegatorValidatorBefore - stakeAmount
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), delegatorTotalBefore - stakeAmount);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore - stakeAmount);

        // Network shall not allow operators or stake owners to withdraw more than the delegator staked on the validator
        stakeAmount = 100 ether;

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.startPrank(operator);
        vm.expectRevert("IPTokenStaking: Insufficient staked amount");
        ipTokenStaking.unstakeOnBehalf(delegatorCmpPubkey, validatorPubkey, stakeAmount);
        vm.stopPrank();

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey),
            delegatorValidatorBefore
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), delegatorTotalBefore);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore);

        // Network shall not allow anyone to withdraw from stake on non-validators’ public keys
        validatorPubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa9888888"; // pragma: allowlist secret
        stakeAmount = 1 ether;

        delegatorValidatorBefore = ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey);
        delegatorTotalBefore = ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        validatorTotalBefore = totalStake;

        vm.startPrank(operator);
        vm.expectRevert("IPTokenStaking: Validator does not exist");
        ipTokenStaking.unstakeOnBehalf(delegatorCmpPubkey, validatorPubkey, stakeAmount);
        vm.stopPrank();

        assertEq(
            ipTokenStaking.delegatorValidatorStakes(delegatorCmpPubkey, validatorPubkey),
            delegatorValidatorBefore
        );
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), delegatorTotalBefore);
        (isActive, moniker, totalStake, commissionRate, maxCommissionRate, maxCommissionChangeRate) = ipTokenStaking
            .validatorMetadata(validatorPubkey);
        assertEq(totalStake, validatorTotalBefore);
    }

    function testIPTokenStaking_SetWithdrawalAddress() public withDefaultValidator {
        bytes memory validatorPubkey = delegatorCmpPubkey;

        vm.deal(delegatorAddr, 50 ether);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: 50 ether }(delegatorUncmpPubkey, validatorPubkey);

        // Network shall allow the delegators to set their withdrawal address
        vm.warp(vm.getBlockTimestamp() + ipTokenStaking.withdrawalAddressChangeInterval() + 1);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.SetWithdrawalAddress(
            delegatorCmpPubkey,
            0x0000000000000000000000000000000000000000000000000000000000000b0b
        );
        vm.prank(delegatorAddr);
        ipTokenStaking.setWithdrawalAddress(delegatorUncmpPubkey, address(0xb0b));
        assertEq(ipTokenStaking.withdrawalAddressChange(delegatorCmpPubkey), vm.getBlockTimestamp());

        // Network shall not allow anyone to set withdrawal address for other delegators
        bytes
            memory delegatorUncmpPubkey1 = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist secret
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid pubkey derived address");
        ipTokenStaking.setWithdrawalAddress(delegatorUncmpPubkey1, address(0xb0b));

        // Network shall not allow anyone to set withdrawal address if cooldown period has not passed
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Withdrawal address change cool-down");
        ipTokenStaking.setWithdrawalAddress(delegatorUncmpPubkey, address(0xb0b));

        // Network shall not allow anyone to set withdrawal address for 0-stake delegators
        vm.prank(delegatorAddr);
        ipTokenStaking.unstake(delegatorUncmpPubkey, validatorPubkey, 51 ether);
        assertEq(ipTokenStaking.delegatorTotalStakes(delegatorCmpPubkey), 0 ether);

        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Delegator must have stake");
        ipTokenStaking.setWithdrawalAddress(delegatorUncmpPubkey, address(0xb0b));
    }

    function testIPTokenStaking_addOperator() public {
        // Network shall not allow others to add operators for a delegator
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        bytes
            memory otherDelegatorUncmpPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist secret
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid pubkey derived address");
        ipTokenStaking.addOperator(otherDelegatorUncmpPubkey, operator);
    }

    function isInArray(address[] memory array, address element) internal pure returns (bool) {
        for (uint256 i = 0; i < array.length; i++) {
            if (array[i] == element) {
                return true;
            }
        }
        return false;
    }

    function testIPTokenStaking_removeOperator() public {
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.prank(delegatorAddr);
        ipTokenStaking.addOperator(delegatorUncmpPubkey, operator);
        assert(isInArray(ipTokenStaking.getOperators(delegatorCmpPubkey), operator));

        // Network shall not allow others to remove operators for a delegator
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenStaking: Invalid pubkey derived address");
        ipTokenStaking.removeOperator(delegatorUncmpPubkey, operator);
        assert(isInArray(ipTokenStaking.getOperators(delegatorCmpPubkey), operator));

        // Network shall allow delegators to remove their operators
        vm.prank(delegatorAddr);
        ipTokenStaking.removeOperator(delegatorUncmpPubkey, operator);
        assert(!isInArray(ipTokenStaking.getOperators(delegatorCmpPubkey), operator));

        // Removing an operator that does not exist reverts
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Operator not found");
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
        vm.expectRevert("IPTokenStaking: minStakeAmount cannot be 0");
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
        vm.expectRevert("IPTokenStaking: minUnstakeAmount cannot be 0");
        ipTokenStaking.setMinUnstakeAmount(0 ether);

        // Set using a non-owner address
        vm.prank(delegatorAddr);
        vm.expectRevert();
        ipTokenStaking.setMinUnstakeAmount(1 ether);
    }

    function testIPTokenStaking_setMinRedelegateAmount() public {
        // Set amount that will be rounded down to 0
        vm.prank(admin);
        ipTokenStaking.setMinRedelegateAmount(5 wei);
        assertEq(ipTokenStaking.minRedelegateAmount(), 0);

        // Set amount that will not be rounded
        vm.prank(admin);
        ipTokenStaking.setMinRedelegateAmount(1 ether);
        assertEq(ipTokenStaking.minRedelegateAmount(), 1 ether);

        // Set 0
        vm.prank(admin);
        vm.expectRevert("IPTokenStaking: minRedelegateAmount cannot be 0");
        ipTokenStaking.setMinRedelegateAmount(0 ether);

        // Set using a non-owner address
        vm.prank(delegatorAddr);
        vm.expectRevert();
        ipTokenStaking.setMinRedelegateAmount(1 ether);
    }
}
