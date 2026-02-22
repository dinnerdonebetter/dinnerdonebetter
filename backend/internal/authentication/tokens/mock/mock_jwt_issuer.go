package mocktokens

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"

	"github.com/stretchr/testify/mock"
)

var _ tokens.Issuer = (*Issuer)(nil)

type Issuer struct {
	mock.Mock
}

func (m *Issuer) IssueToken(ctx context.Context, user tokens.User, expiry time.Duration) (string, error) {
	return m.IssueTokenWithAccount(ctx, user, expiry, "")
}

func (m *Issuer) IssueTokenWithAccount(ctx context.Context, user tokens.User, expiry time.Duration, accountID string) (string, error) {
	returnValues := m.Called(ctx, user, expiry, accountID)
	return returnValues.String(0), returnValues.Error(1)
}

func (m *Issuer) ParseUserIDFromToken(ctx context.Context, token string) (string, error) {
	userID, _, err := m.ParseUserIDAndAccountIDFromToken(ctx, token)
	return userID, err
}

func (m *Issuer) ParseUserIDAndAccountIDFromToken(ctx context.Context, token string) (userID, accountID string, err error) {
	returnValues := m.Called(ctx, token)
	return returnValues.String(0), returnValues.String(1), returnValues.Error(2)
}
