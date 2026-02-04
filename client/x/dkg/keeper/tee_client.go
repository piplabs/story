package keeper

import (
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
)

func CreateTEEClient(endpoint string) (types.TEEClient, error) {
	if endpoint == "" {
		return nil, errors.New("TEE endpoint is required")
	}

	var creds credentials.TransportCredentials
	if strings.HasPrefix(endpoint, "https://") || strings.HasPrefix(endpoint, "tls://") {
		creds = credentials.NewTLS(nil) // TODO: use provided CA certs
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to TEE service")
	}

	return types.NewTEEClient(conn), nil
}
