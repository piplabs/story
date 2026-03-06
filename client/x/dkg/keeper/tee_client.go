package keeper

import (
	"io"
	"strings"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// CreateKernelClient creates a gRPC client for the story-kernel.
// Returns the KernelClient and an io.Closer for the underlying gRPC connection.
func CreateKernelClient(endpoint string) (types.TEEClient, io.Closer, error) {
	if endpoint == "" {
		return nil, nil, errors.New("The endpoint is required")
	}

	var creds credentials.TransportCredentials
	if strings.HasPrefix(endpoint, "https://") || strings.HasPrefix(endpoint, "tls://") {
		creds = credentials.NewTLS(nil) // TODO: use provided CA certs
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to connect to story-kernel client")
	}

	return types.NewTEEClient(conn), conn, nil
}
