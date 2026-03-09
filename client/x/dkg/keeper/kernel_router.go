package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// KernelRouter manages multiple story-kernel clients, routing requests by code commitment.
type KernelRouter struct {
	mu        sync.RWMutex
	endpoints []string                             // configured endpoints
	clients   map[string]types.KernelServiceClient // codeCommitmentHex -> KernelServiceClient
	closers   map[string]io.Closer                 // codeCommitmentHex -> underlying gRPC connection
	ccByEP    map[string]string                    // endpoint -> codeCommitmentHex (reverse lookup)
}

const maxKernelEndpoints = 2

// NewKernelRouter creates a new router with the given endpoint list.
// At most 2 endpoints are supported (old + new binary for upgrade resharing).
func NewKernelRouter(endpoints []string) *KernelRouter {
	if len(endpoints) > maxKernelEndpoints {
		panic(fmt.Sprintf("kernel router supports at most %d endpoints, got %d", maxKernelEndpoints, len(endpoints)))
	}

	return &KernelRouter{
		endpoints: endpoints,
		clients:   make(map[string]types.KernelServiceClient),
		closers:   make(map[string]io.Closer),
		ccByEP:    make(map[string]string),
	}
}

// ConnectAndDiscover connects to an endpoint, calls GetCodeCommitment to discover
// the code commitment, and registers the client keyed by code commitment.
func (r *KernelRouter) ConnectAndDiscover(ctx context.Context, endpoint string) error {
	client, closer, err := CreateKernelClient(endpoint)
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
