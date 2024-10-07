package cmd

import (
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	cosmosk1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/joho/godotenv"
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

func loadValidatorFile(path string) ([]byte, error) {
	keyFileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read validator key file")
	}

	var keyData ValidatorKey
	if err := json.Unmarshal(keyFileBytes, &keyData); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal validator key file")
	}

	privKeyBytes, err := base64.StdEncoding.DecodeString(keyData.PrivKey.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode private key")
	}
	return privKeyBytes, nil
}

func loadPrivKeyFile(path string) ([]byte, error) {
	envMap, err := godotenv.Read(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read .env file")
	}

	privateKey, exists := envMap["PRIVATE_KEY"]
	if !exists || privateKey == "" {
		return nil, errors.New("no private key found in file")
	}

	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode private key")
	}

	return privKeyBytes, nil
}

func privKeyFileToCmpPubKey(path string) ([]byte, error) {
	privKeyBytes, err := loadPrivKeyFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load priv key file")
	}

	return privKeyToCmpPubKey(privKeyBytes)
}

func validatorKeyFileToCmpPubKey(path string) ([]byte, error) {
	privKeyBytes, err := loadValidatorFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load validator key file")
	}
	return privKeyToCmpPubKey(privKeyBytes)
}

func privKeyToCmpPubKey(privateKeyBytes []byte) ([]byte, error) {
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "invalid private key")
	}

	publicKey := &privateKey.PublicKey

	compressedPubKeyBytes := crypto.CompressPubkey(publicKey)

	return compressedPubKeyBytes, nil
}

func cmpPubKeyToUncmpPubKey(compressedPubKeyBytes []byte) ([]byte, error) {
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

func cmpPubKeyToEVMAddress(cmpPubKey []byte) (string, error) {
	if len(cmpPubKey) != secp256k1.PubKeyBytesLenCompressed {
		return "", fmt.Errorf("invalid compressed public key length: %d", len(cmpPubKey))
	}

	pubKey, err := crypto.DecompressPubkey(cmpPubKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to decompress public key")
	}
	evmAddress := crypto.PubkeyToAddress(*pubKey).Hex()
	return evmAddress, nil
}

func cmpPubKeyToDelegatorAddress(cmpPubKey []byte) (string, error) {
	if len(cmpPubKey) != secp256k1.PubKeyBytesLenCompressed {
		return "", fmt.Errorf("invalid compressed public key length: %d", len(cmpPubKey))
	}

	pubKey := &cosmosk1.PubKey{Key: cmpPubKey}
	return cosmostypes.AccAddress(pubKey.Address().Bytes()).String(), nil
}

func cmpPubKeyToValidatorAddress(cmpPubKey []byte) (string, error) {
	if len(cmpPubKey) != secp256k1.PubKeyBytesLenCompressed {
		return "", fmt.Errorf("invalid compressed public key length: %d", len(cmpPubKey))
	}

	pubKey := &cosmosk1.PubKey{Key: cmpPubKey}
	return cosmostypes.ValAddress(pubKey.Address().Bytes()).String(), nil
}

func printKeyFormats(compressedPubKeyBytes []byte) error {
	compressedPubKeyBase64 := base64.StdEncoding.EncodeToString(compressedPubKeyBytes)
	evmAddress, err := cmpPubKeyToEVMAddress(compressedPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to convert compressed pub key to EVM address")
	}

	uncompressedPubKeyHex, err := cmpPubKeyToUncmpPubKey(compressedPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to convert compressed pub key to uncompressed format")
	}

	validatorAddress, err := cmpPubKeyToValidatorAddress(compressedPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to convert compressed pub key to validator address")
	}

	delegatorAddress, err := cmpPubKeyToDelegatorAddress(compressedPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to convert compressed pub key to delegator address")
	}

	fmt.Println("Compressed Public Key (hex):", hex.EncodeToString(compressedPubKeyBytes))
	fmt.Println("Compressed Public Key (base64):", compressedPubKeyBase64)
	fmt.Println("Uncompressed Public Key (hex):", uncompressedPubKeyHex)
	fmt.Println("EVM Address:", evmAddress)
	fmt.Println("Validator Address:", validatorAddress)
	fmt.Println("Delegator Address:", delegatorAddress)

	return nil
}
