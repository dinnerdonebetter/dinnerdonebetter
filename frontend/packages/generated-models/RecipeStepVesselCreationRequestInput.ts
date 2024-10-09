// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVesselCreationRequestInput {
  productOfRecipeStepProductIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vesselID?: string;
  vesselPreposition: string;
  name: string;
  notes: string;
  productOfRecipeStepIndex?: number;
  recipeStepProductID?: string;
}

export class RecipeStepVesselCreationRequestInput implements IRecipeStepVesselCreationRequestInput {
  productOfRecipeStepProductIndex?: number;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vesselID?: string;
  vesselPreposition: string;
  name: string;
  notes: string;
  productOfRecipeStepIndex?: number;
  recipeStepProductID?: string;
  constructor(input: Partial<RecipeStepVesselCreationRequestInput> = {}) {
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.quantity = input.quantity = { min: 0 };
    this.unavailableAfterStep = input.unavailableAfterStep = false;
    this.vesselID = input.vesselID;
    this.vesselPreposition = input.vesselPreposition = '';
    this.name = input.name = '';
    this.notes = input.notes = '';
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.recipeStepProductID = input.recipeStepProductID;
  }
}
