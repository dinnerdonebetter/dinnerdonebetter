package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromValidMeasurementUnitConversions(includeCounts bool, filteredCount uint64, validMeasurementUnitConversions ...*types.ValidMeasurementUnitConversion) *sqlmock.Rows {
	columns := validMeasurementUnitConversionsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validMeasurementUnitConversions {
		ingredient := &types.NullableValidIngredient{}
		if x.OnlyForIngredient != nil {
			ingredient = converters.ConvertValidIngredientToNullableValidIngredient(x.OnlyForIngredient)
		}

		rowValues := []driver.Value{
			x.ID,

			// first valid measurement join
			x.From.ID,
			x.From.Name,
			x.From.Description,
			x.From.Volumetric,
			x.From.IconPath,
			x.From.Universal,
			x.From.Metric,
			x.From.Imperial,
			x.From.Slug,
			x.From.PluralName,
			x.From.CreatedAt,
			x.From.LastUpdatedAt,
			x.From.ArchivedAt,

			// second valid measurement join
			x.To.ID,
			x.To.Name,
			x.To.Description,
			x.To.Volumetric,
			x.To.IconPath,
			x.To.Universal,
			x.To.Metric,
			x.To.Imperial,
			x.To.Slug,
			x.To.PluralName,
			x.To.CreatedAt,
			x.To.LastUpdatedAt,
			x.To.ArchivedAt,

			// valid ingredient join
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Description,
			&ingredient.Warning,
			&ingredient.ContainsEgg,
			&ingredient.ContainsDairy,
			&ingredient.ContainsPeanut,
			&ingredient.ContainsTreeNut,
			&ingredient.ContainsSoy,
			&ingredient.ContainsWheat,
			&ingredient.ContainsShellfish,
			&ingredient.ContainsSesame,
			&ingredient.ContainsFish,
			&ingredient.ContainsGluten,
			&ingredient.AnimalFlesh,
			&ingredient.IsMeasuredVolumetrically,
			&ingredient.IsLiquid,
			&ingredient.IconPath,
			&ingredient.AnimalDerived,
			&ingredient.PluralName,
			&ingredient.RestrictToPreparations,
			&ingredient.MinimumIdealStorageTemperatureInCelsius,
			&ingredient.MaximumIdealStorageTemperatureInCelsius,
			&ingredient.StorageInstructions,
			&ingredient.Slug,
			&ingredient.ContainsAlcohol,
			&ingredient.ShoppingSuggestions,
			&ingredient.IsStarch,
			&ingredient.IsProtein,
			&ingredient.IsGrain,
			&ingredient.IsFruit,
			&ingredient.IsSalt,
			&ingredient.IsFat,
			&ingredient.IsAcid,
			&ingredient.IsHeat,
			&ingredient.CreatedAt,
			&ingredient.LastUpdatedAt,
			&ingredient.ArchivedAt,

			// rest of the valid measurement conversion
			x.Modifier,
			x.Notes,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validMeasurementUnitConversions))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ValidMeasurementUnitConversionExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementUnitConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitConversionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidMeasurementUnitConversionExists(ctx, exampleValidMeasurementUnitConversion.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidMeasurementUnitConversionExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementUnitConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitConversionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidMeasurementUnitConversionExists(ctx, exampleValidMeasurementUnitConversion.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementUnitConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitConversionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidMeasurementUnitConversionExists(ctx, exampleValidMeasurementUnitConversion.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementUnitConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitConversionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnitConversions(false, 0, exampleValidMeasurementUnitConversion))

		actual, err := c.GetValidMeasurementUnitConversion(ctx, exampleValidMeasurementUnitConversion.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitConversion, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidMeasurementUnitConversion(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementUnitConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitConversionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidMeasurementUnitConversion(ctx, exampleValidMeasurementUnitConversion.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementUnitConversionsFromUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleValidMeasurementUnitConversions := fakes.BuildFakeValidMeasurementUnitConversionList().Data

		c, db := buildTestClient(t)

		getValidMeasurementUnitConversionsFromUnitArgs := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitConversionsFromUnitQuery)).
			WithArgs(interfaceToDriverValue(getValidMeasurementUnitConversionsFromUnitArgs)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnitConversions(false, 0, exampleValidMeasurementUnitConversions...))

		actual, err := c.GetValidMeasurementUnitConversionsFromUnit(ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitConversions, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementUnitConversionsToUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleValidMeasurementUnitConversions := fakes.BuildFakeValidMeasurementUnitConversionList().Data

		c, db := buildTestClient(t)

		getValidMeasurementUnitConversionsToUnitArgs := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitConversionsToUnitQuery)).
			WithArgs(interfaceToDriverValue(getValidMeasurementUnitConversionsToUnitArgs)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnitConversions(false, 0, exampleValidMeasurementUnitConversions...))

		actual, err := c.GetValidMeasurementUnitConversionsToUnit(ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitConversions, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
		exampleValidMeasurementUnitConversion.ID = "1"
		exampleValidMeasurementUnitConversion.From = types.ValidMeasurementUnit{ID: exampleValidMeasurementUnitConversion.From.ID}
		exampleValidMeasurementUnitConversion.To = types.ValidMeasurementUnit{ID: exampleValidMeasurementUnitConversion.To.ID}
		if exampleValidMeasurementUnitConversion.OnlyForIngredient != nil {
			exampleValidMeasurementUnitConversion.OnlyForIngredient = &types.ValidIngredient{ID: exampleValidMeasurementUnitConversion.OnlyForIngredient.ID}
		}
		exampleInput := converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionDatabaseCreationInput(exampleValidMeasurementUnitConversion)

		c, db := buildTestClient(t)

		validMeasurementUnitConversionCreationArgs := []any{
			exampleInput.ID,
			exampleInput.From,
			exampleInput.To,
			exampleInput.OnlyForIngredient,
			exampleInput.Modifier,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementUnitConversionCreationQuery)).
			WithArgs(interfaceToDriverValue(validMeasurementUnitConversionCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidMeasurementUnitConversion.CreatedAt
		}

		actual, err := c.CreateValidMeasurementUnitConversion(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitConversion, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidMeasurementUnitConversion(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
		exampleInput := converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionDatabaseCreationInput(exampleValidMeasurementUnitConversion)

		c, db := buildTestClient(t)

		validMeasurementUnitConversionCreationArgs := []any{
			exampleInput.ID,
			exampleInput.From,
			exampleInput.To,
			exampleInput.OnlyForIngredient,
			exampleInput.Modifier,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementUnitConversionCreationQuery)).
			WithArgs(interfaceToDriverValue(validMeasurementUnitConversionCreationArgs)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidMeasurementUnitConversion.CreatedAt
		}

		actual, err := c.CreateValidMeasurementUnitConversion(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		c, db := buildTestClient(t)

		var ingredientID *string
		if exampleValidMeasurementUnitConversion.OnlyForIngredient != nil {
			ingredientID = &exampleValidMeasurementUnitConversion.OnlyForIngredient.ID
		}

		updateValidMeasurementUnitConversionArgs := []any{
			exampleValidMeasurementUnitConversion.From.ID,
			exampleValidMeasurementUnitConversion.To.ID,
			ingredientID,
			exampleValidMeasurementUnitConversion.Modifier,
			exampleValidMeasurementUnitConversion.Notes,
			exampleValidMeasurementUnitConversion.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidMeasurementUnitConversionQuery)).
			WithArgs(interfaceToDriverValue(updateValidMeasurementUnitConversionArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidMeasurementUnitConversion(ctx, exampleValidMeasurementUnitConversion))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidMeasurementUnitConversion(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		c, db := buildTestClient(t)

		var ingredientID *string
		if exampleValidMeasurementUnitConversion.OnlyForIngredient != nil {
			ingredientID = &exampleValidMeasurementUnitConversion.OnlyForIngredient.ID
		}

		updateValidMeasurementUnitConversionArgs := []any{
			exampleValidMeasurementUnitConversion.From.ID,
			exampleValidMeasurementUnitConversion.To.ID,
			ingredientID,
			exampleValidMeasurementUnitConversion.Modifier,
			exampleValidMeasurementUnitConversion.Notes,
			exampleValidMeasurementUnitConversion.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidMeasurementUnitConversionQuery)).
			WithArgs(interfaceToDriverValue(updateValidMeasurementUnitConversionArgs)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidMeasurementUnitConversion(ctx, exampleValidMeasurementUnitConversion))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidMeasurementUnitConversion(ctx, ""))
	})
}
