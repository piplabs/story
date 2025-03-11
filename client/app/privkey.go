package app

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/cometbft/cometbft/crypto"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/privval"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
)

// loadPrivVal returns a privval.FilePV by loading either a CometBFT priv validator key or an Ethereum keystore file.
func loadPrivVal(cfg Config) (*privval.FilePV, error) {
	cmtFile := cfg.Comet.PrivValidatorKeyFile()
	encPrivKeyFile := cfg.EncPrivKeyFile()
	cmtExists := exists(cmtFile)
	encPrivExists := exists(encPrivKeyFile)

	if !cmtExists && !encPrivExists {
		return nil, errors.New("no cometBFT priv validator key file exists", "comet_file", cmtFile, "enc_priv_key_file", encPrivKeyFile)
	}

	var (
		key crypto.PrivKey
		err error
	)
	if encPrivExists {
		password, err := InputPassword(
			PasswordPromptText,
			"",
			false,
			ValidatePasswordInput,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error occurred while input password")
		}

		pv, err := LoadEncryptedPrivKey(password, encPrivKeyFile)
		if err != nil {
			return nil, err
		}
		key = pv.PrivKey
	} else {
		key, err = loadCometFilePV(cmtFile)
		if err != nil {
			return nil, err
		}
	}

	state, err := loadCometPVState(cfg.Comet.PrivValidatorStateFile())
	if err != nil {
		return nil, err
	}

	// Create a new privval.FilePV with the loaded key and state.
	// This is a workaround for the fact that there is no other way
	// to set FilePVLastSignState filePath field.
	resp := privval.NewFilePV(key, "", cfg.Comet.PrivValidatorStateFile())
	resp.LastSignState.Step = state.Step
	resp.LastSignState.Round = state.Round
	resp.LastSignState.Height = state.Height
	resp.LastSignState.Signature = state.Signature
	resp.LastSignState.SignBytes = state.SignBytes

	return resp, nil
}

// loadEthKeystore loads an Ethereum keystore file and returns the private key.
//
//nolint:unused //Ignore unused function temporarily
func loadEthKeystore(keystoreFile string, password string) (crypto.PrivKey, error) {
	bz, err := os.ReadFile(keystoreFile)
	if err != nil {
		return nil, errors.Wrap(err, "read keystore file", "path", keystoreFile)
	}

	key, err := keystore.DecryptKey(bz, password)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt keystore file", "path", keystoreFile)
	}

	return k1util.StdPrivKeyToComet(key.PrivateKey)
}

// loadCometFilePV loads a CometBFT privval file and returns the private key.
func loadCometFilePV(file string) (crypto.PrivKey, error) {
	bz, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read comet privval", "path", file)
	}

	var pvKey privval.FilePVKey
	err = cmtjson.Unmarshal(bz, &pvKey)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal comet privval", "path", file)
	}

	return pvKey.PrivKey, nil
}

// loadCometPVState loads a CometBFT privval state file.
func loadCometPVState(file string) (privval.FilePVLastSignState, error) {
	bz, err := os.ReadFile(file)
	if err != nil {
		return privval.FilePVLastSignState{}, errors.Wrap(err, "read comet privval state", "path", file)
	}

	var state privval.FilePVLastSignState
	err = cmtjson.Unmarshal(bz, &state)
	if err != nil {
		return privval.FilePVLastSignState{}, errors.Wrap(err, "unmarshal comet privval state", "path", file)
	}

	return state, nil
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

// EncryptedKeyRepresentation defines an internal representation of encrypted validator key.
type EncryptedKeyRepresentation struct {
	Crypto  map[string]interface{} `json:"crypto"` //nolint:revive // This is from Prysm.
	Version uint                   `json:"version"`
	Name    string                 `json:"name"`
}

func EncryptAndStoreKey(key privval.FilePVKey, password, filePath string) error {
	encodedKey, err := cmtjson.MarshalIndent(key, "", "\t")
	if err != nil {
		return errors.Wrap(err, "failed to marshal key for encryption")
	}

	encryptor := keystorev4.New()
	encryptedKey, err := encryptor.Encrypt(encodedKey, password)
	if err != nil {
		return errors.Wrap(err, "could not encrypt key")
	}

	encKeyRepr := EncryptedKeyRepresentation{
		Crypto:  encryptedKey,
		Version: encryptor.Version(),
		Name:    encryptor.Name(),
	}

	data, err := json.MarshalIndent(encKeyRepr, "", "\t")
	if err != nil {
		return errors.Wrap(err, "failed to marshal encrypted key")
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return errors.Wrap(err, "failed to write enc_priv_key.json file")
	}

	return nil
}

func LoadEncryptedPrivKey(password, encPrivKeyFile string) (privval.FilePVKey, error) {
	data, err := os.ReadFile(encPrivKeyFile)
	if err != nil {
		return privval.FilePVKey{}, errors.Wrap(err, "failed to read enc_priv_key.json file")
	}

	var encKeyRepr EncryptedKeyRepresentation
	if err := json.Unmarshal(data, &encKeyRepr); err != nil {
		return privval.FilePVKey{}, errors.Wrap(err, "failed to unmarshal enc_priv_key.json data")
	}

	decryptor := keystorev4.New()
	decryptedKey, err := decryptor.Decrypt(encKeyRepr.Crypto, password)
	if err != nil && strings.Contains(err.Error(), "invalid checksum") {
		return privval.FilePVKey{}, errors.Wrap(err, "wrong password for wallet entered")
	} else if err != nil {
		return privval.FilePVKey{}, errors.Wrap(err, "could not decrypt key")
	}

	var key privval.FilePVKey
	if err := cmtjson.Unmarshal(decryptedKey, &key); err != nil {
		return privval.FilePVKey{}, errors.Wrap(err, "failed to unmarshal decrypted key")
	}

	return key, nil
}
