package observability

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

func TestPrepareError(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		descriptionFmt, descriptionArgs := "things and %s", "stuff"
		err := errors.New("blah")
		logger := logging.NewNoopLogger()
		_, span := tracing.StartSpan(ctx)

		assert.Error(t, PrepareError(err, logger, span, descriptionFmt, descriptionArgs))
	})
}

func TestAcknowledgeError(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		descriptionFmt, descriptionArgs := "things and %s", "stuff"
		err := errors.New("blah")
		logger := logging.NewNoopLogger()
		_, span := tracing.StartSpan(ctx)

		AcknowledgeError(err, logger, span, descriptionFmt, descriptionArgs)
	})
}
