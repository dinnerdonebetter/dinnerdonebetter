// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepInstrumentUpdateRequestInput {
  belongsToRecipeStep?: string;
  instrumentID?: string;
  optional?: boolean;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
  name?: string;
  notes?: string;
  optionIndex?: number;
  preferenceRank?: number;
}

export class RecipeStepInstrumentUpdateRequestInput implements IRecipeStepInstrumentUpdateRequestInput {
  belongsToRecipeStep?: string;
  instrumentID?: string;
  optional?: boolean;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
  name?: string;
  notes?: string;
  optionIndex?: number;
  preferenceRank?: number;
  constructor(input: Partial<RecipeStepInstrumentUpdateRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.instrumentID = input.instrumentID;
    this.optional = input.optional;
    this.quantity = input.quantity = {};
    this.recipeStepProductID = input.recipeStepProductID;
    this.name = input.name;
    this.notes = input.notes;
    this.optionIndex = input.optionIndex;
    this.preferenceRank = input.preferenceRank;
  }
}
