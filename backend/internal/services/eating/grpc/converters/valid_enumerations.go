package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

func ConvertGRPCCreateValidIngredientRequestToValidIngredientCreationRequestInput(request *messages.CreateValidIngredientRequest) *types.ValidIngredientCreationRequestInput {
	return &types.ValidIngredientCreationRequestInput{
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: request.StorageTemperatureInCelsius.Max,
			Min: request.StorageTemperatureInCelsius.Min,
		},
		Warning:                request.Warning,
		IconPath:               request.IconPath,
		PluralName:             request.PluralName,
		StorageInstructions:    request.StorageInstructions,
		Name:                   request.Name,
		Description:            request.Description,
		Slug:                   request.Slug,
		ShoppingSuggestions:    request.ShoppingSuggestions,
		ContainsFish:           request.ContainsFish,
		ContainsShellfish:      request.ContainsShellfish,
		AnimalFlesh:            request.AnimalFlesh,
		ContainsEgg:            request.ContainsEgg,
		IsLiquid:               request.IsLiquid,
		ContainsSoy:            request.ContainsSoy,
		ContainsPeanut:         request.ContainsPeanut,
		AnimalDerived:          request.AnimalDerived,
		RestrictToPreparations: request.RestrictToPreparations,
		ContainsDairy:          request.ContainsDairy,
		ContainsSesame:         request.ContainsSesame,
		ContainsTreeNut:        request.ContainsTreeNut,
		ContainsWheat:          request.ContainsWheat,
		ContainsAlcohol:        request.ContainsAlcohol,
		ContainsGluten:         request.ContainsGluten,
		IsStarch:               request.IsStarch,
		IsProtein:              request.IsProtein,
		IsGrain:                request.IsGrain,
		IsFruit:                request.IsFruit,
		IsSalt:                 request.IsSalt,
		IsFat:                  request.IsFat,
		IsAcid:                 request.IsAcid,
		IsHeat:                 request.IsHeat,
	}
}

func ConvertGRPCValidIngredientUpdateRequestInputToValidIngredientUpdateRequestInput(x *messages.ValidIngredientUpdateRequestInput) *types.ValidIngredientUpdateRequestInput {
	return &types.ValidIngredientUpdateRequestInput{
		Name:                   x.Name,
		Description:            x.Description,
		Warning:                x.Warning,
		IconPath:               x.IconPath,
		ContainsDairy:          x.ContainsDairy,
		ContainsPeanut:         x.ContainsPeanut,
		ContainsTreeNut:        x.ContainsTreeNut,
		ContainsEgg:            x.ContainsEgg,
		ContainsWheat:          x.ContainsWheat,
		ContainsShellfish:      x.ContainsShellfish,
		ContainsSesame:         x.ContainsSesame,
		ContainsFish:           x.ContainsFish,
		ContainsGluten:         x.ContainsGluten,
		AnimalFlesh:            x.AnimalFlesh,
		IsLiquid:               x.IsLiquid,
		ContainsSoy:            x.ContainsSoy,
		PluralName:             x.PluralName,
		AnimalDerived:          x.AnimalDerived,
		RestrictToPreparations: x.RestrictToPreparations,
		StorageTemperatureInCelsius: &types.OptionalFloat32Range{
			Max: x.StorageTemperatureInCelsius.Max,
			Min: x.StorageTemperatureInCelsius.Min,
		},
		StorageInstructions: x.StorageInstructions,
		Slug:                x.Slug,
		ContainsAlcohol:     x.ContainsAlcohol,
		ShoppingSuggestions: x.ShoppingSuggestions,
		IsStarch:            x.IsStarch,
		IsProtein:           x.IsProtein,
		IsGrain:             x.IsGrain,
		IsFruit:             x.IsFruit,
		IsSalt:              x.IsSalt,
		IsFat:               x.IsFat,
		IsAcid:              x.IsAcid,
		IsHeat:              x.IsHeat,
	}
}

func ConvertValidIngredientToGRPCValidIngredient(x *types.ValidIngredient) *messages.ValidIngredient {
	return &messages.ValidIngredient{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		StorageTemperatureInCelsius: &messages.OptionalFloat32Range{
			Max: x.StorageTemperatureInCelsius.Max,
			Min: x.StorageTemperatureInCelsius.Min,
		},
		StorageInstructions:    x.StorageInstructions,
		Warning:                x.Warning,
		PluralName:             x.PluralName,
		IconPath:               x.IconPath,
		Name:                   x.Name,
		ID:                     x.ID,
		Description:            x.Description,
		Slug:                   x.Slug,
		ShoppingSuggestions:    x.ShoppingSuggestions,
		ContainsEgg:            x.ContainsEgg,
		ContainsAlcohol:        x.ContainsAlcohol,
		ContainsPeanut:         x.ContainsPeanut,
		ContainsWheat:          x.ContainsWheat,
		ContainsSoy:            x.ContainsSoy,
		AnimalDerived:          x.AnimalDerived,
		RestrictToPreparations: x.RestrictToPreparations,
		ContainsSesame:         x.ContainsSesame,
		ContainsFish:           x.ContainsFish,
		ContainsGluten:         x.ContainsGluten,
		ContainsDairy:          x.ContainsDairy,
		ContainsTreeNut:        x.ContainsTreeNut,
		AnimalFlesh:            x.AnimalFlesh,
		IsStarch:               x.IsStarch,
		IsProtein:              x.IsProtein,
		IsGrain:                x.IsGrain,
		IsFruit:                x.IsFruit,
		IsSalt:                 x.IsSalt,
		IsFat:                  x.IsFat,
		IsAcid:                 x.IsAcid,
		IsHeat:                 x.IsHeat,
		IsLiquid:               x.IsLiquid,
		ContainsShellfish:      x.ContainsShellfish,
	}
}

func ConvertGRPCCreateValidIngredientGroupRequestToValidIngredientGroupCreationRequestInput(request *messages.CreateValidIngredientGroupRequest) *types.ValidIngredientGroupCreationRequestInput {
	members := make([]*types.ValidIngredientGroupMemberCreationRequestInput, len(request.Members))
	for i, member := range request.Members {
		members[i] = &types.ValidIngredientGroupMemberCreationRequestInput{
			ValidIngredientID: member.ValidIngredientID,
		}
	}

	return &types.ValidIngredientGroupCreationRequestInput{
		Name:        request.Name,
		Slug:        request.Slug,
		Description: request.Description,
		Members:     members,
	}
}

func ConvertGRPCValidIngredientGroupUpdateRequestInputToValidIngredientGroupUpdateRequestInput(x *messages.ValidIngredientGroupUpdateRequestInput) *types.ValidIngredientGroupUpdateRequestInput {
	return &types.ValidIngredientGroupUpdateRequestInput{
		Name:        x.Name,
		Slug:        x.Slug,
		Description: x.Description,
	}
}

func ConvertValidIngredientGroupToGRPCValidIngredientGroup(x *types.ValidIngredientGroup) *messages.ValidIngredientGroup {
	members := make([]*messages.ValidIngredientGroupMember, len(x.Members))
	for i, member := range x.Members {
		members[i] = &messages.ValidIngredientGroupMember{
			CreatedAt:       ConvertTimeToPBTimestamp(member.CreatedAt),
			ArchivedAt:      ConvertTimePointerToPBTimestamp(member.ArchivedAt),
			ID:              member.ID,
			BelongsToGroup:  member.BelongsToGroup,
			ValidIngredient: ConvertValidIngredientToGRPCValidIngredient(&member.ValidIngredient),
		}
	}

	return &messages.ValidIngredientGroup{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		ID:            x.ID,
		Name:          x.Name,
		Slug:          x.Slug,
		Description:   x.Description,
		Members:       members,
	}
}

func ConvertGRPCCreateValidIngredientMeasurementUnitRequestToValidIngredientMeasurementUnitCreationRequestInput(request *messages.CreateValidIngredientMeasurementUnitRequest) *types.ValidIngredientMeasurementUnitCreationRequestInput {
	return &types.ValidIngredientMeasurementUnitCreationRequestInput{
		Notes:                  request.Notes,
		ValidMeasurementUnitID: request.ValidMeasurementUnitID,
		ValidIngredientID:      request.ValidIngredientID,
		AllowableQuantity: types.Float32RangeWithOptionalMax{
			Max: request.AllowableQuantity.Max,
			Min: request.AllowableQuantity.Min,
		},
	}
}

func ConvertGRPCValidIngredientMeasurementUnitUpdateRequestInputToValidIngredientMeasurementUnitUpdateRequestInput(x *messages.ValidIngredientMeasurementUnitUpdateRequestInput) *types.ValidIngredientMeasurementUnitUpdateRequestInput {
	return &types.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                  x.Notes,
		ValidMeasurementUnitID: x.ValidMeasurementUnitID,
		ValidIngredientID:      x.ValidIngredientID,
		AllowableQuantity: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(x.AllowableQuantity.Min),
			Max: x.AllowableQuantity.Max,
		},
	}
}

func ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(x *types.ValidIngredientMeasurementUnit) *messages.ValidIngredientMeasurementUnit {
	return &messages.ValidIngredientMeasurementUnit{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		Notes:         x.Notes,
		ID:            x.ID,
		AllowableQuantity: &messages.Float32RangeWithOptionalMax{
			Max: x.AllowableQuantity.Max,
			Min: x.AllowableQuantity.Min,
		},
		MeasurementUnit: ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&x.MeasurementUnit),
		Ingredient:      ConvertValidIngredientToGRPCValidIngredient(&x.Ingredient),
	}
}

func ConvertGRPCCreateValidIngredientPreparationRequestToValidIngredientPreparationCreationRequestInput(x *messages.CreateValidIngredientPreparationRequest) *types.ValidIngredientPreparationCreationRequestInput {
	return &types.ValidIngredientPreparationCreationRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidIngredientID:  x.ValidIngredientID,
	}
}

func ConvertGRPCValidIngredientPreparationUpdateRequestInputToValidIngredientPreparationUpdateRequestInput(x *messages.ValidIngredientPreparationUpdateRequestInput) *types.ValidIngredientPreparationUpdateRequestInput {
	return &types.ValidIngredientPreparationUpdateRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidIngredientID:  x.ValidIngredientID,
	}
}

func ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(x *types.ValidIngredientPreparation) *messages.ValidIngredientPreparation {
	return &messages.ValidIngredientPreparation{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		Notes:         x.Notes,
		ID:            x.ID,
		Preparation:   ConvertValidPreparationToGRPCValidPreparation(&x.Preparation),
		Ingredient:    ConvertValidIngredientToGRPCValidIngredient(&x.Ingredient),
	}
}

func ConvertGRPCCreateValidIngredientStateRequestToValidIngredientStateCreationRequestInput(x *messages.CreateValidIngredientStateRequest) *types.ValidIngredientStateCreationRequestInput {
	return &types.ValidIngredientStateCreationRequestInput{
		Name:          x.Name,
		Slug:          x.Slug,
		PastTense:     x.PastTense,
		Description:   x.Description,
		AttributeType: x.AttributeType,
		IconPath:      x.IconPath,
	}
}

func ConvertGRPCValidIngredientStateUpdateRequestInputToValidIngredientStateUpdateRequestInput(x *messages.ValidIngredientStateUpdateRequestInput) *types.ValidIngredientStateUpdateRequestInput {
	return &types.ValidIngredientStateUpdateRequestInput{
		Name:          x.Name,
		Slug:          x.Slug,
		PastTense:     x.PastTense,
		Description:   x.Description,
		AttributeType: x.AttributeType,
		IconPath:      x.IconPath,
	}
}

func ConvertValidIngredientStateToGRPCValidIngredientState(x *types.ValidIngredientState) *messages.ValidIngredientState {
	return &messages.ValidIngredientState{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		PastTense:     x.PastTense,
		Description:   x.Description,
		IconPath:      x.IconPath,
		ID:            x.ID,
		Name:          x.Name,
		AttributeType: x.AttributeType,
		Slug:          x.Slug,
	}
}

func ConvertGRPCCreateValidIngredientStateIngredientRequestToValidIngredientStateIngredientCreationRequestInput(x *messages.CreateValidIngredientStateIngredientRequest) *types.ValidIngredientStateIngredientCreationRequestInput {
	return &types.ValidIngredientStateIngredientCreationRequestInput{
		Notes:                  x.Notes,
		ValidIngredientStateID: x.ValidIngredientStateID,
		ValidIngredientID:      x.ValidIngredientID,
	}
}

func ConvertGRPCValidIngredientStateIngredientUpdateRequestInputToValidIngredientStateIngredientUpdateRequestInput(x *messages.ValidIngredientStateIngredientUpdateRequestInput) *types.ValidIngredientStateIngredientUpdateRequestInput {
	return &types.ValidIngredientStateIngredientUpdateRequestInput{
		Notes:                  x.Notes,
		ValidIngredientStateID: x.ValidIngredientStateID,
		ValidIngredientID:      x.ValidIngredientID,
	}
}

func ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(x *types.ValidIngredientStateIngredient) *messages.ValidIngredientStateIngredient {
	return &messages.ValidIngredientStateIngredient{
		CreatedAt:       ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:   ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:      ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		Notes:           x.Notes,
		ID:              x.ID,
		IngredientState: ConvertValidIngredientStateToGRPCValidIngredientState(&x.IngredientState),
		Ingredient:      ConvertValidIngredientToGRPCValidIngredient(&x.Ingredient),
	}
}

func ConvertGRPCCreateValidInstrumentRequestToValidInstrumentCreationRequestInput(x *messages.CreateValidInstrumentRequest) *types.ValidInstrumentCreationRequestInput {
	return &types.ValidInstrumentCreationRequestInput{
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		UsableForStorage:               x.UsableForStorage,
	}
}

func ConvertGRPCValidInstrumentUpdateRequestInputToValidInstrumentUpdateRequestInput(x *messages.ValidInstrumentUpdateRequestInput) *types.ValidInstrumentUpdateRequestInput {
	return &types.ValidInstrumentUpdateRequestInput{
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		Slug:                           x.Slug,
		UsableForStorage:               x.UsableForStorage,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
	}
}

func ConvertValidInstrumentToGRPCValidInstrument(x *types.ValidInstrument) *messages.ValidInstrument {
	return &messages.ValidInstrument{
		CreatedAt:                      ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:                  ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:                     ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		Name:                           x.Name,
		ID:                             x.ID,
		IconPath:                       x.IconPath,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		UsableForStorage:               x.UsableForStorage,
	}
}

func ConvertGRPCCreateValidMeasurementUnitRequestToValidMeasurementUnitCreationRequestInput(x *messages.CreateValidMeasurementUnitRequest) *types.ValidMeasurementUnitCreationRequestInput {
	return &types.ValidMeasurementUnitCreationRequestInput{
		Name:        x.Name,
		Description: x.Description,
		IconPath:    x.IconPath,
		PluralName:  x.PluralName,
		Slug:        x.Slug,
		Volumetric:  x.Volumetric,
		Universal:   x.Universal,
		Metric:      x.Metric,
		Imperial:    x.Imperial,
	}
}

func ConvertGRPCValidMeasurementUnitUpdateRequestInputToValidMeasurementUnitUpdateRequestInput(x *messages.ValidMeasurementUnitUpdateRequestInput) *types.ValidMeasurementUnitUpdateRequestInput {
	return &types.ValidMeasurementUnitUpdateRequestInput{
		Name:        x.Name,
		Description: x.Description,
		IconPath:    x.IconPath,
		Volumetric:  x.Volumetric,
		Universal:   x.Universal,
		Metric:      x.Metric,
		Imperial:    x.Imperial,
		PluralName:  x.PluralName,
		Slug:        x.Slug,
	}
}

func ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x *types.ValidMeasurementUnit) *messages.ValidMeasurementUnit {
	return &messages.ValidMeasurementUnit{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		PluralName:    x.PluralName,
		IconPath:      x.IconPath,
		ID:            x.ID,
		Description:   x.Description,
		Name:          x.Name,
		Slug:          x.Slug,
		Volumetric:    x.Volumetric,
		Universal:     x.Universal,
		Metric:        x.Metric,
		Imperial:      x.Imperial,
	}
}

func ConvertGRPCCreateValidMeasurementUnitConversionRequestToValidMeasurementUnitConversionCreationRequestInput(x *messages.CreateValidMeasurementUnitConversionRequest) *types.ValidMeasurementUnitConversionCreationRequestInput {
	return &types.ValidMeasurementUnitConversionCreationRequestInput{
		OnlyForIngredient: x.OnlyForIngredient,
		From:              x.From,
		To:                x.To,
		Notes:             x.Notes,
		Modifier:          x.Modifier,
	}
}

func ConvertGRPCValidMeasurementUnitConversionUpdateRequestInputToValidMeasurementUnitConversionUpdateRequestInput(x *messages.ValidMeasurementUnitConversionUpdateRequestInput) *types.ValidMeasurementUnitConversionUpdateRequestInput {
	return &types.ValidMeasurementUnitConversionUpdateRequestInput{
		OnlyForIngredient: x.OnlyForIngredient,
		From:              x.From,
		To:                x.To,
		Notes:             x.Notes,
		Modifier:          x.Modifier,
	}
}

func ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(x *types.ValidMeasurementUnitConversion) *messages.ValidMeasurementUnitConversion {
	var ingredient *messages.ValidIngredient
	if x.OnlyForIngredient != nil {
		ingredient = ConvertValidIngredientToGRPCValidIngredient(x.OnlyForIngredient)
	}

	y := &messages.ValidMeasurementUnitConversion{
		CreatedAt:         ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:     ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:        ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		From:              ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&x.From),
		To:                ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&x.To),
		Notes:             x.Notes,
		OnlyForIngredient: ingredient,
		ID:                x.ID,
		Modifier:          x.Modifier,
	}

	return y
}

func ConvertGRPCCreateValidPreparationRequestToValidPreparationCreationRequestInput(x *messages.CreateValidPreparationRequest) *types.ValidPreparationCreationRequestInput {
	return &types.ValidPreparationCreationRequestInput{
		InstrumentCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.InstrumentCount.Min),
			Max: ConvertUint32PointerToUint16Pointer(x.InstrumentCount.Max),
		},
		IngredientCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.IngredientCount.Min),
			Max: ConvertUint32PointerToUint16Pointer(x.IngredientCount.Max),
		},
		VesselCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.VesselCount.Min),
			Max: ConvertUint32PointerToUint16Pointer(x.VesselCount.Max),
		},
		IconPath:                    x.IconPath,
		PastTense:                   x.PastTense,
		Slug:                        x.Slug,
		Name:                        x.Name,
		Description:                 x.Description,
		TemperatureRequired:         x.TemperatureRequired,
		TimeEstimateRequired:        x.TimeEstimateRequired,
		ConditionExpressionRequired: x.ConditionExpressionRequired,
		ConsumesVessel:              x.ConsumesVessel,
		OnlyForVessels:              x.OnlyForVessels,
		RestrictToIngredients:       x.RestrictToIngredients,
		YieldsNothing:               x.YieldsNothing,
	}
}

func ConvertGRPCValidPreparationUpdateRequestInputToValidPreparationUpdateRequestInput(x *messages.ValidPreparationUpdateRequestInput) *types.ValidPreparationUpdateRequestInput {
	return &types.ValidPreparationUpdateRequestInput{
		InstrumentCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(x.InstrumentCount.Min)),
			Max: ConvertUint32PointerToUint16Pointer(x.InstrumentCount.Max),
		},
		IngredientCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(x.IngredientCount.Min)),
			Max: ConvertUint32PointerToUint16Pointer(x.IngredientCount.Max),
		},
		VesselCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(x.VesselCount.Min)),
			Max: ConvertUint32PointerToUint16Pointer(x.VesselCount.Max),
		},
		Name:                        x.Name,
		Description:                 x.Description,
		IconPath:                    x.IconPath,
		YieldsNothing:               x.YieldsNothing,
		Slug:                        x.Slug,
		RestrictToIngredients:       x.RestrictToIngredients,
		PastTense:                   x.PastTense,
		TemperatureRequired:         x.TemperatureRequired,
		TimeEstimateRequired:        x.TimeEstimateRequired,
		ConditionExpressionRequired: x.ConditionExpressionRequired,
		ConsumesVessel:              x.ConsumesVessel,
		OnlyForVessels:              x.OnlyForVessels,
	}
}

func ConvertValidPreparationToGRPCValidPreparation(x *types.ValidPreparation) *messages.ValidPreparation {
	return &messages.ValidPreparation{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		InstrumentCount: &messages.Uint16RangeWithOptionalMax{
			Min: uint32(x.InstrumentCount.Min),
			Max: ConvertUint16PointerToUint32Pointer(x.InstrumentCount.Max),
		},
		IngredientCount: &messages.Uint16RangeWithOptionalMax{
			Min: uint32(x.IngredientCount.Min),
			Max: ConvertUint16PointerToUint32Pointer(x.IngredientCount.Max),
		},
		VesselCount: &messages.Uint16RangeWithOptionalMax{
			Min: uint32(x.VesselCount.Min),
			Max: ConvertUint16PointerToUint32Pointer(x.VesselCount.Max),
		},
		Name:                        x.Name,
		ID:                          x.ID,
		IconPath:                    x.IconPath,
		Description:                 x.Description,
		Slug:                        x.Slug,
		PastTense:                   x.PastTense,
		TemperatureRequired:         x.TemperatureRequired,
		ConditionExpressionRequired: x.ConditionExpressionRequired,
		ConsumesVessel:              x.ConsumesVessel,
		OnlyForVessels:              x.OnlyForVessels,
		YieldsNothing:               x.YieldsNothing,
		TimeEstimateRequired:        x.TimeEstimateRequired,
		RestrictToIngredients:       x.RestrictToIngredients,
	}
}

func ConvertGRPCCreateValidPreparationInstrumentRequestToValidPreparationInstrumentCreationRequestInput(x *messages.CreateValidPreparationInstrumentRequest) *types.ValidPreparationInstrumentCreationRequestInput {
	return &types.ValidPreparationInstrumentCreationRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidInstrumentID:  x.ValidInstrumentID,
	}
}

func ConvertGRPCValidPreparationInstrumentUpdateRequestInputToValidPreparationInstrumentUpdateRequestInput(x *messages.ValidPreparationInstrumentUpdateRequestInput) *types.ValidPreparationInstrumentUpdateRequestInput {
	return &types.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidInstrumentID:  x.ValidInstrumentID,
	}
}

func ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(x *types.ValidPreparationInstrument) *messages.ValidPreparationInstrument {
	return &messages.ValidPreparationInstrument{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		ID:            x.ID,
		Notes:         x.Notes,
		Instrument:    ConvertValidInstrumentToGRPCValidInstrument(&x.Instrument),
		Preparation:   ConvertValidPreparationToGRPCValidPreparation(&x.Preparation),
	}
}

func ConvertGRPCCreateValidPreparationVesselRequestToValidPreparationVesselCreationRequestInput(x *messages.CreateValidPreparationVesselRequest) *types.ValidPreparationVesselCreationRequestInput {
	return &types.ValidPreparationVesselCreationRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidVesselID:      x.ValidVesselID,
	}
}

func ConvertGRPCValidPreparationVesselUpdateRequestInputToValidPreparationVesselUpdateRequestInput(x *messages.ValidPreparationVesselUpdateRequestInput) *types.ValidPreparationVesselUpdateRequestInput {
	return &types.ValidPreparationVesselUpdateRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidVesselID:      x.ValidVesselID,
	}
}

func ConvertValidPreparationVesselToGRPCValidPreparationVessel(x *types.ValidPreparationVessel) *messages.ValidPreparationVessel {
	return &messages.ValidPreparationVessel{
		CreatedAt:     ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		ID:            x.ID,
		Notes:         x.Notes,
		Preparation:   ConvertValidPreparationToGRPCValidPreparation(&x.Preparation),
		Vessel:        ConvertValidVesselToGRPCValidVessel(&x.Vessel),
	}
}

func ConvertGRPCCreateValidVesselRequestToValidVesselCreationRequestInput(x *messages.CreateValidVesselRequest) *types.ValidVesselCreationRequestInput {
	return &types.ValidVesselCreationRequestInput{
		CapacityUnitID:                 x.CapacityUnitID,
		Shape:                          x.Shape,
		IconPath:                       x.IconPath,
		PluralName:                     x.PluralName,
		Name:                           x.Name,
		Description:                    x.Description,
		Slug:                           x.Slug,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Capacity:                       x.Capacity,
		WidthInMillimeters:             x.WidthInMillimeters,
		UsableForStorage:               x.UsableForStorage,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
	}
}

func ConvertGRPCValidVesselUpdateRequestInputToValidVesselUpdateRequestInput(x *messages.ValidVesselUpdateRequestInput) *types.ValidVesselUpdateRequestInput {
	return &types.ValidVesselUpdateRequestInput{
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		UsableForStorage:               x.UsableForStorage,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		Capacity:                       x.Capacity,
		CapacityUnitID:                 x.CapacityUnitID,
		WidthInMillimeters:             x.WidthInMillimeters,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Shape:                          x.Shape,
	}
}

func ConvertValidVesselToGRPCValidVessel(x *types.ValidVessel) *messages.ValidVessel {
	return &messages.ValidVessel{
		CreatedAt:                      ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:                  ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:                     ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		CapacityUnit:                   ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x.CapacityUnit),
		Shape:                          x.Shape,
		Description:                    x.Description,
		Name:                           x.Name,
		Slug:                           x.Slug,
		IconPath:                       x.IconPath,
		ID:                             x.ID,
		PluralName:                     x.PluralName,
		WidthInMillimeters:             x.WidthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Capacity:                       x.Capacity,
		LengthInMillimeters:            x.LengthInMillimeters,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		UsableForStorage:               x.UsableForStorage,
	}
}
