// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepInstrumentUpdateRequestInput {
  quantity: OptionalNumberRange;
  instrumentID?: string;
  name?: string;
  notes?: string;
  preferenceRank?: number;
  recipeStepProductID?: string;
  belongsToRecipeStep?: string;
  optionIndex?: number;
  optional?: boolean;
}

export class RecipeStepInstrumentUpdateRequestInput implements IRecipeStepInstrumentUpdateRequestInput {
  quantity: OptionalNumberRange;
  instrumentID?: string;
  name?: string;
  notes?: string;
  preferenceRank?: number;
  recipeStepProductID?: string;
  belongsToRecipeStep?: string;
  optionIndex?: number;
  optional?: boolean;
  constructor(input: Partial<RecipeStepInstrumentUpdateRequestInput> = {}) {
    this.quantity = input.quantity = {};
    this.instrumentID = input.instrumentID;
    this.name = input.name;
    this.notes = input.notes;
    this.preferenceRank = input.preferenceRank;
    this.recipeStepProductID = input.recipeStepProductID;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.optionIndex = input.optionIndex;
    this.optional = input.optional;
  }
}
