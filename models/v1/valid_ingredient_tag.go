package models

import (
	"context"
	"net/http"
)

type (
	// ValidIngredientTag represents a valid ingredient tag.
	ValidIngredientTag struct {
		ID         uint64  `json:"id"`
		Name       string  `json:"name"`
		CreatedOn  uint64  `json:"created_on"`
		UpdatedOn  *uint64 `json:"updated_on"`
		ArchivedOn *uint64 `json:"archived_on"`
	}

	// ValidIngredientTagList represents a list of valid ingredient tags.
	ValidIngredientTagList struct {
		Pagination
		ValidIngredientTags []ValidIngredientTag `json:"valid_ingredient_tags"`
	}

	// ValidIngredientTagCreationInput represents what a user could set as input for creating valid ingredient tags.
	ValidIngredientTagCreationInput struct {
		Name string `json:"name"`
	}

	// ValidIngredientTagUpdateInput represents what a user could set as input for updating valid ingredient tags.
	ValidIngredientTagUpdateInput struct {
		Name string `json:"name"`
	}

	// ValidIngredientTagDataManager describes a structure capable of storing valid ingredient tags permanently.
	ValidIngredientTagDataManager interface {
		ValidIngredientTagExists(ctx context.Context, validIngredientTagID uint64) (bool, error)
		GetValidIngredientTag(ctx context.Context, validIngredientTagID uint64) (*ValidIngredientTag, error)
		GetAllValidIngredientTagsCount(ctx context.Context) (uint64, error)
		GetValidIngredientTags(ctx context.Context, filter *QueryFilter) (*ValidIngredientTagList, error)
		CreateValidIngredientTag(ctx context.Context, input *ValidIngredientTagCreationInput) (*ValidIngredientTag, error)
		UpdateValidIngredientTag(ctx context.Context, updated *ValidIngredientTag) error
		ArchiveValidIngredientTag(ctx context.Context, validIngredientTagID uint64) error
	}

	// ValidIngredientTagDataServer describes a structure capable of serving traffic related to valid ingredient tags.
	ValidIngredientTagDataServer interface {
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

// Update merges an ValidIngredientTagInput with a valid ingredient tag.
func (x *ValidIngredientTag) Update(input *ValidIngredientTagUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}
}

// ToUpdateInput creates a ValidIngredientTagUpdateInput struct for a valid ingredient tag.
func (x *ValidIngredientTag) ToUpdateInput() *ValidIngredientTagUpdateInput {
	return &ValidIngredientTagUpdateInput{
		Name: x.Name,
	}
}
