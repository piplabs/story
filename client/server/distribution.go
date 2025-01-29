//nolint:wrapcheck // The api server is our server, so we don't need to wrap it.
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initDistributionRoute() {
	s.httpMux.HandleFunc("/distribution/params", utils.SimpleWrap(s.aminoCodec, s.GetDistributionParams))

	s.httpMux.HandleFunc("/distribution/validators/{validator_address}", utils.SimpleWrap(s.aminoCodec, s.GetDistributionValidatorByValidatorAddress))
	s.httpMux.HandleFunc("/distribution/validators/{validator_address}/commission", utils.SimpleWrap(s.aminoCodec, s.GetValidatorCommissionByValidatorAddress))
	s.httpMux.HandleFunc("/distribution/validators/{validator_address}/outstanding_rewards", utils.SimpleWrap(s.aminoCodec, s.GetValidatorOutstandingRewardsByValidatorAddress))
	s.httpMux.HandleFunc("/distribution/validators/{validator_address}/slashes", utils.AutoWrap(s.aminoCodec, s.GetValidatorSlashesByValidatorAddress))

	s.httpMux.HandleFunc("/distribution/delegators/{delegator_address}/validators", utils.SimpleWrap(s.aminoCodec, s.GetDistributionValidatorsByDelegatorAddress))
	s.httpMux.HandleFunc("/distribution/delegators/{delegator_address}/rewards", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorRewardsByDelegatorAddress))
	s.httpMux.HandleFunc("/distribution/delegators/{delegator_address}/rewards/{validator_address}", utils.SimpleWrap(s.aminoCodec, s.GetDelegatorRewardsByDelegatorAddressValidatorAddress))
}

// GetDistributionParams queries params of the distribution module.
func (s *Server) GetDistributionParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).Params(queryContext, &distributiontypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDistributionValidatorByValidatorAddress queries validator commission and self-delegation rewards for validator.
func (s *Server) GetDistributionValidatorByValidatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).ValidatorDistributionInfo(queryContext, &distributiontypes.QueryValidatorDistributionInfoRequest{
		ValidatorAddress: bech32ValAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	evmOperatorAddress, err := utils.Bech32DelegatorAddressToEvmAddress(queryResp.OperatorAddress)
	if err != nil {
		return nil, err
	}
	queryResp.OperatorAddress = evmOperatorAddress

	return queryResp, nil
}

// GetValidatorCommissionByValidatorAddress queries accumulated commission for a validator.
func (s *Server) GetValidatorCommissionByValidatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).ValidatorCommission(queryContext, &distributiontypes.QueryValidatorCommissionRequest{
		ValidatorAddress: bech32ValAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidatorOutstandingRewardsByValidatorAddress queries rewards of a validator address.
func (s *Server) GetValidatorOutstandingRewardsByValidatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).ValidatorOutstandingRewards(queryContext, &distributiontypes.QueryValidatorOutstandingRewardsRequest{
		ValidatorAddress: bech32ValAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidatorSlashesByValidatorAddress queries slash events of a validator.
func (s *Server) GetValidatorSlashesByValidatorAddress(req *getValidatorSlashesByValidatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32ValAddress, err := utils.EvmAddressToBech32ValAddress(mux.Vars(r)["validator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).ValidatorSlashes(queryContext, &distributiontypes.QueryValidatorSlashesRequest{
		ValidatorAddress: bech32ValAddress.String(),
		StartingHeight:   req.StartingHeight,
		EndingHeight:     req.EndingHeight,
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

// GetDistributionValidatorsByDelegatorAddress queries the validators of a delegator.
func (s *Server) GetDistributionValidatorsByDelegatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).DelegatorValidators(queryContext, &distributiontypes.QueryDelegatorValidatorsRequest{
		DelegatorAddress: bech32AccAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	for i := range queryResp.Validators {
		evmValidatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(queryResp.Validators[i])
		if err != nil {
			return nil, err
		}

		queryResp.Validators[i] = evmValidatorAddress
	}

	return queryResp, nil
}

// GetDelegatorRewardsByDelegatorAddress queries the total rewards accrued by each validator.
func (s *Server) GetDelegatorRewardsByDelegatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["delegator_address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).DelegationTotalRewards(queryContext, &distributiontypes.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: bech32AccAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDelegatorRewardsByDelegatorAddressValidatorAddress queries the total rewards accrued by a delegation.
func (s *Server) GetDelegatorRewardsByDelegatorAddressValidatorAddress(r *http.Request) (resp any, err error) {
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

	queryResp, err := keeper.NewQuerier(s.store.GetDistrKeeper()).DelegationRewards(queryContext, &distributiontypes.QueryDelegationRewardsRequest{
		DelegatorAddress: bech32AccAddress.String(),
		ValidatorAddress: bech32ValAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
