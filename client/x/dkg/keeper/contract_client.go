package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	maxRetries = 3
)

// ContractClient wraps the DKG contract interaction.
type ContractClient struct {
	ethClient       *ethclient.Client
	dkgContract     *bindings.DKG
	dkgContractAbi  *abi.ABI
	dkgContractAddr common.Address
	cdrContract     *bindings.CDR
	cdrContractAbi  *abi.ABI
	cdrContractAddr common.Address
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

	dkgContractAddr := common.HexToAddress(predeploys.DKG)

	dkgContract, err := bindings.NewDKG(dkgContractAddr, ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DKG contract instance")
	}

	dkgContractAbi, err := bindings.DKGMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DKG contract ABI")
	}

	cdrContractAddr := common.HexToAddress(predeploys.CDR)

	cdrContract, err := bindings.NewCDR(cdrContractAddr, ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CDR contract instance")
	}

	cdrContractAbi, err := bindings.CDRMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CDR contract ABI")
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
		dkgContractAddr: dkgContractAddr,
		cdrContract:     cdrContract,
		cdrContractAbi:  cdrContractAbi,
		cdrContractAddr: cdrContractAddr,
		privateKey:      privateKey,
		fromAddress:     fromAddress,
		chainID:         chainID,
	}

	log.Info(ctx, "Created contract client",
		"dkg_contract_address", dkgContractAddr.Hex(),
		"cdr_contract_address", cdrContractAddr.Hex(),
		"from_address", fromAddress.Hex(),
		"chain_id", chainID,
	)

	return client, nil
}

// Register calls the register contract method.
func (c *ContractClient) Register(ctx context.Context, round uint32, enclaveType [32]byte, startBlockHeight uint64, startBlockHash []byte, dkgPubKey []byte, commPubKey []byte, enclaveReport []byte) (*types.Receipt, error) {
	log.Info(ctx, "Calling register contract method",
		"round", round,
		"enclave_type", hex.EncodeToString(enclaveType[:]),
		"start_block_height", startBlockHeight,
		"start_block_hash", hex.EncodeToString(startBlockHash),
		"dkg_pub_key", hex.EncodeToString(dkgPubKey),
		"comm_pub_key", hex.EncodeToString(commPubKey),
		"raw_quote_len", len(enclaveReport),
	)

	startBlockHash32, err := cast.ToBytes32(startBlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert startBlockHash to bytes32")
	}

	enclaveInstanceData := bindings.IDKGEnclaveInstanceData{
		Round:          round,
		ValidatorAddr:  c.fromAddress,
		EnclaveType:    enclaveType,
		EnclaveCommKey: commPubKey,
		DkgPubKey:      dkgPubKey,
	}

	startBlockHeightBig := new(big.Int).SetUint64(startBlockHeight)

	callData, err := c.dkgContractAbi.Pack("register", enclaveReport, enclaveInstanceData, startBlockHeightBig, startBlockHash32, []byte{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack register call data")
	}

	return c.sendWithRetry(ctx, "Register", c.dkgContractAddr, callData, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.dkgContract.Register(auth, enclaveReport, enclaveInstanceData, startBlockHeightBig, startBlockHash32, []byte{})
	})
}

// Finalize calls the finalize contract method.
func (c *ContractClient) Finalize(
	ctx context.Context,
	round uint32,
	enclaveType [32]byte,
	participantsRoot []byte,
	globalPubKey []byte,
	publicCoeffs [][]byte,
	pubKeyShare []byte,
	signature []byte,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling finalize contract method",
		"enclave_type", hex.EncodeToString(enclaveType[:]),
		"round", round,
		"global_pub_key", hex.EncodeToString(globalPubKey),
		"pub_key_share", hex.EncodeToString(pubKeyShare),
		"signature_len", len(signature),
	)

	participantsRoot32, err := cast.ToBytes32(participantsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert participants root to bytes32")
	}

	callData, err := c.dkgContractAbi.Pack("finalize", round, c.fromAddress, enclaveType, participantsRoot32, globalPubKey, publicCoeffs, pubKeyShare, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack finalize call data")
	}

	return c.sendWithRetry(ctx, "Finalize", c.dkgContractAddr, callData, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.dkgContract.Finalize(auth, round, c.fromAddress, enclaveType, participantsRoot32, globalPubKey, publicCoeffs, pubKeyShare, signature)
	})
}

// SubmitEncryptedPartialDecryption calls the submitEncryptedPartialDecryption contract method.
func (c *ContractClient) SubmitEncryptedPartialDecryption(
	ctx context.Context,
	round uint32,
	pid uint32,
	encryptedPartial []byte,
	ephemeralPubKey []byte,
	pubShare []byte,
	requesterPubKey []byte,
	uuid uint32,
	signature []byte,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling submitEncryptedPartialDecryption contract method",
		"round", round,
		"pid", pid,
		"partial_len", len(encryptedPartial),
		"eph_pub_len", len(ephemeralPubKey),
		"pub_share_len", len(pubShare),
		"requester_pub_key_len", len(requesterPubKey),
		"uuid", uuid,
		"signature_len", len(signature),
	)

	callData, err := c.cdrContractAbi.Pack("submitEncryptedPartialDecryption", round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, uuid, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack submitEncryptedPartialDecryption call data")
	}

	return c.sendWithRetry(ctx, "SubmitEncryptedPartialDecryption", c.cdrContractAddr, callData, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.cdrContract.SubmitEncryptedPartialDecryption(auth, round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, uuid, signature)
	})
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
func (c *ContractClient) estimateGasWithBuffer(ctx context.Context, to common.Address, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From: c.fromAddress,
		To:   &to,
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
	to common.Address,
	callData []byte,
	sendTx func(auth *bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, error) {
	var (
		receipt  *types.Receipt
		gasLimit uint64
		err      error
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		gasLimit, err = c.estimateGasWithBuffer(ctx, to, callData)
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

		log.Info(ctx, methodName+" tx sent",
			"tx_hash", tx.Hash().Hex(),
			"attempt", attempt,
			"gas_limit", gasLimit,
		)

		receipt, err = c.waitForTransaction(ctx, tx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wait for tx", "method_name", methodName)
		}

		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Info(ctx, methodName+" succeeded",
				"tx_hash", tx.Hash().Hex(),
				"gas_used", receipt.GasUsed,
				"attempt", attempt)

			return receipt, nil
		}

		usageRatio := float64(receipt.GasUsed) / float64(gasLimit)
		if usageRatio > 0.95 && attempt < maxRetries {
			log.Warn(ctx, methodName+" likely out-of-gas, retrying",
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
