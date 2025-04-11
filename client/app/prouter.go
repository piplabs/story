package app

import (
	"bytes"
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// processTimeout is the maximum time to process a proposal.
// Timeout results in rejecting the proposal, which could negatively affect liveness.
// But it avoids blocking forever, which also negatively affects liveness.
const processTimeout = time.Second * 10

// makeProcessProposalRouter creates a new process proposal router that only routes
// expected messages to expected modules.
func makeProcessProposalRouter(app *App) *baseapp.MsgServiceRouter {
	router := baseapp.NewMsgServiceRouter()
	router.SetInterfaceRegistry(app.interfaceRegistry)
	app.Keepers.EVMEngKeeper.RegisterProposalService(router) // EVMEngine calls NewPayload on proposals to verify it.

	return router
}

// makeProcessProposalHandler creates a new process proposal handler.
// It ensures all messages included in a cpayload proposal are valid.
// It also updates some external state.
func makeProcessProposalHandler(router *baseapp.MsgServiceRouter, txConfig client.TxConfig) sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		// Only allow 10s to process a proposal. Reject proposal otherwise.
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx.Context(), processTimeout)
		defer timeoutCancel()
		ctx = ctx.WithContext(timeoutCtx)

		if req.Height == 1 {
			if len(req.Txs) > 0 { // First proposal must be empty.
				return rejectProposal(ctx, errors.New("first proposal not empty"))
			}

			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
		} else if len(req.Txs) != 1 {
			return rejectProposal(ctx, errors.New("unexpected transactions in proposal"))
		}

		// Ensure the proposal includes quorum votes.
		var totalPower, votedPower int64
		for _, vote := range req.ProposedLastCommit.Votes {
			totalPower += vote.Validator.Power
			if vote.BlockIdFlag != cmttypes.BlockIDFlagCommit {
				continue
			}
			votedPower += vote.Validator.Power
		}
		if totalPower*2/3 >= votedPower {
			return rejectProposal(ctx, errors.New("proposed doesn't include quorum votes extensions"))
		}

		// Ensure only expected messages types are included the expected number of times.
		expectedMsgCounts := map[string]int{
			sdk.MsgTypeURL(&evmenginetypes.MsgExecutionPayload{}): 1, // Only a single EVM execution payload is allowed.
		}

		rawTX := req.Txs[0]
		tx, err := txConfig.TxDecoder()(rawTX)
		if err != nil {
			return rejectProposal(ctx, errors.Wrap(err, "decode transaction"))
		}

		if err = validateTx(tx); err != nil {
			return rejectProposal(ctx, errors.Wrap(err, "validate tx"))
		}

		for _, msg := range tx.GetMsgs() {
			typeURL := sdk.MsgTypeURL(msg)

			// Ensure the message type is expected and not included too many times.
			if i, ok := expectedMsgCounts[typeURL]; !ok {
				return rejectProposal(ctx, errors.New("unexpected message type", "msg_type", typeURL))
			} else if i <= 0 {
				return rejectProposal(ctx, errors.New("message type included too many times", "msg_type", typeURL))
			}
			expectedMsgCounts[typeURL]--

			handler := router.Handler(msg)
			if handler == nil {
				return rejectProposal(ctx, errors.New("msg handler not found [BUG]", "msg_type", typeURL))
			}

			if _, err := handler(ctx, msg); err != nil {
				return rejectProposal(ctx, errors.Wrap(err, "execute message"))
			}
		}
		// Ensure all expected messages types are included with enough times.
		for typeURL, count := range expectedMsgCounts {
			if count != 0 {
				return rejectProposal(ctx, errors.New("message type not included with expected times", "msg_type", typeURL))
			}
		}

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

// rejectProposal rejects proposal with abci.ResponseProcessProposal_REJECT status
//
//nolint:unparam // Explicitly return nil error
func rejectProposal(ctx context.Context, err error) (*abci.ResponseProcessProposal, error) {
	log.Error(ctx, "Rejecting process proposal", err)
	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
}

type protoTxProvider interface {
	GetProtoTx() *txtypes.Tx
}

// validateTx checks whether the transaction contains any disallowed data.
func validateTx(tx sdk.Tx) error {
	standardTx, ok := tx.(signing.Tx)
	if !ok {
		return errors.New("invalid standard tx message")
	}

	signaturesV2, err := standardTx.GetSignaturesV2()
	if err != nil {
		return errors.Wrap(err, "get signatures from tx")
	}
	if len(signaturesV2) != 0 {
		return errors.New("disallowed signatures in tx")
	}

	if memo := standardTx.GetMemo(); len(memo) != 0 {
		return errors.New("disallowed memo in tx")
	}

	if fee := standardTx.GetFee(); fee != nil {
		return errors.New("disallowed fee in tx")
	}

	if !bytes.Equal(standardTx.FeePayer(), authtypes.NewModuleAddress(evmenginetypes.ModuleName).Bytes()) {
		return errors.New("invalid payer in tx")
	}

	if feeGranter := standardTx.FeeGranter(); feeGranter != nil {
		return errors.New("disallowed fee granter in tx")
	}

	tipTx, ok := tx.(txtypes.TipTx)
	if ok {
		if tip := tipTx.GetTip(); tip != nil {
			return errors.New("disallowed tip in tx")
		}
	}

	protoTx, ok := tx.(protoTxProvider)
	if ok {
		signatures := protoTx.GetProtoTx().Signatures
		if len(signatures) != 0 {
			return errors.New("disallowed signatures in tx")
		}
	}

	return nil
}
