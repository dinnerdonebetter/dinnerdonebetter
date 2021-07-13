package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ValidIngredientPreparation represents a valid ingredient preparation.
	ValidIngredientPreparation struct {
		LastUpdatedOn      *uint64 `json:"lastUpdatedOn"`
		ArchivedOn         *uint64 `json:"archivedOn"`
		Notes              string  `json:"notes"`
		ExternalID         string  `json:"externalID"`
		ValidIngredientID  uint64  `json:"validIngredientID"`
		ValidPreparationID uint64  `json:"validPreparationID"`
		CreatedOn          uint64  `json:"createdOn"`
		ID                 uint64  `json:"id"`
	}

	// ValidIngredientPreparationList represents a list of valid ingredient preparations.
	ValidIngredientPreparationList struct {
		ValidIngredientPreparations []*ValidIngredientPreparation `json:"validIngredientPreparations"`
		Pagination
	}

	// ValidIngredientPreparationCreationInput represents what a user could set as input for creating valid ingredient preparations.
	ValidIngredientPreparationCreationInput struct {
		Notes              string `json:"notes"`
		ValidIngredientID  uint64 `json:"validIngredientID"`
		ValidPreparationID uint64 `json:"validPreparationID"`
	}

	// ValidIngredientPreparationUpdateInput represents what a user could set as input for updating valid ingredient preparations.
	ValidIngredientPreparationUpdateInput struct {
		Notes              string `json:"notes"`
		ValidIngredientID  uint64 `json:"validIngredientID"`
		ValidPreparationID uint64 `json:"validPreparationID"`
	}

	// ValidIngredientPreparationDataManager describes a structure capable of storing valid ingredient preparations permanently.
	ValidIngredientPreparationDataManager interface {
		ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (bool, error)
		GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (*ValidIngredientPreparation, error)
		GetAllValidIngredientPreparationsCount(ctx context.Context) (uint64, error)
		GetAllValidIngredientPreparations(ctx context.Context, resultChannel chan []*ValidIngredientPreparation, bucketSize uint16) error
		GetValidIngredientPreparations(ctx context.Context, filter *QueryFilter) (*ValidIngredientPreparationList, error)
		GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*ValidIngredientPreparation, error)
		CreateValidIngredientPreparation(ctx context.Context, input *ValidIngredientPreparationCreationInput, createdByUser uint64) (*ValidIngredientPreparation, error)
		UpdateValidIngredientPreparation(ctx context.Context, updated *ValidIngredientPreparation, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID, archivedBy uint64) error
		GetAuditLogEntriesForValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) ([]*AuditLogEntry, error)
	}

	// ValidIngredientPreparationDataService describes a structure capable of serving traffic related to valid ingredient preparations.
	ValidIngredientPreparationDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientPreparationUpdateInput with a valid ingredient preparation.
func (x *ValidIngredientPreparation) Update(input *ValidIngredientPreparationUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Notes != "" && input.Notes != x.Notes {
		out = append(out, &FieldChangeSummary{
			FieldName: "Notes",
			OldValue:  x.Notes,
			NewValue:  input.Notes,
		})

		x.Notes = input.Notes
	}

	if input.ValidIngredientID != 0 && input.ValidIngredientID != x.ValidIngredientID {
		out = append(out, &FieldChangeSummary{
			FieldName: "ValidIngredientID",
			OldValue:  x.ValidIngredientID,
			NewValue:  input.ValidIngredientID,
		})

		x.ValidIngredientID = input.ValidIngredientID
	}

	if input.ValidPreparationID != 0 && input.ValidPreparationID != x.ValidPreparationID {
		out = append(out, &FieldChangeSummary{
			FieldName: "ValidPreparationID",
			OldValue:  x.ValidPreparationID,
			NewValue:  input.ValidPreparationID,
		})

		x.ValidPreparationID = input.ValidPreparationID
	}

	return out
}

var _ validation.ValidatableWithContext = (*ValidIngredientPreparationCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientPreparationCreationInput.
func (x *ValidIngredientPreparationCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.ValidPreparationID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientPreparationUpdateInput)(nil)

// ValidateWithContext validates a ValidIngredientPreparationUpdateInput.
func (x *ValidIngredientPreparationUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.ValidPreparationID, validation.Required),
	)
}
