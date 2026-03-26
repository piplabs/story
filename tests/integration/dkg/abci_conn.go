//go:build integration

// Package dkg provides an ABCI query adapter that implements the gogoproto
// grpc.ClientConn interface over CometBFT RPC, allowing Cosmos SDK query
// clients (e.g. dkgtypes.QueryClient) to work without a gRPC endpoint.
package dkg

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
)

// ABCIConn routes gRPC-style Invoke calls through CometBFT RPC ABCI queries.
// It implements the gogoproto grpc.ClientConn interface so that generated
// Cosmos SDK query clients (NewQueryClient) work transparently.
type ABCIConn struct {
	client rpcclient.ABCIClient
}

// NewABCIConn creates a new ABCIConn that sends ABCI queries to the given
// CometBFT RPC endpoint (e.g. "http://localhost:26657").
func NewABCIConn(rpcURL string) (*ABCIConn, error) {
	client, err := rpchttp.New(rpcURL, "")
	if err != nil {
		return nil, fmt.Errorf("create CometBFT RPC client: %w", err)
	}
	return &ABCIConn{client: client}, nil
}

// Invoke sends an ABCI query using the gRPC method path and protobuf-encoded
// request, then unmarshals the response into reply.
func (c *ABCIConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	reqBytes, err := proto.Marshal(args.(proto.Message))
	if err != nil {
		return fmt.Errorf("marshal ABCI query request: %w", err)
	}

	resp, err := c.client.ABCIQueryWithOptions(ctx, method, reqBytes, rpcclient.ABCIQueryOptions{
		Height: 0, // latest
		Prove:  false,
	})
	if err != nil {
		return fmt.Errorf("ABCI query %s: %w", method, err)
	}

	if resp.Response.Code != abci.CodeTypeOK {
		return fmt.Errorf("ABCI query %s failed: code=%d log=%s", method, resp.Response.Code, resp.Response.Log)
	}

	if err := proto.Unmarshal(resp.Response.Value, reply.(proto.Message)); err != nil {
		return fmt.Errorf("unmarshal ABCI query response: %w", err)
	}

	return nil
}

// NewStream is not supported for ABCI queries; Cosmos SDK query clients only use Invoke.
func (c *ABCIConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, fmt.Errorf("streaming not supported over ABCI query")
}
