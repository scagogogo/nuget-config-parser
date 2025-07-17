#!/bin/bash

# NuGet Config Parser - .gitignore validation script
# This script checks if important files are properly ignored

set -e

echo "🔍 Checking .gitignore effectiveness..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a pattern is ignored
check_ignored() {
    local pattern="$1"
    local description="$2"
    
    if git check-ignore "$pattern" >/dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} $description: $pattern"
        return 0
    else
        echo -e "${RED}✗${NC} $description: $pattern (NOT IGNORED)"
        return 1
    fi
}

# Function to create test file and check if ignored
test_ignore() {
    local file="$1"
    local description="$2"
    local dir=$(dirname "$file")

    # Create directory if it doesn't exist
    mkdir -p "$dir" 2>/dev/null || true

    # Create test file
    touch "$file"

    # Check if ignored and get the source
    local ignore_output=$(git check-ignore -v "$file" 2>/dev/null)
    if [ $? -eq 0 ]; then
        local source=$(echo "$ignore_output" | cut -d: -f1)
        if [[ "$source" == ".gitignore" ]]; then
            echo -e "${GREEN}✓${NC} $description: $file (local .gitignore)"
        else
            echo -e "${YELLOW}✓${NC} $description: $file (${source})"
        fi
        rm -f "$file"
        return 0
    else
        echo -e "${RED}✗${NC} $description: $file (NOT IGNORED)"
        rm -f "$file"
        return 1
    fi
}

echo ""
echo "Testing Go-specific files..."

# Test Go files
test_ignore "main" "Go binary"
test_ignore "test.exe" "Windows executable"
test_ignore "coverage.out" "Go coverage file"
test_ignore "test.prof" "Go profile file"
test_ignore "benchmark.bench" "Go benchmark file"

echo ""
echo "Testing documentation files..."

# Test documentation files
if [ -d "docs" ]; then
    test_ignore "docs/node_modules/test" "Node modules"
    test_ignore "docs/.vitepress/dist/index.html" "VitePress build output"
    test_ignore "docs/.vitepress/cache/test" "VitePress cache"
fi

echo ""
echo "Testing IDE files..."

# Test IDE files
test_ignore ".vscode/settings.json" "VS Code settings"
test_ignore ".idea/workspace.xml" "IntelliJ workspace"
test_ignore "test.swp" "Vim swap file"

echo ""
echo "Testing OS files..."

# Test OS files
test_ignore ".DS_Store" "macOS metadata"
test_ignore "Thumbs.db" "Windows thumbnail cache"
test_ignore "desktop.ini" "Windows folder config"

echo ""
echo "Testing project-specific files..."

# Test project-specific files
mkdir -p "examples/test"
test_ignore "examples/test/NuGet.Config" "Example NuGet config"
test_ignore "examples/test/nuget.config" "Example nuget config (lowercase)"
test_ignore "temp-test.config" "Temporary config file"
test_ignore "test.tmp.config" "Test temporary config"
rmdir "examples/test" 2>/dev/null || true

echo ""
echo "Testing files that SHOULD be tracked..."

# Test files that should NOT be ignored
should_not_ignore() {
    local file="$1"
    local description="$2"
    
    if git check-ignore "$file" >/dev/null 2>&1; then
        echo -e "${RED}✗${NC} $description: $file (SHOULD NOT BE IGNORED)"
        return 1
    else
        echo -e "${GREEN}✓${NC} $description: $file"
        return 0
    fi
}

should_not_ignore "README.md" "README file"
should_not_ignore "go.mod" "Go module file"
should_not_ignore "pkg/nuget/api.go" "Go source file"
if [ -f "docs/package.json" ]; then
    should_not_ignore "docs/package.json" "Package.json"
fi
if [ -f "docs/package-lock.json" ]; then
    should_not_ignore "docs/package-lock.json" "Package-lock.json"
fi

echo ""
echo "🎯 Checking for common mistakes..."

# Check if go.sum is ignored (it shouldn't be)
if git check-ignore "go.sum" >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠${NC}  go.sum is ignored - this might not be desired for reproducible builds"
else
    echo -e "${GREEN}✓${NC} go.sum is tracked (good for reproducible builds)"
fi

# Check for any tracked files that might should be ignored
echo ""
echo "🔍 Checking for potentially problematic tracked files..."

problematic_files=$(git ls-files | grep -E '\.(log|tmp|cache|DS_Store)$|node_modules|\.vitepress/dist' || true)
if [ -n "$problematic_files" ]; then
    echo -e "${YELLOW}⚠${NC}  Found potentially problematic tracked files:"
    echo "$problematic_files"
    echo "Consider adding these to .gitignore and removing from tracking with:"
    echo "git rm --cached <file>"
else
    echo -e "${GREEN}✓${NC} No problematic tracked files found"
fi

echo ""
echo "✅ .gitignore validation complete!"
echo ""
echo "💡 Tips:"
echo "   - To see all ignored files: git status --ignored"
echo "   - To check if a specific file is ignored: git check-ignore -v <file>"
echo "   - To force add an ignored file: git add -f <file>"
