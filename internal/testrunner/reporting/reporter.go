package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner/types"
)

// MarkdownReporter generates markdown test reports
type MarkdownReporter struct{}

// NewMarkdownReporter creates a new markdown reporter
func NewMarkdownReporter() *MarkdownReporter {
	return &MarkdownReporter{}
}

// GenerateReport generates a markdown test report
func (r *MarkdownReporter) GenerateReport(suite *types.TestSuite, runDir string) error {
	reportFile := filepath.Join(runDir, "test_report.md")

	content := fmt.Sprintf(`# Test Report

Generated: %s

## Summary

- **Total Tests**: %d
- **Passed**: %d
- **Failed**: %d
- **Success Rate**: %.0f%%

## Test Results

`, time.Now().Format("2006-01-02 15:04:05"), suite.Total, suite.Passed, suite.Failed, float64(suite.Passed)/float64(suite.Total)*100)

	for _, result := range suite.Results {
		status := "❌ FAILED"
		if result.Passed {
			status = "✅ PASSED"
		}
		content += fmt.Sprintf("- **%s**: %s (%.2fs)\n", result.Name, status, result.Duration.Seconds())
	}

	content += "\n## Log Files\n\n"
	for _, result := range suite.Results {
		content += fmt.Sprintf("- [%s](%s)\n", result.Name, filepath.Base(result.LogFile))
	}

	return os.WriteFile(reportFile, []byte(content), 0644)
}

// ConsoleReporter generates console output
type ConsoleReporter struct{}

// NewConsoleReporter creates a new console reporter
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// GenerateReport generates console output
func (r *ConsoleReporter) GenerateReport(suite *types.TestSuite, runDir string) error {
	fmt.Println("\nGenerating Test Report")
	fmt.Println("============================")

	fmt.Printf("\nFinal Test Summary\n")
	fmt.Printf("====================\n")
	fmt.Printf("Total Tests: %d\n", suite.Total)
	fmt.Printf("Passed: %d\n", suite.Passed)
	fmt.Printf("Failed: %d\n", suite.Failed)
	fmt.Printf("Success Rate: %.0f%%\n", float64(suite.Passed)/float64(suite.Total)*100)
	fmt.Printf("\nTest results saved to: %s\n", runDir)

	return nil
}
