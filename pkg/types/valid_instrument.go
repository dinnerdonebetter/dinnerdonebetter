package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidInstrumentDataType indicates an event is related to a valid instrument.
	ValidInstrumentDataType dataType = "valid_instrument"
)

func init() {
	gob.Register(new(ValidInstrument))
	gob.Register(new(ValidInstrumentList))
	gob.Register(new(ValidInstrumentCreationRequestInput))
	gob.Register(new(ValidInstrumentUpdateRequestInput))
}

type (
	// ValidInstrument represents a valid instrument.
	ValidInstrument struct {
		_             struct{}
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		Variant       string  `json:"variant"`
		Description   string  `json:"description"`
		IconPath      string  `json:"iconPath"`
		ID            string  `json:"id"`
		Name          string  `json:"name"`
		CreatedOn     uint64  `json:"createdOn"`
	}

	// ValidInstrumentList represents a list of valid instruments.
	ValidInstrumentList struct {
		_                struct{}
		ValidInstruments []*ValidInstrument `json:"validInstruments"`
		Pagination
	}

	// ValidInstrumentCreationRequestInput represents what a user could set as input for creating valid instruments.
	ValidInstrumentCreationRequestInput struct {
		_ struct{}

		ID          string `json:"-"`
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidInstrumentDatabaseCreationInput represents what a user could set as input for creating valid instruments.
	ValidInstrumentDatabaseCreationInput struct {
		_ struct{}

		ID          string `json:"id"`
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidInstrumentUpdateRequestInput represents what a user could set as input for updating valid instruments.
	ValidInstrumentUpdateRequestInput struct {
		_ struct{}

		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
	}

	// ValidInstrumentDataManager describes a structure capable of storing valid instruments permanently.
	ValidInstrumentDataManager interface {
		ValidInstrumentExists(ctx context.Context, validInstrumentID string) (bool, error)
		GetValidInstrument(ctx context.Context, validInstrumentID string) (*ValidInstrument, error)
		GetTotalValidInstrumentCount(ctx context.Context) (uint64, error)
		GetValidInstruments(ctx context.Context, filter *QueryFilter) (*ValidInstrumentList, error)
		SearchForValidInstruments(ctx context.Context, query string) ([]*ValidInstrument, error)
		GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*ValidInstrument, error)
		CreateValidInstrument(ctx context.Context, input *ValidInstrumentDatabaseCreationInput) (*ValidInstrument, error)
		UpdateValidInstrument(ctx context.Context, updated *ValidInstrument) error
		ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error
	}

	// ValidInstrumentDataService describes a structure capable of serving traffic related to valid instruments.
	ValidInstrumentDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidInstrumentUpdateRequestInput with a valid instrument.
func (x *ValidInstrument) Update(input *ValidInstrumentUpdateRequestInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Variant != "" && input.Variant != x.Variant {
		x.Variant = input.Variant
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.IconPath != "" && input.IconPath != x.IconPath {
		x.IconPath = input.IconPath
	}
}

var _ validation.ValidatableWithContext = (*ValidInstrumentCreationRequestInput)(nil)

// ValidateWithContext validates a ValidInstrumentCreationRequestInput.
func (x *ValidInstrumentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidInstrumentDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidInstrumentDatabaseCreationInput.
func (x *ValidInstrumentDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

// ValidInstrumentDatabaseCreationInputFromValidInstrumentCreationInput creates a DatabaseCreationInput from a CreationInput.
func ValidInstrumentDatabaseCreationInputFromValidInstrumentCreationInput(input *ValidInstrumentCreationRequestInput) *ValidInstrumentDatabaseCreationInput {
	x := &ValidInstrumentDatabaseCreationInput{
		Name:        input.Name,
		Variant:     input.Variant,
		Description: input.Description,
		IconPath:    input.IconPath,
	}

	return x
}

var _ validation.ValidatableWithContext = (*ValidInstrumentUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidInstrumentUpdateRequestInput.
func (x *ValidInstrumentUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
