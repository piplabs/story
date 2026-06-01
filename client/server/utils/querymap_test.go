package utils

import (
	"encoding/base64"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildMap_Normal(t *testing.T) {
	t.Parallel()

	q := url.Values{
		"status": {"BOND_STATUS_BONDED"},
		"page":   {"1"},
	}
	m := buildMap(q)
	require.Equal(t, []string{"BOND_STATUS_BONDED"}, m["status"])
	require.Equal(t, []string{"1"}, m["page"])
}

func TestBuildMap_NestedDot(t *testing.T) {
	t.Parallel()

	q := url.Values{
		"pagination.limit":  {"10"},
		"pagination.offset": {"0"},
	}
	m := buildMap(q)
	sub, ok := m["pagination"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, []string{"10"}, sub["limit"])
	require.Equal(t, []string{"0"}, sub["offset"])
}

func TestBuildMap_ExactlyAtLimit(t *testing.T) {
	t.Parallel()

	// Build a key with depth == maxQueryDepth (should still work)
	segments := make([]string, maxQueryDepth)
	for i := range segments {
		segments[i] = "k"
	}
	deepKey := strings.Join(segments, ".")

	q := url.Values{
		deepKey: {"val"},
	}
	m := buildMap(q)
	require.NotNil(t, m)

	// Traverse to the deepest level
	current := m
	for i := 0; i < maxQueryDepth-1; i++ {
		sub, ok := current["k"].(map[string]any)
		require.True(t, ok, "level %d should be a map", i)
		current = sub
	}
}

func TestQueryMapToVal_DepthExceeded_ReturnsError(t *testing.T) {
	t.Parallel()

	// Key with depth > maxQueryDepth must be rejected up front instead of silently dropped.
	deepKey := strings.Repeat("a.", maxQueryDepth+1) + "a"
	q := url.Values{
		deepKey: {"val"},
	}

	var val map[string]any
	err := QueryMapToVal(q, &val)
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds max nesting depth")
}

func TestQueryMapToVal_MaliciousDepth_ReturnsError(t *testing.T) {
	t.Parallel()

	// A maliciously deep key is rejected up front, before buildMap recurses on it.
	key := strings.Repeat("a.", 6000) + "a"
	q := url.Values{
		key: {"x"},
	}

	var val map[string]any
	err := QueryMapToVal(q, &val)
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds max nesting depth")
}

func TestQueryMapToVal_AtLimit_NoError(t *testing.T) {
	t.Parallel()

	// Key with depth == maxQueryDepth is still accepted.
	atLimitKey := strings.Repeat("a.", maxQueryDepth) + "a"
	q := url.Values{
		atLimitKey: {"val"},
	}

	var val map[string]any
	require.NoError(t, QueryMapToVal(q, &val))
}

func BenchmarkBuildMap_Normal(b *testing.B) {
	q := url.Values{
		"pagination.limit": {"10"},
		"status":           {"BOND_STATUS_BONDED"},
	}
	for b.Loop() {
		buildMap(q)
	}
}

// TestQueryMapToVal_Base64BytesPaginationKey covers the chain-side fix for the
// pagination bug where /staking/validators page 2 returned zero entries:
// PageResponse.NextKey is amino-serialized as a base64 string, but the request
// handler previously cast pagination.key directly to []byte without decoding,
// turning a 21-byte validator iterator key into 28 ASCII bytes that pointed
// outside the prefix store. The hook decodes base64 once at the query-parse
// boundary so handlers can hand the value straight to query.PageRequest.Key.
func TestQueryMapToVal_Base64BytesPaginationKey(t *testing.T) {
	t.Parallel()

	type req struct {
		Key Base64Bytes `mapstructure:"key"`
	}

	// 21-byte validator store key: 0x14 length-prefix + 20-byte address. This
	// is the exact shape emitted by cosmos-sdk x/staking GetValidatorKey.
	raw := []byte{
		0x14,
		0x86, 0x5e, 0xeb, 0x8c, 0x35, 0x9c, 0xc1, 0x13,
		0x77, 0x11, 0xf7, 0xa2, 0xcc, 0x15, 0xd1, 0x95,
		0x18, 0x3f, 0x64, 0x93,
	}
	encoded := base64.StdEncoding.EncodeToString(raw)

	q := url.Values{}
	q.Set("key", encoded)

	var got req
	require.NoError(t, QueryMapToVal(q, &got))
	require.Equal(t, Base64Bytes(raw), got.Key, "base64 query value must decode to raw store key")
}

// TestQueryMapToVal_Base64BytesEmpty ensures the hook treats an absent or empty
// pagination.key as an empty slice (matching how cosmos-sdk paginate() handles
// nil page-request keys).
func TestQueryMapToVal_Base64BytesEmpty(t *testing.T) {
	t.Parallel()

	type req struct {
		Key Base64Bytes `mapstructure:"key"`
	}

	// Empty value.
	q := url.Values{}
	q.Set("key", "")

	var got req
	require.NoError(t, QueryMapToVal(q, &got))
	require.Empty(t, got.Key)

	// Missing field entirely.
	var got2 req
	require.NoError(t, QueryMapToVal(url.Values{}, &got2))
	require.Empty(t, got2.Key)
}

// TestQueryMapToVal_Base64BytesInvalid rejects malformed base64 instead of
// silently casting to raw ASCII bytes (which is what the bug did).
func TestQueryMapToVal_Base64BytesInvalid(t *testing.T) {
	t.Parallel()

	type req struct {
		Key Base64Bytes `mapstructure:"key"`
	}

	q := url.Values{}
	q.Set("key", "not!valid!base64!")

	var got req
	require.Error(t, QueryMapToVal(q, &got))
}
