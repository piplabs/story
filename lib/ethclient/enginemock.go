package ethclient

import (
	"context"
	"crypto/sha256"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	fuzz "github.com/google/gofuzz"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

type payloadArgs struct {
	params     engine.ExecutableData
	beaconRoot *common.Hash
}

//nolint:gochecknoglobals // This is a static mapping.
var stakingAbi = mustGetABI(bindings.IPTokenStakingMetaData)
var depositEvent = stakingAbi.Events["Deposit"]
var withdrawEvent = stakingAbi.Events["Withdraw"]
var slashingAbi = mustGetABI(bindings.IPTokenSlashingMetaData)
var unjailEvent = slashingAbi.Events["Unjail"]

var _ EngineClient = (*engineMock)(nil)

// engineMock mocks the Engine API for testing purposes.
type engineMock struct {
	Client
	fuzzer            *fuzz.Fuzzer
	randomErrs        float64
	failOnBlockHashes map[common.Hash]bool // Field for block hashes that should trigger a failure

	mu          sync.Mutex
	head        *types.Block
	pendingLogs map[common.Address][]types.Log
	logs        map[common.Hash][]types.Log
	payloads    map[engine.PayloadID]payloadArgs
}

// WithMockSelfDelegate returns an option to add a deposit event to the mock.
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
			Data: data,
		}

		mock.pendingLogs[contractAddr] = []types.Log{eventLog}
	}
}

func WithMockUnstake(blkHash common.Hash, contractAddr common.Address, delPubKeyBytes, valPubKeyBytes []byte, ether int64) func(mock *engineMock) {
	return func(mock *engineMock) {
		mock.mu.Lock()
		defer mock.mu.Unlock()

		wei := new(big.Int).Mul(big.NewInt(ether), big.NewInt(params.Ether))
		data, err := withdrawEvent.Inputs.NonIndexed().Pack(delPubKeyBytes, valPubKeyBytes, wei)
		if err != nil {
			panic(errors.Wrap(err, "pack delegate"))
		}

		eventLog := types.Log{
			Address: contractAddr,
			Topics: []common.Hash{
				withdrawEvent.ID,
			},
			Data: data,
		}

		mock.logs[blkHash] = []types.Log{eventLog}

	}
}

func WithMockUnjail(blkHash common.Hash, contractAddr, valEvmAddr common.Address, valPubKeyBytes []byte) func(mock *engineMock) {
	return func(mock *engineMock) {
		mock.mu.Lock()
		defer mock.mu.Unlock()

		data, err := unjailEvent.Inputs.NonIndexed().Pack(valPubKeyBytes)
		if err != nil {
			panic(errors.Wrap(err, "pack delegate"))
		}

		eventLog := types.Log{
			Address: contractAddr,
			Topics: []common.Hash{
				unjailEvent.ID,
				common.BytesToHash(valEvmAddr.Bytes()),
			},
			Data: data,
		}

		mock.logs[blkHash] = []types.Log{eventLog}

	}
}

// TODO: consider making this mocking function more generic for other events also
func WithMockUnstakeAndUnjail(
	blkHash common.Hash, stakingContractAddr, slashingContractAddr, valEvmAddr common.Address,
	delPubKeyBytes, valPubKeyBytes []byte, ether int64) func(mock *engineMock,
) {
	return func(mock *engineMock) {
		mock.mu.Lock()
		defer mock.mu.Unlock()

		wei := new(big.Int).Mul(big.NewInt(ether), big.NewInt(params.Ether))
		data, err := withdrawEvent.Inputs.NonIndexed().Pack(delPubKeyBytes, valPubKeyBytes, wei)
		if err != nil {
			panic(errors.Wrap(err, "pack delegate"))
		}

		eventLog := types.Log{
			Address: stakingContractAddr,
			Topics: []common.Hash{
				withdrawEvent.ID,
			},
			Data: data,
		}

		mock.logs[blkHash] = []types.Log{eventLog}

		data, err = unjailEvent.Inputs.NonIndexed().Pack(valPubKeyBytes)
		if err != nil {
			panic(errors.Wrap(err, "pack delegate"))
		}

		eventLog = types.Log{
			Address: slashingContractAddr,
			Topics: []common.Hash{
				unjailEvent.ID,
				common.BytesToHash(valEvmAddr.Bytes()),
			},
			Data: data,
		}

		mock.logs[blkHash] = append(mock.logs[blkHash], eventLog)
	}
}

func WithMockFailedOnBlockHashes(hashes []common.Hash) func(*engineMock) {
	return func(mock *engineMock) {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		mock.failOnBlockHashes = make(map[common.Hash]bool)
		for _, hash := range hashes {
			mock.failOnBlockHashes[hash] = true
		}
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

	genesisPayload, err := makePayload(fuzzer, height, uint64(timestamp), parentHash, common.Address{}, parentHash, &parentBeaconRoot)
	if err != nil {
		return nil, errors.Wrap(err, "make next payload")
	}
	genesisBlock, err := engine.ExecutableDataToBlock(genesisPayload, nil, &parentBeaconRoot)
	if err != nil {
		return nil, errors.Wrap(err, "executable data to block")
	}

	return genesisBlock, nil
}

// NewEngineMock returns a new mock engine API client.
// Note only some methods are implemented, it will panic if you call an unimplemented method.
func NewEngineMock(opts ...func(mock *engineMock)) (EngineClient, error) {
	genesisBlock, err := MockGenesisBlock()
	if err != nil {
		return nil, err
	}

	m := &engineMock{
		fuzzer:      NewFuzzer(int64(genesisBlock.Time())),
		head:        genesisBlock,
		pendingLogs: make(map[common.Address][]types.Log),
		payloads:    make(map[engine.PayloadID]payloadArgs),
		logs:        make(map[common.Hash][]types.Log),
	}
	for _, opt := range opts {
		opt(m)
	}

	return m, nil
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

func (m *engineMock) FilterLogs(_ context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if q.BlockHash == nil || len(q.Addresses) == 0 {
		return nil, nil
	}

	// If the block hash is in the failOnBlockHashes map, return an error
	for hash := range m.failOnBlockHashes {
		if *q.BlockHash == hash {
			return nil, errors.New("failed to fetch logs")
		}
	}

	// Initialize response slice for collecting logs
	var resp []types.Log

	// Ensure we return the same logs for the same query.
	if eventLogs, ok := m.logs[*q.BlockHash]; ok {
		for _, eventLog := range eventLogs {
			for _, addr := range q.Addresses {
				if eventLog.Address == addr {
					resp = append(resp, eventLog)
					break // Address match found; move to the next log
				}
			}
		}

		return resp, nil
	}

	// If there are no logs in the main map, check pending logs for all addresses
	for _, addr := range q.Addresses {
		if eventLogs, ok := m.pendingLogs[addr]; ok {
			resp = append(resp, eventLogs...)
			delete(m.pendingLogs, addr)
		}
	}

	// Store the collected logs for future queries
	if len(resp) > 0 {
		m.logs[*q.BlockHash] = resp
	}

	return resp, nil
}

func (m *engineMock) BlockNumber(ctx context.Context) (uint64, error) {
	if err := m.maybeErr(ctx); err != nil {
		return 0, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.head.NumberU64(), nil
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

	if hash != m.head.Hash() {
		return nil, errors.New("only head hash supported") // Only support latest block
	}

	return m.head.Header(), nil
}

func (m *engineMock) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	if err := m.maybeErr(ctx); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if number == nil {
		return m.head, nil
	}

	if number.Cmp(m.head.Number()) != 0 {
		return nil, errors.New("block not found") // Only support latest block
	}

	return m.head, nil
}

func (m *engineMock) NewPayloadV3(ctx context.Context, params engine.ExecutableData, _ []common.Hash, beaconRoot *common.Hash) (engine.PayloadStatusV1, error) {
	if err := m.maybeErr(ctx); err != nil {
		return engine.PayloadStatusV1{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

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

	// Maybe update head
	//nolint: nestif // this is a mock it's fine
	if m.head.Hash() != update.HeadBlockHash {
		var found bool
		for _, args := range m.payloads {
			block, err := engine.ExecutableDataToBlock(args.params, nil, args.beaconRoot)
			if err != nil {
				return engine.ForkChoiceResponse{}, errors.Wrap(err, "executable data to block")
			}

			if block.Hash() != update.HeadBlockHash {
				continue
			}

			if err := verifyChild(m.head, block); err != nil {
				return engine.ForkChoiceResponse{}, err
			}

			m.head = block
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
				log.Hex7("forkchoice", m.head.Hash().Bytes()))
		}
	}

	// If we have payload attributes, make a new payload
	if attrs != nil {
		payload, err := makePayload(m.fuzzer, m.head.NumberU64()+1,
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
		"height", m.head.NumberU64(),
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

// TODO(corver): Add support for V3

func (*engineMock) NewPayloadV2(context.Context, engine.ExecutableData) (engine.PayloadStatusV1, error) {
	panic("implement me")
}

func (*engineMock) ForkchoiceUpdatedV2(context.Context, engine.ForkchoiceStateV1, *engine.PayloadAttributes,
) (engine.ForkChoiceResponse, error) {
	panic("implement me")
}

func (*engineMock) GetPayloadV2(context.Context, engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	panic("implement me")
}

// makePayload returns a new fuzzed payload using head as parent if provided.
func makePayload(fuzzer *fuzz.Fuzzer, height uint64, timestamp uint64, parentHash common.Hash,
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
	block := types.NewBlock(&header, nil, nil, trie.NewStackTrie(nil))

	// Convert block to payload
	env := engine.BlockToExecutableData(block, big.NewInt(0), nil)
	payload := *env.ExecutionPayload

	// Ensure the block is valid
	_, err := engine.ExecutableDataToBlock(payload, nil, beaconRoot)
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

	return engine.PayloadID(hash[:8]), nil
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
