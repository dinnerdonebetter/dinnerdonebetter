package chi

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/routing"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

var doNotLog = map[string]struct{}{
	"/_meta_/ready": {}, // ready checks
	"/metrics":      {}, // metrics scrapes
}

// buildLoggingMiddleware builds a logging middleware.
func buildLoggingMiddleware(logger logging.Logger, cfg *routing.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ww := chimiddleware.NewWrapResponseWriter(res, req.ProtoMajor)

			start := time.Now()
			next.ServeHTTP(ww, req)

			if !cfg.SilenceRouteLogging {
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

var doNotTrace = map[string]struct{}{
	"/_meta_/ready": {}, // health checks
}

// buildTracingMiddleware builds a tracing middleware.
func buildTracingMiddleware(tracer tracing.Tracer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			shouldTrace := true
			for route := range doNotTrace {
				if strings.HasPrefix(req.URL.Path, route) || req.URL.Path == route {
					shouldTrace = false
					break
				}
			}

			if shouldTrace {
				ctx, span := tracer.StartCustomSpan(req.Context(), fmt.Sprintf("%s %s", req.Method, req.URL.Path))
				defer span.End()

				tracing.AttachRequestToSpan(span, req)
				req = req.WithContext(ctx)
			}

			ww := chimiddleware.NewWrapResponseWriter(res, req.ProtoMajor)

			next.ServeHTTP(ww, req)
		})
	}
}
