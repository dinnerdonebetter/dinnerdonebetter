package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sync"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/GuiaBolso/darwin"
)

type Migrator struct {
	logger      logging.Logger
	tracer      tracing.Tracer
	config      *databasecfg.Config
	db          *sql.DB
	migrateOnce sync.Once
}

func NewMigrator(logger logging.Logger, tracerProvider tracing.TracerProvider, db *sql.DB, config *databasecfg.Config) *Migrator {
	return &Migrator{
		config: config,
		logger: logging.EnsureLogger(logger),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("postgres_migrator")),
		db:     db,
	}
}

// BuildMigrationFunc returns a sync.Once compatible function closure that will migrate a postgres database.
func (m *Migrator) migrationFunc() {
	infoChan := make(chan darwin.MigrationInfo)
	go func() {
		for x := range infoChan {
			m.logger.WithValues(map[string]any{
				"version": x.Migration.Version,
				"status":  x.Status.String(),
			}).Info("database migration handled")
		}
	}()

	startTime := time.Now()
	if err := darwin.New(darwin.NewGenericDriver(m.db, darwin.PostgresDialect{}), migrations, infoChan).Migrate(); err != nil {
		panic(fmt.Errorf("running migration: %w", err))
	}
	m.logger.WithValue("elapsed", time.Since(startTime).Milliseconds()).Info("migration completed")
}

// Migrate is a simple wrapper around the core querier Migrate call.
func (m *Migrator) Migrate(ctx context.Context) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if m.config == nil {
		return database.ErrNilInputProvided
	}

	if !m.IsReady(ctx) {
		return database.ErrDatabaseNotReady
	}

	m.migrateOnce.Do(m.migrationFunc)

	return nil
}

// IsReady returns whether the database is ready for the querier.
func (m *Migrator) IsReady(ctx context.Context) bool {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithValue("connection_url", m.config.ConnectionDetails.String())

	attemptCount := 0
	for {
		if err := m.db.PingContext(ctx); err != nil {
			logger.WithValue("attempt_count", attemptCount).Info("ping failed, waiting for db")
			time.Sleep(m.config.PingWaitPeriod)

			attemptCount++
			if attemptCount >= int(m.config.MaxPingAttempts) {
				break
			}
		} else {
			return true
		}
	}

	return false
}

func fetchMigration(name string) string {
	file, err := rawMigrations.ReadFile(fmt.Sprintf("migration_files/%s.sql", name))
	if err != nil {
		panic(err)
	}

	return string(file)
}

var (
	//go:embed migration_files/*.sql
	rawMigrations embed.FS

	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "basic infrastructural tables",
			Script:      fetchMigration("00001_baseline"),
		},
		{
			Version:     2,
			Description: "service types and tables",
			Script:      fetchMigration("00002_initial"),
		},
		{
			Version:     3,
			Description: "user notifications table",
			Script:      fetchMigration("00003_user_notifications"),
		},
		{
			Version:     4,
			Description: "audit log table",
			Script:      fetchMigration("00004_audit_log"),
		},
		{
			Version:     5,
			Description: "remove volumetric field",
			Script:      fetchMigration("00005_remove_volumetric_field"),
		},
	}
)
