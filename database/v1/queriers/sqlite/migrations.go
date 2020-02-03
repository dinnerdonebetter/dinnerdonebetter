package sqlite

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
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"password_last_changed_on" INTEGER,
				"two_factor_secret" TEXT NOT NULL,
				"is_admin" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				CONSTRAINT username_unique UNIQUE (username)
			);`,
		},
		{
			Version:     2,
			Description: "create oauth2_clients table",
			Script: `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     3,
			Description: "create webhooks table",
			Script: `
			CREATE TABLE IF NOT EXISTS webhooks (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT NOT NULL,
				"content_type" TEXT NOT NULL,
				"url" TEXT NOT NULL,
				"method" TEXT NOT NULL,
				"events" TEXT NOT NULL,
				"data_types" TEXT NOT NULL,
				"topics" TEXT NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     4,
			Description: "create instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS instruments (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" CHARACTER VARYING NOT NULL,
				"variant" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     5,
			Description: "create ingredients table",
			Script: `
			CREATE TABLE IF NOT EXISTS ingredients (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
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
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     6,
			Description: "create preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS preparations (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" CHARACTER VARYING NOT NULL,
				"variant" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"allergy_warning" CHARACTER VARYING NOT NULL,
				"icon" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     7,
			Description: "create required preparation instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS required_preparation_instruments (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"instrument_id" INTEGER NOT NULL,
				"preparation_id" INTEGER NOT NULL,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     8,
			Description: "create recipes table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipes (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" CHARACTER VARYING NOT NULL,
				"source" CHARACTER VARYING NOT NULL,
				"description" CHARACTER VARYING NOT NULL,
				"inspired_by_recipe_id" INTEGER,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     9,
			Description: "create recipe steps table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_steps (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"index" INTEGER NOT NULL,
				"preparation_id" INTEGER NOT NULL,
				"prerequisite_step" INTEGER NOT NULL,
				"min_estimated_time_in_seconds" INTEGER NOT NULL,
				"max_estimated_time_in_seconds" INTEGER NOT NULL,
				"temperature_in_celsius" INTEGER,
				"notes" CHARACTER VARYING NOT NULL,
				"recipe_id" INTEGER NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     10,
			Description: "create recipe step instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_instruments (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"instrument_id" INTEGER,
				"recipe_step_id" INTEGER NOT NULL,
				"notes" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     11,
			Description: "create recipe step ingredients table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"ingredient_id" INTEGER,
				"quantity_type" CHARACTER VARYING NOT NULL,
				"quantity_value" REAL NOT NULL,
				"quantity_notes" CHARACTER VARYING NOT NULL,
				"product_of_recipe" BOOLEAN NOT NULL,
				"ingredient_notes" CHARACTER VARYING NOT NULL,
				"recipe_step_id" INTEGER NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     12,
			Description: "create recipe step products table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_products (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" CHARACTER VARYING NOT NULL,
				"recipe_step_id" INTEGER NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     13,
			Description: "create recipe iterations table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_iterations (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"recipe_id" INTEGER NOT NULL,
				"end_difficulty_rating" REAL NOT NULL,
				"end_complexity_rating" REAL NOT NULL,
				"end_taste_rating" REAL NOT NULL,
				"end_overall_rating" REAL NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     14,
			Description: "create recipe step events table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_events (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"event_type" CHARACTER VARYING NOT NULL,
				"done" BOOLEAN NOT NULL,
				"recipe_iteration_id" INTEGER NOT NULL,
				"recipe_step_id" INTEGER NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     15,
			Description: "create iteration medias table",
			Script: `
			CREATE TABLE IF NOT EXISTS iteration_medias (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"path" CHARACTER VARYING NOT NULL,
				"mimetype" CHARACTER VARYING NOT NULL,
				"recipe_iteration_id" INTEGER NOT NULL,
				"recipe_step_id" INTEGER,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     16,
			Description: "create invitations table",
			Script: `
			CREATE TABLE IF NOT EXISTS invitations (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"code" CHARACTER VARYING NOT NULL,
				"consumed" BOOLEAN NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
		{
			Version:     17,
			Description: "create reports table",
			Script: `
			CREATE TABLE IF NOT EXISTS reports (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"report_type" CHARACTER VARYING NOT NULL,
				"concern" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`,
		},
	}
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a sqlite database
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.SqliteDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (s *Sqlite) Migrate(ctx context.Context) error {
	s.logger.Info("migrating db")
	if !s.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	s.migrateOnce.Do(buildMigrationFunc(s.db))

	return nil
}
