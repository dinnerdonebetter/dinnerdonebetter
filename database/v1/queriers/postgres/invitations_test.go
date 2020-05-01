package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromInvitation(invitations ...*models.Invitation) *sqlmock.Rows {
	includeCount := len(invitations) > 1
	columns := invitationsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range invitations {
		rowValues := []driver.Value{
			x.ID,
			x.Code,
			x.Consumed,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		if includeCount {
			rowValues = append(rowValues, len(invitations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromInvitation(x *models.Invitation) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(invitationsTableColumns).AddRow(
		x.ArchivedOn,
		x.Code,
		x.Consumed,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToUser,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanInvitations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanInvitations(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanInvitations(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildInvitationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT invitations.id FROM invitations WHERE invitations.id = $1 )"
		expectedArgs := []interface{}{
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := p.buildInvitationExistsQuery(exampleInvitation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_InvitationExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT invitations.id FROM invitations WHERE invitations.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleInvitation.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.InvitationExists(ctx, exampleInvitation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleInvitation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.InvitationExists(ctx, exampleInvitation.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT invitations.id, invitations.code, invitations.consumed, invitations.created_on, invitations.updated_on, invitations.archived_on, invitations.belongs_to_user FROM invitations WHERE invitations.id = $1"
		expectedArgs := []interface{}{
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := p.buildGetInvitationQuery(exampleInvitation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetInvitation(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT invitations.id, invitations.code, invitations.consumed, invitations.created_on, invitations.updated_on, invitations.archived_on, invitations.belongs_to_user FROM invitations WHERE invitations.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleInvitation.ID,
			).
			WillReturnRows(buildMockRowsFromInvitation(exampleInvitation))

		actual, err := p.GetInvitation(ctx, exampleInvitation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleInvitation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleInvitation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetInvitation(ctx, exampleInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllInvitationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(invitations.id) FROM invitations WHERE invitations.archived_on IS NULL"
		actualQuery := p.buildGetAllInvitationsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllInvitationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(invitations.id) FROM invitations WHERE invitations.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllInvitationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetInvitationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT invitations.id, invitations.code, invitations.consumed, invitations.created_on, invitations.updated_on, invitations.archived_on, invitations.belongs_to_user, COUNT(invitations.id) FROM invitations WHERE invitations.archived_on IS NULL AND invitations.created_on > $1 AND invitations.created_on < $2 AND invitations.updated_on > $3 AND invitations.updated_on < $4 GROUP BY invitations.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetInvitationsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetInvitations(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT invitations.id, invitations.code, invitations.consumed, invitations.created_on, invitations.updated_on, invitations.archived_on, invitations.belongs_to_user, COUNT(invitations.id) FROM invitations WHERE invitations.archived_on IS NULL GROUP BY invitations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleInvitationList := fakemodels.BuildFakeInvitationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(
				buildMockRowsFromInvitation(
					&exampleInvitationList.Invitations[0],
					&exampleInvitationList.Invitations[1],
					&exampleInvitationList.Invitations[2],
				),
			)

		actual, err := p.GetInvitations(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleInvitationList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetInvitations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		actual, err := p.GetInvitations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning invitation", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(
				buildErroneousMockRowFromInvitation(exampleInvitation),
			)

		actual, err := p.GetInvitations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO invitations (code,consumed,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleInvitation.Code,
			exampleInvitation.Consumed,
			exampleInvitation.BelongsToUser,
		}
		actualQuery, actualArgs := p.buildCreateInvitationQuery(exampleInvitation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateInvitation(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO invitations (code,consumed,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleInvitation.ID, exampleInvitation.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleInvitation.Code,
				exampleInvitation.Consumed,
				exampleInvitation.BelongsToUser,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateInvitation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleInvitation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleInvitation.Code,
				exampleInvitation.Consumed,
				exampleInvitation.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateInvitation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE invitations SET code = $1, consumed = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleInvitation.Code,
			exampleInvitation.Consumed,
			exampleInvitation.BelongsToUser,
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := p.buildUpdateInvitationQuery(exampleInvitation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateInvitation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE invitations SET code = $1, consumed = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleInvitation.Code,
				exampleInvitation.Consumed,
				exampleInvitation.BelongsToUser,
				exampleInvitation.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateInvitation(ctx, exampleInvitation)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleInvitation.Code,
				exampleInvitation.Consumed,
				exampleInvitation.BelongsToUser,
				exampleInvitation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateInvitation(ctx, exampleInvitation)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE invitations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := p.buildArchiveInvitationQuery(exampleInvitation.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveInvitation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE invitations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleInvitation.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveInvitation(ctx, exampleInvitation.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleInvitation.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveInvitation(ctx, exampleInvitation.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleInvitation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveInvitation(ctx, exampleInvitation.ID, exampleUser.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
