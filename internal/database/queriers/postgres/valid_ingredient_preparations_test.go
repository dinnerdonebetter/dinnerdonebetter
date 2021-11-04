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

	database "github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromValidIngredientPreparations(includeCounts bool, filteredCount uint64, validIngredientPreparations ...*types.ValidIngredientPreparation) *sqlmock.Rows {
	columns := validIngredientPreparationsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredientPreparations {
		rowValues := []driver.Value{
			x.ID,
			x.Notes,
			x.ValidPreparationID,
			x.ValidIngredientID,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validIngredientPreparations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidIngredientPreparations(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidIngredientPreparations(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientPreparationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidIngredientPreparationExists(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientPreparationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientPreparationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidIngredientPreparationExists(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientPreparationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidIngredientPreparationExists(ctx, exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientPreparations(false, 0, exampleValidIngredientPreparation))

		actual, err := c.GetValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparation(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalValidIngredientPreparationCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidIngredientPreparationsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalValidIngredientPreparationCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidIngredientPreparationsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalValidIngredientPreparationCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_preparations", nil, nil, nil, householdOwnershipColumn, validIngredientPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientPreparations(true, exampleValidIngredientPreparationList.FilteredCount, exampleValidIngredientPreparationList.ValidIngredientPreparations...))

		actual, err := c.GetValidIngredientPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()
		exampleValidIngredientPreparationList.Page = 0
		exampleValidIngredientPreparationList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_preparations", nil, nil, nil, householdOwnershipColumn, validIngredientPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientPreparations(true, exampleValidIngredientPreparationList.FilteredCount, exampleValidIngredientPreparationList.ValidIngredientPreparations...))

		actual, err := c.GetValidIngredientPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_preparations", nil, nil, nil, householdOwnershipColumn, validIngredientPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_preparations", nil, nil, nil, householdOwnershipColumn, validIngredientPreparationsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredientPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientPreparationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		var exampleIDs []string
		for _, x := range exampleValidIngredientPreparationList.ValidIngredientPreparations {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidIngredientPreparationsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientPreparations(false, 0, exampleValidIngredientPreparationList.ValidIngredientPreparations...))

		actual, err := c.GetValidIngredientPreparationsWithIDs(ctx, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList.ValidIngredientPreparations, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		var exampleIDs []string
		for _, x := range exampleValidIngredientPreparationList.ValidIngredientPreparations {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidIngredientPreparationsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		var exampleIDs []string
		for _, x := range exampleValidIngredientPreparationList.ValidIngredientPreparations {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidIngredientPreparationsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.ID = "1"
		exampleInput := fakes.BuildFakeValidIngredientPreparationDatabaseCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidPreparationID,
			exampleInput.ValidIngredientID,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientPreparationCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleValidIngredientPreparation.ID))

		c.timeFunc = func() uint64 {
			return exampleValidIngredientPreparation.CreatedOn
		}

		actual, err := c.CreateValidIngredientPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientPreparation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		exampleInput := fakes.BuildFakeValidIngredientPreparationDatabaseCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidPreparationID,
			exampleInput.ValidIngredientID,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientPreparationCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleValidIngredientPreparation.CreatedOn
		}

		actual, err := c.CreateValidIngredientPreparation(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.ValidPreparationID,
			exampleValidIngredientPreparation.ValidIngredientID,
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleValidIngredientPreparation.ID))

		assert.NoError(t, c.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientPreparation(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.ValidPreparationID,
			exampleValidIngredientPreparation.ValidIngredientID,
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleValidIngredientPreparation.ID))

		assert.NoError(t, c.ArchiveValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientPreparation(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientPreparation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientPreparationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
