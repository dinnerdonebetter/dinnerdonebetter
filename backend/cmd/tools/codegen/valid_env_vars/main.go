package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/codemodus/kace"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	structs := parseGoFiles(dir)

	outputLines := []string{}
	if mainAST, found := structs["config.APIServiceConfig"]; found {
		for envVar, fieldPath := range extractEnvVars(mainAST, structs, "main", "", "") {
			outputLines = append(outputLines, fmt.Sprintf(`	// %sEnvVarKey is the environment variable name to set in order to override `+"`"+`config%s`+"`"+`.
	%sEnvVarKey = "%s%s"

`, kace.Pascal(envVar), fieldPath, kace.Pascal(envVar), config.EnvVarPrefix, envVar))
		}
	}

	slices.Sort(outputLines)

	out := fmt.Sprintf(`package envvars

/* 
This file contains a reference of all valid service environment variables.
*/

const (
%s
)
`, strings.Join(outputLines, ""))

	if err = os.WriteFile(filepath.Join(dir, "internal", "config", "envvars", "env_vars.go"), []byte(out), 0o0600); err != nil {
		log.Fatal(err)
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

		node, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", path, err)
			return nil
		}

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

				if _, ok = typeSpec.Type.(*ast.StructType); ok {
					key := fmt.Sprintf("%s.%s", node.Name.Name, typeSpec.Name.Name)
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

// handleIdent handles extracting info from an *ast.Ident node.
func handleIdent(structs map[string]*ast.TypeSpec, fieldType *ast.Ident, envVars map[string]string, currentPackage, prefixValue, fieldNamePrefix, fieldName string) {
	for key, nestedStruct := range structs {
		keyParts := strings.Split(key, ".")
		if len(keyParts) == 2 && keyParts[1] == fieldType.Name {
			if keyParts[0] == currentPackage || currentPackage == "main" {
				for k, v := range extractEnvVars(nestedStruct, structs, keyParts[0], prefixValue, fmt.Sprintf("%s.%s", fieldNamePrefix, fieldName)) {
					envVars[k] = v
				}
			}
		}
	}
}

// handleSelectorExpr handles extracting info from an *ast.SelectorExpr node.
func handleSelectorExpr(structs map[string]*ast.TypeSpec, fieldType *ast.SelectorExpr, envVars map[string]string, prefixValue, fieldNamePrefix, fieldName string) {
	if pkgIdent, isIdentifier := fieldType.X.(*ast.Ident); isIdentifier {
		pkgName := pkgIdent.Name

		fullName := fmt.Sprintf("%s.%s", pkgName, fieldType.Sel.Name)
		if nestedStruct, found := structs[fullName]; found {
			for k, v := range extractEnvVars(nestedStruct, structs, pkgName, prefixValue, fmt.Sprintf("%s.%s", fieldNamePrefix, fieldName)) {
				envVars[k] = v
			}
		}
	}
}

// extractEnvVars traverses a struct definition and collects environment variables, resolving nested structs.
func extractEnvVars(typeSpec *ast.TypeSpec, structs map[string]*ast.TypeSpec, currentPackage, envVarPrefix, fieldNamePrefix string) map[string]string {
	envVars := map[string]string{}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return envVars
	}

	for _, field := range structType.Fields.List {
		if field.Tag == nil {
			continue
		}

		tag := strings.Trim(field.Tag.Value, "`")
		if tag == `json:"-"` || tag == "" {
			continue
		}

		fn := field.Names[0].Name

		if envValue := getTagValue(tag, "env"); envValue != "" {
			if envVarPrefix != "" {
				envValue = envVarPrefix + envValue
			}

			if fieldNamePrefix == "" {
				envVars[envValue] = fn
			} else {
				envVars[envValue] = fmt.Sprintf("%s.%s", fieldNamePrefix, fn)
			}
		}

		if prefixValue := getTagValue(tag, "envPrefix"); prefixValue != "" {
			if envVarPrefix != "" {
				prefixValue = envVarPrefix + prefixValue
			}

			switch fieldType := field.Type.(type) {
			case *ast.Ident:
				handleIdent(structs, fieldType, envVars, currentPackage, prefixValue, fieldNamePrefix, fn)
			case *ast.SelectorExpr:
				handleSelectorExpr(structs, fieldType, envVars, prefixValue, fieldNamePrefix, fn)
			case *ast.StarExpr:
				switch ft := fieldType.X.(type) {
				case *ast.Ident:
					handleIdent(structs, ft, envVars, currentPackage, prefixValue, fieldNamePrefix, fn)
				case *ast.SelectorExpr:
					handleSelectorExpr(structs, ft, envVars, prefixValue, fieldNamePrefix, fn)
				}
			}
		}
	}

	return envVars
}
