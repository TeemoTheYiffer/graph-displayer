#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Log file
LOG_FILE="build.log"

# Function to print status messages
print_status() {
    echo -e "${YELLOW}==> ${NC}$1" | tee -a "$LOG_FILE"
}

print_success() {
    echo -e "${GREEN}==> SUCCESS: ${NC}$1" | tee -a "$LOG_FILE"
}

print_error() {
    echo -e "${RED}==> ERROR: ${NC}$1" | tee -a "$LOG_FILE"
}

# Function to check if required pngs exist
check_pngs() {
    print_status "Checking for required PNG files..."
    required_pngs=("Pie.png")
    missing_pngs=()
    
    for png in "${required_pngs[@]}"; do
        if [ ! -f "ui/examples/$png" ]; then
            missing_pngs+=("$png")
        fi
    done
    
    if [ ${#missing_pngs[@]} -ne 0 ]; then
        print_error "Missing png files in ui/examples directory:"
        for png in "${missing_pngs[@]}"; do
            echo "  - $png" | tee -a "$LOG_FILE"
        done
        return 1
    fi
    
    print_success "All required png files found"
    return 0
}

# Function to run tests
run_tests() {
    print_status "Running tests with native compiler..."
    
    OLD_CC=$CC
    export CC=gcc
    CGO_ENABLED=1 go test -v ./ui 2>&1 | tee -a "$LOG_FILE"
    test_result=$?
    export CC=$OLD_CC
    
    if [ $test_result -ne 0 ]; then
        print_error "Tests failed"
        return 1
    fi
    
    print_success "Tests passed"
    return 0
}

# Function to build for Windows
build_windows() {
    print_status "Building Windows executable..."
    print_status "Current directory structure:"
    ls -R | tee -a "$LOG_FILE"
    
    if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
        print_error "Windows cross-compiler not found. Please install mingw-w64"
        return 1
    fi
    
    print_status "Using compiler: x86_64-w64-mingw32-gcc"
    print_status "Building with verbose output..."
    
    export CC=x86_64-w64-mingw32-gcc
    GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -v . 2>&1 | tee -a "$LOG_FILE"
    build_result=$?
    
    if [ $build_result -ne 0 ]; then
        print_error "Build failed. Check build.log for details"
        print_status "Last few lines of build output:"
        tail -n 20 "$LOG_FILE"
        return 1
    fi
    
    print_success "Build completed successfully"
    return 0
}

# Main execution
main() {
    # Clear or create log file
    > "$LOG_FILE"
    
    print_status "Starting build process..."
    print_status "Go version:"
    go version 2>&1 | tee -a "$LOG_FILE"
    
    if [ ! -d "ui" ]; then
        print_error "Please run this script from the project root directory"
        exit 1
    fi
    
    check_pngs || exit 1
    run_tests || exit 1
    build_windows || exit 1
    
    print_success "All tasks completed successfully!"
}

main