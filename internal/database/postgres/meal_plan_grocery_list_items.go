package postgres

import (
	"context"
	_ "embed"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.MealPlanGroceryListItemDataManager = (*Querier)(nil)

	// mealPlanGroceryListItemsTableColumns are the columns for the meal_plan_grocery_list_items table.
	mealPlanGroceryListItemsTableColumns = []string{
		"meal_plan_grocery_list_items.id",
		"meal_plan_grocery_list_items.belongs_to_meal_plan_option",
		"meal_plan_grocery_list_items.valid_ingredient",
		"meal_plan_grocery_list_items.valid_measurement_unit",
		"meal_plan_grocery_list_items.minimum_quantity_needed",
		"meal_plan_grocery_list_items.maximum_quantity_needed",
		"meal_plan_grocery_list_items.quantity_purchased",
		"meal_plan_grocery_list_items.purchased_measurement_unit",
		"meal_plan_grocery_list_items.purchased_upc",
		"meal_plan_grocery_list_items.purchase_price",
		"meal_plan_grocery_list_items.status_explanation",
		"meal_plan_grocery_list_items.status",
		"meal_plan_grocery_list_items.created_at",
		"meal_plan_grocery_list_items.last_updated_at",
		"meal_plan_grocery_list_items.completed_at",
	}
)

// scanMealPlanGroceryListItem takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan grocery list struct.
func (q *Querier) scanMealPlanGroceryListItem(ctx context.Context, scan database.Scanner) (x *types.MealPlanGroceryListItem, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealPlanGroceryListItem{}

	var (
		purchasedMeasurementUnitID *string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.BelongsToMealPlan,
		&x.Ingredient.ID,
		&x.MeasurementUnit.ID,
		&x.MinimumQuantityNeeded,
		&x.MaximumQuantityNeeded,
		&x.QuantityPurchased,
		&purchasedMeasurementUnitID,
		&x.PurchasedUPC,
		&x.PurchasePrice,
		&x.StatusExplanation,
		&x.Status,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "")
	}

	if purchasedMeasurementUnitID != nil {
		x.PurchasedMeasurementUnit = &types.ValidMeasurementUnit{ID: *purchasedMeasurementUnitID}
	}

	return x, nil
}

// scanMealPlanGroceryListItems takes some database rows and turns them into a slice of meal plan grocery lists.
func (q *Querier) scanMealPlanGroceryListItems(ctx context.Context, rows database.ResultIterator) (mealPlanGroceryListItems []*types.MealPlanGroceryListItem, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, scanErr := q.scanMealPlanGroceryListItem(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		mealPlanGroceryListItems = append(mealPlanGroceryListItems, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "handling rows")
	}

	return mealPlanGroceryListItems, nil
}

//go:embed queries/meal_plan_grocery_list_items/exists.sql
var mealPlanGroceryListItemExistenceQuery string

// MealPlanGroceryListItemExists fetches whether a meal plan grocery list exists from the database.
func (q *Querier) MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	args := []interface{}{
		mealPlanGroceryListItemID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanGroceryListItemExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan grocery list existence check")
	}

	return result, nil
}

func (q *Querier) fleshOutMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItem *types.MealPlanGroceryListItem) (*types.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanGroceryListItem == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItem.ID)
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItem.ID)

	validIngredient, err := q.GetValidIngredient(ctx, mealPlanGroceryListItem.Ingredient.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching grocery list item ingredient")
	}
	mealPlanGroceryListItem.Ingredient = *validIngredient

	validMeasurementUnit, err := q.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.MeasurementUnit.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching grocery list item measurement unit")
	}
	mealPlanGroceryListItem.MeasurementUnit = *validMeasurementUnit

	if mealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnit, getPurchasedMeasurementUnitErr := q.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
		if getPurchasedMeasurementUnitErr != nil {
			return nil, observability.PrepareAndLogError(getPurchasedMeasurementUnitErr, logger, span, "fetching grocery list item purchased measurement unit")
		}
		mealPlanGroceryListItem.PurchasedMeasurementUnit = purchasedMeasurementUnit
	}

	return mealPlanGroceryListItem, nil
}

//go:embed queries/meal_plan_grocery_list_items/get_one.sql
var getMealPlanGroceryListItemQuery string

// GetMealPlanGroceryListItem fetches a meal plan grocery list from the database.
func (q *Querier) GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	args := []interface{}{
		mealPlanID,
		mealPlanGroceryListItemID,
	}

	row := q.getOneRow(ctx, q.db, "meal plan grocery list item", getMealPlanGroceryListItemQuery, args)

	mealPlanGroceryListItem, err := q.scanMealPlanGroceryListItem(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal plan grocery list item")
	}

	// flesh out the data
	mealPlanGroceryListItem, err = q.fleshOutMealPlanGroceryListItem(ctx, mealPlanGroceryListItem)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "augmenting grocery list item data")
	}

	return mealPlanGroceryListItem, nil
}

//go:embed queries/meal_plan_grocery_list_items/get_all_for_meal_plan.sql
var getMealPlanGroceryListItemsForMealPlanQuery string

// GetMealPlanGroceryListItemsForMealPlan fetches a list of meal plan grocery lists from the database that meet a particular filter.
func (q *Querier) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanGroceryListItem, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	x = []*types.MealPlanGroceryListItem{}

	getMealPlanGroceryListItemsForMealPlanArgs := []interface{}{
		mealPlanID,
	}
	rows, err := q.getRows(ctx, q.db, "meal plan grocery list items", getMealPlanGroceryListItemsForMealPlanQuery, getMealPlanGroceryListItemsForMealPlanArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan grocery lists list retrieval query")
	}

	if x, err = q.scanMealPlanGroceryListItems(ctx, rows); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal plan grocery lists")
	}

	for i := range x {
		x[i], err = q.fleshOutMealPlanGroceryListItem(ctx, x[i])
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "augmenting grocery list item data")
		}
	}

	return x, nil
}

//go:embed queries/meal_plan_grocery_list_items/create.sql
var mealPlanGroceryListItemCreationQuery string

// createMealPlanGroceryListItem creates a meal plan grocery list in the database.
func (q *Querier) createMealPlanGroceryListItem(ctx context.Context, querier database.SQLQueryExecutor, input *types.MealPlanGroceryListItemDatabaseCreationInput) (*types.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanGroceryListItemIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.BelongsToMealPlan,
		input.ValidIngredientID,
		input.ValidMeasurementUnitID,
		input.MinimumQuantityNeeded,
		input.MaximumQuantityNeeded,
		input.QuantityPurchased,
		input.PurchasedMeasurementUnitID,
		input.PurchasedUPC,
		input.PurchasePrice,
		input.StatusExplanation,
		input.Status,
	}

	// create the meal plan grocery list.
	if err := q.performWriteQuery(ctx, querier, "meal plan grocery list creation", mealPlanGroceryListItemCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan grocery list creation query")
	}

	x := &types.MealPlanGroceryListItem{
		ID:                    input.ID,
		BelongsToMealPlan:     input.BelongsToMealPlan,
		Ingredient:            types.ValidIngredient{ID: input.ValidIngredientID},
		MeasurementUnit:       types.ValidMeasurementUnit{ID: input.ValidMeasurementUnitID},
		MinimumQuantityNeeded: input.MinimumQuantityNeeded,
		MaximumQuantityNeeded: input.MaximumQuantityNeeded,
		QuantityPurchased:     input.QuantityPurchased,
		PurchasedUPC:          input.PurchasedUPC,
		PurchasePrice:         input.PurchasePrice,
		StatusExplanation:     input.StatusExplanation,
		Status:                input.Status,
		CreatedAt:             q.currentTime(),
	}

	if input.PurchasedMeasurementUnitID != nil {
		x.PurchasedMeasurementUnit = &types.ValidMeasurementUnit{ID: *input.PurchasedMeasurementUnitID}
	}

	tracing.AttachMealPlanGroceryListItemIDToSpan(span, x.ID)
	logger.Info("meal plan grocery list created")

	return x, nil
}

// CreateMealPlanGroceryListItem creates a meal plan grocery list in the database.
func (q *Querier) CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemDatabaseCreationInput) (*types.MealPlanGroceryListItem, error) {
	return q.createMealPlanGroceryListItem(ctx, q.db, input)
}

// CreateMealPlanGroceryListItemsForMealPlan creates a meal plan grocery list in the database.
func (q *Querier) CreateMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string, inputs []*types.MealPlanGroceryListItemDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if inputs == nil {
		return ErrNilInputProvided
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	for _, input := range inputs {
		if _, err = q.createMealPlanGroceryListItem(ctx, tx, input); err != nil {
			q.rollbackTransaction(ctx, tx)
			return observability.PrepareAndLogError(err, logger, span, "updating meal plan grocery list")
		}
	}

	if err = q.MarkMealPlanAsHavingGroceryListInitialized(ctx, mealPlanID); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "marking meal plan grocery list as initialized")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}

//go:embed queries/meal_plan_grocery_list_items/update.sql
var updateMealPlanGroceryListItemQuery string

// UpdateMealPlanGroceryListItem updates a particular meal plan grocery list.
func (q *Querier) UpdateMealPlanGroceryListItem(ctx context.Context, updated *types.MealPlanGroceryListItem) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanGroceryListItemIDKey, updated.ID)
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, updated.ID)

	var purchasedMeasurementUnitID *string
	if updated.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnitID = &updated.PurchasedMeasurementUnit.ID
	}

	args := []interface{}{
		updated.BelongsToMealPlan,
		updated.Ingredient.ID,
		updated.MeasurementUnit.ID,
		updated.MinimumQuantityNeeded,
		updated.MaximumQuantityNeeded,
		updated.QuantityPurchased,
		purchasedMeasurementUnitID,
		updated.PurchasedUPC,
		updated.PurchasePrice,
		updated.StatusExplanation,
		updated.Status,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan grocery list update", updateMealPlanGroceryListItemQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan grocery list")
	}

	logger.Info("meal plan grocery list updated")

	return nil
}

//go:embed queries/meal_plan_grocery_list_items/archive.sql
var archiveMealPlanGroceryListItemQuery string

// ArchiveMealPlanGroceryListItem archives a meal plan grocery list from the database by its ID.
func (q *Querier) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanGroceryListItemID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	args := []interface{}{
		mealPlanGroceryListItemID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan grocery list archive", archiveMealPlanGroceryListItemQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan grocery list")
	}

	logger.Info("meal plan grocery list archived")

	return nil
}
