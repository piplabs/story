package service

import (
	"context"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// queryVerifiedDKGRegistrations queries the x/dkg module for verified registrations using gRPC.
func (s *Service) queryVerifiedDKGRegistrations(ctx context.Context, mrenclave []byte, round uint32) ([]*dkgtypes.DKGRegistration, error) {
	queryClient := dkgtypes.NewQueryClient(s.cosmosClient)

	req := &dkgtypes.QueryGetVerifiedDKGRegistrationsRequest{
		Round:     round,
		Mrenclave: mrenclave,
	}
	resp, err := queryClient.GetVerifiedDKGRegistrations(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query DKG registrations from x/dkg module")
	}

	registrations := make([]*dkgtypes.DKGRegistration, len(resp.Registrations))
	for i := range resp.Registrations {
		registrations[i] = &resp.Registrations[i]
	}

	return registrations, nil
}
