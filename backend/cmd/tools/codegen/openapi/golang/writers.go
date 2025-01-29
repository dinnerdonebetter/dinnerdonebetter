package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
)

const (
	GeneratedDisclaimer = `// GENERATED CODE, DO NOT EDIT MANUALLY`

	componentSchemaPrefix = "#/components/schemas/"

	jsonContentType = "application/json"
	refKey          = "$ref"
	propertiesKey   = "properties"
)

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

func WriteAPIClientFiles(spec *openapi31.Spec, outputPath string) error {
	clientFunctions, err := GenerateClientFunctions(spec)
	if err != nil {
		return fmt.Errorf("failed to generate golang files: %w", err)
	}

	if err = os.MkdirAll(outputPath, 0o0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = purgeGoFiles(outputPath); err != nil {
		return fmt.Errorf("failed to purge golang files: %w", err)
	}

	for filename, function := range clientFunctions {
		fileContents, fileImports, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		if fileContents == "" {
			continue
		}

		imports := strings.Join(fileImports, `"`+"\n\t"+`"`)
		if imports != "" {
			imports = fmt.Sprintf(`"%s"`, imports)
		}

		importStatement := fmt.Sprintf(`
import (
	"context"
	"net/http"
	"github.com/dinnerdonebetter/backend/pkg/types"
	%s
)
`, imports)

		fileContents = GeneratedDisclaimer + "\n\npackage apiclient\n\n" + "\n\n" + importStatement + "\n\n" + fileContents

		if err = os.WriteFile(fmt.Sprintf("%s/%s.gen.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	for filename, fileContents := range baseFiles {
		if err = os.WriteFile(fmt.Sprintf("%s/%s.gen.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	for filename, function := range clientFunctions {
		fileContents, fileImports, renderErr := function.RenderTest()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		if fileContents == "" {
			continue
		}

		imports := strings.Join(fileImports, `"`+"\n\t"+`"`)
		if imports != "" {
			imports = fmt.Sprintf(`"%s"`, imports)
		}

		importStatement := fmt.Sprintf(`
import (
	%s
)
`, imports)

		fileContents = GeneratedDisclaimer + "\n\npackage apiclient\n\n" + "\n\n" + importStatement + "\n\n" + fileContents

		if err = os.WriteFile(fmt.Sprintf("%s/%s_gen_test.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	return nil
}

func WriteAPITypesFiles(spec *openapi31.Spec, outputPath string) error {
	clientFunctions, err := GenerateModelFiles(spec)
	if err != nil {
		return fmt.Errorf("failed to generate golang files: %w", err)
	}

	if err = os.MkdirAll(outputPath, 0o0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = purgeGoFiles(outputPath); err != nil {
		return fmt.Errorf("failed to purge golang files: %w", err)
	}

	for filename, function := range clientFunctions {
		fileContents, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render: %w", renderErr)
		}

		if fileContents == "" {
			continue
		}

		importStatement := `
import (
	"context"
)
`

		fileContents = GeneratedDisclaimer + "\n\npackage apiclient\n\n" + "\n\n" + importStatement + "\n\n" + fileContents

		if err = os.WriteFile(fmt.Sprintf("%s/%s.gen.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	for filename, fileContents := range baseFiles {
		if err = os.WriteFile(fmt.Sprintf("%s/%s.gen.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	return nil
}
