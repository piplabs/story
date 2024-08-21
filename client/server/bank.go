//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initBankRoute() {
	s.httpMux.HandleFunc("/bank/supply/by_denom", utils.AutoWrap(s.aminoCodec, s.GetSupplyByDenom))

	s.httpMux.HandleFunc("/bank/balances/{address}", utils.AutoWrap(s.aminoCodec, s.GetBalancesByAddress))
	s.httpMux.HandleFunc("/bank/balances/{address}/by_denom", utils.AutoWrap(s.aminoCodec, s.GetBalancesByAddressDenom))
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

// GetBalancesByAddress queries the balance of all coins for a single account.
func (s *Server) GetBalancesByAddress(req *getBalancesByAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().AllBalances(queryContext, &banktypes.QueryAllBalancesRequest{
		Address: mux.Vars(r)["address"],
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
		ResolveDenom: req.ResolveDenom,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetBalancesByAddressDenom queries the balance of a single coin for a single account.
func (s *Server) GetBalancesByAddressDenom(req *getBalancesByAddressDenomRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().Balance(queryContext, &banktypes.QueryBalanceRequest{
		Address: mux.Vars(r)["address"],
		Denom:   req.Denom,
	})

	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
