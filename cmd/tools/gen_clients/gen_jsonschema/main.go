package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	codegen "github.com/dinnerdonebetter/backend/cmd/tools/gen_clients"

	"github.com/invopop/jsonschema"
)

func getTypeName(x any) string {
	if t := reflect.TypeOf(x); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func writeSchemaToFile(schema *jsonschema.Schema, outputFilepath string) error {
	encodedSchema, err := json.MarshalIndent(schema, "", "\t")
	if err != nil {
		return err
	}

	if err = os.WriteFile(outputFilepath, encodedSchema, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := os.MkdirAll("artifacts", os.ModeDir); err != nil {
		// it's fine for this to fail
	}

	for _, typesToGenerateFor := range codegen.TypeDefinitionFilesToGenerate {
		for _, typ := range typesToGenerateFor {
			schema := jsonschema.Reflect(typ)

			if err := writeSchemaToFile(schema, fmt.Sprintf("artifacts/%s.jsonschema", getTypeName(typ))); err != nil {
				panic(err)
			}
		}
	}
}
