/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipePrepTaskStepUpdateRequestInput } from './RecipePrepTaskStepUpdateRequestInput';
export type RecipePrepTaskUpdateRequestInput = {
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
  recipeSteps?: Array<RecipePrepTaskStepUpdateRequestInput>;
  storageType?: string;
};
