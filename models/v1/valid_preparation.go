package models

import (
	"context"
	"net/http"
)

type (
	// ValidPreparation represents a valid preparation.
	ValidPreparation struct {
		ID                         uint64  `json:"id"`
		Name                       string  `json:"name"`
		Description                string  `json:"description"`
		Icon                       string  `json:"icon"`
		ApplicableToAllIngredients bool    `json:"applicableToAllIngredients"`
		CreatedOn                  uint64  `json:"createdOn"`
		UpdatedOn                  *uint64 `json:"updatedOn"`
		ArchivedOn                 *uint64 `json:"archivedOn"`
	}

	// ValidPreparationList represents a list of valid preparations.
	ValidPreparationList struct {
		Pagination
		ValidPreparations []ValidPreparation `json:"validPreparations"`
	}

	// ValidPreparationCreationInput represents what a user could set as input for creating valid preparations.
	ValidPreparationCreationInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	// ValidPreparationUpdateInput represents what a user could set as input for updating valid preparations.
	ValidPreparationUpdateInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	// ValidPreparationDataManager describes a structure capable of storing valid preparations permanently.
	ValidPreparationDataManager interface {
		ValidPreparationExists(ctx context.Context, validPreparationID uint64) (bool, error)
		GetValidPreparation(ctx context.Context, validPreparationID uint64) (*ValidPreparation, error)
		GetAllValidPreparationsCount(ctx context.Context) (uint64, error)
		GetValidPreparations(ctx context.Context, filter *QueryFilter) (*ValidPreparationList, error)
		CreateValidPreparation(ctx context.Context, input *ValidPreparationCreationInput) (*ValidPreparation, error)
		UpdateValidPreparation(ctx context.Context, updated *ValidPreparation) error
		ArchiveValidPreparation(ctx context.Context, validPreparationID uint64) error
	}

	// ValidPreparationDataServer describes a structure capable of serving traffic related to valid preparations.
	ValidPreparationDataServer interface {
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

// Update merges an ValidPreparationInput with a valid preparation.
func (x *ValidPreparation) Update(input *ValidPreparationUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.Icon != "" && input.Icon != x.Icon {
		x.Icon = input.Icon
	}
}

// ToUpdateInput creates a ValidPreparationUpdateInput struct for a valid preparation.
func (x *ValidPreparation) ToUpdateInput() *ValidPreparationUpdateInput {
	return &ValidPreparationUpdateInput{
		Name:        x.Name,
		Description: x.Description,
		Icon:        x.Icon,
	}
}
