# go-skeleton
Skeleton for golang project include tools for scafoding code with full strucure for I-energy iot project

## Features

- Create new Go service projects with Uber Fx
- Pre-configured project structure following Clean Architecture
- Built-in templates for common patterns and use cases
- Modern Go project layout with clear separation of concerns
- Integrated logging, configuration, and database support
- Ready-to-use HTTP server with middleware
- Swagger documentation setup

## Installation

```bash
go install github.com/I-Energy-IoT/go-skeleton@latest
```

## Project Structure

The generated service follows this structure:

```
myservice/
├── cmd/                    # Application entry points
│   ├── app/               # Main application server
│   |   └── main.go        # Server entry point with dependency injection
│   └── migrate/           # Database migration tool
│       └── main.go        # Migration entry point
├── config/                # Configuration management
│   └── config.go         # Environment and app configuration loader
├── external/             # External service integrations
│   └── [service]/        # Each external service in its own directory
├── internal/             # Private application code
│   ├── app/             # Application setup and dependency injection
│   ├── delivery/        # Interface adapters layer
│   │   └── http/        # HTTP delivery mechanism
│   │       ├── handler/ # HTTP request handlers
│   │       ├── dto/     # Data transfer objects
│   │       └── router/  # Route definitions and middleware
│   ├── middleware/      # HTTP middleware components
│   ├── model/          # Domain models and entities
│   ├── repository/     # Data access layer
│   └── service/        # Business logic layer
├── pkg/                # Public libraries
│   ├── constant/       # Shared constants
│   ├── database/       # Database connection and utilities
│   ├── errors/         # Custom error types and handling
│   ├── graceful/       # Graceful shutdown utilities
│   ├── logger/         # Logging configuration
│   ├── swagger/        # API documentation
│   ├── util/           # Common utilities
│   └── wrapper/        # Response wrappers
├── test/               # Test suites
│   └── integration/    # Integration tests
├── go.mod             # Go module definition
├── .env              # Environment variables
├── Makefile          # Build and development commands
└── README.md         # Project documentation
```

## Requirements

- Go 1.21 or higher
- PostgreSQL 15 or higher (Can change)
- Git

## Getting Started

1. Create a new service:
```bash
go-skeleton new --name yourservice
```

2. Navigate to your service directory:
```bash
cd yourservice
```

3. Install required tools:
```bash
make install-tools
```

4. Initialize the project:
```bash
make init
```

5. Format code:
```bash
make fmt
```

6. Configure your environment:
- Copy `.env.example` to `.env`
- Update values in `.env` file

7. Build the service:
```bash
make build
```

8. Run the service:
```bash
make run
```

## Development

### Testing
- To create the mock service/repo define with syntax mock at interfaces:
```bash
make generate-mock
```
- To create unit tests:
  1. Click right mouse
  2. Choose Go:Generate Unit Tests For Function
  3. Complete your testcase

### Code Generation
- Generate mocks: `make generate-mock`
- Generate Swagger docs: `make swagger-build`

### Available Make Commands

- `make init` - Initialize project and install dependencies
- `make install-tools` - Install required development tools
- `make build` - Build the service
- `make run` - Run the service
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Run linters
- `make swagger-init` - Generate Swagger documentation
- `make swagger-build` - Update Swagger documentation
- `make generate-mock` - Generate mocks
- `make migrate-create` - Create migration file
- `make migrate-up` - Apply migration
- `make migrate-down` - Revert migration
- `make migrate-version` - Get version migration
- `make migrate-force` - Force migrate
