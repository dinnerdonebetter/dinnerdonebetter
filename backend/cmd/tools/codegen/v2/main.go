package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen/v2/typescript"

	"github.com/swaggest/openapi-go/openapi31"
)

const (
	specFilepath                  = "../openapi_spec.yaml"
	typescriptAPIClientOutputPath = "../frontend/packages/generated-client"
	typescriptModelsOutputPath    = "../frontend/packages/generated-models"
)

func purgeTypescriptFiles(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".ts" {
			if err = os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	})
}

func loadSpec(filepath string) (*openapi31.Spec, error) {
	specBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading spec file: %w", err)
	}

	spec := &openapi31.Spec{}
	if err = spec.UnmarshalYAML(specBytes); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
	}

	return spec, nil
}

func getOpCountForSpec(spec *openapi31.Spec) uint {
	var output uint

	for _, path := range spec.Paths.MapOfPathItemValues {
		if path.Get != nil {
			output++
		}
		if path.Put != nil {
			output++
		}
		if path.Patch != nil {
			output++
		}
		if path.Post != nil {
			output++
		}
		if path.Delete != nil {
			output++
		}
		if path.Head != nil {
			output++
		}
	}

	return output
}

func removeDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if slices.Contains(list, item) == false {
			list = append(list, item)
		}
	}
	return list
}

func writeTypescriptAPIClientFiles(spec *openapi31.Spec) error {
	typescriptClientFiles, err := typescript.GenerateClientFiles(spec)
	if err != nil {
		return fmt.Errorf("failed to generate typescript files: %w", err)
	}

	if err = os.MkdirAll(typescriptAPIClientOutputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = purgeTypescriptFiles(typescriptAPIClientOutputPath); err != nil {
		return fmt.Errorf("failed to purge typescript files: %w", err)
	}

	createdFunctions := []string{}
	imports := map[string][]string{
		"@dinnerdonebetter/models": {
			"QueryFilter",
			"QueryFilteredResult",
		},
	}
	for _, function := range typescriptClientFiles {
		fileContents, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		if function.InputType.Type != "" {
			imports["@dinnerdonebetter/models"] = append(imports["@dinnerdonebetter/models"], function.InputType.Type)
		}

		if function.ResponseType.TypeName != "" && function.ResponseType.TypeName != "string" {
			imports["@dinnerdonebetter/models"] = append(imports["@dinnerdonebetter/models"], function.ResponseType.TypeName)
		}

		if function.ResponseType.GenericContainer != "" {
			imports["@dinnerdonebetter/models"] = append(imports["@dinnerdonebetter/models"], function.ResponseType.GenericContainer)
		}

		createdFunctions = append(createdFunctions, fileContents)
	}

	createdFunctions = removeDuplicates(createdFunctions)
	slices.Sort(createdFunctions)

	modelsImports := imports["@dinnerdonebetter/models"]
	modelsImports = removeDuplicates(modelsImports)
	slices.Sort(modelsImports)

	importBlock := `import axios, {
  AxiosInstance,
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  InternalAxiosRequestConfig,
} from 'axios';
import { Span } from '@opentelemetry/api';

import { buildServerSideLogger, LoggerType } from '@dinnerdonebetter/logger';
` + "import {\n\t" + strings.Join(modelsImports, ",\n\t") + "\n" + `} from "@dinnerdonebetter/models";`

	indexFile := typescript.GeneratedDisclaimer + "\n\n" + importBlock + `

function _curlFromAxiosConfig(config: InternalAxiosRequestConfig): string {
  const method = (config?.method || 'UNKNOWN').toUpperCase();
  const url = config.url;
  const headers = config.headers || {};
  const data = config.data;

  ['get', 'delete', 'head', 'post', 'put', 'patch'].forEach((method) => {
    delete headers[method];
  });

  // iterate through headers["common"], and add each key's value to headers
  const headerDefault = headers as unknown as HeadersDefaults;
  for (const key in headerDefault['common']) {
    if (headerDefault['common'].hasOwnProperty(key)) {
      headers[key] = headerDefault['common'][key];
    }
  }
  delete headers['common'];

  let curlCommand = ` + "`" + `curl -X ${method} "${config?.baseURL || 'MISSING_BASE_URL'}${url}"` + "`" + `;

  for (const key in headers) {
    if (headers.hasOwnProperty(key)) {
      curlCommand += ` + "`" + ` -H "${key}: ${headers[key]}"` + "`" + `;
    }
  }

  if (data) {
    curlCommand += ` + "`" + ` -d '${JSON.stringify(data)}'` + "`" + `;
  }

  return curlCommand;
}

export class DinnerDoneBetterAPIClient {
  baseURL: string;
  client: AxiosInstance;
  requestInterceptorID: number;
  responseInterceptorID: number;
  logger: LoggerType = buildServerSideLogger('api_client');

  constructor(clientName: string = 'DDB-Service-Client', baseURL: string = '', oauth2Token?: string) {
    this.baseURL = baseURL;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Request-Source': 'webapp',
      'X-Service-Client': clientName,
    };

    if (oauth2Token) {
      headers['Authorization'] = ` + "`" + `Bearer ${oauth2Token}` + "`" + `;
    }

    this.client = axios.create({
      baseURL,
      timeout: 10000,
      withCredentials: true,
      crossDomain: true,
      headers,
    } as AxiosRequestConfig);

    this.requestInterceptorID = this.client.interceptors.request.use(
      (request: InternalAxiosRequestConfig) => {
        // this.logger.debug(` + "`" + `Request: ${request.method} ${request.baseURL}${request.url}` + "`" + `);
        console.log(` + "`" + `${_curlFromAxiosConfig(request)}` + "`" + `);

        return request;
      },
      (error) => {
        // Do whatever you want with the response error here
        // But, be SURE to return the rejected promise, so the caller still has
        // the option of additional specialized handling at the call-site:
        return Promise.reject(error);
      },
    );

    this.responseInterceptorID = this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        this.logger.debug(
          ` + "`" + `Response: ${response.status} ${response.config.method} ${response.config.url}` + "`" + `,
          // response.data,
        );

        // console.log(` + "`" + `${response.status} ${_curlFromAxiosConfig(response.config)}` + "`" + `);

        return response;
      },
      (error) => {
        return Promise.reject(error);
      },
    );
  }

  withSpan(span: Span): DinnerDoneBetterAPIClient {
    const spanContext = span.spanContext();
    const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

    this.client.interceptors.request.eject(this.requestInterceptorID);
    this.requestInterceptorID = this.client.interceptors.request.use(
      (request: InternalAxiosRequestConfig) => {
        this.logger.debug(` + "`" + `Request: ${request.method} ${request.url}` + "`" + `, spanLogDetails);

        // console.log(_curlFromAxiosConfig(request));

        if (spanContext.traceId) {
          request.headers.set('traceparent', spanContext.traceId);
        }

        return request;
      },
      (error) => {
        return Promise.reject(error);
      },
    );

    return this;
  }

  // eslint-disable-next-line no-unused-vars
  configureRouterRejectionInterceptor(redirectCallback: (_: Location) => void) {
    this.client.interceptors.response.eject(this.responseInterceptorID);
    this.responseInterceptorID = this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        this.logger.debug(
          ` + "`" + `Response: ${response.status} ${response.config.method} ${response.config.url}${response.config.method === 'POST' || response.config.method === 'PUT' ? ` + "`" + ` ${JSON.stringify(response.config.data)}` + "`" + ` : ''}` + "`" + `,
        );

        return response;
      },
      (error: AxiosError) => {
        console.debug(` + "`" + `Request failed: ${error.response?.status}` + "`" + `);
        if (error.response?.status === 401) {
          redirectCallback(window.location);
        }

        return Promise.reject(error);
      },
    );
  }

`
	for _, createdFile := range createdFunctions {
		indexFile += createdFile + "\n\n"
	}

	indexFile += "\n}\n"

	if err = os.WriteFile(fmt.Sprintf("%s/index.ts", typescriptAPIClientOutputPath), []byte(indexFile), 0o644); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	return nil
}

func writeTypescriptModelFiles(spec *openapi31.Spec) error {
	// first do enums
	enums, err := typescript.GenerateEnumDefinitions(spec)
	if err != nil {
		return fmt.Errorf("could not generate enums file: %w", err)
	}

	enumsFile := typescript.GeneratedDisclaimer + "\n\n"
	for _, enum := range enums {
		def, renderErr := enum.Render()
		if renderErr != nil {
			return fmt.Errorf("could not render enum definition: %w", renderErr)
		}

		enumsFile += fmt.Sprintf("%s\n\n", def)
	}

	// next do models

	typescriptModelFiles, err := typescript.GenerateModelFiles(spec)
	if err != nil {
		return fmt.Errorf("failed to generate typescript models files: %w", err)
	}

	if err = os.MkdirAll(typescriptModelsOutputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = purgeTypescriptFiles(typescriptModelsOutputPath); err != nil {
		return fmt.Errorf("failed to purge typescript models files: %w", err)
	}

	createdFiles := []string{}
	for filename, function := range typescriptModelFiles {
		actualFilepath := fmt.Sprintf("%s/%s.ts", typescriptModelsOutputPath, filename)

		rawFileContents, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		fileContents := fmt.Sprintf("%s\n\n%s", typescript.GeneratedDisclaimer, rawFileContents)
		if err = os.WriteFile(actualFilepath, []byte(fileContents), 0o0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", actualFilepath, err)
		}

		createdFiles = append(createdFiles, actualFilepath)
	}

	slices.Sort(createdFiles)

	fmt.Printf("Wrote %d files, had %d Operations\n", len(createdFiles), getOpCountForSpec(spec))

	indexFile := fmt.Sprintf("%s\n\n", typescript.GeneratedDisclaimer)
	for _, createdFile := range createdFiles {
		indexFile += fmt.Sprintf("export * from './%s';\n", strings.TrimSuffix(strings.TrimPrefix(createdFile, fmt.Sprintf("%s/", typescriptModelsOutputPath)), ".ts"))
	}

	if err = os.WriteFile(fmt.Sprintf("%s/enums.ts", typescriptModelsOutputPath), []byte(enumsFile), 0o644); err != nil {
		return fmt.Errorf("failed to write enums file: %w", err)
	}

	for name := range typescript.StaticFiles {
		indexFile += fmt.Sprintf("export * from './%s';\n", name)
	}
	indexFile += "export * from './enums';\n"

	if err = os.WriteFile(fmt.Sprintf("%s/index.ts", typescriptModelsOutputPath), []byte(indexFile), 0o644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	for name, content := range typescript.StaticFiles {
		if err = os.WriteFile(fmt.Sprintf("%s/%s.ts", typescriptModelsOutputPath, name), []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	return nil
}

func main() {
	spec, err := loadSpec(specFilepath)
	if err != nil {
		log.Fatalf("failed to load spec: %v", err)
	}

	if err = writeTypescriptModelFiles(spec); err != nil {
		log.Fatalf("failed to write typescript model files: %v", err)
	}

	if err = writeTypescriptAPIClientFiles(spec); err != nil {
		log.Fatalf("failed to write typescript API client files: %v", err)
	}
}
