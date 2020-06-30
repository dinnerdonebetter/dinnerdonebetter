package models

import (
	"context"
	"net/http"
)

type (
	// ValidIngredientPreparation represents a valid ingredient preparation.
	ValidIngredientPreparation struct {
		ID                       uint64  `json:"id"`
		Notes                    string  `json:"notes"`
		CreatedOn                uint64  `json:"createdOn"`
		UpdatedOn                *uint64 `json:"updatedOn"`
		ArchivedOn               *uint64 `json:"archivedOn"`
		BelongsToValidIngredient uint64  `json:"belongsToValidIngredient"`
	}

	// ValidIngredientPreparationList represents a list of valid ingredient preparations.
	ValidIngredientPreparationList struct {
		Pagination
		ValidIngredientPreparations []ValidIngredientPreparation `json:"validIngredientPreparations"`
	}

	// ValidIngredientPreparationCreationInput represents what a user could set as input for creating valid ingredient preparations.
	ValidIngredientPreparationCreationInput struct {
		Notes                    string `json:"notes"`
		BelongsToValidIngredient uint64 `json:"-"`
	}

	// ValidIngredientPreparationUpdateInput represents what a user could set as input for updating valid ingredient preparations.
	ValidIngredientPreparationUpdateInput struct {
		Notes                    string `json:"notes"`
		BelongsToValidIngredient uint64 `json:"belongsToValidIngredient"`
	}

	// ValidIngredientPreparationDataManager describes a structure capable of storing valid ingredient preparations permanently.
	ValidIngredientPreparationDataManager interface {
		ValidIngredientPreparationExists(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) (bool, error)
		GetValidIngredientPreparation(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) (*ValidIngredientPreparation, error)
		GetAllValidIngredientPreparationsCount(ctx context.Context) (uint64, error)
		GetValidIngredientPreparations(ctx context.Context, validIngredientID uint64, filter *QueryFilter) (*ValidIngredientPreparationList, error)
		CreateValidIngredientPreparation(ctx context.Context, input *ValidIngredientPreparationCreationInput) (*ValidIngredientPreparation, error)
		UpdateValidIngredientPreparation(ctx context.Context, updated *ValidIngredientPreparation) error
		ArchiveValidIngredientPreparation(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) error
	}

	// ValidIngredientPreparationDataServer describes a structure capable of serving traffic related to valid ingredient preparations.
	ValidIngredientPreparationDataServer interface {
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

// Update merges an ValidIngredientPreparationInput with a valid ingredient preparation.
func (x *ValidIngredientPreparation) Update(input *ValidIngredientPreparationUpdateInput) {
	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

// ToUpdateInput creates a ValidIngredientPreparationUpdateInput struct for a valid ingredient preparation.
func (x *ValidIngredientPreparation) ToUpdateInput() *ValidIngredientPreparationUpdateInput {
	return &ValidIngredientPreparationUpdateInput{
		Notes: x.Notes,
	}
}
