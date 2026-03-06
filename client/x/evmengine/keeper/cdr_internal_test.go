package keeper

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUuidToLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		uuid     uint32
		expected [32]byte
	}{
		{
			name:     "zero uuid",
			uuid:     0,
			expected: [32]byte{},
		},
		{
			name: "uuid 1",
			uuid: 1,
			expected: func() [32]byte {
				var b [32]byte
				b[31] = 1
				return b
			}(),
		},
		{
			name: "uuid 256",
			uuid: 256,
			expected: func() [32]byte {
				var b [32]byte
				b[30] = 1
				return b
			}(),
		},
		{
			name: "max uint32",
			uuid: 0xFFFFFFFF,
			expected: func() [32]byte {
				var b [32]byte
				b[28] = 0xFF
				b[29] = 0xFF
				b[30] = 0xFF
				b[31] = 0xFF
				return b
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := uuidToLabel(tc.uuid)
			require.Equal(t, tc.expected, result)
			require.Len(t, result, 32)

			// Verify round-trip: extract uuid back from label
			extracted := binary.BigEndian.Uint32(result[28:])
			require.Equal(t, tc.uuid, extracted)
		})
	}
}

func TestUuidToLabel_Deterministic(t *testing.T) {
	t.Parallel()

	// Same uuid should always produce the same label
	uuid := uint32(12345)
	label1 := uuidToLabel(uuid)
	label2 := uuidToLabel(uuid)
	require.Equal(t, label1, label2)

	// Different uuids should produce different labels
	label3 := uuidToLabel(uuid + 1)
	require.NotEqual(t, label1, label3)
}

func TestUuidToLabel_BigEndianEncoding(t *testing.T) {
	t.Parallel()

	// Verify that the encoding matches Solidity's abi.encode(uint256(uuid))
	// which puts the value in the last 4 bytes of a 32-byte big-endian word
	uuid := uint32(0x01020304)
	label := uuidToLabel(uuid)

	// First 28 bytes should be zero
	for i := range 28 {
		require.Equal(t, byte(0), label[i], "byte %d should be zero", i)
	}

	// Last 4 bytes should be big-endian representation of uuid
	require.Equal(t, byte(0x01), label[28])
	require.Equal(t, byte(0x02), label[29])
	require.Equal(t, byte(0x03), label[30])
	require.Equal(t, byte(0x04), label[31])
}
