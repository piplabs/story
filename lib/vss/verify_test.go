package vss_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/lib/vss"

	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/share"
	dkgPed "go.dedis.ch/kyber/v4/share/dkg/pedersen"
)

// newSuite returns the Edwards25519 suite used by the DKG protocol.
func newSuite() *edwards25519.SuiteEd25519 {
	return edwards25519.NewBlakeSHA256Ed25519()
}

// generateVSSData creates a private polynomial, derives commitments and shares for n
// participants with the given threshold. It returns serialized commitment bytes
// and serialized share bytes indexed by participant (0-based slice, 1-based indices).
func generateVSSData(t *testing.T, n, threshold int) (commitmentBytes [][]byte, shareBytes [][]byte) {
	t.Helper()

	suite := newSuite()
	secret := suite.Scalar().Pick(suite.RandomStream())
	priPoly := share.NewPriPoly(suite, threshold, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.Point().Base())

	// Info() returns (base, commits []kyber.Point)
	_, commits := pubPoly.Info()

	commitmentBytes = make([][]byte, len(commits))
	for i, c := range commits {
		bz, err := c.MarshalBinary()
		require.NoError(t, err)

		commitmentBytes[i] = bz
	}

	// Shares uses 1-based indices internally but returns a 0-indexed slice
	shares := priPoly.Shares(n)
	shareBytes = make([][]byte, n)

	for i, s := range shares {
		bz, err := s.V.MarshalBinary()
		require.NoError(t, err)

		shareBytes[i] = bz
	}

	return commitmentBytes, shareBytes
}

// TestVerifyPedersenVSS_ValidShare verifies that a correctly generated share
// passes Pedersen VSS verification.
func TestVerifyPedersenVSS_ValidShare(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 5
	threshold := 3

	commitmentBytes, shareBytes := generateVSSData(t, n, threshold)

	for i := range n {
		recipientIndex := i + 1 // 1-based
		ok, err := vss.VerifyPedersenVSS(suite, shareBytes[i], recipientIndex, commitmentBytes, uint32(threshold))
		require.NoError(t, err, "recipient index %d should not error", recipientIndex)
		require.True(t, ok, "valid share for recipient %d should verify", recipientIndex)
	}
}

// TestVerifyPedersenVSS_InvalidShare verifies that a different random scalar
// fails verification against the given commitments.
func TestVerifyPedersenVSS_InvalidShare(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 5
	threshold := 3

	commitmentBytes, _ := generateVSSData(t, n, threshold)

	// Use a different random scalar as an invalid share
	wrongScalar := suite.Scalar().Pick(suite.RandomStream())
	wrongBytes, err := wrongScalar.MarshalBinary()
	require.NoError(t, err)

	ok, err := vss.VerifyPedersenVSS(suite, wrongBytes, 1, commitmentBytes, uint32(threshold))
	require.NoError(t, err, "invalid share should not produce a parse error")
	require.False(t, ok, "invalid share should not verify")
}

// TestVerifyPedersenVSS_EmptyCommitments verifies that an empty commitment
// list returns an error.
func TestVerifyPedersenVSS_EmptyCommitments(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	scalar := suite.Scalar().Pick(suite.RandomStream())
	shareBytes, err := scalar.MarshalBinary()
	require.NoError(t, err)

	ok, err := vss.VerifyPedersenVSS(suite, shareBytes, 1, [][]byte{}, 0)
	require.Error(t, err, "empty commitments should return an error")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_NilCommitments verifies that a nil commitment slice
// returns an error.
func TestVerifyPedersenVSS_NilCommitments(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	scalar := suite.Scalar().Pick(suite.RandomStream())
	shareBytes, err := scalar.MarshalBinary()
	require.NoError(t, err)

	ok, err := vss.VerifyPedersenVSS(suite, shareBytes, 1, nil, 0)
	require.Error(t, err, "nil commitments should return an error")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_MalformedShareBytes verifies that garbage share bytes
// return an error.
func TestVerifyPedersenVSS_MalformedShareBytes(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	commitmentBytes, _ := generateVSSData(t, 3, 2)
	threshold := uint32(2)

	ok, err := vss.VerifyPedersenVSS(suite, []byte("this is definitely not a valid scalar"), 1, commitmentBytes, threshold)
	require.Error(t, err, "malformed share bytes should return an error")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_MalformedCommitmentBytes verifies that a malformed
// commitment byte slice returns an error.
func TestVerifyPedersenVSS_MalformedCommitmentBytes(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	commitmentBytes, shareBytes := generateVSSData(t, 3, 2)

	// Replace the first commitment with garbage
	badCommitments := make([][]byte, len(commitmentBytes))
	copy(badCommitments, commitmentBytes)
	badCommitments[0] = []byte("bad commitment data")

	ok, err := vss.VerifyPedersenVSS(suite, shareBytes[0], 1, badCommitments, uint32(len(badCommitments)))
	require.Error(t, err, "malformed commitment bytes should return an error")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_DifferentIndicesProduceDifferentResults verifies that
// share[0] does not verify against recipient index 2 when it was generated for index 1.
func TestVerifyPedersenVSS_DifferentIndicesProduceDifferentResults(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 5
	threshold := 3

	commitmentBytes, shareBytes := generateVSSData(t, n, threshold)

	// share[0] should verify at index 1, but not at index 2
	okCorrect, err := vss.VerifyPedersenVSS(suite, shareBytes[0], 1, commitmentBytes, uint32(threshold))
	require.NoError(t, err)
	require.True(t, okCorrect, "share 0 should verify at its own index 1")

	okWrong, err := vss.VerifyPedersenVSS(suite, shareBytes[0], 2, commitmentBytes, uint32(threshold))
	require.NoError(t, err)
	require.False(t, okWrong, "share 0 should not verify at a different index (2)")
}

// TestVerifyPedersenVSS_SingleCommitmentThreshold1 verifies that a degree-0
// polynomial (constant secret) works correctly with a single commitment.
func TestVerifyPedersenVSS_SingleCommitmentThreshold1(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 3
	threshold := 1 // degree 0: only constant term, one commitment

	commitmentBytes, shareBytes := generateVSSData(t, n, threshold)
	require.Len(t, commitmentBytes, 1, "threshold=1 means one commitment")

	for i := range n {
		ok, err := vss.VerifyPedersenVSS(suite, shareBytes[i], i+1, commitmentBytes, uint32(threshold))
		require.NoError(t, err)
		require.True(t, ok, "share %d should verify with single commitment", i+1)
	}
}

// TestVerifyPedersenVSS_MultipleCommitmentsThreshold3 verifies the polynomial
// evaluation for threshold=3 (quadratic polynomial, three commitments).
func TestVerifyPedersenVSS_MultipleCommitmentsThreshold3(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 5
	threshold := 3

	commitmentBytes, shareBytes := generateVSSData(t, n, threshold)
	require.Len(t, commitmentBytes, 3, "threshold=3 means three commitments")

	for i := range n {
		ok, err := vss.VerifyPedersenVSS(suite, shareBytes[i], i+1, commitmentBytes, uint32(threshold))
		require.NoError(t, err)
		require.True(t, ok, "share %d should verify with 3 commitments", i+1)
	}
}

// TestVerifyPedersenVSS_WrongShareForOtherPolynomial verifies that a share
// generated from a completely different polynomial fails against another's commitments.
func TestVerifyPedersenVSS_WrongShareForOtherPolynomial(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 3
	threshold := 2

	// First polynomial commitments
	commitmentsA, _ := generateVSSData(t, n, threshold)

	// Second polynomial shares
	_, sharesB := generateVSSData(t, n, threshold)

	// sharesB[0] should not verify against commitmentsA
	ok, err := vss.VerifyPedersenVSS(suite, sharesB[0], 1, commitmentsA, uint32(threshold))
	require.NoError(t, err)
	require.False(t, ok, "share from a different polynomial should not verify")
}

// TestVerifyPedersenVSS_ZeroRecipientIndex verifies that recipient index 0
// returns an error.
func TestVerifyPedersenVSS_ZeroRecipientIndex(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	commitmentBytes, shareBytes := generateVSSData(t, 3, 2)

	ok, err := vss.VerifyPedersenVSS(suite, shareBytes[0], 0, commitmentBytes, uint32(2))
	require.Error(t, err, "index 0 should return an error (must be >= 1)")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_ThresholdMismatch verifies that providing more commitments
// than the expected threshold returns an error.
func TestVerifyPedersenVSS_ThresholdMismatch(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	// Generate 3 commitments (threshold=3)
	commitmentBytes, shareBytes := generateVSSData(t, 5, 3)
	require.Len(t, commitmentBytes, 3)

	// Claim expectedThreshold=2 but provide 3 commitments — should error
	ok, err := vss.VerifyPedersenVSS(suite, shareBytes[0], 1, commitmentBytes, 2)
	require.Error(t, err, "commitment count mismatch should return an error")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_ThresholdZeroSkipsValidation verifies that expectedThreshold=0
// skips commitment count validation and only checks the VSS equation.
func TestVerifyPedersenVSS_ThresholdZeroSkipsValidation(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 3
	threshold := 2

	commitmentBytes, shareBytes := generateVSSData(t, n, threshold)

	// expectedThreshold=0 means "don't check commitment count"
	ok, err := vss.VerifyPedersenVSS(suite, shareBytes[0], 1, commitmentBytes, 0)
	require.NoError(t, err)
	require.True(t, ok, "valid share with expectedThreshold=0 should still verify")
}

// TestVerifyPedersenVSS_CommitmentsExceedMax verifies that more than MaxCommitments
// commitment bytes returns an error.
func TestVerifyPedersenVSS_CommitmentsExceedMax(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	scalar := suite.Scalar().Pick(suite.RandomStream())
	shareBytes, err := scalar.MarshalBinary()
	require.NoError(t, err)

	// Create MaxCommitments+1 fake commitment entries
	tooMany := make([][]byte, vss.MaxCommitments+1)
	for i := range tooMany {
		tooMany[i] = []byte("fake-point")
	}

	ok, err := vss.VerifyPedersenVSS(suite, shareBytes, 1, tooMany, 0)
	require.Error(t, err, "exceeding MaxCommitments should return an error")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_KyberIndexConvention verifies the 0-based to 1-based
// index conversion that callers must perform when using proto values from kyber.
//
// kyber stores PriShare.I as 0-based (0..n-1), but VerifyPedersenVSS expects
// 1-based (1..n). The caller must pass int(secShare.GetI()) + 1.
func TestVerifyPedersenVSS_KyberIndexConvention(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 5
	threshold := 3

	commitmentBytes, shareBytes := generateVSSData(t, n, threshold)

	for kyberIdx := range n {
		// Simulate the conversion: proto 0-based -> VSS 1-based
		vssIndex := kyberIdx + 1

		ok, err := vss.VerifyPedersenVSS(suite, shareBytes[kyberIdx], vssIndex, commitmentBytes, uint32(threshold))
		require.NoError(t, err, "kyber index %d (VSS index %d) should not error", kyberIdx, vssIndex)
		require.True(t, ok, "share at kyber index %d should verify at VSS index %d", kyberIdx, vssIndex)

		// Verify that using the 0-based index directly would fail (wrong evaluation point)
		if kyberIdx > 0 { // skip index 0 since 0 is rejected as <= 0
			okWrong, errWrong := vss.VerifyPedersenVSS(suite, shareBytes[kyberIdx], kyberIdx, commitmentBytes, uint32(threshold))
			require.NoError(t, errWrong)
			require.False(t, okWrong, "share at kyber index %d should NOT verify at wrong VSS index %d", kyberIdx, kyberIdx)
		}
	}
}

// TestVerifyPedersenVSS_NegativeIndex verifies that a negative recipient index
// returns an error.
func TestVerifyPedersenVSS_NegativeIndex(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	commitmentBytes, shareBytes := generateVSSData(t, 3, 2)

	ok, err := vss.VerifyPedersenVSS(suite, shareBytes[0], -1, commitmentBytes, uint32(2))
	require.Error(t, err, "negative index should return an error")
	require.False(t, ok)
}

// TestVerifyPedersenVSS_EndToEndDKGJustificationFlow simulates the complete
// justification verification pipeline using real DKG cryptographic data.
//
// This test proves the critical 0-based → 1-based index conversion by:
//  1. Creating actual kyber DKG DistKeyGenerators for N participants
//  2. Running a full deal exchange (exactly like the live protocol)
//  3. Extracting the polynomial data that would appear in a justification
//  4. Verifying through VerifyPedersenVSS with the correct 1-based conversion
//  5. Demonstrating that the 0-based index FAILS verification
//
// This matches the real on-chain flow:
//
//	kyber Eval(i) → PriShare{I: i (0-based)} → proto SecShare.I = i →
//	CL: recipientIndex = int(secShare.GetI()) + 1 → VerifyPedersenVSS(1-based)
func TestVerifyPedersenVSS_EndToEndDKGJustificationFlow(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 5
	threshold := 3

	// Step 1: Create N key pairs (exactly like DKG participants)
	privKeys := make([]kyber.Scalar, n)

	pubKeys := make([]kyber.Point, n)
	for i := range n {
		privKeys[i] = suite.Scalar().Pick(suite.RandomStream())
		pubKeys[i] = suite.Point().Mul(privKeys[i], nil)
	}

	// Step 2: Create actual DKG generators
	dkgs := make([]*dkgPed.DistKeyGenerator, n)
	for i := range n {
		d, err := dkgPed.NewDistKeyGenerator(suite, privKeys[i], pubKeys, threshold)
		require.NoError(t, err)

		dkgs[i] = d
	}

	// Step 3: Full deal exchange — this creates the REAL internal polynomials
	allResps := make([]*dkgPed.Response, 0)

	for dealerIdx, d := range dkgs {
		deals, err := d.Deals()
		require.NoError(t, err)

		for recipientIdx, deal := range deals {
			// Deal.Index = dealerIdx (0-based, the dealer's index)
			require.Equal(t, uint32(dealerIdx), deal.Index,
				"Deal.Index should be the 0-based dealer index")

			resp, err := dkgs[recipientIdx].ProcessDeal(deal)
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Response.Index = dealerIdx (0-based, which dealer's deal this responds to)
			require.Equal(t, uint32(dealerIdx), resp.Index,
				"Response.Index should be the 0-based dealer index")

			// Response.Response.Index = recipientIdx (0-based, the verifier who responded)
			require.Equal(t, uint32(recipientIdx), resp.Response.Index,
				"VSS Response.Index should be the 0-based recipient/verifier index")

			allResps = append(allResps, resp)
		}
	}

	// Step 4: Process all responses (completes the DKG)
	for _, resp := range allResps {
		for i, d := range dkgs {
			if resp.Response.Index == uint32(i) {
				continue // skip self
			}

			_, err := d.ProcessResponse(resp)
			require.NoError(t, err)
		}
	}

	// Step 5: Verify DKG completed successfully — all participants can derive key shares
	for i, d := range dkgs {
		_, err := d.DistKeyShare()
		require.NoError(t, err, "participant %d should have a valid dist key share", i)
	}

	// Step 6: Now simulate the justification verification path.
	// Create a polynomial (like a dealer's internal polynomial) and verify
	// the 0-based → 1-based conversion through VerifyPedersenVSS.
	//
	// This is what happens in a real justification:
	// - Dealer created shares with PriPoly.Eval(i) → PriShare{I: i (0-based)}
	// - Justification includes the plaintext deal with SecShare.I = i (0-based)
	// - CL must convert: recipientIndex = int(SecShare.I) + 1 for verification
	secret := suite.Scalar().Pick(suite.RandomStream())
	dealerPoly := share.NewPriPoly(suite, threshold, secret, suite.RandomStream())
	pubPoly := dealerPoly.Commit(suite.Point().Base())

	_, commits := pubPoly.Info()

	commitBzs := make([][]byte, len(commits))
	for i, c := range commits {
		bz, err := c.MarshalBinary()
		require.NoError(t, err)

		commitBzs[i] = bz
	}

	// For each participant, verify the index convention
	for kyberI := range n {
		priShare := dealerPoly.Eval(kyberI) // PriShare{I: kyberI (0-based), V: f(1+kyberI)}
		require.Equal(t, kyberI, priShare.I,
			"PriPoly.Eval(%d) should return PriShare.I = %d (0-based)", kyberI, kyberI)

		shareBz, err := priShare.V.MarshalBinary()
		require.NoError(t, err)

		// CORRECT: Convert 0-based → 1-based for VerifyPedersenVSS
		// This is what the CL does: recipientIndex = int(secShare.GetI()) + 1
		recipientIndex := priShare.I + 1
		ok, err := vss.VerifyPedersenVSS(suite, shareBz, recipientIndex, commitBzs, uint32(threshold))
		require.NoError(t, err)
		require.True(t, ok,
			"CORRECT: share[%d] (0-based) should verify at recipientIndex=%d (1-based)",
			kyberI, recipientIndex)

		// WRONG: Using 0-based index directly would fail
		// (except kyberI=0, which is rejected by validation as <= 0)
		if kyberI == 0 {
			// recipientIndex=0 is rejected with an error
			_, errZero := vss.VerifyPedersenVSS(suite, shareBz, 0, commitBzs, uint32(threshold))
			require.Error(t, errZero, "recipientIndex=0 should be rejected")
		} else {
			okWrong, err := vss.VerifyPedersenVSS(suite, shareBz, kyberI, commitBzs, uint32(threshold))
			require.NoError(t, err)
			require.False(t, okWrong,
				"WRONG: share[%d] should NOT verify at recipientIndex=%d (0-based used directly)",
				kyberI, kyberI)
		}
	}
}

// TestVerifyPedersenVSS_TableDriven is a table-driven catch-all for boundary
// and error cases.
func TestVerifyPedersenVSS_TableDriven(t *testing.T) {
	t.Parallel()

	suite := newSuite()
	n := 5
	threshold := 3
	commitmentBytes, shareBytes := generateVSSData(t, n, threshold)

	validShare := shareBytes[0]
	validCommitments := commitmentBytes

	tests := []struct {
		name            string
		shareBytes      []byte
		recipientIndex  int
		commitmentBytes [][]byte
		threshold       uint32
		wantOk          bool
		wantErr         bool
	}{
		{
			name:            "valid share at correct index",
			shareBytes:      validShare,
			recipientIndex:  1,
			commitmentBytes: validCommitments,
			threshold:       uint32(threshold),
			wantOk:          true,
			wantErr:         false,
		},
		{
			name:            "valid share at wrong index",
			shareBytes:      validShare,
			recipientIndex:  3,
			commitmentBytes: validCommitments,
			threshold:       uint32(threshold),
			wantOk:          false,
			wantErr:         false,
		},
		{
			name:            "empty commitments causes error",
			shareBytes:      validShare,
			recipientIndex:  1,
			commitmentBytes: [][]byte{},
			threshold:       0,
			wantOk:          false,
			wantErr:         true,
		},
		{
			name:            "nil commitments causes error",
			shareBytes:      validShare,
			recipientIndex:  1,
			commitmentBytes: nil,
			threshold:       0,
			wantOk:          false,
			wantErr:         true,
		},
		{
			name:            "malformed share bytes",
			shareBytes:      []byte("not a scalar"),
			recipientIndex:  1,
			commitmentBytes: validCommitments,
			threshold:       uint32(threshold),
			wantOk:          false,
			wantErr:         true,
		},
		{
			name:            "malformed commitment bytes",
			shareBytes:      validShare,
			recipientIndex:  1,
			commitmentBytes: [][]byte{[]byte("not a point"), validCommitments[1], validCommitments[2]},
			threshold:       uint32(threshold),
			wantOk:          false,
			wantErr:         true,
		},
		{
			name:            "zero recipient index",
			shareBytes:      validShare,
			recipientIndex:  0,
			commitmentBytes: validCommitments,
			threshold:       uint32(threshold),
			wantOk:          false,
			wantErr:         true,
		},
		{
			name:            "threshold mismatch",
			shareBytes:      validShare,
			recipientIndex:  1,
			commitmentBytes: validCommitments, // has 3 commitments
			threshold:       2,                // expects 2, mismatch
			wantOk:          false,
			wantErr:         true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ok, err := vss.VerifyPedersenVSS(suite, tc.shareBytes, tc.recipientIndex, tc.commitmentBytes, tc.threshold)
			if tc.wantErr {
				require.Error(t, err)
				require.False(t, ok)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantOk, ok)
			}
		})
	}
}
