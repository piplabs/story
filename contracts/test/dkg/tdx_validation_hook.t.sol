// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable max-line-length */

import { Test as ForgeTest } from "forge-std/Test.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { TDXValidationHook } from "../../src/protocol/TDXValidationHook.sol";
import { IAutomataDcapAttestationFee } from "../../src/interfaces/external/IAutomataDcapAttestationFee.sol";
import { TDXBundle } from "../../src/protocol/lib/TDXBundle.sol";
import { TPM2Attest } from "../../src/protocol/lib/TPM2Attest.sol";
import { RSASSAVerify } from "../../src/protocol/lib/RSASSAVerify.sol";

/// @dev Minimal mock of the Automata DCAP attestation contract. Captures the last
///      quote/tcb-num it was called with and returns a tester-controlled success flag.
///      Production behavior (signature verification, TCB lookup) is intentionally not
///      modeled; the hook treats Automata as a black-box trust anchor for the
///      cryptographic chain of trust.
contract MockAutomataDcap is IAutomataDcapAttestationFee {
    bool public shouldSucceed;
    bytes public lastQuote;
    uint32 public lastTcbNum;
    uint256 public callCount;

    constructor() {
        shouldSucceed = true;
    }

    function setShouldSucceed(bool s) external {
        shouldSucceed = s;
    }

    function verifyAndAttestOnChain(
        bytes calldata rawQuote,
        uint32 tcbEvaluationDataNumber
    ) external payable override returns (bool success, bytes memory output) {
        lastQuote = rawQuote;
        lastTcbNum = tcbEvaluationDataNumber;
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

    uint32 internal constant TCB_NUM = 7;

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
    uint256 internal constant OFFSET_REPORT_DATA = 568;
    uint8 internal constant TEE_TYPE_TDX_BYTE0 = 0x81;

    function setUp() public {
        automata = new MockAutomataDcap();
        address impl = address(new TDXValidationHook(dkg));
        bytes memory initData = abi.encodeCall(TDXValidationHook.initialize, (owner, address(automata), TCB_NUM));
        address proxy = address(new ERC1967Proxy(impl, initData));
        hook = TDXValidationHook(proxy);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Initialization                            //
    //////////////////////////////////////////////////////////////////////////*/

    function test_Initialize() public view {
        assertEq(hook.DKG(), dkg);
        assertEq(hook.owner(), owner);
        assertEq(hook.automataValidationAddr(), address(automata));
        assertEq(hook.tcbEvaluationDataNumber(), TCB_NUM);
    }

    function test_Constructor_RevertOnZeroDKG() public {
        vm.expectRevert(bytes("TDXValidationHook: DKG cannot be empty"));
        new TDXValidationHook(address(0));
    }

    function test_Initialize_RevertOnZeroAutomata() public {
        address impl = address(new TDXValidationHook(dkg));
        bytes memory initData = abi.encodeCall(TDXValidationHook.initialize, (owner, address(0), TCB_NUM));
        vm.expectRevert(bytes("TDXValidationHook: Automata Validation cannot be empty"));
        new ERC1967Proxy(impl, initData);
    }

    function test_Initialize_RevertOnReinit() public {
        vm.expectRevert();
        hook.initialize(owner, address(automata), TCB_NUM);
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
        vm.prank(owner);
        hook.setAutomataValidationAddr(next);
        assertEq(hook.automataValidationAddr(), next);
    }

    function test_SetTcbEvaluationDataNumber_OnlyOwner() public {
        vm.prank(stranger);
        vm.expectRevert();
        hook.setTcbEvaluationDataNumber(99);
    }

    function test_SetTcbEvaluationDataNumber_Updates() public {
        vm.prank(owner);
        hook.setTcbEvaluationDataNumber(99);
        assertEq(hook.tcbEvaluationDataNumber(), 99);
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

    function test_ValidateReport_RejectsSGXTeeType() public {
        // Header big enough, but tee_type bytes are 0 (SGX), not 0x81 (TDX).
        bytes memory quote = new bytes(MIN_V4_QUOTE_SIZE);
        // version=4 LE
        quote[0] = bytes1(uint8(4));
        // tee_type left as 0 (SGX)
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Not a TDX quote"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_RejectsBadTeeTypeTrailingByte() public {
        // tee_type[0] is correct (0x81) but tee_type[5] is non-zero — should still fail.
        bytes memory quote = _buildV4QuoteDefault();
        quote[OFFSET_TEE_TYPE + 1] = bytes1(uint8(0x01));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Not a TDX quote"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_RejectsUnsupportedVersion() public {
        bytes memory quote = _buildV4QuoteDefault();
        // overwrite version bytes to 3 (not in {4, 5})
        quote[OFFSET_VERSION] = bytes1(uint8(3));
        quote[OFFSET_VERSION + 1] = 0;
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Unsupported quote version"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_V4TooShortForBody() public {
        // V4 quote: pass header + tee_type checks, but length one byte short.
        bytes memory quote = new bytes(MIN_V4_QUOTE_SIZE - 1);
        quote[0] = bytes1(uint8(4));
        quote[OFFSET_TEE_TYPE] = bytes1(TEE_TYPE_TDX_BYTE0);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Quote too short for V4 body"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    function test_ValidateReport_V5TooShortForBody() public {
        bytes memory quote = new bytes(MIN_V5_QUOTE_SIZE - 1);
        quote[0] = bytes1(uint8(5));
        quote[OFFSET_TEE_TYPE] = bytes1(TEE_TYPE_TDX_BYTE0);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Quote too short for V5 body"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), quote, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Automata + matching                          //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_AutomataFailureReverts() public {
        automata.setShouldSucceed(false);
        bytes memory quote = _buildV4QuoteDefault();
        (bytes32 codeCommit, bytes32 dataCommit) = _expectedDigests(quote);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Attestation failed"));
        hook.validateReport(codeCommit, dataCommit, quote, "");
    }

    function test_ValidateReport_CodeCommitmentMismatch() public {
        bytes memory quote = _buildV4QuoteDefault();
        (, bytes32 dataCommit) = _expectedDigests(quote);
        bytes32 wrongCode = keccak256("wrong-code");
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Code commitment does not match"));
        hook.validateReport(wrongCode, dataCommit, quote, "");
    }

    function test_ValidateReport_DataCommitmentMismatch() public {
        bytes memory quote = _buildV4QuoteDefault();
        (bytes32 codeCommit, ) = _expectedDigests(quote);
        bytes32 wrongData = keccak256("wrong-data");
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Data commitment does not match"));
        hook.validateReport(codeCommit, wrongData, quote, "");
    }

    function test_ValidateReport_HappyPathV4() public {
        bytes memory quote = _buildV4QuoteDefault();
        (bytes32 codeCommit, bytes32 dataCommit) = _expectedDigests(quote);
        vm.prank(dkg);
        bool ok = hook.validateReport(codeCommit, dataCommit, quote, "");
        assertTrue(ok);
        // Automata should have been invoked exactly once with the same raw quote.
        assertEq(automata.callCount(), 1);
        assertEq(automata.lastTcbNum(), TCB_NUM);
    }

    function test_ValidateReport_HappyPathV5() public {
        bytes memory quote = _buildV5QuoteDefault();
        (bytes32 codeCommit, bytes32 dataCommit) = _expectedDigests(quote);
        vm.prank(dkg);
        bool ok = hook.validateReport(codeCommit, dataCommit, quote, "");
        assertTrue(ok);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        Synthetic quote builder                         //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Builds a V4 quote with deterministic, non-zero MRTD/RTMR/reportData fields.
    function _buildV4QuoteDefault() internal pure returns (bytes memory) {
        return _buildQuote(4, MIN_V4_QUOTE_SIZE);
    }

    /// @dev Builds a V5 quote with deterministic, non-zero MRTD/RTMR/reportData fields.
    function _buildV5QuoteDefault() internal pure returns (bytes memory) {
        return _buildQuote(5, MIN_V5_QUOTE_SIZE);
    }

    /// @dev Internal quote builder. Plants header (version + tee_type), MRTD, RTMR0..3,
    ///      and the first 32 bytes of REPORT_DATA. All other body bytes are zero, which
    ///      is sufficient for our extraction tests because the hook does not look at
    ///      any other field once Automata returns success.
    function _buildQuote(uint16 version, uint256 totalLen) internal pure returns (bytes memory) {
        bytes memory quote = new bytes(totalLen);

        // version: little-endian uint16 at offset 0..1
        quote[OFFSET_VERSION] = bytes1(uint8(version & 0xff));
        quote[OFFSET_VERSION + 1] = bytes1(uint8((version >> 8) & 0xff));

        // tee_type: little-endian uint32 = 0x00000081 at offset 4..7
        quote[OFFSET_TEE_TYPE] = bytes1(TEE_TYPE_TDX_BYTE0);
        // bytes 5..7 already zero.

        // Plant 48 deterministic bytes for MRTD. Pattern 0xA0+i avoids both all-zero
        // (which a real TD might produce only on a broken vTPM) and any accidental
        // collision with reportData below.
        for (uint256 i = 0; i < 48; i++) {
            quote[OFFSET_MRTD + i] = bytes1(uint8(0xA0 + (i & 0x0F)));
        }
        // RTMR0..3 each 48 bytes, contiguous from OFFSET_RTMR0.
        for (uint256 r = 0; r < 4; r++) {
            for (uint256 i = 0; i < 48; i++) {
                quote[OFFSET_RTMR0 + r * 48 + i] = bytes1(uint8(0x10 * (r + 1) + (i & 0x0F)));
            }
        }
        // Report data: first 32 bytes used by hook. Use a recognizable pattern.
        bytes32 rd = keccak256("tdx-validation-hook-test-reportdata");
        for (uint256 i = 0; i < 32; i++) {
            quote[OFFSET_REPORT_DATA + i] = rd[i];
        }

        return quote;
    }

    /// @dev Computes the (codeCommitment, dataCommitment) digests that the hook will
    ///      derive from a given quote, so tests can pass matching expected values
    ///      and isolate failure modes to the field they intend to exercise.
    function _expectedDigests(bytes memory quote) internal pure returns (bytes32 codeCommit, bytes32 dataCommit) {
        // Code commitment: keccak256(MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3) — 240 bytes.
        bytes memory ident = new bytes(240);
        for (uint256 i = 0; i < 48; i++) {
            ident[i] = quote[OFFSET_MRTD + i];
        }
        for (uint256 i = 0; i < 192; i++) {
            ident[48 + i] = quote[OFFSET_RTMR0 + i];
        }
        codeCommit = keccak256(ident);

        // Data commitment: first 32 bytes of REPORT_DATA, packed as bytes32.
        bytes32 rd;
        assembly {
            // quote is bytes memory: skip 32-byte length prefix, then offset within data.
            rd := mload(add(add(quote, 32), OFFSET_REPORT_DATA))
        }
        dataCommit = rd;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          STBN bundle-path tests                        //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Wire-format constants exposed to tests so we can mutate
    ///      bundle bytes at specific offsets without re-deriving them.
    uint256 internal constant BUNDLE_OFFSET_MAGIC = 0;
    uint256 internal constant BUNDLE_OFFSET_VERSION = 4;
    uint256 internal constant BUNDLE_OFFSET_FLAGS = 5;
    uint256 internal constant BUNDLE_OFFSET_VENDOR = 6; // 2 bytes BE
    uint256 internal constant BUNDLE_OFFSET_TDX_LEN = 8; // 4 bytes BE
    uint256 internal constant BUNDLE_OFFSET_TDX = 12;

    /// @dev qualifyingData in the embedded test bundles. Matches
    ///      sha256("tdx-bundle-test-data-commitment") computed by the
    ///      Go fixture generator (see contracts/test/dkg/fixtures).
    bytes32 internal constant FIXTURE_QUALIFYING_DATA =
        0x9e80f7b24630d8eba7e21ca0e691225d55563ce9b3cf3b49400240648950719d;

    /// @dev Loads the canonical direct-vendor bundle from disk.
    function _loadDirectBundle() internal view returns (bytes memory) {
        return _loadHex("test/dkg/fixtures/direct_bundle.hex");
    }

    /// @dev Loads the canonical paravisor-vendor bundle from disk.
    function _loadParavisorBundle() internal view returns (bytes memory) {
        return _loadHex("test/dkg/fixtures/paravisor_bundle.hex");
    }

    /// @dev Reads an ASCII hex file from `path` and returns the
    ///      decoded bytes. The file may contain trailing whitespace
    ///      which we strip; otherwise length must be even.
    function _loadHex(string memory path) internal view returns (bytes memory) {
        string memory s = vm.readFile(path);
        bytes memory raw = bytes(s);
        // Strip trailing whitespace (newline) — vm.parseBytes is strict.
        uint256 end = raw.length;
        while (end > 0 && (raw[end - 1] == 0x0A || raw[end - 1] == 0x0D || raw[end - 1] == 0x20)) {
            end--;
        }
        bytes memory clean = new bytes(end);
        for (uint256 i = 0; i < end; i++) {
            clean[i] = raw[i];
        }
        // Re-prepend "0x" so vm.parseBytes accepts the input.
        return vm.parseBytes(string.concat("0x", string(clean)));
    }

    /// @dev Computes the keccak256 code commitment over MRTD||RTMR0..3
    ///      from the V4 quote embedded inside a bundle. The V4 starts
    ///      at offset 12 (BUNDLE_OFFSET_TDX) once we know the length
    ///      is well-formed; the in-V4 offsets are the same as the
    ///      raw-V4 path's constants.
    function _bundleCodeCommitment(bytes memory bundle) internal pure returns (bytes32) {
        // Read the inner V4 length (BE uint32 at offset 8).
        uint256 v4Len = (uint256(uint8(bundle[8])) << 24) |
            (uint256(uint8(bundle[9])) << 16) |
            (uint256(uint8(bundle[10])) << 8) |
            uint256(uint8(bundle[11]));
        require(v4Len >= 632, "fixture v4 too short");
        // Build MRTD||RTMR0..3 by indexing into the bundle directly.
        bytes memory ident = new bytes(240);
        for (uint256 i = 0; i < 48; i++) {
            ident[i] = bundle[BUNDLE_OFFSET_TDX + OFFSET_MRTD + i];
        }
        for (uint256 i = 0; i < 192; i++) {
            ident[48 + i] = bundle[BUNDLE_OFFSET_TDX + OFFSET_RTMR0 + i];
        }
        return keccak256(ident);
    }

    /// @dev Returns a deep copy of `src` with byte at `offset` set
    ///      to `value`.
    function _patch1(bytes memory src, uint256 offset, bytes1 value) internal pure returns (bytes memory) {
        bytes memory out = new bytes(src.length);
        for (uint256 i = 0; i < src.length; i++) {
            out[i] = src[i];
        }
        out[offset] = value;
        return out;
    }

    /// @dev Returns a deep copy of `src` truncated to `newLen` bytes.
    function _truncate(bytes memory src, uint256 newLen) internal pure returns (bytes memory) {
        bytes memory out = new bytes(newLen);
        for (uint256 i = 0; i < newLen; i++) {
            out[i] = src[i];
        }
        return out;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Happy paths                               //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_BundleHappyPath_Direct() public {
        bytes memory bundle = _loadDirectBundle();
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        bool ok = hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
        assertTrue(ok);
        // Automata MUST have been called exactly once over the inner V4.
        assertEq(automata.callCount(), 1);
        // The bundle envelope is stripped before Automata is invoked —
        // assert the recorded quote is the inner V4 (632 bytes), not
        // the full bundle (1325 bytes).
        assertEq(automata.lastQuote().length, 632);
        assertEq(automata.lastTcbNum(), TCB_NUM);
    }

    function test_ValidateReport_BundleHappyPath_Paravisor() public {
        bytes memory bundle = _loadParavisorBundle();
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        bool ok = hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
        assertTrue(ok);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Magic dispatch                            //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_BundleWrongMagic_FallsThroughToRawV4() public {
        // Flip a single byte of the STBN magic. With magic broken,
        // the dispatch falls through to the raw-V4 path; that path
        // sees 0x00 at offset 4 (where TDX expects 0x81), so the
        // raw-V4 TEE-type guard fires.
        bytes memory bundle = _loadDirectBundle();
        bundle[0] = bytes1(uint8(0x53)); // already 'S'; flip last char of magic
        bundle = _patch1(bundle, 3, bytes1(uint8(0x00))); // 'N' -> 0x00
        vm.prank(dkg);
        // The bundle's offset 4..7 carries the bundle version + flags
        // bytes (0x01 0x01 0x00 0x00), not a TDX tee_type, so the
        // raw-V4 path's TEE-type guard rejects.
        vm.expectRevert(bytes("TDXValidationHook: Not a TDX quote"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                        Bundle structural rejection                     //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_BundleWrongVersion() public {
        bytes memory bundle = _patch1(_loadDirectBundle(), BUNDLE_OFFSET_VERSION, bytes1(uint8(0x02)));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXBundle: unsupported version"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    function test_ValidateReport_BundleReservedFlag() public {
        // Set bit1 of flags (reserved) in addition to TPM_PRESENT.
        bytes memory bundle = _patch1(_loadDirectBundle(), BUNDLE_OFFSET_FLAGS, bytes1(uint8(0x03)));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXBundle: reserved flag bits set"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    function test_ValidateReport_BundleMissingTPMPresent() public {
        // Clear TPM_PRESENT (bit0). All TPM sections in the bundle
        // are non-empty so this fails at the missing-TPM_PRESENT
        // check before the cap checks.
        bytes memory bundle = _patch1(_loadDirectBundle(), BUNDLE_OFFSET_FLAGS, bytes1(uint8(0x00)));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXBundle: missing TPM_PRESENT flag"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    function test_ValidateReport_BundleUnknownVendor() public {
        // Set vendor_tag to 0x0042 (not in {0x0000, 0x0001, 0xFFFF}).
        bytes memory bundle = _loadDirectBundle();
        bundle = _patch1(bundle, BUNDLE_OFFSET_VENDOR, bytes1(uint8(0x00)));
        bundle = _patch1(bundle, BUNDLE_OFFSET_VENDOR + 1, bytes1(uint8(0x42)));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXBundle: unknown vendor tag"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    function test_ValidateReport_BundleTruncated() public {
        // Truncate to bundle header + 4 bytes of v4 — not enough for
        // the declared tdx_v4_len.
        bytes memory bundle = _truncate(_loadDirectBundle(), 16);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXBundle: truncated at attest_len"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    function test_ValidateReport_BundleParavisorMissingRuntimeData() public {
        // Take the Direct bundle but flip the vendor tag to PARAVISOR.
        // The runtime_data length is 0 (direct vendor), violating the
        // paravisor invariant.
        bytes memory bundle = _loadDirectBundle();
        bundle = _patch1(bundle, BUNDLE_OFFSET_VENDOR + 1, bytes1(uint8(0x01)));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXBundle: paravisor missing runtime_data"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           V4-shape rejection                           //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_BundleInnerV4_NotTDX() public {
        // Patch the inner V4's tee_type byte at offset 4 (absolute
        // offset 12 + 4 = 16) so the TEE-type guard fires.
        bytes memory bundle = _patch1(_loadDirectBundle(), BUNDLE_OFFSET_TDX + 4, bytes1(uint8(0x00)));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Not a TDX quote"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    function test_ValidateReport_BundleInnerV4_UnsupportedVersion() public {
        // Flip inner V4 version from 4 to 3.
        bytes memory bundle = _patch1(_loadDirectBundle(), BUNDLE_OFFSET_TDX + 0, bytes1(uint8(0x03)));
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Unsupported quote version"));
        hook.validateReport(bytes32(uint256(1)), bytes32(uint256(2)), bundle, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                       Trust-anchor / hash mismatches                   //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_BundleAutomataFailure() public {
        automata.setShouldSucceed(false);
        bytes memory bundle = _loadDirectBundle();
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Attestation failed"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    function test_ValidateReport_BundleCodeCommitmentMismatch() public {
        bytes memory bundle = _loadDirectBundle();
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Code commitment does not match"));
        hook.validateReport(keccak256("wrong-code"), FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    function test_ValidateReport_BundleAKBindingMismatch_Direct() public {
        // Corrupt one byte of the inner V4.report_data so SHA256(AK_pub)
        // no longer matches.
        bytes memory bundle = _loadDirectBundle();
        bundle = _patch1(bundle, BUNDLE_OFFSET_TDX + OFFSET_REPORT_DATA, bytes1(uint8(0xDE)));
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: AK binding (direct) mismatch"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    function test_ValidateReport_BundleAKBindingMismatch_Paravisor() public {
        bytes memory bundle = _loadParavisorBundle();
        bundle = _patch1(bundle, BUNDLE_OFFSET_TDX + OFFSET_REPORT_DATA, bytes1(uint8(0xDE)));
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: AK binding (paravisor) mismatch"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    function test_ValidateReport_BundleReportDataReservedNonZero() public {
        // Tamper a byte in V4.report_data[32:64] (must be zero).
        bytes memory bundle = _loadDirectBundle();
        bundle = _patch1(bundle, BUNDLE_OFFSET_TDX + OFFSET_REPORT_DATA + 32, bytes1(uint8(0xFF)));
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: report_data[32:64] must be zero"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                       TPM signature rejections                         //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Computes the absolute byte offset of the tpm_sig section
    ///      inside the bundle (after tdx_v4 + tpm_attest).
    function _bundleTpmSigOffset(bytes memory bundle) internal pure returns (uint256 off) {
        // tdx_v4 length at offset 8.
        uint256 v4Len = (uint256(uint8(bundle[8])) << 24) |
            (uint256(uint8(bundle[9])) << 16) |
            (uint256(uint8(bundle[10])) << 8) |
            uint256(uint8(bundle[11]));
        // Bundle header (12) + tdx_v4 + 4 (attest_len) + attest_len + 4 (sig_len) = sig start.
        uint256 attestLenOff = 12 + v4Len;
        uint256 attestLen = (uint256(uint8(bundle[attestLenOff])) << 24) |
            (uint256(uint8(bundle[attestLenOff + 1])) << 16) |
            (uint256(uint8(bundle[attestLenOff + 2])) << 8) |
            uint256(uint8(bundle[attestLenOff + 3]));
        // sig section starts after the 4-byte sig_len field.
        off = attestLenOff + 4 + attestLen + 4;
    }

    function _bundleTpmAttestOffset(bytes memory bundle) internal pure returns (uint256 off) {
        uint256 v4Len = (uint256(uint8(bundle[8])) << 24) |
            (uint256(uint8(bundle[9])) << 16) |
            (uint256(uint8(bundle[10])) << 8) |
            uint256(uint8(bundle[11]));
        // tpm_attest starts after bundle header + tdx_v4 + 4-byte attest_len.
        off = 12 + v4Len + 4;
    }

    function test_ValidateReport_BundleSigWrongAlg() public {
        // Replace algId from 0x0014 (RSASSA) to 0x0016 (RSAPSS).
        bytes memory bundle = _loadDirectBundle();
        uint256 sigOff = _bundleTpmSigOffset(bundle);
        bundle = _patch1(bundle, sigOff + 0, bytes1(uint8(0x00)));
        bundle = _patch1(bundle, sigOff + 1, bytes1(uint8(0x16)));
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: TPM sig alg != RSASSA"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    function test_ValidateReport_BundleSigWrongHashAlg() public {
        // Replace hashAlg from 0x000B (SHA-256) to 0x000C (SHA-384).
        bytes memory bundle = _loadDirectBundle();
        uint256 sigOff = _bundleTpmSigOffset(bundle);
        bundle = _patch1(bundle, sigOff + 2, bytes1(uint8(0x00)));
        bundle = _patch1(bundle, sigOff + 3, bytes1(uint8(0x0C)));
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: TPM hash alg != SHA-256"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    function test_ValidateReport_BundleSigCorrupted() public {
        // Flip a byte deep inside the RSA signature blob; modexp still
        // returns 256 bytes but the PKCS#1 v1.5 envelope no longer
        // matches the expected digest.
        bytes memory bundle = _loadDirectBundle();
        uint256 sigOff = _bundleTpmSigOffset(bundle);
        // Sig bytes start 6 bytes after sigOff (algId+hashAlg+sigLen).
        bundle[sigOff + 6 + 100] ^= bytes1(uint8(0x55));
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: TPM sig verify failed"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                  TPMS_ATTEST tampering and qualifyingData              //
    //////////////////////////////////////////////////////////////////////////*/

    function test_ValidateReport_BundleAttestBadMagic() public {
        // Flip the high byte of TPM_GENERATED_VALUE (0xFF) — the
        // signature verifies first, so we ALSO must re-sign… but
        // since we cannot, the test instead expects the signature to
        // fail before the magic check. This locks in the ordering of
        // the verification ladder: signature MUST come before
        // structural parse, so an attacker cannot strip the magic
        // and forge a synthetic attestation.
        bytes memory bundle = _loadDirectBundle();
        uint256 attOff = _bundleTpmAttestOffset(bundle);
        bundle = _patch1(bundle, attOff, bytes1(uint8(0x00)));
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: TPM sig verify failed"));
        hook.validateReport(codeCommit, FIXTURE_QUALIFYING_DATA, bundle, "");
    }

    function test_ValidateReport_BundleQualifyingDataMismatch() public {
        // Happy-path bundle, but caller passes the wrong expected data.
        bytes memory bundle = _loadDirectBundle();
        bytes32 codeCommit = _bundleCodeCommitment(bundle);
        bytes32 wrong = keccak256("not-the-fixture-commitment");
        vm.prank(dkg);
        vm.expectRevert(bytes("TDXValidationHook: Data commitment does not match"));
        hook.validateReport(codeCommit, wrong, bundle, "");
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          Library-level coverage                        //
    //////////////////////////////////////////////////////////////////////////*/

    function test_RSASSAVerify_BadSPKILength() public {
        bytes memory der = new bytes(100);
        vm.expectRevert(bytes("RSASSAVerify: bad DER length"));
        this.rsaExtractModulusTrampoline(der);
    }

    function test_RSASSAVerify_BadSPKIPrefix() public {
        bytes memory der = new bytes(294);
        // Leave prefix all zero — fails the prefix check.
        vm.expectRevert(bytes("RSASSAVerify: bad SPKI prefix"));
        this.rsaExtractModulusTrampoline(der);
    }

    function test_RSASSAVerify_BadSPKIExponent() public {
        bytes memory der = new bytes(294);
        bytes memory good = (hex"30820122300D06092A864886F70D01010105000382010F003082010A0282010100");
        for (uint256 i = 0; i < good.length; i++) {
            der[i] = good[i];
        }
        // Modulus bytes left zero. Exponent suffix at 33 + 256 = 289.
        // Default zero violates the 02 03 01 00 01 expectation.
        vm.expectRevert(bytes("RSASSAVerify: bad SPKI exponent"));
        this.rsaExtractModulusTrampoline(der);
    }

    function test_TPM2Attest_BadMagic() public {
        bytes memory attest = new bytes(64);
        // magic is 0x00... not 0xFF544347.
        vm.expectRevert(bytes("TPM2Attest: bad magic"));
        this.tpmExtractQualifyingDataTrampoline(attest);
    }

    function test_TPM2Attest_NotQuote() public {
        bytes memory attest = new bytes(64);
        // Set magic correctly, type to TPM_ST_ATTEST_CERTIFY (0x8017).
        attest[0] = bytes1(uint8(0xFF));
        attest[1] = bytes1(uint8(0x54));
        attest[2] = bytes1(uint8(0x43));
        attest[3] = bytes1(uint8(0x47));
        attest[4] = bytes1(uint8(0x80));
        attest[5] = bytes1(uint8(0x17));
        vm.expectRevert(bytes("TPM2Attest: not a quote"));
        this.tpmExtractQualifyingDataTrampoline(attest);
    }

    /// @dev External trampoline for RSASSAVerify.extractRSA2048Modulus
    ///      so vm.expectRevert sees the revert at the expected depth.
    function rsaExtractModulusTrampoline(bytes memory der) external pure returns (bytes memory) {
        return RSASSAVerify.extractRSA2048Modulus(der);
    }

    /// @dev External trampoline for TPM2Attest.extractQualifyingData.
    function tpmExtractQualifyingDataTrampoline(bytes memory attest) external pure returns (bytes memory) {
        return TPM2Attest.extractQualifyingData(attest);
    }

    function test_TDXBundle_HasMagic() public {
        bytes memory good = abi.encodePacked(bytes4(0x5354424E), bytes("rest"));
        bytes memory bad = abi.encodePacked(bytes4(0x53544246), bytes("rest"));
        bytes memory tooShort = abi.encodePacked(bytes3(0x535342));
        assertTrue(this.hasMagicTrampoline(good));
        assertFalse(this.hasMagicTrampoline(bad));
        assertFalse(this.hasMagicTrampoline(tooShort));
    }

    /// @dev External trampoline so we can exercise the calldata-only
    ///      `TDXBundle.hasMagic` from a memory-typed call site.
    function hasMagicTrampoline(bytes calldata buf) external pure returns (bool) {
        return TDXBundle.hasMagic(buf);
    }
}
