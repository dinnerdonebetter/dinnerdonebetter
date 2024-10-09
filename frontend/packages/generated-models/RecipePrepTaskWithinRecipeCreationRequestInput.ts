// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepWithinRecipeCreationRequestInput } from './RecipePrepTaskStepWithinRecipeCreationRequestInput';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range';

export interface IRecipePrepTaskWithinRecipeCreationRequestInput {
  belongsToRecipe: string;
  description: string;
  explicitStorageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
  storageType: string;
}

export class RecipePrepTaskWithinRecipeCreationRequestInput implements IRecipePrepTaskWithinRecipeCreationRequestInput {
  belongsToRecipe: string;
  description: string;
  explicitStorageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
  storageType: string;
  constructor(input: Partial<RecipePrepTaskWithinRecipeCreationRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.description = input.description = '';
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
    this.name = input.name = '';
    this.notes = input.notes = '';
    this.optional = input.optional = false;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepWithinRecipeCreationRequestInput();
    this.storageType = input.storageType = '';
  }
}
