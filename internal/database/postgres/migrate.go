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
			Description: "create valid ingredient preparations table",
			Script:      fetchMigration("00005_valid_ingredient_preparations"),
		},
		{
			Version:     6,
			Description: "create recipes table",
			Script:      fetchMigration("00006_recipes"),
		},
		{
			Version:     7,
			Description: "create recipe steps table",
			Script:      fetchMigration("00007_recipe_steps"),
		},
		{
			Version:     8,
			Description: "create recipe step instruments table",
			Script:      fetchMigration("00008_recipe_step_instruments"),
		},
		{
			Version:     9,
			Description: "create recipe step ingredients table",
			Script:      fetchMigration("00009_recipe_step_ingredients"),
		},
		{
			Version:     10,
			Description: "create recipe step products table",
			Script:      fetchMigration("00010_recipe_step_products"),
		},
		{
			Version:     11,
			Description: "create meals table",
			Script:      fetchMigration("00011_meals"),
		},
		{
			Version:     12,
			Description: "create meal plans table",
			Script:      fetchMigration("00012_meal_plans"),
		},
		{
			Version:     13,
			Description: "create meal plan options table",
			Script:      fetchMigration("00013_meal_plan_options"),
		},
		{
			Version:     14,
			Description: "create meal plan option votes table",
			Script:      fetchMigration("00014_meal_plan_option_votes"),
		},
		{
			Version:     15,
			Description: "create meal plan option votes table",
			Script:      fetchMigration("00015_recipe_step_updates"),
		},
		{
			Version:     16,
			Description: "reintroduce recipe step products table",
			Script:      fetchMigration("00016_recipe_step_products"),
		},
		{
			Version:     17,
			Description: "remove yields from recipe steps table",
			Script:      fetchMigration("00017_remove_yields_from_recipe_steps"),
		},
		{
			Version:     18,
			Description: "add birthdate fields to user table",
			Script:      fetchMigration("00018_add_user_birthdate_data"),
		},
		{
			Version:     19,
			Description: "add miscellaneous indices",
			Script:      fetchMigration("00019_indices_catchup"),
		},
		{
			Version:     20,
			Description: "replace invalid uniqueness constraint on valid_ingredients table",
			Script:      fetchMigration("00020_instrument_uniqueness_constraint_fix"),
		},
		{
			Version:     21,
			Description: "replace invalid uniqueness constraint on valid_ingredients table",
			Script:      fetchMigration("00021_rename_user_status_column"),
		},
		{
			Version:     22,
			Description: "modify valid ingredient fields",
			Script:      fetchMigration("00022_modify_valid_ingredient_fields"),
		},
		{
			Version:     23,
			Description: "add quantity fields to recipe step products",
			Script:      fetchMigration("00023_add_quantity_fields_to_recipe_step_products"),
		},
		{
			Version:     24,
			Description: "drop uniqueness constraints on meal plan optinos",
			Script:      fetchMigration("00024_remove_meal_plan_uniqueness_constraint"),
		},
		{
			Version:     25,
			Description: "add password reset tokens",
			Script:      fetchMigration("00025_add_password_reset_tokens"),
		},
		{
			Version:     26,
			Description: "add valid measurement units",
			Script:      fetchMigration("00026_valid_measurement_units"),
		},
		{
			Version:     27,
			Description: "add valid measurement units",
			Script:      fetchMigration("00027_various_bridge_tables"),
		},
		{
			Version:     28,
			Description: "add recipe step ranges",
			Script:      fetchMigration("00028_recipe_step_ranges"),
		},
		{
			Version:     29,
			Description: "add valid measurement unit constraints",
			Script:      fetchMigration("00029_recipe_step_ingredient_units"),
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
		{
			Version:     32,
			Description: "various recipe step product improvements",
			Script:      fetchMigration("00032_recipe_step_products_improvements"),
		},
		{
			Version:     33,
			Description: "bridge table unique constraint improvements",
			Script:      fetchMigration("00033_better_bridge_table_constraints"),
		},
		{
			Version:     34,
			Description: "add recipe step product id to recipe step instruments",
			Script:      fetchMigration("00034_add_recipe_step_product_id_to_recipe_step_instruments"),
		},
		{
			Version:     35,
			Description: "make instrument id a foreign key for recipe step instruments",
			Script:      fetchMigration("00035_make_instrument_id_a_foreign_key_for_recipe_step_instruments"),
		},
		{
			Version:     36,
			Description: "a ton of various additions",
			Script:      fetchMigration("00036_myriad_schema_additions"),
		},
		{
			Version:     37,
			Description: "rename time fields",
			Script:      fetchMigration("00037_better_time_field_names"),
		},
		{
			Version:     38,
			Description: "add advanced prep notifications",
			Script:      fetchMigration("00038_advanced_prep_notifications"),
		},
		{
			Version:     39,
			Description: "refactor meal plans",
			Script:      fetchMigration("00039_meal_plans_refactor"),
		},
		{
			Version:     40,
			Description: "advanced prep step constraints",
			Script:      fetchMigration("00040_advanced_prep_step_constraints"),
		},
		{
			Version:     41,
			Description: "nullable storage temps for products",
			Script:      fetchMigration("00041_nullable_storage_temps_for_products"),
		},
		{
			Version:     42,
			Description: "nullable storage temps for products",
			Script:      fetchMigration("00042_nullable_min_and_max_quantities_for_recipe_steps"),
		},
		{
			Version:     43,
			Description: "remove constraint on advanced prep steps",
			Script:      fetchMigration("00043_remove_constraint_on_advanced_prep_steps"),
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
func (q *Querier) migrationFunc() {
	driver := darwin.NewGenericDriver(q.db, darwin.PostgresDialect{})
	if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
		panic(fmt.Errorf("migrating database: %w", err))
	}
}

// Migrate is a simple wrapper around the core querier Migrate call.
func (q *Querier) Migrate(ctx context.Context, maxAttempts uint8) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	q.logger.Info("migrating db")

	if !q.IsReady(ctx, maxAttempts) {
		return database.ErrDatabaseNotReady
	}

	q.migrateOnce.Do(q.migrationFunc)

	return nil
}
