package keeper

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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

// TestStateManager_PruneOldSessions verifies that after several rounds, sessions
// older than the retention horizon are removed from memory and disk while the active
// round and the retained recent rounds survive.
func TestStateManager_PruneOldSessions(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := context.Background()

	sm, err := NewStateManager(dir)
	require.NoError(t, err)

	// Advance through several rounds.
	const activeRound = 20
	for r := uint32(1); r <= activeRound; r++ {
		require.NoError(t, sm.CreateSession(ctx, newTestSession(r)))
	}

	sm.PruneOldSessions(ctx, activeRound)

	// Retention horizon is activeRound - sessionRetentionRounds. Rounds below it are
	// pruned; the horizon round and everything up to the active round survive.
	horizon := uint32(activeRound - sessionRetentionRounds)

	for r := uint32(1); r < horizon; r++ {
		_, err := sm.GetSession(r)
		require.Error(t, err, "round %d should have been pruned from memory", r)

		_, statErr := os.Stat(filepath.Join(dir, "session_"+strconv.FormatUint(uint64(r), 10)+".json"))
		require.True(t, os.IsNotExist(statErr), "round %d file should have been pruned from disk", r)
	}

	for r := horizon; r <= activeRound; r++ {
		got, err := sm.GetSession(r)
		require.NoError(t, err, "round %d should have survived", r)
		require.Equal(t, r, got.Round)

		_, statErr := os.Stat(filepath.Join(dir, "session_"+strconv.FormatUint(uint64(r), 10)+".json"))
		require.NoError(t, statErr, "round %d file should have survived on disk", r)
	}
}

// TestStateManager_PruneOldSessions_RetainedRoundStillUsable verifies that a retained
// recent round's session remains usable by the decrypt path after pruning: its
// GlobalPubKey is preserved both in memory and after a reload from disk.
func TestStateManager_PruneOldSessions_RetainedRoundStillUsable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := context.Background()

	sm, err := NewStateManager(dir)
	require.NoError(t, err)

	const activeRound = 20
	pubKey := []byte{0x01, 0x02, 0x03, 0x04}

	for r := uint32(1); r <= activeRound; r++ {
		session := newTestSession(r)
		// The round just below the active round is still needed by the decrypt path.
		if r == activeRound-1 {
			session.GlobalPubKey = pubKey
			session.Index = 1
			session.IsFinalized = true
		}
		require.NoError(t, sm.CreateSession(ctx, session))
	}

	sm.PruneOldSessions(ctx, activeRound)

	got, err := sm.GetSession(activeRound - 1)
	require.NoError(t, err, "the round below active must be retained for the decrypt path")
	require.Equal(t, pubKey, got.GlobalPubKey, "GlobalPubKey must be preserved for decryption")

	// The preserved key must also survive a reload from disk.
	sm2, err := NewStateManager(dir)
	require.NoError(t, err)
	reloaded, err := sm2.GetSession(activeRound - 1)
	require.NoError(t, err)
	require.Equal(t, pubKey, reloaded.GlobalPubKey)
}

// TestStateManager_PruneOldSessions_NoopBelowHorizon verifies that pruning is a no-op
// while fewer rounds than the retention window have elapsed (underflow guard).
func TestStateManager_PruneOldSessions_NoopBelowHorizon(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	for r := uint32(1); r <= sessionRetentionRounds; r++ {
		require.NoError(t, sm.CreateSession(ctx, newTestSession(r)))
	}

	sm.PruneOldSessions(ctx, sessionRetentionRounds)

	require.Len(t, sm.ListSessions(), sessionRetentionRounds, "nothing should be pruned below the horizon")
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

// TestStateManager_LoadSessions_SweepsOrphanedTmpFiles verifies that orphaned
// session_<round>.json.tmp files (left by a crash between saveSession's write and
// rename) are removed on construction, while valid .json sessions still load.
func TestStateManager_LoadSessions_SweepsOrphanedTmpFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := context.Background()

	// Persist a valid session to disk with a first StateManager.
	sm1, err := NewStateManager(dir)
	require.NoError(t, err)
	require.NoError(t, sm1.CreateSession(ctx, newTestSession(1)))

	// Simulate crashed writes: orphaned .tmp files for several rounds.
	orphans := []string{
		filepath.Join(dir, "session_2.json.tmp"),
		filepath.Join(dir, "session_3.json.tmp"),
	}
	for _, f := range orphans {
		require.NoError(t, os.WriteFile(f, []byte("{ partial"), 0600))
	}

	// Constructing a new StateManager triggers loadSessions, which sweeps orphans.
	sm2, err := NewStateManager(dir)
	require.NoError(t, err)

	// Orphaned .tmp files must be gone.
	for _, f := range orphans {
		_, statErr := os.Stat(f)
		require.True(t, os.IsNotExist(statErr), "orphaned tmp file %s should be swept", f)
	}

	// The valid session must still load, and no phantom sessions from the .tmp files.
	got, err := sm2.GetSession(1)
	require.NoError(t, err)
	require.Equal(t, uint32(1), got.Round)
	require.Len(t, sm2.ListSessions(), 1, "only the valid session should be loaded")
}

// TestStateManager_DeleteSession_RemovesStrayTmpFile verifies that DeleteSession
// also removes a stray temporary file for the deleted round.
func TestStateManager_DeleteSession_RemovesStrayTmpFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := context.Background()

	sm, err := NewStateManager(dir)
	require.NoError(t, err)

	require.NoError(t, sm.CreateSession(ctx, newTestSession(4)))

	// A stray tmp file appears for the same round after creation.
	tmpFile := filepath.Join(dir, "session_4.json.tmp")
	require.NoError(t, os.WriteFile(tmpFile, []byte("{ partial"), 0600))

	require.NoError(t, sm.DeleteSession(ctx, 4))

	_, statErr := os.Stat(tmpFile)
	require.True(t, os.IsNotExist(statErr), "stray tmp file should be removed on delete")

	_, statErr = os.Stat(filepath.Join(dir, "session_4.json"))
	require.True(t, os.IsNotExist(statErr), "session file should be removed on delete")
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

// TestStateManager_UpdateSession_ConcurrentMarshalAndMutation drives concurrent
// UpdateSession (which marshals the session to disk) against mutex-guarded mutations of the
// same session. Before saveSession marshaled a locked Snapshot, the marshal reflected over
// the live struct without holding session.mu and raced these mutations, so this test fails
// under -race on the pre-fix code. Run with: go test ./client/x/dkg/... -race.
func TestStateManager_UpdateSession_ConcurrentMarshalAndMutation(t *testing.T) {
	t.Parallel()

	sm := newTestStateManager(t)
	ctx := context.Background()

	session := newTestSession(42)
	require.NoError(t, sm.CreateSession(ctx, session))

	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(2)

	// Writer 1: repeatedly persist the session; saveSession marshals a locked Snapshot.
	go func() {
		defer wg.Done()
		for range iterations {
			_ = sm.UpdateSession(ctx, session)
		}
	}()

	// Writer 2: mutate the same session through its mutex-guarded methods.
	go func() {
		defer wg.Done()
		for i := range iterations {
			session.UpdatePhase(types.PhaseDealing)
			session.SetIndex(uint32(i + 1))
			session.SetKeyMaterial(
				[]byte("proot"),
				[]byte("gpk"),
				[]byte("sig"),
				[]byte("share"),
				[][]byte{[]byte("c0"), []byte("c1")},
			)
			session.IncrementRecoveryAttempts()
			session.SetFinalized()
		}
	}()

	wg.Wait()

	// Sanity: the session is still persisted and readable after the concurrent churn.
	got, err := sm.GetSession(42)
	require.NoError(t, err)
	require.Equal(t, uint32(42), got.Round)
}
