package types

import (
	"context"
	"errors"

	storyerr "github.com/piplabs/story/lib/errors"
)

type EpochHooks interface {
	// the first block whose timestamp is after the duration is counted as the end of the epoch
	AfterEpochEnd(ctx context.Context, epochIdentifier string, epochNumber int64) error
	// new epoch is next block of epoch end block
	BeforeEpochStart(ctx context.Context, epochIdentifier string, epochNumber int64) error
	// Returns the name of the module implementing epoch hook.
	GetModuleName() string
}

var _ EpochHooks = MultiEpochHooks{}

// combine multiple hooks, all hook functions are run in array sequence.
type MultiEpochHooks []EpochHooks

// GetModuleName implements EpochHooks.
func (MultiEpochHooks) GetModuleName() string {
	return ModuleName
}

func NewMultiEpochHooks(hooks ...EpochHooks) MultiEpochHooks {
	return hooks
}

// AfterEpochEnd is called when epoch is going to be ended, epochNumber is the number of epoch that is ending.
func (h MultiEpochHooks) AfterEpochEnd(ctx context.Context, epochIdentifier string, epochNumber int64) error {
	var errs error
	for i := range h {
		errs = errors.Join(errs, h[i].AfterEpochEnd(ctx, epochIdentifier, epochNumber))
	}

	if errs != nil {
		return storyerr.Wrap(errs, "after epoch end")
	}

	return nil
}

// BeforeEpochStart is called when epoch is going to be started, epochNumber is the number of epoch that is starting.
func (h MultiEpochHooks) BeforeEpochStart(ctx context.Context, epochIdentifier string, epochNumber int64) error {
	var errs error
	for i := range h {
		errs = errors.Join(errs, h[i].BeforeEpochStart(ctx, epochIdentifier, epochNumber))
	}

	if errs != nil {
		return storyerr.Wrap(errs, "before epoch start")
	}

	return nil
}

// EpochHooksWrapper is a wrapper for modules to inject EpochHooks using depinject.
type EpochHooksWrapper struct{ EpochHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (EpochHooksWrapper) IsOnePerModuleType() {}
