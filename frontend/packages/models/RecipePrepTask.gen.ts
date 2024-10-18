// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskStep } from './RecipePrepTaskStep';
import { NumberRange, NumberRangeWithOptionalMax } from './number_range.gen';

export interface IRecipePrepTask {
  archivedAt: string;
  belongsToRecipe: string;
  createdAt: string;
  description: string;
  explicitStorageInstructions: string;
  id: string;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStep[];
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
}

export class RecipePrepTask implements IRecipePrepTask {
  archivedAt: string;
  belongsToRecipe: string;
  createdAt: string;
  description: string;
  explicitStorageInstructions: string;
  id: string;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  optional: boolean;
  recipeSteps: RecipePrepTaskStep[];
  storageTemperatureInCelsius: NumberRange;
  storageType: string;
  timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
  constructor(input: Partial<RecipePrepTask> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.createdAt = input.createdAt || '';
    this.description = input.description || '';
    this.explicitStorageInstructions = input.explicitStorageInstructions || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.recipeSteps = input.recipeSteps || [];
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.storageType = input.storageType || '';
    this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || { min: 0 };
  }
}
