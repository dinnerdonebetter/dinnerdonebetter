package tokens

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
)

type (
	Issuer interface {
		IssueToken(ctx context.Context, user authentication.User, expiry time.Duration) (string, error)
		ParseUserIDFromToken(ctx context.Context, token string) (string, error)
	}
)
