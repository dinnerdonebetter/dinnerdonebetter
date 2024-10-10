// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVesselCreationRequestInput {
  name: string;
  notes: string;
  productOfRecipeStepIndex: number;
  productOfRecipeStepProductIndex: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vesselID: string;
  vesselPreposition: string;
}

export class RecipeStepVesselCreationRequestInput implements IRecipeStepVesselCreationRequestInput {
  name: string;
  notes: string;
  productOfRecipeStepIndex: number;
  productOfRecipeStepProductIndex: number;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vesselID: string;
  vesselPreposition: string;
  constructor(input: Partial<RecipeStepVesselCreationRequestInput> = {}) {
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex || 0;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex || 0;
    this.quantity = input.quantity || { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.unavailableAfterStep = input.unavailableAfterStep || false;
    this.vesselID = input.vesselID || '';
    this.vesselPreposition = input.vesselPreposition || '';
  }
}
