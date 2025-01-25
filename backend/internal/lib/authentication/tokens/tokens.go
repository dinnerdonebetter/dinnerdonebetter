package tokens

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

type (
	Issuer interface {
		IssueToken(ctx context.Context, user *types.User, expiry time.Duration) (string, error)
		ParseUserIDFromToken(ctx context.Context, token string) (string, error)
	}
)
