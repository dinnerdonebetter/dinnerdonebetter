// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepInstrumentCreationRequestInput {
  notes: string;
  productOfRecipeStepProductIndex?: number;
  productOfRecipeStepIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  instrumentID?: string;
  name: string;
  optionIndex: number;
  optional: boolean;
  preferenceRank: number;
}

export class RecipeStepInstrumentCreationRequestInput implements IRecipeStepInstrumentCreationRequestInput {
  notes: string;
  productOfRecipeStepProductIndex?: number;
  productOfRecipeStepIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  instrumentID?: string;
  name: string;
  optionIndex: number;
  optional: boolean;
  preferenceRank: number;
  constructor(input: Partial<RecipeStepInstrumentCreationRequestInput> = {}) {
    this.notes = input.notes = '';
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.quantity = input.quantity = { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID;
    this.instrumentID = input.instrumentID;
    this.name = input.name = '';
    this.optionIndex = input.optionIndex = 0;
    this.optional = input.optional = false;
    this.preferenceRank = input.preferenceRank = 0;
  }
}
