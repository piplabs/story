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
        bytes32 codeCommitment,
        uint32 round,
        uint64 startBlockHeight,
        bytes32 startBlockHash,
        bytes dkgPubKey,
        bytes commPubKey,
        bytes rawQuote
    );

    event DKGFinalized(
        address indexed msgSender,
        uint32 round,
        bytes32 codeCommitment,
        bytes32 participantsRoot,
        bytes globalPubKey,
        bytes[] publicCoeffs,
        bytes signature
    );

    event UpgradeScheduled(uint32 activationHeight, bytes32 codeCommitment);

    event RemoteAttestationProcessedOnChain(
        address validator,
        ChallengeStatus chalStatus,
        uint32 round,
        bytes32 codeCommitment
    );

    // TODO: remove index and use validator address instead
    event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 codeCommitment);

    event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 codeCommitment);

    event InvalidDeal(uint32 index, uint32 round, bytes32 codeCommitment);

    // Emitted when a client requests TDH2 threshold decryption for a ciphertext/label pair.
    // @param requesterPubKey: secp256k1 uncompressed requester pubkey (65 bytes)
    event ThresholdDecryptRequested(
        address indexed requester,
        uint32 round,
        bytes32 codeCommitment,
        bytes requesterPubKey,
        bytes ciphertext,
        bytes label
    );

    // Emitted when a validator submits a TDH2 partial decryption.
    // @param pid: party ID, 1-based index from DKG registration (used in Kyber polynomial evaluation)
    event PartialDecryptionSubmitted(
        address indexed validator,
        uint32 round,
        bytes32 codeCommitment,
        uint32 pid,
        bytes encryptedPartial,
        bytes ephemeralPubKey,
        bytes pubShare,
        bytes label
    );

    // TDH2 partial decrypt submissions keyed by (codeCommitment, round, labelHash, pid).
    struct PartialDecryptSubmission {
        address validator;
        bytes partialDecryption;
        bytes pubShare;
        bytes label;
        bool exists;
    }
}
