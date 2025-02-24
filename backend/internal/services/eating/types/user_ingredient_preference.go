package types

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	minRating int8 = -10
	maxRating int8 = 10
)

type (
	// IngredientPreference represents an ingredient preference.
	IngredientPreference struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time       `json:"createdAt"`
		LastUpdatedAt *time.Time      `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time      `json:"archivedAt"`
		ID            string          `json:"id"`
		Notes         string          `json:"notes"`
		BelongsToUser string          `json:"belongsToUser"`
		Ingredient    ValidIngredient `json:"ingredient"`
		Rating        int8            `json:"rating"`
		Allergy       bool            `json:"allergy"`
	}

	// IngredientPreferenceCreationRequestInput represents what a user could set as input for creating user ingredient preferences.
	IngredientPreferenceCreationRequestInput struct {
		_ struct{} `json:"-"`

		ValidIngredientGroupID string `json:"validIngredientGroupID"`
		ValidIngredientID      string `json:"validIngredientID"`
		Notes                  string `json:"notes"`
		Rating                 int8   `json:"rating"`
		Allergy                bool   `json:"allergy"`
	}

	// IngredientPreferenceDatabaseCreationInput represents what a user could set as input for creating user ingredient preferences.
	IngredientPreferenceDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ValidIngredientGroupID string `json:"-"`
		ValidIngredientID      string `json:"-"`
		Notes                  string `json:"-"`
		BelongsToUser          string `json:"-"`
		Rating                 int8   `json:"-"`
		Allergy                bool   `json:"-"`
	}

	// IngredientPreferenceUpdateRequestInput represents what a user could set as input for updating user ingredient preferences.
	IngredientPreferenceUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes        *string `json:"notes,omitempty"`
		IngredientID *string `json:"ingredientID"`
		Rating       *int8   `json:"rating"`
		Allergy      *bool   `json:"allergy"`
	}

	// IngredientPreferenceDataManager describes a structure capable of storing user ingredient preferences permanently.
	IngredientPreferenceDataManager interface {
		IngredientPreferenceExists(ctx context.Context, userIngredientPreferenceID, userID string) (bool, error)
		GetIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) (*IngredientPreference, error)
		GetIngredientPreferences(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[IngredientPreference], error)
		CreateIngredientPreference(ctx context.Context, input *IngredientPreferenceDatabaseCreationInput) ([]*IngredientPreference, error)
		UpdateIngredientPreference(ctx context.Context, updated *IngredientPreference) error
		ArchiveIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) error
	}

	// IngredientPreferenceDataService describes a structure capable of serving traffic related to user ingredient preferences.
	IngredientPreferenceDataService interface {
		ListIngredientPreferencesHandler(http.ResponseWriter, *http.Request)
		CreateIngredientPreferenceHandler(http.ResponseWriter, *http.Request)
		UpdateIngredientPreferenceHandler(http.ResponseWriter, *http.Request)
		ArchiveIngredientPreferenceHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an IngredientPreferenceUpdateRequestInput with a ingredient preference.
func (x *IngredientPreference) Update(input *IngredientPreferenceUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.IngredientID != nil && *input.IngredientID != x.Ingredient.ID {
		x.Ingredient = ValidIngredient{ID: *input.IngredientID}
	}

	if input.Rating != nil && *input.Rating != x.Rating {
		x.Rating = *input.Rating
	}

	if input.Allergy != nil && *input.Allergy != x.Allergy {
		x.Allergy = *input.Allergy
	}
}

var _ validation.ValidatableWithContext = (*IngredientPreferenceCreationRequestInput)(nil)

// ValidateWithContext validates a IngredientPreferenceCreationRequestInput.
func (x *IngredientPreferenceCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidIngredientGroupID, validation.When(x.ValidIngredientID == "", validation.Required)),
		validation.Field(&x.ValidIngredientID, validation.When(x.ValidIngredientGroupID == "", validation.Required)),
		validation.Field(&x.Rating, validation.Min(minRating), validation.Max(maxRating)),
	)
}

var _ validation.ValidatableWithContext = (*IngredientPreferenceDatabaseCreationInput)(nil)

// ValidateWithContext validates a IngredientPreferenceDatabaseCreationInput.
func (x *IngredientPreferenceDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.Rating, validation.Min(minRating), validation.Max(maxRating)),
		validation.Field(&x.BelongsToUser, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*IngredientPreferenceUpdateRequestInput)(nil)

// ValidateWithContext validates a IngredientPreferenceUpdateRequestInput.
func (x *IngredientPreferenceUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientID, validation.Required),
		validation.Field(&x.Rating, validation.Min(minRating), validation.Max(maxRating)),
	)
}
