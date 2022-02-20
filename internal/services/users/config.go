package users

type (
	// Config configures the users service.
	Config struct {
		DataChangesTopicName string `json:"dataChangesTopicName,omitempty" mapstructure:"data_changes_topic_name" toml:"data_changes_topic_name,omitempty"`
	}
)
