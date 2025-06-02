package indexing

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestHandleIndexRequest(T *testing.T) {
	T.Parallel()

	T.Run("user index type", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.UserDataManagerMock.On("GetUser", testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)
		dataManager.UserDataManagerMock.On("MarkUserAsIndexed", testutils.ContextMatcher, exampleUser.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleUser.ID,
			IndexType: IndexTypeUsers,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("deleting user index type", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		searchConfig := &textsearchcfg.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.UserDataManagerMock.On("GetUser", testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)
		dataManager.UserDataManagerMock.On("MarkUserAsIndexed", testutils.ContextMatcher, exampleUser.ID).Return(nil)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleUser.ID,
			IndexType: IndexTypeUsers,
			Delete:    true,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider(), searchConfig, dataManager, indexReq))
	})
}
