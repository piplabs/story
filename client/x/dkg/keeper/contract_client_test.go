package keeper

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/contracts/bindings"
)

// mockEthClient implements EthClient for testing.
type mockEthClient struct {
	estimateGasFn        func(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	suggestGasPriceFn    func(ctx context.Context) (*big.Int, error)
	pendingNonceAtFn     func(ctx context.Context, account common.Address) (uint64, error)
	transactionReceiptFn func(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	codeAtFn             func(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
}

func (m *mockEthClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	if m.estimateGasFn != nil {
		return m.estimateGasFn(ctx, msg)
	}

	return 100000, nil
}

func (m *mockEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if m.suggestGasPriceFn != nil {
		return m.suggestGasPriceFn(ctx)
	}

	return big.NewInt(1000000000), nil // 1 Gwei
}

func (m *mockEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if m.pendingNonceAtFn != nil {
		return m.pendingNonceAtFn(ctx, account)
	}

	return 0, nil
}

func (m *mockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if m.transactionReceiptFn != nil {
		return m.transactionReceiptFn(ctx, txHash)
	}

	return &types.Receipt{Status: types.ReceiptStatusSuccessful}, nil
}

func (m *mockEthClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	if m.codeAtFn != nil {
		return m.codeAtFn(ctx, account, blockNumber)
	}

	return []byte{}, nil
}

func (m *mockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, nil
}

// mockDKGContract implements DKGContractBinding for testing.
type mockDKGContract struct {
	feeFn      func(opts *bind.CallOpts) (*big.Int, error)
	registerFn func(opts *bind.TransactOpts, enclaveReport []byte, enclaveInstanceData bindings.IDKGEnclaveInstanceData, startBlockHeight *big.Int, startBlockHash [32]byte, validationContext []byte) (*types.Transaction, error)
	finalizeFn func(opts *bind.TransactOpts, round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) (*types.Transaction, error)
}

func (m *mockDKGContract) Fee(opts *bind.CallOpts) (*big.Int, error) {
	if m.feeFn != nil {
		return m.feeFn(opts)
	}

	return big.NewInt(0), nil
}

func (m *mockDKGContract) Register(opts *bind.TransactOpts, enclaveReport []byte, enclaveInstanceData bindings.IDKGEnclaveInstanceData, startBlockHeight *big.Int, startBlockHash [32]byte, validationContext []byte) (*types.Transaction, error) {
	if m.registerFn != nil {
		return m.registerFn(opts, enclaveReport, enclaveInstanceData, startBlockHeight, startBlockHash, validationContext)
	}

	return makeTx(opts.Nonce.Uint64()), nil
}

func (m *mockDKGContract) Finalize(opts *bind.TransactOpts, round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) (*types.Transaction, error) {
	if m.finalizeFn != nil {
		return m.finalizeFn(opts, round, validatorAddr, enclaveType, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature)
	}

	return makeTx(opts.Nonce.Uint64()), nil
}

// mockCDRContract implements CDRContractBinding for testing.
type mockCDRContract struct {
	baseFeeFn                          func(opts *bind.CallOpts) (*big.Int, error)
	submitEncryptedPartialDecryptionFn func(opts *bind.TransactOpts, round uint32, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, requesterPubKey []byte, ciphertext []byte, uuid uint32, signature []byte) (*types.Transaction, error)
}

func (m *mockCDRContract) BaseFee(opts *bind.CallOpts) (*big.Int, error) {
	if m.baseFeeFn != nil {
		return m.baseFeeFn(opts)
	}

	return big.NewInt(0), nil
}

func (m *mockCDRContract) SubmitEncryptedPartialDecryption(opts *bind.TransactOpts, round uint32, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, requesterPubKey []byte, ciphertext []byte, uuid uint32, signature []byte) (*types.Transaction, error) {
	if m.submitEncryptedPartialDecryptionFn != nil {
		return m.submitEncryptedPartialDecryptionFn(opts, round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, uuid, signature)
	}

	return makeTx(opts.Nonce.Uint64()), nil
}

// Compile-time assertions.
var (
	_ EthClient          = (*mockEthClient)(nil)
	_ DKGContractBinding = (*mockDKGContract)(nil)
	_ CDRContractBinding = (*mockCDRContract)(nil)
)

// testKey generates a deterministic ECDSA key pair for testing.
func testKey(t *testing.T) (*ecdsa.PrivateKey, common.Address) {
	t.Helper()

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	return privateKey, fromAddress
}

// newTestContractClient creates a ContractClient with mock dependencies.
func newTestContractClient(t *testing.T, mock *mockEthClient) *ContractClient {
	t.Helper()

	pk, addr := testKey(t)

	return &ContractClient{
		ethClient:   mock,
		privateKey:  pk,
		fromAddress: addr,
		chainID:     big.NewInt(1),
	}
}

// newFullTestContractClient creates a ContractClient with all mock dependencies
// including DKG/CDR contract bindings and ABI.
func newFullTestContractClient(t *testing.T, ethMock *mockEthClient, dkgMock *mockDKGContract, cdrMock *mockCDRContract) *ContractClient {
	t.Helper()

	pk, addr := testKey(t)

	dkgAbi, err := bindings.DKGMetaData.GetAbi()
	require.NoError(t, err)

	cdrAbi, err := bindings.CDRMetaData.GetAbi()
	require.NoError(t, err)

	return &ContractClient{
		ethClient:       ethMock,
		dkgContract:     dkgMock,
		dkgContractAbi:  dkgAbi,
		dkgContractAddr: common.HexToAddress("0xaaaa"),
		cdrContract:     cdrMock,
		cdrContractAbi:  cdrAbi,
		cdrContractAddr: common.HexToAddress("0xbbbb"),
		privateKey:      pk,
		fromAddress:     addr,
		chainID:         big.NewInt(1),
	}
}

// makeTx creates a minimal legacy transaction for testing.
func makeTx(nonce uint64) *types.Transaction {
	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(1000000000),
		Gas:      21000,
		To:       &common.Address{},
		Value:    big.NewInt(0),
	})
}

// ---------- estimateGasWithBuffer ----------

func TestEstimateGasWithBuffer_AddsBuffer(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
	}
	client := newTestContractClient(t, mock)

	gas, err := client.estimateGasWithBuffer(context.Background(), common.Address{}, []byte{0x01}, nil)
	require.NoError(t, err)
	// 100000 * 12 / 10 = 120000
	require.Equal(t, uint64(120000), gas)
}

func TestEstimateGasWithBuffer_LargeValue(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 500000, nil
		},
	}
	client := newTestContractClient(t, mock)

	gas, err := client.estimateGasWithBuffer(context.Background(), common.Address{}, []byte{0x01}, big.NewInt(1000))
	require.NoError(t, err)
	require.Equal(t, uint64(600000), gas)
}

func TestEstimateGasWithBuffer_Error(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 0, errSentinel
		},
	}
	client := newTestContractClient(t, mock)

	_, err := client.estimateGasWithBuffer(context.Background(), common.Address{}, []byte{0x01}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to estimate gas")
}

func TestEstimateGasWithBuffer_PassesCallMsg(t *testing.T) {
	t.Parallel()

	to := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	data := []byte{0xaa, 0xbb}
	val := big.NewInt(42)

	var captured ethereum.CallMsg
	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, msg ethereum.CallMsg) (uint64, error) {
			captured = msg

			return 100, nil
		},
	}
	client := newTestContractClient(t, mock)

	_, err := client.estimateGasWithBuffer(context.Background(), to, data, val)
	require.NoError(t, err)
	require.Equal(t, client.fromAddress, captured.From)
	require.Equal(t, &to, captured.To)
	require.Equal(t, data, captured.Data)
	require.Equal(t, val, captured.Value)
}

// ---------- createTransactOpts ----------

func TestCreateTransactOpts_Success(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		pendingNonceAtFn: func(_ context.Context, _ common.Address) (uint64, error) {
			return 42, nil
		},
		suggestGasPriceFn: func(_ context.Context) (*big.Int, error) {
			return big.NewInt(5_000_000_000), nil
		},
	}
	client := newTestContractClient(t, mock)

	auth, err := client.createTransactOpts(context.Background(), 300_000, big.NewInt(1000))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(42), auth.Nonce)
	require.Equal(t, big.NewInt(5_000_000_000), auth.GasPrice)
	require.Equal(t, uint64(300_000), auth.GasLimit)
	require.Equal(t, big.NewInt(1000), auth.Value)
}

func TestCreateTransactOpts_NilValueDefaultsToZero(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{}
	client := newTestContractClient(t, mock)

	auth, err := client.createTransactOpts(context.Background(), 21000, nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), auth.Value)
}

func TestCreateTransactOpts_NonceError(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		pendingNonceAtFn: func(_ context.Context, _ common.Address) (uint64, error) {
			return 0, errSentinel
		},
	}
	client := newTestContractClient(t, mock)

	_, err := client.createTransactOpts(context.Background(), 21000, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get pending nonce")
}

func TestCreateTransactOpts_GasPriceError(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		suggestGasPriceFn: func(_ context.Context) (*big.Int, error) {
			return nil, errSentinel
		},
	}
	client := newTestContractClient(t, mock)

	_, err := client.createTransactOpts(context.Background(), 21000, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get gas price")
}

// ---------- waitForTransaction ----------

func TestWaitForTransaction_Success(t *testing.T) {
	t.Parallel()

	tx := makeTx(0)
	expectedReceipt := &types.Receipt{Status: types.ReceiptStatusSuccessful, GasUsed: 21000}

	mock := &mockEthClient{
		transactionReceiptFn: func(_ context.Context, hash common.Hash) (*types.Receipt, error) {
			require.Equal(t, tx.Hash(), hash)

			return expectedReceipt, nil
		},
	}
	client := newTestContractClient(t, mock)

	receipt, err := client.waitForTransaction(context.Background(), tx)
	require.NoError(t, err)
	require.Equal(t, expectedReceipt, receipt)
}

func TestWaitForTransaction_Timeout(t *testing.T) {
	t.Parallel()

	tx := makeTx(0)

	// Simulate a receipt never appearing: return NotFound so WaitMined keeps polling.
	mock := &mockEthClient{
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return nil, ethereum.NotFound
		},
	}
	client := newTestContractClient(t, mock)

	// Use an already-cancelled context so the 60s internal timeout is immediately expired.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.waitForTransaction(ctx, tx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to wait for transaction to be mined")
}

// ---------- sendWithRetry ----------

func TestSendWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	t.Parallel()

	expectedReceipt := &types.Receipt{Status: types.ReceiptStatusSuccessful, GasUsed: 50000}

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return expectedReceipt, nil
		},
	}
	client := newTestContractClient(t, mock)

	sendCalls := 0
	sendTx := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		sendCalls++

		return makeTx(auth.Nonce.Uint64()), nil
	}

	receipt, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.NoError(t, err)
	require.Equal(t, expectedReceipt, receipt)
	require.Equal(t, 1, sendCalls)
}

func TestSendWithRetry_EstimateGasError(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 0, errSentinel
		},
	}
	client := newTestContractClient(t, mock)

	sendTx := func(_ *bind.TransactOpts) (*types.Transaction, error) {
		t.Fatal("sendTx should not be called when gas estimation fails")

		return nil, nil
	}

	_, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to estimate gas")
}

func TestSendWithRetry_SendTxError(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
	}
	client := newTestContractClient(t, mock)

	sendTx := func(_ *bind.TransactOpts) (*types.Transaction, error) {
		return nil, errSentinel
	}

	_, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to send tx")
}

func TestSendWithRetry_RetryOnOutOfGas(t *testing.T) {
	t.Parallel()

	attempt := 0
	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			attempt++
			if attempt <= 2 {
				// gasUsed/gasLimit > 0.95 triggers retry:
				// gasLimit = 100000 * 12/10 = 120000, gasUsed = 119000 => ratio = 0.9917
				return &types.Receipt{Status: types.ReceiptStatusFailed, GasUsed: 119000}, nil
			}

			return &types.Receipt{Status: types.ReceiptStatusSuccessful, GasUsed: 50000}, nil
		},
	}
	client := newTestContractClient(t, mock)

	sendCalls := 0
	sendTx := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		sendCalls++

		return makeTx(auth.Nonce.Uint64()), nil
	}

	receipt, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	require.Equal(t, 3, sendCalls)
}

func TestSendWithRetry_FailAfterMaxRetries(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			// Always return high gas usage to keep retrying
			return &types.Receipt{Status: types.ReceiptStatusFailed, GasUsed: 119000}, nil
		},
	}
	client := newTestContractClient(t, mock)

	sendCalls := 0
	sendTx := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		sendCalls++

		return makeTx(auth.Nonce.Uint64()), nil
	}

	_, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction failed after 3 attempts")
	require.Equal(t, maxRetries, sendCalls)
}

func TestSendWithRetry_FailedWithLowGasUsage(t *testing.T) {
	t.Parallel()

	// If gas usage ratio <= 0.95, do not retry — the failure is not out-of-gas.
	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			// gasLimit = 120000, gasUsed = 10000 => ratio ~0.083, well below 0.95
			return &types.Receipt{Status: types.ReceiptStatusFailed, GasUsed: 10000}, nil
		},
	}
	client := newTestContractClient(t, mock)

	sendCalls := 0
	sendTx := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		sendCalls++

		return makeTx(auth.Nonce.Uint64()), nil
	}

	_, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction failed after 3 attempts")
	// Should NOT retry because gas ratio is low — only 1 attempt
	require.Equal(t, 1, sendCalls)
}

func TestSendWithRetry_WaitForTxError(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return nil, ethereum.NotFound
		},
	}
	client := newTestContractClient(t, mock)

	sendTx := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return makeTx(auth.Nonce.Uint64()), nil
	}

	// Use an already-cancelled context so waitForTransaction returns quickly.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.sendWithRetry(ctx, "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to wait for tx")
}

func TestSendWithRetry_PassesValueToAuth(t *testing.T) {
	t.Parallel()

	expectedValue := big.NewInt(999)

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return &types.Receipt{Status: types.ReceiptStatusSuccessful, GasUsed: 50000}, nil
		},
	}
	client := newTestContractClient(t, mock)

	var capturedAuth *bind.TransactOpts
	sendTx := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		capturedAuth = auth

		return makeTx(auth.Nonce.Uint64()), nil
	}

	_, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, expectedValue, sendTx)
	require.NoError(t, err)
	require.Equal(t, expectedValue, capturedAuth.Value)
}

func TestSendWithRetry_CreateTransactOptsError(t *testing.T) {
	t.Parallel()

	mock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		pendingNonceAtFn: func(_ context.Context, _ common.Address) (uint64, error) {
			return 0, errSentinel
		},
	}
	client := newTestContractClient(t, mock)

	sendTx := func(_ *bind.TransactOpts) (*types.Transaction, error) {
		t.Fatal("sendTx should not be called when createTransactOpts fails")

		return nil, nil
	}

	_, err := client.sendWithRetry(context.Background(), "TestMethod", common.Address{}, []byte{0x01}, nil, sendTx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create transact opts")
}

// ---------- Register ----------

func TestRegister_Success(t *testing.T) {
	t.Parallel()

	successReceipt := &types.Receipt{Status: types.ReceiptStatusSuccessful, GasUsed: 80000}

	ethMock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return successReceipt, nil
		},
	}

	dkgMock := &mockDKGContract{
		feeFn: func(_ *bind.CallOpts) (*big.Int, error) {
			return big.NewInt(1000), nil
		},
	}

	client := newFullTestContractClient(t, ethMock, dkgMock, &mockCDRContract{})

	startBlockHash := make([]byte, 32)
	receipt, err := client.Register(
		context.Background(),
		1,
		[32]byte{0x01},
		100,
		startBlockHash,
		[]byte("dkg-pub-key"),
		[]byte("comm-pub-key"),
		[]byte("enclave-report"),
	)
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}

func TestRegister_FeeQueryError(t *testing.T) {
	t.Parallel()

	ethMock := &mockEthClient{}
	dkgMock := &mockDKGContract{
		feeFn: func(_ *bind.CallOpts) (*big.Int, error) {
			return nil, errSentinel
		},
	}

	client := newFullTestContractClient(t, ethMock, dkgMock, &mockCDRContract{})

	startBlockHash := make([]byte, 32)
	_, err := client.Register(
		context.Background(), 1, [32]byte{0x01}, 100, startBlockHash,
		[]byte("dkg-pub-key"), []byte("comm-pub-key"), []byte("report"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to query DKG fee")
}

func TestRegister_InvalidStartBlockHash(t *testing.T) {
	t.Parallel()

	ethMock := &mockEthClient{}
	dkgMock := &mockDKGContract{}
	client := newFullTestContractClient(t, ethMock, dkgMock, &mockCDRContract{})

	// startBlockHash must be exactly 32 bytes; pass wrong length.
	_, err := client.Register(
		context.Background(), 1, [32]byte{0x01}, 100, []byte("too-short"),
		[]byte("dkg-pub-key"), []byte("comm-pub-key"), []byte("report"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to convert startBlockHash to bytes32")
}

// ---------- Finalize ----------

func TestFinalize_Success(t *testing.T) {
	t.Parallel()

	successReceipt := &types.Receipt{Status: types.ReceiptStatusSuccessful, GasUsed: 90000}

	ethMock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return successReceipt, nil
		},
	}

	dkgMock := &mockDKGContract{
		feeFn: func(_ *bind.CallOpts) (*big.Int, error) {
			return big.NewInt(2000), nil
		},
	}

	client := newFullTestContractClient(t, ethMock, dkgMock, &mockCDRContract{})

	participantsRoot := make([]byte, 32)
	receipt, err := client.Finalize(
		context.Background(),
		1,
		[32]byte{0x01},
		participantsRoot,
		[]byte("global-pub-key"),
		[][]byte{[]byte("coeff1")},
		[]byte("pub-key-share"),
		[]byte("signature"),
	)
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}

func TestFinalize_FeeQueryError(t *testing.T) {
	t.Parallel()

	ethMock := &mockEthClient{}
	dkgMock := &mockDKGContract{
		feeFn: func(_ *bind.CallOpts) (*big.Int, error) {
			return nil, errSentinel
		},
	}
	client := newFullTestContractClient(t, ethMock, dkgMock, &mockCDRContract{})

	participantsRoot := make([]byte, 32)
	_, err := client.Finalize(
		context.Background(), 1, [32]byte{0x01}, participantsRoot,
		[]byte("global-pub-key"), [][]byte{}, []byte("share"), []byte("sig"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to query DKG fee for finalize")
}

func TestFinalize_InvalidParticipantsRoot(t *testing.T) {
	t.Parallel()

	ethMock := &mockEthClient{}
	dkgMock := &mockDKGContract{}
	client := newFullTestContractClient(t, ethMock, dkgMock, &mockCDRContract{})

	// participantsRoot must be exactly 32 bytes.
	_, err := client.Finalize(
		context.Background(), 1, [32]byte{0x01}, []byte("short"),
		[]byte("global-pub-key"), [][]byte{}, []byte("share"), []byte("sig"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to convert participants root to bytes32")
}

// ---------- SubmitEncryptedPartialDecryption ----------

func TestSubmitEncryptedPartialDecryption_Success(t *testing.T) {
	t.Parallel()

	successReceipt := &types.Receipt{Status: types.ReceiptStatusSuccessful, GasUsed: 70000}

	ethMock := &mockEthClient{
		estimateGasFn: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 100000, nil
		},
		transactionReceiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return successReceipt, nil
		},
	}

	cdrMock := &mockCDRContract{
		baseFeeFn: func(_ *bind.CallOpts) (*big.Int, error) {
			return big.NewInt(500), nil
		},
	}

	client := newFullTestContractClient(t, ethMock, &mockDKGContract{}, cdrMock)

	receipt, err := client.SubmitEncryptedPartialDecryption(
		context.Background(),
		1,
		0,
		[]byte("encrypted-partial"),
		[]byte("ephemeral-pub-key"),
		[]byte("pub-share"),
		[]byte("requester-pub-key"),
		[]byte("ciphertext"),
		42,
		[]byte("signature"),
	)
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}

func TestSubmitEncryptedPartialDecryption_BaseFeeError(t *testing.T) {
	t.Parallel()

	ethMock := &mockEthClient{}
	cdrMock := &mockCDRContract{
		baseFeeFn: func(_ *bind.CallOpts) (*big.Int, error) {
			return nil, errSentinel
		},
	}
	client := newFullTestContractClient(t, ethMock, &mockDKGContract{}, cdrMock)

	_, err := client.SubmitEncryptedPartialDecryption(
		context.Background(), 1, 0,
		[]byte("partial"), []byte("eph"), []byte("share"),
		[]byte("req"), []byte("cipher"), 42, []byte("sig"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to query CDR base fee")
}

// ---------- maxRetries constant ----------

func TestMaxRetriesConstant(t *testing.T) {
	t.Parallel()

	require.Equal(t, 3, maxRetries)
}

// errSentinel is a test sentinel error.
var errSentinel = errTestSentinel("sentinel error")

type errTestSentinel string

func (e errTestSentinel) Error() string { return string(e) }
