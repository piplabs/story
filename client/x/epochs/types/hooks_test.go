package types_test

import (
	"context"
	"testing"

	"cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/epochs/types"
)

type KeeperTestSuite struct {
	suite.Suite
	Ctx sdk.Context
}

func TestKeeperTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.Ctx = testutil.DefaultContext(storetypes.NewKVStoreKey(types.StoreKey), storetypes.NewTransientStoreKey("transient_test"))
}

var errDummy = errors.New("9", 9, "dummyError")

// dummyEpochHook is a struct satisfying the epoch hook interface,
// that maintains a counter for how many times its been successfully called,
// and a boolean for whether it should panic during its execution.
type dummyEpochHook struct {
	successCounter int
	shouldError    bool
}

// GetModuleName implements types.EpochHooks.
func (*dummyEpochHook) GetModuleName() string {
	return "dummy"
}

func (hook *dummyEpochHook) AfterEpochEnd(ctx context.Context, epochIdentifier string, epochNumber int64) error {
	if hook.shouldError {
		return errDummy
	}
	hook.successCounter += 1

	return nil
}

func (hook *dummyEpochHook) BeforeEpochStart(ctx context.Context, epochIdentifier string, epochNumber int64) error {
	if hook.shouldError {
		return errDummy
	}
	hook.successCounter += 1

	return nil
}

func (hook *dummyEpochHook) Clone() *dummyEpochHook {
	newHook := dummyEpochHook{successCounter: hook.successCounter, shouldError: hook.shouldError}
	return &newHook
}

var _ types.EpochHooks = &dummyEpochHook{}

func (s *KeeperTestSuite) TestHooksPanicRecovery() {
	errorHook := dummyEpochHook{shouldError: true}
	noErrorHook := dummyEpochHook{shouldError: false}
	simpleHooks := []dummyEpochHook{errorHook, noErrorHook}

	tests := []struct {
		hooks                 []dummyEpochHook
		expectedCounterValues []int
		lenEvents             int
		expErr                bool
	}{
		{[]dummyEpochHook{errorHook}, []int{0}, 0, true},
		{simpleHooks, []int{0, 1, 0, 1}, 2, true},
	}

	for tcIndex, tc := range tests {
		for epochActionSelector := range 2 {
			s.SetupTest()
			hookRefs := []types.EpochHooks{}

			for _, hook := range tc.hooks {
				hookRefs = append(hookRefs, hook.Clone())
			}

			hooks := types.NewMultiEpochHooks(hookRefs...)

			//nolint:nestif // Cosmos style
			if epochActionSelector == 0 {
				err := hooks.BeforeEpochStart(s.Ctx, "id", 0)
				if tc.expErr {
					s.Require().Error(err)
				} else {
					s.Require().NoError(err)
				}
			} else if epochActionSelector == 1 {
				err := hooks.AfterEpochEnd(s.Ctx, "id", 0)
				if tc.expErr {
					s.Require().Error(err)
				} else {
					s.Require().NoError(err)
				}
			}

			for i := range hooks {
				epochHook, ok := hookRefs[i].(*dummyEpochHook)
				s.Require().True(ok)
				s.Require().Equal(tc.expectedCounterValues[i], epochHook.successCounter, "test case index %d", tcIndex)
			}
		}
	}
}
