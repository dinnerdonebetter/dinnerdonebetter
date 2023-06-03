package types

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// MealRatingCreatedCustomerEventType indicates a meal rating was created.
	MealRatingCreatedCustomerEventType CustomerEventType = "meal_rating_created"
	// MealRatingUpdatedCustomerEventType indicates a meal rating was updated.
	MealRatingUpdatedCustomerEventType CustomerEventType = "meal_rating_updated"
	// MealRatingArchivedCustomerEventType indicates a meal rating was archived.
	MealRatingArchivedCustomerEventType CustomerEventType = "meal_rating_archived"
)

func init() {
	gob.Register(new(MealRating))
	gob.Register(new(MealRatingCreationRequestInput))
	gob.Register(new(MealRatingUpdateRequestInput))
}

type (
	// MealRating represents a meal rating.
	MealRating struct {
		_ struct{}

		CreatedAt     time.Time  `json:"createdAt"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		Notes         string     `json:"notes"`
		ID            string     `json:"id"`
		MealID        string     `json:"mealID"`
		ByUser        string     `json:"byUser"`
		Taste         float32    `json:"taste"`
		Instructions  float32    `json:"instructions"`
		Overall       float32    `json:"overall"`
		Cleanup       float32    `json:"cleanup"`
		Difficulty    float32    `json:"difficulty"`
	}

	// MealRatingCreationRequestInput represents what a user could set as input for creating meal ratings.
	MealRatingCreationRequestInput struct {
		_ struct{}

		MealID       string  `json:"mealID"`
		Notes        string  `json:"notes"`
		ByUser       string  `json:"byUser"`
		Taste        float32 `json:"taste"`
		Difficulty   float32 `json:"difficulty"`
		Cleanup      float32 `json:"cleanup"`
		Instructions float32 `json:"instructions"`
		Overall      float32 `json:"overall"`
	}

	// MealRatingDatabaseCreationInput represents what a user could set as input for creating meal ratings.
	MealRatingDatabaseCreationInput struct {
		_ struct{}

		ID           string
		MealID       string
		Notes        string
		ByUser       string
		Taste        float32
		Difficulty   float32
		Cleanup      float32
		Instructions float32
		Overall      float32
	}

	// MealRatingUpdateRequestInput represents what a user could set as input for updating meal ratings.
	MealRatingUpdateRequestInput struct {
		_ struct{}

		MealID       *string  `json:"mealID"`
		Taste        *float32 `json:"taste"`
		Difficulty   *float32 `json:"difficulty"`
		Cleanup      *float32 `json:"cleanup"`
		Instructions *float32 `json:"instructions"`
		Overall      *float32 `json:"overall"`
		Notes        *string  `json:"notes"`
		ByUser       *string  `json:"byUser"`
	}

	// MealRatingDataManager describes a structure capable of storing meal ratings permanently.
	MealRatingDataManager interface {
		MealRatingExists(ctx context.Context, mealRatingID string) (bool, error)
		GetMealRating(ctx context.Context, mealRatingID string) (*MealRating, error)
		GetMealRatings(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[MealRating], error)
		CreateMealRating(ctx context.Context, input *MealRatingDatabaseCreationInput) (*MealRating, error)
		UpdateMealRating(ctx context.Context, updated *MealRating) error
		ArchiveMealRating(ctx context.Context, mealRatingID string) error
	}

	// MealRatingDataService describes a structure capable of serving traffic related to meal ratings.
	MealRatingDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an MealRatingUpdateRequestInput with a meal rating.
func (x *MealRating) Update(input *MealRatingUpdateRequestInput) {
	if input.MealID != nil && *input.MealID != x.MealID {
		x.MealID = *input.MealID
	}

	if input.Taste != nil && *input.Taste != x.Taste {
		x.Taste = *input.Taste
	}

	if input.Difficulty != nil && *input.Difficulty != x.Difficulty {
		x.Difficulty = *input.Difficulty
	}

	if input.Cleanup != nil && *input.Cleanup != x.Cleanup {
		x.Cleanup = *input.Cleanup
	}

	if input.Instructions != nil && *input.Instructions != x.Instructions {
		x.Instructions = *input.Instructions
	}

	if input.Overall != nil && *input.Overall != x.Overall {
		x.Overall = *input.Overall
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
}

var _ validation.ValidatableWithContext = (*MealRatingCreationRequestInput)(nil)

// ValidateWithContext validates a MealRatingCreationRequestInput.
func (x *MealRatingCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var errs *multierror.Error

	if x.Cleanup == 0 && x.Difficulty == 0 && x.Instructions == 0 && x.Overall == 0 && x.Taste == 0 {
		errs = multierror.Append(errs, errors.New("meal rating must have at least one rating"))
	}

	if err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealID, validation.Required),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*MealRatingDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealRatingDatabaseCreationInput.
func (x *MealRatingDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	var errs *multierror.Error

	if x.Cleanup == 0 && x.Difficulty == 0 && x.Instructions == 0 && x.Overall == 0 && x.Taste == 0 {
		errs = multierror.Append(errs, errors.New("meal rating must have at least one rating"))
	}

	if err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MealID, validation.Required),
		validation.Field(&x.ByUser, validation.Required),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*MealRatingUpdateRequestInput)(nil)

// ValidateWithContext validates a MealRatingUpdateRequestInput.
func (x *MealRatingUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ByUser, validation.Required),
		validation.Field(&x.MealID, validation.Required),
	)
}
