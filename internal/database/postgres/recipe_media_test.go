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

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeMedia.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeMediaExists(ctx, exampleRecipeMedia.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeMediaExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeMedia.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeMediaExists(ctx, exampleRecipeMedia.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeMedia.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeMediaExists(ctx, exampleRecipeMedia.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeMedia(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)

		args := []interface{}{
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

		args := []interface{}{
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
		exampleRecipeMediaList := fakes.BuildFakeRecipeMediaList().RecipeMedia

		c, db := buildTestClient(t)

		recipeMediaForRecipeArgs := []interface{}{
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

		recipeMediaForRecipeArgs := []interface{}{
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

		recipeMediaForRecipeArgs := []interface{}{
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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()
		exampleRecipeMedia.ID = "1"
		exampleInput := converters.ConvertRecipeMediaToRecipeMediaDatabaseCreationInput(exampleRecipeMedia)

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.BelongsToRecipe,
			exampleInput.BelongsToRecipeStep,
			exampleInput.MimeType,
			exampleInput.InternalPath,
			exampleInput.ExternalPath,
			exampleInput.Index,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeMediaCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeMedia.CreatedAt
		}

		actual, err := c.CreateRecipeMedia(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeMedia, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeMedia(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()
		exampleInput := converters.ConvertRecipeMediaToRecipeMediaDatabaseCreationInput(exampleRecipeMedia)

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.BelongsToRecipe,
			exampleInput.BelongsToRecipeStep,
			exampleInput.MimeType,
			exampleInput.InternalPath,
			exampleInput.ExternalPath,
			exampleInput.Index,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeMediaCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeMedia.CreatedAt
		}

		actual, err := c.CreateRecipeMedia(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeMedia(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeMedia.BelongsToRecipe,
			exampleRecipeMedia.BelongsToRecipeStep,
			exampleRecipeMedia.MimeType,
			exampleRecipeMedia.InternalPath,
			exampleRecipeMedia.ExternalPath,
			exampleRecipeMedia.Index,
			exampleRecipeMedia.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeMediaQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateRecipeMedia(ctx, exampleRecipeMedia))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeMedia(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeMedia.BelongsToRecipe,
			exampleRecipeMedia.BelongsToRecipeStep,
			exampleRecipeMedia.MimeType,
			exampleRecipeMedia.InternalPath,
			exampleRecipeMedia.ExternalPath,
			exampleRecipeMedia.Index,
			exampleRecipeMedia.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeMediaQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipeMedia(ctx, exampleRecipeMedia))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipeMedia(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeMedia.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeMediaQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveRecipeMedia(ctx, exampleRecipeMedia.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeMedia(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipeMedia := fakes.BuildFakeRecipeMedia()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeMedia.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeMediaQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveRecipeMedia(ctx, exampleRecipeMedia.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
