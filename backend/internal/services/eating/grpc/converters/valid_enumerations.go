package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
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

func ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(x *types.ValidIngredientMeasurementUnit) *messages.ValidIngredientMeasurementUnit {
	return &messages.ValidIngredientMeasurementUnit{}
}

func ConvertGRPCCreateValidIngredientPreparationRequestToValidIngredientPreparationCreationRequestInput(request *messages.CreateValidIngredientPreparationRequest) *types.ValidIngredientPreparationCreationRequestInput {
	return &types.ValidIngredientPreparationCreationRequestInput{}
}

func ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(x *types.ValidIngredientPreparation) *messages.ValidIngredientPreparation {
	return &messages.ValidIngredientPreparation{}
}

func ConvertGRPCCreateValidIngredientStateRequestToValidIngredientStateCreationRequestInput(request *messages.CreateValidIngredientStateRequest) *types.ValidIngredientStateCreationRequestInput {
	return &types.ValidIngredientStateCreationRequestInput{}
}

func ConvertValidIngredientStateToGRPCValidIngredientState(x *types.ValidIngredientState) *messages.ValidIngredientState {
	return &messages.ValidIngredientState{}
}

func ConvertGRPCCreateValidIngredientStateIngredientRequestToValidIngredientStateIngredientCreationRequestInput(request *messages.CreateValidIngredientStateIngredientRequest) *types.ValidIngredientStateIngredientCreationRequestInput {
	return &types.ValidIngredientStateIngredientCreationRequestInput{}
}

func ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(x *types.ValidIngredientStateIngredient) *messages.ValidIngredientStateIngredient {
	return &messages.ValidIngredientStateIngredient{}
}

func ConvertGRPCCreateValidInstrumentRequestToValidInstrumentCreationRequestInput(request *messages.CreateValidInstrumentRequest) *types.ValidInstrumentCreationRequestInput {
	return &types.ValidInstrumentCreationRequestInput{}
}

func ConvertValidInstrumentToGRPCValidInstrument(x *types.ValidInstrument) *messages.ValidInstrument {
	return &messages.ValidInstrument{}
}

func ConvertGRPCCreateValidMeasurementUnitRequestToValidMeasurementUnitCreationRequestInput(request *messages.CreateValidMeasurementUnitRequest) *types.ValidMeasurementUnitCreationRequestInput {
	return &types.ValidMeasurementUnitCreationRequestInput{}
}

func ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x *types.ValidMeasurementUnit) *messages.ValidMeasurementUnit {
	return &messages.ValidMeasurementUnit{}
}

func ConvertGRPCCreateValidMeasurementUnitConversionRequestToValidMeasurementUnitConversionCreationRequestInput(request *messages.CreateValidMeasurementUnitConversionRequest) *types.ValidMeasurementUnitConversionCreationRequestInput {
	return &types.ValidMeasurementUnitConversionCreationRequestInput{}
}

func ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(x *types.ValidMeasurementUnitConversion) *messages.ValidMeasurementUnitConversion {
	return &messages.ValidMeasurementUnitConversion{}
}

func ConvertGRPCCreateValidPreparationRequestToValidPreparationCreationRequestInput(request *messages.CreateValidPreparationRequest) *types.ValidPreparationCreationRequestInput {
	return &types.ValidPreparationCreationRequestInput{}
}

func ConvertValidPreparationToGRPCValidPreparation(x *types.ValidPreparation) *messages.ValidPreparation {
	return &messages.ValidPreparation{}
}

func ConvertGRPCCreateValidPreparationInstrumentRequestToValidPreparationInstrumentCreationRequestInput(request *messages.CreateValidPreparationInstrumentRequest) *types.ValidPreparationInstrumentCreationRequestInput {
	return &types.ValidPreparationInstrumentCreationRequestInput{}
}

func ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(x *types.ValidPreparationInstrument) *messages.ValidPreparationInstrument {
	return &messages.ValidPreparationInstrument{}
}

func ConvertGRPCCreateValidPreparationVesselRequestToValidPreparationVesselCreationRequestInput(request *messages.CreateValidPreparationVesselRequest) *types.ValidPreparationVesselCreationRequestInput {
	return &types.ValidPreparationVesselCreationRequestInput{}
}

func ConvertValidPreparationVesselToGRPCValidPreparationVessel(x *types.ValidPreparationVessel) *messages.ValidPreparationVessel {
	return &messages.ValidPreparationVessel{}
}

func ConvertGRPCCreateValidVesselRequestToValidVesselCreationRequestInput(request *messages.CreateValidVesselRequest) *types.ValidVesselCreationRequestInput {
	return &types.ValidVesselCreationRequestInput{}
}

func ConvertValidVesselToGRPCValidVessel(x *types.ValidVessel) *messages.ValidVessel {
	return &messages.ValidVessel{}
}

/*
ValidIngredientMeasurementUnit
ValidIngredientPreparation
ValidIngredientState
ValidIngredientStateIngredient
ValidInstrument
ValidMeasurementUnit
ValidMeasurementUnitConversion
ValidPreparation
ValidPreparationInstrument
ValidPreparationVessel
ValidVessel
*/
