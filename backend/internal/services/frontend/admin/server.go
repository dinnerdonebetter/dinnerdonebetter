package admin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/frontend/admin/pages"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"

	"github.com/jellydator/ttlcache/v3"
)

type WebappServer struct {
	tracer         tracing.Tracer
	tracerProvider tracing.TracerProvider
	logger         logging.Logger
	apiServerURL   *url.URL
	cookieManager  cookies.Manager
	cookiesConfig  cookies.Config
	pageBuilder    *pages.Builder
	apiClientCache *ttlcache.Cache[string, *apiclient.Client]
	oauth2APIClientID,
	oauth2APIClientSecret string
}

func NewServer(cfg *config.AdminWebappConfig, logger logging.Logger, tracerProvider tracing.TracerProvider) (*WebappServer, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError("frontend admin service")
	}

	parsedURL, err := url.Parse(cfg.APIServiceConnection.APIServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse api server URL: %w", err)
	}

	cookieManager, err := cookies.NewCookieManager(&cfg.Cookies, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("setting up cookie manager: %w", err)
	}

	s := &WebappServer{
		tracerProvider: tracing.EnsureTracerProvider(tracerProvider),
		tracer:         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("admin_webapp")),
		logger:         logging.EnsureLogger(logger).WithName("admin_webapp"),
		apiServerURL:   parsedURL,
		cookiesConfig:  cfg.Cookies,
		cookieManager:  cookieManager,
		pageBuilder:    pages.NewPageBuilder(tracerProvider, logger, parsedURL),
		apiClientCache: ttlcache.New[string, *apiclient.Client](
			ttlcache.WithCapacity[string, *apiclient.Client](cfg.APIClientCache.CacheCapacity),
			ttlcache.WithTTL[string, *apiclient.Client](cfg.APIClientCache.CacheTTL),
		),
		oauth2APIClientID:     cfg.APIServiceConnection.OAuth2APIClientID,
		oauth2APIClientSecret: cfg.APIServiceConnection.OAuth2APIClientSecret,
	}

	return s, nil
}

func (s *WebappServer) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)

		cookie, err := req.Cookie(s.cookiesConfig.CookieName)
		if err != nil {
			logger.Error("fetching request cookie", err)
			next.ServeHTTP(res, req)
			return
		} else if cookie == nil {
			logger.Info("cookie was nil!")
			next.ServeHTTP(res, req)
			return
		}

		var usd *userSessionDetails
		if err = s.cookieManager.Decode(ctx, s.cookiesConfig.CookieName, cookie.Value, &usd); err != nil {
			logger.Error("decoding cookie", err)
			next.ServeHTTP(res, req)
			return
		}

		logger.WithValue("user.id", usd.UserID).Info("user session retrieved from middleware")

		var client *apiclient.Client
		if cachedClientItem := s.apiClientCache.Get(usd.UserID); cachedClientItem == nil {
			client, err = apiclient.NewClient(
				s.apiServerURL,
				tracing.NewNoopTracerProvider(),
				apiclient.UsingOAuth2(
					ctx,
					s.oauth2APIClientID,
					s.oauth2APIClientSecret,
					[]string{authorization.HouseholdAdminRoleName},
					usd.Token,
				),
			)
			if err != nil {
				logger.Error("establishing API client", err)
				next.ServeHTTP(res, req)
				return
			}
		} else {
			client = cachedClientItem.Value()
		}

		req = req.WithContext(context.WithValue(ctx, apiClientContextKey, client))

		next.ServeHTTP(res, req)
	})
}
