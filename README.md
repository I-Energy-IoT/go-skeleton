# Go Skeleton - Clean Architecture REST API Generator

A powerful CLI tool for generating production-ready Go REST API projects with Clean Architecture, dependency injection, and comprehensive tooling for I-Energy IoT projects.

## 🚀 Features

### Core Architecture
- **Clean Architecture** implementation with clear separation of concerns
- **Uber FX** dependency injection framework for robust service management
- **Gin** web framework with middleware support
- **Structured logging** with Zap logger
- **Configuration management** with environment-based settings
- **Graceful shutdown** handling

### Authentication & Security
- **JWT authentication** with token generation and validation
- **Role-based access control** (RBAC) middleware
- **CORS** configuration
- **Request validation** with go-playground/validator
- **Error handling** with custom error types

### Database & Data Layer
- **PostgreSQL** support with GORM ORM
- **Database migrations** with golang-migrate
- **Repository pattern** implementation
- **Soft delete** support
- **Connection pooling** configuration

### API Documentation
- **Swagger/OpenAPI** documentation generation
- **Auto-generated API docs** with examples
- **Interactive API explorer** at `/swagger`

### Development Tools
- **Mock generation** with Mockery
- **Code formatting** with goimports
- **Linting** with golangci-lint
- **Testing** framework with coverage reporting
- **Make commands** for common development tasks

### Deployment & DevOps
- **Kubernetes** deployment templates
- **Docker** support
- **Health checks** and readiness probes
- **Auto-scaling** configuration
- **Prometheus** metrics support
- **Environment-specific** configurations

### Testing & Quality
- **Unit testing** framework
- **Integration testing** setup
- **Mock generation** for testing
- **Code coverage** reporting
- **Git conventions** and workflow guidelines

## 📦 Installation

```bash
go install github.com/I-Energy-IoT/go-skeleton@latest
```

## 🏗️ Project Structure

The generated service follows Clean Architecture principles:

```
myservice/
├── cmd/                    # Application entry points
│   └── app/               # Main application server
│        └── main.go        # Server entry point with FX DI
├── config/                # Configuration management
│   └── config.go         # Environment and app configuration
├── deployment/            # Deployment configurations
│   ├── prd.yaml         # Production Kubernetes config
│   └── stg.yaml         # Staging Kubernetes config
├── docs/                 # Documentation
│   ├── GIT_CONVENTIONS.md # Git workflow guidelines
│   └── TESTING_GUIDELINES.md # Testing best practices
├── internal/             # Private application code
│   ├── adapter/         # Interface adapters layer
│   │   └── http/        # HTTP delivery mechanism
│   │       ├── dto/     # Data Transfer Objects
│   │       ├── handler/ # HTTP request handlers
│   │       └── router/  # Route definitions
│   ├── app/             # Application layer
│   │   ├── middleware/  # Application middlewares
│   │   │   ├── cors.go           # CORS configuration
│   │   │   ├── error_handler.go  # Error handling middleware
│   │   │   ├── jwt_auth.go       # JWT authentication
│   │   │   ├── logging.go        # Request logging
│   │   │   └── recover.go        # Panic recovery
│   │   ├── service/     # Application services
│   │   └── validator/   # Request validation
│   ├── domain/          # Domain models and business logic
│   │   ├── entity/      # Domain entities
│   │   │   └── base_entity.go # Base entity with common fields
│   │   ├── enum/        # Domain enums and constants
│   │   ├── repository/  # Repository interfaces
│   │   └── service/     # Domain services
│   └── infrastructure/  # Infrastructure implementations
│       ├── database/    # Database implementations
│       │   └── postgre.go # PostgreSQL connection
│       ├── external/    # External service integrations
│       └── repository/  # Repository implementations
├── pkg/                # Public libraries
│   ├── authentication/ # JWT authentication utilities
│   │   └── jwt.go     # Token generation and validation
│   ├── errors/         # Custom error types and handling
│   ├── graceful/       # Graceful shutdown utilities
│   ├── logger/         # Structured logging with Zap
│   ├── swagger/        # API documentation setup
│   ├── util/           # Common utilities
│   │   └── context.go  # Context utilities
│   └── wrapper/        # Response wrappers
├── test/              # Test suites
│   └── integration/   # Integration tests
├── migrations/         # Database migration files
├── .github/           # GitHub configuration files
├── go.mod            # Go module definition
├── .env              # Environment variables
├── Makefile          # Build and development commands
└── README.md         # Project documentation
```

## 🛠️ Requirements

- **Go 1.24** or higher
- **PostgreSQL 15** or higher (configurable)
- **Git**
- **Docker** (for containerization)
- **Kubernetes** (for deployment)

## 🚀 Getting Started

### 1. Create a New Service

```bash
go-skeleton new --name yourservice
```

### 2. Navigate to Your Service

```bash
cd yourservice
```

### 3. Install Development Tools

```bash
make install-tools
```

This installs:
- `swag` - Swagger documentation generator
- `golangci-lint` - Code linter
- `goimports` - Code formatter
- `mockery` - Mock generator
- `migrate` - Database migration tool

### 4. Initialize the Project

```bash
make init
```

### 5. Configure Environment

Copy and configure your environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 6. Build and Run

```bash
make build
make run
```

## 🔧 Development

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make install-tools` | Install required development tools |
| `make init` | Initialize project and install dependencies |
| `make build` | Build the service |
| `make run` | Run the service |
| `make test` | Run tests with coverage |
| `make fmt` | Format code with goimports |
| `make lint` | Run linters with golangci-lint |
| `make swagger-init` | Initialize Swagger documentation |
| `make swagger-build` | Update Swagger documentation |
| `make generate-mock` | Generate mocks for interfaces |

### Database Operations

| Command | Description |
|---------|-------------|
| `make migrate-create NAME=migration_name` | Create new migration |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-down-n N=3` | Rollback N migrations |
| `make migrate-force VERSION=1` | Force migration version |
| `make migrate-drop` | Drop all migrations |
| `make migrate-goto VERSION=1` | Apply specific migration |

### Testing

#### Unit Tests
1. Right-click on the interface
2. Choose "Go: Generate Unit Tests For File"
3. Complete your test cases

#### Integration Tests
```bash
make test
```

#### Generate Mocks
```bash
make generate-mock
```

### Code Generation

#### Swagger Documentation
```bash
make swagger-build
```

#### Mock Generation
```bash
make generate-mock
```

## 🏗️ Architecture

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Business entities and core business logic
   - Repository interfaces
   - Domain services

2. **Application Layer** (`internal/app/`)
   - Use cases and application services
   - Middleware implementations
   - Request validation

3. **Adapter Layer** (`internal/adapter/`)
   - HTTP handlers and external service adapters
   - DTOs for data transfer
   - Route definitions

4. **Infrastructure Layer** (`internal/infrastructure/`)
   - Database implementations
   - External service integrations
   - Repository implementations

### Dependency Injection

The project uses **Uber FX** for dependency injection:

```go
app := fx.New(
    config.Module,
    logger.Module,
    internal.Module,
    fx.Invoke(startServer),
)
```

### Configuration Management

Environment-based configuration with support for:
- Development, Staging, UAT, Production environments
- Database configuration
- JWT settings
- Server configuration
- Logging levels

### Authentication & Authorization

- **JWT Token Generation**: Secure token creation with expiration
- **Token Validation**: Middleware for request authentication
- **Role-Based Access**: RBAC middleware for authorization
- **Context Extraction**: Utilities for extracting user info

### Database Layer

- **GORM ORM**: Object-relational mapping
- **PostgreSQL**: Primary database support
- **Migrations**: Version-controlled schema changes
- **Repository Pattern**: Clean data access abstraction
- **Connection Pooling**: Optimized database connections

## 🚀 Deployment

### Kubernetes Deployment

The project includes Kubernetes deployment templates:

- **Production** (`deployment/prd.yaml`): High-availability configuration
- **Staging** (`deployment/stg.yaml`): Testing environment setup

### Docker Support

```bash
# Build Docker image
docker build -t yourservice .

# Run with Docker
docker run -p 8080:8080 yourservice
```

### Environment Configuration

Configure your deployment environment:

```yaml
# Production environment variables
app:
  env:
    - name: APP_ENV
      value: "production"
    - name: DB_HOST
      value: "postgres-service"
    - name: JWT_SECRET
      valueFrom:
        secretKeyRef:
          name: jwt-secret
          key: secret
```

## 📚 API Documentation

### Swagger UI

Access interactive API documentation at:
- **Development**: `http://localhost:8080/swagger/`
- **Production**: `https://yourdomain.com/swagger/`

### Generate Documentation

```bash
make swagger-build
```

## 🧪 Testing

### Test Structure

```
test/
├── integration/     # Integration tests
└── unit/           # Unit tests (generated)
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -v -race -coverprofile=coverage.txt ./...

# View coverage report
go tool cover -html=coverage.txt
```

## 📝 Git Workflow

The project includes comprehensive Git conventions:

### Commit Message Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Commit Types
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions
- `chore`: Maintenance tasks

### Branch Strategy
- `main`: Production-ready code
- `develop`: Development branch
- `feature/*`: Feature branches
- `hotfix/*`: Hotfix branches

## 🔧 Configuration

### Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRATION_TIME=24h
JWT_REFRESH_EXPIRATION_TIME=168h

# Server Configuration
SERVER_URL=localhost
PORT=8080
ENVIRONMENT=development
GIN_MODE=debug
PRODUCTION=false
```

### Configuration Features

- **Environment-based**: Different configs for dev/staging/prod
- **Validation**: Required field validation
- **Defaults**: Sensible default values
- **Hot-reload**: Configuration reloading support

## 🛡️ Security Features

### Authentication
- JWT token-based authentication
- Token expiration and refresh
- Secure token validation
- Context-based user extraction

### Authorization
- Role-based access control
- Middleware for permission checking
- Flexible role system

### Input Validation
- Request validation with go-playground/validator
- Custom validation rules
- Error handling for invalid inputs

### Security Headers
- CORS configuration
- Security middleware
- Panic recovery

## 📊 Monitoring & Observability

### Logging
- Structured logging with Zap
- Request/response logging
- Error tracking
- Performance metrics

### Health Checks
- `/health` - Liveness probe
- `/ready` - Readiness probe
- Custom health check endpoints

### Metrics
- Prometheus metrics support
- Custom business metrics
- Performance monitoring

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow the established Git conventions
- Write comprehensive tests
- Update documentation
- Ensure code passes linting
- Follow Clean Architecture principles

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation in `/docs`
- Review the testing guidelines
- Follow the Git conventions

---

**Built with ❤️ for I-Energy IoT projects**
