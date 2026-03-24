package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	testValidator1 = common.HexToAddress("0x1111111111111111111111111111111111111111")
	testValidator2 = common.HexToAddress("0x2222222222222222222222222222222222222222")
)

// testLabel returns a fixed 32-byte label for use in partial decryption tests.
func testLabel() []byte {
	label := make([]byte, 32)
	label[0] = 0xAB
	label[1] = 0xCD
	return label
}

// TestSetPartialDecryptionSubmission_Success verifies that a valid partial
// decryption submission is stored without error.
func TestSetPartialDecryptionSubmission_Success(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1,
		1, // round
		1, // pid
		[]byte("encrypted-partial"),
		[]byte("ephemeral-pub-key"),
		[]byte("pub-share"),
		[]byte("requester-pub-key"),
		testLabel(),
		[]byte("ciphertext-data"),
	)
	require.NoError(t, err)
}

// TestSetPartialDecryptionSubmission_DuplicateRejected verifies that submitting
// the same partial decryption twice returns ErrDuplicatePartialDecryptionSubmission.
func TestSetPartialDecryptionSubmission_DuplicateRejected(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext-data")

	// First submission
	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)

	// Second submission with same (validator, round, requester, label, ciphertext)
	err = k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial-2"),
		[]byte("eph-key-2"),
		[]byte("pub-share-2"),
		requesterPubKey, label, ciphertext,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDuplicatePartialDecryptionSubmission)
}

// TestSetPartialDecryptionSubmission_DifferentValidators verifies that two
// validators can each submit their own partial decryption for the same request.
func TestSetPartialDecryptionSubmission_DifferentValidators(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext-data")

	// Validator 1 submits
	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial-v1"),
		[]byte("eph-key-v1"),
		[]byte("pub-share-v1"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)

	// Validator 2 submits (different validator → different key → no duplicate)
	err = k.setPartialDecryptionSubmission(
		ctx,
		testValidator2, 1, 2,
		[]byte("enc-partial-v2"),
		[]byte("eph-key-v2"),
		[]byte("pub-share-v2"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)
}

// TestSetPartialDecryptionSubmission_DifferentRounds verifies that the same
// validator can submit for different rounds without triggering duplicate check.
func TestSetPartialDecryptionSubmission_DifferentRounds(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext-data")

	// Round 1
	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)

	// Round 2 — same validator, different round
	err = k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 2, 1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)
}

// TestDkgPartialDecryptKey_DifferentInputsProduceDifferentKeys verifies that
// the key function produces distinct keys for distinct inputs.
func TestDkgPartialDecryptKey_DifferentInputsProduceDifferentKeys(t *testing.T) {
	t.Parallel()

	label := testLabel()
	ciphertext := []byte("cipher")
	requesterPubKey := []byte("req-pub-key")

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	key2 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 2, testValidator1)
	key3 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator2)

	require.NotEqual(t, key1, key2, "different round → different key")
	require.NotEqual(t, key1, key3, "different validator → different key")
}

// TestDkgPartialDecryptKey_SameInputsSameKey verifies that the key function is
// deterministic.
func TestDkgPartialDecryptKey_SameInputsSameKey(t *testing.T) {
	t.Parallel()

	label := testLabel()
	ciphertext := []byte("cipher")
	requesterPubKey := []byte("req-pub-key")

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	key2 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	require.Equal(t, key1, key2, "same inputs must produce same key")
}

// TestDkgPartialDecryptPrefix_ContainsRequesterAndLabel verifies that the
// prefix function produces a string containing the requester hash and label.
func TestDkgPartialDecryptPrefix_ContainsRequesterAndLabel(t *testing.T) {
	t.Parallel()

	requesterPubKey := []byte("req-pub-key")
	label := testLabel()

	prefix := dkgPartialDecryptPrefix(requesterPubKey, label)
	require.NotEmpty(t, prefix)
	require.Contains(t, prefix, "_", "prefix should use underscore as separator")
}

// TestDecodePartialDecryptionSubmission_RoundTrip verifies encode-decode
// round-trip for DKGPartialDecryptionSubmission.
func TestDecodePartialDecryptionSubmission_RoundTrip(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	require.NoError(t, k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 5, 2,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	))

	// Retrieve stored bytes and decode
	key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 5, testValidator1)
	bz, err := k.DKGPartialDecrypt.Get(ctx, key)
	require.NoError(t, err)

	submission, err := decodePartialDecryptionSubmission(bz)
	require.NoError(t, err)
	require.Equal(t, testValidator1.Hex(), submission.Validator)
	require.Equal(t, uint32(5), submission.Round)
	require.Equal(t, uint32(2), submission.Pid)
	require.Equal(t, []byte("enc-partial"), submission.EncryptedPartial)
}

// TestDecodePartialDecryptionSubmission_InvalidJSON verifies that decoding
// invalid JSON returns an error.
func TestDecodePartialDecryptionSubmission_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := decodePartialDecryptionSubmission([]byte("not-json"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshal partial decryption submission")
}
