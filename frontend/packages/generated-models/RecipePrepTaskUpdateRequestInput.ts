// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepUpdateRequestInput } from './RecipePrepTaskStepUpdateRequestInput';
import { NumberRange, OptionalNumberRange } from './number_range';

export interface IRecipePrepTaskUpdateRequestInput {
  storageTemperatureInCelsius: NumberRange;
  description?: string;
  explicitStorageInstructions?: string;
  notes?: string;
  optional?: boolean;
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  belongsToRecipe?: string;
  name?: string;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput;
  storageType?: string;
}

export class RecipePrepTaskUpdateRequestInput implements IRecipePrepTaskUpdateRequestInput {
  storageTemperatureInCelsius: NumberRange;
  description?: string;
  explicitStorageInstructions?: string;
  notes?: string;
  optional?: boolean;
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  belongsToRecipe?: string;
  name?: string;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput;
  storageType?: string;
  constructor(input: Partial<RecipePrepTaskUpdateRequestInput> = {}) {
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.description = input.description;
    this.explicitStorageInstructions = input.explicitStorageInstructions;
    this.notes = input.notes;
    this.optional = input.optional;
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = {};
    this.belongsToRecipe = input.belongsToRecipe;
    this.name = input.name;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepUpdateRequestInput();
    this.storageType = input.storageType;
  }
}
