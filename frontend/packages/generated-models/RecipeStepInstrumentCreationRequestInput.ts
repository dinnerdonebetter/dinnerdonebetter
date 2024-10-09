// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepInstrumentCreationRequestInput {
  productOfRecipeStepIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  name: string;
  preferenceRank: number;
  optionIndex: number;
  optional: boolean;
  productOfRecipeStepProductIndex?: number;
  instrumentID?: string;
  notes: string;
}

export class RecipeStepInstrumentCreationRequestInput implements IRecipeStepInstrumentCreationRequestInput {
  productOfRecipeStepIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  name: string;
  preferenceRank: number;
  optionIndex: number;
  optional: boolean;
  productOfRecipeStepProductIndex?: number;
  instrumentID?: string;
  notes: string;
  constructor(input: Partial<RecipeStepInstrumentCreationRequestInput> = {}) {
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.quantity = input.quantity = { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID;
    this.name = input.name = '';
    this.preferenceRank = input.preferenceRank = 0;
    this.optionIndex = input.optionIndex = 0;
    this.optional = input.optional = false;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.instrumentID = input.instrumentID;
    this.notes = input.notes = '';
  }
}
