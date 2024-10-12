package golang

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
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

func uppercaseFirstLetter(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

var postSkipPaths = map[string]bool{
	"/api/v1/households/{householdID}/default": true,
	"/api/v1/recipes/{recipeID}/clone":         true,
	"/api/v1/meal_plans/{mealPlanID}/finalize": true,
}

var urlParamFinderRegex = regexp.MustCompile(`\{[a-zA-Z\d]+\}`)

func buildFunction(path, method string, op *openapi31.Operation) *APIClientFunction {
	functionName := *op.ID
	containsQF, params := paramsContainQueryFilter(op.Parameters)

	functionParams := []functionParam{}
	for _, p := range params {
		if schemaStr, ok := p.Parameter.Schema["type"].(string); ok {
			functionParams = append(functionParams, functionParam{
				Name: p.Parameter.Name,
				Type: schemaStr,
			})
		}
	}

	var (
		defaultStatusCode uint16
		schema            map[string]any
	)
	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["200"]; ok {
		defaultStatusCode = 200
		if response, ok2 := responseContainer.Response.Content[jsonContentType]; ok2 {
			schema = response.Schema
		}

		// special snowflake
		if response, ok2 := responseContainer.Response.Content["text/mermaid"]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["201"]; ok {
		defaultStatusCode = 201
		if response, ok2 := responseContainer.Response.Content[jsonContentType]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["202"]; ok {
		defaultStatusCode = 202
		if response, ok2 := responseContainer.Response.Content[jsonContentType]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["204"]; ok {
		defaultStatusCode = 204
		if response, ok2 := responseContainer.Response.Content[jsonContentType]; ok2 {
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
			if jsonContent, ok := op.RequestBody.RequestBody.Content[jsonContentType]; ok {
				if refContent, ok2 := jsonContent.Schema[refKey]; ok2 {
					if refContentStr, ok3 := refContent.(string); ok3 {
						ip.Name = lowercaseFirstLetter(strings.TrimPrefix(refContentStr, componentSchemaPrefix))
						ip.Type = strings.TrimPrefix(refContentStr, componentSchemaPrefix)
					}
				}
			}
		}
	}

	returnsList := containsQF

	rt := functionResponseType{}
	if allOf, ok1 := schema["allOf"]; ok1 {
		switch x := allOf.(type) {
		case []any:
			for _, yy := range x {
				if y, ok := yy.(map[string]interface{}); ok {
					if allOfRef, ok2 := y[refKey]; ok2 {
						if allOfRefStr, ok3 := allOfRef.(string); ok3 {
							rt.GenericContainer = allOfRefStr
						}
					}

					if properties, ok2 := y[propertiesKey]; ok2 {
						if props, ok3 := properties.(map[string]any); ok3 {
							if rawData, ok4 := props["data"]; ok4 {
								if data, ok5 := rawData.(map[string]any); ok5 {
									if dataRef, ok6 := data[refKey]; ok6 {
										if z, ok7 := dataRef.(string); ok7 {
											rt.TypeName = z
										}
									} else if itemsRef, ok7 := data["items"]; ok7 {
										rt.IsArray = true
										if !returnsList {
											returnsList = true
										}

										if z, ok9 := itemsRef.(map[string]any); ok9 {
											if itemsDataRef, ok8 := z[refKey]; ok8 {
												if zz, ok0 := itemsDataRef.(string); ok0 {
													rt.TypeName = zz
												}
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
		if typeStr, ok := schema["type"].(string); ok {
			rt.TypeName = typeStr
		}
	}

	rt.GenericContainer = strings.TrimPrefix(rt.GenericContainer, componentSchemaPrefix)
	rt.TypeName = strings.TrimPrefix(rt.TypeName, componentSchemaPrefix)

	return &APIClientFunction{
		Name:              functionName,
		Method:            method,
		QueryFiltered:     containsQF,
		DefaultStatusCode: defaultStatusCode,
		ReturnRawResponse: slices.Contains([]string{"updatePassword", "loginForJWT", "adminLoginForJWT"}, functionName),
		ReturnsList:       returnsList,
		Params:            functionParams,
		InputType:         ip,
		ResponseType:      rt,
		PathTemplate:      urlParamFinderRegex.ReplaceAllString(path, "%s"),
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
}

type APIClientFunction struct {
	InputType    functionInputParam
	ResponseType functionResponseType
	Method,
	Name,
	PathTemplate string
	Params            []functionParam
	DefaultStatusCode uint16
	ReturnRawResponse,
	QueryFiltered,
	ReturnsList bool
}

type functionResponseType struct {
	TypeName         string
	GenericContainer string
	IsArray          bool
}

func (f *APIClientFunction) Render() (string, []string, error) {
	var tmpl string
	imports := []string{}

	switch f.Method {
	case http.MethodGet:
		if f.QueryFiltered {
			// GET routes that return lists

			tmpl = `func (c *Client) {{ if not (contains (lowercase .Name) (lowercase .Method))}}{{ title .Method }}{{ end }}{{ .Name }}(
	ctx context.Context,
	{{ range .Params }}{{ .Name }} {{ .Type }},
	{{ end -}}
	filter *types.QueryFilter,
) (*types.QueryFilteredResult[types.{{ .ResponseType.TypeName }}], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

{{ range .Params }}	if {{ .Name }} == "" {
		return nil, buildInvalidIDError("{{ replace .Name "ID" "" }}")
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})

{{ end }} 

	values := filter.ToValues()
	{{ if paramsContain .Params "q" -}}
	values.Set(types.QueryKeySearch, q)
	{{- end }}

	u := c.BuildURL(ctx, values, fmt.Sprintf("{{ .PathTemplate }}" {{ range .Params }}{{ if ne .Name "q" }}, {{ .Name }}{{ end }}{{ end }}))
	req, err := http.NewRequestWithContext(ctx, http.Method{{ title .Method }}, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of {{ .ResponseType.TypeName }}")
	}
	
	var apiResponse *types.{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}[ {{ end }}[]*types.{{ .ResponseType.TypeName }}{{ if ne .ResponseType.GenericContainer "" }}]{{ end }}
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of {{ .ResponseType.TypeName }}")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.{{ .ResponseType.TypeName }}]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}`
			imports = append(imports,
				"fmt",
				"github.com/dinnerdonebetter/backend/internal/observability",
				"github.com/dinnerdonebetter/backend/internal/observability/tracing",
			)
			if len(f.Params) > 0 {
				imports = append(imports,
					"github.com/dinnerdonebetter/backend/internal/observability/keys")
			}

		} else {
			// GET routes that don't return lists
			tmpl = `func (c *Client) {{ if not (contains (lowercase .Name) (lowercase .Method))}}{{ title .Method }}{{ end }}{{ .Name }}(
	ctx context.Context,
{{ range .Params }}{{ .Name }} {{ .Type }},
{{ end -}}
) ({{ if notNative .ResponseType.TypeName }} *types.{{ end }}{{ .ResponseType.TypeName }}, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

{{ range .Params }}	if {{ .Name }} == "" {
		return {{ if notNative $.ResponseType.TypeName }}nil{{ else }} {{ nativeDefault $.ResponseType.TypeName }}{{ end }}, buildInvalidIDError("{{ replace .Name "ID" "" }}")
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})

{{ end }} 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("{{ .PathTemplate }}" {{ range .Params }}, {{ .Name }} {{ end }}))
	req, err := http.NewRequestWithContext(ctx, http.Method{{ title .Method }}, u, http.NoBody)
	if err != nil {
		return {{ if notNative .ResponseType.TypeName }}nil{{ else }} {{ nativeDefault .ResponseType.TypeName }}{{ end }}, observability.PrepareAndLogError(err, logger, span, "building request to fetch a {{ .ResponseType.TypeName }}")
	}

	var apiResponse *types.{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}[ {{ end }}{{ if notNative .ResponseType.TypeName }} *types.{{ end }}{{ .ResponseType.TypeName }}{{ if ne .ResponseType.GenericContainer "" }}]{{ end }}
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return {{ if notNative .ResponseType.TypeName }}nil{{ else }} {{ nativeDefault .ResponseType.TypeName }}{{ end }}, observability.PrepareAndLogError(err, logger, span, "loading {{ .ResponseType.TypeName }} response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return {{ if notNative .ResponseType.TypeName }}nil{{ else }} {{ nativeDefault .ResponseType.TypeName }}{{ end }}, err
	}


	return apiResponse.Data, nil
}`
			imports = append(imports,
				"fmt",
				"github.com/dinnerdonebetter/backend/internal/observability",
			)
			if len(f.Params) > 0 {
				imports = append(imports,
					"github.com/dinnerdonebetter/backend/internal/observability/tracing",
					"github.com/dinnerdonebetter/backend/internal/observability/keys")
			}

		}

	case http.MethodPost:
		tmpl = `func (c *Client) {{ .Name }}(
	ctx context.Context,
{{ range .Params }}{{ .Name }} {{ .Type }},
{{ end -}}
	{{ if ne .InputType.Type "" }}input *types.{{ .InputType.Type }},{{ end }}
) (*types.{{ .ResponseType.TypeName }}, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

{{ if ne .InputType.Type "" }}
	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}
{{ end }}
{{ range .Params }}	if {{ .Name }} == "" {
		return nil, buildInvalidIDError("{{ replace .Name "ID" "" }}")
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})

{{ end }} 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("{{ .PathTemplate }}" {{ range .Params }}, {{ .Name }} {{ end }}))
	req, err := c.buildDataRequest(ctx, http.Method{{ title .Method }}, u, {{ if ne .InputType.Type "" }}input{{ else }}http.NoBody{{ end }})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	var apiResponse *types.{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}[ {{ end }}*types.{{ .ResponseType.TypeName }}{{ if ne .ResponseType.GenericContainer "" }}]{{ end }}
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading {{ .ResponseType.TypeName }} creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}


	return apiResponse.Data, nil
}`

		imports = append(imports,
			"fmt",
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if len(f.Params) > 0 {
			imports = append(imports,
				"github.com/dinnerdonebetter/backend/internal/observability/tracing",
				"github.com/dinnerdonebetter/backend/internal/observability/keys")
		}

	case http.MethodPut:
		tmpl = `func (c *Client) {{ .Name }}(
	ctx context.Context,
{{ range .Params }}{{ .Name }} {{ .Type }},
{{ end -}}
input *types.{{ .InputType.Type }},
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

{{ range .Params }}	if {{ .Name }} == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})

{{ end }} 


	u := c.BuildURL(ctx, nil, fmt.Sprintf("{{ .PathTemplate }}" {{ range .Params }}, {{ .Name }} {{ end }}))
	req, err := c.buildDataRequest(ctx, http.Method{{ title .Method }}, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	var apiResponse *types.{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}[ {{ end }}*types.{{ .ResponseType.TypeName }}{{ if ne .ResponseType.GenericContainer "" }}]{{ end }}
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading {{ .ResponseType.TypeName }} creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}


	return nil
}`
		imports = append(imports,
			"fmt",
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if len(f.Params) > 0 {
			imports = append(imports,
				"github.com/dinnerdonebetter/backend/internal/observability/tracing",
				"github.com/dinnerdonebetter/backend/internal/observability/keys")
		}

	case http.MethodPatch:
		tmpl = `func (c *Client) {{ .Name }}(
	ctx context.Context,
{{ range .Params }}{{ .Name }} {{ .Type }},
{{ end -}}
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

{{ range .Params }}	if {{ .Name }} == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})

{{ end }} 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("{{ .PathTemplate }}" {{ range .Params }}, {{ .Name }} {{ end }}))
	req, err := http.NewRequestWithContext(ctx, http.Method{{ title .Method }}, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	var apiResponse *types.{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}[ {{ end }}*types.{{ .ResponseType.TypeName }}{{ if ne .ResponseType.GenericContainer "" }}]{{ end }}
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading {{ .ResponseType.TypeName }} creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}

	return nil
}`
		imports = append(imports,
			"fmt",
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if len(f.Params) > 0 {
			imports = append(imports,
				"github.com/dinnerdonebetter/backend/internal/observability/tracing",
				"github.com/dinnerdonebetter/backend/internal/observability/keys")
		}

	case http.MethodDelete:
		// GET routes that don't return lists
		tmpl = `func (c *Client) {{ .Name }}(
	ctx context.Context,
{{ range .Params }}{{ .Name }} {{ .Type }},
{{ end -}}
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

{{ range .Params }}	if {{ .Name }} == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})

{{ end }} 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("{{ .PathTemplate }}" {{ range .Params }}, {{ .Name }} {{ end }}))
	req, err := http.NewRequestWithContext(ctx, http.Method{{ title .Method }}, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	var apiResponse *types.{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}[ {{ end }}*types.{{ .ResponseType.TypeName }}{{ if ne .ResponseType.GenericContainer "" }}]{{ end }}
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading {{ .ResponseType.TypeName }} creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}

	return  nil
}`
		imports = append(imports,
			"fmt",
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if len(f.Params) > 0 {
			imports = append(imports,
				"github.com/dinnerdonebetter/backend/internal/observability/tracing",
				"github.com/dinnerdonebetter/backend/internal/observability/keys")
		}

	}

	if tmpl == "" {
		panic("Unknown template")
	}

	t := template.Must(template.New("function").Funcs(map[string]any{
		"lowercase": strings.ToLower,
		"contains":  strings.Contains,
		"title": func(s string) string {
			return uppercaseFirstLetter(strings.ToLower(s))
		},
		"replace":              strings.ReplaceAll,
		"uppercaseFirstLetter": uppercaseFirstLetter,
		"notNative": func(s string) bool {
			switch s {
			case "string", "bool", "int", "uint64":
				return false
			default:
				return true
			}
		},
		"observabilityKey": func(s string) string {
			out := strings.ReplaceAll(uppercaseFirstLetter(s), "Oauth", "OAuth")

			if out == "Q" {
				out = "SearchQuery"
			}

			return out
		},
		"nativeDefault": func(s string) string {
			switch s {
			case "string":
				return `""`
			default:
				panic(fmt.Sprintf("aaaaaaaaaaaaaaaa bad type: %s", s))
			}
		},
		"paramsContain": func(x []functionParam, y string) bool {
			for _, z := range x {
				if z.Name == types.QueryKeySearch {
					return true
				}
			}

			return false
		},
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, f); err != nil {
		return "", nil, err
	}

	return b.String(), imports, nil
}
