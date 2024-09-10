#!/bin/bash

# Run tests for all Go files in the project
./scripts/run_all_tests.sh

# Check if tests passed
if [ $? -ne 0 ]; then
  echo "Tests failed. Aborting builds."
  exit 1
fi

# Run builds for each lambda
echo "Building createDevice..."
./scripts/build_single_lambda.sh createDevice
if [ $? -ne 0 ]; then
  echo "Build failed for createDevice. Aborting further builds."
  exit 1
fi

echo "Building deleteDevice..."
./scripts/build_single_lambda.sh deleteDevice
if [ $? -ne 0 ]; then
  echo "Build failed for deleteDevice. Aborting further builds."
  exit 1
fi

echo "Building updateDevice..."
./scripts/build_single_lambda.sh updateDevice
if [ $? -ne 0 ]; then
  echo "Build failed for updateDevice. Aborting further builds."
  exit 1
fi

echo "Building getDevice..."
./scripts/build_single_lambda.sh getDevice
if [ $? -ne 0 ]; then
  echo "Build failed for getDevice. Aborting further builds."
  exit 1
fi

echo "Building homeDeviceListener..."
./scripts/build_single_lambda.sh homeDeviceListener
if [ $? -ne 0 ]; then
  echo "Build failed for homeDeviceListener. Aborting further builds."
  exit 1
fi

echo "All builds completed successfully."
