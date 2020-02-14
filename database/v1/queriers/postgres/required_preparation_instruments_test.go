package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowFromRequiredPreparationInstrument(x *models.RequiredPreparationInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(requiredPreparationInstrumentsTableColumns).AddRow(
		x.ID,
		x.InstrumentID,
		x.PreparationID,
		x.Notes,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
	)

	return exampleRows
}

func buildErroneousMockRowFromRequiredPreparationInstrument(x *models.RequiredPreparationInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(requiredPreparationInstrumentsTableColumns).AddRow(
		x.ArchivedOn,
		x.InstrumentID,
		x.PreparationID,
		x.Notes,
		x.CreatedOn,
		x.UpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_buildGetRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleRequiredPreparationInstrumentID := uint64(123)

		expectedArgCount := 1
		expectedQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE id = $1"
		actualQuery, args := p.buildGetRequiredPreparationInstrumentQuery(exampleRequiredPreparationInstrumentID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleRequiredPreparationInstrumentID, args[0].(uint64))
	})
}

func TestPostgres_GetRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE id = $1"
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expected.ID).
			WillReturnRows(buildMockRowFromRequiredPreparationInstrument(expected))

		actual, err := p.GetRequiredPreparationInstrument(context.Background(), expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE id = $1"
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRequiredPreparationInstrument(context.Background(), expected.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRequiredPreparationInstrumentCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedArgCount := 0
		expectedQuery := "SELECT COUNT(id) FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"

		actualQuery, args := p.buildGetRequiredPreparationInstrumentCountQuery(models.DefaultQueryFilter())
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
	})
}

func TestPostgres_GetRequiredPreparationInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetRequiredPreparationInstrumentCount(context.Background(), models.DefaultQueryFilter())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRequiredPreparationInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM required_preparation_instruments WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllRequiredPreparationInstrumentsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRequiredPreparationInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM required_preparation_instruments WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRequiredPreparationInstrumentsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRequiredPreparationInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedArgCount := 0
		expectedQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"
		actualQuery, args := p.buildGetRequiredPreparationInstrumentsQuery(models.DefaultQueryFilter())

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
	})
}

func TestPostgres_GetRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedListQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM required_preparation_instruments WHERE archived_on IS NULL"
		expectedRequiredPreparationInstrument := &models.RequiredPreparationInstrument{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.RequiredPreparationInstrumentList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			RequiredPreparationInstruments: []models.RequiredPreparationInstrument{
				*expectedRequiredPreparationInstrument,
			},
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildMockRowFromRequiredPreparationInstrument(expectedRequiredPreparationInstrument))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetRequiredPreparationInstruments(context.Background(), models.DefaultQueryFilter())

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedListQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRequiredPreparationInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedListQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRequiredPreparationInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning required preparation instrument", func(t *testing.T) {
		expected := &models.RequiredPreparationInstrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildErroneousMockRowFromRequiredPreparationInstrument(expected))

		actual, err := p.GetRequiredPreparationInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expected := &models.RequiredPreparationInstrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, instrument_id, preparation_id, notes, created_on, updated_on, archived_on FROM required_preparation_instruments WHERE archived_on IS NULL LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM required_preparation_instruments WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildMockRowFromRequiredPreparationInstrument(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRequiredPreparationInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RequiredPreparationInstrument{
			ID: 321,
		}
		expectedArgCount := 3
		expectedQuery := "INSERT INTO required_preparation_instruments (instrument_id,preparation_id,notes) VALUES ($1,$2,$3) RETURNING id, created_on"
		actualQuery, args := p.buildCreateRequiredPreparationInstrumentQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.InstrumentID, args[0].(uint64))
		assert.Equal(t, expected.PreparationID, args[1].(uint64))
		assert.Equal(t, expected.Notes, args[2].(string))
	})
}

func TestPostgres_CreateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.RequiredPreparationInstrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RequiredPreparationInstrumentCreationInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO required_preparation_instruments (instrument_id,preparation_id,notes) VALUES ($1,$2,$3) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.PreparationID,
				expected.Notes,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRequiredPreparationInstrument(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expected := &models.RequiredPreparationInstrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RequiredPreparationInstrumentCreationInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		expectedQuery := "INSERT INTO required_preparation_instruments (instrument_id,preparation_id,notes) VALUES ($1,$2,$3) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.PreparationID,
				expected.Notes,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRequiredPreparationInstrument(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RequiredPreparationInstrument{
			ID: 321,
		}
		expectedArgCount := 4
		expectedQuery := "UPDATE required_preparation_instruments SET instrument_id = $1, preparation_id = $2, notes = $3, updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING updated_on"
		actualQuery, args := p.buildUpdateRequiredPreparationInstrumentQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.InstrumentID, args[0].(uint64))
		assert.Equal(t, expected.PreparationID, args[1].(uint64))
		assert.Equal(t, expected.Notes, args[2].(string))
		assert.Equal(t, expected.ID, args[3].(uint64))
	})
}

func TestPostgres_UpdateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.RequiredPreparationInstrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE required_preparation_instruments SET instrument_id = $1, preparation_id = $2, notes = $3, updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.PreparationID,
				expected.Notes,
				expected.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRequiredPreparationInstrument(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expected := &models.RequiredPreparationInstrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE required_preparation_instruments SET instrument_id = $1, preparation_id = $2, notes = $3, updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.PreparationID,
				expected.Notes,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRequiredPreparationInstrument(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RequiredPreparationInstrument{
			ID: 321,
		}
		expectedArgCount := 1
		expectedQuery := "UPDATE required_preparation_instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		actualQuery, args := p.buildArchiveRequiredPreparationInstrumentQuery(expected.ID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.ID, args[0].(uint64))
	})
}

func TestPostgres_ArchiveRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.RequiredPreparationInstrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE required_preparation_instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRequiredPreparationInstrument(context.Background(), expected.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		example := &models.RequiredPreparationInstrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE required_preparation_instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRequiredPreparationInstrument(context.Background(), example.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
