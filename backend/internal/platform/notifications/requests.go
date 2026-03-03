package notifications

// MealPlanTaskNotificationRequest is the message payload for a meal plan task push notification request.
type MealPlanTaskNotificationRequest struct {
	MealPlanTaskID string `json:"mealPlanTaskID"`
}
