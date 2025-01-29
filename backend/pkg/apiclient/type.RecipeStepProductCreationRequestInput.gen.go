// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepProductCreationRequestInput struct {
   Compostable bool `json:"compostable"`
 ContainedInVesselIndex uint64 `json:"containedInVesselIndex"`
 Index uint64 `json:"index"`
 IsLiquid bool `json:"isLiquid"`
 IsWaste bool `json:"isWaste"`
 MeasurementUnitID string `json:"measurementUnitID"`
 Name string `json:"name"`
 Quantity OptionalFloat32Range `json:"quantity"`
 QuantityNotes string `json:"quantityNotes"`
 StorageDurationInSeconds OptionalUint32Range `json:"storageDurationInSeconds"`
 StorageInstructions string `json:"storageInstructions"`
 StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
 Type string `json:"type"`

}
)
