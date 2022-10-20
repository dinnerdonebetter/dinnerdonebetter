package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/panicking"
	"github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serverNamespace = "prixfixe_service"
	loggerName      = "api_server"
)

type (
	// HTTPServer is our API http server.
	HTTPServer struct {
		authService                            types.AuthService
		householdsService                      types.HouseholdDataService
		householdInvitationsService            types.HouseholdInvitationDataService
		usersService                           types.UserDataService
		adminService                           types.AdminService
		apiClientsService                      types.APIClientDataService
		webhooksService                        types.WebhookDataService
		validInstrumentsService                types.ValidInstrumentDataService
		validIngredientsService                types.ValidIngredientDataService
		validPreparationsService               types.ValidPreparationDataService
		validIngredientPreparationsService     types.ValidIngredientPreparationDataService
		recipesService                         types.RecipeDataService
		recipeStepsService                     types.RecipeStepDataService
		recipeStepProductsService              types.RecipeStepProductDataService
		recipeStepInstrumentsService           types.RecipeStepInstrumentDataService
		recipeStepIngredientsService           types.RecipeStepIngredientDataService
		mealsService                           types.MealDataService
		mealPlansService                       types.MealPlanDataService
		mealPlanOptionsService                 types.MealPlanOptionDataService
		mealPlanOptionVotesService             types.MealPlanOptionVoteDataService
		validMeasurementUnitsService           types.ValidMeasurementUnitDataService
		validPreparationInstrumentsService     types.ValidPreparationInstrumentDataService
		validIngredientMeasurementUnitsService types.ValidIngredientMeasurementUnitDataService
		mealPlanEventsService                  types.MealPlanEventDataService
		mealPlanTasksService                   types.MealPlanTaskDataService
		recipePrepTasksService                 types.RecipePrepTaskDataService
		mealPlanGroceryListItemsService        types.MealPlanGroceryListItemDataService
		validMeasurementConversionsService     types.ValidMeasurementConversionDataService
		encoder                                encoding.ServerEncoderDecoder
		logger                                 logging.Logger
		router                                 routing.Router
		tracer                                 tracing.Tracer
		panicker                               panicking.Panicker
		httpServer                             *http.Server
		config                                 Config
	}
)

// ProvideHTTPServer builds a new HTTPServer instance.
func ProvideHTTPServer(
	ctx context.Context,
	serverSettings Config,
	authService types.AuthService,
	usersService types.UserDataService,
	householdsService types.HouseholdDataService,
	householdInvitationsService types.HouseholdInvitationDataService,
	apiClientsService types.APIClientDataService,
	validInstrumentsService types.ValidInstrumentDataService,
	validIngredientsService types.ValidIngredientDataService,
	validPreparationsService types.ValidPreparationDataService,
	validIngredientPreparationsService types.ValidIngredientPreparationDataService,
	mealsService types.MealDataService,
	recipesService types.RecipeDataService,
	recipeStepsService types.RecipeStepDataService,
	recipeStepProductsService types.RecipeStepProductDataService,
	recipeStepInstrumentsService types.RecipeStepInstrumentDataService,
	recipeStepIngredientsService types.RecipeStepIngredientDataService,
	mealPlansService types.MealPlanDataService,
	mealPlanOptionsService types.MealPlanOptionDataService,
	mealPlanOptionVotesService types.MealPlanOptionVoteDataService,
	validMeasurementUnitsService types.ValidMeasurementUnitDataService,
	validPreparationInstrumentsService types.ValidPreparationInstrumentDataService,
	validIngredientMeasurementUnitsService types.ValidIngredientMeasurementUnitDataService,
	mealPlanEventsService types.MealPlanEventDataService,
	mealPlanTasksService types.MealPlanTaskDataService,
	recipePrepTasksService types.RecipePrepTaskDataService,
	mealPlanGroceryListItemsService types.MealPlanGroceryListItemDataService,
	validMeasurementConversionsService types.ValidMeasurementConversionDataService,
	webhooksService types.WebhookDataService,
	adminService types.AdminService,
	logger logging.Logger,
	encoder encoding.ServerEncoderDecoder,
	router routing.Router,
	tracerProvider tracing.TracerProvider,
	metricsHandler metrics.Handler,
) (*HTTPServer, error) {
	srv := &HTTPServer{
		config: serverSettings,

		// infra things,
		tracer:     tracing.NewTracer(tracerProvider.Tracer(loggerName)),
		encoder:    encoder,
		logger:     logging.EnsureLogger(logger).WithName(loggerName),
		panicker:   panicking.NewProductionPanicker(),
		httpServer: provideHTTPServer(serverSettings.HTTPPort),

		// services,
		adminService:                           adminService,
		webhooksService:                        webhooksService,
		usersService:                           usersService,
		householdsService:                      householdsService,
		householdInvitationsService:            householdInvitationsService,
		authService:                            authService,
		validInstrumentsService:                validInstrumentsService,
		validIngredientsService:                validIngredientsService,
		validPreparationsService:               validPreparationsService,
		validIngredientPreparationsService:     validIngredientPreparationsService,
		mealsService:                           mealsService,
		recipesService:                         recipesService,
		recipeStepsService:                     recipeStepsService,
		recipeStepProductsService:              recipeStepProductsService,
		recipeStepInstrumentsService:           recipeStepInstrumentsService,
		recipeStepIngredientsService:           recipeStepIngredientsService,
		mealPlansService:                       mealPlansService,
		mealPlanOptionsService:                 mealPlanOptionsService,
		mealPlanOptionVotesService:             mealPlanOptionVotesService,
		validMeasurementUnitsService:           validMeasurementUnitsService,
		apiClientsService:                      apiClientsService,
		validPreparationInstrumentsService:     validPreparationInstrumentsService,
		validIngredientMeasurementUnitsService: validIngredientMeasurementUnitsService,
		mealPlanEventsService:                  mealPlanEventsService,
		mealPlanTasksService:                   mealPlanTasksService,
		recipePrepTasksService:                 recipePrepTasksService,
		mealPlanGroceryListItemsService:        mealPlanGroceryListItemsService,
		validMeasurementConversionsService:     validMeasurementConversionsService,
	}

	srv.setupRouter(ctx, router, metricsHandler)

	logger.Debug("HTTP server successfully constructed")

	return srv, nil
}

// Serve serves HTTP traffic.
func (s *HTTPServer) Serve() {
	s.logger.Debug("setting up server")

	s.httpServer.Handler = otelhttp.NewHandler(
		s.router.Handler(),
		serverNamespace,
		otelhttp.WithSpanNameFormatter(tracing.FormatSpan),
	)

	http2ServerConf := &http2.Server{}
	if err := http2.ConfigureServer(s.httpServer, http2ServerConf); err != nil {
		s.logger.Error(err, "configuring HTTP2")
		s.panicker.Panic(err)
	}

	s.logger.WithValue("listening_on", s.httpServer.Addr).Debug("Listening for HTTP requests")

	if s.config.HTTPSCertificateFile != "" && s.config.HTTPSCertificateKeyFile != "" {
		// returns ErrServerClosed on graceful close.
		if err := s.httpServer.ListenAndServeTLS(s.config.HTTPSCertificateFile, s.config.HTTPSCertificateKeyFile); err != nil {
			s.logger.Error(err, "server shutting down")

			if errors.Is(err, http.ErrServerClosed) {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop.
				os.Exit(0)
			}
		}
	} else {
		// returns ErrServerClosed on graceful close.
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.logger.Error(err, "server shutting down")

			if errors.Is(err, http.ErrServerClosed) {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop.
				os.Exit(0)
			}
		}
	}
}

const (
	maxTimeout   = 120 * time.Second
	readTimeout  = 5 * time.Second
	writeTimeout = 2 * readTimeout
	idleTimeout  = maxTimeout
)

// provideHTTPServer provides an HTTP httpServer.
func provideHTTPServer(port uint16) *http.Server {
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
