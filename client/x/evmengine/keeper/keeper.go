package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	ormv1alpha1 "cosmossdk.io/api/cosmos/orm/v1alpha1"
	"cosmossdk.io/core/store"
	"cosmossdk.io/orm/model/ormdb"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/comet"
	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/k1util"
)

type Keeper struct {
	cdc             codec.BinaryCodec
	storeService    store.KVStoreService
	headTable       ExecutionHeadTable
	engineCl        ethclient.EngineClient
	txConfig        client.TxConfig
	cmtAPI          comet.API
	buildDelay      time.Duration
	buildOptimistic bool
	validatorAddr   common.Address

	accountKeeper    types.AccountKeeper
	evmstakingKeeper types.EvmStakingKeeper
	upgradeKeeper    types.UpgradeKeeper
	distrKeeper      types.DistrKeeper

	upgradeContract *bindings.UpgradeEntrypoint
	ubiContract     *bindings.UBIPool

	// mutablePayload contains the previous optimistically triggered payload.
	// It is optimistic because the validator set can change,
	// so we might not actually be the next proposer.
	mutablePayload struct {
		sync.Mutex
		ID        engine.PayloadID
		Height    uint64
		UpdatedAt time.Time
	}
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	engineCl ethclient.EngineClient,
	ethCl ethclient.Client,
	txConfig client.TxConfig,
	ak types.AccountKeeper,
	esk types.EvmStakingKeeper,
	uk types.UpgradeKeeper,
	dk types.DistrKeeper,
) (*Keeper, error) {
	schema := &ormv1alpha1.ModuleSchemaDescriptor{SchemaFile: []*ormv1alpha1.ModuleSchemaDescriptor_FileEntry{
		{Id: 1, ProtoFileName: File_client_x_evmengine_keeper_evmengine_proto.Path()},
	}}

	modDB, err := ormdb.NewModuleDB(schema, ormdb.ModuleDBOptions{KVStoreService: storeService})
	if err != nil {
		return nil, errors.Wrap(err, "create module db")
	}

	dbStore, err := NewEvmengineStore(modDB)
	if err != nil {
		return nil, errors.Wrap(err, "create evmengine store")
	}

	upgradeContract, err := bindings.NewUpgradeEntrypoint(common.HexToAddress(predeploys.UpgradeEntrypoint), ethCl)
	if err != nil {
		panic(fmt.Sprintf("failed to bind to the UpgradeEntrypoint contract: %s", err))
	}

	ubiContract, err := bindings.NewUBIPool(common.HexToAddress(predeploys.UBIPool), ethCl)
	if err != nil {
		panic(fmt.Sprintf("failed to bind to the UBIPool contract: %s", err))
	}

	return &Keeper{
		cdc:              cdc,
		storeService:     storeService,
		headTable:        dbStore.ExecutionHeadTable(),
		engineCl:         engineCl,
		txConfig:         txConfig,
		accountKeeper:    ak,
		evmstakingKeeper: esk,
		upgradeKeeper:    uk,
		upgradeContract:  upgradeContract,
		ubiContract:      ubiContract,
		distrKeeper:      dk,
	}, nil
}

// SetCometAPI sets the comet API client.
func (k *Keeper) SetCometAPI(c comet.API) {
	k.cmtAPI = c
}

// SetBuildDelay sets the build delay parameter.
func (k *Keeper) SetBuildDelay(d time.Duration) {
	k.buildDelay = d
}

// SetBuildOptimistic sets the optimistic build parameter.
func (k *Keeper) SetBuildOptimistic(b bool) {
	k.buildOptimistic = b
}

// SetValidatorAddress sets the validator address.
func (k *Keeper) SetValidatorAddress(addr common.Address) {
	k.validatorAddr = addr
}

// RegisterProposalService registers the proposal service on the provided router.
// This implements abci.ProcessProposal verification of new proposals.
func (k *Keeper) RegisterProposalService(server grpc1.Server) {
	types.RegisterMsgServiceServer(server, NewProposalServer(k))
}

// parseAndVerifyProposedPayload parses and returns the proposed payload
// if comparing it against the latest execution block succeeds.
//

func (k *Keeper) parseAndVerifyProposedPayload(ctx context.Context, msg *types.MsgExecutionPayload) (engine.ExecutableData, error) {
	// Parse the payload.
	var payload engine.ExecutableData
	if err := json.Unmarshal(msg.ExecutionPayload, &payload); err != nil {
		return engine.ExecutableData{}, errors.Wrap(err, "unmarshal payload")
	}

	// Fetch the latest execution head from the local keeper DB.
	head, err := k.getExecutionHead(ctx)
	if err != nil {
		return engine.ExecutableData{}, errors.Wrap(err, "latest execution block")
	}

	// Ensure the parent hash and block height matches
	if payload.Number != head.GetBlockHeight()+1 {
		return engine.ExecutableData{}, errors.New("invalid proposed payload number", "proposed", payload.Number, "head", head.GetBlockHeight())
	} else if payload.ParentHash != head.Hash() {
		return engine.ExecutableData{}, errors.New("invalid proposed payload parent hash", "proposed", payload.ParentHash, "head", head.Hash())
	}

	// Ensure the payload timestamp is after latest execution block and before or equaled to the current consensus block.
	minTimestamp := head.GetBlockTime() + 1
	maxTimestamp := uint64(sdk.UnwrapSDKContext(ctx).BlockTime().Unix())
	if maxTimestamp < minTimestamp { // Execution block minimum takes precedence
		maxTimestamp = minTimestamp
	}
	if payload.Timestamp < minTimestamp || payload.Timestamp > maxTimestamp {
		return engine.ExecutableData{}, errors.New("invalid payload timestamp",
			"proposed", payload.Timestamp, "min", minTimestamp, "max", maxTimestamp,
		)
	}

	// Ensure the Randao Digest is equaled to parent hash as this is our workaround at this point.
	if payload.Random != head.Hash() {
		return engine.ExecutableData{}, errors.New("invalid payload random", "proposed", payload.Random, "latest", head.Hash())
	}

	return payload, nil
}

// isNextProposer returns true if the local node is the proposer for the next block
//
// Note that the validator set can change, so this is an optimistic check.
func (k *Keeper) isNextProposer(ctx context.Context, currentHeight int64) (bool, error) {
	// PostFinalize can be called during block replay (performed in newCometNode),
	// but cmtAPI is set only after newCometNode completes (see app.SetCometAPI), so a nil check is necessary.
	if k.cmtAPI == nil {
		return false, nil
	}

	valset, err := k.cmtAPI.Validators(ctx, currentHeight)
	if err != nil {
		return false, err
	}

	nextProposer := valset.CopyIncrementProposerPriority(1).Proposer
	nextAddr, err := k1util.PubKeyToAddress(nextProposer.PubKey) // Convert to EVM address
	if err != nil {
		return false, err
	}

	isNextProposer := nextAddr == k.validatorAddr

	return isNextProposer, nil
}

func (k *Keeper) setOptimisticPayload(id engine.PayloadID, height uint64) {
	k.mutablePayload.Lock()
	defer k.mutablePayload.Unlock()

	k.mutablePayload.ID = id
	k.mutablePayload.Height = height
	k.mutablePayload.UpdatedAt = time.Now()
}

func (k *Keeper) getOptimisticPayload() (engine.PayloadID, uint64, time.Time) {
	k.mutablePayload.Lock()
	defer k.mutablePayload.Unlock()

	return k.mutablePayload.ID, k.mutablePayload.Height, k.mutablePayload.UpdatedAt
}
