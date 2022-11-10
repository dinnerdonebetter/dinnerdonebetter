package observability

import (
	"fmt"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

// PrepareAndLogError standardizes our error handling by logging, tracing, and formatting an error consistently.
func PrepareAndLogError(err error, logger logging.Logger, span tracing.Span, descriptionFmt string, descriptionArgs ...any) error {
	if err == nil {
		return nil
	}

	desc := fmt.Sprintf(descriptionFmt, descriptionArgs...)
	if span != nil {
		tracing.AttachErrorToSpan(span, desc, err)
	}

	if logger != nil {
		logger.Error(err, desc)
	}

	return fmt.Errorf("%s: %w", desc, err)
}

// PrepareError standardizes our error handling by logging, tracing, and formatting an error consistently.
func PrepareError(err error, span tracing.Span, descriptionFmt string, descriptionArgs ...any) error {
	if err == nil {
		return nil
	}

	desc := fmt.Sprintf(descriptionFmt, descriptionArgs...)
	if span != nil {
		tracing.AttachErrorToSpan(span, desc, err)
	}

	return fmt.Errorf("%s: %w", desc, err)
}

// AcknowledgeError standardizes our error handling by logging and tracing consistently.
func AcknowledgeError(err error, logger logging.Logger, span tracing.Span, descriptionFmt string, descriptionArgs ...any) {
	desc := fmt.Sprintf(descriptionFmt, descriptionArgs...)
	logging.EnsureLogger(logger).Error(err, desc)
	if span != nil {
		tracing.AttachErrorToSpan(span, desc, err)
	}
}
