package keeper

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"

	"cosmossdk.io/collections"

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

// GetDKGNetwork queries a DKG network by round.
func (k *Keeper) GetDKGNetwork(ctx context.Context, req *types.QueryGetDKGNetworkRequest) (*types.QueryGetDKGNetworkResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	network, err := k.getDKGNetwork(ctx, req.Round)
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

// GetAllDKGRegistrations queries all DKG registrations for a specific round.
func (k *Keeper) GetAllDKGRegistrations(ctx context.Context, req *types.QueryGetAllDKGRegistrationsRequest) (*types.QueryGetAllDKGRegistrationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	registrations, err := k.getDKGRegistrationsByRound(ctx, req.Round)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Implement pagination
	return &types.QueryGetAllDKGRegistrationsResponse{
		Registrations: registrations,
		Pagination:    nil,
	}, nil
}

// GetAllVerifiedDKGRegistrations queries all verified DKG registrations for a specific round.
func (k *Keeper) GetAllVerifiedDKGRegistrations(ctx context.Context, req *types.QueryGetAllVerifiedDKGRegistrationsRequest) (*types.QueryGetAllVerifiedDKGRegistrationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	registrations, err := k.getDKGRegistrationsByStatus(ctx, req.Round, types.DKGRegStatusVerified)
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

	if latest == nil {
		return nil, status.Error(codes.NotFound, "no active DKG network")
	}

	return &types.QueryGetLatestActiveDKGNetworkResponse{Network: *latest}, nil
}

// GetCDRPartials queries partial decryption submissions for a requester+label pair.
func (k *Keeper) GetCDRPartials(ctx context.Context, req *types.QueryGetCDRPartialsRequest) (*types.QueryGetCDRPartialsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var label [32]byte
	binary.BigEndian.PutUint32(label[28:], req.Uuid)
	requesterPubKey, err := hex.DecodeString(req.RequesterPubKeyHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid requester pubkey hex")
	}

	prefix := dkgPartialDecryptPrefix(requesterPubKey, label[:])
	rangePrefix := (&collections.Range[string]{}).Prefix(prefix)
	iter, err := k.DKGPartialDecrypt.Iterate(ctx, rangePrefix)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer iter.Close()

	grouped := make(map[string]struct {
		round      uint32
		ciphertext []byte
		items      []types.DKGPartialDecryptionSubmission
	})
	for ; iter.Valid(); iter.Next() {
		bz, err := iter.Value()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		submissionTmp, err := decodePartialDecryptionSubmission(bz)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		groupKey := fmt.Sprintf("%d:%s", submissionTmp.Round, hex.EncodeToString(submissionTmp.Ciphertext))
		entry := grouped[groupKey]
		if entry.items == nil {
			entry.round = submissionTmp.Round
			entry.ciphertext = submissionTmp.Ciphertext
		}
		entry.items = append(entry.items, *submissionTmp)
		grouped[groupKey] = entry
	}

	if len(grouped) == 0 {
		return nil, status.Error(codes.NotFound, "partial decryption submission not found")
	}

	groupedResp := make([]types.DKGPartialDecryptionSubmissionsByRound, 0, len(grouped))
	for _, entry := range grouped {
		network, err := k.getDKGNetwork(ctx, entry.round)
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		threshold := network.Threshold
		thresholdMet := uint32(len(entry.items)) >= threshold
		groupedResp = append(groupedResp, types.DKGPartialDecryptionSubmissionsByRound{
			Round:        entry.round,
			Submissions:  entry.items,
			Ciphertext:   entry.ciphertext,
			Threshold:    threshold,
			ThresholdMet: thresholdMet,
		})
	}

	sort.Slice(groupedResp, func(i, j int) bool {
		if groupedResp[i].Round != groupedResp[j].Round {
			return groupedResp[i].Round < groupedResp[j].Round
		}
		return hex.EncodeToString(groupedResp[i].Ciphertext) < hex.EncodeToString(groupedResp[j].Ciphertext)
	})

	return &types.QueryGetCDRPartialsResponse{Submissions: groupedResp}, nil
}
