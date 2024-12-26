//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initAuthRoute() {
	s.httpMux.HandleFunc("/auth/params", utils.SimpleWrap(s.aminoCodec, s.GetAuthParams))
}

// GetAuthParams queries all parameters of auth module.
func (s *Server) GetAuthParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).Params(queryContext, &authtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
