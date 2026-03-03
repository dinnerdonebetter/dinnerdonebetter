package notifications

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption"
	"github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption/salsa20"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications/generated"
)

const (
	o11yName = "notifications_db_client"
)

// Repository is the notifications repository implementation.
// Exported so the manager can wrap it; the manager is the sole provider of notifications.Repository for services.
type Repository struct {
	database.Client
	tracer                tracing.Tracer
	logger                logging.Logger
	generatedQuerier      generated.Querier
	auditLogEntryRepo     audit.Repository
	userDeviceTokenEncDec encryption.EncryptorDecryptor
	readDB                *sql.DB
	writeDB               *sql.DB
}

// ProvideNotificationsRepository provides a new repository.
func ProvideNotificationsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	cfg *databasecfg.Config,
	client database.Client,
) *Repository {
	encDec, err := salsa20.NewEncryptorDecryptor(tracerProvider, logger, []byte(cfg.UserDeviceTokenEncryptionKey))
	if err != nil {
		return nil
	}

	c := &Repository{
		Client:                client,
		readDB:                client.ReadDB(),
		writeDB:               client.WriteDB(),
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:      generated.New(),
		auditLogEntryRepo:     auditLogEntryRepo,
		userDeviceTokenEncDec: encDec,
		logger:                logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
