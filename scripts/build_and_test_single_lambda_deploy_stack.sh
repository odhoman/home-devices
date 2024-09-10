#!/bin/bash

# Run the build and test script for the specified lambda function
./scripts/build_and_test_single_lambda.sh $1

# Check if the build and test script succeeded
if [ $? -ne 0 ]; then
  echo "Build and test for $1 failed. Aborting deployment."
  exit 1
fi

# If the build and test succeeded, proceed with CDK deploy
echo "Build and test for $1 succeeded. Deploying stack..."
cdk deploy

# Check if the deployment was successful
if [ $? -eq 0 ]; then
  echo "CDK deployment completed successfully."
else
  echo "CDK deployment failed."
  exit 1
fi
