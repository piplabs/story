//nolint:wrapcheck // The api server is our server, so we don't need to wrap it
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/piplabs/story/client/server/utils"
)

func (s *Server) initAuthRoute() {
	s.httpMux.HandleFunc("/auth/params", utils.SimpleWrap(s.aminoCodec, s.GetAuthParams))

	s.httpMux.HandleFunc("/auth/accounts", utils.AutoWrap(s.aminoCodec, s.GetAccounts))

	s.httpMux.HandleFunc("/auth/bech32", utils.SimpleWrap(s.aminoCodec, s.GetBech32Prefix))
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
