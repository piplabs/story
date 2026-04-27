package server

import (
	stdmath "math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCDRPartialsRequest_Validate(t *testing.T) {
	const validPubKeyHex = "048a706c17a85f0d1acf8789489568c8f154f415642164397c21c16fa6e380642d" +
		"795e6aceb661560fa5672a73ab9a30252ec4995ed096bde2b378d7991ce1150f"
	require.Len(t, validPubKeyHex, 2*requesterPubKeyByteLen)

	tests := []struct {
		name      string
		req       getCDRPartialsRequest
		expectErr string
	}{
		{
			name: "valid request",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: validPubKeyHex,
			},
		},
		{
			name: "valid request with uuid 0",
			req: getCDRPartialsRequest{
				Uuid:               0,
				RequesterPubKeyHex: validPubKeyHex,
			},
		},
		{
			name: "valid request with uuid at max allowed",
			req: getCDRPartialsRequest{
				Uuid:               maxCDRUuid,
				RequesterPubKeyHex: validPubKeyHex,
			},
		},
		{
			name: "empty pubkey hex",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: "",
			},
			expectErr: "requester_pub_key_hex is required",
		},
		{
			name: "missing pubkey hex (zero value)",
			req: getCDRPartialsRequest{
				Uuid: 42,
			},
			expectErr: "requester_pub_key_hex is required",
		},
		{
			name: "non-hex characters",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: strings.Repeat("z", 2*requesterPubKeyByteLen),
			},
			expectErr: "requester_pub_key_hex must be valid hex",
		},
		{
			name: "odd-length hex",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: validPubKeyHex[:len(validPubKeyHex)-1],
			},
			expectErr: "requester_pub_key_hex must be valid hex",
		},
		{
			name: "valid hex but wrong length (too short)",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: "02a1b2c3",
			},
			expectErr: "uncompressed secp256k1 public key",
		},
		{
			name: "valid hex but wrong length (too long)",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: validPubKeyHex + "00",
			},
			expectErr: "uncompressed secp256k1 public key",
		},
		{
			name: "uuid above max",
			req: getCDRPartialsRequest{
				Uuid:               stdmath.MaxUint32,
				RequesterPubKeyHex: validPubKeyHex,
			},
			expectErr: "uuid out of range",
		},
		{
			name: "valid request with 0x prefix",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: "0x" + validPubKeyHex,
			},
		},
		{
			name: "valid request with 0X prefix",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: "0X" + validPubKeyHex,
			},
		},
		{
			name: "0x prefix only (empty after strip)",
			req: getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: "0x",
			},
			expectErr: "requester_pub_key_hex is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.validate()
			if tt.expectErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectErr)
			}
		})
	}
}

// TestGetCDRPartialsRequest_Validate_StripsHexPrefix verifies that validate()
// normalizes RequesterPubKeyHex by stripping a leading `0x`/`0X` prefix in
// place, so the downstream keeper sees the canonical no-prefix form.
func TestGetCDRPartialsRequest_Validate_StripsHexPrefix(t *testing.T) {
	const validPubKeyHex = "048a706c17a85f0d1acf8789489568c8f154f415642164397c21c16fa6e380642d" +
		"795e6aceb661560fa5672a73ab9a30252ec4995ed096bde2b378d7991ce1150f"

	for _, prefix := range []string{"0x", "0X"} {
		t.Run(prefix, func(t *testing.T) {
			req := getCDRPartialsRequest{
				Uuid:               42,
				RequesterPubKeyHex: prefix + validPubKeyHex,
			}
			require.NoError(t, req.validate())
			require.Equal(t, validPubKeyHex, req.RequesterPubKeyHex)
		})
	}
}
