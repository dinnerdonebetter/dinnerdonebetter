package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepProductDataType indicates an event is related to a recipe step product.
	RecipeStepProductDataType dataType = "recipe_step_product"
)

func init() {
	gob.Register(new(RecipeStepProduct))
	gob.Register(new(RecipeStepProductList))
	gob.Register(new(RecipeStepProductCreationRequestInput))
	gob.Register(new(RecipeStepProductUpdateRequestInput))
}

type (
	// RecipeStepProduct represents a recipe step product.
	RecipeStepProduct struct {
		_                   struct{}
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		ID                  string  `json:"id"`
		Name                string  `json:"name"`
		RecipeStepID        string  `json:"recipeStepID"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		CreatedOn           uint64  `json:"createdOn"`
	}

	// RecipeStepProductList represents a list of recipe step products.
	RecipeStepProductList struct {
		_                  struct{}
		RecipeStepProducts []*RecipeStepProduct `json:"recipeStepProducts"`
		Pagination
	}

	// RecipeStepProductCreationRequestInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductCreationRequestInput struct {
		_ struct{}

		ID                  string `json:"-"`
		Name                string `json:"name"`
		RecipeStepID        string `json:"recipeStepID"`
		BelongsToRecipeStep string `json:"-"`
	}

	// RecipeStepProductDatabaseCreationInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductDatabaseCreationInput struct {
		_ struct{}

		ID                  string `json:"id"`
		Name                string `json:"name"`
		RecipeStepID        string `json:"recipeStepID"`
		BelongsToRecipeStep string `json:"belongsToRecipeStep"`
	}

	// RecipeStepProductUpdateRequestInput represents what a user could set as input for updating recipe step products.
	RecipeStepProductUpdateRequestInput struct {
		_ struct{}

		Name                string `json:"name"`
		RecipeStepID        string `json:"recipeStepID"`
		BelongsToRecipeStep string `json:"belongsToRecipeStep"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently.
	RecipeStepProductDataManager interface {
		RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error)
		GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*RecipeStepProduct, error)
		GetTotalRecipeStepProductCount(ctx context.Context) (uint64, error)
		GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*RecipeStepProductList, error)
		GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*RecipeStepProduct, error)
		CreateRecipeStepProduct(ctx context.Context, input *RecipeStepProductDatabaseCreationInput) (*RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, updated *RecipeStepProduct) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error
	}

	// RecipeStepProductDataService describes a structure capable of serving traffic related to recipe step products.
	RecipeStepProductDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepProductUpdateRequestInput with a recipe step product.
func (x *RecipeStepProduct) Update(input *RecipeStepProductUpdateRequestInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.RecipeStepID != "" && input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = input.RecipeStepID
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepProductCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductCreationRequestInput.
func (x *RecipeStepProductCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepProductDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepProductDatabaseCreationInput.
func (x *RecipeStepProductDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
	)
}

// RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput(input *RecipeStepProductCreationRequestInput) *RecipeStepProductDatabaseCreationInput {
	x := &RecipeStepProductDatabaseCreationInput{
		Name:         input.Name,
		RecipeStepID: input.RecipeStepID,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeStepProductUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductUpdateRequestInput.
func (x *RecipeStepProductUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
	)
}
