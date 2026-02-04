package keeper

import (
	"context"
	"encoding/hex"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/cast"

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

// GetDKGNetwork queries a DKG network by code commitment and round.
func (k *Keeper) GetDKGNetwork(ctx context.Context, req *types.QueryGetDKGNetworkRequest) (*types.QueryGetDKGNetworkResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	codeCommitmentBz, err := hex.DecodeString(req.CodeCommitmentHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid code commitment")
	}

	codeCommitment, err := cast.ToBytes32(codeCommitmentBz)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid length of code commitment")
	}

	network, err := k.getDKGNetwork(ctx, codeCommitment, req.Round)
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

// GetDKGRegistration queries a DKG registration by code commitment, round, and validator address.
func (*Keeper) GetDKGRegistration(_ context.Context, req *types.QueryGetDKGRegistrationRequest) (*types.QueryGetDKGRegistrationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: implement individual registration query using validator address
	// This would require mapping validator address to index or updating the storage key structure
	return nil, status.Error(codes.Unimplemented, "GetDKGRegistration by validator address not implemented")
}

// GetAllDKGRegistrations queries all DKG registrations (registered & verified) for a specific code commitment and round.
func (k *Keeper) GetAllDKGRegistrations(ctx context.Context, req *types.QueryGetAllDKGRegistrationsRequest) (*types.QueryGetAllDKGRegistrationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	codeCommitmentBz, err := hex.DecodeString(req.CodeCommitmentHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid code commitment")
	}

	codeCommitment, err := cast.ToBytes32(codeCommitmentBz)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid length of code commitment")
	}

	registrations, err := k.getDKGRegistrationsByRound(ctx, codeCommitment, req.Round)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Implement pagination
	return &types.QueryGetAllDKGRegistrationsResponse{
		Registrations: registrations,
		Pagination:    nil,
	}, nil
}

// GetAllVerifiedDKGRegistrations queries all verified DKG registrations for a specific code commitment and round.
func (k *Keeper) GetAllVerifiedDKGRegistrations(ctx context.Context, req *types.QueryGetAllVerifiedDKGRegistrationsRequest) (*types.QueryGetAllVerifiedDKGRegistrationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	codeCommitmentBz, err := hex.DecodeString(req.CodeCommitmentHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid code commitment")
	}

	codeCommitment, err := cast.ToBytes32(codeCommitmentBz)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid length of code commitment")
	}

	registrations, err := k.getDKGRegistrationsByStatus(ctx, codeCommitment, req.Round, types.DKGRegStatusVerified)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetAllVerifiedDKGRegistrationsResponse{Registrations: registrations}, nil
}

func (k *Keeper) GetLatestActiveDKGNetwork(ctx context.Context, request *types.QueryGetLatestActiveDKGNetworkRequest) (*types.QueryGetLatestActiveDKGNetworkResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	latest, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetLatestActiveDKGNetworkResponse{Network: *latest}, nil
}
