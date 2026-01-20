package types

import (
	"encoding/hex"
	"fmt"
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

	Mrenclave          []byte    `json:"mrenclave"`
	Round              uint32    `json:"round"`
	GlobalPubKey       []byte    `json:"global_pub_key"`
	DKGPubKey          []byte    `json:"dkg_pub_key"`
	CommPubKey         []byte    `json:"comm_pub_key"`
	RawQuote           []byte    `json:"raw_quote"`
	Phase              DKGPhase  `json:"phase"`
	StartTime          time.Time `json:"start_time"`
	LastUpdate         time.Time `json:"last_update"`
	Index              uint32    `json:"index"`
	SigSetupNetwork    []byte    `json:"sig_setup_network"`
	SigFinalizeNetwork []byte    `json:"sig_finalize_network"`
	PublicCoeffs       [][]byte  `json:"public_coeffs"`

	// Network information
	ActiveValidators []string `json:"active_validators"`
	Total            uint32   `json:"total"`
	Threshold        uint32   `json:"threshold"`

	// DKG state
	Registrations []DKGRegistration `json:"registrations,omitempty"`
	Commitments   []byte            `json:"commitments,omitempty"`
	Deals         map[uint32]Deal   `json:"deals,omitempty"` // deals by dealer index
	Complaints    []Complaint       `json:"complaints,omitempty"`
	IsFinalized   bool              `json:"is_finalized"`

	// Pending threshold decrypt requests (from contract events)
	DecryptRequests []DecryptRequest `json:"decrypt_requests,omitempty"`
}

// NewDKGSession creates a new DKG session from blockchain event data.
func NewDKGSession(mrenclave []byte, round uint32, activeValidators []string) *DKGSession {
	now := time.Now()

	return &DKGSession{
		Mrenclave:        mrenclave,
		Round:            round,
		GlobalPubKey:     make([]byte, 0),
		CommPubKey:       make([]byte, 0),
		Phase:            PhaseInitializing,
		StartTime:        now,
		LastUpdate:       now,
		ActiveValidators: activeValidators,
		Total:            0,
		Threshold:        0,
		Deals:            make(map[uint32]Deal),
		IsFinalized:      false,
		DecryptRequests:  make([]DecryptRequest, 0),
	}
}

// GetMrenclaveString returns the string representation of the mrenclave.
func (s *DKGSession) GetMrenclaveString() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return hex.EncodeToString(s.Mrenclave)
}

// GetSessionKey returns a unique key (mrenclave_round) for this DKG session.
func (s *DKGSession) GetSessionKey() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return fmt.Sprintf("%s_%d", s.GetMrenclaveString(), s.Round)
}

// UpdatePhase updates the session phase and timestamp.
func (s *DKGSession) UpdatePhase(phase DKGPhase) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Phase = phase
	s.LastUpdate = time.Now()
}

// AddDecryptRequest appends a threshold decrypt request to this session.
func (s *DKGSession) AddDecryptRequest(req DecryptRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.DecryptRequests = append(s.DecryptRequests, req)
	s.LastUpdate = time.Now()
}
