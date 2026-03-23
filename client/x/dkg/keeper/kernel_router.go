package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// reconnectBackoff tracks per-endpoint exponential backoff state for kernel reconnection.
type reconnectBackoff struct {
	lastAttempt     time.Time
	backoffDuration time.Duration
}

// KernelRouter manages multiple story-kernel clients, routing requests by code commitment.
type KernelRouter struct {
	mu        sync.RWMutex
	endpoints []string                             // configured endpoints
	tlsCfg    *TLSConfig                           // TLS configuration for client connections (nil = insecure)
	clients   map[string]types.KernelServiceClient // codeCommitmentHex -> KernelServiceClient
	closers   map[string]io.Closer                 // codeCommitmentHex -> underlying gRPC connection
	ccByEP    map[string]string                    // endpoint -> codeCommitmentHex (reverse lookup)
	backoffs  map[string]*reconnectBackoff          // endpoint -> backoff state for reconnection rate limiting
}

const maxKernelEndpoints = 2

// Reconnection backoff constants.
const (
	initialBackoff = 30 * time.Second
	maxBackoff     = 5 * time.Minute
)

// NewKernelRouter creates a new router with the given endpoint list and optional TLS configuration.
// At most 2 endpoints are supported (old + new binary for upgrade resharing).
// Pass nil for tlsCfg to use insecure connections.
func NewKernelRouter(endpoints []string, tlsCfg *TLSConfig) *KernelRouter {
	if len(endpoints) > maxKernelEndpoints {
		panic(fmt.Sprintf("kernel router supports at most %d endpoints, got %d", maxKernelEndpoints, len(endpoints)))
	}

	return &KernelRouter{
		endpoints: endpoints,
		tlsCfg:    tlsCfg,
		clients:   make(map[string]types.KernelServiceClient),
		closers:   make(map[string]io.Closer),
		ccByEP:    make(map[string]string),
		backoffs:  make(map[string]*reconnectBackoff),
	}
}

// ConnectAndDiscover connects to an endpoint, calls GetCodeCommitment to discover
// the code commitment, and registers the client keyed by code commitment.
func (r *KernelRouter) ConnectAndDiscover(ctx context.Context, endpoint string) error {
	client, closer, err := CreateKernelClient(endpoint, r.tlsCfg)
	if err != nil {
		return errors.Wrap(err, "failed to connect to kernel endpoint", "endpoint", endpoint)
	}

	resp, err := client.GetCodeCommitment(ctx, &types.GetCodeCommitmentRequest{})
	if err != nil {
		// Close connection on failure to avoid leak
		if closer != nil {
			_ = closer.Close()
		}

		return errors.Wrap(err, "failed to get code commitment from kernel endpoint", "endpoint", endpoint)
	}

	if len(resp.CodeCommitment) == 0 {
		if closer != nil {
			_ = closer.Close()
		}

		return errors.New("empty code commitment from kernel endpoint", "endpoint", endpoint)
	}

	codeCommitmentHex := hex.EncodeToString(resp.CodeCommitment)

	r.mu.Lock()
	defer r.mu.Unlock()

	// Close existing connection if re-registering same code commitment
	if oldCloser, ok := r.closers[codeCommitmentHex]; ok {
		_ = oldCloser.Close()
	}

	r.ccByEP[endpoint] = codeCommitmentHex

	r.clients[codeCommitmentHex] = client
	if closer != nil {
		r.closers[codeCommitmentHex] = closer
	}

	log.Info(ctx, "Connected to kernel endpoint",
		"endpoint", endpoint,
		"code_commitment", codeCommitmentHex,
	)

	return nil
}

// RegisterClientForEndpoint maps a discovered code commitment to an endpoint's client.
func (r *KernelRouter) RegisterClientForEndpoint(endpoint string, codeCommitment []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	codeCommitmentHex := hex.EncodeToString(codeCommitment)

	// Move the client from endpoint key to code commitment key
	if client, ok := r.clients[endpoint]; ok {
		delete(r.clients, endpoint)
		r.clients[codeCommitmentHex] = client
		r.ccByEP[endpoint] = codeCommitmentHex
	}
}

// RegisterClient maps a code commitment to a kernel client.
func (r *KernelRouter) RegisterClient(codeCommitment []byte, client types.KernelServiceClient) {
	r.mu.Lock()
	defer r.mu.Unlock()

	codeCommitmentHex := hex.EncodeToString(codeCommitment)
	r.clients[codeCommitmentHex] = client
}

// GetClient returns the kernel client for the given code commitment.
// Returns an error if the exact code commitment is not found — no fallback.
func (r *KernelRouter) GetClient(codeCommitment []byte) (types.KernelServiceClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(codeCommitment) == 0 {
		return nil, errors.New("empty code commitment")
	}

	codeCommitmentHex := hex.EncodeToString(codeCommitment)

	client, ok := r.clients[codeCommitmentHex]
	if !ok {
		return nil, errors.New("no kernel client for code commitment",
			"code_commitment", codeCommitmentHex,
		)
	}

	return client, nil
}

// GetAllCodeCommitments returns all connected code commitments.
func (r *KernelRouter) GetAllCodeCommitments() [][]byte {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var ccs [][]byte

	for codeCommitmentHex := range r.clients {
		cc, err := hex.DecodeString(codeCommitmentHex)
		if err != nil {
			continue
		}

		ccs = append(ccs, cc)
	}

	return ccs
}

// Disconnect removes a client by code commitment and closes its underlying gRPC connection.
func (r *KernelRouter) Disconnect(codeCommitment []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	codeCommitmentHex := hex.EncodeToString(codeCommitment)
	if closer, ok := r.closers[codeCommitmentHex]; ok {
		_ = closer.Close()

		delete(r.closers, codeCommitmentHex)
	}

	delete(r.clients, codeCommitmentHex)
}

// HasClients returns true if at least one client is connected.
func (r *KernelRouter) HasClients() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.clients) > 0
}

// reconnectTimeout is the maximum duration for a single kernel reconnection
// attempt (gRPC dial + GetCodeCommitment call).
const reconnectTimeout = 10 * time.Second

// TryReconnect attempts ConnectAndDiscover for any configured endpoints that
// do not yet have an active client connection. It creates its own background
// context with a short timeout so that reconnection is never tied to the
// caller's context (e.g., CometBFT BeginBlocker which gets canceled after
// block processing). Each endpoint is attempted once; failures are logged
// but not returned so that the caller can continue with whatever clients
// are available.
func (r *KernelRouter) TryReconnect() {
	disconnected := r.disconnectedEndpoints()
	if len(disconnected) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), reconnectTimeout)
	defer cancel()

	now := time.Now()

	for _, ep := range disconnected {
		if r.isInCooldown(ep, now) {
			log.Debug(ctx, "Skipping kernel reconnection (in cooldown)", "endpoint", ep)

			continue
		}

		log.Info(ctx, "Attempting kernel reconnection", "endpoint", ep)

		if err := r.ConnectAndDiscover(ctx, ep); err != nil {
			log.Warn(ctx, "Kernel reconnection failed", err, "endpoint", ep)
			r.recordFailedAttempt(ep, now)
		} else {
			r.resetBackoff(ep)
		}
	}
}

// isInCooldown returns true if the endpoint was attempted recently and the
// backoff cooldown has not yet elapsed.
func (r *KernelRouter) isInCooldown(endpoint string, now time.Time) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bo, ok := r.backoffs[endpoint]
	if !ok {
		return false
	}

	return now.Before(bo.lastAttempt.Add(bo.backoffDuration))
}

// recordFailedAttempt updates the backoff state for a failed reconnection attempt,
// doubling the duration up to maxBackoff.
func (r *KernelRouter) recordFailedAttempt(endpoint string, now time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	bo, ok := r.backoffs[endpoint]
	if !ok {
		r.backoffs[endpoint] = &reconnectBackoff{
			lastAttempt:     now,
			backoffDuration: initialBackoff,
		}

		return
	}

	bo.lastAttempt = now
	bo.backoffDuration *= 2
	if bo.backoffDuration > maxBackoff {
		bo.backoffDuration = maxBackoff
	}
}

// resetBackoff clears the backoff state for an endpoint after a successful connection.
func (r *KernelRouter) resetBackoff(endpoint string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.backoffs, endpoint)
}

// disconnectedEndpoints returns configured endpoints that do not have a
// corresponding entry in the ccByEP map (i.e. no successful connection yet).
func (r *KernelRouter) disconnectedEndpoints() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []string
	for _, ep := range r.endpoints {
		if _, ok := r.ccByEP[ep]; !ok {
			out = append(out, ep)
		}
	}

	return out
}
