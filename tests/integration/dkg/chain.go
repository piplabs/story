//go:build integration

package dkg

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/contracts/bindings"
)

// ChainClient 提供链上写操作，用于造景（如 DKG.scheduleUpgrade、后续 CDR 调用等）。
// 未配置时为 nil，仅读用例不受影响。
type ChainClient interface {
	// ScheduleDKGUpgrade 调用 DKG 合约 scheduleUpgrade(activationHeight, upgradeVersion)，用于 IT-BB-03 / upgrade_scheduled 等场景。
	ScheduleDKGUpgrade(ctx context.Context, activationHeight int64, upgradeVersion string) error
	// CancelDKGUpgrade 调用 DKG 合约 cancelUpgrade(upgradeVersion)。
	CancelDKGUpgrade(ctx context.Context, upgradeVersion string) error
	// BlockNumber 返回当前执行层区块高度（用于判断是否 >= activationHeight）。
	BlockNumber(ctx context.Context) (uint64, error)
	// CDRAllocate 调用 CDR 合约 allocate() 创建新 Vault，返回 uuid。
	CDRAllocate(ctx context.Context) (uint32, error)
	// CDRWrite 调用 CDR 合约 write(uuid, accessAuxData, encryptedData)。
	CDRWrite(ctx context.Context, uuid uint32, encryptedData []byte) error
	// CDRRead 调用 CDR 合约 read(uuid, accessAuxData, requesterPubKey)。
	CDRRead(ctx context.Context, uuid uint32, requesterPubKey []byte) error
	// CDRWriteFee 返回当前 CDR write fee。
	CDRWriteFee(ctx context.Context) (*big.Int, error)
	// CDRReadFee 返回当前 CDR read fee。
	CDRReadFee(ctx context.Context) (*big.Int, error)
	// CDRAllocateFee 返回当前 CDR allocate fee。
	CDRAllocateFee(ctx context.Context) (*big.Int, error)
	// RegisterDKGWithParams 调用 DKG.register 并返回错误（用于测试无效参数）。
	RegisterDKGWithParams(ctx context.Context, round uint32, validatorAddr common.Address, enclaveType [32]byte, commKey, dkgPubKey, enclaveReport []byte, startBlockHeight *big.Int, startBlockHash [32]byte) error
	// FinalizeDKGWithParams 调用 DKG.finalize 并返回错误（用于测试无效参数）。
	FinalizeDKGWithParams(ctx context.Context, round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) error
	// GetEthLogs 按 topic 查询 EL logs（用于验证事件）。
	GetEthLogs(ctx context.Context, contractAddr common.Address, topics [][]common.Hash, fromBlock, toBlock uint64) ([]ethtypes.Log, error)
	// SubmitPartialDecryption 直接调用 CDR.submitEncryptedPartialDecryption（用于对抗测试）。
	SubmitPartialDecryption(ctx context.Context, round, pid, uuid uint32, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, signature []byte) error
	// CDRBaseFee 返回当前 CDR base fee。
	CDRBaseFee(ctx context.Context) (*big.Int, error)
}

// NoopChainClient 不执行任何链上写操作。
type NoopChainClient struct{}

func (NoopChainClient) ScheduleDKGUpgrade(context.Context, int64, string) error { return nil }
func (NoopChainClient) CancelDKGUpgrade(context.Context, string) error          { return nil }
func (NoopChainClient) BlockNumber(context.Context) (uint64, error)             { return 0, nil }
func (NoopChainClient) CDRAllocate(context.Context) (uint32, error)             { return 0, nil }
func (NoopChainClient) CDRWrite(context.Context, uint32, []byte) error          { return nil }
func (NoopChainClient) CDRRead(context.Context, uint32, []byte) error           { return nil }
func (NoopChainClient) CDRWriteFee(context.Context) (*big.Int, error)           { return nil, nil }
func (NoopChainClient) CDRReadFee(context.Context) (*big.Int, error)            { return nil, nil }
func (NoopChainClient) CDRAllocateFee(context.Context) (*big.Int, error)        { return nil, nil }
func (NoopChainClient) RegisterDKGWithParams(context.Context, uint32, common.Address, [32]byte, []byte, []byte, []byte, *big.Int, [32]byte) error {
	return nil
}
func (NoopChainClient) FinalizeDKGWithParams(context.Context, uint32, common.Address, [32]byte, [32]byte, []byte, [][]byte, []byte, []byte) error {
	return nil
}
func (NoopChainClient) GetEthLogs(context.Context, common.Address, [][]common.Hash, uint64, uint64) ([]ethtypes.Log, error) {
	return nil, nil
}
func (NoopChainClient) SubmitPartialDecryption(context.Context, uint32, uint32, uint32, []byte, []byte, []byte, []byte, []byte, []byte) error {
	return nil
}
func (NoopChainClient) CDRBaseFee(context.Context) (*big.Int, error) { return nil, nil }

// EthChainClient 使用执行层 JSON-RPC（如 8545）和 DKG/CDR bindings 发送交易。
type EthChainClient struct {
	client        *ethclient.Client
	dkg           *bindings.DKGTransactor
	dkgCaller     *bindings.DKGCaller
	cdr           *bindings.CDRTransactor
	cdrCaller     *bindings.CDRCaller
	auth          *bind.TransactOpts
	dkgAddr       common.Address // DKG contract address
	conditionAddr common.Address // Always-true CDR condition contract
}

// NewEthChainClient 从环境变量创建链上客户端（未配置时返回 nil）：
//   - STORY_ETH_RPC_URL：执行层 RPC，如 http://127.0.0.1:8545
//   - DKG_CONTRACT_ADDRESS：DKG 合约地址，空则用 predeploys.DKG
//   - DKG_SIGNER_PRIVATE_KEY：调用 DKG 的私钥（hex，无 0x 前缀），须为合约 owner；空则返回 nil
func NewEthChainClient() (ChainClient, error) {
	rpcURL := strings.TrimSpace(os.Getenv("STORY_ETH_RPC_URL"))
	if rpcURL == "" {
		return nil, nil
	}
	keyHex := strings.TrimSpace(strings.TrimPrefix(os.Getenv("DKG_SIGNER_PRIVATE_KEY"), "0x"))
	if keyHex == "" {
		return nil, nil
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	addr := os.Getenv("DKG_CONTRACT_ADDRESS")
	if addr == "" {
		addr = predeploys.DKG
	}
	contractAddr := common.HexToAddress(addr)

	dkg, err := bindings.NewDKGTransactor(contractAddr, client)
	if err != nil {
		client.Close()
		return nil, err
	}
	dkgCaller, err := bindings.NewDKGCaller(contractAddr, client)
	if err != nil {
		client.Close()
		return nil, err
	}

	cdrAddr := os.Getenv("CDR_CONTRACT_ADDRESS")
	if cdrAddr == "" {
		cdrAddr = predeploys.CDR
	}
	cdrContractAddr := common.HexToAddress(cdrAddr)
	cdrTransactor, err := bindings.NewCDRTransactor(cdrContractAddr, client)
	if err != nil {
		client.Close()
		return nil, err
	}
	cdrCaller, err := bindings.NewCDRCaller(cdrContractAddr, client)
	if err != nil {
		client.Close()
		return nil, err
	}

	key, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		client.Close()
		return nil, err
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		client.Close()
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		client.Close()
		return nil, err
	}

	// Set gas price for legacy transactions (devnet geth miner minimum is 16 gwei).
	// Default to 17 gwei if DKG_GAS_PRICE is not set.
	gasPriceStr := strings.TrimSpace(os.Getenv("DKG_GAS_PRICE"))
	if gasPriceStr == "" {
		gasPriceStr = "17000000000" // 17 gwei
	}
	gasPrice, err := strconv.ParseInt(gasPriceStr, 10, 64)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("invalid DKG_GAS_PRICE=%q: %w", gasPriceStr, err)
	}
	auth.GasPrice = big.NewInt(gasPrice)

	ec := &EthChainClient{client: client, dkg: dkg, dkgCaller: dkgCaller, cdr: cdrTransactor, cdrCaller: cdrCaller, auth: auth, dkgAddr: contractAddr}

	// Clear any stuck pending transactions (e.g. from DCAP deploy with low gas).
	if clearErr := clearStuckPendingTxs(client, auth, key, chainID); clearErr != nil {
		fmt.Fprintf(os.Stderr, "WARN: clearStuckPendingTxs: %v\n", clearErr)
	}

	// Deploy always-true CDR condition contract for CDR tests.
	condAddr, deployErr := deployAlwaysTrueCondition(client, auth)
	if deployErr == nil {
		ec.conditionAddr = condAddr
	} else {
		// Fallback: use zero address (CDR calls will revert but won't block).
		fmt.Fprintf(os.Stderr, "WARN: deploy condition contract failed: %v (CDR tests will skip)\n", deployErr)
		ec.conditionAddr = contractAddr
	}

	return ec, nil
}

// clearStuckPendingTxs replaces any stuck pending transactions with high-gas self-transfers.
// This prevents DCAP deploy leftovers from blocking CDR test transactions.
func clearStuckPendingTxs(client *ethclient.Client, auth *bind.TransactOpts, key *ecdsa.PrivateKey, chainID *big.Int) error {
	ctx := context.Background()
	latest, err := client.NonceAt(ctx, auth.From, nil)
	if err != nil {
		return err
	}
	pending, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		return err
	}
	if pending <= latest {
		return nil
	}
	fmt.Fprintf(os.Stderr, "[CDR] Clearing %d stuck pending txs (nonce %d..%d)\n", pending-latest, latest, pending-1)
	// Use raw crypto.SignTx with legacy tx (high gas) — works for replacing any tx type
	signer := ethtypes.LatestSignerForChainID(chainID)
	for nonce := latest; nonce < pending; nonce++ {
		// Try EIP-1559 with raw key
		eip1559Tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
			ChainID: chainID, Nonce: nonce, To: &auth.From, Value: big.NewInt(0), Gas: 21000,
			GasTipCap: new(big.Int).Mul(big.NewInt(500), big.NewInt(1e9)),
			GasFeeCap: new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e9)),
		})
		signed, sErr := ethtypes.SignTx(eip1559Tx, signer, key)
		if sErr == nil {
			sErr = client.SendTransaction(ctx, signed)
		}
		if sErr != nil {
			// Fallback: legacy tx with raw key
			legacyTx := ethtypes.NewTx(&ethtypes.LegacyTx{
				Nonce: nonce, To: &auth.From, Value: big.NewInt(0), Gas: 21000,
				GasPrice: new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e9)),
			})
			signed, sErr = ethtypes.SignTx(legacyTx, signer, key)
			if sErr == nil {
				sErr = client.SendTransaction(ctx, signed)
			}
			if sErr != nil {
				fmt.Fprintf(os.Stderr, "[CDR] clear nonce %d: %v\n", nonce, sErr)
				continue
			}
		}
		waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		_, _ = bind.WaitMined(waitCtx, client, signed)
		cancel()
	}
	fmt.Fprintf(os.Stderr, "[CDR] Stuck txs cleared\n")
	return nil
}

// waitNonceStable waits until pending nonce == latest nonce (no pending txs from our account).
// This prevents "replacement transaction underpriced" when sending sequential transactions.
func (c *EthChainClient) waitNonceStable(ctx context.Context) {
	for i := 0; i < 30; i++ {
		pending, err1 := c.client.PendingNonceAt(ctx, c.auth.From)
		latest, err2 := c.client.NonceAt(ctx, c.auth.From, nil)
		if err1 == nil && err2 == nil && pending == latest {
			return
		}
		time.Sleep(time.Second)
	}
}

// sendAndWait sends a transaction via bind and waits for it to be mined.
// If WaitMined fails (e.g. tx reverts and stays in pending pool), it replaces
// the stuck nonce with a high-gas self-transfer to unblock subsequent transactions.
func (c *EthChainClient) sendAndWait(ctx context.Context, tx *ethtypes.Transaction, label string) (*ethtypes.Receipt, error) {
	fmt.Fprintf(os.Stderr, "[CDR] %s: tx=%s, waiting for mine...\n", label, tx.Hash().Hex())
	waitCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(waitCtx, c.client, tx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CDR] %s: WaitMined failed: %v — clearing stuck nonce\n", label, err)
		// Replace stuck tx: use TransactOpts with forced nonce + high gas tip
		nonce := tx.Nonce()
		replaceOpts := *c.auth
		replaceOpts.Nonce = new(big.Int).SetUint64(nonce)
		replaceOpts.GasLimit = 21000
		replaceOpts.Value = big.NewInt(0)
		replaceOpts.GasTipCap = new(big.Int).Mul(big.NewInt(500), big.NewInt(1e9)) // 500 gwei tip
		replaceOpts.GasFeeCap = new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e9)) // 1000 gwei max
		replaceOpts.GasPrice = nil // force EIP-1559
		replaceTx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
			Nonce:     nonce,
			To:        &c.auth.From,
			Value:     big.NewInt(0),
			Gas:       21000,
			GasTipCap: replaceOpts.GasTipCap,
			GasFeeCap: replaceOpts.GasFeeCap,
		})
		signed, sErr := c.auth.Signer(c.auth.From, replaceTx)
		if sErr == nil {
			_ = c.client.SendTransaction(context.Background(), signed)
			replaceCtx, replaceCancel := context.WithTimeout(context.Background(), 30*time.Second)
			_, _ = bind.WaitMined(replaceCtx, c.client, signed)
			replaceCancel()
			fmt.Fprintf(os.Stderr, "[CDR] %s: stuck nonce %d cleared\n", label, nonce)
		}
		return nil, fmt.Errorf("%s WaitMined: %w", label, err)
	}
	fmt.Fprintf(os.Stderr, "[CDR] %s: mined, status=%d\n", label, receipt.Status)
	if receipt.Status == 0 {
		return receipt, fmt.Errorf("%s reverted (tx=%s)", label, tx.Hash().Hex())
	}
	return receipt, nil
}

// deployAlwaysTrueCondition deploys a minimal contract that returns true for any call.
// Used as CDR write/read condition in tests.
func deployAlwaysTrueCondition(client *ethclient.Client, auth *bind.TransactOpts) (common.Address, error) {
	// Runtime: PUSH1 1, PUSH1 0, MSTORE, PUSH1 0x20, PUSH1 0, RETURN (returns uint256(1) = true)
	runtime := []byte{0x60, 0x01, 0x60, 0x00, 0x52, 0x60, 0x20, 0x60, 0x00, 0xf3}
	// Init: copies runtime to memory and returns it
	initLen := byte(12) // length of init code before runtime
	initCode := []byte{
		0x60, byte(len(runtime)), // PUSH1 runtimeLen
		0x60, initLen,            // PUSH1 initLen (offset where runtime starts)
		0x60, 0x00, // PUSH1 0 (dest in memory)
		0x39,                     // CODECOPY
		0x60, byte(len(runtime)), // PUSH1 runtimeLen
		0x60, 0x00, // PUSH1 0
		0xf3, // RETURN
	}
	bytecode := append(initCode, runtime...)

	// Use auth transactor for proper nonce management (avoid conflicts with CDR calls).
	opts := *auth
	opts.GasLimit = 200000
	opts.Value = big.NewInt(0)

	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return common.Address{}, fmt.Errorf("get nonce: %w", err)
	}
	highGasPrice := new(big.Int).Mul(auth.GasPrice, big.NewInt(3)) // 3x gas price to replace any pending tx
	fmt.Fprintf(os.Stderr, "[CDR] Deploy condition: nonce=%d gasPrice=%s\n", nonce, highGasPrice)
	tx := ethtypes.NewContractCreation(nonce, big.NewInt(0), 200000, highGasPrice, bytecode)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("sign tx: %w", err)
	}
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Fprintf(os.Stderr, "[CDR] Deploy condition: send failed: %v\n", err)
		return common.Address{}, fmt.Errorf("send tx: %w", err)
	}
	fmt.Fprintf(os.Stderr, "[CDR] Deploy condition: tx=%s, waiting for mine...\n", signedTx.Hash().Hex())

	waitCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(waitCtx, client, signedTx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CDR] Deploy condition: WaitMined failed: %v\n", err)
		return common.Address{}, fmt.Errorf("wait mine: %w", err)
	}
	fmt.Fprintf(os.Stderr, "[CDR] Deploy condition: deployed at %s\n", receipt.ContractAddress.Hex())
	if receipt.Status == 0 {
		return common.Address{}, fmt.Errorf("deploy tx reverted")
	}

	return receipt.ContractAddress, nil
}

func (c *EthChainClient) ScheduleDKGUpgrade(ctx context.Context, activationHeight int64, upgradeVersion string) error {
	opts := *c.auth
	opts.Context = ctx
	_, err := c.dkg.ScheduleUpgrade(&opts, big.NewInt(activationHeight), upgradeVersion)
	return err
}

func (c *EthChainClient) CancelDKGUpgrade(ctx context.Context, upgradeVersion string) error {
	opts := *c.auth
	opts.Context = ctx
	_, err := c.dkg.CancelUpgrade(&opts, upgradeVersion)
	return err
}

func (c *EthChainClient) BlockNumber(ctx context.Context) (uint64, error) {
	bn, err := c.client.BlockNumber(ctx)
	return bn, err
}

func (c *EthChainClient) CDRAllocate(ctx context.Context) (uint32, error) {
	c.waitNonceStable(ctx)
	fmt.Fprintf(os.Stderr, "[CDR] Allocate: querying fee...\n")
	fee, err := c.cdrCaller.AllocateFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}
	fmt.Fprintf(os.Stderr, "[CDR] Allocate: fee=%s, sending tx...\n", fee)
	opts := *c.auth
	opts.Context = ctx
	opts.Value = fee
	// Use always-true condition contract for both write and read conditions.
	// updatable=true so Write works even if vault has stale data from a prior devnet cycle.
	tx, err := c.cdr.Allocate(&opts, true, c.conditionAddr, c.conditionAddr, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CDR] Allocate: send failed: %v\n", err)
		return 0, err
	}
	receipt, err := c.sendAndWait(ctx, tx, "Allocate")
	if err != nil {
		return 0, err
	}
	for _, log := range receipt.Logs {
		ev, parseErr := bindings.NewCDRFilterer(common.Address{}, nil)
		if parseErr != nil {
			continue
		}
		allocated, parseErr := ev.ParseVaultAllocated(*log)
		if parseErr == nil {
			fmt.Fprintf(os.Stderr, "[CDR] Allocate: uuid=%d\n", allocated.Uuid)
			return allocated.Uuid, nil
		}
	}
	return 0, fmt.Errorf("VaultAllocated event not found in receipt")
}

func (c *EthChainClient) CDRWrite(ctx context.Context, uuid uint32, encryptedData []byte) error {
	c.waitNonceStable(ctx)
	fmt.Fprintf(os.Stderr, "[CDR] Write: uuid=%d, querying fee...\n", uuid)
	fee, err := c.cdrCaller.WriteFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "[CDR] Write: fee=%s, sending tx...\n", fee)
	opts := *c.auth
	opts.Context = ctx
	opts.Value = fee
	tx, err := c.cdr.Write(&opts, uuid, nil, encryptedData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CDR] Write: send failed: %v\n", err)
		return err
	}
	_, err = c.sendAndWait(ctx, tx, "Write")
	return err
}

func (c *EthChainClient) CDRRead(ctx context.Context, uuid uint32, requesterPubKey []byte) error {
	c.waitNonceStable(ctx)
	fmt.Fprintf(os.Stderr, "[CDR] Read: uuid=%d, querying fee...\n", uuid)
	fee, err := c.cdrCaller.ReadFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "[CDR] Read: fee=%s, sending tx...\n", fee)
	opts := *c.auth
	opts.Context = ctx
	opts.Value = fee
	tx, err := c.cdr.Read(&opts, uuid, nil, requesterPubKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CDR] Read: send failed: %v\n", err)
		return err
	}
	_, err = c.sendAndWait(ctx, tx, "Read")
	if err != nil {
		return err
	}
	return nil
}

func (c *EthChainClient) CDRWriteFee(ctx context.Context) (*big.Int, error) {
	return c.cdrCaller.WriteFee(&bind.CallOpts{Context: ctx})
}

func (c *EthChainClient) CDRReadFee(ctx context.Context) (*big.Int, error) {
	return c.cdrCaller.ReadFee(&bind.CallOpts{Context: ctx})
}

func (c *EthChainClient) CDRAllocateFee(ctx context.Context) (*big.Int, error) {
	return c.cdrCaller.AllocateFee(&bind.CallOpts{Context: ctx})
}

func (c *EthChainClient) RegisterDKGWithParams(ctx context.Context, round uint32, validatorAddr common.Address, enclaveType [32]byte, commKey, dkgPubKey, enclaveReport []byte, startBlockHeight *big.Int, startBlockHash [32]byte) error {
	c.waitNonceStable(ctx)
	opts := *c.auth
	opts.Context = ctx

	// Query the registration fee and attach it so the call isn't rejected by fee check.
	fee, err := c.dkgCaller.Fee(&bind.CallOpts{Context: ctx})
	if err == nil && fee != nil {
		opts.Value = fee
	}

	instanceData := bindings.IDKGEnclaveInstanceData{
		Round:           round,
		ValidatorAddr:   validatorAddr,
		EnclaveType:     enclaveType,
		EnclaveCommKey:  commKey,
		DkgPubKey:       dkgPubKey,
	}
	tx, txErr := c.dkg.Register(&opts, enclaveReport, instanceData, startBlockHeight, startBlockHash, nil)
	if txErr != nil {
		return txErr
	}
	_, err = c.sendAndWait(ctx, tx, "Register")
	return err
}

func (c *EthChainClient) FinalizeDKGWithParams(ctx context.Context, round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) error {
	c.waitNonceStable(ctx)
	opts := *c.auth
	opts.Context = ctx
	tx, err := c.dkg.Finalize(&opts, round, validatorAddr, enclaveType, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature)
	if err != nil {
		return err
	}
	_, err = c.sendAndWait(ctx, tx, "Finalize")
	return err
}

func (c *EthChainClient) GetEthLogs(ctx context.Context, contractAddr common.Address, topics [][]common.Hash, fromBlock, toBlock uint64) ([]ethtypes.Log, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics:    topics,
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
	}
	return c.client.FilterLogs(ctx, query)
}

func (c *EthChainClient) SubmitPartialDecryption(ctx context.Context, round, pid, uuid uint32, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, signature []byte) error {
	c.waitNonceStable(ctx)
	baseFee, err := c.cdrCaller.BaseFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("query baseFee: %w", err)
	}
	opts := *c.auth
	opts.Context = ctx
	opts.Value = baseFee
	tx, txErr := c.cdr.SubmitEncryptedPartialDecryption(&opts, round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, uuid, signature)
	if txErr != nil {
		return txErr
	}
	_, err = c.sendAndWait(ctx, tx, "SubmitPartial")
	return err
}

func (c *EthChainClient) CDRBaseFee(ctx context.Context) (*big.Int, error) {
	return c.cdrCaller.BaseFee(&bind.CallOpts{Context: ctx})
}
