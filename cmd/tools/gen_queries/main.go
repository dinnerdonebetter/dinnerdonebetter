package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/mjibson/sqlfmt"
)

const (
	placeholder      = "placeholder"
	idColumn         = "id"
	archivedAtColumn = "archived_at"
)

func formatQuery(query string) string {
	cfg := tree.PrettyCfg{
		LineWidth: 80,
		Align:     tree.PrettyAlignAndDeindent,
		Simplify:  true,
		UseTabs:   true,
		Case:      strings.ToUpper,
		JSONFmt:   true,
	}

	formatted, err := sqlfmt.FmtSQL(cfg, []string{query})
	if err != nil {
		panic(err)
	}

	return formatted
}

func main() {
	queryOutput := map[string][]Query{
		"admin.sql": buildAdminQueries(),
	}

	for filePath, queries := range queryOutput {
		if len(queries) == 0 {
			continue
		}

		existingFile, err := os.ReadFile(path.Join("internal", "database", "postgres", "sqlc_queries", filePath))
		if err != nil {
			panic(err)
		}

		var fileContent string
		for i, query := range queries {
			if i != 0 {
				fileContent += "\n\n"
			}
			fileContent += fmt.Sprintf("-- name: %s %s\n\n%s", query.Annotation.Name, query.Annotation.Type, query.Content)
		}
		fileContent += "\n"

		if string(existingFile) != fileContent {
			log.Fatalf("files don't match: %s", filePath)
		}
	}
}
