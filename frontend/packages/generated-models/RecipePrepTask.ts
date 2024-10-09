// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStep } from './RecipePrepTaskStep';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range';

export interface IRecipePrepTask {
  storageTemperatureInCelsius: NumberRange;
  explicitStorageInstructions: string;
  lastUpdatedAt?: string;
  name: string;
  storageType: string;
  archivedAt?: string;
  id: string;
  notes: string;
  createdAt: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStep;
  belongsToRecipe: string;
  description: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
}

export class RecipePrepTask implements IRecipePrepTask {
  storageTemperatureInCelsius: NumberRange;
  explicitStorageInstructions: string;
  lastUpdatedAt?: string;
  name: string;
  storageType: string;
  archivedAt?: string;
  id: string;
  notes: string;
  createdAt: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStep;
  belongsToRecipe: string;
  description: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  constructor(input: Partial<RecipePrepTask> = {}) {
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.storageType = input.storageType = '';
    this.archivedAt = input.archivedAt;
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.createdAt = input.createdAt = '';
    this.optional = input.optional = false;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStep();
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.description = input.description = '';
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
  }
}
