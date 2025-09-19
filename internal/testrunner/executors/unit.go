package executors

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// UnitTestExecutor handles unit test execution
type UnitTestExecutor struct{}

// NewUnitTestExecutor creates a new unit test executor
func NewUnitTestExecutor() *UnitTestExecutor {
	return &UnitTestExecutor{}
}

// Run executes unit tests
func (e *UnitTestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running Unit Tests")
	fmt.Println("========================")

	start := time.Now()

	args := []string{"test", "./internal/..."}
	if config.Verbose {
		args = append(args, "-v")
	}

	output, err := runCommand("go", args...)
	duration := time.Since(start)

	logFile := filepath.Join(runDir, "unit_tests.log")
	writeLog(logFile, output)

	passed := err == nil
	if passed {
		fmt.Println("PASSED: unit_tests")
	} else {
		fmt.Printf("FAILED: unit_tests\n")
	}

	return types.TestResult{
		Name:     "unit_tests",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}
