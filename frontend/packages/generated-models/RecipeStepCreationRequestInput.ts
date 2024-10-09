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
  notes: string;
  startTimerAutomatically: boolean;
  conditionExpression: string;
  ingredients: RecipeStepIngredientCreationRequestInput;
  products: RecipeStepProductCreationRequestInput;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput;
  estimatedTimeInSeconds: NumberRange;
  instruments: RecipeStepInstrumentCreationRequestInput;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVesselCreationRequestInput;
  optional: boolean;
  preparationID: string;
}

export class RecipeStepCreationRequestInput implements IRecipeStepCreationRequestInput {
  explicitInstructions: string;
  index: number;
  notes: string;
  startTimerAutomatically: boolean;
  conditionExpression: string;
  ingredients: RecipeStepIngredientCreationRequestInput;
  products: RecipeStepProductCreationRequestInput;
  completionConditions: RecipeStepCompletionConditionCreationRequestInput;
  estimatedTimeInSeconds: NumberRange;
  instruments: RecipeStepInstrumentCreationRequestInput;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVesselCreationRequestInput;
  optional: boolean;
  preparationID: string;
  constructor(input: Partial<RecipeStepCreationRequestInput> = {}) {
    this.explicitInstructions = input.explicitInstructions = '';
    this.index = input.index = 0;
    this.notes = input.notes = '';
    this.startTimerAutomatically = input.startTimerAutomatically = false;
    this.conditionExpression = input.conditionExpression = '';
    this.ingredients = input.ingredients = new RecipeStepIngredientCreationRequestInput();
    this.products = input.products = new RecipeStepProductCreationRequestInput();
    this.completionConditions = input.completionConditions = new RecipeStepCompletionConditionCreationRequestInput();
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.instruments = input.instruments = new RecipeStepInstrumentCreationRequestInput();
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
    this.vessels = input.vessels = new RecipeStepVesselCreationRequestInput();
    this.optional = input.optional = false;
    this.preparationID = input.preparationID = '';
  }
}
