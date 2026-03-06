// Package vss provides Pedersen VSS (Verifiable Secret Sharing) verification utilities
// for the DKG complaint/justification system.
package vss

import (
	"fmt"

	"go.dedis.ch/kyber/v4"

	"github.com/piplabs/story/lib/errors"
	"go.dedis.ch/kyber/v4/share"
)

// MaxCommitments is the maximum number of polynomial commitments allowed in a
// justification. This is an arbitrary safety bound (not a VSS protocol limit)
// to prevent resource exhaustion from malicious inputs. The value matches the
// current live network's maximum validator count.
const MaxCommitments = 80

// VerifyPedersenVSS verifies that a secret share is consistent with the public commitments
// using kyber's PubPoly.Check(), which evaluates the Pedersen VSS equation:
//
//	share * G == C_0 + i*C_1 + i^2*C_2 + ... + i^(t-1)*C_{t-1}
//
// Index convention:
//   - recipientIndex must be >= 1 (1-based, Pedersen VSS evaluation point).
//   - Internally, kyber's PubPoly.Check() converts to x-coordinate via x = 1 + PriShare.I.
//   - We pass PriShare.I = recipientIndex - 1, so x = 1 + (recipientIndex - 1) = recipientIndex.
//
// When calling from justification verification, the proto's SecShare.I is 0-based
// (kyber convention), so callers must convert: recipientIndex = int(secShare.GetI()) + 1.
//
// commitmentBytes length must match expectedThreshold (if > 0) and not exceed MaxCommitments.
//
// This function is deterministic and performs only algebraic operations on the
// Edwards25519 curve, making it safe for use in the consensus path.
func VerifyPedersenVSS(suite kyber.Group, shareBytes []byte, recipientIndex int, commitmentBytes [][]byte, expectedThreshold uint32) (bool, error) {
	if len(commitmentBytes) == 0 {
		return false, errors.New("empty commitments")
	}

	if recipientIndex <= 0 {
		return false, fmt.Errorf("recipient index must be >= 1, got %d", recipientIndex)
	}

	if len(commitmentBytes) > MaxCommitments {
		return false, fmt.Errorf("commitment count %d exceeds max %d", len(commitmentBytes), MaxCommitments)
	}

	// If expectedThreshold is specified, commitment count must match
	if expectedThreshold > 0 && uint32(len(commitmentBytes)) != expectedThreshold {
		return false, fmt.Errorf("commitment count %d does not match expected threshold %d", len(commitmentBytes), expectedThreshold)
	}

	// Unmarshal the share scalar
	shareScalar := suite.Scalar()
	if err := shareScalar.UnmarshalBinary(shareBytes); err != nil {
		return false, errors.Wrap(err, "unmarshal share")
	}

	// Unmarshal commitment points
	commits := make([]kyber.Point, len(commitmentBytes))
	for i, cb := range commitmentBytes {
		p := suite.Point()
		if err := p.UnmarshalBinary(cb); err != nil {
			return false, errors.Wrap(err, "unmarshal commitment", "index", i)
		}

		commits[i] = p
	}

	// Use kyber's PubPoly.Check() for share verification.
	// PubPoly uses 0-based indexing internally (Eval(i) evaluates at x = i+1),
	// so we convert our 1-based recipientIndex to 0-based for PriShare.I.
	pubPoly := share.NewPubPoly(suite, suite.Point().Base(), commits)
	priShare := &share.PriShare{
		I: recipientIndex - 1, // convert 1-based to 0-based for kyber
		V: shareScalar,
	}

	return pubPoly.Check(priShare), nil
}
