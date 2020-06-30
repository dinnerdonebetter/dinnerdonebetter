package models

import (
	"context"
	"net/http"
)

type (
	// ValidInstrument represents a valid instrument.
	ValidInstrument struct {
		ID          uint64  `json:"id"`
		Name        string  `json:"name"`
		Variant     string  `json:"variant"`
		Description string  `json:"description"`
		Icon        string  `json:"icon"`
		CreatedOn   uint64  `json:"createdOn"`
		UpdatedOn   *uint64 `json:"updatedOn"`
		ArchivedOn  *uint64 `json:"archivedOn"`
	}

	// ValidInstrumentList represents a list of valid instruments.
	ValidInstrumentList struct {
		Pagination
		ValidInstruments []ValidInstrument `json:"validInstruments"`
	}

	// ValidInstrumentCreationInput represents what a user could set as input for creating valid instruments.
	ValidInstrumentCreationInput struct {
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	// ValidInstrumentUpdateInput represents what a user could set as input for updating valid instruments.
	ValidInstrumentUpdateInput struct {
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	// ValidInstrumentDataManager describes a structure capable of storing valid instruments permanently.
	ValidInstrumentDataManager interface {
		ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (bool, error)
		GetValidInstrument(ctx context.Context, validInstrumentID uint64) (*ValidInstrument, error)
		GetAllValidInstrumentsCount(ctx context.Context) (uint64, error)
		GetValidInstruments(ctx context.Context, filter *QueryFilter) (*ValidInstrumentList, error)
		CreateValidInstrument(ctx context.Context, input *ValidInstrumentCreationInput) (*ValidInstrument, error)
		UpdateValidInstrument(ctx context.Context, updated *ValidInstrument) error
		ArchiveValidInstrument(ctx context.Context, validInstrumentID uint64) error
	}

	// ValidInstrumentDataServer describes a structure capable of serving traffic related to valid instruments.
	ValidInstrumentDataServer interface {
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

// Update merges an ValidInstrumentInput with a valid instrument.
func (x *ValidInstrument) Update(input *ValidInstrumentUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Variant != "" && input.Variant != x.Variant {
		x.Variant = input.Variant
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.Icon != "" && input.Icon != x.Icon {
		x.Icon = input.Icon
	}
}

// ToUpdateInput creates a ValidInstrumentUpdateInput struct for a valid instrument.
func (x *ValidInstrument) ToUpdateInput() *ValidInstrumentUpdateInput {
	return &ValidInstrumentUpdateInput{
		Name:        x.Name,
		Variant:     x.Variant,
		Description: x.Description,
		Icon:        x.Icon,
	}
}
