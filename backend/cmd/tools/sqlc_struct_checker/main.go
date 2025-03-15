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

// noTestFiles excludes _test.go files (and directories)
func noTestFiles(f os.FileInfo) bool {
	return !f.IsDir() && !strings.HasSuffix(f.Name(), "_test.go")
}

// fetchTypesForPackage parses the given package directory and returns a map of type names to their *ast.StructType.
// The filter function can be used to limit which struct types are included.
func fetchTypesForPackage(pkg string, nameFilter func(string) bool) map[string]*ast.StructType {
	fileset := token.NewFileSet()
	astPkg, err := parser.ParseDir(fileset, pkg, noTestFiles, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse package %s: %v", pkg, err)
	}

	// Require exactly one package directory.
	if len(astPkg) != 1 {
		log.Fatalf("expected one package in %s, found %d", pkg, len(astPkg))
	}

	foundTypes := map[string]*ast.StructType{}
	for _, p := range astPkg {
		for _, f := range p.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if structType, ok := ts.Type.(*ast.StructType); ok && nameFilter(ts.Name.Name) {
						foundTypes[ts.Name.Name] = structType
					}
				}
				return true
			})
		}
	}
	return foundTypes
}

// getOrderedFieldNamesForStruct returns the field names of a struct in order.
// It handles both named fields and anonymous fields (using the type name).
func getOrderedFieldNamesForStruct(structType *ast.StructType) []string {
	var fields []string
	for _, field := range structType.Fields.List {
		if field.Names != nil {
			for _, name := range field.Names {
				if name.Name != "_" {
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
func extractSelectorFromExpr(expr ast.Expr) (string, string, bool) {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name, t.Sel.Name, true
		}
	case *ast.StarExpr:
		if sel, ok := t.X.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				return ident.Name, sel.Sel.Name, true
			}
		}
	}
	return "", "", false
}

// checkCompositeLit inspects a composite literal and checks that all fields declared in the struct definition are initialized.
// For key–value initialization, it verifies each field is provided.
// For positional initialization, it checks that the number of elements matches the number of fields.
func checkCompositeLit(lit *ast.CompositeLit, structDef *ast.StructType, fileName string) error {
	expectedFields := getOrderedFieldNamesForStruct(structDef)

	// No elements provided: error for each expected field.
	if len(lit.Elts) == 0 {
		var errs *multierror.Error
		for _, field := range expectedFields {
			errs = multierror.Append(errs, fmt.Errorf("in %s: composite literal of type %s missing field %s", fileName, getTypeAsString(lit.Type), field))
		}
		return errs.ErrorOrNil()
	}

	// Determine if the literal uses key–value initialization.
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
				if ident, ok := kv.Key.(*ast.Ident); ok {
					initialized[ident.Name] = true
				}
			}
		}
		var errs *multierror.Error
		for _, field := range expectedFields {
			if !initialized[field] {
				errs = multierror.Append(errs, fmt.Errorf("in %s: composite literal of type %s missing field %s", fileName, getTypeAsString(lit.Type), field))
			}
		}
		return errs.ErrorOrNil()
	} else {
		// Positional initialization: ensure the number of elements matches the number of expected fields.
		if len(lit.Elts) != len(expectedFields) {
			return fmt.Errorf("in %s: positional composite literal of type %s expects %d fields but got %d", fileName, getTypeAsString(lit.Type), len(expectedFields), len(lit.Elts))
		}
		// If counts match, assume the ordering is correct.
		return nil
	}
}

func evaluatePackageUsage(sourcePackage, implementingPackage, mustContain string) error {
	filterFunc := func(s string) bool {
		if mustContain == "" {
			return true
		}
		return strings.Contains(s, mustContain)
	}

	// Fetch all struct types from the source package.
	sourceTypes := fetchTypesForPackage(sourcePackage, filterFunc)
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

				// Check if the literal’s type is a selector (possibly wrapped in a pointer)
				_, typeName, found := extractSelectorFromExpr(compLit.Type)
				if !found {
					// Not a composite literal with a qualified type – skip.
					return true
				}

				// Only check composite literals whose type (by its unqualified name) comes from the source package.
				structDef, exists := sourceTypes[typeName]
				if !exists {
					return true
				}

				// Check that all fields are initialized.
				if err = checkCompositeLit(compLit, structDef, fileName); err != nil {
					errors = multierror.Append(errors, err)
				}

				return true
			})
		}
	}

	return errors.ErrorOrNil()
}

func main() {
	if err := evaluatePackageUsage("internal/database/postgres/generated", "internal/database/postgres", ""); err != nil {
		log.Fatal(err)
	}
}
