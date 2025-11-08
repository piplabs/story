// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { IPTokenStaking, IIPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { Test } from "../utils/Test.sol";

contract Staker {
    IPTokenStaking private ipTokenStaking;
    uint256 private fee;

    constructor(address ipTokenStakingAddr) {
        ipTokenStaking = IPTokenStaking(ipTokenStakingAddr);
        fee = ipTokenStaking.fee();
    }

    function stake(bytes calldata validatorCmpPubkey) external payable {
        ipTokenStaking.stake{ value: msg.value }(validatorCmpPubkey, IIPTokenStaking.StakingPeriod.FLEXIBLE, "");
    }

    function unstake(bytes calldata validatorCmpPubkey, uint256 amount) external {
        ipTokenStaking.unstake{ value: fee }(validatorCmpPubkey, 0, amount, "");
    }

    function redelegate(
        bytes calldata srcValidatorCmpPubkey,
        bytes calldata dstValidatorCmpPubkey,
        uint256 amount
    ) external {
        ipTokenStaking.redelegate{ value: fee }(srcValidatorCmpPubkey, dstValidatorCmpPubkey, 0, amount);
    }
}

contract IPTokenStakingTest is Test {
    address private delegatorAddr = address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab);

    bytes private validatorCmpPubkey = hex"0381513466154dfc0e702d8325d7a45284bc395f60b813230da299fab640c1eb08"; // pragma: allowlist-secret
    address private validatorAddr = address(0xb5cb887155446f69b5e4D11C30755108AC87e9cD);

    bytes private wrongValidatorCmpPubkey = hex"0342013466154dfc0e702d8325d7a45284bc395f60b813230da299fab640c1eb08"; // pragma: allowlist-secret

    bytes private otherValidatorCmpPubkey = hex"03dd68a0a4923a1b9321d39f01425f7b631066514cb2e5e1b5ed91e5c327d30c53"; // pragma: allowlist-secret
    // Address matching delegatorCmpPubkey
    address private otherValidatorAddr = address(0xf89D606F67a267E9dbCc813c8169988aB8aAeB5E);

    bytes private dataOverMaxLen;

    event Received(address, uint256);

    // For some tests, we need to receive the native token to this contract
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }

    function setUp() public virtual override {
        super.setUp();

        for (uint256 i = 0; i < ipTokenStaking.MAX_DATA_LENGTH() + 1; i++) {
            dataOverMaxLen = abi.encodePacked(dataOverMaxLen, "a");
        }
    }

    function testIPTokenStaking_Constructor() public {
        vm.expectRevert("IPTokenStaking: Invalid default min fee");
        new IPTokenStaking(0 ether, 0);

        address impl;
        IIPTokenStaking.InitializerArgs memory args = IIPTokenStaking.InitializerArgs({
            owner: admin,
            minCreateValidatorAmount: 1 ether,
            minStakeAmount: 0,
            minUnstakeAmount: 1 ether,
            minCommissionRate: 500,
            fee: 1 ether
        });
        impl = address(
            new IPTokenStaking(
                1 ether, // Default min fee charged for adding to CL storage, 1 eth
                256 // maxDataLength
            )
        );
        // IPTokenStaking: minStakeAmount cannot be 0
        vm.expectRevert("IPTokenStaking: Zero min stake amount");
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        // IPTokenStaking: minUnstakeAmount cannot be 0
        vm.expectRevert("IPTokenStaking: Zero min unstake amount");
        args.minStakeAmount = 1 ether;
        args.minUnstakeAmount = 0;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        // IPTokenStaking:   cannot be 0
        vm.expectRevert("IPTokenStaking: Zero min commission rate");
        args.minUnstakeAmount = 1 ether;
        args.minCommissionRate = 0;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));

        vm.expectRevert("IPTokenStaking: Invalid min fee");
        args.minCommissionRate = 10;
        args.fee = 0;
        new ERC1967Proxy(impl, abi.encodeCall(IPTokenStaking.initialize, (args)));
    }

    function testIPTokenStaking_Parameters() public view {
        assertEq(ipTokenStaking.minStakeAmount(), 1024 ether);
        assertEq(ipTokenStaking.minUnstakeAmount(), 1024 ether);
        assertEq(ipTokenStaking.STAKE_ROUNDING(), 1 gwei);
        assertEq(ipTokenStaking.minCommissionRate(), 500);
        assertEq(ipTokenStaking.DEFAULT_MIN_FEE(), 1 ether);
        assertEq(ipTokenStaking.MAX_DATA_LENGTH(), 256);
    }

    function testIPTokenStaking_CreateValidator() public {
        uint256 stakeAmount = 0.5 ether;
        vm.deal(validatorAddr, stakeAmount);
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount under min");
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorCmpPubkey: validatorCmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: false,
            data: ""
        });

        // Network shall allow anyone to create a new validator by staking validator’s own tokens (self-delegation)
        stakeAmount = ipTokenStaking.minStakeAmount();
        vm.deal(validatorAddr, stakeAmount);
        vm.prank(validatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.CreateValidator(
            validatorCmpPubkey,
            "delegator's validator",
            stakeAmount,
            1000,
            5000,
            100,
            1, // supportsUnlocked
            validatorAddr, // self-delegation, validatorAddr = delegatorAddr
            abi.encode("data")
        );
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorCmpPubkey: validatorCmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: true,
            data: abi.encode("data")
        });

        // Network shall not allow a moniker longer than MAX_MONIKER_LENGTH
        string memory moniker;
        for (uint256 i = 0; i < ipTokenStaking.MAX_MONIKER_LENGTH() + 1; i++) {
            moniker = string.concat(moniker, "a");
        }
        stakeAmount = ipTokenStaking.minStakeAmount();
        vm.deal(validatorAddr, stakeAmount);
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Moniker length over max");
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorCmpPubkey: validatorCmpPubkey,
            moniker: moniker,
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: false,
            data: ""
        });

        stakeAmount = ipTokenStaking.minStakeAmount();
        vm.deal(validatorAddr, stakeAmount);
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Data length over max");
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorCmpPubkey: validatorCmpPubkey,
            moniker: "",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: false,
            data: dataOverMaxLen
        });
    }

    function testIPTokenStaking_Stake_Periods() public {
        // Flexible should produce 0 delegationId
        IIPTokenStaking.StakingPeriod stkPeriod = IIPTokenStaking.StakingPeriod.FLEXIBLE;
        vm.deal(delegatorAddr, 10_000 ether);
        vm.prank(delegatorAddr);
        uint256 delegationId = ipTokenStaking.stake{ value: 1024 ether }(validatorCmpPubkey, stkPeriod, "");
        assertEq(delegationId, 0);
        // Staking for short period should produce incremented delegationId and correct duration
        // emitted event
        uint256 stakeAmount = ipTokenStaking.minUnstakeAmount();
        uint256 expectedDelegationId = 1;
        vm.deal(delegatorAddr, 10_000_000_000 ether);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Deposit(
            delegatorAddr,
            validatorCmpPubkey,
            stakeAmount,
            uint32(uint8(IIPTokenStaking.StakingPeriod.SHORT)),
            expectedDelegationId,
            delegatorAddr,
            ""
        );
        delegationId = ipTokenStaking.stake{ value: stakeAmount }(
            validatorCmpPubkey,
            IIPTokenStaking.StakingPeriod.SHORT,
            ""
        );
        assertEq(delegationId, expectedDelegationId);
        expectedDelegationId++;
        // Staking for medium period should produce incremented delegationId and correct duration
        // emitted event
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Deposit(
            delegatorAddr,
            validatorCmpPubkey,
            stakeAmount,
            uint32(uint8(IIPTokenStaking.StakingPeriod.MEDIUM)),
            expectedDelegationId,
            delegatorAddr,
            ""
        );
        delegationId = ipTokenStaking.stake{ value: stakeAmount }(
            validatorCmpPubkey,
            IIPTokenStaking.StakingPeriod.MEDIUM,
            ""
        );
        assertEq(delegationId, expectedDelegationId);
        expectedDelegationId++;
        // Staking for long period should produce incremented delegationId and correct duration
        // emitted event
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Deposit(
            delegatorAddr,
            validatorCmpPubkey,
            stakeAmount,
            uint32(uint8(IIPTokenStaking.StakingPeriod.LONG)),
            expectedDelegationId,
            delegatorAddr,
            ""
        );
        delegationId = ipTokenStaking.stake{ value: stakeAmount }(
            validatorCmpPubkey,
            IIPTokenStaking.StakingPeriod.LONG,
            ""
        );
        assertEq(delegationId, expectedDelegationId);

        // Test revert for invalid validatorCmpPubkey
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        bytes memory invalidValidatorCmpPubkey = hex"1234";
        ipTokenStaking.stake{ value: stakeAmount }(
            invalidValidatorCmpPubkey,
            IIPTokenStaking.StakingPeriod.FLEXIBLE,
            ""
        );
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.stake{ value: stakeAmount }(wrongValidatorCmpPubkey, IIPTokenStaking.StakingPeriod.FLEXIBLE, "");
    }

    function testIPTokenStaking_stake_remainder() public {
        // No remainder if the stake amount has no values under STAKE_ROUNDING
        uint256 stakeAmount = 1024 ether;
        uint256 predeployInitialBalance = 1; // 1 wei, needed to have predeploy at genesis

        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        ipTokenStaking.stake{ value: stakeAmount }(validatorCmpPubkey, IIPTokenStaking.StakingPeriod.FLEXIBLE, "data");
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
        ipTokenStaking.stake{ value: stakeAmount }(validatorCmpPubkey, IIPTokenStaking.StakingPeriod.FLEXIBLE, "data");
        assertEq(address(ipTokenStaking).balance, predeployInitialBalance);
        assertEq(address(delegatorAddr).balance, 1 wei);
    }

    function testIPTokenStaking_Stake_data() public {
        // Network shall not allow a data longer than MAX_DATA_LENGTH
        uint256 stakeAmount = 1024 ether;
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Data length over max");
        ipTokenStaking.stake{ value: stakeAmount }(
            validatorCmpPubkey,
            IIPTokenStaking.StakingPeriod.FLEXIBLE,
            dataOverMaxLen
        );

        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Data length over max");
        ipTokenStaking.stakeOnBehalf{ value: stakeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            IIPTokenStaking.StakingPeriod.FLEXIBLE,
            dataOverMaxLen
        );
    }

    function testIPTokenStaking_Unstake_Flexible() public {
        uint256 feeAmount = ipTokenStaking.fee();

        // Network shall only allow the stake owner to withdraw from their stake pubkey
        uint256 stakeAmount = ipTokenStaking.minUnstakeAmount();
        uint256 delegationId = 1337;
        // Use VM setStorage to set the counter to delegationId + 1
        vm.store(
            address(ipTokenStaking),
            bytes32(uint256(3)), // _delegationIdCounter
            bytes32(uint256(1338))
        );

        vm.deal(delegatorAddr, feeAmount);
        vm.startPrank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Withdraw(delegatorAddr, validatorCmpPubkey, stakeAmount, delegationId, delegatorAddr, "");
        ipTokenStaking.unstake{ value: feeAmount }(validatorCmpPubkey, delegationId, stakeAmount, "");
        vm.stopPrank();

        vm.deal(delegatorAddr, feeAmount);
        vm.startPrank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Unstake amount under min");
        ipTokenStaking.unstake{ value: feeAmount }(validatorCmpPubkey, delegationId, stakeAmount - 1, "");
        vm.stopPrank();

        // Smart contract allows non-operators of a stake owner to withdraw from the stake owner’s public key,
        // but this operation will fail in CL. Testing the event here
        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);

        vm.deal(operator, feeAmount);
        vm.startPrank(operator);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Withdraw(delegatorAddr, validatorCmpPubkey, stakeAmount, delegationId, operator, "");
        ipTokenStaking.unstakeOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            delegationId,
            stakeAmount,
            ""
        );
        vm.stopPrank();

        // Revert if delegationId is invalid
        delegationId++;
        vm.startPrank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid delegation id");
        ipTokenStaking.unstake{ value: feeAmount }(validatorCmpPubkey, delegationId + 2, stakeAmount, "");
        vm.stopPrank();

        // Revert if fee is not paid
        vm.startPrank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unstake{ value: feeAmount - 1 }(validatorCmpPubkey, delegationId + 2, stakeAmount, "");
        vm.stopPrank();

        // Round down to STAKE_ROUNDING if amount is not divisible by STAKE_ROUNDING
        uint256 unroundedAmount = ipTokenStaking.minUnstakeAmount() + ipTokenStaking.STAKE_ROUNDING() + 1 wei;
        uint256 expectedUnstakeAmount = unroundedAmount - 1 wei;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Withdraw(delegatorAddr, validatorCmpPubkey, expectedUnstakeAmount, delegationId, delegatorAddr, "");
        ipTokenStaking.unstake{ value: feeAmount }(validatorCmpPubkey, delegationId, unroundedAmount, "");

        unroundedAmount = 1024000000000999999999 wei;
        expectedUnstakeAmount = 1024 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Withdraw(delegatorAddr, validatorCmpPubkey, expectedUnstakeAmount, delegationId, delegatorAddr, "");
        ipTokenStaking.unstake{ value: feeAmount }(validatorCmpPubkey, delegationId, unroundedAmount, "");

        // Revert if validatorCmpPubkey is invalid
        bytes memory invalidValidatorCmpPubkey = hex"1234";
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        ipTokenStaking.unstake{ value: feeAmount }(invalidValidatorCmpPubkey, delegationId, stakeAmount, "");
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.unstake{ value: feeAmount }(wrongValidatorCmpPubkey, delegationId, stakeAmount, "");
    }

    function testIPTokenStaking_Unstake_data() public {
        uint256 feeAmount = ipTokenStaking.fee();
        uint256 stakeAmount = ipTokenStaking.minUnstakeAmount();
        uint256 delegationId = 1337;
        // Use VM setStorage to set the counter to delegationId + 1
        vm.store(
            address(ipTokenStaking),
            bytes32(uint256(3)), // _delegationIdCounter
            bytes32(uint256(1338))
        );

        // Network shall not allow a data longer than MAX_DATA_LENGTH
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Data length over max");
        ipTokenStaking.unstake{ value: feeAmount }(validatorCmpPubkey, delegationId, stakeAmount, dataOverMaxLen);

        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Data length over max");
        ipTokenStaking.unstakeOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            delegationId,
            stakeAmount,
            dataOverMaxLen
        );
    }

    function testIPTokenStaking_Redelegation() public {
        uint256 stakeAmount = ipTokenStaking.minStakeAmount();
        uint256 delegationId = 1;
        uint256 feeAmount = ipTokenStaking.fee();
        // Use VM setStorage to set the counter to delegationId == 1
        vm.store(
            address(ipTokenStaking),
            bytes32(uint256(3)), // _delegationIdCounter
            bytes32(uint256(1))
        );

        vm.expectEmit(true, true, true, true);
        emit IIPTokenStaking.Redelegate(
            delegatorAddr,
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            delegatorAddr,
            stakeAmount
        );
        vm.deal(delegatorAddr, stakeAmount + feeAmount);
        vm.prank(delegatorAddr);
        ipTokenStaking.redelegate{ value: feeAmount }(
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Redelegating to same validator
        vm.deal(delegatorAddr, stakeAmount + feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Redelegating to same validator");
        ipTokenStaking.redelegate{ value: feeAmount }(
            validatorCmpPubkey,
            validatorCmpPubkey,
            delegationId,
            stakeAmount
        );
        // Invalid source validator address
        bytes memory invalidValidatorCmpPubkey = hex"1234";
        vm.deal(delegatorAddr, stakeAmount + feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        ipTokenStaking.redelegate{ value: feeAmount }(
            invalidValidatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.redelegate{ value: feeAmount }(
            wrongValidatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Invalid destination validator address
        vm.deal(delegatorAddr, stakeAmount + feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        ipTokenStaking.redelegate{ value: feeAmount }(
            validatorCmpPubkey,
            invalidValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.redelegate{ value: feeAmount }(
            validatorCmpPubkey,
            wrongValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Revert if delegationId is invalid
        delegationId++;
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid delegation id");
        ipTokenStaking.redelegate{ value: feeAmount }(
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        delegationId--;

        // Revert if fee is not paid
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.redelegate{ value: feeAmount - 1 }(
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Stake < Min
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount under min");
        ipTokenStaking.redelegate{ value: feeAmount }(
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount - 1
        );
    }

    function testIPTokenStaking_RedelegationOnBehalf() public {
        uint256 stakeAmount = ipTokenStaking.minStakeAmount();
        uint256 delegationId = 1;
        uint256 feeAmount = ipTokenStaking.fee();
        // Use VM setStorage to set the counter to delegationId == 1
        vm.store(
            address(ipTokenStaking),
            bytes32(uint256(3)), // _delegationIdCounter
            bytes32(uint256(1))
        );

        address operator = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);

        vm.expectEmit(true, true, true, true);
        emit IIPTokenStaking.Redelegate(
            delegatorAddr,
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            operator,
            stakeAmount
        );
        vm.deal(operator, stakeAmount + feeAmount);
        vm.prank(operator);
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Redelegating to same validator
        vm.deal(operator, stakeAmount + feeAmount);
        vm.prank(operator);
        vm.expectRevert("IPTokenStaking: Redelegating to same validator");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            validatorCmpPubkey,
            delegationId,
            stakeAmount
        );
        bytes memory invalidValidatorCmpPubkey = hex"1234";
        // Invalid source validator
        vm.deal(operator, stakeAmount + feeAmount);
        vm.prank(operator);
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            invalidValidatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            wrongValidatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Invalid destination validator
        vm.deal(operator, stakeAmount + feeAmount);
        vm.prank(operator);
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            invalidValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            wrongValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Revert if delegationId is invalid
        delegationId++;
        vm.prank(operator);
        vm.expectRevert("IPTokenStaking: Invalid delegation id");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );
        delegationId--;

        // Revert if fee is not paid
        vm.prank(operator);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount - 1 }(
            delegatorAddr,
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount
        );

        // Stake < Min
        vm.deal(delegatorAddr, stakeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Stake amount under min");
        ipTokenStaking.redelegateOnBehalf{ value: feeAmount }(
            delegatorAddr,
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            delegationId,
            stakeAmount - 1
        );
    }

    function testIPTokenStaking_SetWithdrawalAddress() public {
        uint256 feeAmount = ipTokenStaking.fee();
        // Network shall allow the delegators to set their withdrawal address
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.SetWithdrawalAddress(
            delegatorAddr,
            0x0000000000000000000000000000000000000000000000000000000000000b0b
        );
        vm.prank(delegatorAddr);
        ipTokenStaking.setWithdrawalAddress{ value: feeAmount }(address(0xb0b));

        // Network shall not allow anyone to set withdrawal address with insufficient fee
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.setWithdrawalAddress{ value: feeAmount - 1 }(address(0xb0b));
    }

    function testIPTokenStaking_SetRewardsAddress() public {
        uint256 feeAmount = ipTokenStaking.fee();
        // Network shall allow the delegators to set their withdrawal address
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.SetRewardAddress(
            delegatorAddr,
            0x0000000000000000000000000000000000000000000000000000000000000b0b
        );
        vm.prank(delegatorAddr);
        ipTokenStaking.setRewardsAddress{ value: feeAmount }(address(0xb0b));

        // Network shall not allow anyone to set withdrawal address with insufficient fee
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.setRewardsAddress{ value: feeAmount - 1 }(address(0xb0b));
    }

    function testIPTokenStaking_updateValidatorCommission() public {
        uint32 commissionRate = 100_000_000;
        uint256 feeAmount = ipTokenStaking.fee();
        vm.deal(validatorAddr, feeAmount * 10);
        vm.prank(validatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.UpdateValidatorCommission(validatorCmpPubkey, commissionRate);
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(validatorCmpPubkey, commissionRate);

        // Network shall not allow anyone to update the commission rate of a validator if it is less than minCommissionRate.
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Commission rate under min");
        ipTokenStaking.updateValidatorCommission{ value: feeAmount }(validatorCmpPubkey, 0);

        // Network shall not allow anyone to update the commission rate of a validator if the fee is not paid.
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.updateValidatorCommission{ value: feeAmount - 1 }(validatorCmpPubkey, commissionRate);
    }

    function testIPTokenStaking_setOperator() public {
        // Network shall not allow anyone to add operators for a delegator if the fee is not paid.
        address operator = address(0x420);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.setOperator{ value: 0 }(operator);

        // Network shall not allow anyone to add operators for a delegator if the fee is wrong
        uint256 feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount + 1);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.setOperator{ value: feeAmount - 1 }(operator);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.setOperator{ value: feeAmount + 1 }(operator);

        // Network should allow delegators to add operators for themselves
        feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.SetOperator(delegatorAddr, operator);
        ipTokenStaking.setOperator{ value: feeAmount }(operator);
    }

    function testIPTokenStaking_unsetOperator() public {
        uint256 feeAmount = ipTokenStaking.fee();

        // Network shall allow delegators to remove their operators
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.UnsetOperator(delegatorAddr);
        ipTokenStaking.unsetOperator{ value: feeAmount }();

        // Revert if fee is not paid
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unsetOperator{ value: feeAmount - 1 }();
    }

    function testIPTokenStaking_setMinCreateValidatorAmount() public {
        // Set amount that will be rounded down to 1 ether
        performTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinCreateValidatorAmount.selector, 1 ether + 5 wei)
        );
        assertEq(ipTokenStaking.minCreateValidatorAmount(), 1 ether);

        // Set amount that will not be rounded
        schedule(address(ipTokenStaking), abi.encodeWithSelector(IPTokenStaking.setMinCreateValidatorAmount.selector, 1 ether));
        waitForTimelock();
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.MinCreateValidatorAmountSet(1 ether);
        executeTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinCreateValidatorAmount.selector, 1 ether)
        );
        assertEq(ipTokenStaking.minCreateValidatorAmount(), 1 ether);

        // Set 0
        expectRevertTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinCreateValidatorAmount.selector, 0 ether),
            "IPTokenStaking: Zero min create validator amount"
        );

        // Set amount that will be rounded down to 0
        expectRevertTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinCreateValidatorAmount.selector, 5 wei),
            "IPTokenStaking: Zero min create validator amount"
        );

        // Set using a non-owner address
        vm.prank(delegatorAddr);
        vm.expectRevert(
            abi.encodeWithSelector(
                Ownable.OwnableUnauthorizedAccount.selector,
                address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab)
            )
        );
        ipTokenStaking.setMinCreateValidatorAmount(1 ether);
    }

    function testIPTokenStaking_setMinStakeAmount() public {
        // Set amount that will be rounded down to 1 ether
        performTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinStakeAmount.selector, 1 ether + 5 wei)
        );
        assertEq(ipTokenStaking.minStakeAmount(), 1 ether);

        // Set amount that will not be rounded
        schedule(address(ipTokenStaking), abi.encodeWithSelector(IPTokenStaking.setMinStakeAmount.selector, 1 ether));
        waitForTimelock();
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.MinStakeAmountSet(1 ether);
        executeTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinStakeAmount.selector, 1 ether)
        );
        assertEq(ipTokenStaking.minStakeAmount(), 1 ether);

        // Set 0
        expectRevertTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinStakeAmount.selector, 0 ether),
            "IPTokenStaking: Zero min stake amount"
        );

        // Set amount that will be rounded down to 0
        expectRevertTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinStakeAmount.selector, 5 wei),
            "IPTokenStaking: Zero min stake amount"
        );

        // Set using a non-owner address
        vm.prank(delegatorAddr);
        vm.expectRevert(
            abi.encodeWithSelector(
                Ownable.OwnableUnauthorizedAccount.selector,
                address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab)
            )
        );
        ipTokenStaking.setMinStakeAmount(1 ether);
    }

    function testIPTokenStaking_setMinUnstakeAmount() public {
        // Set amount that will be rounded down to 1 ether
        performTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinUnstakeAmount.selector, 1 ether + 5 wei)
        );
        assertEq(ipTokenStaking.minUnstakeAmount(), 1 ether);

        // Set amount that will not be rounded
        schedule(address(ipTokenStaking), abi.encodeWithSelector(IPTokenStaking.setMinUnstakeAmount.selector, 1 ether));
        waitForTimelock();
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.MinUnstakeAmountSet(1 ether);
        executeTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinUnstakeAmount.selector, 1 ether)
        );
        assertEq(ipTokenStaking.minUnstakeAmount(), 1 ether);

        // Set 0
        vm.prank(admin);
        expectRevertTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinUnstakeAmount.selector, 0 ether),
            "IPTokenStaking: Zero min unstake amount"
        );

        // Set amount that will be rounded down to 0 ether
        vm.prank(admin);
        expectRevertTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setMinUnstakeAmount.selector, 5 wei),
            "IPTokenStaking: Zero min unstake amount"
        );

        // Set using a non-owner address
        vm.prank(delegatorAddr);
        vm.expectRevert();
        ipTokenStaking.setMinUnstakeAmount(1 ether);
    }

    function testIPTokenStaking_Unjail() public {
        uint256 feeAmount = 1 ether;
        vm.deal(validatorAddr, feeAmount);

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unjail(validatorCmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(validatorAddr, feeAmount);
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unjail{ value: feeAmount }(validatorCmpPubkey, "");

        // Network shall allow anyone to unjail a validator if the fee is paid.
        feeAmount = 1 ether;
        vm.deal(validatorAddr, feeAmount);
        vm.prank(validatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Unjail(validatorAddr, validatorCmpPubkey, "");
        ipTokenStaking.unjail{ value: feeAmount }(validatorCmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is over.
        feeAmount = 1.1 ether;
        vm.deal(validatorAddr, feeAmount);
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unjail{ value: feeAmount }(validatorCmpPubkey, "");

        // Network shall not allow anyone to unjail with invalid pubkey
        feeAmount = 1 ether;
        vm.deal(validatorAddr, feeAmount);
        vm.prank(validatorAddr);
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        ipTokenStaking.unjail{ value: feeAmount }(hex"", "");
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.unjail{ value: feeAmount }(wrongValidatorCmpPubkey, "");
    }

    function testIPTokenStaking_UnjailOnBehalf() public {
        uint256 feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unjailOnBehalf(validatorCmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(validatorCmpPubkey, "");

        // Network shall allow anyone to unjail a validator if the fee is paid.
        feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Unjail(delegatorAddr, validatorCmpPubkey, "");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(validatorCmpPubkey, "");

        // Network shall not allow anyone to unjail a validator if the fee is over.
        feeAmount = 1.1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Invalid fee amount");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(validatorCmpPubkey, "");

        // Network shall not allow anyone to unjail with invalid pubkey
        feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("Secp256k1Verifier: Invalid cmp pubkey length");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(hex"", "");
        vm.expectRevert("Secp256k1Verifier: pubkey not on curve");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(wrongValidatorCmpPubkey, "");
    }

    function testIPTokenStaking_Unjail_data() public {
        uint256 feeAmount = ipTokenStaking.DEFAULT_MIN_FEE();
        vm.deal(validatorAddr, feeAmount);
        vm.prank(validatorAddr);
        vm.expectRevert("IPTokenStaking: Data length over max");
        ipTokenStaking.unjail{ value: feeAmount }(validatorCmpPubkey, dataOverMaxLen);

        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenStaking: Data length over max");
        ipTokenStaking.unjailOnBehalf{ value: feeAmount }(validatorCmpPubkey, dataOverMaxLen);
    }

    function testIPTokenStaking_SetFee() public {
        // Network shall allow the owner to set the fee charged for adding to CL storage.
        uint256 newFee = 2 ether;
        schedule(address(ipTokenStaking), abi.encodeWithSelector(IPTokenStaking.setFee.selector, newFee));
        waitForTimelock();
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.FeeSet(newFee);
        executeTimelocked(address(ipTokenStaking), abi.encodeWithSelector(IPTokenStaking.setFee.selector, newFee));
        assertEq(ipTokenStaking.fee(), newFee);

        // Network shall not allow non-owner to set the fee charged for adding to CL storage.
        vm.prank(address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA));
        vm.expectRevert();
        ipTokenStaking.setFee(1 ether);
        assertEq(ipTokenStaking.fee(), newFee);

        // Network shall not allow fees < default
        expectRevertTimelocked(
            address(ipTokenStaking),
            abi.encodeWithSelector(IPTokenStaking.setFee.selector, 1),
            "IPTokenStaking: Invalid min fee"
        );
    }

    function testIPTokenStaking_fromSmartContract() public {
        // Network shall   allow anyone to create a new validator by staking validator’s own tokens (self-delegation)
        uint256 stakeAmount = ipTokenStaking.minStakeAmount();
        vm.deal(validatorAddr, stakeAmount);
        vm.prank(validatorAddr);
        ipTokenStaking.createValidator{ value: stakeAmount }({
            validatorCmpPubkey: validatorCmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100,
            supportsUnlocked: true,
            data: abi.encode("data")
        });
        uint256 expectedDelegationId = 0;

        // Deploy Staker contract
        Staker staker = new Staker(address(ipTokenStaking));

        // Test staking
        vm.deal(address(staker), 10_000_000_000 ether);
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Deposit(
            address(staker),
            validatorCmpPubkey,
            stakeAmount,
            uint32(uint8(IIPTokenStaking.StakingPeriod.FLEXIBLE)),
            expectedDelegationId,
            address(staker),
            ""
        );
        staker.stake{ value: stakeAmount }(validatorCmpPubkey);

        vm.deal(address(staker), 1 ether); // fee
        // Test redelegating
        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Redelegate(
            address(staker),
            validatorCmpPubkey,
            otherValidatorCmpPubkey,
            expectedDelegationId,
            address(staker),
            stakeAmount
        );
        staker.redelegate(validatorCmpPubkey, otherValidatorCmpPubkey, stakeAmount);

        // Test unstaking
        vm.deal(address(staker), 1 ether); // fee

        vm.expectEmit(address(ipTokenStaking));
        emit IIPTokenStaking.Withdraw(
            address(staker),
            validatorCmpPubkey,
            stakeAmount,
            expectedDelegationId,
            address(staker),
            ""
        );
        staker.unstake(validatorCmpPubkey, stakeAmount);
    }
}
