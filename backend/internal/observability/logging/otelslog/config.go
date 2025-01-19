package otelslog

import (
	"time"
)

type (
	Config struct {
		CollectorEndpoint string        `env:"ENDPOINT_URL" json:"endpointURL"`
		Insecure          bool          `env:"INSECURE"     json:"insecure"`
		Timeout           time.Duration `env:"TIMEOUT"      json:"timeout"`
	}
)
