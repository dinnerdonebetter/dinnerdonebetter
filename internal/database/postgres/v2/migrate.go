package v2

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

// IsReady is a simple wrapper around the core querier IsReady call.
func (c *DatabaseClient) IsReady(ctx context.Context, waitPeriod time.Duration, maxAttempts uint64) (ready bool) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	attemptCount := 0

	for !ready {
		err := c.db.PingContext(ctx)
		if err != nil {
			c.logger.WithValue("attempt_count", attemptCount).Debug("ping failed, waiting for db")
			time.Sleep(waitPeriod)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		} else {
			ready = true
			return ready
		}
	}

	return false
}

// BuildMigrationFunc returns a sync.Once compatible function closure that will migrate a postgres database.
func (c *DatabaseClient) migrationFunc() {
	driver := darwin.NewGenericDriver(c.db, darwin.PostgresDialect{})
	if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
		panic(fmt.Errorf("running migration: %w", err))
	}
}

// Migrate is a simple wrapper around the core querier Migrate call.
func (c *DatabaseClient) Migrate(ctx context.Context, waitPeriod time.Duration, maxAttempts uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	c.logger.Info("migrating db")

	if !c.IsReady(ctx, waitPeriod, maxAttempts) {
		return database.ErrDatabaseNotReady
	}

	c.migrateOnce.Do(c.migrationFunc)

	return nil
}

var (
	//go:embed migrations/*.sql
	rawMigrations embed.FS

	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "basic infrastructural tables",
			Script:      fetchMigration("00000_initial"),
		},
	}
)
