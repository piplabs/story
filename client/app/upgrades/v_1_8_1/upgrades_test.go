//nolint:staticcheck,revive,nolintlint // versioning
package v_1_8_1

import (
	"context"
	"testing"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/netconf"
)

// stubBank answers balance lookups from a table keyed by module name, so a test
// can state what each retired account holds.
type stubBank struct{ byModule map[string]sdk.Coins }

func (b stubBank) GetAllBalances(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	for name, coins := range b.byModule {
		if authtypes.NewModuleAddress(name).Equals(addr) {
			return coins
		}
	}

	return sdk.NewCoins()
}

// stubDistr records every FundUbi call so a test can assert the direction and
// the exact amount, and can fail on demand.
type stubDistr struct {
	funded map[string]sdk.Coins
	err    error
}

func newStubDistr() *stubDistr { return &stubDistr{funded: map[string]sdk.Coins{}} }

func (d *stubDistr) FundUbi(_ context.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	if d.err != nil {
		return d.err
	}

	for _, name := range []string{dkgModuleName, cdrFeePoolName} {
		if authtypes.NewModuleAddress(name).Equals(sender) {
			d.funded[name] = amount
		}
	}

	return nil
}

// evmstakingModule is a surviving module, used to show the version map keeps
// the modules this binary does have.
const evmstakingModule = "evmstaking"

func coins(amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, amount))
}

func newCtx(t *testing.T) sdk.Context {
	t.Helper()

	key := storetypes.NewKVStoreKey("test")

	return testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient")).Ctx
}

func TestSweepMovesBothRetiredPoolsToTheFeePool(t *testing.T) {
	t.Parallel()

	bank := stubBank{byModule: map[string]sdk.Coins{
		cdrFeePoolName: coins(7_000),
		dkgModuleName:  coins(3),
	}}
	distr := newStubDistr()

	require.NoError(t, sweepRetiredPools(newCtx(t), bank, distr))

	require.Equal(t, coins(7_000), distr.funded[cdrFeePoolName], "CDR fee pool balance was not moved")
	require.Equal(t, coins(3), distr.funded[dkgModuleName], "dkg account balance was not moved")
}

func TestSweepSkipsEmptyAccounts(t *testing.T) {
	t.Parallel()

	// The dkg account is swept every block while the module runs, so an empty
	// one is the expected case and must not produce a zero-amount transfer.
	bank := stubBank{byModule: map[string]sdk.Coins{cdrFeePoolName: coins(5)}}
	distr := newStubDistr()

	require.NoError(t, sweepRetiredPools(newCtx(t), bank, distr))

	require.Equal(t, coins(5), distr.funded[cdrFeePoolName])
	require.NotContains(t, distr.funded, dkgModuleName, "funded the fee pool from an empty account")
}

func TestSweepFailsTheUpgradeWhenTheTransferFails(t *testing.T) {
	t.Parallel()

	// A balance that cannot be moved must stop the upgrade rather than leave
	// the coins stranded in an account no code can reach afterwards.
	bank := stubBank{byModule: map[string]sdk.Coins{cdrFeePoolName: coins(1)}}
	distr := newStubDistr()
	distr.err = errors.New("fee pool rejected the transfer")

	require.Error(t, sweepRetiredPools(newCtx(t), bank, distr))
}

// newHandlerFixture builds the handler against a module manager holding no
// modules. That is enough for the version-map assertions; the sweep is covered
// directly above, where the balances can be stated.
func newHandlerFixture(t *testing.T, k keepersForTest) (sdk.Context, upgradetypes.UpgradeHandler) {
	t.Helper()

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	return newCtx(t), createHandler(module.NewManager(), module.NewConfigurator(cdc, nil, nil), k.bank, k.distr)
}

type keepersForTest struct {
	bank  BankKeeper
	distr DistributionKeeper
}

func emptyKeepers() keepersForTest {
	return keepersForTest{bank: stubBank{byModule: map[string]sdk.Coins{}}, distr: newStubDistr()}
}

func TestHandlerDropsDKGFromVersionMap(t *testing.T) {
	t.Parallel()

	ctx, handler := newHandlerFixture(t, emptyKeepers())

	newVM, err := handler(ctx, upgradetypes.Plan{}, module.VersionMap{
		"dkg":            1,
		evmstakingModule: 1,
		"evmengine":      1,
	})
	require.NoError(t, err)
	require.NotContains(t, newVM, "dkg", "dkg survived the upgrade")
}

func TestHandlerAcceptsAMapThatNeverHadDKG(t *testing.T) {
	t.Parallel()

	ctx, handler := newHandlerFixture(t, emptyKeepers())

	newVM, err := handler(ctx, upgradetypes.Plan{}, module.VersionMap{evmstakingModule: 1})
	require.NoError(t, err)
	require.NotContains(t, newVM, "dkg")
}

func TestHandlerStopsWhenTheSweepFails(t *testing.T) {
	t.Parallel()

	k := keepersForTest{
		bank:  stubBank{byModule: map[string]sdk.Coins{cdrFeePoolName: coins(1)}},
		distr: &stubDistr{funded: map[string]sdk.Coins{}, err: errors.New("boom")},
	}
	ctx, handler := newHandlerFixture(t, k)

	vm := module.VersionMap{"dkg": 1, evmstakingModule: 1}
	got, err := handler(ctx, upgradetypes.Plan{}, vm)
	require.Error(t, err)
	require.Equal(t, vm, got, "a failed upgrade must return the map it was given")
}

// The fork height is what activates the upgrade, so the resolver decides which
// networks take part. The app's generic resolver hands any name it does not
// recognise a height for every chain, so this returning false is the only thing
// keeping the fork off the networks that never ran DKG.
func TestGetUpgradeHeightOnlyResolvesForChainsThatRegisteredOne(t *testing.T) {
	t.Parallel()

	for _, chainID := range []string{netconf.AeneidChainID, netconf.StoryLocalnetID} {
		height, ok := GetUpgradeHeight(sdk.Context{}.WithChainID(chainID))
		require.Truef(t, ok, "%s does not resolve a fork height", chainID)
		require.Positivef(t, height, "%s resolved a non-positive fork height", chainID)
	}

	// Mainnet carries no entry, and neither does a chain nobody registered.
	for _, chainID := range []string{netconf.StoryChainID, "some-other-chain"} {
		height, ok := GetUpgradeHeight(sdk.Context{}.WithChainID(chainID))
		require.Falsef(t, ok, "%s resolved a fork height", chainID)
		require.Zerof(t, height, "%s resolved a fork height", chainID)
	}
}

func TestCreateUpgradeHandlerReturnsAHandler(t *testing.T) {
	t.Parallel()

	// A zero Keepers is enough: the constructor only reads the two keeper
	// fields into the closure and does not call them.
	require.NotNil(t, CreateUpgradeHandler(module.NewManager(), nil, &keepers.Keepers{}))
}
