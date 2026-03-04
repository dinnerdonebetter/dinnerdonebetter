package webappauth

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/pkg/client"
)

type ContextKey string

const apiClientContextKey ContextKey = "api_client"

// ClientFromContext retrieves the authenticated API client from the request context.
func ClientFromContext(ctx context.Context) (client.Client, error) {
	c, ok := ctx.Value(apiClientContextKey).(client.Client)
	if !ok {
		return nil, errors.New("no api client found in context")
	}
	return c, nil
}
