package config

import (
	"encoding/binary"

	"github.com/spf13/pflag"

	"github.com/piplabs/story/lib/errors"
)

const DefaultEnclaveType uint64 = 1 // SGX

type DKGConfig struct {
	// Enable enables or disables the DKG client
	Enable bool

	// KernelEndpoints is the list of story-kernel (TEE) endpoint addresses
	KernelEndpoints []string

	// EngineRPCEndpoint is the RPC endpoint of the execution layer
	EngineRPCEndpoint string

	// EnclaveType is the TEE enclave type identifier (e.g. 1 for SGX), stored as bytes32 on-chain
	EnclaveType uint64

	// DecryptBatchSize is the number of partial decryptions submitted in a single
	// CDR contract call. Must be ≤ the CDR contract's maxBatchSize. Defaults to 20.
	DecryptBatchSize int

	// TLS configuration for gRPC client connections to story-kernel.
	// KernelTLSCAFile is the CA certificate to verify the server.
	// When set, TLS is used for all kernel connections.
	KernelTLSCAFile string

	// KernelTLSCertFile and KernelTLSKeyFile are the client certificate and key
	// for mutual TLS authentication. Both must be set for mTLS.
	KernelTLSCertFile string
	KernelTLSKeyFile  string
}

func DefaultDKGConfig() DKGConfig {
	return DKGConfig{
		Enable:            false,
		KernelEndpoints:   []string{"127.0.0.1:50051"},
		EngineRPCEndpoint: "http://127.0.0.1:8545",
		EnclaveType:       DefaultEnclaveType,
		DecryptBatchSize:  20,
	}
}

func BindDKGFlags(flags *pflag.FlagSet, cfg *DKGConfig) {
	flags.BoolVar(&cfg.Enable, "dkg-enable", cfg.Enable, "DKG client is enabled or not")
	flags.StringSliceVar(&cfg.KernelEndpoints, "dkg-kernel-endpoints", cfg.KernelEndpoints, "Comma-separated list of story-kernel (TEE) endpoints for DKG")
	flags.StringVar(&cfg.EngineRPCEndpoint, "dkg-engine-rpc-endpoint", cfg.EngineRPCEndpoint, "The RPC endpoint of execution layer")
	flags.Uint64Var(&cfg.EnclaveType, "dkg-enc-type", cfg.EnclaveType, "TEE enclave type identifier (e.g. 1 for SGX)")
	flags.IntVar(&cfg.DecryptBatchSize, "dkg-decrypt-batch-size", cfg.DecryptBatchSize, "Number of partial decryptions per CDR batch call (must be ≤ contract maxBatchSize)")
	flags.StringVar(&cfg.KernelTLSCAFile, "dkg-kernel-tls-ca-file", cfg.KernelTLSCAFile, "CA certificate file to verify story-kernel server TLS")
	flags.StringVar(&cfg.KernelTLSCertFile, "dkg-kernel-tls-cert-file", cfg.KernelTLSCertFile, "Client certificate file for mTLS to story-kernel")
	flags.StringVar(&cfg.KernelTLSKeyFile, "dkg-kernel-tls-key-file", cfg.KernelTLSKeyFile, "Client private key file for mTLS to story-kernel")
}

func (c *DKGConfig) Validate() error {
	if !c.Enable {
		return nil
	}

	if len(c.KernelEndpoints) == 0 {
		return errors.New("at least one kernel endpoint is required")
	}

	if len(c.KernelEndpoints) > 2 {
		return errors.New("at most 2 kernel endpoints are supported (old + new binary for upgrade)")
	}

	if c.EngineRPCEndpoint == "" {
		return errors.New("engine rpc endpoint should not be empty")
	}

	if c.EnclaveType == 0 {
		return errors.New("enc-type must not be zero")
	}

	if c.DecryptBatchSize <= 0 {
		return errors.New("dkg-decrypt-batch-size must be > 0")
	}

	// Validate TLS configuration consistency.
	hasCert := c.KernelTLSCertFile != ""
	hasKey := c.KernelTLSKeyFile != ""

	if hasCert != hasKey {
		return errors.New("dkg-kernel-tls-cert-file and dkg-kernel-tls-key-file must both be set or both be empty")
	}

	// Client cert/key without CA file is not useful — the CA file is needed
	// to verify the server, and the cert/key are only for mTLS.
	if hasCert && c.KernelTLSCAFile == "" {
		return errors.New("dkg-kernel-tls-ca-file is required when client certificate is configured for mTLS")
	}

	return nil
}

// EnclaveTypeToBytes32 converts a uint64 enclave type to a [32]byte (big-endian, right-aligned).
func EnclaveTypeToBytes32(v uint64) [32]byte {
	var result [32]byte
	binary.BigEndian.PutUint64(result[24:], v)

	return result
}
