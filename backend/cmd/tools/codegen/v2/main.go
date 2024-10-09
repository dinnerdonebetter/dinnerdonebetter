package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen/v2/typescript"

	"github.com/swaggest/openapi-go/openapi31"
)

const (
	specFilepath                  = "../openapi_spec.yaml"
	typescriptAPIClientOutputPath = "../frontend/packages/generated-client"
	typescriptModelsOutputPath    = "../frontend/packages/generated-models"
	typescriptMockAPIOutputPath   = "../frontend/packages/generated-mock-playwright-api"
)

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

func writeTypescriptFiles(spec *openapi31.Spec) error {
	if err := typescript.WriteModelFiles(spec, typescriptModelsOutputPath); err != nil {
		return fmt.Errorf("failed to write typescript model files: %w", err)
	}

	if err := typescript.WriteAPIClientFiles(spec, typescriptAPIClientOutputPath); err != nil {
		return fmt.Errorf("failed to write typescript API client files: %w", err)
	}

	if err := typescript.WriteMockAPIFiles(spec, typescriptMockAPIOutputPath); err != nil {
		return fmt.Errorf("failed to write typescript mock API files: %w", err)
	}

	return nil
}

func main() {
	spec, err := loadSpec(specFilepath)
	if err != nil {
		log.Fatalf("failed to load spec: %v", err)
	}

	if err = writeTypescriptFiles(spec); err != nil {
		log.Fatalf("failed to write typescript files: %v", err)
	}
}
