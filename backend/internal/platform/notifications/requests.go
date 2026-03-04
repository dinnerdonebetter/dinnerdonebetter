package notifications

const (
	// MobileNotificationRequestTypeMealPlanTask indicates a meal plan task reminder notification.
	MobileNotificationRequestTypeMealPlanTask = "meal_plan_task"
	// MobileNotificationRequestTypeHouseholdInvitationAccepted indicates a household invitation was accepted.
	MobileNotificationRequestTypeHouseholdInvitationAccepted = "household_invitation_accepted"

	// ExcludedUserIDContextKey is the key used in MobileNotificationRequest.Context for the user ID to exclude from recipients.
	ExcludedUserIDContextKey = "excludedUserID"

	// MealPlanTaskIDContextKey is the key used in MobileNotificationRequest.Context for meal plan task ID.
	// When present, the handler uses it for idempotency (MealPlanTaskNotificationHasBeenSent) and
	// marking the notification as sent (MarkMealPlanTaskNotificationSent).
	MealPlanTaskIDContextKey = "mealPlanTaskID"
)

// MobileNotificationRequest is the generic message payload for mobile push notifications.
// RequestType determines which handler processes the request; schedulers format the message.
type MobileNotificationRequest struct {
	RequestType      string            `json:"requestType"`
	Context          map[string]string `json:"context,omitempty"`
	Title            string            `json:"title"`
	Body             string            `json:"body"`
	RecipientUserIDs []string          `json:"recipientUserIDs"`
}
