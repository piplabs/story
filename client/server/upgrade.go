//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initUpgradeRoute() {
	s.httpMux.HandleFunc("/upgrade/module_versions", utils.AutoWrap(s.aminoCodec, s.ModuleVersions))
}

// GetAccounts returns all the existing accounts.
func (s *Server) ModuleVersions(req *getModuleVersionsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetUpgradeKeeper().ModuleVersions(queryContext, &upgradetypes.QueryModuleVersionsRequest{
		ModuleName: req.ModuleName,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
