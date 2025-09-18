// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

contract IDKG {
    enum ChallengeStatus {
        NotChallenged,
        Invalidated,
        Resolved
    }

    struct NodeInfo {
        uint32 index;
        address validator;
        bytes pubKey;
        bytes remoteReport;
        bytes commitments;
        ChallengeStatus chalStatus;
        bool finalized;
    }

    event DKGInitialized(
        address indexed msgSender,
        bytes mrenclave,
        uint32 round,
        uint32 index,
        bytes pubKey,
        bytes remoteReport
    );

    event DKGCommitmentsUpdated(
        address indexed msgSender,
        uint32 round,
        uint32 total,
        uint32 threshold,
        uint32 index,
        bytes commitments,
        bytes signature,
        bytes mrenclave
    );

    event DKGFinalized(
        address indexed msgSender,
        uint32 round,
        uint32 index,
        bool finalized,
        bytes mrenclave,
        bytes signature
    );

    event UpgradeScheduled(
        uint32 activationHeight,
        bytes mrenclave
    );

    event RegistrationChallenged(
        uint32 round,
        bytes mrenclave,
        address indexed challenger
    );

    event InvalidDKGInitialization(
        uint32 round,
        uint32 index,
        address validator,
        bytes mrenclave
    );

    event RemoteAttestationProcessedOnChain(
        uint32 index,
        address validator,
        ChallengeStatus chalStatus,
        uint32 round,
        bytes mrenclave
    );

    event DealComplaintsSubmitted(
        uint32 index,
        uint32[] complainIndexes,
        uint32 round,
        bytes mrenclave
    );

    event DealVerified(
        uint32 index,
        uint32 recipientIndex,
        uint32 round,
        bytes mrenclave
    );

    event InvalidDeal(
        uint32 index,
        uint32 round,
        bytes mrenclave
    );
}