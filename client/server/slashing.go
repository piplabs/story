//nolint:wrapcheck // The api server is our server, so we don't need to wrap it.
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initSlashingRoute() {
	s.httpMux.HandleFunc("/slashing/params", utils.SimpleWrap(s.aminoCodec, s.GetSlashingParams))
	s.httpMux.HandleFunc("/slashing/signing_infos/{pubkey}", utils.SimpleWrap(s.aminoCodec, s.GetSigningInfo))
}

// GetSlashingParams queries the parameters of slashing module.
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

// GetSigningInfo queries the signing info of given cons address.
func (s *Server) GetSigningInfo(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	consAddress, err := utils.CmpPubKeyToBech32ConsAddress(mux.Vars(r)["pubkey"])
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetSlashingKeeper()).SigningInfo(queryContext, &slashingtypes.QuerySigningInfoRequest{
		ConsAddress: consAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
