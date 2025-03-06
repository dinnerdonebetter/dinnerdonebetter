package main

import (
	"bytes"
	"fmt"
	"github.com/verygoodsoftwarenotvirus/typewizard/models"
	"github.com/verygoodsoftwarenotvirus/typewizard/utils"
	"iter"
	"maps"
	"os"
	"slices"
	"strings"
	"text/template"
)

const (
	oneToOneFilepath = "cmd/tools/codegen/converters/artifacts/one_to_one.go" // "internal/services/eating/grpc/converters/one_to_one.go"
)

func main() {
	oneForOne, err := buildOneForOneFile()
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile(oneToOneFilepath, []byte(oneForOne), 0o0644); err != nil {
		panic(err)
	}
}

type converterFunction struct {
	InputType,
	OutputType *models.Struct
}

type fieldPair struct {
	InputField,
	OutputField *models.StructField
}

func (f *converterFunction) CommonFieldsWithOutputType() []*fieldPair {
	output := []*fieldPair{}

	for _, field := range f.OutputType.Fields {
		if x := f.InputType.Fields.AsMapCollection(func(field *models.StructField) string {
			return field.Name
		})[field.Name]; x != nil {
			output = append(output, &fieldPair{
				InputField:  x,
				OutputField: field,
			})
		}
	}

	slices.SortFunc(output, func(a, b *fieldPair) int {
		return strings.Compare(a.InputField.Name, b.InputField.Name)
	})

	return output
}

func (f *converterFunction) Render() (string, error) {
	const funcTemplate = `
func Convert{{.InputType.Name}}ToGRPC{{.OutputType.Name}}(in *{{.InputType.Package.Name}}.{{.InputType.Name}}) *{{.OutputType.Package.Name}}.{{.OutputType.Name}} {
	x := &{{.OutputType.Package.Name}}.{{.OutputType.Name}}{
		{{- range $fieldPair := commonOutputFields }}
			{{ $fieldPair.InputField.Name }}: {{ convertType $fieldPair }},
		{{- end }}
	}

	return x
}
`

	tmpl, err := template.New("convertFunc").Funcs(map[string]any{
		"commonOutputFields": f.CommonFieldsWithOutputType,
		"convertType": func(field *fieldPair) string {
			typePair := fmt.Sprintf("%s.%s | %s.%s", field.InputField.TypePackage, strings.TrimPrefix(field.InputField.Type, "*"), field.OutputField.TypePackage, strings.TrimPrefix(field.OutputField.Type, "*"))

			switch typePair {
			case "time.Time | timestamppb.Timestamp":
				return fmt.Sprintf("in.%s", field.InputField.Name)
			default:
				return fmt.Sprintf("in.%s", field.InputField.Name)
			}
		},
	}).Parse(funcTemplate)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err = tmpl.Execute(&b, f); err != nil {
		return "", err
	}

	return b.String(), nil
}

func toSlice[T any](seq iter.Seq[T]) []T {
	var slice []T
	for v := range seq {
		slice = append(slice, v)
	}
	return slice
}

func buildOneForOneFile() (string, error) {
	grpcTypes, err := utils.GetTypesForPackage("internal/grpc/messages", "messages", nil)
	if err != nil {
		return "", fmt.Errorf("fetching grpc types")
	}

	eatingTypes, err := utils.GetTypesForPackage("internal/services/eating/types", "types", nil)
	if err != nil {
		return "", fmt.Errorf("fetching eating types")
	}

	commonTypes := map[string]bool{}
	allPackages := map[string]bool{}

	// iterate through both map types and fine common types
	for key := range grpcTypes {
		allPackages[grpcTypes[key].Package.Path] = true
		if x, found := eatingTypes[key]; found {
			allPackages[x.Package.Path] = true
			commonTypes[key] = true
		}
	}

	output := fmt.Sprintf(`package converters

import (
 "%s"
)

`, strings.Join(toSlice(maps.Keys(allPackages)), "\"\n\t\""))

	for t := range commonTypes {
		cf := &converterFunction{
			InputType:  eatingTypes[t],
			OutputType: grpcTypes[t],
		}
		out, err := cf.Render()
		if err != nil {
			panic(err)
		}

		output += out
	}

	return output, nil
}
