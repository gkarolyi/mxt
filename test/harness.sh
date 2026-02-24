#!/usr/bin/env bash
#
# Test Harness for mxt Feature Parity Validation
#
# Runs both muxtree and mxt with identical inputs, compares outputs,
# exit codes, and reports differences.
#
# Usage: ./test/harness.sh <command> <args...>
# Example: ./test/harness.sh new feature-branch
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'

# Check if running in TTY
if [[ ! -t 1 ]]; then
    RED='' GREEN='' YELLOW='' BLUE='' CYAN='' BOLD='' DIM='' RESET=''
fi

# Usage message
usage() {
    cat <<EOF
${BOLD}Test Harness for mxt Feature Parity Validation${RESET}

${BOLD}USAGE${RESET}
    $0 <command> [args...]

${BOLD}DESCRIPTION${RESET}
    Runs both muxtree and mxt with identical inputs and compares:
    - Standard output
    - Standard error
    - Exit codes
    - Side effects (where applicable)

${BOLD}EXAMPLES${RESET}
    $0 new feature-branch
    $0 list
    $0 config
    $0 delete feature-branch --force

${BOLD}ENVIRONMENT${RESET}
    MUXTREE_PATH    Path to muxtree binary (default: ./muxtree)
    MXT_PATH        Path to mxt binary (default: ./mxt)
    KEEP_TEMP       Set to 1 to keep temporary files for inspection

${BOLD}EXIT CODES${RESET}
    0   Outputs match (feature parity achieved)
    1   Outputs differ or error occurred
    2   Usage error
EOF
}

# Print colored message
info() { echo -e "${BLUE}▸${RESET} $*"; }
success() { echo -e "${GREEN}✓${RESET} $*"; }
warn() { echo -e "${YELLOW}⚠${RESET} $*"; }
error() { echo -e "${RED}✗${RESET} $*" >&2; }
die() { error "$*"; exit 1; }

# Check arguments
if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    usage
    exit 2
fi

# Paths to binaries
MUXTREE_PATH="${MUXTREE_PATH:-./muxtree}"
MXT_PATH="${MXT_PATH:-./mxt}"

# Check if muxtree exists
if [[ ! -x "$MUXTREE_PATH" ]]; then
    die "muxtree not found at $MUXTREE_PATH (set MUXTREE_PATH to override)"
fi

# Check if mxt exists
if [[ ! -x "$MXT_PATH" ]]; then
    warn "mxt not found at $MXT_PATH - will only run muxtree"
    MXT_MISSING=1
else
    MXT_MISSING=0
fi

# Create temp directory for outputs
TEMP_DIR=$(mktemp -d -t mxt-harness.XXXXXX)
if [[ "${KEEP_TEMP:-0}" != "1" ]]; then
    trap 'rm -rf "$TEMP_DIR"' EXIT
else
    info "Temporary files will be kept at: $TEMP_DIR"
fi

MUXTREE_STDOUT="$TEMP_DIR/muxtree.stdout"
MUXTREE_STDERR="$TEMP_DIR/muxtree.stderr"
MUXTREE_EXIT="$TEMP_DIR/muxtree.exit"

MXT_STDOUT="$TEMP_DIR/mxt.stdout"
MXT_STDERR="$TEMP_DIR/mxt.stderr"
MXT_EXIT="$TEMP_DIR/mxt.exit"

info "Running muxtree $*"
set +e
"$MUXTREE_PATH" "$@" >"$MUXTREE_STDOUT" 2>"$MUXTREE_STDERR"
echo $? >"$MUXTREE_EXIT"
set -e

MUXTREE_EXIT_CODE=$(cat "$MUXTREE_EXIT")
info "muxtree exit code: $MUXTREE_EXIT_CODE"

if [[ "$MXT_MISSING" == "1" ]]; then
    echo
    echo "${BOLD}muxtree output:${RESET}"
    cat "$MUXTREE_STDOUT"
    if [[ -s "$MUXTREE_STDERR" ]]; then
        echo
        echo "${BOLD}muxtree stderr:${RESET}"
        cat "$MUXTREE_STDERR"
    fi
    exit 0
fi

info "Running mxt $*"
set +e
"$MXT_PATH" "$@" >"$MXT_STDOUT" 2>"$MXT_STDERR"
echo $? >"$MXT_EXIT"
set -e

MXT_EXIT_CODE=$(cat "$MXT_EXIT")
info "mxt exit code: $MXT_EXIT_CODE"

echo
echo "════════════════════════════════════════════════════════════════"
echo "${BOLD}COMPARISON RESULTS${RESET}"
echo "════════════════════════════════════════════════════════════════"
echo

# Compare exit codes
if [[ "$MUXTREE_EXIT_CODE" != "$MXT_EXIT_CODE" ]]; then
    error "Exit codes differ!"
    echo "  muxtree: $MUXTREE_EXIT_CODE"
    echo "  mxt:     $MXT_EXIT_CODE"
    DIFFERS=1
else
    success "Exit codes match: $MUXTREE_EXIT_CODE"
    DIFFERS=0
fi

# Compare stdout
if diff -u "$MUXTREE_STDOUT" "$MXT_STDOUT" >/dev/null 2>&1; then
    success "Standard output matches"
else
    error "Standard output differs!"
    echo
    echo "${BOLD}Diff (muxtree vs mxt):${RESET}"
    diff -u "$MUXTREE_STDOUT" "$MXT_STDOUT" || true
    echo
    DIFFERS=1
fi

# Compare stderr
if diff -u "$MUXTREE_STDERR" "$MXT_STDERR" >/dev/null 2>&1; then
    success "Standard error matches"
else
    error "Standard error differs!"
    echo
    echo "${BOLD}Diff (muxtree vs mxt):${RESET}"
    diff -u "$MUXTREE_STDERR" "$MXT_STDERR" || true
    echo
    DIFFERS=1
fi

echo
if [[ "$DIFFERS" == "0" ]]; then
    echo "${GREEN}${BOLD}✓ FEATURE PARITY ACHIEVED${RESET}"
    echo "  All outputs match exactly"
    exit 0
else
    echo "${RED}${BOLD}✗ FEATURE PARITY NOT YET ACHIEVED${RESET}"
    echo "  Review differences above"
    echo
    echo "Temporary files saved at: $TEMP_DIR"
    echo "  muxtree stdout: $MUXTREE_STDOUT"
    echo "  muxtree stderr: $MUXTREE_STDERR"
    echo "  mxt stdout:     $MXT_STDOUT"
    echo "  mxt stderr:     $MXT_STDERR"
    trap - EXIT  # Don't delete temp files on failure
    exit 1
fi
