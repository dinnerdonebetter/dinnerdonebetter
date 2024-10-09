package typescript

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen"

	"github.com/swaggest/openapi-go/openapi31"
)

const (
	componentSchemaPrefix = "#/components/schemas/"

	typeNameNumberRange                = "NumberRange"
	typeNameNumberRangeWithOptionalMax = "NumberRangeWithOptionalMax"
	typeNameOptionalNumberRange        = "OptionalNumberRange"
)

type Field struct {
	Name,
	Type,
	DefaultValue string
	Nullable,
	Enum,
	Array bool
}

type TypeDefinition struct {
	Imports map[string][]string
	Name    string
	Fields  []Field
}

var defaultTypeValues = map[string]string{
	"string":  `''`,
	"number":  `0`,
	"boolean": `false`,
	"integer": `0`,
}

func defaultValueForType(typeName string) string {
	if x, ok := defaultTypeValues[typeName]; ok {
		return x
	}

	switch typeName {
	case "":
		println("whoa")
		return ""
	default:
		return fmt.Sprintf("new %s()", typeName)
	}
}

var typeReplacementMap = map[string]string{
	"integer": "number",
}

func GenerateModelFiles(spec *openapi31.Spec) (map[string]*TypeDefinition, error) {
	output := map[string]*TypeDefinition{}

	enums, err := GenerateEnumDefinitions(spec)
	if err != nil {
		return nil, fmt.Errorf("could not generate enums file: %w", err)
	}
	enumNames := []string{}
	for _, enum := range enums {
		enumNames = append(enumNames, enum.Name)
	}

	for name, component := range spec.Components.Schemas {
		def := &TypeDefinition{
			Name:    name,
			Imports: map[string][]string{},
		}

		// we'll handle these later
		if _, isEnum := component["enum"]; isEnum {
			continue
		}

		if properties, ok := component["properties"]; ok {
			if propMap, ok2 := properties.(map[string]any); ok2 {
				for k, v := range propMap {
					field := Field{
						Name: k,
					}

					if typeMap, ok3 := v.(map[string]any); ok3 {
						if typ, ok5 := typeMap["type"]; ok5 {
							field.Type = typ.(string)

							if field.Type == "array" {
								field.Array = true
								if items, ok6 := typeMap["items"]; ok6 {
									if itemsMap, ok7 := items.(map[string]any); ok7 {
										if typeData, ok8 := itemsMap["type"]; ok8 {
											if typeStr, ok9 := typeData.(string); ok9 {
												field.Type = typeStr
											}
										} else if refData, ok9 := itemsMap["$ref"]; ok9 {
											if typeStr, ok0 := refData.(string); ok0 {
												field.Type = strings.TrimPrefix(typeStr, componentSchemaPrefix)
											}
										}
									}
								}
							}
						}

						if oo, ok4 := typeMap["oneOf"]; ok4 && field.Type == "" {
							if oneOf, ok5 := oo.([]any); ok5 {
								for _, x := range oneOf {
									if y, ok6 := x.(map[string]any); ok6 {
										if z, ok7 := y["type"]; ok7 && z.(string) != "null" {
											field.Type = z.(string)
										} else if z, ok7 = y["type"]; ok7 && z.(string) == "null" {
											field.Nullable = true
										} else if ref, ok8 := typeMap["$ref"]; ok8 {
											field.Type = strings.TrimPrefix(ref.(string), componentSchemaPrefix)
										} else if yRef, ok9 := y["$ref"]; ok9 {
											field.Type = strings.TrimPrefix(yRef.(string), componentSchemaPrefix)
										}
									}
								}
							}
						}

						if field.Type == "" {
							if ref, ok4 := typeMap["$ref"]; ok4 {
								field.Type = strings.TrimPrefix(ref.(string), componentSchemaPrefix)
							}
						}

						if field.Type == "" {
							if val, ok4 := codegen.EnumTypeMap[fmt.Sprintf("%s.%s", name, k)]; ok4 {
								field.Type = val
							}
						}
					}

					if x, ok3 := typeReplacementMap[field.Type]; ok3 {
						field.Type = x
					}

					switch {
					case slices.Contains([]string{
						"OptionalFloat32Range",
						"OptionalUint32Range",
					}, field.Type):
						if _, ok3 := def.Imports["./number_range"]; ok3 {
							def.Imports["./number_range"] = append(def.Imports["./number_range"], typeNameNumberRange)
						} else {
							def.Imports["./number_range"] = []string{typeNameNumberRange}
						}

						field.Type = typeNameNumberRange
						field.DefaultValue = "{ min: 0, max: 0 }"
					case slices.Contains([]string{
						"Float32RangeWithOptionalMax",
						"Uint16RangeWithOptionalMax",
						"Uint32RangeWithOptionalMax",
					}, field.Type):
						if _, ok3 := def.Imports["./number_range"]; ok3 {
							def.Imports["./number_range"] = append(def.Imports["./number_range"], typeNameNumberRangeWithOptionalMax)
						} else {
							def.Imports["./number_range"] = []string{typeNameNumberRangeWithOptionalMax}
						}

						field.Type = typeNameNumberRangeWithOptionalMax
						field.DefaultValue = "{ min: 0 }"
					case slices.Contains([]string{
						"Float32RangeWithOptionalMaxUpdateRequestInput",
						"Uint16RangeWithOptionalMaxUpdateRequestInput",
						"Uint32RangeWithOptionalMaxUpdateRequestInput",
					}, field.Type):
						if _, ok3 := def.Imports["./number_range"]; ok3 {
							def.Imports["./number_range"] = append(def.Imports["./number_range"], typeNameOptionalNumberRange)
						} else {
							def.Imports["./number_range"] = []string{typeNameOptionalNumberRange}
						}

						field.Type = typeNameOptionalNumberRange
						field.DefaultValue = "{}"
					}

					if name == "MealPlan" && k == "status" {
						println("")
					}

					if x, ok3 := codegen.EnumTypeMap[fmt.Sprintf("%s.%s", name, k)]; ok3 {
						field.Type = x
						field.Enum = true
						field.DefaultValue = codegen.DefaultEnumValues[x]
						if _, ok4 := def.Imports["./enums"]; ok4 {
							def.Imports["./enums"] = append(def.Imports["./enums"], x)
						} else {
							def.Imports["./enums"] = []string{x}
						}
					}

					nativeTypes := []string{
						"string",
						"boolean",
						"number",
						"object",
						typeNameNumberRange,
						typeNameOptionalNumberRange,
						typeNameNumberRangeWithOptionalMax,
					}

					if !slices.Contains(nativeTypes, field.Type) && !slices.Contains(enumNames, field.Type) && def.Name != field.Type {
						def.Imports[fmt.Sprintf("./%s", field.Type)] = []string{
							field.Type,
						}
					}

					if field.Type == "object" {
						field.DefaultValue = "{}"
					}

					if field.DefaultValue == "" {
						field.DefaultValue = defaultValueForType(field.Type)
					}

					def.Fields = append(def.Fields, field)
				}
			}
		}

		output[name] = def
	}

	return output, nil
}

func (d *TypeDefinition) Render() (string, error) {
	for k, v := range d.Imports {
		imports := []string{}

		for _, x := range v {
			if !slices.Contains(imports, x) {
				imports = append(imports, x)
			}
		}

		d.Imports[k] = imports
	}

	tmpl := `{{- range $key, $values := .Imports}} import { {{ join $values ", " }} } from '{{ $key }}';
{{ end }}

export interface I{{ .Name }} {
  {{ range .Fields}} {{ .Name }}{{ if .Nullable}}?{{ end }}: {{ .Type }};
{{ end }}
}

export class {{ .Name }} implements I{{ .Name }} {
  {{ range .Fields}} {{ .Name }}{{ if .Nullable}}?{{ end }}: {{ .Type }};
{{ end -}}

  constructor(input: Partial<{{ .Name }}> = {}) {
	{{ range .Fields}} this.{{.Name}} = input.{{.Name}}{{ if not .Nullable }} = {{ .DefaultValue }}{{ end }};
{{ end -}}
  }
}`

	t := template.Must(template.New("model").Funcs(map[string]any{
		"lowercase": strings.ToLower,
		"typeIsNative": func(x string) bool {
			return slices.Contains([]string{
				"string",
			}, x)
		},
		"join": strings.Join,
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, d); err != nil {
		return "", err
	}

	return b.String(), nil
}
