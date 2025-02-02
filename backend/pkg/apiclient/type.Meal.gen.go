// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	Meal struct {
		ArchivedAt           string                      `json:"archivedAt"`
		CreatedAt            string                      `json:"createdAt"`
		CreatedByUser        string                      `json:"createdByUser"`
		Description          string                      `json:"description"`
		ID                   string                      `json:"id"`
		LastUpdatedAt        string                      `json:"lastUpdatedAt"`
		Name                 string                      `json:"name"`
		Components           []MealComponent             `json:"components"`
		EstimatedPortions    Float32RangeWithOptionalMax `json:"estimatedPortions"`
		EligibleForMealPlans bool                        `json:"eligibleForMealPlans"`
	}
)
