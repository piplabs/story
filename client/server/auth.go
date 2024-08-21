//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"

	"github.com/storyprotocol/iliad/client/server/utils"
)

func (s *Server) initAuthRoute() {
	s.httpMux.HandleFunc("/auth/accounts", utils.AutoWrap(s.aminoCodec, s.GetAccounts))
	s.httpMux.HandleFunc("/auth/accounts/{address}", utils.SimpleWrap(s.aminoCodec, s.GetAccountsByAddress))
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
		err = s.prepareUnpackInterfaces(utils.WrapTypeAny[types.AccountI](account))
		if err != nil {
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

	err = s.prepareUnpackInterfaces(utils.WrapTypeAny[types.AccountI](queryResp.Account))
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
