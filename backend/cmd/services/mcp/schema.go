package main

import (
	"strings"

	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/invopop/jsonschema"
)

const (
	jsonSchemaVersion = "https://json-schema.org/draft/2020-12/schema"
	objType           = "object"
	arrType           = "array"
	strType           = "string"
	boolType          = "boolean"
	intType           = "integer"
	dtFmt             = "date-time"
)

// queryFilterSchema returns the JSON schema for a QueryFilter object
func queryFilterSchema() map[string]any {
	return map[string]any{
		"type": objType,
		"properties": map[string]any{
			"SortBy": map[string]any{
				"type":        strType,
				"description": "Field to sort by",
			},
			"CreatedAfter": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results created after this timestamp (ISO 8601)",
			},
			"CreatedBefore": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results created before this timestamp (ISO 8601)",
			},
			"UpdatedAfter": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results updated after this timestamp (ISO 8601)",
			},
			"UpdatedBefore": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results updated before this timestamp (ISO 8601)",
			},
			"PageSize": map[string]any{
				"type":        intType,
				"description": "Maximum number of results to return",
			},
			"IncludeArchived": map[string]any{
				"type":        boolType,
				"description": "Whether to include archived items",
			},
			"Cursor": map[string]any{
				"type":        strType,
				"description": "Pagination cursor for fetching next page",
			},
		},
	}
}

func schemaForType(x any) map[string]any {
	var y map[string]any
	encoding.MustDecodeJSON(encoding.MustEncodeJSON(jsonschema.Reflect(x)), &y)

	// Transform the schema to match the expected format
	result := transformSchema(y)
	return result
}

func transformSchema(schema map[string]any) map[string]any {
	// Keep $defs for reference resolution
	defs, hasDefs := schema["$defs"].(map[string]any)

	// If schema has $ref, resolve it from $defs
	resolvedSchema := schema
	if ref, ok := schema["$ref"].(string); ok && hasDefs {
		if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
			if def, ok := defs[refName].(map[string]any); ok {
				resolvedSchema = def
			}
		}
	}

	// Build the result with $schema and type
	result := map[string]any{
		"$schema": jsonSchemaVersion,
		"type":    objType,
	}

	// Extract properties
	props, ok := resolvedSchema["properties"].(map[string]any)
	if !ok {
		return result
	}

	// Special case: if there's only a "Results" or "results" property that's an array, transform to items format
	if len(props) == 1 {
		var resultsProp map[string]any
		var found bool
		// Check both camelCase and PascalCase
		if resultsProp, found = props["results"].(map[string]any); !found {
			resultsProp, found = props["Results"].(map[string]any)
		}
		if found {
			if propType, ok := resultsProp["type"].(string); ok && propType == "array" {
				// This is the SearchValidIngredientsResult case
				itemsSchema := transformResultsArraySchema(resultsProp, defs, hasDefs)
				result["items"] = itemsSchema
				return result
			}
		}
	}

	// Transform properties
	transformedProps := make(map[string]any)
	for key, value := range props {
		prop, ok := value.(map[string]any)
		if !ok {
			transformedProps[key] = value
			continue
		}

		// Handle $ref in properties (for nested types like QueryFilter)
		if ref, ok := prop["$ref"].(string); ok && hasDefs {
			if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
				if def, ok := defs[refName].(map[string]any); ok {
					prop = transformQueryFilterSchema(def)
				}
			}
		}

		// Convert camelCase to PascalCase for field names
		pascalKey := toPascalCase(key)

		// Add description if not present
		propCopy := make(map[string]any)
		for k, v := range prop {
			propCopy[k] = v
		}
		if _, hasDesc := propCopy["description"]; !hasDesc {
			if desc := getFieldDescription(pascalKey); desc != "" {
				propCopy["description"] = desc
			}
		}

		transformedProps[pascalKey] = propCopy
	}

	result["properties"] = transformedProps
	return result
}

func transformQueryFilterSchema(schema map[string]any) map[string]any {
	result := map[string]any{
		"type": objType,
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		return result
	}

	transformedProps := make(map[string]any)
	for key, value := range props {
		prop, ok := value.(map[string]any)
		if !ok {
			continue
		}

		// Convert camelCase to PascalCase
		pascalKey := toPascalCase(key)

		// Special case: limit -> PageSize
		if key == "limit" {
			pascalKey = "PageSize"
		}

		// Add descriptions based on field name
		propCopy := make(map[string]any)
		for k, v := range prop {
			propCopy[k] = v
		}

		// Add description if not present
		if _, hasDesc := propCopy["description"]; !hasDesc {
			propCopy["description"] = getFieldDescription(pascalKey)
		}

		transformedProps[pascalKey] = propCopy
	}

	result["properties"] = transformedProps
	return result
}

func transformResultsArraySchema(arrayProp map[string]any, defs map[string]any, hasDefs bool) map[string]any {
	result := map[string]any{
		"name": "Results",
		"type": arrType,
	}

	// Get the items schema
	items, ok := arrayProp["items"].(map[string]any)
	if !ok {
		return result
	}

	// Resolve $ref if present
	var itemSchema map[string]any
	if ref, ok := items["$ref"].(string); ok && hasDefs {
		if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
			if def, ok := defs[refName].(map[string]any); ok {
				itemSchema = def
			}
		}
	} else {
		itemSchema = items
	}

	if itemSchema == nil {
		return result
	}

	// Extract properties from the item schema
	props, ok := itemSchema["properties"].(map[string]any)
	if !ok {
		return result
	}

	// Transform properties (keep camelCase as expected)
	transformedProps := make(map[string]any)
	for key, value := range props {
		prop, ok := value.(map[string]any)
		if !ok {
			transformedProps[key] = value
			continue
		}

		propCopy := make(map[string]any)
		for k, v := range prop {
			propCopy[k] = v
		}

		// Handle $ref in properties (for nested types like OptionalFloat32Range)
		if ref, ok := propCopy["$ref"].(string); ok && hasDefs {
			if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
				if def, ok := defs[refName].(map[string]any); ok {
					// Replace $ref with the actual schema
					for k, v := range def {
						propCopy[k] = v
					}
					delete(propCopy, "$ref")
				}
			}
		}

		// Handle nullable types - check if the field is a pointer type
		// For nullable fields, type should be ["string", "null"] or ["number", "null"]
		if propType, ok := propCopy["type"].(string); ok {
			// Check if this is a nullable field based on the expected schema
			nullableFields := map[string]bool{
				"archivedAt":    true,
				"lastUpdatedAt": true,
			}
			if nullableFields[key] {
				propCopy["type"] = []any{propType, "null"}
			}
		}

		// Handle nested objects with nullable number fields
		if propType, ok := propCopy["type"].(string); ok && propType == objType {
			// Remove additionalProperties if present
			delete(propCopy, "additionalProperties")
			if nestedProps, ok := propCopy["properties"].(map[string]any); ok {
				transformedNestedProps := make(map[string]any)
				for nestedKey, nestedValue := range nestedProps {
					nestedProp, ok := nestedValue.(map[string]any)
					if !ok {
						transformedNestedProps[nestedKey] = nestedValue
						continue
					}
					nestedPropCopy := make(map[string]any)
					for k, v := range nestedProp {
						nestedPropCopy[k] = v
					}
					// Make min and max nullable
					if nestedKey == "min" || nestedKey == "max" {
						if nestedType, ok := nestedPropCopy["type"].(string); ok && nestedType == "number" {
							nestedPropCopy["type"] = []any{"number", "null"}
						}
					}
					transformedNestedProps[nestedKey] = nestedPropCopy
				}
				propCopy["properties"] = transformedNestedProps
			}
		}

		transformedProps[key] = propCopy
	}

	result["properties"] = transformedProps
	return result
}

func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getFieldDescription(fieldName string) string {
	descriptions := map[string]string{
		"SortBy":           "Field to sort by",
		"CreatedAfter":     "Filter results created after this timestamp (ISO 8601)",
		"CreatedBefore":    "Filter results created before this timestamp (ISO 8601)",
		"UpdatedAfter":     "Filter results updated after this timestamp (ISO 8601)",
		"UpdatedBefore":    "Filter results updated before this timestamp (ISO 8601)",
		"PageSize":         "Maximum number of results to return",
		"IncludeArchived":  "Whether to include archived items",
		"Cursor":           "Pagination cursor for fetching next page",
		"Query":            "Search query string to match ingredient names or descriptions",
		"UseSearchService": "Whether to use the search service for more advanced search capabilities",
	}
	if desc, ok := descriptions[fieldName]; ok {
		return desc
	}
	return ""
}
