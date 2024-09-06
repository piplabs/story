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
	s.httpMux.HandleFunc("bank/params", utils.SimpleWrap(s.aminoCodec, s.GetBankParams))

	s.httpMux.HandleFunc("/bank/supply", utils.AutoWrap(s.aminoCodec, s.GetSupply))
	s.httpMux.HandleFunc("/bank/supply/by_denom", utils.AutoWrap(s.aminoCodec, s.GetSupplyByDenom))

	s.httpMux.HandleFunc("/bank/balances/{address}", utils.AutoWrap(s.aminoCodec, s.GetBalancesByAddress))
	s.httpMux.HandleFunc("/bank/balances/{address}/by_denom", utils.AutoWrap(s.aminoCodec, s.GetBalancesByAddressDenom))

	s.httpMux.HandleFunc("/bank/denom_owners/{denom}", utils.AutoWrap(s.aminoCodec, s.GetDenomOwners))
	s.httpMux.HandleFunc("/bank/denom_owners_by_query", utils.AutoWrap(s.aminoCodec, s.GetDenomOwnersByQuery))
	s.httpMux.HandleFunc("/bank/denoms_metadata", utils.AutoWrap(s.aminoCodec, s.GetDenomsMetadata))
	s.httpMux.HandleFunc("/bank/denoms_metadata/{denom}", utils.SimpleWrap(s.aminoCodec, s.GetDenomMetadata))
	s.httpMux.HandleFunc("/bank/denoms_metadata_by_query_string", utils.AutoWrap(s.aminoCodec, s.GetDenomMetadataByQueryString))

	s.httpMux.HandleFunc("/bank/send_enabled", utils.AutoWrap(s.aminoCodec, s.GetSendEnabled))
	s.httpMux.HandleFunc("/bank/spendable_balances/{address}", utils.AutoWrap(s.aminoCodec, s.GetSpendableBalances))
	s.httpMux.HandleFunc("/bank/spendable_balances/{address}/by_denom", utils.AutoWrap(s.aminoCodec, s.GetSpendableBalanceByDenom))
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

// GetSupply queries the total supply of all coins.
func (s *Server) GetSupply(req *getSupplyRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().TotalSupply(queryContext, &banktypes.QueryTotalSupplyRequest{
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

	return queryResp, err
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

// GetDenomOwners queries for all account addresses that own a particular token denomination.
func (s *Server) GetDenomOwners(req *getDenomOwnersRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().DenomOwners(queryContext, &banktypes.QueryDenomOwnersRequest{
		Denom: mux.Vars(r)["denom"],
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

	return queryResp, err
}

// GetDenomOwnersByQuery queries for all account addresses that own a particular token denomination.
func (s *Server) GetDenomOwnersByQuery(req *getDenomOwnersByQueryRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().DenomOwnersByQuery(queryContext, &banktypes.QueryDenomOwnersByQueryRequest{
		Denom: req.Denom,
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

	return queryResp, err
}

// GetDenomsMetadata queries the client metadata for all registered coin denominations.
func (s *Server) GetDenomsMetadata(req *getDenomsMetadataRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().DenomsMetadata(queryContext, &banktypes.QueryDenomsMetadataRequest{
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

	return queryResp, err
}

// GetDenomMetadata queries the client metadata of a given coin denomination.
func (s *Server) GetDenomMetadata(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().DenomMetadata(queryContext, &banktypes.QueryDenomMetadataRequest{
		Denom: mux.Vars(r)["denom"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, err
}

// GetDenomMetadataByQueryString queries the client metadata of a given coin denomination.
func (s *Server) GetDenomMetadataByQueryString(req *getDenomMetadataByQueryStringRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().DenomMetadataByQueryString(queryContext, &banktypes.QueryDenomMetadataByQueryStringRequest{
		Denom: req.Denom,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, err
}

// GetSendEnabled queries for SendEnabled entries.
func (s *Server) GetSendEnabled(req *getSendEnabledRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().SendEnabled(queryContext, &banktypes.QuerySendEnabledRequest{
		Denoms: req.Denoms,
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

	return queryResp, err
}

// GetSpendableBalances queries the spendable balance of all coins for a single account.
func (s *Server) GetSpendableBalances(req *getSpendableBalancesRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().SpendableBalances(queryContext, &banktypes.QuerySpendableBalancesRequest{
		Address: mux.Vars(r)["address"],
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

	return queryResp, err
}

// GetSpendableBalanceByDenom queries the spendable balance of a single denom for a single account.
func (s *Server) GetSpendableBalanceByDenom(req *getSpendableBalanceByDenomRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetBankKeeper().SpendableBalanceByDenom(queryContext, &banktypes.QuerySpendableBalanceByDenomRequest{
		Address: mux.Vars(r)["address"],
		Denom:   req.Denom,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, err
}
