// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { CDR } from "../../src/protocol/CDR.sol";
import { ICDR } from "../../src/interfaces/ICDR.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Test } from "../utils/Test.sol";

contract CDRTest is Test {
    CDR internal cdr;

    function setUp() public virtual override {
        super.setUp();
        cdr = CDR(Predeploys.CDR);
    }

    function testCDR_Initialize() public view {
        assertEq(cdr.owner(), address(timelock));
        assertEq(cdr.baseFee(), 0);
        assertEq(cdr.writeFee(), 0);
        assertEq(cdr.readFee(), 0);
        assertEq(cdr.allocateFee(), 0);
        assertEq(cdr.uuid(), 0);
    }

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

    function testCDR_Allocate() public {
        address writeCondition = address(0xAAA);
        address readCondition = address(0xBBB);
        bytes memory writeConditionData = abi.encode("write");
        bytes memory readConditionData = abi.encode("read");

        uint32 vaultUuid = cdr.allocate(true, writeCondition, readCondition, writeConditionData, readConditionData);

        assertEq(vaultUuid, 0);

        ICDR.Vault memory vault = cdr.vaults(vaultUuid);
        assertEq(vault.updatable, true);
        assertEq(vault.writeConditionAddr, writeCondition);
        assertEq(vault.readConditionAddr, readCondition);
        assertEq(vault.writeConditionData, writeConditionData);
        assertEq(vault.readConditionData, readConditionData);
        assertEq(vault.encryptedData.length, 0);
    }

    function testCDR_Allocate_RevertIfBothConditionsZero() public {
        vm.expectRevert("Invalid condition address");
        cdr.allocate(true, address(0), address(0), "", "");
    }

    function testCDR_Allocate_WithFee() public {
        uint256 fee = 0.5 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setAllocateFee.selector, fee));

        // Should revert without fee
        vm.expectRevert("CDR: Invalid fee amount");
        cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");

        // Should succeed with correct fee
        vm.deal(address(this), fee);
        cdr.allocate{ value: fee }(true, address(0xAAA), address(0xBBB), "", "");
    }

    function testCDR_Allocate_IncrementingUuid() public {
        uint32 uuid0 = cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");
        uint32 uuid1 = cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");
        uint32 uuid2 = cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");

        assertEq(uuid0, 0);
        assertEq(uuid1, 1);
        assertEq(uuid2, 2);
        assertEq(cdr.uuid(), 3);
    }

    // -------------------------------------------------------------------------
    // maxBatchSize
    // -------------------------------------------------------------------------

    function testCDR_Initialize_MaxBatchSize() public view {
        assertEq(cdr.maxBatchSize(), 20);
    }

    function testCDR_SetMaxBatchSize() public {
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setMaxBatchSize.selector, uint256(50)));
        assertEq(cdr.maxBatchSize(), 50);
    }

    function testCDR_SetMaxBatchSize_RevertIfNotOwner() public {
        vm.expectRevert();
        cdr.setMaxBatchSize(50);
    }

    function testCDR_SetMaxBatchSize_RevertIfZero() public {
        vm.prank(address(timelock));
        vm.expectRevert("CDR: Max batch size must be > 0");
        cdr.setMaxBatchSize(0);
    }

    // -------------------------------------------------------------------------
    // submitEncryptedPartialDecryptionBatch
    // -------------------------------------------------------------------------

    /// @dev Builds a valid PartialDecryptionRequest for use in batch tests.
    function _makeRequest(uint32 uuid) internal pure returns (ICDR.PartialDecryptionRequest memory) {
        return ICDR.PartialDecryptionRequest({
            round: 1,
            pid: 2,
            encryptedPartial: bytes("encrypted-partial"),
            ephemeralPubKey: bytes("eph-pub"),
            pubShare: bytes("pub-share"),
            requesterPubKey: bytes("requester-pub"),
            ciphertext: bytes("ciphertext"),
            uuid: uuid,
            signature: bytes("sig")
        });
    }

    function testCDR_SubmitBatch_EmptyBatch() public {
        ICDR.PartialDecryptionRequest[] memory requests = new ICDR.PartialDecryptionRequest[](0);
        vm.expectRevert("CDR: Empty batch");
        cdr.submitEncryptedPartialDecryptionBatch(requests);
    }

    function testCDR_SubmitBatch_ExceedsMax() public {
        uint256 max = cdr.maxBatchSize();
        ICDR.PartialDecryptionRequest[] memory requests = new ICDR.PartialDecryptionRequest[](max + 1);
        for (uint256 i = 0; i < requests.length; i++) {
            requests[i] = _makeRequest(uint32(i));
        }
        vm.expectRevert("CDR: Batch exceeds max size");
        cdr.submitEncryptedPartialDecryptionBatch(requests);
    }

    function testCDR_SubmitBatch_WrongFee() public {
        uint256 baseFee = 0.1 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setBaseFee.selector, baseFee));

        ICDR.PartialDecryptionRequest[] memory requests = new ICDR.PartialDecryptionRequest[](2);
        requests[0] = _makeRequest(0);
        requests[1] = _makeRequest(1);

        // Send fee for only 1 item instead of 2
        vm.deal(address(this), baseFee);
        vm.expectRevert("CDR: Invalid fee amount");
        cdr.submitEncryptedPartialDecryptionBatch{ value: baseFee }(requests);
    }

    function testCDR_SubmitBatch_Success() public {
        ICDR.PartialDecryptionRequest[] memory requests = new ICDR.PartialDecryptionRequest[](3);
        for (uint32 i = 0; i < 3; i++) {
            requests[i] = _makeRequest(i);
        }

        vm.expectEmit(true, false, false, false);
        emit ICDR.EncryptedPartialDecryptionSubmitted(
            address(this), 1, 2, bytes("encrypted-partial"), bytes("eph-pub"),
            bytes("pub-share"), bytes("requester-pub"), bytes("ciphertext"), 0, bytes("sig"), 0
        );

        cdr.submitEncryptedPartialDecryptionBatch(requests);
    }

    function testCDR_SubmitBatch_InvalidItem_Skipped() public {
        ICDR.PartialDecryptionRequest[] memory requests = new ICDR.PartialDecryptionRequest[](3);
        requests[0] = _makeRequest(0); // valid
        requests[1] = _makeRequest(1); // will be made invalid
        requests[1].encryptedPartial = bytes(""); // empty → invalid
        requests[2] = _makeRequest(2); // valid

        // Invalid item at index 1 must emit InvalidPartialDecryption and not revert.
        vm.expectEmit(true, false, false, true);
        emit ICDR.InvalidPartialDecryption(address(this), 1, 2, 1, 1);

        cdr.submitEncryptedPartialDecryptionBatch(requests);
    }

    function testCDR_SubmitBatch_InvalidItem_OversizedPartial() public {
        uint256 maxPartialSize = cdr.maxEncryptedPartialSize();

        ICDR.PartialDecryptionRequest[] memory requests = new ICDR.PartialDecryptionRequest[](2);
        requests[0] = _makeRequest(0); // valid
        requests[1] = _makeRequest(1);
        requests[1].encryptedPartial = new bytes(maxPartialSize + 1); // too large → invalid

        vm.expectEmit(true, false, false, true);
        emit ICDR.InvalidPartialDecryption(address(this), 1, 2, 1, 1);

        cdr.submitEncryptedPartialDecryptionBatch(requests);
    }

    function testCDR_SubmitBatch_WithFee() public {
        uint256 baseFee = 0.01 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setBaseFee.selector, baseFee));

        uint256 n = 3;
        ICDR.PartialDecryptionRequest[] memory requests = new ICDR.PartialDecryptionRequest[](n);
        for (uint32 i = 0; i < n; i++) {
            requests[i] = _makeRequest(i);
        }

        uint256 totalFee = baseFee * n;
        vm.deal(address(this), totalFee);

        vm.expectEmit(true, false, false, true);
        emit ICDR.FeeCollected(address(this), totalFee, ICDR.FeeType.SubmitPartial);

        cdr.submitEncryptedPartialDecryptionBatch{ value: totalFee }(requests);
    }
}
