// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepCreationRequestInput } from './RecipePrepTaskStepCreationRequestInput';
import { NumberRangeWithOptionalMax, NumberRange } from './number_range';

export interface IRecipePrepTaskCreationRequestInput {
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
  explicitStorageInstructions: string;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepCreationRequestInput[];
  belongsToRecipe: string;
  storageTemperatureInCelsius: NumberRange;
}

export class RecipePrepTaskCreationRequestInput implements IRecipePrepTaskCreationRequestInput {
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
  explicitStorageInstructions: string;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepCreationRequestInput[];
  belongsToRecipe: string;
  storageTemperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipePrepTaskCreationRequestInput> = {}) {
    this.storageType = input.storageType || '';
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || { min: 0 };
    this.description = input.description || '';
    this.explicitStorageInstructions = input.explicitStorageInstructions || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.recipeSteps = input.recipeSteps || [];
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
  }
}
