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

func TestPostgres_buildGetRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleRecipeStepIngredientID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE belongs_to = $1 AND id = $2"
		actualQuery, args := p.buildGetRecipeStepIngredientQuery(exampleRecipeStepIngredientID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeStepIngredientID, args[1].(uint64))
	})
}

func TestPostgres_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE belongs_to = $1 AND id = $2"
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expected))

		actual, err := p.GetRecipeStepIngredient(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE belongs_to = $1 AND id = $2"
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepIngredient(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepIngredientCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		actualQuery, args := p.buildGetRecipeStepIngredientCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetRecipeStepIngredientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetRecipeStepIngredientCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepIngredientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllRecipeStepIngredientsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepIngredientsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		actualQuery, args := p.buildGetRecipeStepIngredientsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
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

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expectedRecipeStepIngredient))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step ingredient", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepIngredient(expected))

		actual, err := p.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_ingredients WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_GetAllRecipeStepIngredientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeStepIngredient := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepIngredient(expectedRecipeStepIngredient))

		expected := []models.RecipeStepIngredient{*expectedRecipeStepIngredient}
		actual, err := p.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeStepIngredient := &models.RecipeStepIngredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, ingredient_id, quantity_type, quantity_value, quantity_notes, product_of_recipe, ingredient_notes, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepIngredient(exampleRecipeStepIngredient))

		actual, err := p.GetAllRecipeStepIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeStepIngredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 8
		expectedQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,recipe_step_id,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_on"
		actualQuery, args := p.buildCreateRecipeStepIngredientQuery(expected)

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

func TestPostgres_CreateRecipeStepIngredient(T *testing.T) {
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
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,recipe_step_id,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.IngredientID,
				expected.QuantityType,
				expected.QuantityValue,
				expected.QuantityNotes,
				expected.ProductOfRecipe,
				expected.IngredientNotes,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStepIngredient(context.Background(), expectedInput)
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
		expectedQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,recipe_step_id,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		actual, err := p.CreateRecipeStepIngredient(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeStepIngredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 9
		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = $1, quantity_type = $2, quantity_value = $3, quantity_notes = $4, product_of_recipe = $5, ingredient_notes = $6, recipe_step_id = $7, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $8 AND id = $9 RETURNING updated_on"
		actualQuery, args := p.buildUpdateRecipeStepIngredientQuery(expected)

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

func TestPostgres_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = $1, quantity_type = $2, quantity_value = $3, quantity_notes = $4, product_of_recipe = $5, ingredient_notes = $6, recipe_step_id = $7, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $8 AND id = $9 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStepIngredient(context.Background(), expected)
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
		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = $1, quantity_type = $2, quantity_value = $3, quantity_notes = $4, product_of_recipe = $5, ingredient_notes = $6, recipe_step_id = $7, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $8 AND id = $9 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		err := p.UpdateRecipeStepIngredient(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeStepIngredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_step_ingredients SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"
		actualQuery, args := p.buildArchiveRecipeStepIngredientQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestPostgres_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepIngredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_ingredients SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStepIngredient(context.Background(), expected.ID, expectedUserID)
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
		expectedQuery := "UPDATE recipe_step_ingredients SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStepIngredient(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
