/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipeMedia } from './RecipeMedia';
import type { RecipeStepCompletionCondition } from './RecipeStepCompletionCondition';
import type { RecipeStepIngredient } from './RecipeStepIngredient';
import type { RecipeStepInstrument } from './RecipeStepInstrument';
import type { RecipeStepProduct } from './RecipeStepProduct';
import type { RecipeStepVessel } from './RecipeStepVessel';
import type { ValidPreparation } from './ValidPreparation';
export type RecipeStep = {
  archivedAt?: string;
  belongsToRecipe?: string;
  completionConditions?: Array<RecipeStepCompletionCondition>;
  conditionExpression?: string;
  createdAt?: string;
  explicitInstructions?: string;
  id?: string;
  index?: number;
  ingredients?: Array<RecipeStepIngredient>;
  instruments?: Array<RecipeStepInstrument>;
  lastUpdatedAt?: string;
  maximumEstimatedTimeInSeconds?: number;
  maximumTemperatureInCelsius?: number;
  media?: Array<RecipeMedia>;
  minimumEstimatedTimeInSeconds?: number;
  minimumTemperatureInCelsius?: number;
  notes?: string;
  optional?: boolean;
  preparation?: ValidPreparation;
  products?: Array<RecipeStepProduct>;
  startTimerAutomatically?: boolean;
  vessels?: Array<RecipeStepVessel>;
};
