# Turms-Go Server

## Purpose
Go port of the turms Java IM (Instant Messaging) server. Uses MongoDB as primary storage, Redis for caching.

## Tech Stack
- Go 1.25
- MongoDB (go.mongodb.org/mongo-driver)
- Redis (github.com/redis/go-redis/v9)
- Protobuf (google.golang.org/protobuf)

## Code Structure
- `cmd/` - entrypoints
- `internal/domain/` - domain services (conversation, group, user, message, blocklist, etc.)
- `internal/infra/` - infrastructure (MongoDB repos, Redis)
- `internal/storage/` - storage layer
- `internal/pkg/` - shared utilities
- `pkg/` - public packages
- `turms-orig/` - original Java reference code
- `tests/` - integration tests

## Commands
- Build: `go build ./...`
- Test: `go test ./...`
- Format: `gofmt -w .`
- Lint: `.golangci.yml` configured

## Conventions
- Domain-driven design with service/repository layers
- Error codes defined in domain constants
- PO (Persistent Objects) for MongoDB documents
- DTO for request/response objects
