package main

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen"
)

func buildUnionDeclaration(def *codegen.Enum) string {
	if def == nil {
		return ""
	}

	output := ""

	output += fmt.Sprintf("/**\n * %s\n */\n", def.Comment)
	output += fmt.Sprintf("export const ALL_%s: string[] = [\n", def.ConstName)
	for _, value := range def.Values {
		output += fmt.Sprintf("  %q,\n", value)
	}
	output += "];\n"

	output += fmt.Sprintf("type %sTypeTuple = typeof ALL_%s;\n", def.Name, def.ConstName)
	output += fmt.Sprintf("export type %s = %sTypeTuple[number];\n\n", def.Name, def.Name)

	return output
}

func buildUnionsFile() string {
	output := copyString(generatedDisclaimer)

	for _, def := range codegen.Enums {
		output += buildUnionDeclaration(def)
	}

	return output
}
