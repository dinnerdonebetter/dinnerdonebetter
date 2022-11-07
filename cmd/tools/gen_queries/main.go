package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/squirrel"
)

const (
	whatever = "whatever"
)

type (
	queryFileGenerator func(ctx context.Context, psql squirrel.StatementBuilderType, format bool, filePath string) error

	queryFileConfig struct {
		generator queryFileGenerator
		outFile   string
		format    bool
	}
)

func main() {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	ctx := context.Background()

	var files = []*queryFileConfig{
		{
			outFile:   "internal/database/postgres/generated_queries/webhooks_get_one.sql",
			format:    true,
			generator: buildWebhooksGetOne,
		},
	}

	for _, x := range files {
		if err := x.generator(ctx, psql, x.format, x.outFile); err != nil {
			log.Fatalf("error rendering %s: %v", x.outFile, err)
		}
	}
}

func saveFile(_ context.Context, content, outputPath string) error {
	/* #nosec G301 */
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o777); err != nil {
		// okay, who gives a shit?
		_ = err
	}

	/* #nosec G306 */
	return os.WriteFile(outputPath, []byte(content), 0o644)
}

func joinTableName(tableName string, columns []string) []string {
	out := []string{}

	for _, col := range columns {
		out = append(out, fmt.Sprintf("%s.%s", tableName, col))
	}

	return out
}

func mergeStringSlicesAtIndex(s1, s2 []string, i int) []string {
	if i <= 0 {
		i = len(s1)
	}

	out := []string{}
	for j, s := range s1 {
		if i == j {
			out = append(out, s2...)
		}
		out = append(out, s)
	}

	return out
}

func formatQuery(q string) string {
	replacements := map[string]string{
		"SELECT": "SELECT\n\t",
		"FROM":   "\nFROM",
		"JOIN":   "\n\tJOIN",
		"WHERE":  "\nWHERE",
		"AND":    "\n\tAND",
	}

	x := strings.ReplaceAll(q, ",", ",\n\t")

	for to, from := range replacements {
		x = strings.ReplaceAll(x, to, from)
	}

	return x
}
