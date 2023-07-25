package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

func noTestFiles(f os.FileInfo) bool {
	return !f.IsDir() && !strings.HasSuffix(f.Name(), "_test.go")
}

func fetchTypesForPackage(pkg, typeName string) *ast.StructType {
	fileset := token.NewFileSet()

	astPkg, err := parser.ParseDir(fileset, pkg, noTestFiles, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse package: %v", err)
	}

	if len(astPkg) != 1 {
		return nil
	}

	foundTypes := map[string]*ast.StructType{}
	for _, v := range astPkg {
		for _, f := range v.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					if y, ok := x.Type.(*ast.StructType); ok {
						foundTypes[x.Name.Name] = y
					}
				}
				return true
			})
		}
	}

	return foundTypes[typeName]
}

type typeSpec struct {
	typeName string
	packages []string
}

func ensureTypesHaveMatchingFields(allStructDefinitions []*ast.StructType) error {
	allStructFields := []map[string]string{}
	for _, strukt := range allStructDefinitions {
		allStructFields = append(allStructFields, getFieldsForStruct(strukt))
	}

	if len(allStructFields) == 0 {
		return fmt.Errorf("no struct fields found")
	}

	for i := 1; i < len(allStructFields); i++ {
		for k, v := range allStructFields[i] {
			if allStructFields[0][k] != v {
				return fmt.Errorf("field %s has different types: %s vs %s", k, allStructFields[0][k], v)
			}
		}
	}

	return nil
}

func evaluateTypeSpec(spec typeSpec) error {
	allStructDefinitions := []*ast.StructType{}
	for _, p := range spec.packages {
		if fetchedType := fetchTypesForPackage(p, spec.typeName); fetchedType != nil {
			allStructDefinitions = append(allStructDefinitions, fetchedType)
		}
	}

	if len(allStructDefinitions) != len(spec.packages) {
		return fmt.Errorf("failed to find all struct definitions for %s", spec.typeName)
	}

	if err := ensureTypesHaveMatchingFields(allStructDefinitions); err != nil {
		return err
	}

	return nil
}

func getFieldsForStruct(structType *ast.StructType) map[string]string {
	structFields := make(map[string]string)

	for _, field := range structType.Fields.List {
		fieldType := ""
		switch t := field.Type.(type) {
		case *ast.Ident:
			fieldType = t.Name
		case *ast.SelectorExpr:
			fieldType = fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
		}

		for _, name := range field.Names {
			if name.Name != "_" {
				structFields[name.Name] = fieldType
			}
		}
	}

	return structFields
}

func main() {
	specs := []typeSpec{
		{
			typeName: "ValidIngredient",
			packages: []string{
				"internal/database/postgres/v2",
				"pkg/types",
			},
		},
		{
			typeName: "ValidInstrument",
			packages: []string{
				"internal/database/postgres/v2",
				"pkg/types",
			},
		},
		{
			typeName: "ValidPreparation",
			packages: []string{
				"internal/database/postgres/v2",
				"pkg/types",
			},
		},
		{
			typeName: "Webhook",
			packages: []string{
				"internal/database/postgres/v2",
				"pkg/types",
			},
		},
		{
			typeName: "WebhookTriggerEvent",
			packages: []string{
				"internal/database/postgres/v2",
				"pkg/types",
			},
		},
	}

	for _, s := range specs {
		if err := evaluateTypeSpec(s); err != nil {
			log.Fatalf("failed to evaluate type spec: %v", err)
		}
	}
}
