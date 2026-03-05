package main

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/webappauth"
	"github.com/dinnerdonebetter/backend/pkg/client"
)

func (s *AdminFrontendServer) authMiddleware(handler http.Handler) http.Handler {
	buildClient := func(ctx context.Context, accessToken string) (client.Client, error) {
		return webappauth.BuildAuthedClient(ctx, s.config.APIServiceConnection, accessToken, s.developingLocally)
	}
	return webappauth.AuthMiddleware(&webappauth.AuthMiddlewareOpts{
		CookieManager: s.cookieManager,
		CookieConfig:  s.config.Cookies,
		BuildClient:   buildClient,
		RedirectPath:  "/login",
		Logger:        s.logger,
		Tracer:        s.tracer,
	})(handler)
}
