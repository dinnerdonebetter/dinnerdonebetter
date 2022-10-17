package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanGroceryListItemDataType indicates an event is related to a meal plan grocery list item.
	MealPlanGroceryListItemDataType dataType = "valid_preparation"

	// MealPlanGroceryListItemQuantityModifier is what we multiply / divide by for float quantity values.
	MealPlanGroceryListItemQuantityModifier = 100

	// MealPlanGroceryListItemCreatedCustomerEventType indicates a meal plan grocery list item was created.
	MealPlanGroceryListItemCreatedCustomerEventType CustomerEventType = "valid_preparation_created"
	// MealPlanGroceryListItemUpdatedCustomerEventType indicates a meal plan grocery list item was updated.
	MealPlanGroceryListItemUpdatedCustomerEventType CustomerEventType = "valid_preparation_updated"
	// MealPlanGroceryListItemArchivedCustomerEventType indicates a meal plan grocery list item was archived.
	MealPlanGroceryListItemArchivedCustomerEventType CustomerEventType = "valid_preparation_archived"
)

func init() {
	gob.Register(new(MealPlanGroceryListItem))
	gob.Register(new(MealPlanGroceryListItemList))
	gob.Register(new(MealPlanGroceryListItemCreationRequestInput))
	gob.Register(new(MealPlanGroceryListItemUpdateRequestInput))
}

type (
	// MealPlanGroceryListItem represents a meal plan grocery list item.
	MealPlanGroceryListItem struct {
		_                        struct{}
		CreatedAt                time.Time             `json:"createdAt"`
		QuantityPurchased        *float32              `json:"quantityPurchased"`
		LastUpdatedAt            *time.Time            `json:"lastUpdatedAt"`
		PurchasePrice            *float32              `json:"purchasePrice"`
		PurchasedUPC             *string               `json:"purchasedUPC"`
		PurchasedMeasurementUnit *ValidMeasurementUnit `json:"purchasedMeasurementUnit"`
		ArchivedAt               *time.Time            `json:"archivedAt"`
		ID                       string                `json:"id"`
		StatusExplanation        string                `json:"statusExplanation"`
		Status                   string                `json:"status"`
		MeasurementUnit          ValidMeasurementUnit  `json:"measurementUnit"`
		MealPlanOption           MealPlanOption        `json:"mealPlanOption"`
		Ingredient               ValidIngredient       `json:"ingredient"`
		MaximumQuantityNeeded    float32               `json:"maximumQuantityNeeded"`
		MinimumQuantityNeeded    float32               `json:"minimumQuantityNeeded"`
	}

	// MealPlanGroceryListItemList represents a list of meal plan grocery list items.
	MealPlanGroceryListItemList struct {
		_                        struct{}
		MealPlanGroceryListItems []*MealPlanGroceryListItem `json:"data"`
		Pagination
	}

	// MealPlanGroceryListItemCreationRequestInput represents what a user could set as input for creating meal plan grocery list items.
	MealPlanGroceryListItemCreationRequestInput struct {
		_                          struct{}
		PurchasedMeasurementUnitID *string  `json:"purchasedMeasurementUnitID"`
		PurchasedUPC               *string  `json:"purchasedUPC"`
		PurchasePrice              *float32 `json:"purchasePrice"`
		QuantityPurchased          *float32 `json:"quantityPurchased"`
		StatusExplanation          string   `json:"statusExplanation"`
		Status                     string   `json:"status"`
		MealPlanOptionID           string   `json:"mealPlanOptionID"`
		ValidIngredientID          string   `json:"validIngredientID"`
		ValidMeasurementUnitID     string   `json:"validMeasurementUnitID"`
		MinimumQuantityNeeded      float32  `json:"minimumQuantityNeeded"`
		MaximumQuantityNeeded      float32  `json:"maximumQuantityNeeded"`
	}

	// MealPlanGroceryListItemDatabaseCreationInput represents what a user could set as input for creating meal plan grocery list items.
	MealPlanGroceryListItemDatabaseCreationInput struct {
		_                          struct{}
		PurchasePrice              *float32
		PurchasedUPC               *string
		PurchasedMeasurementUnitID *string
		QuantityPurchased          *float32
		Status                     string
		StatusExplanation          string
		ValidMeasurementUnitID     string
		ValidIngredientID          string
		MealPlanOptionID           string
		ID                         string
		MinimumQuantityNeeded      float32
		MaximumQuantityNeeded      float32
	}

	// MealPlanGroceryListItemUpdateRequestInput represents what a user could set as input for updating meal plan grocery list items.
	MealPlanGroceryListItemUpdateRequestInput struct {
		_                          struct{}
		MaximumQuantityNeeded      *float32 `json:"maximumQuantityNeeded"`
		MealPlanOptionID           *string  `json:"mealPlanOptionID"`
		ValidIngredientID          *string  `json:"validIngredientID"`
		ValidMeasurementUnitID     *string  `json:"validMeasurementUnitID"`
		MinimumQuantityNeeded      *float32 `json:"minimumQuantityNeeded"`
		StatusExplanation          *string  `json:"statusExplanation"`
		QuantityPurchased          *float32 `json:"quantityPurchased"`
		PurchasedMeasurementUnitID *string  `json:"purchasedMeasurementUnitID"`
		PurchasedUPC               *string  `json:"purchasedUPC"`
		PurchasePrice              *float32 `json:"purchasePrice"`
		Status                     *string  `json:"status"`
		ID                         string   `json:"-"`
	}

	// MealPlanGroceryListItemDataManager describes a structure capable of storing meal plan grocery list items permanently.
	MealPlanGroceryListItemDataManager interface {
		MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (bool, error)
		GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*MealPlanGroceryListItem, error)
		GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*MealPlanGroceryListItem, error)
		CreateMealPlanGroceryListItem(ctx context.Context, input *MealPlanGroceryListItemDatabaseCreationInput) (*MealPlanGroceryListItem, error)
		UpdateMealPlanGroceryListItem(ctx context.Context, updated *MealPlanGroceryListItem) error
		ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error
	}

	// MealPlanGroceryListItemDataService describes a structure capable of serving traffic related to meal plan grocery list items.
	MealPlanGroceryListItemDataService interface {
		ListByMealPlanHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an MealPlanGroceryListItemUpdateRequestInput with a meal plan grocery list item.
func (x *MealPlanGroceryListItem) Update(input *MealPlanGroceryListItemUpdateRequestInput) {
	if input.MealPlanOptionID != nil && *input.MealPlanOptionID != x.MealPlanOption.ID {
		x.MealPlanOption = MealPlanOption{ID: *input.MealPlanOptionID}
	}
	if input.ValidIngredientID != nil && *input.ValidIngredientID != x.Ingredient.ID {
		x.Ingredient = ValidIngredient{ID: *input.ValidIngredientID}
	}
	if input.ValidMeasurementUnitID != nil && *input.ValidMeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit = ValidMeasurementUnit{ID: *input.ValidMeasurementUnitID}
	}
	if input.MinimumQuantityNeeded != nil && *input.MinimumQuantityNeeded != x.MinimumQuantityNeeded {
		x.MinimumQuantityNeeded = *input.MinimumQuantityNeeded
	}
	if input.MaximumQuantityNeeded != nil && *input.MaximumQuantityNeeded != x.MaximumQuantityNeeded {
		x.MaximumQuantityNeeded = *input.MaximumQuantityNeeded
	}
	// was nil and now will not be nil
	if input.QuantityPurchased != nil && ((x.QuantityPurchased == nil) || (x.QuantityPurchased != nil && (*input.QuantityPurchased != *x.QuantityPurchased))) {
		x.QuantityPurchased = input.QuantityPurchased
	}

	if input.PurchasedMeasurementUnitID != nil && *input.PurchasedMeasurementUnitID != x.PurchasedMeasurementUnit.ID {
		x.PurchasedMeasurementUnit = &ValidMeasurementUnit{ID: *input.PurchasedMeasurementUnitID}
	}

	if input.PurchasedUPC != nil && ((x.PurchasedUPC == nil) || (x.PurchasedUPC != nil && (*input.PurchasedUPC != *x.PurchasedUPC))) {
		x.PurchasedUPC = input.PurchasedUPC
	}
	if input.PurchasePrice != nil && ((x.PurchasePrice == nil) || (x.PurchasePrice != nil && (*input.PurchasePrice != *x.PurchasePrice))) {
		x.PurchasePrice = input.PurchasePrice
	}

	if input.StatusExplanation != nil && *input.StatusExplanation != x.StatusExplanation {
		x.StatusExplanation = *input.StatusExplanation
	}
	if input.Status != nil && *input.Status != x.Status {
		x.Status = *input.Status
	}
}

var _ validation.ValidatableWithContext = (*MealPlanGroceryListItemCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanGroceryListItemCreationRequestInput.
func (x *MealPlanGroceryListItemCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityNeeded, validation.Required),
		validation.Field(&x.MaximumQuantityNeeded, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanGroceryListItemDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanGroceryListItemDatabaseCreationInput.
func (x *MealPlanGroceryListItemDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityNeeded, validation.Required),
		validation.Field(&x.MaximumQuantityNeeded, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanGroceryListItemUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanGroceryListItemUpdateRequestInput.
func (x *MealPlanGroceryListItemUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityNeeded, validation.Required),
		validation.Field(&x.MaximumQuantityNeeded, validation.Required),
	)
}
