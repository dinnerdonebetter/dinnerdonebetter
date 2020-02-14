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

func buildMockRowFromPreparation(x *models.Preparation) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(preparationsTableColumns).AddRow(
		x.ID,
		x.Name,
		x.Variant,
		x.Description,
		x.AllergyWarning,
		x.Icon,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
	)

	return exampleRows
}

func buildErroneousMockRowFromPreparation(x *models.Preparation) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(preparationsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Variant,
		x.Description,
		x.AllergyWarning,
		x.Icon,
		x.CreatedOn,
		x.UpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_buildGetPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		examplePreparationID := uint64(123)

		expectedArgCount := 1
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE id = $1"
		actualQuery, args := p.buildGetPreparationQuery(examplePreparationID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, examplePreparationID, args[0].(uint64))
	})
}

func TestPostgres_GetPreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE id = $1"
		expected := &models.Preparation{
			ID: 123,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expected.ID).
			WillReturnRows(buildMockRowFromPreparation(expected))

		actual, err := p.GetPreparation(context.Background(), expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE id = $1"
		expected := &models.Preparation{
			ID: 123,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetPreparation(context.Background(), expected.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetPreparationCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedArgCount := 0
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL LIMIT 20"

		actualQuery, args := p.buildGetPreparationCountQuery(models.DefaultQueryFilter())
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
	})
}

func TestPostgres_GetPreparationCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetPreparationCount(context.Background(), models.DefaultQueryFilter())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllPreparationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllPreparationsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllPreparationsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedArgCount := 0
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE archived_on IS NULL LIMIT 20"
		actualQuery, args := p.buildGetPreparationsQuery(models.DefaultQueryFilter())

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
	})
}

func TestPostgres_GetPreparations(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE archived_on IS NULL LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL"
		expectedPreparation := &models.Preparation{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.PreparationList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			Preparations: []models.Preparation{
				*expectedPreparation,
			},
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildMockRowFromPreparation(expectedPreparation))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetPreparations(context.Background(), models.DefaultQueryFilter())

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetPreparations(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetPreparations(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning preparation", func(t *testing.T) {
		expected := &models.Preparation{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE archived_on IS NULL LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildErroneousMockRowFromPreparation(expected))

		actual, err := p.GetPreparations(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expected := &models.Preparation{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on FROM preparations WHERE archived_on IS NULL LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildMockRowFromPreparation(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetPreparations(context.Background(), models.DefaultQueryFilter())
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreatePreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Preparation{
			ID: 321,
		}
		expectedArgCount := 5
		expectedQuery := "INSERT INTO preparations (name,variant,description,allergy_warning,icon) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"
		actualQuery, args := p.buildCreatePreparationQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.AllergyWarning, args[3].(string))
		assert.Equal(t, expected.Icon, args[4].(string))
	})
}

func TestPostgres_CreatePreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.Preparation{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.PreparationCreationInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO preparations (name,variant,description,allergy_warning,icon) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
			).WillReturnRows(exampleRows)

		actual, err := p.CreatePreparation(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expected := &models.Preparation{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.PreparationCreationInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		expectedQuery := "INSERT INTO preparations (name,variant,description,allergy_warning,icon) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreatePreparation(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdatePreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Preparation{
			ID: 321,
		}
		expectedArgCount := 6
		expectedQuery := "UPDATE preparations SET name = $1, variant = $2, description = $3, allergy_warning = $4, icon = $5, updated_on = extract(epoch FROM NOW()) WHERE id = $6 RETURNING updated_on"
		actualQuery, args := p.buildUpdatePreparationQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.AllergyWarning, args[3].(string))
		assert.Equal(t, expected.Icon, args[4].(string))
		assert.Equal(t, expected.ID, args[5].(uint64))
	})
}

func TestPostgres_UpdatePreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.Preparation{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE preparations SET name = $1, variant = $2, description = $3, allergy_warning = $4, icon = $5, updated_on = extract(epoch FROM NOW()) WHERE id = $6 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
				expected.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdatePreparation(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expected := &models.Preparation{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE preparations SET name = $1, variant = $2, description = $3, allergy_warning = $4, icon = $5, updated_on = extract(epoch FROM NOW()) WHERE id = $6 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdatePreparation(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchivePreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Preparation{
			ID: 321,
		}
		expectedArgCount := 1
		expectedQuery := "UPDATE preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		actualQuery, args := p.buildArchivePreparationQuery(expected.ID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.ID, args[0].(uint64))
	})
}

func TestPostgres_ArchivePreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := &models.Preparation{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchivePreparation(context.Background(), expected.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		example := &models.Preparation{
			ID:        123,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchivePreparation(context.Background(), example.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
