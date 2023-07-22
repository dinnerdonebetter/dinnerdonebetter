package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/mjibson/sqlfmt"
)

const (
	id                 = "id"
	dummyValue         = "whatever"
	belongsToHousehold = "belongs_to_household"
	createdByUser      = "created_by_user"
	createdAt          = "created_at"
	lastIndexedAt      = "last_indexed_at"
	lastUpdatedAt      = "last_updated_at"
	archivedAt         = "archived_at"
)

var postgresNow = squirrel.Expr("NOW()")

func buildColumnName(table, column string) string {
	return table + "." + column
}

func mergeSlicesAtIndex[T comparable](a, b []T, index uint) []T {
	return append(a[:index], append(b, a[index:]...)...)
}

func buildJoinStatement(join, on, to string) string {
	return fmt.Sprintf("%s ON %s=%s", join, on, to)
}

func writeFileToPath(outputFilepath, content string) error {
	if err := os.Mkdir(filepath.Dir(outputFilepath), os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	// Create the f if it doesn't exist
	f, err := os.Create(outputFilepath)
	if err != nil {
		return fmt.Errorf("failed to create f: %w", err)
	}

	// Write the content to the f
	if _, err = f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to f: %w", err)
	}

	// Close the f
	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to close f: %w", err)
	}

	return nil
}

func formatQuery(query string) (string, error) {
	cfg := tree.DefaultPrettyCfg()
	cfg.UseTabs = true
	cfg.LineWidth = 60
	cfg.TabWidth = 4
	cfg.Simplify = true
	cfg.Align = tree.PrettyAlignOnly
	cfg.Case = strings.ToUpper
	cfg.JSONFmt = true

	return sqlfmt.FmtSQL(cfg, []string{query})
}

type queryFunc func(squirrel.StatementBuilderType) string

const (
	destinationPath = "internal/database/postgres/generated_queries"
)

var (
	fileMap = map[string]queryFunc{
		"webhooks/get_for_user.sql":                        buildGetOneWebhookQuery,
		"webhooks/archive.sql":                             buildArchiveWebhookQuery,
		"webhooks/create.sql":                              buildCreateWebhookQuery,
		"recipes/get_needing_indexing.sql":                 buildSelectRecipesNeedingIndexingQuery,
		"meals/get_needing_indexing.sql":                   buildSelectMealsNeedingIndexingQuery,
		"valid_instruments/get_needing_indexing.sql":       buildSelectValidInstrumentsNeedingIndexingQuery,
		"valid_instruments/get_by_id.sql":                  buildSelectValidInstrumentQuery,
		"valid_ingredients/get_needing_indexing.sql":       buildSelectValidIngredientsNeedingIndexingQuery,
		"valid_ingredients/get_by_id.sql":                  buildSelectValidIngredientQuery,
		"valid_measurement_units/get_needing_indexing.sql": buildSelectValidMeasurementUnitsNeedingIndexingQuery,
		"valid_measurement_units/get_by_id.sql":            buildSelectValidMeasurementUnitQuery,
		"valid_preparations/get_needing_indexing.sql":      buildSelectValidPreparationsNeedingIndexingQuery,
		"valid_preparations/get_by_id.sql":                 buildSelectValidPreparationQuery,
		"valid_ingredient_states/get_needing_indexing.sql": buildSelectValidIngredientStatesNeedingIndexingQuery,
		"valid_ingredient_states/get_by_id.sql":            buildSelectValidIngredientStateQuery,
	}
)

func main() {
	if err := os.MkdirAll(destinationPath, os.ModePerm); err != nil {
		panic(err)
	}

	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	for filename, builder := range fileMap {
		query, err := formatQuery(builder(sqlBuilder))
		if err != nil {
			panic(err)
		}

		if err = writeFileToPath(path.Join(destinationPath, filename), query); err != nil {
			panic(err)
		}
	}
}
