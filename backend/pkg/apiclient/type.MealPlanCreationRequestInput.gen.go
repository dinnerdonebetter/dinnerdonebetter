// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
MealPlanCreationRequestInput struct {
   ElectionMethod string `json:"electionMethod"`
 Events []MealPlanEventCreationRequestInput `json:"events"`
 Notes string `json:"notes"`
 VotingDeadline string `json:"votingDeadline"`

}
)
