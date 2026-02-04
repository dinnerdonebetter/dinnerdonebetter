package dataprivacy

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "webhook_db_client"
)

// repository is the webhook repository client.
type repository struct {
	database.Client
	tracer       tracing.Tracer
	logger       logging.Logger
	identityRepo identity.Repository
	readDB       *sql.DB
	writeDB      *sql.DB
}

// ProvideDataPrivacyRepository provides a new repository.
func ProvideDataPrivacyRepository(logger logging.Logger, tracerProvider tracing.TracerProvider, identityRepo identity.Repository, client database.Client) dataprivacy.Repository {
	c := &repository{
		Client:       client,
		readDB:       client.ReadDB(),
		writeDB:      client.WriteDB(),
		tracer:       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepo: identityRepo,
		logger:       logging.EnsureLogger(logger).WithName(o11yName),
	}

	// these are here for future use
	_, _ = c.readDB, c.writeDB

	return c
}
