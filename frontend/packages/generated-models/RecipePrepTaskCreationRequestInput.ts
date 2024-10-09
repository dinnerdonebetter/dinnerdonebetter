// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepCreationRequestInput } from './RecipePrepTaskStepCreationRequestInput';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range';

export interface IRecipePrepTaskCreationRequestInput {
  explicitStorageInstructions: string;
  recipeSteps: RecipePrepTaskStepCreationRequestInput;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  belongsToRecipe: string;
  description: string;
  optional: boolean;
  storageType: string;
  name: string;
  notes: string;
}

export class RecipePrepTaskCreationRequestInput implements IRecipePrepTaskCreationRequestInput {
  explicitStorageInstructions: string;
  recipeSteps: RecipePrepTaskStepCreationRequestInput;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  belongsToRecipe: string;
  description: string;
  optional: boolean;
  storageType: string;
  name: string;
  notes: string;
  constructor(input: Partial<RecipePrepTaskCreationRequestInput> = {}) {
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepCreationRequestInput();
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.description = input.description = '';
    this.optional = input.optional = false;
    this.storageType = input.storageType = '';
    this.name = input.name = '';
    this.notes = input.notes = '';
  }
}
