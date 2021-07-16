package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ValidPreparationInstrument represents a valid preparation instrument.
	ValidPreparationInstrument struct {
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		Notes         string  `json:"notes"`
		ExternalID    string  `json:"externalID"`
		InstrumentID  uint64  `json:"instrumentID"`
		PreparationID uint64  `json:"preparationID"`
		CreatedOn     uint64  `json:"createdOn"`
		ID            uint64  `json:"id"`
	}

	// ValidPreparationInstrumentList represents a list of valid preparation instruments.
	ValidPreparationInstrumentList struct {
		ValidPreparationInstruments []*ValidPreparationInstrument `json:"validPreparationInstruments"`
		Pagination
	}

	// ValidPreparationInstrumentCreationInput represents what a user could set as input for creating valid preparation instruments.
	ValidPreparationInstrumentCreationInput struct {
		Notes         string `json:"notes"`
		InstrumentID  uint64 `json:"instrumentID"`
		PreparationID uint64 `json:"preparationID"`
	}

	// ValidPreparationInstrumentUpdateInput represents what a user could set as input for updating valid preparation instruments.
	ValidPreparationInstrumentUpdateInput struct {
		Notes         string `json:"notes"`
		InstrumentID  uint64 `json:"instrumentID"`
		PreparationID uint64 `json:"preparationID"`
	}

	// ValidPreparationInstrumentDataManager describes a structure capable of storing valid preparation instruments permanently.
	ValidPreparationInstrumentDataManager interface {
		ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID uint64) (bool, error)
		GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) (*ValidPreparationInstrument, error)
		GetAllValidPreparationInstrumentsCount(ctx context.Context) (uint64, error)
		GetAllValidPreparationInstruments(ctx context.Context, resultChannel chan []*ValidPreparationInstrument, bucketSize uint16) error
		GetValidPreparationInstruments(ctx context.Context, filter *QueryFilter) (*ValidPreparationInstrumentList, error)
		GetValidPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*ValidPreparationInstrument, error)
		CreateValidPreparationInstrument(ctx context.Context, input *ValidPreparationInstrumentCreationInput, createdByUser uint64) (*ValidPreparationInstrument, error)
		UpdateValidPreparationInstrument(ctx context.Context, updated *ValidPreparationInstrument, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID, archivedBy uint64) error
		GetAuditLogEntriesForValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) ([]*AuditLogEntry, error)
	}

	// ValidPreparationInstrumentDataService describes a structure capable of serving traffic related to valid preparation instruments.
	ValidPreparationInstrumentDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidPreparationInstrumentUpdateInput with a valid preparation instrument.
func (x *ValidPreparationInstrument) Update(input *ValidPreparationInstrumentUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.InstrumentID != 0 && input.InstrumentID != x.InstrumentID {
		out = append(out, &FieldChangeSummary{
			FieldName: "InstrumentID",
			OldValue:  x.InstrumentID,
			NewValue:  input.InstrumentID,
		})

		x.InstrumentID = input.InstrumentID
	}

	if input.PreparationID != 0 && input.PreparationID != x.PreparationID {
		out = append(out, &FieldChangeSummary{
			FieldName: "PreparationID",
			OldValue:  x.PreparationID,
			NewValue:  input.PreparationID,
		})

		x.PreparationID = input.PreparationID
	}

	if input.Notes != "" && input.Notes != x.Notes {
		out = append(out, &FieldChangeSummary{
			FieldName: "Notes",
			OldValue:  x.Notes,
			NewValue:  input.Notes,
		})

		x.Notes = input.Notes
	}

	return out
}

var _ validation.ValidatableWithContext = (*ValidPreparationInstrumentCreationInput)(nil)

// ValidateWithContext validates a ValidPreparationInstrumentCreationInput.
func (x *ValidPreparationInstrumentCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.InstrumentID, validation.Required),
		validation.Field(&x.PreparationID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationInstrumentUpdateInput)(nil)

// ValidateWithContext validates a ValidPreparationInstrumentUpdateInput.
func (x *ValidPreparationInstrumentUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.InstrumentID, validation.Required),
		validation.Field(&x.PreparationID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
