#!/bin/bash

# Check if the lambda name was provided
if [ -z "$1" ]; then
  echo "You must provide the name of the lambda (createDevice, deleteDevice, updateDevice, getDevice, homeDeviceListener)."
  exit 1
fi

LAMBDA_NAME=$1

# Define the base directory for the lambda folders
BASE_DIR=$(dirname "$0")/..
LAMBDA_DIR="$BASE_DIR/lambdas/$LAMBDA_NAME"

# Run tests for all Go files within the lambdas folder
echo "Running tests for all Go files within the lambdas folder..."
"$BASE_DIR/scripts/run_all_tests.sh"

# Check if the tests passed
if [ $? -ne 0 ]; then
  echo "Some tests failed. Stopping the script."
  exit 1
fi

# If the tests passed, build the specified lambda
echo "All tests passed, starting build of $LAMBDA_NAME.go..."
cd "$LAMBDA_DIR" || exit
GOOS=linux GOARCH=amd64 go build -o bootstrap "$LAMBDA_NAME.go"

if [ $? -eq 0 ]; then
  echo "Build of $LAMBDA_NAME.go completed successfully."
else
  echo "Error during the build of $LAMBDA_NAME.go."
  exit 1
fi

# Return to the original directory
cd "$CURRENT_DIR" || exit
