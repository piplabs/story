//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	"github.com/piplabs/story/client/server/utils"
	signaltypes "github.com/piplabs/story/client/x/signal/types"
)

func (s *Server) initSignalRoute() {
	s.httpMux.HandleFunc("/signal/current_plan", utils.SimpleWrap(s.aminoCodec, s.GetCurrentPlan))
}

// GetCurrentPlan queries the current upgrade plan.
func (s *Server) GetCurrentPlan(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetSignalKeeper().GetUpgrade(queryContext, &signaltypes.QueryGetUpgradeRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
