package keeper

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"

	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/share"
)

// generateShareAndCommitments generates a valid Pedersen VSS share and commitment
// bytes for a given recipientIndex (1-based) using the Edwards25519 suite.
// The underlying polynomial uses threshold=2 (linear), producing 2 commitments.
func generateShareAndCommitments(t *testing.T, recipientIndex int) (shareBytes []byte, commitmentBytes [][]byte) {
	t.Helper()

	suite := edwards25519.NewBlakeSHA256Ed25519()
	n := 3
	threshold := int(testThreshold)

	secret := suite.Scalar().Pick(suite.RandomStream())
	priPoly := share.NewPriPoly(suite, threshold, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.Point().Base())

	shares := priPoly.Shares(n)
	// Info() returns (base kyber.Point, commits []kyber.Point)
	_, commits := pubPoly.Info()

	idx := recipientIndex - 1
	if idx < 0 || idx >= n {
		idx = 0
	}

	var err error

	shareBytes, err = shares[idx].V.MarshalBinary()
	require.NoError(t, err)

	commitmentBytes = make([][]byte, len(commits))
	for i, c := range commits {
		bz, err := c.MarshalBinary()
		require.NoError(t, err)

		commitmentBytes[i] = bz
	}

	return shareBytes, commitmentBytes
}

// threshold is the number of commitments in these tests.
const testThreshold = uint32(2)

// TestVerifyJustificationVSS_ValidShare verifies that a correctly generated share
// passes the VSS wrapper function.
func TestVerifyJustificationVSS_ValidShare(t *testing.T) {
	t.Parallel()

	shareBytes, commitmentBytes := generateShareAndCommitments(t, 1)

	ok, err := verifyJustificationVSS(shareBytes, 1, commitmentBytes, testThreshold)
	require.NoError(t, err)
	require.True(t, ok, "valid share should pass VSS verification")
}

// TestVerifyJustificationVSS_InvalidShare verifies that a random scalar fails
// verification against commitments it was not derived from.
func TestVerifyJustificationVSS_InvalidShare(t *testing.T) {
	t.Parallel()

	suite := edwards25519.NewBlakeSHA256Ed25519()
	_, commitmentBytes := generateShareAndCommitments(t, 1)

	// Use a completely different random scalar
	wrongScalar := suite.Scalar().Pick(suite.RandomStream())
	wrongBytes, err := wrongScalar.MarshalBinary()
	require.NoError(t, err)

	ok, err := verifyJustificationVSS(wrongBytes, 1, commitmentBytes, testThreshold)
	require.NoError(t, err)
	require.False(t, ok, "wrong scalar should fail VSS verification")
}

// TestVerifyJustificationVSS_EmptyCommitments verifies that empty commitments
// cause verifyJustificationVSS to return an error.
func TestVerifyJustificationVSS_EmptyCommitments(t *testing.T) {
	t.Parallel()

	shareBytes, _ := generateShareAndCommitments(t, 1)

	ok, err := verifyJustificationVSS(shareBytes, 1, [][]byte{}, 0)
	require.Error(t, err, "empty commitments should return an error")
	require.False(t, ok)
}

// TestVerifyJustificationVSS_MalformedShareBytes verifies that malformed share bytes
// produce an error.
func TestVerifyJustificationVSS_MalformedShareBytes(t *testing.T) {
	t.Parallel()

	_, commitmentBytes := generateShareAndCommitments(t, 1)

	ok, err := verifyJustificationVSS([]byte("not-a-scalar"), 1, commitmentBytes, testThreshold)
	require.Error(t, err, "malformed share bytes should return an error")
	require.False(t, ok)
}

// TestVerifyJustificationVSS_MalformedCommitmentBytes verifies that malformed
// commitment bytes produce an error.
func TestVerifyJustificationVSS_MalformedCommitmentBytes(t *testing.T) {
	t.Parallel()

	shareBytes, _ := generateShareAndCommitments(t, 1)

	// threshold=1 but we pass a malformed single commitment
	ok, err := verifyJustificationVSS(shareBytes, 1, [][]byte{[]byte("not-a-point")}, 1)
	require.Error(t, err, "malformed commitment bytes should return an error")
	require.False(t, ok)
}

// TestVerifyJustificationVSS_WrongIndex verifies that a valid share fails
// when checked against a different recipient index.
func TestVerifyJustificationVSS_WrongIndex(t *testing.T) {
	t.Parallel()

	// Generate share for index 1
	shareBytes, commitmentBytes := generateShareAndCommitments(t, 1)

	// Check at index 2 - should fail
	ok, err := verifyJustificationVSS(shareBytes, 2, commitmentBytes, testThreshold)
	require.NoError(t, err)
	require.False(t, ok, "share for index 1 should not verify at index 2")
}

// TestInvalidateDealerByAddr_Success verifies a Verified dealer is transitioned to Invalidated.
func TestInvalidateDealerByAddr_Success(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(5)
	dealerAddr := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerAddr, 3, types.DKGRegStatusVerified)

	require.NoError(t, k.invalidateDealerByAddr(ctx, round, dealerAddr))

	reg, err := k.getDKGRegistration(ctx, round, dealerAddr)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusInvalidated, reg.Status)
}

// TestInvalidateDealerByAddr_AlreadyFinalized verifies a finalized dealer is not invalidated.
func TestInvalidateDealerByAddr_AlreadyFinalized(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(1)
	dealerAddr := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerAddr, 2, types.DKGRegStatusFinalized)

	err := k.invalidateDealerByAddr(ctx, round, dealerAddr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already finalized")
	require.Contains(t, err.Error(), "possible bug")

	reg, err := k.getDKGRegistration(ctx, round, dealerAddr)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusFinalized, reg.Status)
}

// TestInvalidateDealerByAddr_NoRegistration verifies that invalidating an address with no
// registration in the round is a no-op (a dealer that left the committee).
func TestInvalidateDealerByAddr_NoRegistration(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	missing := common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	require.NoError(t, k.invalidateDealerByAddr(ctx, 1, missing))
}

// TestInvalidateDealerByAddr_MultipleRegistrations verifies only the addressed dealer is hit.
func TestInvalidateDealerByAddr_MultipleRegistrations(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(3)
	dealerA := common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
	dealerB := common.HexToAddress("0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")

	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerA, 1, types.DKGRegStatusVerified)
	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerB, 2, types.DKGRegStatusVerified)

	require.NoError(t, k.invalidateDealerByAddr(ctx, round, dealerA))

	regA, err := k.getDKGRegistration(ctx, round, dealerA)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusInvalidated, regA.Status)

	regB, err := k.getDKGRegistration(ctx, round, dealerB)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusVerified, regB.Status)
}

// TestInvalidateDealerByAddr_AlreadyInvalidated verifies idempotency.
func TestInvalidateDealerByAddr_AlreadyInvalidated(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(2)
	dealerAddr := common.HexToAddress("0xEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerAddr, 5, types.DKGRegStatusInvalidated)

	require.NoError(t, k.invalidateDealerByAddr(ctx, round, dealerAddr))

	reg, err := k.getDKGRegistration(ctx, round, dealerAddr)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusInvalidated, reg.Status)
}

// setupDealerRegistrationForInvalidation is a helper that creates a registration
// with a specific status for use in invalidation tests.
func setupDealerRegistrationForInvalidation(t *testing.T, k *Keeper, ctx context.Context, round uint32, addr common.Address, index uint32, status types.DKGRegStatus) {
	t.Helper()

	reg := &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: addr.Hex(),
		Index:         index,
		DkgPubKey:     []byte("dkg-pub"),
		CommPubKey:    []byte("comm-pub"),
		Status:        status,
	}
	require.NoError(t, k.setDKGRegistration(ctx, addr, reg))
}
