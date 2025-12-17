package server

import (
	"encoding/hex"
	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/lib/errors"
	"net/http"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

func (s *Server) initDKGRoute() {
	s.httpMux.HandleFunc("/dkg/global_public_key", utils.SimpleWrap(s.aminoCodec, s.GetDKGGlobalPubKey))
}

func (s *Server) GetDKGGlobalPubKey(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetDKGKeeper().GetLatestDKGNetwork(queryContext, &dkgtypes.QueryGetLatestDKGNetworkRequest{})
	if err != nil {
		return nil, err
	}

	if len(queryResp.Network.GlobalPublicKey) == 0 {
		return nil, errors.New("global public key is not set yet")
	}

	return QueryDKGGlobalPublicKeyResponse{
		PublicKey: hex.EncodeToString(queryResp.Network.GlobalPublicKey),
	}, nil
}
