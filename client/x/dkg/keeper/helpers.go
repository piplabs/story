package keeper

import (
	"context"
	"time"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// extractReportCodeCommitment extracts the MRENCLAVE field from a raw SGX quote.
// Mirrors the Solidity _extractReportCodeCommitment in SGXValidationHook.sol:
//   - SGX quote header: 48 bytes
//   - MRENCLAVE starts at offset 64 within the report body
//   - Total offset from raw quote start: 48 + 64 = 112
//   - returns the 32-byte MRENCLAVE value
func extractReportCodeCommitment(enclaveReport []byte) []byte {
	const offset = 112
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
