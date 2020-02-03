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

func buildMockRowFromRecipeStepIngredient(x *models.RecipeStepIngredient) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepIngredientsTableColumns).AddRow(
		x.ID,
		x.IngredientID,
		x.QuantityType,
		x.QuantityValue,
		x.QuantityNotes,
		x.ProductOfRecipe,
		x.IngredientNotes,
		x.RecipeStepID,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepIngredient(x *models.RecipeStepIngredient) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepIngredientsTableColumns).AddRow(
		x.ArchivedOn,
		x.IngredientID,
		x.QuantityType,
		x.QuantityValue,
		x.QuantityNotes,
		x.ProductOfRecipe,
		x.IngredientNotes,
		x.RecipeStepID,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestMariaDB_buildGetRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleRecipeStepIngredientID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildGetRecipeStepIngredientQuery(exampleRecipeStepIngredientID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeStepIngredientID, args[1].(uint64))
	})
}

func TestMariaDB_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expected))

		actual, err := m.GetRecipeStepIngredient(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepIngredient(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepIngredientCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := m.buildGetRecipeStepIngredientCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepIngredientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetRecipeStepIngredientCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllRecipeStepIngredientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL"

		actualQuery := m.buildGetAllRecipeStepIngredientsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllRecipeStepIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllRecipeStepIngredientsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := m.buildGetRecipeStepIngredientsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL"
		expectedRecipeStepIngredient := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.RecipeStepIngredientList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			RecipeStepIngredients: []models.RecipeStepIngredient{
				*expectedRecipeStepIngredient,
			},
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expectedRecipeStepIngredient))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := m.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step ingredient", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepIngredient(expected))

		actual, err := m.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_GetAllRecipeStepIngredientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeStepIngredient := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expectedRecipeStepIngredient))

		expected := []models.RecipeStepIngredient{*expectedRecipeStepIngredient}
		actual, err := m.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeStepIngredient := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepIngredient(exampleRecipeStepIngredient))

		actual, err := m.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepIngredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 8
		expectedQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())"
		actualQuery, args := m.buildCreateRecipeStepIngredientQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.IngredientID, args[0].(*uint64))
		assert.Equal(t, expected.QuantityType, args[1].(string))
		assert.Equal(t, expected.QuantityValue, args[2].(float32))
		assert.Equal(t, expected.QuantityNotes, args[3].(string))
		assert.Equal(t, expected.ProductOfRecipe, args[4].(bool))
		assert.Equal(t, expected.IngredientNotes, args[5].(string))
		assert.Equal(t, expected.RecipeStepID, args[6].(uint64))
		assert.Equal(t, expected.BelongsTo, args[7].(uint64))
	})
}

func TestMariaDB_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepIngredientCreationInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
			BelongsTo:       expected.BelongsTo,
		}

		m, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.IngredientID,
				expected.QuantityType,
				expected.QuantityValue,
				expected.QuantityNotes,
				expected.ProductOfRecipe,
				expected.IngredientNotes,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM recipe_step_ingredients WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := m.CreateRecipeStepIngredient(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepIngredientCreationInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
			BelongsTo:       expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.IngredientID,
				expected.QuantityType,
				expected.QuantityValue,
				expected.QuantityNotes,
				expected.ProductOfRecipe,
				expected.IngredientNotes,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateRecipeStepIngredient(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepIngredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 9
		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = ?, quantity_type = ?, quantity_value = ?, quantity_notes = ?, product_of_recipe = ?, ingredient_notes = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildUpdateRecipeStepIngredientQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.IngredientID, args[0].(*uint64))
		assert.Equal(t, expected.QuantityType, args[1].(string))
		assert.Equal(t, expected.QuantityValue, args[2].(float32))
		assert.Equal(t, expected.QuantityNotes, args[3].(string))
		assert.Equal(t, expected.ProductOfRecipe, args[4].(bool))
		assert.Equal(t, expected.IngredientNotes, args[5].(string))
		assert.Equal(t, expected.RecipeStepID, args[6].(uint64))
		assert.Equal(t, expected.BelongsTo, args[7].(uint64))
		assert.Equal(t, expected.ID, args[8].(uint64))
	})
}

func TestMariaDB_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = ?, quantity_type = ?, quantity_value = ?, quantity_notes = ?, product_of_recipe = ?, ingredient_notes = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.IngredientID,
				expected.QuantityType,
				expected.QuantityValue,
				expected.QuantityNotes,
				expected.ProductOfRecipe,
				expected.IngredientNotes,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateRecipeStepIngredient(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = ?, quantity_type = ?, quantity_value = ?, quantity_notes = ?, product_of_recipe = ?, ingredient_notes = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.IngredientID,
				expected.QuantityType,
				expected.QuantityValue,
				expected.QuantityNotes,
				expected.ProductOfRecipe,
				expected.IngredientNotes,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateRecipeStepIngredient(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepIngredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_step_ingredients SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := m.buildArchiveRecipeStepIngredientQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestMariaDB_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_ingredients SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := m.ArchiveRecipeStepIngredient(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_ingredients SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveRecipeStepIngredient(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
