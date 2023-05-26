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

func buildMockRowsFromValidIngredientGroups(includeCounts bool, filteredCount uint64, validIngredientGroups ...*types.ValidIngredientGroup) *sqlmock.Rows {
	columns := validIngredientGroupsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredientGroups {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Description,
			x.Slug,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validIngredientGroups))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidIngredientGroups(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidIngredientGroups(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidIngredientGroupExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientGroup.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientGroupExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidIngredientGroupExists(ctx, exampleValidIngredientGroup.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientGroupExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientGroup.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientGroupExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidIngredientGroupExists(ctx, exampleValidIngredientGroup.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientGroup.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientGroupExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidIngredientGroupExists(ctx, exampleValidIngredientGroup.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getValidIngredientGroupArgs := []any{
			exampleValidIngredientGroup.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientGroupQuery)).
			WithArgs(interfaceToDriverValue(getValidIngredientGroupArgs)...).
			WillReturnRows(buildMockRowsFromValidIngredientGroups(false, 0, exampleValidIngredientGroup))

		actual, err := c.GetValidIngredientGroup(ctx, exampleValidIngredientGroup.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientGroup, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientGroup(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientGroup.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientGroupQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientGroup(ctx, exampleValidIngredientGroup.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidIngredientGroups(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroups := fakes.BuildFakeValidIngredientGroupList()

		ctx := context.Background()
		c, db := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientGroupSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientGroups(false, 0, exampleValidIngredientGroups.Data...))

		actual, err := c.SearchForValidIngredientGroups(ctx, exampleQuery, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientGroups.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		actual, err := c.SearchForValidIngredientGroups(ctx, "", filter)
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

		db.ExpectQuery(formatQueryForSQLMock(validIngredientGroupSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidIngredientGroups(ctx, exampleQuery, filter)
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

		db.ExpectQuery(formatQueryForSQLMock(validIngredientGroupSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidIngredientGroups(ctx, exampleQuery, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientGroupList := fakes.BuildFakeValidIngredientGroupList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_groups", nil, nil, nil, householdOwnershipColumn, validIngredientGroupsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientGroups(true, exampleValidIngredientGroupList.FilteredCount, exampleValidIngredientGroupList.Data...))

		actual, err := c.GetValidIngredientGroups(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientGroupList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidIngredientGroupList := fakes.BuildFakeValidIngredientGroupList()
		exampleValidIngredientGroupList.Page = 0
		exampleValidIngredientGroupList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_groups", nil, nil, nil, householdOwnershipColumn, validIngredientGroupsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientGroups(true, exampleValidIngredientGroupList.FilteredCount, exampleValidIngredientGroupList.Data...))

		actual, err := c.GetValidIngredientGroups(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientGroupList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_groups", nil, nil, nil, householdOwnershipColumn, validIngredientGroupsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientGroups(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_groups", nil, nil, nil, householdOwnershipColumn, validIngredientGroupsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredientGroups(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
		exampleValidIngredientGroup.ID = "1"
		exampleInput := converters.ConvertValidIngredientGroupToValidIngredientGroupDatabaseCreationInput(exampleValidIngredientGroup)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Slug,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientGroupCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidIngredientGroup.CreatedAt
		}

		actual, err := c.CreateValidIngredientGroup(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientGroup, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientGroup(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
		exampleInput := converters.ConvertValidIngredientGroupToValidIngredientGroupDatabaseCreationInput(exampleValidIngredientGroup)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Slug,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientGroupCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidIngredientGroup.CreatedAt
		}

		actual, err := c.CreateValidIngredientGroup(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientGroup.Name,
			exampleValidIngredientGroup.Description,
			exampleValidIngredientGroup.Slug,
			exampleValidIngredientGroup.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientGroupQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidIngredientGroup(ctx, exampleValidIngredientGroup))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientGroup(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientGroup.Name,
			exampleValidIngredientGroup.Description,
			exampleValidIngredientGroup.Slug,
			exampleValidIngredientGroup.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientGroupQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidIngredientGroup(ctx, exampleValidIngredientGroup))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientGroup.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientGroupQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidIngredientGroup(ctx, exampleValidIngredientGroup.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientGroup(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientGroup.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientGroupQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidIngredientGroup(ctx, exampleValidIngredientGroup.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
