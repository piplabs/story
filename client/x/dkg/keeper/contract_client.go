package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// EthClient abstracts the Ethereum JSON-RPC methods used by ContractClient.
// This allows test code to inject a mock implementation.
type EthClient interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	// TransactionReceipt and CodeAt satisfy bind.DeployBackend so that
	// bind.WaitMined can be called with this interface directly.
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	BlockNumber(ctx context.Context) (uint64, error)
}

// DKGContractBinding abstracts the DKG smart-contract methods used by ContractClient.
type DKGContractBinding interface {
	Fee(opts *bind.CallOpts) (*big.Int, error)
	Register(opts *bind.TransactOpts, enclaveReport []byte, enclaveInstanceData bindings.IDKGEnclaveInstanceData, startBlockHeight *big.Int, startBlockHash [32]byte, validationContext []byte) (*types.Transaction, error)
	Finalize(opts *bind.TransactOpts, round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) (*types.Transaction, error)
}

// CDRContractBinding abstracts the CDR smart-contract methods used by ContractClient.
type CDRContractBinding interface {
	BaseFee(opts *bind.CallOpts) (*big.Int, error)
	SubmitEncryptedPartialDecryption(opts *bind.TransactOpts, round uint32, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, requesterPubKey []byte, ciphertext []byte, uuid uint32, signature []byte) (*types.Transaction, error)
	SubmitEncryptedPartialDecryptionBatch(opts *bind.TransactOpts, requests []bindings.ICDRPartialDecryptionRequest) (*types.Transaction, error)
}

// Compile-time assertions.
var (
	_ dkgtypes.DKGContractClient = (*ContractClient)(nil)
	_ EthClient                  = (*ethclient.Client)(nil)
	_ DKGContractBinding         = (*bindings.DKG)(nil)
	_ CDRContractBinding         = (*bindings.CDR)(nil)
)

const (
	maxRetries = 3
)

// ContractClient wraps the DKG contract interaction.
type ContractClient struct {
	ethClient       EthClient
	dkgContract     DKGContractBinding
	dkgContractAbi  *abi.ABI
	dkgContractAddr common.Address
	cdrContract     CDRContractBinding
	cdrContractAbi  *abi.ABI
	cdrContractAddr common.Address
	privateKey      *ecdsa.PrivateKey
	fromAddress     common.Address
	chainID         *big.Int

	// pendingTxs tracks the last-sent tx for each operation key (e.g. "register_1").
	// If waitForTransaction times out, the tx may still be in the mempool.
	// The next invocation of the same operation checks this map first to avoid
	// sending a new tx with a higher nonce while the previous one is still pending.
	pendingTxMu sync.Mutex
	pendingTxs  map[string]*types.Transaction
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
		pendingTxs:      make(map[string]*types.Transaction),
	}

	log.Info(ctx, "Created contract client",
		"dkg_contract_address", dkgContractAddr.Hex(),
		"cdr_contract_address", cdrContractAddr.Hex(),
		"from_address", fromAddress.Hex(),
		"chain_id", chainID,
	)

	return client, nil
}

// BlockNumber returns the current block number from the EL client.
func (c *ContractClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
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

	// Query the registration fee from the DKG contract
	fee, err := c.dkgContract.Fee(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query DKG fee")
	}

	log.Info(ctx, "DKG registration fee queried", "fee_wei", fee.String())

	pendingKey := fmt.Sprintf("register_%d", round)
	if receipt, err := c.checkAndResumePendingTx(ctx, pendingKey); receipt != nil || err != nil {
		return receipt, err
	}

	receipt, err := c.sendWithRetry(ctx, "Register", c.dkgContractAddr, callData, fee, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx, txErr := c.dkgContract.Register(auth, enclaveReport, enclaveInstanceData, startBlockHeightBig, startBlockHash32, []byte{})
		if txErr == nil {
			c.storePendingTx(pendingKey, tx)
		}
		return tx, txErr
	})
	if err == nil {
		c.clearPendingTx(pendingKey)
	}
	return receipt, err
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

	// Query the DKG fee — the contract requires it for finalization as well.
	fee, err := c.dkgContract.Fee(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query DKG fee for finalize")
	}

	log.Info(ctx, "DKG finalization fee queried", "fee_wei", fee.String())

	pendingKey := fmt.Sprintf("finalize_%d", round)
	if receipt, err := c.checkAndResumePendingTx(ctx, pendingKey); receipt != nil || err != nil {
		return receipt, err
	}

	receipt, err := c.sendWithRetry(ctx, "Finalize", c.dkgContractAddr, callData, fee, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx, txErr := c.dkgContract.Finalize(auth, round, c.fromAddress, enclaveType, participantsRoot32, globalPubKey, publicCoeffs, pubKeyShare, signature)
		if txErr == nil {
			c.storePendingTx(pendingKey, tx)
		}
		return tx, txErr
	})
	if err == nil {
		c.clearPendingTx(pendingKey)
	}
	return receipt, err
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
	ciphertext []byte,
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
		"ciphertext_len", len(ciphertext),
		"uuid", uuid,
		"signature_len", len(signature),
	)

	callData, err := c.cdrContractAbi.Pack("submitEncryptedPartialDecryption", round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, uuid, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack submitEncryptedPartialDecryption call data")
	}

	fee, err := c.cdrContract.BaseFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query CDR base fee for partial decryption")
	}

	log.Info(ctx, "CDR base fee queried for partial decryption", "fee_wei", fee.String())

	return c.sendWithRetry(ctx, "SubmitEncryptedPartialDecryption", c.cdrContractAddr, callData, fee, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.cdrContract.SubmitEncryptedPartialDecryption(auth, round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, uuid, signature)
	})
}

// SubmitEncryptedPartialDecryptionBatch calls the batch submitEncryptedPartialDecryptionBatch
// contract method, submitting multiple partial decryptions in a single transaction.
func (c *ContractClient) SubmitEncryptedPartialDecryptionBatch(
	ctx context.Context,
	requests []bindings.ICDRPartialDecryptionRequest,
) (*types.Receipt, error) {
	log.Info(ctx, "Calling submitEncryptedPartialDecryptionBatch contract method",
		"batch_size", len(requests),
	)

	callData, err := c.cdrContractAbi.Pack("submitEncryptedPartialDecryptionBatch", requests)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack submitEncryptedPartialDecryptionBatch call data")
	}

	baseFee, err := c.cdrContract.BaseFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query CDR base fee for batch partial decryption")
	}

	totalFee := new(big.Int).Mul(baseFee, big.NewInt(int64(len(requests))))

	log.Info(ctx, "CDR batch fee computed", "base_fee_wei", baseFee.String(), "total_fee_wei", totalFee.String(), "batch_size", len(requests))

	return c.sendWithRetry(ctx, "SubmitEncryptedPartialDecryptionBatch", c.cdrContractAddr, callData, totalFee, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.cdrContract.SubmitEncryptedPartialDecryptionBatch(auth, requests)
	})
}

// createTransactOpts creates transaction options for contract calls.
// If value is nil, it defaults to 0.
func (c *ContractClient) createTransactOpts(ctx context.Context, gasLimit uint64, value *big.Int) (*bind.TransactOpts, error) {
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

	if value == nil {
		value = big.NewInt(0)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	auth.Context = ctx

	return auth, nil
}

// estimateGasWithBuffer estimates gas for a contract transaction and adds a safety buffer.
func (c *ContractClient) estimateGasWithBuffer(ctx context.Context, to common.Address, data []byte, value *big.Int) (uint64, error) {
	msg := ethereum.CallMsg{
		From:  c.fromAddress,
		To:    &to,
		Data:  data,
		Value: value,
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
	value *big.Int,
	sendTx func(auth *bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, error) {
	var (
		receipt  *types.Receipt
		gasLimit uint64
		err      error
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		gasLimit, err = c.estimateGasWithBuffer(ctx, to, callData, value)
		if err != nil {
			return nil, err
		}

		auth, err := c.createTransactOpts(ctx, gasLimit, value)
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

func (c *ContractClient) storePendingTx(key string, tx *types.Transaction) {
	c.pendingTxMu.Lock()
	c.pendingTxs[key] = tx
	c.pendingTxMu.Unlock()
}

func (c *ContractClient) clearPendingTx(key string) {
	c.pendingTxMu.Lock()
	delete(c.pendingTxs, key)
	c.pendingTxMu.Unlock()
}

// checkAndResumePendingTx looks up a previously sent tx that may still be in the
// mempool after a waitForTransaction timeout. It returns:
//   - (receipt, nil)  if the tx was mined successfully — caller should return immediately
//   - (nil, err)      if the tx is still pending — caller should abort to avoid a nonce gap
//   - (nil, nil)      if no pending tx exists or it was mined but failed — caller should proceed normally
func (c *ContractClient) checkAndResumePendingTx(ctx context.Context, key string) (*types.Receipt, error) {
	c.pendingTxMu.Lock()
	tx, ok := c.pendingTxs[key]
	c.pendingTxMu.Unlock()
	if !ok {
		return nil, nil
	}

	receipt, _ := c.ethClient.TransactionReceipt(ctx, tx.Hash())
	if receipt == nil {
		return nil, errors.New("previous tx still pending, skipping to avoid nonce gap",
			"key", key, "tx_hash", tx.Hash().Hex())
	}

	c.clearPendingTx(key)
	if receipt.Status == types.ReceiptStatusSuccessful {
		return receipt, nil
	}
	return nil, nil
}
