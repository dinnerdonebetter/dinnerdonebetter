package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	grpctypes "github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func ConvertGRPCCreateValidIngredientRequestToValidIngredientCreationRequestInput(request *mealplanninggrpc.ValidIngredientCreationRequestInput) *mealplanning.ValidIngredientCreationRequestInput {
	return &mealplanning.ValidIngredientCreationRequestInput{
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

func ConvertGRPCValidIngredientUpdateRequestInputToValidIngredientUpdateRequestInput(x *mealplanninggrpc.ValidIngredientUpdateRequestInput) *mealplanning.ValidIngredientUpdateRequestInput {
	storageTemperatureInCelsius := types.OptionalFloat32Range{}
	if x.StorageTemperatureInCelsius != nil {
		storageTemperatureInCelsius = types.OptionalFloat32Range{
			Max: x.StorageTemperatureInCelsius.Max,
			Min: x.StorageTemperatureInCelsius.Min,
		}
	}

	return &mealplanning.ValidIngredientUpdateRequestInput{
		Name:                        x.Name,
		Description:                 x.Description,
		Warning:                     x.Warning,
		IconPath:                    x.IconPath,
		ContainsDairy:               x.ContainsDairy,
		ContainsPeanut:              x.ContainsPeanut,
		ContainsTreeNut:             x.ContainsTreeNut,
		ContainsEgg:                 x.ContainsEgg,
		ContainsWheat:               x.ContainsWheat,
		ContainsShellfish:           x.ContainsShellfish,
		ContainsSesame:              x.ContainsSesame,
		ContainsFish:                x.ContainsFish,
		ContainsGluten:              x.ContainsGluten,
		AnimalFlesh:                 x.AnimalFlesh,
		IsLiquid:                    x.IsLiquid,
		ContainsSoy:                 x.ContainsSoy,
		PluralName:                  x.PluralName,
		AnimalDerived:               x.AnimalDerived,
		RestrictToPreparations:      x.RestrictToPreparations,
		StorageTemperatureInCelsius: storageTemperatureInCelsius,
		StorageInstructions:         x.StorageInstructions,
		Slug:                        x.Slug,
		ContainsAlcohol:             x.ContainsAlcohol,
		ShoppingSuggestions:         x.ShoppingSuggestions,
		IsStarch:                    x.IsStarch,
		IsProtein:                   x.IsProtein,
		IsGrain:                     x.IsGrain,
		IsFruit:                     x.IsFruit,
		IsSalt:                      x.IsSalt,
		IsFat:                       x.IsFat,
		IsAcid:                      x.IsAcid,
		IsHeat:                      x.IsHeat,
	}
}

func ConvertValidIngredientUpdateRequestInputToGRPCValidIngredientUpdateRequestInput(x *mealplanning.ValidIngredientUpdateRequestInput) *mealplanninggrpc.ValidIngredientUpdateRequestInput {
	return &mealplanninggrpc.ValidIngredientUpdateRequestInput{
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
		StorageTemperatureInCelsius: &grpctypes.OptionalFloat32Range{
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

func ConvertValidIngredientToGRPCValidIngredient(x *mealplanning.ValidIngredient) *mealplanninggrpc.ValidIngredient {
	return &mealplanninggrpc.ValidIngredient{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		StorageTemperatureInCelsius: &grpctypes.OptionalFloat32Range{
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

func ConvertGRPCValidIngredientToValidIngredient(x *mealplanninggrpc.ValidIngredient) *mealplanning.ValidIngredient {
	return &mealplanning.ValidIngredient{
		CreatedAt:     converters.ConvertPBTimestampToTime(x.CreatedAt),
		LastUpdatedAt: converters.ConvertPBTimestampToTimePointer(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertPBTimestampToTimePointer(x.ArchivedAt),
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
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

func ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(x *mealplanning.ValidIngredientCreationRequestInput) *mealplanninggrpc.ValidIngredientCreationRequestInput {
	return &mealplanninggrpc.ValidIngredientCreationRequestInput{
		StorageTemperatureInCelsius: &grpctypes.OptionalFloat32Range{
			Max: x.StorageTemperatureInCelsius.Max,
			Min: x.StorageTemperatureInCelsius.Min,
		},
		Warning:                x.Warning,
		IconPath:               x.IconPath,
		PluralName:             x.PluralName,
		StorageInstructions:    x.StorageInstructions,
		Name:                   x.Name,
		Description:            x.Description,
		Slug:                   x.Slug,
		ShoppingSuggestions:    x.ShoppingSuggestions,
		ContainsPeanut:         x.ContainsPeanut,
		ContainsAlcohol:        x.ContainsAlcohol,
		IsLiquid:               x.IsLiquid,
		ContainsSoy:            x.ContainsSoy,
		AnimalFlesh:            x.AnimalFlesh,
		AnimalDerived:          x.AnimalDerived,
		RestrictToPreparations: x.RestrictToPreparations,
		ContainsDairy:          x.ContainsDairy,
		ContainsSesame:         x.ContainsSesame,
		ContainsTreeNut:        x.ContainsTreeNut,
		ContainsWheat:          x.ContainsWheat,
		ContainsEgg:            x.ContainsEgg,
		ContainsGluten:         x.ContainsGluten,
		IsStarch:               x.IsStarch,
		IsProtein:              x.IsProtein,
		IsGrain:                x.IsGrain,
		IsFruit:                x.IsFruit,
		IsSalt:                 x.IsSalt,
		IsFat:                  x.IsFat,
		IsAcid:                 x.IsAcid,
		IsHeat:                 x.IsHeat,
		ContainsShellfish:      x.ContainsShellfish,
		ContainsFish:           x.ContainsFish,
	}
}

func ConvertGRPCCreateValidIngredientGroupRequestToValidIngredientGroupCreationRequestInput(request *mealplanninggrpc.CreateValidIngredientGroupRequest) *mealplanning.ValidIngredientGroupCreationRequestInput {
	members := make([]*mealplanning.ValidIngredientGroupMemberCreationRequestInput, len(request.Members))
	for i, member := range request.Members {
		members[i] = &mealplanning.ValidIngredientGroupMemberCreationRequestInput{
			ValidIngredientID: member.ValidIngredientID,
		}
	}

	return &mealplanning.ValidIngredientGroupCreationRequestInput{
		Name:        request.Name,
		Slug:        request.Slug,
		Description: request.Description,
		Members:     members,
	}
}

func ConvertGRPCValidIngredientGroupUpdateRequestInputToValidIngredientGroupUpdateRequestInput(x *mealplanninggrpc.ValidIngredientGroupUpdateRequestInput) *mealplanning.ValidIngredientGroupUpdateRequestInput {
	return &mealplanning.ValidIngredientGroupUpdateRequestInput{
		Name:        x.Name,
		Slug:        x.Slug,
		Description: x.Description,
	}
}

func ConvertValidIngredientGroupToGRPCValidIngredientGroup(x *mealplanning.ValidIngredientGroup) *mealplanninggrpc.ValidIngredientGroup {
	members := make([]*mealplanninggrpc.ValidIngredientGroupMember, len(x.Members))
	for i, member := range x.Members {
		members[i] = &mealplanninggrpc.ValidIngredientGroupMember{
			CreatedAt:       converters.ConvertTimeToPBTimestamp(member.CreatedAt),
			ArchivedAt:      converters.ConvertTimePointerToPBTimestamp(member.ArchivedAt),
			ID:              member.ID,
			BelongsToGroup:  member.BelongsToGroup,
			ValidIngredient: ConvertValidIngredientToGRPCValidIngredient(&member.ValidIngredient),
		}
	}

	return &mealplanninggrpc.ValidIngredientGroup{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		ID:            x.ID,
		Name:          x.Name,
		Slug:          x.Slug,
		Description:   x.Description,
		Members:       members,
	}
}

func ConvertGRPCCreateValidIngredientMeasurementUnitRequestToValidIngredientMeasurementUnitCreationRequestInput(request *mealplanninggrpc.CreateValidIngredientMeasurementUnitRequest) *mealplanning.ValidIngredientMeasurementUnitCreationRequestInput {
	return &mealplanning.ValidIngredientMeasurementUnitCreationRequestInput{
		Notes:                  request.Notes,
		ValidMeasurementUnitID: request.ValidMeasurementUnitID,
		ValidIngredientID:      request.ValidIngredientID,
		AllowableQuantity: types.Float32RangeWithOptionalMax{
			Max: request.AllowableQuantity.Max,
			Min: request.AllowableQuantity.Min,
		},
	}
}

func ConvertGRPCValidIngredientMeasurementUnitUpdateRequestInputToValidIngredientMeasurementUnitUpdateRequestInput(x *mealplanninggrpc.ValidIngredientMeasurementUnitUpdateRequestInput) *mealplanning.ValidIngredientMeasurementUnitUpdateRequestInput {
	return &mealplanning.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                  x.Notes,
		ValidMeasurementUnitID: x.ValidMeasurementUnitID,
		ValidIngredientID:      x.ValidIngredientID,
		AllowableQuantity: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Min: x.AllowableQuantity.Min,
			Max: x.AllowableQuantity.Max,
		},
	}
}

func ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(x *mealplanning.ValidIngredientMeasurementUnit) *mealplanninggrpc.ValidIngredientMeasurementUnit {
	return &mealplanninggrpc.ValidIngredientMeasurementUnit{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		Notes:         x.Notes,
		ID:            x.ID,
		AllowableQuantity: &grpctypes.Float32RangeWithOptionalMax{
			Max: x.AllowableQuantity.Max,
			Min: x.AllowableQuantity.Min,
		},
		MeasurementUnit: ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&x.MeasurementUnit),
		Ingredient:      ConvertValidIngredientToGRPCValidIngredient(&x.Ingredient),
	}
}

func ConvertGRPCCreateValidIngredientPreparationRequestToValidIngredientPreparationCreationRequestInput(x *mealplanninggrpc.CreateValidIngredientPreparationRequest) *mealplanning.ValidIngredientPreparationCreationRequestInput {
	return &mealplanning.ValidIngredientPreparationCreationRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidIngredientID:  x.ValidIngredientID,
	}
}

func ConvertGRPCValidIngredientPreparationUpdateRequestInputToValidIngredientPreparationUpdateRequestInput(x *mealplanninggrpc.ValidIngredientPreparationUpdateRequestInput) *mealplanning.ValidIngredientPreparationUpdateRequestInput {
	return &mealplanning.ValidIngredientPreparationUpdateRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidIngredientID:  x.ValidIngredientID,
	}
}

func ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(x *mealplanning.ValidIngredientPreparation) *mealplanninggrpc.ValidIngredientPreparation {
	return &mealplanninggrpc.ValidIngredientPreparation{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		Notes:         x.Notes,
		ID:            x.ID,
		Preparation:   ConvertValidPreparationToGRPCValidPreparation(&x.Preparation),
		Ingredient:    ConvertValidIngredientToGRPCValidIngredient(&x.Ingredient),
	}
}

func ConvertGRPCCreateValidIngredientStateRequestToValidIngredientStateCreationRequestInput(x *mealplanninggrpc.CreateValidIngredientStateRequest) *mealplanning.ValidIngredientStateCreationRequestInput {
	return &mealplanning.ValidIngredientStateCreationRequestInput{
		Name:          x.Name,
		Slug:          x.Slug,
		PastTense:     x.PastTense,
		Description:   x.Description,
		AttributeType: x.AttributeType,
		IconPath:      x.IconPath,
	}
}

func ConvertGRPCValidIngredientStateUpdateRequestInputToValidIngredientStateUpdateRequestInput(x *mealplanninggrpc.ValidIngredientStateUpdateRequestInput) *mealplanning.ValidIngredientStateUpdateRequestInput {
	return &mealplanning.ValidIngredientStateUpdateRequestInput{
		Name:          x.Name,
		Slug:          x.Slug,
		PastTense:     x.PastTense,
		Description:   x.Description,
		AttributeType: x.AttributeType,
		IconPath:      x.IconPath,
	}
}

func ConvertValidIngredientStateToGRPCValidIngredientState(x *mealplanning.ValidIngredientState) *mealplanninggrpc.ValidIngredientState {
	return &mealplanninggrpc.ValidIngredientState{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		PastTense:     x.PastTense,
		Description:   x.Description,
		IconPath:      x.IconPath,
		ID:            x.ID,
		Name:          x.Name,
		AttributeType: x.AttributeType,
		Slug:          x.Slug,
	}
}

func ConvertGRPCCreateValidIngredientStateIngredientRequestToValidIngredientStateIngredientCreationRequestInput(x *mealplanninggrpc.CreateValidIngredientStateIngredientRequest) *mealplanning.ValidIngredientStateIngredientCreationRequestInput {
	return &mealplanning.ValidIngredientStateIngredientCreationRequestInput{
		Notes:                  x.Notes,
		ValidIngredientStateID: x.ValidIngredientStateID,
		ValidIngredientID:      x.ValidIngredientID,
	}
}

func ConvertGRPCValidIngredientStateIngredientUpdateRequestInputToValidIngredientStateIngredientUpdateRequestInput(x *mealplanninggrpc.ValidIngredientStateIngredientUpdateRequestInput) *mealplanning.ValidIngredientStateIngredientUpdateRequestInput {
	return &mealplanning.ValidIngredientStateIngredientUpdateRequestInput{
		Notes:                  x.Notes,
		ValidIngredientStateID: x.ValidIngredientStateID,
		ValidIngredientID:      x.ValidIngredientID,
	}
}

func ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(x *mealplanning.ValidIngredientStateIngredient) *mealplanninggrpc.ValidIngredientStateIngredient {
	return &mealplanninggrpc.ValidIngredientStateIngredient{
		CreatedAt:       converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:   converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:      converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		Notes:           x.Notes,
		ID:              x.ID,
		IngredientState: ConvertValidIngredientStateToGRPCValidIngredientState(&x.IngredientState),
		Ingredient:      ConvertValidIngredientToGRPCValidIngredient(&x.Ingredient),
	}
}

func ConvertGRPCCreateValidInstrumentRequestToValidInstrumentCreationRequestInput(x *mealplanninggrpc.CreateValidInstrumentRequest) *mealplanning.ValidInstrumentCreationRequestInput {
	return &mealplanning.ValidInstrumentCreationRequestInput{
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

func ConvertGRPCValidInstrumentUpdateRequestInputToValidInstrumentUpdateRequestInput(x *mealplanninggrpc.ValidInstrumentUpdateRequestInput) *mealplanning.ValidInstrumentUpdateRequestInput {
	return &mealplanning.ValidInstrumentUpdateRequestInput{
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

func ConvertValidInstrumentToGRPCValidInstrument(x *mealplanning.ValidInstrument) *mealplanninggrpc.ValidInstrument {
	return &mealplanninggrpc.ValidInstrument{
		CreatedAt:                      converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:                  converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:                     converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
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

func ConvertGRPCCreateValidMeasurementUnitRequestToValidMeasurementUnitCreationRequestInput(x *mealplanninggrpc.CreateValidMeasurementUnitRequest) *mealplanning.ValidMeasurementUnitCreationRequestInput {
	return &mealplanning.ValidMeasurementUnitCreationRequestInput{
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

func ConvertGRPCValidMeasurementUnitUpdateRequestInputToValidMeasurementUnitUpdateRequestInput(x *mealplanninggrpc.ValidMeasurementUnitUpdateRequestInput) *mealplanning.ValidMeasurementUnitUpdateRequestInput {
	return &mealplanning.ValidMeasurementUnitUpdateRequestInput{
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

func ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x *mealplanning.ValidMeasurementUnit) *mealplanninggrpc.ValidMeasurementUnit {
	return &mealplanninggrpc.ValidMeasurementUnit{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
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

func ConvertGRPCCreateValidMeasurementUnitConversionRequestToValidMeasurementUnitConversionCreationRequestInput(x *mealplanninggrpc.CreateValidMeasurementUnitConversionRequest) *mealplanning.ValidMeasurementUnitConversionCreationRequestInput {
	return &mealplanning.ValidMeasurementUnitConversionCreationRequestInput{
		OnlyForIngredient: x.OnlyForIngredient,
		From:              x.From,
		To:                x.To,
		Notes:             x.Notes,
		Modifier:          x.Modifier,
	}
}

func ConvertGRPCValidMeasurementUnitConversionUpdateRequestInputToValidMeasurementUnitConversionUpdateRequestInput(x *mealplanninggrpc.ValidMeasurementUnitConversionUpdateRequestInput) *mealplanning.ValidMeasurementUnitConversionUpdateRequestInput {
	return &mealplanning.ValidMeasurementUnitConversionUpdateRequestInput{
		OnlyForIngredient: x.OnlyForIngredient,
		From:              x.From,
		To:                x.To,
		Notes:             x.Notes,
		Modifier:          x.Modifier,
	}
}

func ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(x *mealplanning.ValidMeasurementUnitConversion) *mealplanninggrpc.ValidMeasurementUnitConversion {
	var ingredient *mealplanninggrpc.ValidIngredient
	if x.OnlyForIngredient != nil {
		ingredient = ConvertValidIngredientToGRPCValidIngredient(x.OnlyForIngredient)
	}

	y := &mealplanninggrpc.ValidMeasurementUnitConversion{
		CreatedAt:         converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:     converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:        converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		From:              ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&x.From),
		To:                ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&x.To),
		Notes:             x.Notes,
		OnlyForIngredient: ingredient,
		ID:                x.ID,
		Modifier:          x.Modifier,
	}

	return y
}

func ConvertGRPCCreateValidPreparationRequestToValidPreparationCreationRequestInput(x *mealplanninggrpc.CreateValidPreparationRequest) *mealplanning.ValidPreparationCreationRequestInput {
	return &mealplanning.ValidPreparationCreationRequestInput{
		InstrumentCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.InstrumentCount.Min),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.InstrumentCount.Max),
		},
		IngredientCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.IngredientCount.Min),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.IngredientCount.Max),
		},
		VesselCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.VesselCount.Min),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.VesselCount.Max),
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

func ConvertGRPCValidPreparationUpdateRequestInputToValidPreparationUpdateRequestInput(x *mealplanninggrpc.ValidPreparationUpdateRequestInput) *mealplanning.ValidPreparationUpdateRequestInput {
	return &mealplanning.ValidPreparationUpdateRequestInput{
		InstrumentCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(pointer.Dereference(x.InstrumentCount.Min))),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.InstrumentCount.Max),
		},
		IngredientCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(pointer.Dereference(x.IngredientCount.Min))),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.IngredientCount.Max),
		},
		VesselCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(pointer.Dereference(x.VesselCount.Min))),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.VesselCount.Max),
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

func ConvertValidPreparationToGRPCValidPreparation(x *mealplanning.ValidPreparation) *mealplanninggrpc.ValidPreparation {
	return &mealplanninggrpc.ValidPreparation{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		InstrumentCount: &grpctypes.Uint16RangeWithOptionalMax{
			Min: uint32(x.InstrumentCount.Min),
			Max: converters.ConvertUint16PointerToUint32Pointer(x.InstrumentCount.Max),
		},
		IngredientCount: &grpctypes.Uint16RangeWithOptionalMax{
			Min: uint32(x.IngredientCount.Min),
			Max: converters.ConvertUint16PointerToUint32Pointer(x.IngredientCount.Max),
		},
		VesselCount: &grpctypes.Uint16RangeWithOptionalMax{
			Min: uint32(x.VesselCount.Min),
			Max: converters.ConvertUint16PointerToUint32Pointer(x.VesselCount.Max),
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

func ConvertGRPCValidPreparationToValidPreparation(x *mealplanninggrpc.ValidPreparation) *mealplanning.ValidPreparation {
	return &mealplanning.ValidPreparation{
		InstrumentCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.InstrumentCount.Min),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.InstrumentCount.Max),
		},
		IngredientCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.IngredientCount.Min),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.IngredientCount.Max),
		},
		VesselCount: types.Uint16RangeWithOptionalMax{
			Min: uint16(x.VesselCount.Min),
			Max: converters.ConvertUint32PointerToUint16Pointer(x.VesselCount.Max),
		},
		CreatedAt:                   converters.ConvertPBTimestampToTime(x.CreatedAt),
		ArchivedAt:                  converters.ConvertPBTimestampToTimePointer(x.ArchivedAt),
		LastUpdatedAt:               converters.ConvertPBTimestampToTimePointer(x.LastUpdatedAt),
		IconPath:                    x.IconPath,
		PastTense:                   x.PastTense,
		ID:                          x.ID,
		Name:                        x.Name,
		Description:                 x.Description,
		Slug:                        x.Slug,
		RestrictToIngredients:       x.RestrictToIngredients,
		TemperatureRequired:         x.TemperatureRequired,
		TimeEstimateRequired:        x.TimeEstimateRequired,
		ConditionExpressionRequired: x.ConditionExpressionRequired,
		ConsumesVessel:              x.ConsumesVessel,
		OnlyForVessels:              x.OnlyForVessels,
		YieldsNothing:               x.YieldsNothing,
	}
}

func ConvertGRPCCreateValidPreparationInstrumentRequestToValidPreparationInstrumentCreationRequestInput(x *mealplanninggrpc.CreateValidPreparationInstrumentRequest) *mealplanning.ValidPreparationInstrumentCreationRequestInput {
	return &mealplanning.ValidPreparationInstrumentCreationRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidInstrumentID:  x.ValidInstrumentID,
	}
}

func ConvertGRPCValidPreparationInstrumentUpdateRequestInputToValidPreparationInstrumentUpdateRequestInput(x *mealplanninggrpc.ValidPreparationInstrumentUpdateRequestInput) *mealplanning.ValidPreparationInstrumentUpdateRequestInput {
	return &mealplanning.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidInstrumentID:  x.ValidInstrumentID,
	}
}

func ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(x *mealplanning.ValidPreparationInstrument) *mealplanninggrpc.ValidPreparationInstrument {
	return &mealplanninggrpc.ValidPreparationInstrument{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		ID:            x.ID,
		Notes:         x.Notes,
		Instrument:    ConvertValidInstrumentToGRPCValidInstrument(&x.Instrument),
		Preparation:   ConvertValidPreparationToGRPCValidPreparation(&x.Preparation),
	}
}

func ConvertGRPCCreateValidPreparationVesselRequestToValidPreparationVesselCreationRequestInput(x *mealplanninggrpc.CreateValidPreparationVesselRequest) *mealplanning.ValidPreparationVesselCreationRequestInput {
	return &mealplanning.ValidPreparationVesselCreationRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidVesselID:      x.ValidVesselID,
	}
}

func ConvertGRPCValidPreparationVesselUpdateRequestInputToValidPreparationVesselUpdateRequestInput(x *mealplanninggrpc.ValidPreparationVesselUpdateRequestInput) *mealplanning.ValidPreparationVesselUpdateRequestInput {
	return &mealplanning.ValidPreparationVesselUpdateRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidVesselID:      x.ValidVesselID,
	}
}

func ConvertValidPreparationVesselToGRPCValidPreparationVessel(x *mealplanning.ValidPreparationVessel) *mealplanninggrpc.ValidPreparationVessel {
	return &mealplanninggrpc.ValidPreparationVessel{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
		ID:            x.ID,
		Notes:         x.Notes,
		Preparation:   ConvertValidPreparationToGRPCValidPreparation(&x.Preparation),
		Vessel:        ConvertValidVesselToGRPCValidVessel(&x.Vessel),
	}
}

func ConvertGRPCCreateValidVesselRequestToValidVesselCreationRequestInput(x *mealplanninggrpc.CreateValidVesselRequest) *mealplanning.ValidVesselCreationRequestInput {
	return &mealplanning.ValidVesselCreationRequestInput{
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

func ConvertGRPCValidVesselUpdateRequestInputToValidVesselUpdateRequestInput(x *mealplanninggrpc.ValidVesselUpdateRequestInput) *mealplanning.ValidVesselUpdateRequestInput {
	return &mealplanning.ValidVesselUpdateRequestInput{
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

func ConvertValidVesselToGRPCValidVessel(x *mealplanning.ValidVessel) *mealplanninggrpc.ValidVessel {
	return &mealplanninggrpc.ValidVessel{
		CreatedAt:                      converters.ConvertTimeToPBTimestamp(x.CreatedAt),
		LastUpdatedAt:                  converters.ConvertTimePointerToPBTimestamp(x.LastUpdatedAt),
		ArchivedAt:                     converters.ConvertTimePointerToPBTimestamp(x.ArchivedAt),
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
