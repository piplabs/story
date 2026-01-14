// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

contract IDKG {
    enum ChallengeStatus {
        NotChallenged,
        Invalidated,
        Resolved
    }

    enum NodeStatus {
        Unregistered,
        Registered,
        NetworkSetDone,
        Finalized
    }

    struct NodeInfo {
        bytes dkgPubKey;
        bytes commPubKey;
        bytes rawQuote;
        ChallengeStatus chalStatus;
        NodeStatus nodeStatus;
    }

    event DKGInitialized(
        address indexed msgSender,
        bytes32 mrenclave,
        uint32 round,
        bytes dkgPubKey,
        bytes commPubKey,
        bytes rawQuote
    );

    event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes[] publicCoeffs, bytes signature);

    event DKGNetworkSet(
        address indexed msgSender,
        uint32 round,
        uint32 total,
        uint32 threshold,
        bytes32 mrenclave,
        bytes signature
    );

    event UpgradeScheduled(uint32 activationHeight, bytes32 mrenclave);

    event RemoteAttestationProcessedOnChain(
        address validator,
        ChallengeStatus chalStatus,
        uint32 round,
        bytes32 mrenclave
    );

    // TODO: remove index and use validator address instead
    event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 mrenclave);

    event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 mrenclave);

    event InvalidDeal(uint32 index, uint32 round, bytes32 mrenclave);
}
