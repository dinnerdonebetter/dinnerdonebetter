package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GuiaBolso/darwin"
)

var currentMigration float64 = 0

func incrementMigrationVersion() float64 {
	currentMigration++
	return currentMigration
}

var (
	migrations = []darwin.Migration{
		{
			Version:     incrementMigrationVersion(),
			Description: "create users table",
			Script: `
			CREATE TABLE IF NOT EXISTS users (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"salt" BYTEA NOT NULL,
				"password_last_changed_on" integer,
				"requires_password_change" boolean NOT NULL DEFAULT 'false',
				"two_factor_secret" TEXT NOT NULL,
				"two_factor_secret_verified_on" BIGINT DEFAULT NULL,
				"is_admin" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				UNIQUE ("username")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create sessions table for session manager",
			Script: `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BYTEA NOT NULL,
				expiry TIMESTAMPTZ NOT NULL
			);

			CREATE INDEX sessions_expiry_idx ON sessions (expiry);
		`,
		},
		{
			Version:     incrementMigrationVersion(),
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
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
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
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create valid instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_instruments (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"variant" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
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
				"measurable_by_volume" BOOLEAN NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create valid preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_preparations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"applicable_to_all_ingredients" BOOLEAN NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create valid ingredient preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"notes" CHARACTER VARYING NOT NULL,
				"valid_preparation_id" BIGINT NOT NULL,
				"valid_ingredient_id" BIGINT NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create required preparation instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS required_preparation_instruments (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"instrument_id" BIGINT NOT NULL,
				"preparation_id" BIGINT NOT NULL,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create recipes table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipes (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"source" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"inspired_by_recipe_id" BIGINT,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create recipe steps table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_steps (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"index" INTEGER NOT NULL,
				"preparation_id" BIGINT NOT NULL,
				"prerequisite_step" BIGINT NOT NULL,
				"min_estimated_time_in_seconds" BIGINT NOT NULL,
				"max_estimated_time_in_seconds" BIGINT NOT NULL,
				"temperature_in_celsius" INTEGER,
				"notes" CHARACTER VARYING NOT NULL,
				"recipe_id" BIGINT NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe") REFERENCES "recipes"("id")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create recipe step instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_instruments (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"instrument_id" BIGINT,
				"recipe_step_id" BIGINT NOT NULL,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_step" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_step") REFERENCES "recipe_steps"("id")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create recipe step ingredients table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"ingredient_id" BIGINT,
				"quantity_type" CHARACTER VARYING NOT NULL,
				"quantity_value" DOUBLE PRECISION NOT NULL,
				"quantity_notes" CHARACTER VARYING NOT NULL,
				"product_of_recipe" BOOLEAN NOT NULL,
				"ingredient_notes" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_step" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_step") REFERENCES "recipe_steps"("id")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create recipe step products table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_products (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"recipe_step_id" BIGINT NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_step" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_step") REFERENCES "recipe_steps"("id")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create recipe iterations table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_iterations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"recipe_id" BIGINT NOT NULL,
				"end_difficulty_rating" DOUBLE PRECISION NOT NULL,
				"end_complexity_rating" DOUBLE PRECISION NOT NULL,
				"end_taste_rating" DOUBLE PRECISION NOT NULL,
				"end_overall_rating" DOUBLE PRECISION NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe") REFERENCES "recipes"("id")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create recipe step events table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_events (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"event_type" CHARACTER VARYING NOT NULL,
				"done" BOOLEAN NOT NULL,
				"recipe_iteration_id" BIGINT NOT NULL,
				"recipe_step_id" BIGINT NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_step" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_step") REFERENCES "recipe_steps"("id")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create iteration medias table",
			Script: `
			CREATE TABLE IF NOT EXISTS iteration_medias (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"path" CHARACTER VARYING NOT NULL,
				"mimetype" CHARACTER VARYING NOT NULL,
				"recipe_iteration_id" BIGINT NOT NULL,
				"recipe_step_id" BIGINT,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_recipe_iteration" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_recipe_iteration") REFERENCES "recipe_iterations"("id")
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create invitations table",
			Script: `
			CREATE TABLE IF NOT EXISTS invitations (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"code" CHARACTER VARYING NOT NULL,
				"consumed" BOOLEAN NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);`,
		},
		{
			Version:     incrementMigrationVersion(),
			Description: "create reports table",
			Script: `
			CREATE TABLE IF NOT EXISTS reports (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"report_type" CHARACTER VARYING NOT NULL,
				"concern" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
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
func (p *Postgres) Migrate(ctx context.Context, createDevUser bool) error {
	p.logger.Info("migrating db")

	if createDevUser {
		// this is gross as fuck and you should find a way to get rid of it.
		migrations = append(migrations, darwin.Migration{
			Version:     incrementMigrationVersion(),
			Description: "create example dev user",
			Script: `
			INSERT INTO users (
				"username",
				"hashed_password",
				"two_factor_secret",
				"password_last_changed_on",
				"is_admin",
				"salt",
				"two_factor_secret_verified_on",
				"created_on",
				"last_updated_on",
				"archived_on"
			) 
				VALUES 
			(
				'username', -- username
				-- '$2a$13$UGQH6IAfUtPH2rCfBF6CJumwcr9Gnb3s.J3vOb/bmWPlx2pJ8.zOe', -- hashed_password
				'$2a$10$JzD3CNBqPmwq.IidQuO7eu3zKdu8vEIi3HkLk8/qRjrzb7eNLKlKG', -- hashed_password
				'TS2ZIIDTN6PBS3A6NTMU637VDHAFZR4BQTPDOBJDFY4JCY2A6OY4TQWJ3WN3VXZB45DDAUD4VJTX2JGHFPVJJDQYP6G3YYUQ2W6RYBA=', -- two_factor_secret
				NULL, -- password_last_changed_on
				'true', -- is_admin
				E'\\x56c066aa6cf769f391594a615aad26cb', --salt
				extract(epoch FROM NOW()), -- two_factor_secret_verified_on
				extract(epoch FROM NOW()), -- created_on
				NULL, -- updated_on
				NULL -- archived_on
			);`,
		},
		)
	}

	if !p.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	p.migrateOnce.Do(buildMigrationFunc(p.db))

	return nil
}
