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

func buildMockRowsFromValidIngredientStates(includeCounts bool, filteredCount uint64, validIngredientStates ...*types.ValidIngredientState) *sqlmock.Rows {
	columns := validIngredientStatesTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredientStates {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Description,
			x.IconPath,
			x.Slug,
			x.PastTense,
			x.AttributeType,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validIngredientStates))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidIngredientStates(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidIngredientStates(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidIngredientStateExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidIngredientStateExists(ctx, exampleValidIngredientState.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientStateExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidIngredientStateExists(ctx, exampleValidIngredientState.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidIngredientStateExists(ctx, exampleValidIngredientState.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientStateQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientStates(false, 0, exampleValidIngredientState))

		actual, err := c.GetValidIngredientState(ctx, exampleValidIngredientState.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientState, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientState(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientStateQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientState(ctx, exampleValidIngredientState.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidIngredientStates(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientStates := fakes.BuildFakeValidIngredientStateList()

		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientStates(false, 0, exampleValidIngredientStates.Data...))

		actual, err := c.SearchForValidIngredientStates(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStates.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidIngredientStates(ctx, "")
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

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidIngredientStates(ctx, exampleQuery)
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

		db.ExpectQuery(formatQueryForSQLMock(validIngredientStateSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidIngredientStates(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleValidIngredientStateList := fakes.BuildFakeValidIngredientStateList()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_states", nil, nil, nil, householdOwnershipColumn, validIngredientStatesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientStates(true, exampleValidIngredientStateList.FilteredCount, exampleValidIngredientStateList.Data...))

		actual, err := c.GetValidIngredientStates(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStateList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := (*types.QueryFilter)(nil)
		exampleValidIngredientStateList := fakes.BuildFakeValidIngredientStateList()
		exampleValidIngredientStateList.Page = 0
		exampleValidIngredientStateList.Limit = 0

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_states", nil, nil, nil, householdOwnershipColumn, validIngredientStatesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientStates(true, exampleValidIngredientStateList.FilteredCount, exampleValidIngredientStateList.Data...))

		actual, err := c.GetValidIngredientStates(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStateList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_states", nil, nil, nil, householdOwnershipColumn, validIngredientStatesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientStates(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_states", nil, nil, nil, householdOwnershipColumn, validIngredientStatesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredientStates(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()
		exampleValidIngredientState.ID = "1"
		exampleInput := converters.ConvertValidIngredientStateToValidIngredientStateDatabaseCreationInput(exampleValidIngredientState)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.PastTense,
			exampleInput.Slug,
			exampleInput.AttributeType,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientStateCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidIngredientState.CreatedAt
		}

		actual, err := c.CreateValidIngredientState(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientState, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientState(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()
		exampleInput := converters.ConvertValidIngredientStateToValidIngredientStateDatabaseCreationInput(exampleValidIngredientState)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.PastTense,
			exampleInput.Slug,
			exampleInput.AttributeType,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientStateCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidIngredientState.CreatedAt
		}

		actual, err := c.CreateValidIngredientState(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.Name,
			exampleValidIngredientState.Description,
			exampleValidIngredientState.IconPath,
			exampleValidIngredientState.Slug,
			exampleValidIngredientState.PastTense,
			exampleValidIngredientState.AttributeType,
			exampleValidIngredientState.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientStateQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidIngredientState(ctx, exampleValidIngredientState))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientState(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.Name,
			exampleValidIngredientState.Description,
			exampleValidIngredientState.IconPath,
			exampleValidIngredientState.Slug,
			exampleValidIngredientState.PastTense,
			exampleValidIngredientState.AttributeType,
			exampleValidIngredientState.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientStateQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidIngredientState(ctx, exampleValidIngredientState))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientStateQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidIngredientState(ctx, exampleValidIngredientState.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientState(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientStateQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidIngredientState(ctx, exampleValidIngredientState.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkValidIngredientStateAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientStateLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.MarkValidIngredientStateAsIndexed(ctx, exampleValidIngredientState.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidIngredientStateAsIndexed(ctx, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidIngredientState.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientStateLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkValidIngredientStateAsIndexed(ctx, exampleValidIngredientState.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
