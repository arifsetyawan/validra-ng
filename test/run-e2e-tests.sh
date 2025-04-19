#!/bin/bash

# Colors for terminal output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Log message with timestamp
log() {
    echo -e "${YELLOW}[$(date +"%Y-%m-%d %H:%M:%S")]${NC} $1"
}

# Check if Newman is installed
if ! command -v newman &> /dev/null; then
    log "${RED}Newman is not installed. Please install it using 'npm install -g newman'.${NC}"
    exit 1
fi

# Default port (can be overridden with -p flag)
PORT=8080

# Parse command line arguments
while getopts ":p:" opt; do
  case $opt in
    p) PORT="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    exit 1
    ;;
  esac
done

# Set API base URL
API_URL="http://localhost:${PORT}"

# Update the environment file with the correct port
log "Setting API base URL to ${API_URL}"
tmp=$(mktemp)
jq ".values[0].value = \"$API_URL\"" ./test/validra-environment.json > "$tmp" && mv "$tmp" ./test/validra-environment.json

# Check if the API is running
log "Checking if the API is running at ${API_URL}/health..."
if ! curl -s -o /dev/null -w "%{http_code}" "${API_URL}/health" | grep -q "200"; then
    log "${RED}API is not running at ${API_URL}. Please start the API server before running tests.${NC}"
    exit 1
fi

log "${GREEN}API is running. Starting E2E tests...${NC}"

# Run Newman with the collection and environment files
newman run ./test/validra-api.collection.json \
    --environment ./test/validra-environment.json \
    --reporters cli,junit \
    --reporter-junit-export ./test/test-results.xml

# Check the exit code
if [ $? -eq 0 ]; then
    log "${GREEN}All tests passed successfully!${NC}"
    exit 0
else
    log "${RED}Some tests failed. Please check the test results.${NC}"
    exit 1
fi