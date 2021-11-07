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

func buildMockRowsFromHouseholdInvitations(includeCounts bool, filteredCount uint64, householdInvitations ...*types.HouseholdInvitation) *sqlmock.Rows {
	columns := householdInvitationsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range householdInvitations {
		rowValues := []driver.Value{
			w.ID,
			w.DestinationHousehold,
			w.ToEmail,
			w.ToUser,
			w.FromUser,
			w.Status,
			w.Note,
			w.StatusNote,
			w.Token,
			w.CreatedOn,
			w.LastUpdatedOn,
			w.ArchivedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(householdInvitations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowsFromHouseholdInvitations(includeCounts bool, filteredCount uint64, householdInvitations ...*types.HouseholdInvitation) *sqlmock.Rows {
	columns := householdInvitationsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range householdInvitations {
		rowValues := []driver.Value{
			w.ArchivedOn,
			w.ID,
			w.DestinationHousehold,
			w.ToEmail,
			w.ToUser,
			w.FromUser,
			w.Status,
			w.Note,
			w.StatusNote,
			w.Token,
			w.CreatedOn,
			w.LastUpdatedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(householdInvitations))
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
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(householdInvitationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdInvitationID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.HouseholdInvitationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(householdInvitationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdInvitationID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(householdInvitationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdInvitationID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdInvitationByHouseholdAndID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdID, exampleHouseholdInvitation.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByHouseholdAndIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, exampleHouseholdID, "")
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

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByHouseholdAndIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdInvitationByEmailAndToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdInvitation.ToEmail, exampleHouseholdInvitation.Token}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByEmailAndTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitationByEmailAndToken(ctx, exampleHouseholdInvitation.ToEmail, exampleHouseholdInvitation.Token)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByEmailAndToken(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByEmailAndToken(ctx, exampleHouseholdID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdInvitation.ToEmail, exampleHouseholdInvitation.Token}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByEmailAndTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitationByEmailAndToken(ctx, exampleHouseholdInvitation.ToEmail, exampleHouseholdInvitation.Token)
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

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getAllHouseholdInvitationsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetAllHouseholdInvitationsCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.StatusNote = ""
		exampleInput := fakes.BuildFakeHouseholdInvitationDatabaseCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.FromUser,
			exampleInput.ToUser,
			exampleInput.Note,
			exampleInput.ToEmail,
			exampleInput.Token,
			exampleInput.DestinationHousehold,
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
		exampleHouseholdInvitation.StatusNote = ""
		exampleInput := fakes.BuildFakeHouseholdInvitationDatabaseCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.FromUser,
			exampleInput.ToUser,
			exampleInput.Note,
			exampleInput.ToEmail,
			exampleInput.Token,
			exampleInput.DestinationHousehold,
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

func TestSQLQuerier_BuildGetPendingHouseholdInvitationsFromUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		expectedQuery := "SELECT household_invitations.id, household_invitations.destination_household, household_invitations.to_email, household_invitations.to_user, household_invitations.from_user, household_invitations.status, household_invitations.note, household_invitations.status_note, household_invitations.token, household_invitations.created_on, household_invitations.last_updated_on, household_invitations.archived_on FROM household_invitations WHERE household_invitations.archived_on IS NULL AND household_invitations.from_user = $1 AND household_invitations.status = $2 LIMIT 20"
		expectedArgs := []interface{}{
			userID,
			types.PendingHouseholdInvitationStatus,
		}

		c, _ := buildTestClient(t)

		actualQuery, actualArgs := c.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, userID, filter)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSQLQuerier_GetPendingHouseholdInvitationsFromUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetPendingHouseholdInvitationsFromUser(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetPendingHouseholdInvitationsFromUser(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetPendingHouseholdInvitationsFromUser(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_BuildGetPendingHouseholdInvitationsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		expectedQuery := "SELECT household_invitations.id, household_invitations.destination_household, household_invitations.to_email, household_invitations.to_user, household_invitations.from_user, household_invitations.status, household_invitations.note, household_invitations.status_note, household_invitations.token, household_invitations.created_on, household_invitations.last_updated_on, household_invitations.archived_on FROM household_invitations WHERE household_invitations.archived_on IS NULL AND household_invitations.status = $1 AND household_invitations.to_user = $2 LIMIT 20"
		expectedArgs := []interface{}{
			types.PendingHouseholdInvitationStatus,
			userID,
		}

		c, _ := buildTestClient(t)

		actualQuery, actualArgs := c.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, userID, filter)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSQLQuerier_GetPendingHouseholdInvitationsForUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetPendingHouseholdInvitationsForUser(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetPendingHouseholdInvitationsForUser(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetPendingHouseholdInvitationsForUser(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_setInvitationStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdInvitation.Status,
			exampleHouseholdInvitation.Note,
			exampleHouseholdID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		err := c.setInvitationStatus(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, _ := buildTestClient(t)

		err := c.setInvitationStatus(ctx, "", exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.Error(t, err)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, _ := buildTestClient(t)

		err := c.setInvitationStatus(ctx, exampleHouseholdID, "", exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.Error(t, err)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdInvitation.Status,
			exampleHouseholdInvitation.Note,
			exampleHouseholdID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		err := c.setInvitationStatus(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_CancelHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			types.CancelledHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		err := c.CancelHouseholdInvitation(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_AcceptHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			types.AcceptedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_RejectHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			types.RejectedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		err := c.RejectHouseholdInvitation(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
