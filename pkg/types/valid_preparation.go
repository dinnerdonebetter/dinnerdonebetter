package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationDataType indicates an event is related to a valid preparation.
	ValidPreparationDataType dataType = "valid_preparation"
)

func init() {
	gob.Register(new(ValidPreparation))
	gob.Register(new(ValidPreparationList))
	gob.Register(new(ValidPreparationCreationRequestInput))
	gob.Register(new(ValidPreparationUpdateRequestInput))
}

type (
	// ValidPreparation represents a valid preparation.
	ValidPreparation struct {
		_             struct{}
		ArchivedOn    *uint64 `json:"archivedOn"`
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		IconPath      string  `json:"iconPath"`
		ID            string  `json:"id"`
		CreatedOn     uint64  `json:"createdOn"`
	}

	// ValidPreparationList represents a list of valid preparations.
	ValidPreparationList struct {
		_                 struct{}
		ValidPreparations []*ValidPreparation `json:"validPreparations"`
		Pagination
	}

	// ValidPreparationCreationRequestInput represents what a user could set as input for creating valid preparations.
	ValidPreparationCreationRequestInput struct {
		_ struct{}

		ID          string `json:"-"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidPreparationDatabaseCreationInput represents what a user could set as input for creating valid preparations.
	ValidPreparationDatabaseCreationInput struct {
		_ struct{}

		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidPreparationUpdateRequestInput represents what a user could set as input for updating valid preparations.
	ValidPreparationUpdateRequestInput struct {
		_ struct{}

		Name        string `json:"name"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidPreparationDataManager describes a structure capable of storing valid preparations permanently.
	ValidPreparationDataManager interface {
		ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error)
		GetValidPreparation(ctx context.Context, validPreparationID string) (*ValidPreparation, error)
		GetTotalValidPreparationCount(ctx context.Context) (uint64, error)
		GetValidPreparations(ctx context.Context, filter *QueryFilter) (*ValidPreparationList, error)
		SearchForValidPreparations(ctx context.Context, query string) ([]*ValidPreparation, error)
		GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*ValidPreparation, error)
		CreateValidPreparation(ctx context.Context, input *ValidPreparationDatabaseCreationInput) (*ValidPreparation, error)
		UpdateValidPreparation(ctx context.Context, updated *ValidPreparation) error
		ArchiveValidPreparation(ctx context.Context, validPreparationID string) error
	}

	// ValidPreparationDataService describes a structure capable of serving traffic related to valid preparations.
	ValidPreparationDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidPreparationUpdateRequestInput with a valid preparation.
func (x *ValidPreparation) Update(input *ValidPreparationUpdateRequestInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.IconPath != "" && input.IconPath != x.IconPath {
		x.IconPath = input.IconPath
	}
}

var _ validation.ValidatableWithContext = (*ValidPreparationCreationRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationCreationRequestInput.
func (x *ValidPreparationCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidPreparationDatabaseCreationInput.
func (x *ValidPreparationDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

// ValidPreparationDatabaseCreationInputFromValidPreparationCreationInput creates a DatabaseCreationInput from a CreationInput.
func ValidPreparationDatabaseCreationInputFromValidPreparationCreationInput(input *ValidPreparationCreationRequestInput) *ValidPreparationDatabaseCreationInput {
	x := &ValidPreparationDatabaseCreationInput{
		Name:        input.Name,
		Description: input.Description,
		IconPath:    input.IconPath,
	}

	return x
}

var _ validation.ValidatableWithContext = (*ValidPreparationUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationUpdateRequestInput.
func (x *ValidPreparationUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
