package main

import (
	"fmt"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// OpenAPISchema represents an OpenAPI 3.1 schema model
type OpenAPISchema struct {
	Type       string                    `json:"type,omitempty"`
	Properties map[string]*OpenAPISchema `json:"properties,omitempty"`
	Required   []string                  `json:"required,omitempty"`
}

// RenderYAML generates a YAML representation of the schema
func (s *OpenAPISchema) RenderYAML() (string, error) {
	yamlBytes, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

// Render generates a YAML representation of the schema as a string
func (s *OpenAPISchema) Render(indentation int) string {
	var sb strings.Builder
	indent := strings.Repeat("  ", indentation)

	// Write the type
	if s.Type != "" {
		sb.WriteString(fmt.Sprintf("%stype: %s\n", indent, s.Type))
	}

	// Write the properties if present
	if len(s.Properties) > 0 {
		sb.WriteString(fmt.Sprintf("%sproperties:\n", indent))
		for propName, propSchema := range s.Properties {
			sb.WriteString(fmt.Sprintf("%s  %s:\n", indent, propName))
			sb.WriteString(propSchema.Render(indentation + 2))
		}
	}

	// Write the required fields if any
	if len(s.Required) > 0 {
		sb.WriteString(fmt.Sprintf("%srequired:\n", indent))
		for _, reqField := range s.Required {
			sb.WriteString(fmt.Sprintf("%s  - %s\n", indent, reqField))
		}
	}

	return sb.String()
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
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
	case reflect.Struct:
		schema.Type = "object"
		schema.Properties = make(map[string]*OpenAPISchema)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldSchema := FieldSchema(field.Type)
			schema.Properties[field.Name] = fieldSchema
			if !field.IsExported() {
				continue
			}
		}
	case reflect.Interface:
		println("")
	default:
		panic("unhandled default case")
	}
	return schema
}

// SchemaFromInstance generates an OpenAPI schema for a given Go type instance
func SchemaFromInstance(instance interface{}) *OpenAPISchema {
	if getTypeName(instance) == getTypeName(&types.EmailAddressVerificationRequestInput{}) {
		println("")
	}

	return FieldSchema(reflect.TypeOf(instance))
}
