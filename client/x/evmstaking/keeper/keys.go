package keeper

import (
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/piplabs/story/lib/errors"
)

func UncmpPubKeyToCmpPubKey(uncmpPubKey []byte) ([]byte, error) {
	// Check if the uncompressed public key is 65 bytes and starts with 0x04
	if len(uncmpPubKey) != 65 || uncmpPubKey[0] != 0x04 {
		return nil, errors.New("invalid uncompressed public key length or format")
	}

	// Extract x and y coordinates from the uncompressed public key
	x := uncmpPubKey[1:33]
	y := uncmpPubKey[33:]

	// Determine if y is even or odd by checking the last byte of y
	var prefix byte
	if y[len(y)-1]%2 == 0 {
		prefix = 0x02 // even y-coordinate
	} else {
		prefix = 0x03 // odd y-coordinate
	}

	// Create compressed public key (1 byte prefix + 32 bytes x-coordinate)
	compressedPubKey := append([]byte{prefix}, x...)

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
