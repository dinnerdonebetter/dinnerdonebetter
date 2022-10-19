package workers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/features/grocerylistpreparation"
	"github.com/prixfixeco/api_server/internal/features/recipeanalysis"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestProvideMealPlanGroceryListInitializer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideMealPlanGroceryListInitializer(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&recipeanalysis.MockRecipeAnalyzer{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
			&grocerylistpreparation.MockGroceryListCreator{},
		)
		assert.NotNil(t, actual)
	})
}

func TestMealPlanGroceryListInitializer_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := ProvideMealPlanGroceryListInitializer(
			zerolog.NewZerologLogger(),
			&database.MockDatabase{},
			&recipeanalysis.MockRecipeAnalyzer{},
			&mockpublishers.Publisher{},
			&customerdata.MockCollector{},
			tracing.NewNoopTracerProvider(),
			&grocerylistpreparation.MockGroceryListCreator{},
		)
		assert.NotNil(t, w)

		onion := fakes.BuildFakeValidIngredient()
		carrot := fakes.BuildFakeValidIngredient()
		celery := fakes.BuildFakeValidIngredient()
		salt := fakes.BuildFakeValidIngredient()

		grams := fakes.BuildFakeValidMeasurementUnit()

		expectedMealPlans := []*types.MealPlan{
			{
				ID: fakes.BuildFakeID(),
				Events: []*types.MealPlanEvent{
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Recipes: []*types.Recipe{
										{
											Steps: []*types.RecipeStep{
												{
													Ingredients: []*types.RecipeStepIngredient{
														{
															Ingredient:      onion,
															MinimumQuantity: 100,
															MaximumQuantity: 100,
															MeasurementUnit: *grams,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Recipes: []*types.Recipe{
										{
											Steps: []*types.RecipeStep{
												{
													Ingredients: []*types.RecipeStepIngredient{
														{
															Ingredient:      carrot,
															MinimumQuantity: 100,
															MaximumQuantity: 100,
															MeasurementUnit: *grams,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Recipes: []*types.Recipe{
										{
											Steps: []*types.RecipeStep{
												{
													Ingredients: []*types.RecipeStepIngredient{
														{
															Ingredient:      celery,
															MinimumQuantity: 100,
															MaximumQuantity: 100,
															MeasurementUnit: *grams,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Recipes: []*types.Recipe{
										{
											Steps: []*types.RecipeStep{
												{
													Ingredients: []*types.RecipeStepIngredient{
														{
															Ingredient:      salt,
															MinimumQuantity: 100,
															MaximumQuantity: 100,
															MeasurementUnit: *grams,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Recipes: []*types.Recipe{
										{
											Steps: []*types.RecipeStep{
												{
													Ingredients: []*types.RecipeStepIngredient{
														{
															Ingredient:      onion,
															MinimumQuantity: 100,
															MaximumQuantity: 100,
															MeasurementUnit: *grams,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		ctx := context.Background()
		mdm := database.NewMockDatabase()

		mdm.MealPlanDataManager.On("GetFinalizedMealPlansWithUninitializedGroceryLists", testutils.ContextMatcher).Return(expectedMealPlans, nil)

		firstMealPlanExpectedGroceryListItemInputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      onion.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  200,
				MaximumQuantityNeeded:  200,
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      carrot.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  100,
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      celery.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  100,
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      salt.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  100,
			},
		}

		expectedInputSets := map[string][]*types.MealPlanGroceryListItemDatabaseCreationInput{
			expectedMealPlans[0].ID: firstMealPlanExpectedGroceryListItemInputs,
		}

		mglm := &grocerylistpreparation.MockGroceryListCreator{}
		mglm.On(
			"GenerateGroceryListInputs",
			testutils.ContextMatcher,
			expectedMealPlans[0],
		).Return(firstMealPlanExpectedGroceryListItemInputs, nil)
		w.groceryListCreator = mglm

		for mealPlanID, inputs := range expectedInputSets {
			mdm.MealPlanGroceryListItemDataManager.On("CreateMealPlanGroceryListItemsForMealPlan", testutils.ContextMatcher, mealPlanID, inputs).Return(nil)
		}
		w.dataManager = mdm

		assert.NoError(t, w.HandleMessage(ctx, []byte("{}")))
		mock.AssertExpectationsForObjects(t, mdm)
	})
}
