package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Masterminds/squirrel"
)

const (
	id                 = "id"
	dummyValue         = "whatever"
	belongsToHousehold = "belongs_to_household"
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

func writeFileToPath(filepath, content string) error {
	// Create the f if it doesn't exist
	f, err := os.Create(filepath)
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

func formatQuery(query string) string {
	commonReplacer := strings.NewReplacer(
		",", ",\n\t",
		"FROM", "\nFROM",
		"SET", "\nSET\n\t",
		"WHERE", "\nWHERE",
		") VALUES (", "\n) VALUES (\n\t",
		"AND", "\n\tAND",
		"SELECT", "SELECT\n\t",
		"JOIN", "\n\tJOIN",
	)

	lastArgReplacer := regexp.MustCompile(`\$(\d+)\)`)

	outboundQuery := lastArgReplacer.ReplaceAllString(commonReplacer.Replace(query), "$$$1\n)")

	return fmt.Sprintf("%s;\n", outboundQuery)
}

type queryFunc func(squirrel.StatementBuilderType) string

const (
	destinationPath = "internal/database/postgres/generated_queries"
)

var (
	fileMap = map[string]queryFunc{
		"webhooks/get_for_user.sql": buildGetOneWebhookQuery,
		"webhooks/archive.sql":      buildArchiveWebhookQuery,
		"webhooks/create.sql":       buildCreateWebhookQuery,
	}
)

func main() {
	if err := os.MkdirAll(destinationPath, os.ModePerm); err != nil {
		panic(err)
	}

	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	for filename, builder := range fileMap {
		query := formatQuery(builder(sqlBuilder))
		if err := writeFileToPath(path.Join(destinationPath, filename), query); err != nil {
			panic(err)
		}
	}
}
