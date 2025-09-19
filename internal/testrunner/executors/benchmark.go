package executors

import (
	"fmt"
	"path/filepath"
	"time"

	"go-server/internal/testrunner/types"
)

// BenchmarkTestExecutor handles benchmark test execution
type BenchmarkTestExecutor struct{}

// NewBenchmarkTestExecutor creates a new benchmark test executor
func NewBenchmarkTestExecutor() *BenchmarkTestExecutor {
	return &BenchmarkTestExecutor{}
}

// Run executes benchmark tests
func (e *BenchmarkTestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running Benchmarks")
	fmt.Println("=======================")

	start := time.Now()

	args := []string{"test", "./test", "-bench=.", "-benchmem"}
	if config.Verbose {
		args = append(args, "-v")
	}

	output, err := runCommand("go", args...)
	duration := time.Since(start)

	logFile := filepath.Join(runDir, "benchmarks.log")
	writeLog(logFile, output)

	passed := err == nil
	if passed {
		fmt.Println("PASSED: benchmarks")
	} else {
		fmt.Printf("FAILED: benchmarks\n")
	}

	return types.TestResult{
		Name:     "benchmarks",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}
