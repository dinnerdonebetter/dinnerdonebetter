package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *types.RecipeStep {
	recipeStepID := ksuid.New().String()

	var ingredients []*types.RecipeStepIngredient
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepIngredient()
		ing.BelongsToRecipeStep = recipeStepID
		ingredients = append(ingredients, ing)
	}

	var products []*types.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		p := BuildFakeRecipeStepProduct()
		p.BelongsToRecipeStep = recipeStepID
		products = append(products, p)
	}

	return &types.RecipeStep{
		ID:                        recipeStepID,
		Index:                     fake.Uint32(),
		Preparation:               *BuildFakeValidPreparation(),
		MinEstimatedTimeInSeconds: fake.Uint32(),
		MaxEstimatedTimeInSeconds: fake.Uint32(),
		TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
		Notes:                     fake.LoremIpsumSentence(exampleQuantity),
		Products:                  products,
		Optional:                  false,
		CreatedOn:                 uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe:           ksuid.New().String(),
		Ingredients:               ingredients,
	}
}

// BuildFakeRecipeStepList builds a faked RecipeStepList.
func BuildFakeRecipeStepList() *types.RecipeStepList {
	var examples []*types.RecipeStep
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStep())
	}

	return &types.RecipeStepList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeSteps: examples,
	}
}

// BuildFakeRecipeStepUpdateRequestInput builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInput() *types.RecipeStepUpdateRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return &types.RecipeStepUpdateRequestInput{
		Index:                     recipeStep.Index,
		Preparation:               recipeStep.Preparation,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Products:                  recipeStep.Products,
		Optional:                  recipeStep.Optional,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepUpdateRequestInputFromRecipeStep builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	return &types.RecipeStepUpdateRequestInput{
		Products:                  recipeStep.Products,
		Optional:                  recipeStep.Optional,
		Index:                     recipeStep.Index,
		Preparation:               recipeStep.Preparation,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepCreationRequestInput builds a faked RecipeStepCreationRequestInput.
func BuildFakeRecipeStepCreationRequestInput() *types.RecipeStepCreationRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return BuildFakeRecipeStepCreationRequestInputFromRecipeStep(recipeStep)
}

// BuildFakeRecipeStepCreationRequestInputFromRecipeStep builds a faked RecipeStepCreationRequestInput from a recipe step.
func BuildFakeRecipeStepCreationRequestInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepCreationRequestInput {
	ingredients := []*types.RecipeStepIngredientCreationRequestInput{}
	for _, ingredient := range recipeStep.Ingredients {
		ingredients = append(ingredients, BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(ingredient))
	}

	products := []*types.RecipeStepProductCreationRequestInput{}
	for _, product := range recipeStep.Products {
		products = append(products, BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(product))
	}

	return &types.RecipeStepCreationRequestInput{
		ID:                        recipeStep.ID,
		Products:                  products,
		Optional:                  recipeStep.Optional,
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.Preparation.ID,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
		Ingredients:               ingredients,
	}
}

// BuildFakeRecipeStepDatabaseCreationInput builds a faked RecipeStepDatabaseCreationInput.
func BuildFakeRecipeStepDatabaseCreationInput() *types.RecipeStepDatabaseCreationInput {
	recipeStep := BuildFakeRecipeStep()
	return BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep(recipeStep)
}

// BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep builds a faked RecipeStepDatabaseCreationInput from a recipe step.
func BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepDatabaseCreationInput {
	ingredients := []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, i := range recipeStep.Ingredients {
		ingredients = append(ingredients, BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(i))
	}

	products := []*types.RecipeStepProductDatabaseCreationInput{}
	for _, p := range recipeStep.Products {
		products = append(products, BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct(p))
	}

	return &types.RecipeStepDatabaseCreationInput{
		ID:                        recipeStep.ID,
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.Preparation.ID,
		Optional:                  recipeStep.Optional,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Ingredients:               ingredients,
		Products:                  products,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}
