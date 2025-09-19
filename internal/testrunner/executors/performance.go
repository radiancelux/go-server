package executors

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// PerformanceTestExecutor handles performance test execution
type PerformanceTestExecutor struct{}

// NewPerformanceTestExecutor creates a new performance test executor
func NewPerformanceTestExecutor() *PerformanceTestExecutor {
	return &PerformanceTestExecutor{}
}

// Run executes performance tests
func (e *PerformanceTestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running Performance Tests")
	fmt.Println("===============================")

	start := time.Now()

	args := []string{"test", "./test", "-run", "TestLoadTest|TestMemoryUsage"}
	if config.Verbose {
		args = append(args, "-v")
	}

	output, err := runCommand("go", args...)
	duration := time.Since(start)

	logFile := filepath.Join(runDir, "performance_tests.log")
	writeLog(logFile, output)

	passed := err == nil
	if passed {
		fmt.Println("PASSED: performance_tests")
	} else {
		fmt.Printf("FAILED: performance_tests\n")
	}

	return types.TestResult{
		Name:     "performance_tests",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}
