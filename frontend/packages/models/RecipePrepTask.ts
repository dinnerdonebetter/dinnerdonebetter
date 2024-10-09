// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStep } from './RecipePrepTaskStep';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range';

export interface IRecipePrepTask {
  lastUpdatedAt: string;
  notes: string;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
  explicitStorageInstructions: string;
  storageType: string;
  id: string;
  archivedAt: string;
  belongsToRecipe: string;
  createdAt: string;
  name: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStep[];
}

export class RecipePrepTask implements IRecipePrepTask {
  lastUpdatedAt: string;
  notes: string;
  storageTemperatureInCelsius: NumberRange;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  description: string;
  explicitStorageInstructions: string;
  storageType: string;
  id: string;
  archivedAt: string;
  belongsToRecipe: string;
  createdAt: string;
  name: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStep[];
  constructor(input: Partial<RecipePrepTask> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || { min: 0 };
    this.description = input.description || '';
    this.explicitStorageInstructions = input.explicitStorageInstructions || '';
    this.storageType = input.storageType || '';
    this.id = input.id || '';
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.createdAt = input.createdAt || '';
    this.name = input.name || '';
    this.optional = input.optional || false;
    this.recipeSteps = input.recipeSteps || [];
  }
}
