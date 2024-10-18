// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepInstrumentCreationRequestInput {
  instrumentID: string;
  name: string;
  notes: string;
  optionIndex: number;
  optional: boolean;
  preferenceRank: number;
  productOfRecipeStepIndex: number;
  productOfRecipeStepProductIndex: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
}

export class RecipeStepInstrumentCreationRequestInput implements IRecipeStepInstrumentCreationRequestInput {
  instrumentID: string;
  name: string;
  notes: string;
  optionIndex: number;
  optional: boolean;
  preferenceRank: number;
  productOfRecipeStepIndex: number;
  productOfRecipeStepProductIndex: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
  constructor(input: Partial<RecipeStepInstrumentCreationRequestInput> = {}) {
    this.instrumentID = input.instrumentID || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.optionIndex = input.optionIndex || 0;
    this.optional = input.optional || false;
    this.preferenceRank = input.preferenceRank || 0;
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex || 0;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex || 0;
    this.quantity = input.quantity || { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID || '';
  }
}
