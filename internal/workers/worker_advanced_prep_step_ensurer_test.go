package workers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

func TestProvideAdvancedPrepStepCreationEnsurerWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestAdvancedPrepStepCreationEnsurerWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideAdvancedPrepStepCreationEnsurerWorker(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)

		ctx := context.Background()
		assert.NoError(t, actual.HandleMessage(ctx, []byte("{}")))
	})
}
