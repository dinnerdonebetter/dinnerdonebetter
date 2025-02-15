package otelgrpc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	o11yutils "github.com/dinnerdonebetter/backend/internal/lib/observability/utils"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

// logger is our log wrapper.
type otelSlogLogger struct {
	requestIDFunc logging.RequestIDFunc
	logger        *slog.Logger
}

// NewOtelSlogLogger builds a new otelSlogLogger.
func NewOtelSlogLogger(ctx context.Context, lvl logging.Level, serviceName string, cfg *Config) (logging.Logger, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError("otel slog logger")
	}

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

	logHandlers := []slog.Handler{
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: lvl == logging.DebugLevel,
			Level:     level,
		}),
	}

	if cfg.CollectorEndpoint != "" {
		slog.Info("configuring otelgprc collector handler", slog.String("endpoint", cfg.CollectorEndpoint))

		options := []otlploggrpc.Option{
			otlploggrpc.WithEndpoint(cfg.CollectorEndpoint),
			otlploggrpc.WithTimeout(cfg.Timeout),
		}

		if cfg.Insecure {
			options = append(options, otlploggrpc.WithInsecure())
		}

		// Create the OTLP log exporter that sends logs to configured destination
		logExporter, err := otlploggrpc.New(ctx, options...)
		if err != nil {
			return nil, fmt.Errorf("instantiating otlploggrpc exporter: %w", err)
		}

		// Create the logger provider
		lp := log.NewLoggerProvider(
			log.WithProcessor(log.NewBatchProcessor(logExporter)),
			log.WithResource(o11yutils.MustOtelResource(ctx, serviceName)),
			log.WithAttributeCountLimit(128),
			log.WithAttributeValueLengthLimit(-1),
		)

		// Set the logger provider globally
		global.SetLoggerProvider(lp)

		logHandlers = append(logHandlers, otelslog.NewHandler(
			serviceName,
			otelslog.WithLoggerProvider(lp),
			otelslog.WithVersion("TODO_version"),
			otelslog.WithSource(true),
		))
	}

	logger := &otelSlogLogger{
		logger: slog.New(slogmulti.Fanout(logHandlers...)),
	}

	return logger, nil
}

// WithName is our obligatory contract fulfillment function.
// Slog doesn't support named loggers :( so we have this workaround.
func (l *otelSlogLogger) WithName(name string) logging.Logger {
	l2 := l.logger.With(slog.String(logging.LoggerNameKey, name))
	return &otelSlogLogger{logger: l2}
}

// SetRequestIDFunc sets the request ID retrieval function.
func (l *otelSlogLogger) SetRequestIDFunc(f logging.RequestIDFunc) {
	if f != nil {
		l.requestIDFunc = f
	}
}

// Info satisfies our contract for the logging.Logger Info method.
func (l *otelSlogLogger) Info(input string) {
	l.logger.Info(input)
}

// Debug satisfies our contract for the logging.Logger Debug method.
func (l *otelSlogLogger) Debug(input string) {
	l.logger.Debug(input)
}

// Error satisfies our contract for the logging.Logger Error method.
func (l *otelSlogLogger) Error(whatWasHappeningWhenErrorOccurred string, err error) {
	if err != nil {
		l.logger.Error(fmt.Sprintf("error %s: %s", whatWasHappeningWhenErrorOccurred, err.Error()))
	}
}

// Clone satisfies our contract for the logging.Logger WithValue method.
func (l *otelSlogLogger) Clone() logging.Logger {
	l2 := l.logger.With()
	return &otelSlogLogger{logger: l2}
}

// WithValue satisfies our contract for the logging.Logger WithValue method.
func (l *otelSlogLogger) WithValue(key string, value any) logging.Logger {
	l2 := l.logger.With(slog.Any(key, value))
	return &otelSlogLogger{logger: l2}
}

// WithValues satisfies our contract for the logging.Logger WithValues method.
func (l *otelSlogLogger) WithValues(values map[string]any) logging.Logger {
	var l2 = l.logger.With()

	for key, val := range values {
		l2 = l2.With(slog.Any(key, val))
	}

	return &otelSlogLogger{logger: l2}
}

// WithError satisfies our contract for the logging.Logger WithError method.
func (l *otelSlogLogger) WithError(err error) logging.Logger {
	l2 := l.logger.With(slog.Any("error", err))
	return &otelSlogLogger{logger: l2}
}

// WithSpan satisfies our contract for the logging.Logger WithSpan method.
func (l *otelSlogLogger) WithSpan(span trace.Span) logging.Logger {
	spanCtx := span.SpanContext()
	spanID := spanCtx.SpanID().String()
	traceID := spanCtx.TraceID().String()

	l2 := l.logger.With(slog.String(keys.SpanIDKey, spanID), slog.String(keys.TraceIDKey, traceID))

	return &otelSlogLogger{logger: l2}
}

func (l *otelSlogLogger) attachRequestToLog(req *http.Request) *slog.Logger {
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
func (l *otelSlogLogger) WithRequest(req *http.Request) logging.Logger {
	return &otelSlogLogger{logger: l.attachRequestToLog(req)}
}

// WithResponse satisfies our contract for the logging.Logger WithResponse method.
func (l *otelSlogLogger) WithResponse(res *http.Response) logging.Logger {
	l2 := l.logger.With()
	if res != nil {
		l2 = l.attachRequestToLog(res.Request).With(slog.Int(keys.ResponseStatusKey, res.StatusCode))
	}

	return &otelSlogLogger{logger: l2}
}
