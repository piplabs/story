package ethclient

import (
	"bytes"
	"context"
	"crypto/sha256"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	storetypes "cosmossdk.io/store/types"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	fuzz "github.com/google/gofuzz"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

const (
	// headKey is the key to store the head block.
	headKey = "head"
)

type payloadArgs struct {
	params     engine.ExecutableData
	beaconRoot *common.Hash
}

//nolint:gochecknoglobals // This is a static mapping.
var depositEvent = mustGetABI(bindings.IPTokenStakingMetaData).Events["Deposit"]

var _ EngineClient = (*engineMock)(nil)

// engineMock mocks the Engine API for testing purposes.
type engineMock struct {
	Client

	fuzzer     *fuzz.Fuzzer
	randomErrs float64

	mu sync.Mutex
	// storeKey is added to make engineMock dependent on sdk.Context for better testability.
	// By using storeKey, engineMock's methods can interact with the sdk.Context's store,
	// allowing for independent tests that do not interfere with each other’s store state.
	storeKey *storetypes.KVStoreKey
	// headKey is the key to store the head block.
	headKey      []byte
	genesisBlock *types.Block
	// consider the following maps also dependent on sdk.Context if needed.
	pendingLogs map[common.Address][]types.Log
	logs        map[common.Hash][]types.Log
	payloads    map[engine.PayloadID]payloadArgs
}

// WithMockSelfDelegation returns an option to add a deposit event to the mock.
func WithMockSelfDelegation(pubkey crypto.PubKey, ether int64) func(*engineMock) {
	return func(mock *engineMock) {
		mock.mu.Lock()
		defer mock.mu.Unlock()

		wei := new(big.Int).Mul(big.NewInt(ether), big.NewInt(params.Ether))

		valAddr, err := k1util.PubKeyToAddress(pubkey)
		if err != nil {
			panic(errors.Wrap(err, "pubkey to address"))
		}

		data, err := depositEvent.Inputs.NonIndexed().Pack(wei)
		if err != nil {
			panic(errors.Wrap(err, "pack delegate"))
		}

		// Staking predeploy addr, copied here to avoid import cycle.
		contractAddr := common.HexToAddress("0xcccccc0000000000000000000000000000000001")
		eventLog := types.Log{
			Address: contractAddr,
			Topics: []common.Hash{
				depositEvent.ID,
				common.HexToHash(valAddr.Hex()), // delegator
				common.HexToHash(valAddr.Hex()), // validator
			},
			Data:   data,
			TxHash: common.HexToHash(valAddr.Hex()),
		}

		mock.pendingLogs[contractAddr] = []types.Log{eventLog}
	}
}

type randomErrKey struct{}

// WithRandomErr returns a context that results in random engineMock errors.
// This must only be used for testing.
func WithRandomErr(ctx context.Context, _ *testing.T) context.Context {
	return context.WithValue(ctx, randomErrKey{}, true)
}

func hasRandomErr(ctx context.Context) bool {
	v, ok := ctx.Value(randomErrKey{}).(bool)
	return ok && v
}

// MockGenesisBlock returns a deterministic genesis block for testing.
func MockGenesisBlock() (*types.Block, error) {
	// Deterministic genesis block
	var (
		// Deterministic genesis block
		height           uint64 // 0
		parentHash       common.Hash
		parentBeaconRoot common.Hash
		timestamp        = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC).Unix()
		fuzzer           = NewFuzzer(timestamp)
	)

	genesisPayload, err := MakePayload(fuzzer, height, uint64(timestamp), parentHash, common.Address{}, parentHash, &parentBeaconRoot)
	if err != nil {
		return nil, errors.Wrap(err, "make next payload")
	}

	genesisBlock, err := engine.ExecutableDataToBlock(genesisPayload, nil, &parentBeaconRoot, nil)
	if err != nil {
		return nil, errors.Wrap(err, "executable data to block")
	}

	return genesisBlock, nil
}

// NewEngineMock returns a new mock engine API client.
// Note only some methods are implemented, it will panic if you call an unimplemented method.
func NewEngineMock(key *storetypes.KVStoreKey, opts ...func(mock *engineMock)) (EngineClient, error) {
	genesisBlock, err := MockGenesisBlock()
	if err != nil {
		return nil, err
	}

	m := &engineMock{
		fuzzer:       NewFuzzer(int64(genesisBlock.Time())),
		storeKey:     key,
		headKey:      []byte(headKey),
		genesisBlock: genesisBlock,
		pendingLogs:  make(map[common.Address][]types.Log),
		payloads:     make(map[engine.PayloadID]payloadArgs),
		logs:         make(map[common.Hash][]types.Log),
	}
	for _, opt := range opts {
		opt(m)
	}

	return m, nil
}

func (m *engineMock) FilterLogs(_ context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if q.BlockHash == nil || len(q.Addresses) == 0 {
		return nil, nil
	}

	addr := q.Addresses[0]

	// Ensure we return the same logs for the same query.
	if eventLogs, ok := m.logs[*q.BlockHash]; ok {
		var resp []types.Log

		for _, eventLog := range eventLogs {
			if eventLog.Address == addr {
				resp = append(resp, eventLog)
			}
		}

		return resp, nil
	}

	eventLogs, ok := m.pendingLogs[addr]
	if !ok {
		return nil, nil
	}

	m.logs[*q.BlockHash] = eventLogs
	delete(m.pendingLogs, addr)

	return eventLogs, nil
}

func (m *engineMock) BlockNumber(ctx context.Context) (uint64, error) {
	if err := m.maybeErr(ctx); err != nil {
		return 0, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	headBlock, err := m.getHeadBlock(ctx)
	if err != nil {
		return 0, err
	}

	return headBlock.NumberU64(), nil
}

func (m *engineMock) HeaderByNumber(ctx context.Context, height *big.Int) (*types.Header, error) {
	b, err := m.BlockByNumber(ctx, height)
	if err != nil {
		return nil, err
	}

	return b.Header(), nil
}

func (m *engineMock) HeaderByType(ctx context.Context, typ HeadType) (*types.Header, error) {
	if typ != HeadLatest {
		return nil, errors.New("only support latest block")
	}

	number, err := m.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	return m.HeaderByNumber(ctx, big.NewInt(int64(number)))
}

func (m *engineMock) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	if err := m.maybeErr(ctx); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	head, err := m.getHeadBlock(ctx)
	if err != nil {
		return nil, err
	}

	if hash != head.Hash() {
		return nil, errors.New("only head hash supported") // Only support latest block
	}

	return head.Header(), nil
}

func (m *engineMock) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	if err := m.maybeErr(ctx); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	head, err := m.getHeadBlock(ctx)
	if err != nil {
		return nil, err
	}

	if number == nil {
		return head, nil
	}

	if number.Cmp(head.Number()) != 0 {
		return nil, errors.New("block not found") // Only support latest block
	}

	return head, nil
}

func (m *engineMock) NewPayloadV3(ctx context.Context, params engine.ExecutableData, _ []common.Hash, beaconRoot *common.Hash) (engine.PayloadStatusV1, error) {
	if err := m.maybeErr(ctx); err != nil {
		return engine.PayloadStatusV1{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// if Withdrawals is nil, cannot rlp encode and decode properly.
	params.Withdrawals = make([]*types.Withdrawal, 0)
	args := payloadArgs{
		params:     params,
		beaconRoot: beaconRoot,
	}

	id, err := MockPayloadID(args.params, args.beaconRoot)
	if err != nil {
		return engine.PayloadStatusV1{}, err
	}

	m.payloads[id] = args

	log.Debug(ctx, "Engine mock received new payload from proposer",
		"height", params.Number,
		log.Hex7("hash", params.BlockHash.Bytes()),
	)

	return engine.PayloadStatusV1{
		Status: engine.VALID,
	}, nil
}

func (m *engineMock) ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1,
	attrs *engine.PayloadAttributes,
) (engine.ForkChoiceResponse, error) {
	if err := m.maybeErr(ctx); err != nil {
		return engine.ForkChoiceResponse{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	resp := engine.ForkChoiceResponse{
		PayloadStatus: engine.PayloadStatusV1{
			Status: engine.VALID,
		},
	}

	head, err := m.getHeadBlock(ctx)
	if err != nil {
		return engine.ForkChoiceResponse{}, err
	}

	// Maybe update head
	//nolint: nestif // this is a mock it's fine
	if head.Hash() != update.HeadBlockHash {
		var found bool

		for _, args := range m.payloads {
			block, err := engine.ExecutableDataToBlock(args.params, nil, args.beaconRoot, nil)
			if err != nil {
				return engine.ForkChoiceResponse{}, errors.Wrap(err, "executable data to block")
			}

			if block.Hash() != update.HeadBlockHash {
				continue
			}

			if err := verifyChild(head, block); err != nil {
				return engine.ForkChoiceResponse{}, err
			}

			if err := m.setHeadBlock(ctx, block); err != nil {
				return engine.ForkChoiceResponse{}, err
			}

			found = true

			id, err := MockPayloadID(args.params, args.beaconRoot)
			if err != nil {
				return engine.ForkChoiceResponse{}, err
			}

			resp.PayloadID = &id

			break
		}

		if !found {
			return engine.ForkChoiceResponse{}, errors.New("forkchoice block not found",
				log.Hex7("forkchoice", head.Hash().Bytes()))
		}
	}

	// If we have payload attributes, make a new payload
	if attrs != nil {
		payload, err := MakePayload(m.fuzzer, head.NumberU64()+1,
			attrs.Timestamp, update.HeadBlockHash, attrs.SuggestedFeeRecipient, attrs.Random, attrs.BeaconRoot)
		if err != nil {
			return engine.ForkChoiceResponse{}, err
		}

		args := payloadArgs{params: payload, beaconRoot: attrs.BeaconRoot}

		id, err := MockPayloadID(args.params, args.beaconRoot)
		if err != nil {
			return engine.ForkChoiceResponse{}, err
		}

		m.payloads[id] = args

		resp.PayloadID = &id
	}

	log.Debug(ctx, "Engine mock forkchoice updated",
		"height", head.NumberU64(),
		log.Hex7("forkchoice", update.HeadBlockHash.Bytes()),
	)

	return resp, nil
}

func (m *engineMock) GetPayloadV3(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	if err := m.maybeErr(ctx); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	args, ok := m.payloads[payloadID]
	if !ok {
		return nil, errors.New("payload not found")
	}

	return &engine.ExecutionPayloadEnvelope{
		ExecutionPayload: &args.params,
	}, nil
}

func (m *engineMock) maybeErr(ctx context.Context) error {
	if !hasRandomErr(ctx) {
		return nil
	}
	//nolint:gosec // Test code is fine.
	if rand.Float64() < m.randomErrs {
		return errors.New("test error")
	}

	return nil
}

// getHeadBlock returns the head block from the store.
func (m *engineMock) getHeadBlock(ctx context.Context) (*types.Block, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	headBz := sdkCtx.KVStore(m.storeKey).Get(m.headKey)
	if headBz == nil {
		// Set genesis block as head
		err := m.setHeadBlock(ctx, m.genesisBlock)
		if err != nil {
			return nil, err
		}

		return m.genesisBlock, nil
	}

	var headBlock types.Block

	err := rlp.DecodeBytes(headBz, &headBlock)
	if err != nil {
		return nil, errors.Wrap(err, "decode head")
	}

	return &headBlock, nil
}

// setHeadBlock sets the head block in the store.
func (m *engineMock) setHeadBlock(ctx context.Context, head *types.Block) error {
	buf := new(bytes.Buffer)

	err := head.EncodeRLP(buf)
	if err != nil {
		return errors.Wrap(err, "encode head")
	}

	headBz := buf.Bytes()
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.KVStore(m.storeKey).Set(m.headKey, headBz)

	return nil
}

// MakePayload returns a new fuzzed payload using head as parent if provided.
func MakePayload(fuzzer *fuzz.Fuzzer, height uint64, timestamp uint64, parentHash common.Hash,
	feeRecipient common.Address, randao common.Hash, beaconRoot *common.Hash) (engine.ExecutableData, error) {
	// Build a new header
	var header types.Header
	fuzzer.Fuzz(&header)
	header.Number = big.NewInt(int64(height))
	header.Time = timestamp
	header.ParentHash = parentHash
	header.MixDigest = randao      // this corresponds to Random field in PayloadAttributes
	header.Coinbase = feeRecipient // this corresponds to SuggestedFeeRecipient field in PayloadAttributes
	header.ParentBeaconRoot = beaconRoot

	// Convert header to block
	block := types.NewBlock(&header, &types.Body{Withdrawals: make([]*types.Withdrawal, 0)}, nil, trie.NewStackTrie(nil))

	// Convert block to payload
	env := engine.BlockToExecutableData(block, big.NewInt(0), nil, nil)
	payload := *env.ExecutionPayload

	// Ensure the block is valid
	_, err := engine.ExecutableDataToBlock(payload, nil, beaconRoot, nil)
	if err != nil {
		return engine.ExecutableData{}, errors.Wrap(err, "executable data to block")
	}

	return payload, nil
}

// MockPayloadID returns a deterministic payload id for the given payload.
func MockPayloadID(params engine.ExecutableData, beaconRoot *common.Hash) (engine.PayloadID, error) {
	bz, err := params.MarshalJSON()
	if err != nil {
		return engine.PayloadID{}, errors.Wrap(err, "marshal payload")
	}

	h := sha256.New()
	_, _ = h.Write(bz)
	_, _ = h.Write(beaconRoot.Bytes())
	hash := h.Sum(nil)

	return cast.Array8(hash[:8])
}

// verifyChild returns an error if child is not a valid child of parent.
func verifyChild(parent *types.Block, child *types.Block) error {
	if parent.NumberU64()+1 != child.NumberU64() {
		return errors.New("forkchoice height not following head",
			"head", parent.NumberU64(),
			"forkchoice", child.NumberU64(),
		)
	}

	if parent.Hash() != child.ParentHash() {
		return errors.New("forkchoice parent hash not head",
			log.Hex7("head", parent.Hash().Bytes()),
			log.Hex7("forkchoice", child.Hash().Bytes()),
		)
	}

	return nil
}

// mustGetABI returns the metadata's ABI as an abi.ABI type.
// It panics on error.
func mustGetABI(metadata *bind.MetaData) *abi.ABI {
	abi, err := metadata.GetAbi()
	if err != nil {
		panic(err)
	}

	return abi
}
