package keeper

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// newTestStateManager creates a StateManager backed by a temporary directory
// that is cleaned up when the test finishes.
func newTestStateManager(t *testing.T) *StateManager {
	t.Helper()

	dir := t.TempDir()
	sm, err := NewStateManager(dir)
	require.NoError(t, err)

	return sm
}

// newTestSession builds a minimal DKGSession for the given round.
func newTestSession(round uint32) *types.DKGSession {
	return types.NewDKGSession(round, []string{"0xaabbcc"}, false, [32]byte{})
}

// TestStateManager_NewStateManager_CreatesDir verifies that NewStateManager
// creates the data directory when it does not exist.
func TestStateManager_NewStateManager_CreatesDir(t *testing.T) {
	t.Parallel()

	dir := filepath.Join(t.TempDir(), "subdir", "deep")
	sm, err := NewStateManager(dir)
	require.NoError(t, err)
	require.NotNil(t, sm)

	_, err = os.Stat(dir)
	require.NoError(t, err, "directory should have been created")
}

// TestStateManager_NewStateManager_ExistingDir verifies that NewStateManager
// succeeds when the directory already exists.
func TestStateManager_NewStateManager_ExistingDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir() // already exists
	sm, err := NewStateManager(dir)
	require.NoError(t, err)
	require.NotNil(t, sm)
}

// TestStateManager_CreateAndGetSession verifies the basic create → get cycle.
func TestStateManager_CreateAndGetSession(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()
	session := newTestSession(1)

	require.NoError(t, sm.CreateSession(ctx, session))

	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, uint32(1), got.Round)
}

// TestStateManager_CreateSession_DuplicateIsNoop verifies that creating the
// same session twice does not return an error and preserves the first session.
func TestStateManager_CreateSession_DuplicateIsNoop(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session1 := newTestSession(2)
	require.NoError(t, sm.CreateSession(ctx, session1))

	// Creating the same session again should be a no-op (not an error).
	session1Again := newTestSession(2)
	require.NoError(t, sm.CreateSession(ctx, session1Again))

	sessions := sm.ListSessions()
	require.Len(t, sessions, 1, "duplicate create should not add a second entry")
}

// TestStateManager_GetSession_NotFound verifies that GetSession returns an
// error when the session key does not exist.
func TestStateManager_GetSession_NotFound(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)

	_, err := sm.GetSession(999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "session not found")
}

// TestStateManager_UpdateSession verifies that UpdateSession replaces the
// stored session and persists it to disk.
func TestStateManager_UpdateSession(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session := newTestSession(3)
	require.NoError(t, sm.CreateSession(ctx, session))

	// Update the phase
	session.UpdatePhase(types.PhaseDealing)
	require.NoError(t, sm.UpdateSession(ctx, session))

	got, err := sm.GetSession(3)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase)
}

// TestStateManager_UpdateSession_NotFound verifies that UpdateSession returns
// an error when the session does not exist yet.
func TestStateManager_UpdateSession_NotFound(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session := newTestSession(42)
	err := sm.UpdateSession(ctx, session)
	require.Error(t, err)
	require.Contains(t, err.Error(), "session not found")
}

// TestStateManager_DeleteSession verifies that DeleteSession removes the
// session from memory and disk.
func TestStateManager_DeleteSession(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session := newTestSession(5)
	require.NoError(t, sm.CreateSession(ctx, session))

	require.NoError(t, sm.DeleteSession(ctx, 5))

	_, err := sm.GetSession(5)
	require.Error(t, err)
	require.Contains(t, err.Error(), "session not found")
}

// TestStateManager_DeleteSession_NotFound verifies that DeleteSession returns
// an error when the session does not exist.
func TestStateManager_DeleteSession_NotFound(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	err := sm.DeleteSession(ctx, 100)
	require.Error(t, err)
	require.Contains(t, err.Error(), "session not found")
}

// TestStateManager_ListSessions verifies that ListSessions returns all stored
// sessions.
func TestStateManager_ListSessions(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	require.NoError(t, sm.CreateSession(ctx, newTestSession(1)))
	require.NoError(t, sm.CreateSession(ctx, newTestSession(2)))
	require.NoError(t, sm.CreateSession(ctx, newTestSession(3)))

	sessions := sm.ListSessions()
	require.Len(t, sessions, 3)
}

// TestStateManager_ListSessions_Empty verifies that ListSessions returns an
// empty slice when no sessions exist.
func TestStateManager_ListSessions_Empty(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	sessions := sm.ListSessions()
	require.Empty(t, sessions)
}

// TestStateManager_GetActiveSession_NoCompletedSessions verifies that
// GetActiveSession returns nil when no session is in PhaseCompleted.
func TestStateManager_GetActiveSession_NoCompletedSessions(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session := newTestSession(1)
	require.NoError(t, sm.CreateSession(ctx, session))

	// Session is in PhaseInitializing by default
	got := sm.GetActiveSession("0xvalidator")
	require.Nil(t, got, "no active session when none is Completed")
}

// TestStateManager_GetActiveSession_CompletedSession verifies that
// GetActiveSession returns the session in PhaseCompleted.
func TestStateManager_GetActiveSession_CompletedSession(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session := newTestSession(7)
	require.NoError(t, sm.CreateSession(ctx, session))
	session.UpdatePhase(types.PhaseCompleted)
	require.NoError(t, sm.UpdateSession(ctx, session))

	got := sm.GetActiveSession("0xvalidator")
	require.NotNil(t, got)
	require.Equal(t, uint32(7), got.Round)
}

// TestStateManager_MarkFailed verifies that MarkFailed transitions the phase
// to PhaseFailed and persists the update.
func TestStateManager_MarkFailed(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session := newTestSession(8)
	require.NoError(t, sm.CreateSession(ctx, session))

	sm.MarkFailed(ctx, session)

	got, err := sm.GetSession(8)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase)
}

// TestStateManager_CleanupExpiredSessions verifies that CleanupExpiredSessions
// removes sessions in PhaseCompleted or PhaseFailed.
func TestStateManager_CleanupExpiredSessions(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	// Create 3 sessions: one completed, one failed, one active
	completed := newTestSession(1)
	require.NoError(t, sm.CreateSession(ctx, completed))
	completed.UpdatePhase(types.PhaseCompleted)
	require.NoError(t, sm.UpdateSession(ctx, completed))

	failed := newTestSession(2)
	require.NoError(t, sm.CreateSession(ctx, failed))
	failed.UpdatePhase(types.PhaseFailed)
	require.NoError(t, sm.UpdateSession(ctx, failed))

	active := newTestSession(3)
	require.NoError(t, sm.CreateSession(ctx, active))
	active.UpdatePhase(types.PhaseDealing)
	require.NoError(t, sm.UpdateSession(ctx, active))

	sm.CleanupExpiredSessions(ctx)

	sessions := sm.ListSessions()
	require.Len(t, sessions, 1, "only the dealing session should remain")
	require.Equal(t, uint32(3), sessions[0].Round)
}

// TestStateManager_PersistenceAcrossReload verifies that sessions stored to
// disk are re-loaded by a new StateManager instance pointing at the same dir.
func TestStateManager_PersistenceAcrossReload(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := context.Background()

	// Write session with first StateManager
	sm1, err := NewStateManager(dir)
	require.NoError(t, err)
	session := newTestSession(10)
	require.NoError(t, sm1.CreateSession(ctx, session))

	// Load with a second StateManager pointing at the same directory
	sm2, err := NewStateManager(dir)
	require.NoError(t, err)

	got, err := sm2.GetSession(10)
	require.NoError(t, err)
	require.Equal(t, uint32(10), got.Round)
}

// TestStateManager_DeleteSession_FileRemovedFromDisk verifies that after
// DeleteSession, the corresponding JSON file is removed from disk.
func TestStateManager_DeleteSession_FileRemovedFromDisk(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := context.Background()

	sm, err := NewStateManager(dir)
	require.NoError(t, err)

	session := newTestSession(11)
	require.NoError(t, sm.CreateSession(ctx, session))

	// Verify file exists
	pattern := filepath.Join(dir, "session_11.json")
	_, err = os.Stat(pattern)
	require.NoError(t, err, "session file should exist after create")

	require.NoError(t, sm.DeleteSession(ctx, 11))

	_, err = os.Stat(pattern)
	require.True(t, os.IsNotExist(err), "session file should be removed after delete")
}
