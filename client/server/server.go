//nolint:wrapcheck // The api server is our server, so we don't need to wrap it.
package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"cosmossdk.io/x/tx/signing"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/gogoproto/proto"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	evmstakingkeeper "github.com/piplabs/story/client/x/evmstaking/keeper"
	mintkeeper "github.com/piplabs/story/client/x/mint/keeper"
)

type Store interface {
	CreateQueryContext(height int64, prove bool) (sdk.Context, error)
	GetEvmStakingKeeper() *evmstakingkeeper.Keeper
	GetStakingKeeper() *stakingkeeper.Keeper
	GetSlashingKeeper() slashingkeeper.Keeper
	GetAccountKeeper() authkeeper.AccountKeeper
	GetBankKeeper() bankkeeper.Keeper
	GetDistrKeeper() distrkeeper.Keeper
	GetUpgradeKeeper() *upgradekeeper.Keeper
	GetMintKeeper() mintkeeper.Keeper
}

type Server struct {
	errChan chan error
	store   Store
	cl      rpcclient.Client

	httpMux    *mux.Router
	httpServer *http.Server
	protoCodec *codec.ProtoCodec
	aminoCodec *codec.LegacyAmino
}

func NewServer(store Store, cl rpcclient.Client, listenAddress string, enableUnsafeCORS bool) (*Server, error) {
	s := &Server{
		errChan: make(chan error),
		store:   store,
		httpMux: mux.NewRouter(),
		cl:      cl,
	}

	if err := s.registerCodec(); err != nil {
		return nil, err
	}
	s.registerHandle()

	var svrHandler http.Handler = s.httpMux
	if enableUnsafeCORS {
		svrHandler = handlers.CORS()(s.httpMux)
	}
	s.httpServer = &http.Server{
		Addr:              listenAddress,
		Handler:           svrHandler,
		ReadHeaderTimeout: 60 * time.Second,
	}

	return s, nil
}

func (s *Server) registerCodec() error {
	sdkConfig := sdk.GetConfig()
	reg, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          authcodec.NewBech32Codec(sdkConfig.GetBech32AccountAddrPrefix()),
			ValidatorAddressCodec: authcodec.NewBech32Codec(sdkConfig.GetBech32ValidatorAddrPrefix()),
		},
	})
	if err != nil {
		return err
	}

	s.protoCodec = codec.NewProtoCodec(reg)
	s.aminoCodec = codec.NewLegacyAmino()

	// IMPORTANT: register related types so that we could unpack values from Any.
	std.RegisterInterfaces(s.protoCodec.InterfaceRegistry())
	std.RegisterLegacyAminoCodec(s.aminoCodec)
	authtypes.RegisterInterfaces(s.protoCodec.InterfaceRegistry())
	authtypes.RegisterLegacyAminoCodec(s.aminoCodec)

	return nil
}

func (s *Server) prepareUnpackInterfaces(v codectypes.UnpackInterfacesMessage) error {
	err := codectypes.UnpackInterfaces(v, s.protoCodec)
	if err != nil {
		return err
	}

	return codectypes.UnpackInterfaces(v, codectypes.AminoJSONPacker{Cdc: s.aminoCodec.Amino})
}

func (s *Server) registerHandle() {
	s.initAuthRoute()
	s.initBankRoute()
	s.initCometBFTRoute()
	s.initDistributionRoute()
	s.initEvmStakingRoute()
	s.initSlashingRoute()
	s.initStakingRoute()
	s.initUpgradeRoute()
	s.initMintRoute()
}

func (s *Server) createQueryContextByHeader(r *http.Request) (sdk.Context, error) {
	height, err := strconv.ParseInt(r.Header.Get("X-Height"), 10, 64)
	if err != nil {
		height = 0
	}

	return s.store.CreateQueryContext(height, false)
}

func (s *Server) Start() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errChan <- err
		}
	}()

	select {
	case <-time.After(time.Second):
		return nil
	case err := <-s.errChan:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
