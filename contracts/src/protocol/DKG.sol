// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { IDKG } from "../interfaces/IDKG.sol";
import { ECDSA } from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import { MessageHashUtils } from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

/**
 * @title DKG - Distributed Key Generation Contract
 * @dev Core contract for managing DKG-related state and operations
 */
contract DKG is IDKG {
    uint32 constant UINT32_MAX = type(uint32).max;

    bytes32 public curMrenclave;

    mapping(bytes32 mrenclave => mapping(uint32 round => mapping(address validator => NodeInfo))) public dkgNodeInfos;
    mapping(bytes32 mrenclave => mapping(uint32 round => mapping(address validator => bool))) public valSets;
    mapping(bytes32 mrenclave => mapping(uint32 round => mapping(uint32 index => mapping(address complainant => bool))))
        public dealComplaints;

    mapping(bytes32 mrenclave => mapping(uint32 round => mapping(bytes globalPubKeyCandidates => uint32 votes)))
        public votes;
    mapping(bytes32 mrenclave => mapping(uint32 round => RoundInfo roundInfo)) public roundInfo;

    constructor(bytes32 mrenclave) {
        curMrenclave = mrenclave;
    }

    modifier onlyValidMrenclave(bytes32 mrenclave) {
        require(keccak256(abi.encodePacked(mrenclave)) == keccak256(abi.encodePacked(curMrenclave)), "Invalid mrenclave");
        _;
    }

    function initializeDKG(
        uint32 round,
        bytes32 mrenclave,
        bytes calldata dkgPubKey,
        bytes calldata commPubKey,
        bytes calldata rawQuote
    ) external onlyValidMrenclave(mrenclave) {
        // TODO: round check if it is current active round or not

        require(
            valSets[mrenclave][round][msg.sender] || !_isActiveValSetSubmitted(mrenclave, round),
            "Validator not in active set"
        );

        require(
            _verifyRemoteAttestation(rawQuote, msg.sender, round, dkgPubKey, commPubKey),
            "Invalid remote attestation"
        );

        dkgNodeInfos[mrenclave][round][msg.sender] = NodeInfo({
            dkgPubKey: dkgPubKey,
            commPubKey: commPubKey,
            rawQuote: rawQuote,
            chalStatus: ChallengeStatus.NotChallenged,
            nodeStatus: NodeStatus.Registered
        });

        emit DKGInitialized(msg.sender, mrenclave, round, dkgPubKey, commPubKey, rawQuote);
    }

    function finalizeDKG(
        uint32 round,
        bytes32 mrenclave,
        bytes calldata globalPubKey,
        bytes calldata signature
    ) external onlyValidMrenclave(mrenclave) {
        NodeInfo storage node = dkgNodeInfos[mrenclave][round][msg.sender];
        require(node.chalStatus != ChallengeStatus.Invalidated, "Node was invalidated");
        require(
            _verifyFinalizationSignature(node.commPubKey, round, mrenclave, globalPubKey, signature),
            "Invalid finalization signature"
        );

        node.nodeStatus = NodeStatus.Finalized;

        votes[mrenclave][round][globalPubKey]++;
        if (votes[mrenclave][round][globalPubKey] >= roundInfo[mrenclave][round].threshold) {
            roundInfo[mrenclave][round].globalPubKey = globalPubKey;
        }

        emit DKGFinalized(msg.sender, round, mrenclave, globalPubKey, signature);
    }

    function getGlobalPubKey(bytes32 mrenclave, uint32 round) external view returns (bytes memory) {
        return roundInfo[mrenclave][round].globalPubKey;
    }

    function submitActiveValSet(
        uint32 round,
        bytes32 mrenclave,
        address[] calldata valSet
    ) external onlyValidMrenclave(mrenclave) {
        for (uint256 i = 0; i < valSet.length; i++) {
            // add if validator is not challenged (invalidated)
            // TODO: exclude validators that aren't participating in the DKG system
            if (dkgNodeInfos[mrenclave][round][valSet[i]].chalStatus != ChallengeStatus.Invalidated) {
                valSets[mrenclave][round][valSet[i]] = true;
            }
        }
    }

    function requestRemoteAttestationOnChain(
        address targetValidatorAddr,
        uint32 round,
        bytes32 mrenclave
    ) external onlyValidMrenclave(mrenclave) {
        NodeInfo storage node = dkgNodeInfos[mrenclave][round][targetValidatorAddr];
        require(node.dkgPubKey.length != 0, "Node does not exist");
        require(node.chalStatus == ChallengeStatus.NotChallenged, "Node already challenged");

        bool isValid = _verifyRemoteAttestation(
            node.rawQuote,
            targetValidatorAddr,
            round,
            node.dkgPubKey,
            node.commPubKey
        );
        if (isValid) {
            node.chalStatus = ChallengeStatus.Resolved;
        } else {
            node.chalStatus = ChallengeStatus.Invalidated;
        }

        emit RemoteAttestationProcessedOnChain(targetValidatorAddr, node.chalStatus, round, mrenclave);
    }

    function complainDeals(
        uint32 round,
        uint32 index,
        uint32[] memory complainIndexes,
        bytes32 mrenclave
    ) external onlyValidMrenclave(mrenclave) {
        NodeInfo storage complainant = dkgNodeInfos[mrenclave][round][msg.sender];

        for (uint256 i = 0; i < complainIndexes.length; i++) {
            dealComplaints[mrenclave][round][complainIndexes[i]][msg.sender] = true;
        }

        emit DealComplaintsSubmitted(index, complainIndexes, round, mrenclave);
    }

    function setNetwork(
        uint32 round,
        uint32 total,
        uint32 threshold,
        bytes32 mrenclave,
        bytes calldata signature
    ) external onlyValidMrenclave(mrenclave) {
        NodeInfo storage node = dkgNodeInfos[mrenclave][round][msg.sender];

        require(
            _verifySetNetworkSignature(node.commPubKey, round, total, threshold, mrenclave, signature),
            "Invalid set network signature"
        );

        node.nodeStatus = NodeStatus.NetworkSetDone;
        roundInfo[mrenclave][round].total = total;
        roundInfo[mrenclave][round].threshold = threshold;

        emit DKGNetworkSet(msg.sender, round, total, threshold, mrenclave, signature);
    }

    //////////////////////////////////////////////////////////////
    //                      Getter Functions                    //
    //////////////////////////////////////////////////////////////

    function getNodeInfo(
        bytes32 mrenclave,
        uint32 round,
        address validator
    ) external view returns (NodeInfo memory) {
        return dkgNodeInfos[mrenclave][round][validator];
    }

    function isActiveValidator(bytes32 mrenclave, uint32 round, address validator) external view returns (bool) {
        return valSets[mrenclave][round][validator];
    }

    //////////////////////////////////////////////////////////////
    //                      Internal Functions                  //
    //////////////////////////////////////////////////////////////

    function _isActiveValSetSubmitted(bytes32 /*mrenclave*/, uint32 /*round*/) internal pure returns (bool) {
        // TODO: Implementation
        return false;
    }

    function _verifyFinalizationSignature(
        bytes memory commPubKey,
        uint32 round,
        bytes32 mrenclave,
        bytes calldata globalPubKey,
        bytes calldata signature
    ) internal pure returns (bool) {
        bytes32 msgHash = keccak256(abi.encodePacked( mrenclave, round, globalPubKey));
        address signer = ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(msgHash), signature);
        return signer == address(uint160(uint256(keccak256(commPubKey))));
    }

    function _verifySetNetworkSignature(
        bytes memory commPubKey,
        uint32 round,
        uint32 total,
        uint32 threshold,
        bytes32 mrenclave,
        bytes calldata signature
    ) internal pure returns (bool) {
        bytes32 msgHash = keccak256(abi.encodePacked( mrenclave, round, total, threshold));
        address signer = ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(msgHash), signature);
        return signer == address(uint160(uint256(keccak256(commPubKey))));
    }

    function _verifyRemoteAttestation(
        bytes memory rawQuote,
        address validator,
        uint32 round,
        bytes memory dkgPubKey,
        bytes memory commPubKey
    ) internal pure returns (bool) {
        // TODO: on chain DCAP verification
        require(rawQuote.length > 64, "Invalid raw quote, quote too short");
        require(validator != address(0), "Invalid validator address");
        require(round > 0, "Invalid round of DKG");
        require(dkgPubKey.length > 0, "Invalid DKG public key");
        require(commPubKey.length > 0, "Invalid communication public key");

        bytes32 expectedReportData = _extractReportData(rawQuote);
        return _validateReportData(validator, round, dkgPubKey, commPubKey, expectedReportData);
    }

    function _extractReportData(
        bytes memory rawQuote
    ) internal pure returns (bytes32) {
        // According to Intel’s SGX quote structure:
        // - The SGX quote header is 48 bytes in size
        // - The enclave report body is 384 bytes long
        // - The last 64 bytes of the enclave report body are reserved for report_data
        // Therefore, the starting offset for report_data is: 48 (quote header) + 320 = 368
        // https://github.com/intel/SGX-TDX-DCAP-QuoteVerificationLibrary/blob/16b7291a7a86e486fdfcf1dfb4be885c0cc00b4e/Src/AttestationLibrary/src/QuoteVerification/QuoteConstants.h
        uint256 start = 368;
        bytes32 first32;
        assembly {
            first32 := mload(add(add(rawQuote, 32), start))
        }
        return first32;
    }

    function _validateReportData(
        address validator,
        uint32 round,
        bytes memory dkgPubKey,
        bytes memory commPubKey,
        bytes32 expectedReportData
    ) internal pure returns (bool) {
        bytes32 reportData = keccak256(abi.encodePacked(validator, round, dkgPubKey, commPubKey));
        if (reportData != expectedReportData) {
            return false;
        }
        return true;
    }
}
