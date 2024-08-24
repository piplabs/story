package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

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

func decodeAndUncompressPubKey(compressedPubKeyBase64 string) (string, error) {
	compressedPubKeyBytes, err := base64.StdEncoding.DecodeString(compressedPubKeyBase64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 public key")
	}
	if len(compressedPubKeyBytes) != 33 {
		return "", fmt.Errorf("invalid compressed public key length: %d", len(compressedPubKeyBytes))
	}

	curve := elliptic.P256()
	x, y := elliptic.UnmarshalCompressed(curve, compressedPubKeyBytes)
	if x == nil || y == nil {
		return "", errors.New("failed to unmarshal compressed public key")
	}

	uncompressedPubKeyBytes := elliptic.Marshal(curve, x, y)
	uncompressedPubKeyHex := hex.EncodeToString(uncompressedPubKeyBytes)

	return uncompressedPubKeyHex, nil
}

func deriveUncompressedPublicKeyFromPrivateKey(evmPrivKey *ecdsa.PrivateKey) ([]byte, error) {
	pubKey := evmPrivKey.PublicKey
	uncompressedPubKey := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	if len(uncompressedPubKey) != 65 {
		return nil, fmt.Errorf("invalid uncompressed public key length: %d", len(uncompressedPubKey))
	}

	return uncompressedPubKey, nil
}

func validatorKeyExport(keyFilePath string) error {
	keyFileBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read key file")
	}

	var keyData ValidatorKey
	if err := json.Unmarshal(keyFileBytes, &keyData); err != nil {
		return errors.Wrap(err, "failed to unmarshal key file")
	}

	privKeyBytes, err := base64.StdEncoding.DecodeString(keyData.PrivKey.Value)
	if err != nil {
		return errors.Wrap(err, "failed to decode private key")
	}

	privateKey, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return errors.Wrap(err, "invalid private key")
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ecdsa.PublicKey")
	}
	evmPublicKey := crypto.PubkeyToAddress(*publicKey).Hex()

	// Handle the compressed public key
	compressedPubKeyBytes, err := base64.StdEncoding.DecodeString(keyData.PubKey.Value)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 public key")
	}
	compressedPubKeyHex := hex.EncodeToString(compressedPubKeyBytes)

	// Get the uncompressed public key using the refactored function
	uncompressedPubKeyHex, err := decodeAndUncompressPubKey(keyData.PubKey.Value)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------")
	fmt.Println("EVM Public Key:", evmPublicKey)
	fmt.Println("Compressed Public Key:", compressedPubKeyHex)
	fmt.Println("Uncompressed Public Key:", uncompressedPubKeyHex)
	fmt.Println("------------------------------------------------------")

	return nil
}
