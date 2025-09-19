# Writing Go 

## Testing Conventions

### Subtest Structure

**Always use subtests, even for single test cases.** This provides a consistent structure that makes it easy to add additional test scenarios later.

#### Parameter Naming Convention
- **Main test function**: Use capital `T` for the `*testing.T` parameter
- **Subtest functions**: Use lowercase `t` for the `*testing.T` parameter

```go
func TestMyFunction(T *testing.T) {
    T.Parallel()

    T.Run("standard", func(t *testing.T) {
        t.Parallel()
        // test implementation
    })

    T.Run("with error condition", func(t *testing.T) {
        t.Parallel()
        // error test implementation
    })
}
```

#### Happy Path First
Structure your subtests with the happy path (successful operation) as the first subtest. This pattern makes it easy to add derivative tests that simulate failure conditions:
- Database connection failures
- Queue timeouts
- Invalid input validation
- Permission errors

### Table Tests - Use Sparingly

**Avoid using table tests.** Only use them when there's meaningful code savings or when testing multiple similar scenarios with different inputs. Simple iterations over a slice (like testing multiple algorithms) are often clearer than formal table test structures.

**Prefer explicit subtests when:**
- Each test case has significantly different setup
- Error conditions require different assertions
- The test logic differs meaningfully between cases

**Consider table tests when:**
- Testing the same function with many different input/output pairs
- Validating input parsing with multiple formats
- Testing mathematical operations across ranges

### Parallel Execution

**Always use parallel test execution** where possible:
- Add `T.Parallel()` to main test functions
- Add `t.Parallel()` to subtests
- This significantly improves test suite performance

### Test Helpers

Mark test helper functions with `t.Helper()` to ensure proper error reporting:

```go
func checkEquality(t *testing.T, expected, actual *MyStruct) {
    t.Helper()
    
    assert.NotZero(t, actual.ID)
    assert.Equal(t, expected.Name, actual.Name)
    // ... other assertions
}
```

### Test Organization

- **Unit tests**: Place alongside the code they test (`*_test.go` files)
- **Integration tests**: Located in `tests/integration/` directory
- **Test utilities**: Shared helpers in `internal/platform/testutils/`

### Context Usage

**Always use context in tests** that interact with external services:

```go
func TestMyService(T *testing.T) {
    T.Run("standard", func(t *testing.T) {
        ctx := t.Context()
        // use ctx in service calls
    })
}
```

### Test Reliability

- **No skipped tests**: The linter explicitly forbids `t.SkipNow()` calls
- **No global state**: Each test should be independent
- **Clean up resources**: Properly tear down any created test data

## General Go Conventions

### Error Handling

This codebase follows strict error handling practices:
- Always check errors (enforced by linter configuration)
- Wrap external package errors for context
- Use meaningful error messages

### Code Organization

#### Directory Structure
```
internal/
├── domain/          # Business logic and types
├── platform/        # Infrastructure and utilities
├── repositories/    # Data access layer
└── services/        # Application services
```

#### Package Dependencies
- Domain packages should not import platform packages
- Keep dependencies flowing downward in the architecture
- Use dependency injection (wire.go files)

### Linting and Code Quality

The project uses `golangci-lint` with custom configuration:
- **Strict error checking**: No ignored errors
- **Import organization**: Use `gci` for import sorting
- **Test-specific rules**: Relaxed linting for test files where appropriate
- **Complexity limits**: Enforced cyclomatic complexity bounds

### Performance Considerations

- **Parallel processing**: Leverage goroutines where appropriate
- **Connection pooling**: Reuse database and HTTP connections
- **Circuit breakers**: Implement fault tolerance for external services
- **Caching**: Strategic use of in-memory and distributed caching

## Development Workflow

### Before Submitting Code

1. **Run the full test suite**: `make test`
2. **Check linting**: `make lint`
3. **Verify coverage**: `make coverage`
4. **Run integration tests**: `make integration_tests`

### Code Generation

This project uses extensive code generation:
- **Database queries**: `make querier` (using sqlc)
- **API clients**: `make codegen`
- **Configuration structs**: `make configs`

Always regenerate code after schema changes.

## Naming Conventions

### Interfaces
- Use descriptive suffixes that indicate purpose:
  - `Logger` - for logging abstractions
  - `Manager` - for business logic coordinators  
  - `Repository` - for data access abstractions
  - `Router` - for routing abstractions
  - `Handler` - for request handlers

### Structs
- **Configuration structs**: End with `Config` (e.g., `APIServiceConfig`, `DatabaseConfig`)
- **Request/Response types**: Use descriptive names with purpose (e.g., `APIResponse[T]`, `UpdateRequestInput`)
- **Domain entities**: Use clear business terms (e.g., `MealPlan`, `Recipe`, `ValidIngredient`)

### Constants and Variables
- Use fully descriptive names: `ConfigurationFilePathEnvVarKey`
- Group related constants in `const` blocks
- Use `var` blocks for package-level variables with initialization

### Functions and Methods
- **Constructors**: Use `New` prefix (e.g., `NewService`, `NewGenerator`)
- **Wire providers**: Use `Provide` prefix with descriptive names (e.g., `ProvideDatabaseFromConfig`)
- **Converters**: Use `Convert` prefix describing the transformation

## Struct Design Patterns

### Privacy and Safety
Always include an unexported struct field to prevent accidental struct construction and comparison:

```go
type MyStruct struct {
    _ struct{} `json:"-"`
    
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### Type Definitions
Group related types in type blocks with clear documentation:

```go
type (
    // Level is a simple string alias for dependency injection's sake.
    Level *level
    
    // RequestIDFunc fetches a string ID from a request.
    RequestIDFunc func(*http.Request) string
)
```

### Generics Usage
Leverage Go generics for type safety and reusability:

```go
// APIResponse represents a response we might send to the user.
APIResponse[T any] struct {
    _ struct{} `json:"-"`
    
    Data       T                     `json:"data,omitempty"`
    Pagination *filtering.Pagination `json:"pagination,omitempty"`
    Details    ResponseDetails       `json:"details"`
}

// Range represents a min/max range constraint.
Range[T comparable] struct {
    Min T `json:"min"`
    Max T `json:"max"`
}
```

## Dependency Injection with Wire

### Wire File Structure
Each package that participates in dependency injection should have a `wire.go` file:

```go
package mypackage

import "github.com/google/wire"

var (
    // Providers represents this package's offering to the dependency injector.
    Providers = wire.NewSet(
        NewService,
        ProvideConfigFromEnvironment,
        wire.FieldsOf(new(*Config), "Database", "Logger"),
    )
)
```

### Provider Functions
Wire provider functions should be clearly named and documented:

```go
// ProvideHTTPServerConfigFromAPIServiceConfig extracts HTTP config from API service config.
func ProvideHTTPServerConfigFromAPIServiceConfig(cfg *APIServiceConfig) http.Config {
    return cfg.HTTPServer
}
```

### Wire Sets Organization
- **Package-level providers**: Group all providers a package offers
- **Descriptive variable names**: Use `Providers`, `ConfigProviders`, `ServiceProviders` etc.
- **Logical grouping**: Group related provider functions together

## Common Patterns

### Service Initialization

Services typically follow this pattern:
```go
type Service struct {
    _ struct{} `json:"-"`
    
    logger logger.Logger
    tracer tracing.Tracer
    // ... other dependencies
}

func NewService(
    logger logger.Logger,
    tracer tracing.Tracer,
) *Service {
    return &Service{
        logger: logger,
        tracer: tracer,
    }
}
```

### Interface Design
Keep interfaces focused and purposeful:

```go
// Logger represents a simple logging interface we can build wrappers around.
type Logger interface {
    Info(string)
    Debug(string)
    Error(whatWasHappeningWhenErrorOccurred string, err error)
    
    WithName(string) Logger
    WithValues(map[string]any) Logger
    WithSpan(span trace.Span) Logger
}
```

### Request/Response Handling

- Use structured logging with context
- Implement proper HTTP status codes
- Validate input at service boundaries
- Transform domain objects to/from API representations
- Use consistent response wrapper types (`APIResponse[T]`)

## Project Structure Patterns

### Directory Organization
The codebase follows a clean architecture approach:

```
internal/
├── domain/          # Business logic, entities, and domain services
│   ├── mealplanning/    # Meal planning domain
│   ├── identity/        # User and account management
│   └── auth/           # Authentication domain
├── platform/        # Infrastructure and shared utilities
│   ├── database/       # Database abstractions and implementations
│   ├── observability/ # Logging, tracing, metrics
│   ├── messagequeue/  # Queue abstractions and implementations
│   └── uploads/       # File upload handling
├── repositories/    # Data access implementations
│   └── postgres/      # PostgreSQL-specific implementations
└── services/        # Application services and handlers
    ├── auth/          # Authentication service
    ├── mealplanning/  # Meal planning service
    └── identity/      # Identity management service
```

### Package Import Rules
- **Domain packages** should not import platform packages
- **Services** can import both domain and platform packages
- **Repositories** implement domain interfaces but can use platform utilities
- Keep dependencies flowing downward in the architecture

### Generated Code
Several directories contain generated code:
- `internal/grpc/generated/` - gRPC service definitions
- Database query files (via sqlc)
- Configuration structs (via custom codegen)

Never edit generated files directly; modify the generators instead.

## Getting Started

### For New Contributors

1. **Read the existing tests** in your area of focus to understand patterns
2. **Follow the subtest structure** even for simple tests  
3. **Use the linter** to catch common issues early
4. **Study the wire.go files** to understand dependency relationships
5. **Look at similar implementations** before writing new code
6. **Run the full test suite** before submitting changes

### Development Checklist

Before submitting a pull request:
- [ ] Tests follow the subtest conventions
- [ ] Wire providers are properly defined if needed
- [ ] Structs include privacy protection (`_ struct{}`)
- [ ] Error handling follows project patterns
- [ ] Linting passes (`make lint`)
- [ ] Tests pass (`make test`)
- [ ] Integration tests pass if relevant (`make integration_tests`)

### When to Ask Questions

- If you see patterns that seem inconsistent
- When choosing between multiple valid approaches
- Before making architectural changes
- When adding new dependencies or technologies

## Conclusion

These conventions represent years of learned experience maintaining a large Go codebase. They prioritize:

- **Consistency** over individual preference
- **Testability** over brevity
- **Maintainability** over cleverness
- **Explicitness** over implicit behavior

When in doubt, follow the existing patterns you see in similar code. The goal is to write code that the next developer (including future you) can easily understand, test, and modify.

Remember: Code is read far more often than it's written. These conventions optimize for the reading experience.

