package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	queryOutput := map[string][]*Query{
		"admin.sql":    buildAdminQueries(),
		"webhooks.sql": buildWebhooksQueries(),
		"users.sql":    buildUsersQueries(),
	}

	for filePath, queries := range queryOutput {
		existingFile, err := os.ReadFile(path.Join("internal", "database", "postgres", "sqlc_queries", filePath))
		if err != nil {
			panic(err)
		}

		var fileContent string
		for i, query := range queries {
			if i != 0 {
				fileContent += "\n\n"
			}
			fileContent += fmt.Sprintf("-- name: %s %s\n\n%s", query.Annotation.Name, query.Annotation.Type, formatQuery(query.Content))
		}
		fileContent += "\n"

		if string(existingFile) != fileContent {
			fmt.Printf("files don't match: %s\n", filePath)
		}
	}

	fmt.Println("done")
}
