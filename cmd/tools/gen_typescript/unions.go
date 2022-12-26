package main

import (
	"fmt"
)

type unionDefinition struct {
	Name,
	ConstName,
	Comment string
	Values []string
}

var (
	unions = []*unionDefinition{
		{
			Comment:   "ingredient state attribute types",
			Name:      "ValidIngredientStateAttributeType",
			ConstName: "VALID_INGREDIENT_STATE_ATTRIBUTE_TYPES",
			Values: []string{
				"texture",
				"consistency",
				"color",
				"appearance",
				"odor",
				"taste",
				"sound",
				"other",
			},
		},
		{
			Comment:   "meal plan statuses",
			Name:      "ValidMealPlanStatus",
			ConstName: "VALID_MEAL_PLAN_STATUSES",
			Values: []string{
				"awaiting_votes",
				"finalized",
			},
		},
		{
			Comment:   "meal plan election methods",
			Name:      "ValidMealPlanElectionMethod",
			ConstName: "VALID_MEAL_PLAN_ELECTION_METHODS",
			Values: []string{
				"schulze",
				"instant-runoff",
			},
		},
		{
			Comment:   "recipe step product types",
			Name:      "ValidRecipeStepProductType",
			ConstName: "RECIPE_STEP_PRODUCT_TYPES",
			Values: []string{
				"ingredient",
				"instrument",
			},
		},
	}
)

func buildUnionDeclaration(def *unionDefinition) string {
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

	output += fmt.Sprintf("type %sTypeTuple = typeof %s;\n", def.Name, def.ConstName)
	output += fmt.Sprintf("export type %s = %sTypeTuple[number];\n\n", def.Name, def.Name)

	return output
}

func buildUnionsFile() string {
	output := copyString(generatedDisclaimer)

	for _, def := range unions {
		output += buildUnionDeclaration(def)
	}

	return output
}

const (
	// TODO: generate this file programmatically.
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
export const ALL_VALID_MEAL_PLAN_STATUSES: string[] = [
  'awaiting_votes',
  'finalized',
];
type ValidMealPlanStatusTypeTuple = typeof ALL_VALID_MEAL_PLAN_STATUSES;
export type ValidMealPlanStatus = ValidMealPlanStatusTypeTuple[number];

// meal plan election methods
export const ALL_VALID_MEAL_PLAN_ELECTION_METHODS: string[] = [
  'schulze',
  'instant-runoff',
];
type ValidMealPlanElectionMethodTypeTuple = typeof ALL_VALID_MEAL_PLAN_ELECTION_METHODS;
export type ValidMealPlanElectionMethod = ValidMealPlanElectionMethodTypeTuple[number];

// recipe step product types
export const ALL_RECIPE_STEP_PRODUCT_TYPES = [
  'instrument',
  'ingredient',
];
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
