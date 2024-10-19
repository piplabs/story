//nolint:wrapcheck // The api server is our server, so we don't need to wrap it.
package server

import (
	"github.com/piplabs/story/lib/k1util"
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
	s.httpMux.HandleFunc("/staking/validators/{validator_pub_key}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_pub_key}/delegations", utils.AutoWrap(s.aminoCodec, s.GetValidatorDelegationsByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_pub_key}/delegations/{delegator_pub_key}", utils.SimpleWrap(s.aminoCodec, s.GetDelegationByValidatorAddressDelegatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_pub_key}/unbonding_delegations", utils.AutoWrap(s.aminoCodec, s.GetValidatorUnbondingDelegations))
	s.httpMux.HandleFunc("/staking/validators/{validator_pub_key}/delegations/{delegator_pub_key}/unbonding_delegation", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorUnbondingDelegation))
	s.httpMux.HandleFunc("/staking/validators/{validator_pub_key}/delegators/{delegator_pub_key}/period_delegations", utils.SimpleWrap(s.aminoCodec, s.GetPeriodDelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_pub_key}/delegators/{delegator_pub_key}/period_delegations/{period_delegation_id}", utils.SimpleWrap(s.aminoCodec, s.GetPeriodDelegationByDelegatorAddressAndID))

	s.httpMux.HandleFunc("/staking/delegations/{delegator_pub_key}", utils.AutoWrap(s.aminoCodec, s.GetDelegationsByDelegatorAddress))

	s.httpMux.HandleFunc("/staking/delegators/{delegator_pub_key}", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_pub_key}/redelegations", utils.AutoWrap(s.aminoCodec, s.GetRedelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_pub_key}/unbonding_delegations", utils.AutoWrap(s.aminoCodec, s.GetUnbondingDelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_pub_key}/validators", utils.AutoWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_pub_key}/validators/{validator_pub_key}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddressValidatorAddress))
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

	for _, validator := range queryResp.Validators {
		if err := s.prepareUnpackInterfaces(validator); err != nil {
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

	valAddr, err := k1util.CmpPubKeyToValidatorAddress([]byte(mux.Vars(r)["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Validator(queryContext, &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: valAddr,
	})

	if err != nil {
		return nil, err
	}

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

	valAddr, err := k1util.CmpPubKeyToValidatorAddress([]byte(mux.Vars(r)["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).ValidatorDelegations(queryContext, &stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: valAddr,
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

	return queryResp, nil
}

// GetDelegationByValidatorAddressDelegatorAddress queries delegate info for given validator delegator pair.
func (s *Server) GetDelegationByValidatorAddressDelegatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	valAddr, err := k1util.CmpPubKeyToValidatorAddress([]byte(mux.Vars(r)["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Delegation(queryContext, &stakingtypes.QueryDelegationRequest{
		ValidatorAddr: valAddr,
		DelegatorAddr: delAddr,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidatorUnbondingDelegations queries unbonding delegations of a validator.
func (s *Server) GetValidatorUnbondingDelegations(req *getValidatorUnbondingDelegationsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	valAddr, err := k1util.CmpPubKeyToValidatorAddress([]byte(mux.Vars(r)["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).ValidatorUnbondingDelegations(queryContext, &stakingtypes.QueryValidatorUnbondingDelegationsRequest{
		ValidatorAddr: valAddr,
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

	return queryResp, nil
}

// GetDelegatorUnbondingDelegation queries unbonding info for given validator delegator pair.
func (s *Server) GetDelegatorUnbondingDelegation(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	valAddr, err := k1util.CmpPubKeyToValidatorAddress([]byte(mux.Vars(r)["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).UnbondingDelegation(queryContext, &stakingtypes.QueryUnbondingDelegationRequest{
		ValidatorAddr: valAddr,
		DelegatorAddr: delAddr,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDelegationsByDelegatorAddress queries all delegations of a given delegator address.
func (s *Server) GetDelegationsByDelegatorAddress(req *getDelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorDelegations(queryContext, &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: delAddr,
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

	return queryResp, nil
}

// GetRedelegationsByDelegatorAddress queries redelegations of given address.
func (s *Server) GetRedelegationsByDelegatorAddress(req *getRedelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Redelegations(queryContext, &stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr:    delAddr,
		SrcValidatorAddr: req.SrcValidatorAddr,
		DstValidatorAddr: req.DstValidatorAddr,
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

	return queryResp, nil
}

// GetDelegatorByDelegatorAddress queries delegator info for given delegator address.
func (s *Server) GetDelegatorByDelegatorAddress(r *http.Request) (resp any, err error) {
	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	delWithdrawEvmAddr, err := s.store.GetEvmStakingKeeper().DelegatorWithdrawAddress.Get(r.Context(), delAddr)
	if err != nil {
		return nil, err
	}

	delRewardEvmAddr, err := s.store.GetEvmStakingKeeper().DelegatorRewardAddress.Get(r.Context(), delAddr)
	if err != nil {
		return nil, err
	}

	delOperatorEvmAddr, err := s.store.GetEvmStakingKeeper().DelegatorOperatorAddress.Get(r.Context(), delAddr)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"delegator_pubkey": mux.Vars(r)["delegator_pub_key"],
		"withdraw_address": delWithdrawEvmAddr,
		"reward_address":   delRewardEvmAddr,
		"operator_address": delOperatorEvmAddr,
	}, nil
}

// GetUnbondingDelegationsByDelegatorAddress queries all unbonding delegations of a given delegator address.
func (s *Server) GetUnbondingDelegationsByDelegatorAddress(req *getUnbondingDelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorUnbondingDelegations(queryContext, &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: delAddr,
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

	return queryResp, nil
}

// GetValidatorsByDelegatorAddress queries all validators info for given delegator address.
func (s *Server) GetValidatorsByDelegatorAddress(req *getValidatorsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorValidators(queryContext, &stakingtypes.QueryDelegatorValidatorsRequest{
		DelegatorAddr: delAddr,
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

	for _, validator := range queryResp.Validators {
		if err = s.prepareUnpackInterfaces(validator); err != nil {
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

	valAddr, err := k1util.CmpPubKeyToValidatorAddress([]byte(mux.Vars(r)["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	delAddr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(mux.Vars(r)["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorValidator(queryContext, &stakingtypes.QueryDelegatorValidatorRequest{
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
	})

	if err != nil {
		return nil, err
	}

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

	muxVars := mux.Vars(r)
	valAddrStr, err := k1util.CmpPubKeyToValidatorAddress([]byte(muxVars["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	delAddrStr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(muxVars["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	valAddr, err := s.store.GetAccountKeeper().AddressCodec().StringToBytes(valAddrStr)
	if err != nil {
		return nil, err
	}

	delAddr, err := s.store.GetAccountKeeper().AddressCodec().StringToBytes(delAddrStr)
	if err != nil {
		return nil, err
	}

	return s.store.GetStakingKeeper().GetAllPeriodDelegationsByDelAndValAddr(queryContext, delAddr, valAddr)
}

// GetPeriodDelegationByDelegatorAddressAndID queries period delegation info for given validator delegator pair and period delegation id.
func (s *Server) GetPeriodDelegationByDelegatorAddressAndID(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	muxVars := mux.Vars(r)
	valAddrStr, err := k1util.CmpPubKeyToValidatorAddress([]byte(muxVars["validator_pub_key"]))
	if err != nil {
		return nil, err
	}

	delAddrStr, err := k1util.CmpPubKeyToDelegatorAddress([]byte(muxVars["delegator_pub_key"]))
	if err != nil {
		return nil, err
	}

	valAddr, err := s.store.GetAccountKeeper().AddressCodec().StringToBytes(valAddrStr)
	if err != nil {
		return nil, err
	}

	delAddr, err := s.store.GetAccountKeeper().AddressCodec().StringToBytes(delAddrStr)
	if err != nil {
		return nil, err
	}

	return s.store.GetStakingKeeper().GetPeriodDelegation(queryContext, delAddr, valAddr, muxVars["period_delegation_id"])
}
