package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/GuiaBolso/darwin"

	"github.com/prixfixeco/api_server/internal/database"
)

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

var (
	//go:embed migrations/00001_initial.sql
	initMigration string

	//go:embed migrations/00002_valid_instruments.sql
	validInstrumentsMigration string

	//go:embed migrations/00003_valid_ingredients.sql
	validIngredientsMigration string

	//go:embed migrations/00004_valid_preparations.sql
	validPreparationsMigration string

	//go:embed migrations/00005_valid_ingredient_preparations.sql
	validIngredientPreparationsMigration string

	//go:embed migrations/00006_recipes.sql
	recipesMigration string

	//go:embed migrations/00007_recipe_steps.sql
	recipeStepsMigration string

	//go:embed migrations/00008_recipe_step_instruments.sql
	recipeStepInstrumentsMigration string

	//go:embed migrations/00009_recipe_step_ingredients.sql
	recipeStepIngredientsMigration string

	//go:embed migrations/00010_recipe_step_products.sql
	recipeStepProductsMigration string

	//go:embed migrations/00011_meals.sql
	mealsMigration string

	//go:embed migrations/00012_meal_plans.sql
	mealPlansMigration string

	//go:embed migrations/00013_meal_plan_options.sql
	mealPlanOptionsMigration string

	//go:embed migrations/00014_meal_plan_option_votes.sql
	mealPlanOptionVotesMigration string

	//go:embed migrations/00015_recipe_step_updates.sql
	recipeStepsAdjustmentMigration1 string

	//go:embed migrations/00016_recipe_step_products.sql
	recipeStepProductsReintroductionMigration string

	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "basic infrastructural tables",
			Script:      initMigration,
		},
		{
			Version:     2,
			Description: "create valid instruments table",
			Script:      validInstrumentsMigration,
		},
		{
			Version:     3,
			Description: "create valid ingredients table",
			Script:      validIngredientsMigration,
		},
		{
			Version:     4,
			Description: "create valid preparations table",
			Script:      validPreparationsMigration,
		},
		{
			Version:     5,
			Description: "create valid ingredient preparations table",
			Script:      validIngredientPreparationsMigration,
		},
		{
			Version:     6,
			Description: "create recipes table",
			Script:      recipesMigration,
		},
		{
			Version:     7,
			Description: "create recipe steps table",
			Script:      recipeStepsMigration,
		},
		{
			Version:     8,
			Description: "create recipe step instruments table",
			Script:      recipeStepInstrumentsMigration,
		},
		{
			Version:     9,
			Description: "create recipe step ingredients table",
			Script:      recipeStepIngredientsMigration,
		},
		{
			Version:     10,
			Description: "create recipe step products table",
			Script:      recipeStepProductsMigration,
		},
		{
			Version:     11,
			Description: "create meals table",
			Script:      mealsMigration,
		},
		{
			Version:     12,
			Description: "create meal plans table",
			Script:      mealPlansMigration,
		},
		{
			Version:     13,
			Description: "create meal plan options table",
			Script:      mealPlanOptionsMigration,
		},
		{
			Version:     14,
			Description: "create meal plan option votes table",
			Script:      mealPlanOptionVotesMigration,
		},
		{
			Version:     15,
			Description: "create meal plan option votes table",
			Script:      recipeStepsAdjustmentMigration1,
		},
		{
			Version:     16,
			Description: "reintroduce recipe step products table",
			Script:      recipeStepProductsReintroductionMigration,
		},
	}
)

// BuildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a postgres database.
func (q *SQLQuerier) migrationFunc() {
	driver := darwin.NewGenericDriver(q.db, darwin.PostgresDialect{})
	if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
		panic(fmt.Errorf("migrating database: %w", err))
	}
}
