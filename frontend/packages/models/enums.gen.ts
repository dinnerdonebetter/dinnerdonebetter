// GENERATED CODE, DO NOT EDIT MANUALLY

export const ALL_MEAL_COMPONENT_TYPE: string[] = [
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
type MealComponentTypeTuple = typeof ALL_MEAL_COMPONENT_TYPE;
export type MealComponentType = MealComponentTypeTuple[number];

export const ALL_MEAL_PLAN_TASK_STATUS: string[] = ['unfinished', 'postponed', 'ignored', 'canceled', 'finished'];
type MealPlanTaskStatusTuple = typeof ALL_MEAL_PLAN_TASK_STATUS;
export type MealPlanTaskStatus = MealPlanTaskStatusTuple[number];

export const ALL_VALID_INGREDIENT_STATE_ATTRIBUTE_TYPE: string[] = [
  'texture',
  'consistency',
  'temperature',
  'color',
  'appearance',
  'odor',
  'taste',
  'sound',
  'other',
];
type ValidIngredientStateAttributeTypeTuple = typeof ALL_VALID_INGREDIENT_STATE_ATTRIBUTE_TYPE;
export type ValidIngredientStateAttributeType = ValidIngredientStateAttributeTypeTuple[number];

export const ALL_VALID_MEAL_PLAN_ELECTION_METHOD: string[] = ['schulze', 'instant-runoff'];
type ValidMealPlanElectionMethodTuple = typeof ALL_VALID_MEAL_PLAN_ELECTION_METHOD;
export type ValidMealPlanElectionMethod = ValidMealPlanElectionMethodTuple[number];

export const ALL_VALID_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS: string[] = [
  'unknown',
  'already owned',
  'needs',
  'unavailable',
  'acquired',
];
type ValidMealPlanGroceryListItemStatusTuple = typeof ALL_VALID_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS;
export type ValidMealPlanGroceryListItemStatus = ValidMealPlanGroceryListItemStatusTuple[number];

export const ALL_VALID_MEAL_PLAN_STATUS: string[] = ['awaiting_votes', 'finalized'];
type ValidMealPlanStatusTuple = typeof ALL_VALID_MEAL_PLAN_STATUS;
export type ValidMealPlanStatus = ValidMealPlanStatusTuple[number];

export const ALL_VALID_RECIPE_STEP_PRODUCT_TYPE: string[] = ['ingredient', 'instrument', 'vessel'];
type ValidRecipeStepProductTypeTuple = typeof ALL_VALID_RECIPE_STEP_PRODUCT_TYPE;
export type ValidRecipeStepProductType = ValidRecipeStepProductTypeTuple[number];

export const ALL_VALID_VESSEL_SHAPE_TYPE: string[] = [
  'hemisphere',
  'rectangle',
  'cone',
  'pyramid',
  'cylinder',
  'sphere',
  'cube',
  'other',
];
type ValidVesselShapeTypeTuple = typeof ALL_VALID_VESSEL_SHAPE_TYPE;
export type ValidVesselShapeType = ValidVesselShapeTypeTuple[number];
