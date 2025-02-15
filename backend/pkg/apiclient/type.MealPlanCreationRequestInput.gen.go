// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlanCreationRequestInput struct {
		ElectionMethod string                              `json:"electionMethod"`
		Notes          string                              `json:"notes"`
		VotingDeadline string                              `json:"votingDeadline"`
		Events         []MealPlanEventCreationRequestInput `json:"events"`
	}
)
