// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { IPTokenStaking, IIPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { Secp256k1 } from "../../src/libraries/Secp256k1.sol";

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

        address impl = address(
            new IPTokenStaking(
                0, // stakingRounding
                1000, // defaultCommissionRate, 10%
                5000, // defaultMaxCommissionRate, 50%
                500 // defaultMaxCommissionChangeRate, 5%
            )
        );
        IPTokenStaking staking = IPTokenStaking(address(new ERC1967Proxy(impl, "")));

        IIPTokenStaking.InitializerArgs memory args = IIPTokenStaking.InitializerArgs({
            accessManager: admin,
            minStakeAmount: 0,
            minUnstakeAmount: 1 ether,
            withdrawalAddressChangeInterval: 7 days,
            shortStakingPeriod: 1,
            mediumStakingPeriod: 2,
            longStakingPeriod: 3
        });

        // IPTokenStaking: minStakeAmount cannot be 0
        vm.expectRevert();
        staking.initialize(args);

        // IPTokenStaking: minUnstakeAmount cannot be 0
        vm.expectRevert();
        args.minStakeAmount = 1 ether;
        args.minUnstakeAmount = 0;
        staking.initialize(args);

        // IPTokenStaking: newWithdrawalAddressChangeInterval cannot be 0
        vm.expectRevert();
        args.minUnstakeAmount = 1 ether;
        args.withdrawalAddressChangeInterval = 0;
        staking.initialize(args);

        // TODO test short
        // TODO test medium
        // TODO test long
    }

    function testIPTokenStaking_Parameters() public view {
        assertEq(ipTokenStaking.minStakeAmount(), 1 ether);
        assertEq(ipTokenStaking.minUnstakeAmount(), 1 ether);
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
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount too low");
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: validatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            isLocked: false,
	        data: ""

        });

        // Network shall not allow anyone to create a new validator on behalf if the sender account’s balance is 0.
        bytes
            memory validator1Pubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1000000"; // pragma: allowlist-secret
        stakeAmount = 0 ether;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount too low");
        ipTokenStaking.createValidatorOnBehalf{ value: stakeAmount }({ validatorUncmpPubkey: validator1Pubkey, isLocked: false, data: "" });

        // Network shall allow anyone to create a new validator by staking validator’s own tokens (self-delegation)
        stakeAmount = 1 ether;
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
            true,
            abi.encode("data")
        );
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorUncmpPubkey: delegatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            isLocked: true,
            data: abi.encode("data")
        });

        // Network shall allow anyone to create a new validator on behalf of a validator.
        // Note that the operation stakes sender’s tokens to the validator, and the delegator will still be the validator itself.
        bytes
            memory validator2UncmpPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe222"; // pragma: allowlist-secret
        stakeAmount = 1000 ether;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(
            validator2UncmpPubkey,
            "validator",
            stakeAmount,
            1000,
            5000,
            500,
            false,
            abi.encode("data")
        );
        ipTokenStaking.createValidatorOnBehalf{ value: stakeAmount }({
            validatorUncmpPubkey: validator2UncmpPubkey,
            isLocked: false,
            data: ""
        });

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
            maxCommissionChangeRate: 100,
            isLocked: false,
            data: ""
        });
    }

    modifier withDefaultValidator() {
        vm.deal(delegatorAddr, 1 ether);
        vm.prank(delegatorAddr);
        ipTokenStaking.createValidator{ value: 1 ether }({
            validatorUncmpPubkey: delegatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            isLocked: false,
            data: ""
        });
        _;
    }

    function testIPTokenStaking_Unstake_Flexible() public withDefaultValidator {
        bytes memory validatorPubkey = delegatorUncmpPubkey;
        IIPTokenStaking.StakingPeriod stkPeriod = IIPTokenStaking.StakingPeriod.FLEXIBLE;
        vm.deal(delegatorAddr, 100 ether);
        vm.prank(delegatorAddr);
        uint256 delegationId = ipTokenStaking.stake{ value: 50 ether }(delegatorUncmpPubkey, validatorPubkey, stkPeriod, "");

        assertEq(delegationId, 0);
        // Network shall only allow the stake owner to withdraw from their stake pubkey
        uint256 stakeAmount = 1 ether;

        vm.warp(vm.getBlockTimestamp() + ipTokenStaking.withdrawalAddressChangeInterval() + 1);

        vm.startPrank(delegatorAddr);
        ipTokenStaking.setWithdrawalAddress(delegatorUncmpPubkey, address(0xb0b));
        ipTokenStaking.unstake(delegatorUncmpPubkey, validatorPubkey, stakeAmount, delegationId, "");
        vm.stopPrank();


        // Network shall not allow non-operators of a stake owner to withdraw from the stake owner’s public key
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        stakeAmount = 1 ether;

        vm.startPrank(operator);
        vm.expectRevert("IPTokenStaking: Caller is not an operator");
        ipTokenStaking.unstakeOnBehalf(delegatorUncmpPubkey, validatorPubkey, stakeAmount, delegationId, "");
        vm.stopPrank();

        // Network shall allow operators of a stake owner to withdraw from the stake owner’s public key
        vm.prank(delegatorAddr);
        ipTokenStaking.addOperator(delegatorUncmpPubkey, operator);
        stakeAmount = 1 ether;

        vm.startPrank(operator);
        ipTokenStaking.unstakeOnBehalf(delegatorUncmpPubkey, validatorPubkey, stakeAmount, delegationId, "");
        vm.stopPrank();
    }

    function testIPTokenStaking_SetWithdrawalAddress() public withDefaultValidator {
        bytes memory validatorPubkey = delegatorUncmpPubkey;
        IIPTokenStaking.StakingPeriod stkPeriod = IIPTokenStaking.StakingPeriod.FLEXIBLE;

        vm.deal(delegatorAddr, 50 ether);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: 50 ether }(delegatorUncmpPubkey, validatorPubkey, stkPeriod, "");

        // Network shall allow the delegators to set their withdrawal address
        vm.warp(vm.getBlockTimestamp() + ipTokenStaking.withdrawalAddressChangeInterval() + 1);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.SetWithdrawalAddress(
            delegatorUncmpPubkey,
            0x0000000000000000000000000000000000000000000000000000000000000b0b
        );
        vm.prank(delegatorAddr);
        ipTokenStaking.setWithdrawalAddress(delegatorUncmpPubkey, address(0xb0b));
        assertEq(ipTokenStaking.withdrawalAddressChange(delegatorUncmpPubkey), vm.getBlockTimestamp());

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
        assert(isInArray(ipTokenStaking.getOperators(delegatorUncmpPubkey), operator));

        // Network shall not allow others to remove operators for a delegator
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenStaking: Invalid pubkey derived address");
        ipTokenStaking.removeOperator(delegatorUncmpPubkey, operator);
        assert(isInArray(ipTokenStaking.getOperators(delegatorUncmpPubkey), operator));

        // Network shall allow delegators to remove their operators
        vm.prank(delegatorAddr);
        ipTokenStaking.removeOperator(delegatorUncmpPubkey, operator);
        assert(!isInArray(ipTokenStaking.getOperators(delegatorUncmpPubkey), operator));

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

    function testIPTokenStaking_Unjail() withDefaultValidator public {
        uint256 feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);

        // Network shall not allow anyone to unjail a validator if it is not the validator itself.
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Invalid pubkey derived address");
        ipTokenStaking.unjail(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenStaking.unjail(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
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
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenStaking.unjail{ value: feeAmount }(delegatorUncmpPubkey, "");
    }

    function testIPTokenStaking_UnjailOnBehalf() withDefaultValidator public {
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);

        // Network shall not allow anyone to unjail an non-existing validator.
        uint256 feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);

        // Network shall not allow anyone to unjail with compressed pubkey of incorrect length.
        bytes memory delegatorCmpPubkeyShortLen = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ce"; // pragma: allowlist secret
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Invalid pubkey length");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkeyShortLen, "");

        // Network shall not allow anyone to unjail with compressed pubkey of incorrect prefix.
        bytes
            memory delegatorCmpPubkeyWrongPrefix = hex"05e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced1"; // pragma: allowlist secret
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Invalid pubkey prefix");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkeyWrongPrefix, "");

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenStaking.unjailOnBehalf(delegatorUncmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
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
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(delegatorUncmpPubkey, "");
    }

    function testIPTokenStaking_SetUnjailFee() public {
        // Network shall allow the owner to set the unjail fee.
        uint256 newUnjailFee = 2 ether;
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.UnjailFeeSet(newUnjailFee);
        vm.prank(admin);
        ipTokenStaking.setUnjailFee(newUnjailFee);
        assertEq(ipTokenStaking.unjailFee(), newUnjailFee);

        // Network shall not allow non-owner to set the unjail fee.
        vm.prank(address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA));
        vm.expectRevert();
        ipTokenStaking.setUnjailFee(1 ether);
        assertEq(ipTokenStaking.unjailFee(), newUnjailFee);
    }
}
