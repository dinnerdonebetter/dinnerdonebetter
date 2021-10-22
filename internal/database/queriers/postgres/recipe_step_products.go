package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	recipeStepsOnRecipeStepProductsJoinClause = "recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id"
)

var (
	_ types.RecipeStepProductDataManager = (*SQLQuerier)(nil)

	// recipeStepProductsTableColumns are the columns for the recipe_step_products table.
	recipeStepProductsTableColumns = []string{
		"recipe_step_products.id",
		"recipe_step_products.name",
		"recipe_step_products.recipe_step_id",
		"recipe_step_products.created_on",
		"recipe_step_products.last_updated_on",
		"recipe_step_products.archived_on",
		"recipe_step_products.belongs_to_recipe_step",
	}

	getRecipeStepProductsJoins = []string{
		recipeStepsOnRecipeStepProductsJoinClause,
		recipesOnRecipeStepsJoinClause,
	}
)

// scanRecipeStepProduct takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step product struct.
func (q *SQLQuerier) scanRecipeStepProduct(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStepProduct{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.RecipeStepID,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeStep,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepProducts takes some database rows and turns them into a slice of recipe step products.
func (q *SQLQuerier) scanRecipeStepProducts(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepProducts []*types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepProduct(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		recipeStepProducts = append(recipeStepProducts, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeStepProducts, filteredCount, totalCount, nil
}

const recipeStepProductExistenceQuery = "SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 )"

// RecipeStepProductExists fetches whether a recipe step product exists from the database.
func (q *SQLQuerier) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	args := []interface{}{
		recipeStepID,
		recipeStepProductID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepProductExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step product existence check")
	}

	return result, nil
}

const getRecipeStepProductQuery = "SELECT recipe_step_products.id, recipe_step_products.name, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5"

// GetRecipeStepProduct fetches a recipe step product from the database.
func (q *SQLQuerier) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	args := []interface{}{
		recipeStepID,
		recipeStepProductID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipeStepProduct", getRecipeStepProductQuery, args)

	recipeStepProduct, _, _, err := q.scanRecipeStepProduct(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipeStepProduct")
	}

	return recipeStepProduct, nil
}

const getTotalRecipeStepProductsCountQuery = "SELECT COUNT(recipe_step_products.id) FROM recipe_step_products WHERE recipe_step_products.archived_on IS NULL"

// GetTotalRecipeStepProductCount fetches the count of recipe step products from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalRecipeStepProductCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, getTotalRecipeStepProductsCountQuery, "fetching count of recipe step products")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipe step products")
	}

	return count, nil
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.RecipeStepProductList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	x = &types.RecipeStepProductList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(
		ctx,
		"recipe_step_products",
		getRecipeStepProductsJoins,
		nil,
		householdOwnershipColumn,
		recipeStepProductsTableColumns,
		"",
		false,
		filter,
	)

	rows, err := q.performReadQuery(ctx, q.db, "recipeStepProducts", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe step products list retrieval query")
	}

	if x.RecipeStepProducts, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepProducts(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step products")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetRecipeStepProductsWithIDsQuery(ctx context.Context, recipeStepID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"recipe_step_products.id":                     ids,
		"recipe_step_products.archived_on":            nil,
		"recipe_step_products.belongs_to_recipe_step": recipeStepID,
	}

	subqueryBuilder := q.sqlBuilder.Select(recipeStepProductsTableColumns...).
		From("recipe_step_products").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(recipeStepProductsTableColumns...).
		FromSelect(subqueryBuilder, "recipe_step_products").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetRecipeStepProductsWithIDs fetches recipe step products from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if ids == nil {
		return nil, ErrNilInputProvided
	}

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.buildGetRecipeStepProductsWithIDsQuery(ctx, recipeStepID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipe step products with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step products from database")
	}

	recipeStepProducts, _, _, err := q.scanRecipeStepProducts(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step products")
	}

	return recipeStepProducts, nil
}

const recipeStepProductCreationQuery = "INSERT INTO recipe_step_products (id,name,recipe_step_id,belongs_to_recipe_step) VALUES ($1,$2,$3,$4)"

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *SQLQuerier) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.RecipeStepID,
		input.BelongsToRecipeStep,
	}

	// create the recipe step product.
	if err := q.performWriteQuery(ctx, q.db, "recipe step product creation", recipeStepProductCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating recipe step product")
	}

	x := &types.RecipeStepProduct{
		ID:                  input.ID,
		Name:                input.Name,
		RecipeStepID:        input.RecipeStepID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		CreatedOn:           q.currentTime(),
	}

	tracing.AttachRecipeStepProductIDToSpan(span, x.ID)
	logger.Info("recipe step product created")

	return x, nil
}

const updateRecipeStepProductQuery = "UPDATE recipe_step_products SET name = $1, recipe_step_id = $2, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $3 AND id = $4"

// UpdateRecipeStepProduct updates a particular recipe step product.
func (q *SQLQuerier) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.RecipeStepID,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step product update", updateRecipeStepProductQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product updated")

	return nil
}

const archiveRecipeStepProductQuery = "UPDATE recipe_step_products SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2"

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	args := []interface{}{
		recipeStepID,
		recipeStepProductID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step product archive", archiveRecipeStepProductQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product archived")

	return nil
}
