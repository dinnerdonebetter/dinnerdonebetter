package workers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:            types.HouseholdInvitationDataType,
			HouseholdInvitation: fakes.BuildFakeHouseholdInvitationDatabaseCreationInput(),
		}

		expectedHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			body.HouseholdInvitation,
		).Return(expectedHouseholdInvitation, nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
			&email.MockEmailer{},
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.createHouseholdInvitation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:            types.HouseholdInvitationDataType,
			HouseholdInvitation: fakes.BuildFakeHouseholdInvitationDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			body.HouseholdInvitation,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
			&email.MockEmailer{},
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createHouseholdInvitation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:            types.HouseholdInvitationDataType,
			HouseholdInvitation: fakes.BuildFakeHouseholdInvitationDatabaseCreationInput(),
		}

		expectedHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			body.HouseholdInvitation,
		).Return(expectedHouseholdInvitation, nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
			&email.MockEmailer{},
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createHouseholdInvitation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
