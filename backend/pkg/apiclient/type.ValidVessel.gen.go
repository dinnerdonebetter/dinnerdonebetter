// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidVessel struct {
		ID                             string               `json:"id"`
		PluralName                     string               `json:"pluralName"`
		Slug                           string               `json:"slug"`
		CreatedAt                      string               `json:"createdAt"`
		Description                    string               `json:"description"`
		Shape                          string               `json:"shape"`
		ArchivedAt                     string               `json:"archivedAt"`
		IconPath                       string               `json:"iconPath"`
		Name                           string               `json:"name"`
		LastUpdatedAt                  string               `json:"lastUpdatedAt"`
		CapacityUnit                   ValidMeasurementUnit `json:"capacityUnit"`
		LengthInMillimeters            float64              `json:"lengthInMillimeters"`
		Capacity                       float64              `json:"capacity"`
		HeightInMillimeters            float64              `json:"heightInMillimeters"`
		WidthInMillimeters             float64              `json:"widthInMillimeters"`
		IncludeInGeneratedInstructions bool                 `json:"includeInGeneratedInstructions"`
		DisplayInSummaryLists          bool                 `json:"displayInSummaryLists"`
		UsableForStorage               bool                 `json:"usableForStorage"`
	}
)
