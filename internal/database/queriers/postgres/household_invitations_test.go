package postgres

/*

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromHouseholdInvitations(includeCounts bool, filteredCount uint64, households ...*types.HouseholdInvitation) *sqlmock.Rows {
	columns := append(householdsTableColumns, householdsUserMembershipTableColumns...)

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range households {
		rowValues := []driver.Value{
			x.ID,
			x.Status,
			x.ToUser,
			x.FromUser,
			x.DestinationHousehold,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(households))
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

func TestQuerier_GetHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitation(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with no returned households", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		columns := append(householdsTableColumns, householdsUserMembershipTableColumns...)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows(columns))

		actual, err := c.GetHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetAllHouseholdInvitationsCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getAllHouseholdInvitationsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(exampleCount))

		actual, err := c.GetAllHouseholdInvitationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

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

func TestQuerier_GetHouseholdInvitations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdInvitationsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetHouseholdInvitations(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitations(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		exampleHouseholdInvitationList.Page, exampleHouseholdInvitationList.Limit = 0, 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdInvitationsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetHouseholdInvitations(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdInvitationsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdInvitations(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdInvitationsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdInvitations(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateHouseholdInvitation(T *testing.T) {
	T.SkipNow()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.ID = ""

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdInvitationCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		addUserToHouseholdInvitationArgs := []interface{}{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			strings.Join([]string{authorization.HouseholdAdminRole.String()}, householdMemberRolesSeparator),
		}

		db.ExpectExec(formatQueryForSQLMock(householdInvitationCreationQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdInvitationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleHouseholdInvitation.CreatedOn
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.NotEmpty(t, actual.ID)
		actual.ID = ""

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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.ID = ""

		exampleCreationInput := fakes.BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateHouseholdInvitation(ctx, exampleCreationInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.ID = ""

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHouseholdInvitation.CreatedOn
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error writing user addition to database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.ID = ""

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		addUserToHouseholdInvitationArgs := []interface{}{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			strings.Join([]string{authorization.HouseholdAdminRole.String()}, householdMemberRolesSeparator),
		}

		db.ExpectExec(formatQueryForSQLMock(householdInvitationCreationQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdInvitationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHouseholdInvitation.CreatedOn
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.ID = ""

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		addUserToHouseholdInvitationArgs := []interface{}{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			strings.Join([]string{authorization.HouseholdAdminRole.String()}, householdMemberRolesSeparator),
		}

		db.ExpectExec(formatQueryForSQLMock(householdInvitationCreationQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdInvitationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitation.ID))

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleHouseholdInvitation.CreatedOn
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdInvitationID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdInvitationID))

		assert.NoError(t, c.ArchiveHouseholdInvitation(ctx, exampleHouseholdInvitationID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHouseholdInvitation(ctx, "", exampleUserID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitationID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHouseholdInvitation(ctx, exampleHouseholdInvitationID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdInvitationID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveHouseholdInvitation(ctx, exampleHouseholdInvitationID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

*/
