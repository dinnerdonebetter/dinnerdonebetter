package manager

import (
	"context"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/audit/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/audit/fakes"
	auditmock "github.com/dinnerdonebetter/backend/internal/domain/audit/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildAuditManagerForTest(t *testing.T) (*auditManager, *auditmock.Repository) {
	t.Helper()

	repo := &auditmock.Repository{}
	m := NewAuditDataManager(tracing.NewNoopTracerProvider(), logging.NewNoopLogger(), repo)
	return m.(*auditManager), repo
}

func TestAuditDataManager_GetAuditLogEntry(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		manager, repo := buildAuditManagerForTest(t)

		expected := fakes.BuildFakeAuditLogEntry()
		repo.On(reflection.GetMethodName(repo.GetAuditLogEntry), testutils.ContextMatcher, expected.ID).Return(expected, nil)

		result, err := manager.GetAuditLogEntry(ctx, expected.ID)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mock.AssertExpectationsForObjects(t, repo)
	})
}

func TestAuditDataManager_GetAuditLogEntriesForAccount(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		manager, repo := buildAuditManagerForTest(t)

		accountID := fakes.BuildFakeID()
		filter := filtering.DefaultQueryFilter()
		expected := fakes.BuildFakeAuditLogEntriesList()
		repo.On(reflection.GetMethodName(repo.GetAuditLogEntriesForAccount), testutils.ContextMatcher, accountID, filter).Return(expected, nil)

		result, err := manager.GetAuditLogEntriesForAccount(ctx, accountID, filter)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mock.AssertExpectationsForObjects(t, repo)
	})
}

func TestAuditDataManager_CreateAuditLogEntry(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		manager, repo := buildAuditManagerForTest(t)

		exampleEntry := fakes.BuildFakeAuditLogEntry()
		dbInput := converters.ConvertAuditLogEntryToAuditLogEntryDatabaseCreationInput(exampleEntry)
		querier := &database.MockQueryExecutor{}

		repo.On(reflection.GetMethodName(repo.CreateAuditLogEntry), testutils.ContextMatcher, mock.Anything, mock.MatchedBy(func(in *types.AuditLogEntryDatabaseCreationInput) bool {
			return in.ID == dbInput.ID && in.BelongsToUser == dbInput.BelongsToUser
		})).Return(exampleEntry, nil)

		created, err := manager.CreateAuditLogEntry(ctx, querier, dbInput)

		require.NoError(t, err)
		assert.NotNil(t, created)
		assert.Equal(t, exampleEntry.ID, created.ID)
		mock.AssertExpectationsForObjects(t, repo)
	})
}
