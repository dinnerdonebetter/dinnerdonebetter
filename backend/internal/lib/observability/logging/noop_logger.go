package logging

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// noopLogger is a default Logger we can provide that does nothing in case of dire emergencies.
type noopLogger struct{}

var logger = new(noopLogger)

// NewNoopLogger provides our noop Logger to dependency managers.
func NewNoopLogger() Logger { return logger }

// Info satisfies our interface.
func (l *noopLogger) Info(_ string) {}

// Debug satisfies our interface.
func (l *noopLogger) Debug(_ string) {}

// Warn satisfies our interface.
func (l *noopLogger) Warn(_ string) {}

// Error satisfies our interface.
func (l *noopLogger) Error(_ string, _ error) {}

// Fatal satisfies our interface.
func (l *noopLogger) Fatal(_ error) {}

// Printf satisfies our interface.
func (l *noopLogger) Printf(_ string, _ ...any) {}

// SetLevel satisfies our interface.
func (l *noopLogger) SetLevel(_ Level) {}

// SetRequestIDFunc satisfies our interface.
func (l *noopLogger) SetRequestIDFunc(_ RequestIDFunc) {}

// WithName satisfies our interface.
func (l *noopLogger) WithName(_ string) Logger { return l }

// Clone satisfies our interface.
func (l *noopLogger) Clone() Logger { return l }

// WithValues satisfies our interface.
func (l *noopLogger) WithValues(_ map[string]any) Logger { return l }

// WithValue satisfies our interface.
func (l *noopLogger) WithValue(_ string, _ any) Logger { return l }

// WithRequest satisfies our interface.
func (l *noopLogger) WithRequest(_ *http.Request) Logger { return l }

// WithResponse satisfies our interface.
func (l *noopLogger) WithResponse(_ *http.Response) Logger { return l }

// WithError satisfies our interface.
func (l *noopLogger) WithError(_ error) Logger { return l }

// WithSpan satisfies our interface.
func (l *noopLogger) WithSpan(_ trace.Span) Logger { return l }

func (l *noopLogger) WithContext(_ context.Context) Logger {
	return l
}
