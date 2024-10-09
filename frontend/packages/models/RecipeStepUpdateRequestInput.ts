// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { NumberRange } from './number_range';

export interface IRecipeStepUpdateRequestInput {
  belongsToRecipe: string;
  optional: boolean;
  preparation: ValidPreparation;
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  index: number;
  notes: string;
}

export class RecipeStepUpdateRequestInput implements IRecipeStepUpdateRequestInput {
  belongsToRecipe: string;
  optional: boolean;
  preparation: ValidPreparation;
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  index: number;
  notes: string;
  constructor(input: Partial<RecipeStepUpdateRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.optional = input.optional || false;
    this.preparation = input.preparation || new ValidPreparation();
    this.startTimerAutomatically = input.startTimerAutomatically || false;
    this.temperatureInCelsius = input.temperatureInCelsius || { min: 0, max: 0 };
    this.conditionExpression = input.conditionExpression || '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds || { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions || '';
    this.index = input.index || 0;
    this.notes = input.notes || '';
  }
}
