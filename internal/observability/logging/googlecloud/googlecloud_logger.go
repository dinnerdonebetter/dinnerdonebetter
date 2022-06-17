package googlecloud

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"

	gcplogging "cloud.google.com/go/logging"
	"go.opentelemetry.io/otel/trace"
)

// logger is our log wrapper.
type gcpLogger struct {
	requestIDFunc logging.RequestIDFunc
	loggerClient  *gcplogging.Client
	logger        *gcplogging.Logger
}

// NewLogger builds a new gcpLogger.
func NewLogger(ctx context.Context) (logging.Logger, error) {
	client, err := gcplogging.NewClient(ctx, fmt.Sprintf("projects/%s", os.Getenv("PF_ENVIRONMENT")))
	if err != nil {
		return nil, err
	}

	l := &gcpLogger{
		loggerClient: client,
		logger:       client.Logger(""),
	}

	return l, nil
}

func (l *gcpLogger) Info(s string) {
	l.logger.Log(gcplogging.Entry{
		Severity: gcplogging.Info,
		Payload:  s,
		Labels:   nil,
	})
}

func (l *gcpLogger) Debug(s string) {
	l.logger.Log(gcplogging.Entry{
		Severity: gcplogging.Debug,
		Payload:  s,
		Labels:   nil,
	})
}

func (l *gcpLogger) Error(err error, whatWasHappeningWhenErrorOccurred string) {
	l.logger.Log(gcplogging.Entry{
		Severity: gcplogging.Error,
		Payload:  fmt.Sprintf("%s: %s", whatWasHappeningWhenErrorOccurred, err.Error()),
		Labels:   nil,
	})
}

func (l *gcpLogger) Fatal(err error) {
	l.logger.Log(gcplogging.Entry{
		Severity: gcplogging.Emergency,
		Payload:  err.Error(),
		Labels:   nil,
	})
}

func (l *gcpLogger) Printf(s string, i ...interface{}) {
	l.logger.Log(gcplogging.Entry{
		Severity: gcplogging.Default,
		Payload:  fmt.Sprintf(s, i...),
		Labels:   nil,
	})
}

func (l *gcpLogger) SetLevel(level logging.Level) {
	//
}

func (l *gcpLogger) SetRequestIDFunc(idFunc logging.RequestIDFunc) {
	l.requestIDFunc = idFunc
}

func (l *gcpLogger) Clone() logging.Logger {
	return &gcpLogger{
		requestIDFunc: l.requestIDFunc,
		loggerClient:  l.loggerClient,
		logger:        l.logger,
	}
}

func (l *gcpLogger) WithName(s string) logging.Logger {
	l.logger = l.loggerClient.Logger(s)
	return l
}

func (l *gcpLogger) WithValues(m map[string]interface{}) logging.Logger {
	return newEntryLogger(l).WithValues(m)
}

func (l *gcpLogger) WithValue(s string, i interface{}) logging.Logger {
	return newEntryLogger(l).WithValue(s, i)
}

func (l *gcpLogger) WithRequest(request *http.Request) logging.Logger {
	return newEntryLogger(l).WithRequest(request)
}

func (l *gcpLogger) WithResponse(response *http.Response) logging.Logger {
	return newEntryLogger(l).WithResponse(response)
}

func (l *gcpLogger) WithError(err error) logging.Logger {
	return newEntryLogger(l).WithError(err)
}

func (l *gcpLogger) WithSpan(span trace.Span) logging.Logger {
	return newEntryLogger(l).WithSpan(span)
}

// logger is our log wrapper.
type entryLogger struct {
	logger        *gcplogging.Logger
	entry         gcplogging.Entry
	requestIDFunc logging.RequestIDFunc
}

func newEntryLogger(l *gcpLogger) logging.Logger {
	return &entryLogger{
		logger: l.logger,
		entry: gcplogging.Entry{
			Severity: gcplogging.Default,
			Labels:   map[string]string{},
		},
	}
}

func (l *entryLogger) Info(s string) {
	l.entry.Severity = gcplogging.Info
	l.logger.Log(l.entry)
}

func (l *entryLogger) Debug(s string) {
	l.entry.Severity = gcplogging.Debug
	l.logger.Log(l.entry)
}

func (l *entryLogger) Error(err error, whatWasHappeningWhenErrorOccurred string) {
	l.entry.Severity = gcplogging.Error
	l.entry.Payload = fmt.Sprintf("%s: %s", whatWasHappeningWhenErrorOccurred, err.Error())

	l.logger.Log(l.entry)
}

func (l *entryLogger) Fatal(err error) {
	l.entry.Severity = gcplogging.Emergency
	l.logger.Log(l.entry)
}

func (l *entryLogger) Printf(s string, i ...interface{}) {
	l.logger.Log(l.entry)
}

func (l *entryLogger) SetLevel(level logging.Level) {
	var lvl gcplogging.Severity

	switch level {
	case logging.InfoLevel:
		lvl = gcplogging.Info
	case logging.DebugLevel:
		lvl = gcplogging.Debug
	case logging.WarnLevel:
		lvl = gcplogging.Warning
	case logging.ErrorLevel:
		lvl = gcplogging.Error
	default:
		lvl = gcplogging.Info
	}

	l.entry.Severity = lvl
}

func (l *entryLogger) SetRequestIDFunc(idFunc logging.RequestIDFunc) {
	l.requestIDFunc = idFunc
}

func (l *entryLogger) Clone() logging.Logger {
	return &entryLogger{
		logger:        l.logger,
		entry:         l.entry,
		requestIDFunc: l.requestIDFunc,
	}
}

func (l *entryLogger) WithName(s string) logging.Logger {
	return l
}

func (l *entryLogger) WithValue(s string, i interface{}) logging.Logger {
	l.entry.Labels[s] = fmt.Sprintf("%v", i)

	return l
}

func (l *entryLogger) WithValues(m map[string]interface{}) logging.Logger {
	for k, v := range m {
		l.entry.Labels[k] = fmt.Sprintf("%v", v)
	}

	return l
}

func (l *entryLogger) WithRequest(req *http.Request) logging.Logger {
	l.entry.Labels[keys.RequestMethodKey] = req.Method

	if req.URL != nil {
		l.entry.Labels[keys.RequestURIPathKey] = req.URL.Path
		if req.URL.RawQuery != "" {
			l.entry.Labels[keys.RequestURIQueryKey] = req.URL.RawQuery
		}
	}

	if l.requestIDFunc != nil {
		if reqID := l.requestIDFunc(req); reqID != "" {
			l.entry.Labels[keys.RequestURIQueryKey] = req.URL.RawQuery
		}
	}

	return l
}

func (l *entryLogger) WithResponse(response *http.Response) logging.Logger {
	l.entry.Labels[keys.RequestMethodKey] = response.Request.Method

	return l
}

func (l *entryLogger) WithError(err error) logging.Logger {
	l.entry.Labels["error"] = err.Error()

	return l
}

func (l *entryLogger) WithSpan(span trace.Span) logging.Logger {
	spanCtx := span.SpanContext()
	spanID := spanCtx.SpanID().String()
	traceID := spanCtx.TraceID().String()

	l.entry.SpanID = spanID
	l.entry.Trace = traceID

	return l
}
