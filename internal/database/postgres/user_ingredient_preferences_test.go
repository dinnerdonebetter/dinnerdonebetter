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

func buildMockRowsFromUserIngredientPreferences(includeCounts bool, filteredCount uint64, userIngredientPreferences ...*types.UserIngredientPreference) *sqlmock.Rows {
	columns := userIngredientPreferencesTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range userIngredientPreferences {
		rowValues := []driver.Value{
			&x.ID,
			&x.Ingredient.ID,
			&x.Ingredient.Name,
			&x.Ingredient.Description,
			&x.Ingredient.Warning,
			&x.Ingredient.ContainsEgg,
			&x.Ingredient.ContainsDairy,
			&x.Ingredient.ContainsPeanut,
			&x.Ingredient.ContainsTreeNut,
			&x.Ingredient.ContainsSoy,
			&x.Ingredient.ContainsWheat,
			&x.Ingredient.ContainsShellfish,
			&x.Ingredient.ContainsSesame,
			&x.Ingredient.ContainsFish,
			&x.Ingredient.ContainsGluten,
			&x.Ingredient.AnimalFlesh,
			&x.Ingredient.IsMeasuredVolumetrically,
			&x.Ingredient.IsLiquid,
			&x.Ingredient.IconPath,
			&x.Ingredient.AnimalDerived,
			&x.Ingredient.PluralName,
			&x.Ingredient.RestrictToPreparations,
			&x.Ingredient.MinimumIdealStorageTemperatureInCelsius,
			&x.Ingredient.MaximumIdealStorageTemperatureInCelsius,
			&x.Ingredient.StorageInstructions,
			&x.Ingredient.Slug,
			&x.Ingredient.ContainsAlcohol,
			&x.Ingredient.ShoppingSuggestions,
			&x.Ingredient.IsStarch,
			&x.Ingredient.IsProtein,
			&x.Ingredient.IsGrain,
			&x.Ingredient.IsFruit,
			&x.Ingredient.IsSalt,
			&x.Ingredient.IsFat,
			&x.Ingredient.IsAcid,
			&x.Ingredient.IsHeat,
			&x.Ingredient.CreatedAt,
			&x.Ingredient.LastUpdatedAt,
			&x.Ingredient.ArchivedAt,
			&x.Rating,
			&x.Notes,
			&x.Allergy,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
			&x.BelongsToUser,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(userIngredientPreferences))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanUserIngredientPreferences(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanUserIngredientPreferences(ctx, mockRows, false)
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

		_, _, _, err := q.scanUserIngredientPreferences(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_UserIngredientPreferenceExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)
		args := []any{
			exampleUserIngredientPreference.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(userIngredientPreferenceExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.UserIngredientPreferenceExists(ctx, exampleUserIngredientPreference.ID, exampleUserID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIngredientPreferenceExists(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)
		args := []any{
			exampleUserIngredientPreference.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(userIngredientPreferenceExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.UserIngredientPreferenceExists(ctx, exampleUserIngredientPreference.ID, exampleUserID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)
		args := []any{
			exampleUserIngredientPreference.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(userIngredientPreferenceExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.UserIngredientPreferenceExists(ctx, exampleUserIngredientPreference.ID, exampleUserID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)

		args := []any{
			exampleUserIngredientPreference.ID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getUserIngredientPreferenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUserIngredientPreferences(false, 0, exampleUserIngredientPreference))

		actual, err := c.GetUserIngredientPreference(ctx, exampleUserIngredientPreference.ID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserIngredientPreference, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserIngredientPreference(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)

		args := []any{
			exampleUserIngredientPreference.ID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getUserIngredientPreferenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUserIngredientPreference(ctx, exampleUserIngredientPreference.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUserIngredientPreferences(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()
		exampleUserIngredientPreferenceList := fakes.BuildFakeUserIngredientPreferenceList()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_ingredient_preferences", []string{validIngredientsOnUserIngredientPreferencesJoin}, []string{"user_ingredient_preferences.id", "valid_ingredients.id"}, nil, userOwnershipColumn, userIngredientPreferencesTableColumns, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUserIngredientPreferences(true, exampleUserIngredientPreferenceList.FilteredCount, exampleUserIngredientPreferenceList.Data...))

		actual, err := c.GetUserIngredientPreferences(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserIngredientPreferenceList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		filter := (*types.QueryFilter)(nil)
		exampleUserIngredientPreferenceList := fakes.BuildFakeUserIngredientPreferenceList()
		exampleUserIngredientPreferenceList.Page = 0
		exampleUserIngredientPreferenceList.Limit = 0

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_ingredient_preferences", []string{validIngredientsOnUserIngredientPreferencesJoin}, []string{"user_ingredient_preferences.id", "valid_ingredients.id"}, nil, userOwnershipColumn, userIngredientPreferencesTableColumns, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUserIngredientPreferences(true, exampleUserIngredientPreferenceList.FilteredCount, exampleUserIngredientPreferenceList.Data...))

		actual, err := c.GetUserIngredientPreferences(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserIngredientPreferenceList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_ingredient_preferences", []string{validIngredientsOnUserIngredientPreferencesJoin}, []string{"user_ingredient_preferences.id", "valid_ingredients.id"}, nil, userOwnershipColumn, userIngredientPreferencesTableColumns, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUserIngredientPreferences(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "user_ingredient_preferences", []string{validIngredientsOnUserIngredientPreferencesJoin}, []string{"user_ingredient_preferences.id", "valid_ingredients.id"}, nil, userOwnershipColumn, userIngredientPreferencesTableColumns, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetUserIngredientPreferences(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
		exampleUserIngredientPreference.ID = "1"
		exampleUserIngredientPreference.Ingredient = types.ValidIngredient{ID: exampleUserIngredientPreference.Ingredient.ID}
		exampleInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceDatabaseCreationInput(exampleUserIngredientPreference)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.IngredientID,
			exampleInput.Rating,
			exampleInput.Notes,
			exampleInput.Allergy,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(userIngredientPreferenceCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleUserIngredientPreference.CreatedAt
		}

		actual, err := c.CreateUserIngredientPreference(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserIngredientPreference, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateUserIngredientPreference(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New("blah")
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
		exampleInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceDatabaseCreationInput(exampleUserIngredientPreference)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.IngredientID,
			exampleInput.Rating,
			exampleInput.Notes,
			exampleInput.Allergy,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(userIngredientPreferenceCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleUserIngredientPreference.CreatedAt
		}

		actual, err := c.CreateUserIngredientPreference(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)

		args := []any{
			exampleUserIngredientPreference.Ingredient.ID,
			exampleUserIngredientPreference.Rating,
			exampleUserIngredientPreference.Notes,
			exampleUserIngredientPreference.Allergy,
			exampleUserIngredientPreference.ID,
			exampleUserIngredientPreference.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserIngredientPreferenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateUserIngredientPreference(ctx, exampleUserIngredientPreference))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserIngredientPreference(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)

		args := []any{
			exampleUserIngredientPreference.Ingredient.ID,
			exampleUserIngredientPreference.Rating,
			exampleUserIngredientPreference.Notes,
			exampleUserIngredientPreference.Allergy,
			exampleUserIngredientPreference.ID,
			exampleUserIngredientPreference.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserIngredientPreferenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateUserIngredientPreference(ctx, exampleUserIngredientPreference))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)

		args := []any{
			exampleUserIngredientPreference.ID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveUserIngredientPreferenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveUserIngredientPreference(ctx, exampleUserIngredientPreference.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveUserIngredientPreference(ctx, "", exampleUserID))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		c, db := buildTestClient(t)

		args := []any{
			exampleUserIngredientPreference.ID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveUserIngredientPreferenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveUserIngredientPreference(ctx, exampleUserIngredientPreference.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
