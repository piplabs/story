package cmd

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/piplabs/story/lib/errors"
)

func readContract(ctx context.Context, cfg baseConfig, contractAddress common.Address, data []byte) ([]byte, error) {
	client, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Ethereum client")
	}

	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	result, err := client.CallContract(ctx, callMsg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "contract call failed")
	}

	return result, nil
}

func prepareAndSendTransaction(ctx context.Context, cfg baseConfig, contractAddress common.Address, value *big.Int, data []byte) error {
	client, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Ethereum client")
	}

	evmPrivKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "invalid EVM private key")
	}

	publicKey, ok := evmPrivKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to assert type to *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get nonce")
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to suggest gas price")
	}

	gasPrice.Mul(gasPrice, big.NewInt(120)).Div(gasPrice, big.NewInt(100))

	gasLimit, err := estimateGas(ctx, client, fromAddress, contractAddress, gasPrice, value, data)
	if err != nil {
		return err
	}

	gasTipCap := gasPrice
	gasFeeCap := new(big.Int).Mul(gasPrice, big.NewInt(2))

	gasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasFeeCap)
	totalTxCost := new(big.Int).Add(gasCost, value)

	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to fetch balance")
	}

	if balance.Cmp(totalTxCost) < 0 {
		return errors.New("insufficient funds for gas * price + value", "balance", balance.String(), "totalTxCost", totalTxCost.String())
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(cfg.ChainID),
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &contractAddress,
		Value:     value,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(cfg.ChainID)), evmPrivKey)
	if err != nil {
		return errors.Wrap(err, "failed to sign transaction")
	}

	txHash := signedTx.Hash().Hex()
	fmt.Printf("Transaction hash: %s\n", txHash)
	fmt.Printf("Explorer URL: %s/tx/%s\n", cfg.Explorer, txHash)

	if err = client.SendTransaction(ctx, signedTx); err != nil {
		return errors.Wrap(err, "failed to send transaction")
	}

	fmt.Println("Transaction sent, waiting for confirmation...")

	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return errors.New("transaction failed", "status", receipt.Status)
	}

	fmt.Println("Transaction confirmed successfully!")

	return nil
}

func estimateGas(ctx context.Context, client *ethclient.Client, fromAddress, contractAddress common.Address, gasPrice, value *big.Int, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &contractAddress,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, errors.Wrap(err, "failed to estimate gas")
	}

	return gasLimit, nil
}
