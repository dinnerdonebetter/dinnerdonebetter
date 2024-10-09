// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStep } from './RecipePrepTaskStep';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range';

export interface IRecipePrepTask {
  belongsToRecipe: string;
  lastUpdatedAt?: string;
  recipeSteps: RecipePrepTaskStep;
  storageTemperatureInCelsius: NumberRange;
  archivedAt?: string;
  explicitStorageInstructions: string;
  name: string;
  createdAt: string;
  id: string;
  notes: string;
  optional: boolean;
  storageType: string;
  description: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
}

export class RecipePrepTask implements IRecipePrepTask {
  belongsToRecipe: string;
  lastUpdatedAt?: string;
  recipeSteps: RecipePrepTaskStep;
  storageTemperatureInCelsius: NumberRange;
  archivedAt?: string;
  explicitStorageInstructions: string;
  name: string;
  createdAt: string;
  id: string;
  notes: string;
  optional: boolean;
  storageType: string;
  description: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  constructor(input: Partial<RecipePrepTask> = {}) {
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.recipeSteps = input.recipeSteps = new RecipePrepTaskStep();
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.archivedAt = input.archivedAt;
    this.explicitStorageInstructions = input.explicitStorageInstructions = '';
    this.name = input.name = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.optional = input.optional = false;
    this.storageType = input.storageType = '';
    this.description = input.description = '';
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
  }
}
