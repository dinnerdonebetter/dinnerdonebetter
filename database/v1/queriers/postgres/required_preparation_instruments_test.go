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

func buildMockRowsFromRequiredPreparationInstrument(requiredPreparationInstruments ...*models.RequiredPreparationInstrument) *sqlmock.Rows {
	includeCount := len(requiredPreparationInstruments) > 1
	columns := requiredPreparationInstrumentsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range requiredPreparationInstruments {
		rowValues := []driver.Value{
			x.ID,
			x.ValidInstrumentID,
			x.Notes,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToValidPreparation,
		}

		if includeCount {
			rowValues = append(rowValues, len(requiredPreparationInstruments))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRequiredPreparationInstrument(x *models.RequiredPreparationInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(requiredPreparationInstrumentsTableColumns).AddRow(
		x.ArchivedOn,
		x.ValidInstrumentID,
		x.Notes,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToValidPreparation,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanRequiredPreparationInstruments(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanRequiredPreparationInstruments(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRequiredPreparationInstrumentExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		expectedQuery := "SELECT EXISTS ( SELECT required_preparation_instruments.id FROM required_preparation_instruments JOIN valid_preparations ON required_preparation_instruments.belongs_to_valid_preparation=valid_preparations.id WHERE required_preparation_instruments.belongs_to_valid_preparation = $1 AND required_preparation_instruments.id = $2 AND valid_preparations.id = $3 )"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
			exampleRequiredPreparationInstrument.ID,
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := p.buildRequiredPreparationInstrumentExistsQuery(exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RequiredPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT required_preparation_instruments.id FROM required_preparation_instruments JOIN valid_preparations ON required_preparation_instruments.belongs_to_valid_preparation=valid_preparations.id WHERE required_preparation_instruments.belongs_to_valid_preparation = $1 AND required_preparation_instruments.id = $2 AND valid_preparations.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleRequiredPreparationInstrument.ID,
				exampleValidPreparation.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RequiredPreparationInstrumentExists(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleRequiredPreparationInstrument.ID,
				exampleValidPreparation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RequiredPreparationInstrumentExists(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.valid_instrument_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.updated_on, required_preparation_instruments.archived_on, required_preparation_instruments.belongs_to_valid_preparation FROM required_preparation_instruments JOIN valid_preparations ON required_preparation_instruments.belongs_to_valid_preparation=valid_preparations.id WHERE required_preparation_instruments.belongs_to_valid_preparation = $1 AND required_preparation_instruments.id = $2 AND valid_preparations.id = $3"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
			exampleRequiredPreparationInstrument.ID,
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := p.buildGetRequiredPreparationInstrumentQuery(exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.valid_instrument_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.updated_on, required_preparation_instruments.archived_on, required_preparation_instruments.belongs_to_valid_preparation FROM required_preparation_instruments JOIN valid_preparations ON required_preparation_instruments.belongs_to_valid_preparation=valid_preparations.id WHERE required_preparation_instruments.belongs_to_valid_preparation = $1 AND required_preparation_instruments.id = $2 AND valid_preparations.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleRequiredPreparationInstrument.ID,
				exampleValidPreparation.ID,
			).
			WillReturnRows(buildMockRowsFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument))

		actual, err := p.GetRequiredPreparationInstrument(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleRequiredPreparationInstrument.ID,
				exampleValidPreparation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRequiredPreparationInstrument(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRequiredPreparationInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL"
		actualQuery := p.buildGetAllRequiredPreparationInstrumentsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRequiredPreparationInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRequiredPreparationInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRequiredPreparationInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.valid_instrument_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.updated_on, required_preparation_instruments.archived_on, required_preparation_instruments.belongs_to_valid_preparation, (SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL) FROM required_preparation_instruments JOIN valid_preparations ON required_preparation_instruments.belongs_to_valid_preparation=valid_preparations.id WHERE required_preparation_instruments.archived_on IS NULL AND required_preparation_instruments.belongs_to_valid_preparation = $1 AND valid_preparations.id = $2 AND required_preparation_instruments.created_on > $3 AND required_preparation_instruments.created_on < $4 AND required_preparation_instruments.updated_on > $5 AND required_preparation_instruments.updated_on < $6 ORDER BY required_preparation_instruments.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
			exampleValidPreparation.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRequiredPreparationInstrumentsQuery(exampleValidPreparation.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.valid_instrument_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.updated_on, required_preparation_instruments.archived_on, required_preparation_instruments.belongs_to_valid_preparation, (SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL) FROM required_preparation_instruments JOIN valid_preparations ON required_preparation_instruments.belongs_to_valid_preparation=valid_preparations.id WHERE required_preparation_instruments.archived_on IS NULL AND required_preparation_instruments.belongs_to_valid_preparation = $1 AND valid_preparations.id = $2 ORDER BY required_preparation_instruments.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleValidPreparation.ID,
			).
			WillReturnRows(
				buildMockRowsFromRequiredPreparationInstrument(
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[0],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[1],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[2],
				),
			)

		actual, err := p.GetRequiredPreparationInstruments(ctx, exampleValidPreparation.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrumentList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleValidPreparation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRequiredPreparationInstruments(ctx, exampleValidPreparation.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleValidPreparation.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRequiredPreparationInstruments(ctx, exampleValidPreparation.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning required preparation instrument", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleValidPreparation.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument))

		actual, err := p.GetRequiredPreparationInstruments(ctx, exampleValidPreparation.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		expectedQuery := "INSERT INTO required_preparation_instruments (valid_instrument_id,notes,belongs_to_valid_preparation) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRequiredPreparationInstrument.ValidInstrumentID,
			exampleRequiredPreparationInstrument.Notes,
			exampleRequiredPreparationInstrument.BelongsToValidPreparation,
		}
		actualQuery, actualArgs := p.buildCreateRequiredPreparationInstrumentQuery(exampleRequiredPreparationInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO required_preparation_instruments (valid_instrument_id,notes,belongs_to_valid_preparation) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRequiredPreparationInstrument.ID, exampleRequiredPreparationInstrument.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ValidInstrumentID,
				exampleRequiredPreparationInstrument.Notes,
				exampleRequiredPreparationInstrument.BelongsToValidPreparation,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRequiredPreparationInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ValidInstrumentID,
				exampleRequiredPreparationInstrument.Notes,
				exampleRequiredPreparationInstrument.BelongsToValidPreparation,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRequiredPreparationInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		expectedQuery := "UPDATE required_preparation_instruments SET valid_instrument_id = $1, notes = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_valid_preparation = $3 AND id = $4 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleRequiredPreparationInstrument.ValidInstrumentID,
			exampleRequiredPreparationInstrument.Notes,
			exampleRequiredPreparationInstrument.BelongsToValidPreparation,
			exampleRequiredPreparationInstrument.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRequiredPreparationInstrumentQuery(exampleRequiredPreparationInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE required_preparation_instruments SET valid_instrument_id = $1, notes = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_valid_preparation = $3 AND id = $4 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ValidInstrumentID,
				exampleRequiredPreparationInstrument.Notes,
				exampleRequiredPreparationInstrument.BelongsToValidPreparation,
				exampleRequiredPreparationInstrument.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ValidInstrumentID,
				exampleRequiredPreparationInstrument.Notes,
				exampleRequiredPreparationInstrument.BelongsToValidPreparation,
				exampleRequiredPreparationInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		expectedQuery := "UPDATE required_preparation_instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_valid_preparation = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
			exampleRequiredPreparationInstrument.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRequiredPreparationInstrumentQuery(exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE required_preparation_instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_valid_preparation = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleRequiredPreparationInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRequiredPreparationInstrument(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleRequiredPreparationInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRequiredPreparationInstrument(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
				exampleRequiredPreparationInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRequiredPreparationInstrument(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
