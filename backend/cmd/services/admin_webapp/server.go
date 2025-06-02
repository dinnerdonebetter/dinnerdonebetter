package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/pages"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"

	"github.com/jellydator/ttlcache/v3"
)

type Server struct {
	tracer         tracing.Tracer
	tracerProvider tracing.TracerProvider
	logger         logging.Logger
	apiServerURL   *url.URL
	cookieManager  cookies.Manager
	pageBuilder    *pages.Builder
	apiClientCache *ttlcache.Cache[string, *apiclient.Client]
	oauth2APIClientID,
	oauth2APIClientSecret string
}

type adminWebappCfg struct {
	Cookies                cookies.Config `env:"init"                      envPrefix:"COOKIES_"          json:"cookies"`
	APIServerURL           string         `env:"API_SERVER_URL"            json:"apiServerURL"`
	OAuth2APIClientID      string         `env:"OAUTH2_API_CLIENT_ID"      json:"oauth2APIClientID"`
	OAuth2APIClientSecret  string         `env:"OAUTH2_API_CLIENT_SECRET"  json:"oauth2APIClientSecret"`
	APIClientCacheCapacity uint64         `env:"API_CLIENT_CACHE_CAPACITY" json:"apiClientCacheCapacity"`
	APIClientCacheTTL      time.Duration  `env:"API_CLIENT_CACHE_TTL"      json:"apiClientCacheTTL"`
}

func NewServer(cfg *adminWebappCfg, logger logging.Logger, tracerProvider tracing.TracerProvider) (*Server, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError("frontend admin service")
	}

	parsedURL, err := url.Parse(cfg.APIServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse api Server URL: %w", err)
	}

	cookieManager, err := cookies.NewCookieManager(&cfg.Cookies, tracerProvider)
	if err != nil {
		return nil, fmt.Errorf("setting up cookie manager: %w", err)
	}

	s := &Server{
		tracerProvider: tracing.EnsureTracerProvider(tracerProvider),
		tracer:         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("admin_webapp")),
		logger:         logging.EnsureLogger(logger).WithName("admin_webapp"),
		apiServerURL:   parsedURL,
		cookieManager:  cookieManager,
		pageBuilder:    pages.NewPageBuilder(tracerProvider, logger, parsedURL),
		apiClientCache: ttlcache.New[string, *apiclient.Client](
			ttlcache.WithCapacity[string, *apiclient.Client](cfg.APIClientCacheCapacity),
			ttlcache.WithTTL[string, *apiclient.Client](cfg.APIClientCacheTTL),
		),
		oauth2APIClientID:     cfg.OAuth2APIClientID,
		oauth2APIClientSecret: cfg.OAuth2APIClientSecret,
	}

	return s, nil
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
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
