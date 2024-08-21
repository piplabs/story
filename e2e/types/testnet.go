package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"net"
	"strings"
	"sync/atomic"

	e2e "github.com/cometbft/cometbft/test/e2e/pkg"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/storyprotocol/iliad/lib/evmchain"
	"github.com/storyprotocol/iliad/lib/netconf"
)

// Testnet wraps e2e.Testnet with additional iliad-specific fields.
type Testnet struct {
	*e2e.Testnet
	Network      netconf.ID
	IliadEVMs    []IliadEVM
	AnvilChains  []AnvilChain
	PublicChains []PublicChain
	Perturb      map[string][]Perturb
}

// RandomIliadAddr returns a random iliad address for cometBFT rpc clients.
// It uses the internal IP address of a random node that isn't delayed or a seed.
func (t Testnet) RandomIliadAddr() string {
	var eligible []string
	for _, node := range t.Nodes {
		if node.StartAt != 0 || node.Mode == ModeSeed {
			continue // Skip delayed nodes or seed nodes
		}

		eligible = append(eligible, node.AddressRPC())
	}

	if len(eligible) == 0 {
		return ""
	}

	randIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(eligible))))
	if err != nil {
		return ""
	}

	return eligible[randIdx.Uint64()]
}

// BroadcastIliadEVM returns a Iliad EVM to use for e2e app tx broadcasts.
// It prefers a validator nodes since we have an issue with mempool+p2p+startup where
// txs get stuck in non-validator mempool immediately after startup if not connected to peers yet.
// Also avoid validators that are not started immediately.
func (t Testnet) BroadcastIliadEVM() IliadEVM {
	isDelayed := func(evm string) bool {
		for _, node := range t.Nodes {
			if node.StartAt > 0 && strings.Contains(evm, node.Name) {
				return true
			}
		}

		return false
	}

	for _, evm := range t.IliadEVMs {
		if strings.Contains(evm.InstanceName, "validator") && !isDelayed(evm.InstanceName) {
			return evm
		}
	}

	return t.IliadEVMs[0]
}

// BroadcastNode returns a iliad node to use for RPC queries broadcasts.
// It prefers a validator nodes since we have an issue with mempool+p2p+startup where
// txs get stuck in non-validator mempool immediately after startup if not connected to peers yet.
// Also avoid validators that are not started immediately.
func (t Testnet) BroadcastNode() *e2e.Node {
	for _, node := range t.Nodes {
		if !strings.Contains(node.Name, "validator") {
			continue
		}
		if node.StartAt > 0 {
			continue
		}

		return node
	}

	return t.Nodes[0]
}

// HasPerturbations returns whether the network has any perturbations.
func (t Testnet) HasPerturbations() bool {
	if len(t.Perturb) > 0 {
		return true
	}

	return t.Testnet.HasPerturbations()
}

func (t Testnet) HasIliadEVM() bool {
	return len(t.IliadEVMs) > 0
}

// EVMChain represents a EVM chain in a iliad network.
type EVMChain struct {
	evmchain.Metadata
	IsPublic bool
}

// IliadEVM represents a iliad evm instance in a iliad network. Similar to e2e.Node for iliad instances.
type IliadEVM struct {
	Chain        EVMChain // For netconf (all instances must have the same chain)
	InstanceName string   // For docker container name
	AdvertisedIP net.IP   // For setting up NAT on geth bootnode
	ProxyPort    uint32   // For binding
	InternalRPC  string   // For JSON-RPC queries from client
	ExternalRPC  string   // For JSON-RPC queries from e2e app.
	IsArchive    bool     // Whether this instance is in archive mode
	JWTSecret    string   // JWT secret for authentication

	// P2P networking
	NodeKey *ecdsa.PrivateKey // Private key
	Enode   *enode.Node       // Public key
	Peers   []*enode.Node     // Peer public keys
}

// NodeKeyHex returns the hex-encoded node key. Used for geth's config.
func (o IliadEVM) NodeKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSA(o.NodeKey))
}

// AnvilChain represents an anvil chain instance in a iliad network.
type AnvilChain struct {
	Chain       EVMChain // For netconf
	InternalIP  net.IP   // For docker container IP
	ProxyPort   uint32   // For binding
	InternalRPC string   // For JSON-RPC queries from client
	ExternalRPC string   // For JSON-RPC queries from e2e app.
	LoadState   string   // File path to load anvil state from
}

// PublicChain represents a public chain in a iliad network.
type PublicChain struct {
	chain        EVMChain      // For netconf
	rpcAddresses []string      // For JSON-RPC queries from client/e2e app.
	next         *atomic.Int32 // For round-robin RPC address selection
}

func NewPublicChain(chain EVMChain, rpcAddresses []string) PublicChain {
	return PublicChain{
		chain:        chain,
		rpcAddresses: rpcAddresses,
		next:         new(atomic.Int32),
	}
}

// Chain returns the EVM chain.
func (c PublicChain) Chain() EVMChain {
	return c.chain
}

// NextRPCAddress returns the next RPC address in the list.
func (c PublicChain) NextRPCAddress() string {
	i := c.next.Load()
	defer func() {
		c.next.Store(i + 1)
	}()

	l := len(c.rpcAddresses)
	if l == 0 {
		return ""
	}

	return strings.TrimSpace(c.rpcAddresses[int(i)%l])
}
