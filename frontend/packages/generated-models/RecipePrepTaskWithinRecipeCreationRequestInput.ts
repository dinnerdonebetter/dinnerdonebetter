// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepWithinRecipeCreationRequestInput } from './RecipePrepTaskStepWithinRecipeCreationRequestInput';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range';

export interface IRecipePrepTaskWithinRecipeCreationRequestInput {
  explicitStorageInstructions: string;
  name: string;
  notes: string;
  optional: boolean;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  belongsToRecipe: string;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
}

export class RecipePrepTaskWithinRecipeCreationRequestInput implements IRecipePrepTaskWithinRecipeCreationRequestInput {
  explicitStorageInstructions: string;
  name: string;
  notes: string;
  optional: boolean;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  belongsToRecipe: string;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
  constructor(input: Partial<RecipePrepTaskWithinRecipeCreationRequestInput> = {}) {
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.name = input.name = '';
    this.notes = input.notes = '';
    this.optional = input.optional = false;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.storageType = input.storageType = '';
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepWithinRecipeCreationRequestInput();
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
    this.description = input.description = '';
  }
}
