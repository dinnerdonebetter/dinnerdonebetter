package cookies

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Base64EncodedHashKey  string `env:"HASH_KEY"  json:"base64EncodedHashKey"`
	Base64EncodedBlockKey string `env:"BLOCK_KEY" json:"base64EncodedBlockKey"`
}

func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Base64EncodedHashKey, validation.Required),
		validation.Field(&c.Base64EncodedBlockKey, validation.Required),
	)
}
