package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
)

func (s *Server) initEvmStakingRoute() {
	s.httpMux.HandleFunc("/evmstaking/params", utils.SimpleWrap(s.aminoCodec, s.GetEvmStakingParams))
	s.httpMux.HandleFunc("/evmstaking/delegators/{delegator_address}/operator_address", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorOperatorAddress))
	s.httpMux.HandleFunc("/evmstaking/delegators/{delegator_address}/withdraw_address", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorWithdrawAddress))
	s.httpMux.HandleFunc("/evmstaking/delegators/{delegator_address}/reward_address", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorRewardAddress))
}

// GetEvmStakingParams queries the parameters of evmstaking module.
func (s *Server) GetEvmStakingParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEvmStakingKeeper().Params(queryContext, &evmstakingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDelegatorOperatorAddress queries the operator address of a delegator.
func (s *Server) GetDelegatorOperatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEvmStakingKeeper().GetOperatorAddress(queryContext, &evmstakingtypes.QueryGetOperatorAddressRequest{
		Address: bech32AccAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDelegatorWithdrawAddress queries the withdraw address of a delegator.
func (s *Server) GetDelegatorWithdrawAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEvmStakingKeeper().GetWithdrawAddress(queryContext, &evmstakingtypes.QueryGetWithdrawAddressRequest{
		Address: bech32AccAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDelegatorRewardAddress queries the reward address of a delegator.
func (s *Server) GetDelegatorRewardAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEvmStakingKeeper().GetRewardAddress(queryContext, &evmstakingtypes.QueryGetRewardAddressRequest{
		Address: bech32AccAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
