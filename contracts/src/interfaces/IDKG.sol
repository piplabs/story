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

    struct RoundInfo {
        uint32 total;
        uint32 threshold;
        bytes globalPubKey;
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
        bytes mrenclave,
        uint32 round,
        bytes dkgPubKey,
        bytes commPubKey,
        bytes rawQuote
    );

    event DKGFinalized(address indexed msgSender, uint32 round, bytes mrenclave, bytes globalPubKey, bytes signature);

    event DKGNetworkSet(
        address indexed msgSender,
        uint32 round,
        uint32 total,
        uint32 threshold,
        bytes mrenclave,
        bytes signature
    );

    event UpgradeScheduled(uint32 activationHeight, bytes mrenclave);

    event RemoteAttestationProcessedOnChain(
        address validator,
        ChallengeStatus chalStatus,
        uint32 round,
        bytes mrenclave
    );

    // TODO: remove index and use validator address instead
    event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes mrenclave);

    event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes mrenclave);

    event InvalidDeal(uint32 index, uint32 round, bytes mrenclave);
}
