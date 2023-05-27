package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/internal/proto/types"
)

// Server is our primary GRPC implementation.
type Server struct {
	tracer      tracing.Tracer
	logger      logging.Logger
	dataManager database.DataManager
}

func (s *Server) GetRandomValidIngredient(ctx context.Context, _ *empty.Empty) (*types.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	result, err := s.dataManager.GetRandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, s.logger, span, "getting random valid ingredient")
	}

	converted := &types.ValidIngredient{
		IconPath:                 result.IconPath,
		Warning:                  result.Warning,
		PluralName:               result.PluralName,
		StorageInstructions:      result.StorageInstructions,
		Name:                     result.Name,
		ID:                       result.ID,
		Description:              result.Description,
		Slug:                     result.Slug,
		ShoppingSuggestions:      result.ShoppingSuggestions,
		ContainsShellfish:        result.ContainsShellfish,
		IsMeasuredVolumetrically: result.IsMeasuredVolumetrically,
		IsLiquid:                 result.IsLiquid,
		ContainsPeanut:           result.ContainsPeanut,
		ContainsTreeNut:          result.ContainsTreeNut,
		ContainsEgg:              result.ContainsEgg,
		ContainsWheat:            result.ContainsWheat,
		ContainsSoy:              result.ContainsSoy,
		AnimalDerived:            result.AnimalDerived,
		RestrictToPreparations:   result.RestrictToPreparations,
		ContainsSesame:           result.ContainsSesame,
		ContainsFish:             result.ContainsFish,
		ContainsGluten:           result.ContainsGluten,
		ContainsDairy:            result.ContainsDairy,
		ContainsAlcohol:          result.ContainsAlcohol,
		AnimalFlesh:              result.AnimalFlesh,
		IsStarch:                 result.IsStarch,
		IsProtein:                result.IsProtein,
		IsGrain:                  result.IsGrain,
		IsFruit:                  result.IsFruit,
		IsSalt:                   result.IsSalt,
		IsFat:                    result.IsFat,
		IsAcid:                   result.IsAcid,
		IsHeat:                   result.IsHeat,
		CreatedAt:                timestamppb.New(result.CreatedAt),
	}

	if result.MaximumIdealStorageTemperatureInCelsius != nil {
		converted.MaximumIdealStorageTemperatureInCelsius = pointers.Pointer(float64(*result.MaximumIdealStorageTemperatureInCelsius))
	}

	if result.MinimumIdealStorageTemperatureInCelsius != nil {
		converted.MinimumIdealStorageTemperatureInCelsius = pointers.Pointer(float64(*result.MinimumIdealStorageTemperatureInCelsius))
	}

	if result.LastUpdatedAt != nil {
		converted.LastUpdatedAt = timestamppb.New(*result.LastUpdatedAt)
	}

	if result.ArchivedAt != nil {
		converted.ArchivedAt = timestamppb.New(*result.ArchivedAt)
	}

	return converted, nil
}
