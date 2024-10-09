// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepInstrument {
  preferenceRank: number;
  recipeStepProductID?: string;
  archivedAt?: string;
  belongsToRecipeStep: string;
  createdAt: string;
  name: string;
  notes: string;
  optionIndex: number;
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  id: string;
  instrument?: ValidInstrument;
  lastUpdatedAt?: string;
}

export class RecipeStepInstrument implements IRecipeStepInstrument {
  preferenceRank: number;
  recipeStepProductID?: string;
  archivedAt?: string;
  belongsToRecipeStep: string;
  createdAt: string;
  name: string;
  notes: string;
  optionIndex: number;
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  id: string;
  instrument?: ValidInstrument;
  lastUpdatedAt?: string;
  constructor(input: Partial<RecipeStepInstrument> = {}) {
    this.preferenceRank = input.preferenceRank = 0;
    this.recipeStepProductID = input.recipeStepProductID;
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.createdAt = input.createdAt = '';
    this.name = input.name = '';
    this.notes = input.notes = '';
    this.optionIndex = input.optionIndex = 0;
    this.optional = input.optional = false;
    this.quantity = input.quantity = { min: 0 };
    this.id = input.id = '';
    this.instrument = input.instrument;
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
