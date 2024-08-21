// ADAPTED FROM: https://github.com/Layr-Labs/eigensdk-go/blob/dev/crypto/ecdsa/utils.go
package ecdsa

import (
	"bufio"
	"crypto/ecdsa"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"github.com/storyprotocol/iliad/lib/errors"
)

// WriteKey writes the private key to the given path
// The key is encrypted using the given password
// This function will create the directory if it doesn't exist
// If there's an existing file at the given path, it will be overwritten.
func WriteKey(path string, privateKey *ecdsa.PrivateKey, password string) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "generating UUID")
	}

	// We are using https://github.com/ethereum/go-ethereum/blob/master/accounts/keystore/key.go#L41
	// to store the keys which requires us to have random UUID for encryption
	key := &keystore.Key{
		Id:         uuid,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	encryptedBytes, err := keystore.EncryptKey(key, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return errors.Wrap(err, "encrypting key")
	}

	return writeBytesToFile(path, encryptedBytes)
}

func writeBytesToFile(path string, data []byte) error {
	dir := filepath.Dir(path)

	// create the directory if it doesn't exist. If exists, it does nothing
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(err, "creating directories")
	}

	file, err := os.Create(filepath.Clean(path))
	if err != nil {
		return errors.Wrap(err, "creating file")
	}
	// remember to close the file
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	_, err = file.Write(data)

	return errors.Wrap(err, "writing to file")
}
