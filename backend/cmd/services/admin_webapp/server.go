package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/pages"
	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"

	"github.com/jellydator/ttlcache/v3"
)

type server struct {
	tracer         tracing.Tracer
	logger         logging.Logger
	apiServerURL   *url.URL
	cookieManager  cookies.Manager
	pageBuilder    *pages.Builder
	apiClientCache *ttlcache.Cache[string, *apiclient.Client]
	oauth2APIClientID,
	oauth2APIClientSecret string
}

type config struct {
	Cookies                cookies.Config `env:"init"                      envPrefix:"COOKIES_"          json:"cookies"`
	APIServerURL           string         `env:"API_SERVER_URL"            json:"apiServerURL"`
	OAuth2APIClientID      string         `env:"OAUTH2_API_CLIENT_ID"      json:"oauth2APIClientID"`
	OAuth2APIClientSecret  string         `env:"OAUTH2_API_CLIENT_SECRET"  json:"oauth2APIClientSecret"`
	APIClientCacheCapacity uint64         `env:"API_CLIENT_CACHE_CAPACITY" json:"apiClientCacheCapacity"`
	APIClientCacheTTL      time.Duration  `env:"API_CLIENT_CACHE_TTL"      json:"apiClientCacheTTL"`
}

func newServer(cfg *config, logger logging.Logger, tracerProvider tracing.TracerProvider) (*server, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError("frontend admin service")
	}

	parsedURL, err := url.Parse(cfg.APIServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse api server URL: %w", err)
	}

	cookieManager, err := cookies.NewCookieManager(&cfg.Cookies, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("setting up cookie manager: %w", err)
	}

	s := &server{
		tracer:        tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("admin_webapp")),
		logger:        logging.EnsureLogger(logger).WithName("admin_webapp"),
		apiServerURL:  parsedURL,
		cookieManager: cookieManager,
		pageBuilder:   pages.NewPageBuilder(tracerProvider, logger, parsedURL),
		apiClientCache: ttlcache.New[string, *apiclient.Client](
			ttlcache.WithCapacity[string, *apiclient.Client](cfg.APIClientCacheCapacity),
			ttlcache.WithTTL[string, *apiclient.Client](cfg.APIClientCacheTTL),
		),
		oauth2APIClientID:     cfg.OAuth2APIClientID,
		oauth2APIClientSecret: cfg.OAuth2APIClientSecret,
	}

	return s, nil
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)

		cookie, err := req.Cookie(cookieName)
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
		if err = s.cookieManager.Decode(ctx, cookieName, cookie.Value, &usd); err != nil {
			logger.Error("decoding cookie", err)
			next.ServeHTTP(res, req)
			return
		}

		logger.WithValue("user.id", usd.UserID).Info("user session retrieved from middleware")

		var client *apiclient.Client
		if cachedClientItem := apiClientCache.Get(usd.UserID); cachedClientItem == nil {
			client, err = apiclient.NewClient(
				apiServerURL,
				tracing.NewNoopTracerProvider(),
				apiclient.UsingOAuth2(
					ctx,
					oauth2ClientID,
					oauth2ClientSecret,
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
