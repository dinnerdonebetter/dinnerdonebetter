package elasticsearch

import (
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Address               string        `json:"address" mapstructure:"address" toml:"address,omitempty"`
	Username              string        `json:"username" mapstructure:"username" toml:"username,omitempty"`
	Password              string        `json:"password" mapstructure:"password" toml:"password,omitempty"`
	IndexOperationTimeout time.Duration `json:"indexOperationTimeout" mapstructure:"index_operation_timeout" toml:"index_operation_timeout,omitempty"`
}

func (cfg *Config) provideElasticsearchClient() (*elasticsearch.Client, error) {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			cfg.Address,
		},
		Username:             cfg.Username,
		Password:             cfg.Password,
		RetryOnStatus:        nil,
		EnableRetryOnTimeout: true,
		MaxRetries:           10,
		Transport:            nil,
		Logger:               nil,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing search client: %w", err)
	}

	return c, nil
}
