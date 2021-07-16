package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// RecipeStepProduct represents a recipe step product.
	RecipeStepProduct struct {
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		ExternalID          string  `json:"externalID"`
		QuantityNotes       string  `json:"quantityNotes"`
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		ID                  uint64  `json:"id"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		CreatedOn           uint64  `json:"createdOn"`
		QuantityValue       float32 `json:"quantityValue"`
	}

	// RecipeStepProductList represents a list of recipe step products.
	RecipeStepProductList struct {
		RecipeStepProducts []*RecipeStepProduct `json:"recipeStepProducts"`
		Pagination
	}

	// RecipeStepProductCreationInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductCreationInput struct {
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		BelongsToRecipeStep uint64  `json:"-"`
		QuantityValue       float32 `json:"quantityValue"`
	}

	// RecipeStepProductUpdateInput represents what a user could set as input for updating recipe step products.
	RecipeStepProductUpdateInput struct {
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently.
	RecipeStepProductDataManager interface {
		RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (bool, error)
		GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*RecipeStepProduct, error)
		GetAllRecipeStepProductsCount(ctx context.Context) (uint64, error)
		GetAllRecipeStepProducts(ctx context.Context, resultChannel chan []*RecipeStepProduct, bucketSize uint16) error
		GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepProductList, error)
		GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) ([]*RecipeStepProduct, error)
		CreateRecipeStepProduct(ctx context.Context, input *RecipeStepProductCreationInput, createdByUser uint64) (*RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, updated *RecipeStepProduct, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID, archivedBy uint64) error
		GetAuditLogEntriesForRecipeStepProduct(ctx context.Context, recipeStepProductID uint64) ([]*AuditLogEntry, error)
	}

	// RecipeStepProductDataService describes a structure capable of serving traffic related to recipe step products.
	RecipeStepProductDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepProductUpdateInput with a recipe step product.
func (x *RecipeStepProduct) Update(input *RecipeStepProductUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Name != "" && input.Name != x.Name {
		out = append(out, &FieldChangeSummary{
			FieldName: "Name",
			OldValue:  x.Name,
			NewValue:  input.Name,
		})

		x.Name = input.Name
	}

	if input.QuantityType != "" && input.QuantityType != x.QuantityType {
		out = append(out, &FieldChangeSummary{
			FieldName: "QuantityType",
			OldValue:  x.QuantityType,
			NewValue:  input.QuantityType,
		})

		x.QuantityType = input.QuantityType
	}

	if input.QuantityValue != 0 && input.QuantityValue != x.QuantityValue {
		out = append(out, &FieldChangeSummary{
			FieldName: "QuantityValue",
			OldValue:  x.QuantityValue,
			NewValue:  input.QuantityValue,
		})

		x.QuantityValue = input.QuantityValue
	}

	if input.QuantityNotes != "" && input.QuantityNotes != x.QuantityNotes {
		out = append(out, &FieldChangeSummary{
			FieldName: "QuantityNotes",
			OldValue:  x.QuantityNotes,
			NewValue:  input.QuantityNotes,
		})

		x.QuantityNotes = input.QuantityNotes
	}

	if input.RecipeStepID != 0 && input.RecipeStepID != x.RecipeStepID {
		out = append(out, &FieldChangeSummary{
			FieldName: "RecipeStepID",
			OldValue:  x.RecipeStepID,
			NewValue:  input.RecipeStepID,
		})

		x.RecipeStepID = input.RecipeStepID
	}

	return out
}

var _ validation.ValidatableWithContext = (*RecipeStepProductCreationInput)(nil)

// ValidateWithContext validates a RecipeStepProductCreationInput.
func (x *RecipeStepProductCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
		validation.Field(&x.QuantityNotes, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepProductUpdateInput)(nil)

// ValidateWithContext validates a RecipeStepProductUpdateInput.
func (x *RecipeStepProductUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
		validation.Field(&x.QuantityNotes, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
	)
}
