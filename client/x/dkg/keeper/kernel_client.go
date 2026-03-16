package keeper

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"
	"strings"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// TLSConfig holds the TLS certificate paths for gRPC client connections.
// When CAFile is set, TLS is enabled with server certificate verification.
// When CertFile and KeyFile are also set, mutual TLS (mTLS) is used.
type TLSConfig struct {
	CAFile   string // CA certificate to verify the server
	CertFile string // Client certificate for mTLS
	KeyFile  string // Client private key for mTLS
}

// CreateKernelClient creates a gRPC client for the story-kernel.
// Returns the KernelClient and an io.Closer for the underlying gRPC connection.
//
// TLS behavior:
//   - tlsCfg nil or empty: insecure connection (no TLS)
//   - tlsCfg.CAFile set: TLS with server cert verification against the provided CA
//   - tlsCfg.CertFile+KeyFile also set: mutual TLS (client presents its certificate)
func CreateKernelClient(endpoint string, tlsCfg *TLSConfig) (types.KernelServiceClient, io.Closer, error) {
	if endpoint == "" {
		return nil, nil, errors.New("the endpoint is required")
	}

	creds, err := buildTransportCredentials(endpoint, tlsCfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build TLS credentials", "endpoint", endpoint)
	}

	// Strip scheme prefix (tls://, https://) for gRPC target resolution.
	// gRPC does not understand these schemes — it uses passthrough resolver for direct IP:port.
	target := endpoint
	for _, prefix := range []string{"tls://", "https://"} {
		if strings.HasPrefix(target, prefix) {
			target = strings.TrimPrefix(target, prefix)

			break
		}
	}

	if !strings.Contains(target, "://") {
		target = "passthrough:///" + target
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to connect to story-kernel client")
	}

	return types.NewKernelServiceClient(conn), conn, nil
}

// buildTransportCredentials returns the appropriate gRPC transport credentials
// based on the endpoint scheme and TLS configuration.
func buildTransportCredentials(endpoint string, tlsCfg *TLSConfig) (credentials.TransportCredentials, error) {
	// Explicit TLS config takes precedence over endpoint scheme detection.
	if tlsCfg != nil && tlsCfg.CAFile != "" {
		return loadClientTLSCredentials(tlsCfg)
	}

	// Legacy behavior: detect TLS from endpoint URL scheme.
	if strings.HasPrefix(endpoint, "https://") || strings.HasPrefix(endpoint, "tls://") {
		// Use system CA pool for server verification when no explicit CA is provided.
		return credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		}), nil
	}

	return insecure.NewCredentials(), nil
}

// loadClientTLSCredentials loads TLS credentials for gRPC client connections.
// Requires a CA file for server verification. Optionally loads client cert/key for mTLS.
func loadClientTLSCredentials(tlsCfg *TLSConfig) (credentials.TransportCredentials, error) {
	caCert, err := os.ReadFile(tlsCfg.CAFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CA cert file", "path", tlsCfg.CAFile)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to parse CA cert", "path", tlsCfg.CAFile)
	}

	tlsConfig := &tls.Config{
		RootCAs:    caPool,
		MinVersion: tls.VersionTLS12,
	}

	// Load client certificate for mTLS if both cert and key are provided.
	if tlsCfg.CertFile != "" && tlsCfg.KeyFile != "" {
		clientCert, err := tls.LoadX509KeyPair(tlsCfg.CertFile, tlsCfg.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load client cert/key",
				"cert", tlsCfg.CertFile, "key", tlsCfg.KeyFile)
		}

		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	return credentials.NewTLS(tlsConfig), nil
}
