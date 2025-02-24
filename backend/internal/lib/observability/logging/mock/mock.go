package mocklogging

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"
)

type (
	Logger struct {
		mock.Mock
	}
)

// Info implements the logging.Logger interface.
func (l *Logger) Info(s string) {
	l.Called(s)
}

// Debug implements the logging.Logger interface.
func (l *Logger) Debug(s string) {
	l.Called(s)
}

// Error implements the logging.Logger interface.
func (l *Logger) Error(whatWasHappeningWhenErrorOccurred string, err error) {
	l.Called(whatWasHappeningWhenErrorOccurred, err)
}

// SetRequestIDFunc implements the logging.Logger interface.
func (l *Logger) SetRequestIDFunc(idFunc logging.RequestIDFunc) {
	l.Called(idFunc)
}

// Clone implements the logging.Logger interface.
func (l *Logger) Clone() logging.Logger {
	return l.Called().Get(0).(logging.Logger)
}

// WithName implements the logging.Logger interface.
func (l *Logger) WithName(s string) logging.Logger {
	return l.Called(s).Get(0).(logging.Logger)
}

// WithValues implements the logging.Logger interface.
func (l *Logger) WithValues(m map[string]any) logging.Logger {
	return l.Called(m).Get(0).(logging.Logger)
}

// WithValue implements the logging.Logger interface.
func (l *Logger) WithValue(s string, a any) logging.Logger {
	return l.Called(s, a).Get(0).(logging.Logger)
}

// WithRequest implements the logging.Logger interface.
func (l *Logger) WithRequest(request *http.Request) logging.Logger {
	return l.Called(request).Get(0).(logging.Logger)
}

// WithResponse implements the logging.Logger interface.
func (l *Logger) WithResponse(response *http.Response) logging.Logger {
	return l.Called(response).Get(0).(logging.Logger)
}

// WithError implements the logging.Logger interface.
func (l *Logger) WithError(err error) logging.Logger {
	return l.Called(err).Get(0).(logging.Logger)
}

// WithSpan implements the logging.Logger interface.
func (l *Logger) WithSpan(span trace.Span) logging.Logger {
	return l.Called(span).Get(0).(logging.Logger)
}
