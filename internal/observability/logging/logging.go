package logging

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

const (
	// LoggerNameKey is a key we can use to denote zerologLogger names across implementations.
	LoggerNameKey = "_name_"
)

type (
	level int

	// Level is a simple string alias for dependency injection's sake.
	Level *level

	// RequestIDFunc fetches a string ID from a request.
	RequestIDFunc func(*http.Request) string
)

var (
	// InfoLevel describes a info-level log.
	InfoLevel Level = new(level)
	// DebugLevel describes a debug-level log.
	DebugLevel Level = new(level)
	// ErrorLevel describes a error-level log.
	ErrorLevel Level = new(level)
	// WarnLevel describes a warn-level log.
	WarnLevel Level = new(level)
)

// Logger represents a simple logging interface we can build wrappers around.
// Note that someone, naive and green, may be enticed to add a method to this interface akin to:
// WithQueryFilter(*types.QueryFilter) Logger
// This is a fool's errand, which would introduce a disallowed import cycle.
type Logger interface {
	Info(string)
	Debug(string)
	Error(err error, whatWasHappeningWhenErrorOccurred string)
	Printf(string, ...any)

	SetLevel(Level)
	SetRequestIDFunc(RequestIDFunc)

	Clone() Logger
	WithName(string) Logger
	WithValues(map[string]any) Logger
	WithValue(string, any) Logger
	WithRequest(*http.Request) Logger
	WithResponse(response *http.Response) Logger
	WithError(error) Logger
	WithSpan(span trace.Span) Logger
}

// EnsureLogger guarantees that a Logger is available.
func EnsureLogger(logger Logger) Logger {
	if logger != nil {
		return logger
	}

	return NewNoopLogger()
}
