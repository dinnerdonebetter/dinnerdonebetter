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
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromValidPreparations(includeCounts bool, filteredCount uint64, validPreparations ...*types.ValidPreparation) *sqlmock.Rows {
	columns := validPreparationsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validPreparations {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Description,
			x.IconPath,
			x.YieldsNothing,
			x.RestrictToIngredients,
			x.MinimumIngredientCount,
			x.MaximumIngredientCount,
			x.MinimumInstrumentCount,
			x.MaximumInstrumentCount,
			x.TemperatureRequired,
			x.TimeEstimateRequired,
			x.ConditionExpressionRequired,
			x.ConsumesVessel,
			x.OnlyForVessels,
			x.MinimumVesselCount,
			x.MaximumVesselCount,
			x.Slug,
			x.PastTense,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validPreparations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidPreparations(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidPreparations(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidPreparationExists(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidPreparationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidPreparationExists(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidPreparationExists(ctx, exampleValidPreparation.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparations(false, 0, exampleValidPreparation))

		actual, err := c.GetValidPreparation(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidPreparation(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparation(ctx, exampleValidPreparation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRandomValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparations(false, 0, exampleValidPreparation))

		actual, err := c.GetRandomValidPreparation(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRandomValidPreparation(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidPreparations(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparations := fakes.BuildFakeValidPreparationList()

		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparations(false, 0, exampleValidPreparations.Data...))

		actual, err := c.SearchForValidPreparations(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparations.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidPreparations(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidPreparations(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidPreparations(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparations", nil, nil, nil, householdOwnershipColumn, validPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparations(true, exampleValidPreparationList.FilteredCount, exampleValidPreparationList.Data...))

		actual, err := c.GetValidPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := (*types.QueryFilter)(nil)
		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()
		exampleValidPreparationList.Page = 0
		exampleValidPreparationList.Limit = 0

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparations", nil, nil, nil, householdOwnershipColumn, validPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparations(true, exampleValidPreparationList.FilteredCount, exampleValidPreparationList.Data...))

		actual, err := c.GetValidPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparations", nil, nil, nil, householdOwnershipColumn, validPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparations", nil, nil, nil, householdOwnershipColumn, validPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidPreparationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		exampleIDs := []string{}
		for _, exampleValidPreparation := range exampleValidPreparationList.Data {
			exampleIDs = append(exampleIDs, exampleValidPreparation.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		where := squirrel.Eq{"valid_preparations.id": exampleIDs}
		query, args := c.buildListQuery(ctx, validPreparationsTable, nil, nil, where, householdOwnershipColumn, validPreparationsTableColumns, "", false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparations(false, exampleValidPreparationList.FilteredCount, exampleValidPreparationList.Data...))

		actual, err := c.GetValidPreparationsWithIDs(ctx, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidPreparationThatNeedSearchIndexing(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		c, db := buildTestClient(t)

		exampleIDs := []string{}
		for _, exampleValidPreparation := range exampleValidPreparationList.Data {
			exampleIDs = append(exampleIDs, exampleValidPreparation.ID)
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationsNeedingIndexingQuery)).
			WithArgs(interfaceToDriverValue(nil)...).
			WillReturnRows(buildMockRowsFromIDs(exampleIDs...))

		actual, err := c.GetValidPreparationIDsThatNeedSearchIndexing(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleIDs, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleValidPreparation.ID = "1"
		exampleInput := converters.ConvertValidPreparationToValidPreparationDatabaseCreationInput(exampleValidPreparation)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.YieldsNothing,
			exampleInput.RestrictToIngredients,
			exampleInput.MinimumIngredientCount,
			exampleInput.MaximumIngredientCount,
			exampleInput.MinimumInstrumentCount,
			exampleInput.MaximumInstrumentCount,
			exampleInput.TemperatureRequired,
			exampleInput.TimeEstimateRequired,
			exampleInput.ConditionExpressionRequired,
			exampleInput.ConsumesVessel,
			exampleInput.OnlyForVessels,
			exampleInput.MinimumVesselCount,
			exampleInput.MaximumVesselCount,
			exampleInput.PastTense,
			exampleInput.Slug,
		}

		db.ExpectExec(formatQueryForSQLMock(validPreparationCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidPreparation.CreatedAt
		}

		actual, err := c.CreateValidPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleInput := converters.ConvertValidPreparationToValidPreparationDatabaseCreationInput(exampleValidPreparation)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.YieldsNothing,
			exampleInput.RestrictToIngredients,
			exampleInput.MinimumIngredientCount,
			exampleInput.MaximumIngredientCount,
			exampleInput.MinimumInstrumentCount,
			exampleInput.MaximumInstrumentCount,
			exampleInput.TemperatureRequired,
			exampleInput.TimeEstimateRequired,
			exampleInput.ConditionExpressionRequired,
			exampleInput.ConsumesVessel,
			exampleInput.OnlyForVessels,
			exampleInput.MinimumVesselCount,
			exampleInput.MaximumVesselCount,
			exampleInput.PastTense,
			exampleInput.Slug,
		}

		db.ExpectExec(formatQueryForSQLMock(validPreparationCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidPreparation.CreatedAt
		}

		actual, err := c.CreateValidPreparation(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.IconPath,
			exampleValidPreparation.YieldsNothing,
			exampleValidPreparation.RestrictToIngredients,
			exampleValidPreparation.MinimumIngredientCount,
			exampleValidPreparation.MaximumIngredientCount,
			exampleValidPreparation.MinimumInstrumentCount,
			exampleValidPreparation.MaximumInstrumentCount,
			exampleValidPreparation.TemperatureRequired,
			exampleValidPreparation.TimeEstimateRequired,
			exampleValidPreparation.ConditionExpressionRequired,
			exampleValidPreparation.ConsumesVessel,
			exampleValidPreparation.OnlyForVessels,
			exampleValidPreparation.MinimumVesselCount,
			exampleValidPreparation.MaximumVesselCount,
			exampleValidPreparation.Slug,
			exampleValidPreparation.PastTense,
			exampleValidPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidPreparation(ctx, exampleValidPreparation))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparation(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.IconPath,
			exampleValidPreparation.YieldsNothing,
			exampleValidPreparation.RestrictToIngredients,
			exampleValidPreparation.MinimumIngredientCount,
			exampleValidPreparation.MaximumIngredientCount,
			exampleValidPreparation.MinimumInstrumentCount,
			exampleValidPreparation.MaximumInstrumentCount,
			exampleValidPreparation.TemperatureRequired,
			exampleValidPreparation.TimeEstimateRequired,
			exampleValidPreparation.ConditionExpressionRequired,
			exampleValidPreparation.ConsumesVessel,
			exampleValidPreparation.OnlyForVessels,
			exampleValidPreparation.MinimumVesselCount,
			exampleValidPreparation.MaximumVesselCount,
			exampleValidPreparation.Slug,
			exampleValidPreparation.PastTense,
			exampleValidPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidPreparation(ctx, exampleValidPreparation))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidPreparation(ctx, exampleValidPreparation.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparation(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidPreparation(ctx, exampleValidPreparation.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkValidPreparationAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidPreparationLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.MarkValidPreparationAsIndexed(ctx, exampleValidPreparation.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidPreparationAsIndexed(ctx, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidPreparationLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkValidPreparationAsIndexed(ctx, exampleValidPreparation.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
