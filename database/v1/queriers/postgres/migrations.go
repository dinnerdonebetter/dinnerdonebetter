package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: `
			CREATE TABLE IF NOT EXISTS users (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"password_last_changed_on" integer,
				"two_factor_secret" TEXT NOT NULL,
				"is_admin" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				UNIQUE ("username")
			);`,
		},
		{
			Version:     2,
			Description: "create oauth2_clients table",
			Script: `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     3,
			Description: "create webhooks table",
			Script: `
			CREATE TABLE IF NOT EXISTS webhooks (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" TEXT NOT NULL,
				"content_type" TEXT NOT NULL,
				"url" TEXT NOT NULL,
				"method" TEXT NOT NULL,
				"events" TEXT NOT NULL,
				"data_types" TEXT NOT NULL,
				"topics" TEXT NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     4,
			Description: "create valid instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_instruments (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"variant" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     5,
			Description: "create valid ingredients table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_ingredients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"variant" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"warning" CHARACTER VARYING NOT NULL,
				"contains_egg" BOOLEAN NOT NULL,
				"contains_dairy" BOOLEAN NOT NULL,
				"contains_peanut" BOOLEAN NOT NULL,
				"contains_tree_nut" BOOLEAN NOT NULL,
				"contains_soy" BOOLEAN NOT NULL,
				"contains_wheat" BOOLEAN NOT NULL,
				"contains_shellfish" BOOLEAN NOT NULL,
				"contains_sesame" BOOLEAN NOT NULL,
				"contains_fish" BOOLEAN NOT NULL,
				"contains_gluten" BOOLEAN NOT NULL,
				"animal_flesh" BOOLEAN NOT NULL,
				"animal_derived" BOOLEAN NOT NULL,
				"considered_staple" BOOLEAN NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     6,
			Description: "create valid ingredient tags table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_ingredient_tags (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     7,
			Description: "create ingredient tag mappings table",
			Script: `
			CREATE TABLE IF NOT EXISTS ingredient_tag_mappings (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"valid_ingredient_tag_id" BIGINT NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_valid_ingredient" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_valid_ingredient") REFERENCES "valid_ingredients"("id")
			);`,
		},
		{
			Version:     8,
			Description: "create valid preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_preparations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     9,
			Description: "create required preparation instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS required_preparation_instruments (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"valid_instrument_id" BIGINT NOT NULL,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_valid_preparation" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_valid_preparation") REFERENCES "valid_preparations"("id")
			);`,
		},
		{
			Version:     10,
			Description: "create valid ingredient preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_valid_ingredient" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_valid_ingredient") REFERENCES "valid_ingredients"("id")
			);`,
		},
		{
			Version:     11,
			Description: "create recipes table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipes (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"source" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"inspired_by_recipe_id" BIGINT,
				"private" BOOLEAN NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     12,
			Description: "create recipe tags table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_tags (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe") REFERENCES "recipes"("id")
			);`,
		},
		{
			Version:     13,
			Description: "create recipe steps table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_steps (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"index" INTEGER NOT NULL,
				"valid_preparation_id" BIGINT NOT NULL,
				"prerequisite_step_id" BIGINT,
				"min_estimated_time_in_seconds" BIGINT NOT NULL,
				"max_estimated_time_in_seconds" BIGINT NOT NULL,
				"yields_product_name" CHARACTER VARYING NOT NULL,
				"yields_quantity" INTEGER NOT NULL,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe") REFERENCES "recipes"("id")
			);`,
		},
		{
			Version:     14,
			Description: "create recipe step preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_preparations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"valid_preparation_id" BIGINT NOT NULL,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_step" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_step") REFERENCES "recipe_steps"("id")
			);`,
		},
		{
			Version:     15,
			Description: "create recipe step ingredients table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"valid_ingredient_id" BIGINT NOT NULL,
				"ingredient_notes" CHARACTER VARYING NOT NULL,
				"quantity_type" CHARACTER VARYING NOT NULL,
				"quantity_value" DOUBLE PRECISION NOT NULL,
				"quantity_notes" CHARACTER VARYING NOT NULL,
				"product_of_recipe_step_id" BIGINT,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_step" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_step") REFERENCES "recipe_steps"("id")
			);`,
		},
		{
			Version:     16,
			Description: "create recipe iterations table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_iterations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"end_difficulty_rating" DOUBLE PRECISION NOT NULL,
				"end_complexity_rating" DOUBLE PRECISION NOT NULL,
				"end_taste_rating" DOUBLE PRECISION NOT NULL,
				"end_overall_rating" DOUBLE PRECISION NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe") REFERENCES "recipes"("id")
			);`,
		},
		{
			Version:     17,
			Description: "create recipe iteration steps table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_iteration_steps (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"started_on" BIGINT,
				"ended_on" BIGINT,
				"state" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe") REFERENCES "recipes"("id")
			);`,
		},
		{
			Version:     18,
			Description: "create iteration medias table",
			Script: `
			CREATE TABLE IF NOT EXISTS iteration_medias (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"source" CHARACTER VARYING NOT NULL,
				"mimetype" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_iteration" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_iteration") REFERENCES "recipe_iterations"("id")
			);`,
		},
		{
			Version:     19,
			Description: "create invitations table",
			Script: `
			CREATE TABLE IF NOT EXISTS invitations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"code" CHARACTER VARYING NOT NULL,
				"consumed" BOOLEAN NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     20,
			Description: "create reports table",
			Script: `
			CREATE TABLE IF NOT EXISTS reports (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"report_type" CHARACTER VARYING NOT NULL,
				"concern" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
	}
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a postgres database.
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (p *Postgres) Migrate(ctx context.Context) error {
	p.logger.Info("migrating db")
	if !p.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	p.migrateOnce.Do(buildMigrationFunc(p.db))

	return nil
}
