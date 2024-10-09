// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepInstrument {
  optionIndex: number;
  recipeStepProductID?: string;
  belongsToRecipeStep: string;
  lastUpdatedAt?: string;
  name: string;
  instrument?: ValidInstrument;
  notes: string;
  optional: boolean;
  preferenceRank: number;
  quantity: NumberRangeWithOptionalMax;
  archivedAt?: string;
  createdAt: string;
  id: string;
}

export class RecipeStepInstrument implements IRecipeStepInstrument {
  optionIndex: number;
  recipeStepProductID?: string;
  belongsToRecipeStep: string;
  lastUpdatedAt?: string;
  name: string;
  instrument?: ValidInstrument;
  notes: string;
  optional: boolean;
  preferenceRank: number;
  quantity: NumberRangeWithOptionalMax;
  archivedAt?: string;
  createdAt: string;
  id: string;
  constructor(input: Partial<RecipeStepInstrument> = {}) {
    this.optionIndex = input.optionIndex = 0;
    this.recipeStepProductID = input.recipeStepProductID;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.instrument = input.instrument;
    this.notes = input.notes = '';
    this.optional = input.optional = false;
    this.preferenceRank = input.preferenceRank = 0;
    this.quantity = input.quantity = { min: 0 };
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
  }
}
