/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipeStepCompletionConditionIngredient } from './RecipeStepCompletionConditionIngredient';
import type { ValidIngredientState } from './ValidIngredientState';
export type RecipeStepCompletionCondition = {
  archivedAt?: string;
  belongsToRecipeStep?: string;
  createdAt?: string;
  id?: string;
  ingredientState?: ValidIngredientState;
  ingredients?: Array<RecipeStepCompletionConditionIngredient>;
  lastUpdatedAt?: string;
  notes?: string;
  optional?: boolean;
};
