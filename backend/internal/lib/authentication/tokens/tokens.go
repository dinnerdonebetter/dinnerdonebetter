package tokens

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/users"
)

type (
	Issuer interface {
		IssueToken(ctx context.Context, user users.User, expiry time.Duration) (string, error)
		ParseUserIDFromToken(ctx context.Context, token string) (string, error)
	}

	// TokenResponse is used to return a JWT to a user.
	TokenResponse struct {
		_            struct{}  `json:"-"`
		ExpiresUTC   time.Time `json:"expires"`
		UserID       string    `json:"userID"`
		HouseholdID  string    `json:"householdID"`
		AccessToken  string    `json:"accessToken"`
		RefreshToken string    `json:"refreshToken"`
	}
)
