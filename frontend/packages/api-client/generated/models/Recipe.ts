/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipeMedia } from './RecipeMedia';
import type { RecipePrepTask } from './RecipePrepTask';
import type { RecipeStep } from './RecipeStep';
export type Recipe = {
  archivedAt?: string;
  createdAt?: string;
  createdByUser?: string;
  description?: string;
  eligibleForMeals?: boolean;
  id?: string;
  inspiredByRecipeID?: string;
  lastUpdatedAt?: string;
  maximumEstimatedPortions?: number;
  media?: Array<RecipeMedia>;
  minimumEstimatedPortions?: number;
  name?: string;
  pluralPortionName?: string;
  portionName?: string;
  prepTasks?: Array<RecipePrepTask>;
  sealOfApproval?: boolean;
  slug?: string;
  source?: string;
  steps?: Array<RecipeStep>;
  supportingRecipes?: Array<Recipe>;
  yieldsComponentType?: string;
};
