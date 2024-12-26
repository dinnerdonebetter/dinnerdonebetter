package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	structs := parseGoFiles(dir)

	// Start extraction from the main struct
	if mainAST, found := structs["config.APIServiceConfig"]; found {
		envVars := extractEnvVars(mainAST, structs, "main", "")
		fmt.Println("Environment Variables:", envVars)
	} else {
		fmt.Printf("Main struct not found.\n")
	}
}

// parseGoFiles parses all Go files in the given directory and returns a map of struct names to their AST nodes.
func parseGoFiles(dir string) map[string]*ast.TypeSpec {
	structs := make(map[string]*ast.TypeSpec)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		if strings.Contains(path, "vendor") {
			return filepath.SkipDir
		}

		fs := token.NewFileSet()
		node, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", path, err)
			return nil
		}

		packageName := node.Name.Name

		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
				if !isTypeSpec {
					continue
				}

				// Only store struct types
				if _, ok = typeSpec.Type.(*ast.StructType); ok {
					key := packageName + "." + typeSpec.Name.Name
					structs[key] = typeSpec
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	return structs
}

// extractEnvVars traverses a struct definition and collects environment variables, resolving nested structs.
func extractEnvVars(typeSpec *ast.TypeSpec, structs map[string]*ast.TypeSpec, currentPackage, prefix string) []string {
	var envVars []string

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return envVars
	}

	for _, field := range structType.Fields.List {
		if field.Tag == nil {
			continue
		}

		tag := strings.Trim(field.Tag.Value, "`")
		if tag == `json:"-"` {
			continue
		}

		envValue := getTagValue(tag, "env")
		if envValue != "" {
			if prefix != "" {
				envValue = prefix + envValue
			}
			envVars = append(envVars, envValue)
		}

		if prefix == "OBSERVABILITY_" {
			println("")
		}

		prefixValue := getTagValue(tag, "envPrefix")
		if prefixValue != "" {
			if prefix != "" {
				prefixValue = prefix + prefixValue
			}

			switch fieldType := field.Type.(type) {
			case *ast.Ident:
				for key, nestedStruct := range structs {
					keyParts := strings.Split(key, ".")
					if len(keyParts) == 2 && keyParts[1] == fieldType.Name {
						// Match structs from the same package or external packages
						if keyParts[0] == currentPackage || currentPackage == "main" {
							envVars = append(envVars, extractEnvVars(nestedStruct, structs, keyParts[0], prefixValue)...)
						}
					}
				}
			case *ast.SelectorExpr:
				// Resolve the package and type from the SelectorExpr
				if pkgIdent, isIdentifier := fieldType.X.(*ast.Ident); isIdentifier {
					pkgName := pkgIdent.Name
					typeName := fieldType.Sel.Name

					// Combine package and type to match the key in the structs map
					fullName := pkgName + "." + typeName
					if nestedStruct, found := structs[fullName]; found {
						envVars = append(envVars, extractEnvVars(nestedStruct, structs, pkgName, prefixValue)...)
					} else {
						println("")
					}
				}
			}
		}
	}

	return envVars
}

// getTagValue extracts the value of a specific tag from a struct field tag.
func getTagValue(tag, key string) string {
	tags := strings.Split(tag, " ")
	for _, t := range tags {
		parts := strings.SplitN(t, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			return strings.Trim(parts[1], "\"")
		}
	}
	return ""
}
