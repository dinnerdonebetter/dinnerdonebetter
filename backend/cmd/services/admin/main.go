package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/config/envvars"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	adminServerConfigurationFilepath = "deploy/environments/localdev/config_files/admin_webapp_config.json"

	o11yName      = "admin_frontend"
	adminPassword = "admin_pass"
)

var (
	premadeAdminUser = &identity.User{
		ID:              identifiers.New(),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@example.email",
		Username:        "admin_user",
		HashedPassword:  adminPassword,
	}
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx := context.Background()

	// We don't yet have a way to write these values into the AdminWebappConfig (because they're not present in the root APIConfig struct).
	// This approach is an atrocious hack that I have to employ because I wasn't smart enough to design a better config generation system.
	must(os.Setenv(envvars.CookiesCookieNameEnvVarKey, "dev_admin_frontend"))
	must(os.Setenv(envvars.CookiesHashKeyEnvVarKey, base64.RawURLEncoding.EncodeToString([]byte("HEREISA32CHARSECRETWHICHISMADEUP"))))
	must(os.Setenv(envvars.CookiesBlockKeyEnvVarKey, base64.RawURLEncoding.EncodeToString([]byte("HEREISA32CHARSECRETWHICHISMADEUP"))))
	must(os.Setenv(envvars.APIServiceHTTPAPIServerURLEnvVarKey, "http://localhost:8000"))
	must(os.Setenv(envvars.APIServiceGrpcAPIServerURLEnvVarKey, ":8001"))
	must(os.Setenv(envvars.APIServiceOauth2APIClientIDEnvVarKey, strings.Repeat("A", oauth.ClientIDSize)))
	must(os.Setenv(envvars.APIServiceOauth2APIClientSecretEnvVarKey, strings.Repeat("A", oauth.ClientSecretSize)))
	must(os.Setenv(envvars.ServerPortEnvVarKey, "8888"))
	must(os.Setenv(envvars.ServerStartupDeadlineEnvVarKey, time.Minute.String()))
	must(os.Setenv(envvars.CookiesDomainEnvVarKey, "localhost"))
	must(os.Setenv(envvars.CookiesLifetimeEnvVarKey, time.Hour.String()))

	cfg, err := config.LoadConfigFromPath[config.AdminWebappConfig](ctx, adminServerConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	logger, tracerProvider, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		log.Fatal(err)
	}

	grpcServerAddr := cfg.APIServiceConnection.GRPCAPIServerURL

	code, err := premadeAdminUser.GenerateTOTPCode()
	if err != nil {
		log.Fatal(err)
	}

	loginInput := &authsvc.UserLoginInput{
		Username:  premadeAdminUser.Username,
		Password:  adminPassword,
		TOTPToken: code,
	}

	token, err := localdev.FetchLoginTokenForUser(ctx, grpcServerAddr, loginInput)
	if err != nil {
		log.Fatal(err)
	}

	apiClient, err := localdev.BuildInsecureOAuthedGRPCClient(
		ctx,
		strings.Repeat("A", 16),
		strings.Repeat("A", 16),
		cfg.APIServiceConnection.HTTPAPIServerURL,
		grpcServerAddr,
		token,
	)
	if err != nil {
		log.Fatal(err)
	}

	fs, err := NewAdminFrontendServer(
		ctx,
		apiClient,
		logger,
		tracerProvider,
		encoding.ProvideServerEncoderDecoder(logger, tracerProvider, encoding.ContentTypeJSON),
		cfg,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("serving now")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	// Run server
	go fs.server.Serve()

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
	}()

	cancelCtx, cancelShutdown := context.WithTimeout(ctx, 10*time.Second)
	defer cancelShutdown()

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err = fs.server.Shutdown(cancelCtx); err != nil {
		log.Println("shutting down server", err)
	}
}

///

func header() g.Node {
	return ghtml.Header(
		ghtml.Class("text-center py-6"),
		ghtml.H1(
			ghtml.Class("text-3xl font-bold text-indigo-700"),
			g.Text("My App"),
		),
	)
}

var (
	tailwindImport = ghtml.Script(ghtml.Src("https://cdn.tailwindcss.com?plugins=typography"))

	htmxImport = g.Group{
		ghtml.Script(
			ghtml.Src("https://cdn.jsdelivr.net/npm/htmx.org@2.0.7/dist/htmx.min.js"),
			ghtml.Integrity("sha384-ZBXiYtYQ6hJ2Y0ZNoYuI+Nq5MqWBr+chMrS/RkXpNzQCApHEhOt2aY8EJgqwHLkJ"),
			ghtml.CrossOrigin("anonymous"),
		),

		ghtml.Script(
			ghtml.Src("https://unpkg.com/htmx.org@2.0.7/dist/ext/json-enc.js"),
			ghtml.Integrity("sha384-j+tqxLrwDkbeOdjbpaVekgvQL/J7qm/yh/UqSEs6RjEtnBFHqlJViBWG/jBZ6I2p"),
			ghtml.CrossOrigin("anonymous"),
		),
	}
)

func page(title string, children ...g.Node) g.Node {
	return ghtml.HTML(
		ghtml.Lang("en"),
		ghtml.Head(
			ghtml.Title(title),
			tailwindImport,
			htmxImport,
		),
		ghtml.Body(
			ghtml.Class("bg-gradient-to-b from-white to-indigo-100 min-h-screen flex flex-col"),
			header(),
			ghtml.Main(
				ghtml.Class("flex-grow flex justify-center items-center w-full"),
				ghtml.Div(
					ghtml.Class("w-full max-w-md"),
					g.Group(children),
				),
			),
			footer(),
		),
	)
}

func footer() g.Node {
	return ghtml.Footer(
		ghtml.Class("text-center py-4 text-sm text-gray-600"),
		g.Text("© 2025 My App. All rights reserved."),
	)
}
