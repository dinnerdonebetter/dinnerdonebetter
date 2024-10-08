package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen"

	"github.com/hashicorp/go-multierror"
)

const (
	destinationDirectory = "../frontend/packages/models"

	timeType                 = "time.Time"
	mapStringToBoolType      = "map[string]bool"
	mapStringToChangeLogType = "map[string]ChangeLog"
	stringType               = "string"
	boolType                 = "bool"
	errorCodeType            = "ErrorCode"
)

type CodeLine struct {
	FieldType    string
	FieldName    string
	DefaultValue string
	IsReadonly   bool
	IsPointer    bool
	IsSlice      bool
	CustomType   bool
}

func buildImportMap() map[string]string {
	importMap := map[string]string{
		"NumberRange":                "main.ts",
		"NumberRangeWithOptionalMax": "main.ts",
		"OptionalNumberRange":        "main.ts",
	}
	for _, u := range codegen.Enums {
		importMap[u.Name] = "_unions.ts"
	}

	for filename, typesToGenerateFor := range codegen.TypeDefinitionFilesToGenerate {
		fileImports := []string{}
		for _, typ := range typesToGenerateFor {
			typeName := reflect.TypeOf(typ).Name()
			fileImports = append(fileImports, typeName)
		}

		for _, imp := range fileImports {
			importMap[imp] = fmt.Sprintf("%s.ts", filename)
		}
	}

	return importMap
}

func main() {
	if destinationDirectory == "artifacts/typescript" {
		if err := os.RemoveAll(destinationDirectory); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(destinationDirectory, 0o750); err != nil {
			panic(err)
		}
	}

	var errors *multierror.Error

	if err := os.WriteFile(fmt.Sprintf("%s/%s", destinationDirectory, "_unions.ts"), []byte(buildUnionsFile()), 0o600); err != nil {
		errors = multierror.Append(errors, err)
	}

	indexOutput := `
export * from './_unions';
export * from './main';
`

	importMap := buildImportMap()
	for _, filename := range sortedMapKeys(codegen.TypeDefinitionFilesToGenerate) {
		output := ""
		typesToGenerateFor := codegen.TypeDefinitionFilesToGenerate[filename]
		filesToImportsMapForFile := map[string]map[string]struct{}{}

		for _, typ := range typesToGenerateFor {
			typInterface, importedInterfaceTypes, err := typescriptInterface(typ)
			if err != nil {
				panic(err)
			}

			for _, imp := range importedInterfaceTypes {
				if _, ok := filesToImportsMapForFile[importMap[imp]]; ok {
					filesToImportsMapForFile[importMap[imp]][imp] = struct{}{}
				} else if importMap[imp] != filename {
					if importMap[imp] == "" {
						continue
					}
					filesToImportsMapForFile[importMap[imp]] = map[string]struct{}{imp: {}}
				}
			}

			output += typInterface + "\n"

			typClass, importedClassTypes, err := typescriptClass(typ)
			if err != nil {
				panic(err)
			}

			for _, imp := range importedClassTypes {
				if _, ok := filesToImportsMapForFile[importMap[imp]]; ok {
					filesToImportsMapForFile[importMap[imp]][imp] = struct{}{}
				} else {
					if importMap[imp] == "" {
						continue
					}
					filesToImportsMapForFile[importMap[imp]] = map[string]struct{}{imp: {}}
				}
			}

			output += typClass + "\n"
		}

		fileOutput := copyString(generatedDisclaimer)
		for _, file := range sortedMapKeys(filesToImportsMapForFile) {
			imports := filesToImportsMapForFile[file]
			if file == fmt.Sprintf("%s.ts", filename) {
				continue
			}

			fileOutput += fmt.Sprintf("import { %s } from './%s';\n", strings.Join(sortedMapKeys(imports), ", "), strings.TrimSuffix(file, ".ts"))
		}

		indexOutput += fmt.Sprintf("export * from './%s';\n", strings.TrimSuffix(filename, ".ts"))
		finalOutput := fileOutput + "\n" + output

		if err := os.WriteFile(fmt.Sprintf("%s/%s.ts", destinationDirectory, filename), []byte(finalOutput), 0o600); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if err := os.WriteFile(fmt.Sprintf("%s/index.ts", destinationDirectory), []byte(indexOutput), 0o600); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := errors.ErrorOrNil(); err != nil {
		panic(err)
	}
}
