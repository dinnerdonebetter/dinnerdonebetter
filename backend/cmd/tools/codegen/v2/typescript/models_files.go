package typescript

import "github.com/swaggest/openapi-go/openapi31"

type Field struct {
	Name,
	Type,
	DefaultValue string
	Nullable,
	Array bool
}

type TypeDefinition struct {
	Name   string
	Fields []Field
}

func GenerateModelFiles(spec *openapi31.Spec) (map[string]*TypeDefinition, error) {
	output := map[string]*TypeDefinition{}

	for name, component := range spec.Components.Schemas {
		def := &TypeDefinition{
			Name: name,
		}

		if properties, ok := component["properties"]; ok {
			if propMap, ok2 := properties.(map[string]any); ok2 {
				for k, v := range propMap {
					field := Field{
						Name: k,
					}

					if typeMap, ok3 := v.(map[string]any); ok3 {
						if typ, ok5 := typeMap["type"]; ok5 {
							field.Type = typ.(string)
						}
					}

					def.Fields = append(def.Fields, field)
				}
			}
		}

		output[name] = def
	}

	return output, nil
}

func (d *TypeDefinition) Render() (string, error) {
	return "", nil
}
