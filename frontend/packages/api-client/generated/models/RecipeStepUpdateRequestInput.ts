/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidPreparation } from './ValidPreparation';
export type RecipeStepUpdateRequestInput = {
  belongsToRecipe?: string;
  conditionExpression?: string;
  explicitInstructions?: string;
  index?: number;
  maximumEstimatedTimeInSeconds?: number;
  maximumTemperatureInCelsius?: number;
  minimumEstimatedTimeInSeconds?: number;
  minimumTemperatureInCelsius?: number;
  notes?: string;
  optional?: boolean;
  preparation?: ValidPreparation;
  startTimerAutomatically?: boolean;
};
