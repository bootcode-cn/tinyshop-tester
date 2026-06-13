.PHONY: build test clean run deps fmt test-solution test-pyodide

# Build tester
build:
	go build -o tinyshop-tester .

# Run unit tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f tinyshop-tester
	go clean

# Run tester with ARGS
run:
	go run . $(ARGS)

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Test all stages against local solution (Go tester)
test-solution:
	./scripts/test-solution.sh

# Test all stages in Pyodide (browser simulation)
test-pyodide:
	./scripts/test-pyodide.sh
