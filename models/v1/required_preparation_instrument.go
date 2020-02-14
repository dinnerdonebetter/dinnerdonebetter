package models

import (
	"context"
	"net/http"
)

type (
	// RequiredPreparationInstrument represents a required preparation instrument
	RequiredPreparationInstrument struct {
		ID            uint64  `json:"id"`
		InstrumentID  uint64  `json:"instrument_id"`
		PreparationID uint64  `json:"preparation_id"`
		Notes         string  `json:"notes"`
		CreatedOn     uint64  `json:"created_on"`
		UpdatedOn     *uint64 `json:"updated_on"`
		ArchivedOn    *uint64 `json:"archived_on"`
	}

	// RequiredPreparationInstrumentList represents a list of required preparation instruments
	RequiredPreparationInstrumentList struct {
		Pagination
		RequiredPreparationInstruments []RequiredPreparationInstrument `json:"required_preparation_instruments"`
	}

	// RequiredPreparationInstrumentCreationInput represents what a user could set as input for creating required preparation instruments
	RequiredPreparationInstrumentCreationInput struct {
		InstrumentID  uint64 `json:"instrument_id"`
		PreparationID uint64 `json:"preparation_id"`
		Notes         string `json:"notes"`
	}

	// RequiredPreparationInstrumentUpdateInput represents what a user could set as input for updating required preparation instruments
	RequiredPreparationInstrumentUpdateInput struct {
		InstrumentID  uint64 `json:"instrument_id"`
		PreparationID uint64 `json:"preparation_id"`
		Notes         string `json:"notes"`
	}

	// RequiredPreparationInstrumentDataManager describes a structure capable of storing required preparation instruments permanently
	RequiredPreparationInstrumentDataManager interface {
		GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) (*RequiredPreparationInstrument, error)
		GetRequiredPreparationInstrumentCount(ctx context.Context, filter *QueryFilter) (uint64, error)
		GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (uint64, error)
		GetRequiredPreparationInstruments(ctx context.Context, filter *QueryFilter) (*RequiredPreparationInstrumentList, error)
		CreateRequiredPreparationInstrument(ctx context.Context, input *RequiredPreparationInstrumentCreationInput) (*RequiredPreparationInstrument, error)
		UpdateRequiredPreparationInstrument(ctx context.Context, updated *RequiredPreparationInstrument) error
		ArchiveRequiredPreparationInstrument(ctx context.Context, id uint64) error
	}

	// RequiredPreparationInstrumentDataServer describes a structure capable of serving traffic related to required preparation instruments
	RequiredPreparationInstrumentDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RequiredPreparationInstrumentInput with a required preparation instrument
func (x *RequiredPreparationInstrument) Update(input *RequiredPreparationInstrumentUpdateInput) {
	if input.InstrumentID != x.InstrumentID {
		x.InstrumentID = input.InstrumentID
	}

	if input.PreparationID != x.PreparationID {
		x.PreparationID = input.PreparationID
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

// ToInput creates a RequiredPreparationInstrumentUpdateInput struct for a required preparation instrument
func (x *RequiredPreparationInstrument) ToInput() *RequiredPreparationInstrumentUpdateInput {
	return &RequiredPreparationInstrumentUpdateInput{
		InstrumentID:  x.InstrumentID,
		PreparationID: x.PreparationID,
		Notes:         x.Notes,
	}
}
