//nolint:testpackage // Ignore this linter rule because we want to test private method, so package name should be same with the target package name.
package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/contracts/bindings"
)

func Test_mustGetABI(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       *bind.MetaData
		expectPanic bool
	}{
		{
			name:        "No Panic with valid UpgradeEntrypoint metadata",
			input:       bindings.UpgradeEntrypointMetaData,
			expectPanic: false,
		},
		{
			name:        "No Panic with valid UBIPool metadata",
			input:       bindings.UBIPoolMetaData,
			expectPanic: false,
		},
		{
			name:        "Panics with invalid metadata",
			input:       &bind.MetaData{ABI: "invalid ABI"},
			expectPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.expectPanic {
				require.Panics(t, func() {
					mustGetABI(tc.input)
				})
			} else {
				require.NotPanics(t, func() {
					mustGetABI(tc.input)
				})
			}
		})
	}
}

func Test_mustGetUpgradeEvent(t *testing.T) {
	t.Parallel()

	abi := mustGetABI(bindings.UpgradeEntrypointMetaData)
	testCases := []struct {
		name        string
		eventName   string
		expectPanic bool
	}{
		{
			name:        "No Panic for valid event SoftwareUpgrade",
			eventName:   "SoftwareUpgrade",
			expectPanic: false,
		},
		{
			name:        "Panics for non-existent event",
			eventName:   "NonExistentEvent",
			expectPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.expectPanic {
				require.Panics(t, func() {
					mustGetEvent(abi, tc.eventName)
				})
			} else {
				require.NotPanics(t, func() {
					mustGetEvent(abi, tc.eventName)
				})
			}
		})
	}
}

func Test_mustGetUBIPoolEvent(t *testing.T) {
	t.Parallel()

	abi := mustGetABI(bindings.UBIPoolMetaData)
	testCases := []struct {
		name        string
		eventName   string
		expectPanic bool
	}{
		{
			name:        "No Panic for valid event UBIPercentageSet",
			eventName:   "UBIPercentageSet",
			expectPanic: false,
		},
		{
			name:        "Panics for non-existent event",
			eventName:   "NonExistentEvent",
			expectPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.expectPanic {
				require.Panics(t, func() {
					mustGetEvent(abi, tc.eventName)
				})
			} else {
				require.NotPanics(t, func() {
					mustGetEvent(abi, tc.eventName)
				})
			}
		})
	}
}
