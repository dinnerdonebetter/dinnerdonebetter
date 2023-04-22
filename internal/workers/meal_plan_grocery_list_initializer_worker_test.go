package workers

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/analytics"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/features/grocerylistpreparation"
	"github.com/prixfixeco/backend/internal/features/recipeanalysis"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	testutils "github.com/prixfixeco/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProvideMealPlanGroceryListInitializer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideMealPlanGroceryListInitializer(
			zerolog.NewZerologLogger(logging.DebugLevel),
			&database.MockDatabase{},
			&recipeanalysis.MockRecipeAnalyzer{},
			&mockpublishers.Publisher{},
			&analytics.MockEventReporter{},
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
			zerolog.NewZerologLogger(logging.DebugLevel),
			&database.MockDatabase{},
			&recipeanalysis.MockRecipeAnalyzer{},
			&mockpublishers.Publisher{},
			&analytics.MockEventReporter{},
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
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      onion,
																MinimumQuantity: 100,
																MaximumQuantity: pointers.Pointer(float32(100)),
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
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      carrot,
																MinimumQuantity: 100,
																MaximumQuantity: pointers.Pointer(float32(100)),
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
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      celery,
																MinimumQuantity: 100,
																MaximumQuantity: pointers.Pointer(float32(100)),
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
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      salt,
																MinimumQuantity: 100,
																MaximumQuantity: pointers.Pointer(float32(100)),
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
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      onion,
																MinimumQuantity: 100,
																MaximumQuantity: pointers.Pointer(float32(100)),
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
				MaximumQuantityNeeded:  pointers.Pointer(float32(200)),
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      carrot.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  pointers.Pointer(float32(100)),
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      celery.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  pointers.Pointer(float32(100)),
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      salt.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  pointers.Pointer(float32(100)),
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
