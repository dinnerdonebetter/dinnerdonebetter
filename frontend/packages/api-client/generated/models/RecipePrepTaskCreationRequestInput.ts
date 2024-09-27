/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipePrepTaskStepCreationRequestInput } from './RecipePrepTaskStepCreationRequestInput';
export type RecipePrepTaskCreationRequestInput = {
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
  recipeSteps?: Array<RecipePrepTaskStepCreationRequestInput>;
  storageType?: string;
};
