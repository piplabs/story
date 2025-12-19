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

    event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes signature);

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

    // Emitted when a client requests TDH2 threshold decryption for a ciphertext/label pair.
    // @param requesterPubKey: secp256k1 uncompressed requester pubkey (65 bytes)
    event ThresholdDecryptRequested(
        address indexed requester,
        uint32 round,
        bytes32 mrenclave,
        bytes requesterPubKey,
        bytes ciphertext,
        bytes label
    );

    // Emitted when a validator submits a TDH2 partial decryption.
    event PartialDecryptionSubmitted(
        address indexed validator,
        uint32 round,
        bytes32 mrenclave,
        bytes encryptedPartial,
        bytes ephemeralPubKey,
        bytes pubShare,
        bytes label
    );

    // TDH2 partial decrypt submissions keyed by (mrenclave, round, labelHash, validator).
    struct PartialDecryptSubmission {
        address validator;
        bytes partialDecryption;
        bytes pubShare;
        bytes label;
        bool exists;
    }
}
