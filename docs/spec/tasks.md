# Implementation Plan: Go API Boilerplate (Manual Wiring)

## Overview

This plan implements a reusable Go API boilerplate with a clear, production-ready backend structure. Tasks are ordered for incremental development: foundational packages first (config, errors, application wiring), then infrastructure, then business layers (repository → service → handler), then cross-cutting concerns (routes, response, observability), and finally the gRPC transport and CLI commands.

## Tasks

- [x] 1. Foundation packages — Config, Domain Errors, Response Builder
  - [x] 1.1 Create typed config package with Viper loading and validation
    - Create `internal/config/config.go` with nested structs: `AppConfig`, `ServerConfig`, `DatabaseConfig`, `CacheConfig`, `JWTConfig`, `OTelConfig`, `LogConfig`, `GRPCConfig`
    - Implement `Load() (*Config, error)` using godotenv + Viper with mapstructure tags
    - Implement required field validation (struct tag `required:"true"`) that lists all missing fields in one error
    - Implement precedence: OS env > .env file > Viper defaults
    - Return the loaded config explicitly; do not expose a package-level config singleton
    - Create `.env.example` with all variables documented
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7_

  - [x] 1.2 Create domain error system (`pkg/apperror`)
    - Create `pkg/apperror/error.go` with `AppError` interface: `Error()`, `Code()`, `HTTPStatus()`, `Domain()`, `Unwrap()`
    - Implement `domainError` struct with `New()` and `Wrap()` constructors
    - Ensure `errors.Is()` and `errors.As()` compatibility via `Unwrap()`
    - Create `pkg/apperror/codes.go` with common error code constants
    - _Requirements: 5.1, 5.2, 5.4, 5.5, 5.7_

  - [x] 1.3 Create generic response builder (`pkg/response`)
    - Create `pkg/response/response.go` with `Envelope[T]`, `FieldErr`, `PaginatedData[T]`, `Meta` structs
    - Implement `OK[T]()`, `Created[T]()`, `Error()`, `ValidationError()`, `Paginated[T]()` functions
    - Implement `getRequestID()` extracting from X-Request-ID header or generating UUID v4
    - Implement pagination meta calculation: `totalPages = ceil(total / perPage)`, page 1-based, perPage default 10 max 100
    - _Requirements: 16.1, 16.2, 16.3, 16.4, 16.5, 16.6, 16.7_

  - [x]* 1.4 Write property tests for domain error system
    - **Property 4: Domain error code prefix** — verify code equals uppercase `DOMAIN_SUFFIX`
    - **Property 6: Error wrapping preserves chain** — verify `errors.Is()` and `errors.As()` work through wrapped errors
    - **Validates: Requirements 5.2, 5.5**

  - [x]* 1.5 Write property tests for response builder
    - **Property 14: Response envelope structure** — verify OK produces `success:true`, Error produces `success:false` with correct fields
    - **Property 15: Pagination metadata calculation** — verify `totalPages = ceil(T/PP)` for all valid T, P, PP
    - **Property 16: Validation error completeness** — verify N field errors produce exactly N entries in response
    - **Property 17: Request ID propagation** — verify X-Request-ID header value appears in response
    - **Validates: Requirements 16.1, 16.2, 16.3, 16.5, 16.7**

  - [x]* 1.6 Write property tests for config validation
    - **Property 10: Config precedence resolution** — verify OS env > .env > defaults
    - **Property 11: Required config field validation** — verify all missing required fields listed in error
    - **Validates: Requirements 8.2, 8.3**

- [ ] 2. Application wiring — Manual DI and lifecycle
  - [ ] 2.1 Create `internal/app` composition root
    - Create `internal/app/app.go` with an `App` struct holding config, infrastructure references, and service instances
    - Implement `New(cfg *config.Config) (*App, error)` wiring dependencies explicitly: config → infra → repositories → services → handlers → routes
    - Pass all dependencies via constructor parameters; do not use a DI framework, global container, service locator, or reflection-based lookup
    - Do not create `pkg/container`, `GetInfra[T]()` access, or an `AddInfra()` registry
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [ ] 2.2 Implement App lifecycle (`Start` / `Shutdown`)
    - Implement `Start(ctx)` opening infrastructure in dependency order: database → cache → HTTP server → gRPC server
    - If any component fails to start, close already-opened components in reverse order and return the original error
    - Implement `Shutdown(ctx)` closing infrastructure in reverse order: gRPC → HTTP → cache → database
    - Return joined close errors where appropriate while respecting shutdown context deadlines
    - _Requirements: 2.3, 2.4, 2.7_

  - [ ]* 2.3 Write tests for App wiring and lifecycle
    - **Property 1: Constructor wiring visibility** — verify services and handlers receive dependencies through constructors, with no global service lookup
    - **Property 2: App lifecycle ordering** — verify `Start()` opens components in dependency order, rolls back started components on failure, and `Shutdown()` closes in reverse order
    - **Validates: Requirements 2.1, 2.2, 2.3, 2.4, 2.5, 2.7, 11.3**

- [ ] 3. Infrastructure components — SQL, Cache, API infra
  - [ ] 3.1 Implement SQL infra with GORM + OpenTelemetry plugin
    - Implement `pkg/infra/sql/` with explicit `Open(ctx) error` and `Close(ctx) error` methods
    - Configure connection pooling from config: MaxOpen (10), MaxIdle (5), MaxLifetime (5min)
    - Integrate `otelgorm` plugin for automatic query tracing
    - Receive database config via constructor injection from `app.New(cfg)`
    - _Requirements: 9.4, 9.5, 15.5_

  - [ ] 3.2 Implement Cache (Redis) infra
    - Implement `pkg/infra/cache/` with explicit `Open(ctx) error` and `Close(ctx) error` methods
    - Receive cache config via constructor injection from `app.New(cfg)` (addr, password, db)
    - Support health ping with 2-second timeout for readiness probe
    - _Requirements: 10.4_

  - [ ] 3.3 Implement API (Echo) infra with middleware stack
    - Implement `pkg/infra/api/` with explicit `Start(ctx) error` and `Shutdown(ctx) error` methods
    - Configure Echo with sonic JSON serializer
    - Apply middleware in order: RequestID → OTel → CORS → Recover → Logger → Gzip
    - Configure CORS from config, max body size (default 4MB), read/write timeouts (default 30s)
    - Implement custom Echo error handler mapping `AppError` → JSON response, untyped → 500
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.7, 7.8, 7.9, 5.3, 5.6_

  - [ ]* 3.4 Write property tests for Echo error handler
    - **Property 5: Error handler HTTP status mapping** — verify AppError maps to correct status/code/message, non-AppError maps to 500 with generic message
    - **Validates: Requirements 5.3, 5.6, 5.7**

- [ ] 4. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 5. Structured Logger package
  - [ ] 5.1 Create logger package with context-aware structured logging
    - Create `pkg/logger/logger.go` with `FromContext(ctx) *slog.Logger` extracting trace_id and span_id
    - Implement convenience methods: `Info(ctx, msg, attrs...)`, `Error(ctx, msg, attrs...)`, `Warn(...)`, `Debug(...)`
    - Implement `Setup(cfg)` selecting JSON handler (production) or text handler (development) based on APP_ENV
    - Implement `SetLevel()` for runtime level changes
    - Implement OTel log bridge integration when OTel is configured, graceful fallback when not
    - _Requirements: 18.1, 18.2, 18.4, 18.5, 18.6, 18.7_

  - [ ] 5.2 Create logger middleware for HTTP and gRPC
    - Create HTTP middleware injecting request-scoped logger into context with request_id, method, path
    - Create gRPC unary interceptor injecting grpc_method into context
    - Log request completion with method, path/method, status, latency
    - _Requirements: 18.3_

  - [ ]* 5.3 Write property test for logger trace correlation
    - **Property 19: Logger trace correlation** — verify context with valid span produces logger output containing trace_id and span_id
    - **Validates: Requirements 18.1**

- [ ] 6. Route Registry and Middleware
  - [ ] 6.1 Create declarative route registry (`pkg/route`)
    - Create `pkg/route/registry.go` with `Route`, `RouteGroup`, `Registry` structs
    - Implement `Register(e *echo.Echo) error` applying routes with duplicate detection
    - Support route grouping by prefix with per-group and per-route middleware
    - Support `Public` flag for JWT bypass
    - Log all registered routes at INFO level on startup
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7_

  - [ ] 6.2 Implement JWT middleware
    - Create `internal/middleware/jwt.go` using golang-jwt/jwt/v5
    - Skip authentication for routes marked `Public: true`
    - Extract claims and inject into context
    - Return 401 on invalid/expired token
    - _Requirements: 7.5, 13.4_

  - [ ] 6.3 Implement RBAC middleware
    - Create `internal/middleware/rbac.go` checking role permissions
    - Skip for public routes
    - Return 403 on insufficient permissions
    - _Requirements: 7.1, 7.5_

  - [ ]* 6.4 Write property tests for route registry
    - **Property 7: Route group path composition** — verify prefix + path produces correct full route
    - **Property 8: Duplicate route detection** — verify duplicate method+path returns error
    - **Property 9: Public route JWT bypass** — verify public routes don't require Authorization header
    - **Validates: Requirements 6.2, 6.7, 7.5**

- [ ] 7. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 8. Auth domain module — Entity, Repository, Service
  - [ ] 8.1 Create Auth domain entities and error codes
    - Create `internal/domain/auth/entity.go` with `User` struct (UUID, Email, Password, FullName, IsActive, RefreshToken, timestamps)
    - Create `internal/domain/auth/errors.go` with error code constants and constructors: `ErrInvalidCredentials`, `ErrEmailExists`, `ErrInvalidRefreshToken`, `ErrUserNotFound`
    - Create `internal/domain/auth/constants.go` with domain constants
    - _Requirements: 13.1, 13.6, 13.7, 13.8, 9.1_

  - [ ] 8.2 Create Auth repository interface and GORM Gen implementation
    - Create `internal/repository/auth_repository.go` with interface: `CreateUser`, `FindUserByEmail`, `FindUserByID`, `UpdateRefreshToken`, `WithTx`
    - Create `internal/repository/auth_repository_impl.go` implementing with GORM Gen queries
    - Map GORM errors to domain errors: `ErrRecordNotFound` → `ErrNotFound`, constraint violation → `ErrConflict`
    - Accept `context.Context` as first parameter in all methods
    - Receive DB via constructor injection from `internal/app`
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.8, 13.2_

  - [ ] 8.3 Create Auth service interface and struct implementation
    - Create `internal/service/auth_service.go` with interface: `Register`, `Login`, `RefreshToken`, `Logout`
    - Implement struct with constructor `NewAuthService(repo, hasher, jwt)` injecting dependencies as interfaces
    - Implement password hashing with bcrypt (configurable cost, default 10)
    - Implement JWT token generation/validation with configurable expiry (15min access, 7 days refresh)
    - Return domain errors for invalid credentials, duplicate email, expired refresh token
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 13.1, 13.4, 13.5, 13.6, 13.7, 13.8_

  - [ ]* 8.4 Write property tests for repository error mapping
    - **Property 3: Repository error mapping** — verify GORM errors map to correct domain error types
    - **Validates: Requirements 3.4, 3.5**

  - [ ]* 8.5 Write property tests for JWT and password hashing
    - **Property 12: JWT sign-verify round trip** — verify signing then verifying returns equivalent claims
    - **Property 13: Password hash-compare round trip** — verify hashing then comparing succeeds for any non-empty password
    - **Validates: Requirements 13.4, 13.5**

  - [ ]* 8.6 Write unit tests for Auth service
    - Test Register with valid input, duplicate email
    - Test Login with valid credentials, invalid credentials
    - Test RefreshToken with valid token, expired token
    - Test Logout
    - Use mockgen-generated mocks for repository, hasher, JWT manager
    - _Requirements: 4.7, 13.6, 13.7, 13.8_

- [ ] 9. Auth handler and route wiring
  - [ ] 9.1 Create Auth handler with request validation and swaggo annotations
    - Create `internal/handler/auth_handler.go` with `AuthHandler` struct depending on `AuthService` interface
    - Implement handlers: `Register`, `Login`, `RefreshToken`, `Logout`
    - Add request validation structs with go-playground/validator tags
    - Add swaggo annotations for OpenAPI generation
    - Bind, validate, call service, build response via response builder
    - _Requirements: 4.2, 13.3, 13.9, 12.4_

  - [ ] 9.2 Create route table and wire Auth module in `internal/app`
    - Create `internal/handler/routes.go` with `BuildRoutes()` function defining Auth route group
    - Define routes: POST /v1/auth/register (public), POST /v1/auth/login (public), POST /v1/auth/refresh (public), DELETE /v1/auth/logout (protected)
    - Wire everything in `app.New()`: build infra → build repos → build services → build handlers → register routes
    - _Requirements: 4.8, 6.1, 6.2, 6.3, 6.4, 13.3_

  - [ ] 9.3 Generate mocks with mockgen
    - Create `internal/mock/` directory
    - Generate mocks for `AuthRepository`, `AuthService`, `PasswordHasher`, `JWTManager` interfaces
    - Add `make generate` target for mock regeneration
    - _Requirements: 3.7_

- [ ] 10. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 11. Health probes and graceful shutdown
  - [ ] 11.1 Implement health and readiness endpoints
    - Create `GET /healthz` returning `{"status": "alive"}` with 200, excluded from JWT and request logging
    - Create `GET /readyz` checking database ping (2s timeout) and Redis ping (2s timeout)
    - Return 200 with `{"status": "ready", "checks": {...}}` when all pass
    - Return 503 with `{"status": "not_ready", "checks": {...}}` when any fail
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6_

  - [ ] 11.2 Implement graceful shutdown coordinator
    - Listen for SIGTERM/SIGINT, stop accepting new connections within 1 second
    - Wait for in-flight requests with configurable timeout (default 30s, range 5-120)
    - Close components in reverse dependency order: HTTP → workers → cache → database (max 10s per component)
    - Force-terminate on timeout (exit 1) or second signal (exit 1)
    - Log each component's shutdown result and duration in ms
    - Exit 0 on successful shutdown
    - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 11.6, 11.7_

- [ ] 12. OpenTelemetry observability setup
  - [ ] 12.1 Configure OpenTelemetry SDK with OTLP exporter
    - Set up OTel SDK with OTLP gRPC exporter for distributed tracing
    - Set resource attributes: service.name, service.version, deployment.environment from config
    - Propagate trace context through HTTP middleware, service methods, repository queries
    - Exclude /healthz and /readyz from trace collection
    - Integrate trace-log correlation (trace_id, span_id in slog attributes)
    - _Requirements: 15.1, 15.2, 15.3, 15.4, 15.6, 15.7_

- [ ] 13. gRPC transport layer
  - [ ] 13.1 Create proto files and generation setup
    - Create `proto/auth/v1/auth.proto` with AuthService RPC definitions (Register, Login, RefreshToken, Logout)
    - Create `proto/buf.gen.yaml` or Makefile `proto` target for code generation with protoc-gen-go and protoc-gen-go-grpc
    - Generate Go code from proto files
    - _Requirements: 17.3, 17.4_

  - [ ] 13.2 Implement gRPC server as explicitly wired infra component
    - Create `pkg/infra/grpc/grpc.go` with `Start(ctx) error` and `Shutdown(ctx) error`
    - Configure on separate port from HTTP (`grpc.addr` from config)
    - Add OTel unary interceptor, logging interceptor, recovery interceptor
    - Register standard `grpc.health.v1.Health` service
    - Support enable/disable via `server.grpcEnabled` config
    - _Requirements: 17.1, 17.2, 17.6, 17.7, 17.8, 17.10_

  - [ ] 13.3 Implement gRPC Auth handler reusing service layer
    - Create `internal/grpchandler/auth_handler.go` implementing generated AuthServiceServer
    - Reuse same `AuthService` interface as REST handlers
    - Implement `mapAppErrorToGRPCStatus()` mapping HTTP status → gRPC codes
    - _Requirements: 17.5, 17.9_

  - [ ]* 13.4 Write property test for gRPC error mapping
    - **Property 18: gRPC error mapping consistency** — verify HTTP status maps to canonical gRPC code (400→InvalidArgument, 401→Unauthenticated, etc.)
    - **Property 20: Transport independence** — verify same service call produces same result via REST and gRPC
    - **Validates: Requirements 17.5**

- [ ] 14. CLI commands and project scaffolding
  - [ ] 14.1 Implement Cobra CLI commands
    - Create `cmd/server.go` — start HTTP server with graceful shutdown, Overseer for zero-downtime restarts, SIGUSR1 for log level cycling
    - In `cmd/server.go`, call `config.Load()`, pass the result to `app.New(cfg)`, then call `app.Start(ctx)`
    - Create `cmd/migrate.go` — sub-commands: `up` (AutoMigrate), `down`, `status`; exit non-zero on failure with migration name
    - Create `cmd/seed.go` — populate example data for development
    - Create `cmd/genq.go` — run GORM Gen code generation
    - Create `cmd/docs.go` — generate OpenAPI spec via swaggo/swag
    - _Requirements: 14.1, 14.2, 14.3, 14.4, 14.5, 14.6, 14.7, 14.8_

  - [ ] 14.2 Create GORM Gen configuration and migration
    - Create GORM Gen config generating type-safe query objects to `internal/query/`
    - Create initial migration for Auth User table
    - _Requirements: 9.2, 9.3, 9.6, 9.7_

  - [ ] 14.3 Create Makefile with all targets
    - Targets: `build`, `test`, `lint`, `migrate`, `seed`, `generate` (GORM Gen + mockgen), `docs`, `proto`
    - `make build` produces single statically-linked Go binary
    - _Requirements: 1.4, 1.7, 17.4_

  - [ ] 14.4 Create Dockerfile and project scaffolding files
    - Create multi-stage Dockerfile: builder stage compiling Go binary, final Alpine stage with binary only
    - Create `.env.example` with all environment variables documented
    - Verify directory layout: `cmd/`, `internal/`, `pkg/`, `docs/`, `migrations/`, `proto/`
    - _Requirements: 1.1, 1.2, 1.3, 1.5, 1.6_

- [ ] 15. OpenAPI documentation
  - [ ] 15.1 Configure OpenAPI generation and Swagger UI
    - Ensure swaggo annotations on Auth handlers produce valid OpenAPI 3.0 YAML at `docs/openapi.yaml`
    - Include Bearer JWT security definitions for non-public routes
    - Serve Swagger UI at `/docs` in non-production, return 404 in production
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5, 12.6, 12.7_

- [ ] 16. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests use `pgregory.net/rapid` as specified in the design
- Unit tests use `testing` + `testify/assert` with mockgen-generated mocks
- The Auth module serves as the reference implementation demonstrating all patterns end-to-end
- gRPC transport is optional (disabled by default) but included to demonstrate dual-transport architecture

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "1.2", "1.3"] },
    { "id": 1, "tasks": ["1.4", "1.5", "1.6", "2.1"] },
    { "id": 2, "tasks": ["2.2"] },
    { "id": 3, "tasks": ["2.3", "3.1", "3.2", "3.3"] },
    { "id": 4, "tasks": ["3.4", "5.1"] },
    { "id": 5, "tasks": ["5.2", "5.3", "6.1"] },
    { "id": 6, "tasks": ["6.2", "6.3", "6.4"] },
    { "id": 7, "tasks": ["8.1"] },
    { "id": 8, "tasks": ["8.2", "8.3"] },
    { "id": 9, "tasks": ["8.4", "8.5", "8.6", "9.1", "9.3"] },
    { "id": 10, "tasks": ["9.2"] },
    { "id": 11, "tasks": ["11.1", "11.2", "12.1"] },
    { "id": 12, "tasks": ["13.1"] },
    { "id": 13, "tasks": ["13.2"] },
    { "id": 14, "tasks": ["13.3", "13.4"] },
    { "id": 15, "tasks": ["14.1", "14.2", "14.3", "14.4"] },
    { "id": 16, "tasks": ["15.1"] }
  ]
}
```
