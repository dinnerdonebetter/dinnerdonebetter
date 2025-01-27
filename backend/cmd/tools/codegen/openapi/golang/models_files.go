package golang

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen/openapi/enums"

	"github.com/swaggest/openapi-go/openapi31"
)

const (
	stringType = "string"
	nullType   = "null"

	typeNameNumberRange                = "NumberRange"
	typeNameNumberRangeWithOptionalMax = "NumberRangeWithOptionalMax"
	typeNameOptionalNumberRange        = "OptionalNumberRange"
)

type Field struct {
	Name,
	Type string
	Nullable,
	Array bool
}

type TypeDefinition struct {
	Imports map[string][]string
	Name    string
	Fields  []Field
}

var typeReplacementMap = map[string]string{
	"integer": "number",
}

var skipTypes = map[string]bool{
	"APIResponse": true,
}

func GenerateModelFiles(spec *openapi31.Spec) (map[string]*TypeDefinition, error) {
	output := map[string]*TypeDefinition{}

	for name, component := range spec.Components.Schemas {
		def := &TypeDefinition{
			Name:    name,
			Imports: map[string][]string{},
		}

		if _, ok := skipTypes[name]; ok {
			continue
		}

		// we'll handle these later
		if _, isEnum := component["enum"]; isEnum {
			continue
		}

		if properties, ok := component[propertiesKey]; ok {
			if propMap, ok2 := properties.(map[string]any); ok2 {
				for k, v := range propMap {
					field := Field{
						Name: k,
					}

					if typeMap, ok3 := v.(map[string]any); ok3 {
						if typ, ok5 := typeMap["type"]; ok5 {
							if fieldStr, ok6 := typ.(string); ok6 {
								field.Type = fieldStr
							}

							if field.Type == "array" {
								field.Array = true
								if items, ok6 := typeMap["items"]; ok6 {
									if itemsMap, ok7 := items.(map[string]any); ok7 {
										if typeData, ok8 := itemsMap["type"]; ok8 {
											if typeStr, ok9 := typeData.(string); ok9 {
												field.Type = typeStr
											}
										} else if refData, ok9 := itemsMap[refKey]; ok9 {
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
										if z, ok7 := y["type"]; ok7 {
											if zstr, ok8 := z.(string); ok8 && zstr != nullType {
												field.Type = zstr
											}
										} else if z, ok7 = y["type"]; ok7 {
											if zstr, ok8 := z.(string); ok8 && zstr != nullType {
												field.Nullable = true
											}
										} else if ref, ok8 := typeMap[refKey]; ok8 {
											if refString, ok9 := ref.(string); ok9 {
												field.Type = strings.TrimPrefix(refString, componentSchemaPrefix)
											}
										} else if yRef, ok9 := y[refKey]; ok9 {
											if yRefString, ok0 := yRef.(string); ok0 {
												field.Type = strings.TrimPrefix(yRefString, componentSchemaPrefix)
											}
										}
									}
								}
							}
						}

						if field.Type == "" {
							if ref, ok4 := typeMap[refKey]; ok4 {
								if refString, ok5 := ref.(string); ok5 {
									field.Type = strings.TrimPrefix(refString, componentSchemaPrefix)
								}
							}
						}

						if field.Type == "" {
							if val, ok4 := enums.TypeMap[fmt.Sprintf("%s.%s", name, k)]; ok4 {
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
						field.Type = typeNameNumberRange
					case slices.Contains([]string{
						"Float32RangeWithOptionalMax",
						"Uint16RangeWithOptionalMax",
						"Uint32RangeWithOptionalMax",
					}, field.Type):
						field.Type = typeNameNumberRangeWithOptionalMax
					case slices.Contains([]string{
						"Float32RangeWithOptionalMaxUpdateRequestInput",
						"Uint16RangeWithOptionalMaxUpdateRequestInput",
						"Uint32RangeWithOptionalMaxUpdateRequestInput",
					}, field.Type):
						field.Type = typeNameOptionalNumberRange
					}

					if x, ok3 := enums.TypeMap[fmt.Sprintf("%s.%s", name, k)]; ok3 {
						field.Type = x
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

	slices.SortFunc(d.Fields, func(a, b Field) int {
		switch {
		case a.Name < b.Name:
			return -1
		case a.Name == b.Name:
			return 0
		default:
			return 1
		}
	})

	tmpl := `{{- range $key, $values := .Imports}} import { {{ join (sortStrings $values) ", " }} } from '{{ $key }}';
{{ end }}

export interface I{{ .Name }} {
  {{ range .Fields}} {{ .Name }}{{ if .Nullable}}?{{ end }}: {{ .Type }}{{ if .Array }}[]{{ end }};
{{ end }}
}

export class {{ .Name }} implements I{{ .Name }} {
  {{ range .Fields}} {{ .Name }}{{ if .Nullable}}?{{ end }}: {{ .Type }}{{ if .Array }}[]{{ end }};
{{ end -}}

  constructor(input: Partial<{{ .Name }}> = {}) {
	{{ range .Fields}} this.{{.Name}} = input.{{.Name}}{{ if not .Nullable }} || {{ .DefaultValue }}{{ end }};
{{ end -}}
  }
}`

	t := template.Must(template.New("model").Funcs(map[string]any{
		"lowercase": strings.ToLower,
		"typeIsNative": func(x string) bool {
			return slices.Contains([]string{
				stringType,
			}, x)
		},
		"join": strings.Join,
		"sortStrings": func(s []string) []string {
			slices.Sort(s)
			return s
		},
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, d); err != nil {
		return "", err
	}

	return b.String(), nil
}
