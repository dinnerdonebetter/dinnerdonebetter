/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput } from './RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput';
export type RecipeStepCompletionConditionForExistingRecipeCreationRequestInput = {
  belongsToRecipeStep?: string;
  ingredientStateID?: string;
  ingredients?: Array<RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput>;
  notes?: string;
  optional?: boolean;
};
