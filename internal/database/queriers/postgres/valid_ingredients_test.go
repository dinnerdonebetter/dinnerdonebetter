package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
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
			x.Variant,
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
			x.AnimalDerived,
			x.MeasurableByVolume,
			x.Icon,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
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
		args := []interface{}{
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
		args := []interface{}{
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
		args := []interface{}{
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

		args := []interface{}{
			exampleValidIngredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
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

		args := []interface{}{
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

func TestQuerier_GetTotalValidIngredientCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidIngredientsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalValidIngredientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidIngredientsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalValidIngredientCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

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

		query, args := c.buildListQuery(
			ctx,
			"valid_ingredients",
			nil,
			nil,
			accountOwnershipColumn,
			validIngredientsTableColumns,
			"",
			false,
			filter,
		)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredients(true, exampleValidIngredientList.FilteredCount, exampleValidIngredientList.ValidIngredients...))

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

		query, args := c.buildListQuery(
			ctx,
			"valid_ingredients",
			nil,
			nil,
			accountOwnershipColumn,
			validIngredientsTableColumns,
			"",
			false,
			filter,
		)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredients(true, exampleValidIngredientList.FilteredCount, exampleValidIngredientList.ValidIngredients...))

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

		query, args := c.buildListQuery(
			ctx,
			"valid_ingredients",
			nil,
			nil,
			accountOwnershipColumn,
			validIngredientsTableColumns,
			"",
			false,
			filter,
		)

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

		query, args := c.buildListQuery(
			ctx,
			"valid_ingredients",
			nil,
			nil,
			accountOwnershipColumn,
			validIngredientsTableColumns,
			"",
			false,
			filter,
		)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		var exampleIDs []string
		for _, x := range exampleValidIngredientList.ValidIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidIngredientsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredients(false, 0, exampleValidIngredientList.ValidIngredients...))

		actual, err := c.GetValidIngredientsWithIDs(ctx, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList.ValidIngredients, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientsWithIDs(ctx, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		var exampleIDs []string
		for _, x := range exampleValidIngredientList.ValidIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidIngredientsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		var exampleIDs []string
		for _, x := range exampleValidIngredientList.ValidIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidIngredientsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredientsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		exampleValidIngredient.ID = "1"
		exampleInput := fakes.BuildFakeValidIngredientDatabaseCreationInputFromValidIngredient(exampleValidIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Variant,
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
			exampleInput.AnimalDerived,
			exampleInput.MeasurableByVolume,
			exampleInput.Icon,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleValidIngredient.ID))

		c.timeFunc = func() uint64 {
			return exampleValidIngredient.CreatedOn
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
		exampleInput := fakes.BuildFakeValidIngredientDatabaseCreationInputFromValidIngredient(exampleValidIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Variant,
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
			exampleInput.AnimalDerived,
			exampleInput.MeasurableByVolume,
			exampleInput.Icon,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleValidIngredient.CreatedOn
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

		args := []interface{}{
			exampleValidIngredient.Name,
			exampleValidIngredient.Variant,
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
			exampleValidIngredient.AnimalDerived,
			exampleValidIngredient.MeasurableByVolume,
			exampleValidIngredient.Icon,
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleValidIngredient.ID))

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

		args := []interface{}{
			exampleValidIngredient.Name,
			exampleValidIngredient.Variant,
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
			exampleValidIngredient.AnimalDerived,
			exampleValidIngredient.MeasurableByVolume,
			exampleValidIngredient.Icon,
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

		args := []interface{}{
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleValidIngredient.ID))

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

		args := []interface{}{
			exampleValidIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidIngredient(ctx, exampleValidIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
