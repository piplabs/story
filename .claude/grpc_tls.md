# H-01: gRPC TLS Between Story CL and story-kernel

**Priority**: Low (future PR) — This document describes a planned security enhancement.

## Problem

The Story consensus layer (CL) communicates with `story-kernel` (TEE service) via gRPC. Currently, this connection uses **plaintext (insecure) transport** by default when the endpoint does not carry a `https://` or `tls://` prefix.

While the two processes typically run on the same machine (`127.0.0.1`), this setup presents security concerns when they are deployed on separate hosts or when defence-in-depth is desired:

| Threat | Description |
|--------|-------------|
| **Network-level eavesdropping** | An attacker on the same network can passively capture gRPC traffic containing DKG key shares, VSS deals, attestation reports, and partial decryption results. |
| **Man-in-the-middle (MITM) attacks** | Without server (or mutual) authentication, an attacker can impersonate `story-kernel` and serve malicious responses — e.g., returning a crafted public key or forged attestation. |
| **Credential / key-share theft in transit** | Sealed key material is transmitted from TEE to the CL during dealing and finalization. Plaintext transport exposes these shares to interception. |

## Current State

### story-kernel (server side)

In `story-kernel/server/server.go`, the gRPC server is created without any TLS configuration:

```go
svr := grpc.NewServer()  // no ServerOption for TLS
```

The server listens on a raw TCP socket with no encryption.

### Story CL (client side)

In `client/x/dkg/keeper/tee_client.go`, the client conditionally selects transport credentials:

```go
if strings.HasPrefix(endpoint, "https://") || strings.HasPrefix(endpoint, "tls://") {
    creds = credentials.NewTLS(nil) // TODO: use provided CA certs
} else {
    creds = insecure.NewCredentials()
}
```

The `TODO` comment acknowledges that even the TLS path does not verify a specific CA. In practice, operators use bare `host:port` endpoints, resulting in `insecure.NewCredentials()`.

## Proposed Solution

Implement **mutual TLS (mTLS)** so that both sides authenticate each other.

### Certificate Strategy

| Artifact | Location | Purpose |
|----------|----------|---------|
| CA certificate (`ca.crt`) | Shared between CL and kernel | Root of trust for both sides |
| Server cert + key (`kernel.crt`, `kernel.key`) | `story-kernel` config dir | Authenticates kernel to CL |
| Client cert + key (`cl.crt`, `cl.key`) | Story CL config dir (`~/.story/config/dkg/`) | Authenticates CL to kernel |

### Certificate Generation

Provide a CLI helper (or document `openssl` commands) to generate a self-signed CA and issue leaf certificates:

```bash
# Generate CA
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -out ca.crt -days 3650 -subj "/CN=story-dkg-ca"

# Generate kernel server cert
openssl genrsa -out kernel.key 2048
openssl req -new -key kernel.key -out kernel.csr -subj "/CN=story-kernel"
openssl x509 -req -in kernel.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out kernel.crt -days 365

# Generate CL client cert
openssl genrsa -out cl.key 2048
openssl req -new -key cl.key -out cl.csr -subj "/CN=story-cl"
openssl x509 -req -in cl.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out cl.crt -days 365
```

### Certificate Rotation

- Certificates should have a 1-year validity by default.
- Operators rotate by generating new leaf certs from the same CA and restarting both processes.
- A future enhancement could support automatic rotation via a file watcher or SIGHUP reload.

### Configuration

#### story-kernel (`config.toml`)

```toml
[grpc]
listen_addr = "0.0.0.0:50051"

[grpc.tls]
enable = true
cert_file = "/path/to/kernel.crt"
key_file  = "/path/to/kernel.key"
ca_file   = "/path/to/ca.crt"          # for verifying client certs (mTLS)
```

#### Story CL (`dkg.toml`)

```toml
[dkg]
enable = true
kernel_endpoints = ["127.0.0.1:50051"]

[dkg.tls]
enable  = true
cert_file = "/path/to/cl.crt"
key_file  = "/path/to/cl.key"
ca_file   = "/path/to/ca.crt"          # for verifying kernel server cert
```

When `dkg.tls.enable = false` (default), the current insecure behaviour is preserved for backward compatibility.

## Implementation Plan

### Step 1: story-kernel server-side TLS

1. Add `tls` section to `story-kernel/config/config.go`.
2. In `server/server.go`, load the server certificate and CA, then create `grpc.NewServer(grpc.Creds(tlsCreds))`.
3. When mTLS is enabled, set `tls.Config.ClientAuth = tls.RequireAndVerifyClientCert`.

### Step 2: Story CL client-side TLS

1. Add `tls` fields to the DKG config struct in `client/x/dkg/types/config.go` (or equivalent).
2. In `keeper/tee_client.go`, replace the prefix-based heuristic with explicit config:
   - If `dkg.tls.enable`, load client cert/key and CA, build `credentials.NewTLS(tlsConfig)`.
   - Otherwise, fall back to `insecure.NewCredentials()`.
3. Remove the `TODO: use provided CA certs` comment.

### Step 3: CLI / documentation

1. Add a `story dkg gen-tls-certs` subcommand (optional convenience).
2. Document manual certificate generation steps in operator docs.
3. Update the kernel upgrade guide to note that both old and new kernel binaries must share the same server certificate (or the CA must sign both).

### Step 4: Testing

1. Unit test: verify TLS handshake succeeds with valid certs and fails with invalid/missing certs.
2. Integration test: run local CL + kernel with mTLS enabled and confirm DKG round completes.
3. Backward-compatibility test: ensure `enable = false` still works.

## Risks and Considerations

- **Backward compatibility**: TLS must be opt-in (`enable = false` by default) to avoid breaking existing deployments.
- **Operational complexity**: mTLS adds certificate management overhead. Clear documentation and optional CLI tooling mitigate this.
- **Kernel upgrades**: During a kernel binary upgrade, both old and new binaries need valid server certificates trusted by the same CA.
- **Performance**: TLS adds a small handshake overhead per connection, but gRPC connections are long-lived so the impact is negligible.
