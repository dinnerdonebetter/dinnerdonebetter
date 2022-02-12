package pubsub

// Config configures a PubSub-backed consumer.
type Config struct {
	TopicName string `json:"topicName" mapstructure:"topic_name" toml:"topic_name,omitempty"`
}
