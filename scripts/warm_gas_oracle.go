//go:build ignore

// warm_gas_oracle sends periodic low-fee transfers to keep the go-ethereum gas
// price oracle from getting stuck on a stale cached price (which happens when
// blocks are empty and lastPrice never updates).
//
// Usage:
//
//	go run scripts/warm_gas_oracle.go \
//	  -rpc   https://aeneid.storyrpc.io \
//	  -key   <hex-private-key> \
//	  -to    <recipient-address>   \
//	  -every 10s
package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpc := flag.String("rpc", "https://aeneid.storyrpc.io", "Ethereum JSON-RPC endpoint")
	rawKey := flag.String("key", "", "hex private key (without 0x)")
	toAddr := flag.String("to", "", "recipient address")
	interval := flag.Duration("every", 10*time.Second, "interval between transfers")
	tipGwei := flag.Int64("tip", 1, "priority fee in gwei")
	flag.Parse()

	if *rawKey == "" || *toAddr == "" {
		fmt.Fprintln(os.Stderr, "usage: -key <privkey> -to <address>")
		os.Exit(1)
	}

	privKey, err := crypto.HexToECDSA(*rawKey)
	must(err, "parse private key")

	pubKey := privKey.Public().(*ecdsa.PublicKey)
	from := crypto.PubkeyToAddress(*pubKey)
	to := common.HexToAddress(*toAddr)

	client, err := ethclient.Dial(*rpc)
	must(err, "dial rpc")
	defer client.Close()

	chainID, err := client.ChainID(context.Background())
	must(err, "get chain id")

	tip := new(big.Int).Mul(big.NewInt(*tipGwei), big.NewInt(1e9))
	fmt.Printf("from=%s  to=%s  chain=%s  tip=%d gwei  every=%s\n",
		from.Hex(), to.Hex(), chainID, *tipGwei, *interval)

	for {
		if err := send(client, privKey, from, to, chainID, tip); err != nil {
			fmt.Printf("[%s] ERROR: %v\n", ts(), err)
		}
		time.Sleep(*interval)
	}
}

func send(client *ethclient.Client, key *ecdsa.PrivateKey, from, to common.Address, chainID, tip *big.Int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Query current oracle price for logging.
	oraclePrice, _ := client.SuggestGasPrice(ctx)

	baseFee, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		baseFee = big.NewInt(0)
	}

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return fmt.Errorf("nonce: %w", err)
	}

	// feeCap = baseFee + tip; floor at tip so it's always valid.
	feeCap := new(big.Int).Add(baseFee, tip)
	if feeCap.Cmp(tip) < 0 {
		feeCap = new(big.Int).Set(tip)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tip,
		GasFeeCap: feeCap,
		Gas:       21000,
		To:        &to,
		Value:     big.NewInt(0),
	})

	signer := types.LatestSignerForChainID(chainID)
	signed, err := types.SignTx(tx, signer, key)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	if err := client.SendTransaction(ctx, signed); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	fmt.Printf("[%s] sent  hash=%s  oracle_before=%s gwei  tip=%s gwei\n",
		ts(), signed.Hash().Hex(),
		new(big.Int).Div(oraclePrice, big.NewInt(1e9)),
		new(big.Int).Div(tip, big.NewInt(1e9)),
	)

	return nil
}

func must(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}
}

func ts() string {
	return time.Now().Format("15:04:05")
}
