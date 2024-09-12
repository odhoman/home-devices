#!/bin/bash

# Check if the argument is provided
if [ -z "$1" ]; then
  echo "Please provide the lambda directory name (e.g., createDevice, deleteDevice, etc.)."
  exit 1
fi

LAMBDA_NAME=$1

# Define the path to the lambda folder
LAMBDA_PATH=$(dirname "$0")/../lambdas/cmd/$LAMBDA_NAME

# Change directory to the lambda folder
#cd "$LAMBDA_PATH" || exit

# Build the Go binary for the specified lambda
echo "Building lambda: $LAMBDA_NAME $LAMBDA_PATH"
#GOOS=linux GOARCH=amd64 go build -o "$LAMBDA_PATH/bootstrap" "$LAMBDA_NAME.go"

if [ $? -eq 0 ]; then
  echo "Lambda $LAMBDA_NAME built successfully."
else
  echo "Failed to build lambda $LAMBDA_NAME."
  exit 1
fi
