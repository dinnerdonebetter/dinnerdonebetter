package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromValidIngredients(includeCounts bool, filteredCount uint64, validIngredients ...*types.ValidIngredient) *sqlmock.Rows {
	columns := validIngredientsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredients {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Description,
			x.Warning,
			x.ContainsEgg,
			x.ContainsDairy,
			x.ContainsPeanut,
			x.ContainsTreeNut,
			x.ContainsSoy,
			x.ContainsWheat,
			x.ContainsShellfish,
			x.ContainsSesame,
			x.ContainsFish,
			x.ContainsGluten,
			x.AnimalFlesh,
			x.IsMeasuredVolumetrically,
			x.IsLiquid,
			x.IconPath,
			x.AnimalDerived,
			x.PluralName,
			x.RestrictToPreparations,
			x.MinimumIdealStorageTemperatureInCelsius,
			x.MaximumIdealStorageTemperatureInCelsius,
			x.StorageInstructions,
			x.Slug,
			x.ContainsAlcohol,
			x.ShoppingSuggestions,
			&x.IsStarch,
			&x.IsProtein,
			&x.IsGrain,
			&x.IsFruit,
			&x.IsSalt,
			&x.IsFat,
			&x.IsAcid,
			&x.IsHeat,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validIngredients))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidIngredients(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidIngredients(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidIngredientExists(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidIngredientExists(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidIngredientExists(ctx, exampleValidIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getValidIngredientArgs := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(getValidIngredientArgs)...).
			WillReturnRows(buildMockRowsFromValidIngredients(false, 0, exampleValidIngredient))

		actual, err := c.GetValidIngredient(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredient(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredient(ctx, exampleValidIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRandomValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredients(false, 0, exampleValidIngredient))

		actual, err := c.GetRandomValidIngredient(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRandomValidIngredient(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidIngredients(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredients := fakes.BuildFakeValidIngredientList()

		ctx := context.Background()
		c, db := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredients(false, 0, exampleValidIngredients.Data...))

		actual, err := c.SearchForValidIngredients(ctx, exampleQuery, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredients.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		actual, err := c.SearchForValidIngredients(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidIngredients(ctx, exampleQuery, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidIngredients(ctx, exampleQuery, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidIngredientsForPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		exampleValidIngredients := fakes.BuildFakeValidIngredientList()
		exampleValidIngredients.FilteredCount = 0
		exampleValidIngredients.TotalCount = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		searchByPreparationAndIngredientNameArgs := []any{
			exampleValidPreparation.ID,
			wrapQueryForILIKE(exampleValidIngredient.Name),
		}

		db.ExpectQuery(formatQueryForSQLMock(searchForIngredientsByPreparationAndIngredientNameQuery)).
			WithArgs(interfaceToDriverValue(searchByPreparationAndIngredientNameArgs)...).
			WillReturnRows(buildMockRowsFromValidIngredients(false, 0, exampleValidIngredients.Data...))

		actual, err := c.SearchForValidIngredientsForPreparation(ctx, exampleValidPreparation.ID, exampleValidIngredient.Name, nil)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredients, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidIngredientsForPreparation(ctx, exampleValidPreparation.ID, exampleValidIngredient.Name, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		searchByPreparationAndIngredientNameArgs := []any{
			exampleValidPreparation.ID,
			wrapQueryForILIKE(exampleValidIngredient.Name),
		}

		db.ExpectQuery(formatQueryForSQLMock(searchForIngredientsByPreparationAndIngredientNameQuery)).
			WithArgs(interfaceToDriverValue(searchByPreparationAndIngredientNameArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidIngredientsForPreparation(ctx, exampleValidPreparation.ID, exampleValidIngredient.Name, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredients", nil, nil, nil, householdOwnershipColumn, validIngredientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredients(true, exampleValidIngredientList.FilteredCount, exampleValidIngredientList.Data...))

		actual, err := c.GetValidIngredients(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()
		exampleValidIngredientList.Page = 0
		exampleValidIngredientList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredients", nil, nil, nil, householdOwnershipColumn, validIngredientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredients(true, exampleValidIngredientList.FilteredCount, exampleValidIngredientList.Data...))

		actual, err := c.GetValidIngredients(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredients", nil, nil, nil, householdOwnershipColumn, validIngredientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredients", nil, nil, nil, householdOwnershipColumn, validIngredientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		exampleValidIngredient.ID = "1"
		exampleInput := converters.ConvertValidIngredientToValidIngredientDatabaseCreationInput(exampleValidIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Warning,
			exampleInput.ContainsEgg,
			exampleInput.ContainsDairy,
			exampleInput.ContainsPeanut,
			exampleInput.ContainsTreeNut,
			exampleInput.ContainsSoy,
			exampleInput.ContainsWheat,
			exampleInput.ContainsShellfish,
			exampleInput.ContainsSesame,
			exampleInput.ContainsFish,
			exampleInput.ContainsGluten,
			exampleInput.AnimalFlesh,
			exampleInput.IsMeasuredVolumetrically,
			exampleInput.IsLiquid,
			exampleInput.IconPath,
			exampleInput.AnimalDerived,
			exampleInput.PluralName,
			exampleInput.RestrictToPreparations,
			exampleInput.MinimumIdealStorageTemperatureInCelsius,
			exampleInput.MaximumIdealStorageTemperatureInCelsius,
			exampleInput.StorageInstructions,
			exampleInput.Slug,
			exampleInput.ContainsAlcohol,
			exampleInput.ShoppingSuggestions,
			exampleInput.IsStarch,
			exampleInput.IsProtein,
			exampleInput.IsGrain,
			exampleInput.IsFruit,
			exampleInput.IsSalt,
			exampleInput.IsFat,
			exampleInput.IsAcid,
			exampleInput.IsHeat,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidIngredient.CreatedAt
		}

		actual, err := c.CreateValidIngredient(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		exampleInput := converters.ConvertValidIngredientToValidIngredientDatabaseCreationInput(exampleValidIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Warning,
			exampleInput.ContainsEgg,
			exampleInput.ContainsDairy,
			exampleInput.ContainsPeanut,
			exampleInput.ContainsTreeNut,
			exampleInput.ContainsSoy,
			exampleInput.ContainsWheat,
			exampleInput.ContainsShellfish,
			exampleInput.ContainsSesame,
			exampleInput.ContainsFish,
			exampleInput.ContainsGluten,
			exampleInput.AnimalFlesh,
			exampleInput.IsMeasuredVolumetrically,
			exampleInput.IsLiquid,
			exampleInput.IconPath,
			exampleInput.AnimalDerived,
			exampleInput.PluralName,
			exampleInput.RestrictToPreparations,
			exampleInput.MinimumIdealStorageTemperatureInCelsius,
			exampleInput.MaximumIdealStorageTemperatureInCelsius,
			exampleInput.StorageInstructions,
			exampleInput.Slug,
			exampleInput.ContainsAlcohol,
			exampleInput.ShoppingSuggestions,
			exampleInput.IsStarch,
			exampleInput.IsProtein,
			exampleInput.IsGrain,
			exampleInput.IsFruit,
			exampleInput.IsSalt,
			exampleInput.IsFat,
			exampleInput.IsAcid,
			exampleInput.IsHeat,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidIngredient.CreatedAt
		}

		actual, err := c.CreateValidIngredient(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredient.Name,
			exampleValidIngredient.Description,
			exampleValidIngredient.Warning,
			exampleValidIngredient.ContainsEgg,
			exampleValidIngredient.ContainsDairy,
			exampleValidIngredient.ContainsPeanut,
			exampleValidIngredient.ContainsTreeNut,
			exampleValidIngredient.ContainsSoy,
			exampleValidIngredient.ContainsWheat,
			exampleValidIngredient.ContainsShellfish,
			exampleValidIngredient.ContainsSesame,
			exampleValidIngredient.ContainsFish,
			exampleValidIngredient.ContainsGluten,
			exampleValidIngredient.AnimalFlesh,
			exampleValidIngredient.IsMeasuredVolumetrically,
			exampleValidIngredient.IsLiquid,
			exampleValidIngredient.IconPath,
			exampleValidIngredient.AnimalDerived,
			exampleValidIngredient.PluralName,
			exampleValidIngredient.RestrictToPreparations,
			exampleValidIngredient.MinimumIdealStorageTemperatureInCelsius,
			exampleValidIngredient.MaximumIdealStorageTemperatureInCelsius,
			exampleValidIngredient.StorageInstructions,
			exampleValidIngredient.Slug,
			exampleValidIngredient.ContainsAlcohol,
			exampleValidIngredient.ShoppingSuggestions,
			exampleValidIngredient.IsStarch,
			exampleValidIngredient.IsProtein,
			exampleValidIngredient.IsGrain,
			exampleValidIngredient.IsFruit,
			exampleValidIngredient.IsSalt,
			exampleValidIngredient.IsFat,
			exampleValidIngredient.IsAcid,
			exampleValidIngredient.IsHeat,
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidIngredient(ctx, exampleValidIngredient))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredient(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredient.Name,
			exampleValidIngredient.Description,
			exampleValidIngredient.Warning,
			exampleValidIngredient.ContainsEgg,
			exampleValidIngredient.ContainsDairy,
			exampleValidIngredient.ContainsPeanut,
			exampleValidIngredient.ContainsTreeNut,
			exampleValidIngredient.ContainsSoy,
			exampleValidIngredient.ContainsWheat,
			exampleValidIngredient.ContainsShellfish,
			exampleValidIngredient.ContainsSesame,
			exampleValidIngredient.ContainsFish,
			exampleValidIngredient.ContainsGluten,
			exampleValidIngredient.AnimalFlesh,
			exampleValidIngredient.IsMeasuredVolumetrically,
			exampleValidIngredient.IsLiquid,
			exampleValidIngredient.IconPath,
			exampleValidIngredient.AnimalDerived,
			exampleValidIngredient.PluralName,
			exampleValidIngredient.RestrictToPreparations,
			exampleValidIngredient.MinimumIdealStorageTemperatureInCelsius,
			exampleValidIngredient.MaximumIdealStorageTemperatureInCelsius,
			exampleValidIngredient.StorageInstructions,
			exampleValidIngredient.Slug,
			exampleValidIngredient.ContainsAlcohol,
			exampleValidIngredient.ShoppingSuggestions,
			exampleValidIngredient.IsStarch,
			exampleValidIngredient.IsProtein,
			exampleValidIngredient.IsGrain,
			exampleValidIngredient.IsFruit,
			exampleValidIngredient.IsSalt,
			exampleValidIngredient.IsFat,
			exampleValidIngredient.IsAcid,
			exampleValidIngredient.IsHeat,
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidIngredient(ctx, exampleValidIngredient))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidIngredient(ctx, exampleValidIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredient(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidIngredient(ctx, exampleValidIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkValidIngredientAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.MarkValidIngredientAsIndexed(ctx, exampleValidIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidIngredientAsIndexed(ctx, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkValidIngredientAsIndexed(ctx, exampleValidIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
