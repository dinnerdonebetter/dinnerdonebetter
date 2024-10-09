// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { NumberRange } from './number_range';

export interface IRecipeStepUpdateRequestInput {
  explicitInstructions?: string;
  index?: number;
  optional?: boolean;
  startTimerAutomatically?: boolean;
  temperatureInCelsius: NumberRange;
  belongsToRecipe: string;
  conditionExpression?: string;
  estimatedTimeInSeconds: NumberRange;
  notes?: string;
  preparation?: ValidPreparation;
}

export class RecipeStepUpdateRequestInput implements IRecipeStepUpdateRequestInput {
  explicitInstructions?: string;
  index?: number;
  optional?: boolean;
  startTimerAutomatically?: boolean;
  temperatureInCelsius: NumberRange;
  belongsToRecipe: string;
  conditionExpression?: string;
  estimatedTimeInSeconds: NumberRange;
  notes?: string;
  preparation?: ValidPreparation;
  constructor(input: Partial<RecipeStepUpdateRequestInput> = {}) {
    this.explicitInstructions = input.explicitInstructions;
    this.index = input.index;
    this.optional = input.optional;
    this.startTimerAutomatically = input.startTimerAutomatically;
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.conditionExpression = input.conditionExpression;
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.notes = input.notes;
    this.preparation = input.preparation;
  }
}
