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

func TestBuildMap_DepthLimit(t *testing.T) {
	t.Parallel()

	// Build a key with depth > maxQueryDepth
	segments := make([]string, maxQueryDepth+5)
	for i := range segments {
		segments[i] = "a"
	}
	deepKey := strings.Join(segments, ".")

	q := url.Values{
		deepKey: {"val"},
	}
	m := buildMap(q)
	// Should not panic and should return a map (nested levels beyond limit are nil)
	require.NotNil(t, m)
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

func TestBuildMap_MaliciousDepth_NoOOM(t *testing.T) {
	t.Parallel()

	// Simulate the attack: 6000 dot-separated segments
	key := strings.Repeat("a.", 6000) + "a"
	q := url.Values{
		key: {"x"},
	}

	// This should complete quickly without OOM
	m := buildMap(q)
	require.NotNil(t, m)
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

func BenchmarkBuildMap_DeepKey_6000(b *testing.B) {
	key := strings.Repeat("a.", 6000) + "a"
	q := url.Values{
		key: {"x"},
	}
	for b.Loop() {
		buildMap(q)
	}
}
