//go:build integration

package dkg

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

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

	// Deploy always-true CDR condition contract.
	// Runtime bytecode: returns true (uint256(1)) for any call.
	// PUSH1 1, PUSH1 0, MSTORE, PUSH1 32, PUSH1 0, RETURN
	condAddr, deployErr := deployAlwaysTrueCondition(client, auth)
	if deployErr == nil {
		ec.conditionAddr = condAddr
	} else {
		// Fallback: use DKG contract address
		ec.conditionAddr = contractAddr
	}

	return ec, nil
}

// deployAlwaysTrueCondition deploys a minimal contract that returns true for any call.
// Used as CDR write/read condition in tests.
func deployAlwaysTrueCondition(client *ethclient.Client, auth *bind.TransactOpts) (common.Address, error) {
	// Runtime: PUSH1 1, PUSH1 0, MSTORE, PUSH1 0x20, PUSH1 0, RETURN (returns uint256(1) = true)
	runtime := []byte{0x60, 0x01, 0x60, 0x00, 0x52, 0x60, 0x20, 0x60, 0x00, 0xf3}
	// Init: copies runtime to memory and returns it
	initLen := byte(12) // length of init code before runtime
	init := []byte{
		0x60, byte(len(runtime)), // PUSH1 runtimeLen
		0x60, initLen,            // PUSH1 initLen (offset where runtime starts)
		0x60, 0x00, // PUSH1 0 (dest in memory)
		0x39,                     // CODECOPY
		0x60, byte(len(runtime)), // PUSH1 runtimeLen
		0x60, 0x00, // PUSH1 0
		0xf3, // RETURN
	}
	bytecode := append(init, runtime...)

	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return common.Address{}, err
	}

	tx := ethtypes.NewContractCreation(nonce, big.NewInt(0), 200000, auth.GasPrice, bytecode)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return common.Address{}, err
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return common.Address{}, err
	}

	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		return common.Address{}, err
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
	fee, err := c.cdrCaller.AllocateFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}
	opts := *c.auth
	opts.Context = ctx
	opts.Value = fee
	// Use always-true condition contract for both write and read conditions
	tx, err := c.cdr.Allocate(&opts, false, c.conditionAddr, c.conditionAddr, nil, nil)
	if err != nil {
		return 0, err
	}
	receipt, err := bind.WaitMined(ctx, c.client, tx)
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
			return allocated.Uuid, nil
		}
	}
	return 0, fmt.Errorf("VaultAllocated event not found in receipt")
}

func (c *EthChainClient) CDRWrite(ctx context.Context, uuid uint32, encryptedData []byte) error {
	fee, err := c.cdrCaller.WriteFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	opts := *c.auth
	opts.Context = ctx
	opts.Value = fee
	_, err = c.cdr.Write(&opts, uuid, nil, encryptedData)
	return err
}

func (c *EthChainClient) CDRRead(ctx context.Context, uuid uint32, requesterPubKey []byte) error {
	fee, err := c.cdrCaller.ReadFee(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	opts := *c.auth
	opts.Context = ctx
	opts.Value = fee
	_, err = c.cdr.Read(&opts, uuid, nil, requesterPubKey)
	return err
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
	_, err = c.dkg.Register(&opts, enclaveReport, instanceData, startBlockHeight, startBlockHash, nil)
	return err
}

func (c *EthChainClient) FinalizeDKGWithParams(ctx context.Context, round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) error {
	opts := *c.auth
	opts.Context = ctx
	_, err := c.dkg.Finalize(&opts, round, validatorAddr, enclaveType, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature)
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
