/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipePrepTaskStep } from './RecipePrepTaskStep';
export type RecipePrepTask = {
  archivedAt?: string;
  belongsToRecipe?: string;
  createdAt?: string;
  description?: string;
  explicitStorageInstructions?: string;
  id?: string;
  lastUpdatedAt?: string;
  maximumStorageTemperatureInCelsius?: number;
  maximumTimeBufferBeforeRecipeInSeconds?: number;
  minimumStorageTemperatureInCelsius?: number;
  minimumTimeBufferBeforeRecipeInSeconds?: number;
  name?: string;
  notes?: string;
  optional?: boolean;
  recipeSteps?: Array<RecipePrepTaskStep>;
  storageType?: string;
};
