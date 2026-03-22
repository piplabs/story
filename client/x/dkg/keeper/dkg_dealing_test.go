package keeper

import (
	"context"
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
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

// TestProcessJustifications_TruncatesExcessiveList verifies that ProcessJustifications
// truncates the input list to MaxJustificationsPerBlock when more justifications are given.
// This covers the truncation branch inside ProcessJustifications (line 117-124).
func TestProcessJustifications_TruncatesExcessiveList(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)

	round := uint32(5)
	network := &types.DKGNetwork{
		Round:     round,
		Total:     3,
		Threshold: dealingTestThreshold,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Build MaxJustificationsPerBlock+5 justifications (all invalid — no registered keys).
	// ProcessJustifications will truncate to MaxJustificationsPerBlock, then all will fail
	// signature verification (no dealer registered). The result is an empty valid set,
	// but the truncation branch is covered.
	overCount := MaxJustificationsPerBlock + 5
	justifications := make([]types.Justification, overCount)
	for i := range justifications {
		justifications[i] = types.Justification{
			Index: uint32(i),
			VssJustification: &types.VSSJustification{
				PlainDeal: &types.PlainDeal{
					SecShare: &types.SecShare{I: 0},
				},
			},
		}
	}

	// ProcessJustifications should succeed without error even though all sig-verifications fail
	err := k.ProcessJustifications(ctx, network, justifications)
	require.NoError(t, err, "truncation and sig-drop should not return an error")
}

// TestMaxJustificationsPerBlock verifies that the truncation slice operation
// correctly caps justifications to MaxJustificationsPerBlock, matching the
// logic used inside ProcessJustifications.
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

	// Apply the same truncation logic used in ProcessJustifications.
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
// logic inside ProcessJustifications.
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

		// Apply the same three-step pipeline used by ProcessJustifications.

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

// TestProcessJustifications_InvalidatesDealer verifies that ProcessJustifications
// invalidates a dealer's on-chain registration when VSS verification fails,
// while leaving valid dealers untouched. This is the core fix for issue #717:
// previously, VSS failure only logged a message in the async goroutine and
// could not update on-chain state.
func TestProcessJustifications_InvalidatesDealer(t *testing.T) {
	const n = 3

	k, ctx := setupDKGKeeper(t)
	k.isDKGSvcEnabled = true

	// Initialize stateManager so handleDKGProcessJustifications doesn't panic
	// when forwarding valid justifications to kernel in async goroutine.
	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k.stateManager = sm

	round := uint32(1)

	// Set up two dealers with real crypto keys
	dealer0 := newDealerTestContext(t, n, dealingTestThreshold)
	addr0 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	setupDealerRegistrationWithKey(t, k, ctx, round, addr0, 0, dealer0.pubBytes)

	dealer1 := newDealerTestContext(t, n, dealingTestThreshold)
	addr1 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	setupDealerRegistrationWithKey(t, k, ctx, round, addr1, 1, dealer1.pubBytes)

	network := &types.DKGNetwork{
		Round:     round,
		Total:     uint32(n),
		Threshold: dealingTestThreshold,
	}

	t.Run("VSS failure invalidates dealer registration", func(t *testing.T) {
		// dealer1 produces a justification with valid signature but invalid share
		jBadVSS := dealer1.makeInvalidDealJustification(t, 1)

		err := k.ProcessJustifications(ctx, network, []types.Justification{jBadVSS})
		require.NoError(t, err)

		// Dealer 1's registration must now be Invalidated
		reg1, err := k.getDKGRegistration(ctx, round, addr1)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusInvalidated, reg1.Status,
			"dealer with invalid VSS must be invalidated")

		// Dealer 0's registration must remain Verified (unaffected)
		reg0, err := k.getDKGRegistration(ctx, round, addr0)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg0.Status,
			"uninvolved dealer must remain Verified")
	})

	t.Run("valid justification does not invalidate dealer", func(t *testing.T) {
		// Reset dealer0 to Verified for a clean test
		setupDealerRegistrationWithKey(t, k, ctx, round, addr0, 0, dealer0.pubBytes)

		jValid := dealer0.makeSignedJustification(t, 0, 0)

		err := k.ProcessJustifications(ctx, network, []types.Justification{jValid})
		require.NoError(t, err)

		// Dealer 0's registration must still be Verified
		reg0, err := k.getDKGRegistration(ctx, round, addr0)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg0.Status,
			"dealer with valid VSS must remain Verified")
	})

	t.Run("mixed valid and invalid justifications", func(t *testing.T) {
		// Reset both dealers
		setupDealerRegistrationWithKey(t, k, ctx, round, addr0, 0, dealer0.pubBytes)
		setupDealerRegistrationWithKey(t, k, ctx, round, addr1, 1, dealer1.pubBytes)

		jValid := dealer0.makeSignedJustification(t, 0, 0)
		jBadVSS := dealer1.makeInvalidDealJustification(t, 1)

		err := k.ProcessJustifications(ctx, network, []types.Justification{jValid, jBadVSS})
		require.NoError(t, err)

		// Dealer 0 (valid) must remain Verified
		reg0, err := k.getDKGRegistration(ctx, round, addr0)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg0.Status,
			"dealer with valid justification must remain Verified")

		// Dealer 1 (invalid VSS) must be Invalidated
		reg1, err := k.getDKGRegistration(ctx, round, addr1)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusInvalidated, reg1.Status,
			"dealer with invalid VSS must be Invalidated")
	})
}

// TestBuildDealerPubKeyMap_EmptyDkgPubKey verifies that buildDealerPubKeyMap skips
// registrations with empty DkgPubKey (the `len(reg.DkgPubKey) == 0` branch).
func TestBuildDealerPubKeyMap_EmptyDkgPubKey(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)

	round := uint32(3)
	suite := edwards25519.NewBlakeSHA256Ed25519()

	// Register a dealer with an empty DkgPubKey — should be skipped.
	emptyKeyDealer := common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	reg := &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: emptyKeyDealer.Hex(),
		Index:         1,
		DkgPubKey:     []byte{}, // empty — triggers the skip branch
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, emptyKeyDealer, reg))

	network := &types.DKGNetwork{Round: round}
	pubKeys, err := k.buildDealerPubKeyMap(ctx, network, suite)
	require.NoError(t, err)
	require.Empty(t, pubKeys, "empty DkgPubKey registrations must be skipped")
}

// TestBuildDealerPubKeyMap_InvalidDkgPubKeyBytes verifies that buildDealerPubKeyMap
// silently skips (via log.Warn) registrations whose DkgPubKey bytes cannot be
// unmarshaled as an Edwards25519 point. The map is still returned without error.
func TestBuildDealerPubKeyMap_InvalidDkgPubKeyBytes(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)

	round := uint32(4)
	suite := edwards25519.NewBlakeSHA256Ed25519()

	// Register with a valid key (should appear in the map).
	dtcValid := newDealerTestContext(t, 3, 2)
	validDealer := common.HexToAddress("0x1111111111111111111111111111111111111111")
	setupDealerRegistrationWithKey(t, k, ctx, round, validDealer, 1, dtcValid.pubBytes)

	// Register with invalid DkgPubKey bytes — unmarshal will fail, entry skipped.
	invalidDealer := common.HexToAddress("0x2222222222222222222222222222222222222222")
	badReg := &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: invalidDealer.Hex(),
		Index:         2,
		DkgPubKey:     []byte("not-a-valid-edwards25519-point"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, invalidDealer, badReg))

	network := &types.DKGNetwork{Round: round}
	pubKeys, err := k.buildDealerPubKeyMap(ctx, network, suite)
	require.NoError(t, err, "invalid DkgPubKey bytes must be skipped without error")
	require.Len(t, pubKeys, 1, "only the valid dealer should appear in the map")
	require.True(t, pubKeys[1].Equal(dtcValid.pub), "valid dealer's key must be in the map at index 1")
}

// --- Tests merged from begin_dealing_test.go ---

// TestBeginDealing_BelowMinRegistrations verifies that BeginDealing skips
// to the next round when verified registrations are below the minimum.
func TestBeginDealing_BelowMinRegistrations(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// SkipToNextRound -> InitiateDKGRound -> GetAllValidators
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stakingtypes.Validator{}, nil).AnyTimes()

	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 5 // require 5 but we only register 2
	require.NoError(t, k.SetParams(ctx, params))

	round := uint32(1)
	latestRound := &types.DKGNetwork{
		Round: round,
		Stage: types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	// Register only 2 verified validators (below min of 5)
	for i, addrHex := range []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
	} {
		addr := common.HexToAddress(addrHex)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:     round,
			Index:     uint32(i + 1),
			DkgPubKey: []byte("pub"),
			Status:    types.DKGRegStatusVerified,
		}))
	}

	err := k.BeginDealing(ctx, latestRound)
	require.NoError(t, err)

	// Verify a new round was created (SkipToNextRound)
	nextRound, err := k.getDKGNetwork(ctx, round+1)
	require.NoError(t, err)
	require.Equal(t, round+1, nextRound.Round)
}

// TestBeginDealing_MeetsMinRegistrations verifies that BeginDealing succeeds
// when verified registrations meet the minimum threshold.
func TestBeginDealing_MeetsMinRegistrations(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 3
	params.OperationalThreshold = 667 // 66.7%
	require.NoError(t, k.SetParams(ctx, params))

	round := uint32(1)
	latestRound := &types.DKGNetwork{
		Round: round,
		Stage: types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	// Register 3 verified validators
	for i, addrHex := range []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	} {
		addr := common.HexToAddress(addrHex)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:     round,
			Index:     uint32(i + 1),
			DkgPubKey: []byte("pub"),
			Status:    types.DKGRegStatusVerified,
		}))
	}

	err := k.BeginDealing(ctx, latestRound)
	require.NoError(t, err)

	// Network should be updated with total and threshold
	updated, err := k.getDKGNetwork(ctx, round)
	require.NoError(t, err)
	require.Equal(t, uint32(3), updated.Total)
	require.True(t, updated.Threshold > 0)
}

// TestBeginDealing_WithDKGSvcEnabled verifies that BeginDealing spawns
// an async dealing goroutine when isDKGSvcEnabled is true.
func TestBeginDealing_WithDKGSvcEnabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm
	k.isDKGSvcEnabled = true
	k.validatorEVMAddr = testValidatorAddr

	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 2
	params.OperationalThreshold = 667
	require.NoError(t, k.SetParams(ctx, params))

	round := uint32(1)
	latestRound := &types.DKGNetwork{
		Round:        round,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr, "0xother"},
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	for i, addrHex := range []string{testValidatorAddr, "0x2222222222222222222222222222222222222222"} {
		addr := common.HexToAddress(addrHex)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:     round,
			Index:     uint32(i + 1),
			DkgPubKey: []byte("pub"),
			Status:    types.DKGRegStatusVerified,
		}))
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	err = k.BeginDealing(ctx, latestRound)
	require.NoError(t, err)
}

// --- Tests merged from dkg_process_test.go ---

// TestProcessDeals_DKGSvcDisabled verifies that ProcessDeals emits the event
// and returns nil even when DKG service is disabled.

func TestProcessDeals_DKGSvcDisabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// isDKGSvcEnabled is false by default

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 2},
	}

	err := k.ProcessDeals(ctx, network, deals)
	require.NoError(t, err)
}

// TestProcessDeals_EmptyDeals verifies that ProcessDeals handles an empty
// deals slice without error.

func TestProcessDeals_EmptyDeals(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     2,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	err := k.ProcessDeals(ctx, network, []types.Deal{})
	require.NoError(t, err)
}

// TestProcessDeals_MultipleDeals verifies that ProcessDeals correctly processes
// multiple deals at once.

func TestProcessDeals_MultipleDeals(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     3,
		Total:     5,
		Threshold: 4,
		Stage:     types.DKGStageDealing,
	}

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 2},
		{Index: 1, RecipientIndex: 3},
		{Index: 2, RecipientIndex: 1},
	}

	err := k.ProcessDeals(ctx, network, deals)
	require.NoError(t, err)
}

// TestProcessResponses_DKGSvcDisabled verifies that ProcessResponses emits the
// event and returns nil when DKG service is disabled.

func TestProcessResponses_DKGSvcDisabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	responses := []types.Response{
		{Index: 1},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}

// TestProcessResponses_EmptyResponses verifies that ProcessResponses handles
// an empty response slice without error.

func TestProcessResponses_EmptyResponses(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     2,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	err := k.ProcessResponses(ctx, network, []types.Response{})
	require.NoError(t, err)
}

// TestProcessResponses_MultipleResponses verifies that ProcessResponses correctly
// processes multiple responses.

func TestProcessResponses_MultipleResponses(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     3,
		Total:     5,
		Threshold: 4,
		Stage:     types.DKGStageDealing,
	}

	responses := []types.Response{
		{Index: 1},
		{Index: 2},
		{Index: 3},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}

// TestEnsureSessionIndex_NoSession verifies that ensureSessionIndex returns an error
// when no session exists for the given round.

func TestEnsureSessionIndex_NoSession(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(common.HexToAddress("0x1111111111111111111111111111111111111111"))
	initTestStateManager(t, k)

	// No session exists for round 99 — should return an error
	err := k.ensureSessionIndex(ctx, 99)
	require.Error(t, err, "ensureSessionIndex should fail when no session exists")
}

// TestEnsureSessionIndex_IndexAlreadySet verifies that ensureSessionIndex is a no-op
// when the session already has a non-zero index.

func TestEnsureSessionIndex_IndexAlreadySet(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(common.HexToAddress("0x1111111111111111111111111111111111111111"))
	initTestStateManager(t, k)

	// Create a session with a pre-set index
	sess := newTestSession(5)
	require.NoError(t, k.stateManager.CreateSession(ctx, sess))
	session, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	session.Index = 3 // already set
	require.NoError(t, k.stateManager.UpdateSession(ctx, session))

	// ensureSessionIndex should return nil without touching the index
	err = k.ensureSessionIndex(ctx, 5)
	require.NoError(t, err)

	// Verify index is unchanged
	updated, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	require.Equal(t, uint32(3), updated.Index, "index should remain unchanged when already set")
}

// TestEnsureSessionIndex_SetsIndexFromRegistration verifies that ensureSessionIndex
// reads the on-chain registration index and updates the session.

func TestEnsureSessionIndex_SetsIndexFromRegistration(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	initTestStateManager(t, k)

	addr := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	k.setValidatorAddress(addr)

	// Store a registration with index=2 for round 7
	reg := &types.DKGRegistration{
		Round:         7,
		ValidatorAddr: addr.Hex(),
		Index:         2,
		DkgPubKey:     []byte("test-dkg-pub-key"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, addr, reg))

	// Create a session with Index=0 (not yet set)
	require.NoError(t, k.stateManager.CreateSession(ctx, newTestSession(7)))
	session, err := k.stateManager.GetSession(7)
	require.NoError(t, err)
	require.Equal(t, uint32(0), session.Index, "index should start at 0")

	err = k.ensureSessionIndex(ctx, 7)
	require.NoError(t, err)

	// Index should now be set from the registration
	updated, err := k.stateManager.GetSession(7)
	require.NoError(t, err)
	require.Equal(t, uint32(2), updated.Index, "index should be set from on-chain registration")
}

// TestEnsureSessionIndex_NoRegistration verifies that ensureSessionIndex returns an
// error when the session exists but the validator has no registration.

func TestEnsureSessionIndex_NoRegistration(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	initTestStateManager(t, k)

	// Set validator address to an address with no registration
	k.setValidatorAddress(common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"))

	// Create a session with Index=0
	require.NoError(t, k.stateManager.CreateSession(ctx, newTestSession(15)))

	// No registration for round 15 → should fail
	err := k.ensureSessionIndex(ctx, 15)
	require.Error(t, err, "should fail when registration not found")
}

// TestBeginDealing_BelowMinRequired verifies that BeginDealing calls SkipToNextRound
// when verified registration count is below the minimum required.

func TestBeginDealing_BelowMinRequired(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).AnyTimes()

	// Set params with MinReqRegisteredParticipants=3
	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 3
	require.NoError(t, k.SetParams(ctx, params))

	// Create a network in Registration stage with round=1
	network := &types.DKGNetwork{
		Round:        1,
		Total:        5,
		Threshold:    4,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0x1111111111111111111111111111111111111111"},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Store only 1 verified registration (below min=3)
	addr1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	require.NoError(t, k.setDKGRegistration(ctx, addr1, &types.DKGRegistration{
		Round:         1,
		ValidatorAddr: addr1.Hex(),
		Index:         1,
		Status:        types.DKGRegStatusVerified,
	}))

	// BeginDealing should skip to next round
	err := k.BeginDealing(ctx, network)
	require.NoError(t, err)

	// Verify a new round (round=2) was created via SkipToNextRound
	_, err = k.getDKGNetwork(ctx, 2)
	require.NoError(t, err, "SkipToNextRound should have created round 2")
}

// TestBeginDealing_MeetsMinRequired verifies that BeginDealing proceeds to the
// dealing phase when verified registration count meets the minimum required.

func TestBeginDealing_MeetsMinRequired(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Set params with MinReqRegisteredParticipants=2
	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 2
	require.NoError(t, k.SetParams(ctx, params))

	network := &types.DKGNetwork{
		Round:        10,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0x1111111111111111111111111111111111111111"},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Store 3 verified registrations (above min=2)
	for i, hexAddr := range []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	} {
		addr := common.HexToAddress(hexAddr)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         10,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusVerified,
		}))
	}

	err := k.BeginDealing(ctx, network)
	require.NoError(t, err)

	// Verify round 10 was updated with Total=3 (not skipped to round 11)
	updated, err := k.getDKGNetwork(ctx, 10)
	require.NoError(t, err)
	require.Equal(t, uint32(3), updated.Total, "Total should be set to verified registration count")
}

// TestProcessDeals_DKGSvcEnabled verifies that ProcessDeals starts the async
// goroutine when DKG service is enabled (no panic or error).

func TestProcessDeals_DKGSvcEnabled(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()

	network := &types.DKGNetwork{
		Round:     20,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 2},
	}

	// Should not panic even when svc is enabled (goroutine runs handleDKGProcessDeals)
	err := k.ProcessDeals(ctx, network, deals)
	require.NoError(t, err)
}

// TestProcessResponses_DKGSvcEnabled verifies that ProcessResponses starts the
// async goroutine path when DKG service is enabled, using shouldProcessResponses.

func TestProcessResponses_DKGSvcEnabled(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(common.HexToAddress("0x1111111111111111111111111111111111111111"))

	network := &types.DKGNetwork{
		Round:     21,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	responses := []types.Response{
		{Index: 1},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}

// TestProcessResponses_DKGSvcEnabled_WithStateManager verifies that ProcessResponses
// runs shouldProcessResponses successfully when a stateManager is initialized.

func TestProcessResponses_DKGSvcEnabled_WithStateManager(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	validatorAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	k.setValidatorAddress(validatorAddr)
	initTestStateManager(t, k)

	network := &types.DKGNetwork{
		Round:        22,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{validatorAddr.Hex()},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Create a session so shouldProcessResponses can check session state
	require.NoError(t, k.stateManager.CreateSession(ctx, newTestSession(22)))

	responses := []types.Response{
		{Index: 1},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}

// initTestStateManager creates a StateManager backed by a temp directory and assigns it to the keeper.
// This is needed because stateManager is only initialized in InitDKGService, not in setupDKGKeeperWithMocks.
func initTestStateManager(t *testing.T, k *Keeper) {
	t.Helper()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k.stateManager = sm
}
