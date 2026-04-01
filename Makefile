.PHONY: all clean generate build lint test

# Directories
PROTO_DIR := ./api/proto
PROTO_DEST := ./pkg/protocol

all: build

generate:
	@echo "Generating protobuf files..."
	@mkdir -p $(PROTO_DEST)
	@export PATH=$$PATH:$$(go env GOPATH)/bin && protoc --go_out=$(PROTO_DEST) --go_opt=module=im.turms/server/pkg/protocol -I=$(PROTO_DIR) $$(find $(PROTO_DIR) -name '*.proto')

build:
	@echo "Building gateway..."
	@go build -o bin/turms-gateway ./cmd/turms-gateway
	@echo "Building service..."
	@go build -o bin/turms-service ./cmd/turms-service

lint:
	golangci-lint run

test:
	go test -v -race ./...

clean:
	rm -rf bin/
