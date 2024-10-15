package keeper

import (
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/piplabs/story/lib/errors"
)

func UncmpPubKeyToCmpPubKey(uncmpPubKey []byte) ([]byte, error) {
	if len(uncmpPubKey) != 65 || uncmpPubKey[0] != 0x04 {
		return nil, errors.New("invalid uncompressed public key length or format")
	}

	// Extract the x and y coordinates
	x := new(big.Int).SetBytes(uncmpPubKey[1:33])
	y := new(big.Int).SetBytes(uncmpPubKey[33:])

	// Determine the prefix based on the parity of y
	prefix := byte(0x02) // Even y
	if y.Bit(0) == 1 {
		prefix = 0x03 // Odd y
	}

	// Construct the compressed public key
	compressedPubKey := make([]byte, 33)
	compressedPubKey[0] = prefix
	copy(compressedPubKey[1:], x.Bytes())

	return compressedPubKey, nil
}

func CmpPubKeyToUncmpPubkey(compressedPubKeyBytes []byte) ([]byte, error) {
	if len(compressedPubKeyBytes) != secp256k1.PubKeyBytesLenCompressed {
		return nil, fmt.Errorf("invalid compressed public key length: %d", len(compressedPubKeyBytes))
	}

	pubKey, err := secp256k1.ParsePubKey(compressedPubKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse compressed public key")
	}

	uncompressedPubKeyBytes := pubKey.SerializeUncompressed()

	return uncompressedPubKeyBytes, nil
}

func CmpPubKeyToEVMAddress(cmpPubKey []byte) (common.Address, error) {
	if len(cmpPubKey) != secp256k1.PubKeyBytesLenCompressed {
		return common.Address{}, fmt.Errorf("invalid compressed public key length: %d", len(cmpPubKey))
	}

	pubKey, err := crypto.DecompressPubkey(cmpPubKey)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to decompress public key")
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
