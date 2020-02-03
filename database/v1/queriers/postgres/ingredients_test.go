package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowFromIngredient(x *models.Ingredient) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(ingredientsTableColumns).AddRow(
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
		x.ConsideredStaple,
		x.Icon,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

	return exampleRows
}

func buildErroneousMockRowFromIngredient(x *models.Ingredient) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(ingredientsTableColumns).AddRow(
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
		x.ConsideredStaple,
		x.Icon,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_buildGetIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleIngredientID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE belongs_to = $1 AND id = $2"
		actualQuery, args := p.buildGetIngredientQuery(exampleIngredientID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleIngredientID, args[1].(uint64))
	})
}

func TestPostgres_GetIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE belongs_to = $1 AND id = $2"
		expected := &models.Ingredient{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromIngredient(expected))

		actual, err := p.GetIngredient(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE belongs_to = $1 AND id = $2"
		expected := &models.Ingredient{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIngredient(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIngredientCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		actualQuery, args := p.buildGetIngredientCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetIngredientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetIngredientCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllIngredientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM ingredients WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllIngredientsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM ingredients WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllIngredientsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		actualQuery, args := p.buildGetIngredientsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetIngredients(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM ingredients WHERE archived_on IS NULL"
		expectedIngredient := &models.Ingredient{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.IngredientList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			Ingredients: []models.Ingredient{
				*expectedIngredient,
			},
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromIngredient(expectedIngredient))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning ingredient", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.Ingredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromIngredient(expected))

		actual, err := p.GetIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.Ingredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM ingredients WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromIngredient(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetIngredients(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_GetAllIngredientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedIngredient := &models.Ingredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromIngredient(expectedIngredient))

		expected := []models.Ingredient{*expectedIngredient}
		actual, err := p.GetAllIngredientsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleIngredient := &models.Ingredient{
			ID: 321,
		}
		expectedListQuery := "SELECT id, name, variant, description, warning, contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten, animal_flesh, animal_derived, considered_staple, icon, created_on, updated_on, archived_on, belongs_to FROM ingredients WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromIngredient(exampleIngredient))

		actual, err := p.GetAllIngredientsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Ingredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 19
		expectedQuery := "INSERT INTO ingredients (name,variant,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,animal_derived,considered_staple,icon,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id, created_on"
		actualQuery, args := p.buildCreateIngredientQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.Warning, args[3].(string))
		assert.Equal(t, expected.ContainsEgg, args[4].(bool))
		assert.Equal(t, expected.ContainsDairy, args[5].(bool))
		assert.Equal(t, expected.ContainsPeanut, args[6].(bool))
		assert.Equal(t, expected.ContainsTreeNut, args[7].(bool))
		assert.Equal(t, expected.ContainsSoy, args[8].(bool))
		assert.Equal(t, expected.ContainsWheat, args[9].(bool))
		assert.Equal(t, expected.ContainsShellfish, args[10].(bool))
		assert.Equal(t, expected.ContainsSesame, args[11].(bool))
		assert.Equal(t, expected.ContainsFish, args[12].(bool))
		assert.Equal(t, expected.ContainsGluten, args[13].(bool))
		assert.Equal(t, expected.AnimalFlesh, args[14].(bool))
		assert.Equal(t, expected.AnimalDerived, args[15].(bool))
		assert.Equal(t, expected.ConsideredStaple, args[16].(bool))
		assert.Equal(t, expected.Icon, args[17].(string))
		assert.Equal(t, expected.BelongsTo, args[18].(uint64))
	})
}

func TestPostgres_CreateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Ingredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.IngredientCreationInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
			BelongsTo:         expected.BelongsTo,
		}
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO ingredients (name,variant,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,animal_derived,considered_staple,icon,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Warning,
				expected.ContainsEgg,
				expected.ContainsDairy,
				expected.ContainsPeanut,
				expected.ContainsTreeNut,
				expected.ContainsSoy,
				expected.ContainsWheat,
				expected.ContainsShellfish,
				expected.ContainsSesame,
				expected.ContainsFish,
				expected.ContainsGluten,
				expected.AnimalFlesh,
				expected.AnimalDerived,
				expected.ConsideredStaple,
				expected.Icon,
				expected.BelongsTo,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateIngredient(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Ingredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.IngredientCreationInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
			BelongsTo:         expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO ingredients (name,variant,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,animal_derived,considered_staple,icon,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Warning,
				expected.ContainsEgg,
				expected.ContainsDairy,
				expected.ContainsPeanut,
				expected.ContainsTreeNut,
				expected.ContainsSoy,
				expected.ContainsWheat,
				expected.ContainsShellfish,
				expected.ContainsSesame,
				expected.ContainsFish,
				expected.ContainsGluten,
				expected.AnimalFlesh,
				expected.AnimalDerived,
				expected.ConsideredStaple,
				expected.Icon,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateIngredient(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Ingredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 20
		expectedQuery := "UPDATE ingredients SET name = $1, variant = $2, description = $3, warning = $4, contains_egg = $5, contains_dairy = $6, contains_peanut = $7, contains_tree_nut = $8, contains_soy = $9, contains_wheat = $10, contains_shellfish = $11, contains_sesame = $12, contains_fish = $13, contains_gluten = $14, animal_flesh = $15, animal_derived = $16, considered_staple = $17, icon = $18, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $19 AND id = $20 RETURNING updated_on"
		actualQuery, args := p.buildUpdateIngredientQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Variant, args[1].(string))
		assert.Equal(t, expected.Description, args[2].(string))
		assert.Equal(t, expected.Warning, args[3].(string))
		assert.Equal(t, expected.ContainsEgg, args[4].(bool))
		assert.Equal(t, expected.ContainsDairy, args[5].(bool))
		assert.Equal(t, expected.ContainsPeanut, args[6].(bool))
		assert.Equal(t, expected.ContainsTreeNut, args[7].(bool))
		assert.Equal(t, expected.ContainsSoy, args[8].(bool))
		assert.Equal(t, expected.ContainsWheat, args[9].(bool))
		assert.Equal(t, expected.ContainsShellfish, args[10].(bool))
		assert.Equal(t, expected.ContainsSesame, args[11].(bool))
		assert.Equal(t, expected.ContainsFish, args[12].(bool))
		assert.Equal(t, expected.ContainsGluten, args[13].(bool))
		assert.Equal(t, expected.AnimalFlesh, args[14].(bool))
		assert.Equal(t, expected.AnimalDerived, args[15].(bool))
		assert.Equal(t, expected.ConsideredStaple, args[16].(bool))
		assert.Equal(t, expected.Icon, args[17].(string))
		assert.Equal(t, expected.BelongsTo, args[18].(uint64))
		assert.Equal(t, expected.ID, args[19].(uint64))
	})
}

func TestPostgres_UpdateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Ingredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE ingredients SET name = $1, variant = $2, description = $3, warning = $4, contains_egg = $5, contains_dairy = $6, contains_peanut = $7, contains_tree_nut = $8, contains_soy = $9, contains_wheat = $10, contains_shellfish = $11, contains_sesame = $12, contains_fish = $13, contains_gluten = $14, animal_flesh = $15, animal_derived = $16, considered_staple = $17, icon = $18, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $19 AND id = $20 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Warning,
				expected.ContainsEgg,
				expected.ContainsDairy,
				expected.ContainsPeanut,
				expected.ContainsTreeNut,
				expected.ContainsSoy,
				expected.ContainsWheat,
				expected.ContainsShellfish,
				expected.ContainsSesame,
				expected.ContainsFish,
				expected.ContainsGluten,
				expected.AnimalFlesh,
				expected.AnimalDerived,
				expected.ConsideredStaple,
				expected.Icon,
				expected.BelongsTo,
				expected.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateIngredient(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Ingredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE ingredients SET name = $1, variant = $2, description = $3, warning = $4, contains_egg = $5, contains_dairy = $6, contains_peanut = $7, contains_tree_nut = $8, contains_soy = $9, contains_wheat = $10, contains_shellfish = $11, contains_sesame = $12, contains_fish = $13, contains_gluten = $14, animal_flesh = $15, animal_derived = $16, considered_staple = $17, icon = $18, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $19 AND id = $20 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Variant,
				expected.Description,
				expected.Warning,
				expected.ContainsEgg,
				expected.ContainsDairy,
				expected.ContainsPeanut,
				expected.ContainsTreeNut,
				expected.ContainsSoy,
				expected.ContainsWheat,
				expected.ContainsShellfish,
				expected.ContainsSesame,
				expected.ContainsFish,
				expected.ContainsGluten,
				expected.AnimalFlesh,
				expected.AnimalDerived,
				expected.ConsideredStaple,
				expected.Icon,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateIngredient(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Ingredient{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE ingredients SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"
		actualQuery, args := p.buildArchiveIngredientQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestPostgres_ArchiveIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Ingredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE ingredients SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveIngredient(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.Ingredient{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE ingredients SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveIngredient(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
