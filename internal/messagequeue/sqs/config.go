package sqs

// Config configures a SQS-backed consumer.
type Config struct {
	QueueAddress string `json:"messageQueueAddress" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
}
