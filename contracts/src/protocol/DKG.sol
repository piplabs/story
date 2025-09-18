// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { IDKG } from "../interfaces/IDKG.sol";

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
    
    constructor() {}

    function initializeDKG(
        uint32 round,
        bytes calldata mrenclave,
        bytes calldata pubKey,
        bytes calldata remoteReport
    ) external {
        require(address(uint160(uint256(keccak256(pubKey)))) == msg.sender, "Invalid pubKey for sender");
        require(
            valSets[mrenclave][round][msg.sender] || _isActiveValSetSubmitted(mrenclave, round),
            "Validator not in active set"
        );
        require(_isRemoteReportValid(remoteReport), "Invalid remote report");
        
        uint32 index = nodeCount[mrenclave][round];
        nodeCount[mrenclave][round]++;

        dkgNodeInfos[mrenclave][round][index] = NodeInfo({
            index: index,
            validator: msg.sender,
            pubKey: pubKey,
            remoteReport: remoteReport,
            commitments: "",
            chalStatus: ChallengeStatus.NotChallenged,
            finalized: false
        });

        emit DKGInitialized(
            msg.sender,
            mrenclave,
            round,
            index,
            pubKey,
            remoteReport
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
        require(_verifyCommitmentSignature(msg.sender, commitments, signature), "Invalid signature");

        node.commitments = commitments;

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
        bytes calldata signature
    ) external {
        NodeInfo storage node = dkgNodeInfos[mrenclave][round][index];
        require(node.validator == msg.sender, "Invalid sender");
        require(node.chalStatus != ChallengeStatus.Invalidated, "Node was invalidated");
        require(
            _verifyFinalizationSignature(msg.sender, round, index, finalized, mrenclave, signature), "Invalid signature"
        );

        node.finalized = finalized;

        emit DKGFinalized(
            msg.sender,
            round,
            index,
            finalized,
            mrenclave,
            signature
        );
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

        bool isValid = _verifyRemoteAttestation(node.remoteReport, node.validator, round, node.pubKey);        
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
        address sender,
        bytes calldata commitments,
        bytes calldata signature
    ) internal pure returns (bool) {
        bytes32 msgHash = keccak256(commitments);
        return _recoverSigner(msgHash, signature) == sender;
    }

    function _verifyFinalizationSignature(
    address sender,
    uint32 round,
    uint32 index,
    bool finalized,
    bytes calldata mrenclave,
    bytes calldata signature
    ) internal pure returns (bool) {
        bytes32 msgHash = keccak256(abi.encodePacked(round, index, finalized, mrenclave));
        return _recoverSigner(msgHash, signature) == sender;
    }

    function _recoverSigner(bytes32 msgHash, bytes memory signature) internal pure returns (address) {
        bytes32 ethSignedMessageHash =
            keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", msgHash));
        (bytes32 r, bytes32 s, uint8 v) = _splitSignature(signature);
        return ecrecover(ethSignedMessageHash, v, r, s);
    }

    function _splitSignature(bytes memory sig)
        internal
        pure
        returns (bytes32 r, bytes32 s, uint8 v)
    {
        require(sig.length == 65, "invalid signature length");
        assembly {
            r := mload(add(sig, 32))
            s := mload(add(sig, 64))
            v := byte(0, mload(add(sig, 96)))
        }
    }

    function _isRemoteReportValid(bytes calldata remoteReport) internal pure returns (bool) {
        // TODO: Implementation
        return remoteReport.length > 0;
    }

    function _verifyRemoteAttestation(
        bytes memory remoteReport,
        address validator,
        uint32 round,
        bytes memory pubKey
    ) internal pure returns (bool) {
        // TODO: Implementation
        return remoteReport.length > 0 && validator != address(0) && round > 0 && pubKey.length > 0;
    }
}