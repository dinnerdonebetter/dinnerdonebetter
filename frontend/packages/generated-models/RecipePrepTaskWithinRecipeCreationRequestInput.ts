// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepWithinRecipeCreationRequestInput } from './RecipePrepTaskStepWithinRecipeCreationRequestInput';
import { NumberRangeWithOptionalMax, NumberRange } from './number_range';

export interface IRecipePrepTaskWithinRecipeCreationRequestInput {
  explicitStorageInstructions: string;
  name: string;
  optional: boolean;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  belongsToRecipe: string;
  description: string;
  notes: string;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
}

export class RecipePrepTaskWithinRecipeCreationRequestInput implements IRecipePrepTaskWithinRecipeCreationRequestInput {
  explicitStorageInstructions: string;
  name: string;
  optional: boolean;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  belongsToRecipe: string;
  description: string;
  notes: string;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
  constructor(input: Partial<RecipePrepTaskWithinRecipeCreationRequestInput> = {}) {
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.name = input.name = '';
    this.optional = input.optional = false;
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.storageType = input.storageType = '';
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.description = input.description = '';
    this.notes = input.notes = '';
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepWithinRecipeCreationRequestInput();
  }
}
