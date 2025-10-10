package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
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
	apiConfigurationFilepath = "deploy/environments/localdev/config_files/admin_webapp_config.json"

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

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()

	os.Setenv("DINNER_DONE_BETTER_COOKIES_COOKIE_NAME", "dev_admin_frontend")
	os.Setenv("DINNER_DONE_BETTER_COOKIES_HASH_KEY", base64.RawURLEncoding.EncodeToString([]byte("HEREISA32CHARSECRETWHICHISMADEUP")))
	os.Setenv("DINNER_DONE_BETTER_COOKIES_BLOCK_KEY", base64.RawURLEncoding.EncodeToString([]byte("HEREISA32CHARSECRETWHICHISMADEUP")))
	os.Setenv("DINNER_DONE_BETTER_API_SERVICE_HTTP_API_SERVER_URL", "http://localhost:8000")
	os.Setenv("DINNER_DONE_BETTER_API_SERVICE_GRPC_API_SERVER_URL", ":8001")
	os.Setenv("DINNER_DONE_BETTER_API_SERVICE_OAUTH2_API_CLIENT_ID", strings.Repeat("A", oauth.ClientIDSize))
	os.Setenv("DINNER_DONE_BETTER_API_SERVICE_OAUTH2_API_CLIENT_SECRET", strings.Repeat("A", oauth.ClientSecretSize))
	os.Setenv("DINNER_DONE_BETTER_SERVER_PORT", "8888")
	os.Setenv("DINNER_DONE_BETTER_SERVER_STARTUP_DEADLINE", time.Minute.String())

	cfg, err := config.LoadConfigFromPath[config.AdminWebappConfig](ctx, apiConfigurationFilepath)
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
		apiClient,
		nil,
		nil,
		encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		cfg,
		mux,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("serving now")
	if err = http.ListenAndServe(":8888", fs); err != nil {
		log.Fatal(err)
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
