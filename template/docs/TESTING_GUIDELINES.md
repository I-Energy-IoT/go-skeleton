# Testing Guidelines

## 📋 Overview

This document outlines the testing strategy and guidelines for the IoT Admin Service API project. Following these guidelines ensures code quality, reliability, and maintainability through comprehensive testing practices.

## 🎯 Testing Strategy

### Testing Pyramid

```
                    /\
                   /  \
                  / E2E \
                 /________\
                /          \
               / Integration \
              /______________\
             /                \
            /     Unit Tests    \
           /____________________\
```

- **Unit Tests**: 70% - Fast, isolated, test individual functions
- **Integration Tests**: 20% - Test component interactions
- **End-to-End Tests**: 10% - Test complete workflows

### Test Types Overview

| Test Type | Purpose | Scope | Speed | Tools |
|-----------|---------|-------|-------|-------|
| **Unit Tests** | Test individual functions/methods | Single function/class | Fast | `testing`, `testify` |
| **Integration Tests** | Test component interactions | Multiple components | Medium | `testcontainers`, `httptest` |
| **API Tests** | Test HTTP endpoints | API layer | Medium | `httptest`, Postman |
| **Database Tests** | Test data persistence | Database operations | Slow | `testcontainers` |
| **E2E Tests** | Test complete workflows | Full application | Slow | `httptest`, external tools |

## 🧪 Unit Testing

### Principles

- **Isolation**: Test one function at a time
- **Fast**: Should run in milliseconds
- **Deterministic**: Same result every time
- **No external dependencies**: Mock external calls

### Table-Driven Unit Tests
```go
// internal/application/service/user_service_test.go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/google/uuid"
)

func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name           string
        request        *dto.CreateUserRequest
        setupMocks     func(*MockUserRepository, *MockCustomerRepository)
        expectedResult *dto.UserResponse
        expectedError  error
    }{
        {
            name: "successful user creation",
            request: &dto.CreateUserRequest{
                FirstName: "John",
                LastName:  "Doe",
                Email:     "john.doe@example.com",
                Password:  "securepassword123",
            },
            setupMocks: func(userRepo *MockUserRepository, customerRepo *MockCustomerRepository) {
                expectedUser := &entity.User{
                    ID:        uuid.New(),
                    FirstName: "John",
                    LastName:  "Doe",
                    Email:     "john.doe@example.com",
                }
                userRepo.On("Create", mock.Anything, mock.Anything).Return(expectedUser, nil)
                customerRepo.On("GetByID", mock.Anything, mock.Anything).Return(&entity.Customer{}, nil)
            },
            expectedResult: &dto.UserResponse{
                FirstName: "John",
                LastName:  "Doe",
                Email:     "john.doe@example.com",
            },
            expectedError: nil,
        },
        {
            name: "email already exists",
            request: &dto.CreateUserRequest{
                FirstName: "Jane",
                LastName:  "Smith",
                Email:     "existing@example.com",
                Password:  "password123",
            },
            setupMocks: func(userRepo *MockUserRepository, customerRepo *MockCustomerRepository) {
                existingUser := &entity.User{
                    ID:    uuid.New(),
                    Email: "existing@example.com",
                }
                userRepo.On("GetByEmailWithCustomer", mock.Anything, "existing@example.com", mock.Anything).Return(existingUser, nil)
            },
            expectedResult: nil,
            expectedError:  customErrors.NewConflictError("user with this email already exists"),
        },
        {
            name: "invalid email format",
            request: &dto.CreateUserRequest{
                FirstName: "Invalid",
                LastName:  "User",
                Email:     "invalid-email",
                Password:  "password123",
            },
            setupMocks: func(userRepo *MockUserRepository, customerRepo *MockCustomerRepository) {
                // No mocks needed for validation error
            },
            expectedResult: nil,
            expectedError:  customErrors.NewBadRequestError("invalid email format"),
        },
        {
            name: "customer not found",
            request: &dto.CreateUserRequest{
                FirstName:  "Customer",
                LastName:   "User",
                Email:     "customer@example.com",
                Password:  "password123",
                CustomerID: &uuid.UUID{},
            },
            setupMocks: func(userRepo *MockUserRepository, customerRepo *MockCustomerRepository) {
                customerRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
            },
            expectedResult: nil,
            expectedError:  customErrors.NewNotFoundError("customer not found"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockUserRepo := &MockUserRepository{}
            mockCustomerRepo := &MockCustomerRepository{}
            service := NewUserService(mockUserRepo, mockCustomerRepo, ...)
            
            if tt.setupMocks != nil {
                tt.setupMocks(mockUserRepo, mockCustomerRepo)
            }
            
            // Act
            result, err := service.CreateUser(context.Background(), tt.request)
            
            // Assert
            if tt.expectedError != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedError.Error(), err.Error())
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.expectedResult.FirstName, result.FirstName)
                assert.Equal(t, tt.expectedResult.LastName, result.LastName)
                assert.Equal(t, tt.expectedResult.Email, result.Email)
            }
            
            mockUserRepo.AssertExpectations(t)
            mockCustomerRepo.AssertExpectations(t)
        })
    }
}
```

### Test Naming Convention

```
Test<StructName>_<MethodName>
```

**Examples:**
- `TestUserService_CreateUser`
- `TestUserService_UpdateUser`
- `TestUserService_GetUserByID`

### Mocking Guidelines

```go
// internal/domain/repository/mocks/user_repository.go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user entity.User) (*entity.User, error) {
    args := m.Called(ctx, user)
    return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID, isSystemAdmin bool) (*entity.User, error) {
    args := m.Called(ctx, id, isSystemAdmin)
    return args.Get(0).(*entity.User), args.Error(1)
}
```

### Test Data Management

```go
// test/fixtures/user_fixtures.go
package fixtures

import (
    "admin-api/internal/domain/entity"
    "github.com/google/uuid"
)

func CreateTestUser() *entity.User {
    return &entity.User{
        ID:        uuid.New(),
        FirstName: "Test",
        LastName:  "User",
        Email:     "test@example.com",
        Password:  "hashedpassword",
    }
}

func CreateTestUserWithCustomer(customerID uuid.UUID) *entity.User {
    user := CreateTestUser()
    user.CustomerID = &customerID
    return user
}
```

## 🔗 Integration Testing

### Table-Driven Integration Tests

```go
// test/integration/user_integration_test.go
package integration

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

type UserIntegrationTestSuite struct {
    suite.Suite
    container testcontainers.Container
    db        *gorm.DB
    userRepo  repository.UserRepository
}

func (suite *UserIntegrationTestSuite) SetupSuite() {
    // Start PostgreSQL container
    req := testcontainers.ContainerRequest{
        Image:        "postgres:14",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       "test_db",
            "POSTGRES_USER":     "test_user",
            "POSTGRES_PASSWORD": "test_password",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }
    
    container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    suite.Require().NoError(err)
    suite.container = container
    
    // Setup database connection
    host, _ := container.Host(context.Background())
    port, _ := container.MappedPort(context.Background(), "5432")
    
    dsn := fmt.Sprintf("host=%s port=%s user=test_user password=test_password dbname=test_db sslmode=disable", host, port.Port())
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    suite.Require().NoError(err)
    
    // Run migrations
    err = db.AutoMigrate(&entity.User{}, &entity.Customer{}, &entity.Role{})
    suite.Require().NoError(err)
    
    suite.db = db
    suite.userRepo = repository.NewUserRepository(db)
}

func (suite *UserIntegrationTestSuite) TearDownSuite() {
    suite.container.Terminate(context.Background())
}

func (suite *UserIntegrationTestSuite) TestCreateUser_TableDriven() {
    tests := []struct {
        name           string
        user           *entity.User
        expectedError  error
        expectedFields map[string]interface{}
    }{
        {
            name: "valid user creation",
            user: &entity.User{
                FirstName: "John",
                LastName:  "Doe",
                Email:     "john.doe@example.com",
                Password:  "hashedpassword",
            },
            expectedError: nil,
            expectedFields: map[string]interface{}{
                "FirstName": "John",
                "LastName":  "Doe",
                "Email":     "john.doe@example.com",
            },
        },
        {
            name: "duplicate email",
            user: &entity.User{
                FirstName: "Jane",
                LastName:  "Smith",
                Email:     "john.doe@example.com", // Same email as above
                Password:  "hashedpassword",
            },
            expectedError: gorm.ErrDuplicatedKey,
            expectedFields: nil,
        },
        {
            name: "empty first name",
            user: &entity.User{
                FirstName: "",
                LastName:  "Doe",
                Email:     "empty@example.com",
                Password:  "hashedpassword",
            },
            expectedError: gorm.ErrInvalidData,
            expectedFields: nil,
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Act
            createdUser, err := suite.userRepo.Create(context.Background(), *tt.user)
            
            // Assert
            if tt.expectedError != nil {
                suite.Error(err)
                suite.Nil(createdUser)
            } else {
                suite.NoError(err)
                suite.NotNil(createdUser)
                
                // Verify expected fields
                for field, expectedValue := range tt.expectedFields {
                    suite.Equal(expectedValue, getFieldValue(createdUser, field))
                }
                
                // Verify in database
                var foundUser entity.User
                err = suite.db.First(&foundUser, createdUser.ID).Error
                suite.NoError(err)
                suite.Equal(tt.user.FirstName, foundUser.FirstName)
                suite.Equal(tt.user.Email, foundUser.Email)
            }
        })
    }
}

// Helper function to get field value using reflection
func getFieldValue(user *entity.User, fieldName string) interface{} {
    v := reflect.ValueOf(user).Elem()
    return v.FieldByName(fieldName).Interface()
}
```

### Database Integration Tests

```go
// test/integration/user_integration_test.go
package integration

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

type UserIntegrationTestSuite struct {
    suite.Suite
    container testcontainers.Container
    db        *gorm.DB
    userRepo  repository.UserRepository
}

func (suite *UserIntegrationTestSuite) SetupSuite() {
    // Start PostgreSQL container
    req := testcontainers.ContainerRequest{
        Image:        "postgres:14",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       "test_db",
            "POSTGRES_USER":     "test_user",
            "POSTGRES_PASSWORD": "test_password",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }
    
    container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    suite.Require().NoError(err)
    suite.container = container
    
    // Setup database connection
    host, _ := container.Host(context.Background())
    port, _ := container.MappedPort(context.Background(), "5432")
    
    dsn := fmt.Sprintf("host=%s port=%s user=test_user password=test_password dbname=test_db sslmode=disable", host, port.Port())
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    suite.Require().NoError(err)
    
    // Run migrations
    err = db.AutoMigrate(&entity.User{}, &entity.Customer{}, &entity.Role{})
    suite.Require().NoError(err)
    
    suite.db = db
    suite.userRepo = repository.NewUserRepository(db)
}

func (suite *UserIntegrationTestSuite) TearDownSuite() {
    suite.container.Terminate(context.Background())
}

func (suite *UserIntegrationTestSuite) TestCreateUser_Success() {
    // Arrange
    user := fixtures.CreateTestUser()
    
    // Act
    createdUser, err := suite.userRepo.Create(context.Background(), *user)
    
    // Assert
    suite.NoError(err)
    suite.NotNil(createdUser)
    suite.Equal(user.FirstName, createdUser.FirstName)
    suite.Equal(user.Email, createdUser.Email)
    
    // Verify in database
    var foundUser entity.User
    err = suite.db.First(&foundUser, createdUser.ID).Error
    suite.NoError(err)
    suite.Equal(user.FirstName, foundUser.FirstName)
}
```

### HTTP Integration Tests

```go
// test/integration/api_integration_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/suite"
)

type APIIntegrationTestSuite struct {
    suite.Suite
    router *gin.Engine
    db     *gorm.DB
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
    // Setup test database
    suite.setupTestDatabase()
    
    // Setup router with test dependencies
    suite.router = setupTestRouter(suite.db)
}

func (suite *APIIntegrationTestSuite) TestCreateUserAPI_Success() {
    // Arrange
    reqBody := dto.CreateUserRequest{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john.doe@example.com",
        Password:  "securepassword123",
    }
    
    jsonBody, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/api/v1/admin/users", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+getTestToken())
    
    w := httptest.NewRecorder()
    
    // Act
    suite.router.ServeHTTP(w, req)
    
    // Assert
    suite.Equal(http.StatusCreated, w.Code)
    
    var response dto.UserResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    suite.NoError(err)
    suite.Equal(reqBody.FirstName, response.FirstName)
    suite.Equal(reqBody.Email, response.Email)
}
```

## 🌐 API Testing

### Table-Driven API Testing

```go
// test/api/user_api_test.go
package api

import (
    "testing"
    "net/http"
    "encoding/json"
    "bytes"
    "github.com/stretchr/testify/suite"
)

type UserAPITestSuite struct {
    suite.Suite
    client *http.Client
    baseURL string
}

func (suite *UserAPITestSuite) TestCreateUser_TableDriven() {
    tests := []struct {
        name           string
        requestBody    *dto.CreateUserRequest
        headers        map[string]string
        expectedStatus int
        expectedError  string
        expectedFields map[string]interface{}
    }{
        {
            name: "valid user creation",
            requestBody: &dto.CreateUserRequest{
                FirstName: "John",
                LastName:  "Doe",
                Email:     "john.doe@example.com",
                Password:  "securepassword123",
            },
            headers: map[string]string{
                "Authorization": "Bearer " + getTestToken(),
                "Content-Type":  "application/json",
            },
            expectedStatus: http.StatusCreated,
            expectedError:  "",
            expectedFields: map[string]interface{}{
                "firstName": "John",
                "lastName":  "Doe",
                "email":     "john.doe@example.com",
            },
        },
        {
            name: "missing authorization",
            requestBody: &dto.CreateUserRequest{
                FirstName: "John",
                LastName:  "Doe",
                Email:     "john.doe@example.com",
                Password:  "securepassword123",
            },
            headers: map[string]string{
                "Content-Type": "application/json",
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "unauthorized",
            expectedFields: nil,
        },
        {
            name: "invalid email format",
            requestBody: &dto.CreateUserRequest{
                FirstName: "Invalid",
                LastName:  "User",
                Email:     "invalid-email",
                Password:  "password123",
            },
            headers: map[string]string{
                "Authorization": "Bearer " + getTestToken(),
                "Content-Type":  "application/json",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "validation failed",
            expectedFields: nil,
        },
        {
            name: "empty required fields",
            requestBody: &dto.CreateUserRequest{
                FirstName: "",
                LastName:  "",
                Email:     "",
                Password:  "",
            },
            headers: map[string]string{
                "Authorization": "Bearer " + getTestToken(),
                "Content-Type":  "application/json",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "validation failed",
            expectedFields: nil,
        },
        {
            name: "password too short",
            requestBody: &dto.CreateUserRequest{
                FirstName: "Short",
                LastName:  "Password",
                Email:     "short@example.com",
                Password:  "123", // Too short
            },
            headers: map[string]string{
                "Authorization": "Bearer " + getTestToken(),
                "Content-Type":  "application/json",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "password too short",
            expectedFields: nil,
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Arrange
            jsonBody, err := json.Marshal(tt.requestBody)
            suite.NoError(err)
            
            req, err := http.NewRequest("POST", suite.baseURL+"/api/v1/admin/users", bytes.NewBuffer(jsonBody))
            suite.NoError(err)
            
            // Set headers
            for key, value := range tt.headers {
                req.Header.Set(key, value)
            }
            
            // Act
            resp, err := suite.client.Do(req)
            suite.NoError(err)
            defer resp.Body.Close()
            
            // Assert
            suite.Equal(tt.expectedStatus, resp.StatusCode)
            
            if tt.expectedError != "" {
                var errorResponse wrapper.Response
                err = json.NewDecoder(resp.Body).Decode(&errorResponse)
                suite.NoError(err)
                suite.Contains(errorResponse.Message, tt.expectedError)
            } else {
                var response dto.UserResponse
                err = json.NewDecoder(resp.Body).Decode(&response)
                suite.NoError(err)
                
                // Verify expected fields
                for field, expectedValue := range tt.expectedFields {
                    suite.Equal(expectedValue, getResponseField(&response, field))
                }
            }
        })
    }
}

func (suite *UserAPITestSuite) TestGetUsers_TableDriven() {
    tests := []struct {
        name           string
        queryParams    map[string]string
        headers        map[string]string
        expectedStatus int
        expectedCount  int
        expectedTotal  int64
    }{
        {
            name: "get users with pagination",
            queryParams: map[string]string{
                "page":     "0",
                "pageSize": "10",
            },
            headers: map[string]string{
                "Authorization": "Bearer " + getTestToken(),
            },
            expectedStatus: http.StatusOK,
            expectedCount:  10,
            expectedTotal:  100,
        },
        {
            name: "get users with search",
            queryParams: map[string]string{
                "page":       "0",
                "pageSize":   "10",
                "textSearch": "john",
            },
            headers: map[string]string{
                "Authorization": "Bearer " + getTestToken(),
            },
            expectedStatus: http.StatusOK,
            expectedCount:  5,
            expectedTotal:  5,
        },
        {
            name: "get users with sorting",
            queryParams: map[string]string{
                "page":        "0",
                "pageSize":    "10",
                "sortBy":      "email",
                "sortOrder":   "ASC",
            },
            headers: map[string]string{
                "Authorization": "Bearer " + getTestToken(),
            },
            expectedStatus: http.StatusOK,
            expectedCount:  10,
            expectedTotal:  100,
        },
        {
            name: "unauthorized access",
            queryParams: map[string]string{
                "page":     "0",
                "pageSize": "10",
            },
            headers:        map[string]string{},
            expectedStatus: http.StatusUnauthorized,
            expectedCount:  0,
            expectedTotal:  0,
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Arrange
            req, err := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users", nil)
            suite.NoError(err)
            
            // Add query parameters
            q := req.URL.Query()
            for key, value := range tt.queryParams {
                q.Add(key, value)
            }
            req.URL.RawQuery = q.Encode()
            
            // Set headers
            for key, value := range tt.headers {
                req.Header.Set(key, value)
            }
            
            // Act
            resp, err := suite.client.Do(req)
            suite.NoError(err)
            defer resp.Body.Close()
            
            // Assert
            suite.Equal(tt.expectedStatus, resp.StatusCode)
            
            if tt.expectedStatus == http.StatusOK {
                var response dto.PagedResponse[dto.UserListResponse]
                err = json.NewDecoder(resp.Body).Decode(&response)
                suite.NoError(err)
                suite.Len(response.Data, tt.expectedCount)
                suite.Equal(tt.expectedTotal, response.Total)
            }
        })
    }
}

// Helper function to get response field value
func getResponseField(response *dto.UserResponse, fieldName string) interface{} {
    v := reflect.ValueOf(response).Elem()
    return v.FieldByName(fieldName).Interface()
}
```

### Endpoint Testing

```go
// test/api/user_api_test.go
package api

import (
    "testing"
    "net/http"
    "github.com/stretchr/testify/suite"
)

type UserAPITestSuite struct {
    suite.Suite
    client *http.Client
    baseURL string
}

func (suite *UserAPITestSuite) TestGetUsers_WithPagination() {
    // Arrange
    req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users?page=0&pageSize=10", nil)
    req.Header.Set("Authorization", "Bearer "+getTestToken())
    
    // Act
    resp, err := suite.client.Do(req)
    suite.NoError(err)
    defer resp.Body.Close()
    
    // Assert
    suite.Equal(http.StatusOK, resp.StatusCode)
    
    var response dto.PagedResponse[dto.UserListResponse]
    err = json.NewDecoder(resp.Body).Decode(&response)
    suite.NoError(err)
    suite.Len(response.Data, 10)
    suite.Equal(int64(100), response.Total)
}

func (suite *UserAPITestSuite) TestCreateUser_ValidationError() {
    // Arrange
    reqBody := dto.CreateUserRequest{
        FirstName: "", // Invalid: empty first name
        Email:     "invalid-email", // Invalid email format
    }
    
    jsonBody, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/admin/users", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+getTestToken())
    
    // Act
    resp, err := suite.client.Do(req)
    suite.NoError(err)
    defer resp.Body.Close()
    
    // Assert
    suite.Equal(http.StatusBadRequest, resp.StatusCode)
    
    var errorResponse wrapper.Response
    err = json.NewDecoder(resp.Body).Decode(&errorResponse)
    suite.NoError(err)
    suite.Contains(errorResponse.Message, "validation failed")
}
```

### Authentication Testing

```go
func (suite *UserAPITestSuite) TestGetUsers_Unauthorized() {
    // Arrange
    req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users", nil)
    // No Authorization header
    
    // Act
    resp, err := suite.client.Do(req)
    suite.NoError(err)
    defer resp.Body.Close()
    
    // Assert
    suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (suite *UserAPITestSuite) TestGetUsers_InvalidToken() {
    // Arrange
    req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users", nil)
    req.Header.Set("Authorization", "Bearer invalid-token")
    
    // Act
    resp, err := suite.client.Do(req)
    suite.NoError(err)
    defer resp.Body.Close()
    
    // Assert
    suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}
```

## 🗄️ Database Testing

### Table-Driven Repository Testing

```go
// test/repository/user_repository_test.go
package repository

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go"
)

type UserRepositoryTestSuite struct {
    suite.Suite
    container testcontainers.Container
    db        *gorm.DB
    repo      repository.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
    suite.setupTestDatabase()
    suite.repo = repository.NewUserRepository(suite.db)
}

func (suite *UserRepositoryTestSuite) TestCreateUser_TableDriven() {
    tests := []struct {
        name           string
        user           *entity.User
        expectedError  error
        expectedFields map[string]interface{}
    }{
        {
            name: "valid user creation",
            user: &entity.User{
                FirstName: "John",
                LastName:  "Doe",
                Email:     "john.doe@example.com",
                Password:  "hashedpassword",
            },
            expectedError: nil,
            expectedFields: map[string]interface{}{
                "FirstName": "John",
                "LastName":  "Doe",
                "Email":     "john.doe@example.com",
            },
        },
        {
            name: "user with customer ID",
            user: &entity.User{
                FirstName:  "Customer",
                LastName:   "User",
                Email:     "customer@example.com",
                Password:  "hashedpassword",
                CustomerID: &uuid.UUID{},
            },
            expectedError: nil,
            expectedFields: map[string]interface{}{
                "FirstName": "Customer",
                "LastName":  "User",
                "Email":     "customer@example.com",
            },
        },
        {
            name: "system admin user",
            user: &entity.User{
                FirstName:      "Admin",
                LastName:       "User",
                Email:         "admin@example.com",
                Password:      "hashedpassword",
                IsSystemAdmin: true,
            },
            expectedError: nil,
            expectedFields: map[string]interface{}{
                "FirstName":      "Admin",
                "LastName":       "User",
                "Email":          "admin@example.com",
                "IsSystemAdmin":  true,
            },
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Act
            createdUser, err := suite.repo.Create(context.Background(), *tt.user)
            
            // Assert
            if tt.expectedError != nil {
                suite.Error(err)
                suite.Nil(createdUser)
            } else {
                suite.NoError(err)
                suite.NotNil(createdUser)
                suite.NotEqual(uuid.Nil, createdUser.ID)
                
                // Verify expected fields
                for field, expectedValue := range tt.expectedFields {
                    suite.Equal(expectedValue, getFieldValue(createdUser, field))
                }
                
                // Verify in database
                var foundUser entity.User
                err = suite.db.First(&foundUser, createdUser.ID).Error
                suite.NoError(err)
                suite.Equal(tt.user.FirstName, foundUser.FirstName)
                suite.Equal(tt.user.Email, foundUser.Email)
            }
        })
    }
}

func (suite *UserRepositoryTestSuite) TestGetByID_TableDriven() {
    // Setup test data
    testUsers := []entity.User{
        {
            ID:        uuid.New(),
            FirstName: "John",
            LastName:  "Doe",
            Email:     "john.doe@example.com",
        },
        {
            ID:            uuid.New(),
            FirstName:     "Admin",
            LastName:      "User",
            Email:         "admin@example.com",
            IsSystemAdmin: true,
        },
    }
    
    // Create test users
    for _, user := range testUsers {
        suite.repo.Create(context.Background(), user)
    }
    
    tests := []struct {
        name           string
        userID         uuid.UUID
        isSystemAdmin  bool
        expectedError  error
        expectedUser   *entity.User
    }{
        {
            name:          "get regular user",
            userID:        testUsers[0].ID,
            isSystemAdmin: false,
            expectedError: nil,
            expectedUser:  &testUsers[0],
        },
        {
            name:          "get system admin user",
            userID:        testUsers[1].ID,
            isSystemAdmin: true,
            expectedError: nil,
            expectedUser:  &testUsers[1],
        },
        {
            name:          "user not found",
            userID:        uuid.New(),
            isSystemAdmin: false,
            expectedError: gorm.ErrRecordNotFound,
            expectedUser:  nil,
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Act
            user, err := suite.repo.GetByID(context.Background(), tt.userID, tt.isSystemAdmin)
            
            // Assert
            if tt.expectedError != nil {
                suite.Error(err)
                suite.Nil(user)
                suite.True(errors.Is(err, tt.expectedError))
            } else {
                suite.NoError(err)
                suite.NotNil(user)
                suite.Equal(tt.expectedUser.ID, user.ID)
                suite.Equal(tt.expectedUser.FirstName, user.FirstName)
                suite.Equal(tt.expectedUser.Email, user.Email)
            }
        })
    }
}

func (suite *UserRepositoryTestSuite) TestGetAll_TableDriven() {
    // Setup test data with various scenarios
    testUsers := []entity.User{
        {FirstName: "Alice", LastName: "Johnson", Email: "alice@example.com"},
        {FirstName: "Bob", LastName: "Smith", Email: "bob@example.com"},
        {FirstName: "Charlie", LastName: "Brown", Email: "charlie@example.com"},
        {FirstName: "Diana", LastName: "Prince", Email: "diana@example.com"},
        {FirstName: "Eve", LastName: "Wilson", Email: "eve@example.com"},
    }
    
    customerID := uuid.New()
    for i, user := range testUsers {
        user.ID = uuid.New()
        if i < 3 {
            user.CustomerID = &customerID
        }
        suite.repo.Create(context.Background(), user)
    }
    
    tests := []struct {
        name           string
        offset         int
        limit          int
        sortBy         string
        sortOrder      string
        textSearch     string
        customerID     *uuid.UUID
        isSystemAdmin  bool
        expectedCount  int
        expectedTotal  int64
    }{
        {
            name:          "get all users with pagination",
            offset:        0,
            limit:         3,
            sortBy:        "email",
            sortOrder:     "ASC",
            textSearch:    "",
            customerID:    nil,
            isSystemAdmin: false,
            expectedCount: 3,
            expectedTotal: 5,
        },
        {
            name:          "get users with search",
            offset:        0,
            limit:         10,
            sortBy:        "first_name",
            sortOrder:     "ASC",
            textSearch:    "alice",
            customerID:    nil,
            isSystemAdmin: false,
            expectedCount: 1,
            expectedTotal: 1,
        },
        {
            name:          "get users by customer",
            offset:        0,
            limit:         10,
            sortBy:        "email",
            sortOrder:     "ASC",
            textSearch:    "",
            customerID:    &customerID,
            isSystemAdmin: false,
            expectedCount: 3,
            expectedTotal: 3,
        },
        {
            name:          "get system admin users",
            offset:        0,
            limit:         10,
            sortBy:        "email",
            sortOrder:     "ASC",
            textSearch:    "",
            customerID:    nil,
            isSystemAdmin: true,
            expectedCount: 0, // No system admin users in test data
            expectedTotal: 0,
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Act
            users, total, err := suite.repo.GetAll(
                context.Background(),
                tt.offset,
                tt.limit,
                tt.sortBy,
                tt.sortOrder,
                tt.textSearch,
                tt.customerID,
                tt.isSystemAdmin,
            )
            
            // Assert
            suite.NoError(err)
            suite.Len(users, tt.expectedCount)
            suite.Equal(tt.expectedTotal, total)
        })
    }
}

// Helper function to get field value using reflection
func getFieldValue(user *entity.User, fieldName string) interface{} {
    v := reflect.ValueOf(user).Elem()
    return v.FieldByName(fieldName).Interface()
}
```

### Repository Testing

```go
// test/repository/user_repository_test.go
package repository

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go"
)

type UserRepositoryTestSuite struct {
    suite.Suite
    container testcontainers.Container
    db        *gorm.DB
    repo      repository.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
    suite.setupTestDatabase()
    suite.repo = repository.NewUserRepository(suite.db)
}

func (suite *UserRepositoryTestSuite) TestCreateUser_Success() {
    // Arrange
    user := fixtures.CreateTestUser()
    
    // Act
    createdUser, err := suite.repo.Create(context.Background(), *user)
    
    // Assert
    suite.NoError(err)
    suite.NotNil(createdUser)
    suite.NotEqual(uuid.Nil, createdUser.ID)
    suite.Equal(user.FirstName, createdUser.FirstName)
    suite.Equal(user.Email, createdUser.Email)
}

func (suite *UserRepositoryTestSuite) TestGetByID_NotFound() {
    // Arrange
    nonExistentID := uuid.New()
    
    // Act
    user, err := suite.repo.GetByID(context.Background(), nonExistentID, false)
    
    // Assert
    suite.Error(err)
    suite.Nil(user)
    suite.True(errors.Is(err, gorm.ErrRecordNotFound))
}

func (suite *UserRepositoryTestSuite) TestGetAll_WithPagination() {
    // Arrange
    // Create multiple test users
    for i := 0; i < 15; i++ {
        user := fixtures.CreateTestUser()
        user.Email = fmt.Sprintf("user%d@example.com", i)
        suite.repo.Create(context.Background(), *user)
    }
    
    // Act
    users, total, err := suite.repo.GetAll(context.Background(), 0, 10, "email", "ASC", "", nil, false)
    
    // Assert
    suite.NoError(err)
    suite.Len(users, 10)
    suite.Equal(int64(15), total)
}
```

## 🔧 Test Utilities

### Test Helpers

```go
// test/helpers/test_helpers.go
package helpers

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// CreateTestServer creates a test HTTP server
func CreateTestServer(router *gin.Engine) *httptest.Server {
    return httptest.NewServer(router)
}

// MakeRequest makes an HTTP request for testing
func MakeRequest(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
    var reqBody []byte
    var err error
    
    if body != nil {
        reqBody, err = json.Marshal(body)
        if err != nil {
            return nil, err
        }
    }
    
    req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    for key, value := range headers {
        req.Header.Set(key, value)
    }
    
    client := &http.Client{}
    return client.Do(req)
}

// GenerateTestToken generates a test JWT token
func GenerateTestToken(userID uuid.UUID, customerID *uuid.UUID) string {
    // Implementation for generating test tokens
    return "test-token"
}

// CleanupTestData cleans up test data from database
func CleanupTestData(db *gorm.DB) error {
    return db.Exec("TRUNCATE TABLE users, customers, roles CASCADE").Error
}
```

### Test Configuration

```go
// test/config/test_config.go
package config

import (
    "os"
    "testing"
)

// SetupTestEnvironment sets up test environment variables
func SetupTestEnvironment(t *testing.T) {
    os.Setenv("DB_HOST", "localhost")
    os.Setenv("DB_PORT", "5432")
    os.Setenv("DB_NAME", "test_db")
    os.Setenv("DB_USER", "test_user")
    os.Setenv("DB_PASSWORD", "test_password")
    os.Setenv("JWT_SECRET", "test-secret-key")
    os.Setenv("ENVIRONMENT", "test")
}

// CleanupTestEnvironment cleans up test environment
func CleanupTestEnvironment() {
    os.Unsetenv("DB_HOST")
    os.Unsetenv("DB_PORT")
    os.Unsetenv("DB_NAME")
    os.Unsetenv("DB_USER")
    os.Unsetenv("DB_PASSWORD")
    os.Unsetenv("JWT_SECRET")
    os.Unsetenv("ENVIRONMENT")
}
```

## 📊 Test Coverage

### Coverage Targets

| Component | Target Coverage | Priority |
|-----------|----------------|----------|
| **Domain Layer** | 90% | High |
| **Application Layer** | 85% | High |
| **Infrastructure Layer** | 80% | Medium |
| **Adapter Layer** | 75% | Medium |

### Coverage Commands

```bash
# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

### Coverage Badge

Add to README.md:
```markdown
![Test Coverage](https://img.shields.io/badge/coverage-85%25-green)
```

## 🚀 Performance Testing

### Benchmark Tests

```go
// test/benchmark/user_service_benchmark_test.go
package benchmark

import (
    "context"
    "testing"
    "github.com/google/uuid"
)

func BenchmarkUserService_CreateUser(b *testing.B) {
    service := setupTestUserService()
    req := &dto.CreateUserRequest{
        FirstName: "Benchmark",
        LastName:  "User",
        Email:     "benchmark@example.com",
        Password:  "password123",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req.Email = fmt.Sprintf("benchmark%d@example.com", i)
        _, err := service.CreateUser(context.Background(), req)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkUserRepository_GetAll(b *testing.B) {
    repo := setupTestUserRepository()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _, err := repo.GetAll(context.Background(), 0, 100, "email", "ASC", "", nil, false)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 🔍 Test Data Management

### Test Fixtures

```go
// test/fixtures/user_fixtures.go
package fixtures

import (
    "admin-api/internal/domain/entity"
    "github.com/google/uuid"
)

var (
    TestUsers = []entity.User{
        {
            ID:        uuid.New(),
            FirstName: "John",
            LastName:  "Doe",
            Email:     "john.doe@example.com",
        },
        {
            ID:        uuid.New(),
            FirstName: "Jane",
            LastName:  "Smith",
            Email:     "jane.smith@example.com",
        },
    }
    
    TestCustomers = []entity.Customer{
        {
            ID:   uuid.New(),
            Name: "Test Customer 1",
        },
        {
            ID:   uuid.New(),
            Name: "Test Customer 2",
        },
    }
)

func CreateTestUser() *entity.User {
    return &TestUsers[0]
}

func CreateTestUserWithEmail(email string) *entity.User {
    user := CreateTestUser()
    user.Email = email
    return user
}
```

## 📋 Test Checklist

### Before Writing Tests

- [ ] **Understand the requirement**: What should the code do?
- [ ] **Identify edge cases**: What could go wrong?
- [ ] **Plan test scenarios**: What should be tested?
- [ ] **Choose appropriate test type**: Unit, integration, or E2E?

### Writing Tests

- [ ] **Follow naming convention**: `Test<Struct>_<Method>_<Scenario>`
- [ ] **Use descriptive test names**: Clear what is being tested
- [ ] **Follow AAA pattern**: Arrange, Act, Assert
- [ ] **Test one thing at a time**: Single responsibility
- [ ] **Use meaningful assertions**: Specific, not generic
- [ ] **Mock external dependencies**: Keep tests isolated
- [ ] **Clean up test data**: Don't leave test artifacts

### Running Tests

- [ ] **Run tests locally**: Before committing
- [ ] **Check coverage**: Ensure adequate coverage
- [ ] **Run all test types**: Unit, integration, E2E
- [ ] **Check for race conditions**: Use `-race` flag
- [ ] **Verify CI/CD pipeline**: Tests pass in CI

## 🛠️ Test Tools and Libraries

### Required Tools

```go
// go.mod dependencies for testing
require (
    github.com/stretchr/testify v1.8.4
    github.com/testcontainers/testcontainers-go v0.26.0
    github.com/testcontainers/testcontainers-go/modules/postgres v0.26.0
)
```

### IDE Integration

- **VS Code**: Go Test Explorer extension
- **IntelliJ**: Built-in Go testing support
- **GitHub Actions**: Automated testing in CI/CD

### Test Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./internal/application/service -run TestUserService_CreateUser

# Run benchmark tests
go test -bench=. ./test/benchmark/

# Run tests with race detection
go test -race ./...

# Run integration tests only
go test -tags=integration ./test/integration/
```

## 📚 Best Practices

### Table-Driven Testing Guidelines

#### When to Use Table-Driven Tests
- **Multiple scenarios**: When testing the same function with different inputs
- **Edge cases**: When testing various boundary conditions
- **Error conditions**: When testing different error scenarios
- **Validation**: When testing different validation rules
- **API endpoints**: When testing different HTTP status codes and responses

#### Table Structure Best Practices
```go
tests := []struct {
    name           string        // Descriptive test case name
    input          interface{}   // Input data
    setupMocks     func()        // Mock setup function
    expectedResult interface{}   // Expected output
    expectedError  error         // Expected error
    expectedStatus int           // For API tests
}{
    // Test cases...
}
```

#### Naming Conventions for Table-Driven Tests
- Use descriptive names that explain the scenario
- Include expected outcome in the name
- Use consistent naming patterns

**Good Examples:**
- `TestUserService_CreateUser_TableDriven`
- `TestUserService_CreateUser_ValidationScenarios`
- `TestUserService_CreateUser_ErrorConditions`

#### Table-Driven Test Structure
```go
func TestFunction_TableDriven(t *testing.T) {
    tests := []struct {
        name string
        // ... test case fields
    }{
        // ... test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            // Act
            // Assert
        })
    }
}
```

#### Benefits of Table-Driven Tests
1. **Comprehensive Coverage**: Test multiple scenarios in one test
2. **Maintainability**: Easy to add new test cases
3. **Readability**: Clear structure and easy to understand
4. **Reduced Duplication**: Avoid repetitive test code
5. **Better Organization**: Related test cases grouped together

#### Common Patterns
```go
// Pattern 1: Simple input/output testing
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"empty string", "", ""},
    {"single word", "hello", "hello"},
    {"multiple words", "hello world", "hello world"},
}

// Pattern 2: Error condition testing
tests := []struct {
    name          string
    input         interface{}
    expectedError error
}{
    {"invalid input", nil, ErrInvalidInput},
    {"empty input", "", ErrEmptyInput},
}

// Pattern 3: Mock-based testing
tests := []struct {
    name       string
    setupMocks func(*MockRepository)
    expected   interface{}
}{
    {"success case", func(m *MockRepository) {
        m.On("Get", "id").Return(data, nil)
    }, expectedData},
}
```

### General Guidelines

1. **Test Early, Test Often**
   - Write tests alongside code
   - Don't leave testing for later
   - Test-driven development (TDD) when possible

2. **Keep Tests Simple**
   - One assertion per test
   - Clear test names
   - Minimal setup and teardown

3. **Use Descriptive Names**
   - Test names should describe the scenario
   - Include expected outcome
   - Make failures easy to understand

4. **Maintain Test Data**
   - Use fixtures for test data
   - Clean up after tests
   - Don't rely on test order

5. **Mock External Dependencies**
   - Don't test external services
   - Use interfaces for mocking
   - Keep tests fast and reliable

### Common Anti-patterns

❌ **Don't:**
- Test implementation details
- Write brittle tests
- Ignore test failures
- Skip tests in CI/CD
- Write tests that depend on each other
- **Table-Driven Anti-patterns:**
  - Overly complex test tables with too many fields
  - Unclear test case names
  - Mixing different types of tests in one table
  - Not using subtests for table-driven tests
  - Hardcoding expected values instead of using variables

✅ **Do:**
- Test behavior, not implementation
- Write maintainable tests
- Fix failing tests immediately
- Run all tests in CI/CD
- Keep tests independent
- **Table-Driven Best Practices:**
  - Use descriptive test case names
  - Keep table structure simple and focused
  - Use subtests for better organization
  - Use helper functions for complex setup
  - Group related test cases together

## 📈 Continuous Testing

### CI/CD Integration

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_DB: test_db
          POSTGRES_USER: test_user
          POSTGRES_PASSWORD: test_password
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: go test -v -race ./...
    
    - name: Run integration tests
      run: go test -v -tags=integration ./test/integration/
    
    - name: Generate coverage report
      run: go test -coverprofile=coverage.out ./...
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

---

Following these testing guidelines ensures high code quality, reliability, and maintainability of the IoT Admin Service API project. 