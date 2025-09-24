package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/dkg/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = (*Keeper)(nil)

// Params queries the parameters of the dkg module.
func (k *Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// GetDKGNetwork queries a DKG network by mrenclave and round.
func (k *Keeper) GetDKGNetwork(ctx context.Context, req *types.QueryGetDKGNetworkRequest) (*types.QueryGetDKGNetworkResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	network, err := k.getDKGNetwork(ctx, req.Mrenclave, req.Round)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryGetDKGNetworkResponse{Network: *network}, nil
}

// GetLatestDKGNetwork queries the latest DKG network.
func (k *Keeper) GetLatestDKGNetwork(ctx context.Context, req *types.QueryGetLatestDKGNetworkRequest) (*types.QueryGetLatestDKGNetworkResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	network, err := k.getLatestDKGNetwork(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryGetLatestDKGNetworkResponse{Network: *network}, nil
}

// GetAllDKGNetworks queries all DKG networks.
func (k *Keeper) GetAllDKGNetworks(ctx context.Context, req *types.QueryGetAllDKGNetworksRequest) (*types.QueryGetAllDKGNetworksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	networks, err := k.getAllDKGNetworks(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetAllDKGNetworksResponse{Networks: networks}, nil
}

// GetDKGRegistration queries a DKG registration by mrenclave, round, and validator address.
func (k *Keeper) GetDKGRegistration(ctx context.Context, req *types.QueryGetDKGRegistrationRequest) (*types.QueryGetDKGRegistrationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: implement individual registration query using validator address
	// This would require mapping validator address to index or updating the storage key structure
	return nil, status.Error(codes.Unimplemented, "GetDKGRegistration by validator address not implemented")
}

// GetAllDKGRegistrations queries all DKG registrations (registered & verified) for a specific mrenclave and round.
func (k *Keeper) GetAllDKGRegistrations(ctx context.Context, req *types.QueryGetAllDKGRegistrationsRequest) (*types.QueryGetAllDKGRegistrationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	registrations, err := k.getDKGRegistrationsByRound(ctx, req.Mrenclave, req.Round)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Implement pagination
	return &types.QueryGetAllDKGRegistrationsResponse{
		Registrations: registrations,
		Pagination:    nil,
	}, nil
}

// GetVerifiedDKGRegistrations queries the count of verified DKG registrations (verified == registered) for a specific mrenclave and round.
func (k *Keeper) GetVerifiedDKGRegistrations(ctx context.Context, req *types.QueryGetVerifiedDKGRegistrationsRequest) (*types.QueryGetVerifiedDKGRegistrationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	registrations, err := k.getDKGRegistrationsByStatus(ctx, req.Mrenclave, req.Round, types.DKGRegStatusVerified)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetVerifiedDKGRegistrationsResponse{Registrations: registrations}, nil
}
