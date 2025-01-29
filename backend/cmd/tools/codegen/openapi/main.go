package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func writeGoFiles(spec *openapi31.Spec) error {
	errPrefix := "failed to write golang"

	if err := os.MkdirAll(golangAPIClientOutputPath, 0o0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := purgeGoFiles(golangAPIClientOutputPath); err != nil {
		return fmt.Errorf("failed to purge golang files: %w", err)
	}

	if err := golang.WriteAPIClientFiles(spec, golangAPIClientOutputPath); err != nil {
		return fmt.Errorf("%s API client files: %w", errPrefix, err)
	}

	if err := golang.WriteAPITypesFiles(spec, golangAPIClientOutputPath); err != nil {
		return fmt.Errorf("%s API client files: %w", errPrefix, err)
	}

	return nil
}

func purgeGoFiles(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".go" {
			if err = os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	})
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
				log.Fatalf("writing typescript files: %v", err)
			}
			wg.Done()
		}()
	}

	if *generateGolang {
		wg.Add(1)
		go func() {
			if err = writeGoFiles(spec); err != nil {
				log.Fatalf("writing golang files: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
