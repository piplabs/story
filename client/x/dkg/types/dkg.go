package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"
)

// DKGPhase represents the current phase of the DKG process.
type DKGPhase int32

const (
	PhaseUnknown      DKGPhase = 0
	PhaseInitializing DKGPhase = 1
	PhaseInitialized  DKGPhase = 2
	PhaseDealing      DKGPhase = 3
	PhaseFinalized    DKGPhase = 4
	PhaseCompleted    DKGPhase = 5
	PhaseFailed       DKGPhase = 6
)

func (p DKGPhase) String() string {
	switch p {
	case PhaseUnknown:
		return "Unknown"
	case PhaseInitializing:
		return "Initializing"
	case PhaseInitialized:
		return "Initialized"
	case PhaseDealing:
		return "Dealing"
	case PhaseFinalized:
		return "Finalized"
	case PhaseCompleted:
		return "Completed"
	case PhaseFailed:
		return "Failed"
	default:
		return fmt.Sprintf("Phase(%d)", int(p))
	}
}

// DKGSession represents a local DKG session managed by the service.
type DKGSession struct {
	mu sync.RWMutex

	CodeCommitment     []byte    `json:"code_commitment"`
	Round              uint32    `json:"round"`
	GlobalPubKey       []byte    `json:"global_pub_key"`
	DKGPubKey          []byte    `json:"dkg_pub_key"`
	CommPubKey         []byte    `json:"comm_pub_key"`
	EnclaveReport      []byte    `json:"enclave_report"`
	StartBlockHeight   int64     `json:"start_block_height"`
	StartBlockHash     []byte    `json:"start_block_hash"`
	Phase              DKGPhase  `json:"phase"`
	StartTime          time.Time `json:"start_time"`
	LastUpdate         time.Time `json:"last_update"`
	Index              uint32    `json:"index"`
	SigSetupNetwork    []byte    `json:"sig_setup_network"`
	SigFinalizeNetwork []byte    `json:"sig_finalize_network"`
	PublicCoeffs       [][]byte  `json:"public_coeffs"`
	PubKeyShare        []byte    `json:"pub_key_share"` // validator's own share of the DKG public key
	ParticipantsRoot   []byte    `json:"participants_root"`
	EnclaveType        [32]byte  `json:"enclave_type"`

	// Network information
	ActiveValidators []string `json:"active_validators"`
	Total            uint32   `json:"total"`
	Threshold        uint32   `json:"threshold"`

	// DKG state
	Registrations []DKGRegistration `json:"registrations,omitempty"`
	Commitments   []byte            `json:"commitments,omitempty"`
	Complaints    []Complaint       `json:"complaints,omitempty"`
	IsFinalized   bool              `json:"is_finalized"`
	IsResharing   bool              `json:"is_resharing"`
	IsUpgrade     bool              `json:"is_upgrade"`

	// OldCodeCommitment holds the previous active round's code commitment during upgrade resharing.
	// Dealers use this to route TEE calls to the old binary (which holds the sealed key shares).
	// Empty for non-upgrade rounds.
	OldCodeCommitment []byte `json:"old_code_commitment,omitempty"`

	// Pending threshold decrypt requests (from contract events).
	DecryptRequests []PendingDecryptRequest `json:"decrypt_requests,omitempty"`

	// RecoveryAttempts counts node-local key-material recovery attempts for a keyless
	// active-round session. It bounds the per-block finalize retries in ResumeDKGService.
	// Node-local JSON state only (never consensus/KVStore), so it is determinism-safe.
	RecoveryAttempts uint32 `json:"recovery_attempts,omitempty"`
}

// PendingDecryptRequest wraps a DecryptRequest with a retry counter for the decrypt queue.
type PendingDecryptRequest struct {
	DecryptRequest
	RetryCount int `json:"retry_count"`
}

// NewDKGSession creates a new DKG session from blockchain event data.
func NewDKGSession(round uint32, activeValidators []string, isResharing bool, enclaveType [32]byte) *DKGSession {
	now := time.Now()

	return &DKGSession{
		Round:            round,
		GlobalPubKey:     make([]byte, 0),
		CommPubKey:       make([]byte, 0),
		Phase:            PhaseInitializing,
		StartTime:        now,
		LastUpdate:       now,
		ActiveValidators: activeValidators,
		Total:            0,
		Threshold:        0,
		IsFinalized:      false,
		IsResharing:      isResharing,
		EnclaveType:      enclaveType,

		DecryptRequests: make([]PendingDecryptRequest, 0),
	}
}

// GetCodeCommitmentString returns the string representation of the code commitment.
func (s *DKGSession) GetCodeCommitmentString() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return hex.EncodeToString(s.CodeCommitment)
}

// GetSessionKey returns a unique key (round) for this DKG session.
func (s *DKGSession) GetSessionKey() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return strconv.FormatUint(uint64(s.Round), 10)
}

// UpdatePhase updates the session phase and timestamp.
func (s *DKGSession) UpdatePhase(phase DKGPhase) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Phase = phase
	s.LastUpdate = time.Now()
}

// GetRecoveryAttempts returns the node-local active-round recovery attempt counter.
func (s *DKGSession) GetRecoveryAttempts() uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.RecoveryAttempts
}

// IncrementRecoveryAttempts increments the recovery attempt counter and returns the new value.
func (s *DKGSession) IncrementRecoveryAttempts() uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.RecoveryAttempts++
	s.LastUpdate = time.Now()

	return s.RecoveryAttempts
}

// SetRecoveryAttempts overrides the recovery attempt counter. Used to escalate past the
// retry cap (deterministic divergent-key failures skip straight to the exhausted sentinel).
func (s *DKGSession) SetRecoveryAttempts(n uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.RecoveryAttempts = n
	s.LastUpdate = time.Now()
}

// HasKeyMaterial reports whether the session holds both the global public key and this
// validator's key share. Read under RLock so it never observes a torn slice header while
// the recovery/finalization goroutine writes the fields via SetKeyMaterial.
func (s *DKGSession) HasKeyMaterial() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.GlobalPubKey) > 0 && len(s.PubKeyShare) > 0
}

// GetGlobalPubKey returns a copy of the session global public key under RLock.
func (s *DKGSession) GetGlobalPubKey() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.GlobalPubKey)
}

// GetSigFinalizeNetwork returns a copy of the finalize-network signature under RLock.
func (s *DKGSession) GetSigFinalizeNetwork() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.SigFinalizeNetwork)
}

// SetKeyMaterial atomically stores the key material returned by the kernel finalize call
// (participants root, global public key, finalize signature, this validator's key share, and
// the public coefficients) under a single Lock, so a concurrent reader such as HasKeyMaterial
// on the ABCI thread never observes a partially-written set of slice headers.
func (s *DKGSession) SetKeyMaterial(participantsRoot, globalPubKey, sigFinalizeNetwork, pubKeyShare []byte, publicCoeffs [][]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ParticipantsRoot = participantsRoot
	s.GlobalPubKey = globalPubKey
	s.SigFinalizeNetwork = sigFinalizeNetwork
	s.PubKeyShare = pubKeyShare
	s.PublicCoeffs = publicCoeffs
	s.LastUpdate = time.Now()
}

// GetLastUpdate returns the timestamp of the last session mutation.
func (s *DKGSession) GetLastUpdate() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.LastUpdate
}

// AddDecryptRequest appends a threshold decrypt request to this session.
func (s *DKGSession) AddDecryptRequest(req PendingDecryptRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.DecryptRequests = append(s.DecryptRequests, req)
	s.LastUpdate = time.Now()
}

// GetDecryptRequests returns a copy of the pending decrypt requests.
func (s *DKGSession) GetDecryptRequests() []PendingDecryptRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := make([]PendingDecryptRequest, len(s.DecryptRequests))
	copy(cp, s.DecryptRequests)

	return cp
}

// SetDecryptRequests replaces the decrypt requests slice (used after processing to retain only failed requests).
func (s *DKGSession) SetDecryptRequests(remaining []PendingDecryptRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.DecryptRequests = remaining
	s.LastUpdate = time.Now()
}

// DrainDecryptRequests atomically returns all pending decrypt requests and clears the queue.
// This prevents the TOCTOU race where GetDecryptRequests + SetDecryptRequests could
// overwrite requests added between the two calls by the ABCI thread.
func (s *DKGSession) DrainDecryptRequests() []PendingDecryptRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	reqs := s.DecryptRequests
	s.DecryptRequests = make([]PendingDecryptRequest, 0)
	s.LastUpdate = time.Now()

	return reqs
}

// GetPhase returns the current session phase under RLock.
func (s *DKGSession) GetPhase() DKGPhase {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Phase
}

// GetIndex returns this validator's 1-based on-chain registration index under RLock.
func (s *DKGSession) GetIndex() uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Index
}

// SetIndex stores this validator's 1-based on-chain registration index under Lock.
func (s *DKGSession) SetIndex(index uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Index = index
	s.LastUpdate = time.Now()
}

// GetIsFinalized reports whether the session key material is ready for
// threshold encryption/decryption. Read under RLock.
func (s *DKGSession) GetIsFinalized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.IsFinalized
}

// SetFinalized marks the session ready for threshold encryption/decryption under Lock.
func (s *DKGSession) SetFinalized() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.IsFinalized = true
	s.LastUpdate = time.Now()
}

// GetParticipantsRoot returns a copy of the participants root under RLock.
func (s *DKGSession) GetParticipantsRoot() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.ParticipantsRoot)
}

// GetPubKeyShare returns a copy of this validator's key share under RLock.
func (s *DKGSession) GetPubKeyShare() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.PubKeyShare)
}

// GetPublicCoeffs returns a deep copy of the public coefficients under RLock.
func (s *DKGSession) GetPublicCoeffs() [][]byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return cloneByteSlices(s.PublicCoeffs)
}

// GetRound returns the session round. Immutable after construction, but read under RLock
// for uniform mutex-based access.
func (s *DKGSession) GetRound() uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Round
}

// GetIsResharing reports whether this is a resharing round. Immutable after construction,
// read under RLock for uniformity.
func (s *DKGSession) GetIsResharing() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.IsResharing
}

// GetEnclaveType returns the enclave type. The array is returned by value (copy).
func (s *DKGSession) GetEnclaveType() [32]byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.EnclaveType
}

// GetCodeCommitment returns a copy of the code commitment under RLock.
func (s *DKGSession) GetCodeCommitment() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.CodeCommitment)
}

// GetOldCodeCommitment returns a copy of the previous round's code commitment under RLock.
func (s *DKGSession) GetOldCodeCommitment() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.OldCodeCommitment)
}

// GetDKGPubKey returns a copy of the DKG public key under RLock.
func (s *DKGSession) GetDKGPubKey() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.DKGPubKey)
}

// GetCommPubKey returns a copy of the communication public key under RLock.
func (s *DKGSession) GetCommPubKey() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.CommPubKey)
}

// GetEnclaveReport returns a copy of the enclave report under RLock.
func (s *DKGSession) GetEnclaveReport() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.EnclaveReport)
}

// GetStartBlockHeight returns the setup start block height under RLock.
func (s *DKGSession) GetStartBlockHeight() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.StartBlockHeight
}

// GetStartBlockHash returns a copy of the setup start block hash under RLock.
func (s *DKGSession) GetStartBlockHash() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.Clone(s.StartBlockHash)
}

// HasSetupData reports whether the key-generation setup fields have been populated.
func (s *DKGSession) HasSetupData() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.DKGPubKey) > 0 && len(s.CommPubKey) > 0 && len(s.EnclaveReport) > 0
}

// SetIsUpgrade records whether this is an upgrade round.
func (s *DKGSession) SetIsUpgrade(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.IsUpgrade = v
	s.LastUpdate = time.Now()
}

// SetOldCodeCommitment stores the previous active round's code commitment (upgrade routing).
func (s *DKGSession) SetOldCodeCommitment(oldCC []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.OldCodeCommitment = oldCC
	s.LastUpdate = time.Now()
}

// SetCodeCommitmentIfEmpty sets the code commitment only when it is currently unset,
// as a single atomic check-and-set so a concurrent reader cannot observe a TOCTOU gap.
func (s *DKGSession) SetCodeCommitmentIfEmpty(cc []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.CodeCommitment) == 0 {
		s.CodeCommitment = cc
		s.LastUpdate = time.Now()
	}
}

// SetSetupResult atomically stores the key-generation setup fields returned by the kernel
// GenerateAndSealKey call under a single Lock, so a concurrent reader (e.g. HasSetupData on
// the ABCI thread) never observes a partially-written set of slice headers.
func (s *DKGSession) SetSetupResult(codeCommitment, dkgPubKey, commPubKey, enclaveReport, startBlockHash []byte, startBlockHeight int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CodeCommitment = codeCommitment
	s.DKGPubKey = dkgPubKey
	s.CommPubKey = commPubKey
	s.EnclaveReport = enclaveReport
	s.StartBlockHash = startBlockHash
	s.StartBlockHeight = startBlockHeight
	s.LastUpdate = time.Now()
}

// Snapshot returns a deep copy of the session under RLock, so it can be marshaled
// (e.g. by saveSession) without racing a concurrent mutation. Built field-by-field
// to avoid copying the live mutex (copylocks); do not call other s.mu getters here,
// as the lock is not reentrant.
func (s *DKGSession) Snapshot() *DKGSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &DKGSession{
		CodeCommitment:     bytes.Clone(s.CodeCommitment),
		Round:              s.Round,
		GlobalPubKey:       bytes.Clone(s.GlobalPubKey),
		DKGPubKey:          bytes.Clone(s.DKGPubKey),
		CommPubKey:         bytes.Clone(s.CommPubKey),
		EnclaveReport:      bytes.Clone(s.EnclaveReport),
		StartBlockHeight:   s.StartBlockHeight,
		StartBlockHash:     bytes.Clone(s.StartBlockHash),
		Phase:              s.Phase,
		StartTime:          s.StartTime,
		LastUpdate:         s.LastUpdate,
		Index:              s.Index,
		SigSetupNetwork:    bytes.Clone(s.SigSetupNetwork),
		SigFinalizeNetwork: bytes.Clone(s.SigFinalizeNetwork),
		PublicCoeffs:       cloneByteSlices(s.PublicCoeffs),
		PubKeyShare:        bytes.Clone(s.PubKeyShare),
		ParticipantsRoot:   bytes.Clone(s.ParticipantsRoot),
		EnclaveType:        s.EnclaveType,
		ActiveValidators:   slices.Clone(s.ActiveValidators),
		Total:              s.Total,
		Threshold:          s.Threshold,
		Registrations:      slices.Clone(s.Registrations),
		Commitments:        bytes.Clone(s.Commitments),
		Complaints:         slices.Clone(s.Complaints),
		IsFinalized:        s.IsFinalized,
		IsResharing:        s.IsResharing,
		IsUpgrade:          s.IsUpgrade,
		OldCodeCommitment:  bytes.Clone(s.OldCodeCommitment),
		DecryptRequests:    slices.Clone(s.DecryptRequests),
		RecoveryAttempts:   s.RecoveryAttempts,
	}
}

// cloneByteSlices deep-copies a slice of byte slices. Returns nil for nil input
// to preserve JSON round-trip equivalence.
func cloneByteSlices(in [][]byte) [][]byte {
	if in == nil {
		return nil
	}

	out := make([][]byte, len(in))
	for i, b := range in {
		out[i] = bytes.Clone(b)
	}

	return out
}
