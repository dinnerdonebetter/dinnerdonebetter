// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ServiceSettingConfiguration struct {
		ArchivedAt         string         `json:"archivedAt"`
		BelongsToHousehold string         `json:"belongsToHousehold"`
		BelongsToUser      string         `json:"belongsToUser"`
		CreatedAt          string         `json:"createdAt"`
		ID                 string         `json:"id"`
		LastUpdatedAt      string         `json:"lastUpdatedAt"`
		Notes              string         `json:"notes"`
		Value              string         `json:"value"`
		ServiceSetting     ServiceSetting `json:"serviceSetting"`
	}
)
