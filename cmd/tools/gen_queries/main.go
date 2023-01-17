package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"os"
	"path"
	"regexp"
	"strings"
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

func writeFileToPath(path, content string) error {
	// Create the file if it doesn't exist
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// Close the file
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
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
		"webhooks/get_one.sql": buildGetOneWebhookQuery,
		"webhooks/archive.sql": buildArchiveWebhookQuery,
		"webhooks/create.sql":  buildCreateWebhookQuery,
	}
)

func main() {
	cwd, _ := os.Getwd()
	_ = cwd

	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	if err := os.MkdirAll(destinationPath, 0755); err != nil {
		panic(err)
	}

	for filename, builder := range fileMap {
		query := formatQuery(builder(sqlBuilder))
		if err := writeFileToPath(path.Join(destinationPath, filename), query); err != nil {
			panic(err)
		}
	}
}
