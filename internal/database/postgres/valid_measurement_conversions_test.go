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

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func buildMockRowsFromValidMeasurementConversions(includeCounts bool, filteredCount uint64, validMeasurementConversions ...*types.ValidMeasurementConversion) *sqlmock.Rows {
	columns := validMeasurementConversionsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validMeasurementConversions {
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
			rowValues = append(rowValues, filteredCount, len(validMeasurementConversions))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ValidMeasurementConversionExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementConversionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidMeasurementConversionExists(ctx, exampleValidMeasurementConversion.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidMeasurementConversionExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementConversionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidMeasurementConversionExists(ctx, exampleValidMeasurementConversion.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementConversionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidMeasurementConversionExists(ctx, exampleValidMeasurementConversion.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementConversionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementConversions(false, 0, exampleValidMeasurementConversion))

		actual, err := c.GetValidMeasurementConversion(ctx, exampleValidMeasurementConversion.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementConversion, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidMeasurementConversion(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementConversionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidMeasurementConversion(ctx, exampleValidMeasurementConversion.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementConversionsFromUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleValidMeasurementConversions := fakes.BuildFakeValidMeasurementConversionList().ValidMeasurementConversions

		c, db := buildTestClient(t)

		getValidMeasurementConversionsFromUnitArgs := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementConversionsFromUnitQuery)).
			WithArgs(interfaceToDriverValue(getValidMeasurementConversionsFromUnitArgs)...).
			WillReturnRows(buildMockRowsFromValidMeasurementConversions(false, 0, exampleValidMeasurementConversions...))

		actual, err := c.GetValidMeasurementConversionsFromUnit(ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementConversions, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementConversionsToUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleValidMeasurementConversions := fakes.BuildFakeValidMeasurementConversionList().ValidMeasurementConversions

		c, db := buildTestClient(t)

		getValidMeasurementConversionsToUnitArgs := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementConversionsToUnitQuery)).
			WithArgs(interfaceToDriverValue(getValidMeasurementConversionsToUnitArgs)...).
			WillReturnRows(buildMockRowsFromValidMeasurementConversions(false, 0, exampleValidMeasurementConversions...))

		actual, err := c.GetValidMeasurementConversionsToUnit(ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementConversions, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidMeasurementConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()
		exampleValidMeasurementConversion.ID = "1"
		exampleValidMeasurementConversion.From = types.ValidMeasurementUnit{ID: exampleValidMeasurementConversion.From.ID}
		exampleValidMeasurementConversion.To = types.ValidMeasurementUnit{ID: exampleValidMeasurementConversion.To.ID}
		if exampleValidMeasurementConversion.OnlyForIngredient != nil {
			exampleValidMeasurementConversion.OnlyForIngredient = &types.ValidIngredient{ID: exampleValidMeasurementConversion.OnlyForIngredient.ID}
		}
		exampleInput := converters.ConvertValidMeasurementConversionToValidMeasurementConversionDatabaseCreationInput(exampleValidMeasurementConversion)

		c, db := buildTestClient(t)

		validMeasurementConversionCreationArgs := []any{
			exampleInput.ID,
			exampleInput.From,
			exampleInput.To,
			exampleInput.ForIngredient,
			exampleInput.Modifier,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementConversionCreationQuery)).
			WithArgs(interfaceToDriverValue(validMeasurementConversionCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidMeasurementConversion.CreatedAt
		}

		actual, err := c.CreateValidMeasurementConversion(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementConversion, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidMeasurementConversion(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()
		exampleInput := converters.ConvertValidMeasurementConversionToValidMeasurementConversionDatabaseCreationInput(exampleValidMeasurementConversion)

		c, db := buildTestClient(t)

		validMeasurementConversionCreationArgs := []any{
			exampleInput.ID,
			exampleInput.From,
			exampleInput.To,
			exampleInput.ForIngredient,
			exampleInput.Modifier,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementConversionCreationQuery)).
			WithArgs(interfaceToDriverValue(validMeasurementConversionCreationArgs)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidMeasurementConversion.CreatedAt
		}

		actual, err := c.CreateValidMeasurementConversion(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidMeasurementConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)

		var ingredientID *string
		if exampleValidMeasurementConversion.OnlyForIngredient != nil {
			ingredientID = &exampleValidMeasurementConversion.OnlyForIngredient.ID
		}

		updateValidMeasurementConversionArgs := []any{
			exampleValidMeasurementConversion.From.ID,
			exampleValidMeasurementConversion.To.ID,
			ingredientID,
			exampleValidMeasurementConversion.Modifier,
			exampleValidMeasurementConversion.Notes,
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidMeasurementConversionQuery)).
			WithArgs(interfaceToDriverValue(updateValidMeasurementConversionArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidMeasurementConversion(ctx, exampleValidMeasurementConversion))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidMeasurementConversion(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)

		var ingredientID *string
		if exampleValidMeasurementConversion.OnlyForIngredient != nil {
			ingredientID = &exampleValidMeasurementConversion.OnlyForIngredient.ID
		}

		updateValidMeasurementConversionArgs := []any{
			exampleValidMeasurementConversion.From.ID,
			exampleValidMeasurementConversion.To.ID,
			ingredientID,
			exampleValidMeasurementConversion.Modifier,
			exampleValidMeasurementConversion.Notes,
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidMeasurementConversionQuery)).
			WithArgs(interfaceToDriverValue(updateValidMeasurementConversionArgs)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidMeasurementConversion(ctx, exampleValidMeasurementConversion))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidMeasurementConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidMeasurementConversionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidMeasurementConversion(ctx, exampleValidMeasurementConversion.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidMeasurementConversion(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementConversion.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidMeasurementConversionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidMeasurementConversion(ctx, exampleValidMeasurementConversion.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
