package golang

import (
	"fmt"
	"os"
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

func WriteAPIClientFiles(spec *openapi31.Spec, outputPath string) error {
	clientFunctions, err := GenerateClientFunctions(spec)
	if err != nil {
		return fmt.Errorf("failed to generate golang files: %w", err)
	}

	for filename, function := range clientFunctions {
		fileContents, fileImports, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render client function: %w", renderErr)
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
	%s
)
`, imports)

		fileContents = GeneratedDisclaimer + "\n\npackage apiclient\n\n" + "\n\n" + importStatement + "\n\n" + fileContents

		if err = os.WriteFile(fmt.Sprintf("%s/op.%s.gen.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	for filename, fileContents := range baseFiles {
		if err = os.WriteFile(fmt.Sprintf("%s/core.%s.gen.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	for filename, function := range clientFunctions {
		fileContents, fileImports, renderErr := function.RenderTest()
		if renderErr != nil {
			return fmt.Errorf("failed to render function test: %w", renderErr)
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

		if err = os.WriteFile(fmt.Sprintf("%s/op.%s.gen_test.go", outputPath, filename), []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	return nil
}

func WriteAPITypesFiles(spec *openapi31.Spec, outputPath string) error {
	modelFiles, err := GenerateModelFiles(spec)
	if err != nil {
		return fmt.Errorf("failed to generate golang files: %w", err)
	}

	for filename, function := range modelFiles {
		fileContents, renderErr := function.Render()
		if renderErr != nil {
			return fmt.Errorf("failed to render model: %w", renderErr)
		}

		if fileContents == "" {
			continue
		}

		fileContents = GeneratedDisclaimer + "\n\npackage apiclient\n\n" + "\n\n" + fileContents

		outPath := fmt.Sprintf("%s/type.%s.gen.go", outputPath, filename)
		if err = os.WriteFile(outPath, []byte(fileContents), 0o600); err != nil {
			return fmt.Errorf("failed to write index file: %w", err)
		}
	}

	return nil
}
