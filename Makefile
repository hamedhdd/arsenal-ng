.PHONY: build run test vet lint clean install

# Build the binary into bin/
build:
	go build -ldflags "-s -w" -o bin/arsenal-ng ./cmd/arsenal-ng

# Run directly without producing a binary
run:
	go run ./cmd/arsenal-ng

# Run all tests
test:
	go test -v -count=1 ./...

# Run go vet
vet:
	go vet ./...

# Run golangci-lint (requires golangci-lint to be installed)
lint:
	golangci-lint run ./...

# Remove build artifacts
clean:
	rm -rf bin/ dist/ *.exe

# Install into $GOPATH/bin
install:
	go install ./cmd/arsenal-ng
