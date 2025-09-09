#!/bin/bash
set -e

echo "🔍 Running Go linting tools..."

echo "📝 Formatting code..."
go fmt ./...

echo "🔧 Running go vet..."
go vet ./...

echo "🧹 Cleaning up modules..."
go mod tidy

echo "🏃 Running tests..."
go test ./...

echo "🏁 Running tests with race detector..."
go test -race ./...

echo "📊 Running tests with coverage..."
go test -cover ./...

echo "🚀 Running benchmarks..."
go test -bench=. -run=^$

echo "✅ All linting checks passed!"