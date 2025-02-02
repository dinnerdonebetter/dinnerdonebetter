// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlanTask struct {
		AssignedToUser      string         `json:"assignedToUser"`
		CompletedAt         string         `json:"completedAt"`
		CreatedAt           string         `json:"createdAt"`
		CreationExplanation string         `json:"creationExplanation"`
		ID                  string         `json:"id"`
		LastUpdatedAt       string         `json:"lastUpdatedAt"`
		Status              string         `json:"status"`
		StatusExplanation   string         `json:"statusExplanation"`
		MealPlanOption      MealPlanOption `json:"mealPlanOption"`
		RecipePrepTask      RecipePrepTask `json:"recipePrepTask"`
	}
)
