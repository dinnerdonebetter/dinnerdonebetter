package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fullRecipesColumns are the columns for the recipes table.
var fullRecipesColumns = []string{
	"recipes.id",
	"recipes.name",
	"recipes.slug",
	"recipes.source",
	"recipes.description",
	"recipes.inspired_by_recipe_id",
	"recipes.min_estimated_portions",
	"recipes.max_estimated_portions",
	"recipes.portion_name",
	"recipes.plural_portion_name",
	"recipes.seal_of_approval",
	"recipes.eligible_for_meals",
	"recipes.yields_component_type",
	"recipes.created_at",
	"recipes.last_updated_at",
	"recipes.archived_at",
	"recipes.created_by_user",
	"recipe_steps.id",
	"recipe_steps.index",
	"valid_preparations.id",
	"valid_preparations.name",
	"valid_preparations.description",
	"valid_preparations.icon_path",
	"valid_preparations.yields_nothing",
	"valid_preparations.restrict_to_ingredients",
	"valid_preparations.minimum_ingredient_count",
	"valid_preparations.maximum_ingredient_count",
	"valid_preparations.minimum_instrument_count",
	"valid_preparations.maximum_instrument_count",
	"valid_preparations.temperature_required",
	"valid_preparations.time_estimate_required",
	"valid_preparations.condition_expression_required",
	"valid_preparations.consumes_vessel",
	"valid_preparations.only_for_vessels",
	"valid_preparations.minimum_vessel_count",
	"valid_preparations.maximum_vessel_count",
	"valid_preparations.slug",
	"valid_preparations.past_tense",
	"valid_preparations.created_at",
	"valid_preparations.last_updated_at",
	"valid_preparations.archived_at",
	"recipe_steps.minimum_estimated_time_in_seconds",
	"recipe_steps.maximum_estimated_time_in_seconds",
	"recipe_steps.minimum_temperature_in_celsius",
	"recipe_steps.maximum_temperature_in_celsius",
	"recipe_steps.notes",
	"recipe_steps.explicit_instructions",
	"recipe_steps.condition_expression",
	"recipe_steps.optional",
	"recipe_steps.start_timer_automatically",
	"recipe_steps.created_at",
	"recipe_steps.last_updated_at",
	"recipe_steps.archived_at",
	"recipe_steps.belongs_to_recipe",
}

func buildMockFullRowsFromRecipe(recipe *types.Recipe) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(fullRecipesColumns)

	for _, step := range recipe.Steps {
		exampleRows.AddRow(
			&recipe.ID,
			&recipe.Name,
			&recipe.Slug,
			&recipe.Source,
			&recipe.Description,
			&recipe.InspiredByRecipeID,
			&recipe.MinimumEstimatedPortions,
			&recipe.MaximumEstimatedPortions,
			&recipe.PortionName,
			&recipe.PluralPortionName,
			&recipe.SealOfApproval,
			&recipe.EligibleForMeals,
			&recipe.YieldsComponentType,
			&recipe.CreatedAt,
			&recipe.LastUpdatedAt,
			&recipe.ArchivedAt,
			&recipe.CreatedByUser,
			&step.ID,
			&step.Index,
			&step.Preparation.ID,
			&step.Preparation.Name,
			&step.Preparation.Description,
			&step.Preparation.IconPath,
			&step.Preparation.YieldsNothing,
			&step.Preparation.RestrictToIngredients,
			&step.Preparation.MinimumIngredientCount,
			&step.Preparation.MaximumIngredientCount,
			&step.Preparation.MinimumInstrumentCount,
			&step.Preparation.MaximumInstrumentCount,
			&step.Preparation.TemperatureRequired,
			&step.Preparation.TimeEstimateRequired,
			&step.Preparation.ConditionExpressionRequired,
			&step.Preparation.ConsumesVessel,
			&step.Preparation.OnlyForVessels,
			&step.Preparation.MinimumVesselCount,
			&step.Preparation.MaximumVesselCount,
			&step.Preparation.Slug,
			&step.Preparation.PastTense,
			&step.Preparation.CreatedAt,
			&step.Preparation.LastUpdatedAt,
			&step.Preparation.ArchivedAt,
			&step.MinimumEstimatedTimeInSeconds,
			&step.MaximumEstimatedTimeInSeconds,
			&step.MinimumTemperatureInCelsius,
			&step.MaximumTemperatureInCelsius,
			&step.Notes,
			&step.ExplicitInstructions,
			&step.ConditionExpression,
			&step.Optional,
			&step.StartTimerAutomatically,
			&step.CreatedAt,
			&step.LastUpdatedAt,
			&step.ArchivedAt,
			&step.BelongsToRecipe,
		)
	}

	return exampleRows
}

func TestQuerier_RecipeExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipe(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipe(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipe(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipe(ctx, exampleRecipe.ID, ""))
	})
}

func Test_findCreatedRecipeStepProductsForIngredients(T *testing.T) {
	T.Parallel()

	T.Run("sopa de frijol", func(t *testing.T) {
		t.Parallel()

		soak := fakes.BuildFakeValidPreparation()
		water := fakes.BuildFakeValidIngredient()
		pintoBeans := fakes.BuildFakeValidIngredient()
		garlicPaste := fakes.BuildFakeValidIngredient()
		productName := "soaked pinto beans"

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "first step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
					},
					Index: 0,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "final output",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "second step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(0)),
						},
						{
							IngredientID:      &garlicPaste.ID,
							Name:              "garlic paste",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 1,
				},
			},
		}

		findCreatedRecipeStepProductsForIngredients(exampleRecipeInput)

		require.NotNil(t, exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
	})

	T.Run("slightly more complicated recipe", func(t *testing.T) {
		t.Parallel()

		soak := fakes.BuildFakeValidPreparation()
		water := fakes.BuildFakeValidIngredient()
		pintoBeans := fakes.BuildFakeValidIngredient()
		garlicPaste := fakes.BuildFakeValidIngredient()
		productName := "soaked pinto beans"

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "first step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   5,
						},
					},
					Index: 0,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "pressure cooked beans",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "second step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(0)),
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
						},
						{
							IngredientID:      &garlicPaste.ID,
							Name:              "garlic paste",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 1,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "third step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   5,
						},
					},
					Index: 2,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "final output",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "fourth step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(2)),
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
						},
						{
							Name:              "pressure cooked beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 3,
				},
			},
		}

		findCreatedRecipeStepProductsForIngredients(exampleRecipeInput)

		require.NotNil(t, exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		require.NotNil(t, exampleRecipeInput.Steps[3].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[2].Products[0].ID, *exampleRecipeInput.Steps[3].Ingredients[0].RecipeStepProductID)
	})
}

func Test_findCreatedRecipeStepProductsForInstruments(T *testing.T) {
	T.Parallel()

	T.Run("example", func(t *testing.T) {
		t.Parallel()

		bake := fakes.BuildFakeValidPreparation()
		line := fakes.BuildFakeValidPreparation()
		bakingSheet := fakes.BuildFakeValidInstrument()
		aluminumFoil := fakes.BuildFakeValidIngredient()
		asparagus := fakes.BuildFakeValidIngredient()
		grams := fakes.BuildFakeValidMeasurementUnit()
		sheet := fakes.BuildFakeValidMeasurementUnit()

		productName := "lined baking sheet"

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "example",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:   fakes.BuildFakeID(),
							Name: productName,
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrumentDatabaseCreationInput{
						{
							InstrumentID:        &bakingSheet.ID,
							RecipeStepProductID: nil,
							Name:                "baking sheet",
						},
					},
					Notes:         "first step",
					PreparationID: line.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							RecipeStepProductID: nil,
							IngredientID:        &aluminumFoil.ID,
							Name:                "aluminum foil",
							MeasurementUnitID:   sheet.ID,
							MinimumQuantity:     1,
						},
					},
					Index: 0,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:   fakes.BuildFakeID(),
							Name: "roasted asparagus",
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrumentDatabaseCreationInput{
						{
							InstrumentID:                    &bakingSheet.ID,
							RecipeStepProductID:             nil,
							Name:                            productName,
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(0)),
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
						},
					},
					Notes:         "second step",
					PreparationID: bake.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							RecipeStepProductID: nil,
							IngredientID:        &asparagus.ID,
							Name:                "asparagus",
							MeasurementUnitID:   grams.ID,
							MinimumQuantity:     1000,
						},
					},
					Index: 1,
				},
			},
		}

		findCreatedRecipeStepProductsForInstruments(exampleRecipeInput)

		require.NotNil(t, exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
	})
}

func TestQuerier_MarkRecipeAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkRecipeAsIndexed(ctx, ""))
	})
}
