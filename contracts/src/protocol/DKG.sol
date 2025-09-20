// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { IDKG } from "../interfaces/IDKG.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

/**
 * @title DKG - Distributed Key Generation Contract
 * @dev Core contract for managing DKG-related state and operations
 */
contract DKG is IDKG {
    mapping(bytes mrenclave => mapping(uint32 round => mapping(uint32 index => NodeInfo))) public dkgNodeInfos;
    mapping(bytes mrenclave => mapping(uint32 round => mapping(address validator => bool))) public valSets;
    mapping(bytes mrenclave => mapping(uint32 round => 
    mapping(uint32 index => mapping(address complainant => bool)))) public dealComplaints;
    mapping(bytes mrenclave => mapping(uint32 round => uint32 nodeCount)) public nodeCount;

    mapping(bytes mrenclave => mapping(uint32 round => mapping(bytes globalPubKeyCandidates => uint32 votes))) public votes;
    mapping(bytes mrenclave => mapping(uint32 round => RoundInfo roundInfo)) public roundInfo;

    constructor() {}

    function initializeDKG(
        uint32 round,
        bytes calldata mrenclave,
        bytes calldata dkgPubKey,
        bytes calldata commPubKey,
        bytes calldata rawQuote
    ) external {
        require(
            valSets[mrenclave][round][msg.sender] || _isActiveValSetSubmitted(mrenclave, round),
            "Validator not in active set"
        );

        require(
            _verifyRemoteAttestation(rawQuote, msg.sender, round, dkgPubKey),
            "Invalid remote attestation"
        );

        uint32 index = nodeCount[mrenclave][round];
        nodeCount[mrenclave][round]++;

        dkgNodeInfos[mrenclave][round][index] = NodeInfo({
            index: index,
            validator: msg.sender,
            dkgPubKey: dkgPubKey,
            commPubKey: commPubKey,
            rawQuote: rawQuote,
            commitments: "",
            chalStatus: ChallengeStatus.NotChallenged,
            finalized: false
        });

        emit DKGInitialized(
            msg.sender,
            mrenclave,
            round,
            index,
            dkgPubKey,
            commPubKey,
            rawQuote
        );
    }

    function updateDKGCommitments(
        uint32 round,
        uint32 total,
        uint32 threshold,
        uint32 index,
        bytes calldata mrenclave,
        bytes calldata commitments,
        bytes calldata signature
    ) external {
        NodeInfo storage node = dkgNodeInfos[mrenclave][round][index];
        require(node.validator == msg.sender, "Invalid sender");
        require(node.chalStatus != ChallengeStatus.Invalidated, "Node was invalidated");
        require(_verifyCommitmentSignature(node.commPubKey, round, total, threshold, index, mrenclave, commitments, signature), "Invalid commitment signature");

        node.commitments = commitments;

        // TODO: now we assume all validators submit the same total and threshold
        // in the future, we handle the case where they are different
        roundInfo[mrenclave][round] = RoundInfo({
            total: total,
            threshold: threshold,
            globalPubKey: ""
        });

        emit DKGCommitmentsUpdated(
            msg.sender,
            round,
            total,
            threshold,
            index,
            commitments,
            signature,
            mrenclave
        );
    }

    function finalizeDKG(
        uint32 round,
        uint32 index,
        bool finalized,
        bytes calldata mrenclave,
        bytes calldata globalPubKey,
        bytes calldata signature
    ) external {
        NodeInfo storage node = dkgNodeInfos[mrenclave][round][index];
        require(node.validator == msg.sender, "Invalid sender");
        require(node.chalStatus != ChallengeStatus.Invalidated, "Node was invalidated");
        require(
            _verifyFinalizationSignature(node.commPubKey, round, index, finalized, mrenclave, globalPubKey, signature), "Invalid finalization signature"
        );

        node.finalized = finalized;
        if (finalized) {
            votes[mrenclave][round][globalPubKey]++;
            if (votes[mrenclave][round][globalPubKey]  >= roundInfo[mrenclave][round].threshold ) {
                roundInfo[mrenclave][round].globalPubKey = globalPubKey;
            }
            emit DKGFinalized(
                msg.sender,
                round,
                index,
                finalized,
                mrenclave,
                globalPubKey,
                signature
            );
        }
    }

    function getGlobalPubKey(bytes calldata mrenclave, uint32 round) external view returns (bytes memory) {
        return roundInfo[mrenclave][round].globalPubKey;
    }

    function submitActiveValSet(
        uint32 round,
        bytes calldata mrenclave,
        address[] calldata valSet
    ) external {
        for (uint256 i = 0; i < valSet.length; i++) {
            // add if validator is not challenged (invalidated)
            // TODO: exclude validators that aren't participating in the DKG system
            if (dkgNodeInfos[mrenclave][round][uint32(i)].chalStatus != ChallengeStatus.Invalidated) {
                valSets[mrenclave][round][valSet[i]] = true;
            }
        }
    }

    function requestRemoteAttestationOnChain(
        uint32 targetIndex,
        uint32 round,
        bytes calldata mrenclave
    ) external {
        NodeInfo storage node = dkgNodeInfos[mrenclave][round][targetIndex];
        require(node.validator != address(0), "Node does not exist");
        require(node.chalStatus == ChallengeStatus.NotChallenged, "Node already challenged");

        bool isValid = _verifyRemoteAttestation(node.rawQuote, node.validator, round, node.dkgPubKey);        
        if (isValid) {
            node.chalStatus = ChallengeStatus.Resolved;
        } else {
            node.chalStatus = ChallengeStatus.Invalidated;
        }

        emit RemoteAttestationProcessedOnChain(
            targetIndex,
            node.validator,
            node.chalStatus,
            round,
            mrenclave
        );
    }

    function complainDeals(
        uint32 round,
        uint32 index,
        uint32[] memory complainIndexes,
        bytes calldata mrenclave
    ) external {
        NodeInfo storage complainant = dkgNodeInfos[mrenclave][round][index];
        require(complainant.validator == msg.sender, "Invalid complainant");

        for (uint256 i = 0; i < complainIndexes.length; i++) {
            dealComplaints[mrenclave][round][complainIndexes[i]][msg.sender] = true;
        }

        emit DealComplaintsSubmitted(index, complainIndexes, round, mrenclave);
    }

    //////////////////////////////////////////////////////////////
    //                      Getter Functions                    //
    //////////////////////////////////////////////////////////////

    function getNodeInfo(bytes calldata mrenclave, uint32 round, uint32 index) external view returns (NodeInfo memory) {
        return dkgNodeInfos[mrenclave][round][index];
    }

    function getNodeCount(bytes calldata mrenclave, uint32 round) external view returns (uint32) {
        return nodeCount[mrenclave][round];
    }

    function isActiveValidator(bytes calldata mrenclave, uint32 round, address validator) external view returns (bool) {
        return valSets[mrenclave][round][validator];
    }

    //////////////////////////////////////////////////////////////
    //                      Internal Functions                  //
    //////////////////////////////////////////////////////////////

    function _isActiveValSetSubmitted(bytes calldata /*mrenclave*/, uint32 /*round*/) internal pure returns (bool) {
        // TODO: Implementation
        return true;
    }

    function _verifyCommitmentSignature(
        bytes memory commPubKey,
        uint32 round,
        uint32 total,
        uint32 threshold,
        uint32 index,
        bytes calldata mrenclave,
        bytes calldata commitments,
        bytes calldata signature
    ) internal pure returns (bool) {
        bytes32 msgHash = keccak256(abi.encodePacked(round, total, threshold, index, mrenclave, commitments));
        address signer = ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(msgHash), signature);
        return signer == address(uint160(uint256(keccak256(commPubKey))));
    }

    function _verifyFinalizationSignature(
        bytes memory commPubKey,
        uint32 round,
        uint32 index,
        bool finalized,
        bytes calldata mrenclave,
        bytes calldata globalPubKey,
        bytes calldata signature
    ) internal pure returns (bool) {
        bytes32 msgHash = keccak256(abi.encodePacked(round, index, finalized, mrenclave, globalPubKey));
        address signer = ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(msgHash), signature);
        return signer == address(uint160(uint256(keccak256(commPubKey))));
    }

    function _verifyRemoteAttestation(
        bytes memory rawQuote,
        address validator,
        uint32 round,
        bytes memory dkgPubKey 
    ) internal pure returns (bool) {
        // TODO: Implementation
        return rawQuote.length > 0 && validator != address(0) && round > 0 && dkgPubKey.length > 0;
    }
}