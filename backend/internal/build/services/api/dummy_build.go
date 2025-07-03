package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	"github.com/dinnerdonebetter/backend/internal/platform/server/http"
)

type dummyServer struct{}

func (d dummyServer) Serve() {
}

func (d dummyServer) Shutdown(ctx context.Context) error {
	return nil
}

func (d dummyServer) Router() routing.Router {
	return nil
}

func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (http.Server, error) {
	return &dummyServer{}, nil
}
