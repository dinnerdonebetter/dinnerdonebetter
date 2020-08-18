package models

import (
	"context"
	"net/http"
)

type (
	// IterationMedia represents an iteration media.
	IterationMedia struct {
		ID                       uint64  `json:"id"`
		Path                     string  `json:"path"`
		Mimetype                 string  `json:"mimetype"`
		RecipeIterationID        uint64  `json:"recipeIterationID"`
		RecipeStepID             *uint64 `json:"recipeStepID"`
		CreatedOn                uint64  `json:"createdOn"`
		LastUpdatedOn            *uint64 `json:"lastUpdatedOn"`
		ArchivedOn               *uint64 `json:"archivedOn"`
		BelongsToRecipeIteration uint64  `json:"belongsToRecipeIteration"`
	}

	// IterationMediaList represents a list of iteration medias.
	IterationMediaList struct {
		Pagination
		IterationMedias []IterationMedia `json:"iteration_medias"`
	}

	// IterationMediaCreationInput represents what a user could set as input for creating iteration medias.
	IterationMediaCreationInput struct {
		Path                     string  `json:"path"`
		Mimetype                 string  `json:"mimetype"`
		RecipeIterationID        uint64  `json:"recipeIterationID"`
		RecipeStepID             *uint64 `json:"recipeStepID"`
		BelongsToRecipeIteration uint64  `json:"-"`
	}

	// IterationMediaUpdateInput represents what a user could set as input for updating iteration medias.
	IterationMediaUpdateInput struct {
		Path                     string  `json:"path"`
		Mimetype                 string  `json:"mimetype"`
		RecipeIterationID        uint64  `json:"recipeIterationID"`
		RecipeStepID             *uint64 `json:"recipeStepID"`
		BelongsToRecipeIteration uint64  `json:"belongsToRecipeIteration"`
	}

	// IterationMediaDataManager describes a structure capable of storing iteration medias permanently.
	IterationMediaDataManager interface {
		IterationMediaExists(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (bool, error)
		GetIterationMedia(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*IterationMedia, error)
		GetAllIterationMediasCount(ctx context.Context) (uint64, error)
		GetAllIterationMedias(ctx context.Context, resultChannel chan []IterationMedia) error
		GetIterationMedias(ctx context.Context, recipeID, recipeIterationID uint64, filter *QueryFilter) (*IterationMediaList, error)
		GetIterationMediasWithIDs(ctx context.Context, recipeID, recipeIterationID uint64, limit uint8, ids []uint64) ([]IterationMedia, error)
		CreateIterationMedia(ctx context.Context, input *IterationMediaCreationInput) (*IterationMedia, error)
		UpdateIterationMedia(ctx context.Context, updated *IterationMedia) error
		ArchiveIterationMedia(ctx context.Context, recipeIterationID, iterationMediaID uint64) error
	}

	// IterationMediaDataServer describes a structure capable of serving traffic related to iteration medias.
	IterationMediaDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an IterationMediaInput with an iteration media.
func (x *IterationMedia) Update(input *IterationMediaUpdateInput) {
	if input.Path != "" && input.Path != x.Path {
		x.Path = input.Path
	}

	if input.Mimetype != "" && input.Mimetype != x.Mimetype {
		x.Mimetype = input.Mimetype
	}

	if input.RecipeIterationID != x.RecipeIterationID {
		x.RecipeIterationID = input.RecipeIterationID
	}

	if input.RecipeStepID != nil && input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = input.RecipeStepID
	}
}

// ToUpdateInput creates a IterationMediaUpdateInput struct for an iteration media.
func (x *IterationMedia) ToUpdateInput() *IterationMediaUpdateInput {
	return &IterationMediaUpdateInput{
		Path:              x.Path,
		Mimetype:          x.Mimetype,
		RecipeIterationID: x.RecipeIterationID,
		RecipeStepID:      x.RecipeStepID,
	}
}
