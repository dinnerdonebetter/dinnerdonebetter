package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/panicking"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
)

const (
	serverNamespace = "dinner_done_better_api"
	loggerName      = "api_server"
)

type (
	Server interface {
		Serve()
		Shutdown(context.Context) error
		Router() routing.Router
	}

	// server is our API http server.
	server struct {
		logger         logging.Logger
		router         routing.Router
		panicker       panicking.Panicker
		httpServer     *http.Server
		tracerProvider tracing.TracerProvider
		config         Config
	}
)

// ProvideHTTPServer builds a new server instance.
func ProvideHTTPServer(
	serverSettings Config,
	logger logging.Logger,
	router routing.Router,
	tracerProvider tracing.TracerProvider,
) (Server, error) {
	srv := &server{
		config: serverSettings,

		// infra things,
		router:         router,
		logger:         logging.EnsureLogger(logger).WithName(loggerName),
		panicker:       panicking.NewProductionPanicker(),
		httpServer:     provideStdLibHTTPServer(serverSettings.HTTPPort),
		tracerProvider: tracerProvider,
	}

	logger.Debug("HTTP server successfully constructed")

	return srv, nil
}

// Router returns the router.
func (s *server) Router() routing.Router {
	return s.router
}

// Shutdown shuts down the server.
func (s *server) Shutdown(ctx context.Context) error {
	if err := s.tracerProvider.ForceFlush(ctx); err != nil {
		s.logger.Error("flushing traces", err)
	}

	return s.httpServer.Shutdown(ctx)
}

// Serve serves HTTP traffic.
func (s *server) Serve() {
	s.logger.Debug("setting up server")

	s.httpServer.Handler = otelhttp.NewHandler(s.router.Handler(), serverNamespace, otelhttp.WithSpanNameFormatter(tracing.FormatSpan))

	http2ServerConf := &http2.Server{}
	if err := http2.ConfigureServer(s.httpServer, http2ServerConf); err != nil {
		s.logger.Error("configuring HTTP2", err)
		s.panicker.Panic(err)
	}

	if s.config.TLSCertificateFile != "" && s.config.TLSCertificateKeyFile != "" {
		s.logger.WithValue("listening_on", s.httpServer.Addr).Info("Listening for HTTPS requests")
		// returns ErrServerClosed on graceful close.
		if err := s.httpServer.ListenAndServeTLS(s.config.TLSCertificateFile, s.config.TLSCertificateKeyFile); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop.
				os.Exit(0)
			}

			s.logger.Error("shutting server down", err)
		}
	} else {
		s.logger.WithValue("listening_on", s.httpServer.Addr).Info("Listening for HTTP requests")
		// returns ErrServerClosed on graceful close.
		if err := s.httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop.
				os.Exit(0)
			}

			s.logger.Error("shutting server down", err)
		}
	}
}

const (
	maxTimeout   = 120 * time.Second
	readTimeout  = 5 * time.Second
	writeTimeout = 2 * readTimeout
	idleTimeout  = maxTimeout
)

// provideStdLibHTTPServer provides an HTTP httpServer.
func provideStdLibHTTPServer(port uint16) *http.Server {
	// heavily inspired by https://blog.cloudflare.com/exposing-go-on-the-internet/
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		TLSConfig: &tls.Config{
			// "Only use curves which have assembly implementations"
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	return srv
}
