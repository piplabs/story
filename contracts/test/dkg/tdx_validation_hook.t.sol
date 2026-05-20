// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable max-line-length */

import { Test as ForgeTest } from "forge-std/Test.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { TDXValidationHook } from "../../src/protocol/TDXValidationHook.sol";
import { IAutomataDcapAttestationFee } from "../../src/interfaces/external/IAutomataDcapAttestationFee.sol";

/// @dev Mock of the Automata DCAP attestation contract. Records the last quote and call count,
///      and returns a tester-controlled success flag; no real signature/TCB logic.
contract MockAutomataDcap is IAutomataDcapAttestationFee {
    bool public shouldSucceed;
    bytes public lastQuote;
    uint256 public callCount;
    bool public noArgCalled;

    constructor() {
        shouldSucceed = true;
    }

    function setShouldSucceed(bool s) external {
        shouldSucceed = s;
    }

    function verifyAndAttestOnChain(
        bytes calldata rawQuote
    ) external payable override returns (bool success, bytes memory output) {
        lastQuote = rawQuote;
        callCount++;
        noArgCalled = true;
        return (shouldSucceed, bytes(""));
    }

    /// @dev The hook uses the no-arg overload; this branch only exists to satisfy the interface.
    function verifyAndAttestOnChain(
        bytes calldata rawQuote,
        uint32 /* tcbEvaluationDataNumber */
    ) external payable override returns (bool success, bytes memory output) {
        lastQuote = rawQuote;
        callCount++;
        return (shouldSucceed, bytes(""));
    }
}

contract TDXValidationHookTest is ForgeTest {
    TDXValidationHook internal hook;
    MockAutomataDcap internal automata;

    address internal owner = address(0xA11CE);
    address internal dkg = address(0xDDDDDD);
    address internal stranger = address(0xBEEF);

    /*//////////////////////////////////////////////////////////////////////////
    //                              TDX quote layout                          //
    //                  (mirrors constants in TDXValidationHook)              //
    //////////////////////////////////////////////////////////////////////////*/
    uint256 internal constant QUOTE_HEADER_SIZE = 48;
    uint256 internal constant V4_BODY_SIZE = 584;
    uint256 internal constant V5_BODY_SIZE = 648;
    uint256 internal constant MIN_V4_QUOTE_SIZE = QUOTE_HEADER_SIZE + V4_BODY_SIZE; // 632
    uint256 internal constant MIN_V5_QUOTE_SIZE = QUOTE_HEADER_SIZE + V5_BODY_SIZE; // 696
    uint256 internal constant OFFSET_VERSION = 0;
    uint256 internal constant OFFSET_TEE_TYPE = 4;
    uint256 internal constant OFFSET_MRTD = 184;
    uint256 internal constant OFFSET_RTMR0 = 376;
    uint256 internal constant OFFSET_RTMR1 = 424;
    uint256 internal constant OFFSET_RTMR2 = 472;
    uint256 internal constant OFFSET_REPORT_DATA = 568;
    uint256 internal constant MEASUREMENT_SIZE = 48;
    uint8 internal constant TEE_TYPE_TDX_BYTE0 = 0x81;

    /*//////////////////////////////////////////////////////////////////////////
    //                  Event signatures (mirrored for vm.expectEmit)         //
    //////////////////////////////////////////////////////////////////////////*/
    event AutomataValidationAddrSet(address indexed newAutomataValidationAddr);
    event PlatformApproved(bytes32 indexed platformCommitment, string label);
    event PlatformRevoked(bytes32 indexed platformCommitment);

    function setUp() public {
        automata = new MockAutomataDcap();
        address impl = address(new TDXValidationHook(dkg));
        bytes memory initData = abi.encodeCall(TDXValidationHook.initialize, (owner, address(automata)));
        address proxy = address(new ERC1967Proxy(impl, initData));
        hook = TDXValidationHook(proxy);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                                Helpers                                 //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Returns (keccak256(RTMR2), keccak256(MRTD || RTMR0 || RTMR1)) for a quote.
    function _commitments(bytes memory quote) internal pure returns (bytes32 binary, bytes32 platform) {
        bytes memory rtmr2 = new bytes(MEASUREMENT_SIZE);
        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            rtmr2[i] = quote[OFFSET_RTMR2 + i];
        }
        binary = keccak256(rtmr2);

        bytes memory plat = new bytes(MEASUREMENT_SIZE * 3);
        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            plat[i] = quote[OFFSET_MRTD + i];
            plat[MEASUREMENT_SIZE + i] = quote[OFFSET_RTMR0 + i];
            plat[MEASUREMENT_SIZE * 2 + i] = quote[OFFSET_RTMR1 + i];
        }
        platform = keccak256(plat);
    }

    /// @dev Returns the first 32 bytes of REPORT_DATA from a quote.
    function _dataCommitment(bytes memory quote) internal pure returns (bytes32 rd) {
        assembly {
            rd := mload(add(add(quote, 32), OFFSET_REPORT_DATA))
        }
    }

    function _approvePlatformForQuote(bytes memory quote) internal {
        (, bytes32 platform) = _commitments(quote);
        vm.prank(owner);
        hook.approvePlatform(platform, "test-platform");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Initialization                            //
    //////////////////////////////////////////////////////////////////////////*/

    function test_Initialize() public view {
        assertEq(hook.DKG(), dkg);
        assertEq(hook.owner(), owner);
        assertEq(hook.automataValidationAddr(), address(automata));
    }

    function test_Constructor_RevertOnZeroDKG() public {
        vm.expectRevert(bytes("TDXValidationHook: DKG cannot be empty"));
        new TDXValidationHook(address(0));
    }

    function test_Initialize_RevertOnZeroOwner() public {
        address impl = address(new TDXValidationHook(dkg));
        bytes memory initData = abi.encodeCall(TDXValidationHook.initialize, (address(0), address(automata)));
        vm.expectRevert(bytes("TDXValidationHook: owner cannot be empty"));
        new ERC1967Proxy(impl, initData);
    }

    function test_Initialize_RevertOnZeroAutomata() public {
        address impl = address(new TDXValidationHook(dkg));
        bytes memory initData = abi.encodeCall(TDXValidationHook.initialize, (owner, address(0)));
        vm.expectRevert(bytes("TDXValidationHook: Automata Validation cannot be empty"));
        new ERC1967Proxy(impl, initData);
    }

    function test_Initialize_RevertOnReinit() public {
        vm.expectRevert();
        hook.initialize(owner, address(automata));
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Admin Setters                             //
    //////////////////////////////////////////////////////////////////////////*/

    function test_SetAutomataValidationAddr_OnlyOwner() public {
        vm.prank(stranger);
        vm.expectRevert();
        hook.setAutomataValidationAddr(address(0xCAFE));
    }

    function test_SetAutomataValidationAddr_ZeroReverts() public {
        vm.prank(owner);
        vm.expectRevert(bytes("TDXValidationHook: Automata Validation cannot be empty"));
        hook.setAutomataValidationAddr(address(0));
    }

    function test_SetAutomataValidationAddr_Updates() public {
        address next = address(0xCAFE);
        vm.expectEmit(true, false, false, true, address(hook));
        emit AutomataValidationAddrSet(next);
        vm.prank(owner);
        hook.setAutomataValidationAddr(next);
        assertEq(hook.automataValidationAddr(), next);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                  Hybrid — platform whitelist setters                   //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ApprovePlatform_OnlyOwner() public {
        vm.prank(stranger);
        vm.expectRevert();
        hook.approvePlatform(bytes32(uint256(1)), "x");
    }

    function test_RevokePlatform_OnlyOwner() public {
        vm.prank(stranger);
        vm.expectRevert();
        hook.revokePlatform(bytes32(uint256(1)));
    }

    function test_ApprovePlatform_RejectsZero() public {
        vm.prank(owner);
        vm.expectRevert(bytes("TDXValidationHook: platform commitment cannot be empty"));
        hook.approvePlatform(bytes32(0), "x");
    }

    function test_ApprovePlatform_EmitsEvent() public {
        bytes32 key = keccak256("platform-1");
        vm.expectEmit(true, false, false, true, address(hook));
        emit PlatformApproved(key, "GCP c3-standard-4 / TDVF v1.5 / kernel-1.7.0");
        vm.prank(owner);
        hook.approvePlatform(key, "GCP c3-standard-4 / TDVF v1.5 / kernel-1.7.0");
        assertTrue(hook.isPlatformApproved(key));
    }

    function test_RevokePlatform_EmitsEvent() public {
        bytes32 key = keccak256("platform-1");
        vm.prank(owner);
        hook.approvePlatform(key, "x");
        vm.expectEmit(true, false, false, true, address(hook));
        emit PlatformRevoked(key);
        vm.prank(owner);
        hook.revokePlatform(key);
        assertFalse(hook.isPlatformApproved(key));
    }

    function test_RevokePlatform_RejectsZero() public {
        vm.prank(owner);
        vm.expectRevert(bytes("TDXValidationHook: platform commitment cannot be empty"));
        hook.revokePlatform(bytes32(0));
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          validateReport guards                         //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_OnlyDKG() public {
        bytes memory quote = _buildV4QuoteDefault();
        vm.prank(stranger);
        vm.expectRevert(bytes("TDXValidationHook: Only DKG can call this function"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_EmptyReport() public {
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Empty enclave report"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bytes(""), "");
    }

    function test_ValidateReport_ZeroCodeCommitment() public {
        bytes memory quote = _buildV4QuoteDefault();
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Zero code commitment"));
        hook.validateReport(bytes32(0), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_ZeroDataCommitment() public {
        bytes memory quote = _buildV4QuoteDefault();
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Zero data commitment"));
        hook.validateReport(bytes32(uint256(1)), bytes32(0), quote, "");
    }

    function test_ValidateReport_QuoteTooShortForHeader() public {
        bytes memory quote = new bytes(QUOTE_HEADER_SIZE - 1);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Quote too short for header"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    /// @dev Defense-in-depth length floor: a quote that passes the header check but is shorter
    ///      than the highest field we read (REPORT_DATA[0:32], at absolute offset 600) must
    ///      revert before any assembly extractor runs.
    function test_ValidateReport_QuoteTooShortForBody() public {
        // A quote that satisfies the header floor but is shorter than MIN_QUOTE_SIZE (600).
        // We can't use _buildV4QuoteDefault here because it builds a full V4 minimum quote.
        bytes memory quote = new bytes(QUOTE_HEADER_SIZE + 10);
        // Make the tee_type bytes match TDX so the tee_type guard does not short-circuit us
        // out before the length floor is exercised.
        quote[OFFSET_TEE_TYPE] = bytes1(TEE_TYPE_TDX_BYTE0);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Quote too short for body"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_RejectsSGXTeeType() public {
        // tee_type left as 0 (SGX) — must be rejected before offset extraction.
        bytes memory quote = new bytes(MIN_V4_QUOTE_SIZE);
        quote[0] = bytes1(uint8(4));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Not a TDX quote"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_RejectsBadTeeTypeTrailingByte() public {
        // tee_type[0] is 0x81 but tee_type[5] is non-zero — full uint32 must equal 0x00000081.
        bytes memory quote = _buildV4QuoteDefault();
        quote[OFFSET_TEE_TYPE + 1] = bytes1(uint8(0x01));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Not a TDX quote"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Automata + matching                          //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_AutomataFailureReverts() public {
        automata.setShouldSucceed(false);
        bytes memory quote = _buildV4QuoteDefault();
        _approvePlatformForQuote(quote);
        (bytes32 binary, ) = _commitments(quote);
        bytes32 data = _dataCommitment(quote);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Attestation failed"));
        hook.validateReport(binary, data, quote, "");
    }

    function test_ValidateReport_UnapprovedBinary() public {
        bytes memory quote = _buildV4QuoteDefault();
        _approvePlatformForQuote(quote);
        bytes32 wrongBinary = keccak256("wrong-binary");
        bytes32 data = _dataCommitment(quote);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: unapproved binary"));
        hook.validateReport(wrongBinary, data, quote, "");
    }

    function test_ValidateReport_UnapprovedPlatform() public {
        // No approvePlatform call — platform check must fail after binary check passes.
        bytes memory quote = _buildV4QuoteDefault();
        (bytes32 binary, ) = _commitments(quote);
        bytes32 data = _dataCommitment(quote);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: unapproved platform"));
        hook.validateReport(binary, data, quote, "");
    }

    function test_ValidateReport_DataCommitmentMismatch() public {
        bytes memory quote = _buildV4QuoteDefault();
        _approvePlatformForQuote(quote);
        (bytes32 binary, ) = _commitments(quote);
        bytes32 wrongData = keccak256("wrong-data");
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Data commitment does not match"));
        hook.validateReport(binary, wrongData, quote, "");
    }

    function test_ValidateReport_HappyPathV4() public {
        bytes memory quote = _buildV4QuoteDefault();
        _approvePlatformForQuote(quote);
        (bytes32 binary, ) = _commitments(quote);
        bytes32 data = _dataCommitment(quote);
        vm.prank(dkg);
        bool ok = hook.validateReport(binary, data, quote, "");
        assertTrue(ok);
        // Hook must use the no-arg Automata overload (mirrors PR #816).
        assertEq(automata.callCount(), 1);
        assertTrue(automata.noArgCalled());
    }

    function test_ValidateReport_HappyPathV5() public {
        bytes memory quote = _buildV5QuoteDefault();
        _approvePlatformForQuote(quote);
        (bytes32 binary, ) = _commitments(quote);
        bytes32 data = _dataCommitment(quote);
        vm.prank(dkg);
        bool ok = hook.validateReport(binary, data, quote, "");
        assertTrue(ok);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                       Decomposition invariants                         //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Same RTMR2, different MRTD/RTMR0/RTMR1 → binary matches, platform differs.
    function test_Decomposition_SameBinaryDifferentPlatform() public pure {
        bytes memory q1 = _buildQuoteStatic(4, MIN_V4_QUOTE_SIZE, 0xA0, 0x10, 0xC0);
        bytes memory q2 = _buildQuoteStatic(4, MIN_V4_QUOTE_SIZE, 0xB0, 0x50, 0xC0);

        (bytes32 binary1, bytes32 platform1) = _commitmentsStatic(q1);
        (bytes32 binary2, bytes32 platform2) = _commitmentsStatic(q2);
        require(binary1 == binary2, "binary mismatch");
        require(platform1 != platform2, "platform should differ");
    }

    /// @dev Same MRTD/RTMR0/RTMR1, different RTMR2 → platform matches, binary differs.
    function test_Decomposition_DifferentBinarySamePlatform() public pure {
        bytes memory q1 = _buildQuoteStatic(4, MIN_V4_QUOTE_SIZE, 0xA0, 0x10, 0xC0);
        bytes memory q2 = _buildQuoteStatic(4, MIN_V4_QUOTE_SIZE, 0xA0, 0x10, 0xD0);

        (bytes32 binary1, bytes32 platform1) = _commitmentsStatic(q1);
        (bytes32 binary2, bytes32 platform2) = _commitmentsStatic(q2);
        require(platform1 == platform2, "platform mismatch");
        require(binary1 != binary2, "binary should differ");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        Synthetic quote builder                         //
    //////////////////////////////////////////////////////////////////////////*/

    function _buildV4QuoteDefault() internal pure returns (bytes memory) {
        return _buildQuoteStatic(4, MIN_V4_QUOTE_SIZE, 0xA0, 0x10, 0xC0);
    }

    function _buildV5QuoteDefault() internal pure returns (bytes memory) {
        return _buildQuoteStatic(5, MIN_V5_QUOTE_SIZE, 0xA0, 0x10, 0xC0);
    }

    /// @dev Synthesizes a TDX quote with deterministic field bytes. Seeds let decomposition
    ///      tests vary individual measurement fields while holding others fixed. All non-seeded
    ///      bytes remain zero — the hook ignores them once Automata returns success.
    function _buildQuoteStatic(
        uint16 version,
        uint256 totalLen,
        uint8 mrtdSeed,
        uint8 rtmrPrefix,
        uint8 rtmr2Seed
    ) internal pure returns (bytes memory) {
        bytes memory quote = new bytes(totalLen);

        // version: little-endian uint16 at offset 0..1.
        quote[OFFSET_VERSION] = bytes1(uint8(version & 0xff));
        quote[OFFSET_VERSION + 1] = bytes1(uint8((version >> 8) & 0xff));

        // tee_type: little-endian uint32 = 0x00000081.
        quote[OFFSET_TEE_TYPE] = bytes1(TEE_TYPE_TDX_BYTE0);

        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            quote[OFFSET_MRTD + i] = bytes1(uint8(mrtdSeed + (i & 0x0F)));
        }
        // RTMR0 and RTMR1 share `rtmrPrefix` but are offset by row so they still differ.
        for (uint256 r = 0; r < 2; r++) {
            for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
                quote[OFFSET_RTMR0 + r * MEASUREMENT_SIZE + i] = bytes1(uint8(rtmrPrefix + r * 0x10 + (i & 0x0F)));
            }
        }
        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            quote[OFFSET_RTMR2 + i] = bytes1(uint8(rtmr2Seed + (i & 0x0F)));
        }

        bytes32 rd = keccak256("tdx-validation-hook-test-reportdata");
        for (uint256 i = 0; i < 32; i++) {
            quote[OFFSET_REPORT_DATA + i] = rd[i];
        }

        return quote;
    }

    /// @dev `pure` variant of `_commitments` for use inside `pure` tests.
    function _commitmentsStatic(bytes memory quote) internal pure returns (bytes32 binary, bytes32 platform) {
        bytes memory rtmr2 = new bytes(MEASUREMENT_SIZE);
        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            rtmr2[i] = quote[OFFSET_RTMR2 + i];
        }
        binary = keccak256(rtmr2);

        bytes memory plat = new bytes(MEASUREMENT_SIZE * 3);
        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            plat[i] = quote[OFFSET_MRTD + i];
            plat[MEASUREMENT_SIZE + i] = quote[OFFSET_RTMR0 + i];
            plat[MEASUREMENT_SIZE * 2 + i] = quote[OFFSET_RTMR1 + i];
        }
        platform = keccak256(plat);
    }
}
