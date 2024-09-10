#!/bin/bash

# Define the base path for the lambda folders
BASE_DIR=$(dirname "$0")/..
LAMBDA_DIR="$BASE_DIR/lambdas"

# Save the current directory
CURRENT_DIR=$(pwd)

# Function to recursively traverse directories and execute Go test files
run_tests_in_dir() {
  local dir=$1

  for sub_dir in "$dir"/*; do
    if [ -d "$sub_dir" ]; then
      # If it's a directory, traverse it recursively
      run_tests_in_dir "$sub_dir"
    fi
  done

  # Look for test files in this directory
  test_files=$(find "$dir" -maxdepth 1 -type f -name "*_test.go")
  if [ -n "$test_files" ]; then
    echo "Running tests in $dir"
    cd "$dir" || exit
    go test -v
    if [ $? -ne 0 ]; then
      echo "Tests failed in $dir"
      exit 1
    fi
    cd "$CURRENT_DIR" || exit
  fi
}

# Call the function to run tests in all lambda directories
run_tests_in_dir "$LAMBDA_DIR"
