package executors

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// IntegrationTestExecutor handles integration test execution
type IntegrationTestExecutor struct{}

// NewIntegrationTestExecutor creates a new integration test executor
func NewIntegrationTestExecutor() *IntegrationTestExecutor {
	return &IntegrationTestExecutor{}
}

// Run executes integration tests
func (e *IntegrationTestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running Integration Tests")
	fmt.Println("==============================")

	start := time.Now()

	args := []string{"test", "./test", "-run", "TestServer"}
	if config.Verbose {
		args = append(args, "-v")
	}

	output, err := runCommand("go", args...)
	duration := time.Since(start)

	logFile := filepath.Join(runDir, "integration_tests.log")
	writeLog(logFile, output)

	passed := err == nil
	if passed {
		fmt.Println("PASSED: integration_tests")
	} else {
		fmt.Printf("FAILED: integration_tests\n")
	}

	return types.TestResult{
		Name:     "integration_tests",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}
