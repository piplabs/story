//nolint:wrapcheck // The api server is our server, so we don't need to wrap it.
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initSlashingRoute() {
	s.httpMux.HandleFunc("/slashing/params", utils.SimpleWrap(s.aminoCodec, s.GetSlashingParams))
	s.httpMux.HandleFunc("/slashing/signing_infos", utils.AutoWrap(s.aminoCodec, s.GetSigningInfos))
	s.httpMux.HandleFunc("/slashing/signing_infos/{cons_address}", utils.SimpleWrap(s.aminoCodec, s.GetSigningInfo))
}

func (s *Server) GetSlashingParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetSlashingKeeper()).Params(queryContext, &slashingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

func (s *Server) GetSigningInfos(req *getSigningInfosRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetSlashingKeeper()).SigningInfos(queryContext, &slashingtypes.QuerySigningInfosRequest{
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

func (s *Server) GetSigningInfo(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetSlashingKeeper()).SigningInfo(queryContext, &slashingtypes.QuerySigningInfoRequest{
		ConsAddress: mux.Vars(r)["cons_address"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
