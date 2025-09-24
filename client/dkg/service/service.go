package service

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/piplabs/story/client/dkg/config"
	"github.com/piplabs/story/client/dkg/events"
	dkgpb "github.com/piplabs/story/client/dkg/pb/v1"
	"github.com/piplabs/story/client/dkg/state"
	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// ServiceStatus represents the current status of the DKG service.
//
//nolint:revive // use full name
type ServiceStatus struct {
	IsRunning      bool           `json:"is_running"`
	ValidatorAddr  common.Address `json:"validator_address"`
	ActiveSessions []SessionInfo  `json:"active_sessions"`
}

// SessionInfo represents information about a DKG session.
type SessionInfo struct {
	Mrenclave  string    `json:"mrenclave"`
	Round      uint32    `json:"round"`
	Phase      string    `json:"phase"`
	StartTime  time.Time `json:"start_time"`
	LastUpdate time.Time `json:"last_update"`
}

// Service represents the main DKG service.
type Service struct {
	config           *config.DKGConfig
	eventListener    *events.EventListener
	stateManager     *state.StateManager
	teeClient        dkgpb.TEEClient
	cosmosClient     client.Context
	contractClient   *ContractClient
	stopChan         chan struct{}
	validatorAddress common.Address
}

// NewService creates a new DKG service.
func NewService(
	cfg *config.DKGConfig,
	teeClient dkgpb.TEEClient,
	cosmosClient client.Context,
	cosmosRPCEndpoint string,
	contractConfig *ContractConfig,
) (*Service, error) {
	if !cfg.Enable {
		return nil, errors.New("dkg service is disabled in configuration")
	}

	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid dkg configuration")
	}

	eventListener, err := events.NewEventListener(cosmosRPCEndpoint, cfg.EventPollingInterval)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create event listener")
	}

	stateManager, err := state.NewStateManager(cfg.DataDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create state manager")
	}

	var contractClient *ContractClient
	contractClient, err = NewContractClient(context.Background(), contractConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create contract client")
	}

	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to cast public key to ECDSA")
	}
	validatorAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Service{
		config:           cfg,
		eventListener:    eventListener,
		stateManager:     stateManager,
		teeClient:        teeClient,
		cosmosClient:     cosmosClient,
		contractClient:   contractClient,
		stopChan:         make(chan struct{}),
		validatorAddress: validatorAddress,
	}, nil
}

// Start starts the DKG service.
func (s *Service) Start(ctx context.Context) error {
	log.Info(ctx, "Starting DKG service", "validator", s.validatorAddress)

	// Start event listener
	if err := s.eventListener.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start event listener")
	}

	// Start main event processing loop
	go s.processEvents(ctx)

	// Start periodic cleanup
	go s.periodicCleanup(ctx)

	log.Info(ctx, "DKG service started successfully")

	return nil
}

// Stop stops the DKG service.
func (s *Service) Stop(ctx context.Context) error {
	log.Info(ctx, "Stopping DKG service")

	close(s.stopChan)

	if err := s.eventListener.Stop(ctx); err != nil {
		log.Error(ctx, "Error stopping event listener", err)
	}

	if s.contractClient != nil {
		s.contractClient.Close()
	}

	log.Info(ctx, "DKG service stopped")

	return nil
}

// processEvents is the main event processing loop.
func (s *Service) processEvents(ctx context.Context) {
	eventChan := s.eventListener.EventChannel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case event := <-eventChan:
			if err := s.handleDKGEvent(ctx, event); err != nil {
				log.Error(ctx, "Error handling DKG event", err,
					"event_type", event.EventType,
					"mrenclave", event.Mrenclave,
					"round", event.Round,
				)
			}
		}
	}
}

// handleDKGEvent handles a DKG-related blockchain event.
func (s *Service) handleDKGEvent(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Processing DKG event",
		"type", event.EventType,
		"mrenclave", event.Mrenclave,
		"round", event.Round,
		"height", event.BlockHeight,
	)

	switch event.EventType {
	// TODO: use enums here and in events/helper.go
	case "dkg_begin_initialization":
		return s.handleDKGInitialization(ctx, event)
	case "dkg_begin_network_set":
		return s.handleDKGNetworkSet(ctx, event)
	case "dkg_begin_dealing":
		return s.handleDKGDealing(ctx, event)
	case "dkg_verify_dealings":
		return s.handleDKGDealVerification(ctx, event)
	case "dkg_begin_finalization":
		return s.handleDKGFinalization(ctx, event)
	case "dkg_finalized":
		return s.handleDKGComplete(ctx, event)
	case "dkg_begin_resharing":
		return s.handleDKGResharing(ctx, event)
	case "dkg_begin_tee_upgrade":
		return s.handleTEEUpgrade(ctx, event)
	default:
		log.Warn(ctx, "Unknown DKG event type", nil, "type", event.EventType)
	}

	return nil
}

// periodicCleanup performs periodic cleanup of expired sessions.
func (s *Service) periodicCleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Hour) // Cleanup every hour
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.stateManager.CleanupExpiredSessions(ctx)
		}
	}
}

// GetStatus returns the current status of the DKG service.
func (s *Service) GetStatus() *ServiceStatus {
	sessions := s.stateManager.ListSessions()

	status := &ServiceStatus{
		IsRunning:      true,
		ValidatorAddr:  s.validatorAddress,
		ActiveSessions: make([]SessionInfo, 0, len(sessions)),
	}

	for _, session := range sessions {
		status.ActiveSessions = append(status.ActiveSessions, SessionInfo{
			Mrenclave:  session.GetMrenclaveString(),
			Round:      session.Round,
			Phase:      session.Phase.String(),
			StartTime:  session.StartTime,
			LastUpdate: session.LastUpdate,
		})
	}

	return status
}

// DKGEventHandler defines the interface for handling specific DKG events.
type DKGEventHandler interface {
	HandleEvent(ctx context.Context, service *Service, event *types.DKGEventData) error
	EventType() string
}

// EventHandlerRegistry manages the registration and execution of DKG event handlers.
type EventHandlerRegistry struct {
	handlers map[string]DKGEventHandler
}

// NewEventHandlerRegistry creates a new event handler registry.
func NewEventHandlerRegistry() *EventHandlerRegistry {
	return &EventHandlerRegistry{
		handlers: make(map[string]DKGEventHandler),
	}
}

// RegisterHandler registers a handler for a specific event type.
func (r *EventHandlerRegistry) RegisterHandler(handler DKGEventHandler) {
	r.handlers[handler.EventType()] = handler
}

// HandleEvent routes an event to the appropriate handler.
func (r *EventHandlerRegistry) HandleEvent(ctx context.Context, service *Service, event *types.DKGEventData) error {
	handler, exists := r.handlers[event.EventType]
	if !exists {
		log.Warn(ctx, "No handler registered for DKG event type", nil, "event_type", event.EventType)

		return nil
	}

	return handler.HandleEvent(ctx, service, event)
}
