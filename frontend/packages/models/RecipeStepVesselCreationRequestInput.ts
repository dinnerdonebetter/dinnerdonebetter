// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVesselCreationRequestInput {
  name: string;
  productOfRecipeStepIndex: number;
  recipeStepProductID: string;
  vesselID: string;
  notes: string;
  productOfRecipeStepProductIndex: number;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
}

export class RecipeStepVesselCreationRequestInput implements IRecipeStepVesselCreationRequestInput {
  name: string;
  productOfRecipeStepIndex: number;
  recipeStepProductID: string;
  vesselID: string;
  notes: string;
  productOfRecipeStepProductIndex: number;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
  constructor(input: Partial<RecipeStepVesselCreationRequestInput> = {}) {
    this.name = input.name || '';
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex || 0;
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.vesselID = input.vesselID || '';
    this.notes = input.notes || '';
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex || 0;
    this.quantity = input.quantity || { min: 0 };
    this.unavailableAfterStep = input.unavailableAfterStep || false;
    this.vesselPreposition = input.vesselPreposition || '';
  }
}
