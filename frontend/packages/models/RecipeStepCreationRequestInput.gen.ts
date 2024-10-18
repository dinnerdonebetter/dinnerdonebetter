// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionCreationRequestInput } from './RecipeStepCompletionConditionCreationRequestInput';
import { RecipeStepIngredientCreationRequestInput } from './RecipeStepIngredientCreationRequestInput';
import { RecipeStepInstrumentCreationRequestInput } from './RecipeStepInstrumentCreationRequestInput';
import { RecipeStepProductCreationRequestInput } from './RecipeStepProductCreationRequestInput';
import { RecipeStepVesselCreationRequestInput } from './RecipeStepVesselCreationRequestInput';
import { NumberRange } from './number_range.gen';

export interface IRecipeStepCreationRequestInput {
  completionConditions: RecipeStepCompletionConditionCreationRequestInput[];
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  index: number;
  ingredients: RecipeStepIngredientCreationRequestInput[];
  instruments: RecipeStepInstrumentCreationRequestInput[];
  notes: string;
  optional: boolean;
  preparationID: string;
  products: RecipeStepProductCreationRequestInput[];
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVesselCreationRequestInput[];
}

export class RecipeStepCreationRequestInput implements IRecipeStepCreationRequestInput {
  completionConditions: RecipeStepCompletionConditionCreationRequestInput[];
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  index: number;
  ingredients: RecipeStepIngredientCreationRequestInput[];
  instruments: RecipeStepInstrumentCreationRequestInput[];
  notes: string;
  optional: boolean;
  preparationID: string;
  products: RecipeStepProductCreationRequestInput[];
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVesselCreationRequestInput[];
  constructor(input: Partial<RecipeStepCreationRequestInput> = {}) {
    this.completionConditions = input.completionConditions || [];
    this.conditionExpression = input.conditionExpression || '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds || { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions || '';
    this.index = input.index || 0;
    this.ingredients = input.ingredients || [];
    this.instruments = input.instruments || [];
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.preparationID = input.preparationID || '';
    this.products = input.products || [];
    this.startTimerAutomatically = input.startTimerAutomatically || false;
    this.temperatureInCelsius = input.temperatureInCelsius || { min: 0, max: 0 };
    this.vessels = input.vessels || [];
  }
}
