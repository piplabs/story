// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { DKG} from "../../src/protocol/DKG.sol";
import { Test } from "../utils/Test.sol";

contract DKGTest is Test {
    DKG dkg;

    bytes mrenclave = hex"1234";
    uint32 round = 1;
    uint32 index = 0;
    bytes rawQuote = hex"beef";
    bytes commitments = hex"cafe";
    bytes dkgPubKey = hex"dead";
    bytes globalPubKey = hex"ef01";

    // How to create a new account:
    // cast wallet new & cast wallet public-key --private-key <private-key>

    // private key: 0x1586935bf3a4aa4e40cb0782227d1d302082f727677f2ec4d37054bc92f0079f
    address validator = address(0x061F9f80b3cf1a5cd6769EC0DB77D2Be50A3fa8f);
    bytes commPubKey = hex"ecbfad7a514da8d3bb3dc8d7c5a171ec0ceb04d36f045e519f20cc31ca7c78292d42c8187c9de7238271344bbdf62f5748b5ff339eb0f52cb8463cdc41f595f2";

    // How to generate commitment_signature:
    // 1.cast keccak "0xcafe"
    //   => 0x72318c618151a897569554720f8f1717a3da723042fb73893c064da11b308ae9
    // 2.cast wallet sign --private-key 0x1586935bf3a4aa4e40cb0782227d1d302082f727677f2ec4d37054bc92f0079f 0x72318c618151a897569554720f8f1717a3da723042fb73893c064da11b308ae9
    //   => 0xc70add942114af5f0eb992351aec5461db0705ca5ca7834e6d5483faea7f116f41860f154176c0b0210d3e98796fcd425a688aeadb299719528f3d6c9df58a351b
    bytes commitment_signature = hex"c70add942114af5f0eb992351aec5461db0705ca5ca7834e6d5483faea7f116f41860f154176c0b0210d3e98796fcd425a688aeadb299719528f3d6c9df58a351b";



    // How to generate finalizeDKG_signature:
    // 1.cast abi-encode --packed "tuple(uint32,uint32,bool,bytes,bytes)" 1 0 true 0x1234 0xef01
    //   => 0x0000000100000000011234ef01
    // 2. cast keccak --hex "0x0000000100000000011234ef01"
    //   => 0x21af19238192676c10b75e17db2867695680cee6e09418ebf2ef1e524398cf6e
    // 3. cast wallet sign --private-key 0x1586935bf3a4aa4e40cb0782227d1d302082f727677f2ec4d37054bc92f0079f 0x21af19238192676c10b75e17db2867695680cee6e09418ebf2ef1e524398cf6e
    //   => 0x70a1d8b96be91078aa807ac9f26127147fd147b0dc68e84676761ad2b70a1b604ea07f56e94d77a041119143cf994b78b64f462ef56aa4d6d289eb352f086d9d1b
    bytes finalizeDKG_signature = hex"70a1d8b96be91078aa807ac9f26127147fd147b0dc68e84676761ad2b70a1b604ea07f56e94d77a041119143cf994b78b64f462ef56aa4d6d289eb352f086d9d1b";

    bytes invalid_signature = hex"1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111";

    function setUp() override public {
        dkg = new DKG();
        // Simulate validator in active set
        vm.prank(validator);
        dkg.submitActiveValSet(round, mrenclave, new address[](1));
    }

    function testDKG_Success() public {
        // initialize DKG
        vm.prank(validator);
        dkg.initializeDKG(round, mrenclave, dkgPubKey, commPubKey, rawQuote);
        DKG.NodeInfo memory info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(info.validator, validator);
        assertEq(info.dkgPubKey, dkgPubKey);
        assertEq(info.commPubKey, commPubKey);
        assertEq(info.rawQuote, rawQuote);
        assertEq(info.finalized, false);

        // update commitments
        vm.prank(validator);
        dkg.updateDKGCommitments(round, 1, 1, index, mrenclave, commitments, commitment_signature);
        info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(keccak256(info.commitments), keccak256(commitments));

        // finalize DKG
        vm.prank(validator);
        dkg.finalizeDKG(round, index, true, mrenclave, finalizeDKG_signature, globalPubKey);
        info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(info.finalized, true);
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
        dkg.finalizeDKG(round, index, true, mrenclave, invalid_signature, globalPubKey);
    }
}