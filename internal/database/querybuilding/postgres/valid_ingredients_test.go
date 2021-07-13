package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildValidIngredientExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		expectedQuery := "SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := q.BuildValidIngredientExistsQuery(ctx, exampleValidIngredient.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.external_id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.volumetric, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.id = $1"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := q.BuildGetValidIngredientQuery(ctx, exampleValidIngredient.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllValidIngredientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL"
		actualQuery := q.BuildGetAllValidIngredientsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfValidIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.external_id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.volumetric, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.id > $1 AND valid_ingredients.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfValidIngredientsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.external_id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.volumetric, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on, (SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL) as total_count, (SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.created_on > $1 AND valid_ingredients.created_on < $2 AND valid_ingredients.last_updated_on > $3 AND valid_ingredients.last_updated_on < $4) as filtered_count FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.created_on > $5 AND valid_ingredients.created_on < $6 AND valid_ingredients.last_updated_on > $7 AND valid_ingredients.last_updated_on < $8 GROUP BY valid_ingredients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.BuildGetValidIngredientsQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidIngredientsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT valid_ingredients.id, valid_ingredients.external_id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.volumetric, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM (SELECT valid_ingredients.id, valid_ingredients.external_id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.volumetric, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.id IN ($1,$2,$3)"
		expectedArgs := []interface{}{
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetValidIngredientsWithIDsQuery(ctx, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		exampleInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleValidIngredient.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO valid_ingredients (external_id,name,variant,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,animal_derived,volumetric,icon_path) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id"
		expectedArgs := []interface{}{
			exampleValidIngredient.ExternalID,
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
			exampleValidIngredient.Volumetric,
			exampleValidIngredient.IconPath,
		}
		actualQuery, actualArgs := q.BuildCreateValidIngredientQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		expectedQuery := "UPDATE valid_ingredients SET name = $1, variant = $2, description = $3, warning = $4, contains_egg = $5, contains_dairy = $6, contains_peanut = $7, contains_tree_nut = $8, contains_soy = $9, contains_wheat = $10, contains_shellfish = $11, contains_sesame = $12, contains_fish = $13, contains_gluten = $14, animal_flesh = $15, animal_derived = $16, volumetric = $17, icon_path = $18, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $19"
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
			exampleValidIngredient.Volumetric,
			exampleValidIngredient.IconPath,
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateValidIngredientQuery(ctx, exampleValidIngredient)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredientID := fakes.BuildFakeID()

		expectedQuery := "UPDATE valid_ingredients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleValidIngredientID,
		}
		actualQuery, actualArgs := q.BuildArchiveValidIngredientQuery(ctx, exampleValidIngredientID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForValidIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'valid_ingredient_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForValidIngredientQuery(ctx, exampleValidIngredient.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
