package server

import (
	"net/http"

	"github.com/piplabs/story/client/server/utils"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
)

func (s *Server) initEvmStakingRoute() {
	s.httpMux.HandleFunc("/evmstaking/params", utils.SimpleWrap(s.aminoCodec, s.GetEvmStakingParams))
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
