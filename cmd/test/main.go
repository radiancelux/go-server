package main

import (
	"flag"
	"log"
	"time"

	"github.com/radiancelux/go-server/internal/testrunner"
	"github.com/radiancelux/go-server/internal/testrunner/types"
)

func main() {
	config := parseFlags()

	runner := testrunner.NewTestRunner()

	if err := runner.Run(config); err != nil {
		log.Fatalf("Test execution failed: %v", err)
	}
}

func parseFlags() *types.TestConfig {
	config := &types.TestConfig{}

	flag.StringVar(&config.TestType, "type", "all", "Test type: unit, integration, e2e, performance, benchmark, coverage, lint, postman, all")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&config.Coverage, "coverage", false, "Generate coverage report")
	flag.BoolVar(&config.Benchmark, "bench", false, "Run benchmarks")
	flag.StringVar(&config.OutputDir, "output", "test-results", "Output directory")
	flag.DurationVar(&config.Timeout, "timeout", 5*time.Minute, "Test timeout")

	flag.Parse()

	// Generate test run name
	config.TestRunName = time.Now().Format("20060102_150405")

	return config
}
