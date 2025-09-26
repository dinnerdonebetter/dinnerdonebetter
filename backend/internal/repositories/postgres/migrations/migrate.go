package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	pgmigrations "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/migrations"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var (
	//go:embed migration_files/*.sql
	rawMigrations embed.FS
)

func fetchMigration(name string) string {
	file, err := rawMigrations.ReadFile(fmt.Sprintf("migration_files/%s.sql", name))
	if err != nil {
		panic(err)
	}

	return string(file)
}

func NewMigrator(logger logging.Logger, tracerProvider tracing.TracerProvider, db *sql.DB, config *databasecfg.Config) *pgmigrations.Migrator {
	migrations := []*databasecfg.MigrationSpec{
		{
			Description: "basic infrastructural tables",
			RawQuery:    fetchMigration("00001_baseline"),
		},
		{
			Description: "service types and tables",
			RawQuery:    fetchMigration("00002_initial"),
		},
		{
			Description: "user notifications table",
			RawQuery:    fetchMigration("00003_user_notifications"),
		},
		{
			Description: "audit log table",
			RawQuery:    fetchMigration("00004_audit_log"),
		},
		{
			Description: "remove volumetric field",
			RawQuery:    fetchMigration("00005_remove_volumetric_field"),
		},
	}

	return pgmigrations.NewMigrator(logger, tracerProvider, db, config, migrations)
}
