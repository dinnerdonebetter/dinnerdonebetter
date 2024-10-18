// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument.gen';
import { NumberRangeWithOptionalMax } from './number_range.gen';

export interface IRecipeStepInstrument {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  optionIndex: number;
  optional: boolean;
  preferenceRank: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
}

export class RecipeStepInstrument implements IRecipeStepInstrument {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  optionIndex: number;
  optional: boolean;
  preferenceRank: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
  constructor(input: Partial<RecipeStepInstrument> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.instrument = input.instrument || new ValidInstrument();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.optionIndex = input.optionIndex || 0;
    this.optional = input.optional || false;
    this.preferenceRank = input.preferenceRank || 0;
    this.quantity = input.quantity || { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID || '';
  }
}
