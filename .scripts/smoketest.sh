#!/usr/bin/env bash
set -uo pipefail

# Smoke test for bilingual_pdf CLI
# Run from the project root:
#   bash .scripts/smoketest.sh          # quick tests only (no network)
#   bash .scripts/smoketest.sh --full   # include translation tests (slow, needs network)
#   bash .scripts/smoketest.sh -f       # same as --full
#   bash .scripts/smoketest.sh --keep   # don't clean up generated files
#   bash .scripts/smoketest.sh --clean  # clean up and exit
#   bash .scripts/smoketest.sh --help   # show help

show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "OPTIONS:"
    echo "  -h, --help     Show this help message"
    echo "  -f, --full     Run full tests including network-dependent translation tests"
    echo "  -k, --keep     Keep generated files after tests complete"
    echo "  -c, --clean    Clean up generated files and exit (no tests)"
    echo ""
    echo "Default: Run quick tests only (no network) and clean up files"
}

clean_files() {
    echo "Cleaning up generated files..."
    rm -f testdata/*.pdf testdata/*.html testdata/*.en.md
}

PASS=0
FAIL=0
FULL=false
KEEP=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -f|--full)
            FULL=true
            shift
            ;;
        -k|--keep)
            KEEP=true
            shift
            ;;
        -c|--clean)
            clean_files
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Clean up before running tests
clean_files

run_expect_ok() {
    local desc="$1"; shift
    printf "  %-60s" "$desc"
    if output=$(go run . "$@" 2>&1); then
        echo "PASS"
        ((PASS++))
    else
        echo "FAIL (expected success, got exit code $?)"
        echo "    $output"
        ((FAIL++))
    fi
}

run_expect_fail() {
    local desc="$1"; shift
    printf "  %-60s" "$desc"
    if output=$(go run . "$@" 2>&1); then
        echo "FAIL (expected failure, got success)"
        echo "    $output"
        ((FAIL++))
    else
        echo "PASS"
        ((PASS++))
    fi
}

echo "=== Smoke tests ==="

echo ""
echo "--- Should succeed (quick) ---"
run_expect_ok "--help" --help
run_expect_ok "--list-languages" --list-languages
run_expect_ok "with --translation file (no network)" \
    testdata/sample.fr.md --translation testdata/sample.es.md
run_expect_ok "small fr->es" testdata/sample.fr.md --translation testdata/sample.es.md --font-size small --output testdata/sample.fr.es.s.pdf
run_expect_ok "medium fr->es" testdata/sample.fr.md --translation testdata/sample.es.md --font-size medium --output testdata/sample.fr.es.m.pdf
run_expect_ok "large fr->es" testdata/sample.fr.md --translation testdata/sample.es.md --font-size large --output testdata/sample.fr.es.l.pdf

echo ""
echo "--- Should fail ---"
run_expect_fail "no arguments"
run_expect_fail "nonexistent file" nonexistent.md
run_expect_fail "invalid source lang" testdata/sample.fr.md --source zz
run_expect_fail "invalid target lang" testdata/sample.fr.md --target zz
run_expect_fail "non-.md file" README.txt
run_expect_fail "nonexistent translation" testdata/sample.fr.md --translation missing.md
run_expect_fail "output not .pdf" testdata/sample.fr.md -o out.txt
run_expect_fail "invalid --font-size" testdata/sample.fr.md --font-size huge

if $FULL; then
    echo ""
    echo "--- Should succeed (full, uses network) ---"
    run_expect_ok "default fr->es" testdata/sample.fr.md
    run_expect_ok "explicit --source fr --target es" testdata/sample.fr.md --source fr --target es
    run_expect_ok "with --html --save-translation" testdata/sample.fr.md --html --save-translation
    run_expect_ok "with -o output.pdf" testdata/sample.fr.md -o testdata/smoketest-out.pdf
    run_expect_ok "fr->en" testdata/sample.fr.md --source fr --target en  --save-translation
    run_expect_ok "en->de" testdata/sample.en.md --source en --target de

fi

echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="

# Clean up generated files (unless --keep was specified)
if ! $KEEP; then
    clean_files
fi

exit $FAIL
