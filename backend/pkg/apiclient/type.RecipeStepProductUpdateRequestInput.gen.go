// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepProductUpdateRequestInput struct {
		QuantityNotes               string               `json:"quantityNotes"`
		MeasurementUnitID           string               `json:"measurementUnitID"`
		Type                        string               `json:"type"`
		StorageInstructions         string               `json:"storageInstructions"`
		Name                        string               `json:"name"`
		BelongsToRecipeStep         string               `json:"belongsToRecipeStep"`
		Quantity                    OptionalFloat32Range `json:"quantity"`
		StorageDurationInSeconds    OptionalUint32Range  `json:"storageDurationInSeconds"`
		StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
		Index                       uint64               `json:"index"`
		ContainedInVesselIndex      uint64               `json:"containedInVesselIndex"`
		Compostable                 bool                 `json:"compostable"`
		IsWaste                     bool                 `json:"isWaste"`
		IsLiquid                    bool                 `json:"isLiquid"`
	}
)
