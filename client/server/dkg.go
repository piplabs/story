package server

import (
	"github.com/piplabs/story/client/server/utils"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"net/http"
)

func (s *Server) initDKGRoute() {
	s.httpMux.HandleFunc("/dkg/dkg_network", utils.AutoWrap(s.aminoCodec, s.GetDKGNetwork))
	s.httpMux.HandleFunc("/dkg/registrations/verified", utils.AutoWrap(s.aminoCodec, s.GetVerifiedDKGRegistrations))
	s.httpMux.HandleFunc("/dkg/latest_active", utils.SimpleWrap(s.aminoCodec, s.GetLatestActiveDKGNetwork))
}

func (s *Server) GetDKGNetwork(req *getDKGNetworkRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetDKGKeeper().GetDKGNetwork(queryContext, &dkgtypes.QueryGetDKGNetworkRequest{
		Round:        req.Round,
		MrenclaveHex: req.MrenclaveHex,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

func (s *Server) GetVerifiedDKGRegistrations(req *getVerifiedDKGRegistrationsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetDKGKeeper().GetAllVerifiedDKGRegistrations(queryContext, &dkgtypes.QueryGetAllVerifiedDKGRegistrationsRequest{
		Round:        req.Round,
		MrenclaveHex: req.MrenclaveHex,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

func (s *Server) GetLatestActiveDKGNetwork(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetDKGKeeper().GetLatestActiveDKGNetwork(queryContext, &dkgtypes.QueryGetLatestActiveDKGNetworkRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
