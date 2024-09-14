package main

import (
	"reflect"
)

// OpenAPISchema represents an OpenAPI 3.1 schema model
type OpenAPISchema struct {
	Type       string                    `json:"type,omitempty"`
	Properties map[string]*OpenAPISchema `json:"properties,omitempty"`
	Required   []string                  `json:"required,omitempty"`
}

// RenderSchema generates a YAML representation of the schema
func (s *OpenAPISchema) RenderSchema() map[string]any {
	propertiesMap := map[string]any{}
	output := map[string]any{
		"type": "object",
	}

	for name, prop := range s.Properties {
		propertiesMap[name] = map[string]any{
			"type": prop.Type,
		}
	}

	output["properties"] = propertiesMap

	return output
}

// FieldSchema holds metadata about a field's type
func FieldSchema(t reflect.Type) *OpenAPISchema {
	schema := &OpenAPISchema{}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		schema.Type = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		schema.Type = "integer"
	case reflect.Float32, reflect.Float64:
		schema.Type = "number"
	case reflect.Bool:
		schema.Type = "boolean"
	case reflect.Slice, reflect.Array:
		schema.Type = "array"
		schema.Properties = map[string]*OpenAPISchema{
			"items": FieldSchema(t.Elem()),
		}
	case reflect.Map:
		schema.Type = "object"
	case reflect.Struct:
		schema.Type = "object"
		// TODO: this should probably be a ref or something
		schema.Properties = make(map[string]*OpenAPISchema)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == "_" {
				continue
			}

			if value := field.Tag.Get("json"); value == "" || value == "-" {
				continue
			}

			fieldSchema := FieldSchema(field.Type)
			schema.Properties[field.Name] = fieldSchema
			if !field.IsExported() {
				continue
			}
		}
	case reflect.Interface:
		panic("unhandled interface")
	default:
		panic("unhandled default case")
	}
	return schema
}

// SchemaFromInstance generates an OpenAPI schema for a given Go type instance
func SchemaFromInstance(instance interface{}) *OpenAPISchema {
	return FieldSchema(reflect.TypeOf(instance))
}
