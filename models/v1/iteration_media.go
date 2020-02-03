package models

import (
	"context"
	"net/http"
)

type (
	// IterationMedia represents an iteration media
	IterationMedia struct {
		ID                uint64  `json:"id"`
		Path              string  `json:"path"`
		Mimetype          string  `json:"mimetype"`
		RecipeIterationID uint64  `json:"recipe_iteration_id"`
		RecipeStepID      *uint64 `json:"recipe_step_id"`
		CreatedOn         uint64  `json:"created_on"`
		UpdatedOn         *uint64 `json:"updated_on"`
		ArchivedOn        *uint64 `json:"archived_on"`
		BelongsTo         uint64  `json:"belongs_to"`
	}

	// IterationMediaList represents a list of iteration medias
	IterationMediaList struct {
		Pagination
		IterationMedias []IterationMedia `json:"iteration_medias"`
	}

	// IterationMediaCreationInput represents what a user could set as input for creating iteration medias
	IterationMediaCreationInput struct {
		Path              string  `json:"path"`
		Mimetype          string  `json:"mimetype"`
		RecipeIterationID uint64  `json:"recipe_iteration_id"`
		RecipeStepID      *uint64 `json:"recipe_step_id"`
		BelongsTo         uint64  `json:"-"`
	}

	// IterationMediaUpdateInput represents what a user could set as input for updating iteration medias
	IterationMediaUpdateInput struct {
		Path              string  `json:"path"`
		Mimetype          string  `json:"mimetype"`
		RecipeIterationID uint64  `json:"recipe_iteration_id"`
		RecipeStepID      *uint64 `json:"recipe_step_id"`
		BelongsTo         uint64  `json:"-"`
	}

	// IterationMediaDataManager describes a structure capable of storing iteration medias permanently
	IterationMediaDataManager interface {
		GetIterationMedia(ctx context.Context, iterationMediaID, userID uint64) (*IterationMedia, error)
		GetIterationMediaCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllIterationMediasCount(ctx context.Context) (uint64, error)
		GetIterationMedias(ctx context.Context, filter *QueryFilter, userID uint64) (*IterationMediaList, error)
		GetAllIterationMediasForUser(ctx context.Context, userID uint64) ([]IterationMedia, error)
		CreateIterationMedia(ctx context.Context, input *IterationMediaCreationInput) (*IterationMedia, error)
		UpdateIterationMedia(ctx context.Context, updated *IterationMedia) error
		ArchiveIterationMedia(ctx context.Context, id, userID uint64) error
	}

	// IterationMediaDataServer describes a structure capable of serving traffic related to iteration medias
	IterationMediaDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an IterationMediaInput with an iteration media
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

// ToInput creates a IterationMediaUpdateInput struct for an iteration media
func (x *IterationMedia) ToInput() *IterationMediaUpdateInput {
	return &IterationMediaUpdateInput{
		Path:              x.Path,
		Mimetype:          x.Mimetype,
		RecipeIterationID: x.RecipeIterationID,
		RecipeStepID:      x.RecipeStepID,
	}
}
