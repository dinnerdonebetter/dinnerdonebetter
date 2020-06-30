package models

import (
	"context"
	"net/http"
)

type (
	// IngredientTagMapping represents an ingredient tag mapping.
	IngredientTagMapping struct {
		ID                       uint64  `json:"id"`
		ValidIngredientTagID     uint64  `json:"validIngredientTagID"`
		CreatedOn                uint64  `json:"createdOn"`
		UpdatedOn                *uint64 `json:"updatedOn"`
		ArchivedOn               *uint64 `json:"archivedOn"`
		BelongsToValidIngredient uint64  `json:"belongsToValidIngredient"`
	}

	// IngredientTagMappingList represents a list of ingredient tag mappings.
	IngredientTagMappingList struct {
		Pagination
		IngredientTagMappings []IngredientTagMapping `json:"ingredientTagMappings"`
	}

	// IngredientTagMappingCreationInput represents what a user could set as input for creating ingredient tag mappings.
	IngredientTagMappingCreationInput struct {
		ValidIngredientTagID     uint64 `json:"validIngredientTagId"`
		BelongsToValidIngredient uint64 `json:"-"`
	}

	// IngredientTagMappingUpdateInput represents what a user could set as input for updating ingredient tag mappings.
	IngredientTagMappingUpdateInput struct {
		ValidIngredientTagID     uint64 `json:"validIngredientTagId"`
		BelongsToValidIngredient uint64 `json:"belongsToValidIngredient"`
	}

	// IngredientTagMappingDataManager describes a structure capable of storing ingredient tag mappings permanently.
	IngredientTagMappingDataManager interface {
		IngredientTagMappingExists(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (bool, error)
		GetIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (*IngredientTagMapping, error)
		GetAllIngredientTagMappingsCount(ctx context.Context) (uint64, error)
		GetIngredientTagMappings(ctx context.Context, validIngredientID uint64, filter *QueryFilter) (*IngredientTagMappingList, error)
		CreateIngredientTagMapping(ctx context.Context, input *IngredientTagMappingCreationInput) (*IngredientTagMapping, error)
		UpdateIngredientTagMapping(ctx context.Context, updated *IngredientTagMapping) error
		ArchiveIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) error
	}

	// IngredientTagMappingDataServer describes a structure capable of serving traffic related to ingredient tag mappings.
	IngredientTagMappingDataServer interface {
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

// Update merges an IngredientTagMappingInput with an ingredient tag mapping.
func (x *IngredientTagMapping) Update(input *IngredientTagMappingUpdateInput) {
	if input.ValidIngredientTagID != x.ValidIngredientTagID {
		x.ValidIngredientTagID = input.ValidIngredientTagID
	}
}

// ToUpdateInput creates a IngredientTagMappingUpdateInput struct for an ingredient tag mapping.
func (x *IngredientTagMapping) ToUpdateInput() *IngredientTagMappingUpdateInput {
	return &IngredientTagMappingUpdateInput{
		ValidIngredientTagID: x.ValidIngredientTagID,
	}
}
