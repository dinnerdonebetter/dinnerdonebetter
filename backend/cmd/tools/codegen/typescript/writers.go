package typescript

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
)

const (
	componentSchemaPrefix = "#/components/schemas/"

	enumsFileName       = "./enums"
	numberRangeFileName = "./number_range"

	jsonContentType = "application/json"
	refKey          = "$ref"
	propertiesKey   = "properties"
	modelsPackage   = "@dinnerdonebetter/models"
)

func removeDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if !slices.Contains(list, item) {
			list = append(list, item)
		}
	}
	return list
}

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

func WriteAPIClientFiles(spec *openapi31.Spec, outputPath string) error {
	clientFiles, err := GenerateClientFiles(spec)
	if err != nil {
		return fmt.Errorf("failed to generate typescript files: %w", err)
	}

	if err = os.MkdirAll(outputPath, 0o0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = purgeTypescriptFiles(outputPath); err != nil {
		return fmt.Errorf("failed to purge typescript files: %w", err)
	}

	createdFunctions := []string{}
	imports := map[string][]string{
		modelsPackage: {
			"QueryFilter",
			"QueryFilteredResult",
		},
	}
	for _, function := range clientFiles {
		fileContents, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		if function.InputType.Type != "" {
			imports[modelsPackage] = append(imports[modelsPackage], function.InputType.Type)
		}

		if function.ResponseType.TypeName != "" && function.ResponseType.TypeName != "string" {
			imports[modelsPackage] = append(imports[modelsPackage], function.ResponseType.TypeName)
		}

		if function.ResponseType.GenericContainer != "" {
			imports[modelsPackage] = append(imports[modelsPackage], function.ResponseType.GenericContainer)
		}

		createdFunctions = append(createdFunctions, fileContents)
	}

	createdFunctions = removeDuplicates(createdFunctions)
	slices.Sort(createdFunctions)

	modelsImports := imports[modelsPackage]
	modelsImports = removeDuplicates(modelsImports)
	slices.Sort(modelsImports)

	clientFile := BuildClientFile(modelsImports)
	for _, createdFile := range createdFunctions {
		clientFile += createdFile + "\n\n"
	}

	clientFile += "\n}\n"

	if err = os.WriteFile(fmt.Sprintf("%s/client.ts", outputPath), []byte(clientFile), 0o600); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	if err = os.WriteFile(fmt.Sprintf("%s/index.ts", outputPath), []byte(APIClientIndexFile), 0o600); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	return nil
}

func WriteModelFiles(spec *openapi31.Spec, outputPath string) error {
	// first do enums
	enums, err := GenerateEnumDefinitions(spec)
	if err != nil {
		return fmt.Errorf("could not generate enums file: %w", err)
	}

	enumsFile := GeneratedDisclaimer + "\n\n"
	for _, enum := range enums {
		def, renderErr := enum.Render()
		if renderErr != nil {
			return fmt.Errorf("could not render enum definition: %w", renderErr)
		}

		enumsFile += fmt.Sprintf("%s\n\n", def)
	}

	// next do models

	modelFiles, err := GenerateModelFiles(spec)
	if err != nil {
		return fmt.Errorf("failed to generate typescript models files: %w", err)
	}

	if err = os.MkdirAll(outputPath, 0o0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = purgeTypescriptFiles(outputPath); err != nil {
		return fmt.Errorf("failed to purge typescript models files: %w", err)
	}

	createdFiles := []string{}
	for filename, function := range modelFiles {
		actualFilepath := fmt.Sprintf("%s/%s.ts", outputPath, filename)

		rawFileContents, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		fileContents := fmt.Sprintf("%s\n\n%s", GeneratedDisclaimer, rawFileContents)
		if err = os.WriteFile(actualFilepath, []byte(fileContents), 0o0600); err != nil {
			return fmt.Errorf("failed to write file %s: %w", actualFilepath, err)
		}

		createdFiles = append(createdFiles, actualFilepath)
	}

	slices.Sort(createdFiles)

	indexFile := fmt.Sprintf("%s\n\n", GeneratedDisclaimer)
	for _, createdFile := range createdFiles {
		indexFile += fmt.Sprintf("export * from './%s';\n", strings.TrimSuffix(strings.TrimPrefix(createdFile, fmt.Sprintf("%s/", outputPath)), ".ts"))
	}
	indexFile += "export * from './enums';\n"

	if err = os.WriteFile(fmt.Sprintf("%s/enums.ts", outputPath), []byte(enumsFile), 0o600); err != nil {
		return fmt.Errorf("failed to write enums file: %w", err)
	}

	for _, staticFile := range StaticModelsFiles {
		indexFile += fmt.Sprintf("export * from './%s';\n", staticFile.Name)
		if err = os.WriteFile(fmt.Sprintf("%s/%s.ts", outputPath, staticFile.Name), []byte(staticFile.Content), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	if err = os.WriteFile(fmt.Sprintf("%s/index.ts", outputPath), []byte(indexFile), 0o600); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

func WriteMockAPIFiles(spec *openapi31.Spec, outputPath string) error {
	mockAPIFuncs, err := GenerateMockAPIFunctions(spec)
	if err != nil {
		return fmt.Errorf("failed to generate typescript models files: %w", err)
	}

	if err = os.MkdirAll(outputPath, 0o0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = purgeTypescriptFiles(outputPath); err != nil {
		return fmt.Errorf("failed to purge typescript models files: %w", err)
	}

	createdFiles := []string{}
	for filename, function := range mockAPIFuncs {
		actualFilepath := fmt.Sprintf("%s/%s.ts", outputPath, filename)

		rawFileContents, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		importStatement := "import type { Page, Route } from '@playwright/test';\n\n"
		if function.ResponseType != "" && function.ResponseType != "string" {
			modelsImports := []string{
				function.ResponseType,
			}

			if function.QueryFiltered {
				modelsImports = append(modelsImports, "QueryFilteredResult")
			}

			importStatement += fmt.Sprintf("import { %s } from '@dinnerdonebetter/models'\n\n", strings.Join(modelsImports, ",\n\t"))
		}
		importStatement += "import { assertClient, assertMethod, ResponseConfig } from './helpers';\n\n"

		fileContents := fmt.Sprintf("%s\n\n%s\n\n%s", GeneratedDisclaimer, importStatement, rawFileContents)
		if err = os.WriteFile(actualFilepath, []byte(fileContents), 0o0600); err != nil {
			return fmt.Errorf("failed to write file %s: %w", actualFilepath, err)
		}

		createdFiles = append(createdFiles, actualFilepath)
	}

	slices.Sort(createdFiles)

	indexFile := fmt.Sprintf("%s\n\n", GeneratedDisclaimer)
	for _, createdFile := range createdFiles {
		indexFile += fmt.Sprintf("export * from './%s';\n", strings.TrimSuffix(strings.TrimPrefix(createdFile, fmt.Sprintf("%s/", outputPath)), ".ts"))
	}

	if err = os.WriteFile(fmt.Sprintf("%s/index.ts", outputPath), []byte(indexFile), 0o600); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	if err = os.WriteFile(fmt.Sprintf("%s/helpers.ts", outputPath), []byte(PlaywrightHelpersFile), 0o600); err != nil {
		return fmt.Errorf("failed to write helpers file: %w", err)
	}

	return nil
}
