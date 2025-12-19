package keeper

import (
	"fmt"
	"sync"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/piplabs/story/lib/errors"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/gogoproto/grpc"

	"github.com/piplabs/story/client/x/dkg/types"
)

var (
	// deals and responses store TEE-generated DKG deals and responses that will be broadcast to other validators
	// through the Vote Extension. This queue acts as a temporary buffer between the TEE client and the consensus layer,
	// ensuring that generated deals and responses can be safely enqueued and later dequeued in a thread-safe manner for
	// propagation.
	dealsMu     sync.Mutex
	deals       []types.Deal
	responsesMu sync.Mutex
	responses   []types.Response
)

// Keeper of the dkg store.
type Keeper struct {
	cdc            codec.BinaryCodec
	storeService   storetypes.KVStoreService
	stakingKeeper  types.StakingKeeper
	valStore       baseapp.ValidatorStore
	teeClient      types.TEEClient
	contractClient *ContractClient
	stateManager   *StateManager

	isDKGSvcEnabled  bool
	validatorAddress common.Address

	Schema            collections.Schema
	ParamsStore       collections.Item[types.Params]
	DKGNetworks       collections.Map[string, types.DKGNetwork]          // key: mrenclave_round
	LatestDKGNetwork  collections.Item[string]                           // stores mrenclave key of latest DKG network
	DKGRegistrations  collections.Map[string, types.DKGRegistrationList] // key: mrenclave_round -> list of registrations
	GlobalPubKeyVotes collections.Map[string, uint32]                    // key: mrenclave_round_globalPubKey
	TEEUpgradeInfos   collections.Map[string, types.TEEUpgradeInfo]      // key: mrenclave
}

// NewKeeper creates a new dkg Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	ak types.AccountKeeper,
	sk types.StakingKeeper,
	valStore baseapp.ValidatorStore,
	teeClient types.TEEClient,
	contractClient *ContractClient,
	authority string,
) *Keeper {
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:               cdc,
		storeService:      storeService,
		stakingKeeper:     sk,
		valStore:          valStore,
		teeClient:         teeClient,
		contractClient:    contractClient,
		ParamsStore:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		DKGNetworks:       collections.NewMap(sb, types.DKGNetworkKey, "dkg_networks", collections.StringKey, codec.CollValue[types.DKGNetwork](cdc)),
		LatestDKGNetwork:  collections.NewItem(sb, types.LatestDKGNetworkKey, "latest_dkg_network", collections.StringValue),
		DKGRegistrations:  collections.NewMap(sb, types.DKGRegistrationKey, "dkg_registrations", collections.StringKey, codec.CollValue[types.DKGRegistrationList](cdc)),
		GlobalPubKeyVotes: collections.NewMap(sb, types.GlobalPubKeyVotesKey, "dkg_global_pub_key_votes", collections.StringKey, collections.Uint32Value),
		TEEUpgradeInfos:   collections.NewMap(sb, types.TEEUpgradeInfoKey, "tee_upgrade_infos", collections.StringKey, codec.CollValue[types.TEEUpgradeInfo](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return &k
}

func (k *Keeper) RegisterProposalService(server grpc.Server) {
	types.RegisterMsgServiceServer(server, NewProposalServer(k))
}

func (k *Keeper) InitDKGService(stateDir string, addr common.Address) error {
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(addr)

	stateManager, err := NewStateManager(stateDir)
	if err != nil {
		return errors.Wrap(err, "failed to create state manager")
	}

	k.stateManager = stateManager

	return nil
}

func (k *Keeper) setIsDKGSvcEnabled() {
	k.isDKGSvcEnabled = true
}

func (k *Keeper) setValidatorAddress(addr common.Address) {
	k.validatorAddress = addr
}
