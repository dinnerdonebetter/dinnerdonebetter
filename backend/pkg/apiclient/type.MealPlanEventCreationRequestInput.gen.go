// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlanEventCreationRequestInput struct {
		EndsAt   string                               `json:"endsAt"`
		MealName string                               `json:"mealName"`
		Notes    string                               `json:"notes"`
		StartsAt string                               `json:"startsAt"`
		Options  []MealPlanOptionCreationRequestInput `json:"options"`
	}
)
