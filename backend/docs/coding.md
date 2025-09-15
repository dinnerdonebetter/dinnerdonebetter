# Coding Standards

This is a list, in absolutely no particular order, of practices and guidelines I follow when writing code.

## Any struct that's capable of doing anything needs observability

**Why this matters:** Production systems fail in mysterious ways. When a user reports "the app is slow" or "my data disappeared," you need to trace the execution path to understand what happened. Without observability, you're flying blind. I've spent too many hours staring at stack traces that tell me "something failed" without any context about what the system was trying to do when it failed.

**The problem:** Traditional logging approaches lead to either too much noise (every function logs everything) or too little signal (only errors are logged, but you don't know the execution path that led to the error). You need structured, traceable information that connects business operations to technical execution.

**The solution:** Every struct that performs business logic gets a logger and tracer. This is achieved through the `platform/observability` packages, which provide implementations for tracers, loggers, and metrics collectors. The key insight is that observability isn't just about errors—it's about understanding the complete execution flow.

**Implementation pattern:** If you've completed implementation of business logic methods attached to a struct that is going to be a key part of the service, and you have not yet used the logger attached to the struct, then you haven't completed implementation of the business logic. The logger should be used to record important state changes, decision points, and the flow of data through your system.

There are helper methods in the base level observability package that allow for logging errors, adding the error details to the current span, and returning a wrapped error. While logging an error is not necessary and indeed can lead to noisy logs, marking the span is mandatory—this ensures that even if you don't log the error, it's still captured in your tracing system.

<details>
<summary><strong>For LLMs</strong></summary>

- Every struct with business logic needs `tracer` and `logger` fields
- Use `tracing.NewTracer()` and `logging.EnsureLogger().WithName()` in constructors
- Always start spans with `ctx, span := s.tracer.StartSpan(ctx)` and `defer span.End()`
- Use helper functions: `PrepareAndLogError()`, `PrepareError()`, `AcknowledgeError()`
- Mark spans for all errors, even if not logging them

</details> 

## Golang: Disable Unkeyed Literals

**Why this matters:** Go's struct literals allow both keyed and unkeyed initialization, but unkeyed literals are fragile and error-prone. When you reorder fields (which happens frequently with tools like `fieldalignment`), unkeyed literals break silently, leading to subtle bugs where data gets assigned to the wrong fields.

**The problem:** Consider this struct:

```go
type Person struct {
	Name string `json:"name"`
	Age int `json:"age"`
}
```

It can be initialized like so:

```go
x := &Person{"Gordon", 42}
```

If I later come in and want to change the order of the fields (perhaps to group related fields together or optimize memory layout), I break any of these initializations. The compiler won't catch this—it will just silently assign "Gordon" to the Age field and 42 to the Name field.

**The solution:** Add an unkeyed literal prevention field at the top:

```go
type Person struct {
	_ struct{} `json:"-"`
	
    Name string `json:"name"`
    Age int `json:"age"`
}
```

This tricks the compiler into disabling unkeyed fields in struct literals. Now the code above won't compile, forcing developers to use explicit field names:

```go
x := &Person{Name: "Gordon", Age: 42}
```

**Why this is especially important for me:** I like to use tools like `fieldalignment` to order struct fields automatically for memory optimization. Without this trick, every time the tool reorders fields, it would break existing unkeyed literals throughout the codebase. This pattern ensures that field reordering is safe and that all struct initializations are explicit and readable.

<details>
<summary><strong>For LLMs</strong></summary>

- Add `_ struct{} json:"-"` as first field in all structs
- Forces explicit field names in struct literals
- Prevents silent bugs when field order changes
- Essential when using field reordering tools

</details>

## Context Propagation

**Why this matters:** Context is Go's mechanism for carrying request-scoped data, cancellation signals, and timeouts through your call chain. Without proper context propagation, you lose the ability to cancel long-running operations, set timeouts, or access request-scoped data (like user information) deep in your call stack.

**The problem:** In a typical web request, you might need to:
- Cancel the entire operation if the client disconnects
- Set timeouts for database queries
- Access user information in a repository layer
- Trace the request through multiple service layers

Without context propagation, you end up either passing these values explicitly through every function (cluttering your APIs) or losing the ability to handle these scenarios entirely.

**The solution:** Context should always be the first parameter of any function that might need it. This includes any function that:
- Makes network calls
- Accesses the database
- Performs I/O operations
- Needs to be cancelled or have timeouts
- Needs to access request-scoped data (like session information)

**Pattern:**
```go
func (s *serviceImpl) SomeMethod(ctx context.Context, otherParams ...) (result, error) {
    // Always start with context as first parameter
    ctx, span := s.tracer.StartSpan(ctx)
    defer span.End()
    
    // Use context for all downstream calls
    result, err := s.repository.GetSomething(ctx, id)
    if err != nil {
        return nil, observability.PrepareError(err, span, "fetching something")
    }
    
    return result, nil
}
```

**Key Rules:**
- Context is always the first parameter
- Never store context in structs (except for testing)
- Always pass context through the call chain
- Use `context.WithValue()` to attach request-scoped data
- Use `context.WithTimeout()` and `context.WithCancel()` for timeouts and cancellation

<details>
<summary><strong>For LLMs</strong></summary>

- Context is always first parameter in functions that need it
- Never store context in structs (except testing)
- Pass context through entire call chain
- Use `context.WithValue()` for request-scoped data
- Use `context.WithTimeout()` and `context.WithCancel()` for timeouts

</details>

## Error Wrapping

**Why this matters:** When an error occurs in production, you need to understand not just what failed, but what the system was trying to do when it failed. Raw errors from libraries often lack context about the business operation that triggered them. Without proper error wrapping, you end up with stack traces that tell you "database connection failed" but not "failed while trying to fetch user profile for checkout process."

**The problem:** Consider this scenario: a user tries to place an order, but the database is down. Without error wrapping, you might see:
```
database: connection refused
```

This tells you the database is down, but not that it happened during order processing, which user was affected, or what the system was trying to accomplish.

**The solution:** Always wrap errors with descriptive context using `fmt.Errorf` with `%w` verb. This preserves the original error while adding meaningful context about what operation failed.

**Standard Error Wrapping Functions:**
- `observability.PrepareAndLogError()` - for errors that should be logged and returned
- `observability.PrepareError()` - for errors that should be traced but not logged (reduces noise)
- `observability.AcknowledgeError()` - for errors that should be acknowledged but not returned

**Pattern:**
```go
// ✅ DO THIS
user, err := s.userRepository.GetUserByID(ctx, userID)
if err != nil {
    return nil, observability.PrepareAndLogError(err, logger, span, "fetching user with ID %s", userID)
}

// ❌ DON'T DO THIS
user, err := s.userRepository.GetUserByID(ctx, userID)
if err != nil {
    return nil, err  // Lost context about what failed
}

// ❌ DON'T DO THIS
user, err := s.userRepository.GetUserByID(ctx, userID)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %v", err)  // Use %w, not %v
}
```

**Key Rules:**
- Always use `%w` verb with `fmt.Errorf` to preserve error chain
- Add descriptive context about what operation failed
- Include relevant parameters in the error message
- Use the observability helper functions for consistent error handling
- Don't wrap errors that are already wrapped with the same context

<details>
<summary><strong>For LLMs</strong></summary>

- Always use `fmt.Errorf` with `%w` verb to preserve error chain
- Use `observability.PrepareAndLogError()`, `PrepareError()`, `AcknowledgeError()`
- Add descriptive context about what operation failed
- Include relevant parameters in error messages
- Don't wrap already-wrapped errors with same context

</details>

## Rigid Code Generation

**Why this matters:** Hand-written configuration files and repetitive code are error-prone and inconsistent. When you have 50+ SQL queries, each with slight variations, it's easy to make typos, forget to update one query when you change a column name, or introduce subtle inconsistencies. Generated code eliminates these problems by making the generator the single source of truth.

**The problem:** Consider maintaining SQL queries manually. You have a `users` table with columns `id`, `name`, `email`, `created_at`. You write 20 queries that reference these columns. Then you rename `email` to `email_address`. You have to remember to update all 20 queries, and if you miss one, you get a runtime error. Worse, if you have a typo in a column name, you might not catch it until production.

**The solution:** I strongly prefer code-driven file generation over hand-written configuration files. This approach has several key advantages:

**Type Safety**: Generated code is validated at compile time. If the generator has a bug, the build fails. If a generated file has a bug, it's because the generator has a bug, not because someone made a typo.

**Consistency**: All generated files follow the same patterns. You can't accidentally use `user_id` in one query and `userId` in another.

**Maintainability**: Changes to the generator automatically update all generated files. Rename a column in the generator, and every query gets updated.

**Confidence**: If the generator code is correct, all generated files are correct. You can trust that your generated code is consistent and error-free.

**Examples in this codebase:**

**SQL Queries**: All database queries are generated from Go code in `cmd/tools/codegen/queries/`. I declare table columns in arrays at the top of each file, then programmatically generate every query. This means:
- If I rename a column, I update it in one place and every query gets updated
- If I add a new column, I can easily add it to all relevant queries
- Typos in column names are caught at build time
- I use sqlc, so the generated queries then have type-safe Go functions generated around them

**Service Configs**: Environment-specific configurations are generated in `cmd/tools/codegen/configs/`. This ensures that if the config struct is encodable, then the generated configs are guaranteed to be decodable. No more "config file doesn't match struct" errors.

**Pattern:**
```go
// Generator code (cmd/tools/codegen/queries/main.go)
var usersColumns = []string{
    idColumn,
    nameColumn,
    emailColumn,
    createdAtColumn,
}

func buildUsersQueries(database string) []*Query {
    return []*Query{
        {
            Annotation: QueryAnnotation{
                Name: "GetUserByID",
                Type: OneType,
            },
            Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s FROM %s WHERE %s = sqlc.arg(%s)`,
                strings.Join(usersColumns, ", "),
                usersTableName,
                idColumn, idColumn,
            )),
        },
    }
}
```

**Key Rules:**
- Write generators in Go, not shell scripts or other languages
- Use strong typing for generator inputs and outputs
- Validate generated code at build time
- Make generators idempotent (running multiple times produces same output)
- Include the generator source in version control
- Document the generation process clearly
- Use `go generate` directives where appropriate

<details>
<summary><strong>For LLMs</strong></summary>

- Prefer code-driven generation over hand-written files
- Write generators in Go with strong typing
- Validate generated code at build time
- Make generators idempotent
- Use single source of truth for patterns
- Generate SQL queries from column arrays
- Generate configs to ensure encodable/decodable consistency

</details>

## Golang Unit Tests


