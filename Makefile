.PHONY: lint proto clean build run test

# Variables
GO              := go
GOTEST          := $(GO) test
GOVET           := $(GO) vet
BINARY_NAME     := shortlink-gateway
PROTO_DIR       := proto
PROTO_OUT_DIR   := proto
GO_OUT_DIR      := .

# Tools
GOLINT          := golangci-lint
PROTOC          := protoc
PROTOC_GEN_GO   := protoc-gen-go
PROTOC_GEN_GRPC := protoc-gen-go-grpc

# Detect OS
ifeq ($(OS),Windows_NT)
    WHICH_CMD := where
    NULL_DEV := nul 2>&1
    PROTO_INSTALL_MSG := "protoc not installed. Please download and install from: https://github.com/protocolbuffers/protobuf/releases"
else
    WHICH_CMD := which
    NULL_DEV := /dev/null
    PROTO_INSTALL_MSG := "protoc not installed. Please install using your package manager: apt-get install protobuf-compiler or brew install protobuf"
endif

# Check if tools are installed
lint-check:
	@$(WHICH_CMD) $(GOLINT) >$(NULL_DEV) || (echo "$(GOLINT) not installed. Installing..." && \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

# windows need to install protoc-gen-go-grpc by yourself
proto-check:
	@$(WHICH_CMD) $(PROTOC) >$(NULL_DEV) || (echo $(PROTO_INSTALL_MSG))
	@$(WHICH_CMD) $(PROTOC_GEN_GO) >$(NULL_DEV) || (echo "Installing protoc-gen-go..." && \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)
	@$(WHICH_CMD) $(PROTOC_GEN_GRPC) >$(NULL_DEV) || (echo "Installing protoc-gen-go-grpc..." && \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest)

# Lint the Go code
lint: lint-check
	$(GOLINT) run ./...

# Generate code from protobuf files
proto: proto-check
	$(PROTOC) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# Clean generated files
clean:
ifeq ($(OS),Windows_NT)
	if exist $(BINARY_NAME) del $(BINARY_NAME)
	if exist $(PROTO_OUT_DIR)\*.pb.go del /Q $(PROTO_OUT_DIR)\*.pb.go
else
	rm -f $(BINARY_NAME)
	rm -f $(PROTO_OUT_DIR)/*.pb.go
endif

# Build the application
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/main.go

# Run the application
run:
	$(GO) run ./cmd/main.go

# Run tests
test:
	$(GOTEST) -v ./... 