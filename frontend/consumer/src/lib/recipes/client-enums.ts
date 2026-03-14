/**
 * Client-safe enum constants for recipe creation.
 * Do NOT import from mealplanning_messages or mealplanning_service_types in client code—
 * they transitively import uploaded_media which pulls in @grpc/grpc-js (Node-only).
 */
export const RecipeStepProductType = {
	RECIPE_STEP_PRODUCT_TYPE_INGREDIENT: 0,
	RECIPE_STEP_PRODUCT_TYPE_INSTRUMENT: 1,
	RECIPE_STEP_PRODUCT_TYPE_VESSEL: 2,
	UNRECOGNIZED: -1
} as const;

export type RecipeStepProductType =
	(typeof RecipeStepProductType)[keyof typeof RecipeStepProductType];

export const MealComponentType = {
	MEAL_COMPONENT_TYPE_MAIN: 4
} as const;

export type MealComponentType =
	(typeof MealComponentType)[keyof typeof MealComponentType];
