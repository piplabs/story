package keeper

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"

	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/share"
	vssp "go.dedis.ch/kyber/v4/share/vss/pedersen"
	"go.dedis.ch/kyber/v4/sign/schnorr"
)

// dealingTestThreshold is the threshold used in dealing tests (t=2: linear polynomial).
const dealingTestThreshold = 2

// dealerTestContext holds all the data needed to construct and verify justifications in tests.
type dealerTestContext struct {
	suite      *edwards25519.SuiteEd25519
	longterm   kyber.Scalar // dealer's longterm private key
	pub        kyber.Point  // dealer's public key (dkgPubKey)
	pubBytes   []byte       // marshaled dkgPubKey for registration
	priPoly    *share.PriPoly
	commits    []kyber.Point
	shareBytes [][]byte // per-participant share bytes (0-indexed)
}

// newDealerTestContext generates a dealer key pair and VSS polynomial data for tests.
func newDealerTestContext(t *testing.T, n, threshold int) *dealerTestContext {
	t.Helper()

	suite := edwards25519.NewBlakeSHA256Ed25519()
	longterm := suite.Scalar().Pick(suite.RandomStream())
	pub := suite.Point().Mul(longterm, nil)

	pubBytes, err := pub.MarshalBinary()
	require.NoError(t, err)

	secret := suite.Scalar().Pick(suite.RandomStream())
	priPoly := share.NewPriPoly(suite, threshold, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.Point().Base())
	_, commits := pubPoly.Info()

	shares := priPoly.Shares(n)
	shareBytesList := make([][]byte, n)

	for i, s := range shares {
		bz, err := s.V.MarshalBinary()
		require.NoError(t, err)

		shareBytesList[i] = bz
	}

	return &dealerTestContext{
		suite:      suite,
		longterm:   longterm,
		pub:        pub,
		pubBytes:   pubBytes,
		priPoly:    priPoly,
		commits:    commits,
		shareBytes: shareBytesList,
	}
}

// makeSignedJustification constructs a fully signed justification proto using the dealer's
// longterm key, matching what kyber produces internally.
//
// kyberRecipientIndex is 0-based (kyber convention): the array index in the verifiers list.
// The proto SecShare.I and VSSJustification.Index store this 0-based value.
func (dtc *dealerTestContext) makeSignedJustification(t *testing.T, dealerIndex uint32, kyberRecipientIndex int) types.Justification {
	t.Helper()

	require.True(t, kyberRecipientIndex >= 0 && kyberRecipientIndex < len(dtc.shareBytes), "kyberRecipientIndex out of range")

	shareScalar := dtc.suite.Scalar()
	require.NoError(t, shareScalar.UnmarshalBinary(dtc.shareBytes[kyberRecipientIndex]))

	// Construct the kyber vss.Deal with 0-based index (matching kyber PriPoly.Eval(i))
	deal := &vssp.Deal{
		SessionID: []byte("test-session"),
		SecShare: &share.PriShare{
			I: kyberRecipientIndex, // 0-based (kyber convention)
			V: shareScalar,
		},
		T:           uint32(len(dtc.commits)),
		Commitments: dtc.commits,
	}

	// Construct the kyber vss.Justification with 0-based index
	kyberJust := &vssp.Justification{
		SessionID: []byte("test-session"),
		Index:     uint32(kyberRecipientIndex), // 0-based (verifier index)
		Deal:      deal,
	}

	// Sign using the dealer's longterm key (same as kyber's Dealer.ProcessResponse)
	hash := kyberJust.Hash(dtc.suite)
	sig, err := schnorr.Sign(dtc.suite, dtc.longterm, hash)
	require.NoError(t, err)

	// Convert to proto format (all indices remain 0-based)
	commitmentPoints := make([]*types.Point, len(dtc.commits))
	for i, c := range dtc.commits {
		bz, err := c.MarshalBinary()
		require.NoError(t, err)

		commitmentPoints[i] = &types.Point{Data: bz}
	}

	return types.Justification{
		Index: dealerIndex,
		VssJustification: &types.VSSJustification{
			SessionId: []byte("test-session"),
			Index:     uint32(kyberRecipientIndex), // 0-based
			PlainDeal: &types.PlainDeal{
				SessionId: []byte("test-session"),
				SecShare: &types.SecShare{
					I: uint32(kyberRecipientIndex), // 0-based (kyber convention)
					V: &types.Scalar{Data: dtc.shareBytes[kyberRecipientIndex]},
				},
				Threshold:   uint32(len(dtc.commits)),
				Commitments: commitmentPoints,
			},
			Signature: sig,
		},
	}
}

// makeInvalidDealJustification constructs a signed justification where the share does NOT match
// the commitments (simulating an invalid deal). The signature is still valid (from the dealer).
func (dtc *dealerTestContext) makeInvalidDealJustification(t *testing.T, dealerIndex uint32) types.Justification {
	t.Helper()

	// Use a completely different share that doesn't match dtc.commits
	wrongScalar := dtc.suite.Scalar().Pick(dtc.suite.RandomStream())
	wrongShareBytes, err := wrongScalar.MarshalBinary()
	require.NoError(t, err)

	kyberRecipientIndex := 0 // 0-based

	deal := &vssp.Deal{
		SessionID: []byte("invalid-session"),
		SecShare: &share.PriShare{
			I: kyberRecipientIndex, // 0-based
			V: wrongScalar,
		},
		T:           uint32(len(dtc.commits)),
		Commitments: dtc.commits,
	}

	kyberJust := &vssp.Justification{
		SessionID: []byte("invalid-session"),
		Index:     uint32(kyberRecipientIndex), // 0-based
		Deal:      deal,
	}

	hash := kyberJust.Hash(dtc.suite)
	sig, err := schnorr.Sign(dtc.suite, dtc.longterm, hash)
	require.NoError(t, err)

	commitmentPoints := make([]*types.Point, len(dtc.commits))
	for i, c := range dtc.commits {
		bz, err := c.MarshalBinary()
		require.NoError(t, err)

		commitmentPoints[i] = &types.Point{Data: bz}
	}

	return types.Justification{
		Index: dealerIndex,
		VssJustification: &types.VSSJustification{
			SessionId: []byte("invalid-session"),
			Index:     uint32(kyberRecipientIndex),
			PlainDeal: &types.PlainDeal{
				SessionId: []byte("invalid-session"),
				SecShare: &types.SecShare{
					I: uint32(kyberRecipientIndex), // 0-based
					V: &types.Scalar{Data: wrongShareBytes},
				},
				Threshold:   uint32(len(dtc.commits)),
				Commitments: commitmentPoints,
			},
			Signature: sig,
		},
	}
}

// setupDealerRegistrationWithKey creates a Verified DKG registration with the given
// dkgPubKey (as raw bytes from a real Edwards25519 point).
func setupDealerRegistrationWithKey(t *testing.T, k *Keeper, ctx context.Context, round uint32, dealerAddr common.Address, dealerIndex uint32, dkgPubKey []byte) {
	t.Helper()

	reg := &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: dealerAddr.Hex(),
		Index:         dealerIndex,
		DkgPubKey:     dkgPubKey,
		CommPubKey:    []byte("comm-pub"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, dealerAddr, reg))
}

// TestProcessJustifications_EmptyList verifies that ProcessJustifications returns
// nil immediately when given an empty justification list.
func TestProcessJustifications_EmptyList(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: dealingTestThreshold,
	}

	err := k.ProcessJustifications(ctx, network, []types.Justification{})
	require.NoError(t, err, "empty list should return nil error")
}

// TestProcessJustifications_NilList verifies that ProcessJustifications returns
// nil immediately when given a nil justification list.
func TestProcessJustifications_NilList(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: dealingTestThreshold,
	}

	err := k.ProcessJustifications(ctx, network, nil)
	require.NoError(t, err, "nil list should return nil error")
}

// TestVerifyJustificationSignature_Valid verifies that a correctly signed justification
// passes Schnorr signature verification.
func TestVerifyJustificationSignature_Valid(t *testing.T) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	dtc := newDealerTestContext(t, 3, dealingTestThreshold)
	dealerIndex := uint32(5)

	dealerPubKeys := map[uint32]kyber.Point{
		dealerIndex: dtc.pub,
	}

	j := dtc.makeSignedJustification(t, dealerIndex, 0) // 0-based recipient

	err := verifyJustificationSignature(suite, j, dealerPubKeys)
	require.NoError(t, err)
}

// TestVerifyJustificationSignature_BadSignature verifies that a justification
// signed by a different key than the dealer's registered dkgPubKey is rejected.
func TestVerifyJustificationSignature_BadSignature(t *testing.T) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	dealerIndex := uint32(3)

	// Register with one key
	realDealer := newDealerTestContext(t, 3, dealingTestThreshold)
	dealerPubKeys := map[uint32]kyber.Point{
		dealerIndex: realDealer.pub,
	}

	// Sign with a different key
	attacker := newDealerTestContext(t, 3, dealingTestThreshold)
	j := attacker.makeSignedJustification(t, dealerIndex, 0)

	err := verifyJustificationSignature(suite, j, dealerPubKeys)
	require.Error(t, err, "justification signed by wrong key should fail")
	require.Contains(t, err.Error(), "signature verification failed")
}

// TestVerifyJustificationSignature_DealerNotRegistered verifies that a justification
// from an unregistered dealer index is rejected.
func TestVerifyJustificationSignature_DealerNotRegistered(t *testing.T) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	dtc := newDealerTestContext(t, 3, dealingTestThreshold)

	// Empty map — no dealers registered
	dealerPubKeys := map[uint32]kyber.Point{}

	j := dtc.makeSignedJustification(t, 99, 0)

	err := verifyJustificationSignature(suite, j, dealerPubKeys)
	require.Error(t, err, "unregistered dealer should fail")
	require.Contains(t, err.Error(), "no registration found")
}

// TestVerifyJustificationSignature_NilVSSJustification verifies nil input handling.
func TestVerifyJustificationSignature_NilVSSJustification(t *testing.T) {
	suite := edwards25519.NewBlakeSHA256Ed25519()

	j := types.Justification{Index: 1, VssJustification: nil}

	err := verifyJustificationSignature(suite, j, map[uint32]kyber.Point{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil VSSJustification")
}

// TestVerifyJustificationSignature_EmptySignature verifies empty signature rejection.
func TestVerifyJustificationSignature_EmptySignature(t *testing.T) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	dtc := newDealerTestContext(t, 3, dealingTestThreshold)
	dealerIndex := uint32(1)

	dealerPubKeys := map[uint32]kyber.Point{
		dealerIndex: dtc.pub,
	}

	j := dtc.makeSignedJustification(t, dealerIndex, 0)
	j.VssJustification.Signature = nil

	err := verifyJustificationSignature(suite, j, dealerPubKeys)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty justification signature")
}

// TestVerifyJustification_ValidShare verifies that a correctly generated share
// passes VSS verification via the new standalone verifyJustification function.
func TestVerifyJustification_ValidShare(t *testing.T) {
	dtc := newDealerTestContext(t, 3, dealingTestThreshold)

	network := &types.DKGNetwork{
		Round:     1,
		Threshold: dealingTestThreshold,
	}

	// Test all 3 participants (0-based kyber indices)
	for i := range 3 {
		j := dtc.makeSignedJustification(t, 1, i)

		valid, err := verifyJustification(network, j)
		require.NoError(t, err, "recipient index %d should not error", i)
		require.True(t, valid, "valid share at kyber index %d should verify", i)
	}
}

// TestVerifyJustification_InvalidShare verifies that a wrong share fails VSS verification.
func TestVerifyJustification_InvalidShare(t *testing.T) {
	dtc := newDealerTestContext(t, 3, dealingTestThreshold)

	network := &types.DKGNetwork{
		Round:     1,
		Threshold: dealingTestThreshold,
	}

	j := dtc.makeInvalidDealJustification(t, 1)

	valid, err := verifyJustification(network, j)
	require.NoError(t, err, "invalid deal should not error, just return false")
	require.False(t, valid, "invalid VSS justification should return false")
}

// TestVerifyJustification_NilFields verifies nil field handling.
func TestVerifyJustification_NilFields(t *testing.T) {
	network := &types.DKGNetwork{
		Round:     1,
		Threshold: dealingTestThreshold,
	}

	tests := []struct {
		name string
		j    types.Justification
	}{
		{
			name: "nil VSSJustification",
			j:    types.Justification{Index: 1},
		},
		{
			name: "nil PlainDeal",
			j: types.Justification{
				Index: 1,
				VssJustification: &types.VSSJustification{
					PlainDeal: nil,
				},
			},
		},
		{
			name: "nil SecShare",
			j: types.Justification{
				Index: 1,
				VssJustification: &types.VSSJustification{
					PlainDeal: &types.PlainDeal{
						SecShare: nil,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := verifyJustification(network, tc.j)
			require.Error(t, err)
			require.False(t, valid)
		})
	}
}

// TestBuildDealerPubKeyMap verifies that the dealer public key map is built correctly.
func TestBuildDealerPubKeyMap(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(1)
	suite := edwards25519.NewBlakeSHA256Ed25519()

	dtc1 := newDealerTestContext(t, 3, 2)
	dtc2 := newDealerTestContext(t, 3, 2)

	dealer1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	dealer2 := common.HexToAddress("0x2222222222222222222222222222222222222222")

	setupDealerRegistrationWithKey(t, k, ctx, round, dealer1, 0, dtc1.pubBytes)
	setupDealerRegistrationWithKey(t, k, ctx, round, dealer2, 1, dtc2.pubBytes)

	network := &types.DKGNetwork{
		Round: round,
	}

	pubKeys, err := k.buildDealerPubKeyMap(ctx, network, suite)
	require.NoError(t, err)
	require.Len(t, pubKeys, 2)
	require.True(t, pubKeys[0].Equal(dtc1.pub))
	require.True(t, pubKeys[1].Equal(dtc2.pub))
}

// TestMaxJustificationsPerBlock verifies that the truncation slice operation
// correctly caps justifications to MaxJustificationsPerBlock, matching the
// logic used inside handleDKGProcessJustifications.
func TestMaxJustificationsPerBlock(t *testing.T) {
	const overCount = MaxJustificationsPerBlock + 50

	// Build a list that exceeds the cap. Each entry gets a distinct dealer index
	// so we can verify which ones survive after truncation.
	input := make([]types.Justification, overCount)
	for i := range input {
		input[i] = types.Justification{
			Index: uint32(i),
			VssJustification: &types.VSSJustification{
				PlainDeal: &types.PlainDeal{
					SecShare: &types.SecShare{I: 0},
				},
			},
		}
	}

	// Apply the same truncation logic used in handleDKGProcessJustifications.
	if len(input) > MaxJustificationsPerBlock {
		input = input[:MaxJustificationsPerBlock]
	}

	// After truncation the slice must be exactly MaxJustificationsPerBlock long.
	require.Len(t, input, MaxJustificationsPerBlock,
		"truncated slice must equal MaxJustificationsPerBlock")

	// The surviving entries must be the first MaxJustificationsPerBlock elements
	// (indices 0 .. MaxJustificationsPerBlock-1), not any of the later ones.
	for i, j := range input {
		require.Equal(t, uint32(i), j.Index,
			"entry at position %d should have dealer index %d", i, i)
	}
}

// TestJustificationPipeline_SignatureThenDedupThenVSS verifies the full
// signature → deduplication → VSS verification pipeline that mirrors the
// logic inside handleDKGProcessJustifications.
//
// Three sub-scenarios exercise every branch of the pipeline:
//  1. Valid-signature justifications pass through; invalid-signature ones are dropped.
//  2. Duplicate (dealerIndex, recipientIndex) pairs are deduplicated after signature
//     verification, so a valid duplicate does not produce two results.
//  3. A justification that carries a valid signature but an incorrect share fails
//     VSS verification and is dropped from the final list.
func TestJustificationPipeline_SignatureThenDedupThenVSS(t *testing.T) {
	// n=3 participants, threshold=2 (linear polynomial → 2 commitments).
	const n = 3

	k, ctx := setupDKGKeeper(t)

	round := uint32(1)

	// Dealer 0: valid key pair and VSS data.
	dealer0 := newDealerTestContext(t, n, dealingTestThreshold)
	addr0 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	setupDealerRegistrationWithKey(t, k, ctx, round, addr0, 0, dealer0.pubBytes)

	// Dealer 1: valid key pair and VSS data.
	dealer1 := newDealerTestContext(t, n, dealingTestThreshold)
	addr1 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	setupDealerRegistrationWithKey(t, k, ctx, round, addr1, 1, dealer1.pubBytes)

	// Attacker: a different key pair — its justifications carry bad signatures.
	attacker := newDealerTestContext(t, n, dealingTestThreshold)

	suite := edwards25519.NewBlakeSHA256Ed25519()

	network := &types.DKGNetwork{
		Round:     round,
		Threshold: dealingTestThreshold,
	}

	// Build the dealer public key map once, matching what the handler builds.
	dealerPubKeys, err := k.buildDealerPubKeyMap(ctx, network, suite)
	require.NoError(t, err)
	require.Len(t, dealerPubKeys, 2)

	t.Run("valid signatures pass, invalid signatures are dropped", func(t *testing.T) {
		// dealer0 signs correctly; attacker signs with the wrong key for dealer index 1.
		validJ := dealer0.makeSignedJustification(t, 0, 0)
		invalidSigJ := attacker.makeSignedJustification(t, 1, 0) // wrong key for dealer index 1

		input := []types.Justification{validJ, invalidSigJ}

		var sigVerified []types.Justification

		for _, j := range input {
			if err := verifyJustificationSignature(suite, j, dealerPubKeys); err == nil {
				sigVerified = append(sigVerified, j)
			}
		}

		// Only validJ must survive signature checking.
		require.Len(t, sigVerified, 1,
			"exactly one justification should pass signature verification")
		require.Equal(t, uint32(0), sigVerified[0].Index,
			"surviving justification must be from dealer 0")
	})

	t.Run("duplicate justifications are deduplicated after signature check", func(t *testing.T) {
		// Two valid justifications from dealer0, same recipient → same (dealer,recipient) key.
		j1 := dealer0.makeSignedJustification(t, 0, 0)
		j2 := dealer0.makeSignedJustification(t, 0, 0) // exact duplicate

		input := []types.Justification{j1, j2}

		// Step 1: signature check — both pass because dealer0's key is correct.
		var sigVerified []types.Justification

		for _, j := range input {
			if err := verifyJustificationSignature(suite, j, dealerPubKeys); err == nil {
				sigVerified = append(sigVerified, j)
			}
		}

		require.Len(t, sigVerified, 2, "both should pass signature check before dedup")

		// Step 2: deduplicate — (dealerIndex=0, recipientIndex=0) appears twice.
		deduped := deduplicateJustifications(sigVerified)
		require.Len(t, deduped, 1, "deduplication must collapse the two identical entries to one")
		require.Equal(t, uint32(0), deduped[0].Index)
	})

	t.Run("VSS-invalid justification is dropped after sig check and dedup", func(t *testing.T) {
		// dealer1 provides a justification with a valid signature but an incorrect share.
		invalidVSSJ := dealer1.makeInvalidDealJustification(t, 1)

		// Step 1: signature must pass (dealer1 signed it with its own valid key).
		err := verifyJustificationSignature(suite, invalidVSSJ, dealerPubKeys)
		require.NoError(t, err, "signature from dealer1 should verify against its registered key")

		// Step 2: deduplication is a no-op for a single entry.
		deduped := deduplicateJustifications([]types.Justification{invalidVSSJ})
		require.Len(t, deduped, 1)

		// Step 3: VSS verification must fail for the bad share.
		valid, err := verifyJustification(network, deduped[0])
		require.NoError(t, err, "VSS verification should not error, just return false")
		require.False(t, valid, "justification with invalid share must fail VSS verification")
	})

	t.Run("mixed pipeline: only VSS-valid justifications reach the final list", func(t *testing.T) {
		// Scenario: 3 justifications arrive in a single block.
		//   j_valid      – dealer0, recipient 0, valid signature and valid share → survives
		//   j_bad_sig    – attacker pretending to be dealer1 → dropped at sig check
		//   j_bad_vss    – dealer1, valid signature but wrong share → dropped at VSS
		jValid := dealer0.makeSignedJustification(t, 0, 0)
		jBadSig := attacker.makeSignedJustification(t, 1, 0)  // wrong key for dealer 1
		jBadVSS := dealer1.makeInvalidDealJustification(t, 1) // dealer1, bad share

		input := []types.Justification{jValid, jBadSig, jBadVSS}

		// Apply the same three-step pipeline used by handleDKGProcessJustifications.

		// Step 1: filter by Schnorr signature.
		var sigVerified []types.Justification

		for _, j := range input {
			if err := verifyJustificationSignature(suite, j, dealerPubKeys); err == nil {
				sigVerified = append(sigVerified, j)
			}
		}
		// jBadSig is dropped; jValid and jBadVSS remain.
		require.Len(t, sigVerified, 2, "only sig-valid justifications should pass step 1")

		// Step 2: deduplicate (no duplicates in this scenario).
		deduped := deduplicateJustifications(sigVerified)
		require.Len(t, deduped, 2, "no duplicates, count should be unchanged after dedup")

		// Step 3: filter by Pedersen VSS.
		var finalValid []types.Justification

		for _, j := range deduped {
			ok, err := verifyJustification(network, j)
			if err == nil && ok {
				finalValid = append(finalValid, j)
			}
		}

		// Only jValid (dealer0, valid share) must remain.
		require.Len(t, finalValid, 1, "only the VSS-valid justification should survive")
		require.Equal(t, uint32(0), finalValid[0].Index,
			"surviving justification must come from dealer 0")
	})
}
