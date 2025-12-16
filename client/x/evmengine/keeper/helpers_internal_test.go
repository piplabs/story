package keeper

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetryForever(t *testing.T) {
	t.Parallel()

	attempts := 0

	tests := []struct {
		name             string
		ctxFunc          func() context.Context
		fn               func(ctx context.Context) (bool, error)
		expectedErr      string
		expectedAttempts int
	}{
		{
			name:    "Success after retries",
			ctxFunc: context.Background,
			fn: func(ctx context.Context) (bool, error) {
				attempts++
				if attempts < 3 {
					return false, nil // Retry
				}

				return true, nil // Success
			},
			expectedErr:      "",
			expectedAttempts: 3,
		},
		{
			name: "Context canceled",
			ctxFunc: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately

				return ctx
			},
			fn: func(ctx context.Context) (bool, error) {
				return false, nil
			},
			expectedErr: "retry canceled",
		},
		{
			name:    "Func returns error",
			ctxFunc: context.Background,
			fn: func(ctx context.Context) (bool, error) {
				return false, errors.New("some error")
			},
			expectedErr: "some error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := retryForever(tc.ctxFunc(), tc.fn)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
