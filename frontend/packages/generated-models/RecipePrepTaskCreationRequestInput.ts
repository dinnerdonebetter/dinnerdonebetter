// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepCreationRequestInput } from './RecipePrepTaskStepCreationRequestInput';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range';

export interface IRecipePrepTaskCreationRequestInput {
  explicitStorageInstructions: string;
  name: string;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  belongsToRecipe: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepCreationRequestInput;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
}

export class RecipePrepTaskCreationRequestInput implements IRecipePrepTaskCreationRequestInput {
  explicitStorageInstructions: string;
  name: string;
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  belongsToRecipe: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStepCreationRequestInput;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
  constructor(input: Partial<RecipePrepTaskCreationRequestInput> = {}) {
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.name = input.name = '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.storageType = input.storageType = '';
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.notes = input.notes = '';
    this.optional = input.optional = false;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepCreationRequestInput();
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
    this.description = input.description = '';
  }
}
