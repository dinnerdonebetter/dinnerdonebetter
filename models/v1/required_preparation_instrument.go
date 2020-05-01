package models

import (
	"context"
	"net/http"
)

type (
	// RequiredPreparationInstrument represents a required preparation instrument.
	RequiredPreparationInstrument struct {
		ID                        uint64  `json:"id"`
		ValidInstrumentID         uint64  `json:"valid_instrument_id"`
		Notes                     string  `json:"notes"`
		CreatedOn                 uint64  `json:"created_on"`
		UpdatedOn                 *uint64 `json:"updated_on"`
		ArchivedOn                *uint64 `json:"archived_on"`
		BelongsToValidPreparation uint64  `json:"belongs_to_valid_preparation"`
	}

	// RequiredPreparationInstrumentList represents a list of required preparation instruments.
	RequiredPreparationInstrumentList struct {
		Pagination
		RequiredPreparationInstruments []RequiredPreparationInstrument `json:"required_preparation_instruments"`
	}

	// RequiredPreparationInstrumentCreationInput represents what a user could set as input for creating required preparation instruments.
	RequiredPreparationInstrumentCreationInput struct {
		ValidInstrumentID         uint64 `json:"valid_instrument_id"`
		Notes                     string `json:"notes"`
		BelongsToValidPreparation uint64 `json:"-"`
	}

	// RequiredPreparationInstrumentUpdateInput represents what a user could set as input for updating required preparation instruments.
	RequiredPreparationInstrumentUpdateInput struct {
		ValidInstrumentID         uint64 `json:"valid_instrument_id"`
		Notes                     string `json:"notes"`
		BelongsToValidPreparation uint64 `json:"belongs_to_valid_preparation"`
	}

	// RequiredPreparationInstrumentDataManager describes a structure capable of storing required preparation instruments permanently.
	RequiredPreparationInstrumentDataManager interface {
		RequiredPreparationInstrumentExists(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (bool, error)
		GetRequiredPreparationInstrument(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (*RequiredPreparationInstrument, error)
		GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (uint64, error)
		GetRequiredPreparationInstruments(ctx context.Context, validPreparationID uint64, filter *QueryFilter) (*RequiredPreparationInstrumentList, error)
		CreateRequiredPreparationInstrument(ctx context.Context, input *RequiredPreparationInstrumentCreationInput) (*RequiredPreparationInstrument, error)
		UpdateRequiredPreparationInstrument(ctx context.Context, updated *RequiredPreparationInstrument) error
		ArchiveRequiredPreparationInstrument(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) error
	}

	// RequiredPreparationInstrumentDataServer describes a structure capable of serving traffic related to required preparation instruments.
	RequiredPreparationInstrumentDataServer interface {
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

// Update merges an RequiredPreparationInstrumentInput with a required preparation instrument.
func (x *RequiredPreparationInstrument) Update(input *RequiredPreparationInstrumentUpdateInput) {
	if input.ValidInstrumentID != x.ValidInstrumentID {
		x.ValidInstrumentID = input.ValidInstrumentID
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

// ToUpdateInput creates a RequiredPreparationInstrumentUpdateInput struct for a required preparation instrument.
func (x *RequiredPreparationInstrument) ToUpdateInput() *RequiredPreparationInstrumentUpdateInput {
	return &RequiredPreparationInstrumentUpdateInput{
		ValidInstrumentID: x.ValidInstrumentID,
		Notes:             x.Notes,
	}
}
