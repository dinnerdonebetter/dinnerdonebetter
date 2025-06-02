// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepProductUpdateRequestInput struct {
		QuantityNotes               string               `json:"quantityNotes"`
		BelongsToRecipeStep         string               `json:"belongsToRecipeStep"`
		Type                        string               `json:"type"`
		Name                        string               `json:"name"`
		StorageInstructions         string               `json:"storageInstructions"`
		MeasurementUnitID           string               `json:"measurementUnitID"`
		StorageDurationInSeconds    OptionalUint32Range  `json:"storageDurationInSeconds"`
		StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
		Quantity                    OptionalFloat32Range `json:"quantity"`
		Index                       uint64               `json:"index"`
		ContainedInVesselIndex      uint64               `json:"containedInVesselIndex"`
		Compostable                 bool                 `json:"compostable"`
		IsWaste                     bool                 `json:"isWaste"`
		IsLiquid                    bool                 `json:"isLiquid"`
	}
)
