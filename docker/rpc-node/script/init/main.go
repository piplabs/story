package main

import (
	"crypto/ecdsa"
	"flag"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"path/filepath"
)

const (
	iliadPrivValidatorKeyFile = "priv_validator_key.json"
	iliadNodeKeyFile          = "node_key.json"
	gethNodeKeyFile           = "nodekey"
)

func main() {
	configDir := flag.String("config-dir", "", "config directory")
	flag.Parse()

	pv := privval.NewFilePV(secp256k1.GenPrivKey(), "", "")

	if jsbz, err := json.Marshal(pv.Key); err != nil {
		log.Panicf("marshal %s error: %v", iliadPrivValidatorKeyFile, err)
	} else if err = os.WriteFile(filepath.Join(*configDir, "story", "config", iliadPrivValidatorKeyFile), jsbz, 0644); err != nil {
		log.Panicf("write %s error: %v", iliadPrivValidatorKeyFile, err)
	}

	if _, err := p2p.LoadOrGenNodeKey(filepath.Join(*configDir, "story", "config", iliadNodeKeyFile)); err != nil {
		log.Panicf("gen %s error: %v", iliadNodeKeyFile, err)
	}

	if _, err := loadOrGenGethNodeKey(filepath.Join(*configDir, "geth", "config", gethNodeKeyFile)); err != nil {
		log.Panicf("gen %s error: %v", gethNodeKeyFile, err)
	}
}

func loadOrGenGethNodeKey(path string) (*ecdsa.PrivateKey, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return crypto.LoadECDSA(path)
	} else {
		nodeKey, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		if err := crypto.SaveECDSA(path, nodeKey); err != nil {
			return nil, err
		}
		return nodeKey, nil
	}
}
