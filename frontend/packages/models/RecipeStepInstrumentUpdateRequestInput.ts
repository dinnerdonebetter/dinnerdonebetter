// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepInstrumentUpdateRequestInput {
  name: string;
  optionIndex: number;
  optional: boolean;
  quantity: OptionalNumberRange;
  recipeStepProductID: string;
  belongsToRecipeStep: string;
  notes: string;
  preferenceRank: number;
  instrumentID: string;
}

export class RecipeStepInstrumentUpdateRequestInput implements IRecipeStepInstrumentUpdateRequestInput {
  name: string;
  optionIndex: number;
  optional: boolean;
  quantity: OptionalNumberRange;
  recipeStepProductID: string;
  belongsToRecipeStep: string;
  notes: string;
  preferenceRank: number;
  instrumentID: string;
  constructor(input: Partial<RecipeStepInstrumentUpdateRequestInput> = {}) {
    this.name = input.name || '';
    this.optionIndex = input.optionIndex || 0;
    this.optional = input.optional || false;
    this.quantity = input.quantity || {};
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.notes = input.notes || '';
    this.preferenceRank = input.preferenceRank || 0;
    this.instrumentID = input.instrumentID || '';
  }
}
