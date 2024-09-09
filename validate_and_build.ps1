# Salir si un comando falla
$ErrorActionPreference = "Stop"

# Función para correr tests
function Run-Tests {
    param (
        [string]$Directory
    )
    Write-Host "Running tests in $Directory..."
    Push-Location $Directory
    go test ./...
    Pop-Location
}

# Función para construir el binario
function Build-Lambda {
    param (
        [string]$Directory
    )
    Write-Host "Building $Directory..."
    Push-Location $Directory
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    go build -o bootstrap main.go
    Pop-Location
}

# Directorios de las lambdas
$lambdaDirs = @("lambdas/createDevice", "lambdas/deleteDevice", "lambdas/getDevice", "lambdas/updateDevice")

# Correr tests y construir binarios para cada Lambda
foreach ($dir in $lambdaDirs) {
    Run-Tests -Directory $dir
    Build-Lambda -Directory $dir
}

Write-Host "All tests passed and lambdas are built successfully!"
