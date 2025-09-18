// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { DKG} from "../../src/protocol/DKG.sol";
import { Test } from "../utils/Test.sol";

contract DKGTest is Test {
    DKG dkg;

    bytes mrenclave = hex"1234";
    uint32 round = 1;
    uint32 index = 0;
    bytes remoteReport = hex"beef";
    bytes commitments = hex"cafe";

    // How to create a new account:
    // cast wallet new & cast wallet public-key --private-key <private-key>

    // private key: 0x1586935bf3a4aa4e40cb0782227d1d302082f727677f2ec4d37054bc92f0079f
    address validator = address(0x061F9f80b3cf1a5cd6769EC0DB77D2Be50A3fa8f);
    bytes pubKey = hex"ecbfad7a514da8d3bb3dc8d7c5a171ec0ceb04d36f045e519f20cc31ca7c78292d42c8187c9de7238271344bbdf62f5748b5ff339eb0f52cb8463cdc41f595f2";

    // How to generate commitment_signature:
    // 1.cast keccak "0xcafe"
    //   => 0x72318c618151a897569554720f8f1717a3da723042fb73893c064da11b308ae9
    // 2.cast wallet sign --private-key 0x1586935bf3a4aa4e40cb0782227d1d302082f727677f2ec4d37054bc92f0079f 0x72318c618151a897569554720f8f1717a3da723042fb73893c064da11b308ae9
    //   => 0xc70add942114af5f0eb992351aec5461db0705ca5ca7834e6d5483faea7f116f41860f154176c0b0210d3e98796fcd425a688aeadb299719528f3d6c9df58a351b
    bytes commitment_signature = hex"c70add942114af5f0eb992351aec5461db0705ca5ca7834e6d5483faea7f116f41860f154176c0b0210d3e98796fcd425a688aeadb299719528f3d6c9df58a351b";



    // How to generate finalizeDKG_signature:
    // 1.cast abi-encode --packed "tuple(uint32,uint32,bool,bytes)" 1 0 true 0x1234
    //   => 0x0000000100000000011234
    // 2. cast keccak --hex "0x0000000100000000011234"
    //   => 0x0c47aae97f3d3737a7a2257ea1435b003f5ba1f21168bda5f111b963076a4f1b
    // 3. cast wallet sign --private-key 0x1586935bf3a4aa4e40cb0782227d1d302082f727677f2ec4d37054bc92f0079f 0x0c47aae97f3d3737a7a2257ea1435b003f5ba1f21168bda5f111b963076a4f1b
    //   => a187ec49a839d9f130def4efd2fb0353a5cd7b6a1a32ec0e7a360927f4a4fd4c3bba370f55cfdc053c7a6fda5298da731ec02b3907eacc9be289eb338156fd891c
    bytes finalizeDKG_signature = hex"a187ec49a839d9f130def4efd2fb0353a5cd7b6a1a32ec0e7a360927f4a4fd4c3bba370f55cfdc053c7a6fda5298da731ec02b3907eacc9be289eb338156fd891c";

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
        dkg.initializeDKG(round, mrenclave, pubKey, remoteReport);
        DKG.NodeInfo memory info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(info.validator, validator);
        assertEq(info.pubKey, pubKey);
        assertEq(info.remoteReport, remoteReport);
        assertEq(info.finalized, false);

        // update commitments
        vm.prank(validator);
        dkg.updateDKGCommitments(round, 1, 1, index, mrenclave, commitments, commitment_signature);
        info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(keccak256(info.commitments), keccak256(commitments));

        // finalize DKG
        vm.prank(validator);
        dkg.finalizeDKG(round, index, true, mrenclave, finalizeDKG_signature);
        info = dkg.getNodeInfo(mrenclave, round, index);
        assertEq(info.finalized, true);
    }

    function testInitializeDKG_InvalidPubKey() public {
        vm.prank(address(0x999));
        vm.expectRevert("Invalid pubKey for sender");
        dkg.initializeDKG(round, mrenclave, pubKey, remoteReport);
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
        dkg.initializeDKG(round, mrenclave, pubKey, remoteReport);

        // update DKG commitments with wrong signature
        vm.prank(validator);
        vm.expectRevert("Invalid signature");
        dkg.updateDKGCommitments(round, 1, 1, index, mrenclave, commitments, invalid_signature);
    }

    function testFinalizeDKG_RevertIfInvalidSignature() public {
         // initialize DKG
        vm.prank(validator);
        dkg.initializeDKG(round, mrenclave, pubKey, remoteReport);

        // finalize DKG with wrong signature
        vm.prank(validator);
        vm.expectRevert("Invalid signature");
        dkg.finalizeDKG(round, index, true, mrenclave, invalid_signature);
    }
}