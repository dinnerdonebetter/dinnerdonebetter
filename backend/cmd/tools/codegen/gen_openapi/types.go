package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
)

var nativeTypesMap = map[string]struct{}{
	"int":     {},
	"int8":    {},
	"int16":   {},
	"int32":   {},
	"int64":   {},
	"uint":    {},
	"uint8":   {},
	"uint16":  {},
	"uint32":  {},
	"uint64":  {},
	"float32": {},
	"float64": {},
	"string":  {},
	"bool":    {},
	// these are actually openapi types
	"object":  {},
	"integer": {},
	"number":  {},
	"boolean": {},
	"array":   {},
}

var skipTypes = map[string]bool{
	"QueryFilteredResult":                  true,
	"SessionContextData":                   true,
	"RequesterInfo":                        true,
	"DataChangeMessage":                    true,
	"stringDurationValidator":              true,
	"WebhookExecutionRequest":              true,
	"QueryFilter":                          true,
	"MealPlanTaskDatabaseCreationEstimate": true,
	"FinalizedMealPlanDatabaseResult":      true,
	"MissingVote":                          true,
	"MealUpdateRequestInput":               true,
	"OAuth2ClientToken":                    true,
	"MealComponentUpdateRequestInput":      true,
	"RecipeMediaCreationRequestInput":      true,
	"RecipeMediaUpdateRequestInput":        true,
	// one day...
	"NamedID":    true,
	"FloatRange": true,
	"UintRange":  true,
}

type openapiProperty struct {
	Items    *openapiProperty `json:"items,omitempty"    yaml:"items,omitempty"`
	Type     string           `json:"type,omitempty"     yaml:"type,omitempty"`
	Ref      string           `json:"$ref,omitempty"     yaml:"$ref,omitempty"`
	Format   string           `json:"format,omitempty"   yaml:"format,omitempty"`
	Examples []string         `json:"examples,omitempty" yaml:"examples,omitempty"`
}

type openapiSchema struct {
	Properties map[string]*openapiProperty `json:"properties" yaml:"properties"`
	name       string
	Type       string `json:"type" yaml:"type"`
}

func getJSONTagForField(field *ast.Field) string {
	tag := field.Tag.Value

	tags, err := structtag.Parse(strings.TrimPrefix(strings.TrimSuffix(tag, "`"), "`"))
	if err != nil {
		return ""
	}
	jsonTag, err := tags.Get("json")
	if err != nil {
		return ""
	}

	return jsonTag.Name
}

func parseTypes(pkgDir string) ([]*openapiSchema, error) {
	fileset := token.NewFileSet()

	astPkg, err := parser.ParseDir(fileset, pkgDir, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("parsing package directory: %w", err)
	}

	if len(astPkg) == 0 || astPkg == nil {
		return nil, errors.New("no go files found")
	}

	declaredStructs := []*openapiSchema{
		{
			name: "APIResponseWithError",
			Type: "object",
			Properties: map[string]*openapiProperty{
				"details": {
					Ref: "#/components/schemas/ResponseDetails",
				},
				"error": {
					Ref: "#/components/schemas/APIError",
				},
			},
		},
	}

	for _, file := range astPkg {
		// Traverse the AST
		ast.Inspect(file, func(n ast.Node) bool {
			// Look for type declarations (i.e., structs)
			genDecl, ok := n.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				return true
			}

			// Process each type spec (we're interested in struct types)
			for _, spec := range genDecl.Specs {
				typeSpec, ok1 := spec.(*ast.TypeSpec)
				if !ok1 {
					continue
				}

				typeName := typeSpec.Name.Name
				if _, ok2 := skipTypes[typeName]; ok2 {
					continue
				}

				if strings.Contains(typeName, "DatabaseCreationInput") ||
					strings.Contains(typeName, "DatabaseUpdateInput") ||
					strings.Contains(typeName, "SearchSubset") ||
					strings.Contains(typeName, "Mock") ||
					strings.Contains(typeName, "Nullable") {
					continue
				}

				// Check if it's a struct type
				structType, ok1 := typeSpec.Type.(*ast.StructType)
				if !ok1 {
					continue
				}

				schema := &openapiSchema{
					name:       typeName,
					Type:       "object",
					Properties: map[string]*openapiProperty{},
				}
				for _, field := range structType.Fields.List {
					fieldName := field.Names[0].Name
					if fieldName == "_" {
						continue
					}

					if field.Tag != nil {
						if name := getJSONTagForField(field); name != "" {
							if name == "-" {
								continue
							}
							fieldName = name
						}
					}

					if fieldName == "data" && typeName == "APIResponse" {
						continue
					}

					fieldType, isTime, isArray := deriveOpenAPIFieldType(typeName, fieldName, field)
					property := &openapiProperty{
						Type: fieldType,
					}

					if _, nativeType := nativeTypesMap[fieldType]; !nativeType {
						property.Type = ""
						property.Ref = fmt.Sprintf("#/components/schemas/%s", fieldType)
					}

					if isTime {
						property.Format = "date-time"
					}

					if isArray {
						property.Type = "array"
						property.Ref = ""
						property.Items = &openapiProperty{
							Type: fieldType,
						}
						if _, nativeType := nativeTypesMap[fieldType]; !nativeType {
							property.Items.Type = ""
							property.Items = &openapiProperty{Ref: fmt.Sprintf("#/components/schemas/%s", fieldType)}
						}
					}

					schema.Properties[fieldName] = property
				}

				declaredStructs = append(declaredStructs, schema)
			}

			return true
		})
	}

	return declaredStructs, nil
}

var typeAliases = map[string]string{
	"time.Time":              "string",
	"time.Duration":          "string",
	"AuditLogEntryEventType": "string",
	"ErrorCode":              "string",
}

var openAPITypeMap = map[string]string{
	"int":     "integer",
	"int8":    "integer",
	"int16":   "integer",
	"int32":   "integer",
	"int64":   "integer",
	"uint":    "integer",
	"uint8":   "integer",
	"uint16":  "integer",
	"uint32":  "integer",
	"uint64":  "integer",
	"float32": "number",
	"float64": "number",
	"string":  "string",
	"bool":    "boolean",
	// Add other mappings as needed
}

func deriveOpenAPIFieldType(typeName, fieldName string, field *ast.Field) (value string, isTime, isArray bool) {
	if typeName == "WebhookCreationRequestInput" && fieldName == "events" {
		println("here")
	}

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok && x.Obj == nil {
			value = fmt.Sprintf("%s.%s", x.Name, t.Sel.Name)
		}
	case *ast.StarExpr:
		switch u := t.X.(type) {
		case *ast.SelectorExpr:
			if x, ok := u.X.(*ast.Ident); ok {
				value = fmt.Sprintf("%s.%s", x.Name, u.Sel.Name)
			}
		case *ast.Ident:
			value = u.Name
		}
	case *ast.Ident:
		value = t.Name
	case *ast.ArrayType:
		isArray = true
		switch u := t.Elt.(type) {
		case *ast.Ident:
			value = u.Name
		case *ast.StarExpr:
			switch v := u.X.(type) {
			case *ast.SelectorExpr:
				if x, ok := u.X.(*ast.Ident); ok {
					value = fmt.Sprintf("%s.%s", x.Name, v.Sel.Name)
				}
			case *ast.Ident:
				value = v.Name
			}
		}
	case *ast.MapType:
		// TODO: confirm this is being handled correctly
		if typeName == "AuditLogEntry" && fieldName == "changes" {
			return "ChangeLog", isTime, isArray
		}
		value = "object"
	default:
		panic("unhandled case")
	}

	if value == "" {
		panic("empty value for field")
	}

	if value == "time.Time" {
		isTime = true
	}

	if x, ok := typeAliases[value]; ok {
		return x, isTime, isArray
	}

	if x, ok := openAPITypeMap[value]; ok {
		return x, isTime, isArray
	}

	return value, isTime, isArray
}
