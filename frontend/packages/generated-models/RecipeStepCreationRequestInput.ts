// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionCreationRequestInput } from './RecipeStepCompletionConditionCreationRequestInput';
import { RecipeStepIngredientCreationRequestInput } from './RecipeStepIngredientCreationRequestInput';
import { RecipeStepInstrumentCreationRequestInput } from './RecipeStepInstrumentCreationRequestInput';
import { RecipeStepProductCreationRequestInput } from './RecipeStepProductCreationRequestInput';
import { RecipeStepVesselCreationRequestInput } from './RecipeStepVesselCreationRequestInput';
import { NumberRange } from './number_range';

export interface IRecipeStepCreationRequestInput {
  explicitInstructions: string;
  index: number;
  optional: boolean;
  preparationID: string;
  products: RecipeStepProductCreationRequestInput;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput;
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  temperatureInCelsius: NumberRange;
  instruments: RecipeStepInstrumentCreationRequestInput;
  vessels: RecipeStepVesselCreationRequestInput;
  notes: string;
  ingredients: RecipeStepIngredientCreationRequestInput;
  startTimerAutomatically: boolean;
}

export class RecipeStepCreationRequestInput implements IRecipeStepCreationRequestInput {
  explicitInstructions: string;
  index: number;
  optional: boolean;
  preparationID: string;
  products: RecipeStepProductCreationRequestInput;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput;
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  temperatureInCelsius: NumberRange;
  instruments: RecipeStepInstrumentCreationRequestInput;
  vessels: RecipeStepVesselCreationRequestInput;
  notes: string;
  ingredients: RecipeStepIngredientCreationRequestInput;
  startTimerAutomatically: boolean;
  constructor(input: Partial<RecipeStepCreationRequestInput> = {}) {
    this.explicitInstructions = input.explicitInstructions = '';
    this.index = input.index = 0;
    this.optional = input.optional = false;
    this.preparationID = input.preparationID = '';
    this.products = input.products = new RecipeStepProductCreationRequestInput();
    this.completionConditions = input.completionConditions = new RecipeStepCompletionConditionCreationRequestInput();
    this.conditionExpression = input.conditionExpression = '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
    this.instruments = input.instruments = new RecipeStepInstrumentCreationRequestInput();
    this.vessels = input.vessels = new RecipeStepVesselCreationRequestInput();
    this.notes = input.notes = '';
    this.ingredients = input.ingredients = new RecipeStepIngredientCreationRequestInput();
    this.startTimerAutomatically = input.startTimerAutomatically = false;
  }
}
