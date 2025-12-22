package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepVesselSliceEquality(t *testing.T, stepIndex int, expected, actual []*mealplanning.RecipeStepVessel) {
	t.Helper()
	require.Equal(t, len(expected), len(actual), "expected recipe step %d vessels length", stepIndex)
	for i := range expected {
		checkRecipeStepVesselEquality(t, stepIndex, i, expected[i], actual[i])
	}
}

func checkRecipeStepVesselEquality(t *testing.T, stepIndex, vesselIndex int, expected, actual *mealplanning.RecipeStepVessel) {
	t.Helper()
	assert.NotEmpty(t, actual.ID, "expected step %d vessel %d to have ID", stepIndex, vesselIndex)
	assert.False(t, actual.CreatedAt.IsZero(), "expected step %d vessel %d to have CreatedAt", stepIndex, vesselIndex)
	assert.NotEmpty(t, actual.BelongsToRecipeStep, "expected step %d vessel %d to have BelongsToRecipeStep", stepIndex, vesselIndex)
	assert.Equal(t, expected.Name, actual.Name, "expected step %d vessel %d Name", stepIndex, vesselIndex)
	assert.Equal(t, expected.Notes, actual.Notes, "expected step %d vessel %d Notes", stepIndex, vesselIndex)
	assert.Equal(t, expected.Quantity, actual.Quantity, "expected step %d vessel %d Quantity", stepIndex, vesselIndex)
	assert.Equal(t, expected.VesselPreposition, actual.VesselPreposition, "expected step %d vessel %d VesselPreposition", stepIndex, vesselIndex)
	assert.Equal(t, expected.UnavailableAfterStep, actual.UnavailableAfterStep, "expected step %d vessel %d UnavailableAfterStep", stepIndex, vesselIndex)
	if expected.Vessel != nil {
		require.NotNil(t, actual.Vessel, "expected step %d vessel %d Vessel non-nil", stepIndex, vesselIndex)
		assert.NotEmpty(t, actual.Vessel.ID, "expected step %d vessel %d Vessel.ID", stepIndex, vesselIndex)
		assert.Equal(t, expected.Vessel.ID, actual.Vessel.ID, "expected step %d vessel %d Vessel.ID", stepIndex, vesselIndex)
	}
}

func TestRecipeStepVessels_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should be able to create a recipe step vessel", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		var createdRecipeStepID string
		for _, step := range createdRecipe.Steps {
			createdRecipeStepID = step.ID
			break
		}

		createdValidVessel := createValidVesselForTest(t)

		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
		exampleRecipeStepVessel.BelongsToRecipeStep = createdRecipeStepID
		exampleRecipeStepVessel.Vessel = &mealplanning.ValidVessel{ID: createdValidVessel.ID}
		exampleRecipeStepVesselInput := mpconverters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(exampleRecipeStepVessel)

		createdRecipeStepVesselRes, err := userClient.CreateRecipeStepVessel(ctx, &mealplanninggrpc.CreateRecipeStepVesselRequest{
			RecipeId:     createdRecipe.ID,
			RecipeStepId: createdRecipeStepID,
			Input:        converters.ConvertRecipeStepVesselCreationRequestInputToGRPCRecipeStepVesselCreationRequestInput(exampleRecipeStepVesselInput),
		})
		require.NoError(t, err)

		createdRecipeStepVessel := converters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(createdRecipeStepVesselRes.Created)
		checkRecipeStepVesselEquality(t, -1, -1, exampleRecipeStepVessel, createdRecipeStepVessel)

		retrievedRecipeStepVesselRes, err := userClient.GetRecipeStepVessel(ctx, &mealplanninggrpc.GetRecipeStepVesselRequest{
			RecipeId:           createdRecipe.ID,
			RecipeStepId:       createdRecipeStepID,
			RecipeStepVesselId: createdRecipeStepVessel.ID,
		})
		require.NoError(t, err)

		createdRecipeStepVessel = converters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(retrievedRecipeStepVesselRes.Result)
		checkRecipeStepVesselEquality(t, -1, -1, exampleRecipeStepVessel, createdRecipeStepVessel)

		newVessel := createValidVesselForTest(t)

		newRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
		newRecipeStepVessel.BelongsToRecipeStep = createdRecipeStepID
		newRecipeStepVessel.Vessel = newVessel
		updateInput := mpconverters.ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(newRecipeStepVessel)
		createdRecipeStepVessel.Update(updateInput)
		exampleRecipeStepVessel.Update(updateInput)

		_, err = adminClient.UpdateRecipeStepVessel(ctx, &mealplanninggrpc.UpdateRecipeStepVesselRequest{
			RecipeId:           createdRecipe.ID,
			RecipeStepId:       createdRecipeStepID,
			RecipeStepVesselId: createdRecipeStepVessel.ID,
			Input:              converters.ConvertRecipeStepVesselUpdateRequestInputToGRPCRecipeStepVesselUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)

		retrievedRecipeStepVesselRes, err = userClient.GetRecipeStepVessel(ctx, &mealplanninggrpc.GetRecipeStepVesselRequest{
			RecipeId:           createdRecipe.ID,
			RecipeStepId:       createdRecipeStepID,
			RecipeStepVesselId: createdRecipeStepVessel.ID,
		})
		require.NoError(t, err)

		createdRecipeStepVessel = converters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(retrievedRecipeStepVesselRes.Result)
		checkRecipeStepVesselEquality(t, -1, -1, exampleRecipeStepVessel, createdRecipeStepVessel)

		_, err = userClient.ArchiveRecipeStepVessel(ctx, &mealplanninggrpc.ArchiveRecipeStepVesselRequest{
			RecipeId:           createdRecipe.ID,
			RecipeStepId:       createdRecipeStepID,
			RecipeStepVesselId: createdRecipeStepVessel.ID,
		})
		assert.NoError(t, err)

		_, err = userClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{
			RecipeId:     createdRecipe.ID,
			RecipeStepId: createdRecipeStepID,
		})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
		assert.NoError(t, err)
	})
}

func TestRecipeStepVessels_AsRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("should be able to use a recipe step vessel that was the product of a prior recipe step", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		line := createValidPreparationForTest(t)
		roast := createValidPreparationForTest(t)
		bakingSheet := createValidVesselForTest(t)
		sheets := createValidMeasurementUnitForTest(t)
		head := createValidMeasurementUnitForTest(t)
		unit := createValidMeasurementUnitForTest(t)
		aluminumFoil := createValidIngredientForTest(t)
		garlic := createValidIngredientForTest(t)

		// Create bridge table entries
		// Step 0: aluminumFoil ingredient with line preparation, bakingSheet vessel with line preparation
		vipAluminumFoilLine := createValidIngredientPreparationWithEntitiesForTest(t, aluminumFoil, line)
		vimuAluminumFoilSheets := createValidIngredientMeasurementUnitWithEntitiesForTest(t, aluminumFoil, sheets)
		vpvBakingSheetLine := createValidPreparationVesselWithEntitiesForTest(t, line, bakingSheet)
		// Step 1: garlic ingredient with roast preparation (vessel is a recipe step product, no bridge ID needed)
		vipGarlicRoast := createValidIngredientPreparationWithEntitiesForTest(t, garlic, roast)
		vimuGarlicHead := createValidIngredientMeasurementUnitWithEntitiesForTest(t, garlic, head)

		linedBakingSheetName := "lined baking sheet"

		expected := &mealplanning.Recipe{
			Name:                t.Name(),
			Slug:                "whatever-who-cares-yadda-yadda-vessels",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: nil,
				Min: 1,
			},
			Steps: []*mealplanning.RecipeStep{
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            linedBakingSheetName,
							Type:            mealplanning.RecipeStepProductVesselType,
							MeasurementUnit: unit,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1)),
							},
						},
					},
					Notes:       "first step",
					Preparation: *line,
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient:          aluminumFoil,
							Name:                "aluminum foil",
							MeasurementUnit:     *sheets,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: 3,
							},
						},
					},
					Vessels: []*mealplanning.RecipeStepVessel{
						{
							Vessel: bakingSheet,
							Name:   "baking sheet",
							Quantity: types.Uint16RangeWithOptionalMax{
								Max: nil,
								Min: 1,
							},
						},
					},
					Index: 0,
				},
				{
					Preparation: *roast,
					Vessels: []*mealplanning.RecipeStepVessel{
						{
							Name:   linedBakingSheetName,
							Vessel: nil,
						},
					},
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "roasted garlic",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: head,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1)),
							},
						},
					},
					Notes: "second step",
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient:      garlic,
							Name:            "garlic",
							MeasurementUnit: *head,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: 1,
							},
						},
					},
					Index: 1,
				},
			},
		}

		exampleRecipeInput := mpconverters.ConvertRecipeToRecipeCreationRequestInput(expected)
		exampleRecipeInput.Steps[1].Vessels[0].ProductOfRecipeStepIndex = pointer.To(uint64(0))
		exampleRecipeInput.Steps[1].Vessels[0].ProductOfRecipeStepProductIndex = pointer.To(uint64(0))

		// Set bridge table IDs
		// Step 0: aluminumFoil ingredient and bakingSheet vessel with line preparation
		exampleRecipeInput.Steps[0].Ingredients[0].ValidIngredientPreparationID = &vipAluminumFoilLine.ID
		exampleRecipeInput.Steps[0].Ingredients[0].ValidIngredientMeasurementUnitID = &vimuAluminumFoilSheets.ID
		exampleRecipeInput.Steps[0].Vessels[0].ValidPreparationVesselID = &vpvBakingSheetLine.ID
		// Step 1: garlic ingredient with roast preparation (vessel is a recipe step product, no bridge ID needed)
		exampleRecipeInput.Steps[1].Ingredients[0].ValidIngredientPreparationID = &vipGarlicRoast.ID
		exampleRecipeInput.Steps[1].Ingredients[0].ValidIngredientMeasurementUnitID = &vimuGarlicHead.ID

		createdRes, err := adminClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{Input: converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(exampleRecipeInput)})
		require.NoError(t, err)

		created := converters.ConvertGRPCRecipeToRecipe(createdRes.Created)
		expected.Status = created.Status
		checkRecipeEquality(t, expected, created)

		retrievedRes, err := userClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeId: created.ID})
		require.NotNil(t, retrievedRes)
		require.NoError(t, err)
		checkRecipeEquality(t, expected, created)

		recipeStepProductIndex := -1
		for i, vessel := range created.Steps[1].Vessels {
			if vessel.RecipeStepProductID != nil {
				recipeStepProductIndex = i
			}
		}

		require.NotEqual(t, -1, recipeStepProductIndex)
		require.NotNil(t, created.Steps[1].Vessels[recipeStepProductIndex].RecipeStepProductID)
		assert.Equal(t, created.Steps[0].Products[0].ID, *created.Steps[1].Vessels[recipeStepProductIndex].RecipeStepProductID)
	})
}

func TestRecipeStepVessels_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be able to list recipe step vessels", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		var createdRecipeStepID string
		for _, step := range createdRecipe.Steps {
			createdRecipeStepID = step.ID
			break
		}

		createdValidVessel := createValidVesselForTest(t)

		var expected []*mealplanning.RecipeStepVessel
		for i := 0; i < 5; i++ {
			exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
			exampleRecipeStepVessel.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepVessel.Vessel = &mealplanning.ValidVessel{ID: createdValidVessel.ID}
			exampleRecipeStepVesselInput := mpconverters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(exampleRecipeStepVessel)
			createdRecipeStepVesselRes, err := adminClient.CreateRecipeStepVessel(ctx, &mealplanninggrpc.CreateRecipeStepVesselRequest{
				RecipeId:     createdRecipe.ID,
				RecipeStepId: createdRecipeStepID,
				Input:        converters.ConvertRecipeStepVesselCreationRequestInputToGRPCRecipeStepVesselCreationRequestInput(exampleRecipeStepVesselInput),
			})
			require.NoError(t, err)

			createdRecipeStepVessel := converters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(createdRecipeStepVesselRes.Created)
			checkRecipeStepVesselEquality(t, i, i, exampleRecipeStepVessel, createdRecipeStepVessel)

			retrievedRecipeStepVesselRes, err := userClient.GetRecipeStepVessel(ctx, &mealplanninggrpc.GetRecipeStepVesselRequest{
				RecipeId:           createdRecipe.ID,
				RecipeStepId:       createdRecipeStepID,
				RecipeStepVesselId: createdRecipeStepVessel.ID,
			})
			require.NotNil(t, retrievedRecipeStepVesselRes)
			require.NoError(t, err)

			createdRecipeStepVessel = converters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(retrievedRecipeStepVesselRes.Result)
			require.Equal(t, createdRecipeStepID, createdRecipeStepVessel.BelongsToRecipeStep)

			expected = append(expected, createdRecipeStepVessel)
		}

		// assert recipe step vessel list equality
		actual, err := userClient.GetRecipeStepVessels(ctx, &mealplanninggrpc.GetRecipeStepVesselsRequest{
			RecipeId:     createdRecipe.ID,
			RecipeStepId: createdRecipeStepID,
		})
		require.NotNil(t, actual)
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdRecipeStepVessel := range expected {
			_, err = userClient.ArchiveRecipeStepVessel(ctx, &mealplanninggrpc.ArchiveRecipeStepVesselRequest{
				RecipeId:           createdRecipe.ID,
				RecipeStepId:       createdRecipeStepID,
				RecipeStepVesselId: createdRecipeStepVessel.ID,
			})
			assert.NoError(t, err)
		}

		_, err = userClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{RecipeId: createdRecipe.ID, RecipeStepId: createdRecipeStepID})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
		assert.NoError(t, err)
	})
}
