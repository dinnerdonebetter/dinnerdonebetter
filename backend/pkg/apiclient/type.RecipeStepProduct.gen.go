// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepProduct struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 Compostable bool `json:"compostable"`
 ContainedInVesselIndex uint64 `json:"containedInVesselIndex"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 Index uint64 `json:"index"`
 IsLiquid bool `json:"isLiquid"`
 IsWaste bool `json:"isWaste"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 MeasurementUnit ValidMeasurementUnit `json:"measurementUnit"`
 Name string `json:"name"`
 Quantity OptionalFloat32Range `json:"quantity"`
 QuantityNotes string `json:"quantityNotes"`
 StorageDurationInSeconds OptionalUint32Range `json:"storageDurationInSeconds"`
 StorageInstructions string `json:"storageInstructions"`
 StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
 Type string `json:"type"`

}
)
