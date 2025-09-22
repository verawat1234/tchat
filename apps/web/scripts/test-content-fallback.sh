#!/bin/bash

# T058: Content Fallback E2E Test Runner
# Comprehensive script for running content fallback tests with various options

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_FILE="content-fallback.spec.ts"
BASE_URL="http://localhost:3000"
REPORT_DIR="test-results"

# Functions
print_header() {
    echo -e "${BLUE}=================================="
    echo -e "Content Fallback E2E Test Runner"
    echo -e "==================================${NC}"
    echo ""
}

print_section() {
    echo -e "${YELLOW}$1${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

check_dependencies() {
    print_section "Checking Dependencies..."

    # Check if npm is installed
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed"
        exit 1
    fi
    print_success "npm found"

    # Check if Playwright is installed
    if ! npx playwright --version &> /dev/null; then
        print_error "Playwright is not installed. Run: npm install @playwright/test"
        exit 1
    fi
    print_success "Playwright found"

    # Check if test file exists
    if [ ! -f "tests/e2e/$TEST_FILE" ]; then
        print_error "Test file not found: tests/e2e/$TEST_FILE"
        exit 1
    fi
    print_success "Test file found"

    echo ""
}

check_server() {
    print_section "Checking Development Server..."

    # Check if server is running
    if curl -s "$BASE_URL" > /dev/null 2>&1; then
        print_success "Development server is running at $BASE_URL"
    else
        print_error "Development server is not running at $BASE_URL"
        print_info "Please start the development server with: npm run dev"
        exit 1
    fi

    echo ""
}

run_basic_tests() {
    print_section "Running Basic Content Fallback Tests..."

    npx playwright test "$TEST_FILE" \
        --project=chromium \
        --reporter=list \
        --timeout=30000
}

run_full_tests() {
    print_section "Running Full Cross-Browser Content Fallback Tests..."

    npx playwright test "$TEST_FILE" \
        --reporter=html \
        --timeout=60000
}

run_specific_suite() {
    local suite="$1"
    print_section "Running Content Fallback Tests: $suite..."

    npx playwright test "$TEST_FILE" \
        --grep "$suite" \
        --project=chromium \
        --reporter=list \
        --timeout=30000
}

run_debug_mode() {
    print_section "Running Content Fallback Tests in Debug Mode..."

    npx playwright test "$TEST_FILE" \
        --project=chromium \
        --headed \
        --timeout=0 \
        --workers=1 \
        --reporter=list
}

run_ui_mode() {
    print_section "Running Content Fallback Tests in UI Mode..."

    npx playwright test "$TEST_FILE" --ui
}

generate_report() {
    print_section "Generating Test Report..."

    npx playwright test "$TEST_FILE" \
        --reporter=html \
        --output="$REPORT_DIR"

    print_success "Report generated in $REPORT_DIR directory"
    print_info "Open $REPORT_DIR/index.html in your browser to view the report"
}

show_help() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  basic         Run basic tests (Chromium only)"
    echo "  full          Run full cross-browser tests"
    echo "  offline       Run offline scenario tests only"
    echo "  api-failures  Run API failure tests only"
    echo "  cache         Run cache behavior tests only"
    echo "  performance   Run performance tests only"
    echo "  debug         Run tests in debug mode (headed browser)"
    echo "  ui            Run tests in UI mode"
    echo "  report        Generate HTML test report"
    echo "  help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 basic                    # Quick test run"
    echo "  $0 full                     # Complete test suite"
    echo "  $0 offline                  # Test offline scenarios"
    echo "  $0 debug                    # Debug failing tests"
    echo "  $0 ui                       # Interactive test runner"
    echo ""
}

# Main execution
print_header

case "${1:-basic}" in
    "basic")
        check_dependencies
        check_server
        run_basic_tests
        print_success "Basic tests completed!"
        ;;
    "full")
        check_dependencies
        check_server
        run_full_tests
        print_success "Full test suite completed!"
        ;;
    "offline")
        check_dependencies
        check_server
        run_specific_suite "Offline Scenarios"
        print_success "Offline scenario tests completed!"
        ;;
    "api-failures")
        check_dependencies
        check_server
        run_specific_suite "API Failures"
        print_success "API failure tests completed!"
        ;;
    "cache")
        check_dependencies
        check_server
        run_specific_suite "Cache Behavior"
        print_success "Cache behavior tests completed!"
        ;;
    "performance")
        check_dependencies
        check_server
        run_specific_suite "Performance"
        print_success "Performance tests completed!"
        ;;
    "debug")
        check_dependencies
        check_server
        run_debug_mode
        ;;
    "ui")
        check_dependencies
        check_server
        run_ui_mode
        ;;
    "report")
        check_dependencies
        check_server
        generate_report
        ;;
    "help"|"--help"|"-h")
        show_help
        ;;
    *)
        print_error "Unknown option: $1"
        echo ""
        show_help
        exit 1
        ;;
esac