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

const (
	typeString = "string"
	typeBool   = "bool"
	typeInt    = "int"
	typeUint64 = "uint64"
)

var skipOps = map[string]bool{
	"CheckForLiveness":  true,
	"CheckForReadiness": true,
}

func GenerateClientFunctions(spec *openapi31.Spec) (map[string]*APIClientFunction, error) {
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
		defaultStatusCode = http.StatusOK
		if response, ok2 := responseContainer.Response.Content[jsonContentType]; ok2 {
			schema = response.Schema
		}

		// special snowflake
		if response, ok2 := responseContainer.Response.Content["text/mermaid"]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["201"]; ok {
		defaultStatusCode = http.StatusCreated
		if response, ok2 := responseContainer.Response.Content[jsonContentType]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["202"]; ok {
		defaultStatusCode = http.StatusAccepted
		if response, ok2 := responseContainer.Response.Content[jsonContentType]; ok2 {
			schema = response.Schema
		}
	}

	if responseContainer, ok := op.Responses.MapOfResponseOrReferenceValues["204"]; ok {
		defaultStatusCode = http.StatusNoContent
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
				if y, ok := yy.(map[string]any); ok {
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
		ReturnRawResponse: slices.Contains([]string{"updatePassword", "loginForToken", "adminLoginForToken"}, functionName),
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

func (f *APIClientFunction) Render() (file string, imports []string, err error) {
	var tmpl string
	imports = []string{}

	shouldFormatPath := len(f.Params) > 0 && !(len(f.Params) == 1 && f.Params[0].Name == "q")

	switch f.Method {
	case http.MethodGet:
		if f.QueryFiltered {
			// GET routes that return lists

			tmpl = `func (c *Client) {{ .Name }}(
	ctx context.Context,
	{{ range .Params }}{{ .Name }} {{ .Type }},
	{{ end -}}
	filter *types.QueryFilter,
	reqMods ...RequestModifier,
) (*types.QueryFilteredResult[types.{{ .ResponseType.TypeName }}], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

{{ range .Params }}	{{ if ne .Name "q" }}if {{ .Name }} == "" {
		return nil, buildInvalidIDError("{{ replace .Name "ID" "" }}")
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})

{{ end }} 
{{ end }} 

	values := filter.ToValues()
	{{ if paramsContain .Params "q" -}}
	values.Set(types.QueryKeySearch, q)
	{{- end }}

	u := c.BuildURL(ctx, values, {{ if shouldFormatPath }}fmt.Sprintf({{ end }}"{{ .PathTemplate }}"{{if shouldFormatPath }} {{ range .Params }}{{ if ne .Name "q" }}, {{ .Name }}{{ end }}{{ end }}){{ end }})
	req, err := http.NewRequestWithContext(ctx, http.Method{{ title .Method }}, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of {{ .ResponseType.TypeName }}")
	}

	for _, mod := range reqMods {
		mod(req)
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
				"github.com/dinnerdonebetter/backend/internal/observability",
				"github.com/dinnerdonebetter/backend/internal/observability/tracing",
			)

			if shouldFormatPath {
				imports = append(imports,
					"fmt",
					"github.com/dinnerdonebetter/backend/internal/observability/keys")
			}
		} else {
			// GET routes that don't return lists
			tmpl = `func (c *Client) {{ .Name }}(
	ctx context.Context,
{{ range .Params }}{{ .Name }} {{ .Type }},
{{ end -}} reqMods ...RequestModifier,
) ({{ if notNative .ResponseType.TypeName }} *types.{{ end }}{{ .ResponseType.TypeName }}, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

{{ range .Params }}{{ if ne .Name "q" }}	if {{ .Name }} == "" {
		return {{ if notNative $.ResponseType.TypeName }}nil{{ else }} {{ nativeDefault $.ResponseType.TypeName }}{{ end }}, buildInvalidIDError("{{ replace .Name "ID" "" }}")
	} 
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})
{{ end }}
{{ end }} 

	u := c.BuildURL(ctx, nil, {{ if shouldFormatPath }}fmt.Sprintf({{ end }}"{{ .PathTemplate }}"{{if shouldFormatPath }} {{ range .Params }}{{ if ne .Name "q" }}, {{ .Name }}{{ end }}{{ end }}){{ end }})
	req, err := http.NewRequestWithContext(ctx, http.Method{{ title .Method }}, u, http.NoBody)
	if err != nil {
		return {{ if notNative .ResponseType.TypeName }}nil{{ else }} {{ nativeDefault .ResponseType.TypeName }}{{ end }}, observability.PrepareAndLogError(err, logger, span, "building request to fetch a {{ .ResponseType.TypeName }}")
	}

	for _, mod := range reqMods {
		mod(req)
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
				"github.com/dinnerdonebetter/backend/internal/observability",
			)
			if shouldFormatPath {
				imports = append(imports, "fmt")
			}
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
	reqMods ...RequestModifier,
) ({{ if .ReturnsList }}[]{{ end }}*types.{{ .ResponseType.TypeName }}, error) {
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
{{ range .Params }}	{{ if ne .Name "q" }}if {{ .Name }} == "" {
		return nil, buildInvalidIDError("{{ replace .Name "ID" "" }}")
	}
	logger = logger.WithValue(keys.{{ observabilityKey .Name }}Key, {{ .Name }})
	tracing.AttachToSpan(span, keys.{{ observabilityKey .Name }}Key, {{ .Name }})
{{ end }}

{{ end }} 

	u := c.BuildURL(ctx, nil, {{ if shouldFormatPath }}fmt.Sprintf({{ end }}"{{ .PathTemplate }}"{{if shouldFormatPath }} {{ range .Params }}{{ if ne .Name "q" }}, {{ .Name }}{{ end }}{{ end }}){{ end }})
	req, err := c.buildDataRequest(ctx, http.Method{{ title .Method }}, u, {{ if ne .InputType.Type "" }}input{{ else }}http.NoBody{{ end }})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }}[ {{ end }}{{ if .ReturnsList }}[]{{ end }}*types.{{ .ResponseType.TypeName }}{{ if ne .ResponseType.GenericContainer "" }}]{{ end }}
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading {{ .ResponseType.TypeName }} creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}`

		imports = append(imports,
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if shouldFormatPath {
			imports = append(imports, "fmt")
		}
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
reqMods ...RequestModifier,
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

	u := c.BuildURL(ctx, nil, {{ if shouldFormatPath }}fmt.Sprintf({{ end }}"{{ .PathTemplate }}"{{if shouldFormatPath }} {{ range .Params }}{{ if ne .Name "q" }}, {{ .Name }}{{ end }}{{ end }}){{ end }})
	req, err := c.buildDataRequest(ctx, http.Method{{ title .Method }}, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	for _, mod := range reqMods {
		mod(req)
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
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if shouldFormatPath {
			imports = append(imports, "fmt")
		}
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
input *types.{{ .InputType.Type }},
reqMods ...RequestModifier,
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

	u := c.BuildURL(ctx, nil, {{ if shouldFormatPath }}fmt.Sprintf({{ end }}"{{ .PathTemplate }}"{{if shouldFormatPath }} {{ range .Params }}{{ if ne .Name "q" }}, {{ .Name }}{{ end }}{{ end }}){{ end }})
	req, err := c.buildDataRequest(ctx, http.Method{{ title .Method }}, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	for _, mod := range reqMods {
		mod(req)
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
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if shouldFormatPath {
			imports = append(imports, "fmt")
		}
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
{{ end -}} reqMods ...RequestModifier,
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

	u := c.BuildURL(ctx, nil, {{ if shouldFormatPath }}fmt.Sprintf({{ end }}"{{ .PathTemplate }}"{{if shouldFormatPath }} {{ range .Params }}{{ if ne .Name "q" }}, {{ .Name }}{{ end }}{{ end }}){{ end }})
	req, err := http.NewRequestWithContext(ctx, http.Method{{ title .Method }}, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a {{ .ResponseType.TypeName }}")
	}

	for _, mod := range reqMods {
		mod(req)
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
			"github.com/dinnerdonebetter/backend/internal/observability",
		)
		if shouldFormatPath {
			imports = append(imports, "fmt")
		}
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
			case typeString, typeBool, typeInt, typeUint64:
				return false
			default:
				return true
			}
		},
		"shouldFormatPath": func() bool {
			return shouldFormatPath
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
			case typeString:
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
	if err = t.Execute(&b, f); err != nil {
		return "", nil, err
	}

	return b.String(), imports, nil
}

func (f *APIClientFunction) RenderTest() (file string, imports []string, err error) {
	var tmpl string
	imports = []string{}

	isSearchOp := strings.Contains(f.Name, "TextSearch")

	const dummyTemplate = `

func TestClient_{{ .Name }}(T *testing.T) {
	T.Parallel()

	T.Run("TODO", func(t *testing.T) {
		t.Parallel()

		assert.NotEmpty(t, t.Name())
	})
}
`

	switch f.Method {
	case http.MethodGet:
		if f.QueryFiltered {
			// GET routes that return lists
			tmpl = `
func TestClient_{{ .Name }}(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "{{ .PathTemplate }}"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		list := fakes.BuildFake{{ pluralize (uppercaseFirstLetter .ResponseType.TypeName) }}sList()
		{{- if eq (uppercaseFirstLetter .ResponseType.TypeName) "User" }}
		for i := range list.Data {
			// the hashed passwords is never transmitted over the wire.
			list.Data[i].HashedPassword = ""
			// the two factor secret is transmitted over the wire only on creation.
			list.Data[i].TwoFactorSecret = ""
			// the two factor secret validation is never transmitted over the wire.
			list.Data[i].TwoFactorSecretVerifiedAt = nil
		}
		{{ else if eq (uppercaseFirstLetter .ResponseType.TypeName) "Household" }} 
		for i := range list.Data {
			list.Data[i].WebhookEncryptionKey = ""
		}
		{{ else if or (eq (uppercaseFirstLetter .ResponseType.TypeName) "HouseholdInvitation") }} 
		for i := range list.Data {
			list.Data[i].DestinationHousehold.WebhookEncryptionKey = ""
			list.Data[i].FromUser.TwoFactorSecret = ""
		}
		{{- end }}

		expected := &types.APIResponse[{{ if notNative .ResponseType.TypeName }}[]*types.{{ end }}{{ .ResponseType.TypeName }}]{
			Pagination: &list.Pagination,
			Data:       list.Data,
		}

		spec := newRequestSpec(true, http.Method{{ title .Method }},  {{ if isSearchOp }}fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q){{ else }}"limit=50&page=1&sortBy=asc"{{ end }}, expectedPathFormat, {{ range .Params }}{{ if ne .Name "q" }}{{.Name }},{{ end }}{{ end }})
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }}, {{ end }} nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	{{ range $i, $p := .Params }} {{ if ne $p.Name "q" }}
	T.Run("with empty {{ replace $p.Name "ID" ""  }} ID",  func(t *testing.T) {
		t.Parallel()

		{{ range $j, $p2 := $.Params}}{{ if ne $j $i}}{{ .Name }} := fakes.BuildFakeID(){{ end }}
		{{ end }}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.{{ $.Name }}(ctx, {{ range $j, $p2 := $.Params}}{{ if eq $j $i}}""{{ else }} {{ .Name }} {{ end }}, {{ end }} nil)

		require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual)
		assert.Error(t, err)
	}){{ end }}
{{ end }}

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }},{{ end }} nil)

		require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		spec := newRequestSpec(true, http.MethodGet, {{ if isSearchOp }}fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q){{ else }}"limit=50&page=1&sortBy=asc"{{ end }}, expectedPathFormat, {{ range .Params }}{{ if ne .Name "q" }}{{.Name }},{{ end }}{{ end }})
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }},{{ end }} nil)

		require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual)
		assert.Error(t, err)
	})
}`
			imports = append(imports,
				"testing",
				"context",
				"net/http",
				"github.com/stretchr/testify/assert",
				"github.com/stretchr/testify/require",
				"github.com/dinnerdonebetter/backend/pkg/types",
				"github.com/dinnerdonebetter/backend/pkg/types/fakes",
			)

			if isSearchOp {
				imports = append(imports, "fmt")
			}
		} else {
			// GET routes that don't return lists
			tmpl = `

func TestClient_{{ .Name }}(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "{{ .PathTemplate }}"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		data := fakes.BuildFake{{ uppercaseFirstLetter .ResponseType.TypeName }}()
		{{- if eq (uppercaseFirstLetter .ResponseType.TypeName) "User" }}
		// the hashed passwords is never transmitted over the wire.
		data.HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		data.TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		data.TwoFactorSecretVerifiedAt = nil
		{{ else if eq .Name "FetchUserDataReport" }} 
		data.User.TwoFactorSecret = ""
		data.User.HashedPassword = ""
		data.User.TwoFactorSecretVerifiedAt = nil
		for i := range data.Households {
			data.Households[i].WebhookEncryptionKey = ""
		}
		for i := range data.SentInvites {
			data.SentInvites[i].DestinationHousehold.WebhookEncryptionKey = ""
			data.SentInvites[i].FromUser.TwoFactorSecret = ""
			data.SentInvites[i].FromUser.HashedPassword = ""
			data.SentInvites[i].FromUser.TwoFactorSecretVerifiedAt = nil
		}
		for i := range data.ReceivedInvites {
			data.ReceivedInvites[i].DestinationHousehold.WebhookEncryptionKey = ""
			data.ReceivedInvites[i].FromUser.TwoFactorSecret = ""
			data.ReceivedInvites[i].FromUser.HashedPassword = ""
			data.ReceivedInvites[i].FromUser.TwoFactorSecretVerifiedAt = nil
		}
		{{ else if eq (uppercaseFirstLetter .ResponseType.TypeName) "Household" }} 
		data.WebhookEncryptionKey = ""
		{{ else if or (eq (uppercaseFirstLetter .ResponseType.TypeName) "HouseholdInvitation") }} 
		data.DestinationHousehold.WebhookEncryptionKey = ""
		data.FromUser.TwoFactorSecret = ""
		{{- end }}
		expected := &types.APIResponse[{{ if notNative .ResponseType.TypeName }}*types.{{ end }}{{ .ResponseType.TypeName }}]{
			Data: data,
		}

		spec := newRequestSpec(true, http.Method{{ title .Method }}, "", expectedPathFormat, {{ range .Params }}{{.Name }},{{ end }})
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }}, {{ end }})

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	{{ range $i, $p := .Params }}
	T.Run("with invalid {{ replace $p.Name "ID" "" }} ID",  func(t *testing.T) {
		t.Parallel()

		{{ range $j, $p2 := $.Params}}{{ if ne $j $i}}{{ .Name }} := fakes.BuildFakeID(){{ end }}
		{{ end }}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.{{ $.Name }}(ctx, {{ range $j, $p2 := $.Params}}{{ if eq $j $i}}""{{ else }} {{ .Name }} {{ end }}, {{ end }})

		require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual)
		assert.Error(t, err)
	})
{{ end }}

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }},{{ end }})

		require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, {{ range .Params }}{{.Name }},{{ end }})
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }},{{ end }})

		require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual)
		assert.Error(t, err)
	})
}
`
			imports = append(imports,
				"testing",
				"context",
				"net/http",
				"github.com/stretchr/testify/assert",
				"github.com/stretchr/testify/require",
				"github.com/dinnerdonebetter/backend/pkg/types",
				"github.com/dinnerdonebetter/backend/pkg/types/fakes",
			)
		}

	case http.MethodPost, http.MethodPut, http.MethodPatch:
		tmpl = `

func TestClient_{{ .Name }}(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "{{ .PathTemplate }}"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		data := {{ if .ReturnsList }}[]*types.{{ .ResponseType.TypeName }} { {{ end }}fakes.BuildFake{{ uppercaseFirstLetter .ResponseType.TypeName }}(){{ if .ReturnsList }}}{{ end }}
		{{- if eq (uppercaseFirstLetter .ResponseType.TypeName) "User" }}
		// the hashed passwords is never transmitted over the wire.
		data.HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		data.TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		data.TwoFactorSecretVerifiedAt = nil
		{{ else if eq (uppercaseFirstLetter .ResponseType.TypeName) "Household" }} 
		data.WebhookEncryptionKey = ""
		{{ else if or (eq (uppercaseFirstLetter .ResponseType.TypeName) "HouseholdInvitation") }} 
		data.DestinationHousehold.WebhookEncryptionKey = ""
		data.FromUser.TwoFactorSecret = ""
		{{- end }}
		expected := &types.APIResponse[{{ if notNative .ResponseType.TypeName }}{{ if .ReturnsList }}[]{{ end }}*types.{{ end }}{{ .ResponseType.TypeName }}]{
			Data: data,
		}

		{{ if ne .InputType.Type "" }}exampleInput := fakes.BuildFake{{ .InputType.Type }}({{ if eq .Name "VerifyTOTPSecret" }}data{{ end }}){{ end }}

		spec := newRequestSpec(false, http.Method{{ title .Method }}, "", expectedPathFormat, {{ range .Params }}{{.Name }},{{ end }})
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		{{ if (and (ne .Method "PUT") (ne .Method "PATCH")) }}actual,{{ end }} err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }}, {{ end }} {{ if ne .InputType.Type "" }} exampleInput {{ end }})

		{{ if (and (ne .Method "PUT") (ne .Method "PATCH")) }}require.NotNil(t, actual){{ end }}
		assert.NoError(t, err)
		{{ if (and (ne .Method "PUT") (ne .Method "PATCH")) }}assert.Equal(t, data, actual){{ end }}
	})

	{{ range $i, $p := .Params }}
	T.Run("with invalid {{ replace $p.Name "ID" "" }} ID",  func(t *testing.T) {
		t.Parallel()

		{{ range $j, $p2 := $.Params}}{{ if ne $j $i}}{{ .Name }} := fakes.BuildFakeID(){{ end }}
		{{ end }}

		{{ if eq $.Name "VerifyTOTPSecret" }}data := fakes.BuildFakeUser(){{ end }}
		{{ if ne $.InputType.Type "" }}exampleInput := fakes.BuildFake{{ $.InputType.Type }}({{ if eq .Name "VerifyTOTPSecret" }}data{{ end }}){{ end }}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		{{ if (and (ne $.Method "PUT") (ne $.Method "PATCH")) }}actual,{{ end }} err := c.{{ $.Name }}(ctx, {{ range $j, $p2 := $.Params}}{{ if eq $j $i}}""{{ else }} {{ .Name }} {{ end }}, {{ end }} {{ if ne $.InputType.Type "" }} exampleInput {{ end }})

		{{ if (and (ne $.Method "PUT") (ne $.Method "PATCH")) }}require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual){{ end }}
		assert.Error(t, err)
	})
{{ end }}

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		{{ if eq .Name "VerifyTOTPSecret" }}data := fakes.BuildFakeUser(){{ end }}
		{{ if ne .InputType.Type "" }}exampleInput := fakes.BuildFake{{ .InputType.Type }}({{ if eq .Name "VerifyTOTPSecret" }}data{{ end }}){{ end }}

		c := buildTestClientWithInvalidURL(t)
		{{ if (and (ne .Method "PUT") (ne .Method "PATCH")) }}actual,{{ end }} err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }},{{ end }} {{ if ne .InputType.Type "" }} exampleInput {{ end }})

		{{ if (and (ne .Method "PUT") (ne .Method "PATCH")) }}require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual){{ end }}
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		{{ range .Params }}{{ .Name }} := fakes.BuildFakeID()
		{{ end }}

		{{ if eq .Name "VerifyTOTPSecret" }}data := fakes.BuildFakeUser(){{ end }}
		{{ if ne .InputType.Type "" }}exampleInput := fakes.BuildFake{{ .InputType.Type }}({{ if eq .Name "VerifyTOTPSecret" }}data{{ end }}){{ end }}

		spec := newRequestSpec(false, http.Method{{ title .Method }}, "", expectedPathFormat, {{ range .Params }}{{.Name }},{{ end }})
		c := buildTestClientWithInvalidResponse(t, spec)
		{{ if (and (ne .Method "PUT") (ne .Method "PATCH")) }}actual,{{ end }} err := c.{{ .Name }}(ctx, {{ range .Params }}{{.Name }},{{ end }}{{ if ne .InputType.Type "" }} exampleInput {{ end }})

		{{ if (and (ne .Method "PUT") (ne .Method "PATCH")) }}require.{{ if notNative $.ResponseType.TypeName }}Nil{{ else }} {{ negativeAssertFunc $.ResponseType.TypeName }} {{ end }}(t, actual){{ end }}
		assert.Error(t, err)
	})
}
`
		imports = append(imports,
			"testing",
			"context",
			"net/http",
			"github.com/stretchr/testify/assert",
			"github.com/dinnerdonebetter/backend/pkg/types",
			"github.com/dinnerdonebetter/backend/pkg/types/fakes",
		)

		if f.Method != http.MethodPut && f.Method != http.MethodPatch {
			imports = append(imports, "github.com/stretchr/testify/require")
		}

	case http.MethodDelete:
		tmpl = dummyTemplate
		imports = append(imports,
			"testing",
			"github.com/stretchr/testify/assert",
		)
	}

	if tmpl == "" {
		return "", nil, nil
	}

	t := template.Must(template.New(f.Name).Funcs(map[string]any{
		"lowercase": strings.ToLower,
		"contains":  strings.Contains,
		"title": func(s string) string {
			return uppercaseFirstLetter(strings.ToLower(s))
		},
		"replace":              strings.ReplaceAll,
		"uppercaseFirstLetter": uppercaseFirstLetter,
		"pluralize": func(s string) string {
			switch s {
			case "AuditLogEntry":
				return "AuditLogEntrie"
			default:
				return s
			}
		},
		"notNative": func(s string) bool {
			switch s {
			case typeString, typeBool, typeInt, typeUint64:
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
			case typeString:
				return `""`
			default:
				panic(fmt.Sprintf("aaaaaaaaaaaaaaaa bad type: %s", s))
			}
		},
		"negativeAssertFunc": func(s string) string {
			switch s {
			case typeString:
				return `Empty`
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
		"isSearchOp": func() bool {
			return isSearchOp
		},
	}).Parse(tmpl))

	var b bytes.Buffer
	if err = t.Execute(&b, f); err != nil {
		return "", nil, err
	}

	return b.String(), imports, nil
}
