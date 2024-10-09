// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepVesselUpdateRequestInput {
  vesselPreposition?: string;
  belongsToRecipeStep?: string;
  name?: string;
  notes?: string;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
  unavailableAfterStep?: boolean;
  vesselID?: string;
}

export class RecipeStepVesselUpdateRequestInput implements IRecipeStepVesselUpdateRequestInput {
  vesselPreposition?: string;
  belongsToRecipeStep?: string;
  name?: string;
  notes?: string;
  quantity: OptionalNumberRange;
  recipeStepProductID?: string;
  unavailableAfterStep?: boolean;
  vesselID?: string;
  constructor(input: Partial<RecipeStepVesselUpdateRequestInput> = {}) {
    this.vesselPreposition = input.vesselPreposition;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.name = input.name;
    this.notes = input.notes;
    this.quantity = input.quantity = {};
    this.recipeStepProductID = input.recipeStepProductID;
    this.unavailableAfterStep = input.unavailableAfterStep;
    this.vesselID = input.vesselID;
  }
}
