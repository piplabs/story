package service

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/piplabs/story/client/dkg/config"
	dkgpb "github.com/piplabs/story/client/dkg/pb/v1"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func NewDKGServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dkg-service",
		Short: "Run the DKG (Distributed Key Generation) service",
		Long: `Run the DKG service that participates in distributed key generation.

This service runs alongside the validator node and listens to blockchain events
to participate in DKG operations when enabled. It communicates with TEE clients
and manages DKG state outside of consensus.

The service is opt-in and only runs when --dkg-enable=true is set.`,
		RunE: runDKGService,
	}

	addDKGFlags(cmd)

	return cmd
}

func addDKGFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("dkg-enable", false, "Enable DKG service")
	cmd.Flags().String("dkg-tee-endpoint", "", "TEE client endpoint")
	cmd.Flags().String("dkg-cosmos-rpc-endpoint", "http://localhost:26657", "Cosmos RPC endpoint")
	// TODO: reuse eth rpc endpoint from running the validator node
	cmd.Flags().String("dkg-eth-rpc-endpoint", "http://localhost:8545", "Ethereum RPC endpoint for DKG contract")
	cmd.Flags().Duration("dkg-event-polling-interval", config.DefaultDKGConfig().EventPollingInterval, "Event polling interval")
	cmd.Flags().Int("dkg-max-tee-retries", config.DefaultDKGConfig().MaxTEERetries, "Maximum TEE operation retries")
	cmd.Flags().Duration("dkg-tee-timeout", config.DefaultDKGConfig().TEETimeout, "TEE operation timeout")
	cmd.Flags().String("dkg-data-dir", config.DefaultDKGConfig().DataDir, "DKG data directory")
	cmd.Flags().String("dkg-log-level", config.DefaultDKGConfig().LogLevel, "DKG service log level")

	cmd.Flags().String("dkg-contract-address", "", "DKG contract address")
	// TODO: reuse private key from running the validator node
	cmd.Flags().String("dkg-private-key", "", "Private key for contract transactions (hex-encoded)")
	cmd.Flags().Int64("dkg-chain-id", 31337, "Ethereum chain ID")
}

func runDKGService(cmd *cobra.Command, _ []string) error {
	cfg, err := getDKGConfigFromFlags(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG configuration")
	}

	if !cfg.Enable {
		log.Info(cmd.Context(), "DKG service is disabled")
		return nil
	}

	if err := cfg.Validate(); err != nil {
		return errors.Wrap(err, "invalid DKG configuration")
	}

	cosmosClient, err := createClientContext(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create cosmos client context")
	}

	teeClient, err := createTEEClient(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create TEE client")
	}

	contractConfig, err := getContractConfigFromFlags(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get contract configuration")
	}

	service, err := NewService(cfg, teeClient, cosmosClient, cfg.CosmosRPCEndpoint, contractConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create DKG service")
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := service.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start DKG service")
	}

	log.Info(ctx, "DKG service is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		log.Info(ctx, "Context canceled, shutting down")
	case sig := <-sigChan:
		log.Info(ctx, "Received signal, shutting down", "signal", sig)
		cancel()
	}

	if err := service.Stop(ctx); err != nil {
		log.Error(ctx, "Error stopping DKG service", err)

		return err
	}

	log.Info(ctx, "DKG service stopped successfully")

	return nil
}

func getContractConfigFromFlags(cmd *cobra.Command) (*ContractConfig, error) {
	ethRPCEndpoint, err := cmd.Flags().GetString("dkg-eth-rpc-endpoint")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get eth-rpc-endpoint")
	}

	contractAddr, err := cmd.Flags().GetString("dkg-contract-address")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get contract-address")
	}

	// TODO: reuse private key used by the validator node for EL interactions
	privateKey, err := cmd.Flags().GetString("dkg-private-key")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get private-key")
	}

	chainID, err := cmd.Flags().GetInt64("dkg-chain-id")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain-id")
	}

	if contractAddr == "" || privateKey == "" {
		return nil, errors.New("contract address and private key are required for contract interaction")
	}

	return &ContractConfig{
		EthRPCEndpoint:  ethRPCEndpoint,
		DKGContractAddr: contractAddr,
		PrivateKey:      privateKey,
		ChainID:         chainID,
	}, nil
}

func getDKGConfigFromFlags(cmd *cobra.Command) (*config.DKGConfig, error) {
	cfg := config.DefaultDKGConfig()

	var err error

	if cfg.Enable, err = cmd.Flags().GetBool("dkg-enable"); err != nil {
		return nil, errors.Wrap(err, "get dkg-enable")
	}

	if cfg.TEEEndpoint, err = cmd.Flags().GetString("dkg-tee-endpoint"); err != nil {
		return nil, errors.Wrap(err, "get dkg-tee-endpoint")
	}

	if cfg.CosmosRPCEndpoint, err = cmd.Flags().GetString("dkg-cosmos-rpc-endpoint"); err != nil {
		return nil, errors.Wrap(err, "get dkg-cosmos-rpc-endpoint")
	}

	if cfg.EventPollingInterval, err = cmd.Flags().GetDuration("dkg-event-polling-interval"); err != nil {
		return nil, errors.Wrap(err, "get dkg-event-polling-interval")
	}

	if cfg.MaxTEERetries, err = cmd.Flags().GetInt("dkg-max-tee-retries"); err != nil {
		return nil, errors.Wrap(err, "get dkg-max-tee-retries")
	}

	if cfg.TEETimeout, err = cmd.Flags().GetDuration("dkg-tee-timeout"); err != nil {
		return nil, errors.Wrap(err, "get dkg-tee-timeout")
	}

	if cfg.DataDir, err = cmd.Flags().GetString("dkg-data-dir"); err != nil {
		return nil, errors.Wrap(err, "get dkg-data-dir")
	}

	if cfg.LogLevel, err = cmd.Flags().GetString("dkg-log-level"); err != nil {
		return nil, errors.Wrap(err, "get dkg-log-level")
	}

	return cfg, nil
}

func createClientContext(cfg *config.DKGConfig) (client.Context, error) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	protoCodec := codec.NewProtoCodec(interfaceRegistry)

	rpcClient, err := rpchttp.New(cfg.CosmosRPCEndpoint, "/websocket")
	if err != nil {
		return client.Context{}, errors.Wrap(err, "create RPC client")
	}

	homeDir := cfg.DataDir
	if !filepath.IsAbs(homeDir) {
		homeDir = filepath.Join(os.Getenv("HOME"), ".story", cfg.DataDir)
	}

	kr, err := keyring.New("story", keyring.BackendTest, homeDir, os.Stdin, protoCodec)
	if err != nil {
		return client.Context{}, errors.Wrap(err, "create keyring")
	}

	clientCtx := client.Context{}.
		WithCodec(protoCodec).
		WithInterfaceRegistry(interfaceRegistry).
		WithLegacyAmino(codec.NewLegacyAmino()).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithHomeDir(homeDir).
		WithKeyring(kr).
		WithClient(rpcClient).
		WithChainID(cfg.CosmosChainID)

	return clientCtx, nil
}

func createTEEClient(cfg *config.DKGConfig) (dkgpb.TEEClient, error) {
	if cfg.TEEEndpoint == "" {
		return nil, errors.New("TEE endpoint is required")
	}

	var creds credentials.TransportCredentials
	if strings.HasPrefix(cfg.TEEEndpoint, "https://") || strings.HasPrefix(cfg.TEEEndpoint, "tls://") {
		creds = credentials.NewTLS(nil) // TODO: use provided CA certs
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(cfg.TEEEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to TEE service")
	}

	return dkgpb.NewTEEClient(conn), nil
}
