# Story CDR (Confidential Data Rails) Architecture

> For the overall system architecture, DKG lifecycle, reward distribution, and story-kernel details, see [DKG_DESIGN_DOC.md](DKG).

## Table of Contents

1. [Overview](#1-overview)
2. [Allocate](#2-allocate)
3. [Encryption & Write](#3-encryption--write)
4. [Decryption Request](#4-decryption-request)
5. [Partial Decryption](#5-partial-decryption)
6. [Partial Decryption Submission & CL Verification](#6-partial-decryption-submission--cl-verification)
7. [Threshold Combination](#7-threshold-combination)
8. [Decrypt Request Lifecycle](#8-decrypt-request-lifecycle)
9. [GetCDRPartials Query RPC](#9-getcdrpartials-query-rpc)

---

## 1. Overview

CDR enables users to store encrypted data on-chain and have it decrypted by the DKG committee through threshold decryption. The encryption uses TDH2 (Threshold Decryption with Hybrid encryption) backed by the DKG committee's shared key.

## 2. Allocate

Before writing data, a vault must be allocated. Allocation assigns a unique `uuid` and sets access conditions.

```
User/Application
  |
  CDR.allocate(updatable, writeConditionAddr, readConditionAddr,
               writeConditionData, readConditionData)
  |
  +-- Collects allocateFee (burned)
  +-- Assigns incrementing uuid (uint32)
  +-- Stores vault metadata (conditions, owner)
  |
  v
CDR Contract
  |
  Emits: VaultAllocated event {uuid, updatable, writeConditionAddr,
         readConditionAddr, writeConditionData, readConditionData}
```

## 3. Encryption & Write

```
Off-chain (User/Application)
  |
  1. Fetch GlobalPubKey from DKGNetwork (latest active round)
  2. TDH2.Encrypt(plaintext, GlobalPubKey, label) -> ciphertext
  3. CDR.write(uuid, ciphertext)
     |
     v
CDR Contract
  |
  Emits: VaultWritten event
  Stores: encrypted data in vault identified by uuid
```

The encryption is performed entirely off-chain using the DKG committee's `GlobalPubKey`. The `label` is derived from the vault's `uuid`.

## 4. Decryption Request

```
User calls CDR.read(uuid)
  |
  v
CDR Contract
  |
  Emits: VaultRead event {requester, uuid, ciphertext, requesterPubKey}
  |
  v
evmengine.ProcessCDRVaultRead()
  |
  - Get latest active DKG round
  - Convert uuid to label: label[28:32] = BigEndian(uuid)
  |
  v
dkgKeeper.ThresholdDecryptRequested(round, requesterPubKey, ciphertext, label, blockHeight)
  |
  +-- Consensus: setDecryptRequestHeight(requesterPubKey, label, blockHeight)
  |   (all nodes store this for timeout enforcement)
  |
  +-- If isDKGSvcEnabled AND round is Active AND validator in committee:
        session.AddDecryptRequest(DecryptRequest{round, ciphertext, label, requesterPubKey})
        stateManager.UpdateSession(session)
```

## 5. Partial Decryption

The `DecryptWorker` runs as a background loop (every 3 seconds) on each committee member:

```
StartDecryptWorker()
  |
  +-- decryptWorkerRunning.CompareAndSwap(false, true) [singleton]
  |
  +-- Ticker loop (3s):
        processDecryptQueue()
          |
          +-- For each session with pending DecryptRequests:
                handleDecryptRequest(session, req)
                  |
                  +-- client.PartialDecryptTDH2():
                  |     |
                  |     v
                  |   story-kernel.PartialDecryptTDH2():
                  |   1. Validate code commitment (MRENCLAVE)
                  |   2. Verify round matches latest active network
                  |   3. Load sealed DistKeyShare (from cache or SGX-sealed file)
                  |   4. Convert private share to cb-mpc TDH2PrivateShare
                  |   5. Build TDH2 public key from GlobalPubKey
                  |   6. mpc.TDH2PartialDecrypt(PID, privShare, pubKey, ct, label)
                  |   7. Compute pubShare = priShare.V * G
                  |   8. encryptPartialToRequester():
                  |      - ECDH: ephemeral key + requesterPubKey -> shared secret
                  |      - HKDF(sharedSecret, "dkg-tdh2-partial") -> AES-256 key
                  |      - AES-GCM encrypt partial decryption
                  |   9. Sign response: keccak256(CC || round || encryptedPartial ||
                  |                               ephPubKey || pubShare)
                  |   Return: encryptedPartial, ephemeralPubKey, pubShare, signature
                  |
                  +-- contractClient.SubmitEncryptedPartialDecryption(
                        round, pid, encryptedPartial, ephemeralPubKey,
                        pubShare, requesterPubKey, uuid, signature)
                        |
                        v
                      CDR.submitEncryptedPartialDecryption()
                        |
                        Emits: EncryptedPartialDecryptionSubmitted
```

## 6. Partial Decryption Submission & CL Verification

Each DKG committee member submits its partial decryption to the CDR contract, which emits an event. The CL then verifies and stores the submission:

```
DKG Member (off-chain)
  |
  contractClient.SubmitEncryptedPartialDecryption(
    round, pid, encryptedPartial, ephemeralPubKey,
    pubShare, requesterPubKey, uuid, signature)
  |
  v
CDR.submitEncryptedPartialDecryption()  (collects baseFee)
  |
  Emits: EncryptedPartialDecryptionSubmitted
  |
  v
evmengine.ProcessDKGPartialDecryptionSubmitted()
  |
  +-- Parse event, convert uuid -> label (32-byte, uuid in last 4 bytes)
  |
  v
dkgKeeper.PartialDecryptionSubmitted(validator, round, pid, encryptedPartial,
                                      ephemeralPubKey, pubShare, requesterPubKey,
                                      label, signature)
  |
  +-- Timeout enforcement:
  |     reqHeight = getDecryptRequestHeight(requesterPubKey, label)
  |     if currentHeight - reqHeight > PartialDecryptionTimeoutBlocks:
  |       delete registry entry, return (reject late submissions)
  |
  +-- Get validator's DKG registration
  +-- verifyPartialDecryptionSignature():
  |     Reconstruct message: round(4B) || encryptedPartial || ephPubKey || pubShare
  |     respHash = keccak256(message)
  |     Recover signer from ECDSA signature
  |     Verify recovered address == address(keccak256(commPubKey))
  |
  +-- Verify pubShare == stored pubKeyShare from registration
  +-- Duplicate submission check (by requesterPubKey + label + round + validator)
  +-- Store partial decryption in CL state:
  |     Key: {keccak256(requesterPubKey)}_{labelHex}_{round}_{validator}
  |     Value: JSON(validator, round, pid, encryptedPartial,
  |                  ephemeralPubKey, pubShare, label)
  |
  +-- On success: refund baseFee to submitter
```

> **Note**: Partial decryptions are stored entirely in CL state (Cosmos KV store),
> NOT in the CDR contract. The CDR contract only serves as the event emission
> entry point.

## 7. Threshold Combination

The requester collects partial decryptions from CL via the `GetCDRPartials` RPC and performs TDH2 combination off-chain:

```
Requester (off-chain):
  1. Call GetCDRPartials(uuid, requesterPubKeyHex) on CL
  |     |
  |     v
  |   DKG Module query handler:
  |   - Build prefix: keccak256(requesterPubKey)_labelHex(uuid)_
  |   - Range query all matching entries from CL state
  |   - Group by round, include threshold + threshold_met
  |   - Return DKGPartialDecryptionSubmissionsByRound[]
  |
  2. Check threshold_met == true (enough partials collected)
  3. For each partial:
     ECDH(requesterPrivKey, ephemeralPubKey) -> shared secret
     -> HKDF -> AES key -> AES-GCM decrypt -> raw partial decryption
  4. TDH2Combine(decryptedPartials, pubShares, ciphertext, label) -> plaintext
```

> **Note**: TDH2Combine is performed entirely off-chain by the requester. Each
> partial decryption is encrypted to the requester's public key, so only the
> requester can decrypt and combine them. The plaintext is never exposed on-chain.

## 8. Decrypt Request Lifecycle

```
                    ThresholdDecryptRequested
                    (DecryptRequestRegistry[key] = blockHeight)
                           |
         +-----------------+-----------------+
         |                                   |
   Partials submitted               Timeout (PartialDecryptionTimeoutBlocks)
   within timeout                            |
         |                           pruneTimedOutDecryptRequests()
         v                           (BeginBlocker, every N blocks)
   Enough partials                           |
   collected on-chain                        v
         |                            Registry entry deleted
         v
   TDH2Combine -> plaintext

Registry key format: hex(sha256(requesterPubKey))_hexLabel
Cleanup interval: DecryptRequestRegistryCleanupInterval blocks
```

## 9. GetCDRPartials Query RPC

The `GetCDRPartials` RPC allows requesters to query submitted partial decryption data for a given vault read request.

```
GetCDRPartials(uuid, requester_pub_key_hex)
  |
  +-- Build prefix: sha256(requesterPubKeyBytes)_labelHex(uuid)
  +-- Iterate DKGPartialDecrypt entries matching prefix
  +-- Group submissions by round
  |
  v
QueryGetCDRPartialsResponse:
  submissions_by_round[]:
    - round: uint32
    - threshold: uint32 (from DKGNetwork)
    - threshold_met: bool (len(submissions) >= threshold)
    - submissions[]:
        - validator_addr
        - encrypted_partial_decryption
        - ephemeral_pub_key
        - pub_share
        - requester_pub_key
        - signature
```

> The requester uses this RPC to retrieve all partial decryptions
> needed for TDH2Combine. When `threshold_met` is true, the requester
> has enough partials to reconstruct the plaintext.
