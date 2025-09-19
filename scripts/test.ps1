# Go Server Test Automation Script for Windows PowerShell
# This script runs all tests and generates a comprehensive report

param(
    [string]$TestType = "all",
    [switch]$Coverage = $false,
    [switch]$Performance = $false,
    [switch]$Verbose = $false
)

# Colors for output
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"
$Cyan = "Cyan"

# Test results directory
$TestResultsDir = "test-results"
if (!(Test-Path $TestResultsDir)) {
    New-Item -ItemType Directory -Path $TestResultsDir | Out-Null
}

# Timestamp for test run
$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$TestRunDir = Join-Path $TestResultsDir "test_run_$Timestamp"
New-Item -ItemType Directory -Path $TestRunDir | Out-Null

Write-Host "Starting Go Server Test Automation" -ForegroundColor $Blue
Write-Host "=====================================" -ForegroundColor $Blue
Write-Host "Test run: $Timestamp" -ForegroundColor $Cyan
Write-Host "Results directory: $TestRunDir" -ForegroundColor $Cyan
Write-Host ""

# Function to run tests and capture output
function Run-Test {
    param(
        [string]$TestName,
        [string]$TestCommand,
        [string]$OutputFile
    )
    
    Write-Host "Running $TestName..." -ForegroundColor $Yellow
    
    try {
        $result = Invoke-Expression $TestCommand 2>&1
        $result | Out-File -FilePath $OutputFile -Encoding UTF8
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "PASSED: $TestName" -ForegroundColor $Green
            return $true
        } else {
            Write-Host "FAILED: $TestName" -ForegroundColor $Red
            Write-Host "Check $OutputFile for details" -ForegroundColor $Red
            return $false
        }
    } catch {
        Write-Host "FAILED: $TestName with exception: $($_.Exception.Message)" -ForegroundColor $Red
        $_.Exception.Message | Out-File -FilePath $OutputFile -Encoding UTF8
        return $false
    }
}

# Function to run tests with coverage
function Run-TestWithCoverage {
    param(
        [string]$TestName,
        [string]$TestCommand,
        [string]$OutputFile,
        [string]$CoverageFile
    )
    
    Write-Host "Running $TestName with coverage..." -ForegroundColor $Yellow
    
    try {
        $result = Invoke-Expression "$TestCommand -coverprofile=$CoverageFile" 2>&1
        $result | Out-File -FilePath $OutputFile -Encoding UTF8
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "PASSED: $TestName" -ForegroundColor $Green
            
            # Generate coverage report
            $coverageHtml = Join-Path $TestRunDir "${TestName}_coverage.html"
            go tool cover -html=$CoverageFile -o $coverageHtml
            Write-Host "Coverage report: $coverageHtml" -ForegroundColor $Blue
            return $true
        } else {
            Write-Host "FAILED: $TestName" -ForegroundColor $Red
            Write-Host "Check $OutputFile for details" -ForegroundColor $Red
            return $false
        }
    } catch {
        Write-Host "FAILED: $TestName with exception: $($_.Exception.Message)" -ForegroundColor $Red
        $_.Exception.Message | Out-File -FilePath $OutputFile -Encoding UTF8
        return $false
    }
}

# Test results tracking
$TotalTests = 0
$PassedTests = 0
$FailedTests = 0

# 1. Unit Tests
if ($TestType -eq "all" -or $TestType -eq "unit") {
    Write-Host "Running Unit Tests" -ForegroundColor $Blue
    Write-Host "========================" -ForegroundColor $Blue
    
    # Config tests
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "config_tests.log"
    if (Run-Test "config_tests" "go test ./internal/config -v" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    # Errors tests
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "errors_tests.log"
    if (Run-Test "errors_tests" "go test ./internal/errors -v" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    # Middleware tests
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "middleware_tests.log"
    if (Run-Test "middleware_tests" "go test ./internal/middleware -v" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    # Server tests
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "server_tests.log"
    if (Run-Test "server_tests" "go test ./internal/server -v" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    Write-Host ""
}

# 2. Integration Tests
if ($TestType -eq "all" -or $TestType -eq "integration") {
    Write-Host "Running Integration Tests" -ForegroundColor $Blue
    Write-Host "===============================" -ForegroundColor $Blue
    
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "integration_tests.log"
    if (Run-Test "integration_tests" "go test ./test -v -run TestServer" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    Write-Host ""
}

# 3. End-to-End Tests
if ($TestType -eq "all" -or $TestType -eq "e2e") {
    Write-Host "Running End-to-End Tests" -ForegroundColor $Blue
    Write-Host "=============================" -ForegroundColor $Blue
    
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "e2e_tests.log"
    if (Run-Test "e2e_tests" "go test ./test -v -run TestHealthEndpoint -run TestAPIEndpoint -run TestVersionEndpoint -run TestMetricsEndpoint -run TestCORSEndpoint -run TestRequestSizeLimit" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    Write-Host ""
}

# 4. Performance Tests
if ($TestType -eq "all" -or $TestType -eq "performance" -or $Performance) {
    Write-Host "Running Performance Tests" -ForegroundColor $Blue
    Write-Host "===============================" -ForegroundColor $Blue
    
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "performance_tests.log"
    if (Run-Test "performance_tests" "go test ./test -v -run TestLoadTest -run TestMemoryUsage" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    Write-Host ""
}

# 5. Benchmarks
if ($TestType -eq "all" -or $TestType -eq "benchmark") {
    Write-Host "Running Benchmarks" -ForegroundColor $Blue
    Write-Host "=======================" -ForegroundColor $Blue
    
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "benchmarks.log"
    if (Run-Test "benchmarks" "go test ./test -bench=. -benchmem" $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    Write-Host ""
}

# 6. Coverage Analysis
if ($Coverage -or $TestType -eq "all" -or $TestType -eq "coverage") {
    Write-Host "Running Coverage Analysis" -ForegroundColor $Blue
    Write-Host "=============================" -ForegroundColor $Blue
    
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "coverage.log"
    $coverageFile = Join-Path $TestRunDir "coverage.out"
    if (Run-TestWithCoverage "coverage" "go test ./... -cover" $outputFile $coverageFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    Write-Host ""
}

# 7. Linting
if ($TestType -eq "all" -or $TestType -eq "lint") {
    Write-Host "Running Linting" -ForegroundColor $Blue
    Write-Host "===================" -ForegroundColor $Blue
    
    $TotalTests++
    $outputFile = Join-Path $TestRunDir "linting.log"
    if (Run-Test "linting" "go vet ./..." $outputFile) {
        $PassedTests++
    } else {
        $FailedTests++
    }
    
    Write-Host ""
}

# Generate Test Report
Write-Host "Generating Test Report" -ForegroundColor $Blue
Write-Host "============================" -ForegroundColor $Blue

$reportFile = Join-Path $TestRunDir "test_report.md"
$successRate = if ($TotalTests -gt 0) { [math]::Round(($PassedTests * 100 / $TotalTests), 2) } else { 0 }

$reportContent = @"
# Test Report - $Timestamp

## Summary
Total Tests: $TotalTests
Passed: $PassedTests
Failed: $FailedTests
Success Rate: $successRate%

## Test Results

### Unit Tests
Config Tests: $(if (Test-Path (Join-Path $TestRunDir "config_tests.log")) { if ((Get-Content (Join-Path $TestRunDir "config_tests.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })
Errors Tests: $(if (Test-Path (Join-Path $TestRunDir "errors_tests.log")) { if ((Get-Content (Join-Path $TestRunDir "errors_tests.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })
Middleware Tests: $(if (Test-Path (Join-Path $TestRunDir "middleware_tests.log")) { if ((Get-Content (Join-Path $TestRunDir "middleware_tests.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })
Server Tests: $(if (Test-Path (Join-Path $TestRunDir "server_tests.log")) { if ((Get-Content (Join-Path $TestRunDir "server_tests.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })

### Integration Tests
Integration Tests: $(if (Test-Path (Join-Path $TestRunDir "integration_tests.log")) { if ((Get-Content (Join-Path $TestRunDir "integration_tests.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })

### End-to-End Tests
E2E Tests: $(if (Test-Path (Join-Path $TestRunDir "e2e_tests.log")) { if ((Get-Content (Join-Path $TestRunDir "e2e_tests.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })

### Performance Tests
Performance Tests: $(if (Test-Path (Join-Path $TestRunDir "performance_tests.log")) { if ((Get-Content (Join-Path $TestRunDir "performance_tests.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })

### Benchmarks
Benchmarks: $(if (Test-Path (Join-Path $TestRunDir "benchmarks.log")) { if ((Get-Content (Join-Path $TestRunDir "benchmarks.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })

### Coverage
Coverage: $(if (Test-Path (Join-Path $TestRunDir "coverage.log")) { if ((Get-Content (Join-Path $TestRunDir "coverage.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })

### Linting
Linting: $(if (Test-Path (Join-Path $TestRunDir "linting.log")) { if ((Get-Content (Join-Path $TestRunDir "linting.log") | Select-String "PASS").Count -gt 0) { "PASSED" } else { "FAILED" } } else { "NOT RUN" })

## Files Generated
Test logs: $TestRunDir\*.log
Coverage reports: $TestRunDir\*_coverage.html
This report: $TestRunDir\test_report.md

## Next Steps
$(if ($FailedTests -gt 0) { "Review failed tests and fix issues" } else { "All tests passed! Ready for deployment" })
"@

$reportContent | Out-File -FilePath $reportFile -Encoding UTF8

Write-Host "Test report generated: $reportFile" -ForegroundColor $Green

# Final Summary
Write-Host ""
Write-Host "Final Test Summary" -ForegroundColor $Blue
Write-Host "====================" -ForegroundColor $Blue
Write-Host "Total Tests: $TotalTests"
Write-Host "Passed: $PassedTests" -ForegroundColor $Green
Write-Host "Failed: $FailedTests" -ForegroundColor $Red
Write-Host "Success Rate: $successRate%"
Write-Host ""
Write-Host "Test results saved to: $TestRunDir" -ForegroundColor $Blue
Write-Host ""

if ($FailedTests -eq 0) {
    Write-Host "All tests passed! The Go server is ready for production." -ForegroundColor $Green
    exit 0
} else {
    Write-Host "Some tests failed. Please review the test logs and fix the issues." -ForegroundColor $Red
    exit 1
}
