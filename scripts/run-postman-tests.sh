#!/bin/bash

# Postman Test Runner Script
# Requires Newman (npm install -g newman)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COLLECTION_FILE="postman/Go-Server-API.postman_collection.json"
ENVIRONMENT_FILE="postman/Go-Server-Environment.postman_environment.json"
REPORT_DIR="test-results/postman"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
REPORT_FILE="$REPORT_DIR/postman_report_$TIMESTAMP"

# Create report directory
mkdir -p "$REPORT_DIR"

echo -e "${BLUE}🚀 Starting Postman Test Suite${NC}"
echo "======================================"
echo "Collection: $COLLECTION_FILE"
echo "Environment: $ENVIRONMENT_FILE"
echo "Report Directory: $REPORT_DIR"
echo ""

# Check if Newman is installed
if ! command -v newman &> /dev/null; then
    echo -e "${RED}❌ Newman is not installed. Please install it with:${NC}"
    echo "npm install -g newman"
    exit 1
fi

# Check if collection file exists
if [ ! -f "$COLLECTION_FILE" ]; then
    echo -e "${RED}❌ Collection file not found: $COLLECTION_FILE${NC}"
    exit 1
fi

# Check if environment file exists
if [ ! -f "$ENVIRONMENT_FILE" ]; then
    echo -e "${YELLOW}⚠️  Environment file not found: $ENVIRONMENT_FILE${NC}"
    echo "Running without environment file..."
    ENVIRONMENT_FILE=""
fi

# Run Newman tests
echo -e "${BLUE}📋 Running Postman Collection Tests...${NC}"

if [ -n "$ENVIRONMENT_FILE" ]; then
    newman run "$COLLECTION_FILE" \
        --environment "$ENVIRONMENT_FILE" \
        --reporters cli,json,html \
        --reporter-json-export "$REPORT_FILE.json" \
        --reporter-html-export "$REPORT_FILE.html" \
        --timeout-request 10000 \
        --timeout-script 5000 \
        --delay-request 100
else
    newman run "$COLLECTION_FILE" \
        --reporters cli,json,html \
        --reporter-json-export "$REPORT_FILE.json" \
        --reporter-html-export "$REPORT_FILE.html" \
        --timeout-request 10000 \
        --timeout-script 5000 \
        --delay-request 100
fi

# Check exit code
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ All Postman tests passed!${NC}"
    echo "Reports saved to:"
    echo "  - JSON: $REPORT_FILE.json"
    echo "  - HTML: $REPORT_FILE.html"
else
    echo -e "${RED}❌ Some Postman tests failed!${NC}"
    echo "Check the reports for details:"
    echo "  - JSON: $REPORT_FILE.json"
    echo "  - HTML: $REPORT_FILE.html"
    exit 1
fi
