package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
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

func uncmpPubKeyToCmpPubKey(uncmpPubKey []byte) ([]byte, error) {
	if len(uncmpPubKey) != 65 || uncmpPubKey[0] != 0x04 {
		return nil, errors.New("invalid uncompressed public key length or format")
	}

	x := new(big.Int).SetBytes(uncmpPubKey[1:33])
	y := new(big.Int).SetBytes(uncmpPubKey[33:])

	pubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return crypto.CompressPubkey(&pubKey), nil
}

func printKeyFormats(compressedPubKeyBytes []byte) error {
	compressedPubKeyBase64 := base64.StdEncoding.EncodeToString(compressedPubKeyBytes)
	uncompressedPubKeyBytes, err := cmpPubKeyToUncmpPubKey(compressedPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to convert compressed pub key to uncompressed format")
	}

	uncompressedPubKeyHex := hex.EncodeToString(uncompressedPubKeyBytes)

	evmAddr, err := k1util.CosmosPubkeyToEVMAddress(compressedPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to convert to evm address")
	}

	validatorAddress := cosmostypes.ValAddress(evmAddr.Bytes())
	delegatorAddress := cosmostypes.AccAddress(evmAddr.Bytes())

	fmt.Println("Compressed Public Key (hex):", hex.EncodeToString(compressedPubKeyBytes))
	fmt.Println("Compressed Public Key (base64):", compressedPubKeyBase64)
	fmt.Println("Uncompressed Public Key (hex):", uncompressedPubKeyHex)
	fmt.Println("EVM Address:", evmAddr.String())
	fmt.Println("Validator Address:", validatorAddress.String())
	fmt.Println("Delegator Address:", delegatorAddress.String())

	return nil
}
