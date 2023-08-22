package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromRecipeMedia(recipeMedia ...*types.RecipeMedia) *sqlmock.Rows {
	columns := recipeMediaTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeMedia {
		rowValues := []driver.Value{
			x.ID,
			x.BelongsToRecipe,
			x.BelongsToRecipeStep,
			x.MimeType,
			x.InternalPath,
			x.ExternalPath,
			x.Index,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeMedias(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := q.scanRecipeMedia(ctx, mockRows)
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

		_, err := q.scanRecipeMedia(ctx, mockRows)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeMediaExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeMediaExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeMedia(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeMedia.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeMediaQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeMedia(exampleRecipeMedia))

		actual, err := c.GetRecipeMedia(ctx, exampleRecipeMedia.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeMedia, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeMedia(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeMedia.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeMediaQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeMedia(ctx, exampleRecipeMedia.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeMediaForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeMediaList := fakes.BuildFakeRecipeMediaList().Data

		c, db := buildTestClient(t)

		recipeMediaForRecipeArgs := []any{
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaForRecipeQuery)).
			WithArgs(interfaceToDriverValue(recipeMediaForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeMedia(exampleRecipeMediaList...))

		actual, err := c.getRecipeMediaForRecipe(ctx, exampleRecipeID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeMediaList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		recipeMediaForRecipeArgs := []any{
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaForRecipeQuery)).
			WithArgs(interfaceToDriverValue(recipeMediaForRecipeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.getRecipeMediaForRecipe(ctx, exampleRecipeID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		recipeMediaForRecipeArgs := []any{
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaForRecipeQuery)).
			WithArgs(interfaceToDriverValue(recipeMediaForRecipeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.getRecipeMediaForRecipe(ctx, exampleRecipeID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeMedia(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeMedia(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeMedia(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeMedia(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeMedia(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeMedia(ctx, ""))
	})
}
