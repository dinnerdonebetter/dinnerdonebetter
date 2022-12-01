package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationDataType indicates an event is related to a valid preparation.
	ValidPreparationDataType dataType = "valid_preparation"

	// ValidPreparationCreatedCustomerEventType indicates a valid preparation was created.
	ValidPreparationCreatedCustomerEventType CustomerEventType = "valid_preparation_created"
	// ValidPreparationUpdatedCustomerEventType indicates a valid preparation was updated.
	ValidPreparationUpdatedCustomerEventType CustomerEventType = "valid_preparation_updated"
	// ValidPreparationArchivedCustomerEventType indicates a valid preparation was archived.
	ValidPreparationArchivedCustomerEventType CustomerEventType = "valid_preparation_archived"
)

func init() {
	gob.Register(new(ValidPreparation))
	gob.Register(new(ValidPreparationCreationRequestInput))
	gob.Register(new(ValidPreparationUpdateRequestInput))
}

type (
	// ValidPreparation represents a valid preparation.
	ValidPreparation struct {
		_                        struct{}
		CreatedAt                time.Time  `json:"createdAt"`
		MaximumInstrumentCount   *int32     `json:"maximumInstrumentCount"`
		ArchivedAt               *time.Time `json:"archivedAt"`
		MaximumIngredientCount   *int32     `json:"maximumIngredientCount"`
		LastUpdatedAt            *time.Time `json:"lastUpdatedAt"`
		Description              string     `json:"description"`
		IconPath                 string     `json:"iconPath"`
		PastTense                string     `json:"pastTense"`
		ID                       string     `json:"id"`
		Name                     string     `json:"name"`
		Slug                     string     `json:"slug"`
		MinimumInstrumentCount   int32      `json:"minimumInstrumentCount"`
		MinimumIngredientCount   int32      `json:"minimumIngredientCount"`
		ZeroIngredientsAllowable bool       `json:"zeroIngredientsAllowable"`
		RestrictToIngredients    bool       `json:"restrictToIngredients"`
		YieldsNothing            bool       `json:"yieldsNothing"`
		TemperatureRequired      bool       `json:"temperatureRequired"`
		TimeEstimateRequired     bool       `json:"timeEstimateRequired"`
	}

	// ValidPreparationCreationRequestInput represents what a user could set as input for creating valid preparations.
	ValidPreparationCreationRequestInput struct {
		_                        struct{}
		MaximumInstrumentCount   *int32 `json:"maximumInstrumentCount"`
		MaximumIngredientCount   *int32 `json:"maximumIngredientCount"`
		Description              string `json:"description"`
		IconPath                 string `json:"iconPath"`
		PastTense                string `json:"pastTense"`
		Slug                     string `json:"slug"`
		Name                     string `json:"name"`
		MinimumIngredientCount   int32  `json:"minimumIngredientCount"`
		MinimumInstrumentCount   int32  `json:"minimumInstrumentCount"`
		RestrictToIngredients    bool   `json:"restrictToIngredients"`
		ZeroIngredientsAllowable bool   `json:"zeroIngredientsAllowable"`
		YieldsNothing            bool   `json:"yieldsNothing"`
		TemperatureRequired      bool   `json:"temperatureRequired"`
		TimeEstimateRequired     bool   `json:"timeEstimateRequired"`
	}

	// ValidPreparationDatabaseCreationInput represents what a user could set as input for creating valid preparations.
	ValidPreparationDatabaseCreationInput struct {
		_                        struct{}
		MaximumInstrumentCount   *int32
		MaximumIngredientCount   *int32
		Description              string
		IconPath                 string
		PastTense                string
		Slug                     string
		ID                       string
		Name                     string
		MinimumIngredientCount   int32
		MinimumInstrumentCount   int32
		RestrictToIngredients    bool
		ZeroIngredientsAllowable bool
		YieldsNothing            bool
		TemperatureRequired      bool
		TimeEstimateRequired     bool
	}

	// ValidPreparationUpdateRequestInput represents what a user could set as input for updating valid preparations.
	ValidPreparationUpdateRequestInput struct {
		_ struct{}

		Name                     *string `json:"name"`
		Description              *string `json:"description"`
		IconPath                 *string `json:"iconPath"`
		YieldsNothing            *bool   `json:"yieldsNothing"`
		Slug                     *string `json:"slug"`
		RestrictToIngredients    *bool   `json:"restrictToIngredients"`
		ZeroIngredientsAllowable *bool   `json:"zeroIngredientsAllowable"`
		PastTense                *string `json:"pastTense"`
		MinimumInstrumentCount   *int32  `json:"minimumInstrumentCount"`
		MaximumInstrumentCount   *int32  `json:"maximumInstrumentCount"`
		MinimumIngredientCount   *int32  `json:"minimumIngredientCount"`
		MaximumIngredientCount   *int32  `json:"maximumIngredientCount"`
		TemperatureRequired      *bool   `json:"temperatureRequired"`
		TimeEstimateRequired     *bool   `json:"timeEstimateRequired"`
	}

	// ValidPreparationDataManager describes a structure capable of storing valid preparations permanently.
	ValidPreparationDataManager interface {
		ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error)
		GetValidPreparation(ctx context.Context, validPreparationID string) (*ValidPreparation, error)
		GetRandomValidPreparation(ctx context.Context) (*ValidPreparation, error)
		GetValidPreparations(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidPreparation], error)
		SearchForValidPreparations(ctx context.Context, query string) ([]*ValidPreparation, error)
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
		RandomHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidPreparationUpdateRequestInput with a valid preparation.
func (x *ValidPreparation) Update(input *ValidPreparationUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.IconPath != nil && *input.IconPath != x.IconPath {
		x.IconPath = *input.IconPath
	}

	if input.YieldsNothing != nil && *input.YieldsNothing != x.YieldsNothing {
		x.YieldsNothing = *input.YieldsNothing
	}

	if input.RestrictToIngredients != nil && *input.RestrictToIngredients != x.RestrictToIngredients {
		x.RestrictToIngredients = *input.RestrictToIngredients
	}

	if input.ZeroIngredientsAllowable != nil && *input.ZeroIngredientsAllowable != x.ZeroIngredientsAllowable {
		x.ZeroIngredientsAllowable = *input.ZeroIngredientsAllowable
	}

	if input.MinimumIngredientCount != nil && *input.MinimumIngredientCount != x.MinimumIngredientCount {
		x.MinimumIngredientCount = *input.MinimumIngredientCount
	}
	if input.MaximumIngredientCount != nil && *input.MaximumIngredientCount != *x.MaximumIngredientCount {
		x.MaximumIngredientCount = input.MaximumIngredientCount
	}
	if input.MinimumInstrumentCount != nil && *input.MinimumInstrumentCount != x.MinimumInstrumentCount {
		x.MinimumInstrumentCount = *input.MinimumInstrumentCount
	}
	if input.MaximumInstrumentCount != nil && *input.MaximumInstrumentCount != *x.MaximumInstrumentCount {
		x.MaximumInstrumentCount = input.MaximumInstrumentCount
	}
	if input.TemperatureRequired != nil && *input.TemperatureRequired != x.TemperatureRequired {
		x.TemperatureRequired = *input.TemperatureRequired
	}
	if input.TimeEstimateRequired != nil && *input.TimeEstimateRequired != x.TimeEstimateRequired {
		x.TimeEstimateRequired = *input.TimeEstimateRequired
	}

	if input.PastTense != nil && *input.PastTense != x.PastTense {
		x.PastTense = *input.PastTense
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
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

var _ validation.ValidatableWithContext = (*ValidPreparationUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationUpdateRequestInput.
func (x *ValidPreparationUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
