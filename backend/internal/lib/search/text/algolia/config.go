package algolia

import (
	"time"
)

type Config struct {
	AppID   string        `env:"APP_ID"  json:"appID"`
	APIKey  string        `env:"API_KEY" json:"writeAPIKey"`
	Timeout time.Duration `env:"TIMEOUT" json:"timeout"`
}
