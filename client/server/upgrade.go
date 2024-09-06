//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initUpgradeRoute() {
	s.httpMux.HandleFunc("/upgrade/applied_plan/{name}", utils.SimpleWrap(s.aminoCodec, s.GetAppliedPlan))
	s.httpMux.HandleFunc("/upgrade/authority", utils.SimpleWrap(s.aminoCodec, s.GetAuthority))
	s.httpMux.HandleFunc("/upgrade/current_plan", utils.SimpleWrap(s.aminoCodec, s.GetCurrentPlan))
	s.httpMux.HandleFunc("/upgrade/module_versions", utils.AutoWrap(s.aminoCodec, s.ModuleVersions))
}

// GetAppliedPlan queries a previously applied upgrade plan by its name.
func (s *Server) GetAppliedPlan(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetUpgradeKeeper().AppliedPlan(queryContext, &upgradetypes.QueryAppliedPlanRequest{
		Name: mux.Vars(r)["name"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetAuthority returns the account with authority to conduct upgrades.
func (s *Server) GetAuthority(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetUpgradeKeeper().Authority(queryContext, &upgradetypes.QueryAuthorityRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetCurrentPlan queries the current upgrade plan.
func (s *Server) GetCurrentPlan(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetUpgradeKeeper().CurrentPlan(queryContext, &upgradetypes.QueryCurrentPlanRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// ModuleVersions queries the list of module versions from state.
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
