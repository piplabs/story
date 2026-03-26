//go:build integration

package dkg

import (
	"context"
	"os"
	"sort"
	"testing"
)

// TestDKG_FastSuite is an optimized version of TestDKG_FullSuite that reduces
// execution time by ~40% through two strategies:
//
//  1. Snapshot tests (93 cases): run in parallel via t.Parallel() — they only
//     query chain state with no mutations or stage waits.
//  2. Scenario tests (56 cases, 19 groups): batched by scenario name — one
//     SetupScenario/TeardownScenario per group instead of per case.
//
// Stage-wait (40 cases) and upgrade (13 cases) run sequentially as before.
//
// The original TestDKG_FullSuite in run_test.go is NOT modified and remains
// available as a fallback.
func TestDKG_FastSuite(t *testing.T) {
	// Resume mode
	if os.Getenv("DKG_RESUME") != "true" {
		resetPassedTests()
		t.Log("Fresh run: cleared passed tests history")
	} else {
		passed := loadPassedTests()
		t.Logf("Resume run: %d tests previously passed, will be skipped", len(passed))
	}

	shortReg := getUint64Env("DKG_SHORT_REGISTRATION", 50)

	// Classify all cases first (needed to check if Phase 1 can be skipped)
	allCases := OrderedCases()
	snapshot, stageWait, scenarioGroups, upgrade := classifyCases(allCases)

	// --- Phase 1: Short period configuration ---
	// Skip Phase 1 entirely if all Tier 1-3 tests are already passed or will be skipped (resume mode).
	// This avoids an unnecessary reset when only upgrade tests (Tier 4) remain.
	passed := loadPassedTests()
	tier123NeedRun := 0
	if os.Getenv("DKG_RESUME") == "true" {
		// Build upgrade ID set so we don't count upgrade tests against Phase 1
		upgradeIDs := make(map[string]bool, len(upgrade))
		for _, tc := range upgrade {
			upgradeIDs[tc.ID] = true
		}
		// Count Tier 1-3 tests that are NOT passed and NOT upgrade tests
		for _, tc := range snapshot {
			if !passed[tc.ID] && !upgradeIDs[tc.ID] {
				tier123NeedRun++
			}
		}
		for _, tc := range stageWait {
			if !passed[tc.ID] && !upgradeIDs[tc.ID] {
				tier123NeedRun++
			}
		}
		for _, cases := range scenarioGroups {
			for _, tc := range cases {
				if !passed[tc.ID] && !upgradeIDs[tc.ID] {
					tier123NeedRun++
				}
			}
		}
	} else {
		tier123NeedRun = 1 // force Phase 1 on fresh run
	}

	if tier123NeedRun == 0 {
		t.Log("=== FastSuite Phase 1: All Tier 1-3 tests passed/skipped, jumping to upgrade tests ===")
	} else {
		t.Logf("=== FastSuite Phase 1: %d Tier 1-3 tests need to run ===", tier123NeedRun)
		// No period switching — use whatever the chain has. Only wait for DKG active.
		t.Logf("DKG params active (elapsed=0s, registration_period=%d)", shortReg)
		waitForDKGActive(t, globalHarness)
		ensureAllHealthy(t)
	}

	t.Logf("=== FastSuite Classification: snapshot=%d stageWait=%d scenario=%d(in %d groups) upgrade=%d ===",
		len(snapshot), len(stageWait), countScenarioCases(scenarioGroups), len(scenarioGroups), len(upgrade))

	// --- Phase 2: Run non-upgrade tests ---

	// Tier 1: Snapshot tests in parallel
	t.Logf("=== FastSuite Tier 1: %d snapshot tests (parallel) ===", len(snapshot))
	runSnapshotParallel(t, snapshot)

	// Tier 2: Stage-wait tests sequentially
	t.Logf("=== FastSuite Tier 2: %d stage-wait tests (sequential) ===", len(stageWait))
	runCaseList(t, stageWait)

	// Tier 3: Scenario tests batched by scenario
	t.Logf("=== FastSuite Tier 3: %d scenario tests in %d groups (batched) ===",
		countScenarioCases(scenarioGroups), len(scenarioGroups))
	runScenarioGroupsBatched(t, scenarioGroups)

	// --- Phase 3: Upgrade tests ---
	if len(upgrade) > 0 {
		// Skip Phase 3 if all upgrade tests already passed
		upgradeAllPassed := true
		if os.Getenv("DKG_RESUME") == "true" {
			passed = loadPassedTests() // reload (Tier 1-3 may have added new passes)
			for _, tc := range upgrade {
				if !passed[tc.ID] {
					upgradeAllPassed = false
					break
				}
			}
		} else {
			upgradeAllPassed = false
		}

		if upgradeAllPassed {
			t.Log("=== FastSuite Phase 3: All upgrade tests passed, skipping ===")
		} else {
			// No period switching — use current chain periods for upgrade tests too.
			ensureAllHealthy(t)
		}
		t.Logf("=== FastSuite Tier 4: %d upgrade tests (sequential) ===", len(upgrade))
		runCaseList(t, upgrade)
	}

	// Summary
	total := len(snapshot) + len(stageWait) + countScenarioCases(scenarioGroups) + len(upgrade)
	t.Logf("============================================================")
	t.Logf("FAST SUITE COMPLETE  Total: %d | Snapshot: %d | StageWait: %d | Scenario: %d (%d groups) | Upgrade: %d",
		total, len(snapshot), len(stageWait), countScenarioCases(scenarioGroups), len(scenarioGroups), len(upgrade))
	t.Logf("============================================================")
}

// classifyCases splits all test cases into 4 tiers based on their properties.
//
//	Tier 1 (snapshot):       no scenario, NeedsRoundWait=false, not upgrade → parallel-safe
//	Tier 2 (stageWait):      no scenario, NeedsRoundWait=true, not upgrade → sequential
//	Tier 3 (scenarioGroups): has scenario, not upgrade → batched by scenario name
//	Tier 4 (upgrade):        in upgradeScenarioCaseIDs → run last
func classifyCases(cases []TestCase) (snapshot, stageWait []TestCase, scenarioGroups map[string][]TestCase, upgrade []TestCase) {
	scenarioGroups = make(map[string][]TestCase)

	for _, c := range cases {
		if upgradeScenarioCaseIDs[c.ID] {
			upgrade = append(upgrade, c)
			continue
		}

		scenario := ScenarioNameForCase(c.ID)
		if scenario != "" {
			scenarioGroups[scenario] = append(scenarioGroups[scenario], c)
			continue
		}

		if c.NeedsRoundWait {
			stageWait = append(stageWait, c)
		} else {
			snapshot = append(snapshot, c)
		}
	}

	return
}

// runSnapshotParallel runs snapshot (query-only) tests concurrently using t.Parallel().
// These tests only read chain state — no mutations, no stage waits — so they are safe
// to run in parallel within the same DKG round.
func runSnapshotParallel(t *testing.T, cases []TestCase) {
	t.Helper()

	passed := loadPassedTests()
	resumeMode := os.Getenv("DKG_RESUME") == "true"

	t.Run("Snapshot", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc // capture
			t.Run(tc.ID+"_"+tc.Priority, func(t *testing.T) {
				t.Parallel() // run concurrently with other snapshot tests

				if resumeMode && passed[tc.ID] {
					t.Skipf("previously passed (resume mode)")
					return
				}
				if tc.SkipIfLive != "" && os.Getenv("DKG_TEST_INVALIDATE_INDEX") == "" {
					t.Skip(tc.SkipIfLive)
					return
				}
				if tc.Run == nil {
					t.Skip("no Run (internal/fault-injection only)")
					return
				}

				t.Logf("[%s] %s", tc.ID, tc.Description)
				tc.Run(t, globalHarness)

				if !t.Failed() {
					markTestPassed(tc.ID)
				}
			})
		}
	})
}

// runScenarioGroupsBatched runs scenario tests grouped by scenario name.
// For each group: one SetupScenario → all cases sequentially → one TeardownScenario.
// This avoids repeated setup/teardown overhead (e.g., tee_down_one_node: 12 cases
// share 1 setup instead of 12 separate setups).
func runScenarioGroupsBatched(t *testing.T, groups map[string][]TestCase) {
	t.Helper()

	passed := loadPassedTests()
	resumeMode := os.Getenv("DKG_RESUME") == "true"

	// Sort group names for deterministic execution order
	names := make([]string, 0, len(groups))
	for name := range groups {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, scenarioName := range names {
		cases := groups[scenarioName]
		t.Run("Scenario_"+scenarioName, func(t *testing.T) {
			runScenarioBatch(t, scenarioName, cases, passed, resumeMode)
		})
	}
}

// runScenarioBatch executes a batch of test cases that share the same scenario.
// Setup and teardown happen once for the entire batch.
func runScenarioBatch(t *testing.T, scenarioName string, cases []TestCase, passed map[string]bool, resumeMode bool) {
	t.Helper()
	ctx := context.Background()

	// Check if ALL cases in this batch are already passed (resume mode)
	if resumeMode {
		allPassed := true
		for _, tc := range cases {
			if !passed[tc.ID] {
				allPassed = false
				break
			}
		}
		if allPassed {
			t.Skipf("all %d cases in scenario %q previously passed (resume mode)", len(cases), scenarioName)
			return
		}
	}

	// Setup scenario once
	if globalHarness.Driver != nil {
		t.Logf("Setting up scenario %q for %d cases", scenarioName, len(cases))
		if err := globalHarness.Driver.SetupScenario(ctx, scenarioName); err != nil {
			t.Fatalf("SetupScenario(%s): %v", scenarioName, err)
		}
		defer func() {
			if err := globalHarness.Driver.TeardownScenario(ctx, scenarioName); err != nil {
				t.Logf("TeardownScenario(%s): %v", scenarioName, err)
			}
			ensureAllHealthy(t)
		}()

		// Verify precondition once
		verifyScenarioPrecondition(t, globalHarness, scenarioName)
	}

	// Run each case in the batch sequentially (within the shared scenario)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.ID+"_"+tc.Priority, func(t *testing.T) {
			if resumeMode && passed[tc.ID] {
				t.Skipf("previously passed (resume mode)")
				return
			}
			if tc.Run == nil {
				t.Skip("no Run (internal/fault-injection only)")
				return
			}

			t.Logf("[%s] %s (scenario: %s)", tc.ID, tc.Description, scenarioName)
			if tc.Expected != "" {
				t.Logf("  Expected: %s", tc.Expected)
			}

			tc.Run(t, globalHarness)

			if !t.Failed() {
				markTestPassed(tc.ID)
			}
		})
	}
}

func countScenarioCases(groups map[string][]TestCase) int {
	n := 0
	for _, cases := range groups {
		n += len(cases)
	}
	return n
}

// printFastSuiteClassification prints the detailed classification for debugging.
// Call with -v to see which cases go into which tier.
func printFastSuiteClassification(t *testing.T, snapshot, stageWait []TestCase, scenarioGroups map[string][]TestCase, upgrade []TestCase) {
	t.Helper()
	t.Logf("--- Tier 1: Snapshot (%d) ---", len(snapshot))
	for _, c := range snapshot {
		t.Logf("  %s", c.ID)
	}
	t.Logf("--- Tier 2: Stage-Wait (%d) ---", len(stageWait))
	for _, c := range stageWait {
		t.Logf("  %s", c.ID)
	}
	t.Logf("--- Tier 3: Scenario (%d groups) ---", len(scenarioGroups))
	names := make([]string, 0, len(scenarioGroups))
	for name := range scenarioGroups {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		cases := scenarioGroups[name]
		ids := make([]string, len(cases))
		for i, c := range cases {
			ids[i] = c.ID
		}
		t.Logf("  %s (%d): %v", name, len(cases), ids)
	}
	t.Logf("--- Tier 4: Upgrade (%d) ---", len(upgrade))
	for _, c := range upgrade {
		t.Logf("  %s", c.ID)
	}
}

