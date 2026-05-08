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

	// dkgKernelMu serializes kernel DKG operations (ProcessDeals, ProcessResponses,
	// ProcessJustifications) that all mutate the same cached DistKeyGenerator in
	// story-kernel. Without this, concurrent goroutines corrupt the DKG state and
	// cause "different number of coefficients" errors during finalization.
	dkgKernelMu sync.Mutex

	// pendingIncoming* hold deals/responses that failed kernel processing (e.g., kernel
	// unreachable). They are retried on subsequent blocks when the kernel recovers.
	// In-memory only — lost on process restart (round will fail and retry naturally).
	// Deals MUST be replayed before responses (kyber's ErrNoDealBeforeResponse).
	pendingIncomingDealsMu          sync.Mutex
	pendingIncomingDeals            []types.Deal
	pendingIncomingResponsesMu      sync.Mutex
	pendingIncomingResponses        []types.Response
	pendingIncomingJustificationsMu sync.Mutex
	pendingIncomingJustifications   []types.Justification

	// maxPendingIncoming caps the pending queue to prevent memory exhaustion.
	maxPendingIncoming = 80
)

// Keeper of the dkg store.
type Keeper struct {
	cdc            codec.BinaryCodec
	storeService   storetypes.KVStoreService
	stakingKeeper  types.StakingKeeper
	valStore       baseapp.ValidatorStore
	kernelRouter   *KernelRouter
	contractClient types.DKGContractClient
	stateManager   *StateManager
	authority      string

	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper

	isDKGSvcEnabled  bool
	validatorEVMAddr string   // EVM address of the validator
	enclaveType      [32]byte // TEE enclave type identifier
	decryptBatchSize int      // number of partial decryptions per batch CDR call

	Schema             collections.Schema
	DKGNetworks        collections.Map[string, types.DKGNetwork]        // key: round
	LatestDKGNetwork   collections.Item[string]                         // stores round key of latest DKG network
	LatestActiveRound  collections.Item[string]                         // stores latest active round of DKG network
	DKGRegistrations   collections.Map[string, types.DKGRegistration]   // key: round_address
	GlobalPubKeyVotes  collections.Map[string, uint32]                  // key: round_globalPubKey_hash(publicCoeffs)
	SettlementBalance  collections.Item[string]                         // remaining UBI after committee distribution during FinalizeDKGRound
	KernelUpgradeInfos collections.Map[string, types.KernelUpgradeInfo] // key: upgradeVersion

	DKGPartialDecrypt      collections.Map[string, []byte]               // key: requesterPubKeyHash_label_ciphertextHash_round_validator; value: partial submission
	DecryptRequestRegistry collections.Map[string, types.DecryptRequest] // key: requesterPubKeyHash_label_round_ciphertextHash; value: decrypt request

	CDRPartialSubmitCount collections.Map[string, uint64] // key: validatorAddr; value: valid partial submission count
	CDRFeePoolBalance     collections.Item[string]        // total coins currently held in cdr-fee-pool
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
	contractClient types.DKGContractClient,
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
		cdc:                    cdc,
		storeService:           storeService,
		stakingKeeper:          sk,
		bankKeeper:             bk,
		distributionKeeper:     dk,
		valStore:               valStore,
		kernelRouter:           kernelRouter,
		contractClient:         contractClient,
		authority:              authority,
		decryptBatchSize:       defaultDecryptBatchSize,
		DKGNetworks:            collections.NewMap(sb, types.DKGNetworkKey, "dkg_networks", collections.StringKey, codec.CollValue[types.DKGNetwork](cdc)),
		LatestDKGNetwork:       collections.NewItem(sb, types.LatestDKGNetworkKey, "latest_dkg_network", collections.StringValue),
		LatestActiveRound:      collections.NewItem(sb, types.LatestActiveRoundKey, "latest_active_round", collections.StringValue),
		DKGRegistrations:       collections.NewMap(sb, types.DKGRegistrationKey, "dkg_registrations", collections.StringKey, codec.CollValue[types.DKGRegistration](cdc)),
		GlobalPubKeyVotes:      collections.NewMap(sb, types.GlobalPubKeyVotesKey, "dkg_global_pub_key_votes", collections.StringKey, collections.Uint32Value),
		SettlementBalance:      collections.NewItem(sb, types.SettlementBalanceKey, "settlement_balance", collections.StringValue),
		KernelUpgradeInfos:     collections.NewMap(sb, types.KernelUpgradeInfoKey, "kernel_upgrade_infos", collections.StringKey, codec.CollValue[types.KernelUpgradeInfo](cdc)),
		DKGPartialDecrypt:      collections.NewMap(sb, types.DKGPartialDecryptKey, "dkg_partial_decrypt_submissions", collections.StringKey, collections.BytesValue),
		DecryptRequestRegistry: collections.NewMap(sb, types.DecryptRequestRegistryKey, "decrypt_request_registry", collections.StringKey, codec.CollValue[types.DecryptRequest](cdc)),
		CDRPartialSubmitCount:  collections.NewMap(sb, types.CDRPartialSubmitCountKey, "cdr_partial_submit_count", collections.StringKey, collections.Uint64Value),
		CDRFeePoolBalance:      collections.NewItem(sb, types.CDRFeePoolBalanceKey, "cdr_fee_pool_balance", collections.StringValue),
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

// SetDecryptBatchSize overrides the default batch size used when submitting
// partial decryptions to the CDR contract. Must be called before the DKG
// service processes any decrypt requests.
func (k *Keeper) SetDecryptBatchSize(n int) {
	k.decryptBatchSize = n
}

// GetAuthority returns the module's authority address.
func (k *Keeper) GetAuthority() string {
	return k.authority
}
