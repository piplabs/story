//go:build integration

package dkg

import (
	"os"
	"testing"
)

// TestDKG_PrioritySuite runs all test cases grouped by priority: P0 → P1 → P2.
// Supports resume mode: only runs tests that haven't passed yet.
//
// P0: Happy path + must-work (production daily paths)
// P1: Adversarial + bug regression (contract-layer attacks, known issues)
// P2: Edge/defensive + requires mock kernel (SGX prevents in production)
//
// Usage:
//
//	# Fresh run (clears history):
//	go test -tags=integration -run TestDKG_PrioritySuite/P0 -timeout 60m
//
//	# Resume mode (skip previously passed):
//	DKG_RESUME=true go test -tags=integration -run TestDKG_PrioritySuite/P0 -timeout 60m
//
//	# P0+P1:
//	DKG_RESUME=true go test -tags=integration -run "TestDKG_PrioritySuite/(P0|P1)" -timeout 120m
//
//	# Full:
//	DKG_RESUME=true go test -tags=integration -run TestDKG_PrioritySuite -timeout 180m
func TestDKG_PrioritySuite(t *testing.T) {
	resumeMode := os.Getenv("DKG_RESUME") == "true"

	if !resumeMode {
		resetPassedTests()
		t.Log("Fresh run: cleared passed tests history")
	}

	passed := loadPassedTests()
	cases := collectAllCases()

	counts := map[string]int{}
	skipped := 0
	for _, tc := range cases {
		counts[tc.Priority]++
		if resumeMode && passed[tc.ID] {
			skipped++
		}
	}

	if resumeMode {
		t.Logf("Resume mode: %d previously passed, %d remaining", skipped, len(cases)-skipped)
	}
	t.Logf("Total: %d cases (P0=%d, P1=%d, P2=%d)", len(cases), counts["P0"], counts["P1"], counts["P2"])

	priorities := []struct {
		name  string
		label string
	}{
		{"P0", "P0_HappyPath"},
		{"P1", "P1_Adversarial"},
		{"P2", "P2_Edge"},
	}

	for _, p := range priorities {
		filtered := filterByPriority(cases, p.name)
		t.Run(p.label, func(t *testing.T) {
			// In resume mode, check if entire priority group is done
			if resumeMode {
				remaining := 0
				for _, tc := range filtered {
					if !passed[tc.ID] {
						remaining++
					}
				}
				t.Logf("%s: %d cases (%d remaining)", p.name, len(filtered), remaining)
				if remaining == 0 {
					t.Skipf("all %d %s cases previously passed", len(filtered), p.name)
					return
				}
			} else {
				t.Logf("%s: %d cases", p.name, len(filtered))
			}

			for _, tc := range filtered {
				tc := tc
				t.Run(tc.ID, func(t *testing.T) {
					// Resume: skip already passed
					if resumeMode && passed[tc.ID] {
						t.Skipf("previously passed")
						return
					}

					// Skip mock-kernel-dependent cases in live mode
					if ScenarioNameForCase(tc.ID) == "" && tc.SkipIfLive != "" && os.Getenv("DKG_TEST_INVALIDATE_INDEX") == "" {
						t.Skip(tc.SkipIfLive)
						return
					}
					if tc.Run == nil {
						t.Skip("no Run")
						return
					}

					globalHarness.RunCase(t, tc)

					// Mark passed on success
					if !t.Failed() {
						markTestPassed(tc.ID)
					}
				})
			}
		})
	}
}

// collectAllCases gathers cases from all sources in a stable order.
func collectAllCases() []TestCase {
	var all []TestCase
	all = append(all, P0Cases()...)
	all = append(all, P1Cases()...)
	all = append(all, P2Cases()...)
	all = append(all, P3Cases()...)
	all = append(all, CLCases()...)
	return all
}

// filterByPriority returns cases matching the given priority.
func filterByPriority(cases []TestCase, priority string) []TestCase {
	var out []TestCase
	for _, tc := range cases {
		if tc.Priority == priority {
			out = append(out, tc)
		}
	}
	return out
}
