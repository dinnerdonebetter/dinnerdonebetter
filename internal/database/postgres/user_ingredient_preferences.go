package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	validIngredientsOnUserIngredientPreferencesJoin = "valid_ingredients ON valid_ingredients.id = user_ingredient_preferences.ingredient"
)

var (
	_ types.UserIngredientPreferenceDataManager = (*Querier)(nil)

	// userIngredientPreferencesTableColumns are the columns for the user_ingredient_preferences table.
	userIngredientPreferencesTableColumns = []string{
		"user_ingredient_preferences.id",
		"valid_ingredients.id",
		"valid_ingredients.name",
		"valid_ingredients.description",
		"valid_ingredients.warning",
		"valid_ingredients.contains_egg",
		"valid_ingredients.contains_dairy",
		"valid_ingredients.contains_peanut",
		"valid_ingredients.contains_tree_nut",
		"valid_ingredients.contains_soy",
		"valid_ingredients.contains_wheat",
		"valid_ingredients.contains_shellfish",
		"valid_ingredients.contains_sesame",
		"valid_ingredients.contains_fish",
		"valid_ingredients.contains_gluten",
		"valid_ingredients.animal_flesh",
		"valid_ingredients.volumetric",
		"valid_ingredients.is_liquid",
		"valid_ingredients.icon_path",
		"valid_ingredients.animal_derived",
		"valid_ingredients.plural_name",
		"valid_ingredients.restrict_to_preparations",
		"valid_ingredients.minimum_ideal_storage_temperature_in_celsius",
		"valid_ingredients.maximum_ideal_storage_temperature_in_celsius",
		"valid_ingredients.storage_instructions",
		"valid_ingredients.slug",
		"valid_ingredients.contains_alcohol",
		"valid_ingredients.shopping_suggestions",
		"valid_ingredients.is_starch",
		"valid_ingredients.is_protein",
		"valid_ingredients.is_grain",
		"valid_ingredients.is_fruit",
		"valid_ingredients.is_salt",
		"valid_ingredients.is_fat",
		"valid_ingredients.is_acid",
		"valid_ingredients.is_heat",
		"valid_ingredients.created_at",
		"valid_ingredients.last_updated_at",
		"valid_ingredients.archived_at",
		"user_ingredient_preferences.rating",
		"user_ingredient_preferences.notes",
		"user_ingredient_preferences.allergy",
		"user_ingredient_preferences.created_at",
		"user_ingredient_preferences.last_updated_at",
		"user_ingredient_preferences.archived_at",
		"user_ingredient_preferences.belongs_to_user",
	}
)

// scanUserIngredientPreference takes a database Scanner (i.e. *sql.Row) and scans the result into a user ingredient preference struct.
func (q *Querier) scanUserIngredientPreference(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.UserIngredientPreference, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.UserIngredientPreference{}

	targetVars := []any{
		&x.ID,
		&x.Ingredient.ID,
		&x.Ingredient.Name,
		&x.Ingredient.Description,
		&x.Ingredient.Warning,
		&x.Ingredient.ContainsEgg,
		&x.Ingredient.ContainsDairy,
		&x.Ingredient.ContainsPeanut,
		&x.Ingredient.ContainsTreeNut,
		&x.Ingredient.ContainsSoy,
		&x.Ingredient.ContainsWheat,
		&x.Ingredient.ContainsShellfish,
		&x.Ingredient.ContainsSesame,
		&x.Ingredient.ContainsFish,
		&x.Ingredient.ContainsGluten,
		&x.Ingredient.AnimalFlesh,
		&x.Ingredient.IsMeasuredVolumetrically,
		&x.Ingredient.IsLiquid,
		&x.Ingredient.IconPath,
		&x.Ingredient.AnimalDerived,
		&x.Ingredient.PluralName,
		&x.Ingredient.RestrictToPreparations,
		&x.Ingredient.MinimumIdealStorageTemperatureInCelsius,
		&x.Ingredient.MaximumIdealStorageTemperatureInCelsius,
		&x.Ingredient.StorageInstructions,
		&x.Ingredient.Slug,
		&x.Ingredient.ContainsAlcohol,
		&x.Ingredient.ShoppingSuggestions,
		&x.Ingredient.IsStarch,
		&x.Ingredient.IsProtein,
		&x.Ingredient.IsGrain,
		&x.Ingredient.IsFruit,
		&x.Ingredient.IsSalt,
		&x.Ingredient.IsFat,
		&x.Ingredient.IsAcid,
		&x.Ingredient.IsHeat,
		&x.Ingredient.CreatedAt,
		&x.Ingredient.LastUpdatedAt,
		&x.Ingredient.ArchivedAt,
		&x.Rating,
		&x.Notes,
		&x.Allergy,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.BelongsToUser,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanUserIngredientPreferences takes some database rows and turns them into a slice of user ingredient preferences.
func (q *Querier) scanUserIngredientPreferences(ctx context.Context, rows database.ResultIterator, includeCounts bool) (userIngredientPreferences []*types.UserIngredientPreference, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanUserIngredientPreference(ctx, rows, includeCounts)
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

		userIngredientPreferences = append(userIngredientPreferences, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return userIngredientPreferences, filteredCount, totalCount, nil
}

//go:embed queries/user_ingredient_preferences/exists.sql
var userIngredientPreferenceExistenceQuery string

// UserIngredientPreferenceExists fetches whether a user ingredient preference exists from the database.
func (q *Querier) UserIngredientPreferenceExists(ctx context.Context, userIngredientPreferenceID, userID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userIngredientPreferenceID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, userIngredientPreferenceID)

	if userID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := []any{
		userIngredientPreferenceID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, userIngredientPreferenceExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing user ingredient preference existence check")
	}

	return result, nil
}

//go:embed queries/user_ingredient_preferences/get_one.sql
var getUserIngredientPreferenceQuery string

// GetUserIngredientPreference fetches a user ingredient preference from the database.
func (q *Querier) GetUserIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) (*types.UserIngredientPreference, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userIngredientPreferenceID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, userIngredientPreferenceID)

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := []any{
		userIngredientPreferenceID,
		userID,
	}

	row := q.getOneRow(ctx, q.db, "user ingredient preference", getUserIngredientPreferenceQuery, args)

	userIngredientPreference, _, _, err := q.scanUserIngredientPreference(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning userIngredientPreference")
	}

	return userIngredientPreference, nil
}

// GetUserIngredientPreferences fetches a list of user ingredient preferences from the database that meet a particular filter.
func (q *Querier) GetUserIngredientPreferences(ctx context.Context, userID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.UserIngredientPreference], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	x = &types.QueryFilteredResult[types.UserIngredientPreference]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "user_ingredient_preferences", []string{validIngredientsOnUserIngredientPreferencesJoin}, []string{"user_ingredient_preferences.id", "valid_ingredients.id"}, nil, userOwnershipColumn, userIngredientPreferencesTableColumns, userID, false, filter)

	rows, err := q.getRows(ctx, q.db, "user ingredient preferences", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing user ingredient preferences list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanUserIngredientPreferences(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning user ingredient preferences")
	}

	return x, nil
}

//go:embed queries/user_ingredient_preferences/create.sql
var userIngredientPreferenceCreationQuery string

// CreateUserIngredientPreference creates a user ingredient preference in the database.
func (q *Querier) CreateUserIngredientPreference(ctx context.Context, input *types.UserIngredientPreferenceDatabaseCreationInput) ([]*types.UserIngredientPreference, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.UserIngredientPreferenceIDKey, input.ID)

	validIngredientIDs := []string{}
	if input.ValidIngredientGroupID != "" {
		group, err := q.GetValidIngredientGroup(ctx, input.ValidIngredientGroupID)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient group")
		}

		for _, member := range group.Members {
			validIngredientIDs = append(validIngredientIDs, member.ValidIngredient.ID)
		}
	} else {
		validIngredientIDs = append(validIngredientIDs, input.ValidIngredientID)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	logger = logger.WithValue("valid_ingredient_ids", validIngredientIDs)
	logger.Debug("creating user ingredient preferences")

	output := []*types.UserIngredientPreference{}
	for _, validIngredientID := range validIngredientIDs {
		l := logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
		if validIngredientID == "" {
			continue
		}

		id := identifiers.New()
		args := []any{
			id,
			validIngredientID,
			input.Rating,
			input.Notes,
			input.Allergy,
			input.BelongsToUser,
		}

		// create the user ingredient preference.
		if err = q.performWriteQuery(ctx, tx, "user ingredient preference creation", userIngredientPreferenceCreationQuery, args); err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(err, l, span, "performing user ingredient preference creation query")
		}

		x := &types.UserIngredientPreference{
			ID:            id,
			Rating:        input.Rating,
			Notes:         input.Notes,
			Allergy:       input.Allergy,
			BelongsToUser: input.BelongsToUser,
			Ingredient:    types.ValidIngredient{ID: input.ValidIngredientID},
			CreatedAt:     q.currentTime(),
		}

		tracing.AttachUserIngredientPreferenceIDToSpan(span, x.ID)
		l.Info("user ingredient preference created")

		output = append(output, x)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return output, nil
}

//go:embed queries/user_ingredient_preferences/update.sql
var updateUserIngredientPreferenceQuery string

// UpdateUserIngredientPreference updates a particular user ingredient preference.
func (q *Querier) UpdateUserIngredientPreference(ctx context.Context, updated *types.UserIngredientPreference) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.UserIngredientPreferenceIDKey, updated.ID)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, updated.ID)

	args := []any{
		updated.Ingredient.ID,
		updated.Rating,
		updated.Notes,
		updated.Allergy,
		updated.ID,
		updated.BelongsToUser,
	}

	if err := q.performWriteQuery(ctx, q.db, "user ingredient preference update", updateUserIngredientPreferenceQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user ingredient preference")
	}

	logger.Info("user ingredient preference updated")

	return nil
}

//go:embed queries/user_ingredient_preferences/archive.sql
var archiveUserIngredientPreferenceQuery string

// ArchiveUserIngredientPreference archives a user ingredient preference from the database by its ID.
func (q *Querier) ArchiveUserIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	if userIngredientPreferenceID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, userIngredientPreferenceID)

	args := []any{
		userIngredientPreferenceID,
		userID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user ingredient preference archive", archiveUserIngredientPreferenceQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user ingredient preference")
	}

	logger.Info("user ingredient preference archived")

	return nil
}
