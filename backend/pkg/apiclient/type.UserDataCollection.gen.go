// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	UserDataCollection struct {
		Core     CoreUserDataCollection   `json:"core"`
		User     User                     `json:"user"`
		ReportID string                   `json:"reportID"`
		Eating   EatingUserDataCollection `json:"eating"`
	}
)
