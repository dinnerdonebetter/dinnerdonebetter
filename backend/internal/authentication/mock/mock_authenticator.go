package mockauthn

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"

	"github.com/stretchr/testify/mock"
)

var _ authentication.Authenticator = (*Authenticator)(nil)

// Authenticator is a mock Authenticator.
type Authenticator struct {
	mock.Mock
}

// ValidateLogin satisfies our authenticator interface.
func (m *Authenticator) CredentialsAreValid(ctx context.Context, hash, password, totpSecret, totpCode string) (bool, error) {
	args := m.Called(ctx, hash, password, totpSecret, totpCode)

	return args.Bool(0), args.Error(1)
}

// HashPassword satisfies our authenticator interface.
func (m *Authenticator) HashPassword(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}
