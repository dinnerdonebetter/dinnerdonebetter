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

	spaghettiID               = ksuid.New().String()
	grilledCheeseSandwichesID = ksuid.New().String()
	bakedPotatoID             = ksuid.New().String()
	ramenID                   = ksuid.New().String()
	lasagnaID                 = ksuid.New().String()
	tacosID                   = ksuid.New().String()
	eggFriedRiceID            = ksuid.New().String()
	mashedPotatoesID          = ksuid.New().String()
	collardGreensID           = ksuid.New().String()
	neatballsID               = ksuid.New().String()
)

var recipeCollection = struct {
	MushroomRisotto,
	GrilledChicken,
	CapreseSalad,
	Spaghetti,
	GrilledCheeseSandwiches,
	BakedPotato,
	Ramen,
	Lasagna,
	Tacos,
	EggFriedRice,
	MashedPotatoes,
	CollardGreens,
	Neatballs *types.RecipeDatabaseCreationInput
}{
	MushroomRisotto: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 mushroomRisottoID,
		Name:               "mushroom risotto",
		Source:             "https://www.google.com/search?&q=mushroom+risotto",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
		Steps: []*types.RecipeStepDatabaseCreationInput{
			{
				MinimumTemperatureInCelsius: nil,
				Products:                    nil,
				Notes:                       "",
				PreparationID:               validPreparationCollection.Dice.ID,
				BelongsToRecipe:             mushroomRisottoID,
				ID:                          mushroomRisottoStep1ID,
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
				Index:                         0,
				MinimumEstimatedTimeInSeconds: 0,
				MaximumEstimatedTimeInSeconds: 0,
				Optional:                      false,
			},
			{
				MinimumTemperatureInCelsius: nil,
				Products:                    nil,
				Notes:                       "helps avoid burning the onion",
				PreparationID:               validPreparationCollection.Sautee.ID,
				BelongsToRecipe:             mushroomRisottoID,
				ID:                          mushroomRisottoStep2ID,
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
				Index:                         1,
				MinimumEstimatedTimeInSeconds: 0,
				MaximumEstimatedTimeInSeconds: 0,
				Optional:                      false,
			},
		},
	},
	GrilledChicken: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 grilledChickenID,
		Name:               "grilled chicken",
		Source:             "https://www.google.com/search?&q=grilled+chicken",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
		Steps: []*types.RecipeStepDatabaseCreationInput{
			{
				MinimumTemperatureInCelsius: nil,
				Products:                    nil,
				Notes:                       "",
				PreparationID:               validPreparationCollection.Marinate.ID,
				BelongsToRecipe:             grilledChickenID,
				ID:                          grilledChickenStep1ID,
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
				Index:                         0,
				MinimumEstimatedTimeInSeconds: 0,
				MaximumEstimatedTimeInSeconds: 0,
				Optional:                      false,
			},
			{
				MinimumTemperatureInCelsius: nil,
				Products:                    nil,
				Notes:                       "",
				PreparationID:               validPreparationCollection.Grill.ID,
				BelongsToRecipe:             grilledChickenID,
				ID:                          grilledChickenStep2ID,
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
				Index:                         0,
				MinimumEstimatedTimeInSeconds: 0,
				MaximumEstimatedTimeInSeconds: 0,
				Optional:                      false,
			},
		},
	},
	CapreseSalad: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 capreseSaladID,
		Name:               "caprese salad",
		Source:             "https://www.google.com/search?&q=caprese+salad",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
		Steps: []*types.RecipeStepDatabaseCreationInput{
			{
				MinimumTemperatureInCelsius: nil,
				Products:                    nil,
				Notes:                       "",
				PreparationID:               validPreparationCollection.Slice.ID,
				BelongsToRecipe:             capreseSaladID,
				ID:                          capreseSaladStep1ID,
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
				Index:                         0,
				MinimumEstimatedTimeInSeconds: 120,
				MaximumEstimatedTimeInSeconds: 450,
				Optional:                      false,
			},
			{
				MinimumTemperatureInCelsius: nil,
				Products:                    nil,
				Notes:                       "",
				PreparationID:               validPreparationCollection.Plate.ID,
				BelongsToRecipe:             capreseSaladID,
				ID:                          capreseSaladStep2ID,
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
				Index:                         0,
				MinimumEstimatedTimeInSeconds: 0,
				MaximumEstimatedTimeInSeconds: 0,
				Optional:                      false,
			},
		},
	},
	Spaghetti: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 spaghettiID,
		Name:               "spaghetti",
		Source:             "https://www.google.com/search?&q=spaghetti+with+neatballs",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	Neatballs: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 neatballsID,
		Name:               "neatballs",
		Source:             "https://www.google.com/search?&q=neatballs",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},

	GrilledCheeseSandwiches: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 grilledCheeseSandwichesID,
		Name:               "grilled cheese sandwiches",
		Source:             "https://www.google.com/search?&q=grilled+cheese+sandwiches",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	BakedPotato: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 bakedPotatoID,
		Name:               "baked potato",
		Source:             "https://www.google.com/search?&q=baked+potato",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	Ramen: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 ramenID,
		Name:               "ramen",
		Source:             "https://www.google.com/search?&q=ramen",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	Lasagna: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 lasagnaID,
		Name:               "lasagna",
		Source:             "https://www.google.com/search?&q=lasagna",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	Tacos: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 tacosID,
		Name:               "tacos",
		Source:             "https://www.google.com/search?&q=tacos",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	EggFriedRice: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 eggFriedRiceID,
		Name:               "egg fried rice",
		Source:             "https://www.google.com/search?&q=egg+fried+rice",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	MashedPotatoes: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 mashedPotatoesID,
		Name:               "mashed potatoes",
		Source:             "https://www.google.com/search?&q=mashed+potatoes",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
	CollardGreens: &types.RecipeDatabaseCreationInput{
		InspiredByRecipeID: nil,
		ID:                 collardGreensID,
		Name:               "collard greens",
		Source:             "https://www.google.com/search?&q=collard+greens",
		Description:        "",
		CreatedByUser:      userCollection.MomJones.ID,
	},
}

func scaffoldRecipes(ctx context.Context, db database.DataManager) error {
	recipes := []*types.RecipeDatabaseCreationInput{
		recipeCollection.MushroomRisotto,
		recipeCollection.GrilledChicken,
		recipeCollection.CapreseSalad,
		recipeCollection.Spaghetti,
		recipeCollection.GrilledCheeseSandwiches,
		recipeCollection.BakedPotato,
		recipeCollection.Ramen,
		recipeCollection.Lasagna,
		recipeCollection.Tacos,
		recipeCollection.EggFriedRice,
		recipeCollection.MashedPotatoes,
		recipeCollection.CollardGreens,
		recipeCollection.Neatballs,
	}

	for _, input := range recipes {
		if _, err := db.CreateRecipe(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
