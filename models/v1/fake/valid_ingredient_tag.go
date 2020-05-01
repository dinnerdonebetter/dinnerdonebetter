package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidIngredientTag builds a faked valid ingredient tag.
func BuildFakeValidIngredientTag() *models.ValidIngredientTag {
	return &models.ValidIngredientTag{
		ID:        fake.Uint64(),
		Name:      fake.Word(),
		CreatedOn: uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidIngredientTagList builds a faked ValidIngredientTagList.
func BuildFakeValidIngredientTagList() *models.ValidIngredientTagList {
	exampleValidIngredientTag1 := BuildFakeValidIngredientTag()
	exampleValidIngredientTag2 := BuildFakeValidIngredientTag()
	exampleValidIngredientTag3 := BuildFakeValidIngredientTag()

	return &models.ValidIngredientTagList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		ValidIngredientTags: []models.ValidIngredientTag{
			*exampleValidIngredientTag1,
			*exampleValidIngredientTag2,
			*exampleValidIngredientTag3,
		},
	}
}

// BuildFakeValidIngredientTagUpdateInputFromValidIngredientTag builds a faked ValidIngredientTagUpdateInput from a valid ingredient tag.
func BuildFakeValidIngredientTagUpdateInputFromValidIngredientTag(validIngredientTag *models.ValidIngredientTag) *models.ValidIngredientTagUpdateInput {
	return &models.ValidIngredientTagUpdateInput{
		Name: validIngredientTag.Name,
	}
}

// BuildFakeValidIngredientTagCreationInput builds a faked ValidIngredientTagCreationInput.
func BuildFakeValidIngredientTagCreationInput() *models.ValidIngredientTagCreationInput {
	validIngredientTag := BuildFakeValidIngredientTag()
	return BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(validIngredientTag)
}

// BuildFakeValidIngredientTagCreationInputFromValidIngredientTag builds a faked ValidIngredientTagCreationInput from a valid ingredient tag.
func BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(validIngredientTag *models.ValidIngredientTag) *models.ValidIngredientTagCreationInput {
	return &models.ValidIngredientTagCreationInput{
		Name: validIngredientTag.Name,
	}
}
