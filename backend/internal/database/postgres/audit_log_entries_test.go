package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	auditLogEntriesCreatedForUsersByDefault = 3
)

func createAuditLogEntryForTest(t *testing.T, ctx context.Context, querier database.SQLQueryExecutor, exampleAuditLogEntry *types.AuditLogEntry, user *types.User, household *types.Household, dbc *Querier) *types.AuditLogEntry {
	t.Helper()

	if user == nil {
		user = createUserForTest(t, ctx, nil, dbc)
	}

	if household == nil {
		household = createHouseholdForTest(t, ctx, nil, dbc)
	}

	// create
	if exampleAuditLogEntry == nil {
		exampleAuditLogEntry = fakes.BuildFakeAuditLogEntry()
	}
	exampleAuditLogEntry.BelongsToUser = user.ID
	exampleAuditLogEntry.BelongsToHousehold = &household.ID
	dbInput := converters.ConvertAuditLogEntryToAuditLogEntryDatabaseCreationInput(exampleAuditLogEntry)

	created, err := dbc.createAuditLogEntry(ctx, querier, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleAuditLogEntry.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleAuditLogEntry, created)

	auditLogEntry, err := dbc.GetAuditLogEntry(ctx, created.ID)
	exampleAuditLogEntry.CreatedAt = auditLogEntry.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, auditLogEntry, exampleAuditLogEntry)

	return created
}

func TestQuerier_Integration_AuditLogEntries(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	household := createHouseholdForTest(t, ctx, nil, dbc)

	exampleAuditLogEntry := fakes.BuildFakeAuditLogEntry()
	exampleAuditLogEntry.BelongsToHousehold = &household.ID
	exampleAuditLogEntry.BelongsToUser = user.ID
	createdAuditLogEntries := []*types.AuditLogEntry{}

	// create
	createdAuditLogEntries = append(createdAuditLogEntries, createAuditLogEntryForTest(t, ctx, dbc.db, exampleAuditLogEntry, user, household, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeAuditLogEntry()
		createdAuditLogEntries = append(createdAuditLogEntries, createAuditLogEntryForTest(t, ctx, dbc.db, input, user, household, dbc))
	}

	// fetch as list
	auditLogEntries, err := dbc.GetAuditLogEntriesForUser(ctx, user.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, auditLogEntries.Data)
	assert.Equal(t, len(createdAuditLogEntries)+auditLogEntriesCreatedForUsersByDefault, len(auditLogEntries.Data))

	auditLogEntries, err = dbc.GetAuditLogEntriesForHousehold(ctx, household.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, auditLogEntries.Data)
	assert.Equal(t, len(createdAuditLogEntries)+auditLogEntriesCreatedForUsersByDefault-1, len(auditLogEntries.Data))
}

func TestQuerier_GetAuditLogEntry(T *testing.T) {
	T.Parallel()

	T.Run("with invalid audit log entry ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetAuditLogEntry(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetAuditLogEntryThatNeedSearchIndexing(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, db := buildTestClient(t)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateAuditLogEntry(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.createAuditLogEntry(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}
