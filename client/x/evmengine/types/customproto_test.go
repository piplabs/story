package types_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/types"
)

func TestAddress_MarshalUnmarshal(t *testing.T) {
	original := types.Address(common.HexToAddress("0x1234567890123456789012345678901234567890"))

	// Marshal & Unmarshal binary
	bz, err := original.Marshal()
	require.NoError(t, err)
	require.Len(t, bz, common.AddressLength)

	var unmarshaled types.Address
	err = unmarshaled.Unmarshal(bz)
	require.NoError(t, err)
	require.Equal(t, original, unmarshaled)

	// MarshalTo
	buf := make([]byte, common.AddressLength)
	n, err := original.MarshalTo(buf)
	require.NoError(t, err)
	require.Equal(t, common.AddressLength, n)
	require.Equal(t, bz, buf)

	// Invalid binary length
	err = unmarshaled.Unmarshal([]byte{0x01, 0x02})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid address length")

	// Marshal & Unmarshal JSON
	jsonBz, err := json.Marshal(original)
	require.NoError(t, err)

	var jsonUnmarshaled types.Address
	err = json.Unmarshal(jsonBz, &jsonUnmarshaled)
	require.NoError(t, err)
	require.Equal(t, original, jsonUnmarshaled)

	// Invalid JSON format (not a string)
	invalidJSON := []byte(`12345`)
	err = jsonUnmarshaled.UnmarshalJSON(invalidJSON)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshal address")

	// Invalid hex string (too short)
	err = jsonUnmarshaled.UnmarshalJSON([]byte(`"0x123"`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshal address")
}

func TestHash_MarshalUnmarshal(t *testing.T) {
	original := types.Hash(common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef"))

	// Marshal & Unmarshal binary
	bz, err := original.Marshal()
	require.NoError(t, err)
	require.Len(t, bz, common.HashLength)

	var unmarshaled types.Hash
	err = unmarshaled.Unmarshal(bz)
	require.NoError(t, err)
	require.Equal(t, original, unmarshaled)

	// MarshalTo
	buf := make([]byte, common.HashLength)
	n, err := original.MarshalTo(buf)
	require.NoError(t, err)
	require.Equal(t, common.HashLength, n)
	require.Equal(t, bz, buf)

	// Invalid binary length
	err = unmarshaled.Unmarshal([]byte{0x01, 0x02})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid hash length")

	// Marshal & Unmarshal JSON
	jsonBz, err := json.Marshal(original)
	require.NoError(t, err)

	var jsonUnmarshaled types.Hash
	err = json.Unmarshal(jsonBz, &jsonUnmarshaled)
	require.NoError(t, err)
	require.Equal(t, original, jsonUnmarshaled)

	// Invalid JSON format (not a string)
	invalidJSON := []byte(`true`)
	err = jsonUnmarshaled.UnmarshalJSON(invalidJSON)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshal hash")

	// Invalid hex string (wrong length)
	err = jsonUnmarshaled.UnmarshalJSON([]byte(`"0xdead"`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshal hash")
}

func TestAddress_Size(t *testing.T) {
	var a types.Address
	require.Equal(t, common.AddressLength, a.Size())
}

func TestHash_Size(t *testing.T) {
	var h types.Hash
	require.Equal(t, common.HashLength, h.Size())
}
