package workers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	"github.com/prixfixeco/api_server/pkg/types"
)

func TestProvidePreArchivesWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		actual, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing the first search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, errors.New("blah")
		}

		actual, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing the second search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}

		searchProviderInvocationCount := 0

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			searchProviderInvocationCount++
			if searchProviderInvocationCount == 2 {
				return nil, errors.New("blah")
			}
			return &mocksearch.IndexManager{}, nil
		}

		actual, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing the third search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}

		searchProviderInvocationCount := 0

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			searchProviderInvocationCount++
			if searchProviderInvocationCount == 3 {
				return nil, errors.New("blah")
			}
			return &mocksearch.IndexManager{}, nil
		}

		actual, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing the fourth search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}

		searchProviderInvocationCount := 0

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			searchProviderInvocationCount++
			if searchProviderInvocationCount == 4 {
				return nil, errors.New("blah")
			}
			return &mocksearch.IndexManager{}, nil
		}

		actual, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing the fifth search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}

		searchProviderInvocationCount := 0

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			searchProviderInvocationCount++
			if searchProviderInvocationCount == 5 {
				return nil, errors.New("blah")
			}
			return &mocksearch.IndexManager{}, nil
		}

		actual, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}

func TestArchivesWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := database.NewMockDatabase()
		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, []byte("} bad JSON lol")))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := database.NewMockDatabase()
		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		body := &types.PreArchiveMessage{
			DataType: types.UserMembershipDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
