package types

import (
	"encoding/hex"
	"fmt"
	"time"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// DKGPhase represents the current phase of the DKG process.
type DKGPhase int32

// re-export types from x/dkg/types.
type Deal = dkgtypes.Deal
type EncryptedDeal = dkgtypes.EncryptedDeal
type DealWithCommitments = dkgtypes.DealWithCommitments
type Complaint = dkgtypes.Complaint
type Commitments = []byte

const (
	PhaseUnknown      DKGPhase = 0
	PhaseInitializing DKGPhase = 1
	PhaseRegistering  DKGPhase = 2
	PhaseChallenging  DKGPhase = 3
	PhaseDealing      DKGPhase = 4
	PhaseFinalizing   DKGPhase = 5
	PhaseCompleted    DKGPhase = 6
	PhaseFailed       DKGPhase = 7
)

func (p DKGPhase) String() string {
	switch p {
	case PhaseUnknown:
		return "Unknown"
	case PhaseInitializing:
		return "Initializing"
	case PhaseRegistering:
		return "Registering"
	case PhaseChallenging:
		return "Challenging"
	case PhaseDealing:
		return "Dealing"
	case PhaseFinalizing:
		return "Finalizing"
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
	Mrenclave     []byte    `json:"mrenclave"`
	Round         uint32    `json:"round"`
	Phase         DKGPhase  `json:"phase"`
	StartTime     time.Time `json:"start_time"`
	LastUpdate    time.Time `json:"last_update"`
	ValidatorAddr string    `json:"validator_address"`
	Index         uint32    `json:"index"`

	// Network information
	ActiveValidators []string `json:"active_validators"`
	Total            uint32   `json:"total"`
	Threshold        uint32   `json:"threshold"`

	// DKG state
	Registrations []dkgtypes.DKGRegistration `json:"registrations,omitempty"`
	Commitments   []byte                     `json:"commitments,omitempty"`
	Deals         map[uint32]dkgtypes.Deal   `json:"deals,omitempty"` // deals by dealer index
	Complaints    []dkgtypes.Complaint       `json:"complaints,omitempty"`
	IsFinalized   bool                       `json:"is_finalized"`
}

// NewDKGSession creates a new DKG session from blockchain event data.
func NewDKGSession(mrenclave []byte, round uint32, activeValidators []string) *DKGSession {
	now := time.Now()

	return &DKGSession{
		Mrenclave:        mrenclave,
		Round:            round,
		Phase:            PhaseInitializing,
		StartTime:        now,
		LastUpdate:       now,
		ActiveValidators: activeValidators,
		Total:            0,
		Threshold:        0,
		Deals:            make(map[uint32]dkgtypes.Deal),
		IsFinalized:      false,
	}
}

// GetMrenclaveString returns the string representation of the mrenclave.
func (s *DKGSession) GetMrenclaveString() string {
	return hex.EncodeToString(s.Mrenclave)
}

// GetSessionKey returns a unique key (mrenclave_round) for this DKG session.
func (s *DKGSession) GetSessionKey() string {
	return fmt.Sprintf("%s_%d", s.GetMrenclaveString(), s.Round)
}

// UpdatePhase updates the session phase and timestamp.
func (s *DKGSession) UpdatePhase(phase DKGPhase) {
	s.Phase = phase
	s.LastUpdate = time.Now()
}

// DKGEventData represents data from a DKG-related blockchain event emitted in Cosmos CL (not EL predeploy contract).
type DKGEventData struct {
	EventType        string            `json:"event_type"`
	Mrenclave        string            `json:"mrenclave"`
	Round            uint32            `json:"round"`
	BlockHeight      int64             `json:"block_height"`
	ActiveValidators []string          `json:"active_validators,omitempty"`
	Total            uint32            `json:"total,omitempty"`
	Threshold        uint32            `json:"threshold,omitempty"`
	ValidatorAddr    string            `json:"validator_address,omitempty"`
	Index            uint32            `json:"index,omitempty"`
	Attributes       map[string]string `json:"attributes,omitempty"`
}

// ParseMrenclave converts the hex-encoded mrenclave string to bytes.
func (e *DKGEventData) ParseMrenclave() ([]byte, error) {
	b := []byte(e.Mrenclave)
	if len(b) != 32 { // expect 256-bit digest (32 bytes)
		return nil, errors.New("mrenclave is not a 256-bit digest", "expected_size", 32, "actual_size", len(b))
	}

	return b, nil
}
