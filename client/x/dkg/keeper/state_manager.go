package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/piplabs/story/client/x/dkg/types"
	"os"
	"path/filepath"
	"sync"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// StateManager manages the local state of DKG sessions outside of consensus.
//
//nolint:revive // use full name
type StateManager struct {
	dataDir  string
	mu       sync.RWMutex
	sessions map[string]*types.DKGSession // keyed by session key (mrenclave_round)
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
		log.Info(ctx, "session already exists with the mrenclave and round. skip creating a new session", "mrenclave", session.GetMrenclaveString(), "round", session.Round)

		return nil
	} else {
		sm.sessions[sessionKey] = session

		if err := sm.saveSession(session); err != nil {
			return errors.Wrap(err, "failed to save session to disk")
		}

		log.Info(ctx, "Created DKG session",
			"mrenclave", session.GetMrenclaveString(),
			"round", session.Round,
			"phase", session.Phase.String(),
		)

		return nil
	}
}

// GetSession retrieves a DKG session by its mrenclave and round.
func (sm *StateManager) GetSession(mrenclave []byte, round uint32) (*types.DKGSession, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessionKey := fmt.Sprintf("%x_%d", mrenclave, round)
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
		"mrenclave", session.GetMrenclaveString(),
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
func (sm *StateManager) DeleteSession(ctx context.Context, mrenclave []byte, round uint32) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionKey := fmt.Sprintf("%x_%d", mrenclave, round)

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
		"mrenclave", session.GetMrenclaveString(),
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

// CleanupExpiredSessions removes expired sessions.
func (sm *StateManager) CleanupExpiredSessions(ctx context.Context) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var expired []string
	for key, session := range sm.sessions {
		// TODO: more rigorous expiration check
		if session.Phase == types.PhaseCompleted || session.Phase == types.PhaseFailed {
			expired = append(expired, key)
		}
	}

	for _, key := range expired {
		session := sm.sessions[key]
		delete(sm.sessions, key)

		// Remove from disk
		filename := sm.getSessionFilename(key)
		if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
			log.Error(ctx, "Failed to delete expired session file", err, "filename", filename)
		} else {
			log.Info(ctx, "Cleaned up expired DKG session",
				"mrenclave", session.GetMrenclaveString(),
				"round", session.Round,
				"phase", session.Phase.String(),
			)
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

// saveSession saves a session to disk.
func (sm *StateManager) saveSession(session *types.DKGSession) error {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal session data")
	}

	filename := sm.getSessionFilename(session.GetSessionKey())
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return errors.Wrap(err, "failed to write session file")
	}

	return nil
}

// getSessionFilename returns the filename for a session.
func (sm *StateManager) getSessionFilename(sessionKey string) string {
	return filepath.Join(sm.dataDir, fmt.Sprintf("session_%s.json", sessionKey))
}
