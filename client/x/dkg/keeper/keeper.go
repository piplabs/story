package keeper

import (
	"fmt"
	"strings"
	"sync"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

var (
	// deals, responses, and justifications store TEE-generated DKG data that will be broadcast to other
	// validators through the Vote Extension. These queues act as temporary buffers between the TEE client
	// and the consensus layer, ensuring that generated data can be safely enqueued and later dequeued in
	// a thread-safe manner for propagation.
	dealsMu          sync.Mutex
	deals            []types.Deal
	responsesMu      sync.Mutex
	responses        []types.Response
	justificationsMu sync.Mutex
	justifications   []types.Justification
)

// Keeper of the dkg store.
type Keeper struct {
	cdc            codec.BinaryCodec
	storeService   storetypes.KVStoreService
	stakingKeeper  types.StakingKeeper
	valStore       baseapp.ValidatorStore
	kernelRouter   *KernelRouter
	contractClient *ContractClient
	stateManager   *StateManager

	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper

	isDKGSvcEnabled  bool
	validatorEVMAddr string   // EVM address of the validator
	enclaveType      [32]byte // TEE enclave type identifier

	Schema             collections.Schema
	ParamsStore        collections.Item[types.Params]
	DKGNetworks        collections.Map[string, types.DKGNetwork]        // key: round
	LatestDKGNetwork   collections.Item[string]                         // stores key of latest DKG network
	LatestActiveRound  collections.Item[string]                         // stores latest active round of DKG network
	DKGRegistrations   collections.Map[string, types.DKGRegistration]   // key: round_address
	GlobalPubKeyVotes  collections.Map[string, uint32]                  // key: round_globalPubKey_hash(publicCoeffs)
	KernelUpgradeInfos collections.Map[string, types.KernelUpgradeInfo] // key: upgradeVersion
	SettlementBalance  collections.Item[string]                         // remaining UBI after committee distribution during FinalizeDKGRound
}

// NewKeeper creates a new dkg Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	dk types.DistributionKeeper,
	sk types.StakingKeeper,
	valStore baseapp.ValidatorStore,
	kernelRouter *KernelRouter,
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
		cdc:                cdc,
		storeService:       storeService,
		stakingKeeper:      sk,
		bankKeeper:         bk,
		distributionKeeper: dk,
		valStore:           valStore,
		kernelRouter:       kernelRouter,
		contractClient:     contractClient,
		ParamsStore:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		DKGNetworks:        collections.NewMap(sb, types.DKGNetworkKey, "dkg_networks", collections.StringKey, codec.CollValue[types.DKGNetwork](cdc)),
		LatestDKGNetwork:   collections.NewItem(sb, types.LatestDKGNetworkKey, "latest_dkg_network", collections.StringValue),
		LatestActiveRound:  collections.NewItem(sb, types.LatestActiveRoundKey, "latest_active_round", collections.StringValue),
		DKGRegistrations:   collections.NewMap(sb, types.DKGRegistrationKey, "dkg_registrations", collections.StringKey, codec.CollValue[types.DKGRegistration](cdc)),
		GlobalPubKeyVotes:  collections.NewMap(sb, types.GlobalPubKeyVotesKey, "dkg_global_pub_key_votes", collections.StringKey, collections.Uint32Value),
		KernelUpgradeInfos: collections.NewMap(sb, types.KernelUpgradeInfoKey, "kernel_upgrade_infos", collections.StringKey, codec.CollValue[types.KernelUpgradeInfo](cdc)),
		SettlementBalance:  collections.NewItem(sb, types.SettlementBalanceKey, "settlement_balance", collections.StringValue),
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

func (k *Keeper) InitDKGService(stateDir string, addr common.Address, enclaveType [32]byte) error {
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(addr)
	k.enclaveType = enclaveType

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
	k.validatorEVMAddr = strings.ToLower(addr.Hex())
}
