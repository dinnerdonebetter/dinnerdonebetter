package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	createdFiles := []string{}
	for filename, function := range typescriptClientFiles {
		actualFilepath := fmt.Sprintf("%s/%s.ts", typescriptAPIClientOutputPath, filename)

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

	fmt.Printf("Wrote %d files, had %d Operations\n", len(createdFiles), getOpCountForSpec(spec))

	indexFile := fmt.Sprintf("%s\n\n", typescript.GeneratedDisclaimer)
	for _, createdFile := range createdFiles {
		indexFile += fmt.Sprintf("export * from './%s';\n", strings.TrimSuffix(strings.TrimPrefix(createdFile, fmt.Sprintf("%s/", typescriptAPIClientOutputPath)), ".ts"))
	}

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
