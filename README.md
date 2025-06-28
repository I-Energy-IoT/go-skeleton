# go-skeleton
Skeleton for golang project with tools for scaffolding code with full structure for I-Energy IoT project

## Features

- Create new Go service projects with Uber Fx
- Pre-configured project structure following Clean Architecture
- Built-in templates for common patterns and use cases
- Modern Go project layout with clear separation of concerns
- Integrated logging, configuration, and database support
- Ready-to-use HTTP server with middleware
- Swagger documentation setup
- JWT authentication and authorization
- Database migrations with PostgreSQL support
- Comprehensive testing setup with mocks

## Installation

```bash
go install github.com/I-Energy-IoT/go-skeleton@latest
```

## Project Structure

The generated service follows this structure:

```
myservice/
├── cmd/                    # Application entry points
│   └── app/               # Main application server
│        └── main.go        # Server entry point with dependency injection
├── config/                # Configuration management
│   └── config.go         # Environment and app configuration loader
├── external/             # External service integrations
│   └── [service]/        # Each external service in its own directory
├── internal/             # Private application code
│   ├── adapter/         # Interface adapters layer
│   │   └── http/        # HTTP delivery mechanism
│   │       ├── handler/ # HTTP request handlers
│   │       └── router/  # Route definitions
│   ├── app/             # Application layer
│   │   ├── constant/    # Application constants
│   │   ├── dto/         # Application DTOs
│   │   ├── middleware/  # Application middlewares
│   │   └── service/     # Application services
│   ├── domain/          # Domain models and business logic
│   │   ├── entity/      # Domain entities
│   │   └── enum/        # Domain enums and constants
│   └── infrastructure/  # Infrastructure implementations
│       ├── database/    # Database implementations
│       └── repository/  # Repository implementations
├── pkg/                # Public libraries
│   ├── constant/       # Shared constants
│   ├── errors/         # Custom error types and handling
│   ├── graceful/       # Graceful shutdown utilities
│   ├── logger/         # Logging configuration
│   ├── swagger/        # API documentation
│   ├── util/           # Common utilities
│   └── wrapper/        # Response wrappers
├── migrations/          # Database migration files
├── test/              # Test suites
│   └── integration/   # Integration tests
├── .github/           # GitHub configuration files
├── go.mod            # Go module definition
├── .env              # Environment variables
├── Makefile          # Build and development commands
└── README.md         # Project documentation
```

## Requirements

- Go 1.24 or higher
- PostgreSQL 15 or higher (configurable)
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
- To create mocks for services/repositories, define interfaces with `//go:generate mockgen` comments:
```bash
make generate-mock
```
- To create unit tests:
  1. Right-click on the interface
  2. Choose "Go: Generate Unit Tests For File"
  3. Complete your test cases

### Code Generation
- Generate mocks: `make generate-mock`
- Generate Swagger docs: `make swagger-build`

### Database Operations
- Create migration: `make migrate-create name=migration_name`
- Apply migrations: `make migrate-up`
- Revert migrations: `make migrate-down`
- Check migration version: `make migrate-version`
- Force migration: `make migrate-force version=1`

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make init` | Initialize project and install dependencies |
| `make install-tools` | Install required development tools |
| `make build` | Build the service |
| `make run` | Run the service |
| `make test` | Run tests |
| `make fmt` | Format code |
| `make lint` | Run linters |
| `make swagger-init` | Generate Swagger documentation |
| `make swagger-build` | Update Swagger documentation |
| `make generate-mock` | Generate mocks |
| `make migrate-create` | Create migration file |
| `make migrate-up` | Apply migration |
| `make migrate-down` | Revert migration |
| `make migrate-version` | Get migration version |
| `make migrate-force` | Force migrate |

## Architecture

This skeleton follows Clean Architecture principles with the following layers:

- **Domain Layer**: Contains business entities and core business logic
- **Application Layer**: Contains use cases and application services
- **Adapter Layer**: Contains HTTP handlers and external service adapters
- **Infrastructure Layer**: Contains database implementations and external integrations
