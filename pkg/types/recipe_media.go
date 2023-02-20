package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeMediaDataType indicates an event is related to recipe media.
	RecipeMediaDataType dataType = "recipe_media"

	// RecipeMediaCreatedCustomerEventType indicates recipe media was created.
	RecipeMediaCreatedCustomerEventType CustomerEventType = "recipe_media_created"
	// RecipeMediaUpdatedCustomerEventType indicates recipe media was updated.
	RecipeMediaUpdatedCustomerEventType CustomerEventType = "recipe_media_updated"
	// RecipeMediaArchivedCustomerEventType indicates recipe media was archived.
	RecipeMediaArchivedCustomerEventType CustomerEventType = "recipe_media_archived"
)

func init() {
	gob.Register(new(RecipeMedia))
	gob.Register(new(RecipeMediaCreationRequestInput))
	gob.Register(new(RecipeMediaUpdateRequestInput))
}

type (
	// RecipeMedia represents recipe media.
	RecipeMedia struct {
		_ struct{}

		CreatedAt           time.Time  `json:"createdAt"`
		ArchivedAt          *time.Time `json:"archivedAt"`
		LastUpdatedAt       *time.Time `json:"lastUpdatedAt"`
		ID                  string     `json:"id"`
		BelongsToRecipe     *string    `json:"belongsToRecipe"`
		BelongsToRecipeStep *string    `json:"belongsToRecipeStep"`
		MimeType            string     `json:"mimeType"`
		InternalPath        string     `json:"internalPath"`
		ExternalPath        string     `json:"externalPath"`
		Index               uint16     `json:"index"`
	}

	// RecipeMediaCreationRequestInput represents what a user could set as input for creating valid preparations.
	RecipeMediaCreationRequestInput struct {
		_ struct{}

		BelongsToRecipe     *string `json:"belongsToRecipe"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep"`
		MimeType            string  `json:"mimeType"`
		InternalPath        string  `json:"internalPath"`
		ExternalPath        string  `json:"externalPath"`
		Index               uint16  `json:"index"`
	}

	// RecipeMediaDatabaseCreationInput represents what a user could set as input for creating valid preparations.
	RecipeMediaDatabaseCreationInput struct {
		_ struct{}

		ID                  string
		BelongsToRecipe     *string
		BelongsToRecipeStep *string
		MimeType            string
		InternalPath        string
		ExternalPath        string
		Index               uint16
	}

	// RecipeMediaUpdateRequestInput represents what a user could set as input for updating valid preparations.
	RecipeMediaUpdateRequestInput struct {
		_ struct{}

		BelongsToRecipe     *string `json:"belongsToRecipe"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep"`
		MimeType            *string `json:"mimeType"`
		InternalPath        *string `json:"internalPath"`
		ExternalPath        *string `json:"externalPath"`
		Index               *uint16 `json:"index"`
	}

	// RecipeMediaDataManager describes a structure capable of storing valid preparations permanently.
	RecipeMediaDataManager interface {
		RecipeMediaExists(ctx context.Context, validPreparationID string) (bool, error)
		GetRecipeMedia(ctx context.Context, validPreparationID string) (*RecipeMedia, error)
		CreateRecipeMedia(ctx context.Context, input *RecipeMediaDatabaseCreationInput) (*RecipeMedia, error)
		UpdateRecipeMedia(ctx context.Context, updated *RecipeMedia) error
		ArchiveRecipeMedia(ctx context.Context, validPreparationID string) error
	}

	// RecipeMediaDataService describes a structure capable of serving traffic related to valid preparations.
	RecipeMediaDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		RandomHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeMediaUpdateRequestInput with recipe media.
func (x *RecipeMedia) Update(input *RecipeMediaUpdateRequestInput) {
	if input.BelongsToRecipe != nil && input.BelongsToRecipe != x.BelongsToRecipe {
		x.BelongsToRecipe = input.BelongsToRecipe
	}

	if input.BelongsToRecipeStep != nil && input.BelongsToRecipeStep != x.BelongsToRecipeStep {
		x.BelongsToRecipeStep = input.BelongsToRecipeStep
	}

	if input.MimeType != nil && *input.MimeType != x.MimeType {
		x.MimeType = *input.MimeType
	}

	if input.InternalPath != nil && *input.InternalPath != x.InternalPath {
		x.InternalPath = *input.InternalPath
	}

	if input.ExternalPath != nil && *input.ExternalPath != x.ExternalPath {
		x.ExternalPath = *input.ExternalPath
	}

	if input.Index != nil && *input.Index != x.Index {
		x.Index = *input.Index
	}
}

var _ validation.ValidatableWithContext = (*RecipeMediaCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeMediaCreationRequestInput.
func (x *RecipeMediaCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		// TODO: populate this
	)
}

var _ validation.ValidatableWithContext = (*RecipeMediaDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeMediaDatabaseCreationInput.
func (x *RecipeMediaDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MimeType, validation.Required),
		validation.Field(&x.InternalPath, validation.Required),
		validation.Field(&x.ExternalPath, validation.Required),
		validation.Field(&x.Index, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeMediaUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeMediaUpdateRequestInput.
func (x *RecipeMediaUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		// TODO: populate this
	)
}
