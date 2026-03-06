//nolint:wrapcheck,dupl // The api server is our server, so we don't need to wrap it.
package server

import (
	"net/http"
	"strconv"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initStakingRoute() {
	s.httpMux.HandleFunc("/staking/params", utils.SimpleWrap(s.aminoCodec, s.GetStakingParams))
	s.httpMux.HandleFunc("/staking/pool", utils.SimpleWrap(s.aminoCodec, s.GetStakingPool))
	s.httpMux.HandleFunc("/staking/historical_info/{height}", utils.SimpleWrap(s.aminoCodec, s.GetHistoricalInfoByHeight))
	s.httpMux.HandleFunc("/staking/total_delegators_count", utils.SimpleWrap(s.aminoCodec, s.GetTotalDelegatorsCount))
	s.httpMux.HandleFunc("/staking/total_staked_token", utils.SimpleWrap(s.aminoCodec, s.GetTotalStakedToken))

	s.httpMux.HandleFunc("/staking/validators", utils.AutoWrap(s.aminoCodec, s.GetValidators))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegations", utils.AutoWrap(s.aminoCodec, s.GetValidatorDelegationsByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegations/{delegator_address}", utils.SimpleWrap(s.aminoCodec, s.GetDelegationByValidatorAddressDelegatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/total_delegations_count", utils.SimpleWrap(s.aminoCodec, s.GetValidatorTotalDelegationsCount))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/unbonding_delegations", utils.AutoWrap(s.aminoCodec, s.GetValidatorUnbondingDelegations))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegations/{delegator_address}/unbonding_delegation", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorUnbondingDelegation))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegators/{delegator_address}/period_delegations", utils.AutoWrap(s.aminoCodec, s.GetPeriodDelegations))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegators/{delegator_address}/period_delegations/{period_delegation_id}", utils.SimpleWrap(s.aminoCodec, s.GetPeriodDelegation))

	s.httpMux.HandleFunc("/staking/delegations/{delegator_address}", utils.AutoWrap(s.aminoCodec, s.GetDelegationsByDelegatorAddress))

	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/redelegations", utils.AutoWrap(s.aminoCodec, s.GetRedelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/unbonding_delegations", utils.AutoWrap(s.aminoCodec, s.GetUnbondingDelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/validators", utils.AutoWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/validators/{validator_address}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddressValidatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/staked_token", utils.AutoWrap(s.aminoCodec, s.GetStakedTokenByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/total_staked_token", utils.SimpleWrap(s.aminoCodec, s.GetTotalStakedTokenByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/rewards_token", utils.AutoWrap(s.aminoCodec, s.GetRewardsTokenByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/total_rewards_token", utils.SimpleWrap(s.aminoCodec, s.GetTotalRewardsTokenByDelegatorAddress))
}

// GetStakingParams queries the staking parameters.
func (s *Server) GetStakingParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Params(queryContext, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetStakingPool queries the staking pool info.
func (s *Server) GetStakingPool(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Pool(queryContext, &stakingtypes.QueryPoolRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetHistoricalInfoByHeight queries the historical info for given height.
func (s *Server) GetHistoricalInfoByHeight(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	heightStr := mux.Vars(r)["height"]

	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).HistoricalInfo(queryContext, &stakingtypes.QueryHistoricalInfoRequest{
		Height: height,
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.Hist.Valset {
		evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.Hist.Valset[i].OperatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.Hist.Valset[i].OperatorAddress = evmOperatorAddress
	}

	if err := s.prepareUnpackInterfaces(queryResp.Hist); err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidators queries all validators that match the given status.
func (s *Server) GetValidators(req *getValidatorsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Validators(queryContext, &stakingtypes.QueryValidatorsRequest{
		Status: req.Status,
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.Validators {
		evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.Validators[i].OperatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.Validators[i].OperatorAddress = evmOperatorAddress

		if err := s.prepareUnpackInterfaces(queryResp.Validators[i]); err != nil {
			return nil, err
		}
	}

	return queryResp, nil
}

// GetValidatorByValidatorAddress queries validator info for given validator address.
func (s *Server) GetValidatorByValidatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Validator(queryContext, &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: bech32ValAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.Validator.OperatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.Validator.OperatorAddress = evmOperatorAddress

	if err := s.prepareUnpackInterfaces(queryResp.Validator); err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidatorDelegationsByValidatorAddress queries delegate info for given validator.
func (s *Server) GetValidatorDelegationsByValidatorAddress(req *getValidatorDelegationsByValidatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).ValidatorDelegations(queryContext, &stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: bech32ValAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.DelegationResponses {
		evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.DelegationResponses[i].Delegation.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.DelegationResponses[i].Delegation.DelegatorAddress = evmDelegatorAddress

		evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.DelegationResponses[i].Delegation.ValidatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.DelegationResponses[i].Delegation.ValidatorAddress = evmValidatorAddress
	}

	return queryResp, nil
}

func (s *Server) GetValidatorTotalDelegationsCount(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).ValidatorDelegations(queryContext, &stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: bech32ValAddress.String(),
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &QueryTotalDelegationsCountResponse{
		Total: int(queryResp.Pagination.Total),
	}, nil
}

// GetDelegationByValidatorAddressDelegatorAddress queries delegate info for given validator delegator pair.
func (s *Server) GetDelegationByValidatorAddressDelegatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Delegation(queryContext, &stakingtypes.QueryDelegationRequest{
		ValidatorAddr: bech32ValAddress.String(),
		DelegatorAddr: bech32AccAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.DelegationResponse.Delegation.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.DelegationResponse.Delegation.DelegatorAddress = evmDelegatorAddress

	evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.DelegationResponse.Delegation.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.DelegationResponse.Delegation.ValidatorAddress = evmValidatorAddress

	return queryResp, nil
}

// GetValidatorUnbondingDelegations queries unbonding delegations of a validator.
func (s *Server) GetValidatorUnbondingDelegations(req *getValidatorUnbondingDelegationsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).ValidatorUnbondingDelegations(queryContext, &stakingtypes.QueryValidatorUnbondingDelegationsRequest{
		ValidatorAddr: bech32ValAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.UnbondingResponses {
		evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.UnbondingResponses[i].DelegatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.UnbondingResponses[i].DelegatorAddress = evmDelegatorAddress

		evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.UnbondingResponses[i].ValidatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.UnbondingResponses[i].ValidatorAddress = evmValidatorAddress
	}

	return queryResp, nil
}

// GetDelegatorUnbondingDelegation queries unbonding info for given validator delegator pair.
func (s *Server) GetDelegatorUnbondingDelegation(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).UnbondingDelegation(queryContext, &stakingtypes.QueryUnbondingDelegationRequest{
		ValidatorAddr: bech32ValAddress.String(),
		DelegatorAddr: bech32AccAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.Unbond.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.Unbond.DelegatorAddress = evmDelegatorAddress

	evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.Unbond.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.Unbond.ValidatorAddress = evmValidatorAddress

	return queryResp, nil
}

// GetDelegationsByDelegatorAddress queries all delegations of a given delegator address.
func (s *Server) GetDelegationsByDelegatorAddress(req *getDelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorDelegations(queryContext, &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: bech32AccAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.DelegationResponses {
		evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.DelegationResponses[i].Delegation.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.DelegationResponses[i].Delegation.DelegatorAddress = evmDelegatorAddress

		evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.DelegationResponses[i].Delegation.ValidatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.DelegationResponses[i].Delegation.ValidatorAddress = evmValidatorAddress
	}

	return queryResp, nil
}

func (s *Server) GetTotalDelegatorsCount(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	delegators := make(map[string]struct{})

	err = s.store.GetStakingKeeper().IterateAllDelegations(queryContext, func(delegation stakingtypes.Delegation) bool {
		delegators[delegation.DelegatorAddress] = struct{}{}
		return false // continue iteration
	})
	if err != nil {
		return nil, err
	}

	return &QueryTotalDelegationsCountResponse{
		Total: len(delegators),
	}, nil
}

func (s *Server) GetTotalStakedToken(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	validators, err := s.store.GetStakingKeeper().GetAllValidators(queryContext)
	if err != nil {
		return nil, err
	}

	totalStakedToken := math.ZeroInt()

	for _, val := range validators {
		totalStakedToken = totalStakedToken.Add(val.GetTokens())
	}

	return &QueryTotalStakedTokenResponse{
		TotalStakedToken: totalStakedToken,
	}, nil
}

// GetRedelegationsByDelegatorAddress queries redelegations of given address.
func (s *Server) GetRedelegationsByDelegatorAddress(req *getRedelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	bech32AccAddr := bech32AccAddress.String()

	var bech32SrcValAddr string

	if req.SrcValidatorAddr != "" {
		bech32SrcValAddress, err := utils.EvmAddressToBech32ValAddress(req.SrcValidatorAddr)
		if err != nil {
			return nil, err
		}

		bech32SrcValAddr = bech32SrcValAddress.String()
	}

	var bech32DstValAddr string

	if req.DstValidatorAddr != "" {
		bech32DstValAddress, err := utils.EvmAddressToBech32ValAddress(req.DstValidatorAddr)
		if err != nil {
			return nil, err
		}

		bech32DstValAddr = bech32DstValAddress.String()
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Redelegations(queryContext, &stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr:    bech32AccAddr,
		SrcValidatorAddr: bech32SrcValAddr,
		DstValidatorAddr: bech32DstValAddr,
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.RedelegationResponses {
		evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.RedelegationResponses[i].Redelegation.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.RedelegationResponses[i].Redelegation.DelegatorAddress = evmDelegatorAddress

		evmValidatorSrcAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.RedelegationResponses[i].Redelegation.ValidatorSrcAddress)
		if err != nil {
			return nil, err
		}

		queryResp.RedelegationResponses[i].Redelegation.ValidatorSrcAddress = evmValidatorSrcAddress

		evmValidatorDstAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.RedelegationResponses[i].Redelegation.ValidatorDstAddress)
		if err != nil {
			return nil, err
		}

		queryResp.RedelegationResponses[i].Redelegation.ValidatorDstAddress = evmValidatorDstAddress
	}

	return queryResp, nil
}

// GetUnbondingDelegationsByDelegatorAddress queries all unbonding delegations of a given delegator address.
func (s *Server) GetUnbondingDelegationsByDelegatorAddress(req *getUnbondingDelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorUnbondingDelegations(queryContext, &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: bech32AccAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.UnbondingResponses {
		evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.UnbondingResponses[i].DelegatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.UnbondingResponses[i].DelegatorAddress = evmDelegatorAddress

		evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.UnbondingResponses[i].ValidatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.UnbondingResponses[i].ValidatorAddress = evmValidatorAddress
	}

	return queryResp, nil
}

// GetValidatorsByDelegatorAddress queries all validators info for given delegator address.
func (s *Server) GetValidatorsByDelegatorAddress(req *getValidatorsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorValidators(queryContext, &stakingtypes.QueryDelegatorValidatorsRequest{
		DelegatorAddr: bech32AccAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.Validators {
		evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.Validators[i].OperatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.Validators[i].OperatorAddress = evmOperatorAddress

		if err := s.prepareUnpackInterfaces(queryResp.Validators[i]); err != nil {
			return nil, err
		}
	}

	return queryResp, nil
}

// GetValidatorsByDelegatorAddressValidatorAddress queries validator info for given delegator validator pair.
func (s *Server) GetValidatorsByDelegatorAddressValidatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorValidator(queryContext, &stakingtypes.QueryDelegatorValidatorRequest{
		DelegatorAddr: bech32AccAddress.String(),
		ValidatorAddr: bech32ValAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.Validator.OperatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.Validator.OperatorAddress = evmOperatorAddress

	if err := s.prepareUnpackInterfaces(queryResp.Validator); err != nil {
		return nil, err
	}

	return queryResp, nil
}

func (s *Server) GetStakedTokenByDelegatorAddress(req *getStakedTokenByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorDelegations(queryContext, &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: bech32AccAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	var delStakedToken []DelegationStakedToken

	for _, delResp := range queryResp.DelegationResponses {
		valAddr, err := sdk.ValAddressFromBech32(delResp.Delegation.ValidatorAddress)
		if err != nil {
			return nil, err
		}

		val, err := keeper.NewQuerier(s.store.GetStakingKeeper()).GetValidator(queryContext, valAddr)
		if err != nil {
			return nil, err
		}

		evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(val.OperatorAddress)
		if err != nil {
			return nil, err
		}

		stakedToken := val.TokensFromShares(delResp.Delegation.Shares)
		delStakedToken = append(delStakedToken, DelegationStakedToken{
			ValidatorOperatorAddress: evmOperatorAddress,
			StakedToken:              stakedToken,
		})
	}

	return QueryStakedTokenByDelegatorAddressResponse{
		DelegationStakedToken: delStakedToken,
		Pagination:            queryResp.Pagination,
	}, nil
}

func (s *Server) GetTotalStakedTokenByDelegatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	var nextKey []byte

	totalStakedToken := math.LegacyZeroDec()

	for {
		queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorDelegations(queryContext, &stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: bech32AccAddress.String(),
			Pagination: &query.PageRequest{
				Key:        nextKey,
				Limit:      100,
				CountTotal: false,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, delResp := range queryResp.DelegationResponses {
			valAddr, err := sdk.ValAddressFromBech32(delResp.Delegation.ValidatorAddress)
			if err != nil {
				return nil, err
			}

			val, err := keeper.NewQuerier(s.store.GetStakingKeeper()).GetValidator(queryContext, valAddr)
			if err != nil {
				return nil, err
			}

			stakedToken := val.TokensFromShares(delResp.Delegation.Shares)
			totalStakedToken = totalStakedToken.Add(stakedToken)
		}

		if queryResp.Pagination == nil || len(queryResp.Pagination.NextKey) == 0 {
			break
		}

		nextKey = queryResp.Pagination.NextKey
	}

	return QueryTotalStakedTokenByDelegatorAddressResponse{
		StakedToken: totalStakedToken,
	}, nil
}

func (s *Server) GetRewardsTokenByDelegatorAddress(req *getRewardsTokenByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorDelegations(queryContext, &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: bech32AccAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	var delRewardsToken []DelegationRewardsToken

	for _, delResp := range queryResp.DelegationResponses {
		valAddr, err := sdk.ValAddressFromBech32(delResp.Delegation.ValidatorAddress)
		if err != nil {
			return nil, err
		}

		val, err := keeper.NewQuerier(s.store.GetStakingKeeper()).GetValidator(queryContext, valAddr)
		if err != nil {
			return nil, err
		}

		evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(val.OperatorAddress)
		if err != nil {
			return nil, err
		}

		rewardsToken := val.RewardsTokensFromRewardsShares(delResp.Delegation.RewardsShares)
		delRewardsToken = append(delRewardsToken, DelegationRewardsToken{
			ValidatorOperatorAddress: evmOperatorAddress,
			RewardsToken:             rewardsToken,
		})
	}

	return QueryRewardsTokenByDelegatorAddressResponse{
		DelegationRewardsToken: delRewardsToken,
		Pagination:             queryResp.Pagination,
	}, nil
}

func (s *Server) GetTotalRewardsTokenByDelegatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	var nextKey []byte

	totalRewardsToken := math.LegacyZeroDec()

	for {
		queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorDelegations(queryContext, &stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: bech32AccAddress.String(),
			Pagination: &query.PageRequest{
				Key:        nextKey,
				Limit:      100,
				CountTotal: false,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, delResp := range queryResp.DelegationResponses {
			valAddr, err := sdk.ValAddressFromBech32(delResp.Delegation.ValidatorAddress)
			if err != nil {
				return nil, err
			}

			val, err := keeper.NewQuerier(s.store.GetStakingKeeper()).GetValidator(queryContext, valAddr)
			if err != nil {
				return nil, err
			}

			rewardsToken := val.RewardsTokensFromRewardsShares(delResp.Delegation.RewardsShares)
			totalRewardsToken = totalRewardsToken.Add(rewardsToken)
		}

		if queryResp.Pagination == nil || len(queryResp.Pagination.NextKey) == 0 {
			break
		}

		nextKey = queryResp.Pagination.NextKey
	}

	return QueryTotalRewardsTokenByDelegatorAddressResponse{
		RewardsToken: totalRewardsToken,
	}, nil
}

// GetPeriodDelegations queries period delegations info for given validator delegator pair.
func (s *Server) GetPeriodDelegations(req *getPeriodDelegationsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).PeriodDelegations(queryContext, &stakingtypes.QueryPeriodDelegationsRequest{
		DelegatorAddr: bech32AccAddress.String(),
		ValidatorAddr: bech32ValAddress.String(),
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.PeriodDelegationResponses {
		evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.PeriodDelegationResponses[i].PeriodDelegation.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.PeriodDelegationResponses[i].PeriodDelegation.DelegatorAddress = evmDelegatorAddress

		evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.PeriodDelegationResponses[i].PeriodDelegation.ValidatorAddress)
		if err != nil {
			return nil, err
		}

		queryResp.PeriodDelegationResponses[i].PeriodDelegation.ValidatorAddress = evmValidatorAddress
	}

	return queryResp, nil
}

// GetPeriodDelegation queries period delegation info for given validator delegator pair and period delegation id.
func (s *Server) GetPeriodDelegation(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).PeriodDelegation(queryContext, &stakingtypes.QueryPeriodDelegationRequest{
		DelegatorAddr:      bech32AccAddress.String(),
		ValidatorAddr:      bech32ValAddress.String(),
		PeriodDelegationId: mux.Vars(r)["period_delegation_id"],
	})
	if err != nil {
		return nil, err
	}

	evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.PeriodDelegationResponse.PeriodDelegation.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.PeriodDelegationResponse.PeriodDelegation.DelegatorAddress = evmDelegatorAddress

	evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.PeriodDelegationResponse.PeriodDelegation.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	queryResp.PeriodDelegationResponse.PeriodDelegation.ValidatorAddress = evmValidatorAddress

	return queryResp, nil
}
