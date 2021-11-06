package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromHouseholdInvitations(includeCounts bool, filteredCount uint64, webhooks ...*types.HouseholdInvitation) *sqlmock.Rows {
	columns := webhooksTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range webhooks {
		rowValues := []driver.Value{
			w.ID,
			w.CreatedOn,
			w.LastUpdatedOn,
			w.ArchivedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(webhooks))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowsFromHouseholdInvitations(includeCounts bool, filteredCount uint64, webhooks ...*types.HouseholdInvitation) *sqlmock.Rows {
	columns := webhooksTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range webhooks {
		rowValues := []driver.Value{
			w.ArchivedOn,
			w.ID,
			w.CreatedOn,
			w.LastUpdatedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(webhooks))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanHouseholdInvitations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanHouseholdInvitations(ctx, mockRows, false)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, _, err := q.scanHouseholdInvitations(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_HouseholdInvitationExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdID,
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(webhookExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdID, exampleHouseholdInvitationID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitationID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.HouseholdInvitationExists(ctx, "", exampleHouseholdInvitationID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdID,
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(webhookExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdID, exampleHouseholdInvitationID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdID,
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(webhookExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdID, exampleHouseholdInvitationID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdID, exampleHouseholdInvitation.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitation(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdID, exampleHouseholdInvitation.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetAllHouseholdInvitationsCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := uint64(123)
		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getAllHouseholdInvitationsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expected))

		actual, err := c.GetAllHouseholdInvitationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdInvitations(T *testing.T) {
	T.Parallel()

	exampleHouseholdID := fakes.BuildFakeID()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "webhooks", nil, nil, nil, "belongs_to_household", webhooksTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(
				true,
				exampleHouseholdInvitationList.FilteredCount,
				exampleHouseholdInvitationList.HouseholdInvitations...,
			))

		actual, err := c.GetHouseholdInvitations(ctx, exampleHouseholdID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		exampleHouseholdInvitationList.Page = 0
		exampleHouseholdInvitationList.Limit = 0
		filter := (*types.QueryFilter)(nil)

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "webhooks", nil, nil, nil, "belongs_to_household", webhooksTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(
				true,
				exampleHouseholdInvitationList.FilteredCount,
				exampleHouseholdInvitationList.HouseholdInvitations...,
			))

		actual, err := c.GetHouseholdInvitations(ctx, exampleHouseholdID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitations(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "webhooks", nil, nil, nil, "belongs_to_household", webhooksTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdInvitations(ctx, exampleHouseholdID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous database response", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "webhooks", nil, nil, nil, "belongs_to_household", webhooksTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdInvitations(ctx, exampleHouseholdID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleInput := fakes.BuildFakeHouseholdInvitationDatabaseCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		c.timeFunc = func() uint64 {
			return exampleHouseholdInvitation.CreatedOn
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHouseholdInvitation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing creation query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleInput := fakes.BuildFakeHouseholdInvitationDatabaseCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleHouseholdInvitation.CreatedOn
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}
