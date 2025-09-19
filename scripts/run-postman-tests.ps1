# Postman Test Runner Script for PowerShell
# Requires Newman (npm install -g newman)

param(
    [string]$BaseUrl = "http://localhost:8080",
    [string]$ReportDir = "test-results/postman",
    [switch]$Verbose
)

# Configuration
$CollectionFile = "postman/Go-Server-API.postman_collection.json"
$EnvironmentFile = "postman/Go-Server-Environment.postman_environment.json"
$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$ReportFile = "$ReportDir/postman_report_$Timestamp"

# Create report directory
if (!(Test-Path $ReportDir)) {
    New-Item -ItemType Directory -Path $ReportDir -Force | Out-Null
}

Write-Host "🚀 Starting Postman Test Suite" -ForegroundColor Blue
Write-Host "======================================"
Write-Host "Collection: $CollectionFile"
Write-Host "Environment: $EnvironmentFile"
Write-Host "Base URL: $BaseUrl"
Write-Host "Report Directory: $ReportDir"
Write-Host ""

# Check if Newman is installed
try {
    $null = Get-Command newman -ErrorAction Stop
} catch {
    Write-Host "❌ Newman is not installed. Please install it with:" -ForegroundColor Red
    Write-Host "npm install -g newman" -ForegroundColor Yellow
    exit 1
}

# Check if collection file exists
if (!(Test-Path $CollectionFile)) {
    Write-Host "❌ Collection file not found: $CollectionFile" -ForegroundColor Red
    exit 1
}

# Check if environment file exists
$UseEnvironment = Test-Path $EnvironmentFile
if (!$UseEnvironment) {
    Write-Host "⚠️  Environment file not found: $EnvironmentFile" -ForegroundColor Yellow
    Write-Host "Running without environment file..." -ForegroundColor Yellow
}

# Run Newman tests
Write-Host "📋 Running Postman Collection Tests..." -ForegroundColor Blue

$NewmanArgs = @(
    "run", $CollectionFile
    "--reporters", "cli,json,html"
    "--reporter-json-export", "$ReportFile.json"
    "--reporter-html-export", "$ReportFile.html"
    "--timeout-request", "10000"
    "--timeout-script", "5000"
    "--delay-request", "100"
)

if ($UseEnvironment) {
    $NewmanArgs += @("--environment", $EnvironmentFile)
}

if ($Verbose) {
    $NewmanArgs += "--verbose"
}

try {
    & newman @NewmanArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ All Postman tests passed!" -ForegroundColor Green
        Write-Host "Reports saved to:" -ForegroundColor Green
        Write-Host "  - JSON: $ReportFile.json" -ForegroundColor Green
        Write-Host "  - HTML: $ReportFile.html" -ForegroundColor Green
    } else {
        Write-Host "❌ Some Postman tests failed!" -ForegroundColor Red
        Write-Host "Check the reports for details:" -ForegroundColor Red
        Write-Host "  - JSON: $ReportFile.json" -ForegroundColor Red
        Write-Host "  - HTML: $ReportFile.html" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ Error running Newman: $_" -ForegroundColor Red
    exit 1
}
