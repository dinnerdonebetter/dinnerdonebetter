package models

import (
	"context"
	"net/http"
)

type (
	// RequiredPreparationInstrument represents a required preparation instrument.
	RequiredPreparationInstrument struct {
		ID            uint64  `json:"id"`
		InstrumentID  uint64  `json:"instrumentID"`
		PreparationID uint64  `json:"preparationID"`
		Notes         string  `json:"notes"`
		CreatedOn     uint64  `json:"createdOn"`
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
	}

	// RequiredPreparationInstrumentList represents a list of required preparation instruments.
	RequiredPreparationInstrumentList struct {
		Pagination
		RequiredPreparationInstruments []RequiredPreparationInstrument `json:"requiredPreparationInstruments"`
	}

	// RequiredPreparationInstrumentCreationInput represents what a user could set as input for creating required preparation instruments.
	RequiredPreparationInstrumentCreationInput struct {
		InstrumentID  uint64 `json:"instrumentID"`
		PreparationID uint64 `json:"preparationID"`
		Notes         string `json:"notes"`
	}

	// RequiredPreparationInstrumentUpdateInput represents what a user could set as input for updating required preparation instruments.
	RequiredPreparationInstrumentUpdateInput struct {
		InstrumentID  uint64 `json:"instrumentID"`
		PreparationID uint64 `json:"preparationID"`
		Notes         string `json:"notes"`
	}

	// RequiredPreparationInstrumentDataManager describes a structure capable of storing required preparation instruments permanently.
	RequiredPreparationInstrumentDataManager interface {
		RequiredPreparationInstrumentExists(ctx context.Context, requiredPreparationInstrumentID uint64) (bool, error)
		GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) (*RequiredPreparationInstrument, error)
		GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (uint64, error)
		GetAllRequiredPreparationInstruments(ctx context.Context, resultChannel chan []RequiredPreparationInstrument) error
		GetRequiredPreparationInstruments(ctx context.Context, filter *QueryFilter) (*RequiredPreparationInstrumentList, error)
		GetRequiredPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]RequiredPreparationInstrument, error)
		CreateRequiredPreparationInstrument(ctx context.Context, input *RequiredPreparationInstrumentCreationInput) (*RequiredPreparationInstrument, error)
		UpdateRequiredPreparationInstrument(ctx context.Context, updated *RequiredPreparationInstrument) error
		ArchiveRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) error
	}

	// RequiredPreparationInstrumentDataServer describes a structure capable of serving traffic related to required preparation instruments.
	RequiredPreparationInstrumentDataServer interface {
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

// Update merges an RequiredPreparationInstrumentInput with a required preparation instrument.
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

// ToUpdateInput creates a RequiredPreparationInstrumentUpdateInput struct for a required preparation instrument.
func (x *RequiredPreparationInstrument) ToUpdateInput() *RequiredPreparationInstrumentUpdateInput {
	return &RequiredPreparationInstrumentUpdateInput{
		InstrumentID:  x.InstrumentID,
		PreparationID: x.PreparationID,
		Notes:         x.Notes,
	}
}
