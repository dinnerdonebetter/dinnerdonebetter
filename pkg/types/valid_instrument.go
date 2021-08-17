package types

import (
	"context"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/search"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidInstrumentsSearchIndexName is the name of the index used to search through valid instruments.
	ValidInstrumentsSearchIndexName search.IndexName = "valid_instruments"
)

type (
	// ValidInstrument represents a valid instrument.
	ValidInstrument struct {
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		Variant       string  `json:"variant"`
		Description   string  `json:"description"`
		IconPath      string  `json:"iconPath"`
		ExternalID    string  `json:"externalID"`
		Name          string  `json:"name"`
		ID            uint64  `json:"id"`
		CreatedOn     uint64  `json:"createdOn"`
	}

	// ValidInstrumentList represents a list of valid instruments.
	ValidInstrumentList struct {
		ValidInstruments []*ValidInstrument `json:"validInstruments"`
		Pagination
	}

	// ValidInstrumentCreationInput represents what a user could set as input for creating valid instruments.
	ValidInstrumentCreationInput struct {
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidInstrumentUpdateInput represents what a user could set as input for updating valid instruments.
	ValidInstrumentUpdateInput struct {
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidInstrumentDataManager describes a structure capable of storing valid instruments permanently.
	ValidInstrumentDataManager interface {
		ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (bool, error)
		GetValidInstrument(ctx context.Context, validInstrumentID uint64) (*ValidInstrument, error)
		GetAllValidInstrumentsCount(ctx context.Context) (uint64, error)
		GetAllValidInstruments(ctx context.Context, resultChannel chan []*ValidInstrument, bucketSize uint16) error
		GetValidInstruments(ctx context.Context, filter *QueryFilter) (*ValidInstrumentList, error)
		GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*ValidInstrument, error)
		CreateValidInstrument(ctx context.Context, input *ValidInstrumentCreationInput, createdByUser uint64) (*ValidInstrument, error)
		UpdateValidInstrument(ctx context.Context, updated *ValidInstrument, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveValidInstrument(ctx context.Context, validInstrumentID, archivedBy uint64) error
		GetAuditLogEntriesForValidInstrument(ctx context.Context, validInstrumentID uint64) ([]*AuditLogEntry, error)
	}

	// ValidInstrumentDataService describes a structure capable of serving traffic related to valid instruments.
	ValidInstrumentDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		SearchForValidInstruments(ctx context.Context, sessionCtxData *SessionContextData, query string, filter *QueryFilter) ([]*ValidInstrument, error)
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidInstrumentUpdateInput with a valid instrument.
func (x *ValidInstrument) Update(input *ValidInstrumentUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Name != x.Name {
		out = append(out, &FieldChangeSummary{
			FieldName: "Name",
			OldValue:  x.Name,
			NewValue:  input.Name,
		})

		x.Name = input.Name
	}

	if input.Variant != x.Variant {
		out = append(out, &FieldChangeSummary{
			FieldName: "Variant",
			OldValue:  x.Variant,
			NewValue:  input.Variant,
		})

		x.Variant = input.Variant
	}

	if input.Description != x.Description {
		out = append(out, &FieldChangeSummary{
			FieldName: "Description",
			OldValue:  x.Description,
			NewValue:  input.Description,
		})

		x.Description = input.Description
	}

	if input.IconPath != x.IconPath {
		out = append(out, &FieldChangeSummary{
			FieldName: "IconPath",
			OldValue:  x.IconPath,
			NewValue:  input.IconPath,
		})

		x.IconPath = input.IconPath
	}

	return out
}

var _ validation.ValidatableWithContext = (*ValidInstrumentCreationInput)(nil)

// ValidateWithContext validates a ValidInstrumentCreationInput.
func (x *ValidInstrumentCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidInstrumentUpdateInput)(nil)

// ValidateWithContext validates a ValidInstrumentUpdateInput.
func (x *ValidInstrumentUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
