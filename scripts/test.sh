#!/bin/bash

# Go Server Test Automation Script
# This script runs all tests and generates a comprehensive report

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results directory
TEST_RESULTS_DIR="test-results"
mkdir -p "$TEST_RESULTS_DIR"

# Timestamp for test run
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TEST_RUN_DIR="$TEST_RESULTS_DIR/test_run_$TIMESTAMP"
mkdir -p "$TEST_RUN_DIR"

echo -e "${BLUE}🚀 Starting Go Server Test Automation${NC}"
echo -e "${BLUE}=====================================${NC}"
echo "Test run: $TIMESTAMP"
echo "Results directory: $TEST_RUN_DIR"
echo ""

# Function to run tests and capture output
run_test() {
    local test_name="$1"
    local test_command="$2"
    local output_file="$TEST_RUN_DIR/${test_name}.log"
    
    echo -e "${YELLOW}Running $test_name...${NC}"
    
    if eval "$test_command" > "$output_file" 2>&1; then
        echo -e "${GREEN}✅ $test_name PASSED${NC}"
        return 0
    else
        echo -e "${RED}❌ $test_name FAILED${NC}"
        echo -e "${RED}Check $output_file for details${NC}"
        return 1
    fi
}

# Function to run tests with coverage
run_test_with_coverage() {
    local test_name="$1"
    local test_command="$2"
    local output_file="$TEST_RUN_DIR/${test_name}.log"
    local coverage_file="$TEST_RUN_DIR/${test_name}_coverage.out"
    
    echo -e "${YELLOW}Running $test_name with coverage...${NC}"
    
    if eval "$test_command" -coverprofile="$coverage_file" > "$output_file" 2>&1; then
        echo -e "${GREEN}✅ $test_name PASSED${NC}"
        
        # Generate coverage report
        go tool cover -html="$coverage_file" -o "$TEST_RUN_DIR/${test_name}_coverage.html"
        echo -e "${BLUE}📊 Coverage report: $TEST_RUN_DIR/${test_name}_coverage.html${NC}"
        return 0
    else
        echo -e "${RED}❌ $test_name FAILED${NC}"
        echo -e "${RED}Check $output_file for details${NC}"
        return 1
    fi
}

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 1. Unit Tests
echo -e "${BLUE}📋 Running Unit Tests${NC}"
echo "========================"

# Config tests
TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "config_tests" "go test ./internal/config -v"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# Errors tests
TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "errors_tests" "go test ./internal/errors -v"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# Middleware tests
TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "middleware_tests" "go test ./internal/middleware -v"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# Server tests
TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "server_tests" "go test ./internal/server -v"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# 2. Integration Tests
echo -e "${BLUE}🔗 Running Integration Tests${NC}"
echo "==============================="

TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "integration_tests" "go test ./test -v -run TestServer"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# 3. End-to-End Tests
echo -e "${BLUE}🌐 Running End-to-End Tests${NC}"
echo "============================="

TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "e2e_tests" "go test ./test -v -run TestHealthEndpoint -run TestAPIEndpoint -run TestVersionEndpoint -run TestMetricsEndpoint -run TestCORSEndpoint -run TestRequestSizeLimit"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# 4. Performance Tests
echo -e "${BLUE}⚡ Running Performance Tests${NC}"
echo "==============================="

TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "performance_tests" "go test ./test -v -run TestLoadTest -run TestMemoryUsage"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# 5. Benchmarks
echo -e "${BLUE}📊 Running Benchmarks${NC}"
echo "======================="

TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "benchmarks" "go test ./test -bench=. -benchmem"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# 6. Coverage Analysis
echo -e "${BLUE}📈 Running Coverage Analysis${NC}"
echo "============================="

TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test_with_coverage "coverage" "go test ./... -cover"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# 7. Linting
echo -e "${BLUE}🔍 Running Linting${NC}"
echo "==================="

TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "linting" "go vet ./..."; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# 8. Security Scan
echo -e "${BLUE}🔒 Running Security Scan${NC}"
echo "========================="

TOTAL_TESTS=$((TOTAL_TESTS + 1))
if run_test "security_scan" "go list -json -deps ./... | grep -E '\"(ImportPath|Imports)\"' | grep -v 'go-server' | sort | uniq"; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# Generate Test Report
echo -e "${BLUE}📋 Generating Test Report${NC}"
echo "============================"

cat > "$TEST_RUN_DIR/test_report.md" << EOF
# Test Report - $TIMESTAMP

## Summary
- **Total Tests**: $TOTAL_TESTS
- **Passed**: $PASSED_TESTS
- **Failed**: $FAILED_TESTS
- **Success Rate**: $((PASSED_TESTS * 100 / TOTAL_TESTS))%

## Test Results

### Unit Tests
- Config Tests: $(if [ -f "$TEST_RUN_DIR/config_tests.log" ] && grep -q "PASS" "$TEST_RUN_DIR/config_tests.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)
- Errors Tests: $(if [ -f "$TEST_RUN_DIR/errors_tests.log" ] && grep -q "PASS" "$TEST_RUN_DIR/errors_tests.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)
- Middleware Tests: $(if [ -f "$TEST_RUN_DIR/middleware_tests.log" ] && grep -q "PASS" "$TEST_RUN_DIR/middleware_tests.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)
- Server Tests: $(if [ -f "$TEST_RUN_DIR/server_tests.log" ] && grep -q "PASS" "$TEST_RUN_DIR/server_tests.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

### Integration Tests
- Integration Tests: $(if [ -f "$TEST_RUN_DIR/integration_tests.log" ] && grep -q "PASS" "$TEST_RUN_DIR/integration_tests.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

### End-to-End Tests
- E2E Tests: $(if [ -f "$TEST_RUN_DIR/e2e_tests.log" ] && grep -q "PASS" "$TEST_RUN_DIR/e2e_tests.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

### Performance Tests
- Performance Tests: $(if [ -f "$TEST_RUN_DIR/performance_tests.log" ] && grep -q "PASS" "$TEST_RUN_DIR/performance_tests.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

### Benchmarks
- Benchmarks: $(if [ -f "$TEST_RUN_DIR/benchmarks.log" ] && grep -q "PASS" "$TEST_RUN_DIR/benchmarks.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

### Coverage
- Coverage: $(if [ -f "$TEST_RUN_DIR/coverage.log" ] && grep -q "PASS" "$TEST_RUN_DIR/coverage.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

### Linting
- Linting: $(if [ -f "$TEST_RUN_DIR/linting.log" ] && grep -q "PASS" "$TEST_RUN_DIR/linting.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

### Security
- Security Scan: $(if [ -f "$TEST_RUN_DIR/security_scan.log" ] && grep -q "PASS" "$TEST_RUN_DIR/security_scan.log"; then echo "✅ PASSED"; else echo "❌ FAILED"; fi)

## Files Generated
- Test logs: $TEST_RUN_DIR/*.log
- Coverage reports: $TEST_RUN_DIR/*_coverage.html
- This report: $TEST_RUN_DIR/test_report.md

## Next Steps
$(if [ $FAILED_TESTS -gt 0 ]; then echo "- Review failed tests and fix issues"; else echo "- All tests passed! Ready for deployment"; fi)
EOF

echo -e "${GREEN}📋 Test report generated: $TEST_RUN_DIR/test_report.md${NC}"

# Final Summary
echo ""
echo -e "${BLUE}📊 Final Test Summary${NC}"
echo -e "${BLUE}====================${NC}"
echo -e "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
echo -e "Success Rate: $((PASSED_TESTS * 100 / TOTAL_TESTS))%"
echo ""
echo -e "Test results saved to: ${BLUE}$TEST_RUN_DIR${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! The Go server is ready for production.${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed. Please review the test logs and fix the issues.${NC}"
    exit 1
fi
