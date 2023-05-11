package chi

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/attribute"
)

var doNotLog = map[string]struct{}{
	"/_meta_/ready": {}, // ready checks
	"/metrics":      {}, // metrics scrapes
}

// buildLoggingMiddleware builds a logging middleware.
func buildLoggingMiddleware(logger logging.Logger, silenceRouteLogging bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ww := chimiddleware.NewWrapResponseWriter(res, req.ProtoMajor)

			start := time.Now()
			next.ServeHTTP(ww, req)

			if !silenceRouteLogging {
				shouldLog := true
				for route := range doNotLog {
					if strings.HasPrefix(req.URL.Path, route) || req.URL.Path == route {
						shouldLog = false
						break
					}
				}

				if shouldLog {
					logger.WithRequest(req).WithValues(map[string]any{
						"status":  ww.Status(),
						"elapsed": time.Since(start).Milliseconds(),
						"written": ww.BytesWritten(),
					}).Debug("response served")
				}
			}
		})
	}
}

// buildTracingMiddleware builds a tracing middleware.
func buildTracingMiddleware(tracer tracing.Tracer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			spanName := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
			ctx, span := tracer.StartCustomSpan(req.Context(), spanName)
			defer span.End()

			tracing.AttachRequestToSpan(span, req)
			req = req.WithContext(ctx)

			ww := chimiddleware.NewWrapResponseWriter(res, req.ProtoMajor)
			next.ServeHTTP(ww, req)

			for k, v := range ww.Header() {
				span.SetAttributes(attribute.String(fmt.Sprintf("%s.%s", keys.ResponseHeadersKey, k), strings.Join(v, ";")))
			}

			span.SetAttributes(attribute.Int(keys.ResponseStatusKey, ww.Status()))
			span.SetAttributes(attribute.Int("response.bytes_written", ww.BytesWritten()))
		})
	}
}
