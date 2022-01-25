package main

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"text/template"
	"time"

	"github.com/segmentio/ksuid"

	_ "github.com/lib/pq"
	flag "github.com/spf13/pflag"

	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
)

const defaultDBURL = "postgres://dbuser:hunter2@localhost:5432/prixfixe?sslmode=disable"

var (
	dbString string
	debug    bool

	//go:embed init.sql.tmpl
	initScriptTemplate string
)

func init() {
	flag.StringVarP(&dbString, "dburl", "u", defaultDBURL, "where the database is hosted")
	flag.BoolVarP(&debug, "debug", "d", false, "whether debug mode is enabled")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	logger := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger()

	if dbString == "" {
		logger.Fatal(errors.New("uri must be valid"))
	}

	if dbString == "" {
		logger.Fatal(errors.New("database connection string must be provided"))
	}

	db, err := sql.Open("postgres", dbString)
	if err != nil {
		logger.Fatal(err)
	}

	generatedIDs := map[string]string{}

	templateFuncs := map[string]interface{}{
		"timestamp": func() int64 {
			return time.Now().Unix()
		},
		"getID": func(key string) string {
			if id, ok := generatedIDs[key]; ok {
				return id
			}

			newID := ksuid.New().String()
			generatedIDs[key] = newID

			return newID
		},
	}

	t := template.Must(template.New("init").Funcs(templateFuncs).Parse(initScriptTemplate))

	var query bytes.Buffer
	if err = t.Execute(&query, struct{}{}); err != nil {
		logger.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `
DELETE FROM "users" WHERE id IS NOT NULL;
DELETE FROM "households" WHERE id IS NOT NULL;
DELETE FROM "household_user_memberships" WHERE id IS NOT NULL;
DELETE FROM "valid_ingredients" WHERE id IS NOT NULL;
DELETE FROM "valid_instruments" WHERE id IS NOT NULL;
DELETE FROM "valid_preparations" WHERE id IS NOT NULL;
DELETE FROM "recipes" WHERE id IS NOT NULL;
DELETE FROM "recipe_steps" WHERE id IS NOT NULL;
DELETE FROM "recipe_step_ingredients" WHERE id IS NOT NULL;
DELETE FROM "meals" WHERE id IS NOT NULL;
DELETE FROM "meal_recipes" WHERE id IS NOT NULL;
DELETE FROM "meal_plans" WHERE id IS NOT NULL;
DELETE FROM "meal_plan_options" WHERE id IS NOT NULL;
DELETE FROM "sessions" WHERE data IS NOT NULL;
`)
	if err != nil {
		logger.Fatal(err)
	}

	q := query.String()

	_, err = db.ExecContext(ctx, q)
	if err != nil {
		logger.Fatal(err)
	}
}
