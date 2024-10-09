// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepCreationRequestInput } from './RecipePrepTaskStepCreationRequestInput';
import { NumberRangeWithOptionalMax, NumberRange } from './number_range';

export interface IRecipePrepTaskCreationRequestInput {
  description: string;
  explicitStorageInstructions: string;
  name: string;
  optional: boolean;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  belongsToRecipe: string;
  notes: string;
  recipeSteps: RecipePrepTaskStepCreationRequestInput;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
}

export class RecipePrepTaskCreationRequestInput implements IRecipePrepTaskCreationRequestInput {
  description: string;
  explicitStorageInstructions: string;
  name: string;
  optional: boolean;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  belongsToRecipe: string;
  notes: string;
  recipeSteps: RecipePrepTaskStepCreationRequestInput;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  constructor(input: Partial<RecipePrepTaskCreationRequestInput> = {}) {
    this.description = input.description = '';
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.name = input.name = '';
    this.optional = input.optional = false;
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.notes = input.notes = '';
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepCreationRequestInput();
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.storageType = input.storageType = '';
  }
}
