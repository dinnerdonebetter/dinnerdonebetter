package serverimpl

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (s *Server) ArchiveValidIngredient(ctx context.Context, request *messages.ArchiveValidIngredientRequest) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	validIngredient, err := s.dataManager.GetValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, err
	}

	if err = s.dataManager.ArchiveValidIngredient(ctx, request.ValidIngredientID); err != nil {
		return nil, err
	}

	output := converters.ConvertValidIngredientToProtobuf(validIngredient)

	return output, nil
}

func (s *Server) ArchiveValidIngredientGroup(ctx context.Context, request *messages.ArchiveValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *messages.ArchiveValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientPreparation(ctx context.Context, request *messages.ArchiveValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientState(ctx context.Context, request *messages.ArchiveValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientStateIngredient(ctx context.Context, request *messages.ArchiveValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidInstrument(ctx context.Context, request *messages.ArchiveValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidMeasurementUnit(ctx context.Context, request *messages.ArchiveValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *messages.ArchiveValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparationInstrument(ctx context.Context, request *messages.ArchiveValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparationVessel(ctx context.Context, request *messages.ArchiveValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredient(ctx context.Context, request *messages.GetValidIngredientRequest) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	validIngredient, err := s.dataManager.GetValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting valid ingredient")
	}

	output := converters.ConvertValidIngredientToProtobuf(validIngredient)

	return output, nil
}

func (s *Server) GetValidIngredientGroup(ctx context.Context, request *messages.GetValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientGroups(ctx context.Context, request *messages.GetValidIngredientGroupsRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnits(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByIngredientRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparations(ctx context.Context, request *messages.GetValidIngredientPreparationsRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *messages.GetValidIngredientPreparationsByIngredientRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationsByPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientState(ctx context.Context, request *messages.GetValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredients(ctx context.Context, request *messages.GetValidIngredientStateIngredientsRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientStateRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStates(ctx context.Context, request *messages.GetValidIngredientStatesRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidInstrument(ctx context.Context, request *messages.GetValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidInstruments(ctx context.Context, request *messages.GetValidInstrumentsRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnit(ctx context.Context, request *messages.GetValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversion(ctx context.Context, request *messages.GetValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsFromUnitRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsToUnitRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnits(ctx context.Context, request *messages.GetValidMeasurementUnitsRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstruments(ctx context.Context, request *messages.GetValidPreparationInstrumentsRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentsByInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *messages.GetValidPreparationInstrumentsByPreparationRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVessel(ctx context.Context, request *messages.GetValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVessels(ctx context.Context, request *messages.GetValidPreparationVesselsRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVesselsByPreparation(ctx context.Context, request *messages.GetValidPreparationVesselsByPreparationRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVesselsByVessel(ctx context.Context, request *messages.GetValidPreparationVesselsByVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredient(ctx context.Context, input *messages.ValidIngredientCreationRequestInput) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	created, err := s.dataManager.CreateValidIngredient(ctx, &types.ValidIngredientDatabaseCreationInput{
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Min: input.StorageTemperatureInCelsius.Min,
			Max: input.StorageTemperatureInCelsius.Max,
		},
		ID:                     identifiers.New(),
		Warning:                input.Warning,
		IconPath:               input.IconPath,
		PluralName:             input.PluralName,
		StorageInstructions:    input.StorageInstructions,
		Name:                   input.Name,
		Description:            input.Description,
		Slug:                   input.Slug,
		ShoppingSuggestions:    input.ShoppingSuggestions,
		ContainsFish:           input.ContainsFish,
		ContainsShellfish:      input.ContainsShellfish,
		AnimalFlesh:            input.AnimalFlesh,
		ContainsEgg:            input.ContainsEgg,
		IsLiquid:               input.IsLiquid,
		ContainsSoy:            input.ContainsSoy,
		ContainsPeanut:         input.ContainsPeanut,
		AnimalDerived:          input.AnimalDerived,
		RestrictToPreparations: input.RestrictToPreparations,
		ContainsDairy:          input.ContainsDairy,
		ContainsSesame:         input.ContainsSesame,
		ContainsTreeNut:        input.ContainsTreeNut,
		ContainsWheat:          input.ContainsWheat,
		ContainsAlcohol:        input.ContainsAlcohol,
		ContainsGluten:         input.ContainsGluten,
		IsStarch:               input.IsStarch,
		IsProtein:              input.IsProtein,
		IsGrain:                input.IsGrain,
		IsFruit:                input.IsFruit,
		IsSalt:                 input.IsSalt,
		IsFat:                  input.IsFat,
		IsAcid:                 input.IsAcid,
		IsHeat:                 input.IsHeat,
	})
	if err != nil {
		return nil, err
	}

	output := converters.ConvertValidIngredientToProtobuf(created)

	return output, nil
}

func (s *Server) CreateValidIngredientGroup(ctx context.Context, input *messages.ValidIngredientGroupCreationRequestInput) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientMeasurementUnit(ctx context.Context, input *messages.ValidIngredientMeasurementUnitCreationRequestInput) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientPreparation(ctx context.Context, input *messages.ValidIngredientPreparationCreationRequestInput) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientState(ctx context.Context, input *messages.ValidIngredientStateCreationRequestInput) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientStateIngredient(ctx context.Context, input *messages.ValidIngredientStateIngredientCreationRequestInput) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidInstrument(ctx context.Context, input *messages.ValidInstrumentCreationRequestInput) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidMeasurementUnit(ctx context.Context, input *messages.ValidMeasurementUnitCreationRequestInput) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidMeasurementUnitConversion(ctx context.Context, input *messages.ValidMeasurementUnitConversionCreationRequestInput) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidPreparation(ctx context.Context, input *messages.ValidPreparationCreationRequestInput) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidPreparationInstrument(ctx context.Context, input *messages.ValidPreparationInstrumentCreationRequestInput) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidPreparationVessel(ctx context.Context, input *messages.ValidPreparationVesselCreationRequestInput) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidVessel(ctx context.Context, input *messages.ValidVesselCreationRequestInput) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidIngredient(ctx context.Context, _ *emptypb.Empty) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	ingredient, err := s.dataManager.GetRandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting random valid ingredient")
	}

	output := converters.ConvertValidIngredientToProtobuf(ingredient)

	return output, nil
}

func (s *Server) GetRandomValidPreparation(ctx context.Context, _ *emptypb.Empty) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidVessel(ctx context.Context, _ *emptypb.Empty) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	ingredient, err := s.dataManager.GetRandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting random valid ingredient")
	}

	output := &messages.ValidVessel{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(ingredient.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(ingredient.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(ingredient.ArchivedAt),
		CapacityUnit: &messages.ValidMeasurementUnit{
			CreatedAt:     converters.ConvertTimeToPBTimestamp(ingredient.CapacityUnit.CreatedAt),
			LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(ingredient.CapacityUnit.LastUpdatedAt),
			ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(ingredient.CapacityUnit.ArchivedAt),
			Name:          ingredient.CapacityUnit.Name,
			IconPath:      ingredient.CapacityUnit.IconPath,
			ID:            ingredient.CapacityUnit.ID,
			Description:   ingredient.CapacityUnit.Description,
			PluralName:    ingredient.CapacityUnit.PluralName,
			Slug:          ingredient.CapacityUnit.Slug,
			Volumetric:    ingredient.CapacityUnit.Volumetric,
			Universal:     ingredient.CapacityUnit.Universal,
			Metric:        ingredient.CapacityUnit.Metric,
			Imperial:      ingredient.CapacityUnit.Imperial,
		},
		IconPath:                       ingredient.IconPath,
		PluralName:                     ingredient.PluralName,
		Description:                    ingredient.Description,
		Name:                           ingredient.Name,
		Slug:                           ingredient.Slug,
		Shape:                          ingredient.Shape,
		ID:                             ingredient.ID,
		WidthInMillimeters:             ingredient.WidthInMillimeters,
		LengthInMillimeters:            ingredient.LengthInMillimeters,
		HeightInMillimeters:            ingredient.HeightInMillimeters,
		Capacity:                       ingredient.Capacity,
		IncludeInGeneratedInstructions: ingredient.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          ingredient.DisplayInSummaryLists,
		UsableForStorage:               ingredient.UsableForStorage,
	}

	return output, nil
}

func (s *Server) SearchForValidIngredientGroups(ctx context.Context, request *messages.SearchForValidIngredientGroupsRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidIngredientStates(ctx context.Context, request *messages.SearchForValidIngredientStatesRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidIngredients(ctx context.Context, request *messages.SearchForValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidInstruments(ctx context.Context, request *messages.SearchForValidInstrumentsRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidMeasurementUnits(ctx context.Context, request *messages.SearchForValidMeasurementUnitsRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchValidIngredientsByPreparation(ctx context.Context, request *messages.SearchValidIngredientsByPreparationRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *messages.SearchValidMeasurementUnitsByIngredientRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredient(ctx context.Context, request *messages.UpdateValidIngredientRequest) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	updated := converters.ConvertUpdateValidIngredientRequestToValidIngredient(request)

	if err := s.dataManager.UpdateValidIngredient(ctx, updated); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Server) UpdateValidIngredientGroup(ctx context.Context, request *messages.UpdateValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *messages.UpdateValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientPreparation(ctx context.Context, request *messages.UpdateValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientState(ctx context.Context, request *messages.UpdateValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientStateIngredient(ctx context.Context, request *messages.UpdateValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidInstrument(ctx context.Context, request *messages.UpdateValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidMeasurementUnit(ctx context.Context, request *messages.UpdateValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidMeasurementUnitConversion(ctx context.Context, request *messages.UpdateValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparationInstrument(ctx context.Context, request *messages.UpdateValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparationVessel(ctx context.Context, request *messages.UpdateValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
