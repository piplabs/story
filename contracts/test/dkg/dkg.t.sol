// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { DKG } from "../../src/protocol/DKG.sol";
import { SGXValidationHook } from "../../src/protocol/SGXValidationHook.sol";
import { IDKG } from "../../src/interfaces/IDKG.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { IAttestationReportValidator } from "../../src/interfaces/IAttestationReportValidator.sol";
import { Test } from "../utils/Test.sol";

/// @dev Mock attestation report validator that always succeeds
contract MockAttestationReportValidator is IAttestationReportValidator {
    function validateReport(bytes32, bytes32, bytes calldata, bytes calldata) external pure returns (bool) {
        return true;
    }
}

/// @dev Mock attestation report validator that always fails
contract FailingAttestationReportValidator is IAttestationReportValidator {
    function validateReport(bytes32, bytes32, bytes calldata, bytes calldata) external pure returns (bool) {
        return false;
    }
}

contract DKGTest is Test {
    DKG internal dkg;

    // Genesis enclave type set by GenerateAlloc.setSGXValidationHook()
    bytes32 internal constant GENESIS_ENCLAVE_TYPE = bytes32(uint256(1));

    address internal alice = address(0xA11CE);

    bytes32 internal constant MOCK_ENCLAVE_TYPE = bytes32(uint256(99));
    bytes32 internal constant FAILING_ENCLAVE_TYPE = bytes32(uint256(98));
    // PausableUpgradeable storage slot (OZ v5 ERC-7201)
    bytes32 internal constant PAUSABLE_STORAGE_LOCATION =
        0xcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300;

    MockAttestationReportValidator internal mockValidator;
    FailingAttestationReportValidator internal failingValidator;

    function setUp() public virtual override {
        super.setUp();
        dkg = DKG(Predeploys.DKG);
        mockValidator = new MockAttestationReportValidator();
        failingValidator = new FailingAttestationReportValidator();
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Initialize Tests                             //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_Initialize() public view {
        assertEq(dkg.minReqRegisteredParticipants(), 3);
        assertEq(dkg.minReqFinalizedParticipants(), 3);
        assertEq(dkg.operationalThreshold(), 670);
        assertEq(dkg.fee(), 1 ether);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        Enclave Type Whitelist                          //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_GenesisEnclaveTypeWhitelisted() public view {
        assertTrue(dkg.isEnclaveTypeWhitelisted(GENESIS_ENCLAVE_TYPE));

        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(GENESIS_ENCLAVE_TYPE);
        assertEq(data.codeCommitment, bytes32(uint256(1)));
        assertTrue(data.validationHookAddr != address(0));
    }

    function testDKG_GenesisSGXValidationHook() public view {
        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(GENESIS_ENCLAVE_TYPE);
        SGXValidationHook sgxHook = SGXValidationHook(data.validationHookAddr);

        assertEq(sgxHook.owner(), address(timelock));
        assertEq(sgxHook.DKG(), Predeploys.DKG);
        assertEq(sgxHook.automataValidationAddr(), address(uint160(1000)));
        assertEq(sgxHook.tcbEvaluationDataNumber(), 0);
    }

    function testDKG_WhitelistNewEnclaveType() public {
        address sgxHookImpl = address(new SGXValidationHook(Predeploys.DKG));
        address sgxHookProxy = address(
            new ERC1967Proxy(
                sgxHookImpl,
                abi.encodeCall(SGXValidationHook.initialize, (address(timelock), address(2000), 1))
            )
        );

        bytes32 newEnclaveType = bytes32(uint256(2));
        IDKG.EnclaveTypeData memory enclaveTypeData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(42)),
            validationHookAddr: sgxHookProxy
        });
        performTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, newEnclaveType, enclaveTypeData, true)
        );

        assertTrue(dkg.isEnclaveTypeWhitelisted(newEnclaveType));
        assertEq(dkg.enclaveTypeData(newEnclaveType).validationHookAddr, sgxHookProxy);
        assertEq(dkg.enclaveTypeData(newEnclaveType).codeCommitment, bytes32(uint256(42)));
    }

    function testDKG_WhitelistEnclaveType_RevertIfNotOwner() public {
        IDKG.EnclaveTypeData memory enclaveTypeData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(1)),
            validationHookAddr: address(1)
        });

        vm.expectRevert();
        dkg.whitelistEnclaveType(bytes32(uint256(99)), enclaveTypeData, true);
    }

    function testDKG_WhitelistEnclaveType_RevertIfEmptyEnclaveType() public {
        IDKG.EnclaveTypeData memory data = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(1)),
            validationHookAddr: address(1)
        });

        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, bytes32(0), data, true),
            "DKG: Enclave type cannot be empty"
        );
    }

    function testDKG_WhitelistEnclaveType_RevertIfEmptyCodeCommitment() public {
        IDKG.EnclaveTypeData memory data = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(0),
            validationHookAddr: address(1)
        });

        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, bytes32(uint256(10)), data, true),
            "DKG: Code commitment cannot be empty"
        );
    }

    function testDKG_WhitelistEnclaveType_RevertIfEmptyValidationHook() public {
        IDKG.EnclaveTypeData memory data = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(1)),
            validationHookAddr: address(0)
        });

        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, bytes32(uint256(10)), data, true),
            "DKG: Validation hook cannot be empty"
        );
    }

    function testDKG_WhitelistEnclaveType_Delist() public {
        IDKG.EnclaveTypeData memory data = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(1)),
            validationHookAddr: address(1)
        });

        // Whitelist
        performTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, bytes32(uint256(10)), data, true)
        );
        assertTrue(dkg.isEnclaveTypeWhitelisted(bytes32(uint256(10))));

        // Delist
        performTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, bytes32(uint256(10)), data, false),
            bytes32(keccak256("DELIST_SALT"))
        );
        assertFalse(dkg.isEnclaveTypeWhitelisted(bytes32(uint256(10))));
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Admin Setter Tests                           //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_SetMinReqRegisteredParticipants() public {
        performTimelocked(address(dkg), abi.encodeWithSelector(DKG.setMinReqRegisteredParticipants.selector, 5));
        assertEq(dkg.minReqRegisteredParticipants(), 5);
    }

    function testDKG_SetMinReqRegisteredParticipants_RevertIfZero() public {
        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.setMinReqRegisteredParticipants.selector, 0),
            "DKG: MinReqRegisteredParticipants cannot be zero"
        );
    }

    function testDKG_SetMinReqRegisteredParticipants_RevertIfNotOwner() public {
        vm.expectRevert();
        dkg.setMinReqRegisteredParticipants(5);
    }

    function testDKG_SetMinReqFinalizedParticipants() public {
        performTimelocked(address(dkg), abi.encodeWithSelector(DKG.setMinReqFinalizedParticipants.selector, 5));
        assertEq(dkg.minReqFinalizedParticipants(), 5);
    }

    function testDKG_SetMinReqFinalizedParticipants_RevertIfZero() public {
        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.setMinReqFinalizedParticipants.selector, 0),
            "DKG: MinReqFinalizedParticipants cannot be zero"
        );
    }

    function testDKG_SetMinReqFinalizedParticipants_RevertIfNotOwner() public {
        vm.expectRevert();
        dkg.setMinReqFinalizedParticipants(5);
    }

    function testDKG_SetOperationalThreshold() public {
        performTimelocked(address(dkg), abi.encodeWithSelector(DKG.setOperationalThreshold.selector, 500));
        assertEq(dkg.operationalThreshold(), 500);
    }

    function testDKG_SetOperationalThreshold_RevertIfZero() public {
        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.setOperationalThreshold.selector, 0),
            "DKG: Operational threshold cannot be zero"
        );
    }

    function testDKG_SetOperationalThreshold_RevertIfAboveBasis() public {
        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.setOperationalThreshold.selector, 1001),
            "DKG: Operational threshold cannot be greater than 1000"
        );
    }

    function testDKG_SetOperationalThreshold_BoundaryMax() public {
        performTimelocked(address(dkg), abi.encodeWithSelector(DKG.setOperationalThreshold.selector, 1000));
        assertEq(dkg.operationalThreshold(), 1000);
    }

    function testDKG_SetOperationalThreshold_RevertIfNotOwner() public {
        vm.expectRevert();
        dkg.setOperationalThreshold(500);
    }

    function testDKG_SetFee() public {
        performTimelocked(address(dkg), abi.encodeWithSelector(DKG.setFee.selector, 2 ether));
        assertEq(dkg.fee(), 2 ether);
    }

    function testDKG_SetFee_RevertIfNotOwner() public {
        vm.expectRevert();
        dkg.setFee(2 ether);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        Upgrade Scheduling Tests                        //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_ScheduleUpgrade() public {
        uint256 futureHeight = block.number + 100;
        performTimelocked(address(dkg), abi.encodeWithSelector(DKG.scheduleUpgrade.selector, futureHeight, "v2.0.0"));
    }

    function testDKG_ScheduleUpgrade_RevertIfPastHeight() public {
        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.scheduleUpgrade.selector, 0, "v2.0.0"),
            "DKG: activation must be in future"
        );
    }

    function testDKG_ScheduleUpgrade_RevertIfEmptyVersion() public {
        uint256 futureHeight = block.number + 100;
        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.scheduleUpgrade.selector, futureHeight, ""),
            "DKG: upgrade version cannot be empty"
        );
    }

    function testDKG_ScheduleUpgrade_RevertIfNotOwner() public {
        vm.expectRevert();
        dkg.scheduleUpgrade(block.number + 100, "v2.0.0");
    }

    function testDKG_CancelUpgrade() public {
        performTimelocked(address(dkg), abi.encodeWithSelector(DKG.cancelUpgrade.selector, "v2.0.0"));
    }

    function testDKG_CancelUpgrade_RevertIfEmptyVersion() public {
        expectRevertTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.cancelUpgrade.selector, ""),
            "DKG: upgrade version cannot be empty"
        );
    }

    function testDKG_CancelUpgrade_RevertIfNotOwner() public {
        vm.expectRevert();
        dkg.cancelUpgrade("v2.0.0");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          Register Tests                                //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_Register_RevertIfEmptyReport() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Enclave report cannot be empty");
        dkg.register{ value: fee }("", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_RevertIfRoundZero() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.round = 0;
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Round cannot be zero");
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_RevertIfValidatorAddrZero() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.validatorAddr = address(0);
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Validator address cannot be empty");
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_RevertIfEmptyEnclaveType() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveType = bytes32(0);
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Enclave type cannot be empty");
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_RevertIfEmptyCommKey() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveCommKey = "";
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Enclave communication key cannot be empty");
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_RevertIfEmptyDkgPubKey() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.dkgPubKey = "";
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: DKG public key cannot be empty");
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_RevertIfWrongFee() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        vm.deal(alice, 0.5 ether);
        vm.prank(alice);
        vm.expectRevert("DKG: Invalid fee amount");
        dkg.register{ value: 0.5 ether }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          Finalize Tests                                //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_Finalize_RevertIfRoundZero() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Round cannot be zero");
        dkg.finalize{ value: fee }(
            0,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_RevertIfValidatorAddrZero() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Validator address cannot be empty");
        dkg.finalize{ value: fee }(
            1,
            address(0),
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_RevertIfEnclaveTypeNotWhitelisted() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Enclave type is not whitelisted");
        dkg.finalize{ value: fee }(
            1,
            alice,
            bytes32(uint256(999)),
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_RevertIfEmptyParticipantsRoot() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Participants root cannot be empty");
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(0),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_RevertIfEmptyGlobalPubKey() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Global public key cannot be empty");
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            "",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_RevertIfEmptyPublicCoeffs() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Public coefficients cannot be empty");
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            new bytes[](0),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_RevertIfEmptyPubKeyShare() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Public key share cannot be empty");
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            "",
            hex"cc"
        );
    }

    function testDKG_Finalize_RevertIfEmptySignature() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        vm.prank(alice);
        vm.expectRevert("DKG: Signature cannot be empty");
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            ""
        );
    }

    function testDKG_Finalize_RevertIfWrongFee() public {
        vm.deal(alice, 0.5 ether);
        vm.prank(alice);
        vm.expectRevert("DKG: Invalid fee amount");
        dkg.finalize{ value: 0.5 ether }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_Success() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_EmitsFinalized() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        vm.expectEmit(true, false, false, true);
        emit IDKG.Finalized(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)), // genesis codeCommitment
            bytes32(uint256(1)), // participantsRoot
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_Finalize_FeeBurned() public {
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        uint256 burnBefore = address(0x0).balance;

        vm.prank(alice);
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );

        assertEq(address(0x0).balance, burnBefore + fee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                    Register Happy Path Tests                           //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_Register_Success() public {
        _whitelistMockEnclaveType();

        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveType = MOCK_ENCLAVE_TYPE;
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_EmitsRegistered() public {
        _whitelistMockEnclaveType();

        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveType = MOCK_ENCLAVE_TYPE;
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        vm.expectEmit(true, false, false, true);
        emit IDKG.Registered(
            hex"aa",
            1,
            alice,
            MOCK_ENCLAVE_TYPE,
            hex"aabb",
            hex"ccdd",
            bytes32(uint256(42)), // mock codeCommitment
            100,
            bytes32(uint256(1)),
            ""
        );
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Register_FeeBurned() public {
        _whitelistMockEnclaveType();

        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveType = MOCK_ENCLAVE_TYPE;
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);
        uint256 burnBefore = address(0x0).balance;

        vm.prank(alice);
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");

        assertEq(address(0x0).balance, burnBefore + fee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                  AuthenticateEnclaveReport Tests                        //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_AuthenticateEnclaveReport_Success() public {
        _whitelistMockEnclaveType();

        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveType = MOCK_ENCLAVE_TYPE;
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        dkg.authenticateEnclaveReport{ value: fee }(hex"aa", instanceData, "");
    }

    function testDKG_AuthenticateEnclaveReport_RevertIfNotWhitelisted() public {
        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveType = bytes32(uint256(999));
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        vm.expectRevert("DKG: Enclave type is not whitelisted");
        dkg.authenticateEnclaveReport{ value: fee }(hex"aa", instanceData, "");
    }

    function testDKG_AuthenticateEnclaveReport_RevertIfValidationFails() public {
        _whitelistFailingEnclaveType();

        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        instanceData.enclaveType = FAILING_ENCLAVE_TYPE;
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        vm.expectRevert("DKG: Enclave authentication failed");
        dkg.authenticateEnclaveReport{ value: fee }(hex"aa", instanceData, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        Admin Event Emit Tests                          //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_SetMinReqRegisteredParticipants_EmitsEvent() public {
        bytes memory data = abi.encodeWithSelector(DKG.setMinReqRegisteredParticipants.selector, 5);
        schedule(address(dkg), data);
        waitForTimelock();
        vm.expectEmit(false, false, false, true);
        emit IDKG.MinReqRegisteredParticipantsSet(5);
        executeTimelocked(address(dkg), data);
    }

    function testDKG_SetMinReqFinalizedParticipants_EmitsEvent() public {
        bytes memory data = abi.encodeWithSelector(DKG.setMinReqFinalizedParticipants.selector, 5);
        schedule(address(dkg), data);
        waitForTimelock();
        vm.expectEmit(false, false, false, true);
        emit IDKG.MinReqFinalizedParticipantsSet(5);
        executeTimelocked(address(dkg), data);
    }

    function testDKG_SetOperationalThreshold_EmitsEvent() public {
        bytes memory data = abi.encodeWithSelector(DKG.setOperationalThreshold.selector, 500);
        schedule(address(dkg), data);
        waitForTimelock();
        vm.expectEmit(false, false, false, true);
        emit IDKG.OperationalThresholdSet(500);
        executeTimelocked(address(dkg), data);
    }

    function testDKG_SetFee_EmitsEvent() public {
        bytes memory data = abi.encodeWithSelector(DKG.setFee.selector, 2 ether);
        schedule(address(dkg), data);
        waitForTimelock();
        vm.expectEmit(false, false, false, true);
        emit IDKG.FeeSet(2 ether);
        executeTimelocked(address(dkg), data);
    }

    function testDKG_WhitelistEnclaveType_EmitsEvent() public {
        IDKG.EnclaveTypeData memory etData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(42)),
            validationHookAddr: address(mockValidator)
        });
        bytes memory data = abi.encodeWithSelector(
            DKG.whitelistEnclaveType.selector,
            bytes32(uint256(50)),
            etData,
            true
        );
        schedule(address(dkg), data);
        waitForTimelock();
        vm.expectEmit(false, false, false, true);
        emit IDKG.EnclaveTypeWhitelisted(bytes32(uint256(50)), bytes32(uint256(42)), address(mockValidator), true);
        executeTimelocked(address(dkg), data);
    }

    function testDKG_ScheduleUpgrade_EmitsEvent() public {
        uint256 futureHeight = block.number + 100;
        bytes memory data = abi.encodeWithSelector(DKG.scheduleUpgrade.selector, futureHeight, "v2.0.0");
        schedule(address(dkg), data);
        waitForTimelock();
        vm.expectEmit(false, false, false, true);
        emit IDKG.UpgradeScheduled(futureHeight, "v2.0.0");
        executeTimelocked(address(dkg), data);
    }

    function testDKG_CancelUpgrade_EmitsEvent() public {
        bytes memory data = abi.encodeWithSelector(DKG.cancelUpgrade.selector, "v2.0.0");
        schedule(address(dkg), data);
        waitForTimelock();
        vm.expectEmit(false, false, false, true);
        emit IDKG.UpgradeCancelled("v2.0.0");
        executeTimelocked(address(dkg), data);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        WhenNotPaused Tests                             //
    //////////////////////////////////////////////////////////////////////////*/

    function testDKG_Register_RevertWhenPaused() public {
        _pause(address(dkg));

        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSignature("EnforcedPause()"));
        dkg.register{ value: fee }(hex"aa", instanceData, 100, bytes32(uint256(1)), "");
    }

    function testDKG_Finalize_RevertWhenPaused() public {
        _pause(address(dkg));

        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSignature("EnforcedPause()"));
        dkg.finalize{ value: fee }(
            1,
            alice,
            GENESIS_ENCLAVE_TYPE,
            bytes32(uint256(1)),
            hex"aa",
            _singleBytesArray(),
            hex"bb",
            hex"cc"
        );
    }

    function testDKG_AuthenticateEnclaveReport_RevertWhenPaused() public {
        _pause(address(dkg));

        IDKG.EnclaveInstanceData memory instanceData = _defaultInstanceData();
        uint256 fee = dkg.fee();
        vm.deal(alice, fee);

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSignature("EnforcedPause()"));
        dkg.authenticateEnclaveReport{ value: fee }(hex"aa", instanceData, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Helper Functions                             //
    //////////////////////////////////////////////////////////////////////////*/

    function _defaultInstanceData() internal view returns (IDKG.EnclaveInstanceData memory) {
        return
            IDKG.EnclaveInstanceData({
                round: 1,
                validatorAddr: alice,
                enclaveType: GENESIS_ENCLAVE_TYPE,
                enclaveCommKey: hex"aabb",
                dkgPubKey: hex"ccdd"
            });
    }

    function _singleBytesArray() internal pure returns (bytes[] memory) {
        bytes[] memory arr = new bytes[](1);
        arr[0] = hex"aa";
        return arr;
    }

    function _whitelistMockEnclaveType() internal {
        IDKG.EnclaveTypeData memory data = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(42)),
            validationHookAddr: address(mockValidator)
        });
        performTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, MOCK_ENCLAVE_TYPE, data, true)
        );
    }

    function _whitelistFailingEnclaveType() internal {
        IDKG.EnclaveTypeData memory data = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(42)),
            validationHookAddr: address(failingValidator)
        });
        performTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, FAILING_ENCLAVE_TYPE, data, true)
        );
    }

    /// @dev Directly sets the paused flag via storage manipulation (no public pause function exposed)
    function _pause(address target) internal {
        vm.store(target, PAUSABLE_STORAGE_LOCATION, bytes32(uint256(1)));
    }
}
