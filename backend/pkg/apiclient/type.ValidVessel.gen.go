// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidVessel struct {
		ID                             string               `json:"id"`
		CreatedAt                      string               `json:"createdAt"`
		ArchivedAt                     string               `json:"archivedAt"`
		LastUpdatedAt                  string               `json:"lastUpdatedAt"`
		Description                    string               `json:"description"`
		Slug                           string               `json:"slug"`
		Shape                          string               `json:"shape"`
		IconPath                       string               `json:"iconPath"`
		PluralName                     string               `json:"pluralName"`
		Name                           string               `json:"name"`
		CapacityUnit                   ValidMeasurementUnit `json:"capacityUnit"`
		LengthInMillimeters            float64              `json:"lengthInMillimeters"`
		Capacity                       float64              `json:"capacity"`
		WidthInMillimeters             float64              `json:"widthInMillimeters"`
		HeightInMillimeters            float64              `json:"heightInMillimeters"`
		DisplayInSummaryLists          bool                 `json:"displayInSummaryLists"`
		UsableForStorage               bool                 `json:"usableForStorage"`
		IncludeInGeneratedInstructions bool                 `json:"includeInGeneratedInstructions"`
	}
)
