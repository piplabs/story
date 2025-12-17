package keeper

import (
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
)

func CreateTEEClient(teeEndpoint string) (types.TEEClient, error) {
	if teeEndpoint == "" {
		return nil, errors.New("TEE endpoint is required")
	}

	var creds credentials.TransportCredentials
	if strings.HasPrefix(teeEndpoint, "https://") || strings.HasPrefix(teeEndpoint, "tls://") {
		creds = credentials.NewTLS(nil) // TODO: use provided CA certs
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(teeEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to TEE service")
	}

	return types.NewTEEClient(conn), nil
}
