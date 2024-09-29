package mockauthn

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

var _ authentication.JWTSigner = (*JWTIssuer)(nil)

type JWTIssuer struct {
	mock.Mock
}

func (m *JWTIssuer) IssueJWT(ctx context.Context, user *types.User, expiry time.Duration) (string, error) {
	returnValues := m.Called(ctx, user, expiry)
	return returnValues.String(0), returnValues.Error(1)
}

func (m *JWTIssuer) ParseJWT(ctx context.Context, token string) (*jwt.Token, error) {
	returnValues := m.Called(ctx, token)
	return returnValues.Get(0).(*jwt.Token), returnValues.Error(1)
}
