package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/verygoodsoftwarenotvirus/typewizard/models"
	"github.com/verygoodsoftwarenotvirus/typewizard/utils"
	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

var (
	skipTypes = map[string]bool{
		"Uint16RangeWithOptionalMax":                                                   true,
		"Uint32RangeWithOptionalMax":                                                   true,
		"Float32RangeWithOptionalMax":                                                  true,
		"Float32RangeWithOptionalMaxUpdateRequestInput":                                true,
		"Uint32RangeWithOptionalMaxUpdateRequestInput":                                 true,
		"Uint16RangeWithOptionalMaxUpdateRequestInput":                                 true,
		"RecipePrepTaskWithinRecipeCreationRequestInput":                               true,
		"RecipeStepCompletionConditionForExistingRecipeCreationRequestInput":           true,
		"RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput": true,
		"RecipePrepTaskStepWithinRecipeCreationRequestInput":                           true,
	}

	outputFilenames = map[string]string{
		"ValidVessel":                                    "valid_vessels",
		"ValidIngredientState":                           "valid_ingredient_states",
		"Household":                                      "households",
		"ValidPreparationVessel":                         "valid_preparation_vessels",
		"Meal":                                           "meals",
		"RecipeRating":                                   "recipe_ratings",
		"RecipePrepTask":                                 "recipe_prep_tasks",
		"RecipeStep":                                     "recipe_steps",
		"ValidIngredientStateIngredient":                 "valid_ingredient_state_ingredients",
		"ValidInstrument":                                "valid_instruments",
		"ValidPreparation":                               "valid_preparations",
		"ValidIngredientGroup":                           "valid_ingredient_groups",
		"MealPlanGroceryListItem":                        "meal_plan_grocery_list_items",
		"OAuth2Client":                                   "oauth2_clients",
		"RecipePrepTaskStep":                             "recipe_prep_task_steps",
		"RecipeMedia":                                    "recipe_medias",
		"UserIngredientPreference":                       "user_ingredient_preferences",
		"Recipe":                                         "recipes",
		"Webhook":                                        "webhooks",
		"ValidIngredient":                                "valid_ingredients",
		"HouseholdInvitation":                            "household_invitations",
		"ServiceSetting":                                 "service_settings",
		"ValidIngredientMeasurementUnit":                 "valid_ingredient_measurement_units",
		"RecipeStepVessel":                               "recipe_step_vessels",
		"ServiceSettingConfiguration":                    "service_setting_configurations",
		"ValidIngredientGroupMember":                     "valid_ingredient_group_members",
		"MealPlanOptionVote":                             "meal_plan_option_votes",
		"HouseholdInstrumentOwnership":                   "household_instrument_ownerships",
		"MealPlanEvent":                                  "meal_plan_events",
		"ValidMeasurementUnitConversion":                 "valid_measurement_unit_conversions",
		"PasswordResetToken":                             "password_reset_tokens",
		"RecipeStepIngredient":                           "recipe_step_ingredients",
		"MealPlanOption":                                 "meal_plan_options",
		"RecipeStepProduct":                              "recipe_step_products",
		"MealPlanTask":                                   "meal_plan_tasks",
		"RecipeStepInstrument":                           "recipe_step_instruments",
		"ValidMeasurementUnit":                           "valid_measurement_units",
		"MealPlan":                                       "meal_plans",
		"ValidIngredientPreparation":                     "valid_ingredient_preparations",
		"ValidPreparationInstrument":                     "valid_preparation_instruments",
		"UserNotification":                               "user_notifications",
		"MealComponent":                                  "meal_components",
		"RecipeStepCompletionCondition":                  "recipe_step_completion_conditions",
		"WebhookTriggerEvent":                            "webhook_trigger_events",
		"ValidIngredientCreationRequestInput":            "valid_ingredients",
		"UserIngredientPreferenceCreationRequestInput":   "user_ingredient_preferences",
		"ValidIngredientUpdateRequestInput":              "valid_ingredients",
		"RecipePrepTaskWithinRecipeCreationRequestInput": "recipe_prep_tasks",
	}
)

const (
	miscFile = `package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertFloat32RangeWithOptionalMaxUpdateRequestInputToFloat32RangeWithOptionalMax(input *messages.Float32RangeWithOptionalMaxUpdateRequestInput) *messages.Float32RangeWithOptionalMax {

output := &messages.Float32RangeWithOptionalMax{
    Max: input.Max,
    Min: input.Min,
}

return output
}


func ConvertUint16RangeWithOptionalMaxUpdateRequestInputToUint16RangeWithOptionalMax(input *messages.Uint16RangeWithOptionalMaxUpdateRequestInput) *messages.Uint16RangeWithOptionalMax {

output := &messages.Uint16RangeWithOptionalMax{
    Min: input.Min,
    Max: input.Max,
}

return output
}

func ConvertUint32RangeWithOptionalMaxUpdateRequestInputToUint32RangeWithOptionalMax(input *messages.Uint32RangeWithOptionalMaxUpdateRequestInput) *messages.Uint32RangeWithOptionalMax {

output := &messages.Uint32RangeWithOptionalMax{
    Max: input.Max,
    Min: input.Min,
}

return output
}


`
)

func parseProtoService() (*parser.Proto, error) {
	reader, err := os.Open("internal/services/service.proto")
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	parsed, err := protoparser.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("parsing proto spec: %w", err)
	}

	if err = reader.Close(); err != nil {
		return nil, fmt.Errorf("closing reader: %w", err)
	}

	return parsed, nil
}

func main() {
	allTypes, err := utils.GetTypesForPackage("internal/grpc/messages", "messages", nil)
	if err != nil {
		panic(err)
	}

	updateTypes := make(map[string]*models.Struct)
	creationTypes := make(map[string]*models.Struct)
	baseTypes := make(map[string]*models.Struct)

	for _, t := range allTypes {
		if _, ok := skipTypes[t.Name]; ok {
			continue
		}

		if strings.HasSuffix(t.Name, "UpdateRequestInput") {
			updateTypes[t.Name] = t

			if x := strings.TrimSuffix(t.Name, "UpdateRequestInput"); allTypes[x] != nil {
				baseTypes[t.Name] = allTypes[x]
			}
		}

		if strings.HasSuffix(t.Name, "CreationRequestInput") {
			creationTypes[t.Name] = t

			if x := strings.TrimSuffix(t.Name, "CreationRequestInput"); allTypes[x] != nil {
				baseTypes[t.Name] = allTypes[x]
			}
		}
	}

	customTypePairs := map[string]string{
		"Float32RangeWithOptionalMaxUpdateRequestInput":  "Float32RangeWithOptionalMax",
		"RecipePrepTaskWithinRecipeCreationRequestInput": "RecipePrepTask",
		"Uint32RangeWithOptionalMaxUpdateRequestInput":   "Uint32RangeWithOptionalMax",
		"Uint16RangeWithOptionalMaxUpdateRequestInput":   "Uint16RangeWithOptionalMax",
	}

	outputTypes := map[*models.Struct]string{}

	for t1, t2 := range customTypePairs {
		outputTypes[allTypes[t1]] += GenerateConverter(allTypes[t1], allTypes[t2])
	}

	for _, t := range creationTypes {
		if _, ok := skipTypes[t.Name]; ok {
			continue
		}

		outputTypes[baseTypes[t.Name]] += GenerateConverter(t, baseTypes[t.Name])
	}

	for _, t := range updateTypes {
		if _, ok := skipTypes[t.Name]; ok {
			continue
		}

		outputTypes[baseTypes[t.Name]] += GenerateConverter(t, baseTypes[t.Name])
	}

	if err = os.WriteFile("internal/grpcimpl/converters/generated/misc.go", []byte(miscFile), 0o0644); err != nil {
		panic(err)
	}

	for t, f := range outputTypes {
		if _, ok := skipTypes[t.Name]; ok {
			continue
		}

		outputFilename := outputFilenames[t.Name]
		if outputFilename == "" {
			fmt.Printf("no output filename for %s\n", t.Name)
			continue
		}

		outputFilename = fmt.Sprintf("internal/grpcimpl/converters/generated/%s.go", outputFilename)

		importBlock := buildImportBlock(t.Package)
		outputFile := fmt.Sprintf("package converters\n\n%s\n\n%s", importBlock, f)

		if err = os.WriteFile(
			outputFilename,
			[]byte(outputFile),
			0o0644,
		); err != nil {
			panic(err)
		}
	}
}

func firstLetterIsCapitalized(x string) bool {
	if len(x) == 0 {
		return false
	}
	return strings.ToUpper(string(x[0])) == string(x[0])
}

func buildImportBlock(imports ...models.Package) string {
	var importBlock strings.Builder
	if len(imports) > 0 {
		importBlock.WriteString("import (\n")
		for _, pkg := range imports {
			if pkg.Alias != "" {
				importBlock.WriteString(fmt.Sprintf("\t%s %q\n", pkg.Alias, pkg.Path))
			} else {
				importBlock.WriteString(fmt.Sprintf("\t%q\n", pkg.Path))
			}
		}
		importBlock.WriteString(")")
	}

	return importBlock.String()
}

// BEGIN AI BULLSHIT

func baseType(typ string) string {
	typ = strings.TrimPrefix(typ, "[]")
	typ = strings.TrimPrefix(typ, "*")
	return typ
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func GenerateConverter(source, target *models.Struct) string {
	var code bytes.Buffer

	// Function signature
	fmt.Fprintf(&code, "func Convert%sTo%s(input *%s.%s) *%s.%s {\n",
		source.Name, target.Name,
		source.Package.Name, source.Name,
		target.Package.Name, target.Name)

	var loops []string
	assignments := make(map[string]string)

	// Process each target field
	for _, targetField := range target.Fields {
		if !firstLetterIsCapitalized(targetField.Name) {
			continue
		}

		// Find matching source field
		var sourceField *models.StructField
		for _, f := range source.Fields {
			if f.Name == targetField.Name {
				sourceField = f
				break
			}
		}
		if sourceField == nil {
			continue
		}

		// Handle slice fields
		if sourceField.IsSlice && targetField.IsSlice {
			sourceElem := fmt.Sprintf("%s.%s", sourceField.TypePackage, baseType(sourceField.Type))
			targetElem := fmt.Sprintf("%s.%s", targetField.TypePackage, baseType(targetField.Type))

			if sourceElem != targetElem {
				varName := "converted" + targetField.Name

				if sourceField.BasicType {
					switch sourceField.Type {
					default:
						continue
					}
				}

				converterFunc := fmt.Sprintf("Convert%sTo%s", baseType(sourceField.Type), baseType(targetField.Type))

				loopCode := fmt.Sprintf("%s := make([]*%s, 0, len(input.%s))\n",
					varName, targetElem, sourceField.Name)
				loopCode += fmt.Sprintf("for _, item := range input.%s {\n", sourceField.Name)
				loopCode += fmt.Sprintf("    %s = append(%s, %s(item))\n",
					varName, varName, converterFunc)
				loopCode += "}\n"

				loops = append(loops, loopCode)
				assignments[targetField.Name] = varName
			} else {
				assignments[targetField.Name] = fmt.Sprintf("input.%s", sourceField.Name)
			}
			continue
		}

		// Handle struct pointers
		if !sourceField.BasicType && !sourceField.FromStandardLibrary &&
			!targetField.BasicType && !targetField.FromStandardLibrary {

			sourceType := fmt.Sprintf("%s.%s", sourceField.TypePackage, sourceField.Type)
			targetType := fmt.Sprintf("%s.%s", targetField.TypePackage, targetField.Type)

			if sourceType != targetType {
				converterFunc := fmt.Sprintf("Convert%sTo%s", baseType(sourceField.Type), baseType(targetField.Type))
				assignments[targetField.Name] =
					fmt.Sprintf("%s(input.%s)", converterFunc, sourceField.Name)
			} else {
				assignments[targetField.Name] = fmt.Sprintf("input.%s", sourceField.Name)
			}
			continue
		}

		// Direct assignment for basic/stdlib types
		assignments[targetField.Name] = fmt.Sprintf("input.%s", sourceField.Name)
	}

	// Write conversion loops
	for _, loop := range loops {
		code.WriteString(loop)
	}

	// Build output struct
	code.WriteString("\noutput := &")
	fmt.Fprintf(&code, "%s.%s{\n", target.Package.Name, target.Name)

	for fieldName, value := range assignments {
		fmt.Fprintf(&code, "    %s: %s,\n", fieldName, value)
	}

	code.WriteString("}\n\nreturn output\n}\n")

	return code.String()
}

// END AI BULLSHIT

/*


const (
	inputFieldName  = "input"
	outputFieldName = "output"
)

func writeConversionFunctionForTypes(typeA, typeB *models.Struct) string {
	var sb strings.Builder

	if typeA == nil || typeB == nil {
		if typeA != nil {
			fmt.Println(typeA.Name)
		}

		if typeB != nil {
			fmt.Println(typeB.Name)
		}

		panic("nil type!")
	}

	funcName := fmt.Sprintf("Convert%sTo%s", typeA.Name, typeB.Name)
	sb.WriteString(fmt.Sprintf("func %s(%s *%s.%s) *%s.%s {\n", funcName, inputFieldName, typeA.Package.Name, typeA.Name, typeB.Package.Name, typeB.Name))
	sb.WriteString(fmt.Sprintf("\t%s := &%s.%s{\n", outputFieldName, typeB.Package.Name, typeB.Name))

	for _, fieldA := range typeA.Fields {
		if !firstLetterIsCapitalized(fieldA.Name) {
			continue
		}
		for _, fieldB := range typeB.Fields {
			if !firstLetterIsCapitalized(fieldB.Name) {
				continue
			}
			if fieldA.Name == fieldB.Name {
				if fieldA.Type == fieldB.Type {
					sb.WriteString(fmt.Sprintf("\t\t%s: %s.%s,\n", fieldB.Name, inputFieldName, fieldA.Name))
				} else if (fieldA.Type == "float32" && fieldB.Type == "float64") || (fieldA.Type == "float64" && fieldB.Type == "float32") {
					sb.WriteString(fmt.Sprintf("\t\t%s: %s(%s.%s),\n", fieldB.Name, inputFieldName, fieldB.Type, fieldA.Name))
				}
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\t}\n\n\treturn %s\n", outputFieldName))
	sb.WriteString("}\n\n")

	return sb.String()
}
*/
