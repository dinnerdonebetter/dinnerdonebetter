package mocktokens

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens"

	"github.com/stretchr/testify/mock"
)

var _ tokens.Issuer = (*Issuer)(nil)

type Issuer struct {
	mock.Mock
}

func (m *Issuer) IssueToken(ctx context.Context, user authentication.User, expiry time.Duration) (string, error) {
	returnValues := m.Called(ctx, user, expiry)
	return returnValues.String(0), returnValues.Error(1)
}

func (m *Issuer) ParseUserIDFromToken(ctx context.Context, token string) (string, error) {
	returnValues := m.Called(ctx, token)
	return returnValues.String(0), returnValues.Error(1)
}
