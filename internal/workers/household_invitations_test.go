package workers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postArchivesPublisher

		assert.NoError(t, worker.createHouseholdInvitation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postArchivesPublisher

		assert.Error(t, worker.createHouseholdInvitation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
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

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postArchivesPublisher

		assert.Error(t, worker.createHouseholdInvitation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
