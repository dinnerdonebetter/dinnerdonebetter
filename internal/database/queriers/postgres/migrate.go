package postgres

import (
	"context"
	"embed"
	"fmt"

	"github.com/prixfixeco/api_server/internal/database"

	"github.com/GuiaBolso/darwin"
)

var (
	//go:embed migrations/*.sql
	rawMigrations embed.FS

	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "basic infrastructural tables",
			Script:      fetchMigration("00001_initial"),
		},
		{
			Version:     2,
			Description: "create valid instruments table",
			Script:      fetchMigration("00002_valid_instruments"),
		},
		{
			Version:     3,
			Description: "create valid ingredients table",
			Script:      fetchMigration("00003_valid_ingredients"),
		},
		{
			Version:     4,
			Description: "create valid preparations table",
			Script:      fetchMigration("00004_valid_preparations"),
		},
		{
			Version:     5,
			Description: "create valid measurement units bridge table",
			Script:      fetchMigration("00015_valid_measurement_units"),
		},
		{
			Version:     6,
			Description: "create valid ingredient preparations table",
			Script:      fetchMigration("00005_valid_ingredient_preparations"),
		},
		{
			Version:     7,
			Description: "create recipes table",
			Script:      fetchMigration("00006_recipes"),
		},
		{
			Version:     8,
			Description: "create recipe steps table",
			Script:      fetchMigration("00007_recipe_steps"),
		},
		{
			Version:     9,
			Description: "create recipe step instruments table",
			Script:      fetchMigration("00008_recipe_step_instruments"),
		},
		{
			Version:     10,
			Description: "create recipe step products table",
			Script:      fetchMigration("00010_recipe_step_products"),
		},
		{
			Version:     11,
			Description: "create recipe step ingredients table",
			Script:      fetchMigration("00009_recipe_step_ingredients"),
		},
		{
			Version:     12,
			Description: "create meals table",
			Script:      fetchMigration("00011_meals"),
		},
		{
			Version:     13,
			Description: "create meal plans table",
			Script:      fetchMigration("00012_meal_plans"),
		},
		{
			Version:     14,
			Description: "create meal plan options table",
			Script:      fetchMigration("00013_meal_plan_options"),
		},
		{
			Version:     15,
			Description: "create meal plan option votes table",
			Script:      fetchMigration("00014_meal_plan_option_votes"),
		},
		{
			Version:     16,
			Description: "reintroduce valid preparation instruments bridge table",
			Script:      fetchMigration("00016_valid_preparation_instruments"),
		},
		{
			Version:     17,
			Description: "create valid ingredient measurement units table",
			Script:      fetchMigration("00017_valid_ingredient_measurement_units"),
		},
		{
			Version:     30,
			Description: "myriad recipe step improvements",
			Script:      fetchMigration("00030_recipe_step_improvements"),
		},
		{
			Version:     31,
			Description: "add preference rank to recipe step instruments",
			Script:      fetchMigration("00031_recipe_step_instruments_preference_rank"),
		},
	}
)

func fetchMigration(name string) string {
	file, err := rawMigrations.ReadFile(fmt.Sprintf("migrations/%s.sql", name))
	if err != nil {
		panic(err)
	}

	return string(file)
}

// BuildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a postgres database.
func (q *SQLQuerier) migrationFunc() {
	driver := darwin.NewGenericDriver(q.db, darwin.PostgresDialect{})
	if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
		panic(fmt.Errorf("migrating database: %w", err))
	}
}

// Migrate is a simple wrapper around the core querier Migrate call.
func (q *SQLQuerier) Migrate(ctx context.Context, maxAttempts uint8) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	q.logger.Info("migrating db")

	if !q.IsReady(ctx, maxAttempts) {
		return database.ErrDatabaseNotReady
	}

	q.migrateOnce.Do(q.migrationFunc)

	return nil
}
