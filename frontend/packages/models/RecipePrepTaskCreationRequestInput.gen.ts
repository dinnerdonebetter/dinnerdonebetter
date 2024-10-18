// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepCreationRequestInput } from './RecipePrepTaskStepCreationRequestInput';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range.gen';

export interface IRecipePrepTaskCreationRequestInput {
  belongsToRecipe: string;
  description: string;
  explicitStorageInstructions: string;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepCreationRequestInput[];
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
}

export class RecipePrepTaskCreationRequestInput implements IRecipePrepTaskCreationRequestInput {
  belongsToRecipe: string;
  description: string;
  explicitStorageInstructions: string;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepCreationRequestInput[];
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  constructor(input: Partial<RecipePrepTaskCreationRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.description = input.description || '';
    this.explicitStorageInstructions = input.explicitStorageInstructions || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.recipeSteps = input.recipeSteps || [];
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.storageType = input.storageType || '';
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || { min: 0 };
  }
}
