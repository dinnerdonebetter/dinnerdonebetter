package sqlite

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
		x.BelongsTo,
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
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestSqlite_buildGetPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		examplePreparationID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildGetPreparationQuery(examplePreparationID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, examplePreparationID, args[1].(uint64))
	})
}

func TestSqlite_GetPreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE belongs_to = ? AND id = ?"
		expected := &models.Preparation{
			ID: 123,
		}
		expectedUserID := uint64(321)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromPreparation(expected))

		actual, err := s.GetPreparation(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE belongs_to = ? AND id = ?"
		expected := &models.Preparation{
			ID: 123,
		}
		expectedUserID := uint64(321)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetPreparation(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetPreparationCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := s.buildGetPreparationCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetPreparationCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetPreparationCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllPreparationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL"

		actualQuery := s.buildGetAllPreparationsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllPreparationsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := s.buildGetPreparationsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetPreparations(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
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

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromPreparation(expectedPreparation))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := s.GetPreparations(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetPreparations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetPreparations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning preparation", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.Preparation{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromPreparation(expected))

		actual, err := s.GetPreparations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.Preparation{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM preparations WHERE archived_on IS NULL"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromPreparation(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetPreparations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_GetAllPreparationsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedPreparation := &models.Preparation{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromPreparation(expectedPreparation))

		expected := []models.Preparation{*expectedPreparation}
		actual, err := s.GetAllPreparationsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetAllPreparationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetAllPreparationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		examplePreparation := &models.Preparation{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, allergy_warning, icon, created_on, updated_on, archived_on, belongs_to FROM preparations WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromPreparation(examplePreparation))

		actual, err := s.GetAllPreparationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildCreatePreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Preparation{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 6
		expectedQuery := "INSERT INTO preparations (name,variant,description,allergy_warning,icon,belongs_to) VALUES (?,?,?,?,?,?)"
		actualQuery, args := s.buildCreatePreparationQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.AllergyWarning, args[3].(string))
		assert.Equal(t, expected.Icon, args[4].(string))
		assert.Equal(t, expected.BelongsTo, args[5].(uint64))
	})
}

func TestSqlite_CreatePreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Preparation{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.PreparationCreationInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
			BelongsTo:      expected.BelongsTo,
		}

		s, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO preparations (name,variant,description,allergy_warning,icon,belongs_to) VALUES (?,?,?,?,?,?)"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM preparations WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := s.CreatePreparation(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Preparation{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.PreparationCreationInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
			BelongsTo:      expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO preparations (name,variant,description,allergy_warning,icon,belongs_to) VALUES (?,?,?,?,?,?)"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := s.CreatePreparation(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildUpdatePreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Preparation{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 7
		expectedQuery := "UPDATE preparations SET name = ?, variant = ?, description = ?, allergy_warning = ?, icon = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildUpdatePreparationQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.AllergyWarning, args[3].(string))
		assert.Equal(t, expected.Icon, args[4].(string))
		assert.Equal(t, expected.BelongsTo, args[5].(uint64))
		assert.Equal(t, expected.ID, args[6].(uint64))
	})
}

func TestSqlite_UpdatePreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Preparation{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE preparations SET name = ?, variant = ?, description = ?, allergy_warning = ?, icon = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := s.UpdatePreparation(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Preparation{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE preparations SET name = ?, variant = ?, description = ?, allergy_warning = ?, icon = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.AllergyWarning,
				expected.Icon,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := s.UpdatePreparation(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildArchivePreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Preparation{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE preparations SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := s.buildArchivePreparationQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestSqlite_ArchivePreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Preparation{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE preparations SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.ArchivePreparation(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.Preparation{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE preparations SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := s.ArchivePreparation(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
