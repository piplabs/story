package types_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

func TestDKGPhase_String(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		phase    types.DKGPhase
		expected string
	}{
		{name: "unknown", phase: types.PhaseUnknown, expected: "Unknown"},
		{name: "initializing", phase: types.PhaseInitializing, expected: "Initializing"},
		{name: "initialized", phase: types.PhaseInitialized, expected: "Initialized"},
		{name: "dealing", phase: types.PhaseDealing, expected: "Dealing"},
		{name: "finalized", phase: types.PhaseFinalized, expected: "Finalized"},
		{name: "completed", phase: types.PhaseCompleted, expected: "Completed"},
		{name: "failed", phase: types.PhaseFailed, expected: "Failed"},
		{name: "undefined phase falls through to default", phase: types.DKGPhase(99), expected: "Phase(99)"},
		{name: "negative phase falls through to default", phase: types.DKGPhase(-1), expected: "Phase(-1)"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, tc.phase.String())
		})
	}
}

func TestNewDKGSession(t *testing.T) {
	t.Parallel()

	t.Run("basic creation", func(t *testing.T) {
		t.Parallel()

		validators := []string{"0x1111", "0x2222", "0x3333"}
		enclaveType := [32]byte{0xAA, 0xBB}

		session := types.NewDKGSession(5, validators, true, enclaveType)
		require.NotNil(t, session)
		require.Equal(t, uint32(5), session.Round)
		require.Equal(t, types.PhaseInitializing, session.Phase)
		require.Equal(t, validators, session.ActiveValidators)
		require.True(t, session.IsResharing)
		require.Equal(t, enclaveType, session.EnclaveType)
		require.False(t, session.IsFinalized)
		require.NotNil(t, session.GlobalPubKey)
		require.NotNil(t, session.CommPubKey)
		require.Empty(t, session.DecryptRequests)
		require.Equal(t, uint32(0), session.Total)
		require.Equal(t, uint32(0), session.Threshold)
		require.False(t, session.StartTime.IsZero())
		require.False(t, session.LastUpdate.IsZero())
	})

	t.Run("non-resharing session", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		require.NotNil(t, session)
		require.False(t, session.IsResharing)
		require.Nil(t, session.ActiveValidators)
	})

	t.Run("empty validators", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, []string{}, false, [32]byte{})
		require.NotNil(t, session)
		require.Empty(t, session.ActiveValidators)
	})
}

func TestDKGSession_GetCodeCommitmentString(t *testing.T) {
	t.Parallel()

	t.Run("non-empty code commitment", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.CodeCommitment = []byte{0xDE, 0xAD, 0xBE, 0xEF}
		require.Equal(t, "deadbeef", session.GetCodeCommitmentString())
	})

	t.Run("empty code commitment", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		require.Equal(t, "", session.GetCodeCommitmentString())
	})

	t.Run("32 byte code commitment", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.CodeCommitment = make([]byte, 32)
		for i := range session.CodeCommitment {
			session.CodeCommitment[i] = byte(i)
		}
		expected := "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
		require.Equal(t, expected, session.GetCodeCommitmentString())
	})
}

func TestDKGSession_GetSessionKey(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		round    uint32
		expected string
	}{
		{name: "round 0", round: 0, expected: "0"},
		{name: "round 1", round: 1, expected: "1"},
		{name: "round 42", round: 42, expected: "42"},
		{name: "max uint32", round: ^uint32(0), expected: "4294967295"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			session := types.NewDKGSession(tc.round, nil, false, [32]byte{})
			require.Equal(t, tc.expected, session.GetSessionKey())
		})
	}
}

func TestDKGSession_UpdatePhase(t *testing.T) {
	t.Parallel()

	session := types.NewDKGSession(1, nil, false, [32]byte{})
	require.Equal(t, types.PhaseInitializing, session.Phase)
	beforeUpdate := session.LastUpdate

	session.UpdatePhase(types.PhaseDealing)
	require.Equal(t, types.PhaseDealing, session.Phase)
	require.False(t, session.LastUpdate.Before(beforeUpdate))

	session.UpdatePhase(types.PhaseCompleted)
	require.Equal(t, types.PhaseCompleted, session.Phase)

	session.UpdatePhase(types.PhaseFailed)
	require.Equal(t, types.PhaseFailed, session.Phase)
}

func TestDKGSession_DecryptRequests(t *testing.T) {
	t.Parallel()

	t.Run("add and get decrypt requests", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		require.Empty(t, session.GetDecryptRequests())

		req1 := types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct1")}}
		req2 := types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct2")}}

		session.AddDecryptRequest(req1)
		reqs := session.GetDecryptRequests()
		require.Len(t, reqs, 1)
		require.Equal(t, req1.Ciphertext, reqs[0].Ciphertext)

		session.AddDecryptRequest(req2)
		reqs = session.GetDecryptRequests()
		require.Len(t, reqs, 2)
	})

	t.Run("get decrypt requests returns a copy", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct1")}})

		// Modify the returned copy
		reqs := session.GetDecryptRequests()
		reqs[0].Ciphertext = []byte("modified")

		// Original should be unchanged
		original := session.GetDecryptRequests()
		require.Equal(t, []byte("ct1"), original[0].Ciphertext)
	})

	t.Run("set decrypt requests replaces all", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct1")}})
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct2")}})
		require.Len(t, session.GetDecryptRequests(), 2)

		// Replace with only the failed request
		remaining := []types.PendingDecryptRequest{{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct2")}}}
		session.SetDecryptRequests(remaining)

		reqs := session.GetDecryptRequests()
		require.Len(t, reqs, 1)
		require.Equal(t, []byte("ct2"), reqs[0].Ciphertext)
	})

	t.Run("set decrypt requests to nil", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1}})

		session.SetDecryptRequests(nil)
		reqs := session.GetDecryptRequests()
		require.Empty(t, reqs)
	})

	t.Run("drain decrypt requests clears queue and returns contents", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct1")}})
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("ct2")}})

		drained := session.DrainDecryptRequests()
		require.Len(t, drained, 2)
		require.Equal(t, []byte("ct1"), drained[0].Ciphertext)
		require.Equal(t, []byte("ct2"), drained[1].Ciphertext)

		// Queue should be empty after drain
		require.Empty(t, session.GetDecryptRequests())
	})

	t.Run("drain then add preserves new requests", func(t *testing.T) {
		t.Parallel()

		// Simulates the race: worker drains, ABCI adds a new request,
		// worker re-adds failures — the new request must survive.
		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("uuid50")}})

		// Worker drains
		drained := session.DrainDecryptRequests()
		require.Len(t, drained, 1)

		// ABCI thread adds uuid51 while worker is processing uuid50
		session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{Round: 1, Ciphertext: []byte("uuid51")}})

		// Worker finishes uuid50 successfully — no failures to re-add
		// Queue should still contain uuid51
		reqs := session.GetDecryptRequests()
		require.Len(t, reqs, 1)
		require.Equal(t, []byte("uuid51"), reqs[0].Ciphertext)
	})
}

// TestDKGSession_ByteGettersReturnCopies verifies that the byte-slice getters hand out
// deep copies, so a caller mutating the returned slice cannot corrupt session state (which
// would otherwise be an aliasing data race against concurrent readers/writers).
func TestDKGSession_ByteGettersReturnCopies(t *testing.T) {
	t.Parallel()

	session := types.NewDKGSession(1, nil, false, [32]byte{})
	session.SetKeyMaterial(
		[]byte("participants-root"),
		[]byte("global-pub-key"),
		[]byte("sig-finalize"),
		[]byte("pub-key-share"),
		[][]byte{[]byte("coeff-0"), []byte("coeff-1")},
	)

	t.Run("GetGlobalPubKey", func(t *testing.T) {
		t.Parallel()
		got := session.GetGlobalPubKey()
		require.Equal(t, []byte("global-pub-key"), got)
		got[0] = 'X'
		require.Equal(t, []byte("global-pub-key"), session.GetGlobalPubKey(), "mutating the returned slice must not affect the session")
	})

	t.Run("GetSigFinalizeNetwork", func(t *testing.T) {
		t.Parallel()
		got := session.GetSigFinalizeNetwork()
		require.Equal(t, []byte("sig-finalize"), got)
		got[0] = 'X'
		require.Equal(t, []byte("sig-finalize"), session.GetSigFinalizeNetwork())
	})

	t.Run("GetParticipantsRoot", func(t *testing.T) {
		t.Parallel()
		got := session.GetParticipantsRoot()
		require.Equal(t, []byte("participants-root"), got)
		got[0] = 'X'
		require.Equal(t, []byte("participants-root"), session.GetParticipantsRoot())
	})

	t.Run("GetPubKeyShare", func(t *testing.T) {
		t.Parallel()
		got := session.GetPubKeyShare()
		require.Equal(t, []byte("pub-key-share"), got)
		got[0] = 'X'
		require.Equal(t, []byte("pub-key-share"), session.GetPubKeyShare())
	})

	t.Run("GetPublicCoeffs", func(t *testing.T) {
		t.Parallel()
		got := session.GetPublicCoeffs()
		require.Equal(t, [][]byte{[]byte("coeff-0"), []byte("coeff-1")}, got)
		got[0][0] = 'X'
		require.Equal(t, [][]byte{[]byte("coeff-0"), []byte("coeff-1")}, session.GetPublicCoeffs(), "mutating an inner coefficient slice must not affect the session")
	})
}

// TestDKGSession_ScalarGettersSetters verifies the mutex-guarded scalar accessors round-trip.
func TestDKGSession_ScalarGettersSetters(t *testing.T) {
	t.Parallel()

	session := types.NewDKGSession(1, nil, false, [32]byte{})

	require.Equal(t, types.PhaseInitializing, session.GetPhase())
	session.UpdatePhase(types.PhaseCompleted)
	require.Equal(t, types.PhaseCompleted, session.GetPhase())

	require.Equal(t, uint32(0), session.GetIndex())
	session.SetIndex(5)
	require.Equal(t, uint32(5), session.GetIndex())

	require.False(t, session.GetIsFinalized())
	session.SetFinalized()
	require.True(t, session.GetIsFinalized())
}

// TestDKGSession_Snapshot verifies that Snapshot returns an independent deep copy: field
// values match, and mutating the snapshot's byte slices does not touch the live session.
func TestDKGSession_Snapshot(t *testing.T) {
	t.Parallel()

	session := types.NewDKGSession(7, []string{"0xval"}, true, [32]byte{0xAB})
	session.UpdatePhase(types.PhaseDealing)
	session.SetIndex(3)
	session.SetKeyMaterial(
		[]byte("proot"),
		[]byte("gpk"),
		[]byte("sig"),
		[]byte("share"),
		[][]byte{[]byte("c0")},
	)

	snap := session.Snapshot()

	require.Equal(t, uint32(7), snap.Round)
	require.Equal(t, types.PhaseDealing, snap.Phase)
	require.Equal(t, uint32(3), snap.Index)
	require.True(t, snap.IsResharing)
	require.Equal(t, []byte("gpk"), snap.GlobalPubKey)
	require.Equal(t, [][]byte{[]byte("c0")}, snap.PublicCoeffs)

	// Mutating the snapshot's slices must not affect the live session.
	snap.GlobalPubKey[0] = 'X'
	snap.PublicCoeffs[0][0] = 'X'
	require.Equal(t, []byte("gpk"), session.GetGlobalPubKey())
	require.Equal(t, [][]byte{[]byte("c0")}, session.GetPublicCoeffs())
}

// TestDKGSession_ConcurrentSetupMutationAndRead runs the atomic setup accessors under
// concurrent writers and readers so -race guards the atomicity their docs claim: a
// reader must never observe a torn slice header while SetSetupResult writes the group.
func TestDKGSession_ConcurrentSetupMutationAndRead(t *testing.T) {
	t.Parallel()

	session := types.NewDKGSession(1, nil, false, [32]byte{})

	const (
		goroutines = 8
		iterations = 1000
	)

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				session.SetCodeCommitmentIfEmpty([]byte("cc"))
				session.SetSetupResult(
					[]byte("cc"),
					[]byte("dkgpub"),
					[]byte("commpub"),
					[]byte("report"),
					[]byte("blockhash"),
					int64(42),
				)
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_ = session.HasSetupData()
				_ = session.GetCodeCommitment()
				_ = session.GetDKGPubKey()
			}
		}()
	}

	wg.Wait()

	// After all writers, the setup result is fully populated and self-consistent.
	require.True(t, session.HasSetupData())
	require.Equal(t, []byte("cc"), session.GetCodeCommitment())
	require.Equal(t, []byte("dkgpub"), session.GetDKGPubKey())
	require.Equal(t, int64(42), session.GetStartBlockHeight())
}
