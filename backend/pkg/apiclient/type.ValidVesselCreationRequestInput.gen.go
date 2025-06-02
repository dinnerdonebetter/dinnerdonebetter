// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidVesselCreationRequestInput struct {
		Slug                           string  `json:"slug"`
		CapacityUnitID                 string  `json:"capacityUnitID"`
		Description                    string  `json:"description"`
		IconPath                       string  `json:"iconPath"`
		Name                           string  `json:"name"`
		PluralName                     string  `json:"pluralName"`
		Shape                          string  `json:"shape"`
		Capacity                       float64 `json:"capacity"`
		HeightInMillimeters            float64 `json:"heightInMillimeters"`
		WidthInMillimeters             float64 `json:"widthInMillimeters"`
		LengthInMillimeters            float64 `json:"lengthInMillimeters"`
		DisplayInSummaryLists          bool    `json:"displayInSummaryLists"`
		UsableForStorage               bool    `json:"usableForStorage"`
		IncludeInGeneratedInstructions bool    `json:"includeInGeneratedInstructions"`
	}
)
