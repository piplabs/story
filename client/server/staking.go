//nolint:wrapcheck,dupl // The api server is our server, so we don't need to wrap it.
package server

import (
	"net/http"
	"strconv"

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

	s.httpMux.HandleFunc("/staking/validators", utils.AutoWrap(s.aminoCodec, s.GetValidators))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegations", utils.AutoWrap(s.aminoCodec, s.GetValidatorDelegationsByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegations/{delegator_address}", utils.SimpleWrap(s.aminoCodec, s.GetDelegationByValidatorAddressDelegatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/unbonding_delegations", utils.AutoWrap(s.aminoCodec, s.GetValidatorUnbondingDelegations))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegations/{delegator_address}/unbonding_delegation", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorUnbondingDelegation))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegators/{delegator_address}/period_delegations", utils.SimpleWrap(s.aminoCodec, s.GetPeriodDelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_address}/delegators/{delegator_address}/period_delegations/{period_delegation_id}", utils.SimpleWrap(s.aminoCodec, s.GetPeriodDelegationByDelegatorAddressAndID))

	s.httpMux.HandleFunc("/staking/delegations/{delegator_address}", utils.AutoWrap(s.aminoCodec, s.GetDelegationsByDelegatorAddress))

	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/redelegations", utils.AutoWrap(s.aminoCodec, s.GetRedelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/unbonding_delegations", utils.AutoWrap(s.aminoCodec, s.GetUnbondingDelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/validators", utils.AutoWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_address}/validators/{validator_address}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddressValidatorAddress))
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

// GetPeriodDelegationsByDelegatorAddress queries period delegations info for given validator delegator pair.
func (s *Server) GetPeriodDelegationsByDelegatorAddress(r *http.Request) (resp any, err error) {
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

	queryResp, err := s.store.GetStakingKeeper().GetAllPeriodDelegationsByDelAndValAddr(queryContext, bech32AccAddress, bech32ValAddress)
	if err != nil {
		return nil, err
	}

	for i := range queryResp {
		evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp[i].DelegatorAddress)
		if err != nil {
			return nil, err
		}
		queryResp[i].DelegatorAddress = evmDelegatorAddress

		evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp[i].ValidatorAddress)
		if err != nil {
			return nil, err
		}
		queryResp[i].ValidatorAddress = evmValidatorAddress
	}

	return queryResp, nil
}

// GetPeriodDelegationByDelegatorAddressAndID queries period delegation info for given validator delegator pair and period delegation id.
func (s *Server) GetPeriodDelegationByDelegatorAddressAndID(r *http.Request) (resp any, err error) {
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

	queryResp, err := s.store.GetStakingKeeper().GetPeriodDelegation(queryContext, bech32AccAddress, bech32ValAddress, mux.Vars(r)["period_delegation_id"])
	if err != nil {
		return nil, err
	}

	evmDelegatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	queryResp.DelegatorAddress = evmDelegatorAddress

	evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	queryResp.ValidatorAddress = evmValidatorAddress

	return queryResp, nil
}
