#!/bin/bash

# Salir inmediatamente si un comando falla
set -e

# Función para correr tests
run_tests() {
  echo "Running tests in $1..."
  cd $1
  go test ./...
  cd - > /dev/null
}

# Función para construir el binario
build_lambda() {
  echo "Building $1..."
  cd $1
  GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
  cd - > /dev/null
}

# Directorios de las lambdas
LAMBDA_DIRS=("lambdas/createDevice" "lambdas/deleteDevice" "lambdas/getDevice" "lambdas/updateDevice")

# Correr tests y construir binarios para cada Lambda
for dir in "${LAMBDA_DIRS[@]}"; do
  run_tests $dir
  build_lambda $dir
done

echo "All tests passed and lambdas are built successfully!"
