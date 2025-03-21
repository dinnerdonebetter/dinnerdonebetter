package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
)

func noTestFiles(f os.FileInfo) bool {
	return !f.IsDir() && !strings.HasSuffix(f.Name(), "_test.go")
}

// fetchTypesForPackage parses the given package directory and returns a map of type names to their *ast.StructType.
// The filter function can be used to limit which struct types are included.
func fetchTypesForPackage(pkg string) map[string]*ast.StructType {
	fileset := token.NewFileSet()
	astPkg, err := parser.ParseDir(fileset, pkg, noTestFiles, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse package %s: %v", pkg, err)
	}

	if len(astPkg) != 1 {
		log.Fatalf("expected one package in %s, found %d", pkg, len(astPkg))
	}

	foundTypes := map[string]*ast.StructType{}
	for _, p := range astPkg {
		for _, f := range p.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if structType, ok2 := ts.Type.(*ast.StructType); ok2 && ast.IsExported(ts.Name.Name) {
						foundTypes[ts.Name.Name] = structType
					}
				}
				return true
			})
		}
	}

	return foundTypes
}

// getOrderedFieldNamesForStruct returns the field names of a struct in order. It handles both named fields and anonymous fields (using the type name).
func getOrderedFieldNamesForStruct(structType *ast.StructType) []string {
	var fields []string

	for _, field := range structType.Fields.List {
		if field.Names != nil {
			for _, name := range field.Names {
				if name.Name != "_" && ast.IsExported(name.Name) {
					fields = append(fields, name.Name)
				}
			}
		} else {
			// Anonymous field: derive name from the type.
			switch t := field.Type.(type) {
			case *ast.Ident:
				if t.Name != "_" {
					fields = append(fields, t.Name)
				}
			case *ast.SelectorExpr:
				fields = append(fields, t.Sel.Name)
			}
		}
	}

	return fields
}

// getTypeAsString converts an AST expression representing a type to a string.
func getTypeAsString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name)
		}
		return t.Sel.Name
	case *ast.StarExpr:
		return "*" + getTypeAsString(t.X)
	default:
		return ""
	}
}

// extractSelectorFromExpr extracts the package alias and type name from an expression.
// It supports either a SelectorExpr or a pointer (StarExpr) to a SelectorExpr.
func extractSelectorFromExpr(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		if _, ok := t.X.(*ast.Ident); ok {
			return t.Sel.Name, true
		}
	case *ast.StarExpr:
		if sel, ok := t.X.(*ast.SelectorExpr); ok {
			if _, ok2 := sel.X.(*ast.Ident); ok2 {
				return sel.Sel.Name, true
			}
		}
	}

	return "", false
}

// checkCompositeLit inspects a composite literal and checks that all fields declared in the struct definition are initialized.
// For key–value initialization, it verifies each field is provided.
// For positional initialization, it checks that the number of elements matches the number of fields.
func checkCompositeLit(lit *ast.CompositeLit, structDef *ast.StructType, fileName string, fset *token.FileSet) error {
	pos := fset.Position(lit.Pos())
	expectedFields := getOrderedFieldNamesForStruct(structDef)

	if len(lit.Elts) == 0 {
		var errs *multierror.Error
		for _, field := range expectedFields {
			errs = multierror.Append(errs, fmt.Errorf("in %s:%d: composite literal of type %s missing field %s", fileName, pos.Line, getTypeAsString(lit.Type), field))
		}
		return errs.ErrorOrNil()
	}

	allKV := true
	for _, elt := range lit.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); !ok {
			allKV = false
			break
		}
	}

	if allKV {
		initialized := make(map[string]bool)
		for _, elt := range lit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if ident, ok2 := kv.Key.(*ast.Ident); ok2 {
					initialized[ident.Name] = true
				}
			}
		}
		var errs *multierror.Error
		for _, field := range expectedFields {
			if !initialized[field] {
				errs = multierror.Append(errs, fmt.Errorf("in %s:%d: composite literal of type %s missing field %s", fileName, pos.Line, getTypeAsString(lit.Type), field))
			}
		}
		return errs.ErrorOrNil()
	} else {
		if len(lit.Elts) != len(expectedFields) {
			return fmt.Errorf("in %s:%d: positional composite literal of type %s expects %d fields but got %d", fileName, pos.Line, getTypeAsString(lit.Type), len(expectedFields), len(lit.Elts))
		}

		return nil
	}
}

func evaluatePackageUsage(sourcePackage, implementingPackage string) error {
	// Fetch all struct types from the source package.
	sourceTypes := fetchTypesForPackage(sourcePackage)
	if len(sourceTypes) == 0 {
		log.Fatalf("No types found in source package %s", sourcePackage)
	}

	// Parse the target package.
	fileset := token.NewFileSet()
	targetPkg, err := parser.ParseDir(fileset, implementingPackage, noTestFiles, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse target package %s: %v", implementingPackage, err)
	}

	var errors *multierror.Error

	// Walk through every file in the target package and inspect composite literals.
	for _, pkg := range targetPkg {
		for fileName, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				compLit, ok := n.(*ast.CompositeLit)
				if !ok {
					return true
				}

				typeName, found := extractSelectorFromExpr(compLit.Type)
				if !found {
					return true
				}

				structDef, exists := sourceTypes[typeName]
				if !exists {
					return true
				}

				if err = checkCompositeLit(compLit, structDef, fileName, fileset); err != nil {
					errors = multierror.Append(errors, err)
				}

				return true
			})
		}
	}

	return errors.ErrorOrNil()
}

func main() {
	packagesToCheck := map[string]string{
		"internal/database/postgres/generated": "internal/database/postgres",
		"internal/services/eating/types":       "internal/services/eating/grpc/converters",
		"internal/grpc/messages":               "internal/services/eating/grpc/converters",
	}

	for sourcePackage, implementingPackage := range packagesToCheck {
		if err := evaluatePackageUsage(sourcePackage, implementingPackage); err != nil {
			log.Fatal(err)
		}
	}
}
