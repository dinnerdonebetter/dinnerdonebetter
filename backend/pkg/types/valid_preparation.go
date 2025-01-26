package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationCreatedServiceEventType indicates a valid preparation was created.
	ValidPreparationCreatedServiceEventType = "valid_preparation_created"
	// ValidPreparationUpdatedServiceEventType indicates a valid preparation was updated.
	ValidPreparationUpdatedServiceEventType = "valid_preparation_updated"
	// ValidPreparationArchivedServiceEventType indicates a valid preparation was archived.
	ValidPreparationArchivedServiceEventType = "valid_preparation_archived"
)

func init() {
	gob.Register(new(ValidPreparation))
	gob.Register(new(ValidPreparationCreationRequestInput))
	gob.Register(new(ValidPreparationUpdateRequestInput))
}

type (
	// ValidPreparation represents a valid preparation.
	ValidPreparation struct {
		_ struct{} `json:"-"`

		CreatedAt                   time.Time                  `json:"createdAt"`
		InstrumentCount             Uint16RangeWithOptionalMax `json:"instrumentCount"`
		IngredientCount             Uint16RangeWithOptionalMax `json:"ingredientCount"`
		VesselCount                 Uint16RangeWithOptionalMax `json:"vesselCount"`
		ArchivedAt                  *time.Time                 `json:"archivedAt"`
		LastUpdatedAt               *time.Time                 `json:"lastUpdatedAt"`
		IconPath                    string                     `json:"iconPath"`
		PastTense                   string                     `json:"pastTense"`
		ID                          string                     `json:"id"`
		Name                        string                     `json:"name"`
		Description                 string                     `json:"description"`
		Slug                        string                     `json:"slug"`
		RestrictToIngredients       bool                       `json:"restrictToIngredients"`
		TemperatureRequired         bool                       `json:"temperatureRequired"`
		TimeEstimateRequired        bool                       `json:"timeEstimateRequired"`
		ConditionExpressionRequired bool                       `json:"conditionExpressionRequired"`
		ConsumesVessel              bool                       `json:"consumesVessel"`
		OnlyForVessels              bool                       `json:"onlyForVessels"`
		YieldsNothing               bool                       `json:"yieldsNothing"`
	}

	// ValidPreparationCreationRequestInput represents what a user could set as input for creating valid preparations.
	ValidPreparationCreationRequestInput struct {
		_ struct{} `json:"-"`

		InstrumentCount             Uint16RangeWithOptionalMax `json:"instrumentCount"`
		IngredientCount             Uint16RangeWithOptionalMax `json:"ingredientCount"`
		VesselCount                 Uint16RangeWithOptionalMax `json:"vesselCount"`
		IconPath                    string                     `json:"iconPath"`
		PastTense                   string                     `json:"pastTense"`
		Slug                        string                     `json:"slug"`
		Name                        string                     `json:"name"`
		Description                 string                     `json:"description"`
		TemperatureRequired         bool                       `json:"temperatureRequired"`
		TimeEstimateRequired        bool                       `json:"timeEstimateRequired"`
		ConditionExpressionRequired bool                       `json:"conditionExpressionRequired"`
		ConsumesVessel              bool                       `json:"consumesVessel"`
		OnlyForVessels              bool                       `json:"onlyForVessels"`
		RestrictToIngredients       bool                       `json:"restrictToIngredients"`
		YieldsNothing               bool                       `json:"yieldsNothing"`
	}

	// ValidPreparationDatabaseCreationInput represents what a user could set as input for creating valid preparations.
	ValidPreparationDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		InstrumentCount             Uint16RangeWithOptionalMax `json:"-"`
		IngredientCount             Uint16RangeWithOptionalMax `json:"-"`
		VesselCount                 Uint16RangeWithOptionalMax `json:"-"`
		IconPath                    string                     `json:"-"`
		PastTense                   string                     `json:"-"`
		Slug                        string                     `json:"-"`
		ID                          string                     `json:"-"`
		Name                        string                     `json:"-"`
		Description                 string                     `json:"-"`
		TemperatureRequired         bool                       `json:"-"`
		TimeEstimateRequired        bool                       `json:"-"`
		ConditionExpressionRequired bool                       `json:"-"`
		ConsumesVessel              bool                       `json:"-"`
		OnlyForVessels              bool                       `json:"-"`
		RestrictToIngredients       bool                       `json:"-"`
		YieldsNothing               bool                       `json:"-"`
	}

	// ValidPreparationUpdateRequestInput represents what a user could set as input for updating valid preparations.
	ValidPreparationUpdateRequestInput struct {
		_ struct{} `json:"-"`

		InstrumentCount             Uint16RangeWithOptionalMaxUpdateRequestInput `json:"instrumentCount"`
		IngredientCount             Uint16RangeWithOptionalMaxUpdateRequestInput `json:"ingredientCount"`
		VesselCount                 Uint16RangeWithOptionalMaxUpdateRequestInput `json:"vesselCount"`
		Name                        *string                                      `json:"name,omitempty"`
		Description                 *string                                      `json:"description,omitempty"`
		IconPath                    *string                                      `json:"iconPath,omitempty"`
		YieldsNothing               *bool                                        `json:"yieldsNothing,omitempty"`
		Slug                        *string                                      `json:"slug,omitempty"`
		RestrictToIngredients       *bool                                        `json:"restrictToIngredients,omitempty"`
		PastTense                   *string                                      `json:"pastTense,omitempty"`
		TemperatureRequired         *bool                                        `json:"temperatureRequired,omitempty"`
		TimeEstimateRequired        *bool                                        `json:"timeEstimateRequired,omitempty"`
		ConditionExpressionRequired *bool                                        `json:"conditionExpressionRequired,omitempty"`
		ConsumesVessel              *bool                                        `json:"consumesVessel,omitempty"`
		OnlyForVessels              *bool                                        `json:"onlyForVessels,omitempty"`
	}

	// ValidPreparationDataManager describes a structure capable of storing valid preparations permanently.
	ValidPreparationDataManager interface {
		ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error)
		GetValidPreparation(ctx context.Context, validPreparationID string) (*ValidPreparation, error)
		GetRandomValidPreparation(ctx context.Context) (*ValidPreparation, error)
		GetValidPreparations(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[ValidPreparation], error)
		SearchForValidPreparations(ctx context.Context, query string) ([]*ValidPreparation, error)
		CreateValidPreparation(ctx context.Context, input *ValidPreparationDatabaseCreationInput) (*ValidPreparation, error)
		UpdateValidPreparation(ctx context.Context, updated *ValidPreparation) error
		MarkValidPreparationAsIndexed(ctx context.Context, validPreparationID string) error
		ArchiveValidPreparation(ctx context.Context, validPreparationID string) error
		GetValidPreparationIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetValidPreparationsWithIDs(ctx context.Context, ids []string) ([]*ValidPreparation, error)
	}

	// ValidPreparationDataService describes a structure capable of serving traffic related to valid preparations.
	ValidPreparationDataService interface {
		SearchValidPreparationsHandler(http.ResponseWriter, *http.Request)
		ListValidPreparationsHandler(http.ResponseWriter, *http.Request)
		CreateValidPreparationHandler(http.ResponseWriter, *http.Request)
		ReadValidPreparationHandler(http.ResponseWriter, *http.Request)
		RandomValidPreparationHandler(http.ResponseWriter, *http.Request)
		UpdateValidPreparationHandler(http.ResponseWriter, *http.Request)
		ArchiveValidPreparationHandler(http.ResponseWriter, *http.Request)
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

	if input.IngredientCount.Min != nil && *input.IngredientCount.Min != x.IngredientCount.Min {
		x.IngredientCount.Min = *input.IngredientCount.Min
	}

	if input.IngredientCount.Max != nil && x.IngredientCount.Max != nil && *input.IngredientCount.Max != *x.IngredientCount.Max {
		x.IngredientCount.Max = input.IngredientCount.Max
	}

	if input.InstrumentCount.Min != nil && *input.InstrumentCount.Min != x.InstrumentCount.Min {
		x.InstrumentCount.Min = *input.InstrumentCount.Min
	}

	if input.InstrumentCount.Max != nil && x.InstrumentCount.Max != nil && *input.InstrumentCount.Max != *x.InstrumentCount.Max {
		x.InstrumentCount.Max = input.InstrumentCount.Max
	}

	if input.TemperatureRequired != nil && *input.TemperatureRequired != x.TemperatureRequired {
		x.TemperatureRequired = *input.TemperatureRequired
	}

	if input.TimeEstimateRequired != nil && *input.TimeEstimateRequired != x.TimeEstimateRequired {
		x.TimeEstimateRequired = *input.TimeEstimateRequired
	}

	if input.ConditionExpressionRequired != nil && *input.ConditionExpressionRequired != x.ConditionExpressionRequired {
		x.ConditionExpressionRequired = *input.ConditionExpressionRequired
	}

	if input.PastTense != nil && *input.PastTense != x.PastTense {
		x.PastTense = *input.PastTense
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.ConsumesVessel != nil && *input.ConsumesVessel != x.ConsumesVessel {
		x.ConsumesVessel = *input.ConsumesVessel
	}

	if input.OnlyForVessels != nil && *input.OnlyForVessels != x.OnlyForVessels {
		x.OnlyForVessels = *input.OnlyForVessels
	}

	if input.VesselCount.Min != nil && *input.VesselCount.Min != x.VesselCount.Min {
		x.VesselCount.Min = *input.VesselCount.Min
	}

	if input.VesselCount.Max != nil && x.VesselCount.Max != nil && *input.VesselCount.Max != *x.VesselCount.Max {
		x.VesselCount.Max = input.VesselCount.Max
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
