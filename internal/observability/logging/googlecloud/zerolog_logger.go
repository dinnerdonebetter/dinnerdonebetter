package googlecloud

// import (
// 	"net/http"

// 	"github.com/prixfixeco/api_server/internal/observability/keys"
// 	"github.com/prixfixeco/api_server/internal/observability/logging"

// 	"cloud.google.com/go/logging"
// 	"go.opentelemetry.io/otel/trace"
// )

// const here = "github.com/prixfixeco/api_server/"

// // logger is our log wrapper.
// type gcpLogger struct {
// 	requestIDFunc logging.RequestIDFunc
// }

// // NewZerologLogger builds a new gcpLogger.
// func NewZerologLogger() logging.Logger {
// 	return &gcpLogger{}
// }

// // WithName is our obligatory contract fulfillment function.
// // Zerolog doesn't support named loggers :( so we have this workaround.
// func (l *gcpLogger) WithName(name string) logging.Logger {

// }

// // SetLevel sets the log level for our zerolog logger.
// func (l *gcpLogger) SetLevel(level logging.Level) {

// }

// // SetRequestIDFunc sets the request ID retrieval function.
// func (l *gcpLogger) SetRequestIDFunc(f logging.RequestIDFunc) {
// 	if f != nil {
// 		l.requestIDFunc = f
// 	}
// }

// // Info satisfies our contract for the logging.Logger Info method.
// func (l *gcpLogger) Info(input string) {

// }

// // Debug satisfies our contract for the logging.Logger Debug method.
// func (l *gcpLogger) Debug(input string) {

// }

// // Error satisfies our contract for the logging.Logger Error method.
// func (l *gcpLogger) Error(err error, input string) {
// 	if err != nil {

// 	}
// }

// // Fatal satisfies our contract for the logging.Logger Fatal method.
// func (l *gcpLogger) Fatal(err error) {

// }

// // Printf satisfies our contract for the logging.Logger Printf method.
// func (l *gcpLogger) Printf(format string, args ...interface{}) {

// }

// // Clone satisfies our contract for the logging.Logger WithValue method.
// func (l *gcpLogger) Clone() logging.Logger {

// }

// // WithValue satisfies our contract for the logging.Logger WithValue method.
// func (l *gcpLogger) WithValue(key string, value interface{}) logging.Logger {

// }

// // WithValues satisfies our contract for the logging.Logger WithValues method.
// func (l *gcpLogger) WithValues(values map[string]interface{}) logging.Logger {

// }

// // WithError satisfies our contract for the logging.Logger WithError method.
// func (l *gcpLogger) WithError(err error) logging.Logger {

// }

// // WithSpan satisfies our contract for the logging.Logger WithSpan method.
// func (l *gcpLogger) WithSpan(span trace.Span) logging.Logger {
// 	spanCtx := span.SpanContext()
// 	spanID := spanCtx.SpanID().String()
// 	traceID := spanCtx.TraceID().String()

// 	l = l.WithValue(keys.SpanIDKey, spanID).WithValue(keys.TraceIDKey, traceID)

// 	return l
// }

// // WithRequest satisfies our contract for the logging.Logger WithRequest method.
// func (l *gcpLogger) WithRequest(req *http.Request) logging.Logger {
// 	return &gcpLogger{}
// }

// // WithResponse satisfies our contract for the logging.Logger WithResponse method.
// func (l *gcpLogger) WithResponse(res *http.Response) logging.Logger {
// 	return &gcpLogger{}
// }
