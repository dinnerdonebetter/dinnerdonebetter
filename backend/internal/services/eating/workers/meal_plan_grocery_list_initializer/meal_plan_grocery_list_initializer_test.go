package mealplangrocerylistinitializer

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/grocerylistpreparation"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildNewMealPlanGroceryListInitializerForTest(t *testing.T) *Worker {
	t.Helper()

	cfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProvider{}
	pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	x, err := NewMealPlanGroceryListInitializer(
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		metrics.NewNoopMetricsProvider(),
		pp,
		grocerylistpreparation.NewGroceryListCreator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		cfg,
	)
	require.NoError(t, err)

	return x
}

func TestMealPlanGroceryListInitializer_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := buildNewMealPlanGroceryListInitializerForTest(t)
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
																Ingredient: onion,
																Quantity: types.Float32RangeWithOptionalMax{
																	Max: pointer.To(float32(100)),
																	Min: 100,
																},
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
																Ingredient: carrot,
																Quantity: types.Float32RangeWithOptionalMax{
																	Max: pointer.To(float32(100)),
																	Min: 100,
																},
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
																Ingredient: celery,
																Quantity: types.Float32RangeWithOptionalMax{
																	Max: pointer.To(float32(100)),
																	Min: 100,
																},
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
																Ingredient: salt,
																Quantity: types.Float32RangeWithOptionalMax{
																	Max: pointer.To(float32(100)),
																	Min: 100,
																},
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
																Ingredient: onion,
																Quantity: types.Float32RangeWithOptionalMax{
																	Max: pointer.To(float32(100)),
																	Min: 100,
																},
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

		mdm.MealPlanDataManagerMock.On("GetFinalizedMealPlansWithUninitializedGroceryLists", testutils.ContextMatcher).Return(expectedMealPlans, nil)

		firstMealPlanExpectedGroceryListItemInputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      onion.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(200)),
					Min: 200,
				},
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      carrot.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(100)),
					Min: 100,
				},
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      celery.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(100)),
					Min: 100,
				},
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      salt.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(100)),
					Min: 100,
				},
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

		pup := &mockpublishers.Publisher{}
		for _, inputs := range expectedInputSets {
			for _, input := range inputs {
				mdm.MealPlanGroceryListItemDataManagerMock.On("CreateMealPlanGroceryListItem", testutils.ContextMatcher, input).Return(fakes.BuildFakeMealPlanGroceryListItem(), nil)
				pup.On("Publish", testutils.ContextMatcher, mock.AnythingOfType("*types.DataChangeMessage")).Return(nil)
			}
		}

		w.postUpdatesPublisher = pup
		w.dataManager = mdm

		assert.NoError(t, w.Work(ctx))
		mock.AssertExpectationsForObjects(t, mdm)
	})
}
