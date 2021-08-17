package types

import (
	"context"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/search"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationsSearchIndexName is the name of the index used to search through valid preparations.
	ValidPreparationsSearchIndexName search.IndexName = "valid_preparations"
)

type (
	// ValidPreparation represents a valid preparation.
	ValidPreparation struct {
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		IconPath      string  `json:"iconPath"`
		ExternalID    string  `json:"externalID"`
		CreatedOn     uint64  `json:"createdOn"`
		ID            uint64  `json:"id"`
	}

	// ValidPreparationList represents a list of valid preparations.
	ValidPreparationList struct {
		ValidPreparations []*ValidPreparation `json:"validPreparations"`
		Pagination
	}

	// ValidPreparationCreationInput represents what a user could set as input for creating valid preparations.
	ValidPreparationCreationInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidPreparationUpdateInput represents what a user could set as input for updating valid preparations.
	ValidPreparationUpdateInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidPreparationDataManager describes a structure capable of storing valid preparations permanently.
	ValidPreparationDataManager interface {
		ValidPreparationExists(ctx context.Context, validPreparationID uint64) (bool, error)
		GetValidPreparation(ctx context.Context, validPreparationID uint64) (*ValidPreparation, error)
		GetAllValidPreparationsCount(ctx context.Context) (uint64, error)
		GetAllValidPreparations(ctx context.Context, resultChannel chan []*ValidPreparation, bucketSize uint16) error
		GetValidPreparations(ctx context.Context, filter *QueryFilter) (*ValidPreparationList, error)
		GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*ValidPreparation, error)
		CreateValidPreparation(ctx context.Context, input *ValidPreparationCreationInput, createdByUser uint64) (*ValidPreparation, error)
		UpdateValidPreparation(ctx context.Context, updated *ValidPreparation, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveValidPreparation(ctx context.Context, validPreparationID, archivedBy uint64) error
		GetAuditLogEntriesForValidPreparation(ctx context.Context, validPreparationID uint64) ([]*AuditLogEntry, error)
	}

	// ValidPreparationDataService describes a structure capable of serving traffic related to valid preparations.
	ValidPreparationDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		SearchForValidPreparations(ctx context.Context, sessionCtxData *SessionContextData, query string, filter *QueryFilter) ([]*ValidPreparation, error)
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidPreparationUpdateInput with a valid preparation.
func (x *ValidPreparation) Update(input *ValidPreparationUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Name != x.Name {
		out = append(out, &FieldChangeSummary{
			FieldName: "Name",
			OldValue:  x.Name,
			NewValue:  input.Name,
		})

		x.Name = input.Name
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

var _ validation.ValidatableWithContext = (*ValidPreparationCreationInput)(nil)

// ValidateWithContext validates a ValidPreparationCreationInput.
func (x *ValidPreparationCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationUpdateInput)(nil)

// ValidateWithContext validates a ValidPreparationUpdateInput.
func (x *ValidPreparationUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
