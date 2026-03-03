package oauth

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption"
	encryptioncfg "github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth/generated"
)

const (
	o11yName = "oauth_db_client"
)

// repository is the oauth2 client and token repo implemenation.
type repository struct {
	database.Client
	tracer                  tracing.Tracer
	logger                  logging.Logger
	generatedQuerier        generated.Querier
	auditLogEntryRepo       audit.Repository
	oauth2ClientTokenEncDec encryption.EncryptorDecryptor
	readDB                  *sql.DB
	writeDB                 *sql.DB
}

// ProvideOAuthRepository provides a new repository.
func ProvideOAuthRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	cfg *databasecfg.Config,
	client database.Client,
) oauth.Repository {
	encDec, err := encryptioncfg.ProvideEncryptorDecryptor(&cfg.Encryption, tracerProvider, logger, []byte(cfg.OAuth2TokenEncryptionKey))
	if err != nil {
		return nil
	}

	c := &repository{
		Client:                  client,
		readDB:                  client.ReadDB(),
		writeDB:                 client.WriteDB(),
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:        generated.New(),
		auditLogEntryRepo:       auditLogEntryRepo,
		oauth2ClientTokenEncDec: encDec,
		logger:                  logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
