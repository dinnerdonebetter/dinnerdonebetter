package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromRecipeStepProducts(recipeStepProducts ...*models.RecipeStepProduct) *sqlmock.Rows {
	columns := recipeStepProductsTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepProducts {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.RecipeStepID,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeStep,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepProduct(x *models.RecipeStepProduct) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepProductsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.RecipeStepID,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToRecipeStep,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanRecipeStepProducts(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanRecipeStepProducts(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeStepProductExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepProduct.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeStepProductExistsQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeStepProductExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepProduct.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeStepProductExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepProduct.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeStepProductExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepProduct.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepProductQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepProduct(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepProduct.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeStepProducts(exampleRecipeStepProduct))

		actual, err := p.GetRecipeStepProduct(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProduct, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepProduct.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepProduct(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepProductsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_step_products.id) FROM recipe_step_products WHERE recipe_step_products.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeStepProductsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepProductsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_step_products.id) FROM recipe_step_products WHERE recipe_step_products.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepProductsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfRecipeStepProductsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products WHERE recipe_step_products.id > $1 AND recipe_step_products.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRecipeStepProductsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRecipeStepProducts(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(recipe_step_products.id) FROM recipe_step_products WHERE recipe_step_products.archived_on IS NULL"
	expectedGetQuery := "SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products WHERE recipe_step_products.id > $1 AND recipe_step_products.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepProducts(
					&exampleRecipeStepProductList.RecipeStepProducts[0],
					&exampleRecipeStepProductList.RecipeStepProducts[1],
					&exampleRecipeStepProductList.RecipeStepProducts[2],
				),
			)

		out := make(chan []models.RecipeStepProduct)
		doneChan := make(chan bool, 1)

		err := p.GetAllRecipeStepProducts(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.RecipeStepProduct)

		err := p.GetAllRecipeStepProducts(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []models.RecipeStepProduct)

		err := p.GetAllRecipeStepProducts(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.RecipeStepProduct)

		err := p.GetAllRecipeStepProducts(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepProduct(exampleRecipeStepProduct))

		out := make(chan []models.RecipeStepProduct)

		err := p.GetAllRecipeStepProducts(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepProductsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 AND recipe_step_products.created_on > $5 AND recipe_step_products.created_on < $6 AND recipe_step_products.last_updated_on > $7 AND recipe_step_products.last_updated_on < $8 ORDER BY recipe_step_products.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepProductsQuery(exampleRecipe.ID, exampleRecipeStep.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepProducts(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 ORDER BY recipe_step_products.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepProducts(
					&exampleRecipeStepProductList.RecipeStepProducts[0],
					&exampleRecipeStepProductList.RecipeStepProducts[1],
					&exampleRecipeStepProductList.RecipeStepProducts[2],
				),
			)

		actual, err := p.GetRecipeStepProducts(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepProducts(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepProducts(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step product", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepProduct(exampleRecipeStepProduct))

		actual, err := p.GetRecipeStepProducts(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepProductsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM (SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_products WHERE recipe_step_products.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepProductsWithIDsQuery(exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepProductsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList()
		var exampleIDs []uint64
		for _, recipeStepProduct := range exampleRecipeStepProductList.RecipeStepProducts {
			exampleIDs = append(exampleIDs, recipeStepProduct.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM (SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_products WHERE recipe_step_products.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepProducts(
					&exampleRecipeStepProductList.RecipeStepProducts[0],
					&exampleRecipeStepProductList.RecipeStepProducts[1],
					&exampleRecipeStepProductList.RecipeStepProducts[2],
				),
			)

		actual, err := p.GetRecipeStepProductsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList.RecipeStepProducts, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM (SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_products WHERE recipe_step_products.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepProductsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM (SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_products WHERE recipe_step_products.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepProductsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step product", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM (SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_products WHERE recipe_step_products.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepProduct(exampleRecipeStepProduct))

		actual, err := p.GetRecipeStepProductsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "INSERT INTO recipe_step_products (name,recipe_step_id,belongs_to_recipe_step) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeStepProduct.Name,
			exampleRecipeStepProduct.RecipeStepID,
			exampleRecipeStepProduct.BelongsToRecipeStep,
		}
		actualQuery, actualArgs := p.buildCreateRecipeStepProductQuery(exampleRecipeStepProduct)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_step_products (name,recipe_step_id,belongs_to_recipe_step) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeStepProduct.ID, exampleRecipeStepProduct.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepProduct.Name,
				exampleRecipeStepProduct.RecipeStepID,
				exampleRecipeStepProduct.BelongsToRecipeStep,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStepProduct(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProduct, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepProduct.Name,
				exampleRecipeStepProduct.RecipeStepID,
				exampleRecipeStepProduct.BelongsToRecipeStep,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeStepProduct(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_products SET name = $1, recipe_step_id = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $3 AND id = $4 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRecipeStepProduct.Name,
			exampleRecipeStepProduct.RecipeStepID,
			exampleRecipeStepProduct.BelongsToRecipeStep,
			exampleRecipeStepProduct.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeStepProductQuery(exampleRecipeStepProduct)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_products SET name = $1, recipe_step_id = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $3 AND id = $4 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepProduct.Name,
				exampleRecipeStepProduct.RecipeStepID,
				exampleRecipeStepProduct.BelongsToRecipeStep,
				exampleRecipeStepProduct.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStepProduct(ctx, exampleRecipeStepProduct)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepProduct.Name,
				exampleRecipeStepProduct.RecipeStepID,
				exampleRecipeStepProduct.BelongsToRecipeStep,
				exampleRecipeStepProduct.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeStepProduct(ctx, exampleRecipeStepProduct)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_products SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepProduct.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeStepProductQuery(exampleRecipeStep.ID, exampleRecipeStepProduct.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_products SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepProduct.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStepProduct(ctx, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepProduct.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeStepProduct(ctx, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepProduct.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStepProduct(ctx, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
