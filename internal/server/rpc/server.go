package rpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/internal/proto"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	tracer                     tracing.Tracer
	logger                     logging.Logger
	validIngredientDataManager types.ValidIngredientDataManager
}

// ProvideRPCServer builds a new Server instance.
func ProvideRPCServer(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	validIngredientDataManager types.ValidIngredientDataManager,
) (*Server, error) {
	srv := &Server{
		tracer:                     tracing.NewTracer(tracerProvider.Tracer("rpc-server")),
		logger:                     logger.WithName("rpc-server"),
		validIngredientDataManager: validIngredientDataManager,
	}

	return srv, nil
}

func (s *Server) GetRandomValidIngredient(ctx context.Context, _ *empty.Empty) (*proto.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	ingredient, err := s.validIngredientDataManager.GetRandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, s.logger, span, "failed to get random valid ingredient")
	}

	x := &proto.ValidIngredient{
		CreatedAt:                timestamppb.New(ingredient.CreatedAt),
		IconPath:                 ingredient.IconPath,
		Warning:                  ingredient.Warning,
		PluralName:               ingredient.PluralName,
		StorageInstructions:      ingredient.StorageInstructions,
		Name:                     ingredient.Name,
		ID:                       ingredient.ID,
		Description:              ingredient.Description,
		Slug:                     ingredient.Slug,
		ShoppingSuggestions:      ingredient.ShoppingSuggestions,
		ContainsShellfish:        ingredient.ContainsShellfish,
		IsMeasuredVolumetrically: ingredient.IsMeasuredVolumetrically,
		IsLiquid:                 ingredient.IsLiquid,
		ContainsPeanut:           ingredient.ContainsPeanut,
		ContainsTreeNut:          ingredient.ContainsTreeNut,
		ContainsEgg:              ingredient.ContainsEgg,
		ContainsWheat:            ingredient.ContainsWheat,
		ContainsSoy:              ingredient.ContainsSoy,
		AnimalDerived:            ingredient.AnimalDerived,
		RestrictToPreparations:   ingredient.RestrictToPreparations,
		ContainsSesame:           ingredient.ContainsSesame,
		ContainsFish:             ingredient.ContainsFish,
		ContainsGluten:           ingredient.ContainsGluten,
		ContainsDairy:            ingredient.ContainsDairy,
		ContainsAlcohol:          ingredient.ContainsAlcohol,
		AnimalFlesh:              ingredient.AnimalFlesh,
		IsStarch:                 ingredient.IsStarch,
		IsProtein:                ingredient.IsProtein,
		IsGrain:                  ingredient.IsGrain,
		IsFruit:                  ingredient.IsFruit,
		IsSalt:                   ingredient.IsSalt,
		IsFat:                    ingredient.IsFat,
		IsAcid:                   ingredient.IsAcid,
		IsHeat:                   ingredient.IsHeat,
	}

	if ingredient.LastUpdatedAt != nil {
		x.LastUpdatedAt = timestamppb.New(*ingredient.LastUpdatedAt)
	}

	if ingredient.ArchivedAt != nil {
		x.ArchivedAt = timestamppb.New(*ingredient.ArchivedAt)
	}

	if ingredient.MaximumIdealStorageTemperatureInCelsius != nil {
		x.MaximumIdealStorageTemperatureInCelsius = pointers.Pointer(float64(*ingredient.MaximumIdealStorageTemperatureInCelsius))
	}

	if ingredient.MinimumIdealStorageTemperatureInCelsius != nil {
		x.MinimumIdealStorageTemperatureInCelsius = pointers.Pointer(float64(*ingredient.MinimumIdealStorageTemperatureInCelsius))
	}

	return x, nil
}
