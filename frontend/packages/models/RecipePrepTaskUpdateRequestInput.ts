// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepUpdateRequestInput } from './RecipePrepTaskStepUpdateRequestInput';
import { NumberRange, OptionalNumberRange } from './number_range';

export interface IRecipePrepTaskUpdateRequestInput {
  belongsToRecipe: string;
  explicitStorageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput[];
  storageType: string;
  description: string;
  name: string;
  notes: string;
  optional: boolean;
}

export class RecipePrepTaskUpdateRequestInput implements IRecipePrepTaskUpdateRequestInput {
  belongsToRecipe: string;
  explicitStorageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput[];
  storageType: string;
  description: string;
  name: string;
  notes: string;
  optional: boolean;
  constructor(input: Partial<RecipePrepTaskUpdateRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.explicitStorageInstructions = input.explicitStorageInstructions || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || {};
    this.recipeSteps = input.recipeSteps || [];
    this.storageType = input.storageType || '';
    this.description = input.description || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.optional = input.optional || false;
  }
}
