package ethclient

import (
	"context"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	defaultRPCHTTPTimeout = time.Second * 30

	newPayloadV3 = "engine_newPayloadV3"

	forkchoiceUpdatedV3 = "engine_forkchoiceUpdatedV3"

	getPayloadV3 = "engine_getPayloadV3"
)

// EngineClient defines the Engine API authenticated JSON-RPC endpoints.
// It extends the normal Client interface with the Engine API.
type EngineClient interface {
	Client

	// NewPayloadV3 creates an Eth1 block, inserts it in the chain, and returns the status of the chain.
	NewPayloadV3(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash,
		beaconRoot *common.Hash) (engine.PayloadStatusV1, error)

	// ForkchoiceUpdatedV3 is equivalent to V2 with the addition of parent beacon block root in the payload attributes.
	ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1,
		payloadAttributes *engine.PayloadAttributes) (engine.ForkChoiceResponse, error)

	// GetPayloadV3 returns a cached payload by id.
	GetPayloadV3(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error)
}

// engineClient implements EngineClient using JSON-RPC.
type engineClient struct {
	Wrapper
}

// NewAuthClient returns a new authenticated JSON-RPc engineClient.
func NewAuthClient(ctx context.Context, urlAddr string, jwtSecret []byte) (EngineClient, error) {
	transport := http.DefaultTransport
	if len(jwtSecret) > 0 {
		transport = newJWTRoundTripper(http.DefaultTransport, jwtSecret)
	}

	client := &http.Client{Timeout: defaultRPCHTTPTimeout, Transport: transport}

	rpcClient, err := rpc.DialOptions(ctx, urlAddr, rpc.WithHTTPClient(client))
	if err != nil {
		return engineClient{}, errors.Wrap(err, "rpc dial")
	}

	return engineClient{
		Wrapper: NewClient(rpcClient, "engine", urlAddr),
	}, nil
}

func (c engineClient) NewPayloadV3(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash,
	beaconRoot *common.Hash,
) (engine.PayloadStatusV1, error) {
	const endpoint = "new_payload_v3"
	defer latency(c.chain, endpoint)()

	// isStatusOk returns true if the response status is valid.
	isStatusOk := func(status engine.PayloadStatusV1) bool {
		return map[string]bool{
			engine.VALID:    true,
			engine.INVALID:  true,
			engine.SYNCING:  true,
			engine.ACCEPTED: true,
		}[status.Status]
	}

	var resp engine.PayloadStatusV1
	err := c.cl.Client().CallContext(ctx, &resp, newPayloadV3, params, versionedHashes, beaconRoot)
	if isStatusOk(resp) {
		// Swallow errors when geth returns errors along with proper responses (but at least log it).
		if err != nil {
			log.Warn(ctx, "Ignoring new_payload_v3 error with proper response", err, "status", resp.Status)
		}

		return resp, nil
	} else if err != nil {
		incError(c.chain, endpoint)
		return engine.PayloadStatusV1{}, errors.Wrap(err, "rpc new payload v3")
	} /* else err==nil && status!=ok */

	incError(c.chain, endpoint)

	return engine.PayloadStatusV1{}, errors.New("nil error and unknown status", "status", resp.Status)
}

func (c engineClient) ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1,
	payloadAttributes *engine.PayloadAttributes,
) (engine.ForkChoiceResponse, error) {
	const endpoint = "forkchoice_updated_v3"
	defer latency(c.chain, endpoint)()

	// isStatusOk returns true if the response status is valid.
	isStatusOk := func(resp engine.ForkChoiceResponse) bool {
		return map[string]bool{
			engine.VALID:    true,
			engine.INVALID:  true,
			engine.SYNCING:  true,
			engine.ACCEPTED: false, // Unexpected in ForkchoiceUpdated
		}[resp.PayloadStatus.Status]
	}

	var resp engine.ForkChoiceResponse
	err := c.cl.Client().CallContext(ctx, &resp, forkchoiceUpdatedV3, update, payloadAttributes)
	if isStatusOk(resp) {
		// Swallow errors when geth returns errors along with proper responses (but at least log it).
		if err != nil {
			log.Warn(ctx, "Ignoring forkchoice_updated_v3 error with proper response", err, "status", resp.PayloadStatus.Status)
		}

		return resp, nil
	} else if err != nil {
		incError(c.chain, endpoint)
		return engine.ForkChoiceResponse{}, errors.Wrap(err, "rpc forkchoice updated v3")
	} /* else err==nil && status!=ok */

	incError(c.chain, endpoint)

	return engine.ForkChoiceResponse{}, errors.New("nil error and unknown status", "status", resp.PayloadStatus.Status)
}

func (c engineClient) GetPayloadV3(ctx context.Context, payloadID engine.PayloadID) (
	*engine.ExecutionPayloadEnvelope, error,
) {
	const endpoint = "get_payload_v3"
	defer latency(c.chain, endpoint)()

	var resp engine.ExecutionPayloadEnvelope
	err := c.cl.Client().CallContext(ctx, &resp, getPayloadV3, payloadID)
	if err != nil {
		incError(c.chain, endpoint)
		return nil, errors.Wrap(err, "rpc get payload v3")
	}

	return &resp, nil
}
