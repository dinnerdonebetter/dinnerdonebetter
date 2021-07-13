package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/config"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	"gitlab.com/prixfixe/prixfixe/internal/secrets"
	server "gitlab.com/prixfixe/prixfixe/internal/server"
	auditservice "gitlab.com/prixfixe/prixfixe/internal/services/audit"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	frontendservice "gitlab.com/prixfixe/prixfixe/internal/services/frontend"
	invitationsservice "gitlab.com/prixfixe/prixfixe/internal/services/invitations"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/internal/services/reports"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparationinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/internal/services/webhooks"
	storage "gitlab.com/prixfixe/prixfixe/internal/storage"
	uploads "gitlab.com/prixfixe/prixfixe/internal/uploads"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	defaultPort              = 8888
	defaultCookieDomain      = "localhost"
	debugCookieSecret        = "HEREISA32CHARSECRETWHICHISMADEUP"
	devSqliteConnDetails     = "/tmp/db"
	devMariaDBConnDetails    = "dbuser:hunter2@tcp(database:3306)/prixfixe"
	devPostgresDBConnDetails = "postgres://dbuser:hunter2@database:5432/prixfixe?sslmode=disable"
	defaultCookieName        = authservice.DefaultCookieName

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// database providers.
	postgres = "postgres"

	// test user stuff.
	defaultPassword = "password"

	// search index paths.
	defaultValidInstrumentsSearchIndexPath  = "valid_instruments.bleve"
	defaultValidPreparationsSearchIndexPath = "valid_preparations.bleve"
	defaultValidIngredientsSearchIndexPath  = "valid_ingredients.bleve"

	pasetoSecretSize      = 32
	maxAttempts           = 50
	defaultPASETOLifetime = 1 * time.Minute

	contentTypeJSON = "application/json"
)

var (
	examplePASETOKey = generatePASETOKey()

	noopTracingConfig = tracing.Config{
		Provider:                  "",
		SpanCollectionProbability: 1,
	}

	localServer = server.Config{
		Debug:           true,
		HTTPPort:        defaultPort,
		StartupDeadline: time.Minute,
	}

	localCookies = authservice.CookieConfig{
		Name:       defaultCookieName,
		Domain:     defaultCookieDomain,
		HashKey:    debugCookieSecret,
		SigningKey: debugCookieSecret,
		Lifetime:   authservice.DefaultCookieLifetime,
		SecureOnly: false,
	}

	localTracingConfig = tracing.Config{
		Provider:                  "jaeger",
		SpanCollectionProbability: 1,
		Jaeger: &tracing.JaegerConfig{
			CollectorEndpoint: "http://tracing-server:14268/api/traces",
			ServiceName:       "prixfixe_service",
		},
	}
)

func initializeLocalSecretManager(ctx context.Context) secrets.SecretManager {
	logger := logging.NewNoopLogger()

	cfg := &secrets.Config{
		Provider: secrets.ProviderLocal,
		Key:      "SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU=",
	}

	k, err := secrets.ProvideSecretKeeper(ctx, cfg)
	if err != nil {
		panic(err)
	}

	sm, err := secrets.ProvideSecretManager(logger, k)
	if err != nil {
		panic(err)
	}

	return sm
}

func encryptAndSaveConfig(ctx context.Context, outputPath string, cfg *config.InstanceConfig) error {
	sm := initializeLocalSecretManager(ctx)
	output, err := sm.Encrypt(ctx, cfg)
	if err != nil {
		return fmt.Errorf("encrypting config: %v", err)
	}

	if err = os.MkdirAll(filepath.Dir(outputPath), 0777); err != nil {
		// that's okay
	}

	return os.WriteFile(outputPath, []byte(output), 0644)
}

type configFunc func(ctx context.Context, filePath string) error

var files = map[string]configFunc{
	"environments/local/service.config":                                   localDevelopmentConfig,
	"environments/testing/config_files/frontend-tests.config":             frontendTestsConfig,
	"environments/testing/config_files/integration-tests-postgres.config": buildIntegrationTestForDBImplementation(postgres, devPostgresDBConnDetails),
}

func buildLocalFrontendServiceConfig() frontendservice.Config {
	return frontendservice.Config{
		UseFakeData: false,
	}
}

func mustHashPass(password string) string {
	hashed, err := authentication.ProvideArgon2Authenticator(logging.NewNoopLogger()).
		HashPassword(context.Background(), password)

	if err != nil {
		panic(err)
	}

	return hashed
}

func generatePASETOKey() []byte {
	b := make([]byte, pasetoSecretSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return b
}

func localDevelopmentConfig(ctx context.Context, filePath string) error {
	cfg := &config.InstanceConfig{
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Server: localServer,
		Database: dbconfig.Config{
			Debug:                     true,
			RunMigrations:             true,
			MaxPingAttempts:           maxAttempts,
			Provider:                  postgres,
			ConnectionDetails:         devPostgresDBConnDetails,
			MetricsCollectionInterval: time.Second,
			CreateTestUser: &types.TestUserCreationConfig{
				Username:       "username",
				Password:       defaultPassword,
				HashedPassword: mustHashPass(defaultPassword),
				IsServiceAdmin: true,
			},
		},
		Observability: observability.Config{
			Metrics: metrics.Config{
				Provider:                         "prometheus",
				RouteToken:                       "",
				RuntimeMetricsCollectionInterval: time.Second,
			},
			Tracing: localTracingConfig,
		},
		Uploads: uploads.Config{
			Debug: true,
			Storage: storage.Config{
				UploadFilenameKey: "avatar",
				Provider:          "filesystem",
				BucketName:        "avatars",
				AzureConfig:       nil,
				GCSConfig:         nil,
				S3Config:          nil,
				FilesystemConfig: &storage.FilesystemConfig{
					RootDirectory: "/avatars",
				},
			},
		},
		Search: search.Config{
			Provider: "bleve",
		},
		Services: config.ServicesConfigurations{
			AuditLog: auditservice.Config{
				Debug:   true,
				Enabled: true,
			},
			Auth: authservice.Config{
				PASETO: authservice.PASETOConfig{
					Issuer:       "prixfixe_service",
					Lifetime:     defaultPASETOLifetime,
					LocalModeKey: examplePASETOKey,
				},
				Cookies:               localCookies,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 4,
				MinimumPasswordLength: 8,
			},
			Frontend: buildLocalFrontendServiceConfig(),
			Webhooks: webhooksservice.Config{
				Debug:   true,
				Enabled: false,
			},
			ValidInstruments: validinstrumentsservice.Config{
				SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidInstrumentsSearchIndexPath),
				Logging: logging.Config{
					Name:     "valid_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidPreparations: validpreparationsservice.Config{
				SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidPreparationsSearchIndexPath),
				Logging: logging.Config{
					Name:     "valid_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredients: validingredientsservice.Config{
				SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidIngredientsSearchIndexPath),
				Logging: logging.Config{
					Name:     "valid_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				Logging: logging.Config{
					Name:     "valid_ingredient_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidPreparationInstruments: validpreparationinstrumentsservice.Config{
				Logging: logging.Config{
					Name:     "valid_preparation_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Recipes: recipesservice.Config{
				Logging: logging.Config{
					Name:     "recipes",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeSteps: recipestepsservice.Config{
				Logging: logging.Config{
					Name:     "recipe_steps",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				Logging: logging.Config{
					Name:     "recipe_step_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				Logging: logging.Config{
					Name:     "recipe_step_products",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Invitations: invitationsservice.Config{
				Logging: logging.Config{
					Name:     "invitations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Reports: reportsservice.Config{
				Logging: logging.Config{
					Name:     "reports",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
		},
	}

	return encryptAndSaveConfig(ctx, filePath, cfg)
}

func frontendTestsConfig(ctx context.Context, filePath string) error {
	cfg := &config.InstanceConfig{
		Meta: config.MetaSettings{
			Debug:   false,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Server: localServer,
		Database: dbconfig.Config{
			Debug:                     true,
			RunMigrations:             true,
			Provider:                  postgres,
			ConnectionDetails:         devPostgresDBConnDetails,
			MaxPingAttempts:           maxAttempts,
			MetricsCollectionInterval: time.Second,
		},
		Observability: observability.Config{
			Metrics: metrics.Config{
				Provider:                         "prometheus",
				RouteToken:                       "",
				RuntimeMetricsCollectionInterval: time.Second,
			},
			Tracing: noopTracingConfig,
		},
		Uploads: uploads.Config{
			Debug: true,
			Storage: storage.Config{
				UploadFilenameKey: "avatar",
				Provider:          "memory",
				BucketName:        "avatars",
			},
		},
		Search: search.Config{
			Provider: "bleve",
		},
		Services: config.ServicesConfigurations{
			AuditLog: auditservice.Config{
				Debug:   true,
				Enabled: true,
			},
			Auth: authservice.Config{
				PASETO: authservice.PASETOConfig{
					Issuer:       "prixfixe_service",
					Lifetime:     defaultPASETOLifetime,
					LocalModeKey: examplePASETOKey,
				},
				Cookies:               localCookies,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 4,
				MinimumPasswordLength: 8,
			},
			Frontend: buildLocalFrontendServiceConfig(),
			Webhooks: webhooksservice.Config{
				Debug:   true,
				Enabled: false,
			},
			ValidInstruments: validinstrumentsservice.Config{
				SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidInstrumentsSearchIndexPath),
				Logging: logging.Config{
					Name:     "valid_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidPreparations: validpreparationsservice.Config{
				SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidPreparationsSearchIndexPath),
				Logging: logging.Config{
					Name:     "valid_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredients: validingredientsservice.Config{
				SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidIngredientsSearchIndexPath),
				Logging: logging.Config{
					Name:     "valid_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				Logging: logging.Config{
					Name:     "valid_ingredient_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidPreparationInstruments: validpreparationinstrumentsservice.Config{
				Logging: logging.Config{
					Name:     "valid_preparation_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Recipes: recipesservice.Config{
				Logging: logging.Config{
					Name:     "recipes",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeSteps: recipestepsservice.Config{
				Logging: logging.Config{
					Name:     "recipe_steps",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				Logging: logging.Config{
					Name:     "recipe_step_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				Logging: logging.Config{
					Name:     "recipe_step_products",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Invitations: invitationsservice.Config{
				Logging: logging.Config{
					Name:     "invitations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Reports: reportsservice.Config{
				Logging: logging.Config{
					Name:     "reports",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
		},
	}

	return encryptAndSaveConfig(ctx, filePath, cfg)
}

func buildIntegrationTestForDBImplementation(dbVendor, dbDetails string) configFunc {
	return func(ctx context.Context, filePath string) error {
		startupDeadline := time.Minute

		cfg := &config.InstanceConfig{
			Meta: config.MetaSettings{
				Debug:   false,
				RunMode: testingEnv,
			},
			Encoding: encoding.Config{
				ContentType: contentTypeJSON,
			},
			Server: server.Config{
				Debug:           false,
				HTTPPort:        defaultPort,
				StartupDeadline: startupDeadline,
			},
			Database: dbconfig.Config{
				Debug:                     false,
				RunMigrations:             true,
				Provider:                  dbVendor,
				MaxPingAttempts:           maxAttempts,
				MetricsCollectionInterval: 2 * time.Second,
				ConnectionDetails:         database.ConnectionDetails(dbDetails),
				CreateTestUser: &types.TestUserCreationConfig{
					Username:       "exampleUser",
					Password:       "integration-tests-are-cool",
					HashedPassword: mustHashPass("integration-tests-are-cool"),
					IsServiceAdmin: true,
				},
			},
			Observability: observability.Config{
				Metrics: metrics.Config{
					Provider:                         "",
					RouteToken:                       "",
					RuntimeMetricsCollectionInterval: time.Second,
				},
				Tracing: localTracingConfig,
			},
			Uploads: uploads.Config{
				Debug: false,
				Storage: storage.Config{
					Provider:    "memory",
					BucketName:  "avatars",
					AzureConfig: nil,
					GCSConfig:   nil,
					S3Config:    nil,
				},
			},
			Search: search.Config{
				Provider: "bleve",
			},
			Services: config.ServicesConfigurations{
				AuditLog: auditservice.Config{
					Debug:   false,
					Enabled: true,
				},
				Auth: authservice.Config{
					PASETO: authservice.PASETOConfig{
						Issuer:       "prixfixe_service",
						Lifetime:     defaultPASETOLifetime,
						LocalModeKey: examplePASETOKey,
					},
					Cookies: authservice.CookieConfig{
						Name:       defaultCookieName,
						Domain:     defaultCookieDomain,
						SigningKey: debugCookieSecret,
						Lifetime:   authservice.DefaultCookieLifetime,
						SecureOnly: false,
					},
					Debug:                 false,
					EnableUserSignup:      true,
					MinimumUsernameLength: 4,
					MinimumPasswordLength: 8,
				},
				Frontend: buildLocalFrontendServiceConfig(),
				Webhooks: webhooksservice.Config{
					Debug:   true,
					Enabled: false,
				},
				ValidInstruments: validinstrumentsservice.Config{
					SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidInstrumentsSearchIndexPath),
					Logging: logging.Config{
						Name:     "valid_instruments",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidPreparations: validpreparationsservice.Config{
					SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidPreparationsSearchIndexPath),
					Logging: logging.Config{
						Name:     "valid_preparations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidIngredients: validingredientsservice.Config{
					SearchIndexPath: fmt.Sprintf("/search_indices/%s", defaultValidIngredientsSearchIndexPath),
					Logging: logging.Config{
						Name:     "valid_ingredients",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidIngredientPreparations: validingredientpreparationsservice.Config{
					Logging: logging.Config{
						Name:     "valid_ingredient_preparations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidPreparationInstruments: validpreparationinstrumentsservice.Config{
					Logging: logging.Config{
						Name:     "valid_preparation_instruments",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				Recipes: recipesservice.Config{
					Logging: logging.Config{
						Name:     "recipes",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeSteps: recipestepsservice.Config{
					Logging: logging.Config{
						Name:     "recipe_steps",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeStepIngredients: recipestepingredientsservice.Config{
					Logging: logging.Config{
						Name:     "recipe_step_ingredients",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeStepProducts: recipestepproductsservice.Config{
					Logging: logging.Config{
						Name:     "recipe_step_products",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				Invitations: invitationsservice.Config{
					Logging: logging.Config{
						Name:     "invitations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				Reports: reportsservice.Config{
					Logging: logging.Config{
						Name:     "reports",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
			},
		}

		return encryptAndSaveConfig(ctx, filePath, cfg)
	}
}

func main() {
	ctx := context.Background()

	for filePath, fun := range files {
		if err := fun(ctx, filePath); err != nil {
			log.Fatalf("error rendering %s: %v", filePath, err)
		}
	}
}
