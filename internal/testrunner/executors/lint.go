package executors

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// LintTestExecutor handles linting test execution
type LintTestExecutor struct{}

// NewLintTestExecutor creates a new lint test executor
func NewLintTestExecutor() *LintTestExecutor {
	return &LintTestExecutor{}
}

// Run executes linting tests
func (e *LintTestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running Linting")
	fmt.Println("===================")

	start := time.Now()

	// Run go vet
	vetOutput, vetErr := runCommand("go", "vet", "./...")

	// Run go fmt check
	fmtOutput, fmtErr := runCommand("go", "fmt", "./...")

	output := "Go Vet:\n" + vetOutput + "\n\nGo Fmt:\n" + fmtOutput
	passed := vetErr == nil && fmtErr == nil

	duration := time.Since(start)

	logFile := filepath.Join(runDir, "linting.log")
	writeLog(logFile, output)

	if passed {
		fmt.Println("PASSED: linting")
	} else {
		fmt.Printf("FAILED: linting\n")
	}

	return types.TestResult{
		Name:     "linting",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}
