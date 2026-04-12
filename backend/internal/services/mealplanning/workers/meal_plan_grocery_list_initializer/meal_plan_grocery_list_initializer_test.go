package mealplangrocerylistinitializer

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	grocerylistpreparation2 "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	mealplanningmock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/mocks"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/testutils"

	"github.com/primandproper/platform/messagequeue"
	msgconfig "github.com/primandproper/platform/messagequeue/config"
	mockpublishers "github.com/primandproper/platform/messagequeue/mock"
	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/metrics"
	"github.com/primandproper/platform/observability/tracing"
	"github.com/primandproper/platform/reflection"
	"github.com/primandproper/platform/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildNewMealPlanGroceryListInitializerForTest(t *testing.T) *Worker {
	t.Helper()

	ctx := t.Context()
	cfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProviderMock{
		ProvidePublisherFunc: func(_ context.Context, topic string) (messagequeue.Publisher, error) {
			return &mockpublishers.PublisherMock{
				PublishFunc:      func(_ context.Context, _ any) error { return nil },
				PublishAsyncFunc: func(_ context.Context, _ any) {},
				StopFunc:         func() {},
			}, nil
		},
	}

	x, err := NewMealPlanGroceryListInitializer(
		ctx,
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		metrics.NewNoopMetricsProvider(),
		pp,
		&mealplanningmock.Repository{},
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
																	Max: new(float32(100)),
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
																	Max: new(float32(100)),
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
																	Max: new(float32(100)),
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
																	Max: new(float32(100)),
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
																	Max: new(float32(100)),
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

		mdm.On(reflection.GetMethodName(mdm.GetFinalizedMealPlansWithUninitializedGroceryLists), testutils.ContextMatcher).Return(expectedMealPlans, nil)

		firstMealPlanExpectedGroceryListItemInputs := []*mealplanning.MealPlanGroceryListItemDatabaseCreationInput{
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      onion.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: new(float32(200)),
					Min: 200,
				},
			},
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      carrot.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: new(float32(100)),
					Min: 100,
				},
			},
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      celery.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: new(float32(100)),
					Min: 100,
				},
			},
			{
				Status:                 mealplanning.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      salt.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				QuantityNeeded: types.Float32RangeWithOptionalMax{
					Max: new(float32(100)),
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

		pup := &mockpublishers.PublisherMock{
			PublishFunc: func(_ context.Context, _ any) error { return nil },
		}
		for _, inputs := range expectedInputSets {
			for _, input := range inputs {
				mdm.On(reflection.GetMethodName(mdm.CreateMealPlanGroceryListItem), testutils.ContextMatcher, input).Return(fakes.BuildFakeMealPlanGroceryListItem(), nil)
			}
		}

		mdm.On(reflection.GetMethodName(mdm.MarkMealPlanAsGroceryListInitialized), testutils.ContextMatcher, expectedMealPlans[0].ID).Return(nil)

		w.postUpdatesPublisher = pup
		w.dataManager = mdm

		assert.NoError(t, w.Work(ctx))
		mock.AssertExpectationsForObjects(t, mdm)
	})
}
