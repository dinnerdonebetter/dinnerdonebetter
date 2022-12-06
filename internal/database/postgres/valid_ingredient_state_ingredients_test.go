package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func buildMockRowsFromValidIngredientStateIngredients(includeCounts bool, filteredCount uint64, validIngredientStateIngredients ...*types.ValidIngredientStateIngredient) *sqlmock.Rows {
	columns := validIngredientStateIngredientsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredientStateIngredients {
		rowValues := []driver.Value{
			&x.ID,
			&x.Notes,
			&x.IngredientState.ID,
			&x.IngredientState.Name,
			&x.IngredientState.Description,
			&x.IngredientState.IconPath,
			&x.IngredientState.Slug,
			&x.IngredientState.PastTense,
			&x.IngredientState.AttributeType,
			&x.IngredientState.CreatedAt,
			&x.IngredientState.LastUpdatedAt,
			&x.IngredientState.ArchivedAt,
			&x.Ingredient.ID,
			&x.Ingredient.Name,
			&x.Ingredient.Description,
			&x.Ingredient.Warning,
			&x.Ingredient.ContainsEgg,
			&x.Ingredient.ContainsDairy,
			&x.Ingredient.ContainsPeanut,
			&x.Ingredient.ContainsTreeNut,
			&x.Ingredient.ContainsSoy,
			&x.Ingredient.ContainsWheat,
			&x.Ingredient.ContainsShellfish,
			&x.Ingredient.ContainsSesame,
			&x.Ingredient.ContainsFish,
			&x.Ingredient.ContainsGluten,
			&x.Ingredient.AnimalFlesh,
			&x.Ingredient.IsMeasuredVolumetrically,
			&x.Ingredient.IsLiquid,
			&x.Ingredient.IconPath,
			&x.Ingredient.AnimalDerived,
			&x.Ingredient.PluralName,
			&x.Ingredient.RestrictToPreparations,
			&x.Ingredient.MinimumIdealStorageTemperatureInCelsius,
			&x.Ingredient.MaximumIdealStorageTemperatureInCelsius,
			&x.Ingredient.StorageInstructions,
			&x.Ingredient.Slug,
			&x.Ingredient.ContainsAlcohol,
			&x.Ingredient.ShoppingSuggestions,
			&x.Ingredient.CreatedAt,
			&x.Ingredient.LastUpdatedAt,
			&x.Ingredient.ArchivedAt,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validIngredientStateIngredients))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidIngredientStateIngredients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidIngredientStateIngredients(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidIngredientStateIngredients(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidIngredientStateIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidIngredientStateIngredientExists(ctx, exampleValidIngredientStateIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientStateIngredientExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidIngredientStateIngredientExists(ctx, exampleValidIngredientStateIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidIngredientStateIngredientExists(ctx, exampleValidIngredientStateIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientStateIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientStateIngredients(false, 0, exampleValidIngredientStateIngredient))

		actual, err := c.GetValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStateIngredient, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientStateIngredient(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientStateIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientStateIngredients(T *testing.T) {
	T.Parallel()

	joins := []string{
		validIngredientsOnValidIngredientStateIngredientsJoinClause,
		validPreparationsOnValidIngredientStateIngredientsJoinClause,
	}

	groupBys := []string{
		"valid_ingredients.id",
		"valid_preparations.id",
		"valid_ingredient_state_ingredients.id",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientStateIngredientList := fakes.BuildFakeValidIngredientStateIngredientList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_state_ingredients", joins, groupBys, nil, householdOwnershipColumn, validIngredientStateIngredientsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientStateIngredients(true, exampleValidIngredientStateIngredientList.FilteredCount, exampleValidIngredientStateIngredientList.Data...))

		actual, err := c.GetValidIngredientStateIngredients(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStateIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidIngredientStateIngredientList := fakes.BuildFakeValidIngredientStateIngredientList()
		exampleValidIngredientStateIngredientList.Page = 0
		exampleValidIngredientStateIngredientList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_state_ingredients", joins, groupBys, nil, householdOwnershipColumn, validIngredientStateIngredientsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientStateIngredients(true, exampleValidIngredientStateIngredientList.FilteredCount, exampleValidIngredientStateIngredientList.Data...))

		actual, err := c.GetValidIngredientStateIngredients(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStateIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_state_ingredients", joins, groupBys, nil, householdOwnershipColumn, validIngredientStateIngredientsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientStateIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_state_ingredients", joins, groupBys, nil, householdOwnershipColumn, validIngredientStateIngredientsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredientStateIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
		exampleValidIngredientStateIngredient.ID = "1"
		exampleValidIngredientStateIngredient.IngredientState = types.ValidIngredientState{ID: exampleValidIngredientStateIngredient.IngredientState.ID}
		exampleValidIngredientStateIngredient.Ingredient = types.ValidIngredient{ID: exampleValidIngredientStateIngredient.Ingredient.ID}

		exampleInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientDatabaseCreationInput(exampleValidIngredientStateIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidIngredientStateID,
			exampleInput.ValidIngredientID,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientStateIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidIngredientStateIngredient.CreatedAt
		}

		actual, err := c.CreateValidIngredientStateIngredient(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStateIngredient, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientStateIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
		exampleInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientDatabaseCreationInput(exampleValidIngredientStateIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidIngredientStateID,
			exampleInput.ValidIngredientID,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientStateIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidIngredientStateIngredient.CreatedAt
		}

		actual, err := c.CreateValidIngredientStateIngredient(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientStateIngredient.Notes,
			exampleValidIngredientStateIngredient.IngredientState.ID,
			exampleValidIngredientStateIngredient.Ingredient.ID,
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientStateIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredient))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientStateIngredient(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientStateIngredient.Notes,
			exampleValidIngredientStateIngredient.IngredientState.ID,
			exampleValidIngredientStateIngredient.Ingredient.ID,
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientStateIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredient))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientStateIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientStateIngredient(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientStateIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientStateIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
