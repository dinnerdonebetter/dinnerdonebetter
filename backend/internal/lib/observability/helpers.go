package observability

import (
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

func ObserveValue(span tracing.Span, logger logging.Logger, key string, value any) logging.Logger {
	tracing.AttachToSpan(span, key, value)
	return logger.WithValue(key, value)
}
