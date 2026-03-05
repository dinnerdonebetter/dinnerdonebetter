package webappauth

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/client"
)

// BuildClientFunc builds an authenticated gRPC client from an access token.
type BuildClientFunc func(ctx context.Context, accessToken string) (client.Client, error)

// AuthMiddlewareOpts configures the auth middleware.
type AuthMiddlewareOpts struct {
	CookieManager cookies.Manager
	Logger        logging.Logger
	Tracer        tracing.Tracer
	BuildClient   BuildClientFunc
	RedirectPath  string
	CookieConfig  cookies.Config
}

// AuthMiddleware returns an HTTP middleware that reads the auth cookie, decodes the token,
// builds an authenticated client, and injects it into the request context.
// On any failure, redirects to RedirectPath.
func AuthMiddleware(opts *AuthMiddlewareOpts) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx, span := opts.Tracer.StartSpan(req.Context())
			defer span.End()

			logger := opts.Logger.WithRequest(req)

			cookieName := opts.CookieConfig.CookieName
			cookie, err := req.Cookie(cookieName)
			if err != nil {
				logger.Error("no cookie found", err)
				http.Redirect(res, req, opts.RedirectPath, http.StatusFound)
				return
			}
			if cookie == nil {
				logger.Debug("no cookie found")
				http.Redirect(res, req, opts.RedirectPath, http.StatusFound)
				return
			}

			var payload AuthPayload
			if err = opts.CookieManager.Decode(ctx, cookieName, cookie.Value, &payload); err != nil {
				logger.Error("decoding cookie", err)
				http.Redirect(res, req, opts.RedirectPath, http.StatusFound)
				return
			}

			c, err := opts.BuildClient(ctx, payload.AccessToken)
			if err != nil {
				logger.Error("building client", err)
				http.Redirect(res, req, opts.RedirectPath, http.StatusFound)
				return
			}

			handler.ServeHTTP(res, req.WithContext(context.WithValue(ctx, apiClientContextKey, c)))
		})
	}
}
