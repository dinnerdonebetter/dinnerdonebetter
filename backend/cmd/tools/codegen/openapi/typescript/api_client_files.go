package typescript

import (
	"bytes"
	"net/http"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/swaggest/openapi-go/openapi31"
)

const (
	stringType = "string"
	nullType   = "null"
)

var skipOps = map[string]bool{
	"CheckForLiveness":  true,
	"CheckForReadiness": true,
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

func buildFunction(path, method string, op *openapi31.Operation) *APIClientFunction {
	functionName := lowercaseFirstLetter(*op.ID)
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
		Description:       pointer.Dereference(op.Description),
		ResponseType:      rt,
		PathTemplate:      strings.ReplaceAll(path, "{", "${"),
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
	Description,
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

func (f *APIClientFunction) Render() (string, error) {
	var tmpl string

	switch f.Method {
	case http.MethodGet, http.MethodDelete:
		if f.QueryFiltered {
			// GET routes that return lists

			tmpl = `async {{ .Name }}(
  {{ range .Params }}{{ .Name }}: {{ .Type }}{{ if ne .DefaultValue "" }} = {{ .DefaultValue }}{{ end }},
{{ end -}}
  filter: QueryFilter = QueryFilter.Default(),
): Promise< QueryFilteredResult< {{ .ResponseType.TypeName }} > > {
  let self = this;
  return new Promise(async function (resolve, reject) {
  {{ range .Params }}if ({{ .Name }}.trim() === '') {
        throw new Error('{{ .Name }} is required');
      }

{{ end -}}

    self.client.{{ lowercase .Method }}< {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} >(` + "`" + "{{ .PathTemplate }}" + "`" + `, {
      params: filter.asRecord(),
    })
 		.then((res: AxiosResponse<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}>) => {
		  if (res.data.error) {
            reject(res.data.error);
          } else {
			resolve(new QueryFilteredResult<{{ .ResponseType.TypeName }}>({
			  data: res.data.data,
			  totalCount: res.data.pagination?.totalCount,
			  page: res.data.pagination?.page,
			  limit: res.data.pagination?.limit,
		    }));
		  }          
        })
        .catch((error: AxiosError<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}>) => {
          if (error?.response?.data?.error) {
            reject(error?.response?.data?.error);
          } else {
		    reject({message: error?.message || 'unknown error'});
		  }
        });
  });
}`
		} else {
			// GET routes that don't return lists
			tmpl = `async {{ .Name }}(
  {{ range .Params }}{{ .Name }}: {{ .Type }}{{ if ne .DefaultValue "" }} = {{ .DefaultValue }}{{ end }},
	{{ end -}}): Promise<  {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }} {{ .ResponseType.TypeName }} >  {{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}  {
  let self = this;
  return new Promise(async function (resolve, reject) {
  {{ range .Params }}if ({{ .Name }}.trim() === '') {
        throw new Error('{{ .Name }} is required');
      }

{{ end -}}

    self.client.{{ lowercase .Method }}< {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} >(` + "`" + "{{ .PathTemplate }}" + "`" + `)
 		.then((res: AxiosResponse<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}>) => {
          if (res.data.error) {
            reject(res.data.error);
          } else {
			  resolve(res.data);
		  }
        })
        .catch((error: AxiosError<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}>) => {
          if (error?.response?.data?.error) {
            reject(error?.response?.data?.error);
          } else {
		    reject({message: error?.message || 'unknown error'});
		  }
        });
  });
}`
		}

	case http.MethodPost, http.MethodPut, http.MethodPatch:
		tmpl = `async {{ .Name }}(
  {{ range .Params }}{{ .Name }}: {{ .Type }},{{ end }}
  {{ if ne .InputType.Type "" }}input: {{ .InputType.Type }},{{ end }}
): Promise< {{ if .ReturnRawResponse }} AxiosResponse< {{ end }} {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}  {{ .ResponseType.TypeName }} > {{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }}  {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}{{ if .ReturnRawResponse }} > {{ end }} {
  let self = this;
  return new Promise(async function (resolve, reject) {
  {{ range .Params }}if ({{ .Name }}.trim() === '') {
        throw new Error('{{ .Name }} is required');
      }

{{ end -}}
    self.client.{{ lowercase .Method }}<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} >(` + "`" + "{{ .PathTemplate }}" + "`" + `{{ if ne .InputType.Type "" }}, input{{ end }})
 		.then((res: AxiosResponse<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}>) => {
          if (res.data.error{{ if or (eq .Name "loginForToken") (eq .Name "adminLoginForToken") }} && res.data.error.message.toLowerCase() != "totp required" {{ end }}) {
            reject(res.data.error);
          } else {
	    	resolve(res{{ if not .ReturnRawResponse }}.data{{ end }});
 		  }
        })
        .catch((error: AxiosError<{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}>) => {
          if (error?.response?.data?.error) {
            reject(error?.response?.data?.error);
          } else {
		    reject({message: error?.message || 'unknown error'});
		  }
        });
	  });
}

`
	}

	if tmpl == "" {
		panic("Unknown template")
	}

	t := template.Must(template.New("function").Funcs(map[string]any{
		"lowercase": strings.ToLower,
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, f); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (f *APIClientFunction) RenderTest() (string, error) {
	var tmpl string

	switch f.Method {
	case http.MethodGet, http.MethodDelete:
		if f.QueryFiltered {
			// GET routes that return lists
			tmpl = `
			   it("should {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
{{ end }}
		
			       const exampleResponse = new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}({
			           details: {
			             currentHouseholdID: "test",
			             traceID: "test",
			           },
		               pagination: QueryFilter.Default().toPagination(),
					   data: [
					     new {{ .ResponseType.TypeName }}(),
			           ],
			       });
			       mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply({{ .DefaultStatusCode }}, exampleResponse);
		
			       client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }}).then((response: QueryFilteredResult< {{ .ResponseType.TypeName }}>) => {
					  expect(response.data).toEqual(exampleResponse.data);
					  expect(response.page).toEqual(exampleResponse.pagination?.page);
					  expect(response.limit).toEqual(exampleResponse.pagination?.limit);
			       })
				  .then(() => {
					  expect(mock.history.{{ lowercase .Method }}.length).toBe(1);
					  expect(mock.history.{{ lowercase .Method }}[0].headers).toHaveProperty('Authorization')
					  expect((mock.history.{{ lowercase .Method }}[0].headers || {})['Authorization']).toBe(` + "`" + `Bearer test-token` + "`" + `)
					});
			   });

			   it("should raise errors appropriately when trying to {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
					{{ end }}
			   
	              const expectedError = buildObligatoryError('{{ .Name }} user error');
			      const exampleResponse = new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}(expectedError);
			      mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply({{ .DefaultStatusCode }}, exampleResponse);
			   
			      expect(client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }})).rejects.toEqual(expectedError.error);
			   });

			   it("should raise service errors appropriately when trying to {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
					{{ end }}
			   
                  const expectedError = buildObligatoryError('{{ .Name }} service error');
			  	  const exampleResponse = new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}(expectedError);
			      mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply(500, exampleResponse);
			   
			      expect(client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }})).rejects.toEqual(expectedError.error);
			   });
			`
		} else {
			// GET routes that don't return lists
			tmpl = `
			   it("should {{ testDescription .Description }}", () => {
                   {{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
				   {{ end }}
		
			       const exampleResponse = {{ if not (isNative .ResponseType.TypeName) }}new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} {{ if isNative .ResponseType.TypeName }}({
			           {{ if hasParam "id" }}id: faker.string.uuid(), {{ end }}
			       }){{ end }}{{ else }} {{ defaultValue .ResponseType.TypeName }}{{ end }};
			       mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply({{ .DefaultStatusCode }}, exampleResponse);
		
			       client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }}).then((response: {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}) => {
			           expect(response).toEqual(exampleResponse);
			       })
				  .then(() => {
					  expect(mock.history.{{ lowercase .Method }}.length).toBe(1);
					  expect(mock.history.{{ lowercase .Method }}[0].headers).toHaveProperty('Authorization')
					  expect((mock.history.{{ lowercase .Method }}[0].headers || {})['Authorization']).toBe(` + "`" + `Bearer test-token` + "`" + `)
					});
			   })

			   it("should raise errors appropriately when trying to {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
		 		{{ end }}
			   
			      const expectedError = buildObligatoryError('{{ .Name }} user error');
			  	  const exampleResponse = {{ if not (isNative .ResponseType.TypeName) }}new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} {{ if not (isNative .ResponseType.TypeName) }}(expectedError){{ end }}{{ else }}(expectedError){{ end }};
					mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply({{ .DefaultStatusCode }}, exampleResponse);
			   
			    expect(client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }})).rejects.toEqual(expectedError.error);
			   }) 

			   it("should raise service errors appropriately when trying to {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
		 		{{ end }}
			   
			      const expectedError = buildObligatoryError('{{ .Name }} service error');
			  	  const exampleResponse = {{ if not (isNative .ResponseType.TypeName) }}new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} {{ if not (isNative .ResponseType.TypeName) }}(expectedError){{ end }}{{ else }}(expectedError){{ end }};
					mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply(500, exampleResponse);
			   
			    expect(client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }})).rejects.toEqual(expectedError.error);
			   })
			`
		}

	case http.MethodPost, http.MethodPut, http.MethodPatch:
		tmpl = `
			   it("should {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
{{ end }}
			
    				{{ if ne .InputType.Type "" }}const exampleInput = new {{ .InputType.Type }}();{{ end }}

			       const exampleResponse = {{ if not (isNative .ResponseType.TypeName) }}new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} {{ if isNative .ResponseType.TypeName }}({
			           {{ if hasParam "id" }}id: faker.string.uuid(), {{ end }}
			       }){{ end }}{{ else }} {{ defaultValue .ResponseType.TypeName }}{{ end }};
			       mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply({{ .DefaultStatusCode }}, exampleResponse);
			
			       client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }}{{ if ne .InputType.Type "" }}exampleInput{{ end }}).then((response: {{ if .ReturnRawResponse }}AxiosResponse<{{ end }}{{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }}{{ if .ReturnRawResponse }}>{{ end }}) => {
			           expect(response{{ if .ReturnRawResponse}}.data{{ end }}).toEqual(exampleResponse);
			       })
				  .then(() => {
					  expect(mock.history.{{ lowercase .Method }}.length).toBe(1);
					  {{ if ne .InputType.Type "" }}expect(mock.history.{{ lowercase .Method }}[0].data).toBe(JSON.stringify(exampleInput));{{ end }}
					  expect(mock.history.{{ lowercase .Method }}[0].headers).toHaveProperty('Authorization')
					  expect((mock.history.{{ lowercase .Method }}[0].headers || {})['Authorization']).toBe(` + "`" + `Bearer test-token` + "`" + `)
					});
			   })

			   it("should appropriately raise errors when they occur during {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
{{ end }}
			   
    				{{ if ne .InputType.Type "" }}const exampleInput = new {{ .InputType.Type }}();{{ end }}
			   
			     const expectedError = buildObligatoryError('{{ .Name }} user error');
			  	   const exampleResponse = {{ if not (isNative .ResponseType.TypeName) }}new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} {{ if not (isNative .ResponseType.TypeName) }}(expectedError){{ end }}{{ else }} {{ defaultValue .ResponseType.TypeName }}{{ end }};
			      mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply({{ .DefaultStatusCode }}, exampleResponse);
			   
			      expect(client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }}{{ if ne .InputType.Type "" }}exampleInput{{ end }})).rejects.toEqual(expectedError.error);
			   })

			   it("should appropriately raise service errors when they occur during {{ testDescription .Description }}", () => {
					{{ range .Params }}let {{ .Name }} = {{ fakeFor .Name .Type }};
{{ end }}
			   
    				{{ if ne .InputType.Type "" }}const exampleInput = new {{ .InputType.Type }}();{{ end }}
			   
			    const expectedError = buildObligatoryError('{{ .Name }} service error');
			  	   const exampleResponse = {{ if not (isNative .ResponseType.TypeName) }}new {{ if ne .ResponseType.GenericContainer "" }}{{ .ResponseType.GenericContainer }} < {{ end }}{{ if or .ReturnsList .ResponseType.IsArray }}Array<{{ end }}{{ .ResponseType.TypeName }}{{ if or .ReturnsList .ResponseType.IsArray }}>{{ end }} {{ if ne .ResponseType.GenericContainer "" }} > {{ end }} {{ if not (isNative .ResponseType.TypeName) }}(expectedError){{ end }}{{ else }} {{ defaultValue .ResponseType.TypeName }}{{ end }};
			      mock.on{{ upperFirst .Method  }}(` + "`" + `${baseURL}{{ .PathTemplate }}` + "`" + `).reply(500, exampleResponse);
			   
			      expect(client.{{ .Name }}({{ range .Params }}{{ .Name }},{{ end }}{{ if ne .InputType.Type "" }}exampleInput{{ end }})).rejects.toEqual(expectedError.error);
			   })
			`
	}

	if tmpl == "" {
		panic("Unknown template")
	}

	t := template.Must(template.New("typescript function test").Funcs(map[string]any{
		"fakeFor": func(name, typ string) string {
			switch typ {
			case stringType:
				return "fakeID()"
			default:
				return nullType
			}
		},
		"isNative": func(s string) bool {
			switch s {
			case stringType:
				return true
			default:
				return false
			}
		},
		"defaultValue": func(s string) string {
			switch s {
			case stringType:
				return `""`
			default:
				return ""
			}
		},
		"hasParam": func(name string) bool {
			for _, param := range f.Params {
				if strings.EqualFold(param.Name, name) {
					return true
				}
			}

			return false
		},
		"upperFirst": func(s string) string {
			return uppercaseFirstLetter(strings.ToLower(s))
		},
		"testDescription": func(s string) string {
			return strings.NewReplacer(
				"Creates", "create",
				"Updates", "update",
				"Fetches", "fetch",
				"Gets", "get",
				"Archives", "archive",
				"Deletes", "delete",
				"Cancels", "cancel",
				"Rejects", "reject",
			).Replace(s)
		},
		"lowercase": strings.ToLower,
	}).Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, f); err != nil {
		return "", err
	}

	return b.String(), nil
}
