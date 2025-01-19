package slog

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"go.opentelemetry.io/otel/trace"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

// logger is our log wrapper.
type slogLogger struct {
	requestIDFunc logging.RequestIDFunc
	logger        *slog.Logger
}

// NewSlogLogger builds a new slogLogger.
func NewSlogLogger(lvl logging.Level) logging.Logger {
	var level slog.Leveler
	switch lvl {
	case logging.DebugLevel:
		level = slog.LevelDebug
	case logging.InfoLevel:
		level = slog.LevelInfo
	case logging.WarnLevel:
		level = slog.LevelWarn
	case logging.ErrorLevel:
		level = slog.LevelError
	}

	handlerOptions := &slog.HandlerOptions{
		AddSource: lvl == logging.DebugLevel,
		Level:     level,
	}

	return &slogLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions)),
	}
}

// WithName is our obligatory contract fulfillment function.
// Slog doesn't support named loggers :( so we have this workaround.
func (l *slogLogger) WithName(name string) logging.Logger {
	l2 := l.logger.With(slog.String(logging.LoggerNameKey, name))
	return &slogLogger{logger: l2}
}

// SetRequestIDFunc sets the request ID retrieval function.
func (l *slogLogger) SetRequestIDFunc(f logging.RequestIDFunc) {
	if f != nil {
		l.requestIDFunc = f
	}
}

// Info satisfies our contract for the logging.Logger Info method.
func (l *slogLogger) Info(input string) {
	l.logger.Info(input)
}

// Debug satisfies our contract for the logging.Logger Debug method.
func (l *slogLogger) Debug(input string) {
	l.logger.Debug(input)
}

// Error satisfies our contract for the logging.Logger Error method.
func (l *slogLogger) Error(whatWasHappeningWhenErrorOccurred string, err error) {
	if err != nil {
		l.logger.Error(fmt.Sprintf("error %s: %s", whatWasHappeningWhenErrorOccurred, err.Error()))
	}
}

// Clone satisfies our contract for the logging.Logger WithValue method.
func (l *slogLogger) Clone() logging.Logger {
	l2 := l.logger.With()
	return &slogLogger{logger: l2}
}

// WithValue satisfies our contract for the logging.Logger WithValue method.
func (l *slogLogger) WithValue(key string, value any) logging.Logger {
	l2 := l.logger.With(slog.Any(key, value))
	return &slogLogger{logger: l2}
}

// WithValues satisfies our contract for the logging.Logger WithValues method.
func (l *slogLogger) WithValues(values map[string]any) logging.Logger {
	var l2 = l.logger.With()

	for key, val := range values {
		l2 = l2.With(slog.Any(key, val))
	}

	return &slogLogger{logger: l2}
}

// WithError satisfies our contract for the logging.Logger WithError method.
func (l *slogLogger) WithError(err error) logging.Logger {
	l2 := l.logger.With(slog.Any("error", err))
	return &slogLogger{logger: l2}
}

// WithSpan satisfies our contract for the logging.Logger WithSpan method.
func (l *slogLogger) WithSpan(span trace.Span) logging.Logger {
	spanCtx := span.SpanContext()
	spanID := spanCtx.SpanID().String()
	traceID := spanCtx.TraceID().String()

	l2 := l.logger.With(slog.String(keys.SpanIDKey, spanID), slog.String(keys.TraceIDKey, traceID))

	return &slogLogger{logger: l2}
}

func (l *slogLogger) attachRequestToLog(req *http.Request) *slog.Logger {
	if req != nil {
		l2 := l.logger.With(slog.String(keys.RequestMethodKey, req.Method))

		if req.URL != nil {
			l2 = l2.With(slog.String("path", req.URL.Path))
			if req.URL.RawQuery != "" {
				l2 = l2.With(slog.String(keys.URLQueryKey, req.URL.RawQuery))
			}
		}

		if l.requestIDFunc != nil {
			if reqID := l.requestIDFunc(req); reqID != "" {
				l2 = l2.With(slog.String(keys.RequestIDKey, reqID))
			}
		}

		return l2
	}

	return l.logger
}

// WithRequest satisfies our contract for the logging.Logger WithRequest method.
func (l *slogLogger) WithRequest(req *http.Request) logging.Logger {
	return &slogLogger{logger: l.attachRequestToLog(req)}
}

// WithResponse satisfies our contract for the logging.Logger WithResponse method.
func (l *slogLogger) WithResponse(res *http.Response) logging.Logger {
	l2 := l.logger.With()
	if res != nil {
		l2 = l.attachRequestToLog(res.Request).With(slog.Int(keys.ResponseStatusKey, res.StatusCode))
	}

	return &slogLogger{logger: l2}
}
