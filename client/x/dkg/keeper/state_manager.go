package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// StateManager manages the local state of DKG sessions outside of consensus.
//

type StateManager struct {
	dataDir  string
	mu       sync.RWMutex
	sessions map[string]*types.DKGSession // keyed by session key (round)
}

// NewStateManager creates a new state manager.
func NewStateManager(dataDir string) (*StateManager, error) {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, errors.Wrap(err, "failed to create data directory")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to check data directory")
	}

	sm := &StateManager{
		dataDir:  dataDir,
		sessions: make(map[string]*types.DKGSession),
	}

	// Load existing sessions from disk
	if err := sm.loadSessions(); err != nil {
		return nil, errors.Wrap(err, "failed to load existing sessions")
	}

	return sm, nil
}

// CreateSession creates a new DKG session.
func (sm *StateManager) CreateSession(ctx context.Context, session *types.DKGSession) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionKey := session.GetSessionKey()

	if _, exists := sm.sessions[sessionKey]; exists {
		log.Info(ctx, "Session already exists with the code commitment and round, skip creating a new session", "code_commitment", session.GetCodeCommitmentString(), "round", session.Round)

		return nil
	} else {
		sm.sessions[sessionKey] = session

		if err := sm.saveSession(session); err != nil {
			return errors.Wrap(err, "failed to save session to disk")
		}

		log.Info(ctx, "Created DKG session",
			"code_commitment", session.GetCodeCommitmentString(),
			"round", session.Round,
			"phase", session.Phase.String(),
		)

		return nil
	}
}

// GetSession retrieves a DKG session by round.
func (sm *StateManager) GetSession(round uint32) (*types.DKGSession, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessionKey := strconv.FormatUint(uint64(round), 10)

	session, exists := sm.sessions[sessionKey]
	if !exists {
		return nil, errors.New("session not found", "session_key", sessionKey)
	}

	return session, nil
}

// UpdateSession updates an existing DKG session.
func (sm *StateManager) UpdateSession(ctx context.Context, session *types.DKGSession) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionKey := session.GetSessionKey()

	if _, exists := sm.sessions[sessionKey]; !exists {
		return errors.New("session not found", "session_key", sessionKey)
	}

	sm.sessions[sessionKey] = session

	if err := sm.saveSession(session); err != nil {
		return errors.Wrap(err, "failed to save updated session to disk")
	}

	log.Debug(ctx, "Updated DKG session",
		"code_commitment", session.GetCodeCommitmentString(),
		"round", session.Round,
		"phase", session.Phase.String(),
	)

	return nil
}

func (sm *StateManager) MarkFailed(ctx context.Context, session *types.DKGSession) {
	session.UpdatePhase(types.PhaseFailed)

	if err := sm.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to mark session as failed", err)
	}
}

// ListSessions returns all active DKG sessions.
func (sm *StateManager) ListSessions() []*types.DKGSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*types.DKGSession, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}

	return sessions
}

// DeleteSession removes a DKG session.
func (sm *StateManager) DeleteSession(ctx context.Context, round uint32) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionKey := strconv.FormatUint(uint64(round), 10)

	session, exists := sm.sessions[sessionKey]
	if !exists {
		return errors.New("session not found", "session_key", sessionKey)
	}

	delete(sm.sessions, sessionKey)

	// Remove from disk
	filename := sm.getSessionFilename(sessionKey)
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to delete session file")
	}

	log.Info(ctx, "Deleted DKG session",
		"code_commitment", session.GetCodeCommitmentString(),
		"round", session.Round,
	)

	return nil
}

// GetActiveSession returns the currently active DKG session for a validator.
func (sm *StateManager) GetActiveSession(validatorAddr string) *types.DKGSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, session := range sm.sessions {
		if session.Phase == types.PhaseCompleted {
			return session
		}
	}

	return nil
}

// sessionRetentionRounds bounds how many recent rounds' sessions are kept. The decrypt path
// only needs the latest activated round and the one before it; 10 is a generous margin (each
// session is a small JSON file) that also covers a brief-outage resume, while bounding growth.
const sessionRetentionRounds = 10

// PruneOldSessions deletes local sessions (in-memory + disk) for rounds older than
// (activeRound - sessionRetentionRounds). Rounds activate sequentially, so activeRound is the
// latest active round and is always kept. Node-local only; the round-based criterion is
// deterministic with no consensus or wall-clock dependence.
func (sm *StateManager) PruneOldSessions(ctx context.Context, activeRound uint32) {
	// Guard against underflow: nothing to prune until enough rounds have elapsed.
	if activeRound <= sessionRetentionRounds {
		return
	}

	horizon := activeRound - sessionRetentionRounds

	sm.mu.RLock()
	stale := make([]uint32, 0)
	for _, session := range sm.sessions {
		if session.Round < horizon {
			stale = append(stale, session.Round)
		}
	}
	sm.mu.RUnlock()

	// DeleteSession removes both the in-memory entry and the disk file, taking the
	// lock itself; the result is independent of map-iteration order.
	for _, round := range stale {
		if err := sm.DeleteSession(ctx, round); err != nil {
			log.Error(ctx, "Failed to prune old DKG session", err, "round", round)
		}
	}
}

// loadSessions loads existing sessions from disk.
func (sm *StateManager) loadSessions() error {
	files, err := filepath.Glob(filepath.Join(sm.dataDir, "session_*.json"))
	if err != nil {
		return errors.Wrap(err, "failed to list session files")
	}

	for _, file := range files {
		session, err := sm.loadSessionFromFile(file)
		if err != nil {
			// Log error but continue loading other sessions
			fmt.Printf("Warning: failed to load session from %s: %v\n", file, err)
			continue
		}

		sm.sessions[session.GetSessionKey()] = session
	}

	return nil
}

// loadSessionFromFile loads a single session from a file.
func (*StateManager) loadSessionFromFile(filename string) (*types.DKGSession, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read session file")
	}

	var session types.DKGSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal session data")
	}

	return &session, nil
}

// saveSession saves a session to disk atomically by writing to a temporary file
// and renaming it, preventing corruption from partial writes during crashes.
func (sm *StateManager) saveSession(session *types.DKGSession) error {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal session data")
	}

	filename := sm.getSessionFilename(session.GetSessionKey())
	tmpFilename := filename + ".tmp"

	if err := os.WriteFile(tmpFilename, data, 0600); err != nil {
		return errors.Wrap(err, "failed to write temporary session file")
	}

	if err := os.Rename(tmpFilename, filename); err != nil {
		// Clean up the temp file on rename failure
		_ = os.Remove(tmpFilename)

		return errors.Wrap(err, "failed to atomically rename session file")
	}

	return nil
}

// getSessionFilename returns the filename for a session.
func (sm *StateManager) getSessionFilename(sessionKey string) string {
	return filepath.Join(sm.dataDir, fmt.Sprintf("session_%s.json", sessionKey))
}
