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
	s.httpMux.HandleFunc("/epochs/epoch_infos/{identifier}", utils.SimpleWrap(s.aminoCodec, s.GetEpochInfo))
}

// GetEpochInfos queries running epochInfos.
func (s *Server) GetEpochInfos(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(*s.store.GetEpochsKeeper()).GetEpochInfos(queryContext, &epochstypes.GetEpochInfosRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetEpochInfo queries epoch info of specified identifier.
func (s *Server) GetEpochInfo(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(*s.store.GetEpochsKeeper()).GetEpochInfo(queryContext, &epochstypes.GetEpochInfoRequest{
		Identifier: mux.Vars(r)["identifier"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
