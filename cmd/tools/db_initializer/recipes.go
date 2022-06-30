package main

import (
	"context"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	mushroomRisottoID      = ksuid.New().String()
	mushroomRisottoStep1ID = ksuid.New().String()
	mushroomRisottoStep2ID = ksuid.New().String()

	grilledChickenID      = ksuid.New().String()
	grilledChickenStep1ID = ksuid.New().String()
	grilledChickenStep2ID = ksuid.New().String()

	capreseSaladID      = ksuid.New().String()
	capreseSaladStep1ID = ksuid.New().String()
	capreseSaladStep2ID = ksuid.New().String()

	spaghettiWithNeatballsID      = ksuid.New().String()
	spaghettiWithNeatballsStep1ID = ksuid.New().String()
	spaghettiWithNeatballsStep2ID = ksuid.New().String()
)

var recipeCollection = struct {
	MushroomRisotto,
	GrilledChicken,
	CapreseSalad,
	SpaghettiWithNeatballs *types.RecipeDatabaseCreationInput
}{
	MushroomRisotto: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 mushroomRisottoID,
		Name:               "mushroom risotto",
		Source:             "",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
		Steps: []*types.RecipeStepDatabaseCreationInput{
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "",
				PreparationID:        validPreparationCollection.Dice.ID,
				BelongsToRecipe:      mushroomRisottoID,
				ID:                   mushroomRisottoStep1ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.Onion.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: mushroomRisottoStep1ID,
						QuantityValue:       400,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     0,
				MinEstimatedTimeInSeconds: 0,
				MaxEstimatedTimeInSeconds: 0,
				Optional:                  false,
			},
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "helps avoid burning the onion",
				PreparationID:        validPreparationCollection.Sautee.ID,
				BelongsToRecipe:      mushroomRisottoID,
				ID:                   mushroomRisottoStep2ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.Water.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "milliliters",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: mushroomRisottoStep2ID,
						QuantityValue:       250,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     1,
				MinEstimatedTimeInSeconds: 0,
				MaxEstimatedTimeInSeconds: 0,
				Optional:                  false,
			},
		},
	},
	GrilledChicken: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 grilledChickenID,
		Name:               "grilled chicken",
		Source:             "",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
		Steps: []*types.RecipeStepDatabaseCreationInput{
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "",
				PreparationID:        validPreparationCollection.Marinate.ID,
				BelongsToRecipe:      grilledChickenID,
				ID:                   grilledChickenStep1ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.ChickenBreast.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: grilledChickenStep1ID,
						QuantityValue:       900,
						ProductOfRecipeStep: false,
					},
					{
						IngredientID:        sp(validIngredientCollection.BlackPepper.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: grilledChickenStep1ID,
						QuantityValue:       3,
						ProductOfRecipeStep: false,
					},
					{
						IngredientID:        sp(validIngredientCollection.Garlic.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: grilledChickenStep1ID,
						QuantityValue:       20,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     0,
				MinEstimatedTimeInSeconds: 0,
				MaxEstimatedTimeInSeconds: 0,
				Optional:                  false,
			},
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "",
				PreparationID:        validPreparationCollection.Grill.ID,
				BelongsToRecipe:      grilledChickenID,
				ID:                   grilledChickenStep2ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.ChickenBreast.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: grilledChickenStep2ID,
						QuantityValue:       900,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     0,
				MinEstimatedTimeInSeconds: 0,
				MaxEstimatedTimeInSeconds: 0,
				Optional:                  false,
			},
		},
	},
	CapreseSalad: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 capreseSaladID,
		Name:               "caprese salad",
		Source:             "",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
		Steps: []*types.RecipeStepDatabaseCreationInput{
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "",
				PreparationID:        validPreparationCollection.Slice.ID,
				BelongsToRecipe:      capreseSaladID,
				ID:                   capreseSaladStep1ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.Tomato.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: capreseSaladStep1ID,
						QuantityValue:       500,
						ProductOfRecipeStep: false,
					},
					{
						IngredientID:        sp(validIngredientCollection.Mozzarella.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: capreseSaladStep1ID,
						QuantityValue:       500,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     0,
				MinEstimatedTimeInSeconds: 120,
				MaxEstimatedTimeInSeconds: 450,
				Optional:                  false,
			},
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "",
				PreparationID:        validPreparationCollection.Plate.ID,
				BelongsToRecipe:      capreseSaladID,
				ID:                   capreseSaladStep2ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.Tomato.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: capreseSaladStep2ID,
						QuantityValue:       500,
						ProductOfRecipeStep: false,
					},
					{
						IngredientID:        sp(validIngredientCollection.Mozzarella.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: capreseSaladStep2ID,
						QuantityValue:       500,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     0,
				MinEstimatedTimeInSeconds: 0,
				MaxEstimatedTimeInSeconds: 0,
				Optional:                  false,
			},
		},
	},
	SpaghettiWithNeatballs: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 spaghettiWithNeatballsID,
		Name:               "spaghetti with neatballs",
		Source:             "",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
		Steps: []*types.RecipeStepDatabaseCreationInput{
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "",
				PreparationID:        validPreparationCollection.Boil.ID,
				BelongsToRecipe:      spaghettiWithNeatballsID,
				ID:                   spaghettiWithNeatballsStep1ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.Pasta.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: spaghettiWithNeatballsStep1ID,
						QuantityValue:       420,
						ProductOfRecipeStep: false,
					},
					{
						IngredientID:        sp(validIngredientCollection.Water.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: spaghettiWithNeatballsStep1ID,
						QuantityValue:       420,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     0,
				MinEstimatedTimeInSeconds: 600,
				MaxEstimatedTimeInSeconds: 900,
				Optional:                  false,
			},
			{
				TemperatureInCelsius: nil,
				Products:             nil,
				Notes:                "",
				PreparationID:        validPreparationCollection.Drain.ID,
				BelongsToRecipe:      spaghettiWithNeatballsID,
				ID:                   spaghettiWithNeatballsStep2ID,
				Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
					{
						IngredientID:        sp(validIngredientCollection.Pasta.ID),
						ID:                  ksuid.New().String(),
						QuantityType:        "grams",
						QuantityNotes:       "",
						IngredientNotes:     "",
						BelongsToRecipeStep: spaghettiWithNeatballsStep1ID,
						QuantityValue:       420,
						ProductOfRecipeStep: false,
					},
				},
				Index:                     0,
				MinEstimatedTimeInSeconds: 0,
				MaxEstimatedTimeInSeconds: 0,
				Optional:                  false,
			},
		},
	},
}

func scaffoldRecipes(ctx context.Context, db database.DataManager) error {
	recipes := []*types.RecipeDatabaseCreationInput{
		recipeCollection.MushroomRisotto,
		recipeCollection.GrilledChicken,
		recipeCollection.CapreseSalad,
		recipeCollection.SpaghettiWithNeatballs,
	}

	for _, input := range recipes {
		if _, err := db.CreateRecipe(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
