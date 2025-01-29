// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
MealPlan struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToHousehold string `json:"belongsToHousehold"`
 CreatedAt string `json:"createdAt"`
 CreatedBy string `json:"createdBy"`
 ElectionMethod string `json:"electionMethod"`
 Events []MealPlanEvent `json:"events"`
 GroceryListInitialized bool `json:"groceryListInitialized"`
 ID string `json:"id"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Notes string `json:"notes"`
 Status string `json:"status"`
 TasksCreated bool `json:"tasksCreated"`
 VotingDeadline string `json:"votingDeadline"`

}
)
