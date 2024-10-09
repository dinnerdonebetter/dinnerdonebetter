// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionCreationRequestInput } from './RecipeStepCompletionConditionCreationRequestInput';
import { RecipeStepIngredientCreationRequestInput } from './RecipeStepIngredientCreationRequestInput';
import { RecipeStepInstrumentCreationRequestInput } from './RecipeStepInstrumentCreationRequestInput';
import { RecipeStepProductCreationRequestInput } from './RecipeStepProductCreationRequestInput';
import { RecipeStepVesselCreationRequestInput } from './RecipeStepVesselCreationRequestInput';
import { NumberRange } from './number_range';

export interface IRecipeStepCreationRequestInput {
  conditionExpression: string;
  index: number;
  optional: boolean;
  instruments: RecipeStepInstrumentCreationRequestInput[];
  notes: string;
  vessels: RecipeStepVesselCreationRequestInput[];
  ingredients: RecipeStepIngredientCreationRequestInput[];
  preparationID: string;
  startTimerAutomatically: boolean;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput[];
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  products: RecipeStepProductCreationRequestInput[];
  temperatureInCelsius: NumberRange;
}

export class RecipeStepCreationRequestInput implements IRecipeStepCreationRequestInput {
  conditionExpression: string;
  index: number;
  optional: boolean;
  instruments: RecipeStepInstrumentCreationRequestInput[];
  notes: string;
  vessels: RecipeStepVesselCreationRequestInput[];
  ingredients: RecipeStepIngredientCreationRequestInput[];
  preparationID: string;
  startTimerAutomatically: boolean;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput[];
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  products: RecipeStepProductCreationRequestInput[];
  temperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipeStepCreationRequestInput> = {}) {
    this.conditionExpression = input.conditionExpression || '';
    this.index = input.index || 0;
    this.optional = input.optional || false;
    this.instruments = input.instruments || [];
    this.notes = input.notes || '';
    this.vessels = input.vessels || [];
    this.ingredients = input.ingredients || [];
    this.preparationID = input.preparationID || '';
    this.startTimerAutomatically = input.startTimerAutomatically || false;
    this.completionConditions = input.completionConditions || [];
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds || { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions || '';
    this.products = input.products || [];
    this.temperatureInCelsius = input.temperatureInCelsius || { min: 0, max: 0 };
  }
}
