// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepInstrument {
  preferenceRank: number;
  lastUpdatedAt: string;
  name: string;
  optional: boolean;
  id: string;
  instrument: ValidInstrument;
  notes: string;
  optionIndex: number;
  quantity: NumberRangeWithOptionalMax;
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  recipeStepProductID: string;
}

export class RecipeStepInstrument implements IRecipeStepInstrument {
  preferenceRank: number;
  lastUpdatedAt: string;
  name: string;
  optional: boolean;
  id: string;
  instrument: ValidInstrument;
  notes: string;
  optionIndex: number;
  quantity: NumberRangeWithOptionalMax;
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  recipeStepProductID: string;
  constructor(input: Partial<RecipeStepInstrument> = {}) {
    this.preferenceRank = input.preferenceRank || 0;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.optional = input.optional || false;
    this.id = input.id || '';
    this.instrument = input.instrument || new ValidInstrument();
    this.notes = input.notes || '';
    this.optionIndex = input.optionIndex || 0;
    this.quantity = input.quantity || { min: 0 };
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.recipeStepProductID = input.recipeStepProductID || '';
  }
}
