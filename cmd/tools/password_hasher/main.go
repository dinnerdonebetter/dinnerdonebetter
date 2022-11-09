package main

import (
	"context"
	"os"

	"github.com/prixfixeco/backend/internal/authentication"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

func main() {
	logger := logging.NewNoopLogger()

	hasher := authentication.ProvideArgon2Authenticator(logger, tracing.NewNoopTracerProvider())
	hashed, err := hasher.HashPassword(context.Background(), os.Args[1])

	if err != nil {
		panic(err)
	}

	println(hashed)
}
