// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { CDR } from "../../src/protocol/CDR.sol";
import { ICDR } from "../../src/interfaces/ICDR.sol";
import { ICDRWriteCondition } from "../../src/interfaces/ICDRWriteCondition.sol";
import { ICDRReadCondition } from "../../src/interfaces/ICDRReadCondition.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Test } from "../utils/Test.sol";

/// @dev Mock write condition that allows only a specific writer
contract MockWriteCondition is ICDRWriteCondition {
    address public allowedWriter;

    constructor(address _allowedWriter) {
        allowedWriter = _allowedWriter;
    }

    function checkWriteCondition(uint32, bytes calldata, bytes calldata, address caller) external view returns (bool) {
        return caller == allowedWriter;
    }
}

/// @dev Mock read condition that allows only a specific reader
contract MockReadCondition is ICDRReadCondition {
    address public allowedReader;

    constructor(address _allowedReader) {
        allowedReader = _allowedReader;
    }

    function checkReadCondition(uint32, bytes calldata, bytes calldata, address caller) external view returns (bool) {
        return caller == allowedReader;
    }
}

/// @dev Mock condition that always rejects
contract RejectCondition is ICDRWriteCondition, ICDRReadCondition {
    function checkWriteCondition(uint32, bytes calldata, bytes calldata, address) external pure returns (bool) {
        return false;
    }

    function checkReadCondition(uint32, bytes calldata, bytes calldata, address) external pure returns (bool) {
        return false;
    }
}

contract CDRTest is Test {
    CDR internal cdr;

    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);
    address internal validator1 = address(0xBA1);

    MockWriteCondition internal writeCondition;
    MockReadCondition internal readCondition;
    RejectCondition internal rejectCondition;

    // PausableUpgradeable storage slot (OZ v5 ERC-7201)
    bytes32 internal constant PAUSABLE_STORAGE_LOCATION =
        0xcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300;

    function setUp() public virtual override {
        super.setUp();
        cdr = CDR(Predeploys.CDR);
        writeCondition = new MockWriteCondition(alice);
        readCondition = new MockReadCondition(alice);
        rejectCondition = new RejectCondition();
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Initialize Tests                             //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Initialize() public view {
        assertEq(cdr.owner(), address(timelock));
        assertEq(cdr.baseFee(), 0);
        assertEq(cdr.writeFee(), 0);
        assertEq(cdr.readFee(), 0);
        assertEq(cdr.allocateFee(), 0);
        assertEq(cdr.uuid(), 0);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Fee Setter Tests                             //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_SetBaseFee() public {
        uint256 newBaseFee = 0.1 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setBaseFee.selector, newBaseFee));
        assertEq(cdr.baseFee(), newBaseFee);
    }

    function testCDR_SetWriteFee() public {
        uint256 newWriteFee = 0.2 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setWriteFee.selector, newWriteFee));
        assertEq(cdr.writeFee(), newWriteFee);
    }

    function testCDR_SetReadFee() public {
        uint256 newReadFee = 0.3 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setReadFee.selector, newReadFee));
        assertEq(cdr.readFee(), newReadFee);
    }

    function testCDR_SetAllocateFee() public {
        uint256 newAllocateFee = 0.4 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setAllocateFee.selector, newAllocateFee));
        assertEq(cdr.allocateFee(), newAllocateFee);
    }

    function testCDR_SetFee_RevertIfNotOwner() public {
        vm.expectRevert();
        cdr.setBaseFee(1 ether);

        vm.expectRevert();
        cdr.setWriteFee(1 ether);

        vm.expectRevert();
        cdr.setReadFee(1 ether);

        vm.expectRevert();
        cdr.setAllocateFee(1 ether);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                            Allocate Tests                              //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Allocate() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        assertEq(vaultUuid, 0);

        ICDR.Vault memory vault = cdr.vaults(vaultUuid);
        assertEq(vault.updatable, true);
        assertEq(vault.writeConditionAddr, address(writeCondition));
        assertEq(vault.readConditionAddr, address(readCondition));
        assertEq(vault.encryptedData.length, 0);
    }

    function testCDR_Allocate_WithConditionData() public {
        bytes memory writeData = abi.encode("write-policy");
        bytes memory readData = abi.encode("read-policy");

        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), writeData, readData);

        ICDR.Vault memory vault = cdr.vaults(vaultUuid);
        assertEq(vault.writeConditionData, writeData);
        assertEq(vault.readConditionData, readData);
    }

    function testCDR_Allocate_RevertIfBothConditionsZero() public {
        vm.expectRevert("Invalid condition address");
        cdr.allocate(true, address(0), address(0), "", "");
    }

    function testCDR_Allocate_OneConditionZeroAllowed() public {
        // Only write condition
        uint32 uuid1 = cdr.allocate(true, address(writeCondition), address(0), "", "");
        assertEq(cdr.vaults(uuid1).writeConditionAddr, address(writeCondition));

        // Only read condition
        uint32 uuid2 = cdr.allocate(true, address(0), address(readCondition), "", "");
        assertEq(cdr.vaults(uuid2).readConditionAddr, address(readCondition));
    }

    function testCDR_Allocate_WithFee() public {
        uint256 fee = 0.5 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setAllocateFee.selector, fee));

        vm.expectRevert("CDR: Invalid fee amount");
        cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.deal(address(this), fee);
        cdr.allocate{ value: fee }(true, address(writeCondition), address(readCondition), "", "");
    }

    function testCDR_Allocate_FeeBurned() public {
        uint256 fee = 0.5 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setAllocateFee.selector, fee));

        vm.deal(alice, fee);
        uint256 burnAddrBefore = address(0x0).balance;

        vm.prank(alice);
        cdr.allocate{ value: fee }(true, address(writeCondition), address(readCondition), "", "");

        assertEq(address(0x0).balance, burnAddrBefore + fee);
    }

    function testCDR_Allocate_IncrementingUuid() public {
        uint32 uuid0 = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");
        uint32 uuid1 = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");
        uint32 uuid2 = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        assertEq(uuid0, 0);
        assertEq(uuid1, 1);
        assertEq(uuid2, 2);
        assertEq(cdr.uuid(), 3);
    }

    function testCDR_Allocate_EmitsVaultAllocated() public {
        vm.expectEmit(false, false, false, true);
        emit ICDR.VaultAllocated(0, true, address(writeCondition), address(readCondition), "", "");
        cdr.allocate(true, address(writeCondition), address(readCondition), "", "");
    }

    function testCDR_Allocate_EmitsFeeCollected() public {
        uint256 fee = 0.1 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setAllocateFee.selector, fee));

        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectEmit(true, false, false, true);
        emit ICDR.FeeCollected(alice, fee, ICDR.FeeType.Allocate);
        cdr.allocate{ value: fee }(true, address(writeCondition), address(readCondition), "", "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Write Tests                                //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Write() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");
        bytes memory encData = hex"deadbeef";

        vm.prank(alice);
        cdr.write(vaultUuid, "", encData);

        ICDR.Vault memory vault = cdr.vaults(vaultUuid);
        assertEq(vault.encryptedData, encData);
    }

    function testCDR_Write_EmitsVaultWritten() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");
        bytes memory encData = hex"deadbeef";

        vm.prank(alice);
        vm.expectEmit(false, false, false, true);
        emit ICDR.VaultWritten(vaultUuid, encData);
        cdr.write(vaultUuid, "", encData);
    }

    function testCDR_Write_RevertIfEmptyData() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.prank(alice);
        vm.expectRevert("CDR: Encrypted data cannot be empty");
        cdr.write(vaultUuid, "", "");
    }

    function testCDR_Write_RevertIfWriteConditionNotMet() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.prank(bob); // bob is not the allowed writer
        vm.expectRevert("CDR: Write condition not met");
        cdr.write(vaultUuid, "", hex"deadbeef");
    }

    function testCDR_Write_RevertIfVaultNotUpdatable() public {
        uint32 vaultUuid = cdr.allocate(false, address(writeCondition), address(readCondition), "", "");

        // First write succeeds
        vm.prank(alice);
        cdr.write(vaultUuid, "", hex"deadbeef");

        // Second write fails (not updatable)
        vm.prank(alice);
        vm.expectRevert("CDR: Vault is not updatable");
        cdr.write(vaultUuid, "", hex"cafebabe");
    }

    function testCDR_Write_UpdatableVaultOverwrite() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.prank(alice);
        cdr.write(vaultUuid, "", hex"deadbeef");

        vm.prank(alice);
        cdr.write(vaultUuid, "", hex"cafebabe");

        assertEq(cdr.vaults(vaultUuid).encryptedData, hex"cafebabe");
    }

    function testCDR_Write_WriteConditionAddrCanBypassCheck() public {
        // writeConditionAddr itself can write without condition check
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.prank(address(writeCondition));
        cdr.write(vaultUuid, "", hex"deadbeef");

        assertEq(cdr.vaults(vaultUuid).encryptedData, hex"deadbeef");
    }

    function testCDR_Write_RevertIfWriteConditionAddrNotSet() public {
        // Allocate with only read condition
        uint32 vaultUuid = cdr.allocate(true, address(0), address(readCondition), "", "");

        vm.prank(alice);
        vm.expectRevert("CDR: Write condition address not set");
        cdr.write(vaultUuid, "", hex"deadbeef");
    }

    function testCDR_Write_WithFee() public {
        uint256 fee = 0.2 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setWriteFee.selector, fee));

        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        // Revert without fee
        vm.prank(alice);
        vm.expectRevert("CDR: Invalid fee amount");
        cdr.write(vaultUuid, "", hex"deadbeef");

        // Success with fee
        vm.deal(alice, fee);
        vm.prank(alice);
        cdr.write{ value: fee }(vaultUuid, "", hex"deadbeef");
    }

    function testCDR_Write_RejectCondition() public {
        uint32 vaultUuid = cdr.allocate(true, address(rejectCondition), address(readCondition), "", "");

        vm.prank(alice);
        vm.expectRevert("CDR: Write condition not met");
        cdr.write(vaultUuid, "", hex"deadbeef");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Read Tests                                //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Read() public {
        uint32 vaultUuid = _allocateAndWrite();

        vm.prank(alice);
        cdr.read(vaultUuid, "", hex"04aabbccdd");
    }

    function testCDR_Read_EmitsVaultRead() public {
        uint32 vaultUuid = _allocateAndWrite();
        bytes memory requesterPubKey = hex"04aabbccdd";

        vm.prank(alice);
        vm.expectEmit(true, false, false, true);
        emit ICDR.VaultRead(vaultUuid, alice, hex"deadbeef", requesterPubKey);
        cdr.read(vaultUuid, "", requesterPubKey);
    }

    function testCDR_Read_RevertIfNoData() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.prank(alice);
        vm.expectRevert("CDR: Vault has no data to read");
        cdr.read(vaultUuid, "", hex"04aabbccdd");
    }

    function testCDR_Read_RevertIfReadConditionNotMet() public {
        uint32 vaultUuid = _allocateAndWrite();

        vm.prank(bob); // bob is not the allowed reader
        vm.expectRevert("CDR: Read condition not met");
        cdr.read(vaultUuid, "", hex"04aabbccdd");
    }

    function testCDR_Read_ReadConditionAddrCanBypassCheck() public {
        uint32 vaultUuid = _allocateAndWrite();

        vm.prank(address(readCondition));
        cdr.read(vaultUuid, "", hex"04aabbccdd");
    }

    function testCDR_Read_WithFee() public {
        uint256 fee = 0.3 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setReadFee.selector, fee));

        uint32 vaultUuid = _allocateAndWrite();

        // Revert without fee
        vm.prank(alice);
        vm.expectRevert("CDR: Invalid fee amount");
        cdr.read(vaultUuid, "", hex"04aabbccdd");

        // Success with fee
        vm.deal(alice, fee);
        vm.prank(alice);
        cdr.read{ value: fee }(vaultUuid, "", hex"04aabbccdd");
    }

    function testCDR_Read_RejectCondition() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(rejectCondition), "", "");

        vm.prank(alice);
        cdr.write(vaultUuid, "", hex"deadbeef");

        vm.prank(alice);
        vm.expectRevert("CDR: Read condition not met");
        cdr.read(vaultUuid, "", hex"04aabbccdd");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                    SubmitEncryptedPartialDecryption Tests               //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_SubmitPartial() public {
        vm.prank(validator1);
        cdr.submitEncryptedPartialDecryption(1, 0, hex"aa", hex"bb", hex"cc", hex"dd", hex"ff", 0, hex"ee");
    }

    function testCDR_SubmitPartial_EmitsEvent() public {
        uint256 fee = 0.05 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setBaseFee.selector, fee));

        vm.deal(validator1, fee);
        vm.prank(validator1);
        vm.expectEmit(true, false, false, true);
        emit ICDR.EncryptedPartialDecryptionSubmitted(
            validator1,
            1,
            0,
            hex"aa",
            hex"bb",
            hex"cc",
            hex"dd",
            hex"ff",
            0,
            hex"ee",
            fee
        );
        cdr.submitEncryptedPartialDecryption{ value: fee }(
            1,
            0,
            hex"aa",
            hex"bb",
            hex"cc",
            hex"dd",
            hex"ff",
            0,
            hex"ee"
        );
    }

    function testCDR_SubmitPartial_WithFee() public {
        uint256 fee = 0.05 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setBaseFee.selector, fee));

        // Revert without fee
        vm.prank(validator1);
        vm.expectRevert("CDR: Invalid fee amount");
        cdr.submitEncryptedPartialDecryption(1, 0, hex"aa", hex"bb", hex"cc", hex"dd", hex"ff", 0, hex"ee");

        // Success with fee
        vm.deal(validator1, fee);
        vm.prank(validator1);
        cdr.submitEncryptedPartialDecryption{ value: fee }(
            1,
            0,
            hex"aa",
            hex"bb",
            hex"cc",
            hex"dd",
            hex"ff",
            0,
            hex"ee"
        );
    }

    function testCDR_SubmitPartial_FeeBurned() public {
        uint256 fee = 0.05 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setBaseFee.selector, fee));

        vm.deal(validator1, fee);
        uint256 burnBefore = address(0x0).balance;

        vm.prank(validator1);
        cdr.submitEncryptedPartialDecryption{ value: fee }(
            1,
            0,
            hex"aa",
            hex"bb",
            hex"cc",
            hex"dd",
            hex"ff",
            0,
            hex"ee"
        );

        assertEq(address(0x0).balance, burnBefore + fee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                       Write Fee & Event Tests                          //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Write_EmitsFeeCollected() public {
        uint256 fee = 0.2 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setWriteFee.selector, fee));

        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectEmit(true, false, false, true);
        emit ICDR.FeeCollected(alice, fee, ICDR.FeeType.Write);
        cdr.write{ value: fee }(vaultUuid, "", hex"deadbeef");
    }

    function testCDR_Write_FeeBurned() public {
        uint256 fee = 0.2 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setWriteFee.selector, fee));

        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        vm.deal(alice, fee);
        uint256 burnBefore = address(0x0).balance;

        vm.prank(alice);
        cdr.write{ value: fee }(vaultUuid, "", hex"deadbeef");

        assertEq(address(0x0).balance, burnBefore + fee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                       Read Fee & Event Tests                           //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Read_EmitsFeeCollected() public {
        uint256 fee = 0.3 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setReadFee.selector, fee));

        uint32 vaultUuid = _allocateAndWrite();

        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectEmit(true, false, false, true);
        emit ICDR.FeeCollected(alice, fee, ICDR.FeeType.Read);
        cdr.read{ value: fee }(vaultUuid, "", hex"04aabbccdd");
    }

    function testCDR_Read_FeeBurned() public {
        uint256 fee = 0.3 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setReadFee.selector, fee));

        uint32 vaultUuid = _allocateAndWrite();

        vm.deal(alice, fee);
        uint256 burnBefore = address(0x0).balance;

        vm.prank(alice);
        cdr.read{ value: fee }(vaultUuid, "", hex"04aabbccdd");

        assertEq(address(0x0).balance, burnBefore + fee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                    Read Condition Boundary Tests                        //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Read_RevertIfReadConditionAddrNotSet() public {
        // Allocate with only write condition (readConditionAddr = address(0))
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(0), "", "");
        vm.prank(alice);
        cdr.write(vaultUuid, "", hex"deadbeef");

        // Read reverts because readConditionAddr is address(0) — unlike write(),
        // read() lacks an explicit zero-address check, causing a low-level revert
        vm.prank(alice);
        vm.expectRevert();
        cdr.read(vaultUuid, "", hex"04aabbccdd");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        WhenNotPaused Tests                             //
    //////////////////////////////////////////////////////////////////////////*/

    function testCDR_Allocate_RevertWhenPaused() public {
        _pause(address(cdr));

        vm.expectRevert(abi.encodeWithSignature("EnforcedPause()"));
        cdr.allocate(true, address(writeCondition), address(readCondition), "", "");
    }

    function testCDR_Write_RevertWhenPaused() public {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");

        _pause(address(cdr));

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSignature("EnforcedPause()"));
        cdr.write(vaultUuid, "", hex"deadbeef");
    }

    function testCDR_Read_RevertWhenPaused() public {
        uint32 vaultUuid = _allocateAndWrite();

        _pause(address(cdr));

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSignature("EnforcedPause()"));
        cdr.read(vaultUuid, "", hex"04aabbccdd");
    }

    function testCDR_SubmitPartial_RevertWhenPaused() public {
        _pause(address(cdr));

        vm.prank(validator1);
        vm.expectRevert(abi.encodeWithSignature("EnforcedPause()"));
        cdr.submitEncryptedPartialDecryption(1, 0, hex"aa", hex"bb", hex"cc", hex"dd", hex"ff", 0, hex"ee");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Helper Functions                             //
    //////////////////////////////////////////////////////////////////////////*/

    function _allocateAndWrite() internal returns (uint32) {
        uint32 vaultUuid = cdr.allocate(true, address(writeCondition), address(readCondition), "", "");
        vm.prank(alice);
        cdr.write(vaultUuid, "", hex"deadbeef");
        return vaultUuid;
    }

    /// @dev Directly sets the paused flag via storage manipulation (no public pause function exposed)
    function _pause(address target) internal {
        vm.store(target, PAUSABLE_STORAGE_LOCATION, bytes32(uint256(1)));
    }
}
