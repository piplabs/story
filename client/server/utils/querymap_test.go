package utils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestQueryMapToVal_NumericTypes verifies that single-value query parameters
// decode into all signed and unsigned integer widths. This previously regressed
// for uint32 (and other widths beyond uint64/int) because the
// stringArrayToNative decode hook only handled a subset of integer kinds, so
// any request struct field typed `uint32` (e.g. dkg `round`) would error with
// "unconvertible type '[]string'".
func TestQueryMapToVal_NumericTypes(t *testing.T) {
	t.Parallel()

	type payload struct {
		U8  uint8  `mapstructure:"u8"`
		U16 uint16 `mapstructure:"u16"`
		U32 uint32 `mapstructure:"u32"`
		U64 uint64 `mapstructure:"u64"`
		I8  int8   `mapstructure:"i8"`
		I16 int16  `mapstructure:"i16"`
		I32 int32  `mapstructure:"i32"`
		I64 int64  `mapstructure:"i64"`
		B   bool   `mapstructure:"b"`
		S   string `mapstructure:"s"`
	}

	q := url.Values{}
	q.Set("u8", "200")
	q.Set("u16", "60000")
	q.Set("u32", "4000000000")
	q.Set("u64", "18000000000000000000")
	q.Set("i8", "-100")
	q.Set("i16", "-30000")
	q.Set("i32", "-2000000000")
	q.Set("i64", "-9000000000000000000")
	q.Set("b", "true")
	q.Set("s", "hello")

	var got payload
	require.NoError(t, QueryMapToVal(q, &got))
	require.Equal(t, payload{
		U8:  200,
		U16: 60000,
		U32: 4_000_000_000,
		U64: 18_000_000_000_000_000_000,
		I8:  -100,
		I16: -30000,
		I32: -2_000_000_000,
		I64: -9_000_000_000_000_000_000,
		B:   true,
		S:   "hello",
	}, got)
}

// TestQueryMapToVal_Uint32Round mirrors the dkg request-struct shape that
// originally exposed the missing-Uint32 bug.
func TestQueryMapToVal_Uint32Round(t *testing.T) {
	t.Parallel()

	type req struct {
		Round uint32 `mapstructure:"round"`
	}

	q := url.Values{}
	q.Set("round", "6")

	var got req
	require.NoError(t, QueryMapToVal(q, &got))
	require.Equal(t, uint32(6), got.Round)
}
