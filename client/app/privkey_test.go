package app_test

import (
	"path/filepath"
	"testing"

	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/privval"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/app"
)

func setupTestEnv(t *testing.T) (string, string, string) {
	t.Helper()

	stateFileDir := filepath.Join(t.TempDir(), "stateFileDir")
	encFileDir := filepath.Join(t.TempDir(), "encFileDir")
	password := "testpassword"

	return stateFileDir, encFileDir, password
}

func TestEncryptAndDecrypt_Success(t *testing.T) {
	stateFileDir, encFileDir, password := setupTestEnv(t)

	pv := privval.NewFilePV(k1.GenPrivKey(), "", stateFileDir)

	// Encryption
	err := app.EncryptAndStoreKey(pv.Key, password, encFileDir)
	require.NoError(t, err)

	// Decryption
	loadedKey, err := app.LoadEncryptedPrivKey(password, encFileDir)
	require.NoError(t, err)

	assert.Equal(t, pv.Key, loadedKey, "The decrypted key must match the original.")
}

func TestLoadEncryptedPrivKey_WrongPassword(t *testing.T) {
	stateFileDir, encFileDir, password := setupTestEnv(t)
	wrongPassword := "wrongpassword"

	pv := privval.NewFilePV(k1.GenPrivKey(), "", stateFileDir)

	// Encryption
	err := app.EncryptAndStoreKey(pv.Key, password, encFileDir)
	require.NoError(t, err)

	// Decrypt with wrong password
	_, err = app.LoadEncryptedPrivKey(wrongPassword, encFileDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wrong password for wallet entered")
}

func TestLoadEncryptedPrivKey_FileNotFound(t *testing.T) {
	_, encFileDir, password := setupTestEnv(t)

	_, err := app.LoadEncryptedPrivKey(password, encFileDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read enc_priv_key.json file")
}
