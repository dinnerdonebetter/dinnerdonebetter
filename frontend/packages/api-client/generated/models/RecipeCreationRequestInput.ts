/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipePrepTaskWithinRecipeCreationRequestInput } from './RecipePrepTaskWithinRecipeCreationRequestInput';
import type { RecipeStepCreationRequestInput } from './RecipeStepCreationRequestInput';
export type RecipeCreationRequestInput = {
  alsoCreateMeal?: boolean;
  description?: string;
  eligibleForMeals?: boolean;
  inspiredByRecipeID?: string;
  maximumEstimatedPortions?: number;
  minimumEstimatedPortions?: number;
  name?: string;
  pluralPortionName?: string;
  portionName?: string;
  prepTasks?: Array<RecipePrepTaskWithinRecipeCreationRequestInput>;
  sealOfApproval?: boolean;
  slug?: string;
  source?: string;
  steps?: Array<RecipeStepCreationRequestInput>;
  yieldsComponentType?: string;
};
