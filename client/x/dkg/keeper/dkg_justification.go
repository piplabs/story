package keeper

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/vss"

	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/share"
	vssp "go.dedis.ch/kyber/v4/share/vss/pedersen"
	"go.dedis.ch/kyber/v4/sign/schnorr"
)

// verifyJustificationVSS performs Pedersen VSS verification using the Edwards25519 suite.
// This is deterministic math safe for the consensus path.
// expectedThreshold is used to validate the commitment count matches the round's threshold.
func verifyJustificationVSS(shareBytes []byte, recipientIndex int, commitmentBytes [][]byte, expectedThreshold uint32) (bool, error) {
	suite := edwards25519.NewBlakeSHA256Ed25519()

	return vss.VerifyPedersenVSS(suite, shareBytes, recipientIndex, commitmentBytes, expectedThreshold)
}

// verifyJustificationSignature verifies that the justification's Schnorr signature
// was produced by the dealer identified by dealerIndex. It reconstructs the kyber
// vss.Justification from proto data, computes the canonical hash using kyber's
// Justification.Hash() method, and verifies the Schnorr signature against the
// dealer's registered dkgPubKey.
//
// Index convention in the reconstructed kyber struct:
//   - VSSJustification.Index: 0-based (copied directly from proto, matches kyber's vss.Response.Index)
//   - SecShare.I (PriShare.I): 0-based (copied directly from proto, matches kyber's PriPoly.Eval(i) output)
//
// These values are used as-is for hash reconstruction because kyber's Justification.Hash()
// expects the same 0-based indices that were used during deal generation and signing.
//
// This is deterministic and safe for the consensus path: it only performs
// Edwards25519 scalar/point operations and a hash comparison.
func verifyJustificationSignature(suite *edwards25519.SuiteEd25519, j types.Justification, dealerPubKeys map[uint32]kyber.Point) error {
	vssJ := j.GetVssJustification()
	if vssJ == nil {
		return errors.New("nil VSSJustification")
	}

	plainDeal := vssJ.GetPlainDeal()
	if plainDeal == nil {
		return errors.New("nil PlainDeal")
	}

	secShare := plainDeal.GetSecShare()
	if secShare == nil || secShare.GetV() == nil {
		return errors.New("nil SecShare or scalar value")
	}

	if len(vssJ.GetSignature()) == 0 {
		return errors.New("empty justification signature")
	}

	// Reconstruct the secret share scalar
	shareScalar := suite.Scalar()
	if err := shareScalar.UnmarshalBinary(secShare.GetV().GetData()); err != nil {
		return errors.Wrap(err, "unmarshal share scalar")
	}

	// Reconstruct commitments as kyber points
	commitments := make([]kyber.Point, 0, len(plainDeal.GetCommitments()))
	for _, c := range plainDeal.GetCommitments() {
		p := suite.Point()
		if err := p.UnmarshalBinary(c.GetData()); err != nil {
			return errors.Wrap(err, "unmarshal commitment point")
		}

		commitments = append(commitments, p)
	}

	// Reconstruct the kyber vss.Justification to compute the canonical hash.
	// kyber hashes: "justification" + SessionID + Index(LE uint32) + protobuf.Encode(Deal)
	//
	// All indices here are 0-based, matching kyber's internal representation:
	// - vssJ.GetIndex(): verifier/recipient index (0-based, from vss.Response.Index)
	// - secShare.GetI(): secret share index (0-based, from PriPoly.Eval(i) where i is 0-based)
	kyberJust := &vssp.Justification{
		SessionID: vssJ.GetSessionId(),
		Index:     vssJ.GetIndex(),
		Deal: &vssp.Deal{
			SessionID: plainDeal.GetSessionId(),
			SecShare: &share.PriShare{
				I: int(secShare.GetI()), // 0-based, matches kyber convention
				V: shareScalar,
			},
			T:           plainDeal.GetThreshold(),
			Commitments: commitments,
		},
		Signature: vssJ.GetSignature(),
	}

	// Look up dealer's dkgPubKey from the pre-built map
	dealerPub, ok := dealerPubKeys[j.Index]
	if !ok {
		return errors.New("no registration found for dealer index")
	}

	// Compute the justification hash using kyber's canonical method and verify
	hash := kyberJust.Hash(suite)
	if err := schnorr.Verify(suite, dealerPub, hash, kyberJust.Signature); err != nil {
		return errors.Wrap(err, "schnorr signature verification failed")
	}

	return nil
}

// verifyJustification performs Pedersen VSS verification on a single justification.
// It returns true if the revealed deal is valid (share matches commitments), false otherwise.
//
// This function only performs VSS verification — Schnorr signature verification
// must have been completed before calling this function.
func verifyJustification(latestRound *types.DKGNetwork, j types.Justification) (bool, error) {
	vssJ := j.GetVssJustification()
	if vssJ == nil {
		return false, errors.New("nil VSSJustification")
	}

	plainDeal := vssJ.GetPlainDeal()
	if plainDeal == nil {
		return false, errors.New("nil PlainDeal")
	}

	secShare := plainDeal.GetSecShare()
	if secShare == nil || secShare.GetV() == nil {
		return false, errors.New("nil SecShare or scalar value")
	}

	// VerifyPedersenVSS accepts a 1-based recipient index and internally converts
	// it back to 0-based for kyber's PubPoly.Check(). This +1/-1 round-trip exists
	// so that VerifyPedersenVSS can guard against index <= 0 at its API boundary.
	recipientIndex := int(secShare.GetI()) + 1

	// Extract commitment bytes from the plain deal
	commitmentBytes := make([][]byte, 0, len(plainDeal.GetCommitments()))
	for _, c := range plainDeal.GetCommitments() {
		commitmentBytes = append(commitmentBytes, c.GetData())
	}

	// Pedersen VSS verification: share * G == C_0 + i*C_1 + i^2*C_2 + ... + i^(t-1)*C_{t-1}
	valid, err := verifyJustificationVSS(
		secShare.GetV().GetData(),
		recipientIndex,
		commitmentBytes,
		latestRound.Threshold,
	)
	if err != nil {
		return false, errors.Wrap(err, "vss verification failed")
	}

	return valid, nil
}

// dealerPubKeyMap maps each dealer's 0-based committee position (the kyber justification Index)
// to its dkgPubKey. committee must be dealerCommitteeRegs for the dealer round.
func dealerPubKeyMap(ctx context.Context, committee []types.DKGRegistration, suite kyber.Group) map[uint32]kyber.Point {
	pubKeys := make(map[uint32]kyber.Point, len(committee))
	for i := range committee {
		if len(committee[i].DkgPubKey) == 0 {
			continue
		}

		pub := suite.Point()
		if err := pub.UnmarshalBinary(committee[i].DkgPubKey); err != nil {
			log.Warn(ctx, "Failed to unmarshal dealer dkgPubKey, skipping", err,
				"dealer_index", i,
			)

			continue
		}

		pubKeys[uint32(i)] = pub
	}

	return pubKeys
}

// invalidateDealerByAddr marks dealerAddr's registration in round as Invalidated so it cannot
// finalize. No-op if it has no registration in the round; idempotent.
func (k *Keeper) invalidateDealerByAddr(ctx context.Context, round uint32, dealerAddr common.Address) error {
	has, err := k.hasDKGRegistration(ctx, round, dealerAddr)
	if err != nil {
		return errors.Wrap(err, "failed to check registration")
	}
	if !has {
		log.Debug(ctx, "Dealer has no registration in round, skipping invalidation",
			"validator_addr", dealerAddr.Hex(),
			"round", round,
		)

		return nil
	}

	reg, err := k.getDKGRegistration(ctx, round, dealerAddr)
	if err != nil {
		return errors.Wrap(err, "failed to get registration")
	}

	// Finalized registrations should never appear here: justifications are processed during
	// the dealing stage, before finalization.
	if reg.Status == types.DKGRegStatusFinalized {
		return fmt.Errorf("dealer %s is already finalized, cannot invalidate (possible bug)", dealerAddr.Hex())
	}
	if reg.Status == types.DKGRegStatusInvalidated {
		log.Debug(ctx, "Dealer already invalidated, skipping",
			"validator_addr", dealerAddr.Hex(),
			"round", round,
		)

		return nil
	}

	reg.Status = types.DKGRegStatusInvalidated
	if err := k.setDKGRegistration(ctx, dealerAddr, reg); err != nil {
		return errors.Wrap(err, "failed to update registration status to invalidated")
	}

	log.Info(ctx, "Dealer registration invalidated",
		"validator_addr", dealerAddr.Hex(),
		"round", round,
	)

	return nil
}
