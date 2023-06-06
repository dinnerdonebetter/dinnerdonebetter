package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/panicking"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/services/vendorproxy"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
)

const (
	serverNamespace = "dinner_done_better_service"
	loggerName      = "api_server"
)

type (
	// Server is our API http server.
	Server struct {
		authService                            types.AuthService
		householdsService                      types.HouseholdDataService
		householdInvitationsService            types.HouseholdInvitationDataService
		usersService                           types.UserDataService
		adminService                           types.AdminService
		apiClientsService                      types.APIClientDataService
		webhooksService                        types.WebhookDataService
		validInstrumentsService                types.ValidInstrumentDataService
		validIngredientsService                types.ValidIngredientDataService
		validIngredientGroupsService           types.ValidIngredientGroupDataService
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
		validIngredientStatesService           types.ValidIngredientStateDataService
		validPreparationInstrumentsService     types.ValidPreparationInstrumentDataService
		validIngredientMeasurementUnitsService types.ValidIngredientMeasurementUnitDataService
		mealPlanEventsService                  types.MealPlanEventDataService
		mealPlanTasksService                   types.MealPlanTaskDataService
		recipePrepTasksService                 types.RecipePrepTaskDataService
		mealPlanGroceryListItemsService        types.MealPlanGroceryListItemDataService
		validMeasurementConversionsService     types.ValidMeasurementConversionDataService
		recipeStepCompletionConditionsService  types.RecipeStepCompletionConditionDataService
		validIngredientStateIngredientsService types.ValidIngredientStateIngredientDataService
		recipeStepVesselsService               types.RecipeStepVesselDataService
		serviceSettingsService                 types.ServiceSettingDataService
		serviceSettingConfigurationsService    types.ServiceSettingConfigurationDataService
		userIngredientPreferencesService       types.UserIngredientPreferenceDataService
		mealRatingsService                     types.MealRatingDataService
		householdInstrumentOwnershipService    types.HouseholdInstrumentOwnershipDataService
		vendorProxyService                     vendorproxy.Service
		encoder                                encoding.ServerEncoderDecoder
		logger                                 logging.Logger
		router                                 routing.Router
		tracer                                 tracing.Tracer
		panicker                               panicking.Panicker
		httpServer                             *http.Server
		dataManager                            database.DataManager
		tracerProvider                         tracing.TracerProvider
		config                                 Config
	}
)

// ProvideHTTPServer builds a new Server instance.
func ProvideHTTPServer(
	ctx context.Context,
	serverSettings Config,
	dataManager database.DataManager,
	logger logging.Logger,
	encoder encoding.ServerEncoderDecoder,
	router routing.Router,
	tracerProvider tracing.TracerProvider,
	authService types.AuthService,
	usersService types.UserDataService,
	householdsService types.HouseholdDataService,
	householdInvitationsService types.HouseholdInvitationDataService,
	apiClientsService types.APIClientDataService,
	validInstrumentsService types.ValidInstrumentDataService,
	validIngredientsService types.ValidIngredientDataService,
	validIngredientGroupsService types.ValidIngredientGroupDataService,
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
	validIngredientStatesService types.ValidIngredientStateDataService,
	validPreparationInstrumentsService types.ValidPreparationInstrumentDataService,
	validIngredientMeasurementUnitsService types.ValidIngredientMeasurementUnitDataService,
	mealPlanEventsService types.MealPlanEventDataService,
	mealPlanTasksService types.MealPlanTaskDataService,
	recipePrepTasksService types.RecipePrepTaskDataService,
	mealPlanGroceryListItemsService types.MealPlanGroceryListItemDataService,
	validMeasurementConversionsService types.ValidMeasurementConversionDataService,
	recipeStepCompletionConditionsService types.RecipeStepCompletionConditionDataService,
	validIngredientStateIngredientsService types.ValidIngredientStateIngredientDataService,
	recipeStepVesselsService types.RecipeStepVesselDataService,
	webhooksService types.WebhookDataService,
	adminService types.AdminService,
	vendorProxyService vendorproxy.Service,
	serviceSettingDataService types.ServiceSettingDataService,
	serviceSettingConfigurationsService types.ServiceSettingConfigurationDataService,
	userIngredientPreferencesService types.UserIngredientPreferenceDataService,
	mealRatingsService types.MealRatingDataService,
	householdInstrumentOwnershipService types.HouseholdInstrumentOwnershipDataService,
) (*Server, error) {
	srv := &Server{
		config: serverSettings,

		// infra things,
		tracer:         tracing.NewTracer(tracerProvider.Tracer(loggerName)),
		encoder:        encoder,
		logger:         logging.EnsureLogger(logger).WithName(loggerName),
		panicker:       panicking.NewProductionPanicker(),
		httpServer:     ProvideStdLibHTTPServer(serverSettings.HTTPPort),
		dataManager:    dataManager,
		tracerProvider: tracerProvider,

		// services,
		adminService:                           adminService,
		webhooksService:                        webhooksService,
		usersService:                           usersService,
		householdsService:                      householdsService,
		householdInvitationsService:            householdInvitationsService,
		authService:                            authService,
		validInstrumentsService:                validInstrumentsService,
		validIngredientsService:                validIngredientsService,
		validIngredientGroupsService:           validIngredientGroupsService,
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
		validIngredientStatesService:           validIngredientStatesService,
		recipeStepCompletionConditionsService:  recipeStepCompletionConditionsService,
		validIngredientStateIngredientsService: validIngredientStateIngredientsService,
		recipeStepVesselsService:               recipeStepVesselsService,
		vendorProxyService:                     vendorProxyService,
		serviceSettingsService:                 serviceSettingDataService,
		serviceSettingConfigurationsService:    serviceSettingConfigurationsService,
		userIngredientPreferencesService:       userIngredientPreferencesService,
		mealRatingsService:                     mealRatingsService,
		householdInstrumentOwnershipService:    householdInstrumentOwnershipService,
	}

	srv.setupRouter(ctx, router)

	logger.Debug("HTTP server successfully constructed")

	return srv, nil
}

// Shutdown shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.dataManager.Close()

	if err := s.tracerProvider.ForceFlush(ctx); err != nil {
		s.logger.Error(err, "flushing traces")
	}

	return s.httpServer.Shutdown(ctx)
}

// Serve serves HTTP traffic.
func (s *Server) Serve() {
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
			if errors.Is(err, http.ErrServerClosed) {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop.
				os.Exit(0)
			}

			s.logger.Error(err, "shutting server down")
		}
	} else {
		// returns ErrServerClosed on graceful close.
		if err := s.httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop.
				os.Exit(0)
			}

			s.logger.Error(err, "shutting server down")
		}
	}
}

const (
	maxTimeout   = 120 * time.Second
	readTimeout  = 5 * time.Second
	writeTimeout = 2 * readTimeout
	idleTimeout  = maxTimeout
)

// ProvideStdLibHTTPServer provides an HTTP httpServer.
func ProvideStdLibHTTPServer(port uint16) *http.Server {
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
