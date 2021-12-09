package search

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ElasticsearchProvider represents the elasticsearch search index provider.
	ElasticsearchProvider = "elasticsearch"
)

// Config contains settings regarding search indices.
type Config struct {
	_ struct{}

	Provider string    `json:"provider" mapstructure:"provider" toml:"provider,omitempty"`
	Name     IndexName `json:"indexName" mapstructure:"index_name" toml:"index_name,omitempty"`
	Address  IndexPath `json:"address" mapstructure:"address" toml:"address,omitempty"`
	Username string    `json:"username" mapstructure:"username" toml:"username,omitempty"`
	Password string    `json:"password" mapstructure:"password" toml:"password,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(ElasticsearchProvider)),
	)
}
