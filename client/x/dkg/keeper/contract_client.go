package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/piplabs/story/client/genutil/evm/predeploys"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	maxRetries = 3

	// Node status in DKG contract
	NodeStatusUnregistered   uint8 = 0
	NodeStatusRegistered     uint8 = 1
	NodeStatusNetworkSetDone uint8 = 2
	NodeStatusFinalized      uint8 = 3
)

// ContractClient wraps the DKG contract interaction.
type ContractClient struct {
	ethClient       *ethclient.Client
	dkgContract     *bindings.DKG
	dkgContractAbi  *abi.ABI
	dkgContractAddr common.Address
	privateKey      *ecdsa.PrivateKey
	fromAddress     common.Address
	chainID         *big.Int
}

// ContractConfig holds configuration for contract interaction.
type ContractConfig struct {
	EthRPCEndpoint  string
	DKGContractAddr string
	PrivateKey      string
	ChainID         int64
}

// NewContractClient creates a new contract client.
func NewContractClient(ctx context.Context, engineEndpoint string, engineChainID int64, privKey []byte) (*ContractClient, error) {
	ethClient, err := ethclient.Dial(engineEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Ethereum client")
	}

	contractAddr := common.HexToAddress(predeploys.DKG)
	dkgContract, err := bindings.NewDKG(contractAddr, ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DKG contract instance")
	}

	dkgContractAbi, err := bindings.DKGMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DKG contract ABI")
	}

	privateKey, err := crypto.ToECDSA(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID := big.NewInt(engineChainID)

	client := &ContractClient{
		ethClient:       ethClient,
		dkgContract:     dkgContract,
		dkgContractAbi:  dkgContractAbi,
		dkgContractAddr: contractAddr,
		privateKey:      privateKey,
		fromAddress:     fromAddress,
		chainID:         chainID,
	}

	log.Info(ctx, "Created contract client",
		"dkg_contract_address", contractAddr.Hex(),
		"from_address", fromAddress.Hex(),
		"chain_id", chainID,
	)

	return client, nil
}

// InitializeDKG calls the initializeDKG contract method.
func (c *ContractClient) InitializeDKG(ctx context.Context, round uint32, mrenclave []byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Receipt, error) {
	log.Info(ctx, "Calling initializeDKG contract method",
		"round", round,
		"mrenclave", hex.EncodeToString(mrenclave),
		"dkg_pub_key", hex.EncodeToString(dkgPubKey),
		"comm_pub_key", hex.EncodeToString(commPubKey),
		"raw_quote_len", len(rawQuote),
	)

	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert bytes32")
	}

	callData, err := c.dkgContractAbi.Pack("initializeDKG", round, mrenclave32, dkgPubKey, commPubKey, rawQuote)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack initialize dkg call data")
	}

	return c.sendWithRetry(ctx, "InitializeDKG", callData, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.dkgContract.InitializeDKG(auth, round, mrenclave32, dkgPubKey, commPubKey, rawQuote)
	})
}

// FinalizeDKG calls the finalizeDKG contract method.
func (c *ContractClient) FinalizeDKG(
	ctx context.Context,
	round uint32,
	mrenclave []byte,
	globalPubKey []byte,
	signature []byte,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling finalizeDKG contract method",
		"mrenclave", hex.EncodeToString(mrenclave),
		"round", round,
		"global_pub_key", hex.EncodeToString(globalPubKey),
		"signature_len", len(signature),
	)

	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert bytes32")
	}

	callData, err := c.dkgContractAbi.Pack("finalizeDKG", round, mrenclave32, globalPubKey, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack finalizeDKG call data")
	}

	return c.sendWithRetry(ctx, "FinalizeDKG", callData, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.dkgContract.FinalizeDKG(auth, round, mrenclave32, globalPubKey, signature)
	})
}

// SetNetwork calls the setNetwork contract method.
func (c *ContractClient) SetNetwork(
	ctx context.Context,
	round uint32,
	total uint32,
	threshold uint32,
	mrenclave []byte,
	signature []byte,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling setNetwork contract method",
		"mrenclave", hex.EncodeToString(mrenclave),
		"round", round,
		"total", total,
		"threshold", threshold,
		"signature_len", len(signature),
	)

	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert bytes32")
	}

	callData, err := c.dkgContractAbi.Pack("setNetwork", round, total, threshold, mrenclave32, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack setNetwork call data")
	}

	return c.sendWithRetry(ctx, "SetNetwork", callData, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.dkgContract.SetNetwork(auth, round, total, threshold, mrenclave32, signature)
	})
}

// RequestRemoteAttestationOnChain calls the requestRemoteAttestationOnChain contract method.
func (c *ContractClient) RequestRemoteAttestationOnChain(
	ctx context.Context,
	targetValidatorAddr common.Address,
	round uint32,
	mrenclave []byte,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling requestRemoteAttestationOnChain contract method",
		"mrenclave", string(mrenclave),
		"round", round,
		"target_validator_addr", targetValidatorAddr.Hex(),
	)

	callData, err := c.dkgContractAbi.Pack("requestRemoteAttestationOnChain", targetValidatorAddr, round, mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack requestRemoteAttestationOnChain call data")
	}

	gasLimit, err := c.estimateGasWithBuffer(ctx, callData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to estimate gas for requestRemoteAttestationOnChain")
	}

	auth, err := c.createTransactOpts(ctx, gasLimit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transaction options")
	}

	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert bytes32")
	}

	tx, err := c.dkgContract.RequestRemoteAttestationOnChain(auth, targetValidatorAddr, round, mrenclave32)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call request remote attestation on chain")
	}

	log.Info(ctx, "RequestRemoteAttestationOnChain transaction sent", "tx_hash", tx.Hash().Hex())

	receipt, err := c.waitForTransaction(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for request remote attestation on chain transaction")
	}

	return receipt, nil
}

// ComplainDeals calls the complainDeals contract method.
func (c *ContractClient) ComplainDeals(
	ctx context.Context,
	round uint32,
	index uint32,
	complainIndexes []uint32,
	mrenclave []byte,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling complainDeals contract method",
		"mrenclave", string(mrenclave),
		"round", round,
		"index", index,
		"complain_indexes", complainIndexes,
	)

	callData, err := c.dkgContractAbi.Pack("complainDeals", round, index, complainIndexes, mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack complain deals call data")
	}

	gasLimit, err := c.estimateGasWithBuffer(ctx, callData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to estimate gas for complain deals")
	}

	auth, err := c.createTransactOpts(ctx, gasLimit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transaction options")
	}

	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert bytes32")
	}

	tx, err := c.dkgContract.ComplainDeals(auth, round, index, complainIndexes, mrenclave32)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call complain deals")
	}

	log.Info(ctx, "Complain deals transaction sent", "tx_hash", tx.Hash().Hex())

	receipt, err := c.waitForTransaction(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for complain deals transaction")
	}

	return receipt, nil
}

// SubmitPartialDecryption calls the submitPartialDecryption contract method.
func (c *ContractClient) SubmitPartialDecryption(
	ctx context.Context,
	round uint32,
	mrenclave []byte,
	encryptedPartial []byte,
	ephemeralPubKey []byte,
	pubShare []byte,
	label []byte,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling submitPartialDecryption contract method",
		"mrenclave", hex.EncodeToString(mrenclave),
		"round", round,
		"partial_len", len(encryptedPartial),
		"eph_pub_len", len(ephemeralPubKey),
		"pub_share_len", len(pubShare),
		"label_len", len(label),
	)

	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert bytes32")
	}

	callData, err := c.dkgContractAbi.Pack("submitPartialDecryption", round, mrenclave32, encryptedPartial, ephemeralPubKey, pubShare, label)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack submitPartialDecryption call data")
	}

	return c.sendWithRetry(ctx, "SubmitPartialDecryption", callData, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		bound := bind.NewBoundContract(c.dkgContractAddr, *c.dkgContractAbi, c.ethClient, c.ethClient, c.ethClient)
		return bound.Transact(auth, "submitPartialDecryption", round, mrenclave32, encryptedPartial, ephemeralPubKey, pubShare, label)
	})
}

// GetNodeInfo queries node information from the contract.
func (c *ContractClient) GetNodeInfo(ctx context.Context, mrenclave []byte, round uint32, validatorAddr common.Address) (*bindings.IDKGNodeInfo, error) {
	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert bytes32")
	}

	nodeInfo, err := c.dkgContract.GetNodeInfo(&bind.CallOpts{Context: ctx}, mrenclave32, round, validatorAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get node info")
	}

	return &nodeInfo, nil
}

func (c *ContractClient) IsInitialized(ctx context.Context, round uint32, mrenclave []byte, validator common.Address) (bool, error) {
	nodeInfo, err := c.GetNodeInfo(ctx, mrenclave, round, validator)
	if err != nil {
		return false, err
	}

	return nodeInfo.NodeStatus == NodeStatusRegistered, nil
}

func (c *ContractClient) IsNetworkSet(ctx context.Context, round uint32, mrenclave []byte, validator common.Address) (bool, error) {
	nodeInfo, err := c.GetNodeInfo(ctx, mrenclave, round, validator)
	if err != nil {
		return false, err
	}

	return nodeInfo.NodeStatus == NodeStatusNetworkSetDone, nil
}

func (c *ContractClient) IsFinalized(ctx context.Context, round uint32, mrenclave []byte, validator common.Address) (bool, error) {
	nodeInfo, err := c.GetNodeInfo(ctx, mrenclave, round, validator)
	if err != nil {
		return false, err
	}

	return nodeInfo.NodeStatus == NodeStatusFinalized, nil
}

// createTransactOpts creates transaction options for contract calls.
func (c *ContractClient) createTransactOpts(ctx context.Context, gasLimit uint64) (*bind.TransactOpts, error) {
	nonce, err := c.ethClient.PendingNonceAt(ctx, c.fromAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending nonce")
	}

	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get gas price")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transactor")
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	auth.Context = ctx

	return auth, nil
}

// estimateGasWithBuffer estimates gas for a contract transaction and adds a safety buffer.
func (c *ContractClient) estimateGasWithBuffer(ctx context.Context, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From: c.fromAddress,
		To:   &c.dkgContractAddr,
		Data: data,
	}

	gasLimit, err := c.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, errors.Wrap(err, "failed to estimate gas")
	}

	// 20% buffer
	return gasLimit * 12 / 10, nil
}

// waitForTransaction waits for a transaction to be mined and returns the receipt.
func (c *ContractClient) waitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	// 60 seconds timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)

	receipt, err := bind.WaitMined(timeoutCtx, c.ethClient, tx)
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for transaction to be mined")
	}

	log.Info(ctx, "Transaction mined", "tx_hash", tx.Hash().Hex())

	return receipt, nil
}

func (c *ContractClient) sendWithRetry(
	ctx context.Context,
	methodName string,
	callData []byte,
	sendTx func(auth *bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, error) {
	var (
		receipt  *types.Receipt
		gasLimit uint64
		err      error
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		gasLimit, err = c.estimateGasWithBuffer(ctx, callData)
		if err != nil {
			return nil, err
		}

		auth, err := c.createTransactOpts(ctx, gasLimit)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create transact opts", "method_name", methodName)
		}

		tx, err := sendTx(auth)
		if err != nil {
			return nil, errors.Wrap(err, "failed to send tx", "method_name", methodName)
		}

		log.Info(ctx, fmt.Sprintf("%s tx sent", methodName),
			"tx_hash", tx.Hash().Hex(),
			"attempt", attempt,
			"gas_limit", gasLimit,
		)

		receipt, err = c.waitForTransaction(ctx, tx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wait for tx", "method_name", methodName)
		}

		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Info(ctx, fmt.Sprintf("%s succeeded", methodName),
				"tx_hash", tx.Hash().Hex(),
				"gas_used", receipt.GasUsed,
				"attempt", attempt)

			return receipt, nil
		}

		usageRatio := float64(receipt.GasUsed) / float64(gasLimit)
		if usageRatio > 0.95 && attempt < maxRetries {
			log.Warn(ctx, fmt.Sprintf("%s likely out-of-gas, retrying", methodName),
				nil,
				"old_gas_limit", gasLimit,
				"gas_used", receipt.GasUsed,
				"usage_ratio", usageRatio,
				"attempt", attempt+1,
			)

			continue
		}

		break
	}

	return nil, errors.New(fmt.Sprintf("[%s] transaction failed after %d attempts", methodName, maxRetries))
}
