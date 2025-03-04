package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

var (
	_ types.IngredientPreferenceDataManager = (*Querier)(nil)
)

// IngredientPreferenceExists fetches whether a user ingredient preference exists from the database.
func (q *Querier) IngredientPreferenceExists(ctx context.Context, userIngredientPreferenceID, userID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userIngredientPreferenceID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.IngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachToSpan(span, keys.IngredientPreferenceIDKey, userIngredientPreferenceID)

	if userID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	exists, err = q.generatedQuerier.CheckUserIngredientPreferenceExistence(ctx, q.db, &generated.CheckUserIngredientPreferenceExistenceParams{
		ID:            userIngredientPreferenceID,
		BelongsToUser: userID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing user ingredient preference existence check")
	}

	return exists, nil
}

// GetIngredientPreference fetches a user ingredient preference from the database.
func (q *Querier) GetIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) (*types.IngredientPreference, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userIngredientPreferenceID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.IngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachToSpan(span, keys.IngredientPreferenceIDKey, userIngredientPreferenceID)

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	result, err := q.generatedQuerier.GetUserIngredientPreference(ctx, q.db, &generated.GetUserIngredientPreferenceParams{
		ID:            userIngredientPreferenceID,
		BelongsToUser: userID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning userIngredientPreference")
	}

	userIngredientPreference := &types.IngredientPreference{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		ID:            result.ID,
		Notes:         result.Notes,
		BelongsToUser: result.BelongsToUser,
		Rating:        int8(result.Rating),
		Allergy:       result.Allergy,
		Ingredient: types.ValidIngredient{
			CreatedAt:     result.ValidIngredientCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.ValidIngredientIconPath,
			Warning:                result.ValidIngredientWarning,
			PluralName:             result.ValidIngredientPluralName,
			StorageInstructions:    result.ValidIngredientStorageInstructions,
			Name:                   result.ValidIngredientName,
			ID:                     result.ValidIngredientID,
			Description:            result.ValidIngredientDescription,
			Slug:                   result.ValidIngredientSlug,
			ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
			ContainsShellfish:      result.ValidIngredientContainsShellfish,
			IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:         result.ValidIngredientContainsPeanut,
			ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
			ContainsEgg:            result.ValidIngredientContainsEgg,
			ContainsWheat:          result.ValidIngredientContainsWheat,
			ContainsSoy:            result.ValidIngredientContainsSoy,
			AnimalDerived:          result.ValidIngredientAnimalDerived,
			RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
			ContainsSesame:         result.ValidIngredientContainsSesame,
			ContainsFish:           result.ValidIngredientContainsFish,
			ContainsGluten:         result.ValidIngredientContainsGluten,
			ContainsDairy:          result.ValidIngredientContainsDairy,
			ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
			AnimalFlesh:            result.ValidIngredientAnimalFlesh,
			IsStarch:               result.ValidIngredientIsStarch,
			IsProtein:              result.ValidIngredientIsProtein,
			IsGrain:                result.ValidIngredientIsGrain,
			IsFruit:                result.ValidIngredientIsFruit,
			IsSalt:                 result.ValidIngredientIsSalt,
			IsFat:                  result.ValidIngredientIsFat,
			IsAcid:                 result.ValidIngredientIsAcid,
			IsHeat:                 result.ValidIngredientIsHeat,
		},
	}

	return userIngredientPreference, nil
}

// GetIngredientPreferences fetches a list of user ingredient preferences from the database that meet a particular filter.
func (q *Querier) GetIngredientPreferences(ctx context.Context, userID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.IngredientPreference], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.IngredientPreference]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetUserIngredientPreferencesForUser(ctx, q.db, &generated.GetUserIngredientPreferencesForUserParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
		BelongsToUser:   userID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing user ingredient preferences list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.IngredientPreference{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			ID:            result.ID,
			Notes:         result.Notes,
			BelongsToUser: result.BelongsToUser,
			Rating:        int8(result.Rating),
			Allergy:       result.Allergy,
			Ingredient: types.ValidIngredient{
				CreatedAt:     result.ValidIngredientCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				},
				IconPath:               result.ValidIngredientIconPath,
				Warning:                result.ValidIngredientWarning,
				PluralName:             result.ValidIngredientPluralName,
				StorageInstructions:    result.ValidIngredientStorageInstructions,
				Name:                   result.ValidIngredientName,
				ID:                     result.ValidIngredientID,
				Description:            result.ValidIngredientDescription,
				Slug:                   result.ValidIngredientSlug,
				ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:      result.ValidIngredientContainsShellfish,
				IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:         result.ValidIngredientContainsPeanut,
				ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
				ContainsEgg:            result.ValidIngredientContainsEgg,
				ContainsWheat:          result.ValidIngredientContainsWheat,
				ContainsSoy:            result.ValidIngredientContainsSoy,
				AnimalDerived:          result.ValidIngredientAnimalDerived,
				RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
				ContainsSesame:         result.ValidIngredientContainsSesame,
				ContainsFish:           result.ValidIngredientContainsFish,
				ContainsGluten:         result.ValidIngredientContainsGluten,
				ContainsDairy:          result.ValidIngredientContainsDairy,
				ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
				AnimalFlesh:            result.ValidIngredientAnimalFlesh,
				IsStarch:               result.ValidIngredientIsStarch,
				IsProtein:              result.ValidIngredientIsProtein,
				IsGrain:                result.ValidIngredientIsGrain,
				IsFruit:                result.ValidIngredientIsFruit,
				IsSalt:                 result.ValidIngredientIsSalt,
				IsFat:                  result.ValidIngredientIsFat,
				IsAcid:                 result.ValidIngredientIsAcid,
				IsHeat:                 result.ValidIngredientIsHeat,
			},
		})

		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// CreateIngredientPreference creates a user ingredient preference in the database.
func (q *Querier) CreateIngredientPreference(ctx context.Context, input *types.IngredientPreferenceDatabaseCreationInput) ([]*types.IngredientPreference, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientIDKey, input.ValidIngredientID)

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

	output := []*types.IngredientPreference{}
	for _, validIngredientID := range validIngredientIDs {
		l := logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
		if validIngredientID == "" {
			continue
		}

		id := identifiers.New()
		tracing.AttachToSpan(span, keys.IngredientPreferenceIDKey, id)

		// create the user ingredient preference.
		if err = q.generatedQuerier.CreateUserIngredientPreference(ctx, tx, &generated.CreateUserIngredientPreferenceParams{
			ID:            id,
			Ingredient:    validIngredientID,
			Notes:         input.Notes,
			BelongsToUser: input.BelongsToUser,
			Rating:        int16(input.Rating),
			Allergy:       input.Allergy,
		}); err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(err, l, span, "performing user ingredient preference creation query")
		}

		x := &types.IngredientPreference{
			ID:            id,
			Rating:        input.Rating,
			Notes:         input.Notes,
			Allergy:       input.Allergy,
			BelongsToUser: input.BelongsToUser,
			Ingredient:    types.ValidIngredient{ID: input.ValidIngredientID},
			CreatedAt:     q.currentTime(),
		}

		l.Info("user ingredient preference created")

		output = append(output, x)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return output, nil
}

// UpdateIngredientPreference updates a particular user ingredient preference.
func (q *Querier) UpdateIngredientPreference(ctx context.Context, updated *types.IngredientPreference) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.IngredientPreferenceIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.IngredientPreferenceIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateUserIngredientPreference(ctx, q.db, &generated.UpdateUserIngredientPreferenceParams{
		Ingredient:    updated.Ingredient.ID,
		Notes:         updated.Notes,
		ID:            updated.ID,
		BelongsToUser: updated.BelongsToUser,
		Rating:        int16(updated.Rating),
		Allergy:       updated.Allergy,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user ingredient preference")
	}

	logger.Info("user ingredient preference updated")

	return nil
}

// ArchiveIngredientPreference archives a user ingredient preference from the database by its ID.
func (q *Querier) ArchiveIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if userIngredientPreferenceID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.IngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachToSpan(span, keys.IngredientPreferenceIDKey, userIngredientPreferenceID)

	if _, err := q.generatedQuerier.ArchiveUserIngredientPreference(ctx, q.db, &generated.ArchiveUserIngredientPreferenceParams{
		ID:            userIngredientPreferenceID,
		BelongsToUser: userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving user ingredient preference")
	}

	logger.Info("user ingredient preference archived")

	return nil
}
