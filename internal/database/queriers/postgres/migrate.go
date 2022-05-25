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
			Script:      fetchMigration("migrations/00001_initial.sql"),
		},
		{
			Version:     2,
			Description: "create valid instruments table",
			Script:      fetchMigration("migrations/00002_valid_instruments.sql"),
		},
		{
			Version:     3,
			Description: "create valid ingredients table",
			Script:      fetchMigration("migrations/00003_valid_ingredients.sql"),
		},
		{
			Version:     4,
			Description: "create valid preparations table",
			Script:      fetchMigration("migrations/00004_valid_preparations.sql"),
		},
		{
			Version:     5,
			Description: "create valid ingredient preparations table",
			Script:      fetchMigration("migrations/00005_valid_ingredient_preparations.sql"),
		},
		{
			Version:     6,
			Description: "create recipes table",
			Script:      fetchMigration("migrations/00006_recipes.sql"),
		},
		{
			Version:     7,
			Description: "create recipe steps table",
			Script:      fetchMigration("migrations/00007_recipe_steps.sql"),
		},
		{
			Version:     8,
			Description: "create recipe step instruments table",
			Script:      fetchMigration("migrations/00008_recipe_step_instruments.sql"),
		},
		{
			Version:     9,
			Description: "create recipe step ingredients table",
			Script:      fetchMigration("migrations/00009_recipe_step_ingredients.sql"),
		},
		{
			Version:     10,
			Description: "create recipe step products table",
			Script:      fetchMigration("migrations/00010_recipe_step_products.sql"),
		},
		{
			Version:     11,
			Description: "create meals table",
			Script:      fetchMigration("migrations/00011_meals.sql"),
		},
		{
			Version:     12,
			Description: "create meal plans table",
			Script:      fetchMigration("migrations/00012_meal_plans.sql"),
		},
		{
			Version:     13,
			Description: "create meal plan options table",
			Script:      fetchMigration("migrations/00013_meal_plan_options.sql"),
		},
		{
			Version:     14,
			Description: "create meal plan option votes table",
			Script:      fetchMigration("migrations/00014_meal_plan_option_votes.sql"),
		},
		{
			Version:     15,
			Description: "create meal plan option votes table",
			Script:      fetchMigration("migrations/00015_recipe_step_updates.sql"),
		},
		{
			Version:     16,
			Description: "reintroduce recipe step products table",
			Script:      fetchMigration("migrations/00016_recipe_step_products.sql"),
		},
		{
			Version:     17,
			Description: "remove yields from recipe steps table",
			Script:      fetchMigration("migrations/00017_remove_yields_from_recipe_steps.sql"),
		},
		{
			Version:     18,
			Description: "add birthdate fields to user table",
			Script:      fetchMigration("migrations/00018_add_user_birthdate_data.sql"),
		},
		{
			Version:     19,
			Description: "add miscellaneous indices",
			Script:      fetchMigration("migrations/00019_indices_catchup.sql"),
		},
		{
			Version:     20,
			Description: "replace invalid uniqueness constraint on valid_ingredients table",
			Script:      fetchMigration("migrations/00020_instrument_uniqueness_constraint_fix.sql"),
		},
	}
)

func fetchMigration(name string) string {
	file, err := rawMigrations.ReadFile(name)
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
