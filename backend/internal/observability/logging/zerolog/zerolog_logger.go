package zerolog

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

const here = "github.com/dinnerdonebetter/backend/"

func init() {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		panic(err)
	}

	zerolog.CallerSkipFrameCount += 2
	zerolog.DisableSampling(true)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(loc)
	}
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return strings.TrimPrefix(file, here) + ", line " + strconv.Itoa(line)
	}
	zerolog.LevelFieldName = "severity"
}

// logger is our log wrapper.
type zerologLogger struct {
	requestIDFunc logging.RequestIDFunc
	logger        zerolog.Logger
}

// buildZerologger builds a new zerologger.
func buildZerologger(level logging.Level) zerolog.Logger {
	var lvl zerolog.Level
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	switch level {
	case logging.InfoLevel:
		lvl = zerolog.InfoLevel
	case logging.DebugLevel:
		logger = logger.With().Logger()
		lvl = zerolog.DebugLevel
	case logging.WarnLevel:
		logger = logger.With().Caller().Logger()
		lvl = zerolog.WarnLevel
	case logging.ErrorLevel:
		logger = logger.With().Caller().Logger()
		lvl = zerolog.ErrorLevel
	default:
		lvl = zerolog.InfoLevel
	}

	return logger.Level(lvl)
}

// NewZerologLogger builds a new zerologLogger.
func NewZerologLogger(lvl logging.Level) logging.Logger {
	return &zerologLogger{logger: buildZerologger(lvl)}
}

// WithName is our obligatory contract fulfillment function.
// Zerolog doesn't support named loggers :( so we have this workaround.
func (l *zerologLogger) WithName(name string) logging.Logger {
	l2 := l.logger.With().Str(logging.LoggerNameKey, name).Logger()
	return &zerologLogger{logger: l2}
}

// SetRequestIDFunc sets the request ID retrieval function.
func (l *zerologLogger) SetRequestIDFunc(f logging.RequestIDFunc) {
	if f != nil {
		l.requestIDFunc = f
	}
}

// Info satisfies our contract for the logging.Logger Info method.
func (l *zerologLogger) Info(input string) {
	l.logger.Info().Msg(input)
}

// Debug satisfies our contract for the logging.Logger Debug method.
func (l *zerologLogger) Debug(input string) {
	l.logger.Debug().Msg(input)
}

// Error satisfies our contract for the logging.Logger Error method.
func (l *zerologLogger) Error(err error, whatWasHappeningWhenErrorOccurred string) {
	if err != nil {
		l.logger.Error().Stack().Caller().Msg(fmt.Sprintf("error %s: %s", whatWasHappeningWhenErrorOccurred, err.Error()))
	}
}

// Clone satisfies our contract for the logging.Logger WithValue method.
func (l *zerologLogger) Clone() logging.Logger {
	l2 := l.logger.With().Logger()
	return &zerologLogger{logger: l2}
}

// WithValue satisfies our contract for the logging.Logger WithValue method.
func (l *zerologLogger) WithValue(key string, value any) logging.Logger {
	l2 := l.logger.With().Interface(key, value).Logger()
	return &zerologLogger{logger: l2}
}

// WithValues satisfies our contract for the logging.Logger WithValues method.
func (l *zerologLogger) WithValues(values map[string]any) logging.Logger {
	var l2 = l.logger.With().Logger()

	for key, val := range values {
		l2 = l2.With().Interface(key, val).Logger()
	}

	return &zerologLogger{logger: l2}
}

// WithError satisfies our contract for the logging.Logger WithError method.
func (l *zerologLogger) WithError(err error) logging.Logger {
	l2 := l.logger.With().Err(err).Logger()
	return &zerologLogger{logger: l2}
}

// WithSpan satisfies our contract for the logging.Logger WithSpan method.
func (l *zerologLogger) WithSpan(span trace.Span) logging.Logger {
	spanCtx := span.SpanContext()
	spanID := spanCtx.SpanID().String()
	traceID := spanCtx.TraceID().String()

	l2 := l.logger.With().Str(keys.SpanIDKey, spanID).Str(keys.TraceIDKey, traceID).Logger()

	return &zerologLogger{logger: l2}
}

func (l *zerologLogger) attachRequestToLog(req *http.Request) zerolog.Logger {
	if req != nil {
		l2 := l.logger.With().
			Str("method", req.Method).
			Logger()

		if req.URL != nil {
			l2 = l2.With().Str("path", req.URL.Path).Logger()
			if req.URL.RawQuery != "" {
				l2 = l2.With().Str(keys.URLQueryKey, req.URL.RawQuery).Logger()
			}
		}

		if l.requestIDFunc != nil {
			if reqID := l.requestIDFunc(req); reqID != "" {
				l2 = l2.With().Str("request.id", reqID).Logger()
			}
		}

		return l2
	}

	return l.logger
}

// WithRequest satisfies our contract for the logging.Logger WithRequest method.
func (l *zerologLogger) WithRequest(req *http.Request) logging.Logger {
	return &zerologLogger{logger: l.attachRequestToLog(req)}
}

// WithResponse satisfies our contract for the logging.Logger WithResponse method.
func (l *zerologLogger) WithResponse(res *http.Response) logging.Logger {
	l2 := l.logger.With().Logger()
	if res != nil {
		l2 = l.attachRequestToLog(res.Request).With().Int(keys.ResponseStatusKey, res.StatusCode).Logger()
	}

	return &zerologLogger{logger: l2}
}
