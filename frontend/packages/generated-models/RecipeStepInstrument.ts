// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepInstrument {
  archivedAt?: string;
  belongsToRecipeStep: string;
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  optionIndex: number;
  preferenceRank: number;
  createdAt: string;
  id: string;
  instrument?: ValidInstrument;
  lastUpdatedAt?: string;
  name: string;
  notes: string;
  recipeStepProductID?: string;
}

export class RecipeStepInstrument implements IRecipeStepInstrument {
  archivedAt?: string;
  belongsToRecipeStep: string;
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  optionIndex: number;
  preferenceRank: number;
  createdAt: string;
  id: string;
  instrument?: ValidInstrument;
  lastUpdatedAt?: string;
  name: string;
  notes: string;
  recipeStepProductID?: string;
  constructor(input: Partial<RecipeStepInstrument> = {}) {
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.optional = input.optional = false;
    this.quantity = input.quantity = { min: 0 };
    this.optionIndex = input.optionIndex = 0;
    this.preferenceRank = input.preferenceRank = 0;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.instrument = input.instrument;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.notes = input.notes = '';
    this.recipeStepProductID = input.recipeStepProductID;
  }
}
