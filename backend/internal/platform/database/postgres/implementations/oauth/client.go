package oauth

import (
	"context"
	"database/sql"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption/salsa20"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/implementations/oauth/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
)

// Querier is the oauth2 client and token repo implemenation.
type Querier struct {
	tracer                  tracing.Tracer
	logger                  logging.Logger
	generatedQuerier        generated.Querier
	auditLogEntryRepo       auditlogentries.Repository
	oauth2ClientTokenEncDec encryption.EncryptorDecryptor
	secretGenerator         random.Generator
	timeFunc                func() time.Time
	db                      *sql.DB
}

// ProvideOAuthRepository provides a new client.
func ProvideOAuthRepository(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo auditlogentries.Repository,
	cfg databasecfg.Config,
	db *sql.DB,
) (oauth.Repository, error) {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("auth_db_client"))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	encDec, err := salsa20.NewEncryptorDecryptor(tracerProvider, logger, []byte(cfg.OAuth2TokenEncryptionKey))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating encryptor/decryptor with secret length %d", len(cfg.OAuth2TokenEncryptionKey))
	}

	c := &Querier{
		db:                      db,
		tracer:                  tracer,
		timeFunc:                defaultTimeFunc,
		generatedQuerier:        generated.New(),
		auditLogEntryRepo:       auditLogEntryRepo,
		oauth2ClientTokenEncDec: encDec,
		secretGenerator:         random.NewGenerator(logger, tracerProvider),
		logger:                  logging.EnsureLogger(logger).WithName("querier"),
	}

	return c, nil
}

func defaultTimeFunc() time.Time {
	return time.Now()
}

func (q *Querier) currentTime() time.Time {
	if q == nil || q.timeFunc == nil {
		return defaultTimeFunc()
	}

	return q.timeFunc()
}

func (q *Querier) rollbackTransaction(ctx context.Context, tx database.SQLQueryExecutorAndTransactionManager) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	q.logger.Debug("rolling back transaction")

	if err := tx.Rollback(); err != nil {
		observability.AcknowledgeError(err, q.logger, span, "rolling back transaction")
	}

	q.logger.Debug("transaction rolled back")
}
