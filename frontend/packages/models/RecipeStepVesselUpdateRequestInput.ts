// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepVesselUpdateRequestInput {
  belongsToRecipeStep: string;
  name: string;
  notes: string;
  quantity: OptionalNumberRange;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vesselID: string;
  vesselPreposition: string;
}

export class RecipeStepVesselUpdateRequestInput implements IRecipeStepVesselUpdateRequestInput {
  belongsToRecipeStep: string;
  name: string;
  notes: string;
  quantity: OptionalNumberRange;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vesselID: string;
  vesselPreposition: string;
  constructor(input: Partial<RecipeStepVesselUpdateRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.quantity = input.quantity || {};
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.unavailableAfterStep = input.unavailableAfterStep || false;
    this.vesselID = input.vesselID || '';
    this.vesselPreposition = input.vesselPreposition || '';
  }
}
