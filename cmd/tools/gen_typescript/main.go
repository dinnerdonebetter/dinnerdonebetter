package main

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/prixfixeco/backend/pkg/types"
)

const (
	destinationDirectory = "../frontend/packages/models"

	timeType            = "time.Time"
	mapStringToBoolType = "map[string]bool"
	stringType          = "string"
	boolType            = "bool"

	// TODO: generate this file programmatically
	unionsFile = `// valid ingredient state attribute types
export const ALL_VALID_INGREDIENT_STATE_ATTRIBUTE_TYPES: string[] = [
  'texture',
  'consistency',
  'color',
  'appearance',
  'odor',
  'taste',
  'sound',
  'other',
];

type ValidIngredientStateAttributeTypeTuple = typeof ALL_VALID_INGREDIENT_STATE_ATTRIBUTE_TYPES;
export type ValidIngredientStateAttributeType = ValidIngredientStateAttributeTypeTuple[number];

// meal plan statuses
export const ALL_VALID_MEAL_PLAN_STATUSES: string[] = ['awaiting_votes', 'finalized'];
type ValidMealPlanStatusTypeTuple = typeof ALL_VALID_MEAL_PLAN_STATUSES;
export type ValidMealPlanStatus = ValidMealPlanStatusTypeTuple[number];

// meal plan election methods
export const ALL_VALID_MEAL_PLAN_ELECTION_METHODS: string[] = ['schulze', 'instant-runoff'];
type ValidMealPlanElectionMethodTypeTuple = typeof ALL_VALID_MEAL_PLAN_ELECTION_METHODS;
export type ValidMealPlanElectionMethod = ValidMealPlanElectionMethodTypeTuple[number];

// recipe step product types
export const ALL_RECIPE_STEP_PRODUCT_TYPES = ['instrument', 'ingredient'];
type ValidRecipeStepProductTypeTuple = typeof ALL_RECIPE_STEP_PRODUCT_TYPES;
export type ValidRecipeStepProductType = ValidRecipeStepProductTypeTuple[number];

// meal component types
export const ALL_MEAL_COMPONENT_TYPES = [
  'main',
  'side',
  'appetizer',
  'beverage',
  'dessert',
  'soup',
  'salad',
  'amuse-bouche',
  'unspecified',
];

type MealComponentTypeTuple = typeof ALL_MEAL_COMPONENT_TYPES;
export type MealComponentType = MealComponentTypeTuple[number];
`
)

var (
	filesToGenerate = map[string][]any{
		"admin.ts": {
			types.ModifyUserPermissionsInput{},
		},
		"apiClients.ts": {
			types.APIClient{},
			types.APIClientCreationRequestInput{},
			types.APIClientCreationResponse{},
		},
		"auth.ts": {
			types.ChangeActiveHouseholdInput{},
			types.PASETOCreationInput{},
			types.PASETOResponse{},
			types.PasswordResetToken{},
			types.PasswordResetTokenCreationRequestInput{},
			types.PasswordResetTokenRedemptionRequestInput{},
			types.TOTPSecretRefreshInput{},
			types.TOTPSecretVerificationInput{},
			types.TOTPSecretRefreshResponse{},
			types.PasswordUpdateInput{},
		},
		"errors.ts": {
			types.APIError{},
		},
		"householdInvitations.ts": {
			types.HouseholdInvitation{},
			types.HouseholdInvitationUpdateRequestInput{},
			types.HouseholdInvitationCreationRequestInput{},
		},
		"households.ts": {
			types.Household{},
			types.HouseholdCreationRequestInput{},
			types.HouseholdUpdateRequestInput{},
			types.HouseholdOwnershipTransferInput{},
		},
		"householdUserMemberships.ts": {
			types.HouseholdUserMembership{},
			types.HouseholdUserMembershipWithUser{},
			types.HouseholdUserMembershipCreationRequestInput{},
		},
		"mealPlanEvents.ts": {
			types.MealPlanEvent{},
			types.MealPlanEventCreationRequestInput{},
			types.MealPlanEventUpdateRequestInput{},
		},
		"mealPlanGroceryListItems.ts": {
			types.MealPlanGroceryListItem{},
			types.MealPlanGroceryListItemCreationRequestInput{},
			types.MealPlanGroceryListItemUpdateRequestInput{},
		},
		"mealPlanOptions.ts": {
			types.MealPlanOption{},
			types.MealPlanOptionCreationRequestInput{},
			types.MealPlanOptionUpdateRequestInput{},
		},
		"mealPlanOptionVotes.ts": {
			types.MealPlanOptionVote{},
			types.MealPlanOptionVoteCreationInput{},
			types.MealPlanOptionVoteCreationRequestInput{},
			types.MealPlanOptionVoteUpdateRequestInput{},
		},
		"mealPlans.ts": {
			types.MealPlan{},
			types.MealPlanCreationRequestInput{},
			types.MealPlanUpdateRequestInput{},
		},
		"mealPlanTasks.ts": {
			types.MealPlanTask{},
			types.MealPlanTaskCreationRequestInput{},
			types.MealPlanTaskStatusChangeRequestInput{},
		},
		"meals.ts": {
			types.Meal{},
			types.MealCreationRequestInput{},
			types.MealUpdateRequestInput{},
		},
		"mealComponents.ts": {
			types.MealComponent{},
			types.MealComponentCreationRequestInput{},
			types.MealComponentUpdateRequestInput{},
		},
		"permissions.ts": {
			types.UserPermissionsRequestInput{},
			types.UserPermissionsResponse{},
		},
		"recipeMedia.ts": {
			types.RecipeMedia{},
			types.RecipeMediaCreationRequestInput{},
			types.RecipeMediaUpdateRequestInput{},
		},
		"recipePrepTasks.ts": {
			types.RecipePrepTask{},
			types.RecipePrepTaskCreationRequestInput{},
			types.RecipePrepTaskWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskUpdateRequestInput{},
		},
		"recipePrepTaskSteps.ts": {
			types.RecipePrepTaskStep{},
			types.RecipePrepTaskStepWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskStepCreationRequestInput{},
			types.RecipePrepTaskStepUpdateRequestInput{},
		},
		"recipeStepCompletionConditions.ts": {
			types.RecipeStepCompletionCondition{},
			types.RecipeStepCompletionConditionIngredient{},
			types.RecipeStepCompletionConditionCreationRequestInput{},
			types.RecipeStepCompletionConditionIngredientCreationRequestInput{},
			types.RecipeStepCompletionConditionUpdateRequestInput{},
		},
		"recipeStepIngredients.ts": {
			types.RecipeStepIngredient{},
			types.RecipeStepIngredientCreationRequestInput{},
			types.RecipeStepIngredientUpdateRequestInput{},
		},
		"recipeStepInstruments.ts": {
			types.RecipeStepInstrument{},
			types.RecipeStepInstrumentCreationRequestInput{},
			types.RecipeStepInstrumentUpdateRequestInput{},
		},
		"recipeStepProducts.ts": {
			types.RecipeStepProduct{},
			types.RecipeStepProductCreationRequestInput{},
			types.RecipeStepProductUpdateRequestInput{},
		},
		"recipeSteps.ts": {
			types.RecipeStep{},
			types.RecipeStepCreationRequestInput{},
			types.RecipeStepUpdateRequestInput{},
		},
		"recipes.ts": {
			types.Recipe{},
			types.RecipeCreationRequestInput{},
			types.RecipeUpdateRequestInput{},
		},
		"users.ts": {
			types.UserStatusResponse{},
			types.User{},
			types.UserRegistrationInput{},
			types.UserCreationResponse{},
			types.UserLoginInput{},
			types.UsernameReminderRequestInput{},
			types.UserAccountStatusUpdateInput{},
		},
		"validIngredientMeasurementUnits.ts": {
			types.ValidIngredientMeasurementUnit{},
			types.ValidIngredientMeasurementUnitCreationRequestInput{},
			types.ValidIngredientMeasurementUnitUpdateRequestInput{},
		},
		"validIngredientPreparations.ts": {
			types.ValidIngredientPreparation{},
			types.ValidIngredientPreparationCreationRequestInput{},
			types.ValidIngredientPreparationUpdateRequestInput{},
		},
		"validIngredientStates.ts": {
			types.ValidIngredientState{},
			types.ValidIngredientStateCreationRequestInput{},
			types.ValidIngredientStateUpdateRequestInput{},
		},
		"validIngredientStateIngredients.ts": {
			types.ValidIngredientStateIngredient{},
			types.ValidIngredientStateIngredientCreationRequestInput{},
			types.ValidIngredientStateIngredientUpdateRequestInput{},
		},
		"validIngredients.ts": {
			types.ValidIngredient{},
			types.ValidIngredientCreationRequestInput{},
			types.ValidIngredientUpdateRequestInput{},
		},
		"validInstruments.ts": {
			types.ValidInstrument{},
			types.ValidInstrumentCreationRequestInput{},
			types.ValidInstrumentUpdateRequestInput{},
		},
		"validMeasurementUnitConversions.ts": {
			types.ValidMeasurementUnitConversion{},
			types.ValidMeasurementUnitConversionCreationRequestInput{},
			types.ValidMeasurementUnitConversionUpdateRequestInput{},
		},
		"validMeasurementUnits.ts": {
			types.ValidMeasurementUnit{},
			types.ValidMeasurementUnitCreationRequestInput{},
			types.ValidMeasurementUnitUpdateRequestInput{},
		},
		"validPreparationInstruments.ts": {
			types.ValidPreparationInstrument{},
			types.ValidPreparationInstrumentCreationRequestInput{},
			types.ValidPreparationInstrumentUpdateRequestInput{},
		},
		"validPreparations.ts": {
			types.ValidPreparation{},
			types.ValidPreparationCreationRequestInput{},
			types.ValidPreparationUpdateRequestInput{},
		},
		"webhooks.ts": {
			types.Webhook{},
			types.WebhookTriggerEvent{},
			types.WebhookCreationRequestInput{},
		},
	}
)

var (
	generatedDisclaimer = "// Code generated by gen_typescript. DO NOT EDIT.\n\n"
	customTypeMap       = map[string]string{
		"MealPlan.status":                                        "ValidMealPlanStatus",
		"MealPlan.electionMethod":                                "ValidMealPlanElectionMethod",
		"ValidIngredientState.attributeType":                     "ValidIngredientStateAttributeType",
		"RecipeStepProduct.type":                                 "ValidRecipeStepProductType",
		"MealComponent.componentType":                            "MealComponentType",
		"MealPlanCreationRequestInput.status":                    "ValidMealPlanStatus",
		"MealPlanCreationRequestInput.electionMethod":            "ValidMealPlanElectionMethod",
		"ValidIngredientStateCreationRequestInput.attributeType": "ValidIngredientStateAttributeType",
		"RecipeStepProductCreationRequestInput.type":             "ValidRecipeStepProductType",
		"MealComponentCreationRequestInput.componentType":        "MealComponentType",
		"MealPlanUpdateRequestInput.status":                      "ValidMealPlanStatus",
		"MealPlanUpdateRequestInput.electionMethod":              "ValidMealPlanElectionMethod",
		"ValidIngredientStateUpdateRequestInput.attributeType":   "ValidIngredientStateAttributeType",
		"RecipeStepProductUpdateRequestInput.type":               "ValidRecipeStepProductType",
		"MealComponentUpdateRequestInput.componentType":          "MealComponentType",
	}

	defaultValues = map[string]string{
		"ValidMealPlanStatus":               "'awaiting_votes'",
		"ValidMealPlanElectionMethod":       "'schulze'",
		"ValidIngredientStateAttributeType": "'other'",
		"ValidRecipeStepProductType":        "'ingredient'",
		"MealComponentType":                 "'unspecified'",
	}
)

func buildImportMap() map[string]string {
	importMap := map[string]string{
		"ValidMealPlanStatus":               "_unions.ts",
		"ValidMealPlanElectionMethod":       "_unions.ts",
		"ValidIngredientStateAttributeType": "_unions.ts",
		"ValidRecipeStepProductType":        "_unions.ts",
		"MealComponentType":                 "_unions.ts",
	}

	for filename, typesToGenerateFor := range filesToGenerate {
		fileImports := []string{}
		for _, typ := range typesToGenerateFor {
			fileImports = append(fileImports, reflect.TypeOf(typ).Name())
		}

		for _, imp := range fileImports {
			importMap[imp] = filename
		}
	}

	return importMap
}

func main() {
	if destinationDirectory == "artifacts/typescript" {
		if err := os.RemoveAll(destinationDirectory); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(destinationDirectory, os.ModePerm); err != nil {
			panic(err)
		}
	}

	var errors *multierror.Error

	if err := os.WriteFile(fmt.Sprintf("%s/%s", destinationDirectory, "_unions.ts"), []byte(generatedDisclaimer+unionsFile), 0o600); err != nil {
		errors = multierror.Append(errors, err)
	}

	indexOutput := `
export * from './_unions';
export * from './pagination';
`

	importMap := buildImportMap()
	for filename, typesToGenerateFor := range filesToGenerate {
		output := ""
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

		fileOutput := fmt.Sprintf("%s", generatedDisclaimer)
		for file, imports := range filesToImportsMapForFile {
			if file == filename {
				continue
			}

			fileOutput += fmt.Sprintf("import { %s } from './%s';\n", strings.Join(sortMapKeys(imports), ", "), strings.TrimSuffix(file, ".ts"))
		}

		indexOutput += fmt.Sprintf("export * from './%s';\n", strings.TrimSuffix(filename, ".ts"))
		finalOutput := fileOutput + "\n" + output

		if err := os.WriteFile(fmt.Sprintf("%s/%s", destinationDirectory, filename), []byte(finalOutput), 0o600); err != nil {
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

func sortMapKeys(m map[string]struct{}) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, strings.TrimPrefix(k, "./"))
	}

	sort.Strings(keys)

	return keys
}

var (
	// Times I've tried to optimize this regex before realizing it already accounts for
	// every edge case and there is no value (in either performance or readability terms)
	// in making it smaller: 1.
	numberMatcherRegex = regexp.MustCompile(`((u)?int(8|16|32|64)?|float(32|64))`)
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

func typescriptInterface[T any](x T) (out string, imports []string, err error) {
	typ := reflect.TypeOf(x)
	typeName := typ.Name()
	fieldsForType := reflect.VisibleFields(typ)

	output := fmt.Sprintf("export interface I%s {\n", typeName)

	importedTypes := []string{}
	for i := range fieldsForType {
		field := fieldsForType[i]
		if field.Name == "_" {
			continue
		}

		fieldName := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		if fieldName == "-" {
			continue
		}

		fieldType := strings.Replace(strings.TrimPrefix(strings.Replace(field.Type.String(), "[]", "", 1), "*"), "types.", "", 1)
		isPointer := field.Type.Kind() == reflect.Ptr
		isSlice := field.Type.Kind() == reflect.Slice

		if fieldType == "UserLoginInput" {
			continue
		}

		if isCustomType(fieldType) {
			importedTypes = append(importedTypes, fieldType)
		}

		switch fieldType {
		case timeType:
			fieldType = stringType
		case mapStringToBoolType:
			fieldType = "Record<string, boolean>"
		case boolType:
			fieldType = "boolean"
		}

		if numberMatcherRegex.MatchString(field.Type.String()) {
			fieldType = "number"
		}

		if t, ok := customTypeMap[fmt.Sprintf("%s.%s", typ.Name(), fieldName)]; ok {
			fieldType = t
			importedTypes = append(importedTypes, t)
		}

		line := CodeLine{
			FieldType: fieldType,
			FieldName: fieldName,
			IsPointer: isPointer,
			IsSlice:   isSlice,
		}

		tmpl := template.Must(template.New("").Parse(`	{{.FieldName}}{{if .IsPointer}}?{{end}}: {{if not .IsPointer}}NonNullable<{{end}}{{if .IsSlice}}Array<{{end}}{{.FieldType}}{{if .IsSlice}}>{{end -}}{{if not .IsPointer}}>{{end -}};` + "\n"))

		var b bytes.Buffer
		if err = tmpl.Execute(&b, line); err != nil {
			return "", nil, nil
		}

		output += b.String()
	}

	output += "}\n"

	return output, importedTypes, nil
}

func typescriptClass[T any](x T) (out string, imports []string, err error) {
	typ := reflect.TypeOf(x)
	fieldsForType := reflect.VisibleFields(typ)

	output := fmt.Sprintf("export class %s implements I%s {\n", typ.Name(), typ.Name())
	importedTypes := []string{}

	parsedLines := []CodeLine{}
	for i := range fieldsForType {
		field := fieldsForType[i]
		if field.Name == "_" {
			continue
		}

		fieldName := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		if fieldName == "-" {
			continue
		}

		fieldType := strings.Replace(strings.TrimPrefix(strings.Replace(field.Type.String(), "[]", "", 1), "*"), "types.", "", 1)
		isPointer := field.Type.Kind() == reflect.Ptr
		isSlice := field.Type.Kind() == reflect.Slice
		defaultValue := ""
		customType := isCustomType(fieldType)

		if fieldType == "UserLoginInput" {
			continue
		}

		if isSlice {
			defaultValue = "[]"
		}

		if isCustomType(fieldType) {
			importedTypes = append(importedTypes, fieldType)
		}

		switch fieldType {
		case stringType:
			if !isSlice {
				defaultValue = `''`
				if isPointer {
					defaultValue = ""
				}
			}
		case mapStringToBoolType:
			fieldType = "Record<string, boolean>"
			defaultValue = "{}"
		case timeType:
			fieldType = stringType
			if !isPointer {
				defaultValue = "'1970-01-01T00:00:00Z'"
			}
		case boolType:
			fieldType = "boolean"
			if !isSlice {
				defaultValue = "false"
			}
		}

		if numberMatcherRegex.MatchString(field.Type.String()) {
			fieldType = "number"
			if !isPointer && !isSlice {
				defaultValue = "0"
			}
		}

		if customType && !isSlice {
			defaultValue = fmt.Sprintf("new %s()", fieldType)
		}

		if t, ok := customTypeMap[fmt.Sprintf("%s.%s", typ.Name(), fieldName)]; ok {
			fieldType = t
			importedTypes = append(importedTypes, t)
			defaultValue = defaultValues[t]
		}

		line := CodeLine{
			FieldType:    fieldType,
			FieldName:    fieldName,
			IsPointer:    isPointer,
			IsSlice:      isSlice,
			DefaultValue: defaultValue,
			CustomType:   customType,
		}

		tmpl := template.Must(template.New("").Parse(`	{{.FieldName}}{{if .IsPointer}}?{{end}}: {{if not .IsPointer}}NonNullable<{{end}}{{if .IsSlice}}Array<{{end}}{{.FieldType}}{{if .IsSlice}}>{{end -}}{{if not .IsPointer}}>{{end -}}{{ if ne .DefaultValue "" }} = {{ .DefaultValue }}{{ end -}};` + "\n"))

		var b bytes.Buffer
		if err = tmpl.Execute(&b, line); err != nil {
			return "", nil, nil
		}

		thisLine := b.String()
		output += thisLine
		parsedLines = append(parsedLines, line)
	}

	output += `
	constructor(input: {
`

	for i := range parsedLines {
		line := parsedLines[i]

		sliceDecl := ""
		if line.IsSlice {
			sliceDecl = "[]"
		}

		output += fmt.Sprintf("		%s?: %s%s\n", line.FieldName, line.FieldType, sliceDecl)
	}

	output += `	} = {}) {
`

	for i := range parsedLines {
		line := parsedLines[i]

		dv := ""
		if line.DefaultValue != "" {
			dv = fmt.Sprintf(" ?? %s", line.DefaultValue)
		}

		output += fmt.Sprintf("		this.%s = input.%s%s;\n", line.FieldName, line.FieldName, dv)
	}

	output += `	}
}
`

	return output, importedTypes, nil
}

func isCustomType(x string) bool {
	switch x {
	case "int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"float32",
		"float64",
		boolType,
		mapStringToBoolType,
		timeType,
		stringType:
		return false
	default:
		return true
	}
}
