// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { NumberRange } from './number_range';

export interface IRecipeStepUpdateRequestInput {
  notes?: string;
  optional?: boolean;
  preparation?: ValidPreparation;
  temperatureInCelsius: NumberRange;
  belongsToRecipe: string;
  conditionExpression?: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions?: string;
  index?: number;
  startTimerAutomatically?: boolean;
}

export class RecipeStepUpdateRequestInput implements IRecipeStepUpdateRequestInput {
  notes?: string;
  optional?: boolean;
  preparation?: ValidPreparation;
  temperatureInCelsius: NumberRange;
  belongsToRecipe: string;
  conditionExpression?: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions?: string;
  index?: number;
  startTimerAutomatically?: boolean;
  constructor(input: Partial<RecipeStepUpdateRequestInput> = {}) {
    this.notes = input.notes;
    this.optional = input.optional;
    this.preparation = input.preparation;
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.conditionExpression = input.conditionExpression;
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions;
    this.index = input.index;
    this.startTimerAutomatically = input.startTimerAutomatically;
  }
}
