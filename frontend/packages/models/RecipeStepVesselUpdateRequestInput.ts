// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepVesselUpdateRequestInput {
  name: string;
  notes: string;
  quantity: OptionalNumberRange;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vesselID: string;
  vesselPreposition: string;
  belongsToRecipeStep: string;
}

export class RecipeStepVesselUpdateRequestInput implements IRecipeStepVesselUpdateRequestInput {
  name: string;
  notes: string;
  quantity: OptionalNumberRange;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vesselID: string;
  vesselPreposition: string;
  belongsToRecipeStep: string;
  constructor(input: Partial<RecipeStepVesselUpdateRequestInput> = {}) {
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.quantity = input.quantity || {};
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.unavailableAfterStep = input.unavailableAfterStep || false;
    this.vesselID = input.vesselID || '';
    this.vesselPreposition = input.vesselPreposition || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
  }
}
