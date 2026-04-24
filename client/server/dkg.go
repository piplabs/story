package server

import (
	"encoding/hex"
	"net/http"

	"github.com/piplabs/story/client/server/utils"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

func (s *Server) initDKGRoute() {
	s.httpMux.HandleFunc("/dkg/dkg_network", utils.AutoWrap(s.aminoCodec, s.GetDKGNetwork))
	s.httpMux.HandleFunc("/dkg/registrations", utils.AutoWrap(s.aminoCodec, s.GetAllDKGRegistrations))
	s.httpMux.HandleFunc("/dkg/registrations/verified", utils.AutoWrap(s.aminoCodec, s.GetVerifiedDKGRegistrations))
	s.httpMux.HandleFunc("/dkg/latest_active", utils.SimpleWrap(s.aminoCodec, s.GetLatestActiveDKGNetwork))
	s.httpMux.HandleFunc("/dkg/global_public_key", utils.SimpleWrap(s.aminoCodec, s.GetDKGGlobalPubKey))
}

func (s *Server) GetDKGNetwork(req *getDKGNetworkRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetDKGKeeper().GetDKGNetwork(queryContext, &dkgtypes.QueryGetDKGNetworkRequest{
		Round: req.Round,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetAllDKGRegistrations returns every DKG registration stored for the given
// round regardless of status (Verified, Finalized, or Invalidated). Unlike
// GetVerifiedDKGRegistrations, this works for already-finalized rounds, so
// clients can verify partial decryption signatures using the commPubKey that
// was active when the partial was signed.
func (s *Server) GetAllDKGRegistrations(req *getAllDKGRegistrationsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetDKGKeeper().GetAllDKGRegistrations(queryContext, &dkgtypes.QueryGetAllDKGRegistrationsRequest{
		Round: req.Round,
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
		Round: req.Round,
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

func (s *Server) GetDKGGlobalPubKey(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetDKGKeeper().GetLatestActiveDKGNetwork(queryContext, &dkgtypes.QueryGetLatestActiveDKGNetworkRequest{})
	if err != nil {
		return nil, err
	}

	if len(queryResp.Network.GlobalPublicKey) == 0 {
		return nil, errors.New("global public key is not set yet")
	}

	return QueryDKGGlobalPublicKeyResponse{
		PublicKeyHex: hex.EncodeToString(queryResp.Network.GlobalPublicKey),
	}, nil
}
