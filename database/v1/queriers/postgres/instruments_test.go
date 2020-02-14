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

func buildMockRowFromInstrument(x *models.Instrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(instrumentsTableColumns).AddRow(
		x.ID,
		x.Name,
		x.Variant,
		x.Description,
		x.Icon,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
	)

	return exampleRows
}

func buildErroneousMockRowFromInstrument(x *models.Instrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(instrumentsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Variant,
		x.Description,
		x.Icon,
		x.CreatedOn,
		x.UpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_buildGetInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleInstrumentID := uint64(123)

		expectedArgCount := 1
		expectedQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE id = $1"
		actualQuery, args := p.buildGetInstrumentQuery(exampleInstrumentID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleInstrumentID, args[0].(uint64))
	})
}

func TestPostgres_GetInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE id = $1"
		expected := &models.Instrument{
			ID: 123,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expected.ID).
			WillReturnRows(buildMockRowFromInstrument(expected))

		actual, err := p.GetInstrument(context.Background(), expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE id = $1"
		expected := &models.Instrument{
			ID: 123,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetInstrument(context.Background(), expected.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetInstrumentCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedArgCount := 0
		expectedQuery := "SELECT COUNT(id) FROM instruments WHERE archived_on IS NULL LIMIT 20"

		actualQuery, args := p.buildGetInstrumentCountQuery(models.DefaultQueryFilter())
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
	})
}

func TestPostgres_GetInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM instruments WHERE archived_on IS NULL LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetInstrumentCount(context.Background(), models.DefaultQueryFilter())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM instruments WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllInstrumentsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM instruments WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllInstrumentsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedArgCount := 0
		expectedQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE archived_on IS NULL LIMIT 20"
		actualQuery, args := p.buildGetInstrumentsQuery(models.DefaultQueryFilter())

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
	})
}

func TestPostgres_GetInstruments(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedListQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE archived_on IS NULL LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM instruments WHERE archived_on IS NULL"
		expectedInstrument := &models.Instrument{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.InstrumentList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			Instruments: []models.Instrument{
				*expectedInstrument,
			},
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildMockRowFromInstrument(expectedInstrument))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetInstruments(context.Background(), models.DefaultQueryFilter())

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedListQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedListQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning instrument", func(t *testing.T) {
		expected := &models.Instrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildErroneousMockRowFromInstrument(expected))

		actual, err := p.GetInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expected := &models.Instrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, icon, created_on, updated_on, archived_on FROM instruments WHERE archived_on IS NULL LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM instruments WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildMockRowFromInstrument(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetInstruments(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Instrument{
			ID: 321,
		}
		expectedArgCount := 4
		expectedQuery := "INSERT INTO instruments (name,variant,description,icon) VALUES ($1,$2,$3,$4) RETURNING id, created_on"
		actualQuery, args := p.buildCreateInstrumentQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.Icon, args[3].(string))
	})
}

func TestPostgres_CreateInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.Instrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.InstrumentCreationInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO instruments (name,variant,description,icon) VALUES ($1,$2,$3,$4) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Icon,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateInstrument(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expected := &models.Instrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.InstrumentCreationInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		expectedQuery := "INSERT INTO instruments (name,variant,description,icon) VALUES ($1,$2,$3,$4) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Icon,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateInstrument(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Instrument{
			ID: 321,
		}
		expectedArgCount := 5
		expectedQuery := "UPDATE instruments SET name = $1, variant = $2, description = $3, icon = $4, updated_on = extract(epoch FROM NOW()) WHERE id = $5 RETURNING updated_on"
		actualQuery, args := p.buildUpdateInstrumentQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.Icon, args[3].(string))
		assert.Equal(t, expected.ID, args[4].(uint64))
	})
}

func TestPostgres_UpdateInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.Instrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE instruments SET name = $1, variant = $2, description = $3, icon = $4, updated_on = extract(epoch FROM NOW()) WHERE id = $5 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Icon,
				expected.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateInstrument(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expected := &models.Instrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE instruments SET name = $1, variant = $2, description = $3, icon = $4, updated_on = extract(epoch FROM NOW()) WHERE id = $5 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Icon,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateInstrument(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Instrument{
			ID: 321,
		}
		expectedArgCount := 1
		expectedQuery := "UPDATE instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		actualQuery, args := p.buildArchiveInstrumentQuery(expected.ID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.ID, args[0].(uint64))
	})
}

func TestPostgres_ArchiveInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.Instrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveInstrument(context.Background(), expected.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		example := &models.Instrument{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE instruments SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveInstrument(context.Background(), example.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
