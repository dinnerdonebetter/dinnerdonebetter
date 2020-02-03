package mariadb

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

func buildMockRowFromRecipeStepProduct(x *models.RecipeStepProduct) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepProductsTableColumns).AddRow(
		x.ID,
		x.Name,
		x.RecipeStepID,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepProduct(x *models.RecipeStepProduct) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepProductsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.RecipeStepID,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestMariaDB_buildGetRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleRecipeStepProductID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildGetRecipeStepProductQuery(exampleRecipeStepProductID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeStepProductID, args[1].(uint64))
	})
}

func TestMariaDB_GetRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepProduct{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeStepProduct(expected))

		actual, err := m.GetRecipeStepProduct(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepProduct{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepProduct(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepProductCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := m.buildGetRecipeStepProductCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepProductCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetRecipeStepProductCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllRecipeStepProductsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_products WHERE archived_on IS NULL"

		actualQuery := m.buildGetAllRecipeStepProductsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllRecipeStepProductsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_products WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllRecipeStepProductsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepProductsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := m.buildGetRecipeStepProductsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_products WHERE archived_on IS NULL"
		expectedRecipeStepProduct := &models.RecipeStepProduct{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.RecipeStepProductList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			RecipeStepProducts: []models.RecipeStepProduct{
				*expectedRecipeStepProduct,
			},
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepProduct(expectedRecipeStepProduct))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := m.GetRecipeStepProducts(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepProducts(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepProducts(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step product", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepProduct{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepProduct(expected))

		actual, err := m.GetRecipeStepProducts(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepProduct{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_products WHERE archived_on IS NULL"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepProduct(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepProducts(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_GetAllRecipeStepProductsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeStepProduct := &models.RecipeStepProduct{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepProduct(expectedRecipeStepProduct))

		expected := []models.RecipeStepProduct{*expectedRecipeStepProduct}
		actual, err := m.GetAllRecipeStepProductsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllRecipeStepProductsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllRecipeStepProductsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeStepProduct := &models.RecipeStepProduct{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_products WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepProduct(exampleRecipeStepProduct))

		actual, err := m.GetAllRecipeStepProductsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepProduct{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 3
		expectedQuery := "INSERT INTO recipe_step_products (name,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,UNIX_TIMESTAMP())"
		actualQuery, args := m.buildCreateRecipeStepProductQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.RecipeStepID, args[1].(uint64))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))
	})
}

func TestMariaDB_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepProduct{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepProductCreationInput{
			Name:         expected.Name,
			RecipeStepID: expected.RecipeStepID,
			BelongsTo:    expected.BelongsTo,
		}

		m, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO recipe_step_products (name,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,UNIX_TIMESTAMP())"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.Name,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM recipe_step_products WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := m.CreateRecipeStepProduct(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepProduct{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepProductCreationInput{
			Name:         expected.Name,
			RecipeStepID: expected.RecipeStepID,
			BelongsTo:    expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO recipe_step_products (name,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,UNIX_TIMESTAMP())"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateRecipeStepProduct(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepProduct{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 4
		expectedQuery := "UPDATE recipe_step_products SET name = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildUpdateRecipeStepProductQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.RecipeStepID, args[1].(uint64))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))
		assert.Equal(t, expected.ID, args[3].(uint64))
	})
}

func TestMariaDB_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepProduct{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE recipe_step_products SET name = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateRecipeStepProduct(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepProduct{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_products SET name = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateRecipeStepProduct(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepProduct{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_step_products SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := m.buildArchiveRecipeStepProductQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestMariaDB_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepProduct{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_products SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := m.ArchiveRecipeStepProduct(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.RecipeStepProduct{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_products SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveRecipeStepProduct(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
