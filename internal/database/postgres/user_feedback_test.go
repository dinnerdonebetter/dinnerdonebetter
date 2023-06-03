package postgres

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mustJSON(t *testing.T, x any) string {
	t.Helper()

	out, err := json.Marshal(x)
	require.NoError(t, err)

	return string(out)
}

func buildMockRowsFromUserFeedbacks(t *testing.T, includeCounts bool, filteredCount uint64, userFeedbacks ...*types.UserFeedback) *sqlmock.Rows {
	t.Helper()

	columns := userFeedbacksTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range userFeedbacks {
		rowValues := []driver.Value{
			&x.ID,
			&x.Prompt,
			&x.Feedback,
			&x.Rating,
			mustJSON(t, x.Context),
			&x.ByUser,
			&x.CreatedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(userFeedbacks))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanUserFeedback(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanUserFeedbacks(ctx, mockRows, false)
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

		_, _, _, err := q.scanUserFeedbacks(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_GetUserFeedback(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUserFeedbackList := fakes.BuildFakeUserFeedbackList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_feedback", nil, nil, nil, householdOwnershipColumn, userFeedbacksTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUserFeedbacks(t, true, exampleUserFeedbackList.FilteredCount, exampleUserFeedbackList.Data...))

		actual, err := c.GetUserFeedback(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserFeedbackList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleUserFeedbackList := fakes.BuildFakeUserFeedbackList()
		exampleUserFeedbackList.Page = 0
		exampleUserFeedbackList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_feedback", nil, nil, nil, householdOwnershipColumn, userFeedbacksTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUserFeedbacks(t, true, exampleUserFeedbackList.FilteredCount, exampleUserFeedbackList.Data...))

		actual, err := c.GetUserFeedback(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserFeedbackList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_feedback", nil, nil, nil, householdOwnershipColumn, userFeedbacksTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUserFeedback(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_feedback", nil, nil, nil, householdOwnershipColumn, userFeedbacksTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetUserFeedback(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateUserFeedback(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserFeedback := fakes.BuildFakeUserFeedback()
		exampleUserFeedback.ID = "1"
		exampleInput := converters.ConvertUserFeedbackToUserFeedbackDatabaseCreationInput(exampleUserFeedback)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Prompt,
			exampleInput.Feedback,
			exampleInput.Rating,
			mustJSON(t, exampleInput.Context),
			exampleInput.ByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(userFeedbackCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleUserFeedback.CreatedAt
		}

		actual, err := c.CreateUserFeedback(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserFeedback, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateUserFeedback(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleUserFeedback := fakes.BuildFakeUserFeedback()
		exampleInput := converters.ConvertUserFeedbackToUserFeedbackDatabaseCreationInput(exampleUserFeedback)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Prompt,
			exampleInput.Feedback,
			exampleInput.Rating,
			mustJSON(t, exampleInput.Context),
			exampleInput.ByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(userFeedbackCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleUserFeedback.CreatedAt
		}

		actual, err := c.CreateUserFeedback(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}
