package types

import "time"

// TestConfig represents the configuration for running tests
type TestConfig struct {
	TestType    string
	Verbose     bool
	Coverage    bool
	Benchmark   bool
	OutputDir   string
	Timeout     time.Duration
	TestRunName string
}

// TestResult represents the result of a test execution
type TestResult struct {
	Name     string
	Passed   bool
	Output   string
	LogFile  string
	Duration time.Duration
}

// TestSuite represents a collection of test results
type TestSuite struct {
	Results []TestResult
	Total   int
	Passed  int
	Failed  int
}

// TestExecutor defines the interface for running different types of tests
type TestExecutor interface {
	Run(config *TestConfig, runDir string) TestResult
}

// TestReporter defines the interface for generating test reports
type TestReporter interface {
	GenerateReport(suite *TestSuite, runDir string) error
}
