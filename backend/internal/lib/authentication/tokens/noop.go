package tokens

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/users"
)

type noopTokenIssuer struct{}

// IssueToken implements the interface.
func (n *noopTokenIssuer) IssueToken(context.Context, users.User, time.Duration) (string, error) {
	return "", nil
}

// ParseUserIDFromToken implements the interface.
func (n *noopTokenIssuer) ParseUserIDFromToken(context.Context, string) (string, error) {
	return "", nil
}

func NewNoopTokenIssuer() Issuer {
	return &noopTokenIssuer{}
}
