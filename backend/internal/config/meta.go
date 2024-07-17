package config

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// MetaSettings is primarily used for development.
type MetaSettings struct {
	RunMode runMode `json:"runMode" toml:"run_mode,omitempty"`
	Debug   bool    `json:"debug"   toml:"debug,omitempty"`
}

var _ validation.ValidatableWithContext = (*MetaSettings)(nil)

// ValidateWithContext validates an MetaSettings struct.
func (s MetaSettings) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &s,
		validation.Field(&s.RunMode, validation.In(TestingRunMode, DevelopmentRunMode, ProductionRunMode)),
	)
}
