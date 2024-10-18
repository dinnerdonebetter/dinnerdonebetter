// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation.gen';
import { NumberRange } from './number_range.gen';

export interface IRecipeStepUpdateRequestInput {
  belongsToRecipe: string;
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  index: number;
  notes: string;
  optional: boolean;
  preparation: ValidPreparation;
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
}

export class RecipeStepUpdateRequestInput implements IRecipeStepUpdateRequestInput {
  belongsToRecipe: string;
  conditionExpression: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  index: number;
  notes: string;
  optional: boolean;
  preparation: ValidPreparation;
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipeStepUpdateRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.conditionExpression = input.conditionExpression || '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds || { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions || '';
    this.index = input.index || 0;
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.preparation = input.preparation || new ValidPreparation();
    this.startTimerAutomatically = input.startTimerAutomatically || false;
    this.temperatureInCelsius = input.temperatureInCelsius || { min: 0, max: 0 };
  }
}
