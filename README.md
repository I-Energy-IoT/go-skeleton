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
├── cmd/
│   ├── app/
│   |   └── main.go      # Application entry point
|   ├── config/          # Configuration files
├── external/        # External service integrations
├── internal/
│   ├── app/            # Application bootstrapping
│   ├── http/           # HTTP layer
|   |    ├── handler/   # HTTP handlers
│   |    ├── dto/       # Data transfer object
|   |    └── router/    # Route definitions
|   ├── middleware/     # HTTP middleware
│   ├── model/         # Domain models
│   ├── repository/    # Data access layer
│   └── service/       # Business logic layer
├── pkg/
│   ├── constant/      # Shared constants
│   ├── database/      # Database utilities
│   ├── errors/        # Custom error types
|   ├── graceful/      # Graceful shutdown
|   ├── logger/        # Logging utilities
|   ├── swagger/       # Swagger handler
|   ├── util/          # Common utilities
│   └── wrapper/       # Response wrapper
├── test/
│   └── integration/   # Integration tests
├── go.mod            # Go module file
├── .env              # Environment variables
├── .mockery.yaml          # Build automation
├── Makefile          # Build automation
└── README.md
```

## Requirements

- Go 1.21 or higher
- Git

## Usage

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
- Update values in `.env` file

7. Build the service:
```bash
make build
```

8. Run the service:
```bash
make run
```

## Testing
- To create the mock service/repo define in .mockery.yaml and run
``` bash
make generate-mock
```
- To create the unit test using go tool
1. Click right mouse
2. Choose Go:Generate Unit Tests For Function
3. Complete your testcase

## Available Make Commands

- `make init` - Initialize project and install dependencies
- `make install-tools` - Install required development tools
- `make build` - Build the service
- `make run` - Run the service
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Run linters
- `make swagger-init` - Generate Swagger documentation
- `make swagger-build` - Update Swagger documentation
- `make generata-mock` - Generate mocks
- `make test` - Run test
