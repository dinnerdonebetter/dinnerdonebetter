package elasticsearch

import (
	"time"
)

type Config struct {
	Address               string        `json:"address"               toml:"address,omitempty"`
	Username              string        `json:"username"              toml:"username,omitempty"`
	Password              string        `json:"password"              toml:"password,omitempty"`
	CACert                []byte        `json:"caCert"                toml:"ca_cert,omitempty"`
	IndexOperationTimeout time.Duration `json:"indexOperationTimeout" toml:"index_operation_timeout,omitempty"`
}
