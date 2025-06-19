package server

import (
	"net/http"

	"github.com/piplabs/story/client/server/utils"
	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
)

func (s *Server) initEVMEngineRoute() {
	s.httpMux.HandleFunc("/evmengine/pending_upgrade", utils.SimpleWrap(s.aminoCodec, s.GetPendingUpgrade))
}

func (s *Server) GetPendingUpgrade(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEVMEngineKeeper().GetPendingUpgrade(queryContext, &evmenginetypes.QueryGetPendingUpgradeRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
