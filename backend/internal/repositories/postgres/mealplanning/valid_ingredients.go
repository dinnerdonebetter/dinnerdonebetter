package mealplanning

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ mealplanning.ValidIngredientDataManager = (*repository)(nil)
)

// ValidIngredientExists fetches whether a valid ingredient exists from the database.
func (r *repository) ValidIngredientExists(ctx context.Context, validIngredientID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validIngredientID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	result, err := r.generatedQuerier.CheckValidIngredientExistence(ctx, r.db, validIngredientID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient existence check")
	}

	return result, nil
}

// GetValidIngredient fetches a valid ingredient from the database.
func (r *repository) GetValidIngredient(ctx context.Context, validIngredientID string) (*mealplanning.ValidIngredient, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validIngredientID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	result, err := r.generatedQuerier.GetValidIngredient(ctx, r.db, validIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient")
	}

	validIngredient := &mealplanning.ValidIngredient{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
			Min: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
		},
		IconPath:               result.IconPath,
		Warning:                result.Warning,
		PluralName:             result.PluralName,
		StorageInstructions:    result.StorageInstructions,
		Name:                   result.Name,
		ID:                     result.ID,
		Description:            result.Description,
		Slug:                   result.Slug,
		ShoppingSuggestions:    result.ShoppingSuggestions,
		ContainsShellfish:      result.ContainsShellfish,
		IsLiquid:               database.BoolFromNullBool(result.IsLiquid),
		ContainsPeanut:         result.ContainsPeanut,
		ContainsTreeNut:        result.ContainsTreeNut,
		ContainsEgg:            result.ContainsEgg,
		ContainsWheat:          result.ContainsWheat,
		ContainsSoy:            result.ContainsSoy,
		AnimalDerived:          result.AnimalDerived,
		RestrictToPreparations: result.RestrictToPreparations,
		ContainsSesame:         result.ContainsSesame,
		ContainsFish:           result.ContainsFish,
		ContainsGluten:         result.ContainsGluten,
		ContainsDairy:          result.ContainsDairy,
		ContainsAlcohol:        result.ContainsAlcohol,
		AnimalFlesh:            result.AnimalFlesh,
		IsStarch:               result.IsStarch,
		IsProtein:              result.IsProtein,
		IsGrain:                result.IsGrain,
		IsFruit:                result.IsFruit,
		IsSalt:                 result.IsSalt,
		IsFat:                  result.IsFat,
		IsAcid:                 result.IsAcid,
		IsHeat:                 result.IsHeat,
	}

	return validIngredient, nil
}

// GetRandomValidIngredient fetches a valid ingredient from the database.
func (r *repository) GetRandomValidIngredient(ctx context.Context) (*mealplanning.ValidIngredient, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	result, err := r.generatedQuerier.GetRandomValidIngredient(ctx, r.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching random valid ingredient")
	}

	validIngredient := &mealplanning.ValidIngredient{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
			Min: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
		},
		IconPath:               result.IconPath,
		Warning:                result.Warning,
		PluralName:             result.PluralName,
		StorageInstructions:    result.StorageInstructions,
		Name:                   result.Name,
		ID:                     result.ID,
		Description:            result.Description,
		Slug:                   result.Slug,
		ShoppingSuggestions:    result.ShoppingSuggestions,
		ContainsShellfish:      result.ContainsShellfish,
		IsLiquid:               database.BoolFromNullBool(result.IsLiquid),
		ContainsPeanut:         result.ContainsPeanut,
		ContainsTreeNut:        result.ContainsTreeNut,
		ContainsEgg:            result.ContainsEgg,
		ContainsWheat:          result.ContainsWheat,
		ContainsSoy:            result.ContainsSoy,
		AnimalDerived:          result.AnimalDerived,
		RestrictToPreparations: result.RestrictToPreparations,
		ContainsSesame:         result.ContainsSesame,
		ContainsFish:           result.ContainsFish,
		ContainsGluten:         result.ContainsGluten,
		ContainsDairy:          result.ContainsDairy,
		ContainsAlcohol:        result.ContainsAlcohol,
		AnimalFlesh:            result.AnimalFlesh,
		IsStarch:               result.IsStarch,
		IsProtein:              result.IsProtein,
		IsGrain:                result.IsGrain,
		IsFruit:                result.IsFruit,
		IsSalt:                 result.IsSalt,
		IsFat:                  result.IsFat,
		IsAcid:                 result.IsAcid,
		IsHeat:                 result.IsHeat,
	}

	return validIngredient, nil
}

// SearchForValidIngredients fetches a valid ingredient from the database.
func (r *repository) SearchForValidIngredients(ctx context.Context, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.ValidIngredient], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if query == "" {
		return nil, database.ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, query)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	x := &filtering.QueryFilteredResult[mealplanning.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.SearchForValidIngredients(ctx, r.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient")
	}

	for _, result := range results {
		validIngredient := &mealplanning.ValidIngredient{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.IconPath,
			Warning:                result.Warning,
			PluralName:             result.PluralName,
			StorageInstructions:    result.StorageInstructions,
			Name:                   result.Name,
			ID:                     result.ID,
			Description:            result.Description,
			Slug:                   result.Slug,
			ShoppingSuggestions:    result.ShoppingSuggestions,
			ContainsShellfish:      result.ContainsShellfish,
			IsLiquid:               database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:         result.ContainsPeanut,
			ContainsTreeNut:        result.ContainsTreeNut,
			ContainsEgg:            result.ContainsEgg,
			ContainsWheat:          result.ContainsWheat,
			ContainsSoy:            result.ContainsSoy,
			AnimalDerived:          result.AnimalDerived,
			RestrictToPreparations: result.RestrictToPreparations,
			ContainsSesame:         result.ContainsSesame,
			ContainsFish:           result.ContainsFish,
			ContainsGluten:         result.ContainsGluten,
			ContainsDairy:          result.ContainsDairy,
			ContainsAlcohol:        result.ContainsAlcohol,
			AnimalFlesh:            result.AnimalFlesh,
			IsStarch:               result.IsStarch,
			IsProtein:              result.IsProtein,
			IsGrain:                result.IsGrain,
			IsFruit:                result.IsFruit,
			IsSalt:                 result.IsSalt,
			IsFat:                  result.IsFat,
			IsAcid:                 result.IsAcid,
			IsHeat:                 result.IsHeat,
		}

		x.Data = append(x.Data, validIngredient)
	}

	return x, nil
}

// SearchForValidIngredientsForPreparation fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (r *repository) SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.ValidIngredient], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if preparationID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, preparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, preparationID)

	if query == "" {
		return nil, database.ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	x = &filtering.QueryFilteredResult[mealplanning.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.SearchValidIngredientsByPreparationAndIngredientName(ctx, r.db, &generated.SearchValidIngredientsByPreparationAndIngredientNameParams{
		ValidPreparationID: preparationID,
		NameQuery:          query,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients for preparation")
	}

	for _, result := range results {
		validIngredient := &mealplanning.ValidIngredient{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.IconPath,
			Warning:                result.Warning,
			PluralName:             result.PluralName,
			StorageInstructions:    result.StorageInstructions,
			Name:                   result.Name,
			ID:                     result.ID,
			Description:            result.Description,
			Slug:                   result.Slug,
			ShoppingSuggestions:    result.ShoppingSuggestions,
			ContainsShellfish:      result.ContainsShellfish,
			IsLiquid:               database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:         result.ContainsPeanut,
			ContainsTreeNut:        result.ContainsTreeNut,
			ContainsEgg:            result.ContainsEgg,
			ContainsWheat:          result.ContainsWheat,
			ContainsSoy:            result.ContainsSoy,
			AnimalDerived:          result.AnimalDerived,
			RestrictToPreparations: result.RestrictToPreparations,
			ContainsSesame:         result.ContainsSesame,
			ContainsFish:           result.ContainsFish,
			ContainsGluten:         result.ContainsGluten,
			ContainsDairy:          result.ContainsDairy,
			ContainsAlcohol:        result.ContainsAlcohol,
			AnimalFlesh:            result.AnimalFlesh,
			IsStarch:               result.IsStarch,
			IsProtein:              result.IsProtein,
			IsGrain:                result.IsGrain,
			IsFruit:                result.IsFruit,
			IsSalt:                 result.IsSalt,
			IsFat:                  result.IsFat,
			IsAcid:                 result.IsAcid,
			IsHeat:                 result.IsHeat,
		}

		x.Data = append(x.Data, validIngredient)
	}

	return x, nil
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (r *repository) GetValidIngredients(ctx context.Context, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.ValidIngredient], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[mealplanning.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetValidIngredients(ctx, r.db, &generated.GetValidIngredientsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	for _, result := range results {
		validIngredient := &mealplanning.ValidIngredient{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.IconPath,
			Warning:                result.Warning,
			PluralName:             result.PluralName,
			StorageInstructions:    result.StorageInstructions,
			Name:                   result.Name,
			ID:                     result.ID,
			Description:            result.Description,
			Slug:                   result.Slug,
			ShoppingSuggestions:    result.ShoppingSuggestions,
			ContainsShellfish:      result.ContainsShellfish,
			IsLiquid:               database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:         result.ContainsPeanut,
			ContainsTreeNut:        result.ContainsTreeNut,
			ContainsEgg:            result.ContainsEgg,
			ContainsWheat:          result.ContainsWheat,
			ContainsSoy:            result.ContainsSoy,
			AnimalDerived:          result.AnimalDerived,
			RestrictToPreparations: result.RestrictToPreparations,
			ContainsSesame:         result.ContainsSesame,
			ContainsFish:           result.ContainsFish,
			ContainsGluten:         result.ContainsGluten,
			ContainsDairy:          result.ContainsDairy,
			ContainsAlcohol:        result.ContainsAlcohol,
			AnimalFlesh:            result.AnimalFlesh,
			IsStarch:               result.IsStarch,
			IsProtein:              result.IsProtein,
			IsGrain:                result.IsGrain,
			IsFruit:                result.IsFruit,
			IsSalt:                 result.IsSalt,
			IsFat:                  result.IsFat,
			IsAcid:                 result.IsAcid,
			IsHeat:                 result.IsHeat,
		}

		x.Data = append(x.Data, validIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientsWithIDs fetches a list of valid ingredients from the database that meet a particular filter.
func (r *repository) GetValidIngredientsWithIDs(ctx context.Context, ids []string) ([]*mealplanning.ValidIngredient, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if ids == nil {
		return nil, database.ErrEmptyInputProvided
	}

	results, err := r.generatedQuerier.GetValidIngredientsWithIDs(ctx, r.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients id list retrieval query")
	}

	var ingredients []*mealplanning.ValidIngredient
	for _, result := range results {
		validIngredient := &mealplanning.ValidIngredient{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.IconPath,
			Warning:                result.Warning,
			PluralName:             result.PluralName,
			StorageInstructions:    result.StorageInstructions,
			Name:                   result.Name,
			ID:                     result.ID,
			Description:            result.Description,
			Slug:                   result.Slug,
			ShoppingSuggestions:    result.ShoppingSuggestions,
			ContainsShellfish:      result.ContainsShellfish,
			IsLiquid:               database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:         result.ContainsPeanut,
			ContainsTreeNut:        result.ContainsTreeNut,
			ContainsEgg:            result.ContainsEgg,
			ContainsWheat:          result.ContainsWheat,
			ContainsSoy:            result.ContainsSoy,
			AnimalDerived:          result.AnimalDerived,
			RestrictToPreparations: result.RestrictToPreparations,
			ContainsSesame:         result.ContainsSesame,
			ContainsFish:           result.ContainsFish,
			ContainsGluten:         result.ContainsGluten,
			ContainsDairy:          result.ContainsDairy,
			ContainsAlcohol:        result.ContainsAlcohol,
			AnimalFlesh:            result.AnimalFlesh,
			IsStarch:               result.IsStarch,
			IsProtein:              result.IsProtein,
			IsGrain:                result.IsGrain,
			IsFruit:                result.IsFruit,
			IsSalt:                 result.IsSalt,
			IsFat:                  result.IsFat,
			IsAcid:                 result.IsAcid,
			IsHeat:                 result.IsHeat,
		}

		ingredients = append(ingredients, validIngredient)
	}

	return ingredients, nil
}

// GetValidIngredientIDsThatNeedSearchIndexing fetches a list of valid ingredients from the database that meet a particular filter.
func (r *repository) GetValidIngredientIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	results, err := r.generatedQuerier.GetValidIngredientsNeedingIndexing(ctx, r.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid ingredients list retrieval query")
	}

	return results, err
}

// CreateValidIngredient creates a valid ingredient in the database.
func (r *repository) CreateValidIngredient(ctx context.Context, input *mealplanning.ValidIngredientDatabaseCreationInput) (*mealplanning.ValidIngredient, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, input.ID)
	logger := r.logger.WithValue(keys.ValidIngredientIDKey, input.ID)

	// create the valid ingredient.
	if err := r.generatedQuerier.CreateValidIngredient(ctx, r.db, &generated.CreateValidIngredientParams{
		ID:                                      input.ID,
		Name:                                    input.Name,
		Description:                             input.Description,
		Warning:                                 input.Warning,
		ContainsEgg:                             input.ContainsEgg,
		ContainsDairy:                           input.ContainsDairy,
		ContainsPeanut:                          input.ContainsPeanut,
		ContainsTreeNut:                         input.ContainsTreeNut,
		ContainsSoy:                             input.ContainsSoy,
		ContainsWheat:                           input.ContainsWheat,
		ContainsShellfish:                       input.ContainsShellfish,
		ContainsSesame:                          input.ContainsSesame,
		ContainsFish:                            input.ContainsFish,
		ContainsGluten:                          input.ContainsGluten,
		AnimalFlesh:                             input.AnimalFlesh,
		IsLiquid:                                database.NullBoolFromBool(input.IsLiquid),
		IconPath:                                input.IconPath,
		AnimalDerived:                           input.AnimalDerived,
		PluralName:                              input.PluralName,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MaximumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(input.StorageTemperatureInCelsius.Max),
		MinimumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(input.StorageTemperatureInCelsius.Min),
		StorageInstructions:                     input.StorageInstructions,
		Slug:                                    input.Slug,
		ContainsAlcohol:                         input.ContainsAlcohol,
		ShoppingSuggestions:                     input.ShoppingSuggestions,
		IsStarch:                                input.IsStarch,
		IsProtein:                               input.IsProtein,
		IsGrain:                                 input.IsGrain,
		IsFruit:                                 input.IsFruit,
		IsSalt:                                  input.IsSalt,
		IsFat:                                   input.IsFat,
		IsAcid:                                  input.IsAcid,
		IsHeat:                                  input.IsHeat,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient creation query")
	}

	x := &mealplanning.ValidIngredient{
		ID:                     input.ID,
		Name:                   input.Name,
		Description:            input.Description,
		Warning:                input.Warning,
		ContainsEgg:            input.ContainsEgg,
		ContainsDairy:          input.ContainsDairy,
		ContainsPeanut:         input.ContainsPeanut,
		ContainsTreeNut:        input.ContainsTreeNut,
		ContainsSoy:            input.ContainsSoy,
		ContainsWheat:          input.ContainsWheat,
		ContainsShellfish:      input.ContainsShellfish,
		ContainsSesame:         input.ContainsSesame,
		ContainsFish:           input.ContainsFish,
		ContainsGluten:         input.ContainsGluten,
		AnimalFlesh:            input.AnimalFlesh,
		IsLiquid:               input.IsLiquid,
		IconPath:               input.IconPath,
		AnimalDerived:          input.AnimalDerived,
		PluralName:             input.PluralName,
		IsStarch:               input.IsStarch,
		IsProtein:              input.IsProtein,
		IsGrain:                input.IsGrain,
		IsFruit:                input.IsFruit,
		IsSalt:                 input.IsSalt,
		IsFat:                  input.IsFat,
		IsAcid:                 input.IsAcid,
		IsHeat:                 input.IsHeat,
		RestrictToPreparations: input.RestrictToPreparations,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: input.StorageTemperatureInCelsius.Max,
			Min: input.StorageTemperatureInCelsius.Min,
		},
		StorageInstructions: input.StorageInstructions,
		Slug:                input.Slug,
		ContainsAlcohol:     input.ContainsAlcohol,
		ShoppingSuggestions: input.ShoppingSuggestions,
		CreatedAt:           r.CurrentTime(),
	}

	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, x.ID)
	logger.Info("valid ingredient created")

	return x, nil
}

// UpdateValidIngredient updates a particular valid ingredient.
func (r *repository) UpdateValidIngredient(ctx context.Context, updated *mealplanning.ValidIngredient) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.ValidIngredientIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateValidIngredient(ctx, r.db, &generated.UpdateValidIngredientParams{
		Description:                             updated.Description,
		Warning:                                 updated.Warning,
		ID:                                      updated.ID,
		ShoppingSuggestions:                     updated.ShoppingSuggestions,
		Slug:                                    updated.Slug,
		StorageInstructions:                     updated.StorageInstructions,
		Name:                                    updated.Name,
		PluralName:                              updated.PluralName,
		IconPath:                                updated.IconPath,
		MaximumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(updated.StorageTemperatureInCelsius.Max),
		MinimumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(updated.StorageTemperatureInCelsius.Min),
		IsLiquid:                                database.NullBoolFromBool(updated.IsLiquid),
		ContainsWheat:                           updated.ContainsWheat,
		ContainsPeanut:                          updated.ContainsPeanut,
		ContainsGluten:                          updated.ContainsGluten,
		ContainsFish:                            updated.ContainsFish,
		AnimalDerived:                           updated.AnimalDerived,
		ContainsSesame:                          updated.ContainsSesame,
		RestrictToPreparations:                  updated.RestrictToPreparations,
		ContainsShellfish:                       updated.ContainsShellfish,
		ContainsSoy:                             updated.ContainsSoy,
		ContainsTreeNut:                         updated.ContainsTreeNut,
		AnimalFlesh:                             updated.AnimalFlesh,
		ContainsAlcohol:                         updated.ContainsAlcohol,
		ContainsDairy:                           updated.ContainsDairy,
		IsStarch:                                updated.IsStarch,
		IsProtein:                               updated.IsProtein,
		IsGrain:                                 updated.IsGrain,
		IsFruit:                                 updated.IsFruit,
		IsSalt:                                  updated.IsSalt,
		IsFat:                                   updated.IsFat,
		IsAcid:                                  updated.IsAcid,
		IsHeat:                                  updated.IsHeat,
		ContainsEgg:                             updated.ContainsEgg,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient")
	}

	logger.Info("valid ingredient updated")

	return nil
}

// MarkValidIngredientAsIndexed updates a particular valid ingredient's last_indexed_at value.
func (r *repository) MarkValidIngredientAsIndexed(ctx context.Context, validIngredientID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validIngredientID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	if _, err := r.generatedQuerier.UpdateValidIngredientLastIndexedAt(ctx, r.db, validIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid ingredient as indexed")
	}

	logger.Info("valid ingredient marked as indexed")

	return nil
}

// ArchiveValidIngredient archives a valid ingredient from the database by its ID.
func (r *repository) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validIngredientID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	rowsAffected, err := r.generatedQuerier.ArchiveValidIngredient(ctx, r.db, validIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
