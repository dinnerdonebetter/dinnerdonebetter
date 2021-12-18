package redis

// Config configures a Redis-backed consumer.
type Config struct {
	Username     string `json:"username" mapstructure:"username" toml:"username,omitempty"`
	Password     string `json:"password" mapstructure:"password" toml:"password,omitempty"`
	QueueAddress string `json:"messageQueueAddress" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
	DB           int    `json:"database" mapstructure:"database" toml:"database,omitempty"`
}
