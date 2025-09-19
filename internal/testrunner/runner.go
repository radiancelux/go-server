package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-server/internal/testrunner/executors"
	"go-server/internal/testrunner/reporting"
	"go-server/internal/testrunner/types"
)

// TestRunner orchestrates test execution
type TestRunner struct {
	executors map[string]types.TestExecutor
	reporters []types.TestReporter
}

// NewTestRunner creates a new test runner
func NewTestRunner() *TestRunner {
	runner := &TestRunner{
		executors: make(map[string]types.TestExecutor),
		reporters: []types.TestReporter{
			reporting.NewConsoleReporter(),
			reporting.NewMarkdownReporter(),
		},
	}

	// Register executors
	runner.executors["unit"] = executors.NewUnitTestExecutor()
	runner.executors["integration"] = executors.NewIntegrationTestExecutor()
	runner.executors["e2e"] = executors.NewE2ETestExecutor()
	runner.executors["performance"] = executors.NewPerformanceTestExecutor()
	runner.executors["benchmark"] = executors.NewBenchmarkTestExecutor()
	runner.executors["coverage"] = executors.NewCoverageTestExecutor()
	runner.executors["lint"] = executors.NewLintTestExecutor()
	runner.executors["postman"] = executors.NewPostmanTestExecutor()

	return runner
}

// Run executes tests based on the configuration
func (r *TestRunner) Run(config *types.TestConfig) error {
	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	runDir := filepath.Join(config.OutputDir, fmt.Sprintf("test_run_%s", timestamp))
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return fmt.Errorf("failed to create run directory: %v", err)
	}

	fmt.Printf("Starting Go Server Test Suite\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("Test run: %s\n", timestamp)
	fmt.Printf("Results directory: %s\n\n", runDir)

	var results []types.TestResult

	switch config.TestType {
	case "all":
		results = r.runAllTests(runDir, config)
	default:
		if executor, exists := r.executors[config.TestType]; exists {
			results = append(results, executor.Run(config, runDir))
		} else {
			return fmt.Errorf("unknown test type: %s", config.TestType)
		}
	}

	// Generate summary
	suite := &types.TestSuite{
		Results: results,
		Total:   len(results),
		Passed:  0,
		Failed:  0,
	}

	for _, result := range results {
		if result.Passed {
			suite.Passed++
		} else {
			suite.Failed++
		}
	}

	// Generate reports
	for _, reporter := range r.reporters {
		if err := reporter.GenerateReport(suite, runDir); err != nil {
			fmt.Printf("Warning: Failed to generate report: %v\n", err)
		}
	}

	return nil
}

// runAllTests executes all test types
func (r *TestRunner) runAllTests(runDir string, config *types.TestConfig) []types.TestResult {
	var results []types.TestResult

	// Run all test types
	testTypes := []string{"unit", "integration", "e2e", "performance", "coverage", "lint", "postman"}

	for _, testType := range testTypes {
		if executor, exists := r.executors[testType]; exists {
			results = append(results, executor.Run(config, runDir))
		}
	}

	return results
}
