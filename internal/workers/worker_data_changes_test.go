package workers

import (
	"context"
	"testing"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/email"
)

func TestProvideDataChangesWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideDataChangesWorker(
			zerolog.NewZerologLogger(),
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestDataChangesWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideDataChangesWorker(
			zerolog.NewZerologLogger(),
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)

		ctx := context.Background()
		assert.NoError(t, actual.HandleMessage(ctx, []byte("{}")))
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()

		actual := ProvideDataChangesWorker(
			zerolog.NewZerologLogger(),
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)

		ctx := context.Background()
		assert.Error(t, actual.HandleMessage(ctx, []byte("} bad JSON lol")))
	})
}
