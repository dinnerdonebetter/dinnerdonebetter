package mealplangrocerylistinitializer

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	grocerylistpreparation2 "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	mealplanningmock "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/mocks"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildNewMealPlanGroceryListInitializerForTest(t *testing.T) *Worker {
	t.Helper()

	ctx := t.Context()
	cfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProvider{}
	pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	x, err := NewMealPlanGroceryListInitializer(
		ctx,
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		metrics.NewNoopMetricsProvider(),
		pp,
		grocerylistpreparation2.NewGroceryListCreator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
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

		expectedMealPlans := []*mealplanning.MealPlan{
			{
				ID: fakes.BuildFakeID(),
				Events: []*mealplanning.MealPlanEvent{
					{
						Options: []*mealplanning.MealPlanOption{
							{
								Chosen: true,
								Meal: mealplanning.Meal{
									Components: []*mealplanning.MealComponent{
										{
											Recipe: mealplanning.Recipe{
												Steps: []*mealplanning.RecipeStep{
													{
														Ingredients: []*mealplanning.RecipeStepIngredient{
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
						Options: []*mealplanning.MealPlanOption{
							{
								Chosen: true,
								Meal: mealplanning.Meal{
									Components: []*mealplanning.MealComponent{
										{
											Recipe: mealplanning.Recipe{
												Steps: []*mealplanning.RecipeStep{
													{
														Ingredients: []*mealplanning.RecipeStepIngredient{
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
						Options: []*mealplanning.MealPlanOption{
							{
								Chosen: true,
								Meal: mealplanning.Meal{
									Components: []*mealplanning.MealComponent{
										{
											Recipe: mealplanning.Recipe{
												Steps: []*mealplanning.RecipeStep{
													{
														Ingredients: []*mealplanning.RecipeStepIngredient{
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
						Options: []*mealplanning.MealPlanOption{
							{
								Chosen: true,
								Meal: mealplanning.Meal{
									Components: []*mealplanning.MealComponent{
										{
											Recipe: mealplanning.Recipe{
												Steps: []*mealplanning.RecipeStep{
													{
														Ingredients: []*mealplanning.RecipeStepIngredient{
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
						Options: []*mealplanning.MealPlanOption{
							{
								Chosen: true,
								Meal: mealplanning.Meal{
									Components: []*mealplanning.MealComponent{
										{
											Recipe: mealplanning.Recipe{
												Steps: []*mealplanning.RecipeStep{
													{
														Ingredients: []*mealplanning.RecipeStepIngredient{
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

		ctx := t.Context()
		mdm := &mealplanningmock.Repository{}

		mdm.On("GetFinalizedMealPlansWithUninitializedGroceryLists", testutils.ContextMatcher).Return(expectedMealPlans, nil)

		firstMealPlanExpectedGroceryListItemInputs := []*mealplanning.MealPlanGroceryListItemDatabaseCreationInput{
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      onion.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(200)),
					Min: 200,
				},
			},
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      carrot.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(100)),
					Min: 100,
				},
			},
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      celery.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(100)),
					Min: 100,
				},
			},
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      salt.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: pointer.To(float32(100)),
					Min: 100,
				},
			},
		}

		expectedInputSets := map[string][]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput{
			expectedMealPlans[0].ID: firstMealPlanExpectedGroceryListItemInputs,
		}

		mglm := &grocerylistpreparation2.MockGroceryListCreator{}
		mglm.On(
			"GenerateGroceryListInputs",
			testutils.ContextMatcher,
			expectedMealPlans[0],
		).Return(firstMealPlanExpectedGroceryListItemInputs, nil)
		w.groceryListCreator = mglm

		pup := &mockpublishers.Publisher{}
		for _, inputs := range expectedInputSets {
			for _, input := range inputs {
				mdm.On("CreateMealPlanGroceryListItem", testutils.ContextMatcher, input).Return(fakes.BuildFakeMealPlanGroceryListItem(), nil)
				pup.On("Publish", testutils.ContextMatcher, mock.AnythingOfType("*audit.DataChangeMessage")).Return(nil)
			}
		}

		w.postUpdatesPublisher = pup
		w.dataManager = mdm

		assert.NoError(t, w.Work(ctx))
		mock.AssertExpectationsForObjects(t, mdm)
	})
}
