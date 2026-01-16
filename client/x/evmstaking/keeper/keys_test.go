package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
)

func generatePubKeyWithYParity(t *testing.T, wantEven bool) []byte {
	t.Helper()

	for range 100 {
		priv, err := crypto.GenerateKey()
		require.NoError(t, err)

		pub := crypto.FromECDSAPub(&priv.PublicKey)

		y := pub[33:]
		if (y[len(y)-1]%2 == 0) == wantEven {
			return pub
		}
	}

	t.Fatalf("failed to generate pubkey with desired y parity (even: %v)", wantEven)

	return nil
}

func TestUncmpPubKeyToCmpPubKey_YParityPrefix(t *testing.T) {
	// Even y → prefix 0x02
	evenPubKey := generatePubKeyWithYParity(t, true)
	evenCmp, err := keeper.UncmpPubKeyToCmpPubKey(evenPubKey)
	require.NoError(t, err)
	require.Equal(t, byte(0x02), evenCmp[0])

	// Odd y → prefix 0x03
	oddPubKey := generatePubKeyWithYParity(t, false)
	oddCmp, err := keeper.UncmpPubKeyToCmpPubKey(oddPubKey)
	require.NoError(t, err)
	require.Equal(t, byte(0x03), oddCmp[0])
}

func TestCmpPubKeyToUncmpPubKey_InvalidLength(t *testing.T) {
	// Length is too short (not 33)
	short := make([]byte, 10)
	_, err := keeper.CmpPubKeyToUncmpPubKey(short)
	require.ErrorContains(t, err, "invalid compressed public key length")
}

func TestCmpPubKeyToUncmpPubKey_ParseFailure(t *testing.T) {
	// Valid length but invalid content (not a valid curve point)
	invalid := make([]byte, 33)
	invalid[0] = 0x02 // valid prefix
	// all-zero X coord → not a valid point
	_, err := keeper.CmpPubKeyToUncmpPubKey(invalid)
	require.Error(t, err)
}

func TestCmpPubKeyToEVMAddress_InvalidLength(t *testing.T) {
	// Length is too short (not 33)
	short := make([]byte, 10)
	_, err := keeper.CmpPubKeyToEVMAddress(short)
	require.ErrorContains(t, err, "invalid compressed public key length")
}

func TestCmpPubKeyToEVMAddress_DecompressionError(t *testing.T) {
	// Valid length but not decompressible
	invalid := make([]byte, 33)
	invalid[0] = 0x02
	_, err := keeper.CmpPubKeyToEVMAddress(invalid)
	require.Error(t, err)
}

func TestUncmpPubKeyToCmpPubKey_NotOnCurve(t *testing.T) {
	notOnCurve := make([]byte, 65)
	notOnCurve[0] = 0x04
	// x = 0, y = 0 → definitely not on curve
	copy(notOnCurve[1:], make([]byte, 64))
	_, err := keeper.UncmpPubKeyToCmpPubKey(notOnCurve)
	require.Error(t, err)
}

func TestUncmpPubKeyToCmpPubKey_InvalidFormat(t *testing.T) {
	// wrong prefix
	wrongPrefix := make([]byte, 65)
	wrongPrefix[0] = 0x01
	_, err := keeper.UncmpPubKeyToCmpPubKey(wrongPrefix)
	require.Error(t, err)

	// wrong length
	_, err = keeper.UncmpPubKeyToCmpPubKey([]byte{0x04, 0x01})
	require.Error(t, err)
}

func TestEndToEndConversion(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	uncompressed := crypto.FromECDSAPub(&privKey.PublicKey)
	compressed, err := keeper.UncmpPubKeyToCmpPubKey(uncompressed)
	require.NoError(t, err)

	recovered, err := keeper.CmpPubKeyToUncmpPubKey(compressed)
	require.NoError(t, err)
	require.Equal(t, uncompressed, recovered)

	addr, err := keeper.CmpPubKeyToEVMAddress(compressed)
	require.NoError(t, err)
	require.Equal(t, crypto.PubkeyToAddress(privKey.PublicKey), addr)
}
