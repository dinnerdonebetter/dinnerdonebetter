package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/GuiaBolso/darwin"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// defaultTestUserTwoFactorSecret is the default TwoFactorSecret we give to test users when we initialize them.
	// `otpauth://totp/todo:username?secret=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=&issuer=todo`
	defaultTestUserTwoFactorSecret = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

	testUserExistenceQuery = `
		SELECT users.id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.username = $1 AND users.two_factor_secret_verified_on IS NOT NULL
	`

	testUserCreationQuery = `
		INSERT INTO users (id,username,hashed_password,two_factor_secret,reputation,service_roles,two_factor_secret_verified_on) VALUES ($1,$2,$3,$4,$5,$6,extract(epoch FROM NOW()))
	`
)

// Migrate is a simple wrapper around the core querier Migrate call.
func (q *SQLQuerier) Migrate(ctx context.Context, maxAttempts uint8, testUserConfig *types.TestUserCreationConfig) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	q.logger.Info("migrating db")

	if !q.IsReady(ctx, maxAttempts) {
		return database.ErrDatabaseNotReady
	}

	q.migrateOnce.Do(q.migrationFunc)

	if testUserConfig != nil {
		q.logger.Debug("creating test user")

		testUserExistenceArgs := []interface{}{testUserConfig.Username}

		userRow := q.getOneRow(ctx, q.db, "user", testUserExistenceQuery, testUserExistenceArgs)
		_, _, _, err := q.scanUser(ctx, userRow, false)
		if err != nil {
			if testUserConfig.ID == "" {
				testUserConfig.ID = ksuid.New().String()
			}

			testUserCreationArgs := []interface{}{
				testUserConfig.ID,
				testUserConfig.Username,
				testUserConfig.HashedPassword,
				defaultTestUserTwoFactorSecret,
				types.GoodStandingHouseholdStatus,
				authorization.ServiceAdminRole.String(),
			}

			// these structs will be fleshed out by createUser
			user := &types.User{
				ID:       testUserConfig.ID,
				Username: testUserConfig.Username,
			}
			household := &types.Household{
				ID: ksuid.New().String(),
			}

			if err = q.createUser(ctx, user, household, testUserCreationQuery, testUserCreationArgs); err != nil {
				return observability.PrepareError(err, q.logger, span, "creating test user")
			}
			q.logger.WithValue(keys.UsernameKey, testUserConfig.Username).Debug("created test user and household")
		}
	}

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

	//go:embed migrations/000010_recipe_step_products.sql
	recipeStepProductsMigration string

	//go:embed migrations/000011_meal_plans.sql
	mealPlansMigration string

	//go:embed migrations/000012_meal_plan_options.sql
	mealPlanOptionsMigration string

	//go:embed migrations/000013_meal_plan_option_votes.sql
	mealPlanOptionVotesMigration string

	migrations = []darwin.Migration{
		{
			Version:     0.01,
			Description: "basic infrastructural tables",
			Script:      initMigration,
		},
		{
			Version:     0.02,
			Description: "create valid instruments table",
			Script:      validInstrumentsMigration,
		},
		{
			Version:     0.03,
			Description: "create valid ingredients table",
			Script:      validIngredientsMigration,
		},
		{
			Version:     0.04,
			Description: "create valid preparations table",
			Script:      validPreparationsMigration,
		},
		{
			Version:     0.05,
			Description: "create valid ingredient preparations table",
			Script:      validIngredientPreparationsMigration,
		},
		{
			Version:     0.06,
			Description: "create recipes table",
			Script:      recipesMigration,
		},
		{
			Version:     0.07,
			Description: "create recipe steps table",
			Script:      recipeStepsMigration,
		},
		{
			Version:     0.08,
			Description: "create recipe step instruments table",
			Script:      recipeStepInstrumentsMigration,
		},
		{
			Version:     0.09,
			Description: "create recipe step ingredients table",
			Script:      recipeStepIngredientsMigration,
		},
		{
			Version:     0.1,
			Description: "create recipe step products table",
			Script:      recipeStepProductsMigration,
		},
		{
			Version:     0.11,
			Description: "create meal plans table",
			Script:      mealPlansMigration,
		},
		{
			Version:     0.12,
			Description: "create meal plan options table",
			Script:      mealPlanOptionsMigration,
		},
		{
			Version:     0.13,
			Description: "create meal plan option votes table",
			Script:      mealPlanOptionVotesMigration,
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
