package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	recipesOnRecipeStepsJoinClause = "recipes ON recipe_steps.belongs_to_recipe=recipes.id"
)

var (
	_ types.RecipeStepDataManager = (*SQLQuerier)(nil)

	// recipeStepsTableColumns are the columns for the recipe_steps table.
	recipeStepsTableColumns = []string{
		"recipe_steps.id",
		"recipe_steps.index",
		"valid_preparations.id",
		"valid_preparations.name",
		"valid_preparations.description",
		"valid_preparations.icon_path",
		"valid_preparations.created_on",
		"valid_preparations.last_updated_on",
		"valid_preparations.archived_on",
		"recipe_steps.prerequisite_step",
		"recipe_steps.min_estimated_time_in_seconds",
		"recipe_steps.max_estimated_time_in_seconds",
		"recipe_steps.temperature_in_celsius",
		"recipe_steps.notes",
		"recipe_steps.optional",
		"recipe_steps.created_on",
		"recipe_steps.last_updated_on",
		"recipe_steps.archived_on",
		"recipe_steps.belongs_to_recipe",
	}

	getRecipeStepsJoins = []string{
		recipesOnRecipeStepsJoinClause,
		validPreparationsOnRecipeStepsJoinClause,
	}
)

// scanRecipeStep takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step struct.
func (q *SQLQuerier) scanRecipeStep(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStep{}

	targetVars := []interface{}{
		&x.ID,
		&x.Index,
		&x.Preparation.ID,
		&x.Preparation.Name,
		&x.Preparation.Description,
		&x.Preparation.IconPath,
		&x.Preparation.CreatedOn,
		&x.Preparation.LastUpdatedOn,
		&x.Preparation.ArchivedOn,
		&x.PrerequisiteStep,
		&x.MinEstimatedTimeInSeconds,
		&x.MaxEstimatedTimeInSeconds,
		&x.TemperatureInCelsius,
		&x.Notes,
		&x.Optional,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipe,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeSteps takes some database rows and turns them into a slice of recipe steps.
func (q *SQLQuerier) scanRecipeSteps(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeSteps []*types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStep(ctx, rows, includeCounts)
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

		recipeSteps = append(recipeSteps, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeSteps, filteredCount, totalCount, nil
}

const recipeStepExistenceQuery = "SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.archived_on IS NULL AND recipes.id = $3 )"

// RecipeStepExists fetches whether a recipe step exists from the database.
func (q *SQLQuerier) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	args := []interface{}{
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step existence check")
	}

	return result, nil
}

const getRecipeStepQuery = `SELECT 
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.created_on,
	valid_preparations.last_updated_on,
	valid_preparations.archived_on,
	recipe_steps.prerequisite_step,
	recipe_steps.min_estimated_time_in_seconds,
	recipe_steps.max_estimated_time_in_seconds,
	recipe_steps.temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.optional,
	recipe_steps.created_on,
	recipe_steps.last_updated_on,
	recipe_steps.archived_on,
	recipe_steps.belongs_to_recipe 
FROM recipe_steps
JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_on IS NULL
AND recipe_steps.belongs_to_recipe = $1
AND recipe_steps.id = $2
AND recipes.archived_on IS NULL 
AND recipes.id = $3
`

// GetRecipeStep fetches a recipe step from the database.
func (q *SQLQuerier) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	args := []interface{}{
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipeStep", getRecipeStepQuery, args)

	recipeStep, _, _, err := q.scanRecipeStep(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipeStep")
	}

	return recipeStep, nil
}

const getTotalRecipeStepsCountQuery = "SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL"

// GetTotalRecipeStepCount fetches the count of recipe steps from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalRecipeStepCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalRecipeStepsCountQuery, "fetching count of recipe steps")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipe steps")
	}

	return count, nil
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeSteps(ctx context.Context, recipeID string, filter *types.QueryFilter) (x *types.RecipeStepList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	x = &types.RecipeStepList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "recipe_steps", getRecipeStepsJoins, []string{"valid_preparations.id"}, nil, householdOwnershipColumn, recipeStepsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "recipeSteps", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe steps list retrieval query")
	}

	if x.RecipeSteps, x.FilteredCount, x.TotalCount, err = q.scanRecipeSteps(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe steps")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetRecipeStepsWithIDsQuery(ctx context.Context, recipeID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"recipe_steps.id":                ids,
		"recipe_steps.archived_on":       nil,
		"recipe_steps.belongs_to_recipe": recipeID,
	}

	subqueryBuilder := q.sqlBuilder.Select(recipeStepsTableColumns...).
		From("recipe_steps").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(recipeStepsTableColumns...).
		FromSelect(subqueryBuilder, "recipe_steps").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetRecipeStepsWithIDs fetches recipe steps from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipeStepsWithIDs(ctx context.Context, recipeID string, limit uint8, ids []string) ([]*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

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

	query, args := q.buildGetRecipeStepsWithIDsQuery(ctx, recipeID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipe steps with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe steps from database")
	}

	recipeSteps, _, _, err := q.scanRecipeSteps(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe steps")
	}

	return recipeSteps, nil
}

const recipeStepCreationQuery = "INSERT INTO recipe_steps (id,index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,optional,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"

// CreateRecipeStep creates a recipe step in the database.
func (q *SQLQuerier) createRecipeStep(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Index,
		input.PreparationID,
		input.PrerequisiteStep,
		input.MinEstimatedTimeInSeconds,
		input.MaxEstimatedTimeInSeconds,
		input.TemperatureInCelsius,
		input.Notes,
		input.Optional,
		input.BelongsToRecipe,
	}

	// create the recipe step.
	if err := q.performWriteQuery(ctx, db, "recipe step creation", recipeStepCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing recipe step creation")
	}

	x := &types.RecipeStep{
		ID:                        input.ID,
		Index:                     input.Index,
		PrerequisiteStep:          input.PrerequisiteStep,
		MinEstimatedTimeInSeconds: input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: input.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      input.TemperatureInCelsius,
		Notes:                     input.Notes,
		Optional:                  input.Optional,
		BelongsToRecipe:           input.BelongsToRecipe,
		CreatedOn:                 q.currentTime(),
	}

	for _, ingredientInput := range input.Ingredients {
		ingredientInput.BelongsToRecipeStep = x.ID
		ingredient, createErr := q.createRecipeStepIngredient(ctx, db, ingredientInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step ingredient")
		}

		x.Ingredients = append(x.Ingredients, ingredient)
	}

	for _, productInput := range input.Products {
		productInput.BelongsToRecipeStep = x.ID
		product, createErr := q.createRecipeStepProduct(ctx, db, productInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step product")
		}

		x.Products = append(x.Products, product)
	}

	tracing.AttachRecipeStepIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *SQLQuerier) CreateRecipeStep(ctx context.Context, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	return q.createRecipeStep(ctx, q.db, input)
}

const updateRecipeStepQuery = "UPDATE recipe_steps SET index = $1, preparation_id = $2, prerequisite_step = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, temperature_in_celsius = $6, notes = $7, optional = $8, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $9 AND id = $10"

// UpdateRecipeStep updates a particular recipe step.
func (q *SQLQuerier) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIDKey, updated.ID)
	tracing.AttachRecipeStepIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Index,
		updated.Preparation.ID,
		updated.PrerequisiteStep,
		updated.MinEstimatedTimeInSeconds,
		updated.MaxEstimatedTimeInSeconds,
		updated.TemperatureInCelsius,
		updated.Notes,
		updated.Optional,
		updated.BelongsToRecipe,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step update", updateRecipeStepQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step updated")

	return nil
}

const archiveRecipeStepQuery = "UPDATE recipe_steps SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2"

// ArchiveRecipeStep archives a recipe step from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	args := []interface{}{
		recipeID,
		recipeStepID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step archive", archiveRecipeStepQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step archived")

	return nil
}
