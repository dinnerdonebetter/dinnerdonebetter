package models

import (
	"context"
	"net/http"
)

type (
	// Preparation represents a preparation
	Preparation struct {
		ID             uint64  `json:"id"`
		Name           string  `json:"name"`
		Variant        string  `json:"variant"`
		Description    string  `json:"description"`
		AllergyWarning string  `json:"allergy_warning"`
		Icon           string  `json:"icon"`
		CreatedOn      uint64  `json:"created_on"`
		UpdatedOn      *uint64 `json:"updated_on"`
		ArchivedOn     *uint64 `json:"archived_on"`
		BelongsTo      uint64  `json:"belongs_to"`
	}

	// PreparationList represents a list of preparations
	PreparationList struct {
		Pagination
		Preparations []Preparation `json:"preparations"`
	}

	// PreparationCreationInput represents what a user could set as input for creating preparations
	PreparationCreationInput struct {
		Name           string `json:"name"`
		Variant        string `json:"variant"`
		Description    string `json:"description"`
		AllergyWarning string `json:"allergy_warning"`
		Icon           string `json:"icon"`
		BelongsTo      uint64 `json:"-"`
	}

	// PreparationUpdateInput represents what a user could set as input for updating preparations
	PreparationUpdateInput struct {
		Name           string `json:"name"`
		Variant        string `json:"variant"`
		Description    string `json:"description"`
		AllergyWarning string `json:"allergy_warning"`
		Icon           string `json:"icon"`
		BelongsTo      uint64 `json:"-"`
	}

	// PreparationDataManager describes a structure capable of storing preparations permanently
	PreparationDataManager interface {
		GetPreparation(ctx context.Context, preparationID, userID uint64) (*Preparation, error)
		GetPreparationCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllPreparationsCount(ctx context.Context) (uint64, error)
		GetPreparations(ctx context.Context, filter *QueryFilter, userID uint64) (*PreparationList, error)
		GetAllPreparationsForUser(ctx context.Context, userID uint64) ([]Preparation, error)
		CreatePreparation(ctx context.Context, input *PreparationCreationInput) (*Preparation, error)
		UpdatePreparation(ctx context.Context, updated *Preparation) error
		ArchivePreparation(ctx context.Context, id, userID uint64) error
	}

	// PreparationDataServer describes a structure capable of serving traffic related to preparations
	PreparationDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an PreparationInput with a preparation
func (x *Preparation) Update(input *PreparationUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Variant != "" && input.Variant != x.Variant {
		x.Variant = input.Variant
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.AllergyWarning != "" && input.AllergyWarning != x.AllergyWarning {
		x.AllergyWarning = input.AllergyWarning
	}

	if input.Icon != "" && input.Icon != x.Icon {
		x.Icon = input.Icon
	}
}

// ToInput creates a PreparationUpdateInput struct for a preparation
func (x *Preparation) ToInput() *PreparationUpdateInput {
	return &PreparationUpdateInput{
		Name:           x.Name,
		Variant:        x.Variant,
		Description:    x.Description,
		AllergyWarning: x.AllergyWarning,
		Icon:           x.Icon,
	}
}
