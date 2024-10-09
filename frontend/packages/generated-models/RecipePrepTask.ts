// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStep } from './RecipePrepTaskStep';
import { NumberRangeWithOptionalMax, NumberRange } from './number_range';

export interface IRecipePrepTask {
  explicitStorageInstructions: string;
  id: string;
  description: string;
  name: string;
  optional: boolean;
  storageType: string;
  archivedAt?: string;
  belongsToRecipe: string;
  createdAt: string;
  lastUpdatedAt?: string;
  recipeSteps: RecipePrepTaskStep;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  notes: string;
  storageTemperatureInCelsius: NumberRange;
}

export class RecipePrepTask implements IRecipePrepTask {
  explicitStorageInstructions: string;
  id: string;
  description: string;
  name: string;
  optional: boolean;
  storageType: string;
  archivedAt?: string;
  belongsToRecipe: string;
  createdAt: string;
  lastUpdatedAt?: string;
  recipeSteps: RecipePrepTaskStep;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  notes: string;
  storageTemperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipePrepTask> = {}) {
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.id = input.id = '';
    this.description = input.description = '';
    this.name = input.name = '';
    this.optional = input.optional = false;
    this.storageType = input.storageType = '';
    this.archivedAt = input.archivedAt;
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStep();
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
    this.notes = input.notes = '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
  }
}
