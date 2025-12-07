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
	numberType        = "number"

	dtFmt = "date-time"
)

// queryFilterSchema returns the JSON schema for a QueryFilter object
func queryFilterSchema() map[string]any {
	return objectType(map[string]any{
		"SortBy":          stringField("Field to sort by"),
		"CreatedAfter":    timestampField("Filter results created after this timestamp (ISO 8601)"),
		"CreatedBefore":   timestampField("Filter results created before this timestamp (ISO 8601)"),
		"UpdatedAfter":    timestampField("Filter results updated after this timestamp (ISO 8601)"),
		"UpdatedBefore":   timestampField("Filter results updated before this timestamp (ISO 8601)"),
		"MaxResponseSize": intField("Maximum number of results to return"),
		"IncludeArchived": boolField("Whether to include archived items"),
		"Cursor":          stringField("Pagination cursor for fetching next page"),
	})
}

func optionalFloatRangeSchema() map[string]any {
	return objectType(map[string]any{
		"Min": floatField("Minimum value"),
		"Max": floatField("Maximum value"),
	})
}

func schemaObject(properties map[string]any) map[string]any {
	return map[string]any{
		"$schema":    jsonSchemaVersion,
		"type":       objType,
		"properties": properties,
	}
}

func objectType(fieldSchema map[string]any, requiredFields ...string) map[string]any {
	x := map[string]any{
		"type":       objType,
		"properties": fieldSchema,
	}

	if len(requiredFields) > 0 {
		x["required"] = requiredFields
	}

	return x
}

func floatField(description string) map[string]any {
	return map[string]any{
		"type":        numberType,
		"description": description,
	}
}

func uintField(description string) map[string]any {
	return map[string]any{
		"type":        intType,
		"description": description,
		"minimum":     0,
	}
}

func intField(description string) map[string]any {
	return map[string]any{
		"type":        intType,
		"description": description,
	}
}

func boolField(description string) map[string]any {
	return map[string]any{
		"type":        boolType,
		"description": description,
	}
}

func stringField(description string) map[string]any {
	x := map[string]any{
		"type":        strType,
		"description": description,
	}

	return x
}

func timestampField(description string) map[string]any {
	return stringFieldWithFormat(description, dtFmt)
}

func stringFieldWithFormat(description, format string) map[string]any {
	x := map[string]any{
		"type":        strType,
		"description": description,
		"format":      format,
	}

	return x
}

func schemaForType(x any) map[string]any {
	var y map[string]any
	schema := jsonschema.Reflect(x)
	encoded := encoding.MustEncodeJSON(schema)
	encoding.MustDecodeJSON(encoded, &y)

	// Transform the schema to match the expected format
	result := transformSchema(y)
	return result
}

// extractRefFromPropertyWithNullable extracts a $ref from a property schema and detects if it's nullable.
// It handles:
// - Direct $ref: {"$ref": "..."}
// - allOf with $ref: {"allOf": [{"$ref": "..."}]}
// - oneOf with $ref: {"oneOf": [{"type": "null"}, {"$ref": "..."}]} (returns ref and isNullable=true)
func extractRefFromPropertyWithNullable(prop map[string]any, defs map[string]any, hasDefs bool) (string, bool) {
	// Check for direct $ref
	if ref, ok := prop["$ref"].(string); ok {
		return ref, false
	}

	// Check for allOf with $ref
	if allOf, ok := prop["allOf"].([]any); ok && len(allOf) > 0 {
		for _, item := range allOf {
			if itemMap, ok := item.(map[string]any); ok {
				if ref, ok := itemMap["$ref"].(string); ok {
					return ref, false
				}
			}
		}
	}

	// Check for oneOf with $ref (pointer types often use this)
	// oneOf typically has [{"type": "null"}, {"$ref": "..."}]
	if oneOf, ok := prop["oneOf"].([]any); ok && len(oneOf) > 0 {
		hasNull := false
		var ref string
		for _, item := range oneOf {
			if itemMap, ok := item.(map[string]any); ok {
				// Check if this is the null type
				if itemType, ok := itemMap["type"].(string); ok && itemType == "null" {
					hasNull = true
				}
				// Check if this has a $ref
				if r, ok := itemMap["$ref"].(string); ok {
					ref = r
				}
			}
		}
		if ref != "" {
			return ref, hasNull
		}
	}

	return "", false
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

		// Handle $ref in properties (for nested types like QueryFilter or ValidIngredientCreationRequestInput)
		// This handles direct $ref, nullable pointer types with $ref, and allOf/oneOf with $ref
		ref, isNullable := extractRefFromPropertyWithNullable(prop, defs, hasDefs)
		if ref != "" {
			if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
				if def, ok := defs[refName].(map[string]any); ok {
					var transformed map[string]any
					// Special case for QueryFilter
					if refName == "QueryFilter" || strings.Contains(refName, "QueryFilter") {
						transformed = transformQueryFilterSchema(def)
					} else {
						// For other types, recursively transform the schema
						transformed = transformNestedSchema(def, defs, hasDefs)
					}
					// Handle nullable types (from oneOf with null)
					// Note: For MCP inspector compatibility, we may want to avoid nullable complex objects
					// as they can render as JSON blobs. Only mark as nullable if it's a simple type.
					if isNullable {
						if currentType, ok := transformed["type"].(string); ok {
							// Only make nullable if it's not a complex object (objects should be required)
							// This helps MCP inspector render the form properly
							if currentType != "object" {
								transformed["type"] = []any{currentType, "null"}
							}
							// For objects, we keep them as required (non-nullable) for better MCP inspector UX
						}
					}
					// Completely replace prop with transformed schema to remove any allOf/oneOf/$ref
					// Also ensure no allOf/oneOf/$ref remain, and explicitly set type for objects
					delete(transformed, "allOf")
					delete(transformed, "oneOf")
					delete(transformed, "$ref")
					delete(transformed, "additionalProperties")
					// Ensure type is explicitly set for objects (helps MCP inspector render properly)
					if _, hasProps := transformed["properties"]; hasProps {
						if _, hasType := transformed["type"]; !hasType {
							transformed["type"] = "object"
						}
					}
					prop = transformed
				}
			}
		} else {
			// Even if no $ref, check if prop has allOf/oneOf that should be cleaned up
			if _, hasAllOf := prop["allOf"]; hasAllOf {
				// If we have allOf but couldn't extract a ref, try to flatten it
				if allOf, ok := prop["allOf"].([]any); ok && len(allOf) == 1 {
					if itemMap, ok := allOf[0].(map[string]any); ok {
						// If allOf has only one item, use that item
						prop = itemMap
					}
				}
			}
			if _, hasOneOf := prop["oneOf"]; hasOneOf {
				// If we have oneOf, try to extract the non-null item
				if oneOf, ok := prop["oneOf"].([]any); ok {
					for _, item := range oneOf {
						if itemMap, ok := item.(map[string]any); ok {
							// Skip null type, use the other one
							if itemType, ok := itemMap["type"].(string); ok && itemType != "null" {
								prop = itemMap
								break
							} else if _, hasRef := itemMap["$ref"]; hasRef {
								// If it has a $ref, we should have caught it above, but try again
								if ref, ok := itemMap["$ref"].(string); ok && hasDefs {
									if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
										if def, ok := defs[refName].(map[string]any); ok {
											transformed := transformNestedSchema(def, defs, hasDefs)
											prop = transformed
											break
										}
									}
								}
							}
						}
					}
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

	// For UpdateValidIngredientInvocation and similar tools, ensure complex Input fields are required
	// This helps MCP inspector render them as forms instead of JSON blobs
	if props, ok := result["properties"].(map[string]any); ok {
		if inputProp, ok := props["Input"].(map[string]any); ok {
			// If Input is a complex object with properties, ensure it's explicitly typed as object
			// and doesn't have any nullable/optional markers that confuse MCP inspector
			if inputProps, ok := inputProp["properties"].(map[string]any); ok && len(inputProps) > 0 {
				// Ensure type is explicitly "object" (not nullable)
				// Remove any nullable type markers
				if inputType, ok := inputProp["type"]; ok {
					if typeArray, ok := inputType.([]any); ok {
						// If it's an array like ["object", "null"], extract just "object"
						for _, t := range typeArray {
							if tStr, ok := t.(string); ok && tStr == "object" {
								inputProp["type"] = "object"
								break
							}
						}
					} else if inputTypeStr, ok := inputType.(string); ok && inputTypeStr == "object" {
						// Already correct
						inputProp["type"] = "object"
					}
				} else {
					// Type not set, set it explicitly
					inputProp["type"] = "object"
				}
				// Remove any nullable markers
				delete(inputProp, "nullable")

				// Mark Input as required to force MCP inspector to render it as a form field
				// Some inspectors only render required complex objects as forms
				if required, ok := result["required"].([]string); ok {
					// Check if "Input" is already in required
					found := false
					for _, r := range required {
						if r == "Input" {
							found = true
							break
						}
					}
					if !found {
						result["required"] = append(required, "Input")
					}
				} else {
					// Create required array with Input
					result["required"] = []string{"Input"}
				}
			}
		}
	}

	return result
}

func transformNestedSchema(schema map[string]any, defs map[string]any, hasDefs bool) map[string]any {
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
			transformedProps[key] = value
			continue
		}

		propCopy := make(map[string]any)
		for k, v := range prop {
			propCopy[k] = v
		}

		// Handle $ref in nested properties (e.g., OptionalFloat32Range)
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

		// Handle nested objects with nullable number fields (e.g., OptionalFloat32Range)
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
					// Make min and max nullable for OptionalFloat32Range
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

		// Special case: limit -> MaxResponseSize
		if key == "limit" {
			pascalKey = "MaxResponseSize"
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
		"MaxResponseSize":  "Maximum number of results to return",
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
