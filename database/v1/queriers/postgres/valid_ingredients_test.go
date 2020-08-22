package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromValidIngredients(excludeCount bool, validIngredients ...*models.ValidIngredient) *sqlmock.Rows {
	includeCount := len(validIngredients) > 1 && !excludeCount
	columns := validIngredientsTableColumns

	if includeCount {
		columns = append(columns, "count")
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

		if includeCount {
			rowValues = append(rowValues, len(validIngredients))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromValidIngredient(x *models.ValidIngredient) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(validIngredientsTableColumns).AddRow(
		x.ArchivedOn,
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
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanValidIngredients(mockRows, true)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanValidIngredients(mockRows, true)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildValidIngredientExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		expectedQuery := "SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := p.buildValidIngredientExistsQuery(exampleValidIngredient.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ValidIngredientExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ValidIngredientExists(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ValidIngredientExists(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.id = $1"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientQuery(exampleValidIngredient.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredient(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
			).
			WillReturnRows(buildMockRowsFromValidIngredients(
				false, exampleValidIngredient))

		actual, err := p.GetValidIngredient(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredient(ctx, exampleValidIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllValidIngredientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL"
		actualQuery := p.buildGetAllValidIngredientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllValidIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllValidIngredientsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfValidIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.id > $1 AND valid_ingredients.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfValidIngredientsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllValidIngredients(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL"
	expectedGetQuery := "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.id > $1 AND valid_ingredients.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidIngredientList := fakemodels.BuildFakeValidIngredientList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromValidIngredients(
					true,
					&exampleValidIngredientList.ValidIngredients[0],
					&exampleValidIngredientList.ValidIngredients[1],
					&exampleValidIngredientList.ValidIngredients[2],
				),
			)

		out := make(chan []models.ValidIngredient)
		doneChan := make(chan bool, 1)

		err := p.GetAllValidIngredients(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.ValidIngredient)

		err := p.GetAllValidIngredients(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []models.ValidIngredient)

		err := p.GetAllValidIngredients(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.ValidIngredient)

		err := p.GetAllValidIngredients(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromValidIngredient(exampleValidIngredient))

		out := make(chan []models.ValidIngredient)

		err := p.GetAllValidIngredients(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on, (SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.created_on > $1 AND valid_ingredients.created_on < $2 AND valid_ingredients.last_updated_on > $3 AND valid_ingredients.last_updated_on < $4 ORDER BY valid_ingredients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on, (SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL ORDER BY valid_ingredients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredientList := fakemodels.BuildFakeValidIngredientList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromValidIngredients(
					false,
					&exampleValidIngredientList.ValidIngredients[0],
					&exampleValidIngredientList.ValidIngredients[1],
					&exampleValidIngredientList.ValidIngredients[2],
				),
			)

		actual, err := p.GetValidIngredients(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetValidIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid ingredient", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromValidIngredient(exampleValidIngredient))

		actual, err := p.GetValidIngredients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM (SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredients WHERE valid_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetValidIngredientsWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientList := fakemodels.BuildFakeValidIngredientList()
		var exampleIDs []uint64
		for _, validIngredient := range exampleValidIngredientList.ValidIngredients {
			exampleIDs = append(exampleIDs, validIngredient.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM (SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredients WHERE valid_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromValidIngredients(
					true,
					&exampleValidIngredientList.ValidIngredients[0],
					&exampleValidIngredientList.ValidIngredients[1],
					&exampleValidIngredientList.ValidIngredients[2],
				),
			)

		actual, err := p.GetValidIngredientsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList.ValidIngredients, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM (SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredients WHERE valid_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM (SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredients WHERE valid_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetValidIngredientsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid ingredient", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM (SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.measurable_by_volume, valid_ingredients.icon, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredients WHERE valid_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromValidIngredient(exampleValidIngredient))

		actual, err := p.GetValidIngredientsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		expectedQuery := "INSERT INTO valid_ingredients (name,variant,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,animal_derived,measurable_by_volume,icon) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING id, created_on"
		expectedArgs := []interface{}{
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
		}
		actualQuery, actualArgs := p.buildCreateValidIngredientQuery(exampleValidIngredient)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO valid_ingredients (name,variant,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,animal_derived,measurable_by_volume,icon) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleValidIngredient.ID, exampleValidIngredient.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
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
			).WillReturnRows(exampleRows)

		actual, err := p.CreateValidIngredient(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
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
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateValidIngredient(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		expectedQuery := "UPDATE valid_ingredients SET name = $1, variant = $2, description = $3, warning = $4, contains_egg = $5, contains_dairy = $6, contains_peanut = $7, contains_tree_nut = $8, contains_soy = $9, contains_wheat = $10, contains_shellfish = $11, contains_sesame = $12, contains_fish = $13, contains_gluten = $14, animal_flesh = $15, animal_derived = $16, measurable_by_volume = $17, icon = $18, last_updated_on = extract(epoch FROM NOW()) WHERE id = $19 RETURNING last_updated_on"
		expectedArgs := []interface{}{
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
		actualQuery, actualArgs := p.buildUpdateValidIngredientQuery(exampleValidIngredient)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_ingredients SET name = $1, variant = $2, description = $3, warning = $4, contains_egg = $5, contains_dairy = $6, contains_peanut = $7, contains_tree_nut = $8, contains_soy = $9, contains_wheat = $10, contains_shellfish = $11, contains_sesame = $12, contains_fish = $13, contains_gluten = $14, animal_flesh = $15, animal_derived = $16, measurable_by_volume = $17, icon = $18, last_updated_on = extract(epoch FROM NOW()) WHERE id = $19 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
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
			).WillReturnRows(exampleRows)

		err := p.UpdateValidIngredient(ctx, exampleValidIngredient)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
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
			).WillReturnError(errors.New("blah"))

		err := p.UpdateValidIngredient(ctx, exampleValidIngredient)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		expectedQuery := "UPDATE valid_ingredients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := p.buildArchiveValidIngredientQuery(exampleValidIngredient.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveValidIngredient(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_ingredients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveValidIngredient(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveValidIngredient(ctx, exampleValidIngredient.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveValidIngredient(ctx, exampleValidIngredient.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
