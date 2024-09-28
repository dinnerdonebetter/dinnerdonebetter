/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipePrepTaskStepWithinRecipeCreationRequestInput } from './RecipePrepTaskStepWithinRecipeCreationRequestInput';
export type RecipePrepTaskWithinRecipeCreationRequestInput = {
  belongsToRecipe?: string;
  description?: string;
  explicitStorageInstructions?: string;
  maximumStorageTemperatureInCelsius?: number;
  maximumTimeBufferBeforeRecipeInSeconds?: number;
  minimumStorageTemperatureInCelsius?: number;
  minimumTimeBufferBeforeRecipeInSeconds?: number;
  name?: string;
  notes?: string;
  optional?: boolean;
  recipeSteps?: Array<RecipePrepTaskStepWithinRecipeCreationRequestInput>;
  storageType?: string;
};
