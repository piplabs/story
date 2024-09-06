package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/piplabs/story/client/server/utils"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
)

func (s *Server) initEvmStakingRoute() {
	s.httpMux.HandleFunc("/evmstaking/params", utils.SimpleWrap(s.aminoCodec, s.GetEvmStakingParams))
	s.httpMux.HandleFunc("/evmstaking/withdrawal_queue", utils.AutoWrap(s.aminoCodec, s.GetWithdrawalQueue))
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

// GetWithdrawalQueue queries current withdrawal queue entries.
func (s *Server) GetWithdrawalQueue(req *getWithdrawalQueueRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEvmStakingKeeper().GetWithdrawalQueue(queryContext, &evmstakingtypes.QueryGetWithdrawalQueueRequest{
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
