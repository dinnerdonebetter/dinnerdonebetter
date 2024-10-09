package typescript

import (
	"bytes"
	"net/http"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/swaggest/openapi-go/openapi31"
)

var skipOps = map[string]bool{
	"CheckForLiveness":  true,
	"CheckForReadiness": true,
	"GetRecipeDAG":      true,
}

func GenerateClientFiles(spec *openapi31.Spec) (map[string]*APIClientFunction, error) {
	output := map[string]*APIClientFunction{}

	for path, op := range spec.Paths.MapOfPathItemValues {
		if op.Get != nil {
			opID := *op.Get.ID

			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildFunction(path, http.MethodGet, op.Get)
		}

		if op.Put != nil {
			opID := *op.Put.ID

			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildFunction(path, http.MethodPut, op.Put)
		}

		if op.Patch != nil {
			opID := *op.Patch.ID

			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildFunction(path, http.MethodPatch, op.Patch)
		}

		if op.Post != nil {
			opID := *op.Post.ID

			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildFunction(path, http.MethodPost, op.Post)
		}

		if op.Delete != nil {
			opID := *op.Delete.ID

			if _, ok := skipOps[opID]; ok {
				continue
			}

			output[opID] = buildFunction(path, http.MethodDelete, op.Delete)
		}
	}

	return output, nil
}

// paramsContainQueryFilter returns whether a query filter is in the params or not, and then a QueryFilter-less list of params.
func paramsContainQueryFilter(params []openapi31.ParameterOrReference) (bool, []openapi31.ParameterOrReference) {
	var (
		containsLimit,
		containsPage,
		containsCreatedAfter,
		containsCreatedBefore,
		containsUpdatedAfter,
		containsUpdatedBefore,
		containsSortBy,
		containsIncludeArchived bool
	)

	outParams := []openapi31.ParameterOrReference{}
	for _, p := range params {
		switch p.Parameter.Name {
		case types.QueryKeyLimit:
			containsLimit = true
		case types.QueryKeyPage:
			containsPage = true
		case types.QueryKeyCreatedBefore:
			containsCreatedBefore = true
		case types.QueryKeyCreatedAfter:
			containsCreatedAfter = true
		case types.QueryKeyUpdatedBefore:
			containsUpdatedBefore = true
		case types.QueryKeyUpdatedAfter:
			containsUpdatedAfter = true
		case types.QueryKeyIncludeArchived:
			containsIncludeArchived = true
		case types.QueryKeySortBy:
			containsSortBy = true
		default:
			outParams = append(outParams, p)
		}
	}

	return containsLimit &&
		containsPage &&
		containsCreatedAfter &&
		containsCreatedBefore &&
		containsUpdatedAfter &&
		containsUpdatedBefore &&
		containsSortBy &&
		containsIncludeArchived, outParams
}

func lowercaseFirstLetter(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

var postSkipPaths = map[string]bool{
	"/api/v1/households/{householdID}/default": true,
	"/api/v1/recipes/{recipeID}/clone":         true,
	"/api/v1/meal_plans/{mealPlanID}/finalize": true,
}

func buildFunction(path, method string, op *openapi31.Operation) *APIClientFunction {
	functionName := lowercaseFirstLetter(*op.ID)
	containsQF, params := paramsContainQueryFilter(op.Parameters)

	functionParams := []functionParam{}
	for _, p := range params {
		functionParams = append(functionParams, functionParam{
			Name: p.Parameter.Name,
			Type: p.Parameter.Schema["type"].(string),
		})
	}

	var schema map[string]any
	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["200"]; ok {
		if response, ok2 := responseContainer.Response.Content["application/json"]; ok2 {
			schema = response.Schema
		}

		// special snowflake
		if response, ok2 := responseContainer.Response.Content["text/mermaid"]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["201"]; ok {
		if response, ok2 := responseContainer.Response.Content["application/json"]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["202"]; ok {
		if response, ok2 := responseContainer.Response.Content["application/json"]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["204"]; ok {
		if response, ok2 := responseContainer.Response.Content["application/json"]; ok2 {
			schema = response.Schema
		}
	}

	ip := functionInputParam{}
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if _, skip := postSkipPaths[path]; skip {
			break
		}

		if op.RequestBody != nil {
			if jsonContent, ok := op.RequestBody.RequestBody.Content["application/json"]; ok {
				if refContent, ok2 := jsonContent.Schema["$ref"]; ok2 {
					ip.Name = lowercaseFirstLetter(strings.TrimPrefix(refContent.(string), "#/components/schemas/"))
					ip.Type = strings.TrimPrefix(refContent.(string), "#/components/schemas/")
				}
			}
		}
	}

	rt := functionResponseType{}
	if allOf, ok1 := schema["allOf"]; ok1 {
		switch x := allOf.(type) {
		case []any:
			for _, yy := range x {
				if y, ok := yy.(map[string]interface{}); ok {
					if allOfRef, ok2 := y["$ref"]; ok2 {
						if allOfRefStr, ok3 := allOfRef.(string); ok3 {
							rt.GenericContainer = allOfRefStr
						}
					}

					if properties, ok2 := y["properties"]; ok2 {
						if props, ok3 := properties.(map[string]any); ok3 {
							if rawData, ok4 := props["data"]; ok4 {
								if data, ok5 := rawData.(map[string]any); ok5 {
									if dataRef, ok6 := data["$ref"]; ok6 {
										rt.TypeName = dataRef.(string)
									} else if itemsRef, ok7 := data["items"]; ok7 {
										if z, ok9 := itemsRef.(map[string]any); ok9 {
											if itemsDataRef, ok8 := z["$ref"]; ok8 {
												rt.TypeName = itemsDataRef.(string)
											}
										}
									} else if rawType, ok8 := data["type"]; ok8 {
										if z, ok9 := rawType.(string); ok9 {
											rt.TypeName = z
										}
									}
								}
							}
						}
					}
				}
			}
		}
	} else if schema != nil {
		// if this fails, I want it to panic
		rt.TypeName = schema["type"].(string)
	}

	rt.GenericContainer = strings.TrimPrefix(rt.GenericContainer, "#/components/schemas/")
	rt.TypeName = strings.TrimPrefix(rt.TypeName, "#/components/schemas/")

	return &APIClientFunction{
		Name:          functionName,
		Method:        method,
		QueryFiltered: containsQF,
		ReturnsList:   containsQF,
		Params:        functionParams,
		InputType:     ip,
		ResponseType:  rt,
		PathTemplate:  strings.ReplaceAll(path, "{", "${"),
	}
}

type functionInputParam struct {
	Name,
	Type string
}

type functionParam struct {
	Name,
	Type,
	DefaultValue string
	InQuery bool
}

type APIClientFunction struct {
	InputType     functionInputParam
	ResponseType  functionResponseType
	Method        string
	Name          string
	PathTemplate  string
	Params        []functionParam
	QueryFiltered bool
	ReturnsList   bool
}

type functionResponseType struct {
	TypeName         string
	GenericContainer string
}

func (f *APIClientFunction) Render() (string, error) {
	var tmpl string

	switch f.Method {
	case http.MethodGet, http.MethodDelete:
		if f.QueryFiltered {
			// GET routes that return lists

			tmpl = `import { Axios } from 'axios';

import {
  {{ if and (ne .ResponseType.TypeName "") (not (typeIsNative .ResponseType.TypeName)) }}{{ .ResponseType.TypeName }}, {{ end }}
  QueryFilter,
  QueryFilteredResult,
  {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}, {{ end }}
} from '@dinnerdonebetter/models'; 

export async function {{ .Name }}(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  {{ range .Params }}{{ .Name }}: {{ .Type }}{{ if ne .DefaultValue "" }} = {{ .DefaultValue }}{{ end }},
	{{ end -}}): Promise< QueryFilteredResult< {{ .ResponseType.TypeName }} >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.{{ lowercase .Method }}< {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if .ReturnsList }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if .ReturnsList }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} >(` + "`" + "{{ .PathTemplate }}" + "`" + `, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<{{ .ResponseType.TypeName }}>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}`
		} else {
			// GET routes that don't return lists

			tmpl = `import { Axios } from 'axios';

import {
  {{ if and (ne .ResponseType.TypeName "") (not (typeIsNative .ResponseType.TypeName)) }}{{ .ResponseType.TypeName }}, {{ end }}
  {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}, {{ end }}
} from '@dinnerdonebetter/models'; 

export async function {{ .Name }}(
  client: Axios,
  {{ range .Params }}{{ .Name }}: {{ .Type }}{{ if ne .DefaultValue "" }} = {{ .DefaultValue }}{{ end }},
	{{ end -}}): Promise<  {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }} {{ .ResponseType.TypeName }} >   {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}  {
  return new Promise(async function (resolve, reject) {
    const response = await client.{{ lowercase .Method }}< {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if .ReturnsList }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if .ReturnsList }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} >(` + "`" + "{{ .PathTemplate }}" + "`" + `, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}`
		}

	case http.MethodPost, http.MethodPut, http.MethodPatch:
		tmpl = `import { Axios } from 'axios';

import {
  {{ if and (ne .ResponseType.TypeName "") (not (typeIsNative .ResponseType.TypeName)) }}{{ .ResponseType.TypeName }}, {{ end }}
  {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}, {{ end }}
  {{ if ne .InputType.Type "" }}{{ .InputType.Type }}, {{ end }}
} from '@dinnerdonebetter/models';

export async function {{ .Name }}(
  client: Axios,
  {{ range .Params }}{{ .Name }}: {{ .Type }},{{ end }}
  {{ if ne .InputType.Type "" }}input: {{ .InputType.Type }},{{ end }}
): Promise<  {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }} {{ .ResponseType.TypeName }} > {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} {
  return new Promise(async function (resolve, reject) {
    const response = await client.{{ lowercase .Method }}<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if .ReturnsList }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if .ReturnsList }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} >(` + "`" + "{{ .PathTemplate }}" + "`" + `{{ if ne .InputType.Type "" }}, input{{ end }});
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}`
	}

	if tmpl == "" {
		panic("Unknown template")
	}

	t := template.Must(template.New("function").Funcs(map[string]any{
		"lowercase": strings.ToLower,
		"typeIsNative": func(x string) bool {
			return slices.Contains([]string{
				"string",
			}, x)
		},
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, f); err != nil {
		return "", err
	}

	return b.String(), nil
}
