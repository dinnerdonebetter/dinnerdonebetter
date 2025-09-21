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
	_ mealplanning.RecipeStepProductDataManager = (*repository)(nil)
)

// RecipeStepProductExists fetches whether a recipe step product exists from the database.
func (r *repository) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	result, err := r.generatedQuerier.CheckRecipeStepProductExistence(ctx, r.db, &generated.CheckRecipeStepProductExistenceParams{
		RecipeStepID:        recipeStepID,
		RecipeStepProductID: recipeStepProductID,
		RecipeID:            recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step product existence check")
	}

	return result, nil
}

// GetRecipeStepProduct fetches a recipe step product from the database.
func (r *repository) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*mealplanning.RecipeStepProduct, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	result, err := r.generatedQuerier.GetRecipeStepProduct(ctx, r.db, &generated.GetRecipeStepProductParams{
		RecipeStepID:        recipeStepID,
		RecipeStepProductID: recipeStepProductID,
		RecipeID:            recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step product")
	}

	recipeStepProduct := &mealplanning.RecipeStepProduct{
		CreatedAt: result.CreatedAt,
		Quantity: types.OptionalFloat32Range{
			Max: database.Float32PointerFromNullString(result.MaximumQuantityValue),
			Min: database.Float32PointerFromNullString(result.MinimumQuantityValue),
		},
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: database.Float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
			Min: database.Float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
		},
		StorageDurationInSeconds: types.OptionalUint32Range{
			Max: database.Uint32PointerFromNullInt32(result.MaximumStorageDurationInSeconds),
		},
		ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
		MeasurementUnit:        nil,
		ContainedInVesselIndex: database.Uint16PointerFromNullInt32(result.ContainedInVesselIndex),
		Name:                   result.Name,
		BelongsToRecipeStep:    result.BelongsToRecipeStep,
		Type:                   string(result.Type),
		ID:                     result.ID,
		StorageInstructions:    result.StorageInstructions,
		QuantityNotes:          result.QuantityNotes,
		Index:                  uint16(result.Index),
		IsWaste:                result.IsWaste,
		IsLiquid:               result.IsLiquid,
		Compostable:            result.Compostable,
	}

	if result.ValidMeasurementUnitID.Valid && result.ValidMeasurementUnitID.String != "" {
		recipeStepProduct.MeasurementUnit = &mealplanning.ValidMeasurementUnit{
			CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
			Name:          result.ValidMeasurementUnitName.String,
			IconPath:      result.ValidMeasurementUnitIconPath.String,
			ID:            result.ValidMeasurementUnitID.String,
			Description:   result.ValidMeasurementUnitDescription.String,
			PluralName:    result.ValidMeasurementUnitPluralName.String,
			Slug:          result.ValidMeasurementUnitSlug.String,
			Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
			Universal:     result.ValidMeasurementUnitUniversal.Bool,
			Metric:        result.ValidMeasurementUnitMetric.Bool,
			Imperial:      result.ValidMeasurementUnitImperial.Bool,
		}
	}

	return recipeStepProduct, nil
}

// getRecipeStepProductsForRecipe fetches a list of recipe step products from the database that meet a particular filter.
func (r *repository) getRecipeStepProductsForRecipe(ctx context.Context, recipeID string) ([]*mealplanning.RecipeStepProduct, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := r.generatedQuerier.GetRecipeStepProductsForRecipe(ctx, r.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step products for recipe")
	}

	recipeStepProducts := []*mealplanning.RecipeStepProduct{}
	for _, result := range results {
		recipeStepProduct := &mealplanning.RecipeStepProduct{
			CreatedAt: result.CreatedAt,
			Quantity: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumQuantityValue),
				Min: database.Float32PointerFromNullString(result.MinimumQuantityValue),
			},
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
			},
			StorageDurationInSeconds: types.OptionalUint32Range{
				Max: database.Uint32PointerFromNullInt32(result.MaximumStorageDurationInSeconds),
			},
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			MeasurementUnit:        nil,
			ContainedInVesselIndex: database.Uint16PointerFromNullInt32(result.ContainedInVesselIndex),
			Name:                   result.Name,
			BelongsToRecipeStep:    result.BelongsToRecipeStep,
			Type:                   string(result.Type),
			ID:                     result.ID,
			StorageInstructions:    result.StorageInstructions,
			QuantityNotes:          result.QuantityNotes,
			Index:                  uint16(result.Index),
			IsWaste:                result.IsWaste,
			IsLiquid:               result.IsLiquid,
			Compostable:            result.Compostable,
		}

		if result.ValidMeasurementUnitID.Valid && result.ValidMeasurementUnitID.String != "" {
			recipeStepProduct.MeasurementUnit = &mealplanning.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName.String,
				IconPath:      result.ValidMeasurementUnitIconPath.String,
				ID:            result.ValidMeasurementUnitID.String,
				Description:   result.ValidMeasurementUnitDescription.String,
				PluralName:    result.ValidMeasurementUnitPluralName.String,
				Slug:          result.ValidMeasurementUnitSlug.String,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal.Bool,
				Metric:        result.ValidMeasurementUnitMetric.Bool,
				Imperial:      result.ValidMeasurementUnitImperial.Bool,
			}
		}

		recipeStepProducts = append(recipeStepProducts, recipeStepProduct)
	}

	return recipeStepProducts, nil
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (r *repository) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.RecipeStepProduct], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[mealplanning.RecipeStepProduct]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetRecipeStepProducts(ctx, r.db, &generated.GetRecipeStepProductsParams{
		RecipeStepID:    recipeStepID,
		RecipeID:        recipeID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step products")
	}

	for _, result := range results {
		recipeStepProduct := &mealplanning.RecipeStepProduct{
			CreatedAt: result.CreatedAt,
			Quantity: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumQuantityValue),
				Min: database.Float32PointerFromNullString(result.MinimumQuantityValue),
			},
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
			},
			StorageDurationInSeconds: types.OptionalUint32Range{
				Max: database.Uint32PointerFromNullInt32(result.MaximumStorageDurationInSeconds),
			},
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			MeasurementUnit:        nil,
			ContainedInVesselIndex: database.Uint16PointerFromNullInt32(result.ContainedInVesselIndex),
			Name:                   result.Name,
			BelongsToRecipeStep:    result.BelongsToRecipeStep,
			Type:                   string(result.Type),
			ID:                     result.ID,
			StorageInstructions:    result.StorageInstructions,
			QuantityNotes:          result.QuantityNotes,
			Index:                  uint16(result.Index),
			IsWaste:                result.IsWaste,
			IsLiquid:               result.IsLiquid,
			Compostable:            result.Compostable,
		}

		if result.ValidMeasurementUnitID.Valid && result.ValidMeasurementUnitID.String != "" {
			recipeStepProduct.MeasurementUnit = &mealplanning.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName.String,
				IconPath:      result.ValidMeasurementUnitIconPath.String,
				ID:            result.ValidMeasurementUnitID.String,
				Description:   result.ValidMeasurementUnitDescription.String,
				PluralName:    result.ValidMeasurementUnitPluralName.String,
				Slug:          result.ValidMeasurementUnitSlug.String,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal.Bool,
				Metric:        result.ValidMeasurementUnitMetric.Bool,
				Imperial:      result.ValidMeasurementUnitImperial.Bool,
			}
		}

		x.Data = append(x.Data, recipeStepProduct)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (r *repository) createRecipeStepProduct(ctx context.Context, db database.SQLQueryExecutor, input *mealplanning.RecipeStepProductDatabaseCreationInput) (*mealplanning.RecipeStepProduct, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	// create the recipe step product.
	if err := r.generatedQuerier.CreateRecipeStepProduct(ctx, db, &generated.CreateRecipeStepProductParams{
		QuantityNotes:                      input.QuantityNotes,
		Name:                               input.Name,
		Type:                               generated.RecipeStepProductType(input.Type),
		BelongsToRecipeStep:                input.BelongsToRecipeStep,
		ID:                                 input.ID,
		StorageInstructions:                input.StorageInstructions,
		MinimumQuantityValue:               database.NullStringFromFloat32Pointer(input.Quantity.Min),
		MinimumStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(input.StorageTemperatureInCelsius.Min),
		MaximumStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(input.StorageTemperatureInCelsius.Max),
		MaximumQuantityValue:               database.NullStringFromFloat32Pointer(input.Quantity.Max),
		MeasurementUnit:                    database.NullStringFromStringPointer(input.MeasurementUnitID),
		MaximumStorageDurationInSeconds:    database.NullInt32FromUint32Pointer(input.StorageDurationInSeconds.Max),
		ContainedInVesselIndex:             database.NullInt32FromUint16Pointer(input.ContainedInVesselIndex),
		Index:                              int32(input.Index),
		Compostable:                        input.Compostable,
		IsLiquid:                           input.IsLiquid,
		IsWaste:                            input.IsWaste,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step product creation query")
	}

	x := &mealplanning.RecipeStepProduct{
		ID:                  input.ID,
		Name:                input.Name,
		Type:                input.Type,
		QuantityNotes:       input.QuantityNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Compostable:         input.Compostable,
		Quantity: types.OptionalFloat32Range{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: input.StorageTemperatureInCelsius.Max,
			Min: input.StorageTemperatureInCelsius.Min,
		},
		StorageDurationInSeconds: types.OptionalUint32Range{
			Max: input.StorageDurationInSeconds.Max,
		},
		StorageInstructions:    input.StorageInstructions,
		IsLiquid:               input.IsLiquid,
		IsWaste:                input.IsWaste,
		Index:                  input.Index,
		ContainedInVesselIndex: input.ContainedInVesselIndex,
		CreatedAt:              r.CurrentTime(),
	}

	if input.MeasurementUnitID != nil {
		x.MeasurementUnit = &mealplanning.ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, x.ID)

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (r *repository) CreateRecipeStepProduct(ctx context.Context, input *mealplanning.RecipeStepProductDatabaseCreationInput) (*mealplanning.RecipeStepProduct, error) {
	return r.createRecipeStepProduct(ctx, r.db, input)
}

// UpdateRecipeStepProduct updates a particular recipe step product.
func (r *repository) UpdateRecipeStepProduct(ctx context.Context, updated *mealplanning.RecipeStepProduct) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, updated.ID)

	var measurementUnitID *string
	if updated.MeasurementUnit != nil {
		measurementUnitID = &updated.MeasurementUnit.ID
	}

	if _, err := r.generatedQuerier.UpdateRecipeStepProduct(ctx, r.db, &generated.UpdateRecipeStepProductParams{
		Name:                               updated.Name,
		Type:                               generated.RecipeStepProductType(updated.Type),
		MeasurementUnit:                    database.NullStringFromStringPointer(measurementUnitID),
		MinimumQuantityValue:               database.NullStringFromFloat32Pointer(updated.Quantity.Min),
		MaximumQuantityValue:               database.NullStringFromFloat32Pointer(updated.Quantity.Max),
		QuantityNotes:                      updated.QuantityNotes,
		Compostable:                        updated.Compostable,
		MaximumStorageDurationInSeconds:    database.NullInt32FromUint32Pointer(updated.StorageDurationInSeconds.Max),
		MinimumStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(updated.StorageTemperatureInCelsius.Min),
		MaximumStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(updated.StorageTemperatureInCelsius.Max),
		StorageInstructions:                updated.StorageInstructions,
		IsLiquid:                           updated.IsLiquid,
		IsWaste:                            updated.IsWaste,
		Index:                              int32(updated.Index),
		ContainedInVesselIndex:             database.NullInt32FromUint16Pointer(updated.ContainedInVesselIndex),
		BelongsToRecipeStep:                updated.BelongsToRecipeStep,
		ID:                                 updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product updated")

	return nil
}

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID.
func (r *repository) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeStepID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	rowsAffected, err := r.generatedQuerier.ArchiveRecipeStepProduct(ctx, r.db, &generated.ArchiveRecipeStepProductParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepProductID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
