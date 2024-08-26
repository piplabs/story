package cmd

import (
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/piplabs/story/lib/errors"
)

type ValidatorKey struct {
	Address string  `json:"address"`
	PubKey  KeyInfo `json:"pub_key"`
	PrivKey KeyInfo `json:"priv_key"`
}

type KeyInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func uncompressPubKey(compressedPubKeyBase64 string) (string, error) {
	compressedPubKeyBytes, err := base64.StdEncoding.DecodeString(compressedPubKeyBase64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 public key")
	}
	if len(compressedPubKeyBytes) != secp256k1.PubKeyBytesLenCompressed {
		return "", fmt.Errorf("invalid compressed public key length: %d", len(compressedPubKeyBytes))
	}

	pubKey, err := secp256k1.ParsePubKey(compressedPubKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse compressed public key")
	}

	uncompressedPubKeyBytes := pubKey.SerializeUncompressed()
	uncompressedPubKeyHex := hex.EncodeToString(uncompressedPubKeyBytes)

	return uncompressedPubKeyHex, nil
}

func uncompressPrivateKey(privateKeyHex string) ([]byte, error) {
	evmPrivKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, errors.Wrap(err, "invalid EVM private key")
	}

	pubKey := evmPrivKey.PublicKey
	uncompressedPubKey := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	if len(uncompressedPubKey) != 65 {
		return nil, fmt.Errorf("invalid uncompressed public key length: %d", len(uncompressedPubKey))
	}

	return uncompressedPubKey, nil
}
