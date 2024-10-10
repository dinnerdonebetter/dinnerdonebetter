package typescript

import (
	"bytes"
	"slices"
	"text/template"
	"unicode"

	"github.com/swaggest/openapi-go/openapi31"
)

type EnumDefinition struct {
	Name   string
	Values []string
}

func (e *EnumDefinition) Render() (string, error) {
	tmpl := `export const ALL_{{ capitalize .Name }}: string[] = [
  {{ range .Values }} '{{ . }}', {{ end }}
];
type {{ .Name }}Tuple = typeof ALL_{{ capitalize .Name }};
export type {{ .Name }} = {{ .Name }}Tuple[number];
`

	t := template.Must(template.New("enum").Funcs(map[string]any{
		"capitalize": func(s string) string {
			var result []rune
			for i, r := range s {
				if unicode.IsUpper(r) && i > 0 {
					result = append(result, '_')
				}
				result = append(result, unicode.ToUpper(r))
			}
			return string(result)
		},
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, e); err != nil {
		return "", err
	}

	return b.String(), nil
}

func GenerateEnumDefinitions(spec *openapi31.Spec) ([]EnumDefinition, error) {
	output := []EnumDefinition{}

	for name, component := range spec.Components.Schemas {
		// we'll handle these later
		componentEnum, isEnum := component["enum"]
		if !isEnum {
			continue
		}

		ed := EnumDefinition{
			Name:   name,
			Values: []string{},
		}

		if componentsAny, ok := componentEnum.([]any); ok {
			for _, c := range componentsAny {
				if cString, ok2 := c.(string); ok2 {
					ed.Values = append(ed.Values, cString)
				}
			}
		}

		output = append(output, ed)
	}

	slices.SortFunc(output, func(a, b EnumDefinition) int {
		switch {
		case a.Name < b.Name:
			return -1
		case a.Name == b.Name:
			return 0
		default:
			return 1
		}
	})

	return output, nil
}
