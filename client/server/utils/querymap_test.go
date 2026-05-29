package utils

import (
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
