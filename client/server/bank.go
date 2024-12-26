//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initBankRoute() {
	s.httpMux.HandleFunc("/bank/params", utils.SimpleWrap(s.aminoCodec, s.GetBankParams))

	s.httpMux.HandleFunc("/bank/supply/by_denom", utils.AutoWrap(s.aminoCodec, s.GetSupplyByDenom))

	s.httpMux.HandleFunc("/bank/balances/{address}/by_denom", utils.AutoWrap(s.aminoCodec, s.GetBalancesByAddressDenom))

	s.httpMux.HandleFunc("/bank/spendable_balances/{address}/by_denom", utils.AutoWrap(s.aminoCodec, s.GetSpendableBalancesByAddressDenom))
}

// GetBankParams queries the parameters of x/bank module.
func (s *Server) GetBankParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().Params(queryContext, &banktypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetSupplyByDenom queries the supply of a single coin.
func (s *Server) GetSupplyByDenom(req *getSupplyByDenomRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().SupplyOf(queryContext, &banktypes.QuerySupplyOfRequest{
		Denom: req.Denom,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, err
}

// GetBalancesByAddressDenom queries the balance of a single coin for a single account.
func (s *Server) GetBalancesByAddressDenom(req *getBalancesByAddressDenomRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().Balance(queryContext, &banktypes.QueryBalanceRequest{
		Address: bech32AccAddress.String(),
		Denom:   req.Denom,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetSpendableBalancesByAddressDenom queries the spendable balance of a single coin for a single account.
func (s *Server) GetSpendableBalancesByAddressDenom(req *getSpendableBalancesByAddressDenomRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	bech32AccAddress, err := utils.EvmAddressToBech32AccAddress(mux.Vars(r)["address"])
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().SpendableBalanceByDenom(queryContext, &banktypes.QuerySpendableBalanceByDenomRequest{
		Address: bech32AccAddress.String(),
		Denom:   req.Denom,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
