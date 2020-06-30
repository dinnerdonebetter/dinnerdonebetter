package models

import (
	"context"
	"net/http"
)

type (
	// IterationMedia represents an iteration media.
	IterationMedia struct {
		ID                       uint64  `json:"id"`
		Source                   string  `json:"source"`
		Mimetype                 string  `json:"mimetype"`
		CreatedOn                uint64  `json:"createdOn"`
		UpdatedOn                *uint64 `json:"updatedOn"`
		ArchivedOn               *uint64 `json:"archivedOn"`
		BelongsToRecipeIteration uint64  `json:"belongsToRecipeIteration"`
	}

	// IterationMediaList represents a list of iteration medias.
	IterationMediaList struct {
		Pagination
		IterationMedia []IterationMedia `json:"iterationMedia"`
	}

	// IterationMediaCreationInput represents what a user could set as input for creating iteration medias.
	IterationMediaCreationInput struct {
		Source                   string `json:"source"`
		Mimetype                 string `json:"mimetype"`
		BelongsToRecipeIteration uint64 `json:"-"`
	}

	// IterationMediaUpdateInput represents what a user could set as input for updating iteration medias.
	IterationMediaUpdateInput struct {
		Source                   string `json:"source"`
		Mimetype                 string `json:"mimetype"`
		BelongsToRecipeIteration uint64 `json:"belongsToRecipeIteration"`
	}

	// IterationMediaDataManager describes a structure capable of storing iteration medias permanently.
	IterationMediaDataManager interface {
		IterationMediaExists(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (bool, error)
		GetIterationMedia(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*IterationMedia, error)
		GetAllIterationMediasCount(ctx context.Context) (uint64, error)
		GetIterationMedias(ctx context.Context, recipeID, recipeIterationID uint64, filter *QueryFilter) (*IterationMediaList, error)
		CreateIterationMedia(ctx context.Context, input *IterationMediaCreationInput) (*IterationMedia, error)
		UpdateIterationMedia(ctx context.Context, updated *IterationMedia) error
		ArchiveIterationMedia(ctx context.Context, recipeIterationID, iterationMediaID uint64) error
	}

	// IterationMediaDataServer describes a structure capable of serving traffic related to iteration medias.
	IterationMediaDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ExistenceHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an IterationMediaInput with an iteration media.
func (x *IterationMedia) Update(input *IterationMediaUpdateInput) {
	if input.Source != "" && input.Source != x.Source {
		x.Source = input.Source
	}

	if input.Mimetype != "" && input.Mimetype != x.Mimetype {
		x.Mimetype = input.Mimetype
	}
}

// ToUpdateInput creates a IterationMediaUpdateInput struct for an iteration media.
func (x *IterationMedia) ToUpdateInput() *IterationMediaUpdateInput {
	return &IterationMediaUpdateInput{
		Source:   x.Source,
		Mimetype: x.Mimetype,
	}
}
