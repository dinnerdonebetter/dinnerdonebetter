package postgres

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"

	"github.com/GuiaBolso/darwin"
)

func fetchMigration(name string) string {
	file, err := rawMigrations.ReadFile(fmt.Sprintf("migrations/%s.sql", name))
	if err != nil {
		panic(err)
	}

	return string(file)
}

// BuildMigrationFunc returns a sync.Once compatible function closure that will migrate a postgres database.
func (q *Querier) migrationFunc() {
	q.logger.Info("migrating database")

	infoChan := make(chan darwin.MigrationInfo)
	go func() {
		for x := range infoChan {
			q.logger.WithValues(map[string]any{
				"version": x.Migration.Version,
				"status":  x.Status.String(),
			}).Info("migrating database")
		}
	}()

	startTime := time.Now()
	if err := darwin.New(darwin.NewGenericDriver(q.db, darwin.PostgresDialect{}), migrations, infoChan).Migrate(); err != nil {
		panic(fmt.Errorf("running migration: %w", err))
	}
	q.logger.WithValue("elapsed", time.Since(startTime).Milliseconds()).Info("migration completed")
}

// Migrate is a simple wrapper around the core querier Migrate call.
func (q *Querier) Migrate(ctx context.Context) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if q.config == nil {
		return ErrNilInputProvided
	}

	if !q.IsReady(ctx) {
		return database.ErrDatabaseNotReady
	}

	q.migrateOnce.Do(q.migrationFunc)

	return nil
}

var (
	//go:embed migrations/*.sql
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
