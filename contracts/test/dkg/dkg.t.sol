// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { DKG} from "../../src/protocol/DKG.sol";
import { Test } from "../utils/Test.sol";

contract DKGTest is Test {
    DKG dkg;

    bytes mrenclave = hex"1234";
    uint32 round = 1;
    uint32 index = 0;
    uint32 total = 3;
    uint32 threshold = 2;
    bool finalized = true;
    bytes rawQuote = hex"beef";
    bytes commitments = hex"cafe";
    bytes dkgPubKey = hex"dead";
    bytes globalPubKey = hex"ef01";
    bytes invalid_signature = hex"1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111";

    // validator 0
    address validator = address(0x061F9f80b3cf1a5cd6769EC0DB77D2Be50A3fa8f);
    bytes commPubKey = hex"ecbfad7a514da8d3bb3dc8d7c5a171ec0ceb04d36f045e519f20cc31ca7c78292d42c8187c9de7238271344bbdf62f5748b5ff339eb0f52cb8463cdc41f595f2";
    bytes update_commitment_signature = hex"f1f3b9ac511d2802ed72e09c45a06a967b5b992df6e005a3f011e6cd86f233fb506cbeb948609a36bb48e83022bd06536f9acb55ffba61cca61917bf0946cfed1b";
    bytes finalizeDKG_signature = hex"70a1d8b96be91078aa807ac9f26127147fd147b0dc68e84676761ad2b70a1b604ea07f56e94d77a041119143cf994b78b64f462ef56aa4d6d289eb352f086d9d1b";


    // validator 1
    address validator1 = address(0x8d16983bCa5509C4cca4A73ce7bF5034A152dAb8);
    bytes commPubKey1 = hex"c69a9c874eaf034562948a2c771e3ceed4365dc442c838617fc2b09e634198ba4497f208889c626d40941235e29774e8c290b421f2f8668e2214477986072689";
    bytes update_commitment_signature1 = hex"1d47a7c8164e1d79c7fbb8e44a1fa2222e1863fa1d498ede0d443256cb38225128fa3835e9e563967d5d66a56e40da8cf154e0ee3277e119a6dbc688c4da6e5e1c";
    bytes finalizeDKG_signature1 = hex"1d69f0f6988cbb69c2957fa625dc610639bb0f4380f53024a48d01840d4fce0a07e5b1b6c715784d521977e02cfabbd83235104ccb1a1b1eaec108db968c3c841c";

    // validator 2
    address validator2 = address(0x1BaAfE5C84f0df458362cb05Ac3a1C0dd3585a12);
    bytes commPubKey2 = hex"27b35c5949ef1e72336a61a0252f5a8eb08a9c935b8d93a6f1eba4e176bffbe9899b80725d7c733b5ed91f0035604f0d186992ec2d000fe2c6bd86be5893ab86";
    bytes update_commitment_signature2 = hex"820d5a4d60be9ceef3878a3168f14d5971474979e8db60388a4902e120f0995e3d71342ce16f5c5d02f2885332c9c5bdd0c9a9163addf427a02f4ad1738d9dad1c";
    bytes finalizeDKG_signature2 = hex"ca42aeb5400df43b98d2bd3cc97769aef486ca56a17e526c4dbd76731c88c6ec6e1c547174cd5719ab98cdf19ca556fc65838134f4cab5a4ac039f0f302c5ea81c";

    // validator 3
    address validator3 = address(0x4DdeCF53cafE924A40461ee1AE19d51363c25686);
    bytes commPubKey3 = hex"bafa47a40e19c197ab18f39ef41da4cba37e66a230f2dc7cee9b8ccbef5e9548231809951d60a0ef7c80f7ffdae10fa173ebcef5761f8be6db157e2e3b856951";
    bytes update_commitment_signature3 = hex"776d014812e838218ada7dde2bcc60b3e2f0ff787391c31a1e09d83ff150e95206b1314a7806cbce41777c575627315672de7921a8959da4b00dbb84eaa62a381c";
    bytes finalizeDKG_signature3 = hex"d5114c69a70f800fc3af7459db55039fac78d2301bc685b36b60becd260dad9f17ef076bbe966deb5a2710daca2717aad3da3e0035ca202d26a1d4a61c10b7f41b";

    function setUp() override public {
        dkg = new DKG();
    }

    function testThreeNodeDKG_Success() public {
        // 1.initialize DKG
        // 1.1 validator 0
        vm.prank(validator1);
        dkg.initializeDKG(round, mrenclave, dkgPubKey, commPubKey1, rawQuote);
        DKG.NodeInfo memory info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(info.validator, validator1);
        assertEq(info.dkgPubKey, dkgPubKey);
        assertEq(info.commPubKey, commPubKey1);
        assertEq(info.rawQuote, rawQuote);
        assertEq(info.finalized, false);
        // 1.2 validator 1
        vm.prank(validator2);
        dkg.initializeDKG(round, mrenclave, dkgPubKey, commPubKey2, rawQuote);
        info = dkg.getNodeInfo(mrenclave, round, index+1);
        assertEq(info.validator, validator2);
        assertEq(info.dkgPubKey, dkgPubKey);
        assertEq(info.commPubKey, commPubKey2);
        assertEq(info.rawQuote, rawQuote);
        assertEq(info.finalized, false);
        // 1.3 validator 2
        vm.prank(validator3);
        dkg.initializeDKG(round, mrenclave, dkgPubKey, commPubKey3, rawQuote);
        info = dkg.getNodeInfo(mrenclave, round, index+2);
        assertEq(info.validator, validator3);
        assertEq(info.dkgPubKey, dkgPubKey);
        assertEq(info.commPubKey, commPubKey3);
        assertEq(info.rawQuote, rawQuote);
        assertEq(info.finalized, false);

        assertTrue(keccak256(dkg.getGlobalPubKey(mrenclave, round)) != keccak256(globalPubKey));

        // 2. update commitments
        vm.prank(validator1);
        dkg.updateDKGCommitments(round, total, threshold, index, mrenclave, commitments, update_commitment_signature1);
        info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(keccak256(info.commitments), keccak256(commitments));

        vm.prank(validator2);
        dkg.updateDKGCommitments(round, total, threshold, index+1, mrenclave, commitments, update_commitment_signature2);
        info = dkg.getNodeInfo(mrenclave, round, index+1);
        assertEq(keccak256(info.commitments), keccak256(commitments));  

        vm.prank(validator3);
        dkg.updateDKGCommitments(round, total, threshold, index+2, mrenclave, commitments, update_commitment_signature3);
        info = dkg.getNodeInfo(mrenclave, round, index+2);
        assertEq(keccak256(info.commitments), keccak256(commitments));

        assertTrue(keccak256(dkg.getGlobalPubKey(mrenclave, round)) != keccak256(globalPubKey));

        // 3. finalize DKG
        vm.prank(validator1);
        dkg.finalizeDKG(round, index, finalized, mrenclave, globalPubKey, finalizeDKG_signature1);
        info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(info.finalized, true);
        assertTrue(keccak256(dkg.getGlobalPubKey(mrenclave, round)) != keccak256(globalPubKey));

        vm.prank(validator2);
        dkg.finalizeDKG(round, index+1, finalized, mrenclave, globalPubKey, finalizeDKG_signature2);
        info = dkg.getNodeInfo(mrenclave, round, index+1);
        assertEq(info.finalized, true);
        assertEq(dkg.getGlobalPubKey(mrenclave, round), globalPubKey);

        vm.prank(validator3);
        dkg.finalizeDKG(round, index+2, finalized, mrenclave, globalPubKey, finalizeDKG_signature3);
        info = dkg.getNodeInfo(mrenclave, round, index+2);
        assertEq(info.finalized, true);
        assertEq(dkg.getGlobalPubKey(mrenclave, round), globalPubKey);
    }

    function testUpdateDKGCommitments_RevertIfNotSender() public {
        vm.expectRevert("Invalid sender");
        dkg.updateDKGCommitments(round, 1, 1, 0, mrenclave, hex"cafe", hex"dead");
    }

    function testUpdateDKGCommitments_RevertInvalidatedNode() public {
        // TODO: implement
    }

    function testUpdateDKGCommitments_RevertIfInvalidSignature() public {
        // initialize DKG
        vm.prank(validator);
        dkg.initializeDKG(round, mrenclave,dkgPubKey, commPubKey, rawQuote);

        // update DKG commitments with wrong signature
        vm.prank(validator);
        vm.expectRevert("ECDSAInvalidSignature()");
        dkg.updateDKGCommitments(round, 1, 1, index, mrenclave, commitments, invalid_signature);
    }

    function testFinalizeDKG_RevertIfInvalidSignature() public {
         // initialize DKG
        vm.prank(validator);
        dkg.initializeDKG(round, mrenclave,dkgPubKey, commPubKey, rawQuote);

        // finalize DKG with wrong signature
        vm.prank(validator);
        vm.expectRevert("ECDSAInvalidSignature()");
        dkg.finalizeDKG(round, index, true, mrenclave, globalPubKey, invalid_signature);
    }
}