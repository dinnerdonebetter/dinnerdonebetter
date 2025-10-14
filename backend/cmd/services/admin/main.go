package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/config/envvars"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	adminServerConfigurationFilepath = "deploy/environments/localdev/config_files/admin_webapp_config.json"

	o11yName = "admin_frontend"
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
	must(os.Setenv(envvars.CookiesLifetimeEnvVarKey, (7 * 24 * time.Hour).String()))

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

	fs, err := NewAdminFrontendServer(
		ctx,
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

// LayoutConfig holds configuration for the admin layout
type LayoutConfig struct {
	Palette     *design.Palette
	AppName     string
	MaxWidth    string // e.g., "7xl", "6xl", "full"
	Margin      string // e.g., "4", "6", "8"
	ShowSidebar bool
}

// DefaultLayoutConfig provides sensible defaults for the admin layout
func DefaultLayoutConfig() *LayoutConfig {
	return &LayoutConfig{
		Palette:     &design.StandardPalette,
		AppName:     "Admin Dashboard",
		MaxWidth:    "7xl",
		Margin:      "4",
		ShowSidebar: true,
	}
}

func header(config *LayoutConfig) g.Node {
	if config == nil {
		config = DefaultLayoutConfig()
	}

	return ghtml.Header(
		ghtml.Class(fmt.Sprintf("sticky top-0 z-50 %s border-b %s shadow-sm",
			design.Background(config.Palette.Background),
			design.BorderColor(config.Palette.Text),
		)),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("max-w-%s mx-auto px-%s", config.MaxWidth, config.Margin)),
			ghtml.Div(
				ghtml.Class("flex items-center justify-between h-16"),
				// Left side - Logo and main nav
				ghtml.Div(
					ghtml.Class("flex items-center space-x-8"),
					ghtml.H1(
						ghtml.Class(fmt.Sprintf("text-xl font-bold %s", design.TextColor(config.Palette.Primary))),
						g.Text(config.AppName),
					),
					// Main navigation
					ghtml.Nav(
						ghtml.Class("hidden md:flex space-x-6"),
						navLink("Dashboard", "/", config.Palette),
						navLink("Users", "/users", config.Palette),
						navLink("Settings", "/settings", config.Palette),
					),
				),
				// Right side - User menu and mobile menu button
				ghtml.Div(
					ghtml.Class("flex items-center space-x-4"),
					// Mobile menu button
					ghtml.Button(
						ghtml.Class(fmt.Sprintf("md:hidden p-2 rounded-md %s hover:%s focus:outline-none focus:ring-2 focus:ring-%s",
							design.TextColor(config.Palette.Text),
							design.Background(config.Palette.Secondary),
							config.Palette.Primary.Value,
						)),
						ghtml.Type("button"),
						g.Attr("aria-label", "Open menu"),
						g.Attr("onclick", "toggleMobileMenu()"),
						// Hamburger icon
						ghtml.Div(
							ghtml.Class("w-6 h-6 flex flex-col justify-center items-center"),
							ghtml.Span(ghtml.Class(fmt.Sprintf("block w-5 h-0.5 %s mb-1", design.Background(config.Palette.Text)))),
							ghtml.Span(ghtml.Class(fmt.Sprintf("block w-5 h-0.5 %s mb-1", design.Background(config.Palette.Text)))),
							ghtml.Span(ghtml.Class(fmt.Sprintf("block w-5 h-0.5 %s", design.Background(config.Palette.Text)))),
						),
					),
					// User menu
					ghtml.Div(
						ghtml.Class("relative"),
						ghtml.Button(
							ghtml.Class(fmt.Sprintf("flex items-center space-x-2 px-3 py-2 rounded-md %s hover:%s focus:outline-none focus:ring-2 focus:ring-%s",
								design.TextColor(config.Palette.Text),
								design.Background(config.Palette.Secondary),
								config.Palette.Primary.Value,
							)),
							ghtml.Type("button"),
							g.Text("Admin"),
						),
					),
				),
			),
			// Mobile navigation menu
			ghtml.Div(
				ghtml.ID("mobile-menu"),
				ghtml.Class("hidden md:hidden border-t border-gray-200 py-4"),
				ghtml.Nav(
					ghtml.Class("flex flex-col space-y-2"),
					mobileNavLink("Dashboard", "/", config.Palette),
					mobileNavLink("Users", "/users", config.Palette),
					mobileNavLink("Settings", "/settings", config.Palette),
				),
			),
		),
	)
}

func navLink(text, href string, palette *design.Palette) g.Node {
	return ghtml.A(
		ghtml.Href(href),
		ghtml.Class(fmt.Sprintf("px-3 py-2 rounded-md text-sm font-medium %s hover:%s hover:%s transition-colors duration-200",
			design.TextColor(palette.Text),
			design.TextColor(palette.Primary),
			design.Background(palette.Secondary),
		)),
		g.Text(text),
	)
}

func mobileNavLink(text, href string, palette *design.Palette) g.Node {
	return ghtml.A(
		ghtml.Href(href),
		ghtml.Class(fmt.Sprintf("block px-3 py-2 rounded-md text-base font-medium %s hover:%s hover:%s transition-colors duration-200",
			design.TextColor(palette.Text),
			design.TextColor(palette.Primary),
			design.Background(palette.Secondary),
		)),
		g.Text(text),
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

	// JavaScript for mobile menu toggle
	mobileMenuScript = ghtml.Script(
		g.Raw(`
			function toggleMobileMenu() {
				const menu = document.getElementById('mobile-menu');
				menu.classList.toggle('hidden');
			}
			
			// Close mobile menu when clicking outside
			document.addEventListener('click', function(event) {
				const menu = document.getElementById('mobile-menu');
				const button = event.target.closest('button[aria-label="Open menu"]');
				
				if (!button && !menu.contains(event.target)) {
					menu.classList.add('hidden');
				}
			});
		`),
	)
)

func page(title string, children ...g.Node) g.Node {
	return pageWithConfig(title, nil, children...)
}

func pageWithConfig(title string, config *LayoutConfig, children ...g.Node) g.Node {
	if config == nil {
		config = DefaultLayoutConfig()
	}

	return ghtml.HTML(
		ghtml.Lang("en"),
		ghtml.Head(
			ghtml.Meta(ghtml.Charset("utf-8")),
			ghtml.Meta(ghtml.Name("viewport"), ghtml.Content("width=device-width, initial-scale=1")),
			ghtml.Title(fmt.Sprintf("%s - %s", title, config.AppName)),
			tailwindImport,
			htmxImport,
			mobileMenuScript,
		),
		ghtml.Body(
			ghtml.Class(fmt.Sprintf("min-h-screen flex flex-col %s %s",
				design.Background(config.Palette.Background),
				design.TextColor(config.Palette.Text),
			)),
			header(config),
			ghtml.Main(
				ghtml.Class("flex-1 overflow-hidden"),
				ghtml.Div(
					ghtml.Class(fmt.Sprintf("h-full max-w-%s mx-auto px-%s py-%s",
						config.MaxWidth, config.Margin, config.Margin)),
					ghtml.Div(
						ghtml.Class("h-full overflow-auto"),
						ghtml.Div(
							ghtml.ID("main-content"),
							ghtml.Class("min-h-full"),
							g.Group(children),
						),
					),
				),
			),
			footer(config),
		),
	)
}

func footer(config *LayoutConfig) g.Node {
	if config == nil {
		config = DefaultLayoutConfig()
	}

	return ghtml.Footer(
		ghtml.Class(fmt.Sprintf("border-t %s %s",
			design.BorderColor(config.Palette.Text),
			design.Background(config.Palette.Background),
		)),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("max-w-%s mx-auto px-%s py-4", config.MaxWidth, config.Margin)),
			ghtml.Div(
				ghtml.Class("flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0"),
				ghtml.P(
					ghtml.Class(fmt.Sprintf("text-sm %s", design.TextColor(config.Palette.Text))),
					g.Textf("© 2025 %s. All rights reserved.", config.AppName),
				),
				ghtml.Div(
					ghtml.Class("flex space-x-4 text-sm"),
					ghtml.A(
						ghtml.Href("/privacy"),
						ghtml.Class(fmt.Sprintf("%s hover:%s transition-colors duration-200",
							design.TextColor(config.Palette.Text),
							design.TextColor(config.Palette.Primary),
						)),
						g.Text("Privacy Policy"),
					),
					ghtml.A(
						ghtml.Href("/terms"),
						ghtml.Class(fmt.Sprintf("%s hover:%s transition-colors duration-200",
							design.TextColor(config.Palette.Text),
							design.TextColor(config.Palette.Primary),
						)),
						g.Text("Terms of Service"),
					),
				),
			),
		),
	)
}
