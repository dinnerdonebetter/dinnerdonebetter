// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVesselCreationRequestInput {
  name: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  vesselID?: string;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
}

export class RecipeStepVesselCreationRequestInput implements IRecipeStepVesselCreationRequestInput {
  name: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  vesselID?: string;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
  constructor(input: Partial<RecipeStepVesselCreationRequestInput> = {}) {
    this.name = input.name = '';
    this.notes = input.notes = '';
    this.quantity = input.quantity = { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID;
    this.vesselID = input.vesselID;
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.unavailableAfterStep = input.unavailableAfterStep = false;
    this.vesselPreposition = input.vesselPreposition = '';
  }
}
