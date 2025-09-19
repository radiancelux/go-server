package executors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// PostmanTestExecutor handles Postman collection test execution
type PostmanTestExecutor struct{}

// NewPostmanTestExecutor creates a new Postman test executor
func NewPostmanTestExecutor() *PostmanTestExecutor {
	return &PostmanTestExecutor{}
}

// Run executes Postman collection tests
func (e *PostmanTestExecutor) Run(config *types.TestConfig, runDir string) types.TestResult {
	fmt.Println("Running Postman Collection Tests")
	fmt.Println("=====================================")

	start := time.Now()

	// Check if Newman is installed
	if !isNewmanInstalled() {
		output := "Newman is not installed. Please install it with: npm install -g newman"
		logFile := filepath.Join(runDir, "postman_tests.log")
		writeLog(logFile, output)

		return types.TestResult{
			Name:     "postman_tests",
			Passed:   false,
			Output:   output,
			LogFile:  logFile,
			Duration: time.Since(start),
		}
	}

	// Set up Newman command
	collectionFile := "postman/Go-Server-API.postman_collection.json"
	environmentFile := "postman/Go-Server-Environment.postman_environment.json"

	args := []string{
		"run", collectionFile,
		"--reporters", "cli,json,html",
		"--reporter-json-export", filepath.Join(runDir, "postman_report.json"),
		"--reporter-html-export", filepath.Join(runDir, "postman_report.html"),
		"--timeout-request", "10000",
		"--timeout-script", "5000",
		"--delay-request", "100",
	}

	// Add environment file if it exists
	if _, err := os.Stat(environmentFile); err == nil {
		args = append(args, "--environment", environmentFile)
	}

	// Add verbose flag if requested
	if config.Verbose {
		args = append(args, "--verbose")
	}

	output, err := runCommand("newman", args...)
	duration := time.Since(start)

	logFile := filepath.Join(runDir, "postman_tests.log")
	writeLog(logFile, output)

	passed := err == nil
	if passed {
		fmt.Println("PASSED: postman_tests")
	} else {
		fmt.Printf("FAILED: postman_tests\n")
	}

	return types.TestResult{
		Name:     "postman_tests",
		Passed:   passed,
		Output:   output,
		LogFile:  logFile,
		Duration: duration,
	}
}

// isNewmanInstalled checks if Newman is installed
func isNewmanInstalled() bool {
	_, err := exec.LookPath("newman")
	return err == nil
}
