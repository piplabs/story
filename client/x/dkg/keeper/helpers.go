package keeper

import (
	"context"
	"time"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// extractReportInstanceDataCommitment extracts the report_data field from a raw SGX quote.
// Mirrors the Solidity _extractReportInstanceDataCommitment in SGXValidationHook.sol:
//   - SGX quote header: 48 bytes
//   - report_data starts at offset 48 + 320 = 368 within the raw quote
//   - returns the first 32 bytes of report_data
func extractReportInstanceDataCommitment(enclaveReport []byte) []byte {
	const offset = 368
	if len(enclaveReport) < offset+32 {
		return nil
	}

	return enclaveReport[offset : offset+32]
}

const (
	retryAttempts = 3
	retryDelay   = 2 * time.Second
)

func retry(ctx context.Context, fn func(ctx context.Context) error) error {
	for i := range retryAttempts {
		if err := fn(ctx); err != nil {
			log.Warn(context.Background(), "retry failed", err, "attempt", i+1)

			// Use context-aware sleep so that cancellation can interrupt the delay.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
			}

			continue
		}

		return nil
	}

	return errors.New("all retries failed")
}
