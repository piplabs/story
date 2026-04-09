package keeper

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestRetry_SuccessOnFirstAttempt verifies that retry returns nil when fn
// succeeds on the first call.
func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	t.Parallel()

	calls := 0
	err := retry(context.Background(), func(_ context.Context) error {
		calls++
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 1, calls, "fn should be called exactly once on success")
}

// TestRetry_SuccessOnSecondAttempt verifies that retry keeps calling fn until
// it succeeds and returns nil.
// NOTE: this test sleeps for retryDelay (2s); skipped in -short mode.
func TestRetry_SuccessOnSecondAttempt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping sleep-based test in short mode")
	}

	t.Parallel()

	calls := 0
	err := retry(context.Background(), func(_ context.Context) error {
		calls++
		if calls < 2 {
			return errors.New("transient failure")
		}
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

// TestRetry_AllAttemptsExhausted verifies that retry returns an error after
// all retryAttempts attempts fail.
// NOTE: this test sleeps for retryDelay * (retryAttempts-1) seconds; skipped in -short mode.
func TestRetry_AllAttemptsExhausted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping sleep-based test in short mode")
	}

	t.Parallel()

	calls := 0
	err := retry(context.Background(), func(_ context.Context) error {
		calls++
		return errors.New("permanent failure")
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "all retries failed")
	require.Equal(t, retryAttempts, calls, "fn should be called exactly retryAttempts times")
}

// TestRetry_ContextPassedThrough verifies that the ctx passed to retry is
// forwarded to the fn on each call.
func TestRetry_ContextPassedThrough(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "test-value")

	var receivedValues []interface{}
	err := retry(ctx, func(c context.Context) error {
		receivedValues = append(receivedValues, c.Value(ctxKey{}))
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, []interface{}{"test-value"}, receivedValues)
}

// TestRetry_RetryDelayConstant verifies that retryDelay is set to 2 seconds
// as documented.
func TestRetry_RetryDelayConstant(t *testing.T) {
	t.Parallel()
	require.Equal(t, 2*time.Second, retryDelay)
}

// TestRetry_RetryAttemptsConstant verifies that retryAttempts is 3 as documented.
func TestRetry_RetryAttemptsConstant(t *testing.T) {
	t.Parallel()
	require.Equal(t, 3, retryAttempts)
}
