# Requirements Document

## Introduction

This document defines the requirements for a reusable Go API boilerplate with a clear, production-ready structure. The boilerplate serves as a generic starting template for any backend project — not tied to a specific domain.

The boilerplate uses Echo v4, GORM Gen, Cobra CLI, OpenTelemetry, and Viper config, with manual wiring via explicit constructor injection, a repository interface layer for testability, struct-based services, domain error separation, explicit route registration, health/readiness probes, graceful shutdown coordination, and OpenAPI documentation generation.

Dependencies are wired manually via constructor injection — no DI framework needed. This is the idiomatic Go approach: all dependencies are visible in function signatures, easy to test, and easy to understand.

A single example domain module (Auth) is included to demonstrate all patterns end-to-end.

## Glossary

- **Boilerplate**: A reusable project skeleton providing foundational code structure, patterns, and one example module as a starting template for any Go API project
- **App_Struct**: The central application struct holding all infrastructure and service references, providing `Start()` and `Shutdown()` lifecycle methods with ordered startup/teardown
- **Route_Registry**: Explicit route registration module using a declarative route table, replacing init()-based implicit registration
- **Service_Layer**: Struct-based service implementations with interface contracts for testability and mockability
- **Repository_Layer**: Interface-defined data access layer abstracting GORM Gen queries behind contracts for unit testing with mocks
- **Domain_Error**: Typed error definitions with error code and HTTP status mapping for consistent error handling
- **Health_Probe**: HTTP endpoint reporting application liveness status
- **Readiness_Probe**: HTTP endpoint reporting application readiness to serve traffic
- **Shutdown_Coordinator**: Component orchestrating graceful shutdown in dependency order
- **OpenAPI_Generator**: Tooling producing OpenAPI 3.0 specification from annotated handler code using swaggo/swag
- **Echo_Server**: The labstack/echo v4 HTTP framework instance
- **GORM_Gen**: Type-safe query generator for GORM producing compile-time checked database queries
- **Cobra_CLI**: Command-line interface framework for application entry points
- **Middleware_Stack**: Ordered chain of HTTP middleware: RequestID → OTel → CORS → Recover → Logger → Gzip → JWT → RBAC
- **Config_Manager**: Viper-based configuration system loading from .env files and environment variables into typed structs

## Requirements

### Requirement 1: Project Structure Scaffolding

**User Story:** As a developer, I want a well-organized project directory structure following Go conventions, so that I can navigate and extend the codebase predictably.

#### Acceptance Criteria

1. THE Boilerplate SHALL provide a directory layout with the following top-level packages: `cmd/`, `internal/`, `pkg/`, `docs/`, and `migrations/`
2. THE Boilerplate SHALL organize `internal/` into sub-packages: `config/`, `app/`, `domain/`, `repository/`, `service/`, `handler/`, and `middleware/`
3. THE Boilerplate SHALL organize `internal/domain/` with one example sub-package `auth/` containing entity definitions, error codes, and constants
4. THE Boilerplate SHALL include a `Makefile` with individually invocable targets: `build`, `test`, `lint`, `migrate`, `seed`, `generate` (GORM Gen + mockgen), and `docs`
5. THE Boilerplate SHALL include a `Dockerfile` with multi-stage build: a builder stage compiling the Go binary, and a final stage based on Alpine containing only the binary
6. THE Boilerplate SHALL include a `.env.example` file listing every environment variable with descriptive placeholder values and inline comments
7. WHEN a developer runs `make build`, THE Boilerplate SHALL produce a single statically-linked Go binary

### Requirement 2: Application Wiring and Lifecycle (Manual DI)

**User Story:** As a developer, I want explicit manual wiring with constructor injection and a clear application lifecycle, so that all dependencies are visible and the startup/shutdown order is predictable.

#### Acceptance Criteria

1. THE Boilerplate SHALL define an `App` struct in `internal/app/` that holds all infrastructure references (DB, cache, HTTP server, gRPC server) and service instances as struct fields
2. THE App struct SHALL provide a `New(cfg *config.Config) (*App, error)` constructor that wires all dependencies in explicit order: config → infra → repositories → services → handlers → routes
3. THE App struct SHALL provide a `Start(ctx context.Context) error` method that opens infrastructure connections in dependency order: database → cache → HTTP server → gRPC server
4. THE App struct SHALL provide a `Shutdown(ctx context.Context) error` method that closes infrastructure in reverse order: gRPC → HTTP → cache → database
5. THE Boilerplate SHALL NOT use any DI framework, global container, or reflection-based service locator — all dependencies are passed via constructor parameters
6. THE Boilerplate SHALL wire services via constructor injection: `NewAuthService(repo AuthRepository, hasher PasswordHasher, jwt JWTManager) *AuthService`
7. IF any infrastructure fails to open during `Start()`, THEN the App SHALL close already-opened components in reverse order and return the error

### Requirement 3: Repository Interface Layer

**User Story:** As a developer, I want a repository interface layer between services and GORM Gen queries, so that I can write unit tests with mock implementations.

#### Acceptance Criteria

1. THE Repository_Layer SHALL define one exported Go interface per domain module containing all data access method signatures
2. THE Repository_Layer SHALL provide concrete implementations using GORM Gen type-safe queries internally
3. THE Repository_Layer SHALL accept `context.Context` as the first parameter in every method for tracing propagation and cancellation
4. THE Repository_Layer SHALL return Domain_Error types instead of raw GORM errors, mapping constraint violations to `ErrConflict` and connection failures to `ErrInternal`
5. WHEN a database record is not found, THE Repository_Layer SHALL return a typed `ErrNotFound` error including the entity name and lookup parameters
6. THE Repository_Layer SHALL expose a transaction method accepting a callback function, executing all operations within a single database transaction and rolling back on error
7. THE Boilerplate SHALL include example mock implementations generated using mockgen for the Auth repository interface
8. THE Repository_Layer SHALL receive the `*gorm.DB` instance via constructor parameter injection

### Requirement 4: Struct-Based Service Layer

**User Story:** As a developer, I want services defined as struct methods instead of package-level functions, so that I can inject dependencies and mock services in tests.

#### Acceptance Criteria

1. THE Service_Layer SHALL define each service as a Go struct with all dependency fields typed as Go interfaces to enable mock injection
2. THE Service_Layer SHALL define a Go interface for each service struct, and handler structs SHALL depend on the service interface type
3. THE Service_Layer SHALL accept dependencies through constructor functions: `NewXxxService(repo XxxRepository, ...) *XxxService`
4. THE Service_Layer SHALL accept `context.Context` as the first parameter in every public method and SHALL NOT accept HTTP-specific types
5. THE Service_Layer SHALL return domain-specific error types using Go convention `(result, error)` return signatures
6. THE Boilerplate SHALL include an example Auth service implementation demonstrating the pattern
7. THE Boilerplate SHALL include example unit tests for the Auth service demonstrating mock injection via constructor
8. THE Boilerplate SHALL initialize services in the `App.New()` constructor after infrastructure is created, passing repository instances directly

### Requirement 5: Domain Error System

**User Story:** As a developer, I want typed errors with automatic HTTP status mapping, so that error handling is consistent across the application.

#### Acceptance Criteria

1. THE Domain_Error system SHALL define a base error interface with methods: `Error() string`, `Code() string`, `HTTPStatus() int`, and `Domain() string`
2. THE Domain_Error system SHALL provide error constructors with codes prefixed by uppercase domain name (e.g., "AUTH_INVALID_CREDENTIALS")
3. WHEN a Domain_Error is returned from the Service_Layer, THE Echo_Server SHALL map it to the corresponding HTTP status and structured JSON response automatically
4. THE Domain_Error system SHALL use Go constants for error codes to prevent duplicates at compile time
5. THE Domain_Error system SHALL support error wrapping with `errors.Is()` and `errors.As()` compatibility via `Unwrap() error`
6. IF an untyped error reaches the handler layer, THEN THE Echo_Server SHALL return HTTP 500 with a generic message without exposing internal details
7. WHEN a Domain_Error wraps an underlying error, it SHALL preserve the original for logging while only exposing the domain message in the HTTP response

### Requirement 6: Explicit Route Registration

**User Story:** As a developer, I want explicit route registration via a declarative route table, so that routes are traceable and testable.

#### Acceptance Criteria

1. THE Route_Registry SHALL define routes as a declarative slice containing: HTTP method, path, handler function, middleware list, and a public flag (excluded from JWT)
2. THE Route_Registry SHALL support route grouping by API version and domain module
3. THE Route_Registry SHALL accept handler struct methods as route handlers (not package-level functions)
4. WHEN the server starts, THE Route_Registry SHALL register all routes in a single function call with the complete table in one source file
5. THE Route_Registry SHALL support per-route middleware assignment for authentication and authorization
6. WHEN the application starts, THE Echo_Server SHALL log all registered routes with methods, paths, and middleware names at INFO level
7. IF a duplicate method+path combination is detected, THEN THE Route_Registry SHALL return an error preventing server startup

### Requirement 7: HTTP Server and Middleware Stack

**User Story:** As a developer, I want a pre-configured Echo v4 HTTP server with standard middleware, so that I have consistent request processing.

#### Acceptance Criteria

1. THE Echo_Server SHALL apply middleware in order: RequestID → OpenTelemetry → CORS → Recover → Logger → Gzip → JWT → RBAC
2. THE Echo_Server SHALL use structured slog logging for request logs including method, path, status, latency (ms), and request ID
3. THE Echo_Server SHALL use the sonic JSON serializer for request/response encoding
4. THE Echo_Server SHALL support configurable CORS origins, methods, and headers via Config_Manager
5. IF a route is marked as public, THEN THE Middleware_Stack SHALL skip JWT authentication and RBAC for that route
6. THE Middleware_Stack SHALL propagate OpenTelemetry trace context through all middleware and handler execution
7. IF a panic occurs during request processing, THEN THE Echo_Server SHALL recover, log the stack trace, and return HTTP 500 with a generic message
8. THE Echo_Server SHALL enforce configurable max request body size (default 4MB) returning HTTP 413 if exceeded
9. THE Echo_Server SHALL apply configurable read/write timeouts (default 30 seconds each)

### Requirement 8: Configuration Management

**User Story:** As a developer, I want a typed configuration system loaded from environment variables, so that I can access settings with compile-time safety.

#### Acceptance Criteria

1. THE Config_Manager SHALL load from `.env` file using godotenv and bind to a typed Go struct using Viper with mapstructure tags
2. THE Config_Manager SHALL resolve values in precedence: OS env vars > .env file > Viper defaults
3. THE Config_Manager SHALL validate required fields (marked with `required` struct tag) at startup, listing all missing fields in one error
4. THE Config_Manager SHALL organize config into nested structs: App, Server, Database, Cache, JWT, and OpenTelemetry
5. IF required values are missing, THEN THE Config_Manager SHALL log all missing field names and exit with code 1
6. THE Config_Manager SHALL expose `Load() (*Config, error)` so the CLI entry point can load config explicitly and pass it into `app.New(cfg)`
7. IF the `.env` file does not exist, THEN THE Config_Manager SHALL continue startup relying on OS env vars and defaults

### Requirement 9: Database Layer with GORM Gen

**User Story:** As a developer, I want type-safe database queries with PostgreSQL, so that I have compile-time query validation.

#### Acceptance Criteria

1. THE Boilerplate SHALL include GORM model definitions for the example Auth entity (User)
2. THE Boilerplate SHALL include a GORM Gen configuration generating type-safe query objects to `internal/query/`
3. THE Boilerplate SHALL include a Cobra CLI command `genq` that regenerates query code from model definitions
4. THE Boilerplate SHALL configure connection pooling: max open (default 10), max idle (default 5), max lifetime (default 5 min) from config
5. THE Boilerplate SHALL integrate OpenTelemetry GORM plugin (`otelgorm`) for automatic query tracing
6. THE Boilerplate SHALL include a Cobra CLI `migrate` command executing GORM AutoMigrate for defined models
7. IF a migration fails, THEN the command SHALL print the error and exit without applying partial changes

### Requirement 10: Health Check and Readiness Probes

**User Story:** As a developer, I want standardized health and readiness endpoints for container orchestrators.

#### Acceptance Criteria

1. THE Health_Probe SHALL respond on `GET /healthz` with HTTP 200 and `{"status": "alive"}` within 500ms
2. THE Readiness_Probe SHALL respond on `GET /readyz` with HTTP 200 and `{"status": "ready", "checks": {"database": "up", "redis": "up"}}` when all dependencies pass
3. THE Readiness_Probe SHALL check database connectivity with a 2-second timeout ping
4. THE Readiness_Probe SHALL check Redis connectivity with a 2-second timeout ping
5. IF any check fails, THEN THE Readiness_Probe SHALL respond with HTTP 503 and `{"status": "not_ready", "checks": {"<component>": "down"}}` with error description
6. THE Health_Probe and Readiness_Probe SHALL be excluded from JWT middleware and request logging

### Requirement 11: Graceful Shutdown Coordination

**User Story:** As a developer, I want coordinated graceful shutdown so that in-flight requests complete without data loss.

#### Acceptance Criteria

1. WHEN the application receives SIGTERM or SIGINT, THE Shutdown_Coordinator SHALL initiate graceful shutdown and stop accepting new connections within 1 second
2. THE Shutdown_Coordinator SHALL wait for in-flight requests with configurable timeout (default 30 seconds, range 5-120)
3. THE Shutdown_Coordinator SHALL close components in reverse dependency order: HTTP server → background workers → cache → database, max 10 seconds per component
4. IF the shutdown timeout is exceeded, THEN THE Shutdown_Coordinator SHALL force-terminate and exit with code 1
5. WHEN all components close successfully, THE Shutdown_Coordinator SHALL exit with code 0
6. THE Shutdown_Coordinator SHALL log each component's shutdown result and duration in milliseconds
7. IF a second SIGTERM/SIGINT is received during shutdown, THEN THE Shutdown_Coordinator SHALL force-terminate immediately with code 1

### Requirement 12: OpenAPI Documentation Generation

**User Story:** As a developer, I want auto-generated OpenAPI 3.0 documentation from annotated handler code.

#### Acceptance Criteria

1. THE OpenAPI_Generator SHALL produce an OpenAPI 3.0 YAML file at `docs/openapi.yaml` from swaggo/swag annotations
2. THE Boilerplate SHALL include a Cobra CLI command `docs` that regenerates the specification
3. IF annotations are invalid, THEN the `docs` command SHALL exit with code 1 logging the file and line of the error
4. THE OpenAPI_Generator SHALL document request/response schemas from handler structs including field types and validation tags
5. THE OpenAPI_Generator SHALL include Bearer JWT security definitions applied to non-public routes
6. WHILE running in non-production, THE Echo_Server SHALL serve Swagger UI at `/docs`
7. IF `/docs` is requested in production, THEN THE Echo_Server SHALL respond with HTTP 404

### Requirement 13: Auth Domain Module Example

**User Story:** As a developer, I want a complete Auth module example demonstrating all improved patterns end-to-end.

#### Acceptance Criteria

1. THE Boilerplate SHALL include an Auth service struct with methods: Register, Login, RefreshToken, and Logout
2. THE Boilerplate SHALL include an Auth repository interface with methods: CreateUser, FindUserByEmail, FindUserByID, and UpdateRefreshToken
3. THE Boilerplate SHALL include Auth handler functions with swaggo annotations for: POST /v1/auth/register, POST /v1/auth/login, POST /v1/auth/refresh, DELETE /v1/auth/logout
4. THE Boilerplate SHALL include JWT token generation/validation using golang-jwt/jwt/v5 with configurable expiry (default: 15min access, 7 days refresh)
5. THE Boilerplate SHALL include password hashing using bcrypt with configurable cost (default 10, range 4-31)
6. WHEN login contains invalid credentials, THE Auth service SHALL return Domain_Error "AUTH_INVALID_CREDENTIALS" with HTTP 401
7. WHEN registration email already exists, THE Auth service SHALL return Domain_Error "AUTH_EMAIL_EXISTS" with HTTP 409
8. IF refresh token is expired/invalid, THEN THE Auth service SHALL return Domain_Error "AUTH_INVALID_REFRESH_TOKEN" with HTTP 401
9. THE Boilerplate SHALL include request validation structs with go-playground/validator tags for all Auth endpoints

### Requirement 14: CLI Commands Structure

**User Story:** As a developer, I want Cobra CLI commands for all operational tasks.

#### Acceptance Criteria

1. THE Cobra_CLI SHALL include a `server` command starting the HTTP server with graceful shutdown, listening for SIGINT/SIGTERM/SIGUSR2
2. THE Cobra_CLI SHALL include a `migrate` command with sub-commands: `up`, `down`, and `status`
3. IF `migrate up` or `migrate down` fails, THEN it SHALL exit non-zero with error message indicating the migration name and reason
4. THE Cobra_CLI SHALL include a `seed` command populating the database with example data for development
5. THE Cobra_CLI SHALL include a `genq` command running GORM Gen code generation
6. THE Cobra_CLI SHALL include a `docs` command generating OpenAPI specification
7. THE `server` command SHALL support zero-downtime restarts via Overseer, passing the active listener to the new process
8. IF any command encounters an unrecoverable error, THEN it SHALL exit non-zero with human-readable error to stderr

### Requirement 15: Observability Setup

**User Story:** As a developer, I want pre-configured OpenTelemetry tracing and structured logging for production observability.

#### Acceptance Criteria

1. THE Boilerplate SHALL configure OpenTelemetry SDK with OTLP gRPC exporter for distributed tracing
2. THE Boilerplate SHALL propagate trace context through HTTP middleware, service methods, and repository queries
3. THE Boilerplate SHALL use slog with JSON output in production and text output in development, selected by APP_ENV config
4. THE Boilerplate SHALL include trace-log correlation by injecting trace_id and span_id into slog attributes
5. THE Boilerplate SHALL configure GORM OpenTelemetry plugin for automatic database query span creation with db.name attribute
6. THE Boilerplate SHALL exclude /healthz and /readyz from trace collection
7. THE Boilerplate SHALL set OTel resource attributes: service.name, service.version, and deployment.environment from config

### Requirement 16: Response Format Standardization

**User Story:** As a developer, I want a standardized JSON response format across all endpoints.

#### Acceptance Criteria

1. THE Echo_Server SHALL return successful responses: `{"success": true, "code": "<http_status>", "message": "<msg>", "data": <payload>, "requestId": "<id>"}`
2. THE Echo_Server SHALL return error responses: `{"success": false, "code": "<error_code>", "message": "<msg>", "errors": [{"field": "<name>", "message": "<reason>"}], "requestId": "<id>"}`
3. THE Boilerplate SHALL include a response builder supporting pagination: `{"data": [...], "meta": {"page": N, "perPage": N, "total": N, "totalPages": N}}` with page 1-based, perPage default 10 max 100
4. THE response builder SHALL use Go generics for type-safe response construction
5. WHEN a validation error occurs, THE Echo_Server SHALL return HTTP 422 with per-field error details including field name, message, and rejected value
6. IF X-Request-ID header is absent, THEN THE response builder SHALL generate a UUID v4 as requestId
7. THE response builder SHALL include requestId in every response sourced from X-Request-ID when present

### Requirement 17: gRPC Transport Support

**User Story:** As a developer, I want the boilerplate to support both REST API (Echo) and gRPC as transport layers, so that I can choose the appropriate protocol per project or run both simultaneously.

#### Acceptance Criteria

1. THE Boilerplate SHALL include a gRPC server infra component wired explicitly in `internal/app` alongside the Echo HTTP server
2. THE gRPC server SHALL be configurable to run on a separate port from the HTTP server, defined in Config_Manager (`grpc.addr`)
3. THE Boilerplate SHALL include example `.proto` files for the Auth service defining `Register`, `Login`, `RefreshToken`, and `Logout` RPCs
4. THE Boilerplate SHALL include a `Makefile` target `proto` that generates Go code from `.proto` files using `protoc` with `protoc-gen-go` and `protoc-gen-go-grpc`
5. THE gRPC server SHALL reuse the same Service_Layer structs as the REST handlers — both transports call the same service interface
6. THE gRPC server SHALL include OpenTelemetry interceptors for automatic trace propagation on unary and streaming RPCs
7. THE gRPC server SHALL include a health check service implementing the standard `grpc.health.v1.Health` protocol
8. THE Boilerplate SHALL allow enabling/disabling each transport independently via config: `server.httpEnabled` (default true) and `server.grpcEnabled` (default false)
9. IF both transports are enabled, THEN THE Shutdown_Coordinator SHALL gracefully stop both servers during shutdown in the correct order
10. THE gRPC server SHALL include a unary interceptor for structured slog logging with method name, status code, and latency

### Requirement 18: Structured Logger with Context

**User Story:** As a developer, I want a context-aware structured logger that automatically includes trace correlation and request metadata, so that logs are easy to search and correlate with traces.

#### Acceptance Criteria

1. THE Boilerplate SHALL provide a logger package wrapping slog with a `FromContext(ctx) *slog.Logger` function that extracts trace_id and span_id from the context automatically
2. THE logger SHALL output JSON format when APP_ENV is "production" or "prod", and colored text format otherwise
3. THE logger SHALL include a middleware (HTTP and gRPC) that injects a request-scoped logger into the context with pre-populated attributes: request_id, method, path (or gRPC method)
4. THE logger SHALL support log level configuration via config (`log.level`: debug, info, warn, error) changeable without restart via OS signal (SIGUSR1 cycles levels)
5. THE logger SHALL integrate with OpenTelemetry by forwarding log records to the OTel log bridge when OTel is configured
6. THE logger SHALL provide convenience methods: `logger.Info(ctx, msg, attrs...)`, `logger.Error(ctx, msg, attrs...)` that automatically include context attributes
7. IF OTel is not configured, THEN THE logger SHALL function normally with stdout output only, without errors or degraded behavior
