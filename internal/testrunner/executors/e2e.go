package executors

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// E2ETestExecutor handles end-to-end test execution
type E2ETestExecutor struct{}

// NewE2ETestExecutor creates a new E2E test executor
func NewE2ETestExecutor() *E2ETestExecutor {
	return &E2ETestExecutor{}
}

// Run executes end-to-end tests
func (e *E2ETestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running End-to-End Tests")
	fmt.Println("=============================")

	start := time.Now()

	args := []string{"test", "./test", "-run", "TestHealthEndpoint|TestAPIEndpoint|TestVersionEndpoint|TestMetricsEndpoint|TestCORSEndpoint|TestRequestSizeLimit"}
	if config.Verbose {
		args = append(args, "-v")
	}

	output, err := runCommand("go", args...)
	duration := time.Since(start)

	logFile := filepath.Join(runDir, "e2e_tests.log")
	writeLog(logFile, output)

	passed := err == nil
	if passed {
		fmt.Println("PASSED: e2e_tests")
	} else {
		fmt.Printf("FAILED: e2e_tests\n")
	}

	return types.TestResult{
		Name:     "e2e_tests",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}
