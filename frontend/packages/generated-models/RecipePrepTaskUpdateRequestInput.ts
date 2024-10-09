// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepUpdateRequestInput } from './RecipePrepTaskStepUpdateRequestInput';
import { NumberRange, OptionalNumberRange } from './number_range';

export interface IRecipePrepTaskUpdateRequestInput {
  explicitStorageInstructions?: string;
  optional?: boolean;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput;
  notes?: string;
  storageTemperatureInCelsius: NumberRange;
  storageType?: string;
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  belongsToRecipe?: string;
  description?: string;
  name?: string;
}

export class RecipePrepTaskUpdateRequestInput implements IRecipePrepTaskUpdateRequestInput {
  explicitStorageInstructions?: string;
  optional?: boolean;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput;
  notes?: string;
  storageTemperatureInCelsius: NumberRange;
  storageType?: string;
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  belongsToRecipe?: string;
  description?: string;
  name?: string;
  constructor(input: Partial<RecipePrepTaskUpdateRequestInput> = {}) {
    this.explicitStorageInstructions = input.explicitStorageInstructions;
    this.optional = input.optional;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepUpdateRequestInput();
    this.notes = input.notes;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.storageType = input.storageType;
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = {};
    this.belongsToRecipe = input.belongsToRecipe;
    this.description = input.description;
    this.name = input.name;
  }
}
