package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeIterationMedia builds a faked iteration media.
func BuildFakeIterationMedia() *models.IterationMedia {
	return &models.IterationMedia{
		ID:                       fake.Uint64(),
		Source:                   fake.Word(),
		Mimetype:                 fake.Word(),
		CreatedOn:                uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeIteration: fake.Uint64(),
	}
}

// BuildFakeIterationMediaList builds a faked IterationMediaList.
func BuildFakeIterationMediaList() *models.IterationMediaList {
	exampleIterationMedia1 := BuildFakeIterationMedia()
	exampleIterationMedia2 := BuildFakeIterationMedia()
	exampleIterationMedia3 := BuildFakeIterationMedia()

	return &models.IterationMediaList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		IterationMedia: []models.IterationMedia{
			*exampleIterationMedia1,
			*exampleIterationMedia2,
			*exampleIterationMedia3,
		},
	}
}

// BuildFakeIterationMediaUpdateInputFromIterationMedia builds a faked IterationMediaUpdateInput from an iteration media.
func BuildFakeIterationMediaUpdateInputFromIterationMedia(iterationMedia *models.IterationMedia) *models.IterationMediaUpdateInput {
	return &models.IterationMediaUpdateInput{
		Source:                   iterationMedia.Source,
		Mimetype:                 iterationMedia.Mimetype,
		BelongsToRecipeIteration: iterationMedia.BelongsToRecipeIteration,
	}
}

// BuildFakeIterationMediaCreationInput builds a faked IterationMediaCreationInput.
func BuildFakeIterationMediaCreationInput() *models.IterationMediaCreationInput {
	iterationMedia := BuildFakeIterationMedia()
	return BuildFakeIterationMediaCreationInputFromIterationMedia(iterationMedia)
}

// BuildFakeIterationMediaCreationInputFromIterationMedia builds a faked IterationMediaCreationInput from an iteration media.
func BuildFakeIterationMediaCreationInputFromIterationMedia(iterationMedia *models.IterationMedia) *models.IterationMediaCreationInput {
	return &models.IterationMediaCreationInput{
		Source:                   iterationMedia.Source,
		Mimetype:                 iterationMedia.Mimetype,
		BelongsToRecipeIteration: iterationMedia.BelongsToRecipeIteration,
	}
}
