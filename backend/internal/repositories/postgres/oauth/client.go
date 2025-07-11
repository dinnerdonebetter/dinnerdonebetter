package oauth

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption/salsa20"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	generated2 "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth/generated"
)

const (
	o11yName = "oauth_db_client"
)

// repository is the oauth2 client and token repo implemenation.
type repository struct {
	database.Client
	tracer                  tracing.Tracer
	logger                  logging.Logger
	generatedQuerier        generated2.Querier
	auditLogEntryRepo       audit.Repository
	oauth2ClientTokenEncDec encryption.EncryptorDecryptor
	secretGenerator         random.Generator
	db                      *sql.DB
}

// ProvideOAuthRepository provides a new repository.
func ProvideOAuthRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	cfg databasecfg.Config,
	client database.Client,
) oauth.Repository {
	encDec, err := salsa20.NewEncryptorDecryptor(tracerProvider, logger, []byte(cfg.OAuth2TokenEncryptionKey))
	if err != nil {
		return nil
	}

	c := &repository{
		Client:                  client,
		db:                      client.DB(),
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:        generated2.New(),
		auditLogEntryRepo:       auditLogEntryRepo,
		oauth2ClientTokenEncDec: encDec,
		secretGenerator:         random.NewGenerator(logger, tracerProvider),
		logger:                  logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
