package types_test

import (
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

		req1 := types.DecryptRequest{Round: 1, Ciphertext: []byte("ct1")}
		req2 := types.DecryptRequest{Round: 1, Ciphertext: []byte("ct2")}

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
		session.AddDecryptRequest(types.DecryptRequest{Round: 1, Ciphertext: []byte("ct1")})

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
		session.AddDecryptRequest(types.DecryptRequest{Round: 1, Ciphertext: []byte("ct1")})
		session.AddDecryptRequest(types.DecryptRequest{Round: 1, Ciphertext: []byte("ct2")})
		require.Len(t, session.GetDecryptRequests(), 2)

		// Replace with only the failed request
		remaining := []types.DecryptRequest{{Round: 1, Ciphertext: []byte("ct2")}}
		session.SetDecryptRequests(remaining)

		reqs := session.GetDecryptRequests()
		require.Len(t, reqs, 1)
		require.Equal(t, []byte("ct2"), reqs[0].Ciphertext)
	})

	t.Run("set decrypt requests to nil", func(t *testing.T) {
		t.Parallel()

		session := types.NewDKGSession(1, nil, false, [32]byte{})
		session.AddDecryptRequest(types.DecryptRequest{Round: 1})

		session.SetDecryptRequests(nil)
		reqs := session.GetDecryptRequests()
		require.Empty(t, reqs)
	})
}
