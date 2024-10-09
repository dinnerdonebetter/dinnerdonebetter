// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { NumberRange } from './number_range';

export interface IRecipeStepUpdateRequestInput {
  belongsToRecipe: string;
  index?: number;
  notes?: string;
  preparation?: ValidPreparation;
  conditionExpression?: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions?: string;
  optional?: boolean;
  startTimerAutomatically?: boolean;
  temperatureInCelsius: NumberRange;
}

export class RecipeStepUpdateRequestInput implements IRecipeStepUpdateRequestInput {
  belongsToRecipe: string;
  index?: number;
  notes?: string;
  preparation?: ValidPreparation;
  conditionExpression?: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions?: string;
  optional?: boolean;
  startTimerAutomatically?: boolean;
  temperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipeStepUpdateRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.index = input.index;
    this.notes = input.notes;
    this.preparation = input.preparation;
    this.conditionExpression = input.conditionExpression;
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions;
    this.optional = input.optional;
    this.startTimerAutomatically = input.startTimerAutomatically;
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
  }
}
