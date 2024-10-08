package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/epochs/keeper"
	epochstypes "github.com/piplabs/story/client/x/epochs/types"
)

func (s *Server) initEpochsRoute() {
	s.httpMux.HandleFunc("/epochs/epoch_infos", utils.SimpleWrap(s.aminoCodec, s.GetEpochInfos))
	s.httpMux.HandleFunc("/epochs/current_epoch/{identifier}", utils.SimpleWrap(s.aminoCodec, s.GetCurrentEpoch))
}

// GetEpochInfos queries running epochInfos.
func (s *Server) GetEpochInfos(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(*s.store.GetEpochsKeeper()).EpochInfos(queryContext, &epochstypes.QueryEpochsInfoRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetCurrentEpoch queries current epoch of specified identifier.
func (s *Server) GetCurrentEpoch(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(*s.store.GetEpochsKeeper()).CurrentEpoch(queryContext, &epochstypes.QueryCurrentEpochRequest{
		Identifier: mux.Vars(r)["identifier"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
