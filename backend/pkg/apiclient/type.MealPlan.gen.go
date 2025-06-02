// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlan struct {
		Notes                  string          `json:"notes"`
		BelongsToHousehold     string          `json:"belongsToHousehold"`
		CreatedAt              string          `json:"createdAt"`
		CreatedBy              string          `json:"createdBy"`
		ElectionMethod         string          `json:"electionMethod"`
		ID                     string          `json:"id"`
		LastUpdatedAt          string          `json:"lastUpdatedAt"`
		ArchivedAt             string          `json:"archivedAt"`
		Status                 string          `json:"status"`
		VotingDeadline         string          `json:"votingDeadline"`
		Events                 []MealPlanEvent `json:"events"`
		GroceryListInitialized bool            `json:"groceryListInitialized"`
		TasksCreated           bool            `json:"tasksCreated"`
	}
)
