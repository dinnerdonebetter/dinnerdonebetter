// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepVesselUpdateRequestInput {
  vesselID?: string;
  vesselPreposition?: string;
  belongsToRecipeStep?: string;
  name?: string;
  notes?: string;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
  unavailableAfterStep?: boolean;
}

export class RecipeStepVesselUpdateRequestInput implements IRecipeStepVesselUpdateRequestInput {
  vesselID?: string;
  vesselPreposition?: string;
  belongsToRecipeStep?: string;
  name?: string;
  notes?: string;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
  unavailableAfterStep?: boolean;
  constructor(input: Partial<RecipeStepVesselUpdateRequestInput> = {}) {
    this.vesselID = input.vesselID;
    this.vesselPreposition = input.vesselPreposition;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.name = input.name;
    this.notes = input.notes;
    this.quantity = input.quantity = {};
    this.recipeStepProductID = input.recipeStepProductID;
    this.unavailableAfterStep = input.unavailableAfterStep;
  }
}
