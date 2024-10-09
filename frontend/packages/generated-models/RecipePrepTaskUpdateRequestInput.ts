// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStepUpdateRequestInput } from './RecipePrepTaskStepUpdateRequestInput';
import { OptionalNumberRange, NumberRange } from './number_range';

export interface IRecipePrepTaskUpdateRequestInput {
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  belongsToRecipe?: string;
  explicitStorageInstructions?: string;
  optional?: boolean;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput;
  storageTemperatureInCelsius: NumberRange;
  storageType?: string;
  description?: string;
  name?: string;
  notes?: string;
}

export class RecipePrepTaskUpdateRequestInput implements IRecipePrepTaskUpdateRequestInput {
  timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
  belongsToRecipe?: string;
  explicitStorageInstructions?: string;
  optional?: boolean;
  recipeSteps: RecipePrepTaskStepUpdateRequestInput;
  storageTemperatureInCelsius: NumberRange;
  storageType?: string;
  description?: string;
  name?: string;
  notes?: string;
  constructor(input: Partial<RecipePrepTaskUpdateRequestInput> = {}) {
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = {};
    this.belongsToRecipe = input.belongsToRecipe;
    this.explicitStorageInstructions = input.explicitStorageInstructions;
    this.optional = input.optional;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepUpdateRequestInput();
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.storageType = input.storageType;
    this.description = input.description;
    this.name = input.name;
    this.notes = input.notes;
  }
}
