package oauth

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption/salsa20"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/oauth/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
)

const (
	o11yName = "oauth_db_client"
)

// Querier is the oauth2 client and token repo implemenation.
type Querier struct {
	database.Client
	tracer                  tracing.Tracer
	logger                  logging.Logger
	generatedQuerier        generated.Querier
	auditLogEntryRepo       auditlogentries.Repository
	oauth2ClientTokenEncDec encryption.EncryptorDecryptor
	secretGenerator         random.Generator
	db                      *sql.DB
}

// ProvideOAuthRepository provides a new client.
func ProvideOAuthRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo auditlogentries.Repository,
	cfg databasecfg.Config,
	client database.Client,
) oauth.Repository {
	encDec, err := salsa20.NewEncryptorDecryptor(tracerProvider, logger, []byte(cfg.OAuth2TokenEncryptionKey))
	if err != nil {
		return nil
	}

	c := &Querier{
		Client:                  client,
		db:                      client.DB(),
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:        generated.New(),
		auditLogEntryRepo:       auditLogEntryRepo,
		oauth2ClientTokenEncDec: encDec,
		secretGenerator:         random.NewGenerator(logger, tracerProvider),
		logger:                  logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
