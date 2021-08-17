package postgres

import (
	"database/sql"
	"fmt"

	"github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create sessions table for session manager",
			Script: `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BYTEA NOT NULL,
				expiry TIMESTAMPTZ NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())
			);`,
		},
		{
			Version:     2,
			Description: "create sessions table for session manager",
			Script:      `CREATE INDEX sessions_expiry_idx ON sessions (expiry);`,
		},
		{
			Version:     3,
			Description: "create audit log table",
			Script: `
			CREATE TABLE IF NOT EXISTS audit_log (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				event_type TEXT NOT NULL,
				context JSONB NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())
			);`,
		},
		{
			Version:     4,
			Description: "create users table",
			Script: `
			CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				username TEXT NOT NULL,
				avatar_src TEXT,
				hashed_password TEXT NOT NULL,
				password_last_changed_on INTEGER,
				requires_password_change BOOLEAN NOT NULL DEFAULT 'false',
				two_factor_secret TEXT NOT NULL,
				two_factor_secret_verified_on BIGINT DEFAULT NULL,
				service_roles TEXT NOT NULL DEFAULT 'service_user',
				reputation TEXT NOT NULL DEFAULT 'unverified',
				reputation_explanation TEXT NOT NULL DEFAULT '',
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				UNIQUE("username")
			);`,
		},
		{
			Version:     5,
			Description: "create accounts table",
			Script: `
			CREATE TABLE IF NOT EXISTS accounts (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				billing_status TEXT NOT NULL DEFAULT 'unpaid',
				contact_email TEXT NOT NULL DEFAULT '',
				contact_phone TEXT NOT NULL DEFAULT '',
				payment_processor_customer_id TEXT NOT NULL DEFAULT '',
				subscription_plan_id TEXT,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_user BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				UNIQUE("belongs_to_user", "name")
			);`,
		},
		{
			Version:     6,
			Description: "create account user memberships table",
			Script: `
			CREATE TABLE IF NOT EXISTS account_user_memberships (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
				belongs_to_user BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				default_account BOOLEAN NOT NULL DEFAULT 'false',
				account_roles TEXT NOT NULL DEFAULT 'account_user',
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				UNIQUE("belongs_to_account", "belongs_to_user")
			);`,
		},
		{
			Version:     7,
			Description: "create API clients table",
			Script: `
			CREATE TABLE IF NOT EXISTS api_clients (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT DEFAULT '',
				client_id TEXT NOT NULL,
				secret_key BYTEA NOT NULL,
				permissions BIGINT NOT NULL DEFAULT 0,
				admin_permissions BIGINT NOT NULL DEFAULT 0,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_user BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
			);`,
		},
		{
			Version:     8,
			Description: "create webhooks table",
			Script: `
			CREATE TABLE IF NOT EXISTS webhooks (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				content_type TEXT NOT NULL,
				url TEXT NOT NULL,
				method TEXT NOT NULL,
				events TEXT NOT NULL,
				data_types TEXT NOT NULL,
				topics TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`,
		},
		{
			Version:     9,
			Description: "create valid instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_instruments (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				variant TEXT NOT NULL,
				description TEXT NOT NULL,
				icon_path TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				UNIQUE("name", "variant")
			);`,
		},
		{
			Version:     10,
			Description: "create valid preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_preparations (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				description TEXT NOT NULL,
				icon_path TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				UNIQUE("name")
			);`,
		},
		{
			Version:     11,
			Description: "create valid ingredients table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_ingredients (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				variant TEXT NOT NULL,
				description TEXT NOT NULL,
				warning TEXT NOT NULL,
				contains_egg BOOLEAN NOT NULL,
				contains_dairy BOOLEAN NOT NULL,
				contains_peanut BOOLEAN NOT NULL,
				contains_tree_nut BOOLEAN NOT NULL,
				contains_soy BOOLEAN NOT NULL,
				contains_wheat BOOLEAN NOT NULL,
				contains_shellfish BOOLEAN NOT NULL,
				contains_sesame BOOLEAN NOT NULL,
				contains_fish BOOLEAN NOT NULL,
				contains_gluten BOOLEAN NOT NULL,
				animal_flesh BOOLEAN NOT NULL,
				animal_derived BOOLEAN NOT NULL,
				volumetric BOOLEAN NOT NULL,
				icon_path TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				UNIQUE("name", "variant")
			);`,
		},
		{
			Version:     12,
			Description: "create valid ingredient preparations table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				notes TEXT NOT NULL,
				valid_ingredient_id BIGINT NOT NULL REFERENCES valid_ingredients(id),
				valid_preparation_id BIGINT NOT NULL REFERENCES valid_preparations(id),
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     13,
			Description: "create valid preparation instruments table",
			Script: `
			CREATE TABLE IF NOT EXISTS valid_preparation_instruments (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				instrument_id BIGINT NOT NULL,
				preparation_id BIGINT NOT NULL,
				notes TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL
			);`,
		},
		{
			Version:     14,
			Description: "create recipes table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipes (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				source TEXT NOT NULL,
				description TEXT NOT NULL,
				inspired_by_recipe_id BIGINT,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`,
		},
		{
			Version:     15,
			Description: "create recipe steps table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_steps (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				index INTEGER NOT NULL,
				preparation_id BIGINT NOT NULL REFERENCES valid_preparations(id),
				prerequisite_step BIGINT NOT NULL,
				min_estimated_time_in_seconds BIGINT NOT NULL,
				max_estimated_time_in_seconds BIGINT NOT NULL,
				temperature_in_celsius INTEGER,
				notes TEXT NOT NULL,
				why TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_recipe BIGINT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE
			);`,
		},
		{
			Version:     16,
			Description: "create recipe step ingredients table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				ingredient_id BIGINT NOT NULL REFERENCES valid_ingredients(id),
				name TEXT NOT NULL,
				quantity_type TEXT NOT NULL,
				quantity_value DOUBLE PRECISION NOT NULL,
				quantity_notes TEXT NOT NULL,
				product_of_recipe_step BOOLEAN NOT NULL,
				ingredient_notes TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_recipe_step BIGINT NOT NULL REFERENCES recipe_steps(id) ON DELETE CASCADE
			);`,
		},
		{
			Version:     17,
			Description: "create recipe step products table",
			Script: `
			CREATE TABLE IF NOT EXISTS recipe_step_products (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				quantity_type TEXT NOT NULL,
				quantity_value DOUBLE PRECISION NOT NULL,
				quantity_notes TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_recipe_step BIGINT NOT NULL REFERENCES recipe_steps(id) ON DELETE CASCADE
			);`,
		},
		{
			Version:     18,
			Description: "create invitations table",
			Script: `
			CREATE TABLE IF NOT EXISTS invitations (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				code TEXT NOT NULL,
				consumed BOOLEAN NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`,
		},
		{
			Version:     19,
			Description: "create reports table",
			Script: `
			CREATE TABLE IF NOT EXISTS reports (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				report_type TEXT NOT NULL,
				concern TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`,
		},
	}
)

// BuildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a postgres database.
func (b *Postgres) BuildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(fmt.Errorf("migrating database: %w", err))
		}
	}
}
