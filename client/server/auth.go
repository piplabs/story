//nolint:wrapcheck,dupl // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initAuthRoute() {
	s.httpMux.HandleFunc("/auth/account_info/{address}", utils.SimpleWrap(s.aminoCodec, s.GetAccountInfo))
	s.httpMux.HandleFunc("/auth/accounts", utils.AutoWrap(s.aminoCodec, s.GetAccounts))
	s.httpMux.HandleFunc("/auth/accounts/{address}", utils.SimpleWrap(s.aminoCodec, s.GetAccountsByAddress))
	s.httpMux.HandleFunc("/auth/address_by_id", utils.AutoWrap(s.aminoCodec, s.GetAccountAddressByID))
	s.httpMux.HandleFunc("/auth/bech32", utils.SimpleWrap(s.aminoCodec, s.GetBech32Prefix))
	s.httpMux.HandleFunc("/auth/bech32/{address_bytes}", utils.SimpleWrap(s.aminoCodec, s.Bech32AddressBytesToString))
	s.httpMux.HandleFunc("/auth/bech32/{address_string}", utils.SimpleWrap(s.aminoCodec, s.Bech32AddressStringToBytes))
	s.httpMux.HandleFunc("/auth/module_accounts", utils.SimpleWrap(s.aminoCodec, s.GetModuleAccounts))
	s.httpMux.HandleFunc("/auth/module_accounts/{name}", utils.SimpleWrap(s.aminoCodec, s.GetModuleAccountByName))
	s.httpMux.HandleFunc("/auth/params", utils.SimpleWrap(s.aminoCodec, s.GetAuthParams))
}

// GetAccountInfo returns account info which is common to all account types.
func (s *Server) GetAccountInfo(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).AccountInfo(queryContext, &authtypes.QueryAccountInfoRequest{
		Address: mux.Vars(r)["address"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetAccounts returns all the existing accounts.
func (s *Server) GetAccounts(req *getAccountsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).Accounts(queryContext, &authtypes.QueryAccountsRequest{
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

	for _, account := range queryResp.Accounts {
		if err := s.prepareUnpackInterfaces(utils.WrapTypeAny[types.AccountI](account)); err != nil {
			return nil, err
		}
	}

	return queryResp, nil
}

// GetAccountsByAddress returns account details based on address.
func (s *Server) GetAccountsByAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).Account(queryContext, &authtypes.QueryAccountRequest{
		Address: mux.Vars(r)["address"],
	})
	if err != nil {
		return nil, err
	}

	if err := s.prepareUnpackInterfaces(utils.WrapTypeAny[types.AccountI](queryResp.Account)); err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetAccountAddressByID returns account address based on account number.
func (s *Server) GetAccountAddressByID(req *getAccountAddressByIDRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).AccountAddressByID(queryContext, &authtypes.QueryAccountAddressByIDRequest{
		AccountId: req.AccountID,
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetBech32Prefix queries bech32Prefix.
func (s *Server) GetBech32Prefix(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).Bech32Prefix(queryContext, &authtypes.Bech32PrefixRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// Bech32AddressBytesToString converts Account Address bytes to string.
func (s *Server) Bech32AddressBytesToString(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).AddressBytesToString(queryContext, &authtypes.AddressBytesToStringRequest{
		AddressBytes: []byte(mux.Vars(r)["address_bytes"]),
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// Bech32AddressStringToBytes converts Address string to bytes.
func (s *Server) Bech32AddressStringToBytes(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).AddressStringToBytes(queryContext, &authtypes.AddressStringToBytesRequest{
		AddressString: mux.Vars(r)["address_string"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetModuleAccounts returns all the existing module accounts.
func (s *Server) GetModuleAccounts(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).ModuleAccounts(queryContext, &authtypes.QueryModuleAccountsRequest{})
	if err != nil {
		return nil, err
	}

	for _, account := range queryResp.Accounts {
		if err := s.prepareUnpackInterfaces(utils.WrapTypeAny[types.ModuleAccountI](account)); err != nil {
			return nil, err
		}
	}

	return queryResp, nil
}

// GetModuleAccountByName returns the module account info by module name.
func (s *Server) GetModuleAccountByName(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).ModuleAccountByName(queryContext, &authtypes.QueryModuleAccountByNameRequest{
		Name: mux.Vars(r)["name"],
	})
	if err != nil {
		return nil, err
	}

	if err := s.prepareUnpackInterfaces(utils.WrapTypeAny[types.ModuleAccountI](queryResp.Account)); err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetAuthParams queries all parameters of auth module.
func (s *Server) GetAuthParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQueryServer(s.store.GetAccountKeeper()).Params(queryContext, &authtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
