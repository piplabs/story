package server

import (
	"net/http"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/mint/keeper"
	minttypes "github.com/piplabs/story/client/x/mint/types"
)

func (s *Server) initMintRoute() {
	s.httpMux.HandleFunc("/mint/params", utils.SimpleWrap(s.aminoCodec, s.GetMintParams))
}

// GetMintParams queries params of the mint module.
func (s *Server) GetMintParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(*s.store.GetMintKeeper()).Params(queryContext, &minttypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
