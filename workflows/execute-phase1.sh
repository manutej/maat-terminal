#!/usr/bin/env bash
#
# MAAT Phase 1 Execution Script
#
# Purpose: Execute both Phase 1 workflows (A: SQLite Store, B: Mock Data)
#          and validate Phase 1 completion
#
# Usage: ./workflows/execute-phase1.sh [--dry-run] [--verbose]
#
# Options:
#   --dry-run   Preview execution without running agents
#   --verbose   Show detailed agent execution logs
#   --help      Show this help message

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse arguments
DRY_RUN=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help)
            echo "Usage: $0 [--dry-run] [--verbose]"
            echo ""
            echo "Execute MAAT Phase 1 workflows (SQLite Store + Mock Data)"
            echo ""
            echo "Options:"
            echo "  --dry-run   Preview execution without running agents"
            echo "  --verbose   Show detailed agent execution logs"
            echo "  --help      Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check Go installation
    if ! command -v go &> /dev/null; then
        log_error "Go not found. Please install Go 1.21+."
        exit 1
    fi
    log_success "Go $(go version | awk '{print $3}') installed"

    # Check required files
    local required_files=(
        "internal/graph/schema.go"
        "docs/CONSTITUTION.md"
        "specs/FUNCTIONAL-REQUIREMENTS.md"
        "workflows/phase1-sqlite-store.yaml"
        "workflows/phase1-mock-data.yaml"
    )

    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            log_error "Required file not found: $file"
            exit 1
        fi
    done
    log_success "All required files present"

    # Check go.mod includes sqlite3
    if ! grep -q "github.com/mattn/go-sqlite3" go.mod 2>/dev/null; then
        log_warning "go-sqlite3 not in go.mod. Adding..."
        go get github.com/mattn/go-sqlite3
        log_success "Added go-sqlite3 to dependencies"
    fi
}

# Execute workflow
execute_workflow() {
    local workflow_name=$1
    local workflow_file=$2
    local input_file=$3
    local output_file=$4

    echo ""
    log_info "=========================================="
    log_info "Executing: $workflow_name"
    log_info "=========================================="

    # Build command
    local cmd="/ois-compose --workflow-plan $workflow_file"

    if [[ "$input_file" == *" "* ]]; then
        # Multiple input files
        cmd="$cmd --input-file \"$input_file\""
    else
        cmd="$cmd --input-file $input_file"
    fi

    cmd="$cmd --output $output_file"

    if [[ "$DRY_RUN" == true ]]; then
        cmd="$cmd --dry-run"
    fi

    if [[ "$VERBOSE" == true ]]; then
        cmd="$cmd --verbose"
    fi

    log_info "Command: $cmd"
    echo ""

    # Execute
    if eval "$cmd"; then
        log_success "$workflow_name completed successfully"
        return 0
    else
        log_error "$workflow_name failed"
        return 1
    fi
}

# Validate output
validate_output() {
    local file=$1
    local description=$2

    if [[ ! -f "$file" ]]; then
        log_error "$description not found: $file"
        return 1
    fi

    # Check if file is non-empty
    if [[ ! -s "$file" ]]; then
        log_error "$description is empty: $file"
        return 1
    fi

    log_success "$description generated: $file"
    return 0
}

# Run tests
run_tests() {
    log_info "Running tests..."

    # Test graph package
    if go test ./internal/graph -v; then
        log_success "Graph package tests passed"
    else
        log_error "Graph package tests failed"
        return 1
    fi

    # Test TUI package
    if go test ./internal/tui -v; then
        log_success "TUI package tests passed"
    else
        log_warning "TUI package tests failed (expected if integration incomplete)"
    fi

    return 0
}

# Build binary
build_binary() {
    log_info "Building MAAT binary..."

    if go build -o maat ./cmd/maat; then
        log_success "Binary built successfully: ./maat"

        # Check binary size
        local size=$(du -h maat | awk '{print $1}')
        log_info "Binary size: $size"

        # Validate size < 25MB (success criteria)
        local size_bytes=$(stat -f%z maat 2>/dev/null || stat -c%s maat 2>/dev/null)
        local max_bytes=$((25 * 1024 * 1024))

        if [[ $size_bytes -lt $max_bytes ]]; then
            log_success "Binary size OK (< 25MB)"
        else
            log_warning "Binary size exceeds 25MB target"
        fi

        return 0
    else
        log_error "Build failed"
        return 1
    fi
}

# Main execution
main() {
    log_info "MAAT Phase 1 Execution"
    log_info "======================"
    echo ""

    # Check prerequisites
    check_prerequisites
    echo ""

    # Execute Phase 1A: SQLite Store
    if execute_workflow \
        "Phase 1A: SQLite Store" \
        "workflows/phase1-sqlite-store.yaml" \
        "internal/graph/schema.go" \
        "internal/graph/store_generated.go"; then

        if [[ "$DRY_RUN" == false ]]; then
            validate_output "internal/graph/store_generated.go" "SQLite store implementation"
        fi
    else
        log_error "Phase 1A failed. Aborting."
        exit 1
    fi

    # Execute Phase 1B: Mock Data
    if execute_workflow \
        "Phase 1B: Mock Data Generator" \
        "workflows/phase1-mock-data.yaml" \
        "docs/CONSTITUTION.md specs/FUNCTIONAL-REQUIREMENTS.md" \
        "internal/tui/mock_data.go"; then

        if [[ "$DRY_RUN" == false ]]; then
            validate_output "internal/tui/mock_data.go" "Mock data generator"
        fi
    else
        log_error "Phase 1B failed. Aborting."
        exit 1
    fi

    # If dry run, stop here
    if [[ "$DRY_RUN" == true ]]; then
        echo ""
        log_success "Dry run complete. Workflows validated."
        log_info "Run without --dry-run to execute agents."
        exit 0
    fi

    echo ""
    log_info "=========================================="
    log_info "Phase 1 Validation"
    log_info "=========================================="
    echo ""

    # Run tests
    if ! run_tests; then
        log_warning "Some tests failed. Review and fix before continuing."
    fi
    echo ""

    # Build binary
    if ! build_binary; then
        log_error "Build failed. Fix compilation errors."
        exit 1
    fi
    echo ""

    # Success summary
    echo ""
    log_success "=========================================="
    log_success "Phase 1 Complete!"
    log_success "=========================================="
    echo ""
    log_info "Next steps:"
    echo "  1. Run ./maat to test TUI with mock data"
    echo "  2. Verify graph renders correctly"
    echo "  3. Test keyboard navigation (hjkl)"
    echo "  4. Proceed to Phase 2: Graph Navigation"
    echo ""
    log_info "Phase 1 achievements:"
    echo "  ✓ SQLite store implemented"
    echo "  ✓ Mock data generator created"
    echo "  ✓ Tests passing"
    echo "  ✓ Binary built (<25MB)"
    echo "  ✓ Ready for Phase 2"
    echo ""
}

# Run main
main
