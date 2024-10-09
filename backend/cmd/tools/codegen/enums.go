package codegen

type Enum struct {
	Name,
	ConstName,
	Comment string
	Values []string
}

var (
	EnumTypeMap = map[string]string{
		"ValidIngredientState.attributeType":                     "ValidIngredientStateAttributeType",
		"ValidIngredientStateCreationRequestInput.attributeType": "ValidIngredientStateAttributeType",
		"ValidIngredientStateUpdateRequestInput.attributeType":   "ValidIngredientStateAttributeType",
		"RecipeStepProduct.type":                                 "ValidRecipeStepProductType",
		"RecipeStepProductCreationRequestInput.type":             "ValidRecipeStepProductType",
		"RecipeStepProductUpdateRequestInput.type":               "ValidRecipeStepProductType",
		"MealPlanGroceryListItem.status":                         "ValidMealPlanGroceryListItemStatus",
		"MealPlanGroceryListItemCreationRequestInput.status":     "ValidMealPlanGroceryListItemStatus",
		"MealPlanGroceryListItemUpdateRequestInput.status":       "ValidMealPlanGroceryListItemStatus",
		"MealComponent.componentType":                            "MealComponentType",
		"MealComponentCreationRequestInput.componentType":        "MealComponentType",
		"MealComponentUpdateRequestInput.componentType":          "MealComponentType",
		"MealPlanTaskStatusChangeRequestInput.status":            "MealPlanTaskStatus",
		"MealPlanTask.status":                                    "MealPlanTaskStatus",
		"MealPlanTaskCreationRequestInput.status":                "MealPlanTaskStatus",
		"MealPlanTaskUpdateRequestInput.status":                  "MealPlanTaskStatus",
		"MealPlan.status":                                        "ValidMealPlanStatus",
		"MealPlanCreationRequestInput.status":                    "ValidMealPlanStatus",
		"MealPlanUpdateRequestInput.status":                      "ValidMealPlanStatus",
		"MealPlan.electionMethod":                                "ValidMealPlanElectionMethod",
		"MealPlanUpdateRequestInput.electionMethod":              "ValidMealPlanElectionMethod",
		"MealPlanCreationRequestInput.electionMethod":            "ValidMealPlanElectionMethod",
		"ValidVessel.shape":                                      "ValidVesselShapeType",
		"ValidVesselUpdateRequestInput.shape":                    "ValidVesselShapeType",
		"ValidVesselCreationRequestInput.shape":                  "ValidVesselShapeType",
	}

	DefaultEnumValues = map[string]string{
		"ValidMealPlanStatus":                "'awaiting_votes'",
		"ValidMealPlanGroceryListItemStatus": "'unknown'",
		"ValidMealPlanElectionMethod":        "'schulze'",
		"ValidIngredientStateAttributeType":  "'other'",
		"ValidRecipeStepProductType":         "'ingredient'",
		"MealComponentType":                  "'unspecified'",
		"MealPlanTaskStatus":                 "'unfinished'",
		"ValidVesselShapeType":               "'other'",
	}

	// Enums is every defined union in the app. This list is unfortunately important.
	Enums = []*Enum{
		{
			Comment:   "ingredient state attribute types",
			Name:      "ValidIngredientStateAttributeType",
			ConstName: "VALID_INGREDIENT_STATE_ATTRIBUTE_TYPES",
			Values: []string{
				"texture",
				"consistency",
				"temperature",
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
			Comment:   "meal plan grocery list item statuses",
			Name:      "ValidMealPlanGroceryListItemStatus",
			ConstName: "VALID_MEAL_PLAN_GROCERY_LIST_ITEM_STATUSES",
			Values: []string{
				"unknown",
				"already owned",
				"needs",
				"unavailable",
				"acquired",
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
				"vessel",
			},
		},
		{
			Comment:   "meal component types",
			Name:      "MealComponentType",
			ConstName: "MEAL_COMPONENT_TYPES",
			Values: []string{
				"main",
				"side",
				"appetizer",
				"beverage",
				"dessert",
				"soup",
				"salad",
				"amuse-bouche",
				"unspecified",
			},
		},
		{
			Comment:   "meal plan task statuses",
			Name:      "MealPlanTaskStatus",
			ConstName: "MEAL_PLAN_TASK_STATUSES",
			Values: []string{
				"unfinished",
				"postponed",
				"ignored",
				"canceled",
				"finished",
			},
		},
		{
			Comment:   "valid vessel shape types",
			Name:      "ValidVesselShapeType",
			ConstName: "VALID_VESSEL_SHAPE_TYPES",
			Values: []string{
				"hemisphere",
				"rectangle",
				"cone",
				"pyramid",
				"cylinder",
				"sphere",
				"cube",
				"other",
			},
		},
	}
)
