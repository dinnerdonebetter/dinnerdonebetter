package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeIngredientTagMapping builds a faked ingredient tag mapping.
func BuildFakeIngredientTagMapping() *models.IngredientTagMapping {
	return &models.IngredientTagMapping{
		ID:                       fake.Uint64(),
		ValidIngredientTagID:     uint64(fake.Uint32()),
		CreatedOn:                uint64(uint32(fake.Date().Unix())),
		BelongsToValidIngredient: fake.Uint64(),
	}
}

// BuildFakeIngredientTagMappingList builds a faked IngredientTagMappingList.
func BuildFakeIngredientTagMappingList() *models.IngredientTagMappingList {
	exampleIngredientTagMapping1 := BuildFakeIngredientTagMapping()
	exampleIngredientTagMapping2 := BuildFakeIngredientTagMapping()
	exampleIngredientTagMapping3 := BuildFakeIngredientTagMapping()

	return &models.IngredientTagMappingList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		IngredientTagMappings: []models.IngredientTagMapping{
			*exampleIngredientTagMapping1,
			*exampleIngredientTagMapping2,
			*exampleIngredientTagMapping3,
		},
	}
}

// BuildFakeIngredientTagMappingUpdateInputFromIngredientTagMapping builds a faked IngredientTagMappingUpdateInput from an ingredient tag mapping.
func BuildFakeIngredientTagMappingUpdateInputFromIngredientTagMapping(ingredientTagMapping *models.IngredientTagMapping) *models.IngredientTagMappingUpdateInput {
	return &models.IngredientTagMappingUpdateInput{
		ValidIngredientTagID:     ingredientTagMapping.ValidIngredientTagID,
		BelongsToValidIngredient: ingredientTagMapping.BelongsToValidIngredient,
	}
}

// BuildFakeIngredientTagMappingCreationInput builds a faked IngredientTagMappingCreationInput.
func BuildFakeIngredientTagMappingCreationInput() *models.IngredientTagMappingCreationInput {
	ingredientTagMapping := BuildFakeIngredientTagMapping()
	return BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(ingredientTagMapping)
}

// BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping builds a faked IngredientTagMappingCreationInput from an ingredient tag mapping.
func BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(ingredientTagMapping *models.IngredientTagMapping) *models.IngredientTagMappingCreationInput {
	return &models.IngredientTagMappingCreationInput{
		ValidIngredientTagID:     ingredientTagMapping.ValidIngredientTagID,
		BelongsToValidIngredient: ingredientTagMapping.BelongsToValidIngredient,
	}
}
