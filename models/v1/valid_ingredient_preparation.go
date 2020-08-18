package models

import (
	"context"
	"net/http"
)

type (
	// ValidIngredientPreparation represents a valid ingredient preparation.
	ValidIngredientPreparation struct {
		ID                 uint64  `json:"id"`
		Notes              string  `json:"notes"`
		ValidPreparationID uint64  `json:"validPreparationID"`
		ValidIngredientID  uint64  `json:"validIngredientID"`
		CreatedOn          uint64  `json:"createdOn"`
		LastUpdatedOn      *uint64 `json:"lastUpdatedOn"`
		ArchivedOn         *uint64 `json:"archivedOn"`
	}

	// ValidIngredientPreparationList represents a list of valid ingredient preparations.
	ValidIngredientPreparationList struct {
		Pagination
		ValidIngredientPreparations []ValidIngredientPreparation `json:"valid_ingredient_preparations"`
	}

	// ValidIngredientPreparationCreationInput represents what a user could set as input for creating valid ingredient preparations.
	ValidIngredientPreparationCreationInput struct {
		Notes              string `json:"notes"`
		ValidPreparationID uint64 `json:"validPreparationID"`
		ValidIngredientID  uint64 `json:"validIngredientID"`
	}

	// ValidIngredientPreparationUpdateInput represents what a user could set as input for updating valid ingredient preparations.
	ValidIngredientPreparationUpdateInput struct {
		Notes              string `json:"notes"`
		ValidPreparationID uint64 `json:"validPreparationID"`
		ValidIngredientID  uint64 `json:"validIngredientID"`
	}

	// ValidIngredientPreparationDataManager describes a structure capable of storing valid ingredient preparations permanently.
	ValidIngredientPreparationDataManager interface {
		ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (bool, error)
		GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (*ValidIngredientPreparation, error)
		GetAllValidIngredientPreparationsCount(ctx context.Context) (uint64, error)
		GetAllValidIngredientPreparations(ctx context.Context, resultChannel chan []ValidIngredientPreparation) error
		GetValidIngredientPreparations(ctx context.Context, filter *QueryFilter) (*ValidIngredientPreparationList, error)
		GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]ValidIngredientPreparation, error)
		CreateValidIngredientPreparation(ctx context.Context, input *ValidIngredientPreparationCreationInput) (*ValidIngredientPreparation, error)
		UpdateValidIngredientPreparation(ctx context.Context, updated *ValidIngredientPreparation) error
		ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) error
	}

	// ValidIngredientPreparationDataServer describes a structure capable of serving traffic related to valid ingredient preparations.
	ValidIngredientPreparationDataServer interface {
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

// Update merges an ValidIngredientPreparationInput with a valid ingredient preparation.
func (x *ValidIngredientPreparation) Update(input *ValidIngredientPreparationUpdateInput) {
	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}

	if input.ValidPreparationID != x.ValidPreparationID {
		x.ValidPreparationID = input.ValidPreparationID
	}

	if input.ValidIngredientID != x.ValidIngredientID {
		x.ValidIngredientID = input.ValidIngredientID
	}
}

// ToUpdateInput creates a ValidIngredientPreparationUpdateInput struct for a valid ingredient preparation.
func (x *ValidIngredientPreparation) ToUpdateInput() *ValidIngredientPreparationUpdateInput {
	return &ValidIngredientPreparationUpdateInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidIngredientID:  x.ValidIngredientID,
	}
}
