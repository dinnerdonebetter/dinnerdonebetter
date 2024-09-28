/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { RecipeStepCompletionConditionCreationRequestInput } from './RecipeStepCompletionConditionCreationRequestInput';
import type { RecipeStepIngredientCreationRequestInput } from './RecipeStepIngredientCreationRequestInput';
import type { RecipeStepInstrumentCreationRequestInput } from './RecipeStepInstrumentCreationRequestInput';
import type { RecipeStepProductCreationRequestInput } from './RecipeStepProductCreationRequestInput';
import type { RecipeStepVesselCreationRequestInput } from './RecipeStepVesselCreationRequestInput';
export type RecipeStepCreationRequestInput = {
  completionConditions?: Array<RecipeStepCompletionConditionCreationRequestInput>;
  conditionExpression?: string;
  explicitInstructions?: string;
  index?: number;
  ingredients?: Array<RecipeStepIngredientCreationRequestInput>;
  instruments?: Array<RecipeStepInstrumentCreationRequestInput>;
  maximumEstimatedTimeInSeconds?: number;
  maximumTemperatureInCelsius?: number;
  minimumEstimatedTimeInSeconds?: number;
  minimumTemperatureInCelsius?: number;
  notes?: string;
  optional?: boolean;
  preparationID?: string;
  products?: Array<RecipeStepProductCreationRequestInput>;
  startTimerAutomatically?: boolean;
  vessels?: Array<RecipeStepVesselCreationRequestInput>;
};
