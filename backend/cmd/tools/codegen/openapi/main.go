package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen/openapi/golang"
	"github.com/dinnerdonebetter/backend/cmd/tools/codegen/openapi/typescript"

	"github.com/spf13/pflag"
	"github.com/swaggest/openapi-go/openapi31"
)

var (
	generateTypescript = pflag.BoolP("typescript", "", false, "generate typescript code")
	generateGolang     = pflag.BoolP("golang", "", false, "generate go code")
)

const (
	golangAPIClientOutputPath     = "pkg/apiclient"
	specFilepath                  = "../openapi_spec.yaml"
	typescriptModelsOutputPath    = "../frontend/packages/models"
	typescriptAPIClientOutputPath = "../frontend/packages/api-client"
	typescriptMockAPIOutputPath   = "../frontend/packages/mock-playwright-api"
)

func loadSpec(filepath string) (*openapi31.Spec, error) {
	specBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading spec file: %w", err)
	}

	spec := &openapi31.Spec{}
	if err = spec.UnmarshalYAML(specBytes); err != nil {
		return nil, fmt.Errorf("unmarshaling spec: %w", err)
	}

	return spec, nil
}

func writeTypescriptFiles(spec *openapi31.Spec) error {
	errPrefix := "failed to write typescript"
	if err := typescript.WriteModelFiles(spec, typescriptModelsOutputPath); err != nil {
		return fmt.Errorf("%s model files: %w", errPrefix, err)
	}

	if err := typescript.WriteAPIClientFiles(spec, typescriptAPIClientOutputPath); err != nil {
		return fmt.Errorf("%s API client files: %w", errPrefix, err)
	}

	if err := typescript.WriteMockAPIFiles(spec, typescriptMockAPIOutputPath); err != nil {
		return fmt.Errorf("%s mock API files: %w", errPrefix, err)
	}

	return nil
}

func writeGoFiles(outPath string, spec *openapi31.Spec) error {
	errPrefix := "failed to write golang"

	if err := os.RemoveAll(fmt.Sprintf("%s/*.go", outPath)); err != nil {
		return fmt.Errorf("removing go files: %w", err)
	}

	if err := golang.WriteAPIClientFiles(spec, outPath); err != nil {
		return fmt.Errorf("%s API client files: %w", errPrefix, err)
	}

	if err := golang.WriteAPITypesFiles(spec, outPath); err != nil {
		return fmt.Errorf("%s type files: %w", errPrefix, err)
	}

	return nil
}

func main() {
	pflag.Parse()

	spec, err := loadSpec(specFilepath)
	if err != nil {
		log.Fatalf("failed to load spec: %v", err)
	}

	var wg sync.WaitGroup

	if *generateTypescript {
		wg.Add(1)
		go func() {
			if err = writeTypescriptFiles(spec); err != nil {
				log.Fatalf("failed to write typescript files: %v", err)
			}
			wg.Done()
		}()
	}

	if *generateGolang {
		wg.Add(1)
		go func() {
			if err = writeGoFiles(golangAPIClientOutputPath, spec); err != nil {
				log.Fatalf("failed to write typescript files: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
