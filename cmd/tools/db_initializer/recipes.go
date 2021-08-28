package main

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

func u16p(u uint16) *uint16 {
	return &u
}

func u64p(u uint64) *uint64 {
	return &u
}

var (
	recipes = []*types.RecipeCreationInput{
		{
			Name:        "morning coffee",
			Description: "what I make in the morning to wake up",
			Steps: []*types.RecipeStepCreationInput{
				{
					Index:                0,
					TemperatureInCelsius: u16p(200),
					Ingredients: []*types.RecipeStepIngredientCreationInput{
						{
							IngredientID:  u64p(1),
							QuantityType:  "grams",
							QuantityValue: 900,
							QuantityNotes: "",
						},
					},
					PreparationID: 1,
				},
				{
					Index: 1,
					Ingredients: []*types.RecipeStepIngredientCreationInput{
						{
							IngredientID:  u64p(2),
							QuantityType:  "grams",
							QuantityValue: 30,
							QuantityNotes: "",
						},
					},
					PreparationID: 2,
				},
				{
					Index: 2,
					Ingredients: []*types.RecipeStepIngredientCreationInput{
						{
							IngredientID:  u64p(1),
							QuantityType:  "grams",
							QuantityValue: 900,
							QuantityNotes: "",
						},
						{
							IngredientID:  u64p(2),
							QuantityType:  "grams",
							QuantityValue: 30,
							QuantityNotes: "",
						},
					},
					PreparationID: 3,
				},
			},
		},
	}
)
