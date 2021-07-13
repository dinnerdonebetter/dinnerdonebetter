package querier

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromRecipeStepIngredients(includeCounts bool, filteredCount uint64, recipeStepIngredients ...*types.RecipeStepIngredient) *sqlmock.Rows {
	columns := querybuilding.RecipeStepIngredientsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepIngredients {
		rowValues := []driver.Value{
			x.ID,
			x.ExternalID,
			x.IngredientID,
			x.Name,
			x.QuantityType,
			x.QuantityValue,
			x.QuantityNotes,
			x.ProductOfRecipeStep,
			x.IngredientNotes,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeStep,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeStepIngredients))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepIngredients(ctx, mockRows, false)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepIngredients(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildRecipeStepIngredientExistsQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, 0, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, 0, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildRecipeStepIngredientExistsQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildRecipeStepIngredientExistsQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredient))

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, 0, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, 0, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAllRecipeStepIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		exampleCount := uint64(123)

		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAllRecipeStepIngredientsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetAllRecipeStepIngredientsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAllRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		results := make(chan []*types.RecipeStepIngredient)
		doneChan := make(chan bool, 1)
		expectedCount := uint64(20)
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()
		exampleBatchSize := uint16(1000)

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAllRecipeStepIngredientsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetBatchOfRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		assert.NoError(t, c.GetAllRecipeStepIngredients(ctx, results, exampleBatchSize))

		stillQuerying := true
		for stillQuerying {
			select {
			case batch := <-results:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with nil results channel", func(t *testing.T) {
		t.Parallel()

		exampleBatchSize := uint16(1000)
		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.GetAllRecipeStepIngredients(ctx, nil, exampleBatchSize))
	})

	T.Run("with now rows returned", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.RecipeStepIngredient)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAllRecipeStepIngredientsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetBatchOfRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnError(sql.ErrNoRows)

		assert.NoError(t, c.GetAllRecipeStepIngredients(ctx, results, exampleBatchSize))

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.RecipeStepIngredient)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAllRecipeStepIngredientsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		err := c.GetAllRecipeStepIngredients(ctx, results, exampleBatchSize)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.RecipeStepIngredient)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAllRecipeStepIngredientsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetBatchOfRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnError(errors.New("blah"))

		assert.NoError(t, c.GetAllRecipeStepIngredients(ctx, results, exampleBatchSize))

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.RecipeStepIngredient)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAllRecipeStepIngredientsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetBatchOfRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		assert.NoError(t, c.GetAllRecipeStepIngredients(ctx, results, exampleBatchSize))

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(true, exampleRecipeStepIngredientList.FilteredCount, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredients(ctx, 0, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, 0, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()
		exampleRecipeStepIngredientList.Page = 0
		exampleRecipeStepIngredientList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(true, exampleRecipeStepIngredientList.FilteredCount, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsQuery",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetRecipeStepIngredientsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		var exampleIDs []uint64
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsWithIDsQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList.RecipeStepIngredients, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		var exampleIDs []uint64
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, 0, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("sets limit if not present", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()
		var exampleIDs []uint64
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsWithIDsQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList.RecipeStepIngredients, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()
		var exampleIDs []uint64
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsWithIDsQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()
		var exampleIDs []uint64
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetRecipeStepIngredientsWithIDsQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ExternalID = ""
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildCreateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeCreationQuery, fakeCreationArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleRecipeStepIngredient.CreatedOn
		}
		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ExternalID = ""

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(ctx, nil, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ExternalID = ""
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ExternalID = ""
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildCreateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleRecipeStepIngredient.CreatedOn
		}

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error creating audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ExternalID = ""
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildCreateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeCreationQuery, fakeCreationArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ExternalID = ""
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildCreateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeCreationQuery, fakeCreationArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleRecipeStepIngredient.CreatedOn
		}

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildUpdateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectCommit()

		assert.NoError(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, nil, exampleUser.ID, nil))
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient, 0, nil))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient, exampleUser.ID, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildUpdateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing audit log entry to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildUpdateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildUpdateRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildArchiveRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.NoError(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, 0, exampleRecipeStepIngredient.ID, exampleUserID))
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, 0, exampleUserID))
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID, 0))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID, exampleUserID))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildArchiveRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildArchiveRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildArchiveRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleRecipeStepIngredient.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAuditLogEntriesForRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleAuditLogEntriesList := fakes.BuildFakeAuditLogEntryList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromAuditLogEntries(false, exampleAuditLogEntriesList.Entries...))

		actual, err := c.GetAuditLogEntriesForRecipeStepIngredient(ctx, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntriesList.Entries, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetAuditLogEntriesForRecipeStepIngredient(ctx, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetAuditLogEntriesForRecipeStepIngredient(ctx, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.RecipeStepIngredientSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForRecipeStepIngredientQuery",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetAuditLogEntriesForRecipeStepIngredient(ctx, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}
