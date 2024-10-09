// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionCreationRequestInput } from './RecipeStepCompletionConditionCreationRequestInput';
import { RecipeStepIngredientCreationRequestInput } from './RecipeStepIngredientCreationRequestInput';
import { RecipeStepInstrumentCreationRequestInput } from './RecipeStepInstrumentCreationRequestInput';
import { RecipeStepProductCreationRequestInput } from './RecipeStepProductCreationRequestInput';
import { RecipeStepVesselCreationRequestInput } from './RecipeStepVesselCreationRequestInput';
import { NumberRange } from './number_range';

export interface IRecipeStepCreationRequestInput {
  temperatureInCelsius: NumberRange;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  ingredients: RecipeStepIngredientCreationRequestInput;
  optional: boolean;
  preparationID: string;
  vessels: RecipeStepVesselCreationRequestInput;
  index: number;
  notes: string;
  startTimerAutomatically: boolean;
  conditionExpression: string;
  instruments: RecipeStepInstrumentCreationRequestInput;
  products: RecipeStepProductCreationRequestInput;
}

export class RecipeStepCreationRequestInput implements IRecipeStepCreationRequestInput {
  temperatureInCelsius: NumberRange;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  ingredients: RecipeStepIngredientCreationRequestInput;
  optional: boolean;
  preparationID: string;
  vessels: RecipeStepVesselCreationRequestInput;
  index: number;
  notes: string;
  startTimerAutomatically: boolean;
  conditionExpression: string;
  instruments: RecipeStepInstrumentCreationRequestInput;
  products: RecipeStepProductCreationRequestInput;
  constructor(input: Partial<RecipeStepCreationRequestInput> = {}) {
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
    this.completionConditions = input.completionConditions = new RecipeStepCompletionConditionCreationRequestInput();
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions = '';
    this.ingredients = input.ingredients = new RecipeStepIngredientCreationRequestInput();
    this.optional = input.optional = false;
    this.preparationID = input.preparationID = '';
    this.vessels = input.vessels = new RecipeStepVesselCreationRequestInput();
    this.index = input.index = 0;
    this.notes = input.notes = '';
    this.startTimerAutomatically = input.startTimerAutomatically = false;
    this.conditionExpression = input.conditionExpression = '';
    this.instruments = input.instruments = new RecipeStepInstrumentCreationRequestInput();
    this.products = input.products = new RecipeStepProductCreationRequestInput();
  }
}
