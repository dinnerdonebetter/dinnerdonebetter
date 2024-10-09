// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVesselCreationRequestInput {
  name: string;
  productOfRecipeStepProductIndex?: number;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
  notes: string;
  productOfRecipeStepIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  vesselID?: string;
}

export class RecipeStepVesselCreationRequestInput implements IRecipeStepVesselCreationRequestInput {
  name: string;
  productOfRecipeStepProductIndex?: number;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
  notes: string;
  productOfRecipeStepIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  vesselID?: string;
  constructor(input: Partial<RecipeStepVesselCreationRequestInput> = {}) {
    this.name = input.name = '';
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.unavailableAfterStep = input.unavailableAfterStep = false;
    this.vesselPreposition = input.vesselPreposition = '';
    this.notes = input.notes = '';
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.quantity = input.quantity = { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID;
    this.vesselID = input.vesselID;
  }
}
