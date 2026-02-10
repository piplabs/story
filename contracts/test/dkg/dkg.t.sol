// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { DKG } from "../../src/protocol/DKG.sol";
import { Test } from "../utils/Test.sol";

contract DKGTest is Test {
    DKG dkg;

    bytes32 codeCommitment = hex"4d53ef0428afd0bc343e4c0ca19efd05ad5d5747b4b230491c5e1237ca294739";
    uint32 round = 1;
    uint32 index = 0;
    uint32 total = 3;
    uint32 threshold = 2;
    uint64 startBlockHeight = 1000;
    bytes32 startBlockHash = hex"abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd";
    bytes rawQuote =
        hex"beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef";
    bytes dkgPubKey = hex"dead";
    bytes globalPubKey = hex"ef01";
    bytes32 participantsRoot = hex"7c3d9a2f41e8b6d0f52c1b9d84a6e3f1c7b02d9e5a1f4c8b3e6d19f0a7b4c2d1";
    bytes invalid_signature =
        hex"1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111";

    // validator 0
    address validator = address(0x061F9f80b3cf1a5cd6769EC0DB77D2Be50A3fa8f);
    bytes commPubKey =
        hex"ecbfad7a514da8d3bb3dc8d7c5a171ec0ceb04d36f045e519f20cc31ca7c78292d42c8187c9de7238271344bbdf62f5748b5ff339eb0f52cb8463cdc41f595f2";
    bytes finalizeDKG_signature =
        hex"70a1d8b96be91078aa807ac9f26127147fd147b0dc68e84676761ad2b70a1b604ea07f56e94d77a041119143cf994b78b64f462ef56aa4d6d289eb352f086d9d1b";

    // validator 1
    address validator1 = address(0x591942F67Cf7d6104aD2f41Be899713fCb60dceB);
    bytes commPubKey1 =
        hex"55a921e6dce5eef957df0b9a322122ed81d0c2a98deadd20e6ff17fdc714143cbf376aab658436baf321b8cce34a06078f2082c9762445d98e51d6808950f600";
    bytes finalizeDKG_signature1 =
        hex"c945b681a8c3872b2d50b9ab2db56f5eead49902d02e99a8054b4b0015fb5aff319d4d704e415e87a89abfc1be1ec9f1fa73db92842fb3967d6d2e780e24c9601b";

    // validator 2
    address validator2 = address(0xEB8e62E11504B3961DF295bA363385B241606355);
    bytes commPubKey2 =
        hex"7f28a5e6d5ad9b315b7241b73c4c5f9d68d3fd2051d52e0ec39b70a557acdd3d5d7e53f87bb3c6c74be59775247fc280f827f36e3fa45bcf3e95aa6f94e3f8c0";
    bytes finalizeDKG_signature2 =
        hex"b3ab103cecd44f3202f1e68c795493148c7eeb1e30b9b821ffe9e85cd6d3ea48258b541e76d75e810bb2382ce0b2a010585f1fb1bdf338f0e7d3def01c0964ae1c";

    // validator 3
    address validator3 = address(0xd80E4c5D255c28305572D452D6f381f3DA06fE0b);
    bytes commPubKey3 =
        hex"4971528c66918f5eb181a70ffd455569aadec5aeee4dbd01d8dc208241381b137f92101076e3b7cf2c5c2e0fb63793abef2360570429a1b597da0422a4aaf10e";
    bytes finalizeDKG_signature3 =
        hex"31ab1b5bfc86648bba3c5783a3795df750c8517009b623a0281eda3520ad240b505f2a83ecde1dc9a205defcaebe0ab20cf41122f596f4ad4e188d4fde7321101b";

    function setUp() public override {
        dkg = new DKG(codeCommitment);
    }

    function testThreeNodeDKG_Success() public {
        bytes[] memory publicCoeffs = new bytes[](3);
        publicCoeffs[0] = hex"123456789abcdef0";
        publicCoeffs[1] = hex"aabbccddeeff11223344";
        publicCoeffs[2] = hex"deadbeefcafebabe01020304";

        // 1.initialize DKG
        // 1.1 validator 0
        vm.prank(validator1);
        dkg.initializeDKG(round, codeCommitment, startBlockHeight, startBlockHash, dkgPubKey, commPubKey1, rawQuote);
        DKG.NodeInfo memory info = dkg.getNodeInfo(codeCommitment, round, validator1);
        assertEq(info.dkgPubKey, dkgPubKey);
        assertEq(info.commPubKey, commPubKey1);
        assertEq(info.rawQuote, rawQuote);
        assertEq(uint8(info.nodeStatus), 1); // Registered
        // 1.2 validator 1
        vm.prank(validator2);
        dkg.initializeDKG(round, codeCommitment, startBlockHeight, startBlockHash, dkgPubKey, commPubKey2, rawQuote);
        info = dkg.getNodeInfo(codeCommitment, round, validator2);
        assertEq(info.dkgPubKey, dkgPubKey);
        assertEq(info.commPubKey, commPubKey2);
        assertEq(info.rawQuote, rawQuote);
        assertEq(uint8(info.nodeStatus), 1); // Registered
        // 1.3 validator 2
        vm.prank(validator3);
        dkg.initializeDKG(round, codeCommitment, startBlockHeight, startBlockHash, dkgPubKey, commPubKey3, rawQuote);
        info = dkg.getNodeInfo(codeCommitment, round, validator3);
        assertEq(info.dkgPubKey, dkgPubKey);
        assertEq(info.commPubKey, commPubKey3);
        assertEq(info.rawQuote, rawQuote);
        assertEq(uint8(info.nodeStatus), 1); // Registered

        // 3. finalize DKG
        vm.prank(validator1);
        dkg.finalizeDKG(round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, finalizeDKG_signature1);
        info = dkg.getNodeInfo(codeCommitment, round, validator1);
        assertEq(uint8(info.nodeStatus), 3); // Finalized

        vm.prank(validator2);
        dkg.finalizeDKG(round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, finalizeDKG_signature2);
        info = dkg.getNodeInfo(codeCommitment, round, validator2);
        assertEq(uint8(info.nodeStatus), 3); // Finalized

        vm.prank(validator3);
        dkg.finalizeDKG(round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, finalizeDKG_signature3);
        info = dkg.getNodeInfo(codeCommitment, round, validator3);
        assertEq(uint8(info.nodeStatus), 3); // Finalized
    }

    function testFinalizeDKG_RevertIfInvalidSignature() public {
        bytes[] memory publicCoeffs = new bytes[](3);
        publicCoeffs[0] = hex"123456789abcdef0";
        publicCoeffs[1] = hex"aabbccddeeff11223344";
        publicCoeffs[2] = hex"deadbeefcafebabe01020304";

        // initialize DKG
        vm.prank(validator);
        dkg.initializeDKG(round, codeCommitment, startBlockHeight, startBlockHash, dkgPubKey, commPubKey, rawQuote);

        // finalize DKG with wrong signature
        vm.prank(validator);
        vm.expectRevert("ECDSAInvalidSignature()");
        dkg.finalizeDKG(round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, invalid_signature);
    }
}
