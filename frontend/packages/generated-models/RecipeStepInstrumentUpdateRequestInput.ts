// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepInstrumentUpdateRequestInput {
  belongsToRecipeStep?: string;
  notes?: string;
  optional?: boolean;
  preferenceRank?: number;
  instrumentID?: string;
  name?: string;
  optionIndex?: number;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
}

export class RecipeStepInstrumentUpdateRequestInput implements IRecipeStepInstrumentUpdateRequestInput {
  belongsToRecipeStep?: string;
  notes?: string;
  optional?: boolean;
  preferenceRank?: number;
  instrumentID?: string;
  name?: string;
  optionIndex?: number;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
  constructor(input: Partial<RecipeStepInstrumentUpdateRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.notes = input.notes;
    this.optional = input.optional;
    this.preferenceRank = input.preferenceRank;
    this.instrumentID = input.instrumentID;
    this.name = input.name;
    this.optionIndex = input.optionIndex;
    this.quantity = input.quantity = {};
    this.recipeStepProductID = input.recipeStepProductID;
  }
}
