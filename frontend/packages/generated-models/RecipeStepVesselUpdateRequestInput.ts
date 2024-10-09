// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepVesselUpdateRequestInput {
  recipeStepProductID?: string;
  unavailableAfterStep?: boolean;
  vesselID?: string;
  vesselPreposition?: string;
  belongsToRecipeStep?: string;
  name?: string;
  notes?: string;
  quantity: OptionalNumberRange;
}

export class RecipeStepVesselUpdateRequestInput implements IRecipeStepVesselUpdateRequestInput {
  recipeStepProductID?: string;
  unavailableAfterStep?: boolean;
  vesselID?: string;
  vesselPreposition?: string;
  belongsToRecipeStep?: string;
  name?: string;
  notes?: string;
  quantity: OptionalNumberRange;
  constructor(input: Partial<RecipeStepVesselUpdateRequestInput> = {}) {
    this.recipeStepProductID = input.recipeStepProductID;
    this.unavailableAfterStep = input.unavailableAfterStep;
    this.vesselID = input.vesselID;
    this.vesselPreposition = input.vesselPreposition;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.name = input.name;
    this.notes = input.notes;
    this.quantity = input.quantity = {};
  }
}
