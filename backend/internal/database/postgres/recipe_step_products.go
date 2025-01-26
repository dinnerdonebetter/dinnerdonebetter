package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipeStepProductDataManager = (*Querier)(nil)
)

// RecipeStepProductExists fetches whether a recipe step product exists from the database.
func (q *Querier) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	result, err := q.generatedQuerier.CheckRecipeStepProductExistence(ctx, q.db, &generated.CheckRecipeStepProductExistenceParams{
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
func (q *Querier) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	result, err := q.generatedQuerier.GetRecipeStepProduct(ctx, q.db, &generated.GetRecipeStepProductParams{
		RecipeStepID:        recipeStepID,
		RecipeStepProductID: recipeStepProductID,
		RecipeID:            recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step product")
	}

	recipeStepProduct := &types.RecipeStepProduct{
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
		recipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{
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
func (q *Querier) getRecipeStepProductsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := q.generatedQuerier.GetRecipeStepProductsForRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step products for recipe")
	}

	recipeStepProducts := []*types.RecipeStepProduct{}
	for _, result := range results {
		recipeStepProduct := &types.RecipeStepProduct{
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
			recipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{
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
func (q *Querier) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.RecipeStepProduct], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.RecipeStepProduct]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeStepProducts(ctx, q.db, &generated.GetRecipeStepProductsParams{
		RecipeStepID:    recipeStepID,
		RecipeID:        recipeID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step products")
	}

	for _, result := range results {
		recipeStepProduct := &types.RecipeStepProduct{
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
			recipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{
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
func (q *Querier) createRecipeStepProduct(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// create the recipe step product.
	if err := q.generatedQuerier.CreateRecipeStepProduct(ctx, db, &generated.CreateRecipeStepProductParams{
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

	x := &types.RecipeStepProduct{
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
		CreatedAt:              q.currentTime(),
	}

	if input.MeasurementUnitID != nil {
		x.MeasurementUnit = &types.ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, x.ID)

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *Querier) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	return q.createRecipeStepProduct(ctx, q.db, input)
}

// UpdateRecipeStepProduct updates a particular recipe step product.
func (q *Querier) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, updated.ID)

	var measurementUnitID *string
	if updated.MeasurementUnit != nil {
		measurementUnitID = &updated.MeasurementUnit.ID
	}

	if _, err := q.generatedQuerier.UpdateRecipeStepProduct(ctx, q.db, &generated.UpdateRecipeStepProductParams{
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
func (q *Querier) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	if _, err := q.generatedQuerier.ArchiveRecipeStepProduct(ctx, q.db, &generated.ArchiveRecipeStepProductParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepProductID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product archived")

	return nil
}
