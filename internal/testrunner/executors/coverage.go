package executors

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// CoverageTestExecutor handles coverage test execution
type CoverageTestExecutor struct{}

// NewCoverageTestExecutor creates a new coverage test executor
func NewCoverageTestExecutor() *CoverageTestExecutor {
	return &CoverageTestExecutor{}
}

// Run executes coverage tests
func (e *CoverageTestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running Coverage Analysis")
	fmt.Println("=============================")

	start := time.Now()

	coverageFile := filepath.Join(runDir, "coverage.out")
	args := []string{"test", "./internal/...", "-coverprofile=" + coverageFile, "-covermode=atomic"}

	output, err := runCommand("go", args...)

	// Generate HTML coverage report
	if err == nil {
		htmlFile := filepath.Join(runDir, "coverage.html")
		htmlOutput, htmlErr := runCommand("go", "tool", "cover", "-html="+coverageFile, "-o", htmlFile)
		if htmlErr != nil {
			output += "\n\nHTML Coverage Generation:\n" + htmlOutput
		}
	}

	duration := time.Since(start)

	logFile := filepath.Join(runDir, "coverage.log")
	writeLog(logFile, output)

	passed := err == nil
	if passed {
		fmt.Println("PASSED: coverage")
	} else {
		fmt.Printf("FAILED: coverage\n")
	}

	return types.TestResult{
		Name:     "coverage",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}
