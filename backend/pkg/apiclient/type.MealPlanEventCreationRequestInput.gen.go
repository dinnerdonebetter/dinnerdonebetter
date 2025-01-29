// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
MealPlanEventCreationRequestInput struct {
   EndsAt string `json:"endsAt"`
 MealName string `json:"mealName"`
 Notes string `json:"notes"`
 Options []MealPlanOptionCreationRequestInput `json:"options"`
 StartsAt string `json:"startsAt"`

}
)
