package algolia

import (
	"time"
)

type Config struct {
	AppID   string        `json:"appID"       mapstructure:"app_id"        toml:"app_id,omitempty"`
	APIKey  string        `json:"writeAPIKey" mapstructure:"write_api_key" toml:"write_api_key,omitempty"`
	Timeout time.Duration `json:"timeout"     mapstructure:"timeout"       toml:"timeout,omitempty"`
}
