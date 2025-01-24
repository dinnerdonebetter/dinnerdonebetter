package typescript

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"

	"github.com/swaggest/openapi-go/openapi31"
)

func GenerateMockAPIFunctions(spec *openapi31.Spec) (map[string]*MockAPIFunction, error) {
	output := map[string]*MockAPIFunction{}

	for path, op := range spec.Paths.MapOfPathItemValues {
		if op.Get != nil {
			opID := *op.Get.ID
			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildMockFunction(path, http.MethodGet, op.Get)
		}

		if op.Put != nil {
			opID := *op.Put.ID

			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildMockFunction(path, http.MethodPut, op.Put)
		}

		if op.Patch != nil {
			opID := *op.Patch.ID
			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildMockFunction(path, http.MethodPatch, op.Patch)
		}

		if op.Post != nil {
			opID := *op.Post.ID
			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildMockFunction(path, http.MethodPost, op.Post)
		}

		if op.Delete != nil {
			opID := *op.Delete.ID
			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildMockFunction(path, http.MethodDelete, op.Delete)
		}
	}

	return output, nil
}

func buildMockFunction(path, method string, op *openapi31.Operation) *MockAPIFunction {
	f := buildFunction(path, method, op)

	params := []MockAPIPathParam{}
	for _, param := range f.Params {
		params = append(params, MockAPIPathParam{
			Name: param.Name,
			Type: param.Type,
		})
	}

	return &MockAPIFunction{
		Name:                 uppercaseFirstLetter(f.Name),
		ResponseType:         f.ResponseType.TypeName,
		Method:               f.Method,
		PathTemplate:         strings.ReplaceAll(path, "{", "${resCfg."),
		PathParams:           params,
		QueryFiltered:        f.QueryFiltered,
		Search:               strings.Contains(path, "search"),
		ExpectedResponseCode: f.DefaultStatusCode,
	}
}

type MockAPIPathParam struct {
	Name, Type string
}

type MockAPIFunction struct {
	Name,
	ResponseType,
	Method,
	PathTemplate string
	PathParams []MockAPIPathParam
	QueryFiltered,
	Search bool
	ExpectedResponseCode uint16
}

func (f *MockAPIFunction) Render() (string, error) {
	tmpl := `export class Mock{{ .Name }}ResponseConfig extends ResponseConfig<{{ if .QueryFiltered }}QueryFilteredResult<{{ end }}{{ .ResponseType }}{{ if .QueryFiltered }}>{{ end }}> {
		  {{ range .PathParams}} {{ .Name }}: {{ .Type }};
		{{ end }}

 		  constructor({{ range .PathParams}} {{ .Name }}: {{ .Type }}, {{ end }}status: number = {{ .ExpectedResponseCode }}, body{{ if not (or .Search .QueryFiltered) }}?{{ end }}: {{ .ResponseType }}{{ if or .Search .QueryFiltered }}[] = []{{ end }}) {
		    super();

		{{ range .PathParams}} this.{{ .Name }} = {{ .Name }};
		{{ end }}
		    this.status = status;
			if (this.body) {
			  this.body{{ if .QueryFiltered }}.data{{ end }} = body;
			}
		  }
}

export const mock{{ .Name }}{{ if .QueryFiltered }}s{{ end }} = (resCfg: Mock{{ .Name }}ResponseConfig) => {
  return (page: Page) =>
    page.route(
      ` + "`" + `**{{ .PathTemplate }}` + "`" + `,
      (route: Route) => {
        const req = route.request();

        assertMethod('{{ uppercase .Method }}', route);
        assertClient(route);

		{{ if .QueryFiltered }}
        if (resCfg.body && resCfg.filter) resCfg.body.limit = resCfg.filter.limit;
        if (resCfg.body && resCfg.filter) resCfg.body.page = resCfg.filter.page;
		{{ end }}

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};`

	t := template.Must(template.New("function").Funcs(map[string]any{
		"uppercase": strings.ToUpper,
		"title": func(s string) string {
			return uppercaseFirstLetter(strings.ToLower(s))
		},
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, f); err != nil {
		return "", err
	}

	result := b.String()

	return result, nil
}
