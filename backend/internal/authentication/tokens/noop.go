package tokens

import (
	"context"
	"time"
)

type noopTokenIssuer struct{}

// IssueToken implements the interface.
func (n *noopTokenIssuer) IssueToken(context.Context, User, time.Duration) (string, error) {
	return "", nil
}

// IssueTokenWithAccount implements the interface.
func (n *noopTokenIssuer) IssueTokenWithAccount(context.Context, User, time.Duration, string) (string, error) {
	return "", nil
}

// ParseUserIDFromToken implements the interface.
func (n *noopTokenIssuer) ParseUserIDFromToken(context.Context, string) (string, error) {
	return "", nil
}

// ParseUserIDAndAccountIDFromToken implements the interface.
func (n *noopTokenIssuer) ParseUserIDAndAccountIDFromToken(context.Context, string) (userID, accountID string, err error) {
	return "", "", nil
}

func NewNoopTokenIssuer() Issuer {
	return &noopTokenIssuer{}
}
