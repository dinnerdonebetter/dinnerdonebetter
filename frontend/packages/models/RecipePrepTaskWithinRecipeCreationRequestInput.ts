// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepWithinRecipeCreationRequestInput } from './RecipePrepTaskStepWithinRecipeCreationRequestInput';
import { NumberRangeWithOptionalMax, NumberRange } from './number_range';

export interface IRecipePrepTaskWithinRecipeCreationRequestInput {
  name: string;
  notes: string;
  optional: boolean;
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  belongsToRecipe: string;
  description: string;
  explicitStorageInstructions: string;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput[];
  storageTemperatureInCelsius: NumberRange;
}

export class RecipePrepTaskWithinRecipeCreationRequestInput implements IRecipePrepTaskWithinRecipeCreationRequestInput {
  name: string;
  notes: string;
  optional: boolean;
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  belongsToRecipe: string;
  description: string;
  explicitStorageInstructions: string;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput[];
  storageTemperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipePrepTaskWithinRecipeCreationRequestInput> = {}) {
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.storageType = input.storageType || '';
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || { min: 0 };
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.description = input.description || '';
    this.explicitStorageInstructions = input.explicitStorageInstructions || '';
    this.recipeSteps = input.recipeSteps || [];
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
  }
}
