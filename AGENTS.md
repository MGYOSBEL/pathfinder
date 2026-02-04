# AGENTS.md - Developer Guide for Pathfinder

## Project Overview

Pathfinder is a Go-based data platform project focused on MQTT message processing and forwarding.
It integrates with message brokers (RabbitMQ), MQTT brokers (HiveMQ, VerneMQ), and time-series databases (TimescaleDB, InfluxDB).

## Project Structure

```
pathfinder/
├── cmd/                    # Application entry points
│   └── metrics-injector/   # Main service application
├── internal/               # Private application code
│   ├── config/            # Configuration handling
│   └── processor/         # Message processing logic
├── pkg/                   # Public libraries (importable by external projects)
│   ├── message/           # Message types and structures
│   ├── mqtt/              # MQTT client implementation
│   └── pubsub/            # Publisher/Subscriber interfaces
├── deploy/                # Deployment configurations
│   └── docker/           # Docker compose files
└── benthos/              # Benthos stream processing configs
```

## Build, Test & Lint Commands

### Building
```bash
# Build the main application
go build -o bin/metrics-injector ./cmd/metrics-injector

# Build all packages
go build ./...

# Cross-compile for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/metrics-injector-linux ./cmd/metrics-injector
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run a single test
go test -v -run TestName ./path/to/package

# Run tests for a specific package
go test -v ./pkg/mqtt
go test -v ./internal/processor

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
go test -bench=BenchmarkName -benchmem ./path/to/package
```

### Linting & Formatting
```bash
# Format code (always run before committing)
go fmt ./...
gofmt -s -w .

# Run go vet for static analysis
go vet ./...

# Install and run golangci-lint (recommended)
# Install: brew install golangci-lint or go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
golangci-lint run --fix

# Check imports
goimports -w .

# Tidy dependencies
go mod tidy
go mod verify
```

### Docker & Deployment
```bash
# Deploy with default services (HiveMQ, RabbitMQ, TimescaleDB)
make deploy

# Deploy with custom services
MQTT_BROKER=vernemq TIMESERIES_DB=influxdb make deploy

# Stop containers
make stop

# Remove containers
make down

# Remove containers and volumes
make destroy

# View composed configuration
make config
```

## Code Style Guidelines

### Package Organization
- `cmd/`: Executable entry points, minimal logic, only main packages
- `internal/`: Private application code, not importable by external projects
- `pkg/`: Public libraries that external projects can import
- Each package should have a clear, single responsibility

### Imports
- Group imports in this order: standard library, external packages, internal packages
- Use blank line to separate groups
- Example:
```go
import (
    "context"
    "fmt"
    "os"

    "github.com/eclipse/paho.mqtt.golang"

    "github.com/MGYOSBEL/pathfinder/internal/processor"
    "github.com/MGYOSBEL/pathfinder/pkg/mqtt"
)
```
- Avoid dot imports (import . "package")
- Use descriptive aliases when needed: `MQTT "github.com/eclipse/paho.mqtt.golang"`

### Naming Conventions
- **Packages**: lowercase, single word, no underscores (e.g., `mqtt`, `pubsub`, `processor`)
- **Types**: PascalCase (e.g., `MqttClient`, `Options`, `Processor`)
- **Interfaces**: PascalCase, often noun or noun+er (e.g., `Publisher`, `Subscriber`, `Handler`)
- **Functions/Methods**: camelCase (private) or PascalCase (exported)
- **Variables**: camelCase (e.g., `client`, `inputTopic`)
- **Constants**: PascalCase or UPPER_CASE for exported constants
- **Acronyms**: Keep consistent (MQTT, not Mqtt; QoS, not Qos)

### Types & Structs
- Define Options structs for complex configuration:
```go
type Options struct {
    Server string
    Topic  string
    QoS    byte
}
```
- Use embedded types for composition (e.g., `Metadata` embedded in `Message`)
- Export fields that need to be accessed outside the package
- Constructor pattern: `NewClient(opts Options) *MqttClient`

### Interfaces
- Keep interfaces small and focused (1-3 methods is ideal)
- Define interfaces where they are used, not where they are implemented
- Example:
```go
type Publisher interface {
    Publish(topic string, m message.Message) error
}
```

### Error Handling
- Always check and handle errors explicitly
- Return errors as the last return value
- Use `fmt.Errorf()` for wrapping errors with context
- Don't panic in library code; reserve for truly unrecoverable errors
- In main/cmd, it's acceptable to panic on setup failures
- Pattern:
```go
if err != nil {
    return fmt.Errorf("failed to connect: %w", err)
}
```

### Context Usage
- Pass `context.Context` as the first parameter in functions that perform I/O or long operations
- Use `context.Background()` as top-level context in main
- Use `signal.NotifyContext()` for graceful shutdown
- Example:
```go
ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
defer cancel()
```

### Concurrency
- Use goroutines for async operations
- Always ensure goroutines can exit cleanly (context cancellation, channels)
- Use `defer` for cleanup (e.g., `defer client.Disconnect()`)
- Protect shared state with mutexes or use channels for communication

### Comments & Documentation
- Add package-level comments describing the package purpose
- Document exported types, functions, and methods with godoc-style comments
- Start comments with the name being documented: `// NewClient creates a new MQTT client`
- Keep TODOs in comments for tracking future work
- Avoid obvious comments; comment "why" not "what"

### Testing (When Adding Tests)
- Test files: `*_test.go` in the same package
- Test function names: `TestFunctionName(t *testing.T)`
- Table-driven tests for multiple scenarios
- Use `t.Helper()` in test helper functions
- Use subtests with `t.Run()` for grouping related tests

## Common Patterns

### Client/Service Pattern
```go
type Client struct {
    options Options
    client  ThirdPartyClient
}

func NewClient(opts Options) *Client {
    return &Client{options: opts}
}

func (c *Client) Connect() error { /* ... */ }
func (c *Client) Disconnect() { /* ... */ }
```

### Handler/Callback Pattern
```go
type Handler func(msg message.Message)

func (c *Client) Subscribe(h Handler) error { /* ... */ }
```

## Environment & Dependencies

- **Go Version**: 1.23.1+
- **Module Path**: `github.com/MGYOSBEL/pathfinder`
- **Key Dependencies**:
  - `github.com/eclipse/paho.mqtt.golang` - MQTT client
- **Docker**: Required for running the full stack

## Git Workflow

- Keep commits atomic and focused
- Write clear commit messages
- Run `go fmt ./...` and `go vet ./...` before committing
- Ensure `go mod tidy` is run to keep go.mod/go.sum clean

## Notes for AI Agents

1. **No tests exist yet** - when adding functionality, tests are welcomed but not required unless specified
2. **Focus on simplicity** - this is a data platform project; prefer clear, maintainable code over clever solutions
3. **Interface-based design** - the codebase uses small interfaces (Publisher, Subscriber) for loose coupling
4. **Context propagation** - always propagate context for cancellation and timeouts
5. **Docker-first deployment** - changes should be compatible with Docker Compose deployment model
6. **Configuration** - prefer Options structs over many function parameters
