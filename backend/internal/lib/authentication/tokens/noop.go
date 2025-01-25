package tokens

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
)

type noopTokenIssuer struct{}

// IssueToken implements the interface.
func (n *noopTokenIssuer) IssueToken(context.Context, authentication.User, time.Duration) (string, error) {
	return "", nil
}

// ParseUserIDFromToken implements the interface.
func (n *noopTokenIssuer) ParseUserIDFromToken(context.Context, string) (string, error) {
	return "", nil
}

func NewNoopTokenIssuer() Issuer {
	return &noopTokenIssuer{}
}
